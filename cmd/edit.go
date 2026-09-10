package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"todo/internal/task"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <id> <new-title>",
	Short: "Edit a task's title",
	Long:  `Update the title of an existing task by specifying its ID and the new title.`,
	Example: `  todo edit 1 "Learn advanced Go"
  todo edit 3 "Write thorough unit tests"`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("task ID is required. Usage: todo edit <id> <new-title>")
		}
		if len(args) < 2 {
			return fmt.Errorf("new task title is required. Usage: todo edit <id> <new-title>")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			return fmt.Errorf("invalid task ID %q, must be a positive integer", args[0])
		}

		newTitle := strings.TrimSpace(strings.Join(args[1:], " "))
		if newTitle == "" {
			return fmt.Errorf("new task title cannot be empty")
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

		tasks[idx].Title = newTitle
		tasks[idx].UpdatedAt = time.Now()

		if err := store.SaveTasks(tasks); err != nil {
			return fmt.Errorf("failed to save updated task: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Task #%d updated successfully.\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
