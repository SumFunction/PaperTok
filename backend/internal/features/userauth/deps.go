package userauth

import (
	"context"

	"github.com/rrlian/papertok/backend/internal/core/auth"
	"github.com/rrlian/papertok/backend/internal/repository/user"
)

// authService defines the authentication service capability required by this feature.
type authService interface {
	// HashPassword creates a bcrypt hash of the given password.
	HashPassword(ctx context.Context, password string) (string, error)

	// VerifyPassword checks if the provided password matches the hash.
	VerifyPassword(ctx context.Context, hashedPassword, password string) error

	// GenerateToken creates a new JWT token for the given user.
	GenerateToken(ctx context.Context, userID int64, username, email string) (*auth.TokenInfo, error)

	// ValidateToken validates a JWT token and returns its claims.
	ValidateToken(ctx context.Context, token string) (*auth.Claims, error)

	// RefreshToken generates a new token from an existing valid token.
	RefreshToken(ctx context.Context, token string) (*auth.TokenInfo, error)
}

// userRepository defines the user repository capability required by this feature.
type userRepository interface {
	// Create creates a new user in the database.
	Create(ctx context.Context, user *user.User) error

	// FindByEmail retrieves a user by email address.
	FindByEmail(ctx context.Context, email string) (*user.User, error)

	// FindByUsername retrieves a user by username.
	FindByUsername(ctx context.Context, username string) (*user.User, error)

	// FindByID retrieves a user by their ID.
	FindByID(ctx context.Context, id int64) (*user.User, error)

	// ExistsByEmail checks if a user with the given email exists.
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername checks if a user with the given username exists.
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
