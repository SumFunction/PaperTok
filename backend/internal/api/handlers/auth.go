package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rrlian/papertok/backend/internal/features/userauth"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authSvc *userauth.Impl
}

// NewAuthHandler creates a new auth handler instance.
func NewAuthHandler(authSvc *userauth.Impl) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

// RegisterRequest represents the request body for user registration.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

// LoginRequest represents the request body for user login.
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// RegisterHandler handles POST /api/v1/auth/register
// @Summary Register a new user
// @Description Create a new user account with username, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 200 {object} APIResponse{data=userauth.AuthResponse}
// @Failure 400 {object} APIResponse{error=ErrorInfo}
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "VALIDATION_ERROR",
				Message: "Invalid request format",
				Details: err.Error(),
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Convert to feature request
	featureReq := &userauth.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// Call service
	resp, err := h.authSvc.Register(c.Request.Context(), featureReq)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      resp,
		Timestamp: time.Now().Unix(),
	})
}

// LoginHandler handles POST /api/v1/auth/login
// @Summary Login user
// @Description Authenticate a user with email/username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} APIResponse{data=userauth.AuthResponse}
// @Failure 400 {object} APIResponse{error=ErrorInfo}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "VALIDATION_ERROR",
				Message: "Invalid request format",
				Details: err.Error(),
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Convert to feature request
	featureReq := &userauth.LoginRequest{
		Identifier: req.Identifier,
		Password:   req.Password,
	}

	// Call service
	resp, err := h.authSvc.Login(c.Request.Context(), featureReq)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      resp,
		Timestamp: time.Now().Unix(),
	})
}

// GetProfileHandler handles GET /api/v1/auth/profile
// @Summary Get current user profile
// @Description Get the profile of the authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} APIResponse{data=userauth.ProfileResponse}
// @Failure 401 {object} APIResponse{error=ErrorInfo}
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfileHandler(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "INVALID_USER_ID",
				Message: "Invalid user ID in context",
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Call service
	profile, err := h.authSvc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      profile,
		Timestamp: time.Now().Unix(),
	})
}

// RefreshTokenHandler handles POST /api/v1/auth/refresh
// @Summary Refresh access token
// @Description Get a new access token using a valid token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body map[string]string true "Refresh token"
// @Success 200 {object} APIResponse{data=userauth.AuthResponse}
// @Failure 401 {object} APIResponse{error=ErrorInfo}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "VALIDATION_ERROR",
				Message: "Token is required",
				Details: err.Error(),
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Call service
	resp, err := h.authSvc.RefreshToken(c.Request.Context(), req.Token)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      resp,
		Timestamp: time.Now().Unix(),
	})
}

// handleError handles service errors and returns appropriate HTTP responses.
func (h *AuthHandler) handleError(c *gin.Context, err error) {
	code := userauth.GetErrorCode(err)
	message := userauth.GetErrorMessage(err)

	statusCode := http.StatusInternalServerError
	switch code {
	case "USER_EXISTS":
		statusCode = http.StatusConflict
	case "INVALID_CREDENTIALS", "UNAUTHORIZED":
		statusCode = http.StatusUnauthorized
	case "USER_NOT_FOUND":
		statusCode = http.StatusNotFound
	case "VALIDATION_ERROR", "VALIDATION_FAILED":
		statusCode = http.StatusBadRequest
	case "INVALID_EMAIL", "INVALID_USERNAME", "WEAK_PASSWORD":
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
		Timestamp: time.Now().Unix(),
	})
}
