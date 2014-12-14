package wundergo

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var NewHTTPTransport = func() HTTPTransport {
	return newDefaultHTTPTransport()
}

type HTTPHelper interface {
	Get(url string) ([]byte, error)
	Post(url string, body string) ([]byte, error)
	Put(url string, body string) ([]byte, error)
	Patch(url string, body string) ([]byte, error)
}

type OauthClientHTTPHelper struct {
	accessToken   string
	clientID      string
	httpTransport HTTPTransport
}

func NewOauthClientHTTPHelper(accessToken string, clientID string) *OauthClientHTTPHelper {
	return &OauthClientHTTPHelper{
		accessToken:   accessToken,
		clientID:      clientID,
		httpTransport: NewHTTPTransport(),
	}
}

func (h OauthClientHTTPHelper) Get(url string) ([]byte, error) {
	return h.performHTTPAction(
		url,
		"GET",
		"",
		nil,
	)
}

func (h OauthClientHTTPHelper) Put(url string, body string) ([]byte, error) {
	return h.performHTTPAction(
		url,
		"PUT",
		body,
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	)
}

func (h OauthClientHTTPHelper) Post(url string, body string) ([]byte, error) {
	return h.performHTTPAction(
		url,
		"POST",
		body,
		map[string]string{"Content-Type": "application/json"})
}

func (h OauthClientHTTPHelper) Patch(url string, body string) ([]byte, error) {
	return h.performHTTPAction(
		url,
		"PATCH",
		body,
		map[string]string{"Content-Type": "application/json"})
}

func (h OauthClientHTTPHelper) performHTTPAction(
	url string,
	action string,
	body string,
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

	if body != "" {
		req.Body = ioutil.NopCloser(strings.NewReader(body))
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

type HTTPTransport interface {
	NewRequest(method, urlStr string, body io.Reader) (*http.Request, error)
	DoRequest(req *http.Request) (resp *http.Response, err error)
}

type defaultHTTPTransport struct {
	client http.Client
}

func newDefaultHTTPTransport() *defaultHTTPTransport {
	return &defaultHTTPTransport{
		client: http.Client{},
	}
}

func (h defaultHTTPTransport) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, urlStr, body)
}

func (h defaultHTTPTransport) DoRequest(req *http.Request) (resp *http.Response, err error) {
	return h.client.Do(req)
}
