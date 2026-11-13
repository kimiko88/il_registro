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

	"registro-backend/pkg/crypto"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInvalidRole  = errors.New("invalid role")
)

// validRoles is the exhaustive set of registered system roles.
var validRoles = map[string]bool{
	"superadmin": true, "admin": true, "principal": true, "vice_principal": true,
	"secretary": true, "teacher": true, "student": true, "parent": true,
	"system_auditor": true, "dsga": true, "collaboratore_ds": true, "collaboratore_scolastico": true,
	"assistente_amministrativo": true, "assistente_alunni": true, "assistente_personale": true,
	"assistente_contabilita": true, "assistente_protocollo": true, "assistente_sportello": true,
	"assistente_tecnico": true, "responsabile_servizio": true,
	"responsabile_gestione_documentale": true, "responsabile_conservazione": true, "dpo": true,
	"coordinator": true, "coordinatore_classe": true, "segretario_consiglio": true,
	"referente_progetto": true, "referente_inclusione": true, "responsabile_dipartimento": true,
	"tutor_orientatore": true, "animatore_digitale": true,
}

// IsValidRole returns true if the provided role string is a registered valid role.
func IsValidRole(role string) bool {
	return validRoles[role]
}

// bcryptCost è il cost factor usato per tutti gli hash bcrypt nel package users.
// Usare cost 12 invece di bcrypt.DefaultCost (10) per maggiore resistenza al brute-force.
const bcryptCost = 12

// allowedCreators maps each role to the set of roles it is allowed to create.
var allowedCreators = map[string]map[string]bool{
	"superadmin": {
		"superadmin": true, "admin": true, "principal": true, "vice_principal": true,
		"secretary": true, "teacher": true, "student": true, "parent": true,
		"dsga": true, "collaboratore_ds": true, "collaboratore_scolastico": true,
		"assistente_amministrativo": true, "assistente_alunni": true, "assistente_personale": true,
		"assistente_contabilita": true, "assistente_protocollo": true, "assistente_sportello": true,
		"assistente_tecnico": true, "responsabile_servizio": true,
		"responsabile_gestione_documentale": true, "responsabile_conservazione": true, "dpo": true,
		"coordinator": true, "coordinatore_classe": true, "system_auditor": true,
	},
	"admin": {
		"admin": true, "principal": true, "vice_principal": true,
		"secretary": true, "teacher": true, "student": true, "parent": true,
		"dsga": true, "collaboratore_ds": true, "collaboratore_scolastico": true,
		"assistente_amministrativo": true, "assistente_alunni": true, "assistente_personale": true,
		"assistente_contabilita": true, "assistente_protocollo": true, "assistente_sportello": true,
		"assistente_tecnico": true, "responsabile_servizio": true,
		"responsabile_gestione_documentale": true, "responsabile_conservazione": true, "dpo": true,
		"coordinator": true, "coordinatore_classe": true,
	},
	"principal": {
		"vice_principal": true, "secretary": true, "teacher": true, "student": true, "parent": true,
		"dsga": true, "collaboratore_ds": true, "collaboratore_scolastico": true,
		"assistente_amministrativo": true, "assistente_alunni": true, "assistente_personale": true,
		"assistente_contabilita": true, "assistente_protocollo": true, "assistente_sportello": true,
		"assistente_tecnico": true, "responsabile_servizio": true,
		"responsabile_gestione_documentale": true, "responsabile_conservazione": true, "dpo": true,
		"coordinator": true, "coordinatore_classe": true,
	},
	"secretary": {
		"teacher": true, "student": true, "parent": true, "collaboratore_ds": true,
		"assistente_amministrativo": true, "assistente_alunni": true, "assistente_personale": true,
		"assistente_contabilita": true, "assistente_protocollo": true, "assistente_sportello": true,
		"collaboratore_scolastico": true, "assistente_tecnico": true, "coordinator": true,
	},
	"dsga": {
		"assistente_amministrativo": true, "assistente_alunni": true, "assistente_personale": true,
		"assistente_contabilita": true, "assistente_protocollo": true, "assistente_sportello": true,
		"collaboratore_scolastico": true, "assistente_tecnico": true, "responsabile_servizio": true,
	},
	"assistente_amministrativo": {
		"student": true, "parent": true,
	},
	"assistente_alunni": {
		"student": true, "parent": true,
	},
	"assistente_personale": {
		"teacher": true, "collaboratore_scolastico": true, "assistente_tecnico": true,
	},
}

