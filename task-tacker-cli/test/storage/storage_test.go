package storage_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"task-tacker-cli/storage"
	"task-tacker-cli/task"
	"testing"
)

func TestLoadTasks_FileNotExist(t *testing.T) {
	tasks, err := storage.LoadTasks("nonexistent_file.json")
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestLoadTasks_EmptyFile(t *testing.T) {
	file, err := os.CreateTemp("", "empty*.json")
	assert.NoError(t, err)
	defer func(name string) {
		err := file.Close()
		if err != nil {
			t.Error(err)
		}
		err = os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(file.Name())

	tasks, err := storage.LoadTasks(file.Name())
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestLoadTasks_ValidFile(t *testing.T) {
	file, err := os.CreateTemp("", "tasks*.json")
	assert.NoError(t, err)
	defer func(name string) {
		err := file.Close()
		if err != nil {
			t.Error(err)
		}
		err = os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(file.Name())

	tasks := []task.Task{
		{ID: 1, Description: "Task 1", Status: task.StatusTodo},
		{ID: 2, Description: "Task 2", Status: task.StatusInProgress},
	}
	data, err := json.Marshal(tasks)
	assert.NoError(t, err)
	err = os.WriteFile(file.Name(), data, 0644)
	assert.NoError(t, err)

	loadedTasks, err := storage.LoadTasks(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, tasks, loadedTasks)
}

func TestSaveTasks(t *testing.T) {
	file, err := os.CreateTemp("", "tasks*.json")
	assert.NoError(t, err)
	defer func(name string) {
		err := file.Close()
		if err != nil {
			t.Error(err)
		}
		err = os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(file.Name())

	tasks := []task.Task{
		{ID: 1, Description: "Task 1", Status: task.StatusTodo},
		{ID: 2, Description: "Task 2", Status: task.StatusInProgress},
	}

	err = storage.SaveTasks(file.Name(), tasks)
	assert.NoError(t, err)

	data, err := os.ReadFile(file.Name())
	assert.NoError(t, err)

	var loadedTasks []task.Task
	err = json.Unmarshal(data, &loadedTasks)
	assert.NoError(t, err)
	assert.Equal(t, tasks, loadedTasks)
}
