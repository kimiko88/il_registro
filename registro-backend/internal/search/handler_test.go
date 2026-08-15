package search

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Returns 401 when unauthenticated", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		h := NewHandler(svc)

		r := gin.New()
		h.RegisterRoutes(r.Group("/api/v1"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Returns 400 when school_id is missing for non-superadmin", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		h := NewHandler(svc)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_id", "user-1")
			c.Set("role", "teacher")
			c.Next()
		})
		h.RegisterRoutes(r.Group("/api/v1"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Returns 500 when service fails", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("GlobalSearch", mock.Anything, "teacher", "school-1", "test", "").Return(nil, errors.New("db down")).Once()

		svc := NewService(mockRepo)
		h := NewHandler(svc)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_id", "user-1")
			c.Set("school_id", "school-1")
			c.Set("role", "teacher")
			c.Next()
		})
		h.RegisterRoutes(r.Group("/api/v1"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
