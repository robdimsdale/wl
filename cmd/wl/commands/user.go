package commands

import "github.com/spf13/cobra"

const (
	nameLongFlag = "name"
)

var (
	// Flags
	name string

	// Commands
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

	cmdUpdateUser = &cobra.Command{
		Use:   "update-user",
		Short: "updates the user",
		Long: `update-user obtains the current state of the logged-in user,
and updates fields with the provided flags.
`,
		Run: func(cmd *cobra.Command, args []string) {
			client := newClient(cmd)
			user, err := client.User()
			if err != nil {
				handleError(err)
			}

			if name != "" {
				user.Name = name
			}
			renderOutput(client.UpdateUser(user))
		},
	}
)

func init() {
	cmdUsers.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
	cmdUpdateUser.Flags().StringVar(&name, nameLongFlag, "", "name")
}
