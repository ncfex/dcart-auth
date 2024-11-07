package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/domain"
	"github.com/ncfex/dcart/auth-service/internal/infrastructure/database/postgres"
	database "github.com/ncfex/dcart/auth-service/internal/infrastructure/database/sqlc"
)

var (
	ErrTokenNotFound   = errors.New("token not found")
	ErrTokenExpired    = errors.New("token expired")
	ErrTokenRevoked    = errors.New("token revoked")
	ErrInvalidToken    = errors.New("invalid token")
	ErrStoringToken    = errors.New("error storing token")
	ErrValidatingToken = errors.New("error validating token")
)

type tokenRepository struct {
	queries   *database.Queries
	expiresIn time.Duration
}

func NewTokenRepository(db *postgres.Database, expiresIn time.Duration) ports.TokenRepository {
	return &tokenRepository{
		queries:   database.New(db.DB),
		expiresIn: expiresIn,
	}
}

func (r *tokenRepository) StoreToken(ctx context.Context, user *domain.User, token string) error {
	params := database.CreateRefreshTokenParams{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(r.expiresIn),
	}

	_, err := r.queries.CreateRefreshToken(ctx, params)
	if err != nil {
		return errors.Join(ErrStoringToken, err)
	}

	return nil
}

func (r *tokenRepository) GetTokenByTokenString(ctx context.Context, token string) (*domain.RefreshToken, error) {
	refreshToken, err := r.queries.GetTokenByTokenString(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	return domain.NewRefreshTokenFromDB(&refreshToken), nil
}

func (r *tokenRepository) GetUserFromToken(ctx context.Context, token string) (*domain.User, error) {
	user, err := r.queries.GetUserFromRefreshToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTokenNotFound
		}
		return nil, errors.Join(ErrValidatingToken, err)
	}

	return domain.NewUserFromDB(&user), nil
}

func (r *tokenRepository) RevokeToken(ctx context.Context, token string) error {
	_, err := r.queries.RevokeRefreshToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTokenNotFound
		}
		return err
	}

	return nil
}
