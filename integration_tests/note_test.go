package wl_integration_test

import (
	"time"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
)

var _ = Describe("basic note functionality", func() {
	var (
		newList wl.List
		newTask wl.Task
		err     error
	)

	BeforeEach(func() {
		By("Creating a new list")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle := uuid1.String()

		Eventually(func() error {
			newList, err = client.CreateList(newListTitle)
			return err
		}).Should(Succeed())

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
				time.Date(1971, 12, 31, 0, 0, 0, 0, time.UTC),
				false,
			)
			return err
		}).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Deleting task")
		Eventually(func() error {
			t, err := client.Task(newTask.ID)
			if err != nil {
				return err
			}
			return client.DeleteTask(t)
		}).Should(Succeed())

		var tasks []wl.Task
		Eventually(func() (bool, error) {
			tasks, err = client.TasksForListID(newList.ID)
			return taskContains(tasks, newTask), err
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

	It("can perform note CRUD", func() {
		var note wl.Note
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
		var n wl.Note
		Eventually(func() error {
			n, err = client.UpdateNote(note)
			return err
		}).Should(Succeed())
		note = n

		By("Getting note")
		var newNote wl.Note
		Eventually(func() error {
			newNote, err = client.Note(note.ID)
			return err
		}).Should(Succeed())
		Expect(newNote.Content).To(Equal("newNoteContent"))

		By("Deleting note")
		Eventually(func() error {
			n, err := client.Note(note.ID)
			if err != nil {
				return err
			}
			return client.DeleteNote(n)
		}).Should(Succeed())
	})
})

func noteContains(notes []wl.Note, note wl.Note) bool {
	for _, n := range notes {
		if n.ID == note.ID {
			return true
		}
	}
	return false
}
