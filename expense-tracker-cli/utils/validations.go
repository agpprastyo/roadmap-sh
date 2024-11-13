package utils

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateAmount checks if the given amount is valid (greater than zero).
func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return nil
}

// ValidateDescription checks if the given description is valid (non-empty and within length limits).
func ValidateDescription(description string) error {
	if len(strings.TrimSpace(description)) == 0 {
		return errors.New("description cannot be empty")
	}
	if len(description) > 255 {
		return errors.New("description cannot exceed 255 characters")
	}
	return nil
}

// ValidateDate checks if the given date string matches the specified layout.
func ValidateDate(dateStr, layout string) error {
	_, err := ParseDate(dateStr, layout)
	if err != nil {
		return errors.New("invalid date format")
	}
	return nil
}

// ValidateCategory checks if the category name matches a specific pattern (e.g., letters and spaces only).
func ValidateCategory(category string) error {
	if len(category) > 50 {
		return errors.New("category name cannot exceed 50 characters")
	}

	re := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !re.MatchString(category) {
		return errors.New("category name can only contain letters and spaces")
	}

	return nil
}
