package outbound

import (
	"time"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type TokenGenerator interface {
	Generate(subjectString string, expiresIn time.Duration) (tokenString string, err error)
	Validate(tokenString string) (subjectString string, err error)
}

type RefreshTokenGenerator interface {
	Generate() (string, error)
}
