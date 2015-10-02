package wl_integration_test

import (
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

		newList, err = client.CreateList(newListTitle)
		Expect(err).NotTo(HaveOccurred())

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
				"1970-01-01",
				false,
			)
			return err
		}).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Deleting task")
		Eventually(func() error {
			newTask, err := client.Task(newTask.ID)
			if err != nil {
				return err
			}
			return client.DeleteTask(newTask)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			tasks, err := client.TasksForListID(newList.ID)
			return taskContains(tasks, newTask), err
		}).Should(BeFalse())

		By("Deleting new list")
		Eventually(func() error {
			newList, err := client.List(newList.ID)
			if err != nil {
				return err
			}
			return client.DeleteList(newList)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			lists, err := client.Lists()
			return listContains(lists, newList), err
		}).Should(BeFalse())
	})

	It("correctly creates and deletes a task comment", func() {
		By("Creating a task comment")
		taskComment, err := client.CreateTaskComment("someText", newTask.ID)
		Expect(err).NotTo(HaveOccurred())

		By("Verifying task comment is present in all task comments")
		Eventually(func() bool {
			// It is statistically probable that one of the lists will
			// be deleted, so we ignore error here.
			taskComments, _ := client.TaskComments()
			return taskCommentContains(taskComments, taskComment)
		}).Should(BeTrue())

		By("Verifying task comment is present in task comments for list")
		taskCommentsForList, err := client.TaskCommentsForListID(newList.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentContains(taskCommentsForList, taskComment)).To(BeTrue())

		By("Verifying task comment is present in task comments for task")
		taskCommentsForTask, err := client.TaskCommentsForTaskID(newTask.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentContains(taskCommentsForTask, taskComment)).To(BeTrue())

		By("Getting task comment")
		taskCommentAgain, err := client.TaskComment(taskComment.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentAgain.ID).To(Equal(taskComment.ID))

		By("Deleting task comment")
		err = client.DeleteTaskComment(taskComment)
		Expect(err).NotTo(HaveOccurred())

		By("Verifying task comment is not present in task comments for list")
		taskCommentsForList, err = client.TaskCommentsForListID(newList.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentContains(taskCommentsForList, taskComment)).To(BeFalse())

		By("Verifying task comment is not present in task comments for task")
		taskCommentsForTask, err = client.TaskCommentsForTaskID(newTask.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentContains(taskCommentsForTask, taskComment)).To(BeFalse())
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
