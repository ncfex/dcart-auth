package shared

import "time"

type Event interface {
	GetAggregateID() string
	GetAggregateType() string
	GetEventType() string
	GetVersion() int
	GetTimestamp() time.Time
	GetPayload() interface{}
}

type BaseEvent struct {
	AggregateID   string
	AggregateType string
	EventType     string
	Version       int
	Timestamp     time.Time
	Payload       interface{}
}

func (e BaseEvent) GetAggregateID() string   { return e.AggregateID }
func (e BaseEvent) GetAggregateType() string { return e.AggregateType }
func (e BaseEvent) GetEventType() string     { return e.EventType }
func (e BaseEvent) GetVersion() int          { return e.Version }
func (e BaseEvent) GetTimestamp() time.Time  { return e.Timestamp }
func (e BaseEvent) GetPayload() interface{}  { return e.Payload }
