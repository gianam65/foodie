package messaging

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Example usage of RabbitMQ publisher and consumer.
//
// This file demonstrates how to use RabbitMQ for event-driven architecture.
// DO NOT import this in production code - it's just for reference.

func ExampleRabbitMQUsage() {
	// 1. Create publisher
	publisher, err := NewPublisher()
	if err != nil {
		log.Fatal("Failed to create publisher:", err)
	}
	defer publisher.(*RabbitMQPublisher).Close()

	// 2. Publish events
	ctx := context.Background()

	// Order created event
	err = publisher.Publish(ctx, Event{
		Type:        "order.created",
		AggregateID: "order-123",
		Payload: map[string]interface{}{
			"order_id":    "order-123",
			"user_id":     "user-456",
			"total":       50.00,
			"items_count": 3,
		},
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		log.Printf("Failed to publish order.created: %v", err)
	}

	// Order delivered event
	err = publisher.Publish(ctx, Event{
		Type:        "order.delivered",
		AggregateID: "order-123",
		Payload: map[string]interface{}{
			"order_id":        "order-123",
			"delivery_time":   time.Now().Format(time.RFC3339),
			"delivery_person": "driver-789",
		},
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		log.Printf("Failed to publish order.delivered: %v", err)
	}

	// 3. Create consumer
	consumer, err := NewConsumer()
	if err != nil {
		log.Fatal("Failed to create consumer:", err)
	}
	defer consumer.Close()

	// 4. Define handler
	orderHandler := func(ctx context.Context, event Event) error {
		log.Printf("Processing event: %s for aggregate: %s", event.Type, event.AggregateID)

		switch event.Type {
		case "order.created":
			// Send confirmation email
			log.Printf("Sending confirmation email for order: %s", event.AggregateID)
			return nil

		case "order.delivered":
			// Update delivery status
			log.Printf("Order delivered: %s", event.AggregateID)
			return nil

		default:
			return fmt.Errorf("unknown event type: %s", event.Type)
		}
	}

	// 5. Start consuming
	err = consumer.Consume(
		ctx,
		"order_processing_queue", // queue name
		"order.*",                // routing pattern: all order events
		orderHandler,
	)
	if err != nil {
		log.Fatal("Failed to start consumer:", err)
	}

	// Keep running
	select {}
}

// ExampleNotificationWorker demonstrates a notification worker.
func ExampleNotificationWorker() {
	consumer, err := NewConsumer()
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	ctx := context.Background()

	// Handler for notification events
	notificationHandler := func(ctx context.Context, event Event) error {
		switch event.Type {
		case "notification.email":
			log.Printf("Sending email notification: %v", event.Payload)
			// Implement email sending logic
			return nil

		case "notification.sms":
			log.Printf("Sending SMS notification: %v", event.Payload)
			// Implement SMS sending logic
			return nil

		default:
			return fmt.Errorf("unknown notification type: %s", event.Type)
		}
	}

	// Consume all notification events
	err = consumer.Consume(
		ctx,
		"notification_queue",
		"notification.*",
		notificationHandler,
	)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
