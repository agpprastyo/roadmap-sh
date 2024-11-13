package expense

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Storage provides methods to save and load expenses from a file.
type Storage interface {
	SaveExpenses(expenses []Expense) error
	LoadExpenses() ([]Expense, error)
}

// FileStorage is an implementation of Storage that uses a JSON file.
type FileStorage struct {
	filePath string
}

// NewFileStorage creates a new FileStorage instance.
func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		filePath: filePath,
	}
}

// SaveExpenses saves the list of expenses to a JSON file.
func (f *FileStorage) SaveExpenses(expenses []Expense) error {
	data, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal expenses: %w", err)
	}

	if err := os.WriteFile(f.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// LoadExpenses loads the list of expenses from a JSON file.
func (f *FileStorage) LoadExpenses() ([]Expense, error) {
	if _, err := os.Stat(f.filePath); errors.Is(err, os.ErrNotExist) {
		return []Expense{}, nil // Return an empty list if the file does not exist
	}

	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var expenses []Expense
	if err := json.Unmarshal(data, &expenses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal expenses: %w", err)
	}

	return expenses, nil
}
