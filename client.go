package wundergo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	apiUrl = "https://a.wunderlist.com/api/v1"
)

type Client interface {
	User() (User, error)
}

type OauthClient struct {
	accessToken string
	clientID    string
}

func NewClient(accessToken string, clientID string) *OauthClient {
	return &OauthClient{
		accessToken: accessToken,
		clientID:    clientID,
	}
}

func (c OauthClient) getRequest(url string) ([]byte, error) {
	client := &http.Client{}
	userRequest, err := http.NewRequest("GET", url, nil)
	userRequest.Header.Add("X-Access-Token", c.accessToken)
	userRequest.Header.Add("X-Client-ID", c.clientID)

	resp, err := client.Do(userRequest)
	if err != nil {
		log.Printf("Error making request: %s\n", err.Error())
	}

	if resp == nil {
		return nil, errors.New("Nil body returned")
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c OauthClient) User() (User, error) {
	b, err := c.getRequest(fmt.Sprintf("%s/user", apiUrl))
	if err != nil {
		return User{}, err
	}

	var u User
	err = json.Unmarshal(b, &u)
	if err != nil {
		return User{}, err
	}
	return u, nil
}
