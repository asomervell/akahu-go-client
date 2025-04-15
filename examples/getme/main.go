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

	// Get user information
	user, err := akahuClient.GetMe(ctx)
	if err != nil {
		log.Fatalf("Failed to get user info: %v", err)
	}

	// Print user details
	fmt.Printf("User ID: %s\n", user.ID)
	fmt.Printf("Email: %s\n", user.Email)
	if user.PreferredName != "" {
		fmt.Printf("Preferred Name: %s\n", user.PreferredName)
	}
	fmt.Printf("Access Granted: %s\n", user.AccessGrantedAt.Format("2006-01-02 15:04:05"))
}
