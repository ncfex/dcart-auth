package user

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		expectedError error
	}{
		{
			name:          "valid user",
			username:      "testuser",
			password:      "validpass123",
			expectedError: nil,
		},
		{
			name:          "empty username",
			username:      "",
			password:      "validpass123",
			expectedError: ErrInvalidCredentials,
		},
		{
			name:          "password too short",
			username:      "testuser",
			password:      "short",
			expectedError: ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewUser("test", tt.username, tt.password)
			if err != tt.expectedError {
				t.Errorf("NewUser() error = %v, expected error %v", err, tt.expectedError)
				return
			}

			if err == nil {
				if u.Username != tt.username {
					t.Errorf("NewUser() username = %v, expected %v", u.Username, tt.username)
				}

				if u.ID == "" {
					t.Error("NewUser() ID is empty")
				}

				if u.PasswordHash == "" {
					t.Error("NewUser() PasswordHash is empty")
				}

				if u.CreatedAt.IsZero() {
					t.Error("NewUser() CreatedAt is zero")
				}

				if u.UpdatedAt.IsZero() {
					t.Error("NewUser() UpdatedAt is zero")
				}

				if u.CreatedAt != u.UpdatedAt {
					t.Error("NewUser() CreatedAt and UpdatedAt should be equal for new user")
				}
			}
		})
	}
}

func TestUser_Authenticate(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		testPassword  string
		shouldSucceed bool
	}{
		{
			name:          "correct password",
			username:      "testuser",
			password:      "validpass123",
			testPassword:  "validpass123",
			shouldSucceed: true,
		},
		{
			name:          "incorrect password",
			username:      "testuser",
			password:      "validpass123",
			testPassword:  "wrongpass123",
			shouldSucceed: false,
		},
		{
			name:          "empty password",
			username:      "testuser",
			password:      "validpass123",
			testPassword:  "",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewUser("test", tt.username, tt.password)
			if err != nil {
				t.Fatalf("Failed to create user: %v", err)
			}

			if authenticated := u.Authenticate(tt.testPassword); authenticated != tt.shouldSucceed {
				t.Errorf("Authenticate() = %v, expected %v", authenticated, tt.shouldSucceed)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		expectedError error
	}{
		{
			name:          "valid username",
			username:      "testuser",
			expectedError: nil,
		},
		{
			name:          "empty username",
			username:      "",
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUser("test", tt.username, "validpass123")
			if err != tt.expectedError {
				t.Errorf("validateUsername() error = %v, expected error %v", err, tt.expectedError)
			}
		})
	}
}
