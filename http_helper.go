package wundergo

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type HTTPHelper interface {
	Get(url string) ([]byte, error)
	Put(url string, body string) ([]byte, error)
}

type oauthClientHTTPHelper struct {
	accessToken string
	clientID    string
}

func newOauthClientHTTPHelper(accessToken string, clientID string) *oauthClientHTTPHelper {
	return &oauthClientHTTPHelper{
		accessToken: accessToken,
		clientID:    clientID,
	}
}

func (h oauthClientHTTPHelper) Get(url string) ([]byte, error) {
	return h.performHTTPAction(url, "GET", "")
}

func (h oauthClientHTTPHelper) Put(url string, body string) ([]byte, error) {
	return h.performHTTPAction(url, "PUT", body)
}

func (h oauthClientHTTPHelper) performHTTPAction(
	url string,
	action string,
	body string) ([]byte, error) {

	req, err := http.NewRequest(action, url, nil)
	if err != nil {
		log.Printf("Error constructing http request: %s\n", err.Error())
	}
	client := &http.Client{}

	req.Header.Add("X-Access-Token", h.accessToken)
	req.Header.Add("X-Client-ID", h.clientID)

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
