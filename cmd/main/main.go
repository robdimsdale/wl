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

	log.Printf("User info\n")
	user, err := client.User()
	if err != nil {
		log.Printf("Error getting user: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", user)
	fmt.Printf("----------------------------------\n")

	log.Printf("Creating new list\n")
	newList, err := client.CreateList("newListTitle")
	if err != nil {
		log.Printf("Error creating list: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", newList)
	fmt.Printf("----------------------------------\n")

	log.Printf("Getting lists\n")
	lists, err := client.Lists()
	if err != nil {
		log.Printf("Error getting lists: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", lists)
	fmt.Printf("----------------------------------\n")

	log.Printf("Deleting new list\n")
	err = client.DeleteList(newList)
	if err != nil {
		log.Printf("Error deleting list: %s\n", err.Error())
	}
	fmt.Printf("----------------------------------\n")

	log.Printf("Getting lists\n")
	lists, err = client.Lists()
	if err != nil {
		log.Printf("Error getting lists: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", lists)
	fmt.Printf("----------------------------------\n")
}
