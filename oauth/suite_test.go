package oauth_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/fakes"
	"github.com/robdimsdale/wundergo/oauth"

	"testing"
)

const (
	dummyAccessToken = "dummyAccessToken"
	dummyClientID    = "dummyClientID"
)

var (
	fakeLogger fakes.FakeLogger

	client wundergo.Client

	server *ghttp.Server
	apiURL string
)

func TestWundergo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wundergo Suite")
}

var _ = BeforeEach(func() {
	fakeLogger = fakes.FakeLogger{}

	server = ghttp.NewServer()
	apiURL = server.URL()
})

var _ = AfterEach(func() {
	server.Close()
})

var initializeClient = func() {
	oauth.NewLogger = func() wundergo.Logger {
		return &fakeLogger
	}

	client = oauth.NewClient(
		dummyAccessToken,
		dummyClientID,
		apiURL,
	)
}
