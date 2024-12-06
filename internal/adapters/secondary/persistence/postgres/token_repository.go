package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ncfex/dcart-auth/internal/adapters/secondary/persistence/postgres/db"
	"github.com/ncfex/dcart-auth/internal/application/ports/secondary"
	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"

	"github.com/google/uuid"
)

var (
	ErrStoringToken    = errors.New("error storing token")
	ErrValidatingToken = errors.New("error validating token")
)

type tokenRepository struct {
	queries   *db.Queries
	expiresIn time.Duration
}

func NewTokenRepository(database *database, expiresIn time.Duration) secondary.TokenRepository {
	return &tokenRepository{
		queries:   db.New(database.DB),
		expiresIn: expiresIn,
	}
}

func (r *tokenRepository) Add(ctx context.Context, token *tokenDomain.RefreshToken) error {
	userID, err := uuid.Parse(token.UserID)
	if err != nil {
		return userDomain.ErrInvalidCredentials
	}

	params := db.CreateRefreshTokenParams{
		Token:     token.Token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(r.expiresIn),
	}

	_, err = r.queries.CreateRefreshToken(ctx, params)
	if err != nil {
		return errors.Join(ErrStoringToken, err)
	}

	return nil
}

func (r *tokenRepository) GetByToken(ctx context.Context, tokenString string) (*tokenDomain.RefreshToken, error) {
	refreshToken, err := r.queries.GetTokenByTokenString(ctx, tokenString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, tokenDomain.ErrTokenNotFound
		}
		return nil, err
	}

	return db.ToRefreshTokenDomain(&refreshToken), nil
}

func (r *tokenRepository) Revoke(ctx context.Context, tokenString string) error {
	_, err := r.queries.RevokeRefreshToken(ctx, tokenString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tokenDomain.ErrTokenNotFound
		}
		return err
	}

	return nil
}

func (r *tokenRepository) Save(ctx context.Context, token *tokenDomain.RefreshToken) error {
	userID, err := uuid.Parse(token.UserID)
	if err != nil {
		return ErrStoringToken
	}

	revokedAt := sql.NullTime{
		Time:  token.RevokedAt,
		Valid: !token.RevokedAt.IsZero(),
	}

	params := db.SaveTokenParams{
		Token:     token.Token,
		UserID:    userID,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
		ExpiresAt: token.ExpiresAt,
		RevokedAt: revokedAt,
	}

	err = r.queries.SaveToken(ctx, params)
	if err != nil {
		return ErrStoringToken
	}

	return nil
}
