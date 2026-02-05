package userauth

import "context"

// Service defines the interface for user authentication operations.
// This includes registration, login, and user profile management.
type Service interface {
	// Register creates a new user account.
	// Returns ErrUserAlreadyExists if the email or username is already taken.
	// Returns ErrValidationFailed if input validation fails.
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)

	// Login authenticates a user with their credentials.
	// The identifier can be either an email or username.
	// Returns ErrInvalidCredentials if the credentials are incorrect.
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)

	// GetProfile retrieves a user's profile by their ID.
	// Returns ErrUserNotFound if the user doesn't exist.
	GetProfile(ctx context.Context, userID int64) (*ProfileResponse, error)

	// ValidateToken validates a JWT token and returns the associated user.
	// Returns ErrInvalidToken for invalid tokens.
	// Returns ErrExpiredToken for expired tokens.
	ValidateToken(ctx context.Context, token string) (*User, error)

	// RefreshToken generates a new access token from a valid token.
	RefreshToken(ctx context.Context, token string) (*AuthResponse, error)
}
