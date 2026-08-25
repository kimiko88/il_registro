package attendance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockServiceForFixes is a mock implementing Service interface for attendance fix tests
type MockServiceForFixes struct {
	mock.Mock
}

func (m *MockServiceForFixes) MarkAttendance(ctx context.Context, teacherID, schoolID string, req CreateAttendanceRequest) error {
	args := m.Called(ctx, teacherID, schoolID, req)
	return args.Error(0)
}

func (m *MockServiceForFixes) MarkBulk(ctx context.Context, teacherID, schoolID string, req BulkAttendanceRequest) error {
	args := m.Called(ctx, teacherID, schoolID, req)
	return args.Error(0)
}

func (m *MockServiceForFixes) UpdateAttendance(ctx context.Context, teacherID, schoolID, id string, req UpdateAttendanceRequest) error {
	args := m.Called(ctx, teacherID, schoolID, id, req)
	return args.Error(0)
}

func (m *MockServiceForFixes) DeleteClassAttendanceHour(ctx context.Context, actorID, actorRole, schoolID, classID, dateStr string, hour int) error {
	args := m.Called(ctx, actorID, actorRole, schoolID, classID, dateStr, hour)
	return args.Error(0)
}

func (m *MockServiceForFixes) GetClassAttendance(ctx context.Context, actorID, actorRole, schoolID, classID string, date string) (*ClassDailyAttendance, error) {
	args := m.Called(ctx, actorID, actorRole, schoolID, classID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ClassDailyAttendance), args.Error(1)
}

func (m *MockServiceForFixes) GetStudentAttendance(ctx context.Context, actorID, actorRole, schoolID, studentID string, from, to time.Time) ([]AttendanceResponse, error) {
	args := m.Called(ctx, actorID, actorRole, schoolID, studentID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]AttendanceResponse), args.Error(1)
}

func (m *MockServiceForFixes) RequestJustification(ctx context.Context, parentID string, req JustificationRequest) error {
	args := m.Called(ctx, parentID, req)
	return args.Error(0)
}

func (m *MockServiceForFixes) ProcessJustification(ctx context.Context, teacherID, actorRole, justificationID string, approve bool) error {
	args := m.Called(ctx, teacherID, actorRole, justificationID, approve)
	return args.Error(0)
}

func (m *MockServiceForFixes) GetPendingJustifications(ctx context.Context, actorID, actorRole, classID, schoolID string) ([]JustificationResponse, error) {
	args := m.Called(ctx, actorID, actorRole, classID, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]JustificationResponse), args.Error(1)
}

func (m *MockServiceForFixes) DeleteJustification(ctx context.Context, actorID string, justificationID string) error {
	args := m.Called(ctx, actorID, justificationID)
	return args.Error(0)
}

func (m *MockServiceForFixes) GetStudentSummary(ctx context.Context, actorID, actorRole, schoolID, studentID string) (*SummaryResponse, error) {
	args := m.Called(ctx, actorID, actorRole, schoolID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SummaryResponse), args.Error(1)
}

func (m *MockServiceForFixes) GetSchoolAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AnalyticsResponse), args.Error(1)
}

func (m *MockServiceForFixes) GetChildAttendance(ctx context.Context, parentID, studentID string, from, to time.Time) ([]AttendanceResponse, error) {
	args := m.Called(ctx, parentID, studentID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]AttendanceResponse), args.Error(1)
}

func (m *MockServiceForFixes) GetChildSummary(ctx context.Context, parentID, studentID, schoolID string) (*SummaryResponse, error) {
	args := m.Called(ctx, parentID, studentID, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SummaryResponse), args.Error(1)
}

func (m *MockServiceForFixes) GetChildAttendanceTrends(ctx context.Context, parentID, studentID string) (*TrendsResponse, error) {
	args := m.Called(ctx, parentID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TrendsResponse), args.Error(1)
}

func (m *MockServiceForFixes) GetMonthlyBreakdown(ctx context.Context, actorID, actorRole, schoolID, studentID, schoolYear string) (*MonthlyBreakdownResponse, error) {
	args := m.Called(ctx, actorID, actorRole, schoolID, studentID, schoolYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MonthlyBreakdownResponse), args.Error(1)
}

func (m *MockServiceForFixes) GetChildMonthlyBreakdown(ctx context.Context, parentID, studentID, schoolYear string) (*MonthlyBreakdownResponse, error) {
	args := m.Called(ctx, parentID, studentID, schoolYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MonthlyBreakdownResponse), args.Error(1)
}

