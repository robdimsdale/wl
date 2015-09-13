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

// Folders returns Folders created by the current user.
func (c oauthClient) Folders() ([]wundergo.Folder, error) {
	url := fmt.Sprintf(
		"%s/folders",
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

	folders := []wundergo.Folder{}
	err = json.NewDecoder(resp.Body).Decode(&folders)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return nil, err
	}
	return folders, nil
}

type folderCreateConfig struct {
	Title   string `json:"title"`
	ListIDs []uint `json:"list_ids"`
}

// CreateFolder creates a new folder with the provided parameters.
func (c oauthClient) CreateFolder(
	title string,
	listIDs []uint,
) (wundergo.Folder, error) {
	if title == "" {
		return wundergo.Folder{}, errors.New("title must be non-empty")
	}

	if listIDs == nil {
		return wundergo.Folder{}, errors.New("listIDs must be non-nil")
	}

	fcc := folderCreateConfig{
		Title:   title,
		ListIDs: listIDs,
	}

	body, err := json.Marshal(fcc)
	if err != nil {
		return wundergo.Folder{}, err
	}

	reqURL := fmt.Sprintf("%s/folders", c.apiURL)

	req, err := c.newPostRequest(reqURL, body)
	if err != nil {
		return wundergo.Folder{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.Folder{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		if resp.Body != nil {
			b, _ := ioutil.ReadAll(resp.Body)
			c.logger.Debug("", lager.Data{"response.Body": string(b)})
		}
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.Folder{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	folder := wundergo.Folder{}
	err = json.NewDecoder(resp.Body).Decode(&folder)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.Folder{}, err
	}
	return folder, nil
}

// Folder returns the Folder for the corresponding folderID.
func (c oauthClient) Folder(folderID uint) (wundergo.Folder, error) {
	if folderID == 0 {
		return wundergo.Folder{}, errors.New("folderID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/folders/%d",
		c.apiURL,
		folderID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.Folder{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.Folder{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Folder{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	folder := wundergo.Folder{}
	err = json.NewDecoder(resp.Body).Decode(&folder)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.Folder{}, err
	}
	return folder, nil
}
