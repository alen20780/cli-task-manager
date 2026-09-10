package cmd

import (
	"fmt"
	"strconv"
	"time"
	"todo/internal/task"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done <id>",
	Short: "Mark a task as completed",
	Long:  `Mark a task as completed by specifying its ID. If already completed, no redundant update occurs.`,
	Example: `  todo done 1
  todo done 42`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("task ID is required. Usage: todo done <id>")
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

		if tasks[idx].Completed {
			fmt.Fprintf(cmd.OutOrStdout(), "Task #%d is already completed.\n", id)
			return nil
		}

		tasks[idx].Completed = true
		tasks[idx].UpdatedAt = time.Now()

		if err := store.SaveTasks(tasks); err != nil {
			return fmt.Errorf("failed to save task: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Task #%d marked as completed.\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
