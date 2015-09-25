package commands

import "github.com/spf13/cobra"

var (
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
)
