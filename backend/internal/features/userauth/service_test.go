package userauth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rrlian/papertok/backend/internal/core/auth"
	"github.com/rrlian/papertok/backend/internal/repository/user"
)

// mockAuthService is a mock implementation of authService for testing.
type mockAuthService struct {
	hashPassword   func(ctx context.Context, password string) (string, error)
	verifyPassword func(ctx context.Context, hashedPassword, password string) error
	generateToken  func(ctx context.Context, userID int64, username, email string) (*auth.TokenInfo, error)
	validateToken  func(ctx context.Context, token string) (*auth.Claims, error)
	refreshToken   func(ctx context.Context, token string) (*auth.TokenInfo, error)
}

func (m *mockAuthService) HashPassword(ctx context.Context, password string) (string, error) {
	return m.hashPassword(ctx, password)
}

func (m *mockAuthService) VerifyPassword(ctx context.Context, hashedPassword, password string) error {
	return m.verifyPassword(ctx, hashedPassword, password)
}

func (m *mockAuthService) GenerateToken(ctx context.Context, userID int64, username, email string) (*auth.TokenInfo, error) {
	return m.generateToken(ctx, userID, username, email)
}

func (m *mockAuthService) ValidateToken(ctx context.Context, token string) (*auth.Claims, error) {
	return m.validateToken(ctx, token)
}

func (m *mockAuthService) RefreshToken(ctx context.Context, token string) (*auth.TokenInfo, error) {
	return m.refreshToken(ctx, token)
}

// mockUserRepo is a mock implementation of userRepository for testing.
type mockUserRepo struct {
	create           func(ctx context.Context, u *user.User) error
	findByEmail      func(ctx context.Context, email string) (*user.User, error)
	findByUsername   func(ctx context.Context, username string) (*user.User, error)
	findByID         func(ctx context.Context, id int64) (*user.User, error)
	existsByEmail    func(ctx context.Context, email string) (bool, error)
	existsByUsername func(ctx context.Context, username string) (bool, error)
}

func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error {
	return m.create(ctx, u)
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return m.findByEmail(ctx, email)
}

func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	return m.findByUsername(ctx, username)
}

func (m *mockUserRepo) FindByID(ctx context.Context, id int64) (*user.User, error) {
	return m.findByID(ctx, id)
}

func (m *mockUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return m.existsByEmail(ctx, email)
}

func (m *mockUserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return m.existsByUsername(ctx, username)
}

