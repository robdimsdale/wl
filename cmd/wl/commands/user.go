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

	cmdUsers = &cobra.Command{
		Use:   "users",
		Short: "gets users the user can access",
		Long: `users gets the users the logged-in user can access.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if listID != 0 {
				renderOutput(newClient(cmd).UsersForListID(listID))
			} else {
				renderOutput(newClient(cmd).Users())
			}
		},
	}
)

func init() {
	cmdUsers.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
}
