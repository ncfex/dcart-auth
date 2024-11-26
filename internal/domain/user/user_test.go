package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		expectError   bool
		expectedError error
	}{
		{
			name:          "Valid user creation",
			username:      "testuser",
			password:      "password123",
			expectError:   false,
			expectedError: nil,
		},
		{
			name:          "Empty username",
			username:      "",
			password:      "password123",
			expectError:   true,
			expectedError: ErrInvalidCredentials,
		},
		{
			name:          "Empty password",
			username:      "testuser",
			password:      "",
			expectError:   true,
			expectedError: ErrInvalidCredentials,
		},
		{
			name:          "Short password",
			username:      "testuser",
			password:      "short",
			expectError:   true,
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.username, tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.NotEqual(t, "", user.ID)
				assert.Equal(t, tt.username, user.Username)
				assert.Empty(t, user.PasswordHash)
				assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
			}
		})
	}
}

func TestUser_SetHashedPassword(t *testing.T) {
	user, err := NewUser("testuser", "password123")
	assert.NoError(t, err)

	originalCreatedAt := user.CreatedAt
	originalUpdatedAt := user.UpdatedAt

	time.Sleep(time.Millisecond)

	hashedPassword := "hashedpassword123"
	user.SetHashedPassword(hashedPassword)

	assert.Equal(t, hashedPassword, user.PasswordHash)
	assert.Equal(t, originalCreatedAt, user.CreatedAt)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		expectError   bool
		expectedError error
	}{
		{
			name:          "Valid password",
			password:      "password123",
			expectError:   false,
			expectedError: nil,
		},
		{
			name:          "Empty password",
			password:      "",
			expectError:   true,
			expectedError: ErrInvalidCredentials,
		},
		{
			name:          "Short password",
			password:      "short",
			expectError:   true,
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		expectError   bool
		expectedError error
	}{
		{
			name:          "Valid username",
			username:      "testuser",
			expectError:   false,
			expectedError: nil,
		},
		{
			name:          "Empty username",
			username:      "",
			expectError:   true,
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserName(tt.username)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
