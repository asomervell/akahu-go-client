package client

import "time"

// Account represents a bank account
type Account struct {
	ID               string                 `json:"_id"`
	Credentials      string                 `json:"_credentials"`
	Connection       map[string]interface{} `json:"connection"`
	Name             string                 `json:"name"`
	FormattedAccount string                 `json:"formatted_account"`
	Status           string                 `json:"status"`
	Type             string                 `json:"type"`
	Attributes       []string               `json:"attributes"`
	Balance          Balance                `json:"balance"`
	Meta             map[string]interface{} `json:"meta"`
	Refreshed        map[string]string      `json:"refreshed"`
}

// Balance represents an account balance
type Balance struct {
	Currency  string  `json:"currency"`
	Current   float64 `json:"current"`
	Available float64 `json:"available"`
	Limit     float64 `json:"limit"`
	Overdrawn bool    `json:"overdrawn"`
}

// Transaction represents a bank transaction
type Transaction struct {
	ID          string    `json:"_id"`
	AccountID   string    `json:"account"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Balance     float64   `json:"balance"`
	Category    Category  `json:"category"`
	Merchant    Merchant  `json:"merchant"`
	Meta        Meta      `json:"meta"`
	Status      string    `json:"status"`
	IsPending   bool      `json:"pending"`
}

// Category represents a transaction category
type Category struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

// Merchant represents a transaction merchant
type Merchant struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Logo        string `json:"logo,omitempty"`
	Description string `json:"description,omitempty"`
}

// Meta represents transaction metadata
type Meta struct {
	Particulars string `json:"particulars,omitempty"`
	Code        string `json:"code,omitempty"`
	Reference   string `json:"reference,omitempty"`
	OtherParty  string `json:"other_party,omitempty"`
}

// GenieSearchQuery represents a query to the Genie API
type GenieSearchQuery struct {
	ID          string  `json:"id,omitempty"`
	Description string  `json:"description"`
	Connection  string  `json:"_connection,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	Direction   string  `json:"direction,omitempty"`
}

// GenieCategoryGroup represents a category group in the Genie API
type GenieCategoryGroup struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

// GenieCategory represents a category in the Genie API
type GenieCategory struct {
	ID     string                         `json:"_id"`
	Name   string                         `json:"name"`
	Groups map[string]*GenieCategoryGroup `json:"groups"`
}

// GenieMerchant represents a merchant in the Genie API
type GenieMerchant struct {
	ID      string `json:"_id"`
	Name    string `json:"name"`
	Logo    string `json:"logo,omitempty"`
	Website string `json:"website,omitempty"`
}

// GenieSearchResult represents a search result from the Genie API
type GenieSearchResult struct {
	Confidence float64       `json:"confidence"`
	Category   GenieCategory `json:"category"`
	Merchant   GenieMerchant `json:"merchant,omitempty"`
}

// GenieSearchResponseItem represents a single item in the Genie API response
type GenieSearchResponseItem struct {
	ID      string              `json:"id,omitempty"`
	Query   string              `json:"query"`
	Results []GenieSearchResult `json:"results"`
}

// GenieSearchResponse represents the response from the Genie API
type GenieSearchResponse struct {
	Success bool                      `json:"success"`
	Items   []GenieSearchResponseItem `json:"items"`
}

// User represents an authenticated user
type User struct {
	ID            string    `json:"_id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	PreferredName string    `json:"preferred_name,omitempty"`
	MobileNumber  string    `json:"mobile_number,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// API response types
type accountsResponse struct {
	Items []Account `json:"items"`
}

type transactionsResponse struct {
	Items []Transaction `json:"items"`
}
