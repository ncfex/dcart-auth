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

func NewUser(username, password string) (*User, error) {
	if err := validateUserName(username); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	id := uuid.New().String()
	now := time.Now()
	return &User{
		ID:        id,
		Username:  username,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) SetHashedPassword(hashedPassword string) {
	u.PasswordHash = hashedPassword
	u.UpdatedAt = time.Now()
}

func validatePassword(password string) error {
	if password == "" || len(password) < 8 {
		return ErrInvalidCredentials
	}
	return nil
}

func validateUserName(username string) error {
	if username == "" {
		return ErrInvalidCredentials
	}
	return nil
}
