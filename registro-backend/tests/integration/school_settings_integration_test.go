package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/schoolsettings"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock ──────────────────────────────────────────────────────────────────

type mockSchoolSettingsRepo struct {
	data map[string]*schoolsettings.SchoolSettings
}

func (m *mockSchoolSettingsRepo) GetBySchoolID(_ context.Context, schoolID string) (*schoolsettings.SchoolSettings, error) {
	if s, ok := m.data[schoolID]; ok {
		return s, nil
	}
	// Return defaults
	return &schoolsettings.SchoolSettings{
		SchoolID:                       schoolID,
		AllowParentsViewGrades:         true,
		AllowStudentsViewClassAverages: false,
		RequireMFAForStaff:             false,
		UpdatedAt:                      time.Now(),
	}, nil
}

func (m *mockSchoolSettingsRepo) Upsert(_ context.Context, s *schoolsettings.SchoolSettings) error {
	s.UpdatedAt = time.Now()
	m.data[s.SchoolID] = s
	return nil
}

func setupSchoolSettingsRouter(repo *mockSchoolSettingsRepo, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := schoolsettings.NewService(repo)
	h := schoolsettings.NewHandler(svc)
	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", "user-1")
		c.Set("role", role)
		c.Set("school_id", schoolID)
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

// ─── Tests ─────────────────────────────────────────────────────────────────

// SS01 — Admin can GET school settings
func TestSchoolSettings_GetSettings_Admin(t *testing.T) {
	repo := &mockSchoolSettingsRepo{data: map[string]*schoolsettings.SchoolSettings{}}
	r := setupSchoolSettingsRouter(repo, "admin", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/school-settings", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var s schoolsettings.SchoolSettings
	err := json.Unmarshal(w.Body.Bytes(), &s)
	require.NoError(t, err)
	assert.Equal(t, "school-1", s.SchoolID)
}

// SS02 — Teacher can GET school settings (read-only)
func TestSchoolSettings_GetSettings_Teacher(t *testing.T) {
	repo := &mockSchoolSettingsRepo{data: map[string]*schoolsettings.SchoolSettings{}}
	r := setupSchoolSettingsRouter(repo, "teacher", "school-2")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/school-settings", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// SS03 — Admin can UPDATE school settings
func TestSchoolSettings_UpdateSettings_Admin(t *testing.T) {
	repo := &mockSchoolSettingsRepo{data: map[string]*schoolsettings.SchoolSettings{}}
	r := setupSchoolSettingsRouter(repo, "admin", "school-1")

	trueVal := true
	falseVal := false
	payload := schoolsettings.UpdateSchoolSettingsRequest{
		AllowParentsViewGrades:         &trueVal,
		AllowStudentsViewClassAverages: &falseVal,
		RequireMFAForStaff:             &trueVal,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/school-settings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	saved := repo.data["school-1"]
	require.NotNil(t, saved)
	assert.True(t, saved.AllowParentsViewGrades)
	assert.False(t, saved.AllowStudentsViewClassAverages)
	assert.True(t, saved.RequireMFAForStaff)
}

// SS04 — Teacher is forbidden from updating school settings
func TestSchoolSettings_UpdateSettings_ForbiddenForTeacher(t *testing.T) {
	repo := &mockSchoolSettingsRepo{data: map[string]*schoolsettings.SchoolSettings{}}
	r := setupSchoolSettingsRouter(repo, "teacher", "school-1")

	trueVal := true
	payload := schoolsettings.UpdateSchoolSettingsRequest{AllowParentsViewGrades: &trueVal}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/school-settings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// SS05 — Student is forbidden from updating school settings
func TestSchoolSettings_UpdateSettings_ForbiddenForStudent(t *testing.T) {
	repo := &mockSchoolSettingsRepo{data: map[string]*schoolsettings.SchoolSettings{}}
	r := setupSchoolSettingsRouter(repo, "student", "school-1")

	trueVal := true
	payload := schoolsettings.UpdateSchoolSettingsRequest{AllowParentsViewGrades: &trueVal}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/school-settings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// SS06 — Principal can update school settings
func TestSchoolSettings_UpdateSettings_AllowedForPrincipal(t *testing.T) {
	repo := &mockSchoolSettingsRepo{data: map[string]*schoolsettings.SchoolSettings{}}
	r := setupSchoolSettingsRouter(repo, "principal", "school-3")

	falseVal := false
	payload := schoolsettings.UpdateSchoolSettingsRequest{LockScrutinyEditingAfterValidation: &falseVal}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/school-settings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// SS07 — GET without school_id returns 400
func TestSchoolSettings_GetSettings_MissingSchoolID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockSchoolSettingsRepo{data: map[string]*schoolsettings.SchoolSettings{}}
	r := gin.New()
	svc := schoolsettings.NewService(repo)
	h := schoolsettings.NewHandler(svc)
	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", "user-1")
		c.Set("role", "admin")
		// Intentionally NOT setting school_id
		c.Next()
	})
	h.RegisterRoutes(api)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/school-settings", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// SS08 — Superadmin can override school_id via query parameter
func TestSchoolSettings_SuperadminOverrideSchoolID(t *testing.T) {
	repo := &mockSchoolSettingsRepo{
		data: map[string]*schoolsettings.SchoolSettings{
			"school-remote": {SchoolID: "school-remote", AllowParentsViewGrades: true, UpdatedAt: time.Now()},
		},
	}
	r := setupSchoolSettingsRouter(repo, "superadmin", "school-local")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/school-settings?school_id=school-remote", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var s schoolsettings.SchoolSettings
	_ = json.Unmarshal(w.Body.Bytes(), &s)
	assert.Equal(t, "school-remote", s.SchoolID)
}
