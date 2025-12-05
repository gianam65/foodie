package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQPublisher implements Publisher interface using RabbitMQ.
type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  RabbitMQConfig
}

// RabbitMQConfig holds RabbitMQ connection configuration.
type RabbitMQConfig struct {
	URL      string
	Exchange string
	// Queue configurations
	QueueDurable    bool
	QueueAutoDelete bool
	// Message configurations
	Mandatory   bool
	Immediate   bool
	ContentType string
}

// DefaultRabbitMQConfig returns default RabbitMQ configuration.
func DefaultRabbitMQConfig() RabbitMQConfig {
	return RabbitMQConfig{
		URL:             "amqp://guest:guest@localhost:5672/",
		Exchange:        "foodie_events",
		QueueDurable:    true,
		QueueAutoDelete: false,
		Mandatory:       false,
		Immediate:       false,
		ContentType:     "application/json",
	}
}

// NewRabbitMQPublisher creates a new RabbitMQ publisher.
func NewRabbitMQPublisher(config RabbitMQConfig) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange (topic exchange for flexible routing)
	err = ch.ExchangeDeclare(
		config.Exchange, // name
		"topic",         // type: topic exchange for routing by pattern
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQPublisher{
		conn:    conn,
		channel: ch,
		config:  config,
	}, nil
}

// Publish publishes an event to RabbitMQ exchange.
// Routing key format: "{event.Type}" e.g., "order.created", "order.delivered"
func (p *RabbitMQPublisher) Publish(ctx context.Context, event Event) error {
	// Marshal event payload to JSON
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Routing key is the event type (e.g., "order.created", "notification.email")
	routingKey := event.Type

	// Publish message
	err = p.channel.PublishWithContext(
		ctx,
		p.config.Exchange, // exchange
		routingKey,        // routing key
		p.config.Mandatory,
		p.config.Immediate,
		amqp.Publishing{
			ContentType:  p.config.ContentType,
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Timestamp:    time.Now(),
			MessageId:    fmt.Sprintf("%s-%d", event.AggregateID, time.Now().UnixNano()),
			Headers: amqp.Table{
				"event_type":      event.Type,
				"aggregate_id":    event.AggregateID,
				"event_timestamp": event.Timestamp,
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published event: %s (aggregate: %s)", event.Type, event.AggregateID)
	return nil
}

// Close closes the RabbitMQ connection and channel.
func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return err
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

// RabbitMQConsumer handles consuming messages from RabbitMQ queues.
type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  RabbitMQConfig
}

// NewRabbitMQConsumer creates a new RabbitMQ consumer.
func NewRabbitMQConsumer(config RabbitMQConfig) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		config.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: ch,
		config:  config,
	}, nil
}

// ConsumerHandler is a function that processes consumed messages.
type ConsumerHandler func(ctx context.Context, event Event) error

// Consume creates a queue, binds it to exchange with routing pattern, and starts consuming.
// routingPattern examples:
//   - "order.*" - all order events
//   - "order.created" - only order.created events
//   - "notification.#" - all notification events
func (c *RabbitMQConsumer) Consume(
	ctx context.Context,
	queueName string,
	routingPattern string,
	handler ConsumerHandler,
) error {
	// Declare queue
	q, err := c.channel.QueueDeclare(
		queueName, // name
		c.config.QueueDurable,
		c.config.QueueAutoDelete,
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	err = c.channel.QueueBind(
		q.Name,            // queue name
		routingPattern,    // routing key pattern
		c.config.Exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	// Set QoS (quality of service) - prefetch count
	err = c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Start consuming
	msgs, err := c.channel.Consume(
		q.Name, // queue
		"",     // consumer tag
		false,  // auto-ack (manual ack for reliability)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Started consuming queue: %s with pattern: %s", queueName, routingPattern)

	// Process messages
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Stopping consumer for queue: %s", queueName)
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("Message channel closed for queue: %s", queueName)
					return
				}

				// Parse event
				var event Event
				if err := json.Unmarshal(msg.Body, &event); err != nil {
					log.Printf("Failed to unmarshal event: %v", err)
					msg.Nack(false, false) // Reject message without requeue
					continue
				}

				// Process event
				if err := handler(ctx, event); err != nil {
					log.Printf("Failed to process event %s: %v", event.Type, err)
					msg.Nack(false, true) // Reject and requeue for retry
				} else {
					msg.Ack(false) // Acknowledge message
				}
			}
		}
	}()

	return nil
}

// Close closes the RabbitMQ connection and channel.
func (c *RabbitMQConsumer) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
