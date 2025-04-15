# Akahu API Client

A Go client library for the [Akahu API](https://developers.akahu.nz/).

## Features

- Complete implementation of all Akahu API endpoints
- Transaction enrichment via Akahu Genie API
- Environment-based configuration
- Context support for cancellation and timeouts
- Error handling and logging
- Type-safe models
- CLI tool for interacting with the API

## Installation

```bash
go get github.com/asomervell/akahu-go-client
```

## Configuration

The client is configured using environment variables. You can set them directly or use a `.env` file:

```bash
AKAHU_APP_ID="your-app-id"
AKAHU_APP_SECRET="your-app-secret"
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

	"github.com/asomervell/akahu-go-client/internal/client"
)

func main() {
	// Create a new client
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

	// Type assert the response
	if accs, ok := accounts.([]models.Account); ok {
		for _, acc := range accs {
			fmt.Printf("Account: %s, Balance: $%.2f %s\n", 
				acc.Name, 
				acc.Balance.Current,
				acc.Balance.Currency)
		}
	}

	// Enrich a transaction with merchant and category data
	if tx, ok := transaction.(*models.Transaction); ok {
		enriched, err := akahuClient.EnrichTransaction(ctx, tx)
		if err != nil {
			log.Fatalf("Failed to enrich transaction: %v", err)
		}
		fmt.Printf("Enriched transaction: %+v\n", enriched)
	}
}
```

## CLI Tool

The client includes a CLI tool for interacting with the API. Run with `-h` to see available commands:

```bash
go run main.go -h
```

Example commands:

```bash
# List all accounts
go run main.go -cmd accounts:list

# Get transactions for an account
go run main.go -cmd transactions:by-account -id ACCOUNT_ID

# Enrich a transaction
go run main.go -cmd transactions:enrich -id TRANSACTION_ID
```

## Available Methods

### Accounts
- `GetAccounts(ctx context.Context) (interface{}, error)`
- `GetAccount(ctx context.Context, id string) (*models.Account, error)`
- `RevokeAccountAccess(ctx context.Context, id string) error`

### Auth
- `ExchangeAuthorizationCode(ctx context.Context, code string) (string, error)`
- `RevokeToken(ctx context.Context, token string) error`

### Categories
- `GetCategories(ctx context.Context) ([]models.Category, error)`
- `GetCategory(ctx context.Context, id string) (*models.Category, error)`

### Connections
- `GetConnections(ctx context.Context) ([]models.Connection, error)`
- `GetConnection(ctx context.Context, id string) (*models.Connection, error)`

### Data Refresh
- `RefreshAllAccounts(ctx context.Context) error`
- `RefreshAccounts(ctx context.Context, accountIDs []string) error`

### Payments
- `GetPayments(ctx context.Context) ([]models.Payment, error)`
- `CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error)`
- `GetPayment(ctx context.Context, id string) (*models.Payment, error)`
- `CancelPayment(ctx context.Context, id string) error`

### Transactions
- `GetTransactions(ctx context.Context) ([]models.Transaction, error)`
- `GetPendingTransactions(ctx context.Context) ([]models.Transaction, error)`
- `GetTransaction(ctx context.Context, id string) (*models.Transaction, error)`
- `GetTransactionsByAccount(ctx context.Context, accountID string) ([]models.Transaction, error)`
- `GetTransactionsByIDs(ctx context.Context, ids []string) ([]models.Transaction, error)`
- `EnrichTransaction(ctx context.Context, tx *models.Transaction) (*GenieSearchResponse, error)`
- `EnrichTransactions(ctx context.Context, txs []models.Transaction) (*GenieSearchResponse, error)`

### Transfers
- `GetTransfers(ctx context.Context) ([]models.Transfer, error)`
- `CreateTransfer(ctx context.Context, transfer *models.Transfer) (*models.Transfer, error)`
- `GetTransfer(ctx context.Context, id string) (*models.Transfer, error)`

### User
- `GetCurrentUser(ctx context.Context) (*models.User, error)`

### Webhooks
- `GetWebhooks(ctx context.Context) ([]models.Webhook, error)`
- `CreateWebhook(ctx context.Context, webhook *models.Webhook) (*models.Webhook, error)`
- `GetWebhookPublicKey(ctx context.Context) (string, error)`
- `DeleteWebhook(ctx context.Context, id string) error`
- `GetWebhookEvents(ctx context.Context) ([]string, error)`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin feature/my-new-feature`)
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 