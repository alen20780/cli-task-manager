package cmd

import (
	"fmt"
	"strconv"
	"time"
	"todo/internal/task"

	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo <id>",
	Short: "Mark a completed task as pending",
	Long:  `Revert a completed task back to pending state by specifying its ID.`,
	Example: `  todo undo 1
  todo undo 42`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("task ID is required. Usage: todo undo <id>")
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

		if !tasks[idx].Completed {
			fmt.Fprintf(cmd.OutOrStdout(), "Task #%d is already pending.\n", id)
			return nil
		}

		tasks[idx].Completed = false
		tasks[idx].UpdatedAt = time.Now()

		if err := store.SaveTasks(tasks); err != nil {
			return fmt.Errorf("failed to save task: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Task #%d marked as pending.\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
}
