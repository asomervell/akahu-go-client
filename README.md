# Akahu API Client

A Go client library for the [Akahu API](https://developers.akahu.nz/).

## Features

- Complete implementation of all Akahu API endpoints
- SQLite database support using GORM
- Configuration using YAML/environment variables
- Context support for cancellation and timeouts
- Error handling and logging
- Type-safe models

## Installation

```bash
go get github.com/andrewsomervell/akahu-client
```

## Configuration

Create a `config.yaml` file in your project root:

```yaml
akahu_api:
  base_url: "https://api.akahu.io/v1"
  app_id: "your-app-id"
  app_secret: "your-app-secret"
  user_token: "your-user-token"

database:
  path: "akahu.db"
```

Or set environment variables:

```bash
export AKAHU_API_BASE_URL="https://api.akahu.io/v1"
export AKAHU_API_APP_ID="your-app-id"
export AKAHU_API_APP_SECRET="your-app-secret"
export AKAHU_API_USER_TOKEN="your-user-token"
export DATABASE_PATH="akahu.db"
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andrewsomervell/akahu-client/internal/client"
	"github.com/andrewsomervell/akahu-client/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create a new client
	akahuClient := client.New(cfg)

	// Create context
	ctx := context.Background()

	// Get current user
	user, err := akahuClient.GetCurrentUser(ctx)
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}
	fmt.Printf("Current user: %+v\n", user)

	// Get accounts
	accounts, err := akahuClient.GetAccounts(ctx)
	if err != nil {
		log.Fatalf("Failed to get accounts: %v", err)
	}
	fmt.Printf("Found %d accounts\n", len(accounts))
}
```

## Available Methods

### Accounts
- `GetAccounts(ctx context.Context) ([]models.Account, error)`
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
- `GetPendingTransactionsByAccount(ctx context.Context, accountID string) ([]models.Transaction, error)`
- `GetTransactionsByIDs(ctx context.Context, ids []string) ([]models.Transaction, error)`

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