package budget

import (
	"errors"
	"fmt"
	"time"
)

// Service provides methods to manage budgets.
type Service struct {
	budgets map[string]Budget
}

// NewService creates a new Service instance.
func NewService() *Service {
	return &Service{
		budgets: make(map[string]Budget),
	}
}

// SetBudget sets a budget for a specific month, year, and category.
func (s *Service) SetBudget(month time.Month, year int, category string, amount float64) (*Budget, error) {
	if amount <= 0 {
		return nil, errors.New("budget amount must be greater than zero")
	}

	key := generateBudgetKey(month, year, category)
	budget := NewBudget(month, year, category, amount)
	s.budgets[key] = *budget

	return budget, nil
}

// GetBudget retrieves a budget for a specific month, year, and category.
func (s *Service) GetBudget(month time.Month, year int, category string) (*Budget, error) {
	key := generateBudgetKey(month, year, category)
	if budget, exists := s.budgets[key]; exists {
		return &budget, nil
	}

	return nil, errors.New("budget not found")
}

// CheckIfExceeded checks if the total expenses have exceeded the budget.
func (s *Service) CheckIfExceeded(month time.Month, year int, category string, totalExpenses float64) (bool, error) {
	budget, err := s.GetBudget(month, year, category)
	if err != nil {
		return false, err
	}

	return totalExpenses > budget.Amount, nil
}

// generateBudgetKey creates a unique key for a budget based on month, year, and category.
func generateBudgetKey(month time.Month, year int, category string) string {
	return fmt.Sprintf("%04d-%02d-%s", year, month, category)
}
