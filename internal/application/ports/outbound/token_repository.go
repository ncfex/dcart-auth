package outbound

import (
	"context"

	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type TokenRepository interface {
	StoreToken(ctx context.Context, user *userDomain.User, token string) error
	GetTokenByTokenString(ctx context.Context, token string) (*tokenDomain.RefreshToken, error)
	GetUserFromToken(ctx context.Context, token string) (*userDomain.User, error)
	RevokeToken(ctx context.Context, token string) error
}
