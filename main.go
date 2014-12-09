package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	client := &http.Client{}
	userRequest, err := http.NewRequest("GET", "https://a.wunderlist.com/api/v1/user", nil)
	userRequest.Header.Add("X-Access-Token", *accessToken)
	userRequest.Header.Add("X-Client-ID", *clientID)

	resp, err := client.Do(userRequest)
	if err != nil {
		log.Printf("Error making request: %s\n", err.Error())
	}
	if resp != nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %s\n", resp.Body)
		}
		fmt.Printf("%s\n", string(body))
	}
}
