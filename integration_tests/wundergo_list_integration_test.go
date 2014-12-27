package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("Basic list functionality", func() {
	It("performs CRUD for lists", func() {
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle1 := uuid1.String()

		uuid2, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle2 := uuid2.String()

		var originalLists *[]wundergo.List
		Eventually(func() error {
			originalLists, err = client.Lists()
			return err
		}).Should(Succeed())

		var newList *wundergo.List
		Eventually(func() error {
			newList, err = client.CreateList(newListTitle1)
			return err
		}).Should(Succeed())

		var newLists *[]wundergo.List
		Eventually(func() error {
			newLists, err = client.Lists()
			return err
		}).Should(Succeed())
		Expect(listContains(newLists, newList)).To(BeTrue())

		newList.Title = newListTitle2
		var updatedList *wundergo.List
		Eventually(func() error {
			updatedList, err = client.UpdateList(*newList)
			return err
		}).Should(Succeed())
		newList.Revision = newList.Revision + 1
		Expect(updatedList).To(Equal(newList))

		Eventually(func() error {
			newList, err = client.List(newList.ID)
			return client.DeleteList(*newList)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			lists, err := client.Lists()
			return listContains(lists, newList), err
		}).Should(BeFalse())
	})
})
