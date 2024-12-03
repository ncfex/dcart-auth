package infra

import (
	"context"

	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type TokenRepository interface {
	Add(ctx context.Context, token *tokenDomain.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*tokenDomain.RefreshToken, error)
	GetUserByToken(ctx context.Context, token string) (*userDomain.User, error)
	Revoke(ctx context.Context, token string) error
	Save(ctx context.Context, token *tokenDomain.RefreshToken) error
}
