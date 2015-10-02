package wl_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
)

var _ = Describe("basic task position functionality", func() {
	It("reorders task positions", func() {

		By("Creating a new list")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle := uuid1.String()

		newList, err := client.CreateList(newListTitle)
		Expect(err).NotTo(HaveOccurred())

		By("Creating new tasks")
		uuid2, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle1 := uuid2.String()

		uuid3, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle2 := uuid3.String()

		var newTask1 wl.Task
		Eventually(func() error {
			newTask1, err = client.CreateTask(
				newTaskTitle1,
				newList.ID,
				0,
				false,
				"",
				0,
				"",
				false,
			)
			return err
		}).Should(Succeed())

		var newTask2 wl.Task
		Eventually(func() error {
			newTask2, err = client.CreateTask(
				newTaskTitle2,
				newList.ID,
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

		By("Reordering tasks")
		var taskPosition wl.Position

		Eventually(func() error {
			taskPositions, err := client.TaskPositionsForListID(newList.ID)
			tp := taskPositions
			taskPosition = tp[0]
			return err
		}).Should(Succeed())

		taskPosition.Values = append(taskPosition.Values, newTask1.ID, newTask2.ID)

		Eventually(func() (bool, error) {
			taskPosition, err := client.UpdateTaskPosition(taskPosition)
			task1Contained := positionContainsValue(taskPosition, newTask1.ID)
			task2Contained := positionContainsValue(taskPosition, newTask2.ID)
			return task1Contained && task2Contained, err
		}).Should(BeTrue())

		By("Deleting tasks")
		Eventually(func() error {
			newTask1, err = client.Task(newTask1.ID)
			return client.DeleteTask(newTask1)
		}).Should(Succeed())

		Eventually(func() error {
			newTask2, err = client.Task(newTask2.ID)
			return client.DeleteTask(newTask2)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			tasks, err := client.TasksForListID(newList.ID)
			return taskContains(tasks, newTask1), err
		}).Should(BeFalse())

		Eventually(func() (bool, error) {
			tasks, err := client.TasksForListID(newList.ID)
			return taskContains(tasks, newTask2), err
		}).Should(BeFalse())

		By("Deleting new list")
		Eventually(func() error {
			newList, err = client.List(newList.ID)
			return client.DeleteList(newList)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			lists, err := client.Lists()
			return listContains(lists, newList), err
		}).Should(BeFalse())
	})
})
