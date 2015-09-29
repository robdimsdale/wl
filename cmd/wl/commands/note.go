package commands

import "github.com/spf13/cobra"

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
)

func init() {
	cmdNotes.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "filter by taskID")
	cmdNotes.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
}
