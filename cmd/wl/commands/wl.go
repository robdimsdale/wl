package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/robdimsdale/wl"
	"github.com/robdimsdale/wl/logger"
	"github.com/robdimsdale/wl/oauth"
	"github.com/spf13/cobra"
)

const (
	// Global flags
	accessTokenEnvVariable = "WL_ACCESS_TOKEN"
	clientIDEnvVariable    = "WL_CLIENT_ID"

	accessTokenLongFlag = "accessToken"
	clientIDLongFlag    = "clientID"

	verboseLongFlag  = "verbose"
	verboseShortFlag = "v"

	useJSONLongFlag  = "useJSON"
	useJSONShortFlag = "j"

	// Shared, non-global flags
	listIDLongFlag  = "listID"
	listIDShortFlag = "l"

	taskIDLongFlag  = "taskID"
	taskIDShortFlag = "t"

	titleLongFlag = "title"

	completedLongFlag = "completed"

	listIDsLongFlag = "listIDs"
)

var (
	// Global flags
	accessToken string
	clientID    string
	verbose     bool
	useJSON     bool

	// Non-global, shared flags
	taskID    uint
	listID    uint
	title     string
	completed bool
	listIDs   string

	// WLCmd is the root command. All other commands are subcommands of it.
	WLCmd = &cobra.Command{Use: "wl"}
)

// Execute adds all child commands to the root command WLCmd,
// and executes the root command.
func Execute() {
	addCommands()
	WLCmd.Execute()
}

// Sets global flags
func init() {
	WLCmd.PersistentFlags().BoolVarP(&verbose, verboseLongFlag, verboseShortFlag, false, "verbose output")
	WLCmd.PersistentFlags().StringVarP(&accessToken, accessTokenLongFlag, "", "", `Wunderlist access token. 
                      	Required, but can be provided via WL_ACCESS_TOKEN environment variable instead.`)
	WLCmd.PersistentFlags().StringVarP(&clientID, clientIDLongFlag, "", "", `Wunderlist client ID. 
                     Required, but can be provided via WL_CLIENT_ID environment variable instead.`)
	WLCmd.PersistentFlags().BoolVarP(&useJSON, useJSONLongFlag, useJSONShortFlag, false, "render output as JSON instead of YAML.")
}

func addCommands() {
	WLCmd.AddCommand(cmdInbox)
	WLCmd.AddCommand(cmdRoot)
	WLCmd.AddCommand(cmdLists)
	WLCmd.AddCommand(cmdList)
	WLCmd.AddCommand(cmdCreateList)
	WLCmd.AddCommand(cmdUpdateList)
	WLCmd.AddCommand(cmdDeleteList)
	WLCmd.AddCommand(cmdDeleteAllLists)

	WLCmd.AddCommand(cmdFolders)
	WLCmd.AddCommand(cmdFolder)
	WLCmd.AddCommand(cmdCreateFolder)
	WLCmd.AddCommand(cmdUpdateFolder)
	WLCmd.AddCommand(cmdDeleteFolder)
	WLCmd.AddCommand(cmdDeleteAllFolders)

	WLCmd.AddCommand(cmdTasks)
	WLCmd.AddCommand(cmdTask)
	WLCmd.AddCommand(cmdCreateTask)
	WLCmd.AddCommand(cmdUpdateTask)
	WLCmd.AddCommand(cmdDeleteTask)
	WLCmd.AddCommand(cmdDeleteAllTasks)

	WLCmd.AddCommand(cmdUploadFile)
	WLCmd.AddCommand(cmdCreateFile)
	WLCmd.AddCommand(cmdFile)
	WLCmd.AddCommand(cmdFiles)
	WLCmd.AddCommand(cmdDestroyFile)
	WLCmd.AddCommand(cmdFilePreview)

	WLCmd.AddCommand(cmdUsers)
	WLCmd.AddCommand(cmdUser)
	WLCmd.AddCommand(cmdUpdateUser)
	WLCmd.AddCommand(cmdAvatarURL)

	WLCmd.AddCommand(cmdNotes)
	WLCmd.AddCommand(cmdNote)
	WLCmd.AddCommand(cmdCreateNote)
	WLCmd.AddCommand(cmdUpdateNote)
	WLCmd.AddCommand(cmdDeleteNote)

	WLCmd.AddCommand(cmdSubtasks)
	WLCmd.AddCommand(cmdCreateSubtask)
	WLCmd.AddCommand(cmdSubtask)
	WLCmd.AddCommand(cmdUpdateSubtask)
	WLCmd.AddCommand(cmdDeleteSubtask)

	WLCmd.AddCommand(cmdWebhooks)
	WLCmd.AddCommand(cmdWebhook)
	WLCmd.AddCommand(cmdCreateWebhook)
	WLCmd.AddCommand(cmdDeleteWebhook)

	WLCmd.AddCommand(cmdReminders)
	WLCmd.AddCommand(cmdCreateReminder)
	WLCmd.AddCommand(cmdReminder)
	WLCmd.AddCommand(cmdUpdateReminder)
	WLCmd.AddCommand(cmdDeleteReminder)

	WLCmd.AddCommand(cmdTaskComments)
	WLCmd.AddCommand(cmdCreateTaskComment)
	WLCmd.AddCommand(cmdTaskComment)
	WLCmd.AddCommand(cmdDeleteTaskComment)

	WLCmd.AddCommand(cmdMemberships)
	WLCmd.AddCommand(cmdMembership)
	WLCmd.AddCommand(cmdRejectMembership)
	WLCmd.AddCommand(cmdAcceptMembership)
	WLCmd.AddCommand(cmdRemoveMembership)
	WLCmd.AddCommand(cmdInviteMember)

	WLCmd.AddCommand(cmdListPositions)
	WLCmd.AddCommand(cmdListPosition)
	WLCmd.AddCommand(cmdUpdateListPosition)

	WLCmd.AddCommand(cmdTaskPositions)
	WLCmd.AddCommand(cmdTaskPosition)
	WLCmd.AddCommand(cmdUpdateTaskPosition)

	WLCmd.AddCommand(cmdSubtaskPositions)
	WLCmd.AddCommand(cmdSubtaskPosition)
	WLCmd.AddCommand(cmdUpdateSubtaskPosition)
}

func newClient(cmd *cobra.Command) wl.Client {
	var l logger.Logger
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

	return oauth.NewClient(accessToken, clientID, wl.APIURL, l)
}

func handleError(err error) {
	fmt.Printf("exiting - error: %v\n", err)
	os.Exit(1)
}

func renderOutput(output interface{}, err error) {
	if err != nil {
		handleError(err)
	}

	var data []byte
	if useJSON {
		data, err = json.Marshal(output)
		data = append(data, '\n')
	} else {
		data, err = yaml.Marshal(output)
	}

	if err != nil {
		fmt.Printf("exiting - failed to render output - error: %v\n", err)
		os.Exit(1)
	}

	// The JSON package escapes & which we do not want.
	// It also escapes < and > but those are not present in URLs
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)

	fmt.Printf("%s", string(data))
}

func splitStringToUints(input string) ([]uint, error) {
	split := strings.Split(input, ",")
	splitUints := make([]uint, len(split))

	for i, s := range split {
		s = strings.TrimSpace(s)
		idInt, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("%v at index %d", err, i)
		}
		splitUints[i] = uint(idInt)
	}

	return splitUints, nil
}
