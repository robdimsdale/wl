package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robdimsdale/wl"
	"github.com/spf13/cobra"
)

const (
	assigneeIDLongFlag      = "assingeeID"
	recurrenceTypeLongFlag  = "recurrenceType"
	recurrenceCountLongFlag = "recurrenceCount"
	dueDateLongFlag         = "dueDate"
	starredLongFlag         = "starred"
)

var (
	// Flags
	assigneeID      uint
	recurrenceType  string
	recurrenceCount uint
	dueDate         string
	starred         bool

	// Commands
	cmdTasks = &cobra.Command{
		Use:   "tasks",
		Short: "gets all tasks",
		Long: `tasks gets the user's tasks.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			// Currently sending completed=false is the same as not sending completed
			// Checking for whether the flag has changed protects us from potentially
			// breaking changes, i.e. if the tasks endpoint decides to return all tasks,
			// not just non-completed ones.

			if listID == 0 {
				if cmd.Flags().Changed(completedLongFlag) {
					renderOutput(newClient(cmd).CompletedTasks(completed))
				} else {
					renderOutput(newClient(cmd).Tasks())
				}
			} else {
				if cmd.Flags().Changed(completedLongFlag) {
					renderOutput(newClient(cmd).CompletedTasksForListID(listID, completed))
				} else {
					renderOutput(newClient(cmd).TasksForListID(listID))
				}
			}
		},
	}

	cmdTask = &cobra.Command{
		Use:   "task <task-id>",
		Short: "gets the task for the provided task id",
		Long: `task gets a task specified by <task-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(task(cmd, args))
		},
	}

	cmdCreateTask = &cobra.Command{
		Use:   "create-task",
		Short: "creates a task with the specified args",
		Long: `create-task creates a task with the specified args
        `,
		Run: func(cmd *cobra.Command, args []string) {
			parsedDueDate, err := parseDueDate(dueDate)
			if err != nil {
				handleError(err)
			}

			renderOutput(newClient(cmd).CreateTask(
				title,
				listID,
				assigneeID,
				completed,
				recurrenceType,
				recurrenceCount,
				parsedDueDate,
				starred,
			))
		},
	}

	cmdUpdateTask = &cobra.Command{
		Use:   "update-task",
		Short: "updates a task with the specified args",
		Long: `update-task obtains the current state of the task,
and updates fields with the provided flags.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			var parsedDueDate time.Time
			if cmd.Flags().Changed(dueDateLongFlag) {
				var err error
				parsedDueDate, err = parseDueDate(dueDate)
				if err != nil {
					handleError(err)
				}
			}

			task, err := task(cmd, args)
			if err != nil {
				handleError(err)
			}

			if cmd.Flags().Changed(listIDLongFlag) {
				task.ListID = listID
			}

			if cmd.Flags().Changed(titleLongFlag) {
				task.Title = title
			}

			if cmd.Flags().Changed(assigneeIDLongFlag) {
				task.AssigneeID = assigneeID
			}

			if cmd.Flags().Changed(completedLongFlag) {
				task.Completed = completed
			}

			if cmd.Flags().Changed(recurrenceTypeLongFlag) {
				task.RecurrenceType = recurrenceType
			}

			if cmd.Flags().Changed(recurrenceCountLongFlag) {
				task.RecurrenceCount = recurrenceCount
			}

			if cmd.Flags().Changed(dueDateLongFlag) {
				task.DueDate = parsedDueDate
			}

			if cmd.Flags().Changed(starredLongFlag) {
				task.Starred = starred
			}

			renderOutput(newClient(cmd).UpdateTask(task))
		},
	}

	cmdDeleteTask = &cobra.Command{
		Use:   "delete-task <task-id>",
		Short: "deletes the task for the provided task id",
		Long: `delete-task deletes the task specified by <task-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			task, err := task(cmd, args)
			if err != nil {
				fmt.Printf("error getting task: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = newClient(cmd).DeleteTask(task)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("task %d deleted successfully\n", task.ID)
		},
	}

	cmdDeleteAllTasks = &cobra.Command{
		Use:   "delete-all-tasks",
		Short: "deletes all tasks",
		Long: `delete-all-tasks deletes all tasks.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			err := newClient(cmd).DeleteAllTasks()
			if err != nil {
				handleError(err)
			}
			fmt.Printf("all tasks deleted successfully\n")
		},
	}
)

func init() {
	cmdTasks.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")
	cmdTasks.Flags().BoolVar(&completed, completedLongFlag, false, "filter for completed tasks")

	cmdCreateTask.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "id of list to which task will belong")
	cmdCreateTask.Flags().StringVar(&title, titleLongFlag, "", "title of task")
	cmdCreateTask.Flags().UintVar(&assigneeID, assigneeIDLongFlag, 0, "id of task assignee")
	cmdCreateTask.Flags().BoolVar(&completed, completedLongFlag, false, "whether task is completed")
	cmdCreateTask.Flags().StringVar(&recurrenceType, recurrenceTypeLongFlag, "", "recurrence type")
	cmdCreateTask.Flags().UintVar(&recurrenceCount, recurrenceCountLongFlag, 0, "id of task assignee")
	cmdCreateTask.Flags().StringVar(&dueDate, dueDateLongFlag, "", "due date of task")
	cmdCreateTask.Flags().BoolVar(&starred, starredLongFlag, false, "whether task is starred")

	cmdUpdateTask.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "id of list to which task will belong")
	cmdUpdateTask.Flags().StringVar(&title, titleLongFlag, "", "title of task")
	cmdUpdateTask.Flags().UintVar(&assigneeID, assigneeIDLongFlag, 0, "id of task assignee")
	cmdUpdateTask.Flags().BoolVar(&completed, completedLongFlag, false, "whether task is completed")
	cmdUpdateTask.Flags().StringVar(&recurrenceType, recurrenceTypeLongFlag, "", "recurrence type")
	cmdUpdateTask.Flags().UintVar(&recurrenceCount, recurrenceCountLongFlag, 0, "id of task assignee")
	cmdUpdateTask.Flags().StringVar(&dueDate, dueDateLongFlag, "", "due date of task")
	cmdUpdateTask.Flags().BoolVar(&starred, starredLongFlag, false, "whether task is starred")
}

func task(cmd *cobra.Command, args []string) (wl.Task, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing taskID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).Task(id)
}

func parseDueDate(dueDate string) (time.Time, error) {
	splitDate := strings.Split(dueDate, "-")
	if len(splitDate) < 3 {
		return time.Now(), fmt.Errorf("Failed to parse dueDate into expected YYYY-MM-DD format: %s", dueDate)
	}

	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return time.Now(), err
	}

	monthInt, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return time.Now(), err
	}
	month := time.Month(monthInt)

	day, err := strconv.Atoi(splitDate[2])
	if err != nil {
		return time.Now(), err
	}

	hour := 0
	minute := 0
	second := 0
	nano := 0

	return time.Date(year, month, day, hour, minute, second, nano, time.Local), nil
}
