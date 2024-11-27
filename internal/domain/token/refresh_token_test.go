package token

import (
	"testing"
	"time"
)

func TestNewRefreshToken(t *testing.T) {
	tests := []struct {
		name          string
		tokenString   string
		userID        string
		expectedError error
	}{
		{
			name:          "valid token",
			tokenString:   "valid-token-string",
			userID:        "user-123",
			expectedError: nil,
		},
		{
			name:          "empty token string",
			tokenString:   "",
			userID:        "user-123",
			expectedError: ErrTokenInvalid,
		},
		{
			name:          "empty user ID",
			tokenString:   "valid-token-string",
			userID:        "",
			expectedError: ErrTokenInvalid,
		},
		{
			name:          "both empty",
			tokenString:   "",
			userID:        "",
			expectedError: ErrTokenInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt, err := NewRefreshToken(tt.tokenString, tt.userID)
			if err != tt.expectedError {
				t.Errorf("NewRefreshToken() error = %v, expected error %v", err, tt.expectedError)
				return
			}

			if err == nil {
				if rt.Token != tt.tokenString {
					t.Errorf("Token = %v, expected %v", rt.Token, tt.tokenString)
				}
				if rt.UserID != tt.userID {
					t.Errorf("UserID = %v, expected %v", rt.UserID, tt.userID)
				}
				if rt.CreatedAt.IsZero() {
					t.Error("CreatedAt should not be zero")
				}
				if rt.UpdatedAt.IsZero() {
					t.Error("UpdatedAt should not be zero")
				}
				if rt.ExpiresAt.IsZero() {
					t.Error("ExpiresAt should not be zero")
				}
				if !rt.RevokedAt.IsZero() {
					t.Error("RevokedAt should be zero for new token")
				}

				if rt.CreatedAt != rt.UpdatedAt {
					t.Error("CreatedAt and UpdatedAt should be equal for new token")
				}
			}
		})
	}
}

func TestRefreshToken_Revoke(t *testing.T) {
	rt, err := NewRefreshToken("valid-token", "user-123")
	if err != nil {
		t.Fatalf("Failed to create refresh token: %v", err)
	}

	originalUpdatedAt := rt.UpdatedAt
	time.Sleep(time.Millisecond * 100)

	rt.Revoke()

	if rt.RevokedAt.IsZero() {
		t.Error("RevokedAt should not be zero after revocation")
	}
	if rt.UpdatedAt == originalUpdatedAt {
		t.Error("UpdatedAt should be updated after revocation")
	}
	if rt.RevokedAt != rt.UpdatedAt {
		t.Error("RevokedAt and UpdatedAt should be same")
	}
}

func TestRefreshToken_Expire(t *testing.T) {
	rt, err := NewRefreshToken("valid-token", "user-123")
	if err != nil {
		t.Fatalf("Failed to create refresh token: %v", err)
	}

	originalUpdatedAt := rt.UpdatedAt
	originalExpiresAt := rt.ExpiresAt
	time.Sleep(time.Millisecond * 100)

	rt.Expire()

	if rt.ExpiresAt == originalExpiresAt {
		t.Error("ExpiresAt should be updated after expiration")
	}
	if rt.UpdatedAt == originalUpdatedAt {
		t.Error("UpdatedAt should be updated after expiration")
	}
	if rt.ExpiresAt != rt.UpdatedAt {
		t.Error("ExpiresAt and UpdatedAt should be same")
	}
}

func TestRefreshToken_IsValid(t *testing.T) {
	tests := []struct {
		name          string
		setupToken    func() *RefreshToken
		expectedError error
	}{
		{
			name: "valid token",
			setupToken: func() *RefreshToken {
				rt, _ := NewRefreshToken("valid-token", "user-123")
				rt.ExpiresAt = time.Now().Add(time.Hour)
				return rt
			},
			expectedError: nil,
		},
		{
			name: "revoked token",
			setupToken: func() *RefreshToken {
				rt, _ := NewRefreshToken("valid-token", "user-123")
				rt.ExpiresAt = time.Now().Add(time.Hour)
				rt.Revoke()
				return rt
			},
			expectedError: ErrTokenRevoked,
		},
		{
			name: "expired token",
			setupToken: func() *RefreshToken {
				rt, _ := NewRefreshToken("valid-token", "user-123")
				rt.Expire()
				return rt
			},
			expectedError: ErrTokenExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := tt.setupToken()
			err := rt.IsValid()
			if err != tt.expectedError {
				t.Errorf("IsValid() error = %v, expected error %v", err, tt.expectedError)
			}
		})
	}
}
