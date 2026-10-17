package may15

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMay15Router(svc Service, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "school-1")
		c.Set("user_id", "teacher-coord-1")
		c.Set("role", role)
		c.Next()
	})
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestMay15Handler_GetAndSaveDocument(t *testing.T) {
	repo := new(MockMay15Repo)
	svc := NewService(repo)
	r := setupMay15Router(svc, "coordinator")

	// 1. GetDocument
	expectedDoc := &ClassMay15Document{
		ID:           "doc-5a",
		ClassID:      "class-5a",
		AcademicYear: "2025/2026",
		Status:       StatusBozza,
	}
	repo.On("GetByClassAndYear", mock.Anything, "class-5a", "2025/2026").Return(expectedDoc, nil).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/may15/class/class-5a?academic_year=2025/2026", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. SaveDocument
	repo.On("Upsert", mock.Anything, mock.Anything).Return(nil).Once()
	saveReq := SaveMay15Request{
		AcademicYear:      "2025/2026",
		ClassPresentation: "Ottima classe",
		Status:            StatusApprovatoCdC,
	}
	body, _ := json.Marshal(saveReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/may15/class/class-5a", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. PublishDocument
	pubDoc := &ClassMay15Document{
		ID:           "doc-5a",
		ClassID:      "class-5a",
		AcademicYear: "2025/2026",
		Status:       StatusPubblicato,
	}
	repo.On("Publish", mock.Anything, "class-5a", "2025/2026").Return(pubDoc, nil).Once()
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/may15/class/class-5a/publish?academic_year=2025/2026", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
