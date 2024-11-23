package postgres

import (
	"context"
	"database/sql"
	"errors"

	"time"

	"github.com/ncfex/dcart-auth/internal/adapters/secondary/postgres/db"
	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

var (
	ErrStoringToken    = errors.New("error storing token")
	ErrValidatingToken = errors.New("error validating token")
)

type tokenRepository struct {
	queries   *db.Queries
	expiresIn time.Duration
}

func NewTokenRepository(database *database, expiresIn time.Duration) outbound.TokenRepository {
	return &tokenRepository{
		queries:   db.New(database.DB),
		expiresIn: expiresIn,
	}
}

func (r *tokenRepository) StoreToken(ctx context.Context, user *userDomain.User, token string) error {
	params := db.CreateRefreshTokenParams{
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

func (r *tokenRepository) GetTokenByTokenString(ctx context.Context, tokenString string) (*tokenDomain.RefreshToken, error) {
	refreshToken, err := r.queries.GetTokenByTokenString(ctx, tokenString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, tokenDomain.ErrTokenNotFound
		}
		return nil, err
	}

	return db.ToRefreshTokenDomain(&refreshToken), nil
}

func (r *tokenRepository) GetUserFromToken(ctx context.Context, tokenString string) (*userDomain.User, error) {
	user, err := r.queries.GetUserFromRefreshToken(ctx, tokenString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, tokenDomain.ErrTokenNotFound
		}
		return nil, errors.Join(ErrValidatingToken, err)
	}

	return db.ToUserDomain(&user), nil
}

func (r *tokenRepository) RevokeToken(ctx context.Context, tokenString string) error {
	_, err := r.queries.RevokeRefreshToken(ctx, tokenString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tokenDomain.ErrTokenNotFound
		}
		return err
	}

	return nil
}
