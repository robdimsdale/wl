package wundergo_integration_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("basic upload and file functionality", func() {
	var (
		localFilePath  string
		remoteFileName string
		contentType    string
		md5sum         string
		tempDirPath    string

		firstList wundergo.List
		task      wundergo.Task
	)

	BeforeEach(func() {
		By("Creating temporary fixtures")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		remoteFileName = uuid1.String()

		tempDirPath, err = ioutil.TempDir(os.TempDir(), "wundergo-integration-test")
		Expect(err).NotTo(HaveOccurred())

		localFilePath = filepath.Join(tempDirPath, "test-file")

		fileContent := []byte("some-text")
		err = ioutil.WriteFile(localFilePath, fileContent, os.ModePerm)

		contentType = "text"
		md5sum = ""

		By("Creating a task")
		var lists []wundergo.List
		Eventually(func() error {
			lists, err = client.Lists()
			return err
		}).Should(Succeed())
		firstList = lists[0]

		uuid2, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle := uuid2.String()

		Eventually(func() error {
			task, err = client.CreateTask(
				newTaskTitle,
				firstList.ID,
				0,
				false,
				"",
				0,
				"1970-01-01",
				false,
			)
			return err
		}).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		var err error
		By("Deleting task")
		Eventually(func() error {
			task, err = client.Task(task.ID)
			return client.DeleteTask(task)
		}).Should(Succeed())

		var tasks []wundergo.Task
		Eventually(func() (bool, error) {
			tasks, err = client.TasksForListID(firstList.ID)
			return taskContains(tasks, task), err
		}).Should(BeFalse())

		err = os.RemoveAll(tempDirPath)
		Expect(err).ToNot(HaveOccurred())
	})

	It("can upload a local file", func() {
		By("Uploading a local file")
		upload, err := client.UploadFile(
			localFilePath,
			remoteFileName,
			contentType,
			md5sum,
		)

		Expect(err).NotTo(HaveOccurred())

		By("Creating a file to bind the upload to a task")
		file, err := client.CreateFile(upload.ID, task.ID)
		Expect(err).NotTo(HaveOccurred())

		By("Validating the file returns correctly")
		Eventually(func() (wundergo.File, error) {
			return client.File(file.ID)
		}).Should(Equal(file))

		By("Validating the file is correctly associated with the task")
		Expect(file.TaskID).To(Equal(task.ID))

		Eventually(func() (bool, error) {
			filesForTask, err := client.FilesForTaskID(task.ID)
			return fileContains(filesForTask, file), err
		}).Should(BeTrue())

		By("Validating the file is correctly associated with the list")
		Eventually(func() (bool, error) {
			filesForFirstList, err := client.FilesForListID(firstList.ID)
			return fileContains(filesForFirstList, file), err
		}).Should(BeTrue())

		By("Validating the file can be destroyed successfully")
		err = client.DestroyFile(file)
		Expect(err).NotTo(HaveOccurred())

		By("Validating the new file is not present in list of files")
		Eventually(func() (bool, error) {
			filesForTask, err := client.FilesForTaskID(task.ID)
			return fileContains(filesForTask, file), err
		}).Should(BeFalse())
	})
})

func fileContains(files []wundergo.File, file wundergo.File) bool {
	for _, f := range files {
		if f.ID == file.ID {
			return true
		}
	}
	return false
}
