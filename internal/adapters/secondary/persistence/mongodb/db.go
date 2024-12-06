package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	URI            string
	Database       string
	ConnectTimeout time.Duration
	MaxPoolSize    uint64
	MinPoolSize    uint64
}

var (
	ErrConnect       = errors.New("failed to establish mongodb connection")
	ErrInvalidConfig = errors.New("invalid mongodb configuration")
	ErrHealth        = errors.New("mongodb health check failed")
)

type Client struct {
	mongo  *mongo.Client
	db     *mongo.Database
	config Config
}

func NewClient(config Config) (*Client, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	return &Client{
		config: config,
	}, nil
}

func validateConfig(config Config) error {
	if config.URI == "" || config.Database == "" {
		return errors.New("missing required configuration fields")
	}
	return nil
}

func (c *Client) Connect(ctx context.Context) error {
	opts := options.Client().
		ApplyURI(c.config.URI).
		SetConnectTimeout(c.config.ConnectTimeout).
		SetMaxPoolSize(c.config.MaxPoolSize).
		SetMinPoolSize(c.config.MinPoolSize)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrConnect, err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("%w: %v", ErrHealth, err)
	}

	c.mongo = client
	c.db = client.Database(c.config.Database)
	return nil
}

func (c *Client) Database() *mongo.Database {
	return c.db
}

func (c *Client) Disconnect(ctx context.Context) error {
	if c.mongo != nil {
		return c.mongo.Disconnect(ctx)
	}
	return nil
}
