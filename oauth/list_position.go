package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wundergo"
)

// ListPositions returns the positions of all Lists the client can access.
// The returned ListPosition.Values might be empty if the Lists have never been reordered.
func (c oauthClient) ListPositions() (*[]wundergo.Position, error) {
	url := fmt.Sprintf(
		"%s/list_positions",
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

	listPositions := &[]wundergo.Position{}
	err = json.NewDecoder(resp.Body).Decode(listPositions)
	if err != nil {
		c.logger.Println(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return listPositions, nil
}

// ListPosition returns the ListPosition associated with the provided listPositionID.
func (c oauthClient) ListPosition(listPositionID uint) (*wundergo.Position, error) {
	url := fmt.Sprintf(
		"%s/list_positions/%d",
		c.apiURL,
		listPositionID,
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

	listPosition := &wundergo.Position{}
	err = json.NewDecoder(resp.Body).Decode(listPosition)
	if err != nil {
		c.logger.Println(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return listPosition, nil
}

// UpdateListPosition updates the provided ListPosition.
// This will reorder the Lists.
func (c oauthClient) UpdateListPosition(listPosition wundergo.Position) (*wundergo.Position, error) {
	body, err := json.Marshal(listPosition)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"%s/list_positions/%d",
		c.apiURL,
		listPosition.ID,
	)

	req, err := c.newPatchRequest(url, body)
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

	returnedListPosition := &wundergo.Position{}
	err = json.NewDecoder(resp.Body).Decode(returnedListPosition)
	if err != nil {
		c.logger.Println(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return returnedListPosition, nil
}
