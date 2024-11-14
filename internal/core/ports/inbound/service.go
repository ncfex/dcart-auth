package inbound

import (
	"context"
	"time"

	"github.com/google/uuid"
	tokenDomain "github.com/ncfex/dcart-auth/internal/core/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/core/domain/user"
)

type AuthenticationService interface {
	Register(ctx context.Context, username string, password string) (*userDomain.User, error)
	Login(ctx context.Context, username string, password string) (*tokenDomain.TokenPair, error)
	Refresh(ctx context.Context, token string) (*tokenDomain.TokenPair, error)
	Logout(ctx context.Context, token string) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type TokenGenerator interface {
	Generate(userID uuid.UUID, expiresIn time.Duration) (string, error)
	Validate(token string) (uuid.UUID, error)
}

type RefreshTokenGenerator interface {
	Generate() (string, error)
}
