package wl_integration_test

import (
	"time"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
)

var _ = Describe("basic task functionality", func() {
	var (
		newList wl.List
		newTask wl.Task
		err     error
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

		var tasks []wl.Task
		Eventually(func() (bool, error) {
			tasks, err = client.TasksForListID(newList.ID)
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

		var lists []wl.List
		Eventually(func() (bool, error) {
			lists, err = client.Lists()
			return listContains(lists, newList), err
		}).Should(BeFalse())
	})

	Describe("moving a tasks between lists", func() {
		var secondList wl.List

		BeforeEach(func() {
			By("Creating a second list")
			uuid1, err := uuid.NewV4()
			Expect(err).NotTo(HaveOccurred())
			secondListTitle := uuid1.String()

			Eventually(func() error {
				secondList, err = client.CreateList(secondListTitle)
				return err
			}).Should(Succeed())
		})

		AfterEach(func() {
			By("Deleting second list")
			Eventually(func() error {
				l, err := client.List(secondList.ID)
				if err != nil {
					return err
				}
				return client.DeleteList(l)
			}).Should(Succeed())

			var lists []wl.List
			Eventually(func() (bool, error) {
				lists, err = client.Lists()
				return listContains(lists, secondList), err
			}).Should(BeFalse())
		})

		It("can move a task between lists", func() {
			By("Moving task to second list")
			newTask.ListID = secondList.ID
			var t wl.Task
			Eventually(func() error {
				t, err = client.UpdateTask(newTask)
				return err
			}).Should(Succeed())
			newTask = t

			By("Verifying task appears in tasks for second list")
			var completedTasksForSecondList []wl.Task
			Eventually(func() (bool, error) {
				showCompletedTasks := false
				completedTasksForSecondList, err =
					client.CompletedTasksForListID(secondList.ID, showCompletedTasks)
				return taskContains(completedTasksForSecondList, newTask), err
			}).Should(BeTrue())

			By("Verifying task does not appear in tasks for first list")
			var completedTasksForFirstList []wl.Task
			Eventually(func() (bool, error) {
				showCompletedTasks := false
				completedTasksForFirstList, err =
					client.CompletedTasksForListID(newList.ID, showCompletedTasks)
				return taskContains(completedTasksForFirstList, newTask), err
			}).Should(BeFalse())

			By("Moving task back to first list")
			newTask.ListID = newList.ID
			Eventually(func() error {
				t, err = client.UpdateTask(newTask)
				return err
			}).Should(Succeed())
			newTask = t

			By("Verifying task does not appear in tasks for second list")
			Eventually(func() (bool, error) {
				showCompletedTasks := false
				completedTasksForSecondList, err =
					client.CompletedTasksForListID(secondList.ID, showCompletedTasks)
				return taskContains(completedTasksForSecondList, newTask), err
			}).Should(BeFalse())

			By("Verifying task does appear in tasks for first list")
			Eventually(func() (bool, error) {
				showCompletedTasks := false
				completedTasksForFirstList, err =
					client.CompletedTasksForListID(newList.ID, showCompletedTasks)
				return taskContains(completedTasksForFirstList, newTask), err
			}).Should(BeTrue())
		})
	})

	It("can complete tasks", func() {
		var completedTasksForList []wl.Task
		By("Ensuring task is not already in completed tasks")
		showCompletedTasks := true
		Eventually(func() (bool, error) {
			completedTasksForList, err =
				client.CompletedTasksForListID(newList.ID, showCompletedTasks)
			return taskContains(completedTasksForList, newTask), err
		}).Should(BeFalse())

		By("Completing task")
		newTask.Completed = true
		var t wl.Task
		Eventually(func() error {
			t, err = client.UpdateTask(newTask)
			return err
		}).Should(Succeed())
		newTask = t

		By("Verifying task appears in completed tasks for list")
		Eventually(func() (bool, error) {
			completedTasksForList, err =
				client.CompletedTasksForListID(newList.ID, showCompletedTasks)
			return taskContains(completedTasksForList, newTask), err
		}).Should(BeTrue())

		By("Verifying task appears in completed tasks")
		var completedTasks []wl.Task
		Eventually(func() bool {
			// It is statistically probable that one of the lists will
			// be deleted, so we ignore error here.
			completedTasks, _ = client.CompletedTasks(showCompletedTasks)
			return taskContains(completedTasks, newTask)
		}).Should(BeTrue())
	})

	It("can update tasks", func() {
		By("Setting properties")
		newTask.Starred = true
		newTask.Completed = true
		newTask.RecurrenceType = "week"
		newTask.RecurrenceCount = 2

		By("Updating task")
		var t wl.Task
		Eventually(func() error {
			t, err = client.UpdateTask(newTask)
			return err
		}).Should(Succeed())
		newTask = t

		By("Getting task again")
		var taskAgain wl.Task
		Eventually(func() error {
			taskAgain, err = client.Task(newTask.ID)
			return err
		}).Should(Succeed())

		By("Ensuring properties are set")
		Expect(taskAgain.Starred).Should(BeTrue())
		Expect(taskAgain.Completed).Should(BeTrue())
		Expect(taskAgain.RecurrenceType).Should(Equal("week"))
		Expect(taskAgain.RecurrenceCount).Should(Equal(uint(2)))

		By("Resetting properties")
		taskAgain.Starred = false
		taskAgain.Completed = false
		taskAgain.RecurrenceType = ""
		taskAgain.RecurrenceCount = 0

		By("Updating task")
		Eventually(func() error {
			t, err = client.UpdateTask(taskAgain)
			return err
		}).Should(Succeed())
		taskAgain = t

		By("Verifying properties are reset")
		Expect(taskAgain.Starred).Should(BeFalse())
		Expect(taskAgain.Completed).Should(BeFalse())
		Expect(taskAgain.RecurrenceType).Should(Equal(""))
		Expect(taskAgain.RecurrenceCount).Should(Equal(uint(0)))
	})

	It("can update the due date", func() {
		By("Setting properties")
		firstDate := time.Date(1968, 1, 2, 0, 0, 0, 0, time.UTC)
		newTask.DueDate = firstDate

		By("Updating task")
		var t wl.Task
		Eventually(func() error {
			t, err = client.UpdateTask(newTask)
			return err
		}).Should(Succeed())
		newTask = t

		By("Ensuring due date is set")
		Expect(newTask.DueDate).Should(Equal(firstDate))

		By("Updating properties")
		newDate := time.Date(1972, 2, 3, 0, 0, 0, 0, time.UTC)
		newTask.DueDate = newDate

		By("Updating task")
		Eventually(func() error {
			t, err = client.UpdateTask(newTask)
			return err
		}).Should(Succeed())
		newTask = t

		By("Ensuring due date is set")
		Expect(newTask.DueDate).Should(Equal(newDate))

		By("Removing due date")
		newTask.DueDate = time.Time{}

		By("Updating task")
		Eventually(func() error {
			t, err = client.UpdateTask(newTask)
			return err
		}).Should(Succeed())
		newTask = t

		By("Verifying due date is removed")
		Expect(newTask.DueDate).Should(Equal(time.Time{}))
	})

	It("can perform subtask CRUD", func() {
		By("Creating subtask")
		var subtask wl.Subtask
		subtaskComplete := false
		Eventually(func() error {
			subtask, err =
				client.CreateSubtask("mySubtaskTitle", newTask.ID, subtaskComplete)
			return err
		}).Should(Succeed())

		By("Getting subtask")
		Eventually(func() (wl.Subtask, error) {
			return client.Subtask(subtask.ID)
		}).Should(Equal(subtask))

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
		var s wl.Subtask
		Eventually(func() error {
			s, err = client.UpdateSubtask(subtask)
			return err
		}).Should(Succeed())
		subtask = s

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
			subtasksForList, err :=
				client.CompletedSubtasksForListID(newList.ID, showCompletedSubtasks)
			return subtaskContains(subtasksForList, subtask), err
		}).Should(BeTrue())

		By("Validating subtask exists in completed subtasks for task")
		Eventually(func() (bool, error) {
			subtasksForTask, err :=
				client.CompletedSubtasksForTaskID(newTask.ID, showCompletedSubtasks)
			return subtaskContains(subtasksForTask, subtask), err
		}).Should(BeTrue())

		By("Deleting subtask")
		Eventually(func() error {
			s, err := client.Subtask(subtask.ID)
			if err != nil {
				return err
			}
			return client.DeleteSubtask(s)
		}).Should(Succeed())
	})
})

func subtaskContains(subtasks []wl.Subtask, subtask wl.Subtask) bool {
	for _, t := range subtasks {
		if t.ID == subtask.ID {
			return true
		}
	}
	return false
}
