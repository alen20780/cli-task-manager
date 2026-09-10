package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"todo/internal/task"

	"github.com/spf13/cobra"
)

var (
	listAll       bool
	listCompleted bool
	listPending   bool
	listSort      string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long: `Display all or filtered tasks in a formatted terminal table.
By default, all tasks are displayed. You can filter by status using flags.`,
	Example: `  todo list
  todo list --all
  todo list --pending
  todo list --completed
  todo list --sort created
  todo list --sort updated`,
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
			fmt.Fprintln(cmd.OutOrStdout(), "No tasks found.")
			return nil
		}

		listCompleted, _ := cmd.Flags().GetBool("completed")
		listPending, _ := cmd.Flags().GetBool("pending")
		listSort, _ := cmd.Flags().GetString("sort")

		// Filter tasks
		filterOpt := task.FilterAll
		if listCompleted && !listPending {
			filterOpt = task.FilterCompleted
		} else if listPending && !listCompleted {
			filterOpt = task.FilterPending
		}
		// If both --completed and --pending or neither/all are set, show all

		filtered := task.Filter(tasks, filterOpt)
		if len(filtered) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No tasks found matching criteria.")
			return nil
		}

		// Sort tasks
		sortOpt := task.SortID
		switch strings.ToLower(listSort) {
		case "created":
			sortOpt = task.SortCreated
		case "updated":
			sortOpt = task.SortUpdated
		case "id", "":
			sortOpt = task.SortID
		default:
			return fmt.Errorf("invalid sort option %q. Valid options: id, created, updated", listSort)
		}

		task.Sort(filtered, sortOpt)

		// Render table
		out := cmd.OutOrStdout()
		w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATUS\tTASK")
		for _, t := range filtered {
			status := "[ ]"
			if t.Completed {
				status = "[x]"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\n", t.ID, status, t.Title)
		}
		return w.Flush()
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all tasks (default)")
	listCmd.Flags().BoolVarP(&listCompleted, "completed", "c", false, "Show completed tasks only")
	listCmd.Flags().BoolVarP(&listPending, "pending", "p", false, "Show pending tasks only")
	listCmd.Flags().StringVarP(&listSort, "sort", "s", "id", "Sort tasks by (id, created, updated)")

	rootCmd.AddCommand(listCmd)
}
