package wundergo

import (
	"bytes"
	"errors"
	"io/ioutil"
)

var NewHTTPTransportHelper = func() HTTPTransportHelper {
	return NewDefaultHTTPTransportHelper()
}

type HTTPHelper interface {
	Get(url string) ([]byte, error)
	Post(url string, body []byte) ([]byte, error)
	Put(url string, body []byte) ([]byte, error)
	Patch(url string, body []byte) ([]byte, error)
	Delete(url string) error
}

type OauthClientHTTPHelper struct {
	accessToken   string
	clientID      string
	httpTransport HTTPTransportHelper
}

func NewOauthClientHTTPHelper(accessToken string, clientID string) *OauthClientHTTPHelper {
	return &OauthClientHTTPHelper{
		accessToken:   accessToken,
		clientID:      clientID,
		httpTransport: NewHTTPTransportHelper(),
	}
}

func (h OauthClientHTTPHelper) Get(url string) ([]byte, error) {
	return h.performHTTPAction(
		url,
		"GET",
		nil,
		nil,
	)
}

func (h OauthClientHTTPHelper) Put(url string, body []byte) ([]byte, error) {
	return h.performHTTPAction(
		url,
		"PUT",
		body,
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	)
}

func (h OauthClientHTTPHelper) Post(url string, body []byte) ([]byte, error) {
	return h.performHTTPAction(
		url,
		"POST",
		body,
		map[string]string{"Content-Type": "application/json"})
}

func (h OauthClientHTTPHelper) Patch(url string, body []byte) ([]byte, error) {
	return h.performHTTPAction(
		url,
		"PATCH",
		body,
		map[string]string{"Content-Type": "application/json"})
}

func (h OauthClientHTTPHelper) Delete(url string) error {
	_, err := h.performHTTPAction(
		url,
		"DELETE",
		nil,
		nil)
	return err
}

func (h OauthClientHTTPHelper) performHTTPAction(
	url string,
	action string,
	body []byte,
	headers map[string]string,
) ([]byte, error) {

	req, err := h.httpTransport.NewRequest(action, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Access-Token", h.accessToken)
	req.Header.Add("X-Client-ID", h.clientID)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if body != nil && len(body) != 0 {
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}

	resp, err := h.httpTransport.DoRequest(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("Nil response returned")
	}

	if resp.Body == nil {
		return []byte{}, nil
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

