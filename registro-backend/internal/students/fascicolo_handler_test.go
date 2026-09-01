package students

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetFascicolo_PermissionsAndStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewFascicoloHandler(nil)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "student-123")
		c.Set("role", "student")
		c.Next()
	})
	handler.RegisterRoutes(r.Group("/"))

	req, _ := http.NewRequest("GET", "/students/student-123/fascicolo?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp StudentFascicolo
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "student-123", resp.StudentID)
	assert.Equal(t, 1, resp.Semester)
	assert.NotNil(t, resp.Voti)
	assert.NotNil(t, resp.Presenze)
	assert.NotNil(t, resp.Note)
	assert.NotNil(t, resp.PCTO)
	assert.NotNil(t, resp.Compiti)
	assert.NotNil(t, resp.Documenti)
}

func TestGetFascicolo_ForbiddenForOtherStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewFascicoloHandler(nil)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "student-123")
		c.Set("role", "student")
		c.Next()
	})
	handler.RegisterRoutes(r.Group("/"))

	req, _ := http.NewRequest("GET", "/students/other-student-456/fascicolo", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetFascicolo_AdminAndSecretaryAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewFascicoloHandler(nil)

	// Admin test
	rAdmin := gin.New()
	rAdmin.Use(func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Next()
	})
	handler.RegisterRoutes(rAdmin.Group("/"))

	reqAdmin, _ := http.NewRequest("GET", "/students/student-999/fascicolo", nil)
	wAdmin := httptest.NewRecorder()
	rAdmin.ServeHTTP(wAdmin, reqAdmin)
	assert.Equal(t, http.StatusOK, wAdmin.Code)

	// Secretary test
	rSec := gin.New()
	rSec.Use(func(c *gin.Context) {
		c.Set("user_id", "sec-1")
		c.Set("role", "secretary")
		c.Next()
	})
	handler.RegisterRoutes(rSec.Group("/"))

	reqSec, _ := http.NewRequest("GET", "/students/student-999/fascicolo", nil)
	wSec := httptest.NewRecorder()
	rSec.ServeHTTP(wSec, reqSec)
	assert.Equal(t, http.StatusOK, wSec.Code)
}
