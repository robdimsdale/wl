package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("Basic task functionality", func() {
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
		Eventually(func() (bool, error) {
			completed := true
			completedTasks, err = client.CompletedTasksForListID(list.ID, completed)
			return taskContains(completedTasks, task), err
		}).Should(BeFalse())

		task.DueDate = "1971-01-01"
		task.Completed = true
		Eventually(func() error {
			task, err = client.UpdateTask(*task)
			return err
		}).Should(Succeed())

		completed := true
		Eventually(func() (bool, error) {
			completedTasks, err = client.CompletedTasksForListID(list.ID, completed)
			return taskContains(completedTasks, task), err
		}).Should(BeTrue())

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

		var subtask *wundergo.Subtask
		subtaskComplete := false
		Eventually(func() error {
			subtask, err = client.CreateSubtask("mySubtaskTitle", task.ID, subtaskComplete)
			return err
		}).Should(Succeed())

		subtask.Title = "newSubtaskTitle"
		Eventually(func() error {
			subtask, err = client.UpdateSubtask(*subtask)
			return err
		}).Should(Succeed())

		Eventually(func() error {
			return client.DeleteSubtask(*subtask)
		}).Should(Succeed())

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
