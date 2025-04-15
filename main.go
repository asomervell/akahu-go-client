package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"akahu-test/internal/client"
	"akahu-test/internal/models"
)

var (
	// Command flags
	command = flag.String("cmd", "", `Command to execute. Available commands:
	accounts:list            - List all accounts
	accounts:get            - Get account by ID
	accounts:revoke         - Revoke access to account
	auth:exchange           - Exchange authorization code for token
	auth:revoke            - Revoke token
	categories:list         - List all categories
	categories:get         - Get category by ID
	connections:list        - List all connections
	connections:get        - Get connection by ID
	refresh:all            - Refresh all accounts
	refresh:accounts       - Refresh specific accounts
	payments:list          - List all payments
	payments:create        - Create a payment
	payments:get          - Get payment by ID
	payments:cancel        - Cancel payment
	transactions:list      - List all transactions
	transactions:pending   - List pending transactions
	transactions:get      - Get transaction by ID
	transactions:by-account - Get transactions by account
	transactions:by-ids    - Get transactions by IDs
	transactions:enrich    - Enrich transaction with merchant and category data
	transactions:enrich-batch - Enrich multiple transactions
	transfers:list         - List all transfers
	transfers:create       - Create a transfer
	transfers:get         - Get transfer by ID
	user:me               - Get current user
	webhooks:list         - List all webhooks
	webhooks:create       - Create a webhook
	webhooks:key          - Get webhook public key
	webhooks:delete       - Delete webhook
	webhooks:events       - List webhook events`)

	// Common flags
	id     = flag.String("id", "", "Resource ID")
	format = flag.String("format", "json", "Output format (json or pretty)")

	// Specific command flags
	accountIDs = flag.String("account-ids", "", "Comma-separated list of account IDs")
	code       = flag.String("code", "", "Authorization code")
	token      = flag.String("token", "", "Token to revoke")

	// Payment/Transfer flags
	fromAccount = flag.String("from", "", "Source account ID")
	toAccount   = flag.String("to", "", "Destination account ID")
	amount      = flag.Float64("amount", 0.0, "Amount to transfer/pay")
	description = flag.String("description", "", "Payment/Transfer description")
	reference   = flag.String("reference", "", "Payment reference")
	particulars = flag.String("particulars", "", "Payment particulars")

	// Webhook flags
	webhookURL = flag.String("url", "", "Webhook URL")
	events     = flag.String("events", "", "Comma-separated list of webhook events")
)

