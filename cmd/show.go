package cmd

import (
	"fmt"
	"strconv"
	"todo/internal/task"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show detailed information about a task",
	Long:  `Display complete details including Title, Status, and Timestamps for a task by its ID.`,
	Example: `  todo show 1
  todo show 42`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("task ID is required. Usage: todo show <id>")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			return fmt.Errorf("invalid task ID %q, must be a positive integer", args[0])
		}

		store, err := GetStorage()
		if err != nil {
			return err
		}

		tasks, err := store.LoadTasks()
		if err != nil {
			return err
		}

		idx := task.FindIndex(tasks, id)
		if idx == -1 {
			return fmt.Errorf("task with ID %d not found", id)
		}

		t := tasks[idx]
		status := "Pending"
		if t.Completed {
			status = "Completed"
		}

		const timeFormat = "2006-01-02 15:04:05"

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Task #%d\n\n", t.ID)
		fmt.Fprintf(out, "Title:      %s\n", t.Title)
		fmt.Fprintf(out, "Status:     %s\n", status)
		fmt.Fprintf(out, "Created:    %s\n", t.CreatedAt.Local().Format(timeFormat))
		fmt.Fprintf(out, "Updated:    %s\n", t.UpdatedAt.Local().Format(timeFormat))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
