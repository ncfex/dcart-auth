package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ncfex/dcart-auth/internal/core/domain"
	"github.com/ncfex/dcart-auth/internal/ports"
)

type service struct {
	issuer      string
	tokenSecret string
}

func NewJWTService(issuer, tokenSecret string) ports.TokenManager {
	return &service{
		issuer:      issuer,
		tokenSecret: tokenSecret,
	}
}

func (s *service) Make(userID *uuid.UUID, expiresIn time.Duration) (string, error) {
	currentTime := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    s.issuer,
		IssuedAt:  jwt.NewNumericDate(currentTime.UTC()),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(expiresIn)),
		Subject:   userID.String(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.tokenSecret))
	if err != nil {
		return "", domain.ErrTokenSigningFailed
	}
	return token, nil
}

func (s *service) Validate(tokenString string) (*uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) { return []byte(s.tokenSecret), nil },
	)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return nil, domain.ErrTokenInvalidClaims
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return nil, domain.ErrTokenInvalidClaims
	}

	if issuer != string(s.issuer) {
		return nil, domain.ErrTokenInvalidIssuer
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return nil, domain.ErrTokenInvalidClaims
	}
	return &userID, nil
}
