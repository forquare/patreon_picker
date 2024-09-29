package main

import (
	"github.com/forquare/patreon_picker/config"
	"github.com/forquare/patreon_picker/handlers"
	"github.com/forquare/patreon_picker/utils"
	"html/template"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

var Version = " DeVeLoPmEnT"

func getVersion(c *gin.Context) {
	c.Set("version", Version)
	c.Next()
}

func main() {
	logger.SetLevel(logger.DebugLevel)
	l, _ := logger.ParseLevel(config.GetConfig().LogLevel)
	logger.SetLevel(l)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//err := router.SetTrustedProxies([]string{""})
	//if err != nil {
	//	logger.Panic(err)
	//	return
	//}

	store := cookie.NewStore([]byte(config.GetConfig().Session.CookieKey))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 2, // Cookie valid for two days
	})

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions(config.GetConfig().Session.Name, store))
	router.Use(getVersion)

	router.Static("/css", "./static/css")
	router.Static("/img", "./static/img")
	router.SetFuncMap(template.FuncMap{
		"add": utils.Add,
	})
	router.LoadHTMLFiles("templates/index.tmpl", "templates/error.tmpl", "templates/auth.tmpl")

	router.GET("/login", handlers.LoginHandler)
	router.GET("/auth", handlers.AuthHandler)

	authorized := router.Group("/")
	authorized.Use(handlers.AuthorizeRequest())
	{
		authorized.GET("/", handlers.IndexHandler)
	}

	err := router.Run(config.GetConfig().Connection.Address + ":" + config.GetConfig().Connection.Port)
	if err != nil {
		logger.Panic(err)
		return
	}
	logger.Info("Patreon Picker" + Version)
}
