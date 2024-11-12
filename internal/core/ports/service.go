package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ncfex/dcart-auth/internal/core/domain"
)

type UserAuthenticator interface {
	Register(ctx context.Context, username string, password string) (*domain.User, error)
	Login(ctx context.Context, username string, password string) (*domain.TokenPair, error)
	Refresh(ctx context.Context, token string) (*domain.TokenPair, error)
	Logout(ctx context.Context, token string) error
}

type PasswordEncrypter interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

// TODO - more generic interface handle jwt,hex etc.
type TokenManager interface {
	Make(userID *uuid.UUID, expiresIn time.Duration) (string, error)
	Validate(token string) (*uuid.UUID, error)
}
