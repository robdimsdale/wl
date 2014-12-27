package wundergo_test

import (
	"log"
	"os"
	"time"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

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
		It("performs CRUD for lists", func() {
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
			}).Should(Succeed())

			var newList *wundergo.List
			Eventually(func() error {
				newList, err = client.CreateList(newListTitle1)
				return err
			}).Should(Succeed())

			var newLists *[]wundergo.List
			Eventually(func() error {
				newLists, err = client.Lists()
				return err
			}).Should(Succeed())
			Expect(listContains(newLists, newList)).To(BeTrue())

			newList.Title = newListTitle2
			var updatedList *wundergo.List
			Eventually(func() error {
				updatedList, err = client.UpdateList(*newList)
				return err
			}).Should(Succeed())
			newList.Revision = newList.Revision + 1
			Expect(updatedList).To(Equal(newList))

			Eventually(func() error {
				newList, err = client.List(newList.ID)
				return err
			}).Should(Succeed())

			Eventually(func() error {
				return client.DeleteList(*newList)
			}).Should(Succeed())

			Eventually(client.Lists, SERVER_CONSISTENCY_TIMEOUT, POLLING_INTERVAL).Should(Equal(originalLists))
		})
	})

	Describe("Basic task functionality", func() {
		It("can perform CRUD for tasks", func() {
			var lists []wundergo.List
			Eventually(func() error {
				l, err := client.Lists()
				lists = *l
				return err
			}).Should(Succeed())
			list := lists[0]

			uuid, err := uuid.NewV4()
			Expect(err).NotTo(HaveOccurred())
			newTaskTitle := uuid.String()

			var originalTasks *[]wundergo.Task
			Eventually(func() error {
				originalTasks, err = client.TasksForListID(list.ID)
				return err
			}).ShouldNot(HaveOccurred())

			var task *wundergo.Task
			Eventually(func() error {
				task, err = client.CreateTask(
					newTaskTitle,
					list.ID,
					0,
					false,
					"",
					0,
					"1970-01-01",
					false,
				)
				return err
			}).ShouldNot(HaveOccurred())

			var completedTasks *[]wundergo.Task
			Eventually(func() error {
				completed := true
				completedTasks, err = client.CompletedTasksForListID(list.ID, completed)
				return err
			}).Should(Succeed())
			Expect(taskContains(completedTasks, task)).To(BeFalse())

			task.DueDate = "1971-01-01"
			task.Completed = true
			Eventually(func() error {
				task, err = client.UpdateTask(*task)
				return err
			}).Should(Succeed())

			Eventually(func() error {
				completed := true
				completedTasks, err = client.CompletedTasksForListID(list.ID, completed)
				return err
			}).ShouldNot(HaveOccurred())
			Expect(taskContains(completedTasks, task)).To(BeTrue())

			var note *wundergo.Note
			Eventually(func() error {
				note, err = client.CreateNote("myContent", task.ID)
				return err
			}).Should(Succeed())

			note.Content = "newContent"
			Eventually(func() error {
				note, err = client.UpdateNote(*note)
				return err
			}).Should(Succeed())

			Eventually(func() error {
				return client.DeleteNote(*note)
			}).Should(Succeed())

			Eventually(func() error {
				task, err = client.Task(task.ID)
				return err
			}).Should(Succeed())

			Eventually(func() error {
				return client.DeleteTask(*task)
			}).Should(Succeed())

			Eventually(func() (*[]wundergo.Task, error) {
				return client.TasksForListID(list.ID)
			}, SERVER_CONSISTENCY_TIMEOUT, POLLING_INTERVAL).Should(Equal(originalTasks))
		})
	})
})
