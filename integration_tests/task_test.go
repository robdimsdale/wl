package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("basic task functionality", func() {
	var (
		newList wundergo.List
		newTask wundergo.Task
		err     error
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
			newTask, err = client.Task(newTask.ID)
			return client.DeleteTask(newTask)
		}).Should(Succeed())

		var tasks []wundergo.Task
		Eventually(func() (bool, error) {
			tasks, err = client.TasksForListID(newList.ID)
			return taskContains(tasks, newTask), err
		}).Should(BeFalse())

		By("Deleting new list")
		Eventually(func() error {
			newList, err = client.List(newList.ID)
			return client.DeleteList(newList)
		}).Should(Succeed())

		var lists []wundergo.List
		Eventually(func() (bool, error) {
			lists, err = client.Lists()
			return listContains(lists, newList), err
		}).Should(BeFalse())
	})

	Describe("moving a tasks between lists", func() {
		var secondList wundergo.List

		BeforeEach(func() {
			By("Creating a second list")
			uuid1, err := uuid.NewV4()
			Expect(err).NotTo(HaveOccurred())
			secondListTitle := uuid1.String()

			secondList, err = client.CreateList(secondListTitle)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			By("Deleting second list")
			Eventually(func() error {
				secondList, err = client.List(secondList.ID)
				return client.DeleteList(secondList)
			}).Should(Succeed())

			var lists []wundergo.List
			Eventually(func() (bool, error) {
				lists, err = client.Lists()
				return listContains(lists, secondList), err
			}).Should(BeFalse())
		})

		It("can move a task between lists", func() {
			By("Moving task to second list")
			newTask.ListID = secondList.ID
			Eventually(func() error {
				newTask, err = client.UpdateTask(newTask)
				return err
			}).Should(Succeed())

			By("Verifying task appears in tasks for second list")
			var completedTasksForSecondList []wundergo.Task
			Eventually(func() (bool, error) {
				showCompletedTasks := false
				completedTasksForSecondList, err = client.CompletedTasksForListID(secondList.ID, showCompletedTasks)
				return taskContains(completedTasksForSecondList, newTask), err
			}).Should(BeTrue())

			By("Verifying task does not appear in tasks for first list")
			var completedTasksForFirstList []wundergo.Task
			Eventually(func() (bool, error) {
				showCompletedTasks := false
				completedTasksForFirstList, err = client.CompletedTasksForListID(newList.ID, showCompletedTasks)
				return taskContains(completedTasksForFirstList, newTask), err
			}).Should(BeFalse())

			By("Moving task back to first list")
			newTask.ListID = newList.ID
			Eventually(func() error {
				newTask, err = client.UpdateTask(newTask)
				return err
			}).Should(Succeed())

			By("Verifying task does not appear in tasks for second list")
			Eventually(func() (bool, error) {
				showCompletedTasks := false
				completedTasksForSecondList, err = client.CompletedTasksForListID(secondList.ID, showCompletedTasks)
				return taskContains(completedTasksForSecondList, newTask), err
			}).Should(BeFalse())

			By("Verifying task does appear in tasks for first list")
			Eventually(func() (bool, error) {
				showCompletedTasks := false
				completedTasksForFirstList, err = client.CompletedTasksForListID(newList.ID, showCompletedTasks)
				return taskContains(completedTasksForFirstList, newTask), err
			}).Should(BeTrue())
		})
	})

	It("can complete tasks", func() {
		var completedTasksForList []wundergo.Task
		showCompletedTasks := true
		Eventually(func() (bool, error) {
			completedTasksForList, err = client.CompletedTasksForListID(newList.ID, showCompletedTasks)
			return taskContains(completedTasksForList, newTask), err
		}).Should(BeFalse())

		By("Completing task")
		newTask.Completed = true
		Eventually(func() error {
			newTask, err = client.UpdateTask(newTask)
			return err
		}).Should(Succeed())

		By("Verifying task appears in completed tasks for list")
		Eventually(func() (bool, error) {
			completedTasksForList, err = client.CompletedTasksForListID(newList.ID, showCompletedTasks)
			return taskContains(completedTasksForList, newTask), err
		}).Should(BeTrue())

		By("Verifying task appears in completed tasks")
		var completedTasks []wundergo.Task
		Eventually(func() (bool, error) {
			completedTasks, err = client.CompletedTasks(showCompletedTasks)
			return taskContains(completedTasks, newTask), err
		}).Should(BeTrue())
	})

	It("can update tasks", func() {
		By("Setting properties")
		newTask.DueDate = "1971-01-01"
		newTask.Starred = true
		newTask.Completed = true
		newTask.RecurrenceType = "week"
		newTask.RecurrenceCount = 2

		By("Updating task")
		Eventually(func() error {
			newTask, err = client.UpdateTask(newTask)
			return err
		}).Should(Succeed())

		By("Getting task again")
		var taskAgain wundergo.Task
		Eventually(func() error {
			taskAgain, err = client.Task(newTask.ID)
			return err
		}).Should(Succeed())

		By("Ensuring properties are set")
		Expect(taskAgain.DueDate).Should(Equal("1971-01-01"))
		Expect(taskAgain.Starred).Should(BeTrue())
		Expect(taskAgain.Completed).Should(BeTrue())
		Expect(taskAgain.RecurrenceType).Should(Equal("week"))
		Expect(taskAgain.RecurrenceCount).Should(Equal(uint(2)))

		By("Resetting properties")
		taskAgain.DueDate = ""
		taskAgain.Starred = false
		taskAgain.Completed = false
		taskAgain.RecurrenceType = ""
		taskAgain.RecurrenceCount = 0

		By("Updating task")
		Eventually(func() error {
			taskAgain, err = client.UpdateTask(taskAgain)
			return err
		}).Should(Succeed())

		By("Verifying properties are reset")
		Expect(taskAgain.DueDate).Should(Equal(""))
		Expect(taskAgain.Starred).Should(BeFalse())
		Expect(taskAgain.Completed).Should(BeFalse())
		Expect(taskAgain.RecurrenceType).Should(Equal(""))
		Expect(taskAgain.RecurrenceCount).Should(Equal(uint(0)))
	})

	It("can perform subtask CRUD", func() {
		By("Creating subtask")
		var subtask wundergo.Subtask
		subtaskComplete := false
		Eventually(func() error {
			subtask, err = client.CreateSubtask("mySubtaskTitle", newTask.ID, subtaskComplete)
			return err
		}).Should(Succeed())

		By("Getting subtask")
		var aSubtask wundergo.Subtask
		Eventually(func() error {
			aSubtask, err = client.Subtask(subtask.ID)
			return err
		}).Should(Succeed())
		Expect(aSubtask.ID).To(Equal(subtask.ID))

		By("Validating subtask exists in all subtasks")
		Eventually(func() bool {
			// It is statistically probable that one of the lists will
			// be deleted, so we ignore error here.
			subtasks, _ := client.Subtasks()
			return subtaskContains(subtasks, subtask)
		}).Should(BeTrue())

		By("Validating subtask exists in subtasks for list")
		Eventually(func() (bool, error) {
			subtasksForList, err := client.SubtasksForListID(newList.ID)
			return subtaskContains(subtasksForList, subtask), err
		}).Should(BeTrue())

		By("Validating subtask exists in subtasks for task")
		Eventually(func() (bool, error) {
			subtasksForTask, err := client.SubtasksForTaskID(newTask.ID)
			return subtaskContains(subtasksForTask, subtask), err
		}).Should(BeTrue())

		By("Completing subtask")
		subtask.Completed = true
		Eventually(func() error {
			subtask, err = client.UpdateSubtask(subtask)
			return err
		}).Should(Succeed())

		By("Validating subtask exists in all completed subtasks")
		showCompletedSubtasks := true
		Eventually(func() bool {
			// It is statistically probable that one of the lists will
			// be deleted, so we ignore error here.
			subtasks, _ := client.CompletedSubtasks(showCompletedSubtasks)
			return subtaskContains(subtasks, subtask)
		}).Should(BeTrue())

		By("Validating subtask exists in completed subtasks for list")
		Eventually(func() (bool, error) {
			subtasksForList, err := client.CompletedSubtasksForListID(newList.ID, showCompletedSubtasks)
			return subtaskContains(subtasksForList, subtask), err
		}).Should(BeTrue())

		By("Validating subtask exists in completed subtasks for task")
		Eventually(func() (bool, error) {
			subtasksForTask, err := client.CompletedSubtasksForTaskID(newTask.ID, showCompletedSubtasks)
			return subtaskContains(subtasksForTask, subtask), err
		}).Should(BeTrue())

		By("Deleting subtask")
		Eventually(func() error {
			return client.DeleteSubtask(subtask)
		}).Should(Succeed())
	})
})

func subtaskContains(subtasks []wundergo.Subtask, subtask wundergo.Subtask) bool {
	for _, t := range subtasks {
		if t.ID == subtask.ID {
			return true
		}
	}
	return false
}
