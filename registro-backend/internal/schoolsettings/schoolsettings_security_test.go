package schoolsettings

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSchoolSettingsSecurity_AdminCannotOverrideSchoolIDViaQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepository)
	svc := NewService(repo)
	h := NewHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		c.Next()
	})
	h.RegisterRoutes(r.Group("/api"))

	initialSettings := &SchoolSettings{
		SchoolID:  "school-1",
		UpdatedAt: time.Now(),
	}
	repo.On("GetBySchoolID", mock.Anything, "school-1").Return(initialSettings, nil)

	// Admin attempts to query school-2 via ?school_id=school-2
	req := httptest.NewRequest(http.MethodGet, "/api/school-settings?school_id=school-2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Must have fetched school-1, NOT school-2
	repo.AssertCalled(t, "GetBySchoolID", mock.Anything, "school-1")
	repo.AssertNotCalled(t, "GetBySchoolID", mock.Anything, "school-2")
}

func TestSchoolSettingsSecurity_UpdateForbiddenForTeacher(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepository)
	svc := NewService(repo)
	h := NewHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		c.Next()
	})
	h.RegisterRoutes(r.Group("/api"))

	body := bytes.NewBufferString(`{"allow_parents_view_grades":false}`)
	req := httptest.NewRequest(http.MethodPut, "/api/school-settings", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
