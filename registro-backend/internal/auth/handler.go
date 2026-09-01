package auth

import (
	"fmt"
	"net/http"
	"os"

	"registro-backend/pkg/wsticket"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for authentication
type Handler struct {
	service       *Service
	wsTicketStore *wsticket.Store
}

// NewHandler creates a new auth handler
func NewHandler(service *Service, wsTicketStore ...*wsticket.Store) *Handler {
	var store *wsticket.Store
	if len(wsTicketStore) > 0 {
		store = wsTicketStore[0]
	}
	return &Handler{
		service:       service,
		wsTicketStore: store,
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

func isHTTPSRequest(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	if os.Getenv("TRUST_PROXY_HEADERS") == "true" {
		if c.GetHeader("X-Forwarded-Proto") == "https" || c.GetHeader("X-Forwarded-Ssl") == "on" {
			return true
		}
	}
	return false
}

func setRefreshTokenCookie(c *gin.Context, token string, maxAge int) {
	appEnv := os.Getenv("APP_ENV")
	ginMode := os.Getenv("GIN_MODE")
	cookieSecure := os.Getenv("COOKIE_SECURE")

	isSecure := cookieSecure == "true" || appEnv == "production" || ginMode == "release" || isHTTPSRequest(c)
	if cookieSecure == "false" && appEnv != "production" && ginMode != "release" {
		isSecure = false
	}
	domain := os.Getenv("COOKIE_DOMAIN")
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refreshToken",
		Value:    token,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   domain,
		Secure:   isSecure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
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
		var statusCode int
		switch err {
		case ErrTooManyAttempts:
			statusCode = http.StatusTooManyRequests
		case ErrMFARequired:
			statusCode = http.StatusPreconditionRequired
		default:
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	if authResp != nil && authResp.RefreshToken != "" {
		setRefreshTokenCookie(c, authResp.RefreshToken, 604800)
		authResp.RefreshToken = ""
	}

	c.JSON(http.StatusOK, authResp)
}

// RefreshToken handles token refresh
// POST /auth/refresh-token
func (h *Handler) RefreshToken(c *gin.Context) {
	rt := ""
	if cookieToken, err := c.Cookie("refreshToken"); err == nil && cookieToken != "" {
		rt = cookieToken
	} else if c.Request.Body != nil && c.Request.ContentLength != 0 {
		var req RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Message: err.Error()})
			return
		}
		rt = req.RefreshToken
	}

	if rt == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "missing refresh_token"})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	tokens, err := h.service.RefreshToken(c.Request.Context(), rt, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	if tokens != nil && tokens.RefreshToken != "" {
		setRefreshTokenCookie(c, tokens.RefreshToken, 604800)
		tokens.RefreshToken = ""
	}

	c.JSON(http.StatusOK, tokens)
}

// Logout handles user logout
// POST /auth/logout
func (h *Handler) Logout(c *gin.Context) {
	rt := ""
	if cookieToken, err := c.Cookie("refreshToken"); err == nil && cookieToken != "" {
		rt = cookieToken
	} else if c.Request.Body != nil && c.Request.ContentLength != 0 {
		var req RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Message: err.Error()})
			return
		}
		rt = req.RefreshToken
	}

	callerUserID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.service.Logout(c.Request.Context(), rt, callerUserID); err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "token does not belong to the authenticated user"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	setRefreshTokenCookie(c, "", -1)
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
		IsStaff:       user.IsStaff,
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

	codes, err := h.service.VerifyMFA(c.Request.Context(), userID, req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "MFA verified successfully",
		"recovery_codes": codes,
	})
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

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case ErrInvalidCredentials:
			statusCode = http.StatusUnauthorized
		case ErrPasswordTooShort, ErrPasswordTooLong, ErrPasswordNoUppercase, ErrPasswordNoLowercase, ErrPasswordNoNumber, ErrPasswordNoSpecial:
			statusCode = http.StatusBadRequest
		case ErrPasswordReused:
			statusCode = http.StatusUnprocessableEntity
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "password changed successfully"})
}

// IssueWSTicket issues a single-use opaque ticket to authenticate a WebSocket connection.
// Requires a valid JWT in the Authorization header.
func (h *Handler) IssueWSTicket(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}
	if h.wsTicketStore == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "ws ticket store not configured"})
		return
	}

	emailStr := ""
	if emailVal, ok := c.Get("email"); ok && emailVal != nil {
		emailStr = fmt.Sprint(emailVal)
	}
	role, roleExists := GetUserRole(c)
	schoolIDStr := ""
	if sID, ok := GetSchoolID(c); ok {
		schoolIDStr = sID
	}

	if emailStr == "" || !roleExists || role == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "incomplete user metadata for ws ticket"})
		return
	}

	ticket, err := h.wsTicketStore.Issue(
		userID,
		emailStr,
		role,
		schoolIDStr,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "could not issue ws ticket"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

// RegisterRoutes registers all auth routes.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup, middleware *Middleware) {
	auth := router.Group("/auth")
	{
		// Fully public routes (no JWT required)
		auth.POST("/login", h.Login)
		auth.POST("/refresh-token", h.RefreshToken)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/password-reset", h.RequestPasswordReset)
		auth.POST("/password-reset/confirm", h.ConfirmPasswordReset)

		// Protected routes (JWT required)
		protected := auth.Group("")
		protected.Use(middleware.Authenticate())
		{
			protected.GET("/me", h.GetCurrentUser)
			protected.POST("/logout", h.Logout)
			protected.POST("/change-password", h.ChangePassword)
			protected.POST("/mfa/setup", h.SetupMFA)
			protected.POST("/mfa/verify", h.VerifyMFA)
			protected.POST("/ws-ticket", h.IssueWSTicket)

			// Registration is protected: caller must be superadmin, admin or segreteria.
			protected.POST("/register", h.Register)
		}
	}
}
