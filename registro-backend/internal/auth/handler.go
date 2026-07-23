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

// Register handles user registration by an authenticated privileged user.
// The caller's role is extracted from the JWT (set by Authenticate middleware)
// and enforced against the role permission matrix in ValidateRegisterRequest.
// This endpoint is no longer public: it requires a valid JWT.
//
// POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	// callerRole comes from the validated JWT — never trust the request body for this.
	callerRole, exists := GetUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := ValidateRegisterRequest(callerRole, &req); err != nil {
		statusCode := http.StatusBadRequest
		switch err {
		case ErrInsufficientRole:
			statusCode = http.StatusForbidden
		case ErrCannotCreateRole:
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == ErrEmailAlreadyExists {
			statusCode = http.StatusConflict
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

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	tokens, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
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

	_ = h.service.RequestPasswordReset(c.Request.Context(), req.Email)
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
		switch err {
		case ErrInvalidToken:
			statusCode = http.StatusBadRequest
		case ErrPasswordTooShort:
			statusCode = http.StatusBadRequest
		case ErrPasswordReused:
			// 422 Unprocessable Entity: request ben formata ma password
			// viola la policy di riutilizzo.
			statusCode = http.StatusUnprocessableEntity
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "password reset successfully"})
}

// RegisterRoutes registers all auth routes.
//
// BREAKING CHANGE: POST /auth/register è ora un endpoint PROTETTO.
// Richiede un JWT valido nel header Authorization: Bearer <token>.
// Il caller deve avere ruolo superadmin, admin o segreteria.
// La matrice dei permessi di creazione ruoli è applicata in ValidateRegisterRequest.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup, middleware *Middleware) {
	auth := router.Group("/auth")
	{
		// Fully public routes (no JWT required)
		auth.POST("/login", h.Login)
		auth.POST("/refresh-token", h.RefreshToken)
		auth.POST("/password-reset", h.RequestPasswordReset)
		auth.POST("/password-reset/confirm", h.ConfirmPasswordReset)

		// Protected routes (JWT required)
		protected := auth.Group("")
		protected.Use(middleware.Authenticate())
		{
			protected.GET("/me", h.GetCurrentUser)
			protected.POST("/logout", h.Logout)
			protected.POST("/mfa/setup", h.SetupMFA)
			protected.POST("/mfa/verify", h.VerifyMFA)

			// Registration is protected: caller must be superadmin, admin or segreteria.
			// Role-level permission checks are enforced inside the handler via
			// ValidateRegisterRequest (validator.go).
			protected.POST("/register", h.Register)
		}
	}
}
