package password

import (
	"errors"

	"github.com/ncfex/dcart-auth/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyPassword = errors.New("password cannot be empty")
)

type service struct {
	cost int
}

func NewPasswordService(cost int) ports.PasswordEncrypter {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &service{cost: cost}
}

func (s *service) Hash(password string) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	data, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *service) Compare(hashedPassword, password string) error {
	if password == "" {
		return ErrEmptyPassword
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
