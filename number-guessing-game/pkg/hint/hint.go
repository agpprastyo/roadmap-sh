package hint

import (
	"fmt"
)

func ProvideHint(attempts int, guess int, target int) string {
	if attempts > 2 {
		difference := target - guess
		if difference < 0 {
			difference = -difference
		}
		if difference <= 10 {
			return fmt.Sprintf("Hint: You are very close! The difference is within 10.")
		} else if difference <= 20 {
			return fmt.Sprintf("Hint: Getting warmer! The difference is within 20.")
		}
	}
	return ""
}
