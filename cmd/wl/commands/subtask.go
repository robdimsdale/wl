package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robdimsdale/wl"
	"github.com/spf13/cobra"
)

var (
	// Commands
	cmdSubtasks = &cobra.Command{
		Use:   "subtasks",
		Short: "gets all subtasks",
		Long: `subtasks gets the user's subtasks.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			// Currently sending completed=false is the same as not sending completed
			// Checking for whether the flag has changed protects us from potentially
			// breaking changes, i.e. if the subtasks endpoint decides to return all tasks,
			// not just non-completed ones.

			if taskID != 0 {
				if cmd.Flags().Changed(completedLongFlag) {
					renderOutput(newClient(cmd).CompletedSubtasksForTaskID(taskID, completed))
				} else {
					renderOutput(newClient(cmd).SubtasksForTaskID(taskID))
				}
			} else if listID != 0 {
				if cmd.Flags().Changed(completedLongFlag) {
					renderOutput(newClient(cmd).CompletedSubtasksForListID(listID, completed))
				} else {
					renderOutput(newClient(cmd).SubtasksForListID(listID))
				}
			} else {
				if cmd.Flags().Changed(completedLongFlag) {
					renderOutput(newClient(cmd).CompletedSubtasks(completed))
				} else {
					renderOutput(newClient(cmd).Subtasks())
				}
			}
		},
	}

	cmdSubtask = &cobra.Command{
		Use:   "subtask <subtask-id>",
		Short: "gets the subtask for the provided subtask id",
		Long: `subtask gets a subtask specified by <subtask-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(subtask(cmd, args))
		},
	}

	cmdCreateSubtask = &cobra.Command{
		Use:   "create-subtask",
		Short: "creates a subtask with the specified args",
		Long: `create-subtask creates a subtask with the specified args
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).CreateSubtask(
				title,
				taskID,
				completed,
			))
		},
	}

	cmdUpdateSubtask = &cobra.Command{
		Use:   "update-subtask",
		Short: "updates a subtask with the specified args",
		Long: `update-subtask obtains the current state of the subtask,
and updates fields with the provided flags.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			subtask, err := subtask(cmd, args)
			if err != nil {
				handleError(err)
			}

			if cmd.Flags().Changed(titleLongFlag) {
				subtask.Title = title
			}

			if cmd.Flags().Changed(completedLongFlag) {
				subtask.Completed = completed
			}

			renderOutput(newClient(cmd).UpdateSubtask(subtask))
		},
	}

	cmdDeleteSubtask = &cobra.Command{
		Use:   "delete-subtask <subtask-id>",
		Short: "deletes the subtask for the provided subtask id",
		Long: `delete-subtask deletes the subtask specified by <subtask-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			subtask, err := subtask(cmd, args)
			if err != nil {
				fmt.Printf("error getting subtask: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = newClient(cmd).DeleteSubtask(subtask)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("subtask %d deleted successfully\n", subtask.ID)
		},
	}
)

func init() {
	cmdSubtasks.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
	cmdSubtasks.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "filter by taskID")
	cmdSubtasks.Flags().BoolVar(&completed, completedLongFlag, false, "filter for completed tasks")

	cmdCreateSubtask.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "id of task to which subtask belongs")
	cmdCreateSubtask.Flags().BoolVar(&completed, completedLongFlag, false, "whether subtask is completed")
	cmdCreateSubtask.Flags().StringVar(&title, titleLongFlag, "", "subtask title")

	cmdUpdateSubtask.Flags().BoolVar(&completed, completedLongFlag, false, "whether subtask is completed")
	cmdUpdateSubtask.Flags().StringVar(&title, titleLongFlag, "", "subtask title")
}

func subtask(cmd *cobra.Command, args []string) (wl.Subtask, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing subtaskID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).Subtask(id)
}
