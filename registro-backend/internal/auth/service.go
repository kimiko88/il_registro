package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"registro-backend/pkg/crypto"
	"registro-backend/pkg/jwt"
	"registro-backend/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// Service handles authentication business logic
type Service struct {
	repo         Repository
	tokenManager *jwt.TokenManager
	mfaService   *MFAService
	bcryptCost   int
}

// NewService creates a new auth service
func NewService(repo Repository, tokenManager *jwt.TokenManager, mfaService *MFAService) *Service {
	return &Service{
		repo:         repo,
		tokenManager: tokenManager,
		mfaService:   mfaService,
		bcryptCost:   12,
	}
}

// Pre-computed dummy bcrypt hash (cost 12) used to normalize timing when user email does not exist.
const dummyBcryptHash = "$2a$12$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeg6Lruj3vjPGga31lW"

// Register creates a new user account.
// Input validation and RBAC checks are performed by the handler layer
// (ValidateRegisterRequest in validator.go) before this method is called.
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	if err := NewPasswordValidator().Validate(req.Password); err != nil {
		return nil, err
	}

	// Check if email already exists
	_, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
	}
	if req.SchoolID != "" {
		user.SchoolID = &req.SchoolID
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Seed password history so ResetPassword can detect immediate re-use.
	_ = s.repo.AddPasswordHistory(ctx, user.ID, string(passwordHash))

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *Service) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*AuthResponse, error) {
	// --- Rate limiting ---
	// 1. Per (email, IP): max 5 attempts in 15 minutes — blocks single-IP bursts.
	// 2. Per email only: max 20 attempts in 1 hour — blocks distributed IP-rotation attacks.
	since15m := time.Now().Add(-15 * time.Minute)
	attemptsPerIP, err := s.repo.GetRecentLoginAttempts(ctx, req.Email, ipAddress, since15m)
	if err != nil {
		return nil, err
	}
	if attemptsPerIP >= 5 {
		return nil, ErrTooManyAttempts
	}

	since1h := time.Now().Add(-1 * time.Hour)
	attemptsPerEmail, err := s.repo.GetRecentLoginAttemptsByEmail(ctx, req.Email, since1h)
	if err != nil {
		return nil, err
	}
	if attemptsPerEmail >= 20 {
		return nil, ErrTooManyAttempts
	}

	// Get user
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Run bcrypt against dummy hash to prevent timing attacks that enumerate valid user accounts
		_ = bcrypt.CompareHashAndPassword([]byte(dummyBcryptHash), []byte(req.Password))
		s.recordFailedAttempt(ctx, req.Email, ipAddress)
		return nil, ErrInvalidCredentials
	}

	// Verify password BEFORE active check to maintain constant response timing
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if bcryptErr != nil {
		s.recordFailedAttempt(ctx, req.Email, ipAddress)
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		s.recordFailedAttempt(ctx, req.Email, ipAddress)
		return nil, ErrUserInactive
	}

	// Check if password has expired (90 days for privileged roles: superadmin, admin, segreteria, teacher)
	if user.Role == RoleSuperAdmin || user.Role == RoleAdmin || user.Role == RoleSegreteria || user.Role == RoleTeacher {
		if user.PasswordChangedAt != nil {
			if time.Since(*user.PasswordChangedAt) > 90*24*time.Hour {
				return nil, ErrPasswordExpired
			}
		} else if !user.CreatedAt.IsZero() && time.Since(user.CreatedAt) > 90*24*time.Hour {
			return nil, ErrPasswordExpired
		}
	}

	// Check MFA — decrypt stored secret before verifying the TOTP code.
	if user.MFAEnabled {
		if req.MFAToken == "" {
			return nil, ErrMFARequired
		}
		encryptedSecret, err := s.repo.GetMFASecret(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		secret, err := crypto.DecryptString(encryptedSecret)
		if err != nil {
			return nil, err
		}
		if !s.mfaService.VerifyTOTPUser(user.ID, secret, req.MFAToken) {
			return nil, ErrInvalidMFAToken
		}
	}

	// Generate tokens
	accessToken, err := s.tokenManager.GenerateAccessToken(
		user.ID, user.Email, user.Role,
		func() string {
			if user.SchoolID != nil {
				return *user.SchoolID
			}
			return ""
		}(),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	rt := &RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	if err := s.repo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	// Update last login
	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	// Record successful attempt
	s.recordSuccessfulAttempt(ctx, req.Email, ipAddress)

	return &AuthResponse{
		User: &UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Role:          user.Role,
			SchoolID:      user.SchoolID,
			EmailVerified: user.EmailVerified,
			MFAEnabled:    user.MFAEnabled,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.tokenManager.GetAccessTokenTTL(),
	}, nil
}

// RefreshToken generates new access token from refresh token.
// It also verifies that the user account is still active before issuing
// a new token — prevents disabled accounts from silently regaining access.
func (s *Service) RefreshToken(ctx context.Context, refreshToken, ipAddress, userAgent string) (*TokenPair, error) {
	// Get refresh token from DB
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// Check if revoked
	if rt.Revoked {
		return nil, ErrTokenRevoked
	}

	// Check if expired
	if time.Now().After(rt.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	// Verify account is still active before issuing a new access token.
	// A disabled account must not be able to obtain new tokens via refresh.
	if !user.IsActive {
		// Revoke the refresh token so it cannot be retried.
		_ = s.repo.RevokeRefreshToken(ctx, rt.ID)
		return nil, ErrUserInactive
	}

	// Generate new access token
	accessToken, err := s.tokenManager.GenerateAccessToken(
		user.ID, user.Email, user.Role,
		func() string {
			if user.SchoolID != nil {
				return *user.SchoolID
			}
			return ""
		}(),
	)
	if err != nil {
		return nil, err
	}

	// Revoke the old refresh token
	if err := s.repo.RevokeRefreshToken(ctx, rt.ID); err != nil {
		return nil, err
	}

	// Generate a new refresh token
	newRefreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Update IP and UserAgent from current request if provided, falling back to previous token's values
	newIp := ipAddress
	if newIp == "" {
		newIp = rt.IPAddress
	}
	newUserAgent := userAgent
	if newUserAgent == "" {
		newUserAgent = rt.UserAgent
	}

	// Store the new refresh token in DB
	newRt := &RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IPAddress: newIp,
		UserAgent: newUserAgent,
	}
	if err := s.repo.CreateRefreshToken(ctx, newRt); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.tokenManager.GetAccessTokenTTL(),
	}, nil
}

