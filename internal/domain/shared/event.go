package shared

import "time"

type Event interface {
	GetAggregateID() string
	GetAggregateType() string
	GetEventType() string
	GetVersion() int
	GetTimestamp() time.Time
}

type BaseEvent struct {
	AggregateID   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	EventType     string    `json:"event_type"`
	Version       int       `json:"version"`
	Timestamp     time.Time `json:"timestamp"`
}

func (e BaseEvent) GetAggregateID() string   { return e.AggregateID }
func (e BaseEvent) GetAggregateType() string { return e.AggregateType }
func (e BaseEvent) GetEventType() string     { return e.EventType }
func (e BaseEvent) GetVersion() int          { return e.Version }
func (e BaseEvent) GetTimestamp() time.Time  { return e.Timestamp }
