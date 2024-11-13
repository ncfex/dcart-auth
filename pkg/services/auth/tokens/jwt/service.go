package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrTokenSigningFailed = errors.New("token signing failed")
	ErrTokenInvalid       = errors.New("token invalid")
	ErrTokenInvalidClaims = errors.New("token invalid claims")
	ErrTokenInvalidIssuer = errors.New("token invalid issuer")
)

type service struct {
	issuer      string
	tokenSecret string
}

func NewJWTService(issuer, tokenSecret string) *service {
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
		return "", ErrTokenSigningFailed
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
		return nil, ErrTokenInvalid
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return nil, ErrTokenInvalidClaims
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return nil, ErrTokenInvalidClaims
	}

	if issuer != string(s.issuer) {
		return nil, ErrTokenInvalidIssuer
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return nil, ErrTokenInvalidClaims
	}
	return &userID, nil
}
