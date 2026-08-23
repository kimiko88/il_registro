package attendance

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPendingJustifications_TeacherEmptyClassID_Forbidden(t *testing.T) {
	svc := &service{
		repo: &mockAuditAttendanceRepo{},
	}

	_, err := svc.GetPendingJustifications(context.Background(), "teacher-1", "teacher", "", "school-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden: classID obbligatorio per i docenti")
}

func TestDeleteClassAttendanceHour_DifferentSchool_Forbidden(t *testing.T) {
	mockRepo := new(mockAuditAttendanceRepo)
	svc := &service{repo: mockRepo}

	// Mock class not in user's school
	mockRepo.isClassInSchool = false

	err := svc.DeleteClassAttendanceHour(context.Background(), "admin-1", "admin", "school-1", "class-2", "2026-03-01", 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "la classe non appartiene alla scuola dell'utente")
}

func TestProcessJustification_EmptyTeacherID_Unauthorized(t *testing.T) {
	svc := &service{}
	err := svc.ProcessJustification(context.Background(), "", "justification-1", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized: teacherID mancante")
}

func TestParseWindowParams_ExtendedItalianSchoolYear_Allowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	// 400 days range (e.g. Sept 1 2025 to Oct 5 2026)
	req, _ := http.NewRequest("GET", "/test?from=2025-09-01&to=2026-10-05", nil)
	c.Request = req

	from, to, err := parseWindowParams(c)
	assert.NoError(t, err)
	assert.Equal(t, "2025-09-01", from.Format("2006-01-02"))
	assert.Equal(t, "2026-10-05", to.Format("2006-01-02"))
}

func TestExportAttendance_ErrorMapping_ForbiddenAndNotFound(t *testing.T) {
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

	// Test 403 response when teacher is not assigned to class
	mockSvc.On("GetClassAttendance", mock.Anything, "teacher-1", "teacher", "school-1", "class-forbidden", "2026-03-01").
		Return(nil, errors.New("forbidden: docente non assegnato alla classe")).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/attendance/export?class_id=class-forbidden&date=2026-03-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

type mockAuditAttendanceRepo struct {
	Repository
	pendingJustifications []Justification
	isClassInSchool       bool
}

func (m *mockAuditAttendanceRepo) FindPendingJustifications(classID, schoolID string) ([]Justification, error) {
	if schoolID == "" {
		return m.pendingJustifications, nil
	}
	var filtered []Justification
	for _, j := range m.pendingJustifications {
		if j.SchoolID == schoolID {
			filtered = append(filtered, j)
		}
	}
	return filtered, nil
}

func (m *mockAuditAttendanceRepo) IsTeacherAssignedToClass(_ context.Context, _, _ string) (bool, error) {
	return true, nil
}

func (m *mockAuditAttendanceRepo) IsClassInSchool(_ context.Context, _, _ string) (bool, error) {
	return m.isClassInSchool, nil
}

func TestValidator_StatusAndDateValidation(t *testing.T) {
	v := NewValidator()

	// Test Status validation
	assert.True(t, v.IsValidStatus(StatusPresent))
	assert.True(t, v.IsValidStatus(StatusAbsent))
	assert.True(t, v.IsValidStatus(StatusLate))
	assert.False(t, v.IsValidStatus("present"))
	assert.False(t, v.IsValidStatus("invalid_status_xyz"))

	// Test Invalid Status Entry
	invalidStatusEntry := &Attendance{
		Status: "invalid_status_xyz",
		Date:   time.Now(),
	}
	assert.Error(t, v.ValidateEntry(invalidStatusEntry))
	assert.Contains(t, v.ValidateEntry(invalidStatusEntry).Error(), "invalid attendance status")

	// Test Future Date Entry
	futureEntry := &Attendance{
		Status: StatusPresent,
		Date:   time.Now().Add(48 * time.Hour),
	}
	assert.Error(t, v.ValidateEntry(futureEntry))
	assert.Contains(t, v.ValidateEntry(futureEntry).Error(), "cannot mark attendance in future")

	// Test Obsolete Limit (over 2 years)
	oldEntry := &Attendance{
		Status: StatusPresent,
		Date:   time.Now().AddDate(-3, 0, 0),
	}
	assert.Error(t, v.ValidateEntry(oldEntry))
	assert.Contains(t, v.ValidateEntry(oldEntry).Error(), "cannot mark attendance for dates older than 2 years")
}
