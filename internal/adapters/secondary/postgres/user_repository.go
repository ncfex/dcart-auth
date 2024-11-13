package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	userDomain "github.com/ncfex/dcart-auth/internal/core/domain/user"
	"github.com/ncfex/dcart-auth/internal/core/ports/outbound"
	"github.com/ncfex/dcart-auth/internal/infrastructure/database/postgres"
	database "github.com/ncfex/dcart-auth/internal/infrastructure/database/postgres/sqlc"
)

type userRepository struct {
	queries *database.Queries
}

func NewUserRepository(db *postgres.Database) outbound.UserRepository {
	return &userRepository{
		queries: database.New(db.DB),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, userObj *userDomain.User) (*userDomain.User, error) {
	params := database.CreateUserParams{
		Username:     userObj.Username,
		PasswordHash: userObj.PasswordHash,
	}

	dbUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userDomain.ErrUserAlreadyExists
		}
		return nil, err
	}

	return postgres.ToUserDomain(&dbUser), nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID *uuid.UUID) (*userDomain.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, *userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userDomain.ErrUserNotFound
		}
		return nil, err
	}
	return postgres.ToUserDomain(&dbUser), nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*userDomain.User, error) {
	dbUser, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userDomain.ErrUserNotFound
		}
		return nil, err
	}
	return postgres.ToUserDomain(&dbUser), nil
}
