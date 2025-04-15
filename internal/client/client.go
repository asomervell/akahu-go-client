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

	"github.com/asomervell/akahu-go-client/internal/models"

	"github.com/joho/godotenv"
)

const (
	defaultBaseURL = "https://api.akahu.io/v1"
	genieBaseURL   = "https://api.genie.akahu.io/v1"
)

type Client struct {
	baseURL    string
	appID      string
	appSecret  string
	genieToken string
	httpClient *http.Client
}

func New() (AkahuClient, error) {
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

type accountsResponse struct {
	Items []models.Account `json:"items"`
}

type categoriesResponse struct {
	Items []models.Category `json:"items"`
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

func (c *Client) GetAccounts(ctx context.Context) (interface{}, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/accounts", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	req.Header.Set("Authorization", "Bearer "+c.appSecret)
	req.Header.Set("X-Akahu-Id", c.appID)

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Debug: Print raw response
	fmt.Printf("Raw response: %s\n", string(body))

	// Parse response
	var response accountsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response.Items, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (interface{}, error) {
	resp, err := c.makeRequest(ctx, "GET", "/accounts/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var account models.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &account, nil
}

func (c *Client) RevokeAccountAccess(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) ExchangeAuthorizationCode(ctx context.Context, code string) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) RevokeToken(ctx context.Context, token string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) GetCategories(ctx context.Context) (interface{}, error) {
	resp, err := c.makeRequest(ctx, "GET", "/categories", nil)
	if err != nil {
		return nil, err
	}

	var response categoriesResponse
	if err := handleResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

func (c *Client) GetCategory(ctx context.Context, id string) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) GetConnections(ctx context.Context) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) GetConnection(ctx context.Context, id string) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) RefreshAllAccounts(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) RefreshAccounts(ctx context.Context, ids []string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) GetPayments(ctx context.Context) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) CreatePayment(ctx context.Context, payment *models.Payment) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) GetPayment(ctx context.Context, id string) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) CancelPayment(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}

type apiResponse struct {
	Success bool                 `json:"success"`
	Items   []models.Transaction `json:"items"`
}

func (c *Client) GetTransactions(ctx context.Context) (interface{}, error) {
	resp, err := c.makeRequest(ctx, "GET", "/transactions", nil)
	if err != nil {
		return nil, err
	}

	var response apiResponse
	if err := handleResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

func (c *Client) GetPendingTransactions(ctx context.Context) (interface{}, error) {
	resp, err := c.makeRequest(ctx, "GET", "/transactions/pending", nil)
	if err != nil {
		return nil, err
	}

	var response apiResponse
	if err := handleResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

func (c *Client) GetTransaction(ctx context.Context, id string) (interface{}, error) {
	resp, err := c.makeRequest(ctx, "GET", "/transactions/"+id, nil)
	if err != nil {
		return nil, err
	}

	var transaction models.Transaction
	if err := handleResponse(resp, &transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (c *Client) GetTransactionsByAccount(ctx context.Context, accountID string) (interface{}, error) {
	resp, err := c.makeRequest(ctx, "GET", "/accounts/"+accountID+"/transactions", nil)
	if err != nil {
		return nil, err
	}

	var response apiResponse
	if err := handleResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

func (c *Client) GetTransactionsByIDs(ctx context.Context, ids []string) (interface{}, error) {
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

	var response apiResponse
	if err := handleResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

func (c *Client) GetTransfers(ctx context.Context) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) CreateTransfer(ctx context.Context, transfer *models.Transfer) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) GetTransfer(ctx context.Context, id string) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) GetCurrentUser(ctx context.Context) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) GetWebhooks(ctx context.Context) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) CreateWebhook(ctx context.Context, webhook *models.Webhook) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) GetWebhookPublicKey(ctx context.Context) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
}

func (c *Client) DeleteWebhook(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) GetWebhookEvents(ctx context.Context) (interface{}, error) {
	return map[string]string{"status": "not implemented"}, nil
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

// Genie API types
type GenieSearchQuery struct {
	ID          string  `json:"id,omitempty"`
	Description string  `json:"description"`
	Connection  string  `json:"_connection,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	Direction   string  `json:"direction,omitempty"`
}

type GenieCategoryGroup struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type GenieCategory struct {
	ID     string                         `json:"_id"`
	Name   string                         `json:"name"`
	Groups map[string]*GenieCategoryGroup `json:"groups"`
}

type GenieMerchant struct {
	ID      string `json:"_id"`
	Name    string `json:"name"`
	Logo    string `json:"logo,omitempty"`
	Website string `json:"website,omitempty"`
}

type GenieSearchResult struct {
	Confidence float64       `json:"confidence"`
	Category   GenieCategory `json:"category"`
	Merchant   GenieMerchant `json:"merchant,omitempty"`
}

type GenieSearchResponseItem struct {
	ID      string              `json:"id,omitempty"`
	Query   string              `json:"query"`
	Results []GenieSearchResult `json:"results"`
}

type GenieSearchResponse struct {
	Success bool                      `json:"success"`
	Items   []GenieSearchResponseItem `json:"items"`
}

// EnrichTransaction enriches a single transaction with merchant and category data
func (c *Client) EnrichTransaction(ctx context.Context, tx *models.Transaction) (*GenieSearchResponse, error) {
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
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response GenieSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// EnrichTransactions enriches multiple transactions with merchant and category data
func (c *Client) EnrichTransactions(ctx context.Context, txs []models.Transaction) (*GenieSearchResponse, error) {
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
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response GenieSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
