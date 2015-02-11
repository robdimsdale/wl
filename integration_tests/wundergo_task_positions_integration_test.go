package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("basic task position functionality", func() {
	It("reorders task positions", func() {

		// Create lists

		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle1 := uuid1.String()

		uuid2, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle2 := uuid2.String()

		var firstList *wundergo.List
		Eventually(func() error {
			l, err := client.Lists()
			lists := *l
			firstList = &lists[0]
			return err
		}).Should(Succeed())

		var newTask1 *wundergo.Task
		Eventually(func() error {
			newTask1, err = client.CreateTask(
				newTaskTitle1,
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

		var newTask2 *wundergo.Task
		Eventually(func() error {
			newTask2, err = client.CreateTask(
				newTaskTitle2,
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

		// We have to reorder the tasks before they are present in the
		// returned response. This seems like a bug in Wunderlist API

		// Assume tasks are in first TaskPosition

		var taskPosition *wundergo.Position

		Eventually(func() error {
			taskPositions, err := client.TaskPositionsForListID(firstList.ID)
			tp := *taskPositions
			taskPosition = &tp[0]
			return err
		}).Should(Succeed())

		taskPosition.Values = append(taskPosition.Values, newTask1.ID, newTask2.ID)

		Eventually(func() (bool, error) {
			taskPosition, err := client.UpdateTaskPosition(*taskPosition)
			task1Contained := positionContainsValue(taskPosition, newTask1.ID)
			task2Contained := positionContainsValue(taskPosition, newTask2.ID)
			return task1Contained && task2Contained, err
		}).Should(BeTrue())

		// Delete tasks

		Eventually(func() error {
			newTask1, err = client.Task(newTask1.ID)
			return client.DeleteTask(*newTask1)
		}).Should(Succeed())

		Eventually(func() error {
			newTask2, err = client.Task(newTask2.ID)
			return client.DeleteTask(*newTask2)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			tasks, err := client.TasksForListID(firstList.ID)
			return taskContains(tasks, newTask1), err
		}).Should(BeFalse())

		Eventually(func() (bool, error) {
			tasks, err := client.TasksForListID(firstList.ID)
			return taskContains(tasks, newTask2), err
		}).Should(BeFalse())
	})
})
