package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robdimsdale/wl"
	"github.com/spf13/cobra"
)

const (
	subtaskIDsLongFlag = "subtaskIDs"
)

var (
	// Flags
	subtaskIDs string

	// Commands
	cmdSubtaskPositions = &cobra.Command{
		Use:   "subtask-positions",
		Short: "gets all subtask positions",
		Long: `subtask-positions gets the positions of the user's subtasks.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if taskID != 0 {
				renderOutput(newClient(cmd).SubtaskPositionsForTaskID(taskID))
			} else if listID != 0 {
				renderOutput(newClient(cmd).SubtaskPositionsForListID(listID))
			} else {
				renderOutput(newClient(cmd).SubtaskPositions())
			}
		},
	}

	cmdSubtaskPosition = &cobra.Command{
		Use:   "subtask-position <subtask-position-id>",
		Short: "gets the subtask-position for the provided subtask-position id",
		Long: `subtask-position gets a subtask-position specified by <subtask-position-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(subtaskPosition(cmd, args))
		},
	}

	cmdUpdateSubtaskPosition = &cobra.Command{
		Use:   "update-subtask-position <subtask-position-id> [flags]",
		Short: "updates the subtask-position",
		Long: `update-subtask-position obtains the current state of the subtask-position specified by <subtask-position-id>,
and updates fields with the provided flags.
`,
		Run: func(cmd *cobra.Command, args []string) {
			subtaskPosition, err := subtaskPosition(cmd, args)
			if err != nil {
				fmt.Printf("error getting sub task position: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			if cmd.Flags().Changed(subtaskIDsLongFlag) {
				subtaskIDsInts, err := splitStringToUints(subtaskIDs)
				if err != nil {
					fmt.Printf("error parsing subtaskIDs: %v\n\n", err)
					cmd.Usage()
					os.Exit(2)
				}
				subtaskPosition.Values = subtaskIDsInts
			}

			renderOutput(newClient(cmd).UpdateSubtaskPosition(subtaskPosition))
		},
	}
)

func init() {
	cmdSubtaskPositions.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "filter by taskID")
	cmdSubtaskPositions.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")

	cmdUpdateSubtaskPosition.Flags().StringVar(&subtaskIDs, subtaskIDsLongFlag, "", "comma-separated subtask IDs (required)")
}

func subtaskPosition(cmd *cobra.Command, args []string) (wl.Position, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing subtask position ID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).SubtaskPosition(id)
}
