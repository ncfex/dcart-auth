package infra

import (
	"context"

	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
)

type TokenRepository interface {
	Add(ctx context.Context, token *tokenDomain.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*tokenDomain.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	Save(ctx context.Context, token *tokenDomain.RefreshToken) error
}
