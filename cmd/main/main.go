package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robdimsdale/wundergo"
)

const (
	WL_CLIENT_ID    = "WL_CLIENT_ID"
	WL_ACCESS_TOKEN = "WL_ACCESS_TOKEN"
)

func main() {
	accessToken := os.Getenv(WL_ACCESS_TOKEN)
	clientID := os.Getenv(WL_CLIENT_ID)

	if accessToken == "" {
		log.Fatal("Access Token must be provided")
	}

	if clientID == "" {
		log.Fatal("Client ID must be provided")
	}

	client := wundergo.NewOauthClient(accessToken, clientID)

	user, err := client.User()
	if err != nil {
		log.Printf("Error getting user: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", user)
}
