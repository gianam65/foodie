package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"foodie/backend/internal/infrastructure/database"
	"foodie/backend/internal/infrastructure/scheduler"
	"foodie/backend/internal/infrastructure/scheduler/tasks"
	"foodie/backend/pkg/config"
)

func main() {
	logger := log.New(os.Stdout, "scheduler ", log.LstdFlags|log.Lshortfile)

	// Load environment variables
	config.MustLoad()

	// Initialize database connection
	db, err := database.NewConnectionFromEnv()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	repos, err := database.NewRepositories(db)
	if err != nil {
		logger.Fatalf("Failed to initialize repositories: %v", err)
	}

	// Create scheduler
	sched := scheduler.NewScheduler(logger)

	// Register scheduled tasks
	registerTasks(sched, repos, logger)

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start scheduler
	sched.Start()

	logger.Printf("Scheduler started. Registered tasks: %v", sched.ListTasks())

	// Wait for interrupt signal
	<-ctx.Done()
	logger.Println("Shutdown signal received")

	// Stop scheduler
	sched.Stop()
	logger.Println("Scheduler stopped")
}

// registerTasks registers all scheduled tasks.
func registerTasks(sched *scheduler.Scheduler, repos *database.Repositories, logger *log.Logger) {
	// Health check every 5 minutes
	healthTask := tasks.NewHealthCheckTask(logger)
	if err := sched.AddTask("0 */5 * * * *", healthTask); err != nil {
		logger.Printf("Failed to register health check task: %v", err)
	}

	// Order cleanup every day at 2 AM
	orderCleanupTask := tasks.NewOrderCleanupTask(repos.Order, logger)
	if err := sched.AddTask("0 0 2 * * *", orderCleanupTask); err != nil {
		logger.Printf("Failed to register order cleanup task: %v", err)
	}

	// Cleanup completed orders every Sunday at 3 AM
	completedOrdersTask := tasks.NewCleanupCompletedOrdersTask(repos.Order, logger)
	if err := sched.AddTask("0 0 3 * * 0", completedOrdersTask); err != nil {
		logger.Printf("Failed to register completed orders cleanup task: %v", err)
	}

	// TODO: Add more tasks as needed:
	// - Send order reminders
	// - Update order statuses (auto-complete after delivery time)
	// - Generate daily reports
	// - Sync with external services
}
