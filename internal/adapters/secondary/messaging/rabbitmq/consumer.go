package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ncfex/dcart-auth/internal/adapters/secondary/persistence/mongodb"
	"github.com/ncfex/dcart-auth/internal/domain/shared"
	"github.com/ncfex/dcart-auth/internal/domain/user"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerConfig struct {
	URI               string
	Exchange          string
	ExchangeType      string
	Queue             string
	RoutingKey        string
	PrefetchCount     int
	ReconnectDelay    time.Duration
	ProcessingTimeout time.Duration
}

type Consumer struct {
	config    ConsumerConfig
	conn      *amqp.Connection
	channel   *amqp.Channel
	projector *mongodb.MongoProjector
	mu        sync.RWMutex
	connected bool
	stopChan  chan struct{}
	done      chan struct{}
}

func NewConsumer(config ConsumerConfig, projector *mongodb.MongoProjector) (*Consumer, error) {
	consumer := &Consumer{
		config:    config,
		projector: projector,
		stopChan:  make(chan struct{}),
		done:      make(chan struct{}),
	}

	if err := consumer.initialize(); err != nil {
		return nil, fmt.Errorf("consumer initialization failed: %w", err)
	}

	return consumer, nil
}

func (c *Consumer) initialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := amqp.Dial(c.config.URI)
	if err != nil {
		return fmt.Errorf("connection establishment failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("channel creation failed: %w", err)
	}

	// topic
	if err := ch.ExchangeDeclare(
		c.config.Exchange,
		c.config.ExchangeType,
		true, // durable
		false,
		false,
		false,
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("exchange declaration failed: %w", err)
	}

	// dlx
	dlxName := c.config.Exchange + ".dlx"
	if err := ch.ExchangeDeclare(
		dlxName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("DLX declaration failed: %w", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange": dlxName,
		"x-message-ttl":          259200000, // 72hrs
	}
	if _, err := ch.QueueDeclare(
		c.config.Queue,
		true,
		false,
		false,
		false,
		args,
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("queue declaration failed: %w", err)
	}

	// dlq
	dlqName := c.config.Queue + ".dlq"
	if _, err := ch.QueueDeclare(
		dlqName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("DLQ declaration failed: %w", err)
	}

	if err := ch.QueueBind(
		c.config.Queue,
		c.config.RoutingKey,
		c.config.Exchange,
		false,
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("queue binding failed: %w", err)
	}

	if err := ch.QueueBind(
		dlqName,
		c.config.RoutingKey,
		dlxName,
		false,
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("DLQ binding failed: %w", err)
	}

	// qos
	if err := ch.Qos(
		c.config.PrefetchCount,
		0,
		false,
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("QoS setting failed: %w", err)
	}

	c.conn = conn
	c.channel = ch
	c.connected = true

	return nil
}

func (c *Consumer) Start(ctx context.Context) error {
	c.mu.RLock()
	if !c.connected {
		c.mu.RUnlock()
		return fmt.Errorf("consumer not connected")
	}
	c.mu.RUnlock()

	deliveries, err := c.channel.Consume(
		c.config.Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue consume setup failed: %w", err)
	}

	go c.consume(ctx, deliveries)

	return nil
}

func (c *Consumer) consume(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			c.shutdown()
			return
		case <-c.stopChan:
			c.shutdown()
			return
		case delivery, ok := <-deliveries:
			if !ok {
				c.reconnect(ctx)
				return
			}
			c.processDelivery(ctx, delivery)
		}
	}
}

func (c *Consumer) processDelivery(ctx context.Context, delivery amqp.Delivery) {
	processCtx, cancel := context.WithTimeout(ctx, c.config.ProcessingTimeout)
	defer cancel()

	var eventMsg EventMessage
	if err := json.Unmarshal(delivery.Body, &eventMsg); err != nil {
		c.handleProcessingError(delivery, fmt.Errorf("payload deserialization failed: %w", err), "invalid message format")
		return
	}

	event := shared.BaseEvent{
		AggregateID:   eventMsg.AggregateID,
		AggregateType: eventMsg.AggregateType,
		EventType:     eventMsg.EventType,
		Version:       eventMsg.Version,
		Timestamp:     eventMsg.Timestamp,
	}

	// todo improve
	switch eventMsg.EventType {
	case "USER_REGISTERED":
		var payload user.UserRegisteredEventPayload
		if err := json.Unmarshal(eventMsg.Payload, &payload); err != nil {
			c.handleProcessingError(delivery, fmt.Errorf("payload deserialization failed: %w", err), "invalid payload format")
			return
		}
		event.Payload = payload

	case "USER_PASSWORD_CHANGED":
		var payload user.UserPasswordChangedEventPayload
		if err := json.Unmarshal(eventMsg.Payload, &payload); err != nil {
			c.handleProcessingError(delivery, fmt.Errorf("payload deserialization failed: %w", err), "invalid payload format")
			return
		}
		event.Payload = payload

	default:
		c.handleProcessingError(
			delivery,
			fmt.Errorf("unsupported event type: %s", eventMsg.EventType),
			"unknown event type",
		)
		return
	}

	if err := c.projector.ProjectEvent(processCtx, event); err != nil {
		c.handleProcessingError(delivery, fmt.Errorf("projection failed: %w", err), "projection error")
		return
	}

	if err := delivery.Ack(false); err != nil {
		fmt.Printf("failed to ack message: %v\n", err)
	}
}

func (c *Consumer) handleProcessingError(delivery amqp.Delivery, err error, msg string) {
	fmt.Printf("%s: %v\n", msg, err)

	headers := delivery.Headers
	if headers == nil {
		headers = amqp.Table{}
	}

	retryCount, _ := headers["x-retry-count"].(int32)
	if retryCount < 3 {
		headers["x-retry-count"] = retryCount + 1
		delivery.Headers = headers
		delivery.Nack(false, true) // requeue
	} else {
		delivery.Nack(false, false) // DLQ
	}
}

func (c *Consumer) reconnect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopChan:
			return
		case <-time.After(c.config.ReconnectDelay):
			if err := c.initialize(); err != nil {
				fmt.Printf("failed to reconnect: %v\n", err)
				continue
			}
			if err := c.Start(ctx); err != nil {
				fmt.Printf("failed to restart consumer: %v\n", err)
				continue
			}
			return
		}
	}
}

func (c *Consumer) shutdown() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	c.connected = false
	close(c.done)
}

func (c *Consumer) Stop() error {
	close(c.stopChan)
	<-c.done
	return nil
}
