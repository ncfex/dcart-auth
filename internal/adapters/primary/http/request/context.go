package request

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/core/domain"
	userDomain "github.com/ncfex/dcart-auth/internal/core/domain/user"
)

func SetValueToContext(ctx context.Context, key domain.ContextKey, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetUserFromContext(ctx context.Context) (*userDomain.User, bool) {
	user, ok := ctx.Value(domain.ContextUserKey).(*userDomain.User)
	return user, ok
}
