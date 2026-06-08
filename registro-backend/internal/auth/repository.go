package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Repository handles auth-related database operations
type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	UpdateLastLogin(ctx context.Context, userID string) error

	// MFA operations
	EnableMFA(ctx context.Context, userID, secret string) error
	DisableMFA(ctx context.Context, userID string) error
	GetMFASecret(ctx context.Context, userID string) (string, error)

	// Recovery codes
	CreateRecoveryCodes(ctx context.Context, userID string, codes []string) error
	GetRecoveryCodes(ctx context.Context, userID string) ([]*MFARecoveryCode, error)
	UseRecoveryCode(ctx context.Context, userID, code string) error

	// Refresh tokens
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error

	// Password reset
	CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error)
	UsePasswordResetToken(ctx context.Context, tokenID string) error
	UpdatePassword(ctx context.Context, userID, passwordHash string) error

	// Rate limiting
	RecordLoginAttempt(ctx context.Context, attempt *LoginAttempt) error
	GetRecentLoginAttempts(ctx context.Context, email, ipAddress string, since time.Time) (int, error)
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, is_active, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true
	user.EmailVerified = false

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Role, user.SchoolID, user.IsActive, user.EmailVerified,
		user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, school_id, 
		       is_active, email_verified, mfa_enabled, mfa_secret, created_at, updated_at, last_login
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	user := &User{}
	var mfaSecret sql.NullString
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.SchoolID, &user.IsActive, &user.EmailVerified,
		&user.MFAEnabled, &mfaSecret, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user.MFASecret = mfaSecret.String
	return user, nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, school_id,
		       is_active, email_verified, mfa_enabled, mfa_secret, created_at, updated_at, last_login
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	user := &User{}
	var mfaSecret sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.SchoolID, &user.IsActive, &user.EmailVerified,
		&user.MFAEnabled, &mfaSecret, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user.MFASecret = mfaSecret.String
	return user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User) error {
	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, updated_at = $3
		WHERE id = $4
	`
	user.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.UpdatedAt, user.ID)
	return err
}

func (r *repository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `UPDATE users SET last_login = $1 WHERE id = $2`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, userID)
	return err
}

func (r *repository) EnableMFA(ctx context.Context, userID, secret string) error {
	query := `UPDATE users SET mfa_enabled = true, mfa_secret = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, secret, userID)
	return err
}

func (r *repository) DisableMFA(ctx context.Context, userID string) error {
	query := `UPDATE users SET mfa_enabled = false, mfa_secret = NULL WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *repository) GetMFASecret(ctx context.Context, userID string) (string, error) {
	query := `SELECT mfa_secret FROM users WHERE id = $1 AND mfa_enabled = true`
	var secret string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&secret)
	if err == sql.ErrNoRows {
		return "", ErrMFANotEnabled
	}
	return secret, err
}

func (r *repository) CreateRecoveryCodes(ctx context.Context, userID string, codes []string) error {
	query := `INSERT INTO mfa_recovery_codes (id, user_id, code, used, created_at) VALUES ($1, $2, $3, false, $4)`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	for _, code := range codes {
		_, err := tx.ExecContext(ctx, query, uuid.New().String(), userID, code, time.Now())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *repository) GetRecoveryCodes(ctx context.Context, userID string) ([]*MFARecoveryCode, error) {
	query := `SELECT id, user_id, code, used, used_at, created_at FROM mfa_recovery_codes WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []*MFARecoveryCode
	for rows.Next() {
		code := &MFARecoveryCode{}
		if err := rows.Scan(&code.ID, &code.UserID, &code.Code, &code.Used, &code.UsedAt, &code.CreatedAt); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, rows.Err()
}

func (r *repository) UseRecoveryCode(ctx context.Context, userID, code string) error {
	query := `UPDATE mfa_recovery_codes SET used = true, used_at = $1 WHERE user_id = $2 AND code = $3 AND used = false`
	result, err := r.db.ExecContext(ctx, query, time.Now(), userID, code)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return ErrInvalidMFAToken
	}
	return nil
}

func (r *repository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, revoked, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	token.ID = uuid.New().String()
	token.CreatedAt = time.Now()
	token.Revoked = false

	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, token.Token, token.ExpiresAt,
		token.CreatedAt, token.Revoked, token.IPAddress, token.UserAgent,
	)
	return err
}

func (r *repository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, revoked, ip_address, user_agent
		FROM refresh_tokens
		WHERE token = $1
	`
	rt := &RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt,
		&rt.CreatedAt, &rt.Revoked, &rt.IPAddress, &rt.UserAgent,
	)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidToken
	}
	return rt, err
}

func (r *repository) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	return err
}

func (r *repository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *repository) CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (id, user_id, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, false, $5)
	`
	token.ID = uuid.New().String()
	token.CreatedAt = time.Now()
	token.Used = false

	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt,
	)
	return err
}

func (r *repository) GetPasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token = $1
	`
	prt := &PasswordResetToken{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&prt.ID, &prt.UserID, &prt.Token, &prt.ExpiresAt, &prt.Used, &prt.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidToken
	}
	return prt, err
}

func (r *repository) UsePasswordResetToken(ctx context.Context, tokenID string) error {
	query := `UPDATE password_reset_tokens SET used = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	return err
}

func (r *repository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, passwordHash, time.Now(), userID)
	return err
}

func (r *repository) RecordLoginAttempt(ctx context.Context, attempt *LoginAttempt) error {
	query := `
		INSERT INTO login_attempts (email, ip_address, success, attempted_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, attempt.Email, attempt.IPAddress, attempt.Success, attempt.AttemptedAt)
	return err
}

func (r *repository) GetRecentLoginAttempts(ctx context.Context, email, ipAddress string, since time.Time) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM login_attempts
		WHERE email = $1 AND ip_address = $2 AND success = false AND attempted_at > $3
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, email, ipAddress, since).Scan(&count)
	return count, err
}
