package attendance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockServiceForHardenTest struct {
	mock.Mock
}

func (m *mockServiceForHardenTest) MarkAttendance(ctx context.Context, teacherID, schoolID string, req CreateAttendanceRequest) error {
	return m.Called(ctx, teacherID, schoolID, req).Error(0)
}
func (m *mockServiceForHardenTest) MarkBulk(ctx context.Context, teacherID, schoolID string, req BulkAttendanceRequest) error {
	return m.Called(ctx, teacherID, schoolID, req).Error(0)
}
func (m *mockServiceForHardenTest) UpdateAttendance(ctx context.Context, teacherID, schoolID, id string, req UpdateAttendanceRequest) error {
	return m.Called(ctx, teacherID, schoolID, id, req).Error(0)
}
func (m *mockServiceForHardenTest) DeleteClassAttendanceHour(ctx context.Context, actorID, actorRole, schoolID, classID, dateStr string, hour int) error {
	return m.Called(ctx, actorID, actorRole, schoolID, classID, dateStr, hour).Error(0)
}
func (m *mockServiceForHardenTest) GetClassAttendance(ctx context.Context, actorID, actorRole, schoolID, classID string, date string) (*ClassDailyAttendance, error) {
	args := m.Called(ctx, actorID, actorRole, schoolID, classID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ClassDailyAttendance), args.Error(1)
}
func (m *mockServiceForHardenTest) GetStudentAttendance(ctx context.Context, actorID, actorRole, schoolID, studentID string, from, to time.Time) ([]AttendanceResponse, error) {
	args := m.Called(ctx, actorID, actorRole, schoolID, studentID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]AttendanceResponse), args.Error(1)
}
func (m *mockServiceForHardenTest) RequestJustification(ctx context.Context, parentID string, req JustificationRequest) error {
	return m.Called(ctx, parentID, req).Error(0)
}
func (m *mockServiceForHardenTest) ProcessJustification(ctx context.Context, teacherID, actorRole, justificationID string, approve bool) error {
	return m.Called(ctx, teacherID, actorRole, justificationID, approve).Error(0)
}
func (m *mockServiceForHardenTest) GetPendingJustifications(ctx context.Context, actorID, actorRole, classID, schoolID string) ([]JustificationResponse, error) {
	args := m.Called(ctx, actorID, actorRole, classID, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]JustificationResponse), args.Error(1)
}
func (m *mockServiceForHardenTest) DeleteJustification(ctx context.Context, actorID string, justificationID string) error {
	return m.Called(ctx, actorID, justificationID).Error(0)
}
func (m *mockServiceForHardenTest) GetStudentSummary(ctx context.Context, actorID, actorRole, schoolID, studentID string) (*SummaryResponse, error) {
	args := m.Called(ctx, actorID, actorRole, schoolID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SummaryResponse), args.Error(1)
}
func (m *mockServiceForHardenTest) GetSchoolAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AnalyticsResponse), args.Error(1)
}
func (m *mockServiceForHardenTest) GetChildAttendance(ctx context.Context, parentID, studentID string, from, to time.Time) ([]AttendanceResponse, error) {
	args := m.Called(ctx, parentID, studentID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]AttendanceResponse), args.Error(1)
}
func (m *mockServiceForHardenTest) GetChildSummary(ctx context.Context, parentID, studentID, schoolID string) (*SummaryResponse, error) {
	args := m.Called(ctx, parentID, studentID, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SummaryResponse), args.Error(1)
}
func (m *mockServiceForHardenTest) GetChildAttendanceTrends(ctx context.Context, parentID, studentID string) (*TrendsResponse, error) {
	args := m.Called(ctx, parentID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TrendsResponse), args.Error(1)
}
func (m *mockServiceForHardenTest) GetMonthlyBreakdown(ctx context.Context, actorID, actorRole, schoolID, studentID, schoolYear string) (*MonthlyBreakdownResponse, error) {
	args := m.Called(ctx, studentID, schoolYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MonthlyBreakdownResponse), args.Error(1)
}
func (m *mockServiceForHardenTest) GetChildMonthlyBreakdown(ctx context.Context, parentID, studentID, schoolYear string) (*MonthlyBreakdownResponse, error) {
	args := m.Called(ctx, parentID, studentID, schoolYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MonthlyBreakdownResponse), args.Error(1)
}
func (m *mockServiceForHardenTest) GetChildUnjustified(ctx context.Context, parentID, studentID string) ([]Attendance, error) {
	args := m.Called(ctx, parentID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Attendance), args.Error(1)
}
func (m *mockServiceForHardenTest) JustifyChildAbsence(ctx context.Context, parentID, studentID, attendanceID string, req JustifyAbsenceRequest) error {
	return m.Called(ctx, parentID, studentID, attendanceID, req).Error(0)
}
func (m *mockServiceForHardenTest) GetChildAttendanceStats(ctx context.Context, parentID, studentID string) (*AttendanceStats, error) {
	args := m.Called(ctx, parentID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AttendanceStats), args.Error(1)
}

func TestGetAnalytics_RequiresUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForHardenTest)
	h := NewHandler(mockSvc)

	r := gin.New()
	// No user_id set in context
	r.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
	})
	h.RegisterRoutes(r.Group("/api/v1"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/attendance/analytics", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestExportAttendance_HeaderSanitization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForHardenTest)
	h := NewHandler(mockSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
	})
	h.RegisterRoutes(r.Group("/api/v1"))

	mockSvc.On("GetClassAttendance", mock.Anything, "teacher-1", "teacher", "school-1", "class\"1;malicious\r\n", "2026-03-01").
		Return(&ClassDailyAttendance{Records: []AttendanceResponse{}}, nil).Once()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/attendance/export?class_id=class%221%3Bmalicious%0D%0A&date=2026-03-01", nil)
	assert.NoError(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	disp := w.Header().Get("Content-Disposition")
	assert.NotContains(t, disp, "\r")
	assert.NotContains(t, disp, "\n")
	assert.NotContains(t, disp, ";filename=") // No extra semicolon injection
	assert.Contains(t, disp, "class_1_malicious__")
}
