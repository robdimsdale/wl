package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robdimsdale/wl"
	"github.com/spf13/cobra"
)

const (
	dateLongFlag = "date"
)

var (
	// Flags
	date string

	// Commands
	cmdReminders = &cobra.Command{
		Use:   "reminders",
		Short: "gets all reminders",
		Long: `reminders gets the user's reminders.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			// Currently sending completed=false is the same as not sending completed
			// Checking for whether the flag has changed protects us from potentially
			// breaking changes, i.e. if the reminders endpoint decides to return all tasks,
			// not just non-completed ones.

			if taskID != 0 {
				renderOutput(newClient(cmd).RemindersForTaskID(taskID))
			} else if listID != 0 {
				renderOutput(newClient(cmd).RemindersForListID(listID))
			} else {
				renderOutput(newClient(cmd).Reminders())
			}
		},
	}

	cmdReminder = &cobra.Command{
		Use:   "reminder <reminder-id>",
		Short: "gets the reminder for the provided reminder id",
		Long: `reminder gets a reminder specified by <reminder-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(reminder(cmd, args))
		},
	}

	cmdCreateReminder = &cobra.Command{
		Use:   "create-reminder",
		Short: "creates a reminder with the specified args",
		Long: `create-reminder creates a reminder with the specified args
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).CreateReminder(
				date,
				taskID,
				"",
			))
		},
	}

	cmdUpdateReminder = &cobra.Command{
		Use:   "update-reminder",
		Short: "updates a reminder with the specified args",
		Long: `update-reminder obtains the current state of the reminder,
and updates fields with the provided flags.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			reminder, err := reminder(cmd, args)
			if err != nil {
				handleError(err)
			}

			if cmd.Flags().Changed(dateLongFlag) {
				reminder.Date = date
			}

			renderOutput(newClient(cmd).UpdateReminder(reminder))
		},
	}

	cmdDeleteReminder = &cobra.Command{
		Use:   "delete-reminder <reminder-id>",
		Short: "deletes the reminder for the provided reminder id",
		Long: `delete-reminder deletes the reminder specified by <reminder-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			reminder, err := reminder(cmd, args)
			if err != nil {
				fmt.Printf("error getting reminder: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = newClient(cmd).DeleteReminder(reminder)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("reminder %d deleted successfully\n", reminder.ID)
		},
	}
)

func init() {
	cmdReminders.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
	cmdReminders.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "filter by taskID")

	cmdCreateReminder.Flags().UintVarP(&taskID, taskIDLongFlag, taskIDShortFlag, 0, "id of task to which reminder belongs")
	cmdCreateReminder.Flags().StringVar(&date, dateLongFlag, "", "reminder date")

	cmdUpdateReminder.Flags().StringVar(&date, dateLongFlag, "", "reminder date")
}

func reminder(cmd *cobra.Command, args []string) (wl.Reminder, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing reminderID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).Reminder(id)
}
