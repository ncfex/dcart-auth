package mqtest

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
)

type EventPublisher struct {
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{}
}

func (ep *EventPublisher) PublishEvent(ctx context.Context, event shared.Event) error {
	fmt.Printf("%s event published for aggregate ID: %s\naggregate type: %s\n at time %s\n", event.GetEventType(), event.GetAggregateID(), event.GetAggregateType(), event.GetTimestamp().UTC())

	return nil
}
