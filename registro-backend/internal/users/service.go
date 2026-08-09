package users

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

// allowedCreators maps each role to the set of roles it is allowed to create.
var allowedCreators = map[string]map[string]bool{
	"superadmin": {"superadmin": true, "admin": true, "secretary": true, "teacher": true, "student": true, "parent": true},
	"admin":      {"secretary": true, "teacher": true, "student": true, "parent": true},
	"secretary":  {"student": true, "parent": true},
}

// Service handles business logic for user management.
type Service struct {
	repo Repository
}

// NewService creates a new Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ─── Auth helpers ────────────────────────────────────────────────────────────

func isPrivileged(role string) bool {
	return role == "admin" || role == "superadmin" || role == "secretary"
}

// ─── CRUD ────────────────────────────────────────────────────────────────────

// CreateUser creates a new user, enforcing role-based permission checks.
func (s *Service) CreateUser(ctx context.Context, actorRole string, req CreateUserRequest) (*User, error) {
	allowed, ok := allowedCreators[actorRole]
	if !ok {
		return nil, ErrUnauthorized
	}
	if !allowed[req.Role] {
		return nil, ErrUnauthorized
	}

	if err := validatePasswordComplexity(req.Password); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	var fiscalCode *string
	if req.FiscalCode != "" {
		fc := req.FiscalCode
		fiscalCode = &fc
	}
	var phoneNumber *string
	if req.PhoneNumber != "" {
		ph := req.PhoneNumber
		phoneNumber = &ph
	}
	var jobTitle *string
	if req.JobTitle != "" {
		jt := req.JobTitle
		jobTitle = &jt
	}

	now := time.Now()
	user := &User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		FiscalCode:   fiscalCode,
		Role:         req.Role,
		SchoolID:     req.SchoolID,
		ClassID:      req.ClassID,
		PhoneNumber:  phoneNumber,
		JobTitle:     jobTitle,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// ListUsers returns a paginated, filtered list of users.
func (s *Service) ListUsers(ctx context.Context, actorRole, actorSchoolID string, filter UserFilter) ([]User, int, error) {
	if !isPrivileged(actorRole) && actorRole != "teacher" && actorRole != "principal" && actorRole != "vice_principal" {
		return nil, 0, ErrUnauthorized
	}
	if actorRole != "superadmin" && actorSchoolID != "" {
		filter.SchoolID = &actorSchoolID
	}
	return s.repo.List(ctx, filter)
}

// GetUser returns a single user by ID.
func (s *Service) GetUser(ctx context.Context, actorRole, actorSchoolID string, id string) (*User, error) {
	if !isPrivileged(actorRole) && actorRole != "teacher" {
		return nil, ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if actorRole != "superadmin" && actorSchoolID != "" && user.SchoolID != nil && *user.SchoolID != actorSchoolID {
		return nil, ErrUnauthorized
	}
	return user, nil
}

// UpdateUser updates an existing user's fields.
func (s *Service) UpdateUser(ctx context.Context, actorRole, actorSchoolID, id string, req UpdateUserRequest) (*User, error) {
	if !isPrivileged(actorRole) {
		return nil, ErrUnauthorized
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if actorRole != "superadmin" && actorSchoolID != "" && user.SchoolID != nil && *user.SchoolID != actorSchoolID {
		return nil, ErrUnauthorized
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Role != nil && *req.Role != user.Role {
		allowed, ok := allowedCreators[actorRole]
		if !ok || !allowed[*req.Role] {
			return nil, ErrUnauthorized
		}
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.JobTitle != nil {
		user.JobTitle = req.JobTitle
	}
	if req.ClassID != nil {
		user.ClassID = req.ClassID
	}
	if req.SchoolID != nil {
		if actorRole == "superadmin" || (actorSchoolID != "" && *req.SchoolID == actorSchoolID) || (user.SchoolID != nil && *req.SchoolID == *user.SchoolID) {
			user.SchoolID = req.SchoolID
		} else {
			return nil, errors.New("forbidden: only superadmin can modify user school_id")
		}
	}
	if req.FiscalCode != nil {
		user.FiscalCode = req.FiscalCode
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser soft-deletes a user.
func (s *Service) DeleteUser(ctx context.Context, actorRole, actorSchoolID, id string) error {
	if !isPrivileged(actorRole) {
		return ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if actorRole != "superadmin" && actorSchoolID != "" && user.SchoolID != nil && *user.SchoolID != actorSchoolID {
		return ErrUnauthorized
	}
	return s.repo.Delete(ctx, id)
}

// BulkDeleteUsers soft-deletes multiple users at once.
func (s *Service) BulkDeleteUsers(ctx context.Context, actorRole, actorSchoolID string, ids []string) (int, error) {
	if !isPrivileged(actorRole) {
		return 0, ErrUnauthorized
	}
	if actorRole == "admin" || actorRole == "secretary" {
		targetUsers, err := s.repo.ListByIDs(ctx, ids)
		if err != nil {
			return 0, fmt.Errorf("failed to verify target users: %w", err)
		}
		var safeIDs []string
		for _, tu := range targetUsers {
			if tu.Role != "admin" && tu.Role != "superadmin" {
				if actorSchoolID == "" || (tu.SchoolID != nil && *tu.SchoolID == actorSchoolID) {
					safeIDs = append(safeIDs, tu.ID)
				}
			}
		}
		ids = safeIDs
	}
	if len(ids) == 0 {
		return 0, nil
	}
	return s.repo.BulkDelete(ctx, ids)
}

// RestoreUser un-deletes a soft-deleted user.
func (s *Service) RestoreUser(ctx context.Context, actorRole string, id string) error {
	if !isPrivileged(actorRole) {
		return ErrUnauthorized
	}
	return s.repo.Restore(ctx, id)
}

// ─── Password ────────────────────────────────────────────────────────────────

// ChangePassword allows a user to change their own password.
// It enforces that the new password differs from the last 5 used passwords.
func (s *Service) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := validatePasswordComplexity(req.NewPassword); err != nil {
		return err
	}

	// Field is CurrentPassword in ChangePasswordRequest DTO
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Check password history
	history, err := s.repo.GetPasswordHistory(ctx, userID)
	if err != nil {
		return err
	}
	for _, old := range history {
		if bcrypt.CompareHashAndPassword([]byte(old), []byte(req.NewPassword)) == nil {
			return errors.New("new password must not match any of the last 5 passwords")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	user.PasswordHash = string(hash)
	user.PasswordChangedAt = &now
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}
	if err := s.repo.AddPasswordHistory(ctx, userID, string(hash)); err != nil {
		return err
	}
	_ = s.repo.RevokeAllUserTokens(ctx, userID)
	return nil
}

func validatePasswordComplexity(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 128 {
		return errors.New("password is too long (max 128 characters)")
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}
	return nil
}

// ResetPassword allows an admin to force-reset a user's password.
func (s *Service) ResetPassword(ctx context.Context, actorRole, actorSchoolID, userID, newPassword string) error {
	if actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if actorRole != "superadmin" && actorSchoolID != "" && user.SchoolID != nil && *user.SchoolID != actorSchoolID {
		return ErrUnauthorized
	}

	if err := validatePasswordComplexity(newPassword); err != nil {
		return err
	}

	history, err := s.repo.GetPasswordHistory(ctx, userID)
	if err == nil {
		for _, old := range history {
			if bcrypt.CompareHashAndPassword([]byte(old), []byte(newPassword)) == nil {
				return errors.New("new password must not match any of the last 5 passwords")
			}
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	user.PasswordHash = string(hash)
	user.PasswordChangedAt = &now
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}
	return s.repo.AddPasswordHistory(ctx, userID, string(hash))
}

// DisableMFA disables MFA for the given user (admin only).
func (s *Service) DisableMFA(ctx context.Context, actorRole string, userID string) error {
	if actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	user.MFAEnabled = false
	user.MFASecret = ""
	return s.repo.Update(ctx, user)
}

// ─── Bulk Import ─────────────────────────────────────────────────────────────

// BulkImport imports users from a CSV or XLSX file.
// Returns an ImportResult (the type defined in dto.go).
func (s *Service) BulkImport(ctx context.Context, actorRole, actorSchoolID string, file multipart.File, filename string) (*ImportResult, error) {
	if !isPrivileged(actorRole) {
		return nil, ErrUnauthorized
	}

	// NewImporter takes no arguments
	importer := NewImporter()
	users, parseErrors, err := importer.ParseFile(file, filename)
	if err != nil {
		return nil, err
	}

	if actorRole != "superadmin" && actorSchoolID != "" {
		for i := range users {
			scID := actorSchoolID
			users[i].SchoolID = &scID
		}
	}

	result := &ImportResult{
		Total:  len(users) + len(parseErrors),
		Failed: len(parseErrors),
		Errors: parseErrors,
	}

	if len(users) > 0 {
		count, bulkErrors, err := s.repo.BulkCreate(ctx, users)
		result.Created = count
		result.Failed += len(bulkErrors)
		result.Errors = append(result.Errors, bulkErrors...)
		if err != nil && len(bulkErrors) == 0 {
			result.Errors = append(result.Errors, err.Error())
			result.Failed += len(users)
		}
	}

	return result, nil
}

// ─── GDPR ────────────────────────────────────────────────────────────────────

// GDPRDataExport exports all personal data for a user (GDPR right of access).
func (s *Service) GDPRDataExport(ctx context.Context, actorID, actorRole, targetID string) (map[string]interface{}, error) {
	if actorID != targetID && actorRole != "admin" && actorRole != "superadmin" {
		return nil, ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}
	// Fetch audit logs for the export (max 1000, offset 0)
	logs, _, err := s.repo.GetAuditLogs(ctx, targetID, 1000, 0)
	if err != nil {
		return nil, err
	}
	// NewExporter takes no arguments; GenerateDataExport is on GDPRHandler
	gdpr := NewGDPRHandler(s.repo)
	return gdpr.GenerateDataExport(user, logs), nil
}

// GDPRDelete pseudonymizes a user (GDPR right to erasure).
func (s *Service) GDPRDelete(ctx context.Context, actorRole string, targetID string) error {
	if actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return err
	}
	gdpr := NewGDPRHandler(s.repo)
	gdpr.PseudonymizeUser(user)
	// Persist the pseudonymized data
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}
	_ = s.repo.RevokeAllUserTokens(ctx, targetID)
	return nil
}

// ─── Audit ───────────────────────────────────────────────────────────────────

// GetAuditLogs returns the audit log for a specific user (paginated).
// Access: the user themselves, or admin/superadmin.
func (s *Service) GetAuditLogs(ctx context.Context, actorID, actorRole, targetID string, limit, offset int) ([]AuditLog, int, error) {
	if actorID != targetID && actorRole != "admin" && actorRole != "superadmin" {
		return nil, 0, ErrUnauthorized
	}
	return s.repo.GetAuditLogs(ctx, targetID, limit, offset)
}

// ─── Relationships ───────────────────────────────────────────────────────────

// GetChildren returns the list of students linked to a parent.
func (s *Service) GetChildren(ctx context.Context, actorRole, actorID, parentUserID string) ([]StudentChild, error) {
	if actorID != parentUserID && !isPrivileged(actorRole) {
		return nil, ErrUnauthorized
	}
	return s.repo.GetChildren(ctx, parentUserID)
}

// GetGuardians returns all guardians linked to a student.
func (s *Service) GetGuardians(ctx context.Context, actorRole string, studentUserID string) ([]GuardianInfo, error) {
	if !isPrivileged(actorRole) {
		return nil, ErrUnauthorized
	}
	studentProfileID, err := s.repo.GetStudentProfile(ctx, studentUserID)
	if err != nil {
		return nil, fmt.Errorf("student profile not found: %w", err)
	}
	return s.repo.GetGuardians(ctx, studentProfileID)
}

// AddGuardian links a parent to a student.
func (s *Service) AddGuardian(ctx context.Context, actorRole, studentUserID, parentUserID, relationship string) error {
	if !isPrivileged(actorRole) {
		return ErrUnauthorized
	}
	studentProfileID, err := s.repo.GetStudentProfile(ctx, studentUserID)
	if err != nil {
		return fmt.Errorf("student profile not found: %w", err)
	}
	parentProfileID, err := s.repo.GetParentProfile(ctx, parentUserID)
	if err != nil {
		return fmt.Errorf("parent profile not found: %w", err)
	}
	return s.repo.AddGuardian(ctx, studentProfileID, parentProfileID, relationship)
}

// RemoveGuardian unlinks a parent from a student.
func (s *Service) RemoveGuardian(ctx context.Context, actorRole, studentUserID, parentUserID string) error {
	if !isPrivileged(actorRole) {
		return ErrUnauthorized
	}
	studentProfileID, err := s.repo.GetStudentProfile(ctx, studentUserID)
	if err != nil {
		return fmt.Errorf("student profile not found: %w", err)
	}
	parentProfileID, err := s.repo.GetParentProfile(ctx, parentUserID)
	if err != nil {
		return fmt.Errorf("parent profile not found: %w", err)
	}
	return s.repo.RemoveGuardian(ctx, studentProfileID, parentProfileID)
}

// IsGuardian checks if parentUserID is a guardian for studentUserID.
func (s *Service) IsGuardian(ctx context.Context, parentUserID, studentUserID string) (bool, error) {
	return s.repo.IsGuardian(ctx, parentUserID, studentUserID)
}

// SwitchChildContext verifies that parentUserID is guardian of studentUserID and returns the target student's profile context.
func (s *Service) SwitchChildContext(ctx context.Context, parentUserID, targetStudentID string) (*User, error) {
	isGuard, err := s.repo.IsGuardian(ctx, parentUserID, targetStudentID)
	if err != nil {
		return nil, err
	}
	if !isGuard {
		return nil, errors.New("forbidden: target user is not a child of this parent")
	}
	return s.repo.GetByID(ctx, targetStudentID)
}

func sanitizeCSV(s string) string {
	if len(s) > 0 && (s[0] == '=' || s[0] == '+' || s[0] == '-' || s[0] == '@' || s[0] == '\t' || s[0] == '\r') {
		return "'" + s
	}
	return s
}

type UserExportDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	SchoolID  *string   `json:"school_id,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// ExportUsers exports users matching filter in specified format ("csv" or "json").
func (s *Service) ExportUsers(ctx context.Context, actorRole, actorSchoolID string, filter UserFilter, format string) ([]byte, error) {
	if !isPrivileged(actorRole) {
		return nil, ErrUnauthorized
	}
	if actorRole != "superadmin" && actorSchoolID != "" {
		filter.SchoolID = &actorSchoolID
	}
	users, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(format) == "csv" {
		var buf bytes.Buffer
		w := csv.NewWriter(&buf)
		_ = w.Write([]string{"ID", "Email", "FirstName", "LastName", "Role"})
		for _, u := range users {
			_ = w.Write([]string{
				sanitizeCSV(u.ID),
				sanitizeCSV(u.Email),
				sanitizeCSV(u.FirstName),
				sanitizeCSV(u.LastName),
				sanitizeCSV(u.Role),
			})
		}
		w.Flush()
		return buf.Bytes(), nil
	}

	var dtos []UserExportDTO
	for _, u := range users {
		dtos = append(dtos, UserExportDTO{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
			SchoolID:  u.SchoolID,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt,
		})
	}
	return json.Marshal(dtos)
}

func (s *Service) GetStudentFascicolo(ctx context.Context, actorID, actorRole, studentID string) (*StudentFascicolo, error) {
	if !isPrivileged(actorRole) {
		if actorRole == "parent" {
			isGuard, err := s.repo.IsGuardian(ctx, actorID, studentID)
			if err != nil || !isGuard {
				return nil, ErrUnauthorized
			}
		} else if actorID != studentID {
			return nil, ErrUnauthorized
		}
	}

	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	var guardians []GuardianInfo
	if studentProfileID, err := s.repo.GetStudentProfile(ctx, studentID); err == nil && studentProfileID != "" {
		guardians, _ = s.repo.GetGuardians(ctx, studentProfileID)
	}

	summary, err := s.repo.GetFascicoloSummary(ctx, studentID, student.IsActive)
	if err != nil {
		statusStr := "Inactive"
		if student.IsActive {
			statusStr = "Active"
		}
		summary = map[string]interface{}{
			"status":          statusStr,
			"documents_count": 0,
			"notes_count":     0,
			"pcto_hours":      0,
		}
	}

	return &StudentFascicolo{
		Student:     student,
		Guardians:   guardians,
		Summary:     summary,
		GeneratedAt: time.Now(),
	}, nil
}