func TestRegister(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		req           *RegisterRequest
		setupMock     func(*mockAuthService, *mockUserRepo)
		wantErr       error
		checkResponse func(*testing.T, *AuthResponse)
	}{
		{
			name: "successful registration",
			req: &RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "SecurePass123",
			},
			setupMock: func(authSvc *mockAuthService, userRepo *mockUserRepo) {
				userRepo.existsByEmail = func(ctx context.Context, email string) (bool, error) { return false, nil }
				userRepo.existsByUsername = func(ctx context.Context, username string) (bool, error) { return false, nil }
				authSvc.hashPassword = func(ctx context.Context, password string) (string, error) { return "hashed", nil }
				userRepo.create = func(ctx context.Context, u *user.User) error {
					u.ID = 1
					u.CreatedAt = time.Now()
					u.UpdatedAt = time.Now()
					return nil
				}
				authSvc.generateToken = func(ctx context.Context, userID int64, username, email string) (*auth.TokenInfo, error) {
					return &auth.TokenInfo{Token: "jwt-token", ExpiresAt: time.Now().Add(time.Hour)}, nil
				}
			},
			wantErr: nil,
			checkResponse: func(t *testing.T, resp *AuthResponse) {
				if resp.User.Username != "testuser" {
					t.Errorf("Username = %v, want testuser", resp.User.Username)
				}
				if resp.User.Email != "test@example.com" {
					t.Errorf("Email = %v, want test@example.com", resp.User.Email)
				}
				if resp.Token != "jwt-token" {
					t.Errorf("Token = %v, want jwt-token", resp.Token)
				}
			},
		},
		{
			name: "user already exists by email",
			req: &RegisterRequest{
				Username: "testuser",
				Email:    "existing@example.com",
				Password: "SecurePass123",
			},
			setupMock: func(authSvc *mockAuthService, userRepo *mockUserRepo) {
				userRepo.existsByEmail = func(ctx context.Context, email string) (bool, error) { return true, nil }
			},
			wantErr: ErrUserAlreadyExists,
		},
		{
			name: "user already exists by username",
			req: &RegisterRequest{
				Username: "existinguser",
				Email:    "test@example.com",
				Password: "SecurePass123",
			},
			setupMock: func(authSvc *mockAuthService, userRepo *mockUserRepo) {
				userRepo.existsByEmail = func(ctx context.Context, email string) (bool, error) { return false, nil }
				userRepo.existsByUsername = func(ctx context.Context, username string) (bool, error) { return true, nil }
			},
			wantErr: ErrUserAlreadyExists,
		},
		{
			name: "invalid email format",
			req: &RegisterRequest{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "SecurePass123",
			},
			setupMock: func(authSvc *mockAuthService, userRepo *mockUserRepo) {},
			wantErr:   ErrValidationFailed,
		},
		{
			name: "weak password",
			req: &RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "weak",
			},
			setupMock: func(authSvc *mockAuthService, userRepo *mockUserRepo) {},
			wantErr:   ErrValidationFailed,
		},
		{
			name: "invalid username - too short",
			req: &RegisterRequest{
				Username: "ab",
				Email:    "test@example.com",
				Password: "SecurePass123",
			},
			setupMock: func(authSvc *mockAuthService, userRepo *mockUserRepo) {},
			wantErr:   ErrValidationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvc := &mockAuthService{}
			userRepo := &mockUserRepo{}
			tt.setupMock(authSvc, userRepo)

			svc := New(authSvc, userRepo)
			resp, err := svc.Register(ctx, tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Register() expected error %v, got nil", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Fatalf("Register() unexpected error = %v", err)
				}
				if tt.checkResponse != nil {
					tt.checkResponse(t, resp)
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	existingUser := &user.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tests := []struct {
		name          string
		req           *LoginRequest
		setupMock     func(*mockAuthService, *mockUserRepo)
		wantErr       error
		checkResponse func(*testing.T, *AuthResponse)
	}{
		{
			name: "successful login with email",
			req: &LoginRequest{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMock: func(authSvc *mockAuthService, userRepo *mockUserRepo) {
				userRepo.findByEmail = func(ctx context.Context, email string) (*user.User, error) {
					return existingUser, nil
				}
				authSvc.verifyPassword = func(ctx context.Context, hash, password string) error { return nil }
				authSvc.generateToken = func(ctx context.Context, userID int64, username, email string) (*auth.TokenInfo, error) {
					return &auth.TokenInfo{Token: "jwt-token", ExpiresAt: time.Now().Add(time.Hour)}, nil
				}
			},
			wantErr: nil,
			checkResponse: func(t *testing.T, resp *AuthResponse) {
				if resp.User.Username != "testuser" {
					t.Errorf("Username = %v, want testuser", resp.User.Username)
				}
				if resp.Token != "jwt-token" {
					t.Errorf("Token = %v, want jwt-token", resp.Token)
				}
			},
		},
		{
			name: "successful login with username",
			req: &LoginRequest{
				Identifier: "testuser",
				Password:   "password123",
			},
			setupMock: func(authSvc *mockAuthService, userRepo *mockUserRepo) {
				userRepo.findByEmail = func(ctx context.Context, email string) (*user.User, error) {
					return nil, user.ErrUserNotFound
				}
				userRepo.findByUsername = func(ctx context.Context, username string) (*user.User, error) {
					return existingUser, nil
				}
				authSvc.verifyPassword = func(ctx context.Context, hash, password string) error { return nil }
				authSvc.generateToken = func(ctx context.Context, userID int64, username, email string) (*auth.TokenInfo, error) {
					return &auth.TokenInfo{Token: "jwt-token", ExpiresAt: time.Now().Add(time.Hour)}, nil
				}
			},
			wantErr: nil,
		},
		{
			name: "invalid credentials - wrong password",
			req: &LoginRequest{
				Identifier: "test@example.com",
				Password:   "wrongpassword",
			},
			setupMock: func(authSvc *mockAuthService, userRepo *mockUserRepo) {
				userRepo.findByEmail = func(ctx context.Context, email string) (*user.User, error) {
					return existingUser, nil
				}
				authSvc.verifyPassword = func(ctx context.Context, hash, password string) error {
					return auth.ErrInvalidPassword
				}
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "user not found",
			req: &LoginRequest{
				Identifier: "nonexistent@example.com",
				Password:   "password123",
			},
			setupMock: func(authSvc *mockAuthService, userRepo *mockUserRepo) {
				userRepo.findByEmail = func(ctx context.Context, email string) (*user.User, error) {
					return nil, user.ErrUserNotFound
				}
				userRepo.findByUsername = func(ctx context.Context, username string) (*user.User, error) {
					return nil, user.ErrUserNotFound
				}
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "empty identifier",
			req: &LoginRequest{
				Identifier: "",
				Password:   "password123",
			},
			setupMock: func(authSvc *mockAuthService, userRepo *mockUserRepo) {},
			wantErr:   ErrValidationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvc := &mockAuthService{}
			userRepo := &mockUserRepo{}
			tt.setupMock(authSvc, userRepo)

			svc := New(authSvc, userRepo)
			resp, err := svc.Login(ctx, tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Login() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Login() unexpected error = %v", err)
				}
				if tt.checkResponse != nil {
					tt.checkResponse(t, resp)
				}
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	ctx := context.Background()
	existingUser := &user.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		userID  int64
		setup   func(*mockUserRepo)
		wantErr error
	}{
		{
			name:   "successful get profile",
			userID: 1,
			setup: func(repo *mockUserRepo) {
				repo.findByID = func(ctx context.Context, id int64) (*user.User, error) {
					return existingUser, nil
				}
			},
			wantErr: nil,
		},
		{
			name:   "user not found",
			userID: 999,
			setup: func(repo *mockUserRepo) {
				repo.findByID = func(ctx context.Context, id int64) (*user.User, error) {
					return nil, user.ErrUserNotFound
				}
			},
			wantErr: ErrUserNotFound,
		},
		{
			name:    "invalid user ID",
			userID:  -1,
			setup:   func(repo *mockUserRepo) {},
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mockUserRepo{}
			tt.setup(userRepo)

			svc := New(nil, userRepo)
			resp, err := svc.GetProfile(ctx, tt.userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetProfile() expected error %v, got nil", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Fatalf("GetProfile() unexpected error = %v", err)
				}
				if resp.Username != "testuser" {
					t.Errorf("Username = %v, want testuser", resp.Username)
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	ctx := context.Background()
	existingUser := &user.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		token   string
		setup   func(*mockAuthService, *mockUserRepo)
		wantErr error
	}{
		{
			name:  "valid token",
			token: "valid-token",
			setup: func(authSvc *mockAuthService, userRepo *mockUserRepo) {
				authSvc.validateToken = func(ctx context.Context, token string) (*auth.Claims, error) {
					return &auth.Claims{UserID: 1, Username: "testuser", Email: "test@example.com"}, nil
				}
				userRepo.findByID = func(ctx context.Context, id int64) (*user.User, error) {
					return existingUser, nil
				}
			},
			wantErr: nil,
		},
		{
			name:  "invalid token",
			token: "invalid-token",
			setup: func(authSvc *mockAuthService, userRepo *mockUserRepo) {
				authSvc.validateToken = func(ctx context.Context, token string) (*auth.Claims, error) {
					return nil, auth.ErrInvalidToken
				}
			},
			wantErr: auth.ErrInvalidToken,
		},
		{
			name:  "user not found",
			token: "valid-token-but-no-user",
			setup: func(authSvc *mockAuthService, userRepo *mockUserRepo) {
				authSvc.validateToken = func(ctx context.Context, token string) (*auth.Claims, error) {
					return &auth.Claims{UserID: 999}, nil
				}
				userRepo.findByID = func(ctx context.Context, id int64) (*user.User, error) {
					return nil, user.ErrUserNotFound
				}
			},
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvc := &mockAuthService{}
			userRepo := &mockUserRepo{}
			tt.setup(authSvc, userRepo)

			svc := New(authSvc, userRepo)
			_, err := svc.ValidateToken(ctx, tt.token)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidateToken() expected error %v, got nil", tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("ValidateToken() unexpected error = %v", err)
			}
		})
	}
}

// Integration test with real auth service
func TestIntegrationWithRealAuthService(t *testing.T) {
	ctx := context.Background()

	// Create real auth service
	cfg := auth.Config{
		Secret:             "test-secret-key-for-integration-test",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
		Issuer:             "papertok-test",
		PasswordCost:       10,
	}
	realAuthSvc, err := auth.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create auth service: %v", err)
	}

	// Create memory user repository
	memUserRepo := user.NewMemoryRepository()

	// Create userauth service
	svc := New(realAuthSvc, memUserRepo)

	// Test registration flow
	t.Run("register and login flow", func(t *testing.T) {
		// Register new user
		regReq := &RegisterRequest{
			Username: "integrationuser",
			Email:    "integration@test.com",
			Password: "SecurePassword123",
		}

		authResp, err := svc.Register(ctx, regReq)
		if err != nil {
			t.Fatalf("Register() failed: %v", err)
		}

		if authResp.User.Username != regReq.Username {
			t.Errorf("Username = %v, want %v", authResp.User.Username, regReq.Username)
		}

		if authResp.Token == "" {
			t.Error("Token is empty")
		}

		// Login with the same user
		loginReq := &LoginRequest{
			Identifier: regReq.Email,
			Password:   regReq.Password,
		}

		loginResp, err := svc.Login(ctx, loginReq)
		if err != nil {
			t.Fatalf("Login() failed: %v", err)
		}

		if loginResp.User.ID != authResp.User.ID {
			t.Errorf("Login User ID = %v, want %v", loginResp.User.ID, authResp.User.ID)
		}

		// Validate token
		validatedUser, err := svc.ValidateToken(ctx, loginResp.Token)
		if err != nil {
			t.Fatalf("ValidateToken() failed: %v", err)
		}

		if validatedUser.ID != loginResp.User.ID {
			t.Errorf("Validated User ID = %v, want %v", validatedUser.ID, loginResp.User.ID)
		}

		// Test invalid password
		invalidLoginReq := &LoginRequest{
			Identifier: regReq.Email,
			Password:   "WrongPassword123",
		}

		_, err = svc.Login(ctx, invalidLoginReq)
		if err != ErrInvalidCredentials {
			t.Errorf("Expected ErrInvalidCredentials, got: %v", err)
		}

		// Test duplicate registration
		dupRegReq := &RegisterRequest{
			Username: "integrationuser",
			Email:    "another@test.com",
			Password: "SecurePassword123",
		}

		_, err = svc.Register(ctx, dupRegReq)
		if err != ErrUserAlreadyExists {
			t.Errorf("Expected ErrUserAlreadyExists, got: %v", err)
		}
	})
}

func TestJWTStandardClaims(t *testing.T) {
	// This test ensures we're using standard JWT claims correctly
	now := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  int64(123),
		"username": "testuser",
		"email":    "test@example.com",
		"iat":      now,
		"exp":      now + 3600, // 1 hour from now
	})

	tokenString, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	if tokenString == "" {
		t.Error("Token string is empty")
	}

	// Verify the token can be parsed
	parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !parsed.Valid {
		t.Error("Parsed token is not valid")
	}
}
