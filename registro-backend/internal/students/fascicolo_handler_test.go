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
