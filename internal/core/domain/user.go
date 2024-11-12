package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	database "github.com/ncfex/dcart-auth/internal/infrastructure/database/sqlc"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) FromDB(dbUser *database.User) {
	u.ID = dbUser.ID
	u.Username = dbUser.Username
	u.PasswordHash = dbUser.PasswordHash
	u.CreatedAt = dbUser.CreatedAt
	u.UpdatedAt = dbUser.UpdatedAt
}

func (u *User) ToDB() *database.User {
	return &database.User{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func NewUserFromDB(dbUser *database.User) *User {
	user := &User{}
	user.FromDB(dbUser)
	return user
}
