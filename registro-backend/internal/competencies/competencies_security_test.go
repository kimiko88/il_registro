package competencies

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestCompetenciesSecurity_Unauthenticated(t *testing.T) {
	svc := NewService(nil)
	h := NewHandler(svc)

	r := setupTestRouter()
	r.Use(func(c *gin.Context) {
		// No user_id set
		c.Next()
	})
	h.RegisterRoutes(r.Group("/api"))

	// 1. GetClassEvaluations unauthenticated
	req := httptest.NewRequest(http.MethodGet, "/api/competencies?class_id=cls1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 2. GetStudentEvaluations unauthenticated
	req = httptest.NewRequest(http.MethodGet, "/api/competencies/student/stu1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 3. SaveEvaluation unauthenticated
	req = httptest.NewRequest(http.MethodPost, "/api/competencies/evaluations", bytes.NewBufferString(`{}`))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCompetenciesSecurity_StudentAccessOther(t *testing.T) {
	svc := NewService(nil)
	h := NewHandler(svc)

	r := setupTestRouter()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "student-1")
		c.Set("role", "student")
		c.Next()
	})
	h.RegisterRoutes(r.Group("/api"))

	req := httptest.NewRequest(http.MethodGet, "/api/competencies/student/student-2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCompetenciesSecurity_BatchSaveInvalidJSON(t *testing.T) {
	svc := NewService(nil)
	h := NewHandler(svc)

	r := setupTestRouter()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		c.Next()
	})
	h.RegisterRoutes(r.Group("/api"))

	req := httptest.NewRequest(http.MethodPost, "/api/competencies/batch", bytes.NewBufferString(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
