package cmd

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"todo/internal/task"

	"github.com/spf13/cobra"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a task",
	Long: `Delete a task by specifying its ID. 
Prompts for confirmation unless the --force (-f) flag is passed.`,
	Example: `  todo delete 1
  todo delete 1 --force
  todo delete 1 -f`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("task ID is required. Usage: todo delete <id>")
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

		targetTask := tasks[idx]

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Delete task #%d %q? [y/N]: ", targetTask.ID, targetTask.Title)
			reader := bufio.NewReader(cmd.InOrStdin())
			response, _ := reader.ReadString('\n')
			trimmed := strings.ToLower(strings.TrimSpace(response))
			if trimmed != "y" && trimmed != "yes" {
				fmt.Fprintln(out, "Deletion canceled.")
				return nil
			}
		}

		// Remove element
		tasks = append(tasks[:idx], tasks[idx+1:]...)

		if err := store.SaveTasks(tasks); err != nil {
			return fmt.Errorf("failed to save tasks after deletion: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Task #%d deleted successfully.\n", id)
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Delete without confirmation")
	rootCmd.AddCommand(deleteCmd)
}
