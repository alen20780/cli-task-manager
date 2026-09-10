package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"todo/internal/task"
)

var (
	// ErrStorageCorrupted is returned when tasks storage file cannot be parsed.
	ErrStorageCorrupted = errors.New("corrupted tasks data file")
)

// Storage defines the interface for persisting and retrieving tasks.
type Storage interface {
	LoadTasks() ([]task.Task, error)
	SaveTasks(tasks []task.Task) error
	GetPath() string
}

// ResolveDataDir determines the storage directory path following OS conventions
// and environment variable overrides (TODO_DATA_DIR).
func ResolveDataDir() (string, error) {
	if envDir := os.Getenv("TODO_DATA_DIR"); envDir != "" {
		return envDir, nil
	}

	if runtime.GOOS == "windows" {
		// Prefer %APPDATA%\todo
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "todo"), nil
		}
		// Fallback to UserConfigDir
		configDir, err := os.UserConfigDir()
		if err == nil && configDir != "" {
			return filepath.Join(configDir, "todo"), nil
		}
	}

	// Linux/macOS and fallback: ~/.todo
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".todo"), nil
}
