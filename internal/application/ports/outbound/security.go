package outbound

import (
	"time"

	"github.com/google/uuid"
)

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
