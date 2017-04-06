package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wl"
)

// Root returns the Root for the current user.
func (c oauthClient) Root() (wl.Root, error) {
	url := fmt.Sprintf(
		"%s/root",
		c.apiURL,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return wl.Root{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return wl.Root{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.Root{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	root := wl.Root{}
	err = json.NewDecoder(resp.Body).Decode(&root)
	if err != nil {
		return root, err
	}
	return root, nil
}
