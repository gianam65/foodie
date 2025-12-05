package messaging

import (
	"context"
	"fmt"
)

// Event represents a domain event to be published.
type Event struct {
	Type        string
	AggregateID string
	Payload     interface{}
	Timestamp   int64
}

// Publisher publishes domain events to message brokers.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

// InMemoryPublisher is a simple in-memory implementation for development/testing.
// In production, replace with Kafka, RabbitMQ, or other message broker adapter.
type InMemoryPublisher struct {
	events []Event
}

// NewInMemoryPublisher creates a new in-memory event publisher.
func NewInMemoryPublisher() *InMemoryPublisher {
	return &InMemoryPublisher{
		events: make([]Event, 0),
	}
}

// Publish publishes an event to in-memory store.
func (p *InMemoryPublisher) Publish(ctx context.Context, event Event) error {
	p.events = append(p.events, event)
	return nil
}

// GetEvents returns all published events (for testing).
func (p *InMemoryPublisher) GetEvents() []Event {
	return p.events
}

// Clear clears all events (for testing).
func (p *InMemoryPublisher) Clear() {
	p.events = make([]Event, 0)
}

// OutboxPublisher implements outbox pattern for reliable event publishing.
// Events are stored in database first, then published asynchronously.
type OutboxPublisher struct {
	db interface{} // TODO: Add database connection
}

// NewOutboxPublisher creates a new outbox-based event publisher.
func NewOutboxPublisher(db interface{}) *OutboxPublisher {
	return &OutboxPublisher{db: db}
}

// Publish stores event in outbox table within the same transaction.
func (p *OutboxPublisher) Publish(ctx context.Context, event Event) error {
	// TODO: Insert into outbox table
	// This should be called within a database transaction
	return fmt.Errorf("not implemented")
}

// ProcessOutbox processes events from outbox table and publishes them.
// This should run as a background worker.
func (p *OutboxPublisher) ProcessOutbox(ctx context.Context) error {
	// TODO: Select events from outbox, publish, mark as processed
	return fmt.Errorf("not implemented")
}
