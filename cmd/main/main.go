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
	logger := wundergo.NewPrintlnLogger()

	fmt.Printf("----------------------------------\n")
	logger.LogLine("User info\n")
	user, err := client.User()
	if err != nil {
		logger.LogLine(fmt.Sprintf("Error getting user: %s\n", err.Error()))
	}
	fmt.Printf("%+v\n", user)
	fmt.Printf("----------------------------------\n")

	logger.LogLine("Creating new list\n")
	newList, err := client.CreateList("newListTitle")
	if err != nil {
		logger.LogLine(fmt.Sprintf("Error creating list: %s\n", err.Error()))
	}
	fmt.Printf("%+v\n", newList)
	fmt.Printf("----------------------------------\n")

	logger.LogLine("Getting lists\n")
	lists, err := client.Lists()
	if err != nil {
		logger.LogLine(fmt.Sprintf("Error getting lists: %s\n", err.Error()))
	}
	fmt.Printf("%+v\n", lists)
	fmt.Printf("----------------------------------\n")

	logger.LogLine("Deleting new list\n")
	err = client.DeleteList(newList)
	if err != nil {
		logger.LogLine(fmt.Sprintf("Error deleting list: %s\n", err.Error()))
	}
	fmt.Printf("----------------------------------\n")

	logger.LogLine("Getting lists\n")
	lists, err = client.Lists()
	if err != nil {
		logger.LogLine(fmt.Sprintf("Error getting lists: %s\n", err.Error()))
	}
	fmt.Printf("%+v\n", lists)
	fmt.Printf("----------------------------------\n")
}
