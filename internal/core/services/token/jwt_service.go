package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ncfex/dcart/auth-service/internal/core/ports"
)

var (
	ErrTokenExpired  = errors.New("token expired")
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidIssuer = errors.New("invalid issuer")
	ErrInvalidUserID = errors.New("invalid user ID")
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

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.tokenSecret))
}

func (s *service) Validate(tokenString string) (*uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) { return []byte(s.tokenSecret), nil },
	)
	if err != nil {
		return nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return nil, err
	}

	if issuer != string(s.issuer) {
		return nil, ErrInvalidIssuer
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return nil, ErrInvalidUserID
	}
	return &userID, nil
}
