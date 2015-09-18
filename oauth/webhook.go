package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wundergo"
)

// WebhooksForListID returns Webhooks for the provided listID.
func (c oauthClient) WebhooksForListID(listID uint) ([]wundergo.Webhook, error) {
	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/webhooks?list_id=%d",
		c.apiURL,
		listID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	webhooks := []wundergo.Webhook{}
	err = json.NewDecoder(resp.Body).Decode(&webhooks)
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

// CreateWebhook creates a new webhook with the provided parameters.
// listID must be non-zero; the remaining parameters are not validated.
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

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Webhook{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wundergo.Webhook{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	webhook := wundergo.Webhook{}
	err = json.NewDecoder(resp.Body).Decode(&webhook)
	if err != nil {
		return wundergo.Webhook{}, err
	}
	return webhook, nil
}

// DeleteNote deletes the provided webhook.
func (c oauthClient) DeleteWebhook(webhook wundergo.Webhook) error {
	url := fmt.Sprintf(
		"%s/webhooks/%d",
		c.apiURL,
		webhook.ID,
	)

	req, err := c.newDeleteRequest(url)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
	}

	return nil
}
