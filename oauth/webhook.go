package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pivotal-golang/lager"
	"github.com/robdimsdale/wundergo"
)

func (c oauthClient) CreateWebhook(
	listID uint,
	url string,
	processorType string,
	configuration string,
) (wundergo.Webhook, error) {
	if listID == 0 {
		return wundergo.Webhook{}, errors.New("listID must be > 0")
	}

	body := []byte(fmt.Sprintf(`{
		"list_id":%d,
		"url":"%s",
		"processor_type":"%s",
		"configuration":"%s"
	}`,
		listID,
		url,
		processorType,
		configuration,
	))

	reqURL := fmt.Sprintf("%s/webhooks", c.apiURL)

	req, err := c.newPostRequest(reqURL, body)
	if err != nil {
		return wundergo.Webhook{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.Webhook{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		if resp.Body != nil {
			b, _ := ioutil.ReadAll(resp.Body)
			c.logger.Debug("", lager.Data{"response.Body": string(b)})
		}
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.Webhook{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	note := wundergo.Webhook{}
	err = json.NewDecoder(resp.Body).Decode(&note)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.Webhook{}, err
	}
	return note, nil

	return wundergo.Webhook{}, nil
}
