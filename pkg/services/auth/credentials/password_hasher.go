package credentials

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooLong = errors.New("password exceeds maximum length of 72 bytes")
	ErrHashingFailed   = errors.New("hashing failed")
	ErrComparingFailed = errors.New("comparing failed")
)

const maxPasswordLength = 72

type bcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *bcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &bcryptHasher{cost: cost}
}

func (s *bcryptHasher) Hash(password string) (string, error) {
	if len(password) > maxPasswordLength {
		return "", ErrPasswordTooLong
	}

	data, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return "", ErrPasswordTooLong
		}
		return "", ErrHashingFailed
	}
	return string(data), nil
}

func (s *bcryptHasher) Compare(hashedPassword, password string) error {
	if len(password) > maxPasswordLength {
		return ErrPasswordTooLong
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrComparingFailed
		}
		return err
	}
	return nil
}
