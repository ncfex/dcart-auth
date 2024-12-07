// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error)
	GetTokenByTokenString(ctx context.Context, token string) (RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) (RefreshToken, error)
	SaveToken(ctx context.Context, arg SaveTokenParams) error
}

var _ Querier = (*Queries)(nil)