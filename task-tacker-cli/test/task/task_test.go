package task_test

import (
	"task-tacker-cli/task"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	id := 1
	description := "Test task"
	task := task.NewTask(id, description)

	assert.Equal(t, id, task.ID)
	assert.Equal(t, description, task.Description)
	assert.Equal(t, task.Status, task.Status)
	assert.WithinDuration(t, time.Now(), task.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), task.UpdatedAt, time.Second)
}

func TestUpdateDescription(t *testing.T) {
	task := task.NewTask(1, "Initial description")
	newDescription := "Updated description"
	task.UpdateDescription(newDescription)

	assert.Equal(t, newDescription, task.Description)
	assert.WithinDuration(t, time.Now(), task.UpdatedAt, time.Second)
}

func TestMarkInProgress(t *testing.T) {
	task := task.NewTask(1, "Test task")
	task.MarkInProgress()

	assert.Equal(t, task.Status, task.Status)
	assert.WithinDuration(t, time.Now(), task.UpdatedAt, time.Second)
}

func TestMarkDone(t *testing.T) {
	task := task.NewTask(1, "Test task")
	task.MarkDone()

	assert.Equal(t, task.Status, task.Status)
	assert.WithinDuration(t, time.Now(), task.UpdatedAt, time.Second)
}
