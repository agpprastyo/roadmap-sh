package cmd

import (
	"expense-tracker-cli/config"
	"expense-tracker-cli/internals/expense"
	"expense-tracker-cli/utils"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new expense",
	Long:  `Add a new expense with a description, amount, and optional date.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Please provide both description and amount")
			return
		}

		description := args[0]
		amount, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Println("Invalid amount format")
			return
		}

		// Validate the description and amount
		if err := utils.ValidateDescription(description); err != nil {
			fmt.Println(err)
			return
		}
		if err := utils.ValidateAmount(amount); err != nil {
			fmt.Println(err)
			return
		}

		// Optional: Use current date if no date is provided
		date := utils.CurrentDate()
		if len(args) >= 3 {
			date = args[2]
			if err := utils.ValidateDate(date, "2006-01-02"); err != nil {
				fmt.Println("Invalid date format. Expected format is YYYY-MM-DD.")
				return
			}
		}

		// Create a new storage and service instance
		expenseStorage := expense.NewFileStorage(config.ExpensesFilePath)
		expenseService := expense.NewService(expenseStorage)

		// Add the expense
		exp, err := expenseService.AddExpense(description, amount, date)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Expense added successfully (ID: %d)\n", exp.ID)
	},
}
