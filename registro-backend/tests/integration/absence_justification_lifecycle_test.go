package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/attendance"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Absence_Justification_Lifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockAbsenceRepoForLifecycle{}
	mockUsers := &mockUserRepoForComms{}
	svc := attendance.NewService(mockRepo, mockUsers, nil, nil)
	handler := attendance.NewHandler(svc)

	r := gin.New()
	r.GET("/attendance/justifications/pending", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.GetPendingJustifications(c)
	})

	r.POST("/attendance/justifications/process", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.ProcessJustification(c)
	})

	// 1. Get pending justifications requiring class_id for teacher
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/attendance/justifications/pending", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusForbidden, w1.Code)

	// 2. Get pending justifications with class_id
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/attendance/justifications/pending?class_id=class-1", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

type mockAbsenceRepoForLifecycle struct{}

func (m *mockAbsenceRepoForLifecycle) Create(att *attendance.Attendance) error         { return nil }
func (m *mockAbsenceRepoForLifecycle) BatchCreate(atts []*attendance.Attendance) error { return nil }
func (m *mockAbsenceRepoForLifecycle) Update(att *attendance.Attendance) error         { return nil }
func (m *mockAbsenceRepoForLifecycle) DeleteByClassDateHour(schoolID, classID string, date time.Time, hour int) error {
	return nil
}
func (m *mockAbsenceRepoForLifecycle) FindByID(id string) (*attendance.Attendance, error) {
	return &attendance.Attendance{ID: id}, nil
}
func (m *mockAbsenceRepoForLifecycle) FindByClassAndDate(classID string, date time.Time) ([]attendance.Attendance, error) {
	return []attendance.Attendance{}, nil
}
func (m *mockAbsenceRepoForLifecycle) FindByStudent(studentID string, startDate, endDate time.Time) ([]attendance.Attendance, error) {
	return []attendance.Attendance{}, nil
}
func (m *mockAbsenceRepoForLifecycle) GetStats(studentID string) (*attendance.SummaryResponse, error) {
	return &attendance.SummaryResponse{}, nil
}
func (m *mockAbsenceRepoForLifecycle) GetStatsBatch(ctx context.Context, studentIDs []string) (map[string]*attendance.SummaryResponse, error) {
	return map[string]*attendance.SummaryResponse{}, nil
}
func (m *mockAbsenceRepoForLifecycle) CountDistinctDays(studentID string) (int, error) {
	return 0, nil
}
func (m *mockAbsenceRepoForLifecycle) GetAnalytics(ctx context.Context, schoolID string) (*attendance.AnalyticsResponse, error) {
	return &attendance.AnalyticsResponse{}, nil
}
func (m *mockAbsenceRepoForLifecycle) CreateJustification(j *attendance.Justification) error {
	return nil
}
func (m *mockAbsenceRepoForLifecycle) UpdateJustification(j *attendance.Justification) error {
	return nil
}
func (m *mockAbsenceRepoForLifecycle) ProcessJustificationTx(ctx context.Context, j *attendance.Justification, teacherID string, approve bool) error {
	return nil
}
func (m *mockAbsenceRepoForLifecycle) FindJustificationByID(id string) (*attendance.Justification, error) {
	return &attendance.Justification{ID: id}, nil
}
func (m *mockAbsenceRepoForLifecycle) FindPendingJustifications(classID, schoolID string) ([]attendance.Justification, error) {
	return []attendance.Justification{}, nil
}
func (m *mockAbsenceRepoForLifecycle) DeleteJustification(id string) error        { return nil }
func (m *mockAbsenceRepoForLifecycle) DeletePendingJustification(id string) error { return nil }
func (m *mockAbsenceRepoForLifecycle) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	return true, nil
}
func (m *mockAbsenceRepoForLifecycle) IsTeacherSubstitute(ctx context.Context, teacherID, classID string, date time.Time, hour int) (bool, error) {
	return false, nil
}
func (m *mockAbsenceRepoForLifecycle) JustifyAbsenceByParent(attendanceID string, parentID string, reason string, documentURL string) error {
	return nil
}
func (m *mockAbsenceRepoForLifecycle) IsClassInSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	return true, nil
}
func (m *mockAbsenceRepoForLifecycle) IsStudentInClass(ctx context.Context, studentID, classID string) (bool, error) {
	return true, nil
}
func (m *mockAbsenceRepoForLifecycle) AreStudentsInClass(ctx context.Context, studentIDs []string, classID string) (map[string]bool, error) {
	res := make(map[string]bool)
	for _, id := range studentIDs {
		res[id] = true
	}
	return res, nil
}
func (m *mockAbsenceRepoForLifecycle) FindUnjustifiedByStudent(studentID string) ([]attendance.Attendance, error) {
	return []attendance.Attendance{}, nil
}
func (m *mockAbsenceRepoForLifecycle) GetMonthlyBreakdown(ctx context.Context, studentID string, schoolID string) ([]attendance.MonthlyBreakdownRow, error) {
	return []attendance.MonthlyBreakdownRow{}, nil
}
func (m *mockAbsenceRepoForLifecycle) GetStudentAttendanceStats(studentID string) (*attendance.AttendanceStats, error) {
	return &attendance.AttendanceStats{}, nil
}
func (m *mockAbsenceRepoForLifecycle) HasOverlappingJustification(ctx context.Context, studentID string, startDate, endDate time.Time) (bool, error) {
	return false, nil
}
