package outbound

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
)

type EventStore interface {
	SaveEvents(ctx context.Context, aggregateID string, events []shared.Event) error
	GetEvents(ctx context.Context, aggregateID string) ([]shared.Event, error)
	GetEventsByType(ctx context.Context, aggregateType string) ([]shared.Event, error)
}
