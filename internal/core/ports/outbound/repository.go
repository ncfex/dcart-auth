package outbound

import (
	"context"

	"github.com/google/uuid"
	tokenDomain "github.com/ncfex/dcart-auth/internal/core/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/core/domain/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *userDomain.User) (*userDomain.User, error)
	GetUserByID(ctx context.Context, userID *uuid.UUID) (*userDomain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*userDomain.User, error)
}

type TokenRepository interface {
	StoreToken(ctx context.Context, user *userDomain.User, token string) error
	GetTokenByTokenString(ctx context.Context, token string) (*tokenDomain.RefreshToken, error)
	GetUserFromToken(ctx context.Context, token string) (*userDomain.User, error)
	RevokeToken(ctx context.Context, token string) error
}
