package cmd_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"task-tacker-cli/cmd"
	"task-tacker-cli/task"
)

func setup() (*cmd.CommandLineApp, *os.File) {
	file, err := os.CreateTemp("", "tasks*.json")
	if err != nil {
		panic(err)
	}
	app := cmd.NewCommandLineApp(file.Name())
	return app, file
}

func teardown(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
	err = os.Remove(file.Name())
	if err != nil {
		panic(err)
	}
}

func TestNewCommandLineApp(t *testing.T) {
	app, file := setup()
	defer teardown(file)

	assert.NotNil(t, app)
	assert.NotNil(t, app.TaskManager)
	assert.Equal(t, file.Name(), app.FilePath)
}

func TestCommandLineApp_Run_Add(t *testing.T) {
	app, file := setup()
	defer teardown(file)

	os.Args = []string{"task-cli", "add", "Test task"}
	app.Run()

	assert.Len(t, app.TaskManager.Tasks, 1)
	assert.Equal(t, "Test task", app.TaskManager.Tasks[0].Description)
}

func TestCommandLineApp_Run_Update(t *testing.T) {
	app, file := setup()
	defer teardown(file)

	os.Args = []string{"task-cli", "add", "Initial description"}
	app.Run()

	os.Args = []string{"task-cli", "update", "1", "Updated description"}
	app.Run()

	assert.Equal(t, "Updated description", app.TaskManager.Tasks[0].Description)
}

func TestCommandLineApp_Run_Delete(t *testing.T) {
	app, file := setup()
	defer teardown(file)

	os.Args = []string{"task-cli", "add", "Test task"}
	app.Run()

	os.Args = []string{"task-cli", "delete", "1"}
	app.Run()

	assert.Empty(t, app.TaskManager.Tasks)
}

func TestCommandLineApp_Run_MarkInProgress(t *testing.T) {
	app, file := setup()
	defer teardown(file)

	os.Args = []string{"task-cli", "add", "Test task"}
	app.Run()

	os.Args = []string{"task-cli", "mark-in-progress", "1"}
	app.Run()

	assert.Equal(t, task.StatusInProgress, app.TaskManager.Tasks[0].Status)
}

func TestCommandLineApp_Run_MarkDone(t *testing.T) {
	app, file := setup()
	defer teardown(file)

	os.Args = []string{"task-cli", "add", "Test task"}
	app.Run()

	os.Args = []string{"task-cli", "mark-done", "1"}
	app.Run()

	assert.Equal(t, task.StatusDone, app.TaskManager.Tasks[0].Status)
}

func TestCommandLineApp_Run_List(t *testing.T) {
	app, file := setup()
	defer teardown(file)

	os.Args = []string{"task-cli", "add", "Test task 1"}
	app.Run()
	os.Args = []string{"task-cli", "add", "Test task 2"}
	app.Run()

	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf("failed to create pipe: %v", err)
		return
	}

	old := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()

	os.Args = []string{"task-cli", "list"}
	app.Run()

	err = w.Close()
	if err != nil {
		t.Errorf("failed to close pipe: %v", err)
		return
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Errorf("failed to copy pipe to buffer: %v", err)
		return
	}

	output := buf.String()
	assert.Contains(t, output, "Test task 1")
	assert.Contains(t, output, "Test task 2")
}

func TestCommandLineApp_Run_ListByStatus(t *testing.T) {
	app, file := setup()
	defer teardown(file)

	os.Args = []string{"task-cli", "add", "Test task 1"}
	app.Run()
	os.Args = []string{"task-cli", "add", "Test task 2"}
	app.Run()
	os.Args = []string{"task-cli", "mark-in-progress", "1"}
	app.Run()

	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf("failed to create pipe: %v", err)
		return
	}

	old := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()

	os.Args = []string{"task-cli", "list", "in-progress"}
	app.Run()

	err = w.Close()
	if err != nil {
		t.Errorf("failed to close pipe: %v", err)
		return
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Errorf("failed to copy pipe to buffer: %v", err)
		return
	}

	output := buf.String()
	assert.Contains(t, output, "Test task 1")
	assert.NotContains(t, output, "Test task 2")
}
