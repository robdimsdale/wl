package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wl"
)

// ListPositions returns the positions of all Lists the client can access.
// The returned ListPosition.Values might be empty if the Lists have never been reordered.
func (c oauthClient) ListPositions() ([]wl.Position, error) {
	url := fmt.Sprintf(
		"%s/list_positions",
		c.apiURL,
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

	listPositions := []wl.Position{}
	err = json.NewDecoder(resp.Body).Decode(&listPositions)
	if err != nil {
		return nil, err
	}
	return listPositions, nil
}

// ListPosition returns the ListPosition associated with the provided listPositionID.
func (c oauthClient) ListPosition(listPositionID uint) (wl.Position, error) {
	url := fmt.Sprintf(
		"%s/list_positions/%d",
		c.apiURL,
		listPositionID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wl.Position{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Position{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.Position{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	listPosition := wl.Position{}
	err = json.NewDecoder(resp.Body).Decode(&listPosition)
	if err != nil {
		return wl.Position{}, err
	}
	return listPosition, nil
}

// UpdateListPosition updates the provided ListPosition.
// This will reorder the Lists.
func (c oauthClient) UpdateListPosition(listPosition wl.Position) (wl.Position, error) {
	body, err := json.Marshal(listPosition)
	if err != nil {
		return wl.Position{}, err
	}

	url := fmt.Sprintf(
		"%s/list_positions/%d",
		c.apiURL,
		listPosition.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wl.Position{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Position{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.Position{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	returnedListPosition := wl.Position{}
	err = json.NewDecoder(resp.Body).Decode(&returnedListPosition)
	if err != nil {
		return wl.Position{}, err
	}
	return returnedListPosition, nil
}
