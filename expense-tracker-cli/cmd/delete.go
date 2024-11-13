package cmd

import (
	"expense-tracker-cli/config"
	"expense-tracker-cli/internals/expense"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an expense",
	Long:  `Delete an expense by its ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide the expense ID")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid ID format")
			return
		}

		// Create a new storage and service instance
		expenseStorage := expense.NewFileStorage(config.ExpensesFilePath)
		expenseService := expense.NewService(expenseStorage)

		err = expenseService.DeleteExpense(id)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Expense deleted successfully")
	},
}
