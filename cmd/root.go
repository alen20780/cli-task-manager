package cmd

import (
	"fmt"
	"io"
	"os"
	"todo/internal/storage"

	"github.com/spf13/cobra"
)

var (
	appStorage storage.Storage

	// Version holds the current application version.
	Version = "1.0.0"

	// InReader and OutWriter can be overridden for testing.
	InReader  io.Reader = os.Stdin
	OutWriter io.Writer = os.Stdout
	ErrWriter io.Writer = os.Stderr
)

// SetStorage allows overriding the storage provider (used in testing or dependency injection).
func SetStorage(s storage.Storage) {
	appStorage = s
}

// GetStorage returns the active storage provider, initializing a default one if unset.
func GetStorage() (storage.Storage, error) {
	if appStorage != nil {
		return appStorage, nil
	}

	store, err := storage.NewJSONStorage("")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	appStorage = store
	return appStorage, nil
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "todo is a clean and lightweight CLI task manager",
	Long: `todo is a production-quality terminal application for managing your daily tasks
with persistent local JSON storage, simple ID tracking, and rich filtering.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	rootCmd.SetIn(InReader)
	rootCmd.SetOut(OutWriter)
	rootCmd.SetErr(ErrWriter)
	return rootCmd.Execute()
}

// RootCmd returns the root command instance (useful for testing or doc generation).
func RootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("todo version %s\n", Version))
}
