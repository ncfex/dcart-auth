package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URI          string
	Exchange     string
	ExchangeType string
	RoutingKey   string
	Timeout      time.Duration
}

type RabbitMQAdapter struct {
	config    RabbitMQConfig
	conn      *amqp.Connection
	channel   *amqp.Channel
	mu        sync.RWMutex
	connected bool
}

func NewRabbitMQAdapter(config RabbitMQConfig) (*RabbitMQAdapter, error) {
	adapter := &RabbitMQAdapter{
		config: config,
	}

	if err := adapter.initialize(); err != nil {
		return nil, fmt.Errorf("adapter initialization failed: %w", err)
	}

	return adapter, nil
}

func (a *RabbitMQAdapter) PublishEvent(ctx context.Context, event shared.Event) error {
	a.mu.RLock()
	if !a.connected {
		a.mu.RUnlock()
		return fmt.Errorf("rabbitmq connection unavailable")
	}
	a.mu.RUnlock()

	// todo use gob instead of json
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("event serialization failed: %w", err)
	}

	confirms := a.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
		Timestamp:   event.GetTimestamp(),
		Headers: amqp.Table{
			"aggregate_id":   event.GetAggregateID(),
			"aggregate_type": event.GetAggregateType(),
			"event_type":     event.GetEventType(),
			"version":        event.GetVersion(),
		},
		DeliveryMode: amqp.Persistent,
	}

	if err := a.channel.PublishWithContext(
		ctx,
		a.config.Exchange,
		a.config.RoutingKey,
		true,
		false,
		msg,
	); err != nil {
		return fmt.Errorf("message publishing failed: %w", err)
	}

	select {
	case confirm := <-confirms:
		if !confirm.Ack {
			return fmt.Errorf("message delivery unconfirmed by broker")
		}
	case <-ctx.Done():
		return fmt.Errorf("publish confirmation timeout: %w", ctx.Err())
	}

	return nil
}

func (a *RabbitMQAdapter) initialize() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	conn, err := amqp.Dial(a.config.URI)
	if err != nil {
		return fmt.Errorf("connection establishment failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("channel creation failed: %w", err)
	}

	if err := ch.Confirm(false); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("publisher confirms activation failed: %w", err)
	}

	if err := ch.ExchangeDeclare(
		a.config.Exchange,
		a.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("exchange declaration failed: %w", err)
	}

	a.conn = conn
	a.channel = ch
	a.connected = true

	return nil
}

func (a *RabbitMQAdapter) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.channel != nil {
		if err := a.channel.Close(); err != nil {
			return fmt.Errorf("channel closure failed: %w", err)
		}
	}

	if a.conn != nil {
		if err := a.conn.Close(); err != nil {
			return fmt.Errorf("connection closure failed: %w", err)
		}
	}

	return nil
}
