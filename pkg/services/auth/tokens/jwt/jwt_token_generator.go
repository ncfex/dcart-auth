package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	ttl         time.Duration
}

func NewJWTService(issuer, tokenSecret string, ttl time.Duration) *service {
	return &service{
		issuer:      issuer,
		tokenSecret: tokenSecret,
		ttl:         ttl,
	}
}

func (s *service) Generate(subjectString string) (string, error) {
	currentTime := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    s.issuer,
		IssuedAt:  jwt.NewNumericDate(currentTime.UTC()),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(s.ttl)),
		Subject:   subjectString,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.tokenSecret))
	if err != nil {
		return "", ErrTokenSigningFailed
	}
	return token, nil
}

func (s *service) Validate(tokenString string) (string, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) { return []byte(s.tokenSecret), nil },
	)
	if err != nil {
		return "", ErrTokenInvalid
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", ErrTokenInvalidClaims
	}

	if issuer, err := token.Claims.GetIssuer(); err != nil || issuer != string(s.issuer) {
		return "", ErrTokenInvalidClaims
	}

	return userIDString, nil
}
