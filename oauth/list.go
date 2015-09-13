package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pivotal-golang/lager"
	"github.com/robdimsdale/wundergo"
)

// Lists returns all lists the client has permission to access.
func (c oauthClient) Lists() ([]wundergo.List, error) {
	url := fmt.Sprintf(
		"%s/lists",
		c.apiURL,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	lists := []wundergo.List{}
	err = json.NewDecoder(resp.Body).Decode(&lists)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return lists, err
	}
	return lists, nil
}

// List returns the list for the corresponding listID.
func (c oauthClient) List(listID uint) (wundergo.List, error) {
	url := fmt.Sprintf(
		"%s/lists/%d",
		c.apiURL,
		listID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.List{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.List{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.List{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	list := wundergo.List{}
	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.List{}, err
	}
	return list, nil
}

// CreateList creates a list with the provided title.
func (c oauthClient) CreateList(title string) (wundergo.List, error) {
	if title == "" {
		return wundergo.List{}, fmt.Errorf("title must be non-empty")
	}

	url := fmt.Sprintf("%s/lists", c.apiURL)
	body := []byte(fmt.Sprintf(`{"title":"%s"}`, title))

	req, err := c.newPostRequest(url, body)
	if err != nil {
		return wundergo.List{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.List{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wundergo.List{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	list := wundergo.List{}
	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.List{}, err
	}
	return list, nil
}

// UpdateList updates the provided List.
func (c oauthClient) UpdateList(list wundergo.List) (wundergo.List, error) {
	body, err := json.Marshal(list)
	if err != nil {
		return wundergo.List{}, err
	}

	url := fmt.Sprintf(
		"%s/lists/%d",
		c.apiURL,
		list.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wundergo.List{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.List{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.List{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	returnedList := wundergo.List{}
	err = json.NewDecoder(resp.Body).Decode(&returnedList)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.List{}, err
	}
	return returnedList, nil
}

// DeleteList deletes the provided list.
func (c oauthClient) DeleteList(list wundergo.List) error {
	url := fmt.Sprintf(
		"%s/lists/%d?revision=%d",
		c.apiURL,
		list.ID,
		list.Revision,
	)

	req, err := c.newDeleteRequest(url)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
	}

	return nil
}
