package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func FormatElapsedTime(elapsedTime time.Duration) string {
	totalSeconds := int(elapsedTime.Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("Time taken: %d minutes %d seconds", minutes, seconds)
}
