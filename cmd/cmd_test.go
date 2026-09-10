package cmd

import (
	"bytes"
	"strings"
	"testing"
	"todo/internal/storage"
)

func setupTestStorage(t *testing.T) storage.Storage {
	t.Helper()
	tempDir := t.TempDir()
	store, err := storage.NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("failed to create test storage: %v", err)
	}
	SetStorage(store)
	return store
}

func executeCommand(args ...string) (string, error) {
	return executeCommandWithInput("", args...)
}

func executeCommandWithInput(input string, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	inBuf := strings.NewReader(input)
	root := RootCmd()
	root.SetIn(inBuf)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	// Reset flags that may persist across executions in tests
	deleteCmd.Flags().Set("force", "false")
	clearCmd.Flags().Set("force", "false")
	listCmd.Flags().Set("all", "false")
	listCmd.Flags().Set("completed", "false")
	listCmd.Flags().Set("pending", "false")
	listCmd.Flags().Set("sort", "id")

	err := root.Execute()
	return buf.String(), err
}

func TestCLI_AddAndList(t *testing.T) {
	setupTestStorage(t)

	// List on empty tasks
	out, err := executeCommand("list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No tasks found.") {
		t.Fatalf("expected 'No tasks found.', got: %s", out)
	}

	// Add a task
	out, err = executeCommand("add", "Learn Go")
	if err != nil {
		t.Fatalf("unexpected error adding task: %v", err)
	}
	if !strings.Contains(out, "Task added successfully.") || !strings.Contains(out, "ID:    1") {
		t.Fatalf("unexpected output adding task: %s", out)
	}

	// Add second task
	out, err = executeCommand("add", "Build CLI")
	if err != nil {
		t.Fatalf("unexpected error adding second task: %v", err)
	}
	if !strings.Contains(out, "ID:    2") {
		t.Fatalf("unexpected output: %s", out)
	}

	// List tasks
	out, err = executeCommand("list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Learn Go") || !strings.Contains(out, "Build CLI") {
		t.Fatalf("expected tasks in list, got: %s", out)
	}
}

func TestCLI_DoneAndUndo(t *testing.T) {
	setupTestStorage(t)

	_, _ = executeCommand("add", "Task to complete")

	// Mark done
	out, err := executeCommand("done", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "marked as completed") {
		t.Fatalf("unexpected output: %s", out)
	}

	// Mark done again (idempotent notification)
	out, err = executeCommand("done", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "already completed") {
		t.Fatalf("unexpected output: %s", out)
	}

	// Check status in list
	out, _ = executeCommand("list")
	if !strings.Contains(out, "[x]") {
		t.Fatalf("expected [x] in list output: %s", out)
	}

	// Undo
	out, err = executeCommand("undo", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "marked as pending") {
		t.Fatalf("unexpected output: %s", out)
	}

	// Undo again
	out, err = executeCommand("undo", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "already pending") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestCLI_ShowAndEdit(t *testing.T) {
	setupTestStorage(t)

	_, _ = executeCommand("add", "Initial Title")

	// Show
	out, err := executeCommand("show", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Title:      Initial Title") || !strings.Contains(out, "Status:     Pending") {
		t.Fatalf("unexpected show output: %s", out)
	}

	// Edit
	out, err = executeCommand("edit", "1", "Updated Title")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Task #1 updated successfully.") {
		t.Fatalf("unexpected edit output: %s", out)
	}

	// Verify show reflects updated title
	out, err = executeCommand("show", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Title:      Updated Title") {
		t.Fatalf("unexpected show output after edit: %s", out)
	}
}

func TestCLI_DeleteAndIDRetention(t *testing.T) {
	setupTestStorage(t)

	_, _ = executeCommand("add", "Task 1")
	_, _ = executeCommand("add", "Task 2")
	_, _ = executeCommand("add", "Task 3")

	// Delete task 2 with force
	out, err := executeCommand("delete", "2", "--force")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Task #2 deleted successfully.") {
		t.Fatalf("unexpected output: %s", out)
	}

	// Add a new task - ID must be 4, NOT 3 (prevent ID collision with existing task 3)
	out, err = executeCommand("add", "Task 4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "ID:    4") {
		t.Fatalf("expected ID 4, got: %s", out)
	}

	// Interactive cancellation
	out, err = executeCommandWithInput("n\n", "delete", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Deletion canceled.") {
		t.Fatalf("expected deletion canceled, got: %s", out)
	}

	// Interactive confirmation
	out, err = executeCommandWithInput("y\n", "delete", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Task #1 deleted successfully.") {
		t.Fatalf("expected task 1 deleted, got: %s", out)
	}
}

func TestCLI_Clear(t *testing.T) {
	setupTestStorage(t)

	_, _ = executeCommand("add", "Task 1")
	_, _ = executeCommand("add", "Task 2")

	// Interactive cancel
	out, err := executeCommandWithInput("n\n", "clear")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Clear canceled.") {
		t.Fatalf("expected clear canceled, got: %s", out)
	}

	// Clear with force
	out, err = executeCommand("clear", "--force")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "All tasks cleared successfully.") {
		t.Fatalf("unexpected output: %s", out)
	}

	// List should now be empty
	out, _ = executeCommand("list")
	if !strings.Contains(out, "No tasks found.") {
		t.Fatalf("expected no tasks found, got: %s", out)
	}
}

func TestCLI_Version(t *testing.T) {
	out, err := executeCommand("version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "todo version 1.0.0") {
		t.Fatalf("unexpected version output: %s", out)
	}
}

func TestCLI_ValidationErrors(t *testing.T) {
	setupTestStorage(t)

	// Add with empty title
	_, err := executeCommand("add", "   ")
	if err == nil || !strings.Contains(err.Error(), "cannot be empty") {
		t.Fatalf("expected empty title error, got: %v", err)
	}

	// Show non-existent ID
	_, err = executeCommand("show", "999")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got: %v", err)
	}

	// Show invalid ID format
	_, err = executeCommand("show", "abc")
	if err == nil || !strings.Contains(err.Error(), "invalid task ID") {
		t.Fatalf("expected invalid ID error, got: %v", err)
	}
}
