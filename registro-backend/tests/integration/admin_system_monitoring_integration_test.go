package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/admin"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Admin Repository ─────────────────────────────────────────────────

type mockMonitoringAdminRepo struct{}

func (m *mockMonitoringAdminRepo) CountSchools(_ context.Context, _ *string) (int64, error) { return 1, nil }
func (m *mockMonitoringAdminRepo) CountUsers(_ context.Context, _ *string) (int64, error)   { return 10, nil }
func (m *mockMonitoringAdminRepo) CountUsersByRole(_ context.Context, _ string, _ *string) (int64, error) {
	return 5, nil
}
func (m *mockMonitoringAdminRepo) CountActiveUsers24h(_ context.Context, _ *string) (int64, error) {
	return 2, nil
}
func (m *mockMonitoringAdminRepo) CountDocuments(_ context.Context, _ *string) (int64, error) {
	return 10, nil
}
func (m *mockMonitoringAdminRepo) CountPendingDocuments(_ context.Context, _ *string) (int64, error) {
	return 1, nil
}
func (m *mockMonitoringAdminRepo) CountCommunications(_ context.Context, _ *string) (int64, error) {
	return 5, nil
}
func (m *mockMonitoringAdminRepo) GetRecentEvents(_ context.Context, _ int, _ *string) ([]admin.RecentEvent, error) {
	return nil, nil
}
func (m *mockMonitoringAdminRepo) GetSystemHealth(_ context.Context) (*admin.SystemHealthStatus, error) {
	return &admin.SystemHealthStatus{
		OverallStatus: "healthy",
		Database:      admin.HealthCheck{Status: "healthy"},
		Storage:       admin.HealthCheck{Status: "healthy"},
		API:           admin.HealthCheck{Status: "healthy"},
	}, nil
}
func (m *mockMonitoringAdminRepo) GetUserGrowth(_ context.Context, _ *string) ([]admin.UserGrowthPoint, error) {
	return nil, nil
}
func (m *mockMonitoringAdminRepo) ListSchools(_ context.Context, _ *admin.SchoolListRequest, _ int, _ *string) ([]admin.SchoolResponse, int64, error) {
	return nil, 0, nil
}
func (m *mockMonitoringAdminRepo) GetSchool(_ context.Context, _ string) (*admin.SchoolResponse, error) {
	return nil, nil
}
func (m *mockMonitoringAdminRepo) CreateSchool(_ context.Context, _ *admin.CreateSchoolRequest) (*admin.SchoolResponse, error) {
	return nil, nil
}
func (m *mockMonitoringAdminRepo) UpdateSchool(_ context.Context, _ string, _ *admin.UpdateSchoolRequest) (*admin.SchoolResponse, error) {
	return nil, nil
}
func (m *mockMonitoringAdminRepo) DeleteSchool(_ context.Context, _ string) error { return nil }
func (m *mockMonitoringAdminRepo) SchoolCodeExists(_ context.Context, _ string) (bool, error) {
	return false, nil
}
func (m *mockMonitoringAdminRepo) ListAdminUsers(_ context.Context, _, _ int, _ *string) ([]admin.AdminUserResponse, int64, error) {
	return nil, 0, nil
}
func (m *mockMonitoringAdminRepo) GetAdminUserByID(_ context.Context, _ string) (*admin.AdminUserResponse, error) {
	return nil, nil
}
func (m *mockMonitoringAdminRepo) CreateAdminUser(_ context.Context, _ *admin.CreateAdminRequest) (*admin.AdminUserResponse, error) {
	return nil, nil
}
func (m *mockMonitoringAdminRepo) UpdateAdminUser(_ context.Context, _ string, _ *admin.UpdateAdminRequest) (*admin.AdminUserResponse, error) {
	return nil, nil
}
func (m *mockMonitoringAdminRepo) DeleteAdminUser(_ context.Context, _ string) error { return nil }
func (m *mockMonitoringAdminRepo) UserEmailExists(_ context.Context, _ string) (bool, error) {
	return false, nil
}
func (m *mockMonitoringAdminRepo) GetAdminActivity(_ context.Context, _ string, _ int) ([]admin.ActivityLogEntry, error) {
	return nil, nil
}
func (m *mockMonitoringAdminRepo) ListAuditLogs(_ context.Context, _ *admin.AuditLogListRequest, _ int) ([]admin.ActivityLogEntry, int64, error) {
	return nil, 0, nil
}
func (m *mockMonitoringAdminRepo) LogAdminAction(_ context.Context, _, _, _ string, _, _ *string, _ string) error {
	return nil
}
func (m *mockMonitoringAdminRepo) GetSetting(_ context.Context, _, _ string) (string, error) {
	return "", nil
}
func (m *mockMonitoringAdminRepo) UpdateSetting(_ context.Context, _, _, _ string) error { return nil }

// ─── Setup Router Helper ───────────────────────────────────────────────────

func setupAdminMonitoringRouter(role, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := admin.NewService(&mockMonitoringAdminRepo{})
	h := admin.NewHandler(svc)

	api := r.Group("/api/v1/admin")
	api.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	{
		api.GET("/system/health", func(c *gin.Context) {
			if c.GetString("user_id") == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			if c.GetString("role") != "superadmin" && c.GetString("role") != "admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
			h.GetSystemHealth(c)
		})
		api.GET("/system/metrics", func(c *gin.Context) {
			if c.GetString("user_id") == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			if c.GetString("role") != "superadmin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
			h.GetSystemMetrics(c)
		})
	}
	return r
}

// ─── Tests ─────────────────────────────────────────────────────────────────

// MON01 — Superadmin retrieves live system health status
func TestAdmin_GetSystemHealth_SuperAdmin(t *testing.T) {
	r := setupAdminMonitoringRouter("superadmin", "superadmin-1", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/system/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var health map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &health)
	require.NoError(t, err)
	assert.Equal(t, "healthy", health["status"])
	assert.NotNil(t, health["services"])
	assert.NotNil(t, health["metrics"])
}

// MON02 — Superadmin retrieves Go runtime memory and goroutine metrics
func TestAdmin_GetSystemMetrics_SuperAdmin(t *testing.T) {
	r := setupAdminMonitoringRouter("superadmin", "superadmin-1", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/system/metrics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var metrics map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &metrics)
	require.NoError(t, err)
	assert.Contains(t, metrics, "goroutines")
	assert.Contains(t, metrics, "memory_alloc_mb")
	assert.Contains(t, metrics, "uptime_seconds")
}

// MON03 — Unauthenticated request receives 401
func TestAdmin_SystemHealth_Unauthorized(t *testing.T) {
	r := setupAdminMonitoringRouter("", "", "")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/system/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// MON04 — Teacher is forbidden from accessing system metrics
func TestAdmin_SystemMetrics_TeacherForbidden(t *testing.T) {
	r := setupAdminMonitoringRouter("teacher", "teacher-1", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/system/metrics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
