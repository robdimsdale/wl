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
	cmdListPositions = &cobra.Command{
		Use:   "list-positions",
		Short: "gets all list positions",
		Long: `list-positions gets the positions of the user's lists.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).ListPositions())
		},
	}

	cmdListPosition = &cobra.Command{
		Use:   "list-position <list-position-id>",
		Short: "gets the list-position for the provided list-position id",
		Long: `list-position gets a list-position specified by <list-position-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(listPosition(cmd, args))
		},
	}

	cmdUpdateListPosition = &cobra.Command{
		Use:   "update-list-position <list-position-id> [flags]",
		Short: "updates the list-position",
		Long: `update-list-position obtains the current state of the list-position specified by <list-position-id>,
and updates fields with the provided flags.
`,
		Run: func(cmd *cobra.Command, args []string) {
			listPosition, err := listPosition(cmd, args)
			if err != nil {
				fmt.Printf("error getting list position: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			if cmd.Flags().Changed(listIDsLongFlag) {
				listIDsUints, err := splitStringToUints(listIDs)
				if err != nil {
					fmt.Printf("error parsing listIDs: %v\n\n", err)
					cmd.Usage()
					os.Exit(2)
				}
				listPosition.Values = listIDsUints
			}

			renderOutput(newClient(cmd).UpdateListPosition(listPosition))
		},
	}
)

func init() {
	cmdUpdateListPosition.Flags().StringVar(&listIDs, listIDsLongFlag, "", "comma-separated list IDs (required)")
}

func listPosition(cmd *cobra.Command, args []string) (wl.Position, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing list position ID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).ListPosition(id)
}
