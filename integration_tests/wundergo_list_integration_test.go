package wundergo_integration_test

import (
	"fmt"
	"reflect"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("basic list functionality", func() {
	It("performs CRUD for lists", func() {

		By("Creating a new list")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle1 := fmt.Sprintf("%s-original", uuid1.String())

		newList, err := client.CreateList(newListTitle1)
		Expect(err).NotTo(HaveOccurred())

		By("Updating a list")
		uuid2, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle2 := fmt.Sprintf("%s-updated", uuid2.String())

		newList.Revision = newList.Revision + 1
		newList.Title = newListTitle2
		var updatedList *wundergo.List
		updatedList, err = client.UpdateList(*newList)
		Expect(err).NotTo(HaveOccurred())

		newList.Revision = newList.Revision + 1
		Eventually(func() (bool, error) {
			aList, err := client.List(newList.ID)
			return reflect.DeepEqual(updatedList, aList), err
		}).Should(BeTrue())

		var newLists *[]wundergo.List
		Eventually(func() (bool, error) {
			newLists, err = client.Lists()
			return listContains(newLists, newList), err
		}).Should(BeTrue())

		By("Deleting a list")
		newList, err = client.List(newList.ID)
		err = client.DeleteList(*newList)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() (bool, error) {
			lists, err := client.Lists()
			return listContains(lists, newList), err
		}).Should(BeFalse())
	})
})
