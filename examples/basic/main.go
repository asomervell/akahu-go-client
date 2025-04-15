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
		if err != nil {
			log.Fatalf("Failed to get transactions: %v", err)
		}

		for _, tx := range transactions {
			fmt.Printf("Transaction: %s - $%.2f\n", tx.Description, tx.Amount)

			// Enrich the transaction with merchant and category data
			enriched, err := akahuClient.EnrichTransaction(ctx, &tx)
			if err != nil {
				log.Printf("Failed to enrich transaction: %v", err)
				continue
			}

			if len(enriched.Items) > 0 && len(enriched.Items[0].Results) > 0 {
				result := enriched.Items[0].Results[0]
				fmt.Printf("  Category: %s\n", result.Category.Name)
				if result.Merchant.Name != "" {
					fmt.Printf("  Merchant: %s\n", result.Merchant.Name)
				}
			}
		}
	}
}
