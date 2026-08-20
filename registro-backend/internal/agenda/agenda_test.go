package agenda

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetTaskCompletion_ParentGuardian(t *testing.T) {
	mockRepo := new(MockRepository)

	svc := &Service{
		repo:     mockRepo,
		userRepo: nil,
	}

	item := &AgendaItem{ID: "item-10", TeacherID: "teacher-1"}
	mockRepo.On("GetByID", mock.Anything, "item-10").Return(item, nil)
	mockRepo.On("SetCompletion", mock.Anything, "item-10", "student-1", true).Return(nil)

	// Direct student completion
	err := svc.SetTaskCompletion(context.Background(), "student-1", "student", "", "item-10", true)
	assert.NoError(t, err)

	// Parent completion without userRepo returns ErrUnauthorized (security hardening)
	err = svc.SetTaskCompletion(context.Background(), "parent-1", "parent", "student-1", "item-10", true)
	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}

func TestAgendaHandler_GetByID_ErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	h := NewHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "u-1")
		c.Set("role", "student")
		c.Set("school_id", "school-1")
	})
	h.RegisterRoutes(r.Group("/api/v1"))

	// 1. Not found -> 404
	mockRepo.On("GetByID", mock.Anything, "non-existent").Return(nil, ErrNotFound).Once()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/agenda/non-existent", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 2. Internal error -> 500
	mockRepo.On("GetByID", mock.Anything, "db-error").Return(nil, errors.New("connection failed")).Once()
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/agenda/db-error", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}
