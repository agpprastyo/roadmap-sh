package task

import (
	"time"
)

// TaskStatus represents the status of a task.
type TaskStatus string

const (
	// StatusTodo represents a task that is yet to be started.
	StatusTodo TaskStatus = "todo"
	// StatusInProgress represents a task that is currently in progress.
	StatusInProgress TaskStatus = "in-progress"
	// StatusDone represents a task that is completed.
	StatusDone TaskStatus = "done"
)

// Task represents a task with its properties.
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// NewTask creates a new task with the given description.
func NewTask(id int, description string) *Task {
	currentTime := time.Now()
	return &Task{
		ID:          id,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}
}

// UpdateDescription updates the task's description.
func (t *Task) UpdateDescription(newDescription string) {
	t.Description = newDescription
	t.UpdatedAt = time.Now()
}

// MarkInProgress sets the task's status to "in-progress".
func (t *Task) MarkInProgress() {
	t.Status = StatusInProgress
	t.UpdatedAt = time.Now()
}

// MarkDone sets the task's status to "done".
func (t *Task) MarkDone() {
	t.Status = StatusDone
	t.UpdatedAt = time.Now()
}