func (m *MockServiceForFixes) GetChildUnjustified(ctx context.Context, parentID, studentID string) ([]Attendance, error) {
	args := m.Called(ctx, parentID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Attendance), args.Error(1)
}

func (m *MockServiceForFixes) JustifyChildAbsence(ctx context.Context, parentID, studentID, attendanceID string, req JustifyAbsenceRequest) error {
	args := m.Called(ctx, parentID, studentID, attendanceID, req)
	return args.Error(0)
}

func (m *MockServiceForFixes) GetChildAttendanceStats(ctx context.Context, parentID, studentID string) (*AttendanceStats, error) {
	args := m.Called(ctx, parentID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AttendanceStats), args.Error(1)
}

func TestDeleteClassAttendanceHour_ForbiddenReturns403(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockServiceForFixes)
	h := NewHandler(mockSvc)

	mockSvc.On("DeleteClassAttendanceHour", mock.Anything, "t1", "teacher", "sch1", "c1", "2026-08-22", 1).
		Return(ErrForbidden)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", "t1")
	c.Set("role", "teacher")
	c.Set("school_id", "sch1")
	c.Params = gin.Params{{Key: "id", Value: "c1"}, {Key: "hour", Value: "1"}}
	c.Request, _ = http.NewRequest("DELETE", "/class/c1/hour/1?date=2026-08-22", nil)

	h.DeleteClassAttendanceHour(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestProcessJustification_AlreadyProcessedReturns409(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockServiceForFixes)
	h := NewHandler(mockSvc)

	mockSvc.On("ProcessJustification", mock.Anything, "t1", "teacher", "j123", true).
		Return(ErrAlreadyProcessed)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", "t1")
	c.Set("role", "teacher")
	c.Params = gin.Params{{Key: "id", Value: "j123"}}
	body := `{"approve": true}`
	c.Request, _ = http.NewRequest("POST", "/justification/j123/process", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.ProcessJustification(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestExportAttendance_IncludesDateHourTimes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockServiceForFixes)
	h := NewHandler(mockSvc)

	mockRes := &ClassDailyAttendance{
		ClassID: "c1",
		Date:    "2026-08-22",
		Records: []AttendanceResponse{
			{
				ID:          "a1",
				StudentID:   "s1",
				Date:        "2026-08-22",
				Hour:        2,
				Status:      StatusPresent,
				IsJustified: false,
				EntryTime:   "08:00",
				ExitTime:    "13:00",
				Notes:       "On time",
			},
		},
	}

	mockSvc.On("GetClassAttendance", mock.Anything, "t1", "teacher", "sch1", "c1", "2026-08-22").
		Return(mockRes, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", "t1")
	c.Set("role", "teacher")
	c.Set("school_id", "sch1")
	c.Request, _ = http.NewRequest("GET", "/export?class_id=c1&date=2026-08-22", nil)

	h.ExportAttendance(c)

	assert.Equal(t, http.StatusOK, w.Code)
	csvBody := w.Body.String()
	assert.Contains(t, csvBody, "StudentID,Date,Hour,Status,IsJustified,EntryTime,ExitTime,Notes")
	assert.Contains(t, csvBody, "s1,2026-08-22,2,Present,false,08:00,13:00,On time")
}

func TestCountWeekdays(t *testing.T) {
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC) // Tuesday
	end := time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC)  // Monday (14 days total)
	count := countWeekdays(start, end)
	assert.Equal(t, 10, count)
}

func TestAttendance_MarkBulk_HourValidation(t *testing.T) {
	s := &service{}

	ctx := context.Background()
	req0 := BulkAttendanceRequest{
		ClassID: "c1",
		Date:    "2026-03-15",
		Hour:    0, // Invalid hour 0
		Statuses: []StudentStatusRequest{
			{StudentID: "s1", Status: StatusPresent},
		},
	}
	err0 := s.MarkBulk(ctx, "t1", "sch1", req0)
	assert.Error(t, err0)
	assert.Contains(t, err0.Error(), "ora lezione non valida")

	req13 := BulkAttendanceRequest{
		ClassID: "c1",
		Date:    "2026-03-15",
		Hour:    13, // Invalid hour > 12
		Statuses: []StudentStatusRequest{
			{StudentID: "s1", Status: StatusPresent},
		},
	}
	err13 := s.MarkBulk(ctx, "t1", "sch1", req13)
	assert.Error(t, err13)
	assert.Contains(t, err13.Error(), "ora lezione non valida")
}

