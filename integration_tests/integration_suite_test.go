package wundergo_integration_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/logger"
	"github.com/robdimsdale/wundergo/oauth"

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
	RunSpecs(t, "Wundergo Integration Test Suite")
}

var (
	client    wundergo.Client
	wlBinPath string

	wlAccessToken string
	wlClientID    string
)

func listContains(lists []wundergo.List, list wundergo.List) bool {
	for _, l := range lists {
		if l.ID == list.ID {
			return true
		}
	}
	return false
}

func taskContains(tasks []wundergo.Task, task wundergo.Task) bool {
	for _, t := range tasks {
		if t.ID == task.ID {
			return true
		}
	}
	return false
}

func positionsContainValue(position []wundergo.Position, id uint) bool {
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

func positionContainsValue(position wundergo.Position, id uint) bool {
	for _, v := range position.Values {
		if v == id {
			return true
		}
	}
	return false
}

func taskCommentsContain(taskComments []wundergo.TaskComment, taskComment wundergo.TaskComment) bool {
	for _, t := range taskComments {
		if t.ID == taskComment.ID {
			return true
		}
	}
	return false
}

var _ = BeforeSuite(func() {
	By("Compiling binary")
	var err error
	wlBinPath, err = gexec.Build("github.com/robdimsdale/wundergo/cmd/wl", "-race")
	Expect(err).ShouldNot(HaveOccurred())

	SetDefaultEventuallyTimeout(5 * time.Second)

	By("Obtaining credentials from environment")
	wlAccessToken = os.Getenv(wlAccessTokenEnvKey)
	wlClientID = os.Getenv(wlClientIDEnvKey)

	if wlAccessToken == "" {
		Fail(fmt.Sprintf("Error - %s must be provided", wlAccessTokenEnvKey))
	}

	if wlClientID == "" {
		Fail(fmt.Sprintf("Error - %s must be provided", wlClientIDEnvKey))
	}

	By("Creating client")
	testLogger := logger.NewTestLogger(GinkgoWriter)
	client = oauth.NewClient(
		wlAccessToken,
		wlClientID,
		apiURL,
		testLogger,
	)

	By("Logging in")
	_, err = client.User()
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
