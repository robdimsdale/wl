package oauth

import (
	"fmt"
	"net/http"
)

// AvatarURL returns the URL of the user associated with userID.
// Non-positive sizes are ignored, positive sizes are validated according to
// the values at https://developer.wunderlist.com/documentation/endpoints/avatar.
func (c oauthClient) AvatarURL(userID uint, size int, fallback bool) (string, error) {
	url := fmt.Sprintf(
		"%s/avatar?user_id=%d",
		c.apiURL,
		userID,
	)

	if size > 0 {
		if !c.validSize(size) {
			return "", fmt.Errorf("Invalid size: %d", size)
		}
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
	c.logRequest(req)

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return "", err
	}

	if resp != nil {
		c.logResponse(resp)
	}

	if fallback {
		if resp.StatusCode != http.StatusFound {
			return "", fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusFound)
		}
	} else {
		if resp.StatusCode == http.StatusNoContent {
			return "", nil
		}

		if resp.StatusCode != http.StatusFound {
			return "", fmt.Errorf("Unexpected response code %d - expected either %d or %d", resp.StatusCode, http.StatusNoContent, http.StatusFound)
		}
	}

	location, err := resp.Location()
	if err != nil {
		return "", err
	}
	return location.String(), nil
}

func (c oauthClient) validSize(size int) bool {
	validSizes := []int{25, 28, 30, 32, 50, 54, 56, 60, 64, 108, 128, 135, 256, 270, 512}

	for _, s := range validSizes {
		if s == size {
			return true
		}
	}

	return false
}