// Service handles business logic for user management.
type Service struct {
	repo Repository
}

// NewService creates a new Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// hasSchoolAccess verifies that a privileged actor can operate on a target user.
// Superadmins are global; every other actor must have an explicit matching tenant.
func hasSchoolAccess(actorRole, actorSchoolID string, target *User) bool {
	if actorRole == "superadmin" {
		return true
	}
	return actorSchoolID != "" && target != nil && target.SchoolID != nil && *target.SchoolID == actorSchoolID
}

func actorSchoolScope(scopes []string) string {
	if len(scopes) == 0 {
		return ""
	}
	return scopes[0]
}

// ─── Auth helpers ────────────────────────────────────────────────────────────

func isPrivileged(role string) bool {
	switch role {
	case "admin", "superadmin", "principal", "vice_principal", "secretary", "dsga",
		"assistente_amministrativo", "assistente_alunni", "assistente_personale",
		"assistente_contabilita", "assistente_protocollo", "assistente_sportello":
		return true
	default:
		return false
	}
}

// ─── CRUD ────────────────────────────────────────────────────────────────────

// CreateUserForSchool prevents non-superadmin callers from selecting or omitting
// the tenant. The HTTP handler must pass the authenticated school from context.
func (s *Service) CreateUserForSchool(ctx context.Context, actorRole, actorSchoolID string, req CreateUserRequest) (*User, error) {
	if actorRole != "superadmin" {
		if actorSchoolID == "" {
			return nil, ErrUnauthorized
		}
		if req.SchoolID != nil && *req.SchoolID != actorSchoolID {
			return nil, ErrUnauthorized
		}
		schoolID := actorSchoolID
		req.SchoolID = &schoolID
	}
	return s.CreateUser(ctx, actorRole, req)
}

// CreateUser creates a new user, enforcing role-based permission checks.
func (s *Service) CreateUser(ctx context.Context, actorRole string, req CreateUserRequest) (*User, error) {
	allowed, ok := allowedCreators[actorRole]
	if !ok {
		return nil, ErrUnauthorized
	}
	if !IsValidRole(req.Role) {
		return nil, ErrInvalidRole
	}
	if !allowed[req.Role] {
		return nil, ErrUnauthorized
	}

	if err := validatePasswordComplexity(req.Password); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword(crypto.PrehashPassword(req.Password), bcryptCost)
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
		IsStaff:      req.IsStaff != nil && *req.IsStaff,
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
	if !isPrivileged(actorRole) && actorRole != "teacher" && actorRole != "principal" && actorRole != "vice_principal" && actorRole != "collaboratore_ds" {
		if (actorRole == "parent" || actorRole == "student") && filter.Role == "teacher" {
			// Allowed to query teachers in their school for booking colloqui or communications
		} else if (actorRole == "collaboratore_scolastico" || actorRole == "responsabile_servizio") && (filter.Role == "student" || filter.Role == "teacher") {
			// Allowed to query students (for visitor registry, early exits) or teachers in their school
		} else {
			return nil, 0, ErrUnauthorized
		}
	}
	if actorRole != "superadmin" {
		if actorSchoolID == "" {
			return nil, 0, ErrUnauthorized
		}
		filter.SchoolID = &actorSchoolID
	}
	return s.repo.List(ctx, filter)
}

