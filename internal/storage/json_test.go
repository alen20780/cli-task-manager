package storage

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
	"todo/internal/task"
)

func TestJSONStorage_LoadSave(t *testing.T) {
	tempDir := t.TempDir()
	store, err := NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	// 1. Loading from non-existent file should yield empty tasks and no error
	tasks, err := store.LoadTasks()
	if err != nil {
		t.Fatalf("unexpected error loading non-existent file: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected empty tasks, got %d", len(tasks))
	}

	// 2. Save tasks
	now := time.Date(2026, 9, 10, 10, 0, 0, 0, time.UTC)
	initialTasks := []task.Task{
		{ID: 1, Title: "Task 1", Completed: false, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Title: "Task 2", Completed: true, CreatedAt: now, UpdatedAt: now},
	}

	if err := store.SaveTasks(initialTasks); err != nil {
		t.Fatalf("failed to save tasks: %v", err)
	}

	// 3. Load tasks back
	loaded, err := store.LoadTasks()
	if err != nil {
		t.Fatalf("failed to load saved tasks: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(loaded))
	}
	if loaded[0].Title != "Task 1" || loaded[1].Title != "Task 2" {
		t.Fatalf("loaded tasks do not match: %+v", loaded)
	}
	if loaded[0].Completed != false || loaded[1].Completed != true {
		t.Fatalf("loaded tasks status mismatch: %+v", loaded)
	}

	// 4. Save empty/nil slice
	if err := store.SaveTasks(nil); err != nil {
		t.Fatalf("failed to save nil tasks: %v", err)
	}
	loadedEmpty, err := store.LoadTasks()
	if err != nil {
		t.Fatalf("failed to load empty tasks: %v", err)
	}
	if len(loadedEmpty) != 0 {
		t.Fatalf("expected 0 tasks after saving nil, got %d", len(loadedEmpty))
	}
}

func TestJSONStorage_CorruptedJSON(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "tasks.json")

	// Write invalid JSON
	if err := os.WriteFile(filePath, []byte("NOT_JSON_DATA"), 0o644); err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	store, err := NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	_, err = store.LoadTasks()
	if err == nil {
		t.Fatal("expected error for corrupted JSON, got nil")
	}
	if !errors.Is(err, ErrStorageCorrupted) {
		t.Fatalf("expected ErrStorageCorrupted, got %v", err)
	}
}

func TestJSONStorage_AutoCreateDir(t *testing.T) {
	tempDir := t.TempDir()
	nestedDir := filepath.Join(tempDir, "nested", "storage", "dir")

	store, err := NewJSONStorage(nestedDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	tasks := []task.Task{
		{ID: 1, Title: "Nested task"},
	}

	if err := store.SaveTasks(tasks); err != nil {
		t.Fatalf("failed to save tasks in nested non-existent directory: %v", err)
	}

	loaded, err := store.LoadTasks()
	if err != nil {
		t.Fatalf("failed to load tasks from nested directory: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Title != "Nested task" {
		t.Fatalf("unexpected loaded tasks: %+v", loaded)
	}
}

func TestResolveDataDir_EnvOverride(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("TODO_DATA_DIR", tempDir)

	resolved, err := ResolveDataDir()
	if err != nil {
		t.Fatalf("unexpected error resolving dir: %v", err)
	}
	if resolved != tempDir {
		t.Fatalf("expected %q, got %q", tempDir, resolved)
	}
}
