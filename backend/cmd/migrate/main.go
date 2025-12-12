package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"foodie/backend/internal/infrastructure/database"
	"foodie/backend/pkg/config"
	"foodie/backend/pkg/migrate"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	var (
		command        = flag.String("command", "up", "Migration command: up, down")
		driver         = flag.String("driver", "postgres", "Database driver: postgres, mysql")
		migrationsPath = flag.String("path", "./migrations", "Path to migrations directory")
	)
	flag.Parse()

	// Load .env
	config.MustLoad()

	// Get database connection string from environment
	dsn := database.BuildDSN()
	if dsn == "" {
		log.Fatalf("SQL_DSN or (DB_USER/POSTGRES_USER, DB_PASSWORD/POSTGRES_PASSWORD, DB_NAME/POSTGRES_DB) environment variables are required")
	}

	// Open database connection
	db, err := sql.Open(*driver, dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// Run migrations
	switch *command {
	case "up":
		log.Printf("Running migrations from %s...", *migrationsPath)
		if err := migrate.RunMigrations(db, *driver, *migrationsPath); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
		log.Println("Migrations completed successfully")

	case "down":
		log.Printf("Rolling back last migration from %s...", *migrationsPath)
		if err := migrate.DownMigrations(db, *driver, *migrationsPath); err != nil {
			log.Fatalf("rollback failed: %v", err)
		}
		log.Println("Rollback completed successfully")

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s. Use 'up' or 'down'\n", *command)
		os.Exit(1)
	}
}
