package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the application version",
	Long:    `Print the current version of the todo CLI application.`,
	Example: `  todo version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "todo version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
