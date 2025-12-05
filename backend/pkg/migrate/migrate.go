package migrate

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs all pending migrations from the migrations directory.
// It supports both PostgreSQL and MySQL.
func RunMigrations(db *sql.DB, driver string, migrationsPath string) error {
	var instance database.Driver
	var err error

	switch driver {
	case "postgres":
		instance, err = postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return fmt.Errorf("create postgres driver instance: %w", err)
		}
	case "mysql":
		instance, err = mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			return fmt.Errorf("create mysql driver instance: %w", err)
		}
	default:
		return fmt.Errorf("unsupported driver: %s (supported: postgres, mysql)", driver)
	}

	// Get absolute path to migrations
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("get absolute path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		driver,
		instance,
	)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

// DownMigrations rolls back the last migration.
func DownMigrations(db *sql.DB, driver string, migrationsPath string) error {
	var instance database.Driver
	var err error

	switch driver {
	case "postgres":
		instance, err = postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return fmt.Errorf("create postgres driver instance: %w", err)
		}
	case "mysql":
		instance, err = mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			return fmt.Errorf("create mysql driver instance: %w", err)
		}
	default:
		return fmt.Errorf("unsupported driver: %s", driver)
	}

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("get absolute path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		driver,
		instance,
	)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("rollback migration: %w", err)
	}

	return nil
}
