package cmd

import (
	"fmt"
	"os"
	"strconv"
	"task-tacker-cli/storage"
	"task-tacker-cli/task"
)

// CommandLineApp is the main structure to run the CLI application.
type CommandLineApp struct {
	TaskManager *task.TaskManager
	FilePath    string
}

// NewCommandLineApp creates a new instance of the CommandLineApp.
func NewCommandLineApp(filePath string) *CommandLineApp {
	// Load existing tasks from the file
	tasks, err := storage.LoadTasks(filePath)
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		os.Exit(1)
	}

	// Initialize the TaskManager with the loaded tasks
	taskManager := task.NewTaskManager()
	taskManager.Tasks = tasks

	return &CommandLineApp{
		TaskManager: taskManager,
		FilePath:    filePath,
	}
}

// Run parses the command-line arguments and executes the corresponding function.
func (app *CommandLineApp) Run() {
	if len(os.Args) < 2 {
		app.printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task description.")
			return
		}
		description := os.Args[2]
		app.TaskManager.AddTask(description)

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Please provide a task ID and new description.")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID.")
			return
		}
		description := os.Args[3]
		if err := app.TaskManager.UpdateTask(id, description); err != nil {
			fmt.Println(err)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task ID.")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID.")
			return
		}
		if err := app.TaskManager.DeleteTask(id); err != nil {
			fmt.Println(err)
		}

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task ID.")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID.")
			return
		}
		if err := app.TaskManager.MarkTaskInProgress(id); err != nil {
			fmt.Println(err)
		}

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task ID.")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID.")
			return
		}
		if err := app.TaskManager.MarkTaskDone(id); err != nil {
			fmt.Println(err)
		}

	case "list":
		if len(os.Args) == 2 {
			app.TaskManager.ListTasks()
		} else if len(os.Args) == 3 {
			status := os.Args[2]
			switch status {
			case "todo":
				app.TaskManager.ListTasksByStatus(task.StatusTodo)
			case "in-progress":
				app.TaskManager.ListTasksByStatus(task.StatusInProgress)
			case "done":
				app.TaskManager.ListTasksByStatus(task.StatusDone)
			default:
				fmt.Println("Invalid status. Use 'todo', 'in-progress', or 'done'.")
			}
		}

	default:
		app.printUsage()
	}

	// Save tasks after performing the operation
	if err := storage.SaveTasks(app.FilePath, app.TaskManager.Tasks); err != nil {
		fmt.Println("Error saving tasks:", err)
	}
}

// printUsage prints the usage information for the CLI.
func (app *CommandLineApp) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  task-cli add <description>                 # Add a new task")
	fmt.Println("  task-cli update <id> <new description>     # Update an existing task")
	fmt.Println("  task-cli delete <id>                       # Delete a task")
	fmt.Println("  task-cli mark-in-progress <id>             # Mark a task as in-progress")
	fmt.Println("  task-cli mark-done <id>                    # Mark a task as done")
	fmt.Println("  task-cli list [status]                     # List all tasks or filter by status")
}
