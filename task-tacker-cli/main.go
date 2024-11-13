package main

import (
	"task-tacker-cli/cmd"
)

func main() {
	// Define the path for the tasks JSON file
	filePath := "task.json"

	// Create a new CommandLineApp instance
	app := cmd.NewCommandLineApp(filePath)

	// Run the application
	app.Run()

}
