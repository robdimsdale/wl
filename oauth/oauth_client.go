package oauth

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/robdimsdale/wl"
	"github.com/robdimsdale/wl/logger"
)

// oauthClient is an implementation of wl.Client.
type oauthClient struct {
	apiURL      string
	accessToken string
	clientID    string
	logger      logger.Logger
}

// NewClient is a utility method to simplify initialization
// of a new oauthClient.
func NewClient(
	accessToken string,
	clientID string,
	apiURL string,
	logger logger.Logger,
) wl.Client {
	return &oauthClient{
		apiURL:      apiURL,
		accessToken: accessToken,
		clientID:    clientID,
		logger:      logger,
	}
}

func (c oauthClient) validateRecurrence(recurrenceType string, recurrenceCount uint) error {
	if recurrenceType == "" && recurrenceCount > 0 {
		return errors.New("recurrenceCount must be zero if provided recurrenceType is not provided")
	}

	if recurrenceCount == 0 && recurrenceType != "" {
		return errors.New("recurrenceType must be valid if provided recurrenceCount is non-zero")
	}

	return nil
}

func (c oauthClient) addAuthHeaders(req *http.Request) {
	req.Header.Add("X-Access-Token", c.accessToken)
	req.Header.Add("X-Client-ID", c.clientID)
}

func (c oauthClient) newGetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.addAuthHeaders(req)
	return req, nil
}

func (c oauthClient) do(req *http.Request) (*http.Response, error) {
	c.logRequest(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		c.logResponse(resp)
	}
	return resp, err
}

func (c oauthClient) logRequest(req *http.Request) {
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		c.logger.Error("received error while dumping HTTP request", err)
	} else {
		if reqDump != nil {
			c.logger.Debug(
				" - sending request",
				map[string]interface{}{"request": string(reqDump)})
		}
	}
}

func (c oauthClient) logResponse(resp *http.Response) {
	if resp == nil {
		c.logger.Debug(" - nil response received")
		return
	}

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		c.logger.Error("received error while dumping HTTP Response", err)
	} else {
		if respDump != nil {
			c.logger.Debug(" - received response", map[string]interface{}{"response": string(respDump)})
		}
	}
}

func (c oauthClient) newPostRequest(url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	c.addAuthHeaders(req)
	c.addBody(req, body)

	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c oauthClient) addBody(req *http.Request, body []byte) {
	if body != nil && len(body) > 0 {
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
}

func (c oauthClient) newPutRequest(url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}

	c.addAuthHeaders(req)
	c.addBody(req, body)

	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c oauthClient) newPatchRequest(url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return nil, err
	}

	c.addAuthHeaders(req)
	c.addBody(req, body)

	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c oauthClient) newDeleteRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	c.addAuthHeaders(req)
	return req, nil
}
