package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	completedTasksLongFlag = "completed"
)

var (
	// Flags
	completedTasks bool

	// Commands
	cmdTasks = &cobra.Command{
		Use:   "tasks",
		Short: "gets all tasks",
		Long: `tasks gets the user's tasks.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			// Currently sending completed=false is the same as not sending completed
			// Checking for whether the flag has changed protects us from potentially
			// breaking changes, i.e. if the tasks endpoint decides to run all tasks,
			// not just non-completed ones.

			if listID == 0 {
				if cmd.Flags().Changed(completedTasksLongFlag) {
					renderOutput(newClient(cmd).CompletedTasks(completedTasks))
				} else {
					renderOutput(newClient(cmd).Tasks())
				}
			} else {
				if cmd.Flags().Changed(completedTasksLongFlag) {
					renderOutput(newClient(cmd).CompletedTasksForListID(listID, completedTasks))
				} else {
					renderOutput(newClient(cmd).TasksForListID(listID))
				}
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

	cmdDeleteTask = &cobra.Command{
		Use:   "delete-task <task-id>",
		Short: "deletes the task for the provided task id",
		Long: `delete-task deletes the task specified by <task-id>
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

			client := newClient(cmd)
			task, err := client.Task(id)
			if err != nil {
				fmt.Printf("error getting task: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = client.DeleteTask(task)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("task %d deleted successfully\n", id)
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
	cmdTasks.Flags().BoolVar(&completedTasks, completedTasksLongFlag, false, "filter for completed tasks")
}
