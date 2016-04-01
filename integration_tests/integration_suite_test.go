package wl_integration_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/robdimsdale/wl"
	"github.com/robdimsdale/wl/logger"
	"github.com/robdimsdale/wl/oauth"

	"testing"
	"time"
)

const (
	apiURL = "https://a.wunderlist.com/api/v1"

	wlAccessTokenEnvKey = "WL_ACCESS_TOKEN"
	wlClientIDEnvKey    = "WL_CLIENT_ID"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WL Integration Test Suite")
}

var (
	client    wl.Client
	wlBinPath string

	wlAccessToken string
	wlClientID    string
)

func listContains(lists []wl.List, list wl.List) bool {
	for _, l := range lists {
		if l.ID == list.ID {
			return true
		}
	}
	return false
}

func taskContains(tasks []wl.Task, task wl.Task) bool {
	for _, t := range tasks {
		if t.ID == task.ID {
			return true
		}
	}
	return false
}

func positionsContainValue(position []wl.Position, id uint) bool {
	if position == nil {
		return false
	}

	for _, p := range position {
		if positionContainsValue(p, id) {
			return true
		}
	}
	return false
}

func positionContainsValue(position wl.Position, id uint) bool {
	for _, v := range position.Values {
		if v == id {
			return true
		}
	}
	return false
}

var _ = BeforeSuite(func() {
	By("Setting Eventually defaults")
	SetDefaultEventuallyTimeout(1 * time.Minute)
	SetDefaultEventuallyPollingInterval(1 * time.Second)

	By("Obtaining credentials from environment")
	wlAccessToken = os.Getenv(wlAccessTokenEnvKey)
	wlClientID = os.Getenv(wlClientIDEnvKey)

	if wlAccessToken == "" {
		Fail(fmt.Sprintf("Error - %s must be provided", wlAccessTokenEnvKey))
	}

	if wlClientID == "" {
		Fail(fmt.Sprintf("Error - %s must be provided", wlClientIDEnvKey))
	}

	By("Compiling binary")
	var err error
	wlBinPath, err = gexec.Build("github.com/robdimsdale/wl/cmd/wl", "-race")
	Expect(err).ShouldNot(HaveOccurred())

	By("Creating client")
	testLogger := logger.NewTestLogger(GinkgoWriter)
	client = oauth.NewClient(
		wlAccessToken,
		wlClientID,
		apiURL,
		testLogger,
	)

	By("Logging in")
	Eventually(func() error {
		_, err := client.User()
		return err
	}).Should(Succeed())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
