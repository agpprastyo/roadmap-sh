package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"number-guessing-game/internal/game"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "number-guessing-game",
	Short: "A simple number guessing game",
	Run: func(cmd *cobra.Command, args []string) {
		game.Play()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
