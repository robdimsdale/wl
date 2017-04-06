package oauth

import (
	"fmt"
	"net/http"
)

func (c oauthClient) Curl(method string, url string, body []byte) (*http.Response, error) {
	reqURL := fmt.Sprintf(
		"%s/%s",
		c.apiURL,
		url,
	)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	c.addBody(req, body)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
