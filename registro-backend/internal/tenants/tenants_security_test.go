package tenants

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTenantsSecurity_GetTenantCrossSchoolForbidden(t *testing.T) {
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

	// Admin of school-1 attempts to get details of school-2
	req := httptest.NewRequest(http.MethodGet, "/api/tenants/school-2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestTenantsSecurity_GetTenantOwnSchoolAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepository)
	repo.On("GetTenantByID", mock.Anything, "school-1").Return(&Tenant{ID: "school-1", Name: "My School"}, nil)

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

	req := httptest.NewRequest(http.MethodGet, "/api/tenants/school-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
