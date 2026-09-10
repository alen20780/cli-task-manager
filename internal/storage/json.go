package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"todo/internal/task"
)

// JSONStorage implements the Storage interface persisting tasks as JSON.
type JSONStorage struct {
	mu       sync.Mutex
	filePath string
}

// NewJSONStorage creates a new JSONStorage instance at the specified data directory.
func NewJSONStorage(dataDir string) (*JSONStorage, error) {
	if dataDir == "" {
		resolved, err := ResolveDataDir()
		if err != nil {
			return nil, err
		}
		dataDir = resolved
	}

	filePath := filepath.Join(dataDir, "tasks.json")
	return &JSONStorage{filePath: filePath}, nil
}

// GetPath returns the full path to the storage file.
func (s *JSONStorage) GetPath() string {
	return s.filePath
}

// LoadTasks reads tasks from the JSON file.
// If the file does not exist, it returns an empty slice without an error.
func (s *JSONStorage) LoadTasks() ([]task.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []task.Task{}, nil
		}
		return nil, fmt.Errorf("failed to open storage file %q: %w", s.filePath, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage file: %w", err)
	}

	// Empty file is considered empty tasks list
	if len(data) == 0 {
		return []task.Task{}, nil
	}

	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStorageCorrupted, err)
	}

	return tasks, nil
}

// SaveTasks writes tasks to the JSON file atomically using a temporary file.
func (s *JSONStorage) SaveTasks(tasks []task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create storage directory %q: %w", dir, err)
	}

	// If tasks is nil, serialize as empty slice `[]` rather than `null`
	if tasks == nil {
		tasks = []task.Task{}
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks to JSON: %w", err)
	}
	data = append(data, '\n')

	// Write to temporary file in the same directory for atomic rename
	tempFile, err := os.CreateTemp(dir, "tasks-*.json.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary storage file: %w", err)
	}
	tempName := tempFile.Name()

	// Ensure temp file cleanup in case of error
	success := false
	defer func() {
		if !success {
			_ = os.Remove(tempName)
		}
	}()

	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to write data to temporary file: %w", err)
	}

	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to sync temporary file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// On Windows, os.Rename will overwrite if the file exists in newer Go versions,
	// but replacing safely is handled cleanly by os.Rename or remove & rename fallback.
	if err := os.Rename(tempName, s.filePath); err != nil {
		// Fallback for Windows file locking edge cases
		_ = os.Remove(s.filePath)
		if err := os.Rename(tempName, s.filePath); err != nil {
			return fmt.Errorf("failed to atomically replace storage file: %w", err)
		}
	}

	success = true
	return nil
}
