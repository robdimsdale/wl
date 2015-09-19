package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wundergo"
)

// Root returns the Root for the current user.
func (c oauthClient) Root() (wundergo.Root, error) {
	url := fmt.Sprintf(
		"%s/root",
		c.apiURL,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.Root{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Root{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Root{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	root := wundergo.Root{}
	err = json.NewDecoder(resp.Body).Decode(&root)
	if err != nil {
		return root, err
	}
	return root, nil
}
