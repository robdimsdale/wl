package oauth

import (
	"encoding/json"
	"fmt"
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
