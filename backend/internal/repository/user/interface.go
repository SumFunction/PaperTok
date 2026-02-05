package user

import (
	"context"
	"time"
)

// User represents a user entity in the system.
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Repository defines the interface for user data access operations.
// This abstraction allows for different storage implementations
// (MySQL, PostgreSQL, MongoDB, in-memory, etc.)
type Repository interface {
	// Create creates a new user in the database.
	// Returns ErrUserAlreadyExists if a user with the same email or username exists.
	Create(ctx context.Context, user *User) error

	// FindByEmail retrieves a user by email address.
	// Returns ErrUserNotFound if no user exists with the given email.
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByUsername retrieves a user by username.
	// Returns ErrUserNotFound if no user exists with the given username.
	FindByUsername(ctx context.Context, username string) (*User, error)

	// FindByID retrieves a user by their ID.
	// Returns ErrUserNotFound if no user exists with the given ID.
	FindByID(ctx context.Context, id int64) (*User, error)

	// ExistsByEmail checks if a user with the given email exists.
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername checks if a user with the given username exists.
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
