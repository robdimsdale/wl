package wundergo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

func (c OauthClient) performHTTPAction(
	url string,
	action string,
	body string) ([]byte, error) {

	req, err := http.NewRequest(action, url, nil)
	if err != nil {
		log.Printf("Error constructing http request: %s\n", err.Error())
	}
	client := &http.Client{}

	req.Header.Add("X-Access-Token", c.accessToken)
	req.Header.Add("X-Client-ID", c.clientID)

	if body != "" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Body = ioutil.NopCloser(strings.NewReader(body))
	}

	resp, err := client.Do(req)
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
	b, err := c.performHTTPAction(fmt.Sprintf("%s/user", apiUrl), "GET", "")
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

func (c OauthClient) UpdateUser(user User) (User, error) {
	body := fmt.Sprintf("revision=%d&name=%s", user.Revision, user.Name)
	b, err := c.performHTTPAction(fmt.Sprintf("%s/user", apiUrl), "PUT", body)
	log.Printf("Body received: %s\n", string(b))
	if err != nil {
		return User{}, err
	}

	var u User
	err = json.Unmarshal(b, &u)
	if err != nil {
		log.Printf("Body received: %s\n", string(b))
		return User{}, err
	}
	return u, nil
}
