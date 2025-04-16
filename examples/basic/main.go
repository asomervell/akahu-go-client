package main

import (
	"context"
	"fmt"
	"log"

	"github.com/asomervell/akahu-go-client/client"
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

	// Get all accounts
	accounts, err := akahuClient.GetAccounts(ctx)
	if err != nil {
		log.Fatalf("Failed to get accounts: %v", err)
	}

	for _, acc := range accounts {
		fmt.Printf("Account: %s - Balance: $%.2f %s\n",
			acc.Name,
			acc.Balance.Current,
			acc.Balance.Currency)
	}

	// Get transactions for the first account
	if len(accounts) > 0 {

		transactions, err := akahuClient.GetTransactionsByAccount(ctx, accounts[0].ID)

		fmt.Printf("%v transactions\n", len(transactions))

		if err != nil {
			log.Fatalf("Failed to get transactions: %v", err)
		}

		for _, tx := range transactions {
			fmt.Printf("\nTransaction: %s - $%.2f\n", tx.Description, tx.Amount)
			fmt.Printf("Meta: %+v\n", tx.Meta)
			if tx.Meta.Logo != "" {
				fmt.Printf("Meta Logo: %s\n", tx.Meta.Logo)
			}
			if tx.Merchant.ID != "" {
				fmt.Printf("Merchant: ID=%s, Name=%s, Website=%s, Logo=%s, NZBN=%s\n",
					tx.Merchant.ID,
					tx.Merchant.Name,
					tx.Merchant.Website,
					tx.Merchant.Logo,
					tx.Merchant.NZBN)
			}
			if tx.Category.ID != "" {
				fmt.Printf("Category: ID=%s, Name=%s, Group=%s\n",
					tx.Category.ID,
					tx.Category.Name,
					tx.Category.Group)
				if len(tx.Category.Groups) > 0 {
					fmt.Printf("Category Groups: %+v\n", tx.Category.Groups)
				}
			}
		}
	}
}
