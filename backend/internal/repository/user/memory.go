package user

import (
	"context"
	"sync"
	"time"
)

// MemoryRepository implements the Repository interface using in-memory storage.
// This is primarily intended for testing purposes.
type MemoryRepository struct {
	mu     sync.RWMutex
	users  map[int64]*User
	nextID int64
}

// Ensure MemoryRepository implements Repository interface.
var _ Repository = (*MemoryRepository)(nil)

// NewMemoryRepository creates a new in-memory user repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users:  make(map[int64]*User),
		nextID: 1,
	}
}

// Create creates a new user in memory.
func (r *MemoryRepository) Create(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate email
	for _, u := range r.users {
		if u.Email == user.Email {
			return ErrUserAlreadyExists
		}
		if u.Username == user.Username {
			return ErrUserAlreadyExists
		}
	}

	now := time.Now()
	user.ID = r.nextID
	user.CreatedAt = now
	user.UpdatedAt = now

	r.users[user.ID] = user
	r.nextID++

	return nil
}

// FindByEmail retrieves a user by email address.
func (r *MemoryRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			// Return a copy to avoid mutation
			return r.copyUser(user), nil
		}
	}

	return nil, ErrUserNotFound
}

// FindByUsername retrieves a user by username.
func (r *MemoryRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			// Return a copy to avoid mutation
			return r.copyUser(user), nil
		}
	}

	return nil, ErrUserNotFound
}

// FindByID retrieves a user by their ID.
func (r *MemoryRepository) FindByID(ctx context.Context, id int64) (*User, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}

	return r.copyUser(user), nil
}

// ExistsByEmail checks if a user with the given email exists.
func (r *MemoryRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return true, nil
		}
	}

	return false, nil
}

// ExistsByUsername checks if a user with the given username exists.
func (r *MemoryRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return true, nil
		}
	}

	return false, nil
}

// Clear removes all users from the repository.
// This is useful for testing.
func (r *MemoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = make(map[int64]*User)
	r.nextID = 1
}

// copyUser creates a deep copy of a user.
func (r *MemoryRepository) copyUser(u *User) *User {
	return &User{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// UpdatePassword updates a user's password hash.
func (r *MemoryRepository) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return ErrUserNotFound
	}

	user.PasswordHash = passwordHash
	user.UpdatedAt = time.Now()

	return nil
}

// Update updates an existing user.
func (r *MemoryRepository) Update(ctx context.Context, user *User) error {
	if user == nil || user.ID <= 0 {
		return ErrInvalidID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return ErrUserNotFound
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user

	return nil
}

// Delete removes a user from the repository.
func (r *MemoryRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

// List returns all users.
// This is primarily for testing purposes.
func (r *MemoryRepository) List(ctx context.Context) ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, r.copyUser(user))
	}

	return users, nil
}

// Count returns the total number of users.
func (r *MemoryRepository) Count(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.users), nil
}

// CreateWithID creates a new user with a specific ID.
// This is useful for testing when you want to control the ID.
func (r *MemoryRepository) CreateWithID(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate email
	for _, u := range r.users {
		if u.Email == user.Email {
			return ErrUserAlreadyExists
		}
		if u.Username == user.Username {
			return ErrUserAlreadyExists
		}
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	r.users[user.ID] = user
	if user.ID >= r.nextID {
		r.nextID = user.ID + 1
	}

	return nil
}

// FindByCredentials retrieves a user by either email or username.
// This is useful for login where users can provide either.
func (r *MemoryRepository) FindByCredentials(ctx context.Context, identifier string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == identifier || user.Username == identifier {
			return r.copyUser(user), nil
		}
	}

	return nil, ErrUserNotFound
}
