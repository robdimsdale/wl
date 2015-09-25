package commands

import "github.com/spf13/cobra"

var (
	cmdUser = &cobra.Command{
		Use:   "user",
		Short: "gets user",
		Long: `user gets the logged-in user's information.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).User())
		},
	}
)
