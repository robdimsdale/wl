package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cmdFolders = &cobra.Command{
		Use:   "folders",
		Short: "gets all folders",
		Long: `folders gets the user's folders.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).Folders())
		},
	}

	cmdDeleteAllFolders = &cobra.Command{
		Use:   "delete-all-folders",
		Short: "deletes all folders",
		Long: `delete-all-folders deletes all folders.
        Lists that are present in folders are not deleted.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			err := newClient(cmd).DeleteAllFolders()
			if err != nil {
				handleError(err)
			}
			fmt.Printf("all folders deleted successfully\n")
		},
	}
)
