package task

import (
	"errors"
	"sort"
	"strings"
	"time"
)

var (
	// ErrEmptyTitle is returned when a task title is blank or whitespace only.
	ErrEmptyTitle = errors.New("task title cannot be empty")
	// ErrTaskNotFound is returned when a requested task does not exist.
	ErrTaskNotFound = errors.New("task not found")
	// ErrInvalidTaskID is returned when an invalid task ID (<= 0) is specified.
	ErrInvalidTaskID = errors.New("invalid task ID, must be greater than 0")
)

// Task represents a single todo item.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NextID calculates the next available ID for a new task.
// It ensures predictability and prevents ID collision when items are deleted.
func NextID(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// New creates a new Task instance with normalized title and timestamps.
func New(id int, title string, now time.Time) (*Task, error) {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return nil, ErrEmptyTitle
	}
	if id <= 0 {
		return nil, ErrInvalidTaskID
	}

	return &Task{
		ID:        id,
		Title:     trimmed,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// FindIndex locates a task in the slice by its ID and returns its index, or -1 if not found.
func FindIndex(tasks []Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

// FilterOption defines filtering criteria for task lists.
type FilterOption string

const (
	FilterAll       FilterOption = "all"
	FilterPending   FilterOption = "pending"
	FilterCompleted FilterOption = "completed"
)

// SortOption defines sorting criteria for task lists.
type SortOption string

const (
	SortID      SortOption = "id"
	SortCreated SortOption = "created"
	SortUpdated SortOption = "updated"
)

// Filter filters a slice of tasks according to the specified option.
func Filter(tasks []Task, opt FilterOption) []Task {
	if opt == FilterAll {
		result := make([]Task, len(tasks))
		copy(result, tasks)
		return result
	}

	var filtered []Task
	for _, t := range tasks {
		if opt == FilterCompleted && t.Completed {
			filtered = append(filtered, t)
		} else if opt == FilterPending && !t.Completed {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

// Sort sorts the tasks in-place based on the given sort option.
func Sort(tasks []Task, opt SortOption) {
	switch opt {
	case SortCreated:
		sort.Slice(tasks, func(i, j int) bool {
			if tasks[i].CreatedAt.Equal(tasks[j].CreatedAt) {
				return tasks[i].ID < tasks[j].ID
			}
			return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		})
	case SortUpdated:
		sort.Slice(tasks, func(i, j int) bool {
			if tasks[i].UpdatedAt.Equal(tasks[j].UpdatedAt) {
				return tasks[i].ID < tasks[j].ID
			}
			return tasks[i].UpdatedAt.Before(tasks[j].UpdatedAt)
		})
	case SortID:
		fallthrough
	default:
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].ID < tasks[j].ID
		})
	}
}
