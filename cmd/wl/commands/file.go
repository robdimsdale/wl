package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	filePreviewSizeLongFlag = "size"

	filePreviewPlatformLongFlag = "platform"
)

var (
	// Flags
	filePreviewPlatform string
	filePreviewSize     string

	// Commands
	cmdUploadFile = &cobra.Command{
		Use:   "upload-file <local-path> <remote-file-name> <content-type>",
		Short: "uploads a file",
		Long: `upload-file uploads the file at <local-path> to the remote name <remote-file-name>
        and giving it the content-type <content-type>.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				fmt.Printf("incorrect number of arguments provided\n\n")
				cmd.Usage()
				os.Exit(2)
			}

			localFilePath := args[0]
			remoteName := args[1]
			contentType := args[2]

			if localFilePath == "" || remoteName == "" || contentType == "" {
				fmt.Printf("invalid arguments provided\n\n")
				cmd.Usage()
				os.Exit(2)
			}

			renderOutput(newClient(cmd).UploadFile(
				localFilePath,
				remoteName,
				contentType,
				"",
			))
		},
	}

	cmdCreateFile = &cobra.Command{
		Use:   "create-file <upload-id> <task-id>",
		Short: "creates a file from the specified upload in the specified task",
		Long: `create-file creates a file from the upload specified by <upload-id>
        in the task specified by <task-id>.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Printf("incorrect number of arguments provided\n\n")
				cmd.Usage()
				os.Exit(2)
			}

			uploadIDInt, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("error parsing uploadID: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}
			uploadID := uint(uploadIDInt)

			taskIDInt, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Printf("error parsing taskID: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}
			taskID := uint(taskIDInt)

			renderOutput(newClient(cmd).CreateFile(
				uploadID,
				taskID,
			))
		},
	}

	cmdFile = &cobra.Command{
		Use:   "file <file-id>",
		Short: "gets the file for the provided file id",
		Long: `file gets a file specified by <file-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Printf("incorrect number of arguments provided\n\n")
				cmd.Usage()
				os.Exit(2)
			}

			fileIDInt, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("error parsing fileID: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}
			fileID := uint(fileIDInt)

			renderOutput(newClient(cmd).File(
				fileID,
			))
		},
	}

	cmdFiles = &cobra.Command{
		Use:   "files",
		Short: "gets all files",
		Long: `files gets the user's files.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if taskID != 0 {
				renderOutput(newClient(cmd).FilesForTaskID(taskID))
			} else if listID != 0 {
				renderOutput(newClient(cmd).FilesForListID(listID))
			} else {
				renderOutput(newClient(cmd).Files())
			}
		},
	}

	cmdDestroyFile = &cobra.Command{
		Use:   "destroy-file <file-id>",
		Short: "destroys (deletes) the file for the provided file id",
		Long: `destroy-file destroys (deletes) the file specified by <file-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Printf("incorrect number of arguments provided\n\n")
				cmd.Usage()
				os.Exit(2)
			}

			fileIDInt, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("error parsing fileID: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}
			fileID := uint(fileIDInt)

			client := newClient(cmd)
			file, err := client.File(fileID)
			if err != nil {
				fmt.Printf("error getting file: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = client.DestroyFile(file)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("file %d destroyed successfully\n", fileID)
		},
	}

	cmdFilePreview = &cobra.Command{
		Use:   "file-preview <file-id>",
		Short: "gets a preview of the file for the provided file id",
		Long: `file preview gets a preview of the file specified by <file-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Printf("incorrect number of arguments provided\n\n")
				cmd.Usage()
				os.Exit(2)
			}

			fileIDInt, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("error parsing fileID: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}
			fileID := uint(fileIDInt)

			renderOutput(newClient(cmd).FilePreview(
				fileID,
				filePreviewPlatform,
				filePreviewSize,
			))
		},
	}
)

func init() {
	cmdFiles.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
	cmdFiles.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "filter by taskID")
	cmdFilePreview.Flags().StringVar(&filePreviewSize, filePreviewSizeLongFlag, "", "obtain preview for specific size")
	cmdFilePreview.Flags().StringVar(&filePreviewPlatform, filePreviewPlatformLongFlag, "", "obtain preview for specific platform")
}
