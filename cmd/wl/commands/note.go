package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robdimsdale/wl"
	"github.com/spf13/cobra"
)

const (
	contentLongFlag = "content"
)

var (
	// Flags
	content string

	// Commands
	cmdNotes = &cobra.Command{
		Use:   "notes",
		Short: "gets all notes",
		Long: `notes gets the user's notes.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if taskID != 0 {
				renderOutput(newClient(cmd).NotesForTaskID(taskID))
			} else if listID != 0 {
				renderOutput(newClient(cmd).NotesForListID(listID))
			} else {
				renderOutput(newClient(cmd).Notes())
			}
		},
	}

	cmdNote = &cobra.Command{
		Use:   "note <note-id>",
		Short: "gets the note for the provided note id",
		Long: `note gets a note specified by <note-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(note(cmd, args))
		},
	}

	cmdCreateNote = &cobra.Command{
		Use:   "create-note",
		Short: "creates a note with the specified args",
		Long: `create-note creates a note with the specified args
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).CreateNote(
				content,
				taskID,
			))
		},
	}

	cmdUpdateNote = &cobra.Command{
		Use:   "update-note",
		Short: "updates a note with the specified args",
		Long: `update-note obtains the current state of the note,
and updates fields with the provided flags.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			note, err := note(cmd, args)
			if err != nil {
				handleError(err)
			}

			if cmd.Flags().Changed(contentLongFlag) {
				note.Content = content
			}

			renderOutput(newClient(cmd).UpdateNote(note))
		},
	}

	cmdDeleteNote = &cobra.Command{
		Use:   "delete-note <note-id>",
		Short: "deletes the note for the provided note id",
		Long: `delete-note deletes the note specified by <note-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			note, err := note(cmd, args)
			if err != nil {
				fmt.Printf("error getting note: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = newClient(cmd).DeleteNote(note)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("note %d deleted successfully\n", note.ID)
		},
	}
)

func init() {
	cmdNotes.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "filter by taskID")
	cmdNotes.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")

	cmdCreateNote.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "id of task to which note belongs")
	cmdCreateNote.Flags().StringVar(&content, contentLongFlag, "", "note content")

	cmdUpdateNote.Flags().StringVar(&content, contentLongFlag, "", "note content")
}

func note(cmd *cobra.Command, args []string) (wl.Note, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing noteID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).Note(id)
}
