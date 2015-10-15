package wl_integration_test

import (
	"errors"
	"time"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
)

var _ = Describe("basic subtask position functionality", func() {
	var (
		newList wl.List
		newTask wl.Task
	)

	BeforeEach(func() {
		By("Creating a new list")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle1 := uuid1.String()

		Eventually(func() error {
			newList, err = client.CreateList(newListTitle1)
			return err
		}).Should(Succeed())

		By("Creating a new task")
		uuidTask, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle := uuidTask.String()

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
		}).Should(Succeed())
	})

	AfterEach(func() {
		By("Deleting task (and hence associated subtasks)")
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

		By("Deleting lists")
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

	It("reorders subtask positions", func() {
		By("Creating associated subtasks")
		uuid2, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newSubtaskTitle1 := uuid2.String()

		uuid3, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newSubtaskTitle2 := uuid3.String()

		var newSubtask1 wl.Subtask
		Eventually(func() error {
			newSubtask1, err = client.CreateSubtask(
				newSubtaskTitle1,
				newTask.ID,
				false,
			)
			return err
		}).Should(Succeed())

		var newSubtask2 wl.Subtask
		Eventually(func() error {
			newSubtask2, err = client.CreateSubtask(
				newSubtaskTitle2,
				newTask.ID,
				false,
			)
			return err
		}).Should(Succeed())

		// We have to reorder the subtasks before they are present in the
		// returned response. This seems like a bug in Wunderlist API

		By("Reordering subtasks")
		var firstListTasks []wl.Task
		Eventually(func() error {
			firstListTasks, err = client.TasksForListID(newList.ID)
			return err
		}).Should(Succeed())

		var index int
		for i, task := range firstListTasks {
			if task.ID == newTask.ID {
				index = i
			}
		}

		var subtaskPosition wl.Position

		Eventually(func() error {
			subtaskPositions, err := client.SubtaskPositionsForListID(newList.ID)
			if err != nil {
				return err
			}

			tp := subtaskPositions
			if len(tp) < index {
				return errors.New("subtasks not long enough to contain expected subtask")
			}
			subtaskPosition = tp[index]
			return nil
		}).Should(Succeed())

		subtaskPosition.Values = append(subtaskPosition.Values, newSubtask1.ID, newSubtask2.ID)

		Eventually(func() (bool, error) {
			sp, err := client.UpdateSubtaskPosition(subtaskPosition)
			if err != nil {
				return false, err
			}
			task1Contained := positionContainsValue(sp, newSubtask1.ID)
			task2Contained := positionContainsValue(sp, newSubtask2.ID)
			return task1Contained && task2Contained, err
		}).Should(BeTrue())

		Eventually(func() (bool, error) {
			firstListSubtaskPositions, err := client.SubtaskPositionsForListID(newList.ID)
			task1Contained := positionsContainValue(firstListSubtaskPositions, newSubtask1.ID)
			task2Contained := positionsContainValue(firstListSubtaskPositions, newSubtask2.ID)
			return task1Contained && task2Contained, err
		}).Should(BeTrue())
	})
})
