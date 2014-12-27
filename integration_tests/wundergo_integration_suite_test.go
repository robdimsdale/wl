package wundergo_integration_test

import (
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"

	"testing"
	"time"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wundergo Integration Test Suite")
}

const (
	SERVER_CONSISTENCY_TIMEOUT = 30 * time.Second
	POLLING_INTERVAL           = 10 * time.Millisecond
)

var (
	client wundergo.Client
)

func listContains(lists *[]wundergo.List, list *wundergo.List) bool {
	for _, l := range *lists {
		if l == *list {
			return true
		}
	}
	return false
}

func taskContains(tasks *[]wundergo.Task, task *wundergo.Task) bool {
	for _, t := range *tasks {
		if t == *task {
			return true
		}
	}
	return false
}

var _ = BeforeSuite(func() {
	accessToken := os.Getenv("WL_ACCESS_TOKEN")
	clientID := os.Getenv("WL_CLIENT_ID")

	if accessToken == "" {
		log.Fatal("Error - WL_ACCESS_TOKEN must be provided")
	}

	if clientID == "" {
		log.Fatal("Error - WL_CLIENT_ID must be provided")
	}

	client = wundergo.NewOauthClient(accessToken, clientID)
})
