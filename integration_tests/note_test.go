package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("basic note functionality", func() {
	var (
		newList wundergo.List
		newTask wundergo.Task
		err     error
	)

	BeforeEach(func() {
		By("Creating a new list")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle := uuid1.String()

		newList, err = client.CreateList(newListTitle)
		Expect(err).NotTo(HaveOccurred())

		By("Creating task in new list")
		uuid, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle := uuid.String()

		Eventually(func() error {
			newTask, err = client.CreateTask(
				newTaskTitle,
				newList.ID,
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
		By("Deleting task")
		Eventually(func() error {
			newTask, err = client.Task(newTask.ID)
			return client.DeleteTask(newTask)
		}).Should(Succeed())

		var tasks []wundergo.Task
		Eventually(func() (bool, error) {
			tasks, err = client.TasksForListID(newList.ID)
			return taskContains(tasks, newTask), err
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

	It("can perform note CRUD", func() {
		var note wundergo.Note
		Eventually(func() error {
			note, err = client.CreateNote("myNoteContent", newTask.ID)
			return err
		}).Should(Succeed())

		By("Verifying note appears in all notes")
		Eventually(func() bool {
			// It is statistically probable that one of the lists will
			// be deleted, so we ignore error here.
			allNotes, _ := client.Notes()
			return noteContains(allNotes, note)
		}).Should(BeTrue())

		By("Verifying note appears in notes for task")
		Eventually(func() (bool, error) {
			taskNotes, err := client.NotesForTaskID(newTask.ID)
			return noteContains(taskNotes, note), err
		}).Should(BeTrue())

		By("Verifying note appears in notes for list")
		Eventually(func() (bool, error) {
			listNotes, err := client.NotesForListID(newList.ID)
			return noteContains(listNotes, note), err
		}).Should(BeTrue())

		By("Updating note")
		note.Content = "newNoteContent"
		Eventually(func() error {
			note, err = client.UpdateNote(note)
			return err
		}).Should(Succeed())

		By("Getting note")
		newNote, err := client.Note(note.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(newNote.Content).To(Equal("newNoteContent"))

		By("Deleting note")
		Eventually(func() error {
			return client.DeleteNote(note)
		}).Should(Succeed())
	})
})

func noteContains(notes []wundergo.Note, note wundergo.Note) bool {
	for _, n := range notes {
		if n.ID == note.ID {
			return true
		}
	}
	return false
}
