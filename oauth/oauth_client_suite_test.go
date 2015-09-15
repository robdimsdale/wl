package oauth_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/logger"
	"github.com/robdimsdale/wundergo/oauth"

	"testing"
)

const (
	dummyAccessToken = "dummyAccessToken"
	dummyClientID    = "dummyClientID"
)

var (
	client wundergo.Client

	server *ghttp.Server
	apiURL string

	testLogger logger.Logger
)

func TestWundergo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wundergo Suite")
}

var _ = BeforeEach(func() {
	server = ghttp.NewServer()
	apiURL = server.URL()

	testLogger = logger.NewTestLogger()
	client = oauth.NewClient(
		dummyAccessToken,
		dummyClientID,
		apiURL,
		testLogger,
	)
})

var _ = AfterEach(func() {
	server.Close()
})
