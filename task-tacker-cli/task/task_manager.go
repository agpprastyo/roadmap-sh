package task

import (
	"errors"
	"fmt"
	"time"
)

// TaskManager manages a collection of tasks.
type TaskManager struct {
	Tasks []Task
}

// NewTaskManager creates a new TaskManager.
func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks: []Task{},
	}
}

// AddTask adds a new task to the manager.
func (tm *TaskManager) AddTask(description string) {
	id := len(tm.Tasks) + 1
	task := NewTask(id, description)
	tm.Tasks = append(tm.Tasks, *task)
	fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
}

// UpdateTask updates the description of a task by its ID.
func (tm *TaskManager) UpdateTask(id int, newDescription string) error {
	for i, task := range tm.Tasks {
		if task.ID == id {
			tm.Tasks[i].UpdateDescription(newDescription)
			fmt.Printf("Task (ID: %d) updated successfully\n", id)
			return nil
		}
	}
	return errors.New("task not found")
}

// DeleteTask deletes a task by its ID.
func (tm *TaskManager) DeleteTask(id int) error {
	for i, task := range tm.Tasks {
		if task.ID == id {
			tm.Tasks = append(tm.Tasks[:i], tm.Tasks[i+1:]...)
			fmt.Printf("Task (ID: %d) deleted successfully\n", id)
			return nil
		}
	}
	return errors.New("task not found")
}

// MarkTaskInProgress marks a task as in-progress by its ID.
func (tm *TaskManager) MarkTaskInProgress(id int) error {
	for i, task := range tm.Tasks {
		if task.ID == id {
			tm.Tasks[i].MarkInProgress()
			fmt.Printf("Task (ID: %d) marked as in-progress\n", id)
			return nil
		}
	}
	return errors.New("task not found")
}

// MarkTaskDone marks a task as done by its ID.
func (tm *TaskManager) MarkTaskDone(id int) error {
	for i, task := range tm.Tasks {
		if task.ID == id {
			tm.Tasks[i].MarkDone()
			fmt.Printf("Task (ID: %d) marked as done\n", id)
			return nil
		}
	}
	return errors.New("task not found")
}

// ListTasks lists all tasks.
func (tm *TaskManager) ListTasks() {
	fmt.Println("All Tasks:")
	for _, task := range tm.Tasks {
		printTask(task)
	}
}

// ListTasksByStatus lists tasks by their status (todo, in-progress, done).
func (tm *TaskManager) ListTasksByStatus(status TaskStatus) {
	fmt.Printf("Tasks with status '%s':\n", status)
	for _, task := range tm.Tasks {
		if task.Status == status {
			printTask(task)
		}
	}
}

// printTask is a helper function to print a task's details.
func printTask(task Task) {
	fmt.Printf("ID: %d, Description: %s, Status: %s, CreatedAt: %s, UpdatedAt: %s\n",
		task.ID, task.Description, task.Status, task.CreatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339))
}
