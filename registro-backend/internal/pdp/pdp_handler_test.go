package pdp_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/pdp"
)

type mockPdpRepo struct {
	mock.Mock
}

func (m *mockPdpRepo) Create(ctx context.Context, plan *pdp.PdpPlan) (*pdp.PdpPlan, error) {
	args := m.Called(ctx, plan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pdp.PdpPlan), args.Error(1)
}

func (m *mockPdpRepo) GetByID(ctx context.Context, id string) (*pdp.PdpPlan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pdp.PdpPlan), args.Error(1)
}

func (m *mockPdpRepo) GetByStudent(ctx context.Context, studentID, year string) ([]*pdp.PdpPlan, error) {
	args := m.Called(ctx, studentID, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*pdp.PdpPlan), args.Error(1)
}

func (m *mockPdpRepo) GetByClass(ctx context.Context, classID, year string) ([]*pdp.PdpPlan, error) {
	args := m.Called(ctx, classID, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*pdp.PdpPlan), args.Error(1)
}

func (m *mockPdpRepo) Update(ctx context.Context, id string, req *pdp.UpdatePdpRequest) (*pdp.PdpPlan, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pdp.PdpPlan), args.Error(1)
}

func (m *mockPdpRepo) SetSharedWithFamily(ctx context.Context, id string, shared bool) error {
	args := m.Called(ctx, id, shared)
	return args.Error(0)
}

func (m *mockPdpRepo) ApproveByFamily(ctx context.Context, id, approverID string) error {
	args := m.Called(ctx, id, approverID)
	return args.Error(0)
}

func (m *mockPdpRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestPdpRouter(svc *pdp.Service, actorID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if actorID != "" {
			c.Set(pdp.ContextKeyUserID, actorID)
		}
		if role != "" {
			c.Set(pdp.ContextKeyUserRole, role)
		}
		if schoolID != "" {
			c.Set(pdp.ContextKeySchoolID, schoolID)
		}
		c.Next()
	})
	h := pdp.NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestPDPHandler_Endpoints(t *testing.T) {
	repo := new(mockPdpRepo)
	svc := pdp.NewService(repo)

	// Router for teacher/coordinator
	r := setupTestPdpRouter(svc, "teacher-1", "teacher", "sch-1")

	// 1. Unauthorized check (when no actorID in context)
	unauthRouter := setupTestPdpRouter(svc, "", "", "")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/pdp/student/s1", nil)
	w := httptest.NewRecorder()
	unauthRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 2. GET /pdp/student/:studentId
	repo.On("GetByStudent", mock.Anything, "s1", "").Return([]*pdp.PdpPlan{{ID: "plan-1", SchoolID: "sch-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pdp/student/s1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. GET /pdp/class/:classId
	repo.On("GetByClass", mock.Anything, "cls-1", "").Return([]*pdp.PdpPlan{{ID: "plan-1", SchoolID: "sch-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pdp/class/cls-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. GET /pdp/:id - Found & Not Found
	repo.On("GetByID", mock.Anything, "plan-1").Return(&pdp.PdpPlan{ID: "plan-1", SchoolID: "sch-1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pdp/plan-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetByID", mock.Anything, "missing").Return(nil, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pdp/missing", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 5. POST /pdp (Create)
	repo.On("Create", mock.Anything, mock.Anything).Return(&pdp.PdpPlan{ID: "new-plan", SchoolID: "sch-1"}, nil).Once()
	createPayload := `{"student_id":"s1","class_id":"cls-1","academic_year":"2025/2026","plan_type":"pdp","diagnosis":"Dislessia"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pdp", bytes.NewBufferString(createPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pdp", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 6. PUT /pdp/:id (Update)
	repo.On("GetByID", mock.Anything, "plan-1").Return(&pdp.PdpPlan{ID: "plan-1", SchoolID: "sch-1"}, nil).Once()
	repo.On("Update", mock.Anything, "plan-1", mock.Anything).Return(&pdp.PdpPlan{ID: "plan-1", SchoolID: "sch-1"}, nil).Once()
	updPayload := `{"diagnosis":"Aggiornamento diagnosi"}`
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/pdp/plan-1", bytes.NewBufferString(updPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. DELETE /pdp/:id
	repo.On("GetByID", mock.Anything, "plan-1").Return(&pdp.PdpPlan{ID: "plan-1", SchoolID: "sch-1"}, nil).Once()
	repo.On("Delete", mock.Anything, "plan-1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/pdp/plan-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// 8. POST /pdp/:id/share
	repo.On("GetByID", mock.Anything, "plan-1").Return(&pdp.PdpPlan{ID: "plan-1", SchoolID: "sch-1"}, nil).Twice()
	repo.On("SetSharedWithFamily", mock.Anything, "plan-1", true).Return(nil).Once()
	sharePayload := `{"share":true}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pdp/plan-1/share", bytes.NewBufferString(sharePayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 9. POST /pdp/:id/approve - Forbidden for teacher, allowed for parent
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pdp/plan-1/approve", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	parentRouter := setupTestPdpRouter(svc, "parent-1", "parent", "sch-1")
	repo.On("GetByID", mock.Anything, "plan-1").Return(&pdp.PdpPlan{ID: "plan-1", SchoolID: "sch-1", SharedWithFamily: true}, nil).Once()
	repo.On("ApproveByFamily", mock.Anything, "plan-1", "parent-1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pdp/plan-1/approve", nil)
	w = httptest.NewRecorder()
	parentRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 10. ListCompensative & ListDispensative
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pdp/measures/compensative", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pdp/measures/dispensative", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.AssertExpectations(t)
}

func TestPDPHandler_ErrorMappings(t *testing.T) {
	repo := new(mockPdpRepo)
	svc := pdp.NewService(repo)
	parentRouter := setupTestPdpRouter(svc, "parent-1", "parent", "sch-1")

	// Not shared yet -> 403 Forbidden
	repo.On("GetByID", mock.Anything, "p-not-shared").Return(&pdp.PdpPlan{ID: "p-not-shared", SchoolID: "sch-1", SharedWithFamily: false}, nil).Once()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/pdp/p-not-shared/approve", nil)
	w := httptest.NewRecorder()
	parentRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Already approved -> 409 Conflict
	now := time.Now()
	repo.On("GetByID", mock.Anything, "p-approved").Return(&pdp.PdpPlan{ID: "p-approved", SchoolID: "sch-1", SharedWithFamily: true, FamilyApprovedAt: &now}, nil).Once()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/pdp/p-approved/approve", nil)
	w = httptest.NewRecorder()
	parentRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)

	// Database internal error -> 500 Internal Server Error
	repo.On("GetByID", mock.Anything, "p-err").Return(nil, errors.New("db connection lost")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/pdp/p-err", nil)
	w = httptest.NewRecorder()
	parentRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	repo.AssertExpectations(t)
}