// GetUser returns a single user by ID.
func (s *Service) GetUser(ctx context.Context, actorRole, actorSchoolID string, id string) (*User, error) {
	if !isPrivileged(actorRole) && actorRole != "teacher" && actorRole != "principal" && actorRole != "vice_principal" && actorRole != "collaboratore_ds" && actorRole != "collaboratore_scolastico" && actorRole != "responsabile_servizio" {
		return nil, ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if actorRole != "superadmin" && (actorSchoolID == "" || user.SchoolID == nil || *user.SchoolID != actorSchoolID) {
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

	if actorRole != "superadmin" && (actorSchoolID == "" || user.SchoolID == nil || *user.SchoolID != actorSchoolID) {
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
		if !ok {
			return nil, ErrUnauthorized
		}
		if !IsValidRole(*req.Role) {
			return nil, ErrInvalidRole
		}
		if !allowed[*req.Role] {
			return nil, ErrUnauthorized
		}
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.IsStaff != nil {
		user.IsStaff = *req.IsStaff
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

	// Changes to role, activation state or tenant invalidate all previously issued sessions.
	if req.Role != nil || req.IsActive != nil || req.SchoolID != nil {
		if err := s.repo.RevokeAllUserTokens(ctx, id); err != nil {
			return nil, fmt.Errorf("failed to revoke sessions after security-sensitive account change: %w", err)
		}
	}
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
	if actorRole != "superadmin" && (actorSchoolID == "" || user.SchoolID == nil || *user.SchoolID != actorSchoolID) {
		return ErrUnauthorized
	}
	if (actorRole == "secretary" || actorRole == "dsga" || actorRole == "assistente_amministrativo") && (user.Role == "admin" || user.Role == "superadmin" || (actorRole != "dsga" && user.Role == "dsga")) {
		return ErrUnauthorized
	}
	if err := s.repo.RevokeAllUserTokens(ctx, id); err != nil {
		return fmt.Errorf("failed to revoke sessions before deleting user: %w", err)
	}
	return s.repo.Delete(ctx, id)
}

// BulkDeleteUsers soft-deletes users after verifying tenant ownership and revoking sessions.
func (s *Service) BulkDeleteUsers(ctx context.Context, actorRole, actorSchoolID string, ids []string) (int, error) {
	if !isPrivileged(actorRole) {
		return 0, ErrUnauthorized
	}
	if len(ids) == 0 {
		return 0, nil
	}

	// Resolve every target before mutation. Never trust caller-supplied IDs as
	// proof of tenant ownership, regardless of the caller's privileged role.
	targetUsers, err := s.repo.ListByIDs(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to verify target users: %w", err)
	}

	safeIDs := make([]string, 0, len(targetUsers))
	seen := make(map[string]struct{}, len(targetUsers))
	for _, target := range targetUsers {
		if target.ID == "" {
			continue
		}
		if _, duplicate := seen[target.ID]; duplicate {
			continue
		}
		if actorRole != "superadmin" {
			if actorSchoolID == "" || target.SchoolID == nil || *target.SchoolID != actorSchoolID {
				continue
			}
		}
		if (actorRole == "admin" || actorRole == "secretary" || actorRole == "dsga" || actorRole == "assistente_amministrativo") &&
			(target.Role == "admin" || target.Role == "superadmin" || (actorRole != "dsga" && target.Role == "dsga")) {
			continue
		}
		seen[target.ID] = struct{}{}
		safeIDs = append(safeIDs, target.ID)
	}
	if len(safeIDs) == 0 {
		return 0, nil
	}

	// Revocation must succeed before the soft-delete. A partial revocation
	// leaves some sessions invalidated but never deletes users with live tokens.
	for _, userID := range safeIDs {
		if err := s.repo.RevokeAllUserTokens(ctx, userID); err != nil {
			return 0, fmt.Errorf("failed to revoke sessions before bulk deleting user %s: %w", userID, err)
		}
	}
	return s.repo.BulkDelete(ctx, safeIDs)
}

// RestoreUser un-deletes a soft-deleted user.
func (s *Service) RestoreUser(ctx context.Context, actorRole, actorSchoolID string, id string) error {
	if !isPrivileged(actorRole) {
		return ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if actorRole != "superadmin" && (actorSchoolID == "" || user.SchoolID == nil || *user.SchoolID != actorSchoolID) {
		return ErrUnauthorized
	}
	return s.repo.Restore(ctx, id)
}

// ─── Password Error Keys ──────────────────────────────────────────────────────

var (
	ErrCurrentPasswordIncorrect = errors.New("ERR_CURRENT_PASSWORD_INCORRECT")
	ErrPasswordComplexity       = errors.New("ERR_PASSWORD_COMPLEXITY")
	ErrPasswordTooShort         = errors.New("ERR_PASSWORD_TOO_SHORT")
	ErrPasswordTooLong          = errors.New("ERR_PASSWORD_TOO_LONG")
	ErrPasswordRecentlyUsed     = errors.New("ERR_PASSWORD_RECENTLY_USED")
)

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

	// Field is CurrentPassword in ChangePasswordRequest DTO (supports pre-hashed and legacy raw)
	preCurrent := crypto.PrehashPassword(req.CurrentPassword)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), preCurrent); err != nil {
		legacyPepperedErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), crypto.PrehashPasswordLegacy(req.CurrentPassword))
		if legacyPepperedErr != nil {
			if !crypto.LegacyPasswordMigrationEnabled() ||
				bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)) != nil {
				return ErrCurrentPasswordIncorrect
			}
		}
	}

	// Check password history
	history, err := s.repo.GetPasswordHistory(ctx, userID)
	if err != nil {
		return err
	}
	preNew := crypto.PrehashPassword(req.NewPassword)
	for _, old := range history {
		if bcrypt.CompareHashAndPassword([]byte(old), preNew) == nil ||
			(crypto.LegacyPasswordMigrationEnabled() &&
				(bcrypt.CompareHashAndPassword([]byte(old), crypto.PrehashPasswordLegacy(req.NewPassword)) == nil ||
					bcrypt.CompareHashAndPassword([]byte(old), []byte(req.NewPassword)) == nil)) {
			return ErrPasswordRecentlyUsed
		}
	}

	hash, err := bcrypt.GenerateFromPassword(preNew, bcryptCost)
	if err != nil {
		return err
	}

	// ChangePasswordTx aggiorna password e history in un'unica transazione atomica.
	// Questo garantisce che se la history fallisce, anche la password non viene salvata.
	if err := s.repo.RevokeAllUserTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke existing sessions: %w", err)
	}
	if err := s.repo.ChangePasswordTx(ctx, userID, string(hash)); err != nil {
		return err
	}
	return nil
}

