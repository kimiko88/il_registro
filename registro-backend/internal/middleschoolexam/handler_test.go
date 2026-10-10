package middleschoolexam

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupExamTestRouterWithRole(svc *Service, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "school-test-1")
		c.Set("role", role)
		c.Next()
	})
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestHandler_RoleProtection(t *testing.T) {
	repo := newMockExamRepo()
	svc := NewService(repo)
	// Role "student" must be forbidden
	r := setupExamTestRouterWithRole(svc, "student")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/middle-school-exam/classes/cls-1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden status 403, got %d", w.Code)
	}
}

func TestHandler_GetOrCreateExam_Success(t *testing.T) {
	repo := newMockExamRepo()
	svc := NewService(repo)
	r := setupExamTestRouterWithRole(svc, "teacher")

	payload := map[string]string{
		"president_name": "Prof. Rossi",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/middle-school-exam/classes/cls-3b", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_SaveAdmission_Success(t *testing.T) {
	repo := newMockExamRepo()
	svc := NewService(repo)
	r := setupExamTestRouterWithRole(svc, "coordinator")

	payload := map[string]interface{}{
		"grade":    8,
		"judgment": "Ammesso con giudizio positivo",
		"admitted": true,
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/middle-school-exam/exam-123/candidates/st-456/admission", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GenerateDiploma_Success(t *testing.T) {
	repo := newMockExamRepo()
	svc := NewService(repo)
	r := setupExamTestRouterWithRole(svc, "secretary")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/middle-school-exam/candidates/cand-1/diploma", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("expected text/plain content type, got %s", w.Header().Get("Content-Type"))
	}
}
