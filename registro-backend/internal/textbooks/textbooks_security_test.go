package textbooks_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/textbooks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTextbooksRepo struct {
	mock.Mock
}

func (m *mockTextbooksRepo) Create(ctx context.Context, t *textbooks.Textbook) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *mockTextbooksRepo) GetByID(ctx context.Context, id string) (*textbooks.Textbook, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*textbooks.Textbook), args.Error(1)
}
func (m *mockTextbooksRepo) Update(ctx context.Context, t *textbooks.Textbook) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *mockTextbooksRepo) List(ctx context.Context, schoolID string) ([]textbooks.Textbook, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]textbooks.Textbook), args.Error(1)
}
func (m *mockTextbooksRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockTextbooksRepo) AssignToClass(ctx context.Context, classID string, subjectID string, textbookID string, optional bool) error {
	args := m.Called(ctx, classID, subjectID, textbookID, optional)
	return args.Error(0)
}
func (m *mockTextbooksRepo) RemoveFromClass(ctx context.Context, assignmentID string) error {
	args := m.Called(ctx, assignmentID)
	return args.Error(0)
}
func (m *mockTextbooksRepo) ListByClass(ctx context.Context, classID string) ([]textbooks.ClassTextbook, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]textbooks.ClassTextbook), args.Error(1)
}

func buildTextbooksEngine(repo textbooks.Repository, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	if userID != "" {
		r.Use(func(c *gin.Context) {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
			c.Next()
		})
	}

	svc := textbooks.NewService(repo)
	h := textbooks.NewHandler(svc)
	h.RegisterRoutes(r.Group("/api"))

	return r
}

func TestTextbooks_List_TenantIsolation(t *testing.T) {
	t.Run("Regular user cannot bypass tenant via school_id query parameter", func(t *testing.T) {
		mockRepo := new(mockTextbooksRepo)
		// Should query "school-my" (from context), NOT "school-other" (from query parameter)
		mockRepo.On("List", mock.Anything, "school-my").Return([]textbooks.Textbook{}, nil).Once()

		r := buildTextbooksEngine(mockRepo, "user-1", "teacher", "school-my")
		req := httptest.NewRequest(http.MethodGet, "/api/textbooks?school_id=school-other", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Superadmin CAN specify target school_id query parameter", func(t *testing.T) {
		mockRepo := new(mockTextbooksRepo)
		mockRepo.On("List", mock.Anything, "school-other").Return([]textbooks.Textbook{}, nil).Once()

		r := buildTextbooksEngine(mockRepo, "superadmin-1", "superadmin", "school-default")
		req := httptest.NewRequest(http.MethodGet, "/api/textbooks?school_id=school-other", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})
}
