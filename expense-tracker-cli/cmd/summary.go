package cmd

import (
	"expense-tracker-cli/config"
	"expense-tracker-cli/internals/expense"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show a summary of expenses",
	Long:  `Show a summary of total expenses, optionally filtered by month.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new storage and service instance
		expenseStorage := expense.NewFileStorage(config.ExpensesFilePath)
		expenseService := expense.NewService(expenseStorage)

		if len(args) == 0 {
			// Show total summary
			total, err := expenseService.TotalExpenses()
			if err != nil {
				fmt.Println("Error calculating total expenses:", err)
				return
			}
			fmt.Printf("Total expenses: $%.2f\n", total)
		} else if len(args) == 1 {
			// Show summary for a specific month
			month, err := strconv.Atoi(args[0])
			if err != nil || month < 1 || month > 12 {
				fmt.Println("Invalid month. Please provide a month as a number (1-12).")
				return
			}

			currentYear := time.Now().Year()
			total, err := expenseService.TotalExpensesByMonth(time.Month(month), currentYear)
			if err != nil {
				fmt.Println("Error calculating expenses for the month:", err)
				return
			}
			fmt.Printf("Total expenses for %s: $%.2f\n", time.Month(month).String(), total)
		} else {
			fmt.Println("Invalid usage. Use 'summary' or 'summary [month]'.")
		}
	},
}
