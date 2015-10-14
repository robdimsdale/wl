package wl_integration_test

import (
	"strings"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
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

		var newList1 wl.List
		Eventually(func() error {
			newList1, err = client.CreateList(newListTitle1)
			return err
		}).Should(Succeed())

		var newList2 wl.List
		Eventually(func() error {
			newList2, err = client.CreateList(newListTitle2)
			return err
		}).Should(Succeed())

		// We have to reorder the lists before they are present in the
		// returned response. This seems like a bug in Wunderlist API

		By("Reordering lists")
		var listPosition wl.Position

		for {
			Eventually(func() error {
				listPositions, err := client.ListPositions()
				listPosition = listPositions[0]
				return err
			}).Should(Succeed())

			listPosition.Values = []uint{newList1.ID, newList2.ID}

			listPosition, err = client.UpdateListPosition(listPosition)
			if err != nil {
				if strings.Contains(err.Error(), "409") {
					err = nil
					continue
				}
				break // Unexpected error
			}
			break // No error
		}
		Expect(err).NotTo(HaveOccurred())

		list1Contained := positionContainsValue(listPosition, newList1.ID)
		list2Contained := positionContainsValue(listPosition, newList2.ID)
		Expect(list1Contained).To(BeTrue())
		Expect(list2Contained).To(BeTrue())

		By("Deleting lists")
		Eventually(func() error {
			l, err := client.List(newList1.ID)
			if err != nil {
				return err
			}
			return client.DeleteList(l)
		}).Should(Succeed())

		Eventually(func() error {
			l, err := client.List(newList2.ID)
			if err != nil {
				return err
			}
			return client.DeleteList(l)
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
