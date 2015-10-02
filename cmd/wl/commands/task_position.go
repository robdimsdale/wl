package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robdimsdale/wl"
	"github.com/spf13/cobra"
)

const (
	taskIDsLongFlag = "taskIDs"
)

var (
	// Flags
	taskIDs string

	// Commands
	cmdTaskPositions = &cobra.Command{
		Use:   "task-positions",
		Short: "gets all task positions",
		Long: `task-positions gets the positions of the user's tasks.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if listID == 0 {
				renderOutput(newClient(cmd).TaskPositions())
			} else {
				renderOutput(newClient(cmd).TaskPositionsForListID(listID))
			}
		},
	}

	cmdTaskPosition = &cobra.Command{
		Use:   "task-position <task-position-id>",
		Short: "gets the task-position for the provided task-position id",
		Long: `task-position gets a task-position specified by <task-position-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(taskPosition(cmd, args))
		},
	}

	cmdUpdateTaskPosition = &cobra.Command{
		Use:   "update-task-position <task-position-id> [flags]",
		Short: "updates the task-position",
		Long: `update-task-position obtains the current state of the task-position specified by <task-position-id>,
and updates fields with the provided flags.
`,
		Run: func(cmd *cobra.Command, args []string) {
			taskPosition, err := taskPosition(cmd, args)
			if err != nil {
				fmt.Printf("error getting task position: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			if cmd.Flags().Changed(taskIDsLongFlag) {
				taskIDsUints, err := splitStringToUints(taskIDs)
				if err != nil {
					fmt.Printf("error parsing taskIDs: %v\n\n", err)
					cmd.Usage()
					os.Exit(2)
				}
				taskPosition.Values = taskIDsUints
			}

			renderOutput(newClient(cmd).UpdateTaskPosition(taskPosition))
		},
	}
)

func init() {
	cmdTaskPositions.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")

	cmdUpdateTaskPosition.Flags().StringVar(&taskIDs, taskIDsLongFlag, "", "comma-separated task IDs (required)")
}

func taskPosition(cmd *cobra.Command, args []string) (wl.Position, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing task position ID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).TaskPosition(id)
}
