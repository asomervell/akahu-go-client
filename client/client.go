// Package client provides a Go client for the Akahu API
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	defaultBaseURL = "https://api.akahu.io/v1"
	genieBaseURL   = "https://api.genie.akahu.io/v1"
)

// Client represents an Akahu API client
type Client struct {
	baseURL    string
	appID      string
	appSecret  string
	genieToken string
	httpClient *http.Client
}

// New creates a new Akahu API client
func New() (*Client, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Get required environment variables
	appID := os.Getenv("AKAHU_APP_ID")
	if appID == "" {
		return nil, fmt.Errorf("AKAHU_APP_ID environment variable is required")
	}

	appSecret := os.Getenv("AKAHU_APP_SECRET")
	if appSecret == "" {
		return nil, fmt.Errorf("AKAHU_APP_SECRET environment variable is required")
	}

	genieToken := os.Getenv("AKAHU_GENIE_TOKEN")
	if genieToken == "" {
		return nil, fmt.Errorf("AKAHU_GENIE_TOKEN environment variable is required")
	}

	// Get optional environment variables with defaults
	baseURL := os.Getenv("AKAHU_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &Client{
		baseURL:    baseURL,
		appID:      appID,
		appSecret:  appSecret,
		genieToken: genieToken,
		httpClient: &http.Client{},
	}, nil
}

// Helper function to make authenticated requests
func (c *Client) makeRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-App-Token", c.appSecret)
	req.Header.Set("Authorization", "Bearer "+c.appID)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

// Helper function to handle API responses
func handleResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}

// GetAccounts retrieves all accounts for the authenticated user
func (c *Client) GetAccounts(ctx context.Context) ([]Account, error) {
	resp, err := c.makeRequest(ctx, "GET", "/accounts", nil)
	if err != nil {
		return nil, err
	}

	var response accountsResponse
	if err := handleResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// GetAccount retrieves a specific account by ID
func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	resp, err := c.makeRequest(ctx, "GET", "/accounts/"+id, nil)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := handleResponse(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// GetTransactions retrieves all transactions
func (c *Client) GetTransactions(ctx context.Context) ([]Transaction, error) {
	resp, err := c.makeRequest(ctx, "GET", "/transactions", nil)
	if err != nil {
		return nil, err
	}

	var response transactionsResponse
	if err := handleResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// GetTransaction retrieves a specific transaction by ID
func (c *Client) GetTransaction(ctx context.Context, id string) (*Transaction, error) {
	resp, err := c.makeRequest(ctx, "GET", "/transactions/"+id, nil)
	if err != nil {
		return nil, err
	}

	var transaction Transaction
	if err := handleResponse(resp, &transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// GetTransactionsByAccount retrieves all transactions for a specific account
func (c *Client) GetTransactionsByAccount(ctx context.Context, accountID string) ([]Transaction, error) {
	resp, err := c.makeRequest(ctx, "GET", "/accounts/"+accountID+"/transactions", nil)
	if err != nil {
		return nil, err
	}

	var response transactionsResponse
	if err := handleResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// GetTransactionsByIDs retrieves transactions by their IDs
func (c *Client) GetTransactionsByIDs(ctx context.Context, ids []string) ([]Transaction, error) {
	payload := map[string][]string{
		"ids": ids,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := c.makeRequest(ctx, "POST", "/transactions/get", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	var response transactionsResponse
	if err := handleResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// EnrichTransaction enriches a single transaction with merchant and category data
func (c *Client) EnrichTransaction(ctx context.Context, tx *Transaction) (*GenieSearchResponse, error) {
	query := GenieSearchQuery{
		Description: tx.Description,
		Amount:      tx.Amount,
	}

	jsonData, err := json.Marshal([]GenieSearchQuery{query})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search query: %w", err)
	}

	resp, err := c.makeGenieRequest(ctx, "POST", "/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var response GenieSearchResponse
	if err := handleResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// EnrichTransactions enriches multiple transactions with merchant and category data
func (c *Client) EnrichTransactions(ctx context.Context, txs []Transaction) (*GenieSearchResponse, error) {
	queries := make([]GenieSearchQuery, len(txs))
	for i, tx := range txs {
		queries[i] = GenieSearchQuery{
			ID:          tx.ID,
			Description: tx.Description,
			Amount:      tx.Amount,
		}
	}

	jsonData, err := json.Marshal(queries)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search queries: %w", err)
	}

	resp, err := c.makeGenieRequest(ctx, "POST", "/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var response GenieSearchResponse
	if err := handleResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Helper function to make Genie API requests
func (c *Client) makeGenieRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, genieBaseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.genieToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}
