package secondary

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
)

type EventPublisher interface {
	PublishEvent(ctx context.Context, event shared.Event) error
}
