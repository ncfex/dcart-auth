package rabbitmq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
)

type EventMessage struct {
	AggregateID   string          `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	EventType     string          `json:"event_type"`
	Version       int             `json:"version"`
	Timestamp     time.Time       `json:"timestamp"`
	Payload       json.RawMessage `json:"payload"`
}

func SerializeEvent(event shared.Event) (*EventMessage, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return &EventMessage{
		AggregateID:   event.GetAggregateID(),
		AggregateType: event.GetAggregateType(),
		EventType:     event.GetEventType(),
		Version:       event.GetVersion(),
		Timestamp:     event.GetTimestamp(),
		Payload:       payload,
	}, nil
}

func DeserializeEvent(msg *EventMessage, registry shared.EventRegistry) (shared.Event, error) {
	event, ok := registry.CreateEvent(shared.EventType(msg.EventType))
	if !ok {
		return nil, fmt.Errorf("unknown event type: %s", msg.EventType)
	}

	if err := json.Unmarshal(msg.Payload, event); err != nil {
		return nil, err
	}

	return event, nil
}
