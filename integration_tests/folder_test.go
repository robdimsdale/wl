package wl_integration_test

import (
	"fmt"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
)

var _ = Describe("basic folder functionality", func() {

	It("creates folders", func() {
		By("Creating a new list")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle := uuid1.String()

		var newList wl.List
		Eventually(func() error {
			newList, err = client.CreateList(newListTitle)
			return err
		}).Should(Succeed())

		By("Creating a new folder")
		uuid2, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newFolderTitle := uuid2.String()

		folderListIDs := []uint{newList.ID}
		var newFolder wl.Folder
		Eventually(func() error {
			newFolder, err = client.CreateFolder(newFolderTitle, folderListIDs)
			return err
		}).Should(Succeed())

		By("Verifying folder exists")
		var folders []wl.Folder
		Eventually(func() (bool, error) {
			folders, err = client.Folders()
			return folderContains(folders, newFolder), err
		}).Should(BeTrue())

		Eventually(func() (wl.Folder, error) {
			return client.Folder(newFolder.ID)
		}).Should(Equal(newFolder))

		By("Verifying folder revisions exist")
		var folderRevisions []wl.FolderRevision
		Eventually(func() (bool, error) {
			folderRevisions, err = client.FolderRevisions()
			return folderRevisionContainsFolder(folderRevisions, newFolder), err
		}).Should(BeTrue())

		By("Updating a folder")
		uuid3, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newFolderTitle2 := fmt.Sprintf("%s-updated", uuid3.String())

		newFolder.Title = newFolderTitle2
		var updatedFolder wl.Folder
		Eventually(func() error {
			updatedFolder, err = client.UpdateFolder(newFolder)
			return err
		}).Should(Succeed())

		Eventually(func() (wl.Folder, error) {
			return client.Folder(newFolder.ID)
		}).Should(Equal(updatedFolder))

		By("Deleting new folder")
		Eventually(func() error {
			n, err := client.Folder(newFolder.ID)
			if err != nil {
				return err
			}
			return client.DeleteFolder(n)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			folders, err := client.Folders()
			return folderContains(folders, newFolder), err
		}).Should(BeFalse())

		By("Deleting new list")
		Eventually(func() error {
			l, err := client.List(newList.ID)
			if err != nil {
				return err
			}
			return client.DeleteList(l)
		}).Should(Succeed())

		var lists []wl.List
		Eventually(func() (bool, error) {
			lists, err = client.Lists()
			return listContains(lists, newList), err
		}).Should(BeFalse())
	})
})

func folderContains(folders []wl.Folder, folder wl.Folder) bool {
	for _, f := range folders {
		if f.ID == folder.ID {
			return true
		}
	}
	return false
}

func folderRevisionContainsFolder(folderRevisions []wl.FolderRevision, folder wl.Folder) bool {
	for _, f := range folderRevisions {
		if f.ID == folder.ID {
			return true
		}
	}
	return false
}
