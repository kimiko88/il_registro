package schools

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSchoolsSecurity_CreateSuperAdminOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepository)
	svc := NewService(repo)
	h := NewHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin") // regular admin, not superadmin
		c.Next()
	})
	h.RegisterRoutes(r.Group("/api"))

	body := bytes.NewBufferString(`{"name":"New School","code":"CODE123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/schools/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestSchoolsSecurity_AdminUpdateOtherSchoolForbidden(t *testing.T) {
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

	// Admin trying to update school-2
	body := bytes.NewBufferString(`{"name":"Hacked Name"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/schools/school-2", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
