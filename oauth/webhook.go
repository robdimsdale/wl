package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wl"
)

// Webhooks gets all webhooks for all lists.
func (c oauthClient) Webhooks() ([]wl.Webhook, error) {
	lists, err := c.Lists()
	if err != nil {
		return nil, err
	}

	listCount := len(lists)
	c.logger.Debug(
		"webhooks",
		map[string]interface{}{"listCount": listCount},
	)

	webhooksChan := make(chan []wl.Webhook, listCount)
	idErrChan := make(chan idErr, listCount)
	for _, l := range lists {
		go func(list wl.List) {
			c.logger.Debug(
				"webhooks - getting webhooks for list",
				map[string]interface{}{"listID": list.ID},
			)
			webhooks, err := c.WebhooksForListID(list.ID)
			idErrChan <- idErr{idType: "list", id: list.ID, err: err}
			webhooksChan <- webhooks
		}(l)
	}

	e := multiIDErr{}
	for i := 0; i < listCount; i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug(
				"webhooks - error received getting webhooks for list",
				map[string]interface{}{"listID": idErr.id, "err": err},
			)
			e.addError(idErr)
		}
	}

	totalWebhooks := []wl.Webhook{}
	for i := 0; i < listCount; i++ {
		webhooks := <-webhooksChan
		totalWebhooks = append(totalWebhooks, webhooks...)
	}

	if len(e.errors()) > 0 {
		return totalWebhooks, e
	}

	return totalWebhooks, nil
}

// WebhooksForListID returns Webhooks for the provided listID.
func (c oauthClient) WebhooksForListID(listID uint) ([]wl.Webhook, error) {
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

	webhooks := []wl.Webhook{}
	err = json.NewDecoder(resp.Body).Decode(&webhooks)
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

// Webhook returns the Webhook for the corresponding webhookID.
func (c oauthClient) Webhook(webhookID uint) (wl.Webhook, error) {
	allWebhooks, err := c.Webhooks()
	for _, w := range allWebhooks {
		if w.ID == webhookID {
			return w, err
		}
	}

	return wl.Webhook{}, fmt.Errorf("webhook not found")
}

// CreateWebhook creates a new webhook with the provided parameters.
// listID must be non-zero; the remaining parameters are not validated.
func (c oauthClient) CreateWebhook(
	listID uint,
	url string,
	processorType string,
	configuration string,
) (wl.Webhook, error) {
	if listID == 0 {
		return wl.Webhook{}, errors.New("listID must be > 0")
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
		return wl.Webhook{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Webhook{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wl.Webhook{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	webhook := wl.Webhook{}
	err = json.NewDecoder(resp.Body).Decode(&webhook)
	if err != nil {
		return wl.Webhook{}, err
	}
	return webhook, nil
}

// DeleteNote deletes the provided webhook.
func (c oauthClient) DeleteWebhook(webhook wl.Webhook) error {
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