func validatePasswordComplexity(password string) error {
	if len(password) < 10 {
		return ErrPasswordTooShort
	}
	if len(password) > 128 {
		return ErrPasswordTooLong
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
		return ErrPasswordComplexity
	}
	return nil
}

// ResetPassword allows an admin, superadmin, or secretary to force-reset a user's password.
// Secretary is strictly limited to resetting passwords for teachers, students, and parents.
func (s *Service) ResetPassword(ctx context.Context, actorRole, actorSchoolID, userID, newPassword string) error {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" {
		return ErrUnauthorized
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if actorRole == "secretary" {
		if user.Role != "teacher" && user.Role != "student" && user.Role != "parent" {
			return errors.New("forbidden: la segreteria può resettare solo le password di docenti, studenti e genitori")
		}
	}

	if actorRole != "superadmin" && (actorSchoolID == "" || user.SchoolID == nil || *user.SchoolID != actorSchoolID) {
		return ErrUnauthorized
	}

	if err := validatePasswordComplexity(newPassword); err != nil {
		return err
	}

	preNew := crypto.PrehashPassword(newPassword)
	history, err := s.repo.GetPasswordHistory(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch password history: %w", err)
	}
	{
		for _, old := range history {
			if bcrypt.CompareHashAndPassword([]byte(old), preNew) == nil ||
				(crypto.LegacyPasswordMigrationEnabled() &&
					(bcrypt.CompareHashAndPassword([]byte(old), crypto.PrehashPasswordLegacy(newPassword)) == nil ||
						bcrypt.CompareHashAndPassword([]byte(old), []byte(newPassword)) == nil)) {
				return errors.New("la nuova password non può essere uguale a una delle ultime 5 password utilizzate")
			}
		}
	}

	hash, err := bcrypt.GenerateFromPassword(preNew, bcryptCost)
	if err != nil {
		return err
	}

	now := time.Now()
	user.PasswordHash = string(hash)
	user.PasswordChangedAt = &now
	if err := s.repo.RevokeAllUserTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke existing sessions: %w", err)
	}
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}
	return s.repo.AddPasswordHistory(ctx, userID, string(hash))
}

