package lessons

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLessonsService struct {
	mock.Mock
}

func (m *mockLessonsService) CreateLesson(teacherID string, req CreateLessonRequest) (*LessonResponse, error) {
	args := m.Called(teacherID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LessonResponse), args.Error(1)
}

func (m *mockLessonsService) GetLessonByID(id string) (*LessonResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LessonResponse), args.Error(1)
}

func (m *mockLessonsService) UpdateLesson(teacherID, role, id string, req UpdateLessonRequest) (*LessonResponse, error) {
	args := m.Called(teacherID, role, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LessonResponse), args.Error(1)
}

func (m *mockLessonsService) DeleteLesson(teacherID, role, id string) error {
	return m.Called(teacherID, role, id).Error(0)
}

func (m *mockLessonsService) GetLessons(classID, subjectID string, date string) ([]LessonResponse, error) {
	args := m.Called(classID, subjectID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]LessonResponse), args.Error(1)
}

func (m *mockLessonsService) GetLessonsByGroup(groupID string, date string) ([]LessonResponse, error) {
	args := m.Called(groupID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]LessonResponse), args.Error(1)
}

func (m *mockLessonsService) GetTeacherDiary(teacherID string, fromDate, toDate string) ([]LessonResponse, error) {
	args := m.Called(teacherID, fromDate, toDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]LessonResponse), args.Error(1)
}

func (m *mockLessonsService) CreateHomework(teacherID string, req CreateHomeworkRequest) (*HomeworkResponse, error) {
	args := m.Called(teacherID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*HomeworkResponse), args.Error(1)
}

func (m *mockLessonsService) UpdateHomework(teacherID, role, id string, req UpdateHomeworkRequest) (*HomeworkResponse, error) {
	args := m.Called(teacherID, role, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*HomeworkResponse), args.Error(1)
}

func (m *mockLessonsService) DeleteHomework(teacherID, role, id string) error {
	return m.Called(teacherID, role, id).Error(0)
}

func (m *mockLessonsService) GetHomeworks(classID string) ([]HomeworkResponse, error) {
	args := m.Called(classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]HomeworkResponse), args.Error(1)
}

func TestGetActivityHours_RangeLimits(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockLessonsService)
	h := NewHandler(mockSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "u-1")
		c.Set("role", "teacher")
	})
	h.RegisterRoutes(r.Group("/api/v1"))

	// 1. To before From -> 400
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodGet, "/api/v1/lessons/class/c-1/activity-hours?from=2026-06-01&to=2026-01-01", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusBadRequest, w1.Code)
	assert.Contains(t, w1.Body.String(), "non può essere precedente")

	// 2. Over 1-year range -> 400
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/lessons/class/c-1/activity-hours?from=2024-01-01&to=2026-06-01", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
	assert.Contains(t, w2.Body.String(), "range di date troppo ampio")
}

func TestGetLessons_ErrorSanitization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockLessonsService)
	h := NewHandler(mockSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "u-1")
		c.Set("role", "teacher")
	})
	h.RegisterRoutes(r.Group("/api/v1"))

	mockSvc.On("GetLessons", "c-1", "", "").Return(nil, errors.New("sensitive db query syntax error at line 42")).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/lessons/class/c-1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NotContains(t, w.Body.String(), "sensitive db query")
	assert.Contains(t, w.Body.String(), "internal server error")
}

func TestGetActivityHours_JointPCTOOrientamentoCounting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockLessonsService)
	h := NewHandler(mockSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "u-1")
		c.Set("role", "teacher")
	})
	h.RegisterRoutes(r.Group("/api/v1"))

	// 10h pure PCTO, 5h pure Orientamento, 20h joint pcto_orientamento
	mockLessons := []LessonResponse{
		{ID: "l1", ActivityType: "pcto", Duration: 10},
		{ID: "l2", ActivityType: "orientamento", Duration: 5},
		{ID: "l3", ActivityType: "pcto_orientamento", Duration: 20},
	}
	mockSvc.On("GetLessons", "c-100", "", "").Return(mockLessons, nil).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/lessons/class/c-100/activity-hours", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var counters map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &counters)
	assert.NoError(t, err)

	// Joint hours should be clamped to 15h
	assert.Equal(t, 20, counters["pcto_orientamento"])
	assert.Equal(t, 15, counters["pcto_orientamento_effective"])
	// pcto = 10 + 15 = 25
	assert.Equal(t, 25, counters["pcto"])
	// orientamento = 5 + 15 = 20
	assert.Equal(t, 20, counters["orientamento"])
}
