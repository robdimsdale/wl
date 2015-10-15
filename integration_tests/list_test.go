package wl_integration_test

import (
	"fmt"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
)

var _ = Describe("basic list functionality", func() {
	It("performs CRUD for lists", func() {
		By("Creating a new list")
		var newList wl.List

		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle1 := uuid1.String()

		Eventually(func() error {
			newList, err = client.CreateList(newListTitle1)
			return err
		}).Should(Succeed())

		By("Verifying list exists in lists")
		var newLists []wl.List
		Eventually(func() (bool, error) {
			newLists, err = client.Lists()
			return listContains(newLists, newList), err
		}).Should(BeTrue())

		By("Updating a list")
		uuid2, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle2 := fmt.Sprintf("%s-updated", uuid2.String())

		newList.Title = newListTitle2
		var updatedList wl.List
		Eventually(func() error {
			updatedList, err = client.UpdateList(newList)
			return err
		}).Should(Succeed())

		newList.Revision = newList.Revision + 1
		Eventually(func() (wl.List, error) {
			return client.List(newList.ID)
		}).Should(Equal(updatedList))

		By("Deleting a list")
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

	It("retrieves inbox", func() {
		var err error

		var inboxList wl.List
		Eventually(func() error {
			inboxList, err = client.Inbox()
			return err
		}).Should(Succeed())

		Expect(inboxList.Title).To(Equal("inbox"))
	})
})
