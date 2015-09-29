package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robdimsdale/wundergo"
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

	cmdFolder = &cobra.Command{
		Use:   "folder <folder-id>",
		Short: "gets the folder for the provided folder id",
		Long: `folder gets a folder specified by <folder-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(folder(cmd, args))
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

func folder(cmd *cobra.Command, args []string) (wundergo.Folder, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing folderID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).Folder(id)
}
