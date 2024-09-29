package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"github.com/forquare/patreon_picker/config"
	"github.com/forquare/patreon_picker/picker"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gopkg.in/mxpv/patreon-go.v1"
)

var conf *oauth2.Config

func init() {
	conf = &oauth2.Config{
		ClientID:     config.GetConfig().Credentials.Id,
		ClientSecret: config.GetConfig().Credentials.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.patreon.com/oauth2/authorize",
			TokenURL: "https://www.patreon.com/api/oauth2/token",
		},
		Scopes:      []string{"users", "pledges-to-me", "my-campaign"},
		RedirectURL: config.GetConfig().Credentials.RedirectURL,
	}
	gob.Register(oauth2.Token{})
	logger.Trace("Handlers initialized")
}

func AuthorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("token")
		if v == nil {
			LoginHandler(c)
		}
		c.Next()
	}
}

func RandToken(l int) string {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func getClient(token *oauth2.Token) (*patreon.Client, error) {
	tc := conf.Client(context.Background(), token)

	client := patreon.NewClient(tc)
	_, err := client.FetchUser()

	return client, err
}

func AuthHandler(c *gin.Context) {
	// Handle the exchange code to initiate transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")
	if retrievedState != queryState {
		log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	code := c.Request.URL.Query().Get("code")

	tok, err := conf.Exchange(context.Background(), code)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, err = getClient(tok)
	if err != nil {
		log.Printf("Error making new client.\n")
		c.Redirect(http.StatusSeeOther, "/login")
	}

	session.Set("token", tok)
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func LoginHandler(c *gin.Context) {
	state := RandToken(32)
	session := sessions.Default(c)
	session.Set("state", state)
	err := session.Save()
	if err != nil {
		return
	}
	link := conf.AuthCodeURL(state)
	c.HTML(http.StatusOK, "auth.tmpl", gin.H{"link": link, "version": c.MustGet("version")})
}

func IndexHandler(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("token")
	if v != nil {
		tok := session.Get("token").(oauth2.Token)
		client, err := getClient(&tok)
		if err != nil {
			log.Printf("Error making new client.\n")
			c.Redirect(http.StatusSeeOther, "/login")
		}
		mentions := picker.GetPatreonMentions(client)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"mentions": mentions, "version": c.MustGet("version")})
	}
}
