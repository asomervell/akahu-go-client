package models

import (
	"time"

	"gorm.io/gorm"
)

// Base model for common fields
type Base struct {
	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

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
	Balance          struct {
		Currency  string  `json:"currency"`
		Current   float64 `json:"current"`
		Available float64 `json:"available"`
		Limit     float64 `json:"limit"`
		Overdrawn bool    `json:"overdrawn"`
	} `json:"balance"`
	Meta      map[string]interface{} `json:"meta"`
	Refreshed map[string]string      `json:"refreshed"`
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
	Category    struct {
		ID          string `json:"_id"`
		Name        string `json:"name"`
		Group       string `json:"group"`
		Description string `json:"description"`
	} `json:"category"`
	Merchant struct {
		ID          string `json:"_id"`
		Name        string `json:"name"`
		Logo        string `json:"logo,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"merchant"`
	Meta struct {
		Particulars string `json:"particulars,omitempty"`
		Code        string `json:"code,omitempty"`
		Reference   string `json:"reference,omitempty"`
		OtherParty  string `json:"other_party,omitempty"`
	} `json:"meta"`
	Status    string `json:"status"`
	IsPending bool   `json:"pending"`
}

// Category represents a transaction category
type Category struct {
	Base
	Name        string `json:"name"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

// Connection represents a bank connection
type Connection struct {
	Base
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// Payment represents a payment transaction
type Payment struct {
	FromAccountID string  `json:"from_account_id"`
	ToAccountID   string  `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	Reference     string  `json:"reference"`
	Particulars   string  `json:"particulars"`
}

// Transfer represents an internal transfer
type Transfer struct {
	FromAccountID string  `json:"from_account_id"`
	ToAccountID   string  `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
}

// User represents the current authenticated user
type User struct {
	Base
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Webhook represents a webhook subscription
type Webhook struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}
