package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/robdimsdale/wundergo"
)

var (
	accessToken = flag.String("accessToken", "", "Access Token")
	clientID    = flag.String("clientID", "", "Client ID")
)

func main() {
	flag.Parse()

	if *accessToken == "" {
		log.Fatal("Access Token must be provided")
	}

	if *clientID == "" {
		log.Fatal("Client ID must be provided")
	}

	client := wundergo.NewClient(*accessToken, *clientID)
	user, err := client.User()
	if err != nil {
		log.Printf("Error getting user: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", user)
}
