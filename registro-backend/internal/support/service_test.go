package support

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSupportRepo struct {
	mock.Mock
}

func (m *MockSupportRepo) CreateDiaryEntry(ctx context.Context, entry *SupportDiaryEntry) error {
	args := m.Called(ctx, entry)
	entry.ID = "entry-1"
	return args.Error(0)
}

func (m *MockSupportRepo) ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]SupportDiaryEntry, error) {
	args := m.Called(ctx, schoolID, studentID, teacherID, classID, isFamily)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SupportDiaryEntry), args.Error(1)
}

func (m *MockSupportRepo) DeleteDiaryEntry(ctx context.Context, id, teacherID string) error {
	args := m.Called(ctx, id, teacherID)
	return args.Error(0)
}

func (m *MockSupportRepo) CreatePeiGoal(ctx context.Context, goal *SupportPeiGoal) error {
	args := m.Called(ctx, goal)
	goal.ID = "goal-1"
	return args.Error(0)
}

func (m *MockSupportRepo) UpdatePeiGoalProgress(ctx context.Context, id, progressStatus string) error {
	args := m.Called(ctx, id, progressStatus)
	return args.Error(0)
}

func (m *MockSupportRepo) ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]SupportPeiGoal, error) {
	args := m.Called(ctx, schoolID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SupportPeiGoal), args.Error(1)
}

func (m *MockSupportRepo) DeletePeiGoal(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSupportService_CreateDiaryEntry(t *testing.T) {
	repo := new(MockSupportRepo)
	svc := NewService(repo)

	req := CreateDiaryEntryRequest{
		StudentID:          "s-1",
		ClassID:            "c-1",
		EntryDate:          "2025-03-10",
		TimeSlot:           "1ª Ora (08:00-09:00)",
		ActivityType:       "in_classe",
		TopicAndActivities: "Attività di supporto in matematica sulle equazioni con materiale visivo",
		StudentResponses:   "Buona concentrazione e partecipazione",
		EducatorNotes:      "L'educatore ha fornito supporto durante la seconda parte dell'ora",
		IsSharedWithFamily: true,
	}

	repo.On("CreateDiaryEntry", mock.Anything, mock.MatchedBy(func(e *SupportDiaryEntry) bool {
		return e.StudentID == "s-1" && e.TopicAndActivities == req.TopicAndActivities && e.IsSharedWithFamily
	})).Return(nil).Once()

	entry, err := svc.CreateDiaryEntry(context.Background(), "school-1", "t-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, "entry-1", entry.ID)
	assert.Equal(t, true, entry.IsSharedWithFamily)

	// Validation on empty topic
	reqEmpty := req
	reqEmpty.TopicAndActivities = ""
	_, errEmpty := svc.CreateDiaryEntry(context.Background(), "school-1", "t-1", reqEmpty)
	assert.ErrorIs(t, errEmpty, ErrEmptyTopic)
}

func TestSupportService_PeiGoals(t *testing.T) {
	repo := new(MockSupportRepo)
	svc := NewService(repo)

	req := CreatePeiGoalRequest{
		StudentID:      "s-1",
		PeiType:        "equipollente",
		Axis:           "cognitiva",
		Title:          "Risoluzione problemi lineari guidati",
		Description:    "Sviluppo capacità di comprensione del testo e impostazione equazione",
		ExpectedTerm:   "q1",
		ProgressStatus: "iniziale",
	}

	repo.On("CreatePeiGoal", mock.Anything, mock.MatchedBy(func(g *SupportPeiGoal) bool {
		return g.StudentID == "s-1" && g.Axis == "cognitiva" && g.PeiType == "equipollente"
	})).Return(nil).Once()

	goal, err := svc.CreatePeiGoal(context.Background(), "school-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, goal)
	assert.Equal(t, "goal-1", goal.ID)

	repo.On("UpdatePeiGoalProgress", mock.Anything, "goal-1", "raggiunto").Return(nil).Once()
	err = svc.UpdateGoalProgress(context.Background(), "goal-1", "raggiunto")
	assert.NoError(t, err)
}

type MockUsersRepo struct {
	mock.Mock
	users.Repository
}

func (m *MockUsersRepo) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

func TestSupportHandler_Security(t *testing.T) {
	repo := new(MockSupportRepo)
	mockUsers := new(MockUsersRepo)
	svc := NewService(repo, mockUsers)
	h := NewHandler(svc)

	t.Run("CreateDiaryEntry_ForbiddenForStudent", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "student-1")
		c.Set("role", "student")
		c.Request, _ = http.NewRequest("POST", "/support/diaries", strings.NewReader(`{"topic_and_activities":"test"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateDiaryEntry(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("ListPeiGoals_ParentWithoutStudentID_BadRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "parent-1")
		c.Set("role", "parent")
		c.Request, _ = http.NewRequest("GET", "/support/pei-goals", nil)

		h.ListPeiGoals(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ListPeiGoals_ParentNotGuardian_Forbidden", func(t *testing.T) {
		mockUsers.On("IsGuardian", mock.Anything, "parent-1", "student-2").Return(false, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "parent-1")
		c.Set("role", "parent")
		c.Request, _ = http.NewRequest("GET", "/support/pei-goals?student_id=student-2", nil)

		h.ListPeiGoals(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
		mockUsers.AssertExpectations(t)
	})
}
