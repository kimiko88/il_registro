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

type mockAuditLogLifecycleRepo struct {
	events []auditlog.AuditEvent
}

func newMockAuditLogLifecycleRepo() *mockAuditLogLifecycleRepo {
	return &mockAuditLogLifecycleRepo{
		events: []auditlog.AuditEvent{},
	}
}

func (m *mockAuditLogLifecycleRepo) Insert(ctx context.Context, event *auditlog.AuditEvent) error {
	if event.ID == "" {
		event.ID = "audit-" + time.Now().Format("150405.000000")
	}
	event.CreatedAt = time.Now()
	m.events = append(m.events, *event)
	return nil
}

func (m *mockAuditLogLifecycleRepo) InsertAsync(event auditlog.AuditEvent) {
	if event.ID == "" {
		event.ID = "audit-" + time.Now().Format("150405.000000")
	}
	event.CreatedAt = time.Now()
	m.events = append(m.events, event)
}

func (m *mockAuditLogLifecycleRepo) List(ctx context.Context, p auditlog.FilterParams) ([]auditlog.AuditEvent, int, error) {
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
	return filtered, len(filtered), nil
}

func setupAuditLogRouter(repo auditlog.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := auditlog.NewService(repo)
	handler := auditlog.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = "admin"
		}
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "admin-1"
		}
		c.Set("role", role)
		c.Set("user_id", userID)
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_AuditLog_Lifecycle_And_Export_Workflow(t *testing.T) {
	repo := newMockAuditLogLifecycleRepo()
	r := setupAuditLogRouter(repo)

	// Pre-populate with events
	repo.InsertAsync(auditlog.AuditEvent{
		SchoolID:   "school-1",
		ActorID:    "teacher-1",
		ActorRole:  "teacher",
		ActorName:  "Prof. Mario Rossi",
		Action:     "CREATE_GRADE",
		EntityType: "grades",
		EntityID:   "grade-101",
		Details:    "Inserito voto 8.5 per studente Mario Bianchi",
		IPAddress:  "192.168.1.50",
	})
	repo.InsertAsync(auditlog.AuditEvent{
		SchoolID:   "school-1",
		ActorID:    "secretary-1",
		ActorRole:  "secretary",
		ActorName:  "Anna Bianchi",
		Action:     "GENERATE_CERTIFICATE",
		EntityType: "certificates",
		EntityID:   "cert-202",
		Details:    "Emesso certificato iscrizione PROT-2026-0089",
		IPAddress:  "192.168.1.10",
	})

	// 1. Admin lists audit log entries
	reqList, _ := http.NewRequest("GET", "/api/v1/audit-log", nil)
	reqList.Header.Set("X-Role", "admin")
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)
	require.Equal(t, http.StatusOK, wList.Code)

	var listResp auditlog.PaginatedAuditLogs
	err := json.Unmarshal(wList.Body.Bytes(), &listResp)
	require.NoError(t, err)
	assert.Equal(t, 2, listResp.Total)

	// 2. Filter audit log by action
	reqFilter, _ := http.NewRequest("GET", "/api/v1/audit-log?action=CREATE_GRADE", nil)
	reqFilter.Header.Set("X-Role", "admin")
	wFilter := httptest.NewRecorder()
	r.ServeHTTP(wFilter, reqFilter)
	require.Equal(t, http.StatusOK, wFilter.Code)

	var filterResp auditlog.PaginatedAuditLogs
	err = json.Unmarshal(wFilter.Body.Bytes(), &filterResp)
	require.NoError(t, err)
	assert.Equal(t, 1, filterResp.Total)
	assert.Equal(t, "CREATE_GRADE", filterResp.Data[0].Action)

	// 3. Export audit log to CSV
	reqExport, _ := http.NewRequest("GET", "/api/v1/audit-log/export", nil)
	reqExport.Header.Set("X-Role", "admin")
	wExport := httptest.NewRecorder()
	r.ServeHTTP(wExport, reqExport)
	require.Equal(t, http.StatusOK, wExport.Code)
	assert.Contains(t, wExport.Header().Get("Content-Type"), "text/csv")
	assert.Contains(t, wExport.Body.String(), "CREATE_GRADE")

	// 4. Unauthorized student access should be forbidden
	reqForbidden, _ := http.NewRequest("GET", "/api/v1/audit-log", nil)
	reqForbidden.Header.Set("X-Role", "student")
	wForbidden := httptest.NewRecorder()
	r.ServeHTTP(wForbidden, reqForbidden)
	assert.Equal(t, http.StatusForbidden, wForbidden.Code)
}
