package wundergo_integration_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/logger"
	"github.com/robdimsdale/wundergo/oauth"

	"testing"
	"time"
)

const (
	apiURL = "https://a.wunderlist.com/api/v1"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wundergo Integration Test Suite")
}

var (
	client wundergo.Client
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
	SetDefaultEventuallyTimeout(5 * time.Second)

	By("Logging in")

	accessToken := os.Getenv("WL_ACCESS_TOKEN")
	clientID := os.Getenv("WL_CLIENT_ID")

	if accessToken == "" {
		Fail("Error - WL_ACCESS_TOKEN must be provided")
	}

	if clientID == "" {
		Fail("Error - WL_CLIENT_ID must be provided")
	}

	logger := logger.NewTestLogger(GinkgoWriter)
	client = oauth.NewClient(
		accessToken,
		clientID,
		apiURL,
		logger,
	)
	_, err := client.Lists()
	Expect(err).To(BeNil())
})
