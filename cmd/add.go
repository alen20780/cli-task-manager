package cmd

import (
	"fmt"
	"strings"
	"time"
	"todo/internal/task"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "Add a new task",
	Long: `Add a new task with the specified title to your task list.
The task is assigned the next sequential ID and saved to local storage.`,
	Example: `  todo add "Learn Go"
  todo add "Build CLI application"`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("task title is required. Usage: todo add <title>")
		}
		if strings.TrimSpace(strings.Join(args, " ")) == "" {
			return fmt.Errorf("task title cannot be empty")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := GetStorage()
		if err != nil {
			return err
		}

		tasks, err := store.LoadTasks()
		if err != nil {
			return err
		}

		title := strings.TrimSpace(strings.Join(args, " "))
		nextID := task.NextID(tasks)
		newTask, err := task.New(nextID, title, time.Now())
		if err != nil {
			return err
		}

		tasks = append(tasks, *newTask)
		if err := store.SaveTasks(tasks); err != nil {
			return fmt.Errorf("failed to save task: %w", err)
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Task added successfully.\n\n")
		fmt.Fprintf(out, "ID:    %d\n", newTask.ID)
		fmt.Fprintf(out, "Title: %s\n", newTask.Title)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
