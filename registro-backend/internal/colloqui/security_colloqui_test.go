package colloqui

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


type mockColloquiRepo struct {
	mock.Mock
	Repository
}

func (m *mockColloquiRepo) GetSlotByID(ctx context.Context, id string) (*ColloquioSlot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ColloquioSlot), args.Error(1)
}

func (m *mockColloquiRepo) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *mockColloquiRepo) CancelSlot(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestColloqui_CancelSlot_CrossTenantAdminRejection(t *testing.T) {
	repo := new(mockColloquiRepo)
	svc := NewService(repo)

	slotID := "slot-school-A"
	repo.On("GetSlotByID", mock.Anything, slotID).Return(&ColloquioSlot{
		ID:        slotID,
		SchoolID:  "school-A",
		TeacherID: "teacher-1",
	}, nil)

	// Admin with empty actorSchoolID MUST be rejected
	errEmptySchool := svc.CancelSlot(context.Background(), "admin-1", "admin", "", slotID)
	assert.ErrorIs(t, errEmptySchool, ErrUnauthorized)

	// Admin with mismatched schoolID MUST be rejected
	errDiffSchool := svc.CancelSlot(context.Background(), "admin-1", "admin", "school-B", slotID)
	assert.ErrorIs(t, errDiffSchool, ErrUnauthorized)

	repo.AssertExpectations(t)
}

func TestColloqui_ListSlots_RequiresFilter(t *testing.T) {
	repo := new(mockColloquiRepo)
	svc := NewService(repo)

	_, err := svc.ListSlots(context.Background(), "", "", time.Time{}, time.Time{}, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "obbligatorio")
}

func (m *mockColloquiRepo) BookQueueTicket(ctx context.Context, t *GeneralMeetingQueueTicket) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}


func (m *mockColloquiRepo) ListQueueTickets(ctx context.Context, meetingID, teacherID, parentID string) ([]GeneralMeetingQueueTicket, error) {
	args := m.Called(ctx, meetingID, teacherID, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]GeneralMeetingQueueTicket), args.Error(1)
}

func TestColloqui_QueueTicket_Security(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Unauthenticated BookQueueTicket returns 401", func(t *testing.T) {
		r := gin.New()
		h := NewHandler(NewService(new(mockColloquiRepo)))
		h.RegisterRoutes(r.Group("/api"))

		req := httptest.NewRequest("POST", "/api/colloqui/general-meetings/book-ticket", bytes.NewReader([]byte(`{}`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Student BookQueueTicket returns 403", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_id", "student-1")
			c.Set("role", "student")
			c.Next()
		})
		h := NewHandler(NewService(new(mockColloquiRepo)))
		h.RegisterRoutes(r.Group("/api"))

		req := httptest.NewRequest("POST", "/api/colloqui/general-meetings/book-ticket", bytes.NewReader([]byte(`{"meeting_id":"m-1","teacher_id":"t-1"}`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Parent ListQueueTickets forces parentID to current user (Anti-IDOR)", func(t *testing.T) {
		mockRepo := new(mockColloquiRepo)
		mockRepo.On("ListQueueTickets", mock.Anything, "m-1", "", "parent-me").Return([]GeneralMeetingQueueTicket{}, nil).Once()

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_id", "parent-me")
			c.Set("role", "parent")
			c.Next()
		})
		h := NewHandler(NewService(mockRepo))
		h.RegisterRoutes(r.Group("/api"))

		// Even if parent specifies parent_id=victim-parent, the handler must force parent_id=parent-me
		req := httptest.NewRequest("GET", "/api/colloqui/general-meetings/m-1/tickets?parent_id=victim-parent", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})
}


