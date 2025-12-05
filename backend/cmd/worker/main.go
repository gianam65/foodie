package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"foodie/backend/internal/infrastructure/messaging"
)

// WorkerConfig holds configuration for the worker.
type WorkerConfig struct {
	QueueName      string
	RoutingPattern string
	WorkerType     string // "order", "notification", "email", etc.
}

func main() {
	logger := log.New(os.Stdout, "worker ", log.LstdFlags|log.Lshortfile)

	// Parse worker type from command line args
	if len(os.Args) < 2 {
		logger.Fatalf("Usage: %s <worker-type> [queue-name] [routing-pattern]", os.Args[0])
		logger.Fatalf("Examples:")
		logger.Fatalf("  %s order", os.Args[0])
		logger.Fatalf("  %s notification notification_queue notification.*", os.Args[0])
	}

	workerType := os.Args[1]
	queueName := getEnvOrDefault("QUEUE_NAME", fmt.Sprintf("%s_queue", workerType))
	routingPattern := getEnvOrDefault("ROUTING_PATTERN", fmt.Sprintf("%s.*", workerType))

	if len(os.Args) >= 3 {
		queueName = os.Args[2]
	}
	if len(os.Args) >= 4 {
		routingPattern = os.Args[3]
	}

	config := WorkerConfig{
		QueueName:      queueName,
		RoutingPattern: routingPattern,
		WorkerType:     workerType,
	}

	logger.Printf("Starting worker: type=%s, queue=%s, pattern=%s", config.WorkerType, config.QueueName, config.RoutingPattern)

	// Create consumer
	consumer, err := messaging.NewConsumer()
	if err != nil {
		logger.Fatalf("Failed to create consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Printf("Error closing consumer: %v", err)
		}
	}()

	// Create handler based on worker type
	handler := createHandler(config.WorkerType, logger)

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start consuming
	if err := consumer.Consume(ctx, config.QueueName, config.RoutingPattern, handler); err != nil {
		logger.Fatalf("Failed to start consumer: %v", err)
	}

	logger.Printf("Worker started successfully. Waiting for messages...")

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Println("Shutdown signal received, stopping worker...")

	// Give some time for in-flight messages to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	<-shutdownCtx.Done()
	logger.Println("Worker stopped")
}

// createHandler creates a handler function based on worker type.
func createHandler(workerType string, logger *log.Logger) messaging.ConsumerHandler {
	switch workerType {
	case "order":
		return createOrderHandler(logger)
	case "notification", "email":
		return createNotificationHandler(logger)
	case "sms":
		return createSMSHandler(logger)
	default:
		logger.Fatalf("Unknown worker type: %s", workerType)
		return nil
	}
}

// createOrderHandler creates a handler for order events.
func createOrderHandler(logger *log.Logger) messaging.ConsumerHandler {
	return func(ctx context.Context, event messaging.Event) error {
		logger.Printf("Processing order event: %s for order: %s", event.Type, event.AggregateID)

		// Parse payload
		payloadBytes, err := json.Marshal(event.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		switch event.Type {
		case "order.created":
			logger.Printf("Order created: %s", event.AggregateID)
			// TODO: Implement order processing logic
			// - Calculate delivery fee
			// - Assign delivery driver
			// - Update inventory
			// - Send notifications

		case "order.confirmed":
			logger.Printf("Order confirmed: %s", event.AggregateID)
			// TODO: Implement order confirmation logic

		case "order.delivered":
			logger.Printf("Order delivered: %s", event.AggregateID)
			// TODO: Implement delivery completion logic

		case "order.cancelled":
			logger.Printf("Order cancelled: %s", event.AggregateID)
			// TODO: Implement cancellation logic

		default:
			logger.Printf("Unknown order event type: %s", event.Type)
		}

		return nil
	}
}

// createNotificationHandler creates a handler for notification events.
func createNotificationHandler(logger *log.Logger) messaging.ConsumerHandler {
	return func(ctx context.Context, event messaging.Event) error {
		logger.Printf("Processing notification event: %s", event.Type)

		payloadBytes, err := json.Marshal(event.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		switch event.Type {
		case "notification.email":
			to, _ := payload["to"].(string)
			subject, _ := payload["subject"].(string)
			logger.Printf("Sending email to: %s, subject: %s", to, subject)
			// TODO: Implement email sending logic
			// Use service like SendGrid, AWS SES, etc.

		case "notification.sms":
			phone, _ := payload["phone"].(string)
			message, _ := payload["message"].(string)
			logger.Printf("Sending SMS to: %s, message: %s", phone, message)
			// TODO: Implement SMS sending logic
			// Use service like Twilio, AWS SNS, etc.

		default:
			logger.Printf("Unknown notification event type: %s", event.Type)
		}

		return nil
	}
}

// createSMSHandler creates a handler specifically for SMS notifications.
func createSMSHandler(logger *log.Logger) messaging.ConsumerHandler {
	return func(ctx context.Context, event messaging.Event) error {
		logger.Printf("Processing SMS event: %s", event.Type)

		payloadBytes, err := json.Marshal(event.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		phone, _ := payload["phone"].(string)
		message, _ := payload["message"].(string)

		logger.Printf("Sending SMS to %s: %s", phone, message)
		// TODO: Implement SMS sending logic

		return nil
	}
}

// getEnvOrDefault returns environment variable value or default.
func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
