package main

import (
	"fmt"

	"github.com/robdimsdale/wl/cmd/wl/commands"

	"github.com/spf13/cobra"
)

var (
	// version is deliberately left uninitialized so it can be set at compile-time
	version string

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
)

func main() {
	if version == "" {
		version = "dev"
	}

	commands.WLCmd.AddCommand(cmdVersion)
	commands.Execute()
}
