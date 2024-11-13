package budget

import (
	"time"
)

// Budget represents a monthly budget for a specific category or overall.
type Budget struct {
	Month    time.Month `json:"month"`    // The month the budget applies to
	Year     int        `json:"year"`     // The year the budget applies to
	Category string     `json:"category"` // The category the budget is for, or empty for overall budget
	Amount   float64    `json:"amount"`   // The budgeted amount
}

// NewBudget creates a new budget instance.
func NewBudget(month time.Month, year int, category string, amount float64) *Budget {
	return &Budget{
		Month:    month,
		Year:     year,
		Category: category,
		Amount:   amount,
	}
}
