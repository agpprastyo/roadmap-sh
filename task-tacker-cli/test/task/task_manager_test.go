package task_test

import (
	"github.com/stretchr/testify/assert"
	"task-tacker-cli/task"
	"testing"
)

func TestNewTaskManager(t *testing.T) {
	tm := task.NewTaskManager()
	assert.NotNil(t, tm)
	assert.Empty(t, tm.Tasks)
}

func TestAddTask(t *testing.T) {
	tm := task.NewTaskManager()
	tm.AddTask("Test task")
	assert.Len(t, tm.Tasks, 1)
	assert.Equal(t, "Test task", tm.Tasks[0].Description)
	assert.Equal(t, task.StatusTodo, tm.Tasks[0].Status)
}

func TestUpdateTask(t *testing.T) {
	tm := task.NewTaskManager()
	tm.AddTask("Initial description")
	err := tm.UpdateTask(1, "Updated description")
	assert.NoError(t, err)
	assert.Equal(t, "Updated description", tm.Tasks[0].Description)
}

func TestDeleteTask(t *testing.T) {
	tm := task.NewTaskManager()
	tm.AddTask("Test task")
	err := tm.DeleteTask(1)
	assert.NoError(t, err)
	assert.Empty(t, tm.Tasks)
}

func TestMarkTaskInProgress(t *testing.T) {
	tm := task.NewTaskManager()
	tm.AddTask("Test task")
	err := tm.MarkTaskInProgress(1)
	assert.NoError(t, err)
	assert.Equal(t, task.StatusInProgress, tm.Tasks[0].Status)
}

func TestMarkTaskDone(t *testing.T) {
	tm := task.NewTaskManager()
	tm.AddTask("Test task")
	err := tm.MarkTaskDone(1)
	assert.NoError(t, err)
	assert.Equal(t, task.StatusDone, tm.Tasks[0].Status)
}

func TestListTasks(t *testing.T) {
	tm := task.NewTaskManager()
	tm.AddTask("Test task 1")
	tm.AddTask("Test task 2")
	tm.ListTasks()
	assert.Len(t, tm.Tasks, 2)
}

func TestListTasksByStatus(t *testing.T) {
	tm := task.NewTaskManager()
	tm.AddTask("Test task 1")
	tm.AddTask("Test task 2")
	err := tm.MarkTaskInProgress(1)
	if err != nil {
		t.Error(err)
	}
	tm.ListTasksByStatus(task.StatusInProgress)
	assert.Equal(t, task.StatusInProgress, tm.Tasks[0].Status)
}
