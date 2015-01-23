package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("Basic subtask position functionality", func() {
	It("reorders subtask positions", func() {

		// Create task and subtasks

		uuidTask, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle := uuidTask.String()

		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newSubtaskTitle1 := uuid1.String()

		uuid2, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newSubtaskTitle2 := uuid2.String()

		var firstList *wundergo.List
		Eventually(func() error {
			l, err := client.Lists()
			lists := *l
			firstList = &lists[0]
			return err
		}).Should(Succeed())

		var newTask *wundergo.Task
		Eventually(func() error {
			newTask, err = client.CreateTask(
				newTaskTitle,
				firstList.ID,
				0,
				false,
				"",
				0,
				"",
				false,
			)
			return err
		}).Should(Succeed())

		var newSubtask1 *wundergo.Subtask
		Eventually(func() error {
			newSubtask1, err = client.CreateSubtask(
				newSubtaskTitle1,
				newTask.ID,
				false,
			)
			return err
		}).Should(Succeed())

		var newSubtask2 *wundergo.Subtask
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

		var firstListTasks []wundergo.Task
		Eventually(func() error {
			flt, err := client.TasksForListID(firstList.ID)
			firstListTasks = *flt
			return err
		}).Should(Succeed())

		var index int
		for i, task := range firstListTasks {
			if task.ID == newTask.ID {
				index = i
			}
		}

		var subtaskPosition *wundergo.Position

		Eventually(func() error {
			subtaskPositions, err := client.SubtaskPositionsForListID(firstList.ID)
			tp := *subtaskPositions
			subtaskPosition = &tp[index]
			return err
		}).Should(Succeed())

		subtaskPosition.Values = append(subtaskPosition.Values, newSubtask1.ID, newSubtask2.ID)

		Eventually(func() (bool, error) {
			subtaskPosition, err := client.UpdateSubtaskPosition(*subtaskPosition)
			task1Contained := positionContainsValue(subtaskPosition, newSubtask1.ID)
			task2Contained := positionContainsValue(subtaskPosition, newSubtask2.ID)
			return task1Contained && task2Contained, err
		}).Should(BeTrue())

		Eventually(func() (bool, error) {
			firstListSubtaskPositions, err := client.SubtaskPositionsForListID(firstList.ID)
			task1Contained := positionsContainValue(firstListSubtaskPositions, newSubtask1.ID)
			task2Contained := positionsContainValue(firstListSubtaskPositions, newSubtask2.ID)
			return task1Contained && task2Contained, err
		}).Should(BeTrue())

		// Delete task

		Eventually(func() error {
			newTask, err = client.Task(newTask.ID)
			return client.DeleteTask(*newTask)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			tasks, err := client.TasksForListID(firstList.ID)
			return taskContains(tasks, newTask), err
		}).Should(BeFalse())
	})
})
