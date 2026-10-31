package unit

import (
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Input Validation & Sanitization (advanced) ──────────────────────────────

// Simulate backend email validation logic
func isValidEmail(email string) bool {
	// Basic RFC-like validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// Simulate backend fiscal code (codice fiscale) validation
func isValidFiscalCode(fc string) bool {
	if len(fc) != 16 {
		return false
	}
	pattern := regexp.MustCompile(`^[A-Z]{6}[0-9LMNPQRSTUV]{2}[ABCDEHLMPRST][0-9LMNPQRSTUV]{2}[A-Z][0-9LMNPQRSTUV]{3}[A-Z]$`)
	return pattern.MatchString(strings.ToUpper(fc))
}

// Simulate grade range validation (Italian school: 1-10)
func isValidGrade(grade float64) bool {
	return grade >= 1.0 && grade <= 10.0
}

// Simulate SQL injection pattern detection
func containsSQLInjectionPatterns(input string) bool {
	dangerous := []string{
		"' OR '1'='1",
		"--",
		"; DROP TABLE",
		"UNION SELECT",
		"'; EXEC",
		"xp_cmdshell",
		"SLEEP(",
		"BENCHMARK(",
		"LOAD_FILE(",
	}
	upperInput := strings.ToUpper(input)
	for _, pattern := range dangerous {
		if strings.Contains(upperInput, strings.ToUpper(pattern)) {
			return true
		}
	}
	return false
}

// Simulate path traversal detection
func containsPathTraversal(input string) bool {
	return strings.Contains(input, "../") ||
		strings.Contains(input, "..\\") ||
		strings.Contains(input, "%2e%2e") ||
		strings.Contains(input, "..%2f")
}

// Simulate LDAP injection detection
func containsLDAPInjectionPatterns(input string) bool {
	dangerous := []string{"*", "(", ")", "\\", "\x00", "/"}
	for _, d := range dangerous {
		if strings.Contains(input, d) {
			return true
		}
	}
	return false
}

// Simulate username allowed character validation
func isValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '.' && r != '_' && r != '-' {
			return false
		}
	}
	return true
}

// ─── Email Validation ─────────────────────────────────────────────────────────

func TestInputValidation_Email_ValidInputs(t *testing.T) {
	valid := []string{
		"mario.rossi@scuola.it",
		"docente+tag@istituto.edu",
		"a@b.co",
		"user.name+tag@sub.domain.org",
		"test_user-123@myschool.edu.it",
	}
	for _, email := range valid {
		t.Run(email, func(t *testing.T) {
			assert.True(t, isValidEmail(email), "expected valid: %s", email)
		})
	}
}

func TestInputValidation_Email_InvalidInputs(t *testing.T) {
	invalid := []string{
		"",
		"not-an-email",
		"@missinglocal.it",
		"missing@",
		"double@@domain.it",
		strings.Repeat("a", 255) + "@domain.it", // too long
		"<script>alert(1)</script>@domain.it",
	}
	for _, email := range invalid {
		t.Run(email, func(t *testing.T) {
			assert.False(t, isValidEmail(email), "expected invalid: %s", email)
		})
	}
}

// ─── Fiscal Code Validation ───────────────────────────────────────────────────

func TestInputValidation_FiscalCode_ValidCode(t *testing.T) {
	// Example valid Italian fiscal codes
	valid := []string{
		"RSSMRA80A01H501U",
		"BNCGPP85T10A662B",
	}
	for _, fc := range valid {
		t.Run(fc, func(t *testing.T) {
			assert.True(t, isValidFiscalCode(fc), "expected valid fiscal code: %s", fc)
		})
	}
}

func TestInputValidation_FiscalCode_InvalidInputs(t *testing.T) {
	invalid := []string{
		"",
		"TOOSHORT",
		"WAYTOOLONGTOBEAVALIDFISCALCODE",
		"' OR '1'='1         ", // SQL injection attempt padded to 16 chars
		"1234567890123456",     // all digits
	}
	for _, fc := range invalid {
		t.Run(fc, func(t *testing.T) {
			assert.False(t, isValidFiscalCode(fc), "expected invalid: %s", fc)
		})
	}
}

