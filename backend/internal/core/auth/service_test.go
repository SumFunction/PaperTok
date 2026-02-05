package auth

import (
	"context"
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	ctx := context.Background()
	cfg := TestConfig()
	svc, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create auth service: %v", err)
	}

	tests := []struct {
		name      string
		password  string
		wantErr   bool
		checkHash bool
	}{
		{
			name:      "valid password",
			password:  "SecurePassword123!",
			wantErr:   false,
			checkHash: true,
		},
		{
			name:      "short password",
			password:  "pass",
			wantErr:   false,
			checkHash: true,
		},
		{
			name:      "empty password",
			password:  "",
			wantErr:   false,
			checkHash: true,
		},
		{
			name:      "long password",
			password:  "ThisIsAVeryLongPasswordWithSpecialCharacters!@#$%^&*()_+",
			wantErr:   false,
			checkHash: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := svc.HashPassword(ctx, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkHash {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}

				// Verify hash starts with bcrypt prefix
				if len(hash) < 60 {
					t.Errorf("HashPassword() hash too short: got %d chars, want at least 60", len(hash))
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	ctx := context.Background()
	cfg := TestConfig()
	svc, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create auth service: %v", err)
	}

	password := "TestPassword123!"
	hash, err := svc.HashPassword(ctx, password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        error
	}{
		{
			name:           "correct password",
			hashedPassword: hash,
			password:       password,
			wantErr:        nil,
		},
		{
			name:           "incorrect password",
			hashedPassword: hash,
			password:       "WrongPassword",
			wantErr:        ErrInvalidPassword,
		},
		{
			name:           "empty password",
			hashedPassword: hash,
			password:       "",
			wantErr:        ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.VerifyPassword(ctx, tt.hashedPassword, tt.password)
			if err != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	ctx := context.Background()
	cfg := TestConfig()
	svc, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create auth service: %v", err)
	}

	tokenInfo, err := svc.GenerateToken(ctx, 123, "testuser", "test@example.com")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if tokenInfo.Token == "" {
		t.Error("GenerateToken() returned empty token")
	}

	if tokenInfo.ExpiresAt.Before(time.Now()) {
		t.Error("GenerateToken() expiration time is in the past")
	}
}

func TestValidateToken(t *testing.T) {
	ctx := context.Background()
	cfg := TestConfig()
	svc, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create auth service: %v", err)
	}

	// Generate a token
	tokenInfo, err := svc.GenerateToken(ctx, 456, "validuser", "valid@example.com")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr error
		check   func(*testing.T, *Claims)
	}{
		{
			name:    "valid token",
			token:   tokenInfo.Token,
			wantErr: nil,
			check: func(t *testing.T, claims *Claims) {
				if claims.UserID != 456 {
					t.Errorf("ValidateToken() UserID = %v, want 456", claims.UserID)
				}
				if claims.Username != "validuser" {
					t.Errorf("ValidateToken() Username = %v, want validuser", claims.Username)
				}
				if claims.Email != "valid@example.com" {
					t.Errorf("ValidateToken() Email = %v, want valid@example.com", claims.Email)
				}
			},
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: ErrInvalidToken,
			check:   nil,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: ErrInvalidToken,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := svc.ValidateToken(ctx, tt.token)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidateToken() expected error %v, got nil", tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("ValidateToken() unexpected error = %v", err)
			}

			if err == nil && tt.check != nil {
				tt.check(t, claims)
			}
		})
	}
}

func TestTokenLifecycle(t *testing.T) {
	ctx := context.Background()
	cfg := TestConfig()
	svc, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create auth service: %v", err)
	}

	// Generate token
	tokenInfo, err := svc.GenerateToken(ctx, 789, "lifecycle", "lifecycle@example.com")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// Validate token
	claims, err := svc.ValidateToken(ctx, tokenInfo.Token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != 789 {
		t.Errorf("UserID = %v, want 789", claims.UserID)
	}

	// Refresh token
	newTokenInfo, err := svc.RefreshToken(ctx, tokenInfo.Token)
	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}

	// Validate new token
	newClaims, err := svc.ValidateToken(ctx, newTokenInfo.Token)
	if err != nil {
		t.Fatalf("ValidateToken() on refreshed token error = %v", err)
	}

	if newClaims.UserID != claims.UserID {
		t.Errorf("Refreshed token UserID = %v, want %v", newClaims.UserID, claims.UserID)
	}

	if newClaims.Username != claims.Username {
		t.Errorf("Refreshed token Username = %v, want %v", newClaims.Username, claims.Username)
	}
}

func TestExpiredToken(t *testing.T) {
	ctx := context.Background()
	// For testing expired tokens, we bypass validation and create Impl directly
	svc := &Impl{
		cfg: Config{
			Secret:             "test-secret-for-testing-only",
			AccessTokenExpiry:  -1 * time.Hour, // Already expired
			RefreshTokenExpiry: 7 * 24 * time.Hour,
			Issuer:             "papertok-test",
			PasswordCost:       10,
		},
	}

	tokenInfo, err := svc.GenerateToken(ctx, 999, "expired", "expired@example.com")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, err = svc.ValidateToken(ctx, tokenInfo.Token)
	if err != ErrExpiredToken {
		t.Errorf("ValidateToken() on expired token error = %v, want ErrExpiredToken", err)
	}
}
