package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("basic task functionality", func() {

	var firstList *wundergo.List
	var newTask *wundergo.Task
	var err error

	BeforeEach(func() {
		By("Getting first list")
		var lists []wundergo.List
		Eventually(func() error {
			l, err := client.Lists()
			lists = *l
			return err
		}).Should(Succeed())
		firstList = &lists[0]

		By("Creating task in first list")
		uuid, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle := uuid.String()

		Eventually(func() error {
			newTask, err = client.CreateTask(
				newTaskTitle,
				firstList.ID,
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
			newTask, err = client.Task(newTask.ID)
			return client.DeleteTask(*newTask)
		}).Should(Succeed())

		var tasks *[]wundergo.Task
		Eventually(func() (bool, error) {
			tasks, err = client.TasksForListID(firstList.ID)
			return taskContains(tasks, newTask), err
		}).Should(BeFalse())
	})

	It("can update tasks", func() {

		var completedTasks *[]wundergo.Task
		showCompletedTasks := true
		Eventually(func() (bool, error) {
			completedTasks, err = client.CompletedTasksForListID(firstList.ID, showCompletedTasks)
			return taskContains(completedTasks, newTask), err
		}).Should(BeFalse())

		By("Updating task")
		newTask.DueDate = "1971-01-01"
		newTask.Completed = true
		Eventually(func() error {
			newTask, err = client.UpdateTask(*newTask)
			return err
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			completedTasks, err = client.CompletedTasksForListID(firstList.ID, showCompletedTasks)
			return taskContains(completedTasks, newTask), err
		}).Should(BeTrue())
	})

	It("can perform subtask CRUD", func() {
		var subtask *wundergo.Subtask
		subtaskComplete := false
		Eventually(func() error {
			subtask, err = client.CreateSubtask("mySubtaskTitle", newTask.ID, subtaskComplete)
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
	})

	It("can perform reminder CRUD", func() {
		var reminder *wundergo.Reminder
		reminderDate := "1970-08-30T08:29:46.203Z"
		createdByDeviceUdid := ""
		Eventually(func() error {
			reminder, err = client.CreateReminder(reminderDate, newTask.ID, createdByDeviceUdid)
			return err
		}).Should(Succeed())

		reminder.Date = "1971-08-30T08:29:46.203Z"
		Eventually(func() error {
			reminder, err = client.UpdateReminder(*reminder)
			return err
		}).Should(Succeed())

		Eventually(func() error {
			return client.DeleteReminder(*reminder)
		}).Should(Succeed())
	})
})
