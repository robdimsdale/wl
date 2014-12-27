package wundergo_test

import (
	"log"
	"os"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var (
	client wundergo.Client
)

func contains(lists *[]wundergo.List, list *wundergo.List) bool {
	for _, l := range *lists {
		if l == *list {
			return true
		}
	}
	return false
}

var _ = Describe("Wundergo library", func() {
	BeforeEach(func() {
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

	Describe("Basic list functionality", func() {
		It("creates, updates and deletes new list", func() {
			uuid1, err := uuid.NewV4()
			Expect(err).NotTo(HaveOccurred())
			newListTitle1 := uuid1.String()

			uuid2, err := uuid.NewV4()
			Expect(err).NotTo(HaveOccurred())
			newListTitle2 := uuid2.String()

			var originalLists *[]wundergo.List
			Eventually(func() error {
				originalLists, err = client.Lists()
				return err
			}).ShouldNot(HaveOccurred())

			var newList *wundergo.List
			Eventually(func() error {
				newList, err = client.CreateList(newListTitle1)
				return err
			}).ShouldNot(HaveOccurred())

			var newLists *[]wundergo.List
			Eventually(func() error {
				newLists, err = client.Lists()
				return err
			}).ShouldNot(HaveOccurred())
			Expect(contains(newLists, newList)).To(BeTrue())

			newList.Title = newListTitle2
			var updatedList *wundergo.List
			Eventually(func() error {
				updatedList, err = client.UpdateList(*newList)
				return err
			}).ShouldNot(HaveOccurred())
			newList.Revision = newList.Revision + 1
			Expect(updatedList).To(Equal(newList))

			Eventually(func() error {
				_, err = client.TasksForListID(newList.ID)
				return err
			}).ShouldNot(HaveOccurred())

			var task *wundergo.Task
			Eventually(func() error {
				task, err = client.CreateTask(
					"myTask",
					newList.ID,
					0,
					false,
					"",
					0,
					"1970-01-01",
					false,
				)
				return err
			}).ShouldNot(HaveOccurred())
			newList.Revision = newList.Revision + 1

			task.DueDate = "1971-01-01"
			Eventually(func() error {
				task, err = client.UpdateTask(*task)
				return err
			}).ShouldNot(HaveOccurred())
			newList.Revision = newList.Revision + 1

			Eventually(func() error {
				_, err := client.CreateNote("myContent", task.ID)
				return err
			}).ShouldNot(HaveOccurred())
			newList.Revision = newList.Revision + 1

			Eventually(func() error {
				return client.DeleteList(*newList)
			}).ShouldNot(HaveOccurred())

			var afterDeleteLists *[]wundergo.List
			Eventually(func() error {
				afterDeleteLists, err = client.Lists()
				return err
			}).ShouldNot(HaveOccurred())

			Expect(afterDeleteLists).To(Equal(originalLists))
		})
	})
})
