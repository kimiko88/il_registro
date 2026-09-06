package teachers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTeachersRouter(svc *Service, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		if role != "" {
			c.Set("role", role)
		}
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	handler := NewHandler(svc)
	handler.RegisterRoutes(r.Group("/api/v1"))
	return r
}

func TestGetSchoolID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. From Query
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test?school_id=school-from-query", nil)
	assert.Equal(t, "school-from-query", getSchoolID(c))

	// 2. From Header
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("X-School-ID", "school-from-header")
	assert.Equal(t, "school-from-header", getSchoolID(c))

	// 3. From Context
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Set("school_id", "school-from-context")
	assert.Equal(t, "school-from-context", getSchoolID(c))

	// 4. Default fallback
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	assert.Equal(t, "162737ff-081f-436c-8874-11cd57bc60f1", getSchoolID(c))
}

func TestHandler_Teachers_List_And_Get(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	// Unauthorized (no user_id)
	rAnon := setupTeachersRouter(svc, "", "", "")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/teachers", nil)
	w := httptest.NewRecorder()
	rAnon.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	rAuth := setupTeachersRouter(svc, "admin-1", "admin", "school-1")

	// 1. List with subject_id filter
	mockRepo.On("GetBySubject", mock.Anything, "sub-1").Return([]Teacher{{ID: "t-1", FirstName: "Giuseppe"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/teachers?subject_id=sub-1", nil)
	w = httptest.NewRecorder()
	rAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. List without filter
	mockRepo.On("List", mock.Anything, "school-1").Return([]Teacher{{ID: "t-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/teachers", nil)
	w = httptest.NewRecorder()
	rAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. Get single teacher
	mockRepo.On("Get", mock.Anything, "t-1").Return(&Teacher{ID: "t-1", SchoolID: "school-1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/teachers/t-1", nil)
	w = httptest.NewRecorder()
	rAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Get subjects
	mockRepo.On("GetSubjects", mock.Anything, "t-1").Return([]TeacherSubject{{ID: "ts-1", SubjectName: "Matematica"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/teachers/t-1/subjects", nil)
	w = httptest.NewRecorder()
	rAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockRepo.AssertExpectations(t)
}

func TestHandler_Teachers_Assign_Remove_Dashboard(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	rAuth := setupTeachersRouter(svc, "admin-1", "admin", "school-1")

	// 1. AssignSubject - success
	mockRepo.On("Get", mock.Anything, "t-1").Return(&Teacher{ID: "t-1", SchoolID: "school-1"}, nil).Once()
	mockRepo.On("AssignSubject", mock.Anything, "t-1", "sub-1").Return(nil).Once()
	body, _ := json.Marshal(AssignSubjectRequest{SubjectID: "sub-1"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/teachers/t-1/subjects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	rAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. AssignSubject - bad JSON
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/teachers/t-1/subjects", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 3. RemoveSubject - success
	mockRepo.On("Get", mock.Anything, "t-1").Return(&Teacher{ID: "t-1", SchoolID: "school-1"}, nil).Once()
	mockRepo.On("RemoveSubject", mock.Anything, "t-1", "sub-1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/teachers/t-1/subjects/sub-1", nil)
	w = httptest.NewRecorder()
	rAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// 4. RemoveSubject - service error
	mockRepo.On("Get", mock.Anything, "t-1").Return(&Teacher{ID: "t-1", SchoolID: "school-1"}, nil).Once()
	mockRepo.On("RemoveSubject", mock.Anything, "t-1", "sub-2").Return(errors.New("db error")).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/teachers/t-1/subjects/sub-2", nil)
	w = httptest.NewRecorder()
	rAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 5. GetDashboardStats - success
	mockRepo.On("GetDashboardStats", mock.Anything, "admin-1").Return(map[string]interface{}{"classes_count": 3, "students_count": 65}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/teachers/dashboard/stats", nil)
	w = httptest.NewRecorder()
	rAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockRepo.AssertExpectations(t)
}
