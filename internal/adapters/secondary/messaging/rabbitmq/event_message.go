package rabbitmq

import (
	"encoding/json"
	"time"
)

type EventMessage struct {
	AggregateID   string
	AggregateType string
	EventType     string
	Version       int
	Timestamp     time.Time
	Payload       json.RawMessage
}
