package commands

import "github.com/spf13/cobra"

var (
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
)
