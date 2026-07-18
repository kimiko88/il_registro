package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for authentication
type Handler struct {
	service *Service
}

// NewHandler creates a new auth handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Register handles user registration
// POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == ErrEmailAlreadyExists {
			statusCode = http.StatusConflict
		} else if err == ErrInvalidEmail || err == ErrPasswordTooShort || err == ErrInvalidRole {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Role:          user.Role,
		SchoolID:      user.SchoolID,
		EmailVerified: user.EmailVerified,
		MFAEnabled:    user.MFAEnabled,
	})
}

// Login handles user login
// POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	authResp, err := h.service.Login(c.Request.Context(), &req, ipAddress, userAgent)
	if err != nil {
		statusCode := http.StatusUnauthorized
		if err == ErrTooManyAttempts {
			statusCode = http.StatusTooManyRequests
		} else if err == ErrMFARequired {
			statusCode = http.StatusPreconditionRequired
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResp)
}

// RefreshToken handles token refresh
// POST /auth/refresh-token
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	tokens, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		statusCode := http.StatusUnauthorized
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Logout handles user logout
// POST /auth/logout
func (h *Handler) Logout(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	// Extract the authenticated user's ID from the JWT (set by Authenticate middleware).
	// This ensures a user can only revoke their own sessions.
	callerUserID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken, callerUserID); err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "token does not belong to the authenticated user"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "logged out successfully"})
}

// GetCurrentUser returns the current authenticated user's profile
// GET /auth/me
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Role:          user.Role,
		SchoolID:      user.SchoolID,
		EmailVerified: user.EmailVerified,
		MFAEnabled:    user.MFAEnabled,
	})
}

// SetupMFA initiates MFA setup
// POST /auth/mfa/setup
func (h *Handler) SetupMFA(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	mfaSetup, err := h.service.SetupMFA(c.Request.Context(), userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == ErrMFAAlreadyEnabled {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, mfaSetup)
}

// VerifyMFA verifies MFA token
// POST /auth/mfa/verify
func (h *Handler) VerifyMFA(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	var req MFAVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	if err := h.service.VerifyMFA(c.Request.Context(), userID, req.Token); err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "MFA verified successfully"})
}

// RequestPasswordReset handles password reset request
// POST /auth/password-reset
func (h *Handler) RequestPasswordReset(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	if err := h.service.RequestPasswordReset(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "password reset email sent"})
}

// ConfirmPasswordReset handles password reset confirmation
// POST /auth/password-reset/confirm
func (h *Handler) ConfirmPasswordReset(c *gin.Context) {
	var req PasswordResetConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		statusCode := http.StatusInternalServerError
		if err == ErrInvalidToken {
			statusCode = http.StatusBadRequest
		} else if err == ErrPasswordTooShort {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "password reset successfully"})
}

// RegisterRoutes registers all auth routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup, middleware *Middleware) {
	auth := router.Group("/auth")
	{
		// Public routes
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh-token", h.RefreshToken)
		auth.POST("/password-reset", h.RequestPasswordReset)
		auth.POST("/password-reset/confirm", h.ConfirmPasswordReset)

		// Protected routes
		protected := auth.Group("")
		protected.Use(middleware.Authenticate())
		{
			protected.GET("/me", h.GetCurrentUser)
			protected.POST("/logout", h.Logout)
			protected.POST("/mfa/setup", h.SetupMFA)
			protected.POST("/mfa/verify", h.VerifyMFA)
		}
	}
}
