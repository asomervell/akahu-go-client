package main

import (
	"context"
	"fmt"
	"log"

	"github.com/asomervell/akahu-go-client/internal/client"
	"github.com/asomervell/akahu-go-client/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file (optional)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Create a new Akahu client
	akahuClient, err := client.New()
	if err != nil {
		log.Fatalf("Failed to create Akahu client: %v", err)
	}

	// Create a context
	ctx := context.Background()

	// Example 1: Get all categories
	categories, err := akahuClient.GetCategories(ctx)
	if err != nil {
		log.Fatalf("Failed to get categories: %v", err)
	}
	// Type assert the response to []models.Category
	if cats, ok := categories.([]models.Category); ok {
		for _, cat := range cats {
			fmt.Printf("Category: %s (%s)\n", cat.Name, cat.Group)
		}
	}

	// Example 2: Get transactions
	transactions, err := akahuClient.GetTransactions(ctx)
	if err != nil {
		log.Fatalf("Failed to get transactions: %v", err)
	}
	// Type assert the response to []models.Transaction
	if txs, ok := transactions.([]models.Transaction); ok {
		for _, tx := range txs {
			fmt.Printf("Transaction: %s - $%.2f\n", tx.Description, tx.Amount)
		}
	}

	// Example 3: Get a specific account
	accountID := "your-account-id"
	account, err := akahuClient.GetAccount(ctx, accountID)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}
	// Type assert the response to *models.Account
	if acc, ok := account.(*models.Account); ok {
		fmt.Printf("Account: %s - Balance: $%.2f\n", acc.Name, acc.Balance.Current)
	}

	// Example 4: Create a payment
	payment := &models.Payment{
		FromAccountID: "source-account-id",
		ToAccountID:   "destination-account-id",
		Amount:        100.00,
		Description:   "Test payment",
		Reference:     "REF123",
		Particulars:   "PART123",
	}
	result, err := akahuClient.CreatePayment(ctx, payment)
	if err != nil {
		log.Fatalf("Failed to create payment: %v", err)
	}
	fmt.Printf("Payment result: %+v\n", result)
}
