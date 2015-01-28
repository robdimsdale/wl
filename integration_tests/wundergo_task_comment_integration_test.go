package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("Basic task comment functionality", func() {
	It("correctly creates and deletes a task comment", func() {
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

		taskComment, err := client.CreateTaskComment("someText", task.ID)
		Expect(err).NotTo(HaveOccurred())

		taskCommentsForList, err := client.TaskCommentsForListID(list.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentsContain(taskCommentsForList, taskComment)).To(BeTrue())

		taskCommentsForTask, err := client.TaskCommentsForTaskID(task.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentsContain(taskCommentsForTask, taskComment)).To(BeTrue())

		taskCommentAgain, err := client.TaskComment(taskComment.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentAgain.ID).To(Equal(taskComment.ID))

		err = client.DeleteTaskComment(*taskComment)
		Expect(err).NotTo(HaveOccurred())

		taskCommentsForList, err = client.TaskCommentsForListID(list.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentsContain(taskCommentsForList, taskComment)).To(BeFalse())

		taskCommentsForTask, err = client.TaskCommentsForTaskID(task.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(taskCommentsContain(taskCommentsForTask, taskComment)).To(BeFalse())

		// Delete task

		Eventually(func() error {
			task, err = client.Task(task.ID)
			return client.DeleteTask(*task)
		}).Should(Succeed())

		var tasks *[]wundergo.Task
		Eventually(func() (bool, error) {
			tasks, err = client.TasksForListID(list.ID)
			return taskContains(tasks, task), err
		}).Should(BeFalse())
	})
})
