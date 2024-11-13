package expense

import (
	"time"
)

// Expense represents a single financial expense.
type Expense struct {
	ID          int       `json:"id"`          // Unique identifier for the expense
	Description string    `json:"description"` // Description of the expense
	Amount      float64   `json:"amount"`      // Amount of money spent
	Date        time.Time `json:"date"`        // Date of the expense
	Category    string    `json:"category"`    // Optional category for the expense
}

// NewExpense creates a new expense instance.
func NewExpense(id int, description string, amount float64, category string) *Expense {
	return &Expense{
		ID:          id,
		Description: description,
		Amount:      amount,
		Date:        time.Now(), // Default to current date and time
		Category:    category,
	}
}

// Update updates the fields of the expense.
func (e *Expense) Update(description string, amount float64, category string) {
	e.Description = description
	e.Amount = amount
	e.Category = category
	e.Date = time.Now() // Update the date to current time
}
