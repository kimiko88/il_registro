package users

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// GDPRHandler handles privacy-related tasks
type GDPRHandler struct {
	repo Repository
}

func NewGDPRHandler(repo Repository) *GDPRHandler {
	return &GDPRHandler{repo: repo}
}

// PseudonymizeUser replaces PII with random/hashed values
// This allows keeping the user ID for relational integrity (grades, etc.)
// while removing personal data.
func (h *GDPRHandler) PseudonymizeUser(user *User) {
	// Generate a secure hash of the ID + time as key
	hash := sha256.Sum256([]byte(user.ID + time.Now().String()))
	pseudoID := hex.EncodeToString(hash[:])[:12]

	user.FirstName = "Deleted"
	user.LastName = "User-" + pseudoID
	user.Email = fmt.Sprintf("deleted-%s@anonymized.local", pseudoID)
	anon := "ANONYMIZED"
	user.FiscalCode = &anon
	empty := ""
	user.PhoneNumber = &empty
	user.JobTitle = &empty
	user.PasswordHash = "" // Clear password
	user.MFASecret = ""
	user.MFAEnabled = false
	user.IsActive = false
	user.EmailVerified = false

	now := time.Now()
	user.PseudonymizedAt = &now
	user.DeletedAt = &now // Also mark as deleted
}

// GenerateDataExport creates a map of all user data for portability
func (h *GDPRHandler) GenerateDataExport(user *User, logs []AuditLog) map[string]interface{} {
	return map[string]interface{}{
		"profile": map[string]interface{}{
			"id":          user.ID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"fiscal_code": user.FiscalCode,
			"role":        user.Role,
			"created_at":  user.CreatedAt,
		},
		"audit_logs": logs,
		// In a real system, we'd fetch grades, attendance, etc. here too
		"generated_at": time.Now(),
	}
}
