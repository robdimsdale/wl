package oauth

import (
	"fmt"
	"net/http"
)

// AvatarURL returns the URL of the user associated with userID
// sized is checked to ensure it is positive, but is not validate otherwise.
func (c oauthClient) AvatarURL(userID uint, size int, fallback bool) (string, error) {
	url := fmt.Sprintf(
		"%s/avatar?user_id=%d",
		c.apiURL,
		userID,
	)

	if size > 0 {
		url = fmt.Sprintf(
			"%s&size=%d",
			url,
			size,
		)
	}

	if !fallback {
		url = fmt.Sprintf(
			"%s&fallback=%t",
			url,
			fallback,
		)

	}

	req, err := c.newGetRequest(url)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return "", err
	}

	if !fallback && resp.StatusCode == http.StatusNoContent {
		return "", nil
	}

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusFound)
	}

	location, err := resp.Location()
	if err != nil {
		return "", err
	}
	return location.String(), nil
}
