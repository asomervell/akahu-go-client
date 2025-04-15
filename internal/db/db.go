package db

import (
	"fmt"

	"github.com/asomervell/akahu-go-client/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB represents the database connection
type DB struct {
	*gorm.DB
}

// New creates a new database connection
func New(path string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.Account{},
		&models.Transaction{},
		&models.Category{},
		&models.Connection{},
		&models.Payment{},
		&models.Transfer{},
		&models.User{},
		&models.Webhook{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &DB{DB: db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database: %w", err)
	}
	return sqlDB.Close()
}
