package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	cmdLists = &cobra.Command{
		Use:   "lists",
		Short: "gets all lists",
		Long: `lists gets the user's lists.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).Lists())
		},
	}

	cmdList = &cobra.Command{
		Use:   "list <list-id>",
		Short: "gets the list for the provided list id",
		Long: `list gets a list specified by <list-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Printf("incorrect number of arguments provided\n\n")
				cmd.Usage()
				os.Exit(2)
			}

			idInt, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("error parsing listID: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}
			id := uint(idInt)

			renderOutput(newClient(cmd).List(
				id,
			))
		},
	}

	cmdDeleteAllLists = &cobra.Command{
		Use:   "delete-all-lists",
		Short: "deletes all lists",
		Long: `delete-all-lists deletes all lists.
        This will not delete the inbox, as it cannot be deleted.
				This deletes all tasks that are not in the inbox,
        and all folders that the inbox is not a member of.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			err := newClient(cmd).DeleteAllLists()
			if err != nil {
				handleError(err)
			}
			fmt.Printf("all lists deleted successfully\n")
		},
	}
)
