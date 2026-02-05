package userauth

import (
	"context"
	"fmt"
	"regexp"
	"unicode"

	"github.com/rrlian/papertok/backend/internal/core/auth"
	"github.com/rrlian/papertok/backend/internal/repository/user"
)

// Impl implements the Service interface.
type Impl struct {
	authSvc  authService
	userRepo userRepository
}

// Ensure Impl implements Service interface.
var _ Service = (*Impl)(nil)

// New creates a new user authentication service instance.
func New(authSvc auth.Service, userRepo user.Repository) *Impl {
	return &Impl{
		authSvc:  authSvc,
		userRepo: userRepo,
	}
}

// Register creates a new user account.
func (s *Impl) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Validate input
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidationFailed, err)
	}

	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	passwordHash, err := s.authSvc.HashPassword(ctx, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	newUser := &user.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token
	tokenInfo, err := s.authSvc.GenerateToken(ctx, newUser.ID, newUser.Username, newUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		User:  s.convertToAuthUser(newUser),
		Token: tokenInfo.Token,
	}, nil
}

// Login authenticates a user with their credentials.
func (s *Impl) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// Validate input
	if req.Identifier == "" || req.Password == "" {
		return nil, ErrValidationFailed
	}

	// Find user by email or username
	var u *user.User
	var err error

	// Try email first
	if isEmail(req.Identifier) {
		u, err = s.userRepo.FindByEmail(ctx, req.Identifier)
	} else {
		// Try username
		u, err = s.userRepo.FindByUsername(ctx, req.Identifier)
	}

	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verify password
	if err := s.authSvc.VerifyPassword(ctx, u.PasswordHash, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	tokenInfo, err := s.authSvc.GenerateToken(ctx, u.ID, u.Username, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		User:  s.convertToAuthUser(u),
		Token: tokenInfo.Token,
	}, nil
}

// GetProfile retrieves a user's profile by their ID.
func (s *Impl) GetProfile(ctx context.Context, userID int64) (*ProfileResponse, error) {
	if userID <= 0 {
		return nil, ErrUserNotFound
	}

	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &ProfileResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

// ValidateToken validates a JWT token and returns the associated user.
func (s *Impl) ValidateToken(ctx context.Context, token string) (*User, error) {
	claims, err := s.authSvc.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	u, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return s.convertToAuthUser(u), nil
}

// RefreshToken generates a new access token from a valid token.
func (s *Impl) RefreshToken(ctx context.Context, token string) (*AuthResponse, error) {
	tokenInfo, err := s.authSvc.RefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Get user info from new token
	claims, err := s.authSvc.ValidateToken(ctx, tokenInfo.Token)
	if err != nil {
		return nil, err
	}

	u, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &AuthResponse{
		User:  s.convertToAuthUser(u),
		Token: tokenInfo.Token,
	}, nil
}

// validateRegisterRequest validates the registration request.
func (s *Impl) validateRegisterRequest(req *RegisterRequest) error {
	if req.Username == "" {
		return ErrInvalidUsername
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		return ErrInvalidUsername
	}

	// Username should only contain alphanumeric characters and underscores
	if !isValidUsername(req.Username) {
		return ErrInvalidUsername
	}

	if req.Email == "" {
		return ErrInvalidEmail
	}

	if !isEmail(req.Email) {
		return ErrInvalidEmail
	}

	if req.Password == "" {
		return ErrWeakPassword
	}

	if len(req.Password) < 8 {
		return ErrWeakPassword
	}

	// Check password strength
	if !isStrongPassword(req.Password) {
		return ErrWeakPassword
	}

	return nil
}

// convertToAuthUser converts a repository User to an auth User.
func (s *Impl) convertToAuthUser(u *user.User) *User {
	return &User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// isEmail checks if a string is a valid email format.
func isEmail(s string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(s)
}

// isValidUsername checks if a username is valid.
func isValidUsername(s string) bool {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(s)
}

// isStrongPassword checks if a password meets minimum security requirements.
func isStrongPassword(p string) bool {
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range p {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Require at least 3 of 4 character types
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	return count >= 2
}
