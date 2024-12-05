package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrHashingPassword  = errors.New("hashing failed")
)

type Password string

func NewPassword(rawPassword string) (Password, error) {
	if len(rawPassword) < 8 {
		return "", ErrPasswordTooShort
	}
	return Password(rawPassword), nil
}

func (p Password) Hash() (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashingPassword
	}
	return string(hashedBytes), nil
}

func (p Password) Matches(hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(p))
	return err == nil
}