// DisableMFA disables MFA only for an administrator operating within their own school.
func (s *Service) DisableMFA(ctx context.Context, actorRole string, userID string, actorSchoolIDs ...string) error {
	if actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !hasSchoolAccess(actorRole, actorSchoolScope(actorSchoolIDs), user) {
		return ErrUnauthorized
	}
	if err := s.repo.RevokeAllUserTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke sessions before disabling MFA: %w", err)
	}
	if err := s.repo.ClearTempMFASecret(ctx, userID); err != nil {
		return fmt.Errorf("failed to clear temporary MFA secret: %w", err)
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

	if actorRole != "superadmin" {
		if actorSchoolID == "" {
			return nil, ErrUnauthorized
		}
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

	allowed, ok := allowedCreators[actorRole]

	var validUsers []User
	for _, u := range users {
		if !ok || !allowed[u.Role] {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("utente %s: ruolo %s non consentito per il ruolo %s", u.Email, u.Role, actorRole))
			continue
		}
		if !strings.HasPrefix(u.PasswordHash, "$2a$") && !strings.HasPrefix(u.PasswordHash, "$2b$") {
			if err := validatePasswordComplexity(u.PasswordHash); err != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("utente %s: password non conforme (%v)", u.Email, err))
				continue
			}
			hash, err := bcrypt.GenerateFromPassword(crypto.PrehashPassword(u.PasswordHash), bcryptCost)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("utente %s: errore hashing password", u.Email))
				continue
			}
			u.PasswordHash = string(hash)
		}
		validUsers = append(validUsers, u)
	}

	if len(validUsers) > 0 {
		count, bulkErrors, err := s.repo.BulkCreate(ctx, validUsers)
		result.Created = count
		result.Failed += len(bulkErrors)
		result.Errors = append(result.Errors, bulkErrors...)
		if err != nil && len(bulkErrors) == 0 {
			result.Errors = append(result.Errors, err.Error())
			result.Failed += len(validUsers)
		}
	}

	return result, nil
}

// ─── GDPR ────────────────────────────────────────────────────────────────────

// GDPRDataExport exports personal data for the user or an administrator's own school.
func (s *Service) GDPRDataExport(ctx context.Context, actorID, actorRole, targetID string, actorSchoolIDs ...string) (map[string]interface{}, error) {
	if actorID != targetID && actorRole != "admin" && actorRole != "superadmin" {
		return nil, ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}
	if actorID != targetID && !hasSchoolAccess(actorRole, actorSchoolScope(actorSchoolIDs), user) {
		return nil, ErrUnauthorized
	}
	logs, _, err := s.repo.GetAuditLogs(ctx, targetID, 1000, 0)
	if err != nil {
		return nil, err
	}
	var personalData map[string]interface{}
	if exporter, ok := s.repo.(interface {
		GetUserPersonalDataExport(context.Context, string, string) (map[string]interface{}, error)
	}); ok {
		personalData, err = exporter.GetUserPersonalDataExport(ctx, targetID, user.Role)
		if err != nil { return nil, err }
	} else {
		// Non-Postgres repositories used by isolated tests may not expose domain exports.
		personalData = map[string]interface{}{}
	}
	gdpr := NewGDPRHandler(s.repo)
	return gdpr.GenerateDataExport(user, logs, personalData), nil
}

