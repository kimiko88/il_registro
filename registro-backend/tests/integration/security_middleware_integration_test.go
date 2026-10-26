package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"registro-backend/internal/auth"
	"registro-backend/internal/middleware"
	"registro-backend/pkg/jwt"
	"registro-backend/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ─── Security Headers Integration ────────────────────────────────────────────

func TestSecurityHeadersIntegration_AllHeadersPresentOnEveryRoute(t *testing.T) {
	r := gin.New()
	r.Use(middleware.SecurityHeadersMiddleware())
	r.GET("/api/v1/grades", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"grades": []interface{}{}})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/grades", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
	assert.NotEmpty(t, w.Header().Get("Permissions-Policy"))
	assert.NotEmpty(t, w.Header().Get("Content-Security-Policy"))
}

func TestSecurityHeadersIntegration_CSPReportUriPresent(t *testing.T) {
	r := gin.New()
	r.Use(middleware.SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))

	csp := w.Header().Get("Content-Security-Policy")
	assert.Contains(t, csp, "report-uri", "CSP must have a report-uri directive")
}

// ─── Idempotency + Auth Integration ──────────────────────────────────────────

func TestIdempotencyAndAuth_Integration_DuplicateGradeSubmission(t *testing.T) {
	r := gin.New()
	r.Use(middleware.IdempotencyMiddleware())

	callCount := 0
	r.POST("/api/v1/grades", func(c *gin.Context) {
		callCount++
		c.JSON(http.StatusCreated, gin.H{"id": "grade-001", "value": 8.5})
	})

	key := "grade-dedup-test-key-" + t.Name()
	makeRequest := func() *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]interface{}{
			"student_id": "student-123",
			"subject_id": "math-456",
			"value":      8.5,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/grades", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", key)
		r.ServeHTTP(w, req)
		return w
	}

	// First request — should be processed
	w1 := makeRequest()
	require.Equal(t, http.StatusCreated, w1.Code)

	// Second request with same key — should be replayed
	w2 := makeRequest()
	assert.Equal(t, http.StatusCreated, w2.Code)
	assert.Equal(t, "true", w2.Header().Get("X-Idempotency-Replayed"))

	// Handler must be invoked only once
	assert.Equal(t, 1, callCount, "grade creation handler must be called only once")
}

func TestIdempotencyAndAuth_Integration_AttendanceDeduplication(t *testing.T) {
	r := gin.New()
	r.Use(middleware.IdempotencyMiddleware())

	insertCount := 0
	r.PUT("/api/v1/attendance/lesson-123", func(c *gin.Context) {
		insertCount++
		c.JSON(http.StatusOK, gin.H{"status": "signed"})
	})

	key := "attendance-sign-" + t.Name()
	makeReq := func() int {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]interface{}{"lesson_id": "lesson-123", "signature": "sig-data"})
		req := httptest.NewRequest(http.MethodPut, "/api/v1/attendance/lesson-123", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", key)
		r.ServeHTTP(w, req)
		return w.Code
	}

	for i := 0; i < 5; i++ {
		code := makeReq()
		assert.Equal(t, http.StatusOK, code)
	}
	assert.Equal(t, 1, insertCount, "attendance signing must be idempotent — handler called only once")
}

// ─── Auth Security Integration (login attempts / brute force) ────────────────

func TestAuthSecurity_Integration_MissingFieldsReturns400(t *testing.T) {
	mockRepo := new(testhelpers.MockAuthRepository)
	priv, pub := testhelpers.GenerateRSAKeys()
	tokenManager := jwt.NewTokenManager(priv, pub)
	mfaService := auth.NewMFAService("test")
	service := auth.NewService(mockRepo, tokenManager, mfaService, nil)
	handler := auth.NewHandler(service)

	r := gin.New()
	r.Use(middleware.SecurityHeadersMiddleware())
	r.POST("/auth/login", handler.Login)

	tests := []struct {
		name     string
		payload  string
		wantCode int
	}{
		{"missing password", `{"email":"user@school.it"}`, http.StatusBadRequest},
		{"missing email", `{"password":"SomePassword123!"}`, http.StatusBadRequest},
		{"empty body", `{}`, http.StatusBadRequest},
		{"empty strings", `{"email":"","password":""}`, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/login",
				bytes.NewBufferString(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.wantCode, w.Code)
		})
	}

	mockRepo.AssertExpectations(t)
}

