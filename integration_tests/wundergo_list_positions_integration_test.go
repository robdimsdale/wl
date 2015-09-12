package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("basic list position functionality", func() {
	It("reorders list positions", func() {

		By("Creating new lists")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle1 := uuid1.String()

		uuid2, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle2 := uuid2.String()

		var newList1 wundergo.List
		Eventually(func() error {
			newList1, err = client.CreateList(newListTitle1)
			return err
		}).Should(Succeed())

		var newList2 wundergo.List
		Eventually(func() error {
			newList2, err = client.CreateList(newListTitle2)
			return err
		}).Should(Succeed())

		// We have to reorder the lists before they are present in the
		// returned response. This seems like a bug in Wunderlist API

		By("Reordering lists")
		var listPosition wundergo.Position

		Eventually(func() error {
			listPositions, err := client.ListPositions()
			lp := listPositions
			listPosition = lp[0]
			return err
		}).Should(Succeed())

		listPosition.Values = append(listPosition.Values, newList1.ID, newList2.ID)

		Eventually(func() (bool, error) {
			listPosition, err := client.UpdateListPosition(listPosition)
			list1Contained := positionContainsValue(listPosition, newList1.ID)
			list2Contained := positionContainsValue(listPosition, newList2.ID)
			return list1Contained && list2Contained, err
		}).Should(BeTrue())

		By("Deleting lists")
		Eventually(func() error {
			newList1, err = client.List(newList1.ID)
			return client.DeleteList(newList1)
		}).Should(Succeed())

		Eventually(func() error {
			newList2, err = client.List(newList2.ID)
			return client.DeleteList(newList2)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			lists, err := client.Lists()
			return listContains(lists, newList1), err
		}).Should(BeFalse())

		Eventually(func() (bool, error) {
			lists, err := client.Lists()
			return listContains(lists, newList2), err
		}).Should(BeFalse())
	})
})
