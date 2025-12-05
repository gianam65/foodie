package messaging

import (
	"fmt"

	"foodie/backend/pkg/config"
)

// NewPublisher creates a publisher based on configuration.
// It reads from environment variables:
//   - MESSAGE_BROKER_TYPE: "rabbitmq", "memory", or "kafka" (default: "memory")
//   - RABBITMQ_URL: RabbitMQ connection URL (default: "amqp://guest:guest@localhost:5672/")
//   - RABBITMQ_EXCHANGE: Exchange name (default: "foodie_events")
func NewPublisher() (Publisher, error) {
	brokerType := config.Get("MESSAGE_BROKER_TYPE", "memory")

	switch brokerType {
	case "rabbitmq":
		return newRabbitMQPublisherFromConfig()
	case "memory":
		return NewInMemoryPublisher(), nil
	case "kafka":
		return nil, fmt.Errorf("Kafka publisher not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported message broker type: %s", brokerType)
	}
}

// newRabbitMQPublisherFromConfig creates a RabbitMQ publisher from environment variables.
func newRabbitMQPublisherFromConfig() (*RabbitMQPublisher, error) {
	cfg := DefaultRabbitMQConfig()

	// Override with environment variables
	if url := config.Get("RABBITMQ_URL", ""); url != "" {
		cfg.URL = url
	}
	if exchange := config.Get("RABBITMQ_EXCHANGE", ""); exchange != "" {
		cfg.Exchange = exchange
	}

	return NewRabbitMQPublisher(cfg)
}

// NewConsumer creates a RabbitMQ consumer from configuration.
func NewConsumer() (*RabbitMQConsumer, error) {
	cfg := DefaultRabbitMQConfig()

	// Override with environment variables
	if url := config.Get("RABBITMQ_URL", ""); url != "" {
		cfg.URL = url
	}
	if exchange := config.Get("RABBITMQ_EXCHANGE", ""); exchange != "" {
		cfg.Exchange = exchange
	}

	return NewRabbitMQConsumer(cfg)
}