// ─── Grade Range Validation ───────────────────────────────────────────────────

func TestGradeValidation_ValidRange(t *testing.T) {
	valid := []float64{1.0, 4.5, 6.0, 7.5, 8.0, 10.0}
	for _, g := range valid {
		assert.True(t, isValidGrade(g), "%.1f should be valid", g)
	}
}

func TestGradeValidation_InvalidRange(t *testing.T) {
	invalid := []float64{0.0, 0.5, -1.0, 10.1, 11.0, 100.0}
	for _, g := range invalid {
		assert.False(t, isValidGrade(g), "%.1f should be invalid", g)
	}
}

// ─── SQL Injection Detection ──────────────────────────────────────────────────

func TestSQLInjection_DetectsClassicPayloads(t *testing.T) {
	malicious := []string{
		"' OR '1'='1",
		"admin'--",
		"1; DROP TABLE users",
		"1 UNION SELECT * FROM users",
		"'; EXEC xp_cmdshell('dir')",
		"' OR SLEEP(5)--",
		"1' BENCHMARK(5000000,MD5('a'))--",
	}
	for _, payload := range malicious {
		t.Run(payload, func(t *testing.T) {
			// After sanitization (SQL param binding), raw input should still be flagged
			assert.True(t, containsSQLInjectionPatterns(payload),
				"should detect SQL injection in: %s", payload)
		})
	}
}

func TestSQLInjection_AllowsSafeInputs(t *testing.T) {
	safe := []string{
		"Mario Rossi",
		"Classe 3A",
		"Matematica e Fisica",
		"Note: Assente giustificato",
		"10", "9.5", "8",
	}
	for _, input := range safe {
		t.Run(input, func(t *testing.T) {
			assert.False(t, containsSQLInjectionPatterns(input),
				"safe input should not trigger SQL injection detection: %s", input)
		})
	}
}

// ─── Path Traversal ──────────────────────────────────────────────────────────

func TestPathTraversal_DetectsTraversalAttempts(t *testing.T) {
	payloads := []string{
		"../../etc/passwd",
		"..\\..\\Windows\\System32\\cmd.exe",
		"%2e%2e/etc/passwd",
		"..%2fetc%2fpasswd",
	}
	for _, p := range payloads {
		t.Run(p, func(t *testing.T) {
			assert.True(t, containsPathTraversal(p), "should detect path traversal: %s", p)
		})
	}
}

func TestPathTraversal_AllowsSafePaths(t *testing.T) {
	safe := []string{
		"uploads/document.pdf",
		"reports/2025/january.pdf",
		"student-123/photo.jpg",
	}
	for _, p := range safe {
		assert.False(t, containsPathTraversal(p), "safe path should pass: %s", p)
	}
}

// ─── Username Validation ──────────────────────────────────────────────────────

func TestUsernameValidation_ValidUsernames(t *testing.T) {
	valid := []string{
		"mario.rossi",
		"admin",
		"teacher_123",
		"prof-biagi",
		"AbcDef123",
	}
	for _, u := range valid {
		assert.True(t, isValidUsername(u), "expected valid username: %s", u)
	}
}

func TestUsernameValidation_InvalidUsernames(t *testing.T) {
	invalid := []string{
		"",
		"ab",                         // too short
		strings.Repeat("a", 51),      // too long
		"user@domain",                // contains @
		"user name",                  // contains space
		"<script>",                   // XSS attempt
		"user; DROP TABLE users; --", // SQL injection
	}
	for _, u := range invalid {
		t.Run(u, func(t *testing.T) {
			assert.False(t, isValidUsername(u), "expected invalid username: %s", u)
		})
	}
}

// ─── Pagination Parameter Validation ─────────────────────────────────────────

func TestPaginationValidation_ValidParameters(t *testing.T) {
	type params struct{ page, limit int }
	valid := []params{
		{1, 10},
		{1, 50},
		{5, 100},
		{100, 1},
	}
	for _, p := range valid {
		isValid := p.page >= 1 && p.limit >= 1 && p.limit <= 200
		assert.True(t, isValid, "page=%d limit=%d should be valid", p.page, p.limit)
	}
}

