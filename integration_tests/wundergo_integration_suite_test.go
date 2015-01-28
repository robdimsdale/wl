package wundergo_integration_test

import (
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"

	"testing"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wundergo Integration Test Suite")
}

var (
	client wundergo.Client
)

func listContains(lists *[]wundergo.List, list *wundergo.List) bool {
	if lists == nil || list == nil {
		return false
	}

	for _, l := range *lists {
		if l.ID == list.ID {
			return true
		}
	}
	return false
}

func taskContains(tasks *[]wundergo.Task, task *wundergo.Task) bool {
	if tasks == nil || task == nil {
		return false
	}
	for _, t := range *tasks {
		if t.ID == task.ID {
			return true
		}
	}
	return false
}

func positionsContainValue(position *[]wundergo.Position, id uint) bool {
	if position == nil {
		return false
	}

	for _, p := range *position {
		if positionContainsValue(&p, id) {
			return true
		}
	}
	return false
}

func positionContainsValue(position *wundergo.Position, id uint) bool {
	if position == nil {
		return false
	}

	for _, v := range position.Values {
		if v == id {
			return true
		}
	}
	return false
}

func taskCommentsContain(taskComments *[]wundergo.TaskComment, taskComment *wundergo.TaskComment) bool {
	if taskComments == nil || taskComment == nil {
		return false
	}
	for _, t := range *taskComments {
		if t.ID == taskComment.ID {
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
