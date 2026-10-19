package groups

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter(svc *Service, userID, role, schoolID string) *gin.Engine {
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
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestGroupsService_ExtendedEdgeCases(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()

	// ListGroupsBySchool empty schoolID
	_, err := svc.ListGroupsBySchool(ctx, "")
	assert.Error(t, err)

	// ListGroupsByTeacher empty teacherID
	_, err = svc.ListGroupsByTeacher(ctx, "")
	assert.Error(t, err)

	// ListGroupsByStudent empty studentID
	_, err = svc.ListGroupsByStudent(ctx, "")
	assert.Error(t, err)

	// ListGroupsBySchool success
	repo.On("ListBySchool", ctx, "sch-1").Return([]Group{{ID: "g1", Name: "Coro"}}, nil).Once()
	res, err := svc.ListGroupsBySchool(ctx, "sch-1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	// ListGroupsByTeacher success
	repo.On("ListByTeacher", ctx, "teach-1").Return([]Group{{ID: "g1", Name: "Teatro"}}, nil).Once()
	res, err = svc.ListGroupsByTeacher(ctx, "teach-1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	// ListGroupsByStudent success
	repo.On("ListByStudent", ctx, "stud-1").Return([]Group{{ID: "g1", Name: "Robotica"}}, nil).Once()
	res, err = svc.ListGroupsByStudent(ctx, "stud-1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	repo.AssertExpectations(t)
}

func TestGroupsHandler_Endpoints(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	// 1. Forbidden for unauthorized roles (e.g. student or missing user_id)
	unauthRouter := setupTestRouter(svc, "stud-1", "student", "sch-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	unauthRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Authorized router for admin
	adminRouter := setupTestRouter(svc, "admin-1", "admin", "sch-1")

	// 2. POST /groups - Success & bad json
	repo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
	repo.On("GetByID", mock.Anything, mock.Anything).Return(&Group{ID: "g1", Name: "Orchestra", SchoolID: "sch-1"}, nil).Once()
	createPayload := `{"name":"Orchestra","academic_year":"2025/2026"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewBufferString(createPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 3. GET /groups
	repo.On("ListBySchool", mock.Anything, "sch-1").Return([]Group{{ID: "g1", Name: "Orchestra"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/groups", nil)
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /groups with teacher_id query
	repo.On("ListByTeacher", mock.Anything, "teach-1").Return([]Group{{ID: "g1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/groups?teacher_id=teach-1", nil)
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /groups with student_id query
	repo.On("ListByStudent", mock.Anything, "stud-1").Return([]Group{{ID: "g1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/groups?student_id=stud-1", nil)
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. GET /groups/:id - Found & Not Found
	repo.On("GetByID", mock.Anything, "g1").Return(&Group{ID: "g1", Name: "Orchestra"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/groups/g1", nil)
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetByID", mock.Anything, "missing").Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/groups/missing", nil)
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 5. PUT /groups/:id (Update)
	repo.On("GetByID", mock.Anything, "g1").Return(&Group{ID: "g1", SchoolID: "sch-1", Name: "Orchestra"}, nil).Twice()
	repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
	updPayload := `{"name":"Orchestra d'Istituto","description":"Aggiornato"}`
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/groups/g1", bytes.NewBufferString(updPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodPut, "/api/v1/groups/g1", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 6. DELETE /groups/:id
	repo.On("GetByID", mock.Anything, "g1").Return(&Group{ID: "g1", SchoolID: "sch-1"}, nil).Once()
	repo.On("Delete", mock.Anything, "g1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/groups/g1", nil)
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. POST /groups/:id/students (Add)
	repo.On("GetByID", mock.Anything, "g1").Return(&Group{ID: "g1", SchoolID: "sch-1"}, nil).Once()
	repo.On("AddStudents", mock.Anything, "g1", []string{"s1", "s2"}).Return(nil).Once()
	addPayload := `{"student_ids":["s1","s2"]}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/groups/g1/students", bytes.NewBufferString(addPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 8. DELETE /groups/:id/students/:student_id (Remove)
	repo.On("GetByID", mock.Anything, "g1").Return(&Group{ID: "g1", SchoolID: "sch-1"}, nil).Once()
	repo.On("RemoveStudent", mock.Anything, "g1", "s1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/groups/g1/students/s1", nil)
	w = httptest.NewRecorder()
	adminRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.AssertExpectations(t)
}
