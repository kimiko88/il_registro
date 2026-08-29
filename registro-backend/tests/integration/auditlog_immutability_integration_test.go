package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/auditlog"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Auditlog Repository ──────────────────────────────────────────────

type mockAuditlogRepo struct {
	events []auditlog.AuditEvent
}

func newMockAuditlogRepo() *mockAuditlogRepo {
	return &mockAuditlogRepo{events: make([]auditlog.AuditEvent, 0)}
}

func (m *mockAuditlogRepo) Insert(_ context.Context, event *auditlog.AuditEvent) error {
	m.InsertAsync(*event)
	return nil
}

func (m *mockAuditlogRepo) InsertAsync(event auditlog.AuditEvent) {
	if event.ID == "" {
		event.ID = "log-" + time.Now().Format("20060102150405")
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	m.events = append(m.events, event)
}

func (m *mockAuditlogRepo) List(_ context.Context, p auditlog.FilterParams) ([]auditlog.AuditEvent, int, error) {
	var filtered []auditlog.AuditEvent
	for _, e := range m.events {
		if p.SchoolID != "" && e.SchoolID != p.SchoolID {
			continue
		}
		if p.Action != "" && e.Action != p.Action {
			continue
		}
		if p.EntityType != "" && e.EntityType != p.EntityType {
			continue
		}
		filtered = append(filtered, e)
	}
	total := len(filtered)
	start := (p.Page - 1) * p.Limit
	if start > total {
		return []auditlog.AuditEvent{}, total, nil
	}
	end := start + p.Limit
	if end > total {
		end = total
	}
	return filtered[start:end], total, nil
}

// ─── Router Setup Helper ───────────────────────────────────────────────────

func setupAuditlogRouter(repo *mockAuditlogRepo, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := auditlog.NewService(repo)
	h := auditlog.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		if role != "" {
			c.Set("role", role)
			c.Set("school_id", schoolID)
			c.Set("user_id", "admin-1")
		}
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

// ─── Tests ─────────────────────────────────────────────────────────────────

// AUD01 — Admin can list audit logs for their school
func TestAuditLog_List_AdminSuccess(t *testing.T) {
	repo := newMockAuditlogRepo()
	repo.InsertAsync(auditlog.AuditEvent{
		SchoolID:   "school-1",
		ActorID:    "u-1",
		ActorName:  "Prof. Rossi",
		Action:     "GRADE_CREATE",
		EntityType: "grade",
		Details:    "Voto 8 registrato per Mario Rossi",
	})
	repo.InsertAsync(auditlog.AuditEvent{
		SchoolID:   "school-1",
		ActorID:    "u-2",
		ActorName:  "Segreteria",
		Action:     "DOCUMENT_PUBLISH",
		EntityType: "document",
		Details:    "Circolare n. 42 pubblicata",
	})

	r := setupAuditlogRouter(repo, "admin", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/audit-log?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp auditlog.PaginatedAuditLogs
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 2, resp.Total)
	assert.Len(t, resp.Data, 2)
}

// AUD02 — Audit logs filtered by action type
func TestAuditLog_List_FilterByAction(t *testing.T) {
	repo := newMockAuditlogRepo()
	repo.InsertAsync(auditlog.AuditEvent{
		SchoolID:   "school-1",
		Action:     "USER_LOGIN",
		EntityType: "auth",
	})
	repo.InsertAsync(auditlog.AuditEvent{
		SchoolID:   "school-1",
		Action:     "PASSWORD_RESET",
		EntityType: "auth",
	})

	r := setupAuditlogRouter(repo, "admin", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/audit-log?action=PASSWORD_RESET", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp auditlog.PaginatedAuditLogs
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Total)
	assert.Equal(t, "PASSWORD_RESET", resp.Data[0].Action)
}

// AUD03 — Teachers and students are forbidden from accessing audit logs
func TestAuditLog_List_ForbiddenForTeacherAndStudent(t *testing.T) {
	repo := newMockAuditlogRepo()

	rTeacher := setupAuditlogRouter(repo, "teacher", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/audit-log", nil)
	wTeacher := httptest.NewRecorder()
	rTeacher.ServeHTTP(wTeacher, req)
	assert.Equal(t, http.StatusForbidden, wTeacher.Code)

	rStudent := setupAuditlogRouter(repo, "student", "school-1")
	wStudent := httptest.NewRecorder()
	rStudent.ServeHTTP(wStudent, req)
	assert.Equal(t, http.StatusForbidden, wStudent.Code)
}

// AUD04 — Superadmin and system auditor have authorized access
func TestAuditLog_List_SystemAuditorAccess(t *testing.T) {
	repo := newMockAuditlogRepo()
	repo.InsertAsync(auditlog.AuditEvent{
		SchoolID: "school-A",
		Action:   "ENCRYPTION_KEY_ROTATE",
	})

	r := setupAuditlogRouter(repo, "system_auditor", "")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/audit-log", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// AUD05 — CSV export streams CSV with proper headers
func TestAuditLog_ExportCSV_Admin(t *testing.T) {
	repo := newMockAuditlogRepo()
	repo.InsertAsync(auditlog.AuditEvent{
		SchoolID:   "school-1",
		ActorName:  "Admin Rossi",
		Action:     "USER_SUSPEND",
		EntityType: "user",
		Details:    "Account sospeso per violazione policy",
	})

	r := setupAuditlogRouter(repo, "admin", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/audit-log/export", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")
	assert.Contains(t, w.Body.String(), "Data/Ora")
	assert.Contains(t, w.Body.String(), "USER_SUSPEND")
}
