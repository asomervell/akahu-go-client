package client

import (
	"context"

	"akahu-test/internal/models"
)

// AkahuClient defines the interface for interacting with the Akahu API
type AkahuClient interface {
	// Accounts
	GetAccounts(ctx context.Context) (interface{}, error)
	GetAccount(ctx context.Context, id string) (interface{}, error)
	RevokeAccountAccess(ctx context.Context, id string) error

	// Auth
	ExchangeAuthorizationCode(ctx context.Context, code string) (interface{}, error)
	RevokeToken(ctx context.Context, token string) error

	// Categories
	GetCategories(ctx context.Context) (interface{}, error)
	GetCategory(ctx context.Context, id string) (interface{}, error)

	// Connections
	GetConnections(ctx context.Context) (interface{}, error)
	GetConnection(ctx context.Context, id string) (interface{}, error)

	// Data Refresh
	RefreshAllAccounts(ctx context.Context) error
	RefreshAccounts(ctx context.Context, ids []string) error

	// Payments
	GetPayments(ctx context.Context) (interface{}, error)
	CreatePayment(ctx context.Context, payment *models.Payment) (interface{}, error)
	GetPayment(ctx context.Context, id string) (interface{}, error)
	CancelPayment(ctx context.Context, id string) error

	// Transactions
	GetTransactions(ctx context.Context) (interface{}, error)
	GetPendingTransactions(ctx context.Context) (interface{}, error)
	GetTransaction(ctx context.Context, id string) (interface{}, error)
	GetTransactionsByAccount(ctx context.Context, accountID string) (interface{}, error)
	GetTransactionsByIDs(ctx context.Context, ids []string) (interface{}, error)
	EnrichTransaction(ctx context.Context, tx *models.Transaction) (*GenieSearchResponse, error)
	EnrichTransactions(ctx context.Context, txs []models.Transaction) (*GenieSearchResponse, error)

	// Transfers
	GetTransfers(ctx context.Context) (interface{}, error)
	CreateTransfer(ctx context.Context, transfer *models.Transfer) (interface{}, error)
	GetTransfer(ctx context.Context, id string) (interface{}, error)

	// User
	GetCurrentUser(ctx context.Context) (interface{}, error)

	// Webhooks
	GetWebhooks(ctx context.Context) (interface{}, error)
	CreateWebhook(ctx context.Context, webhook *models.Webhook) (interface{}, error)
	GetWebhookPublicKey(ctx context.Context) (interface{}, error)
	DeleteWebhook(ctx context.Context, id string) error
	GetWebhookEvents(ctx context.Context) (interface{}, error)
}
