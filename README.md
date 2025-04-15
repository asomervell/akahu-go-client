# Akahu API Client

A Go client library for the [Akahu API](https://developers.akahu.nz/).

## Features

- Complete implementation of all Akahu API endpoints
- Transaction enrichment via Akahu Genie API
- Environment-based configuration
- Context support for cancellation and timeouts
- Error handling and logging
- Type-safe models

## Installation

```bash
go get github.com/asomervell/akahu-go-client
```

## Configuration

The client is configured using environment variables. You can set them directly or use a `.env` file:

```bash
AKAHU_APP_TOKEN="your-app-token"
AKAHU_USER_TOKEN="your-user-token"
AKAHU_GENIE_TOKEN="your-genie-token"  # Optional: for transaction enrichment
AKAHU_BASE_URL="https://api.akahu.io/v1"  # Optional: defaults to production API
```

## Usage

```go
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
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Get accounts
	accounts, err := akahuClient.GetAccounts(ctx)
	if err != nil {
		log.Fatalf("Failed to get accounts: %v", err)
	}

	// Get authenticated user information
	user, err := akahuClient.GetMe(ctx)
	if err != nil {
		log.Fatalf("Failed to get user info: %v", err)
	}
	fmt.Printf("Logged in as: %s (%s)\n", user.Name, user.Email)

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
```

See the [examples](./examples) directory for more usage examples.

## Available Methods

### User
- `GetMe(ctx context.Context) (*User, error)`

### Accounts
- `GetAccounts(ctx context.Context) ([]Account, error)`
- `GetAccount(ctx context.Context, id string) (*Account, error)`

### Transactions
- `GetTransactions(ctx context.Context) ([]Transaction, error)`
- `GetTransaction(ctx context.Context, id string) (*Transaction, error)`
- `GetTransactionsByAccount(ctx context.Context, accountID string) ([]Transaction, error)`
- `GetTransactionsByIDs(ctx context.Context, ids []string) ([]Transaction, error)`
- `EnrichTransaction(ctx context.Context, tx *Transaction) (*GenieSearchResponse, error)`
- `EnrichTransactions(ctx context.Context, txs []Transaction) (*GenieSearchResponse, error)`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin feature/my-new-feature`)
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 