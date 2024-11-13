package utils

import (
	"fmt"
	"time"
)

// FormatDate formats a time.Time object into a string with the given layout.
func FormatDate(t time.Time, layout string) string {
	return t.Format(layout)
}

// ParseDate parses a date string into a time.Time object according to the given layout.
func ParseDate(dateStr, layout string) (time.Time, error) {
	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date: %w", err)
	}
	return parsedTime, nil
}

// CurrentDate returns the current date formatted as YYYY-MM-DD.
func CurrentDate() string {
	return time.Now().Format("2006-01-02")
}
