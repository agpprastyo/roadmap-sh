package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"task-tacker-cli/task"
)

func LoadTasks(filePath string) ([]task.Task, error) {
	var tasks []task.Task

	// Check if the file exists
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		// File does not exist, return an empty task list
		return tasks, nil
	}

	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Handle empty file case
	if len(data) == 0 {
		return tasks, nil
	}

	// Unmarshal the JSON data into the tasks slice
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// SaveTasks saves the tasks to the JSON file.
func SaveTasks(filePath string, tasks []task.Task) error {
	// Marshal the tasks slice into JSON data
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	return nil
}
