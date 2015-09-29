package commands

import "github.com/spf13/cobra"

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
)

func init() {
	cmdSubtasks.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
	cmdSubtasks.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "filter by taskID")
	cmdSubtasks.Flags().BoolVar(&completed, completedLongFlag, false, "filter for completed tasks")
}
