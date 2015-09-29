package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/robdimsdale/wundergo"
	"github.com/spf13/cobra"
)

const (
	listIDsLongFlag = "listIDs"
)

var (
	// Flags
	listIDs string

	// Commands
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

	cmdCreateFolder = &cobra.Command{
		Use:   "create-folder",
		Short: "creates a folder with the specified args",
		Long: `create-folder creates a folder with the specified args
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if listIDs == "" {
				cmd.Usage()
				os.Exit(2)
			}

			splitListIDs := strings.Split(listIDs, ",")
			listIDsUints := make([]uint, len(splitListIDs))

			for i, s := range splitListIDs {
				idInt, err := strconv.Atoi(s)
				if err != nil {
					fmt.Printf("error parsing listIDs: %v at index %d\n\n", err, i)
					cmd.Usage()
					os.Exit(2)
				}
				listIDsUints[i] = uint(idInt)
			}

			renderOutput(newClient(cmd).CreateFolder(
				title,
				listIDsUints,
			))
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

func init() {
	cmdCreateFolder.Flags().StringVar(&title, titleLongFlag, "", "title of folder (required)")
	cmdCreateFolder.Flags().StringVar(&listIDs, listIDsLongFlag, "", "comma-separated list IDs (required)")
}

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
