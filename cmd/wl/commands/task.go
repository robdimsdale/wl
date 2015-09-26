package commands

import (
	"fmt"
	"os"
	"strconv"

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

	cmdTask = &cobra.Command{
		Use:   "task <task-id>",
		Short: "gets the task for the provided task id",
		Long: `task gets a task specified by <task-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Printf("incorrect number of arguments provided\n\n")
				cmd.Usage()
				os.Exit(2)
			}

			idInt, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("error parsing taskID: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}
			id := uint(idInt)

			renderOutput(newClient(cmd).Task(
				id,
			))
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
