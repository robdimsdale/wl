package oauth

import (
	"fmt"
	"net/http"
	"strings"
)

func (c oauthClient) Curl(method string, url string, body []byte) (*http.Response, error) {
	url = strings.TrimPrefix(url, "/")

	reqURL := fmt.Sprintf(
		"%s/%s",
		c.apiURL,
		url,
	)

	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, err
	}

	c.addBody(req, body)

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
