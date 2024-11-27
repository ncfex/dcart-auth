package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/ncfex/dcart-auth/internal/adapters/secondary/postgres/db"
	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

// TODO USE DTO, DON'T USE DOMAIN MODEL DIRECTLY
// TODO DONT USE DOMAIN ERRORS HERE
type userRepository struct {
	queries *db.Queries
}

func NewUserRepository(database *database) outbound.UserRepository {
	return &userRepository{
		queries: db.New(database.DB),
	}
}

func (r *userRepository) Add(ctx context.Context, user *userDomain.User) error {
	params := db.CreateUserParams{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}

	_, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userDomain.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userIDString string) (*userDomain.User, error) {
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return nil, err
	}

	dbUser, err := r.queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userDomain.ErrUserNotFound
		}
		return nil, err
	}
	return db.ToUserDomain(&dbUser), nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*userDomain.User, error) {
	dbUser, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userDomain.ErrUserNotFound
		}
		return nil, err
	}
	return db.ToUserDomain(&dbUser), nil
}
