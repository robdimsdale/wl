package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robdimsdale/wl"
	"github.com/spf13/cobra"
)

const (
	textLongFlag = "text"
)

var (
	// Flags
	text string

	// Commands
	cmdTaskComments = &cobra.Command{
		Use:   "task-comments",
		Short: "gets all task-comments",
		Long: `task-comments gets the user's task-comments.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			// Currently sending completed=false is the same as not sending completed
			// Checking for whether the flag has changed protects us from potentially
			// breaking changes, i.e. if the taskComments endpoint decides to return all tasks,
			// not just non-completed ones.

			if taskID != 0 {
				renderOutput(newClient(cmd).TaskCommentsForTaskID(taskID))
			} else if listID != 0 {
				renderOutput(newClient(cmd).TaskCommentsForListID(listID))
			} else {
				renderOutput(newClient(cmd).TaskComments())
			}
		},
	}

	cmdTaskComment = &cobra.Command{
		Use:   "task-comment <task-comment-id>",
		Short: "gets the task-comment for the provided task-comment id",
		Long: `task-comment gets a task-comment specified by <task-comment-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(taskComment(cmd, args))
		},
	}

	cmdCreateTaskComment = &cobra.Command{
		Use:   "create-task-comment",
		Short: "creates a task-comment with the specified args",
		Long: `create-task-comment creates a task-comment with the specified args
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).CreateTaskComment(
				text,
				taskID,
			))
		},
	}

	cmdDeleteTaskComment = &cobra.Command{
		Use:   "delete-task-comment <task-comment-id>",
		Short: "deletes the task-comment for the provided task-comment id",
		Long: `delete-task-comment deletes the task-comment specified by <task-comment-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			taskComment, err := taskComment(cmd, args)
			if err != nil {
				fmt.Printf("error getting task-comment: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = newClient(cmd).DeleteTaskComment(taskComment)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("task-comment %d deleted successfully\n", taskComment.ID)
		},
	}
)

func init() {
	cmdTaskComments.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
	cmdTaskComments.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "filter by taskID")

	cmdCreateTaskComment.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "id of task to which task-comment belongs")
	cmdCreateTaskComment.Flags().StringVar(&text, textLongFlag, "", "task-comment text")
}

func taskComment(cmd *cobra.Command, args []string) (wl.TaskComment, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing task-comment ID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).TaskComment(id)
}