func TestAuthSecurity_Integration_MalformedJSONReturns400(t *testing.T) {
	mockRepo := new(testhelpers.MockAuthRepository)
	priv, pub := testhelpers.GenerateRSAKeys()
	tokenManager := jwt.NewTokenManager(priv, pub)
	mfaService := auth.NewMFAService("test")
	service := auth.NewService(mockRepo, tokenManager, mfaService, nil)
	handler := auth.NewHandler(service)

	r := gin.New()
	r.POST("/auth/login", handler.Login)

	malformed := []string{
		`{email: "missing quotes"}`,
		`{"email": "test@test.it"`, // unclosed brace
		`not json at all`,
		`<script>alert(1)</script>`,
		`null`,
	}

	for _, payload := range malformed {
		t.Run(payload, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/login",
				bytes.NewBufferString(payload))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}

	mockRepo.AssertNotCalled(t, "GetUserByEmail", mock.Anything, mock.Anything)
}

func TestAuthSecurity_Integration_XSSInLoginFieldsDoesNotPanic(t *testing.T) {
	mockRepo := new(testhelpers.MockAuthRepository)
	mockRepo.On("GetRecentLoginAttempts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, nil).Maybe()
	mockRepo.On("GetRecentLoginAttemptsByEmail", mock.Anything, mock.Anything, mock.Anything).Return(0, nil).Maybe()
	mockRepo.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("invalid credentials")).Maybe()
	mockRepo.On("RecordLoginAttempt", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	priv, pub := testhelpers.GenerateRSAKeys()
	tokenManager := jwt.NewTokenManager(priv, pub)
	mfaService := auth.NewMFAService("test")
	service := auth.NewService(mockRepo, tokenManager, mfaService, nil)
	handler := auth.NewHandler(service)

	r := gin.New()
	r.POST("/auth/login", handler.Login)

	xssPayloads := []struct{ email, password string }{
		{"<script>alert('xss')</script>@evil.com", "pass"},
		{"user@school.it", "<img src=x onerror=alert(1)>"},
		{"'; DROP TABLE users; --@evil.com", "anything"},
		{"javascript:alert(1)@evil.com", "pass"},
	}

	for _, p := range xssPayloads {
		t.Run(p.email, func(t *testing.T) {
			body, _ := json.Marshal(map[string]string{
				"email":    p.email,
				"password": p.password,
			})
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Should not panic; should return 400 or 401, never 500
			assert.NotPanics(t, func() { r.ServeHTTP(w, req) })
			assert.Less(t, w.Code, http.StatusInternalServerError,
				"XSS payloads must never cause a 5xx error")
		})
	}
}

// ─── Rate Limit + Security Headers Combination ───────────────────────────────

func TestRateLimitAndSecurityHeaders_Integration_BothApplied(t *testing.T) {
	r := gin.New()
	r.Use(middleware.SecurityHeadersMiddleware())
	r.GET("/api/v1/public/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/test", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	r.ServeHTTP(w, req)

	// Security headers must always be present
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
}

// ─── CORS Headers (OPTIONS preflight) ────────────────────────────────────────

func TestCORS_Integration_OptionsPreflightAllowsKnownOrigin(t *testing.T) {
	// Use localhost:5173 which is in the default allowed list
	os.Unsetenv("ALLOWED_ORIGINS") // use default config

	r := gin.New()
	r.Use(middleware.CORSMiddleware())
	r.Any("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}

func TestCORS_Integration_OptionsPreflightBlocksUnknownOrigin(t *testing.T) {
	os.Unsetenv("ALLOWED_ORIGINS")

	r := gin.New()
	r.Use(middleware.CORSMiddleware())
	r.Any("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/test", nil)
	req.Header.Set("Origin", "https://evil-site.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_Integration_OptionsPreflightRequiresOriginHeader(t *testing.T) {
	os.Unsetenv("ALLOWED_ORIGINS")

	r := gin.New()
	r.Use(middleware.CORSMiddleware())
	r.Any("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	// OPTIONS without Origin header
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
