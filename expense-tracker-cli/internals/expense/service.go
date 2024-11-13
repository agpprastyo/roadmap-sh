package expense

import (
	"errors"
	"fmt"
	"time"
)

// Service provides methods to manage expenses.
type Service struct {
	storage Storage
}

// NewService creates a new Service instance.
func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// AddExpense adds a new expense to the system.
func (s *Service) AddExpense(description string, amount float64, category string) (*Expense, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	expenses, err := s.storage.LoadExpenses()
	if err != nil {
		return nil, fmt.Errorf("failed to load expenses: %w", err)
	}

	// Generate a new ID based on the last expense ID
	newID := len(expenses) + 1
	expense := NewExpense(newID, description, amount, category)

	expenses = append(expenses, *expense)
	if err := s.storage.SaveExpenses(expenses); err != nil {
		return nil, fmt.Errorf("failed to save expense: %w", err)
	}

	return expense, nil
}

// ListExpenses returns all the expenses.
func (s *Service) ListExpenses() ([]Expense, error) {
	return s.storage.LoadExpenses()
}

// DeleteExpense deletes an expense by ID.
func (s *Service) DeleteExpense(id int) error {
	expenses, err := s.storage.LoadExpenses()
	if err != nil {
		return fmt.Errorf("failed to load expenses: %w", err)
	}

	for i, expense := range expenses {
		if expense.ID == id {
			expenses = append(expenses[:i], expenses[i+1:]...)
			if err := s.storage.SaveExpenses(expenses); err != nil {
				return fmt.Errorf("failed to save expenses: %w", err)
			}
			return nil
		}
	}

	return errors.New("expense not found")
}

// UpdateExpense updates an existing expense by ID.
func (s *Service) UpdateExpense(id int, description string, amount float64, category string) error {
	expenses, err := s.storage.LoadExpenses()
	if err != nil {
		return fmt.Errorf("failed to load expenses: %w", err)
	}

	for i, expense := range expenses {
		if expense.ID == id {
			expenses[i].Update(description, amount, category)
			if err := s.storage.SaveExpenses(expenses); err != nil {
				return fmt.Errorf("failed to save expenses: %w", err)
			}
			return nil
		}
	}

	return errors.New("expense not found")
}

// GetSummary returns the total amount spent.
func (s *Service) GetSummary() (float64, error) {
	expenses, err := s.storage.LoadExpenses()
	if err != nil {
		return 0, fmt.Errorf("failed to load expenses: %w", err)
	}

	total := 0.0
	for _, expense := range expenses {
		total += expense.Amount
	}

	return total, nil
}

// TotalExpenses calculates the total expenses.
func (s *Service) TotalExpenses() (float64, error) {
	expenses, err := s.storage.LoadExpenses()
	if err != nil {
		return 0, err
	}

	var total float64
	for _, exp := range expenses {
		total += exp.Amount
	}

	return total, nil
}

// TotalExpensesByMonth calculates the total expenses for a specific month and year.
func (s *Service) TotalExpensesByMonth(month time.Month, year int) (float64, error) {
	expenses, err := s.storage.LoadExpenses()
	if err != nil {
		return 0, err
	}

	var total float64
	for _, exp := range expenses {
		expDate := exp.Date
		if expDate.Month() == month && expDate.Year() == year {
			total += exp.Amount
		}
	}

	return total, nil
}
