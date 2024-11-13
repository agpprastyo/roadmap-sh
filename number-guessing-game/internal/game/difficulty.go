package game

import "fmt"

type Difficulty struct {
	Name    string
	Chances int
}

func SelectDifficulty(choice int) Difficulty {
	switch choice {
	case 1:
		return Difficulty{Name: "Easy", Chances: 10}
	case 2:
		return Difficulty{Name: "Medium", Chances: 5}
	case 3:
		return Difficulty{Name: "Hard", Chances: 3}
	default:
		fmt.Println("Invalid choice, defaulting to Medium difficulty.")
		return Difficulty{Name: "Medium", Chances: 5}
	}
}
