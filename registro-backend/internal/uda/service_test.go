package uda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUdaRepository struct {
	mock.Mock
}

func (m *MockUdaRepository) Create(ctx context.Context, plan *UdaPlan) error {
	args := m.Called(ctx, plan)
	return args.Error(0)
}

func (m *MockUdaRepository) ListByClass(ctx context.Context, classID string) ([]*UdaPlan, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*UdaPlan), args.Error(1)
}

func (m *MockUdaRepository) GetByID(ctx context.Context, id string) (*UdaPlan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UdaPlan), args.Error(1)
}

func (m *MockUdaRepository) ListBySchool(ctx context.Context, schoolID string) ([]*UdaPlan, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*UdaPlan), args.Error(1)
}

func (m *MockUdaRepository) ListAll(ctx context.Context) ([]*UdaPlan, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*UdaPlan), args.Error(1)
}

func (m *MockUdaRepository) Update(ctx context.Context, id string, req UpdateUdaRequest) (*UdaPlan, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UdaPlan), args.Error(1)
}

func (m *MockUdaRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUdaPlanModel(t *testing.T) {
	plan := UdaPlan{
		ID:           "uda-1",
		Title:        "Unità di Apprendimento 1: La Rivoluzione Industriale",
		Period:       "primo_quadrimestre",
		Competencies: []string{"COMP_CIVIC", "COMP_DIGITAL"},
		Status:       "draft",
	}

	assert.NotEmpty(t, plan.Title)
	assert.Len(t, plan.Competencies, 2)
}

func TestService_CreateUda(t *testing.T) {
	mockRepo := new(MockUdaRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	// 1. Success case
	mockRepo.On("Create", ctx, mock.MatchedBy(func(p *UdaPlan) bool {
		return p.Title == "Storia Contemporanea" && p.Period == "annuale"
	})).Return(nil).Once()

	req := CreateUdaRequest{
		ClassID:   "c-1",
		SubjectID: "s-1",
		Title:     "Storia Contemporanea",
		StartDate: "2026-10-01",
		EndDate:   "2026-12-01",
	}
	plan, err := svc.CreateUda(ctx, "school-1", "teacher-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, "draft", plan.Status)
	assert.Equal(t, "annuale", plan.Period)

	// 2. Error case: end date before start date
	invalidReq := CreateUdaRequest{
		Title:     "Errata",
		StartDate: "2026-12-01",
		EndDate:   "2026-10-01",
	}
	_, err = svc.CreateUda(ctx, "school-1", "teacher-1", invalidReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "end_date cannot be before start_date")

	// 3. Error case: repo error
	mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("db insert fail")).Once()
	req2 := CreateUdaRequest{Title: "Fail DB"}
	_, err = svc.CreateUda(ctx, "school-1", "teacher-1", req2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create UDA plan")

	mockRepo.AssertExpectations(t)
}

func TestService_ListAll_And_ListByClass(t *testing.T) {
	mockRepo := new(MockUdaRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	plans := []*UdaPlan{{ID: "uda-1", Title: "UDA 1"}}

	// Superadmin sees all
	mockRepo.On("ListAll", ctx).Return(plans, nil).Once()
	res, err := svc.ListAll(ctx, "school-1", "superadmin")
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	// Regular user sees school only
	mockRepo.On("ListBySchool", ctx, "school-1").Return(plans, nil).Once()
	res, err = svc.ListAll(ctx, "school-1", "teacher")
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	// ListByClass
	mockRepo.On("ListByClass", ctx, "class-10").Return(plans, nil).Once()
	res, err = svc.ListByClass(ctx, "class-10")
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	mockRepo.AssertExpectations(t)
}

func TestService_UpdateUda(t *testing.T) {
	mockRepo := new(MockUdaRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	existing := &UdaPlan{
		ID:        "uda-1",
		SchoolID:  "school-1",
		TeacherID: "teacher-1",
	}
	req := UpdateUdaRequest{Title: "Titolo Aggiornato"}

	// 1. Not found
	mockRepo.On("GetByID", ctx, "missing").Return(nil, errors.New("not found")).Once()
	_, err := svc.UpdateUda(ctx, "teacher-1", "teacher", "school-1", "missing", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uda plan non trovato")

	// 2. Cross school forbidden
	mockRepo.On("GetByID", ctx, "uda-1").Return(existing, nil).Once()
	_, err = svc.UpdateUda(ctx, "teacher-1", "teacher", "school-2", "uda-1", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "appartiene ad un'altra scuola")

	// 3. Not owner forbidden
	mockRepo.On("GetByID", ctx, "uda-1").Return(existing, nil).Once()
	_, err = svc.UpdateUda(ctx, "other-teacher", "teacher", "school-1", "uda-1", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non sei il docente proprietario")

	// 4. Admin bypass ownership
	updated := &UdaPlan{ID: "uda-1", Title: "Titolo Aggiornato"}
	mockRepo.On("GetByID", ctx, "uda-1").Return(existing, nil).Once()
	mockRepo.On("Update", ctx, "uda-1", req).Return(updated, nil).Once()
	res, err := svc.UpdateUda(ctx, "admin-1", "admin", "school-1", "uda-1", req)
	assert.NoError(t, err)
	assert.Equal(t, "Titolo Aggiornato", res.Title)

	// 5. Superadmin bypass school & ownership
	mockRepo.On("GetByID", ctx, "uda-1").Return(existing, nil).Once()
	mockRepo.On("Update", ctx, "uda-1", req).Return(updated, nil).Once()
	res, err = svc.UpdateUda(ctx, "super-1", "superadmin", "", "uda-1", req)
	assert.NoError(t, err)
	assert.Equal(t, "Titolo Aggiornato", res.Title)

	mockRepo.AssertExpectations(t)
}

func TestService_DeleteUda(t *testing.T) {
	mockRepo := new(MockUdaRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	existing := &UdaPlan{
		ID:        "uda-1",
		SchoolID:  "school-1",
		TeacherID: "teacher-1",
	}

	// 1. Not found
	mockRepo.On("GetByID", ctx, "missing").Return(nil, errors.New("not found")).Once()
	err := svc.DeleteUda(ctx, "teacher-1", "teacher", "school-1", "missing")
	assert.Error(t, err)

	// 2. Cross school forbidden
	mockRepo.On("GetByID", ctx, "uda-1").Return(existing, nil).Once()
	err = svc.DeleteUda(ctx, "teacher-1", "teacher", "school-2", "uda-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "appartiene ad un'altra scuola")

	// 3. Not owner forbidden
	mockRepo.On("GetByID", ctx, "uda-1").Return(existing, nil).Once()
	err = svc.DeleteUda(ctx, "other-teacher", "teacher", "school-1", "uda-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non sei il docente proprietario")

	// 4. Success delete
	mockRepo.On("GetByID", ctx, "uda-1").Return(existing, nil).Once()
	mockRepo.On("Delete", ctx, "uda-1").Return(nil).Once()
	err = svc.DeleteUda(ctx, "teacher-1", "teacher", "school-1", "uda-1")
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func setupUdaRouter(svc *Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("school_id", "school-1")
		c.Set("role", "teacher")
		c.Next()
	})
	handler := NewHandler(svc)
	rg := r.Group("/api/v1")
	handler.RegisterRoutes(rg)
	return r
}

func TestHandler_UdaRoutes(t *testing.T) {
	mockRepo := new(MockUdaRepository)
	svc := NewService(mockRepo)
	router := setupUdaRouter(svc)

	// 1. ListAll
	mockRepo.On("ListBySchool", mock.Anything, "school-1").Return([]*UdaPlan{{ID: "uda-1"}}, nil).Once()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/uda", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. ListByClass
	mockRepo.On("ListByClass", mock.Anything, "c-10").Return([]*UdaPlan{{ID: "uda-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/uda/class/c-10", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. CreateUda - Unauthorized (bare router with no user_id)
	bareRouter := gin.New()
	bareHandler := NewHandler(svc)
	bareHandler.RegisterRoutes(bareRouter.Group("/api/v1"))
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/uda", bytes.NewBufferString(`{}`))
	w = httptest.NewRecorder()
	bareRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 4. CreateUda - Success with middleware context
	r2 := gin.New()
	handler := NewHandler(svc)
	r2.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("school_id", "school-1")
		c.Set("role", "teacher")
		c.Next()
	})
	handler.RegisterRoutes(r2.Group("/api/v1"))

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
	body, _ := json.Marshal(CreateUdaRequest{ClassID: "c-1", SubjectID: "s-1", Title: "Nuova UDA"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/uda", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r2.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 5. UpdateUda - Success
	mockRepo.On("GetByID", mock.Anything, "uda-1").Return(&UdaPlan{ID: "uda-1", SchoolID: "school-1", TeacherID: "teacher-1"}, nil).Once()
	mockRepo.On("Update", mock.Anything, "uda-1", mock.Anything).Return(&UdaPlan{ID: "uda-1", Title: "Aggiornato"}, nil).Once()
	body, _ = json.Marshal(UpdateUdaRequest{Title: "Aggiornato"})
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/uda/uda-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r2.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. DeleteUda - Success
	mockRepo.On("GetByID", mock.Anything, "uda-1").Return(&UdaPlan{ID: "uda-1", SchoolID: "school-1", TeacherID: "teacher-1"}, nil).Once()
	mockRepo.On("Delete", mock.Anything, "uda-1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/uda/uda-1", nil)
	w = httptest.NewRecorder()
	r2.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. CreateUda - Forbidden role (e.g. student)
	rForbidden := gin.New()
	rForbidden.Use(func(c *gin.Context) {
		c.Set("user_id", "student-1")
		c.Set("role", "student")
		c.Next()
	})
	handler.RegisterRoutes(rForbidden.Group("/api/v1"))
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/uda", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rForbidden.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 8. CreateUda - Bad JSON
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/uda", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r2.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 9. UpdateUda - Bad JSON
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/uda/uda-1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r2.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 10. UpdateUda - Service error
	mockRepo.On("GetByID", mock.Anything, "uda-err").Return(nil, errors.New("db error")).Once()
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/uda/uda-err", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r2.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 11. DeleteUda - Service error
	mockRepo.On("GetByID", mock.Anything, "uda-err").Return(nil, errors.New("db error")).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/uda/uda-err", nil)
	w = httptest.NewRecorder()
	r2.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 12. ListAll - Service error
	mockRepo.On("ListBySchool", mock.Anything, "school-1").Return(nil, errors.New("list error")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/uda", nil)
	w = httptest.NewRecorder()
	r2.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 13. ListByClass - Service error
	mockRepo.On("ListByClass", mock.Anything, "c-err").Return(nil, errors.New("class error")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/uda/class/c-err", nil)
	w = httptest.NewRecorder()
	r2.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockRepo.AssertExpectations(t)
}
