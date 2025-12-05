package database

import (
	"database/sql"
	"fmt"
	"time"

	"foodie/backend/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Config holds database configuration.
type Config struct {
	Driver          string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewConnection creates a new database connection with proper configuration.
// This function handles connection pooling and settings.
func NewConnection(cfg Config) (*sql.DB, error) {
	if cfg.Driver == "" {
		return nil, fmt.Errorf("database driver is required")
	}
	if cfg.DSN == "" {
		return nil, fmt.Errorf("database DSN is required")
	}

	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// NewConnectionFromEnv creates a database connection from environment variables.
func NewConnectionFromEnv() (*sql.DB, error) {
	driver := config.Get("SQL_DRIVER", "postgres")
	dsn := config.Get("SQL_DSN", "")

	if dsn == "" {
		return nil, fmt.Errorf("SQL_DSN environment variable is required")
	}

	cfg := Config{
		Driver:          driver,
		DSN:             dsn,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	return NewConnection(cfg)
}
