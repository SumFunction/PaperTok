package auth

import "context"

// Service defines the interface for authentication operations.
// This service handles password hashing, JWT token generation,
// and token validation.
type Service interface {
	// HashPassword creates a bcrypt hash of the given password.
	// The cost factor determines the computational cost of hashing.
	HashPassword(ctx context.Context, password string) (string, error)

	// VerifyPassword checks if the provided password matches the hash.
	// Returns ErrInvalidPassword if the passwords don't match.
	VerifyPassword(ctx context.Context, hashedPassword, password string) error

	// GenerateToken creates a new JWT token for the given user.
	// The token contains user ID, username, and email in its claims.
	GenerateToken(ctx context.Context, userID int64, username, email string) (*TokenInfo, error)

	// ValidateToken validates a JWT token and returns its claims.
	// Returns ErrInvalidToken for malformed tokens,
	// ErrExpiredToken for expired tokens.
	ValidateToken(ctx context.Context, token string) (*Claims, error)

	// RefreshToken generates a new token from an existing valid token.
	// This is useful for implementing token refresh logic.
	RefreshToken(ctx context.Context, token string) (*TokenInfo, error)
}
