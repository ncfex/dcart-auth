package token_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ncfex/dcart-auth/internal/core/services/token"
	"github.com/stretchr/testify/assert"
)

func TestJWTService_MakeJWT(t *testing.T) {
	jwtService := token.NewJWTService("test", "secret")
	userID := uuid.New()

	tests := []struct {
		name        string
		userID      uuid.UUID
		duration    time.Duration
		shouldError bool
	}{
		{
			name:        "valid token",
			userID:      userID,
			duration:    time.Hour,
			shouldError: false,
		},
		{
			name:        "zero duration",
			userID:      userID,
			duration:    0,
			shouldError: false,
		},
		{
			name:        "negative duration",
			userID:      userID,
			duration:    -time.Hour,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwtService.Make(&tt.userID, tt.duration)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestJWTService_ValidateJWT(t *testing.T) {
	jwtService := token.NewJWTService("test", "secret")
	userID := uuid.New()

	tests := []struct {
		name        string
		setupToken  func() string
		shouldError bool
	}{
		{
			name: "valid token",
			setupToken: func() string {
				token, _ := jwtService.Make(&userID, time.Hour)
				return token
			},
			shouldError: false,
		},
		{
			name: "expired token",
			setupToken: func() string {
				token, _ := jwtService.Make(&userID, -time.Hour)
				return token
			},
			shouldError: true,
		},
		{
			name: "invalid token format",
			setupToken: func() string {
				return "invalid.token.string"
			},
			shouldError: true,
		},
		{
			name: "empty token",
			setupToken: func() string {
				return ""
			},
			shouldError: true,
		},
		{
			name: "token with wrong issuer",
			setupToken: func() string {
				wrongIssuerService := token.NewJWTService("wrong-issuer", "secret")
				token, _ := wrongIssuerService.Make(&userID, time.Hour)
				return token
			},
			shouldError: true,
		},
		{
			name: "token with invalid uuid",
			setupToken: func() string {
				claims := jwt.RegisteredClaims{
					Issuer:    "test",
					Subject:   "not-a-uuid",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte("secret"))
				return tokenString
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			validatedUserID, err := jwtService.Validate(token)

			if tt.shouldError {
				assert.Error(t, err)
				assert.Equal(t, uuid.Nil, validatedUserID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, userID, validatedUserID)
			}
		})
	}
}
