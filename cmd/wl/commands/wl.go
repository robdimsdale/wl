package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

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
)

var (
	// Global flags
	accessToken string
	clientID    string
	verbose     bool
	useJSON     bool

	// Non-global, shared flags
	taskID uint
	listID uint
	title  string

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
	WundergoCmd.AddCommand(cmdCreateList)
	WundergoCmd.AddCommand(cmdUpdateList)
	WundergoCmd.AddCommand(cmdDeleteList)
	WundergoCmd.AddCommand(cmdDeleteAllLists)
	WundergoCmd.AddCommand(cmdList)

	WundergoCmd.AddCommand(cmdFolders)
	WundergoCmd.AddCommand(cmdDeleteAllFolders)

	WundergoCmd.AddCommand(cmdTasks)
	WundergoCmd.AddCommand(cmdTask)
	WundergoCmd.AddCommand(cmdCreateTask)
	WundergoCmd.AddCommand(cmdDeleteTask)
	WundergoCmd.AddCommand(cmdDeleteAllTasks)

	WundergoCmd.AddCommand(cmdUploadFile)
	WundergoCmd.AddCommand(cmdCreateFile)
	WundergoCmd.AddCommand(cmdFile)
	WundergoCmd.AddCommand(cmdFiles)
	WundergoCmd.AddCommand(cmdDestroyFile)
	WundergoCmd.AddCommand(cmdFilePreview)

	WundergoCmd.AddCommand(cmdUser)
	WundergoCmd.AddCommand(cmdUsers)
	WundergoCmd.AddCommand(cmdUpdateUser)
	WundergoCmd.AddCommand(cmdAvatarURL)
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
