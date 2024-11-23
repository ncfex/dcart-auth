package inbound

import (
	"context"

	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type AuthenticationService interface {
	Register(ctx context.Context, username string, password string) (*userDomain.User, error)
	Login(ctx context.Context, username string, password string) (*tokenDomain.TokenPair, error)
	Refresh(ctx context.Context, token string) (*tokenDomain.TokenPair, error)
	Logout(ctx context.Context, token string) error
	Validate(ctx context.Context, token string) (*userDomain.User, error)
}
