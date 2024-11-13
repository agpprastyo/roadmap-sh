package game

import (
	"fmt"
	"number-guessing-game/internal/utils"
	"number-guessing-game/pkg/hint"
	"time"
)

func Play() {
	fmt.Println("Welcome to the Number Guessing Game!")
	fmt.Println("I'm thinking of a number between 1 and 100.")
	fmt.Println("Please select the difficulty level:")
	fmt.Println("1. Easy (10 chances)")
	fmt.Println("2. Medium (5 chances)")
	fmt.Println("3. Hard (3 chances)")

	var choice int
	fmt.Print("Enter your choice: ")
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("Invalid input, defaulting to Medium difficulty.")
		choice = 2
	}

	level := SelectDifficulty(choice)
	fmt.Printf("Great! You have selected the %s difficulty level.\n", level.Name)

	target := utils.GenerateRandomNumber(1, 100)
	startTime := time.Now()

	for attempts := 1; attempts <= level.Chances; attempts++ {
		var guess int
		fmt.Printf("Enter your guess (%d chances left): ", level.Chances-attempts+1)
		_, err := fmt.Scan(&guess)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid number.")
			attempts--
			continue
		}

		if guess == target {
			fmt.Printf("Congratulations! You guessed the correct number in %d attempts.\n", attempts)
			break
		} else if guess < target {
			fmt.Println("Incorrect! The number is greater than", guess)
		} else {
			fmt.Println("Incorrect! The number is less than", guess)
		}

		if attempts == level.Chances {
			fmt.Println("Sorry, you've run out of chances. The correct number was", target)
		} else if hintNeeded := hint.ProvideHint(attempts, guess, target); hintNeeded != "" {
			fmt.Println(hintNeeded)
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Println(utils.FormatElapsedTime(elapsedTime))
}
