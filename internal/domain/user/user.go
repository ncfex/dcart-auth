package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidUser        = errors.New("invalid user")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewUser(username, rawPassword string) (*User, error) {
	if err := validateUserName(username); err != nil {
		return nil, err
	}

	password, err := NewPassword(rawPassword)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := password.Hash()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u *User) Authenticate(rawPassword string) bool {
	password := Password(rawPassword)
	return password.Matches(u.PasswordHash)
}

// todo - use value object
func validateUserName(username string) error {
	if username == "" {
		return ErrInvalidCredentials
	}
	return nil
}
