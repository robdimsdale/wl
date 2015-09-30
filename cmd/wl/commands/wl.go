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

	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/logger"
	"github.com/robdimsdale/wundergo/oauth"
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

	// WundergoCmd is the root command. All other commands are subcommands of it.
	WundergoCmd = &cobra.Command{Use: "wl"}
)

// Execute adds all child commands to the root command WundergoCmd,
// and executes the root command.
func Execute() {
	addCommands()
	WundergoCmd.Execute()
}

// Sets global flags
func init() {
	WundergoCmd.PersistentFlags().BoolVarP(&verbose, verboseLongFlag, verboseShortFlag, false, "verbose output")
	WundergoCmd.PersistentFlags().StringVarP(&accessToken, accessTokenLongFlag, "", "", `Wunderlist access token. 
                      	Required, but can be provided via WL_ACCESS_TOKEN environment variable instead.`)
	WundergoCmd.PersistentFlags().StringVarP(&clientID, clientIDLongFlag, "", "", `Wunderlist client ID. 
                     Required, but can be provided via WL_CLIENT_ID environment variable instead.`)
	WundergoCmd.PersistentFlags().BoolVarP(&useJSON, useJSONLongFlag, useJSONShortFlag, false, "render output as JSON instead of YAML.")
}

func addCommands() {
	WundergoCmd.AddCommand(cmdInbox)
	WundergoCmd.AddCommand(cmdRoot)
	WundergoCmd.AddCommand(cmdLists)
	WundergoCmd.AddCommand(cmdList)
	WundergoCmd.AddCommand(cmdCreateList)
	WundergoCmd.AddCommand(cmdUpdateList)
	WundergoCmd.AddCommand(cmdDeleteList)
	WundergoCmd.AddCommand(cmdDeleteAllLists)

	WundergoCmd.AddCommand(cmdFolders)
	WundergoCmd.AddCommand(cmdFolder)
	WundergoCmd.AddCommand(cmdCreateFolder)
	WundergoCmd.AddCommand(cmdUpdateFolder)
	WundergoCmd.AddCommand(cmdDeleteFolder)
	WundergoCmd.AddCommand(cmdDeleteAllFolders)

	WundergoCmd.AddCommand(cmdTasks)
	WundergoCmd.AddCommand(cmdTask)
	WundergoCmd.AddCommand(cmdCreateTask)
	WundergoCmd.AddCommand(cmdUpdateTask)
	WundergoCmd.AddCommand(cmdDeleteTask)
	WundergoCmd.AddCommand(cmdDeleteAllTasks)

	WundergoCmd.AddCommand(cmdUploadFile)
	WundergoCmd.AddCommand(cmdCreateFile)
	WundergoCmd.AddCommand(cmdFile)
	WundergoCmd.AddCommand(cmdFiles)
	WundergoCmd.AddCommand(cmdDestroyFile)
	WundergoCmd.AddCommand(cmdFilePreview)

	WundergoCmd.AddCommand(cmdUsers)
	WundergoCmd.AddCommand(cmdUser)
	WundergoCmd.AddCommand(cmdUpdateUser)
	WundergoCmd.AddCommand(cmdAvatarURL)

	WundergoCmd.AddCommand(cmdNotes)
	WundergoCmd.AddCommand(cmdNote)
	WundergoCmd.AddCommand(cmdCreateNote)
	WundergoCmd.AddCommand(cmdUpdateNote)
	WundergoCmd.AddCommand(cmdDeleteNote)

	WundergoCmd.AddCommand(cmdSubtasks)
	WundergoCmd.AddCommand(cmdCreateSubtask)
	WundergoCmd.AddCommand(cmdSubtask)
	WundergoCmd.AddCommand(cmdUpdateSubtask)
	WundergoCmd.AddCommand(cmdDeleteSubtask)

	WundergoCmd.AddCommand(cmdWebhooks)
	WundergoCmd.AddCommand(cmdWebhook)
	WundergoCmd.AddCommand(cmdCreateWebhook)
	WundergoCmd.AddCommand(cmdDeleteWebhook)

	WundergoCmd.AddCommand(cmdReminders)
	WundergoCmd.AddCommand(cmdCreateReminder)
	WundergoCmd.AddCommand(cmdReminder)
	WundergoCmd.AddCommand(cmdUpdateReminder)
	WundergoCmd.AddCommand(cmdDeleteReminder)

	WundergoCmd.AddCommand(cmdTaskComments)
	WundergoCmd.AddCommand(cmdCreateTaskComment)
	WundergoCmd.AddCommand(cmdTaskComment)
	WundergoCmd.AddCommand(cmdDeleteTaskComment)

	WundergoCmd.AddCommand(cmdMemberships)
	WundergoCmd.AddCommand(cmdMembership)
	WundergoCmd.AddCommand(cmdRejectMembership)
	WundergoCmd.AddCommand(cmdAcceptMembership)
	WundergoCmd.AddCommand(cmdRemoveMembership)
	WundergoCmd.AddCommand(cmdInviteMember)

	WundergoCmd.AddCommand(cmdListPositions)
	WundergoCmd.AddCommand(cmdListPosition)
	WundergoCmd.AddCommand(cmdUpdateListPosition)

	WundergoCmd.AddCommand(cmdTaskPositions)
	WundergoCmd.AddCommand(cmdTaskPosition)
	WundergoCmd.AddCommand(cmdUpdateTaskPosition)

	WundergoCmd.AddCommand(cmdSubtaskPositions)
	WundergoCmd.AddCommand(cmdSubtaskPosition)
	WundergoCmd.AddCommand(cmdUpdateSubtaskPosition)
}

func newClient(cmd *cobra.Command) wundergo.Client {
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

	return oauth.NewClient(accessToken, clientID, wundergo.APIURL, l)
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
		idInt, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("%v at index %d", err, i)
		}
		splitUints[i] = uint(idInt)
	}

	return splitUints, nil
}
