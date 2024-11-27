package user

import (
	"testing"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name          string
		rawPassword   string
		expectedError error
	}{
		{
			name:          "valid password",
			rawPassword:   "validpass123",
			expectedError: nil,
		},
		{
			name:          "password too short",
			rawPassword:   "short",
			expectedError: ErrPasswordTooShort,
		},
		{
			name:          "empty password",
			rawPassword:   "",
			expectedError: ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := NewPassword(tt.rawPassword)
			if err != tt.expectedError {
				t.Errorf("NewPassword() error = %v, expected error %v", err, tt.expectedError)
				return
			}

			if err == nil && string(password) != tt.rawPassword {
				t.Errorf("NewPassword() = %v, expected %v", password, tt.rawPassword)
			}
		})
	}
}

func TestPassword_Hash(t *testing.T) {
	password, err := NewPassword("validpass123")
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}

	hash1, err := password.Hash()
	if err != nil {
		t.Errorf("Hash() error = %v", err)
		return
	}

	if hash1 == "" {
		t.Error("Hash() returned empty hash")
	}

	// Test that hashing the same password twice produces different hashes
	hash2, err := password.Hash()
	if err != nil {
		t.Errorf("Hash() error = %v", err)
		return
	}

	if hash1 == hash2 {
		t.Error("Hash() produced identical hashes for same password")
	}
}

func TestPassword_Matches(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		matchWith   string
		shouldMatch bool
	}{
		{
			name:        "matching passwords",
			password:    "validpass123",
			matchWith:   "validpass123",
			shouldMatch: true,
		},
		{
			name:        "non-matching passwords",
			password:    "validpass123",
			matchWith:   "differentpass123",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := NewPassword(tt.password)
			if err != nil {
				t.Fatalf("Failed to create password: %v", err)
			}

			hashedPassword, err := password.Hash()
			if err != nil {
				t.Fatalf("Failed to hash password: %v", err)
			}

			comparePassword, err := NewPassword(tt.matchWith)
			if err != nil {
				t.Fatalf("Failed to create comparison password: %v", err)
			}

			if matches := comparePassword.Matches(hashedPassword); matches != tt.shouldMatch {
				t.Errorf("Matches() = %v, expected %v", matches, tt.shouldMatch)
			}
		})
	}
}
