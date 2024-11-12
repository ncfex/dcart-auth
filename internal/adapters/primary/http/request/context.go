package request

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/core/domain"
)

func SetValueToContext(ctx context.Context, key domain.ContextKey, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(domain.ContextUserKey).(*domain.User)
	return user, ok
}