func TestCountWeekdays_MidnightNormalization(t *testing.T) {
	// Single day test: 2026-09-01 (Tuesday)
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	endSameDayEndOfDay := time.Date(2026, 9, 1, 23, 59, 59, 999999999, time.UTC)
	endSameDayMidnight := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, 1, countWeekdays(start, endSameDayEndOfDay))
	assert.Equal(t, 1, countWeekdays(start, endSameDayMidnight))
}

func TestAttendance_ProcessJustification_RoleEnforcement(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	jid := "just-1"
	j := &Justification{ID: jid, StudentID: "stud-1", Status: JustificationPending}
	classID := "cls-1"
	mockRepo.On("FindJustificationByID", jid).Return(j, nil)
	mockUserRepo.On("GetByID", mock.Anything, "stud-1").Return(&users.User{ID: "stud-1", ClassID: &classID}, nil)

	// Parent role must be forbidden
	mockUserRepo.On("GetByID", mock.Anything, "parent-1").Return(&users.User{ID: "parent-1", Role: "parent"}, nil).Once()
	errParent := svc.ProcessJustification(context.Background(), "parent-1", "parent", jid, true)
	assert.Error(t, errParent)
	assert.Contains(t, errParent.Error(), "non è autorizzato")

	// Student role must be forbidden
	mockUserRepo.On("GetByID", mock.Anything, "student-1").Return(&users.User{ID: "student-1", Role: "student"}, nil).Once()
	errStudent := svc.ProcessJustification(context.Background(), "student-1", "student", jid, true)
	assert.Error(t, errStudent)
	assert.Contains(t, errStudent.Error(), "non è autorizzato")
}

func TestAttendance_Validator_EntryChecks(t *testing.T) {
	v := NewValidator()

	// Canonical valid statuses
	assert.True(t, v.IsValidStatus(StatusPresent))
	assert.True(t, v.IsValidStatus(StatusAbsent))
	assert.True(t, v.IsValidStatus(StatusLate))
	assert.True(t, v.IsValidStatus(StatusEarlyExit))
	assert.False(t, v.IsValidStatus("unknown_status"))
	assert.False(t, v.IsValidStatus("left_early"))

	// Historical entry within current/recent school year (60 days ago) must be allowed
	h := 1
	validPastEntry := &Attendance{
		Status: StatusPresent,
		Date:   time.Now().AddDate(0, 0, -60),
		Hour:   &h,
	}
	assert.NoError(t, v.ValidateEntry(validPastEntry))

	// Future entry (tomorrow) must be rejected
	tomorrowEntry := &Attendance{
		Status: StatusPresent,
		Date:   time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour).Add(12 * time.Hour),
		Hour:   &h,
	}
	assert.Error(t, v.ValidateEntry(tomorrowEntry))
}

func TestAttendance_ProcessJustification_AdminBypassesClassAssignment(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	jid := "just-admin-1"
	j := &Justification{ID: jid, StudentID: "stud-1", Status: JustificationPending}
	classID := "cls-1"
	schoolID := "school-1"

	mockRepo.On("FindJustificationByID", jid).Return(j, nil)
	mockUserRepo.On("GetByID", mock.Anything, "stud-1").Return(&users.User{ID: "stud-1", ClassID: &classID, SchoolID: &schoolID}, nil)
	mockUserRepo.On("GetByID", mock.Anything, "admin-1").Return(&users.User{ID: "admin-1", Role: "admin", SchoolID: &schoolID}, nil)
	mockRepo.On("ProcessJustificationTx", mock.Anything, j, "admin-1", true).Return(nil)

	err := svc.ProcessJustification(context.Background(), "admin-1", "admin", jid, true)
	assert.NoError(t, err)
}

func TestAttendance_GetStudentAttendance_UnassignedStudent_Forbidden(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockUserRepo := new(MockUserRepo)
	svc := NewService(mockRepo, mockUserRepo, nil, nil)

	schoolID := "school-1"
	mockUserRepo.On("GetByID", mock.Anything, "teacher-1").Return(&users.User{ID: "teacher-1", Role: "teacher", SchoolID: &schoolID}, nil)
	// Student with nil ClassID
	mockUserRepo.On("GetByID", mock.Anything, "stud-unassigned").Return(&users.User{ID: "stud-unassigned", Role: "student", SchoolID: &schoolID, ClassID: nil}, nil)

	_, err := svc.GetStudentAttendance(context.Background(), "teacher-1", "teacher", schoolID, "stud-unassigned", time.Now().Add(-24*time.Hour), time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non è assegnato ad alcuna classe")
}