// GDPRDelete pseudonymizes a user after verifying tenant ownership and revoking sessions.
func (s *Service) GDPRDelete(ctx context.Context, actorRole string, targetID string, actorSchoolIDs ...string) error {
	if actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return err
	}
	if !hasSchoolAccess(actorRole, actorSchoolScope(actorSchoolIDs), user) {
		return ErrUnauthorized
	}
	if err := s.repo.RevokeAllUserTokens(ctx, targetID); err != nil {
		return fmt.Errorf("failed to revoke sessions before GDPR deletion: %w", err)
	}
	gdpr := NewGDPRHandler(s.repo)
	gdpr.PseudonymizeUser(user)
	return s.repo.Update(ctx, user)
}

// ─── Audit ───────────────────────────────────────────────────────────────────

// GetAuditLogs returns audit logs only to the user or an administrator in the same school.
func (s *Service) GetAuditLogs(ctx context.Context, actorID, actorRole, targetID string, limit, offset int, actorSchoolIDs ...string) ([]AuditLog, int, error) {
	if actorID != targetID && actorRole != "admin" && actorRole != "superadmin" {
		return nil, 0, ErrUnauthorized
	}
	if actorID != targetID {
		target, err := s.repo.GetByID(ctx, targetID)
		if err != nil {
			return nil, 0, err
		}
		if !hasSchoolAccess(actorRole, actorSchoolScope(actorSchoolIDs), target) {
			return nil, 0, ErrUnauthorized
		}
	}
	return s.repo.GetAuditLogs(ctx, targetID, limit, offset)
}

// ─── Relationships ───────────────────────────────────────────────────────────

// GetChildren exposes a parent's linked students only to that authenticated parent.
// Staff-facing child searches must use a dedicated school-scoped endpoint.
func (s *Service) GetChildren(ctx context.Context, actorRole, actorID, parentUserID string, _ ...string) ([]StudentChild, error) {
	if actorRole != "parent" || actorID == "" || actorID != parentUserID {
		return nil, ErrUnauthorized
	}
	return s.repo.GetChildren(ctx, parentUserID)
}

func (s *Service) GetGuardians(ctx context.Context, actorRole string, studentUserID string, actorSchoolIDs ...string) ([]GuardianInfo, error) {
	if !isPrivileged(actorRole) {
		return nil, ErrUnauthorized
	}
	student, err := s.repo.GetByID(ctx, studentUserID)
	if err != nil {
		return nil, err
	}
	if !hasSchoolAccess(actorRole, actorSchoolScope(actorSchoolIDs), student) {
		return nil, ErrUnauthorized
	}
	studentProfileID, err := s.repo.GetStudentProfile(ctx, studentUserID)
	if err != nil {
		return nil, fmt.Errorf("student profile not found: %w", err)
	}
	return s.repo.GetGuardians(ctx, studentProfileID)
}