// Logout revokes ALL refresh tokens for the user.
// It verifies that the refresh token presented belongs to the authenticated
// user (from the JWT access token in context) to prevent cross-user logout.
func (s *Service) Logout(ctx context.Context, refreshToken string, callerUserID string) error {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		// Token not found or already revoked — still a valid logout intent.
		return nil
	}
	// Ownership check: refuse to revoke another user's sessions.
	if rt.UserID != callerUserID {
		return ErrUnauthorized
	}
	// Revoke ALL sessions for this user for maximum security.
	return s.repo.RevokeAllUserTokens(ctx, rt.UserID)
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// SetupMFA initiates MFA setup for a user.
// The TOTP secret is encrypted with AES-256-GCM before being stored so that
// a database compromise does not expose raw secrets.
func (s *Service) SetupMFA(ctx context.Context, userID string) (*MFASetupResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.MFAEnabled {
		return nil, ErrMFAAlreadyEnabled
	}

	// Generate plaintext secret
	secret, err := s.mfaService.GenerateSecret()
	if err != nil {
		return nil, err
	}

	// Generate QR code URL using the plaintext secret
	qrURL := s.mfaService.GenerateQRCodeURL(user.Email, secret)

	// Generate recovery codes (plaintext — shown once to the user)
	plainRecoveryCodes, err := s.mfaService.GenerateRecoveryCodes(10)
	if err != nil {
		return nil, err
	}

	// Hash each recovery code before persisting.
	// bcrypt is used so that a DB compromise does not leak usable codes.
	hashedCodes := make([]string, len(plainRecoveryCodes))
	for i, code := range plainRecoveryCodes {
		h, err := bcrypt.GenerateFromPassword([]byte(code), s.bcryptCost)
		if err != nil {
			return nil, err
		}
		hashedCodes[i] = string(h)
	}
	// Encrypt and store the TOTP secret FIRST, before saving recovery codes.
	// This prevents orphan recovery codes if secret encryption or saving fails.
	encryptedSecret, err := crypto.EncryptString(secret)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveTempMFASecret(ctx, userID, encryptedSecret); err != nil {
		return nil, err
	}

	if err := s.repo.CreateRecoveryCodes(ctx, userID, hashedCodes); err != nil {
		return nil, err
	}

	return &MFASetupResponse{
		Secret:        secret,             // plaintext shown once to the user for manual entry
		QRCodeURL:     qrURL,
		RecoveryCodes: plainRecoveryCodes, // shown once; hashed copy already in DB
	}, nil
}

