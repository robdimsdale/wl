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

	if *accessToken == "" {
		logger.Error("exiting", errors.New("accessToken must be provided"))
		os.Exit(2)
	}

	if *clientID == "" {
		logger.Error("exiting", errors.New("clientID must be provided"))
		os.Exit(2)
	}

	client := oauth.NewClient(*accessToken, *clientID, wundergo.APIURL, logger)

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
}