func main() {
	flag.Parse()

	if *command == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Create client
	akahuClient, err := client.New()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Execute command
	ctx := context.Background()
	var result interface{}
	var cmdErr error

	switch *command {
	// Accounts
	case "accounts:list":
		result, cmdErr = akahuClient.GetAccounts(ctx)
	case "accounts:get":
		requireFlag("id", *id)
		result, cmdErr = akahuClient.GetAccount(ctx, *id)
	case "accounts:revoke":
		requireFlag("id", *id)
		cmdErr = akahuClient.RevokeAccountAccess(ctx, *id)

	// Auth
	case "auth:exchange":
		requireFlag("code", *code)
		result, cmdErr = akahuClient.ExchangeAuthorizationCode(ctx, *code)
	case "auth:revoke":
		requireFlag("token", *token)
		cmdErr = akahuClient.RevokeToken(ctx, *token)

	// Categories
	case "categories:list":
		result, cmdErr = akahuClient.GetCategories(ctx)
	case "categories:get":
		requireFlag("id", *id)
		result, cmdErr = akahuClient.GetCategory(ctx, *id)

	// Connections
	case "connections:list":
		result, cmdErr = akahuClient.GetConnections(ctx)
	case "connections:get":
		requireFlag("id", *id)
		result, cmdErr = akahuClient.GetConnection(ctx, *id)

	// Data Refresh
	case "refresh:all":
		cmdErr = akahuClient.RefreshAllAccounts(ctx)
	case "refresh:accounts":
		requireFlag("account-ids", *accountIDs)
		ids := strings.Split(*accountIDs, ",")
		cmdErr = akahuClient.RefreshAccounts(ctx, ids)

	// Payments
	case "payments:list":
		result, cmdErr = akahuClient.GetPayments(ctx)
	case "payments:create":
		payment := &models.Payment{
			FromAccountID: *fromAccount,
			ToAccountID:   *toAccount,
			Amount:        *amount,
			Description:   *description,
			Reference:     *reference,
			Particulars:   *particulars,
		}
		result, cmdErr = akahuClient.CreatePayment(ctx, payment)
	case "payments:get":
		requireFlag("id", *id)
		result, cmdErr = akahuClient.GetPayment(ctx, *id)
	case "payments:cancel":
		requireFlag("id", *id)
		cmdErr = akahuClient.CancelPayment(ctx, *id)

	// Transactions
	case "transactions:list":
		result, cmdErr = akahuClient.GetTransactions(ctx)
	case "transactions:pending":
		result, cmdErr = akahuClient.GetPendingTransactions(ctx)
	case "transactions:get":
		requireFlag("id", *id)
		result, cmdErr = akahuClient.GetTransaction(ctx, *id)
	case "transactions:by-account":
		requireFlag("id", *id)
		result, cmdErr = akahuClient.GetTransactionsByAccount(ctx, *id)
	case "transactions:by-ids":
		requireFlag("id", *id)
		ids := strings.Split(*id, ",")
		result, cmdErr = akahuClient.GetTransactionsByIDs(ctx, ids)
	case "transactions:enrich":
		requireFlag("id", *id)
		tx, err := akahuClient.GetTransaction(ctx, *id)
		if err != nil {
			log.Fatalf("Failed to get transaction: %v", err)
		}
		transaction, ok := tx.(*models.Transaction)
		if !ok {
			log.Fatalf("Failed to convert transaction")
		}
		result, cmdErr = akahuClient.EnrichTransaction(ctx, transaction)
	case "transactions:enrich-batch":
		requireFlag("id", *id)
		ids := strings.Split(*id, ",")
		txs, err := akahuClient.GetTransactionsByIDs(ctx, ids)
		if err != nil {
			log.Fatalf("Failed to get transactions: %v", err)
		}
		transactions, ok := txs.([]models.Transaction)
		if !ok {
			log.Fatalf("Failed to convert transactions")
		}
		result, cmdErr = akahuClient.EnrichTransactions(ctx, transactions)

	// Transfers
	case "transfers:list":
		result, cmdErr = akahuClient.GetTransfers(ctx)
	case "transfers:create":
		transfer := &models.Transfer{
			FromAccountID: *fromAccount,
			ToAccountID:   *toAccount,
			Amount:        *amount,
			Description:   *description,
		}
		result, cmdErr = akahuClient.CreateTransfer(ctx, transfer)
	case "transfers:get":
		requireFlag("id", *id)
		result, cmdErr = akahuClient.GetTransfer(ctx, *id)

	// User
	case "user:me":
		result, cmdErr = akahuClient.GetCurrentUser(ctx)

	// Webhooks
	case "webhooks:list":
		result, cmdErr = akahuClient.GetWebhooks(ctx)
	case "webhooks:create":
		webhook := &models.Webhook{
			URL:    *webhookURL,
			Events: strings.Split(*events, ","),
		}
		result, cmdErr = akahuClient.CreateWebhook(ctx, webhook)
	case "webhooks:key":
		result, cmdErr = akahuClient.GetWebhookPublicKey(ctx)
	case "webhooks:delete":
		requireFlag("id", *id)
		cmdErr = akahuClient.DeleteWebhook(ctx, *id)
	case "webhooks:events":
		result, cmdErr = akahuClient.GetWebhookEvents(ctx)

	default:
		log.Fatalf("Unknown command: %s", *command)
	}

	if cmdErr != nil {
		log.Fatalf("Command failed: %v", cmdErr)
	}

	if result != nil {
		if *format == "pretty" {
			prettyPrint(result)
		} else {
			jsonPrint(result)
		}
	}
}

func requireFlag(name, value string) {
	if value == "" {
		log.Fatalf("-%s flag is required", name)
	}
}

func prettyPrint(v interface{}) {
	switch val := v.(type) {
	case string:
		fmt.Println(val)
	default:
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			log.Fatalf("Failed to format output: %v", err)
		}
		fmt.Println(string(b))
	}
}

func jsonPrint(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Failed to format output: %v", err)
	}
	fmt.Println(string(b))
}
