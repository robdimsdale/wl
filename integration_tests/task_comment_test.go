package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("basic task comment functionality", func() {
	It("correctly creates and deletes a task comment", func() {

		By("Creating a task")
		var lists []wundergo.List
		Eventually(func() error {
			l, err := client.Lists()
			lists = l
			return err
		}).Should(Succeed())
		list := lists[0]

		uuid, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle := uuid.String()

		var task wundergo.Task
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

		By("Creating an associated task comment")
		taskComment, err := client.CreateTaskComment("someText", task.ID)
		Expect(err).NotTo(HaveOccurred())

		By("Verifying task comment is present in all task comments")
		taskComments, err := client.TaskComments()
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentsContain(taskComments, taskComment)).To(BeTrue())

		By("Verifying task comment is present in task comments for list")
		taskCommentsForList, err := client.TaskCommentsForListID(list.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentsContain(taskCommentsForList, taskComment)).To(BeTrue())

		By("Verifying task comment is present in task comments for task")
		taskCommentsForTask, err := client.TaskCommentsForTaskID(task.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentsContain(taskCommentsForTask, taskComment)).To(BeTrue())

		By("Getting task comment")
		taskCommentAgain, err := client.TaskComment(taskComment.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentAgain.ID).To(Equal(taskComment.ID))

		By("Deleting task comment")
		err = client.DeleteTaskComment(taskComment)
		Expect(err).NotTo(HaveOccurred())

		By("Verifying task comment is not present in task comments for list")
		taskCommentsForList, err = client.TaskCommentsForListID(list.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentsContain(taskCommentsForList, taskComment)).To(BeFalse())

		By("Verifying task comment is not present in task comments for task")
		taskCommentsForTask, err = client.TaskCommentsForTaskID(task.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentsContain(taskCommentsForTask, taskComment)).To(BeFalse())

		By("Deleting task (and hence associated subtasks)")
		Eventually(func() error {
			task, err = client.Task(task.ID)
			return client.DeleteTask(task)
		}).Should(Succeed())

		By("Verifying task is not present in tasks for list")
		var tasks []wundergo.Task
		Eventually(func() (bool, error) {
			tasks, err = client.TasksForListID(list.ID)
			return taskContains(tasks, task), err
		}).Should(BeFalse())
	})
})
