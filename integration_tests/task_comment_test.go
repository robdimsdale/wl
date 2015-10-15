package wl_integration_test

import (
	"time"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
)

var _ = Describe("basic task comment functionality", func() {
	var (
		newList wl.List
		newTask wl.Task
	)

	BeforeEach(func() {
		By("Creating a new list")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle := uuid1.String()

		Eventually(func() error {
			newList, err = client.CreateList(newListTitle)
			return err
		}).Should(Succeed())

		By("Creating task in new list")
		uuid, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle := uuid.String()

		Eventually(func() error {
			newTask, err = client.CreateTask(
				newTaskTitle,
				newList.ID,
				0,
				false,
				"",
				0,
				time.Date(1971, 12, 31, 0, 0, 0, 0, time.UTC),
				false,
			)
			return err
		}).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Deleting task")
		Eventually(func() error {
			t, err := client.Task(newTask.ID)
			if err != nil {
				return err
			}
			return client.DeleteTask(t)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			tasks, err := client.TasksForListID(newList.ID)
			return taskContains(tasks, newTask), err
		}).Should(BeFalse())

		By("Deleting new list")
		Eventually(func() error {
			l, err := client.List(newList.ID)
			if err != nil {
				return err
			}
			return client.DeleteList(l)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			lists, err := client.Lists()
			return listContains(lists, newList), err
		}).Should(BeFalse())
	})

	It("correctly creates and deletes a task comment", func() {
		By("Creating a task comment")
		var err error
		var taskComment wl.TaskComment
		Eventually(func() error {
			taskComment, err = client.CreateTaskComment("someText", newTask.ID)
			return err
		}).Should(Succeed())

		By("Verifying task comment is present in all task comments")
		Eventually(func() bool {
			// It is statistically probable that one of the lists will
			// be deleted, so we ignore error here.
			taskComments, _ := client.TaskComments()
			return taskCommentContains(taskComments, taskComment)
		}).Should(BeTrue())

		By("Verifying task comment is present in task comments for list")
		Eventually(func() (bool, error) {
			taskCommentsForList, err := client.TaskCommentsForListID(newList.ID)
			return taskCommentContains(taskCommentsForList, taskComment), err
		}).Should(BeTrue())

		By("Verifying task comment is present in task comments for task")
		Eventually(func() (bool, error) {
			taskCommentsForTask, err := client.TaskCommentsForTaskID(newTask.ID)
			return taskCommentContains(taskCommentsForTask, taskComment), err
		}).Should(BeTrue())

		By("Getting task comment")
		var taskCommentAgain wl.TaskComment
		Eventually(func() error {
			taskCommentAgain, err = client.TaskComment(taskComment.ID)
			return err
		}).Should(Succeed())
		Expect(taskCommentAgain.ID).To(Equal(taskComment.ID))

		By("Deleting task comment")
		Eventually(func() error {
			return client.DeleteTaskComment(taskComment)
		}).Should(Succeed())

		By("Verifying task comment is not present in task comments for list")
		Eventually(func() (bool, error) {
			taskCommentsForList, err := client.TaskCommentsForListID(newList.ID)
			return taskCommentContains(taskCommentsForList, taskComment), err
		}).Should(BeFalse())

		By("Verifying task comment is not present in task comments for task")
		Eventually(func() (bool, error) {
			taskCommentsForTask, err := client.TaskCommentsForTaskID(newTask.ID)
			return taskCommentContains(taskCommentsForTask, taskComment), err
		}).Should(BeFalse())
	})
})

func taskCommentContains(taskComments []wl.TaskComment, taskComment wl.TaskComment) bool {
	for _, t := range taskComments {
		if t.ID == taskComment.ID {
			return true
		}
	}
	return false
}
