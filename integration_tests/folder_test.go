package wundergo_integration_test

import (
	"fmt"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("basic webhook functionality", func() {

	It("creates folders", func() {
		By("Creating a new list")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle := uuid1.String()

		newList, err := client.CreateList(newListTitle)
		Expect(err).NotTo(HaveOccurred())

		By("Creating a new folder")
		uuid2, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newFolderTitle := uuid2.String()

		folderListIDs := []uint{newList.ID}
		newFolder, err := client.CreateFolder(newFolderTitle, folderListIDs)
		Expect(err).NotTo(HaveOccurred())

		By("Verifying folder exists")
		var folders []wundergo.Folder
		Eventually(func() (bool, error) {
			folders, err = client.Folders()
			return folderContains(folders, newFolder), err
		}).Should(BeTrue())

		Eventually(func() (wundergo.Folder, error) {
			return client.Folder(newFolder.ID)
		}).Should(Equal(newFolder))

		By("Updating a folder")
		uuid3, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newFolderTitle2 := fmt.Sprintf("%s-updated", uuid3.String())

		newFolder.Title = newFolderTitle2
		var updatedFolder wundergo.Folder
		updatedFolder, err = client.UpdateFolder(newFolder)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() (wundergo.Folder, error) {
			return client.Folder(newFolder.ID)
		}).Should(Equal(updatedFolder))

		By("Deleting new folder")
		Eventually(func() error {
			newFolder, err = client.Folder(newFolder.ID)
			return client.DeleteFolder(newFolder)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			folders, err := client.Folders()
			return folderContains(folders, newFolder), err
		}).Should(BeFalse())

		By("Deleting new list")
		Eventually(func() error {
			newList, err = client.List(newList.ID)
			return client.DeleteList(newList)
		}).Should(Succeed())

		var lists []wundergo.List
		Eventually(func() (bool, error) {
			lists, err = client.Lists()
			return listContains(lists, newList), err
		}).Should(BeFalse())
	})
})

func folderContains(folders []wundergo.Folder, folder wundergo.Folder) bool {
	for _, f := range folders {
		if f.ID == folder.ID {
			return true
		}
	}
	return false
}
