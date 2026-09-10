package task

import (
	"errors"
	"testing"
	"time"
)

func TestNextID(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []Task
		expected int
	}{
		{
			name:     "empty list",
			tasks:    []Task{},
			expected: 1,
		},
		{
			name: "sequential tasks",
			tasks: []Task{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
			expected: 4,
		},
		{
			name: "tasks with gaps due to deletion",
			tasks: []Task{
				{ID: 1},
				{ID: 2},
				{ID: 5},
			},
			expected: 6,
		},
		{
			name: "unordered tasks",
			tasks: []Task{
				{ID: 8},
				{ID: 2},
				{ID: 4},
			},
			expected: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextID(tt.tasks)
			if got != tt.expected {
				t.Fatalf("expected next ID %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestNew(t *testing.T) {
	now := time.Date(2026, 9, 10, 10, 0, 0, 0, time.UTC)

	t.Run("valid task", func(t *testing.T) {
		task, err := New(1, "  Learn Go  ", now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if task.ID != 1 {
			t.Errorf("expected ID 1, got %d", task.ID)
		}
		if task.Title != "Learn Go" {
			t.Errorf("expected trimmed title 'Learn Go', got '%s'", task.Title)
		}
		if task.Completed {
			t.Errorf("expected completed to be false, got true")
		}
		if !task.CreatedAt.Equal(now) || !task.UpdatedAt.Equal(now) {
			t.Errorf("expected timestamps to equal %v", now)
		}
	})

	t.Run("empty title error", func(t *testing.T) {
		_, err := New(1, "   ", now)
		if !errors.Is(err, ErrEmptyTitle) {
			t.Fatalf("expected ErrEmptyTitle, got %v", err)
		}
	})

	t.Run("invalid task ID error", func(t *testing.T) {
		_, err := New(0, "Test", now)
		if !errors.Is(err, ErrInvalidTaskID) {
			t.Fatalf("expected ErrInvalidTaskID, got %v", err)
		}

		_, err = New(-1, "Test", now)
		if !errors.Is(err, ErrInvalidTaskID) {
			t.Fatalf("expected ErrInvalidTaskID, got %v", err)
		}
	})
}

func TestFindIndex(t *testing.T) {
	tasks := []Task{
		{ID: 1, Title: "Task 1"},
		{ID: 5, Title: "Task 5"},
		{ID: 10, Title: "Task 10"},
	}

	if idx := FindIndex(tasks, 5); idx != 1 {
		t.Errorf("expected index 1, got %d", idx)
	}
	if idx := FindIndex(tasks, 99); idx != -1 {
		t.Errorf("expected index -1 for missing task, got %d", idx)
	}
}

func TestFilter(t *testing.T) {
	tasks := []Task{
		{ID: 1, Title: "Pending 1", Completed: false},
		{ID: 2, Title: "Done 1", Completed: true},
		{ID: 3, Title: "Pending 2", Completed: false},
	}

	all := Filter(tasks, FilterAll)
	if len(all) != 3 {
		t.Errorf("expected 3 tasks for FilterAll, got %d", len(all))
	}

	pending := Filter(tasks, FilterPending)
	if len(pending) != 2 || pending[0].ID != 1 || pending[1].ID != 3 {
		t.Errorf("unexpected pending tasks result: %+v", pending)
	}

	completed := Filter(tasks, FilterCompleted)
	if len(completed) != 1 || completed[0].ID != 2 {
		t.Errorf("unexpected completed tasks result: %+v", completed)
	}
}

func TestSort(t *testing.T) {
	t1 := time.Date(2026, 9, 10, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 9, 10, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)

	tasks := []Task{
		{ID: 3, Title: "Three", CreatedAt: t2, UpdatedAt: t3},
		{ID: 1, Title: "One", CreatedAt: t3, UpdatedAt: t1},
		{ID: 2, Title: "Two", CreatedAt: t1, UpdatedAt: t2},
	}

	// Sort by ID
	tasksCopy := append([]Task(nil), tasks...)
	Sort(tasksCopy, SortID)
	if tasksCopy[0].ID != 1 || tasksCopy[1].ID != 2 || tasksCopy[2].ID != 3 {
		t.Errorf("SortID failed: %+v", tasksCopy)
	}

	// Sort by CreatedAt
	tasksCopy = append([]Task(nil), tasks...)
	Sort(tasksCopy, SortCreated)
	if tasksCopy[0].ID != 2 || tasksCopy[1].ID != 3 || tasksCopy[2].ID != 1 {
		t.Errorf("SortCreated failed: %+v", tasksCopy)
	}

	// Sort by UpdatedAt
	tasksCopy = append([]Task(nil), tasks...)
	Sort(tasksCopy, SortUpdated)
	if tasksCopy[0].ID != 1 || tasksCopy[1].ID != 2 || tasksCopy[2].ID != 3 {
		t.Errorf("SortUpdated failed: %+v", tasksCopy)
	}
}
