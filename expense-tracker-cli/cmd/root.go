package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "expense-tracker",
	Short: "Expense Tracker is a CLI tool to manage your expenses",
	Long:  `A simple CLI application to track and manage your personal expenses.`,
}

// Execute runs the root command and all its child commands.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands to the root command
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(updateCmd)
	RootCmd.AddCommand(summaryCmd)
}
