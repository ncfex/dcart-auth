package password_test

import (
	"testing"

	"github.com/ncfex/dcart-auth/internal/core/services/password"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewPasswordService(t *testing.T) {
	tests := []struct {
		name         string
		cost         int
		expectedCost int
	}{
		{
			name:         "zero cost uses default",
			cost:         0,
			expectedCost: bcrypt.DefaultCost,
		},
		{
			name:         "custom cost",
			cost:         12,
			expectedCost: 12,
		},
		{
			name:         "minimum cost",
			cost:         bcrypt.MinCost,
			expectedCost: bcrypt.MinCost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := password.NewPasswordService(tt.cost)
			assert.NotNil(t, service)

			hash, err := service.Hash("testpassword")
			assert.NoError(t, err)

			cost, err := bcrypt.Cost([]byte(hash))
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCost, cost)
		})
	}
}

func TestPasswordService_HashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		cost        int
		shouldError bool
	}{
		{
			name:        "valid password default cost",
			password:    "password123",
			cost:        0,
			shouldError: false,
		},
		{
			name:        "valid password custom cost",
			password:    "password123",
			cost:        12,
			shouldError: false,
		},
		{
			name:        "empty password",
			password:    "",
			cost:        0,
			shouldError: true,
		},
		{
			name:        "very long password",
			password:    string(make([]byte, 72)),
			cost:        0,
			shouldError: false,
		},
		{
			name:        "too long password",
			password:    string(make([]byte, 73)),
			cost:        0,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := password.NewPasswordService(tt.cost)
			hash, err := service.Hash(tt.password)

			if tt.shouldError {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)

				err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
				assert.NoError(t, err)
			}
		})
	}
}

func TestPasswordService_CheckPasswordHash(t *testing.T) {
	service := password.NewPasswordService(0)
	validPassword := "password123"
	hash, _ := service.Hash(validPassword)

	tests := []struct {
		name        string
		password    string
		hash        string
		shouldError bool
	}{
		{
			name:        "correct password",
			password:    validPassword,
			hash:        hash,
			shouldError: false,
		},
		{
			name:        "incorrect password",
			password:    "wrongpassword",
			hash:        hash,
			shouldError: true,
		},
		{
			name:        "empty password",
			password:    "",
			hash:        hash,
			shouldError: true,
		},
		{
			name:        "invalid hash",
			password:    validPassword,
			hash:        "invalid_hash",
			shouldError: true,
		},
		{
			name:        "empty hash",
			password:    validPassword,
			hash:        "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Compare(tt.hash, tt.password)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
