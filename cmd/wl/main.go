package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/logger"
	"github.com/robdimsdale/wundergo/oauth"

	"github.com/spf13/cobra"
)

const (
	accessTokenEnvVariable = "WL_ACCESS_TOKEN"
	clientIDEnvVariable    = "WL_CLIENT_ID"

	accessTokenLongFlag = "accessToken"
	clientIDLongFlag    = "clientID"

	verboseLongFlag  = "verbose"
	verboseShortFlag = "v"

	useJSONLongFlag  = "useJSON"
	useJSONShortFlag = "j"

	listIDLongFlag  = "listID"
	listIDShortFlag = "l"

	taskIDLongFlag  = "taskID"
	taskIDShortFlag = "t"
)

var (
	// version is deliberately left uninitialized so it can be set at compile-time
	version string

	l logger.Logger

	// global flags
	accessToken string
	clientID    string
	verbose     bool
	useJSON     bool

	// non-global flags
	listID uint
	taskID uint

	cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "shows application version",
		Long: `version shows the version of the application.
        The version will be 'dev' if the application has been compiled
        without providing an explicit version.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	cmdInbox = &cobra.Command{
		Use:   "inbox",
		Short: "gets inbox",
		Long: `inbox gets the user's inbox.
        The inbox is a specific list, identified by the list_type field having value of 'inbox'.
        It cannot be deleted.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).Inbox())
		},
	}

	cmdRoot = &cobra.Command{
		Use:   "root",
		Short: "gets root",
		Long: `root gets the user's root.
        Root is the top of the list,task etc hierarchy'.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).Root())
		},
	}

	cmdFolders = &cobra.Command{
		Use:   "folders",
		Short: "gets all folders",
		Long: `folders gets the user's folders.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).Folders())
		},
	}

	cmdLists = &cobra.Command{
		Use:   "lists",
		Short: "gets all lists",
		Long: `lists gets the user's lists.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).Lists())
		},
	}

	cmdTasks = &cobra.Command{
		Use:   "tasks",
		Short: "gets all tasks",
		Long: `tasks gets the user's tasks.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if listID == 0 {
				renderOutput(newClient(cmd).Tasks())
			} else {
				renderOutput(newClient(cmd).TasksForListID(listID))
			}
		},
	}

	cmdDeleteAllFolders = &cobra.Command{
		Use:   "delete-all-folders",
		Short: "deletes all folders",
		Long: `delete-all-folders deletes all folders.
        Lists that are present in folders are not deleted.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			handleError(newClient(cmd).DeleteAllFolders())
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
			handleError(newClient(cmd).DeleteAllLists())
		},
	}

	cmdDeleteAllTasks = &cobra.Command{
		Use:   "delete-all-tasks",
		Short: "deletes all tasks",
		Long: `delete-all-tasks deletes all tasks.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			handleError(newClient(cmd).DeleteAllLists())
		},
	}

	cmdUploadFile = &cobra.Command{
		Use:   "upload-file <local-path> <remote-file-name> <content-type>",
		Short: "uploads a file",
		Long: `upload-file uploads the file at <local-path> to the remote name <remote-file-name>
        and giving it the content-type <content-type>.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				fmt.Printf("insufficient number of arguments provided\n\n")
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
				fmt.Printf("insufficient number of arguments provided\n\n")
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
				fmt.Printf("insufficient number of arguments provided\n\n")
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
)

func newClient(cmd *cobra.Command) wundergo.Client {
	if verbose {
		l = logger.NewLogger(logger.DEBUG)
	} else {
		l = logger.NewLogger(logger.INFO)
	}

	if accessToken == "" {
		accessToken = os.Getenv(accessTokenEnvVariable)
	}

	if accessToken == "" {
		l.Error(
			"exiting",
			errors.New("accessToken not found. Either provide the flag -"+accessTokenLongFlag+" or set the environment variable "+accessTokenEnvVariable))
		os.Exit(2)
	}

	if clientID == "" {
		clientID = os.Getenv(clientIDEnvVariable)
	}

	if clientID == "" {
		l.Error(
			"exiting",
			errors.New("clientID not found. Either provide the flag -"+clientIDLongFlag+" or set the environment variable "+clientIDEnvVariable))
		os.Exit(2)
	}

	return oauth.NewClient(accessToken, clientID, wundergo.APIURL, l)
}

func main() {
	if version == "" {
		version = "dev"
	}

	var rootCmd = &cobra.Command{Use: "wl"}

	rootCmd.PersistentFlags().BoolVarP(&verbose, verboseLongFlag, verboseShortFlag, false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&accessToken, accessTokenLongFlag, "", "", `Wunderlist access token. 
                      	Required, but can be provided via WL_ACCESS_TOKEN environment variable instead.`)
	rootCmd.PersistentFlags().StringVarP(&clientID, clientIDLongFlag, "", "", `Wunderlist client ID. 
                     Required, but can be provided via WL_CLIENT_ID environment variable instead.`)
	rootCmd.PersistentFlags().BoolVarP(&useJSON, useJSONLongFlag, useJSONShortFlag, false, "render output as JSON instead of YAML.")

	rootCmd.AddCommand(cmdVersion)
	rootCmd.AddCommand(cmdInbox)
	rootCmd.AddCommand(cmdRoot)
	rootCmd.AddCommand(cmdLists)
	rootCmd.AddCommand(cmdFolders)
	rootCmd.AddCommand(cmdTasks)
	rootCmd.AddCommand(cmdDeleteAllLists)
	rootCmd.AddCommand(cmdDeleteAllFolders)
	rootCmd.AddCommand(cmdDeleteAllTasks)
	rootCmd.AddCommand(cmdUploadFile)
	rootCmd.AddCommand(cmdCreateFile)
	rootCmd.AddCommand(cmdFile)
	rootCmd.AddCommand(cmdFiles)

	cmdTasks.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
	cmdFiles.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
	cmdFiles.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "filter by taskID")

	rootCmd.Execute()
}

func handleError(err error) {
	if err != nil {
		l.Error("exiting", err)
		os.Exit(1)
	}
}

func renderOutput(output interface{}, err error) {
	handleError(err)
	if useJSON {
		json.NewEncoder(os.Stdout).Encode(output)
	} else {
		data, err := yaml.Marshal(output)
		if err != nil {
			l.Error("exiting - failed to render yaml", err)
			os.Exit(1)
		}
		fmt.Printf(string(data))
	}

}