// VerifyMFA verifies MFA token and enables it on successful confirmation.
// The stored encrypted secret is decrypted before TOTP validation.
func (s *Service) VerifyMFA(ctx context.Context, userID, token string) error {
	encryptedSecret, err := s.repo.GetMFASecret(ctx, userID)
	if err != nil {
		return err
	}

	secret, err := crypto.DecryptString(encryptedSecret)
	if err != nil {
		return err
	}

	if !s.mfaService.VerifyTOTPUser(userID, secret, token) {
		return ErrInvalidMFAToken
	}

	// Token is correct, enable MFA officially
	return s.repo.ConfirmMFA(ctx, userID)
}

// RequestPasswordReset creates a password reset token.
// To prevent email flooding attacks, at most 3 reset requests are allowed
// per email address in a 15-minute window.
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal whether the email exists — always return success.
		return nil
	}

	// Rate limit: max 3 password reset requests per email in 15 minutes.
	since15m := time.Now().Add(-15 * time.Minute)
	recentResets, err := s.repo.GetRecentPasswordResets(ctx, user.ID, since15m)
	if err != nil {
		return err
	}
	if recentResets >= 3 {
		logger.Log.Warnf("SECURITY AUDIT: Password reset rate limit hit for email %s (user %s)", email, user.ID)
		// Silently drop the request — do not reveal the throttle to the caller.
		return nil
	}

	// Generate reset token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}
	token := hex.EncodeToString(tokenBytes)

	// Store token
	prt := &PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if err := s.repo.CreatePasswordResetToken(ctx, prt); err != nil {
		return err
	}

	// Log the reset request (no sensitive data in log)
	logger.Log.Infof("PASSWORD RESET REQUEST for user %s", user.ID)

	return nil
}

// ResetPassword resets user password with token
func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate new password
	passwordValidator := NewPasswordValidator()
	if err := passwordValidator.Validate(newPassword); err != nil {
		return err
	}

	// Get reset token
	prt, err := s.repo.GetPasswordResetToken(ctx, token)
	if err != nil {
		return err
	}

	// Check if already used
	if prt.Used {
		return ErrInvalidToken
	}

	// Check if expired
	if time.Now().After(prt.ExpiresAt) {
		return ErrInvalidToken
	}

	// Check password history (prevent reuse of last 5)
	history, err := s.repo.GetPasswordHistory(ctx, prt.UserID)
	if err != nil {
		return fmt.Errorf("failed to fetch password history: %w", err)
	}
	for _, oldHash := range history {
		if bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(newPassword)) == nil {
			return ErrPasswordReused
		}
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.bcryptCost)
	if err != nil {
		return err
	}

	// Mark token as used FIRST to prevent race conditions or double execution
	if err := s.repo.UsePasswordResetToken(ctx, prt.ID); err != nil {
		return err
	}

	// Update password
	if err := s.repo.UpdatePassword(ctx, prt.UserID, string(passwordHash)); err != nil {
		return err
	}

	// Add to password history
	_ = s.repo.AddPasswordHistory(ctx, prt.UserID, string(passwordHash))

	// Revoke all refresh tokens for security
	_ = s.repo.RevokeAllUserTokens(ctx, prt.UserID)

	return nil
}

func (s *Service) recordFailedAttempt(ctx context.Context, email, ipAddress string) {
	attempt := &LoginAttempt{
		Email:       email,
		IPAddress:   ipAddress,
		Success:     false,
		AttemptedAt: time.Now(),
	}
	_ = s.repo.RecordLoginAttempt(ctx, attempt)
}

func (s *Service) recordSuccessfulAttempt(ctx context.Context, email, ipAddress string) {
	attempt := &LoginAttempt{
		Email:       email,
		IPAddress:   ipAddress,
		Success:     true,
		AttemptedAt: time.Now(),
	}
	_ = s.repo.RecordLoginAttempt(ctx, attempt)
}
