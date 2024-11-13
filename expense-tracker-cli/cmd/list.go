package cmd

import (
	"expense-tracker-cli/config"
	"expense-tracker-cli/internals/expense"
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all expenses",
	Long:  `List all expenses with their IDs, dates, descriptions, and amounts.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Create a new storage and service instance
		expenseStorage := expense.NewFileStorage(config.ExpensesFilePath)
		expenseService := expense.NewService(expenseStorage)

		// List expenses
		expenses, err := expenseService.ListExpenses()
		if err != nil {
			fmt.Println("Error listing expenses:", err)
			return
		}

		fmt.Println("ID  Date       Description  Amount")
		for _, exp := range expenses {
			fmt.Printf("%d   %s  %s  $%.2f\n", exp.ID, exp.Date.Format("2006-01-02"), exp.Description, exp.Amount)
		}
	},
}
