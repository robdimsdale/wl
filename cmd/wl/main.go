package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/logger"
	"github.com/robdimsdale/wundergo/oauth"
)

var (
	// version is deliberately left uninitialized so it can be set at compile-time
	version string

	accessToken = flag.String("accessToken", "", "Wunderlist access token")
	clientID    = flag.String("clientID", "", "Wunderlist client ID")

	logLevel = flag.String("logLevel", "info", "log level: debug, info, error or fatal")
)

func main() {
	if version == "" {
		version = "dev"
	}

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "version" || arg == "-v" || arg == "--version" {
			fmt.Printf("%s\n", version)
			os.Exit(0)
		}
	}

	flag.Parse()

	logger := logger.NewLogger(logger.LogLevelFromString(*logLevel))

	var wlAccessToken string
	if *accessToken != "" {
		wlAccessToken = *accessToken
	} else {
		wlAccessToken = os.Getenv("WL_ACCESS_TOKEN")
	}

	if wlAccessToken == "" {
		logger.Error(
			"exiting",
			errors.New("accessToken not found. Either provide the flag -accessToken or set the environment variable WL_ACCESS_TOKEN"))
		os.Exit(2)
	}

	var wlClientID string
	if *clientID != "" {
		wlClientID = *clientID
	} else {
		wlClientID = os.Getenv("WL_CLIENT_ID")
	}

	if wlClientID == "" {
		logger.Error(
			"exiting",
			errors.New("clientID not found. Either provide the flag -clientID or set the environment variable WL_CLIENT_ID"))
		os.Exit(2)
	}

	client := oauth.NewClient(
		wlAccessToken,
		wlClientID,
		wundergo.APIURL,
		logger,
	)

	args := flag.Args()
	if len(args) == 0 {
		logger.Info("no command specified - exiting")
		os.Exit(2)
	}

	if args[0] == "folders" {
		folders, err := client.Folders()
		if err != nil {
			logger.Error("exiting", err)
			os.Exit(1)
		}
		json.NewEncoder(os.Stdout).Encode(folders)
	}

	if args[0] == "delete-all-folders" {
		err := client.DeleteAllFolders()
		if err != nil {
			logger.Error("exiting", err)
			os.Exit(1)
		}
		fmt.Printf("All folders deleted successfully")
	}

	if args[0] == "lists" {
		lists, err := client.Lists()
		if err != nil {
			logger.Error("exiting", err)
			os.Exit(1)
		}
		json.NewEncoder(os.Stdout).Encode(lists)
	}

	if args[0] == "delete-all-lists" {
		err := client.DeleteAllLists()
		if err != nil {
			logger.Error("exiting", err)
			os.Exit(1)
		}
		fmt.Printf("All lists deleted successfully")
	}

	if args[0] == "tasks" {
		tasks, err := client.Tasks()
		if err != nil {
			logger.Error("exiting", err)
			os.Exit(1)
		}
		json.NewEncoder(os.Stdout).Encode(tasks)
	}

	if args[0] == "delete-all-tasks" {
		err := client.DeleteAllTasks()
		if err != nil {
			logger.Error("exiting", err)
			os.Exit(1)
		}
		fmt.Printf("All tasks deleted successfully")
	}
}
