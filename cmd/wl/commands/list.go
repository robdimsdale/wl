package commands

import (
	"fmt"

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