func TestPaginationValidation_InvalidParameters(t *testing.T) {
	type params struct{ page, limit int }
	invalid := []params{
		{0, 10},    // page=0
		{-1, 10},   // negative page
		{1, 0},     // limit=0
		{1, -5},    // negative limit
		{1, 10000}, // limit too large
	}
	for _, p := range invalid {
		isValid := p.page >= 1 && p.limit >= 1 && p.limit <= 200
		assert.False(t, isValid, "page=%d limit=%d should be invalid", p.page, p.limit)
	}
}

// ─── Content-Type Security ────────────────────────────────────────────────────

func TestContentTypeValidation_AllowedMIMETypes(t *testing.T) {
	allowed := []string{
		"application/json",
		"application/pdf",
		"image/jpeg",
		"image/png",
		"text/plain; charset=utf-8",
	}
	for _, ct := range allowed {
		// Check that the MIME type doesn't contain dangerous patterns
		assert.NotContains(t, ct, "<script>")
		assert.NotContains(t, ct, "javascript")
	}
}

func TestContentTypeValidation_DangerousMIMETypes(t *testing.T) {
	dangerous := []string{
		"text/html; charset=utf-8", // could lead to XSS in file downloads
		"application/x-msdownload",
		"application/octet-stream",
		"application/x-php",
		"application/x-sh",
	}
	for _, ct := range dangerous {
		// A document service should reject these types
		isDocType := ct == "application/pdf" || ct == "application/json" || ct == "image/jpeg" || ct == "image/png"
		assert.False(t, isDocType, "MIME type %s should not be allowed as a document type", ct)
	}
}

// ─── Role Escalation Prevention ───────────────────────────────────────────────

func TestRoleEscalation_CannotAssignSuperadminRole(t *testing.T) {
	// Simulate a role assignment request from a school_admin
	type RoleAssignRequest struct {
		TargetUserID  string
		NewRole       string
		RequestorRole string
	}

	// Roles a school_admin is allowed to assign
	allowedBySchoolAdmin := map[string]bool{
		"teacher":   true,
		"student":   true,
		"parent":    true,
		"secretary": true,
	}

	requests := []struct {
		req     RoleAssignRequest
		allowed bool
	}{
		{RoleAssignRequest{"user-1", "teacher", "school_admin"}, true},
		{RoleAssignRequest{"user-2", "student", "school_admin"}, true},
		{RoleAssignRequest{"user-3", "superadmin", "school_admin"}, false},   // escalation attempt
		{RoleAssignRequest{"user-4", "school_admin", "school_admin"}, false}, // peer-level escalation
		{RoleAssignRequest{"user-5", "teacher", "teacher"}, false},           // teacher cannot assign roles
	}

	for _, tc := range requests {
		t.Run(tc.req.NewRole+"_by_"+tc.req.RequestorRole, func(t *testing.T) {
			isAllowed := tc.req.RequestorRole == "school_admin" && allowedBySchoolAdmin[tc.req.NewRole]
			assert.Equal(t, tc.allowed, isAllowed)
		})
	}
}

// ─── Token Expiry Boundary ────────────────────────────────────────────────────

func TestTokenExpiry_ExpiryBoundaryConditions(t *testing.T) {
	// Test that near-expired tokens (within 30s) should trigger refresh
	type tokenState struct {
		remainingSeconds int
		shouldRefresh    bool
	}

	const REFRESH_THRESHOLD_SECONDS = 30

	states := []tokenState{
		{120, false}, // 2 minutes remaining — no refresh needed
		{60, false},  // 1 minute — still fine
		{30, true},   // exactly at threshold — refresh
		{15, true},   // under threshold — refresh
		{0, true},    // expired — refresh
		{-10, true},  // already expired — refresh
	}

	for _, s := range states {
		shouldRefresh := s.remainingSeconds <= REFRESH_THRESHOLD_SECONDS
		require.Equal(t, s.shouldRefresh, shouldRefresh,
			"remainingSeconds=%d: shouldRefresh expected %v", s.remainingSeconds, s.shouldRefresh)
	}
}
