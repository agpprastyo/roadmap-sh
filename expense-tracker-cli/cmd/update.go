package cmd

import (
	"expense-tracker-cli/config"
	"expense-tracker-cli/internals/expense"
	"expense-tracker-cli/utils"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing expense",
	Long:  `Update an existing expense by providing its ID, new description, new amount, and an optional date.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			fmt.Println("Please provide the ID, new description, and new amount")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid ID format")
			return
		}

		description := args[1]
		amount, err := strconv.ParseFloat(args[2], 64)
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
		if len(args) >= 4 {
			date = args[3]
			if err := utils.ValidateDate(date, "2006-01-02"); err != nil {
				fmt.Println("Invalid date format. Expected format is YYYY-MM-DD.")
				return
			}
		}

		// Create a new storage and service instance
		expenseStorage := expense.NewFileStorage(config.ExpensesFilePath)
		expenseService := expense.NewService(expenseStorage)

		// Update the expense
		err = expenseService.UpdateExpense(id, description, amount, date)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Expense updated successfully")
	},
}
