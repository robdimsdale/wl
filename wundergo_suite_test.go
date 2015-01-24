package wundergo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/fakes"

	"testing"
)

const (
	dummyAccessToken = "dummyAccessToken"
	dummyClientID    = "dummyClientID"

	apiURL = "https://a.wunderlist.com/api/v1"
)

var (
	fakeHTTPHelper fakes.FakeHTTPHelper
	fakeLogger     fakes.FakeLogger
	fakeJSONHelper fakes.FakeJSONHelper

	client wundergo.Client
)

func TestWundergo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wundergo Suite")
}

var initializeFakes = func() {
	fakeHTTPHelper = fakes.FakeHTTPHelper{}
	fakeLogger = fakes.FakeLogger{}
	fakeJSONHelper = fakes.FakeJSONHelper{}
}

var initializeClient = func() {
	wundergo.NewHTTPHelper = func(accessToken string, clientID string) wundergo.HTTPHelper {
		return &fakeHTTPHelper
	}

	wundergo.NewLogger = func() wundergo.Logger {
		return &fakeLogger
	}

	wundergo.NewJSONHelper = func() wundergo.JSONHelper {
		return &fakeJSONHelper
	}

	client = wundergo.NewOauthClient(dummyAccessToken, dummyClientID)
}

type erroringReadCloser struct {
	readError  error
	closeError error
}

func (e erroringReadCloser) Read([]byte) (int, error) {
	return 0, e.readError
}

func (e erroringReadCloser) Close() error {
	return e.closeError
}
