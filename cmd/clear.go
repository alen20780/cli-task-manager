package cmd

import (
	"bufio"
	"fmt"
	"strings"
	"todo/internal/task"

	"github.com/spf13/cobra"
)

var clearForce bool

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete all tasks",
	Long: `Remove all tasks from storage.
Prompts for confirmation unless the --force (-f) flag is passed.`,
	Example: `  todo clear
  todo clear --force
  todo clear -f`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := GetStorage()
		if err != nil {
			return err
		}

		tasks, err := store.LoadTasks()
		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No tasks to clear.")
			return nil
		}

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Are you sure you want to delete all %d tasks? [y/N]: ", len(tasks))
			reader := bufio.NewReader(cmd.InOrStdin())
			response, _ := reader.ReadString('\n')
			trimmed := strings.ToLower(strings.TrimSpace(response))
			if trimmed != "y" && trimmed != "yes" {
				fmt.Fprintln(out, "Clear canceled.")
				return nil
			}
		}

		if err := store.SaveTasks([]task.Task{}); err != nil {
			return fmt.Errorf("failed to clear tasks: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "All tasks cleared successfully.")
		return nil
	},
}

func init() {
	clearCmd.Flags().BoolVarP(&clearForce, "force", "f", false, "Clear all tasks without confirmation")
	rootCmd.AddCommand(clearCmd)
}
