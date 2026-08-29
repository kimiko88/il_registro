package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/tenants"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Tenants Repository ───────────────────────────────────────────────

type mockTenantsRepo struct {
	tenants map[string]*tenants.Tenant
}

func newMockTenantsRepo() *mockTenantsRepo {
	return &mockTenantsRepo{
		tenants: make(map[string]*tenants.Tenant),
	}
}

func (m *mockTenantsRepo) CreateTenant(_ context.Context, t *tenants.Tenant) error {
	t.ID = "school-" + t.Code
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	t.Status = "active"
	m.tenants[t.ID] = t
	return nil
}

func (m *mockTenantsRepo) GetTenantByID(_ context.Context, id string) (*tenants.Tenant, error) {
	t, ok := m.tenants[id]
	if !ok {
		return nil, errors.New("tenant not found")
	}
	return t, nil
}

func (m *mockTenantsRepo) ListTenants(_ context.Context) ([]*tenants.Tenant, error) {
	var list []*tenants.Tenant
	for _, t := range m.tenants {
		list = append(list, t)
	}
	return list, nil
}

// ── Setup Router Helper ───────────────────────────────────────────────────

func setupTenantsRouter(repo *mockTenantsRepo, role, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := tenants.NewService(repo)
	h := tenants.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

// ── Tests ─────────────────────────────────────────────────────────────────

// TNT01 — Superadmin successfully provisions a new tenant
func TestTenants_CreateTenant_SuperadminSuccess(t *testing.T) {
	repo := newMockTenantsRepo()
	r := setupTenantsRouter(repo, "superadmin", "super-1", "")

	body, _ := json.Marshal(tenants.CreateTenantRequest{
		Name:         "Liceo Scientifico Galileo Galilei",
		Code:         "RMPS01000P",
		MaxStudents:  1200,
		MaxTeachers:  100,
		MaxStorageMB: 51200,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var tenant tenants.Tenant
	err := json.Unmarshal(w.Body.Bytes(), &tenant)
	require.NoError(t, err)
	assert.Equal(t, "Liceo Scientifico Galileo Galilei", tenant.Name)
	assert.Equal(t, "RMPS01000P", tenant.Code)
	assert.Equal(t, 1200, tenant.Quota.MaxStudents)
}

// TNT02 — Regular school admin is forbidden from provisioning tenants
func TestTenants_CreateTenant_ForbiddenForAdmin(t *testing.T) {
	repo := newMockTenantsRepo()
	r := setupTenantsRouter(repo, "admin", "admin-1", "school-1")

	body, _ := json.Marshal(tenants.CreateTenantRequest{
		Name: "Nuova Scuola",
		Code: "SCUOLA01",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TNT03 — Teacher is forbidden from provisioning tenants
func TestTenants_CreateTenant_ForbiddenForTeacher(t *testing.T) {
	repo := newMockTenantsRepo()
	r := setupTenantsRouter(repo, "teacher", "teacher-1", "school-1")

	body, _ := json.Marshal(tenants.CreateTenantRequest{
		Name: "Nuova Scuola",
		Code: "SCUOLA02",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TNT04 — Superadmin can list all system tenants
func TestTenants_ListTenants_Superadmin(t *testing.T) {
	repo := newMockTenantsRepo()
	_ = repo.CreateTenant(context.Background(), &tenants.Tenant{Name: "Scuola 1", Code: "SC1"})
	_ = repo.CreateTenant(context.Background(), &tenants.Tenant{Name: "Scuola 2", Code: "SC2"})

	r := setupTenantsRouter(repo, "superadmin", "super-1", "")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var list []*tenants.Tenant
	err := json.Unmarshal(w.Body.Bytes(), &list)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

// TNT05 — Secretary cannot list all tenants
func TestTenants_ListTenants_ForbiddenForSecretary(t *testing.T) {
	repo := newMockTenantsRepo()
	r := setupTenantsRouter(repo, "secretary", "sec-1", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TNT06 — School admin can get their own tenant details
func TestTenants_GetTenant_SelfAccess(t *testing.T) {
	repo := newMockTenantsRepo()
	_ = repo.CreateTenant(context.Background(), &tenants.Tenant{
		Name:  "Scuola Dante",
		Code:  "DANTE01",
		Quota: tenants.TenantQuota{MaxStudents: 800},
	})
	tenantID := "school-DANTE01"

	r := setupTenantsRouter(repo, "admin", "admin-1", tenantID)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/"+tenantID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var tnt tenants.Tenant
	err := json.Unmarshal(w.Body.Bytes(), &tnt)
	require.NoError(t, err)
	assert.Equal(t, "Scuola Dante", tnt.Name)
}

// TNT07 — Admin cannot view other schools without superadmin role
func TestTenants_GetTenant_CrossSchoolForbidden(t *testing.T) {
	repo := newMockTenantsRepo()
	_ = repo.CreateTenant(context.Background(), &tenants.Tenant{
		Name: "Altra Scuola",
		Code: "OTHER01",
	})
	otherSchoolID := "school-OTHER01"

	r := setupTenantsRouter(repo, "admin", "admin-1", "my-school-id")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/"+otherSchoolID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TNT08 — Unauthenticated requests receive 401
func TestTenants_Unauthorized_NoAuth(t *testing.T) {
	repo := newMockTenantsRepo()
	r := setupTenantsRouter(repo, "", "", "")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/some-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
