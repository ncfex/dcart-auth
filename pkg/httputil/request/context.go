package request

import (
	"context"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

func SetValueToContext(ctx context.Context, key ContextKey, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetDataFromContext[T any](ctx context.Context, ctxKey ContextKey) (*T, bool) {
	dat, ok := ctx.Value(ContextUserKey).(*T)
	return dat, ok
}
