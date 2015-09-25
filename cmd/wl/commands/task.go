package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Commands
	cmdTasks = &cobra.Command{
		Use:   "tasks",
		Short: "gets all tasks",
		Long: `tasks gets the user's tasks.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if listID == 0 {
				renderOutput(newClient(cmd).Tasks())
			} else {
				renderOutput(newClient(cmd).TasksForListID(listID))
			}
		},
	}

	cmdDeleteAllTasks = &cobra.Command{
		Use:   "delete-all-tasks",
		Short: "deletes all tasks",
		Long: `delete-all-tasks deletes all tasks.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			err := newClient(cmd).DeleteAllTasks()
			if err != nil {
				handleError(err)
			}
			fmt.Printf("all tasks deleted successfully\n")
		},
	}
)

func init() {
	cmdTasks.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
}
