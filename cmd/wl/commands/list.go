package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robdimsdale/wl"
	"github.com/spf13/cobra"
)

var (
	// Commands
	cmdLists = &cobra.Command{
		Use:   "lists",
		Short: "gets all lists",
		Long: `lists gets the user's lists.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).Lists())
		},
	}

	cmdList = &cobra.Command{
		Use:   "list <list-id>",
		Short: "gets the list for the provided list id",
		Long: `list gets a list specified by <list-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(list(cmd, args))
		},
	}

	cmdCreateList = &cobra.Command{
		Use:   "create-list <title>",
		Short: "creates a list with the specified <title>",
		Long: `create-list creates a list with the specified <title>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Printf("incorrect number of arguments provided\n\n")
				cmd.Usage()
				os.Exit(2)
			}

			title := args[0]

			renderOutput(newClient(cmd).CreateList(
				title,
			))
		},
	}

	cmdUpdateList = &cobra.Command{
		Use:   "update-list <list-id> [flags]",
		Short: "updates the list",
		Long: `update-list obtains the current state of the list specified by <list-id>,
and updates fields with the provided flags.
`,
		Run: func(cmd *cobra.Command, args []string) {
			list, err := list(cmd, args)
			if err != nil {
				fmt.Printf("error getting list: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			if title != "" {
				list.Title = title
			}

			renderOutput(newClient(cmd).UpdateList(list))
		},
	}

	cmdDeleteList = &cobra.Command{
		Use:   "delete-list <list-id>",
		Short: "deletes the list for the provided list id",
		Long: `delete-list deletes the list specified by <list-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			list, err := list(cmd, args)
			if err != nil {
				fmt.Printf("error getting list: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = newClient(cmd).DeleteList(list)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("list %d deleted successfully\n", list.ID)
		},
	}

	cmdDeleteAllLists = &cobra.Command{
		Use:   "delete-all-lists",
		Short: "deletes all lists",
		Long: `delete-all-lists deletes all lists.
        This will not delete the inbox, as it cannot be deleted.
				This deletes all tasks that are not in the inbox,
        and all folders that the inbox is not a member of.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			err := newClient(cmd).DeleteAllLists()
			if err != nil {
				handleError(err)
			}
			fmt.Printf("all lists deleted successfully\n")
		},
	}
)

func init() {
	cmdUpdateList.Flags().StringVar(&title, titleLongFlag, "", "title of list")
}

func list(cmd *cobra.Command, args []string) (wl.List, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing list ID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).List(id)
}
