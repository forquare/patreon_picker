package utils

import (
	"gopkg.in/mxpv/patreon-go.v1"
	"sync"
)

var (
	once sync.Once
	user *patreon.User
)

func GetAuthenticatedUser(client *patreon.Client) *patreon.User {
	once.Do(func() {
		u, err := client.FetchUser()
		if err != nil {
			panic(err)
		}
		user = &u.Data
	})
	return user
}

func GetAuthenticatedUserName(client *patreon.Client) string {
	return GetAuthenticatedUser(client).Attributes.FullName
}
