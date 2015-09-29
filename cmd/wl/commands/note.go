package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robdimsdale/wundergo"
	"github.com/spf13/cobra"
)

var (
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
)

func init() {
	cmdNotes.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "filter by taskID")
	cmdNotes.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
}

func note(cmd *cobra.Command, args []string) (wundergo.Note, error) {
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
