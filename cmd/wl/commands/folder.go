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
			listIDsUints, err := splitStringToUints(listIDs)
			if err != nil {
				fmt.Printf("error parsing listIDs: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			renderOutput(newClient(cmd).CreateFolder(
				title,
				listIDsUints,
			))
		},
	}

	cmdUpdateFolder = &cobra.Command{
		Use:   "update-folder",
		Short: "updates a folder with the specified args",
		Long: `update-folder obtains the current state of the folder,
and updates fields with the provided flags.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			folder, err := folder(cmd, args)
			if err != nil {
				handleError(err)
			}

			if cmd.Flags().Changed(titleLongFlag) {
				folder.Title = title
			}

			if cmd.Flags().Changed(listIDsLongFlag) {
				listIDsUints, err := splitStringToUints(listIDs)
				if err != nil {
					fmt.Printf("error parsing listIDs: %v\n\n", err)
					cmd.Usage()
					os.Exit(2)
				}
				folder.ListIDs = listIDsUints
			}

			renderOutput(newClient(cmd).UpdateFolder(folder))
		},
	}

	cmdDeleteFolder = &cobra.Command{
		Use:   "delete-folder <folder-id>",
		Short: "deletes the folder for the provided folder id",
		Long: `delete-folder deletes the folder specified by <folder-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			folder, err := folder(cmd, args)
			if err != nil {
				fmt.Printf("error getting folder: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = newClient(cmd).DeleteFolder(folder)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("folder %d deleted successfully\n", folder.ID)
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

	cmdUpdateFolder.Flags().StringVar(&title, titleLongFlag, "", "title of folder (required)")
	cmdUpdateFolder.Flags().StringVar(&listIDs, listIDsLongFlag, "", "comma-separated list IDs (required)")
}

func folder(cmd *cobra.Command, args []string) (wl.Folder, error) {
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