func (s *Service) AddGuardian(ctx context.Context, actorRole, studentUserID, parentUserID, relationship string, actorSchoolIDs ...string) error {
	if !isPrivileged(actorRole) {
		return ErrUnauthorized
	}
	student, err := s.repo.GetByID(ctx, studentUserID)
	if err != nil {
		return err
	}
	parent, err := s.repo.GetByID(ctx, parentUserID)
	if err != nil {
		return err
	}
	scope := actorSchoolScope(actorSchoolIDs)
	if !hasSchoolAccess(actorRole, scope, student) || !hasSchoolAccess(actorRole, scope, parent) ||
		student.SchoolID == nil || parent.SchoolID == nil || *student.SchoolID != *parent.SchoolID {
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

func (s *Service) RemoveGuardian(ctx context.Context, actorRole, studentUserID, parentUserID string, actorSchoolIDs ...string) error {
	if !isPrivileged(actorRole) {
		return ErrUnauthorized
	}
	student, err := s.repo.GetByID(ctx, studentUserID)
	if err != nil {
		return err
	}
	parent, err := s.repo.GetByID(ctx, parentUserID)
	if err != nil {
		return err
	}
	scope := actorSchoolScope(actorSchoolIDs)
	if !hasSchoolAccess(actorRole, scope, student) || !hasSchoolAccess(actorRole, scope, parent) {
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
	if actorRole != "superadmin" {
		if actorSchoolID == "" {
			return nil, ErrUnauthorized
		}
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

func (s *Service) GetStudentFascicolo(ctx context.Context, actorID, actorRole, studentID string, actorSchoolIDs ...string) (*StudentFascicolo, error) {
	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if actorRole != "superadmin" && !hasSchoolAccess(actorRole, actorSchoolScope(actorSchoolIDs), student) {
		return nil, ErrUnauthorized
	}
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

// ─── Incarichi Aggiuntivi & Funzionali ──────────────────────────────────────────

func canAssignDuty(actorRole string, dutyType string) bool {
	switch actorRole {
	case "superadmin", "admin", "principal", "vice_principal", "collaboratore_ds":
		return true
	case "dsga":
		// DSGA governs ATA staff duties and administrative assignments
		return dutyType == "responsabile_servizio" ||
			dutyType == "addetto_sicurezza" ||
			dutyType == "assistente_alunni" ||
			dutyType == "assistente_personale" ||
			dutyType == "assistente_contabilita" ||
			dutyType == "assistente_protocollo" ||
			dutyType == "assistente_sportello"
	case "secretary", "assistente_personale":
		// Secretary / Personnel staff can record class coordinator or meeting secretary assignments
		return dutyType == "coordinatore_classe" ||
			dutyType == "coordinator" ||
			dutyType == "segretario_consiglio"
	default:
		return false
	}
}

// GetUserAssignments returns assignments only to the user or authorized staff in the same school.
func (s *Service) GetUserAssignments(ctx context.Context, actorRole, targetUserID string, actorContext ...string) ([]UserAssignment, error) {
	if !isPrivileged(actorRole) && actorRole != "teacher" && actorRole != "student" && actorRole != "parent" {
		return nil, ErrUnauthorized
	}
	actorID, schoolID := "", ""
	if len(actorContext) > 0 {
		actorID = actorContext[0]
	}
	if len(actorContext) > 1 {
		schoolID = actorContext[1]
	}
	if actorRole != "superadmin" && actorID == "" {
		return nil, ErrUnauthorized
	}
	target, err := s.repo.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	switch {
	case isPrivileged(actorRole), actorRole == "teacher":
		if !hasSchoolAccess(actorRole, schoolID, target) {
			return nil, ErrUnauthorized
		}
	case actorRole == "parent":
		if !hasSchoolAccess(actorRole, schoolID, target) {
			return nil, ErrUnauthorized
		}
		isGuardian, err := s.repo.IsGuardian(ctx, actorID, targetUserID)
		if err != nil || !isGuardian {
			return nil, ErrUnauthorized
		}
	case actorRole == "student":
		if actorID != targetUserID {
			return nil, ErrUnauthorized
		}
	default:
		return nil, ErrUnauthorized
	}
	return s.repo.GetAssignments(ctx, targetUserID)
}

// AddUserAssignment assigns an extra duty/assignment to a user enforcing authority rules.
func (s *Service) AddUserAssignment(ctx context.Context, actorRole, actorUserID, actorSchoolID, targetUserID string, req CreateAssignmentRequest) (*UserAssignment, error) {
	if !canAssignDuty(actorRole, req.AssignmentType) {
		return nil, fmt.Errorf("%w: non hai i permessi per attribuire questo incarico", ErrForbidden)
	}

	targetUser, err := s.repo.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	schoolID := actorSchoolID
	if schoolID == "" && targetUser.SchoolID != nil {
		schoolID = *targetUser.SchoolID
	}
	if actorRole != "superadmin" && (actorSchoolID == "" || targetUser.SchoolID == nil || *targetUser.SchoolID != actorSchoolID) {
		return nil, ErrUnauthorized
	}

	assignment := &UserAssignment{
		ID:             uuid.New().String(),
		SchoolID:       schoolID,
		UserID:         targetUserID,
		AssignmentType: req.AssignmentType,
		ScopeType:      req.ScopeType,
		ScopeID:        req.ScopeID,
		Title:          req.Title,
		AssignedBy:     &actorUserID,
		Metadata:       req.Metadata,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if assignment.ScopeType == "" {
		assignment.ScopeType = "school"
	}
	if assignment.Title == "" {
		assignment.Title = req.AssignmentType
	}

	if err := s.repo.CreateAssignment(ctx, assignment); err != nil {
		return nil, err
	}

	// Sincronizzazione automatica se incarico è coordinatore di classe
	if req.AssignmentType == "coordinatore_classe" && req.ScopeID != nil && *req.ScopeID != "" {
		_ = s.repo.SetCoordinatedClasses(ctx, schoolID, targetUserID, []string{*req.ScopeID})
	}

	return assignment, nil
}

// DeleteUserAssignment revokes an assignment from a user.
func (s *Service) DeleteUserAssignment(ctx context.Context, actorRole, actorSchoolID, targetUserID, assignmentID string) error {
	targetUser, err := s.repo.GetByID(ctx, targetUserID)
	if err != nil {
		return err
	}
	if actorRole != "superadmin" && (actorSchoolID == "" || targetUser.SchoolID == nil || *targetUser.SchoolID != actorSchoolID) {
		return ErrUnauthorized
	}

	assignments, err := s.repo.GetAssignments(ctx, targetUserID)
	if err != nil {
		return err
	}

	var targetAssignment *UserAssignment
	for _, a := range assignments {
		if a.ID == assignmentID {
			targetAssignment = &a
			break
		}
	}
	if targetAssignment == nil {
		return errors.New("assignment not found")
	}

	if !canAssignDuty(actorRole, targetAssignment.AssignmentType) {
		return fmt.Errorf("%w: non hai i permessi per revocare questo incarico", ErrForbidden)
	}

	return s.repo.DeleteAssignment(ctx, assignmentID)
}

// SetCoordinatedClasses sets all classes coordinated by a teacher (1 or more classes).
func (s *Service) SetCoordinatedClasses(ctx context.Context, actorRole, actorSchoolID, teacherUserID string, classIDs []string) error {
	if !canAssignDuty(actorRole, "coordinatore_classe") {
		return fmt.Errorf("%w: solo Dirigente Scolastico, Admin o Segreteria possono assegnare coordinatori di classe", ErrForbidden)
	}

	teacher, err := s.repo.GetByID(ctx, teacherUserID)
	if err != nil {
		return err
	}
	schoolID := actorSchoolID
	if schoolID == "" && teacher.SchoolID != nil {
		schoolID = *teacher.SchoolID
	}
	if actorRole != "superadmin" && (actorSchoolID == "" || teacher.SchoolID == nil || *teacher.SchoolID != actorSchoolID) {
		return ErrUnauthorized
	}

	return s.repo.SetCoordinatedClasses(ctx, schoolID, teacherUserID, classIDs)
}
