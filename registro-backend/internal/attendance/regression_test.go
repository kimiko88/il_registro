package attendance

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAttRepo struct {
	mock.Mock
}

var _ Repository = (*MockAttRepo)(nil)

func (m *MockAttRepo) Create(a *Attendance) error              { return m.Called(a).Error(0) }
func (m *MockAttRepo) BatchCreate(atts []*Attendance) error    { return m.Called(atts).Error(0) }
func (m *MockAttRepo) FindByID(id string) (*Attendance, error) { return nil, nil }
func (m *MockAttRepo) Update(a *Attendance) error              { return nil }
func (m *MockAttRepo) DeleteByClassDateHour(schoolID, classID string, date time.Time, hour int) error {
	return nil
}
func (m *MockAttRepo) FindByClassAndDate(classID string, date time.Time) ([]Attendance, error) {
	return nil, nil
}
func (m *MockAttRepo) FindByStudent(studentID string, start, end time.Time) ([]Attendance, error) {
	return nil, nil
}
func (m *MockAttRepo) CreateJustification(j *Justification) error              { return nil }
func (m *MockAttRepo) FindJustificationByID(id string) (*Justification, error) { return nil, nil }
func (m *MockAttRepo) UpdateJustification(j *Justification) error              { return nil }
func (m *MockAttRepo) ProcessJustificationTx(ctx context.Context, j *Justification, teacherID string, approve bool) error {
	return m.Called(ctx, j, teacherID, approve).Error(0)
}
func (m *MockAttRepo) FindPendingJustifications(classID, schoolID string) ([]Justification, error) {
	return nil, nil
}
func (m *MockAttRepo) FindPendingJustificationsForTeacher(ctx context.Context, teacherID, schoolID string) ([]Justification, error) {
	return nil, nil
}
func (m *MockAttRepo) DeletePendingJustification(id string) error { return nil }
func (m *MockAttRepo) IsStudentInClass(ctx context.Context, studentID, classID string) (bool, error) {
	return true, nil
}
func (m *MockAttRepo) AreStudentsInClass(ctx context.Context, studentIDs []string, classID string) (map[string]bool, error) {
	res := make(map[string]bool)
	for _, id := range studentIDs {
		res[id] = true
	}
	return res, nil
}
func (m *MockAttRepo) IsClassInSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	return true, nil
}
func (m *MockAttRepo) GetStats(studentID string) (*SummaryResponse, error) { return nil, nil }
func (m *MockAttRepo) GetStatsBatch(ctx context.Context, studentIDs []string) (map[string]*SummaryResponse, error) {
	return make(map[string]*SummaryResponse), nil
}
func (m *MockAttRepo) CountDistinctDays(studentID string) (int, error) { return 0, nil }
func (m *MockAttRepo) GetAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error) {
	return nil, nil
}
func (m *MockAttRepo) DeleteJustification(id string) error { return nil }
func (m *MockAttRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	return true, nil
}
func (m *MockAttRepo) GetMonthlyBreakdown(ctx context.Context, studentID, schoolYear string) ([]MonthlyBreakdownRow, error) {
	return nil, nil
}
func (m *MockAttRepo) FindUnjustifiedByStudent(studentID string) ([]Attendance, error) {
	return nil, nil
}
func (m *MockAttRepo) JustifyAbsenceByParent(attendanceID string, studentID string, reason string, notes string) error {
	return nil
}
func (m *MockAttRepo) GetStudentAttendanceStats(studentID string) (*AttendanceStats, error) {
	return &AttendanceStats{}, nil
}
func (m *MockAttRepo) IsTeacherSubstitute(ctx context.Context, teacherID, classID string, date time.Time, hour int) (bool, error) {
	return false, nil
}
func (m *MockAttRepo) HasOverlappingJustification(ctx context.Context, studentID string, startDate, endDate time.Time) (bool, error) {
	return false, nil
}

func TestRegression_FutureAttendance(t *testing.T) {
	repo := &MockAttRepo{}
	mockUserRepo := new(MockUserRepo)
	svc := NewService(repo, mockUserRepo, nil, nil)
	ctx := context.Background()

	// Scenario: Marking attendance for way in future (> 24h allowed buffer)
	future := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	req := CreateAttendanceRequest{
		StudentID: "std-1",
		ClassID:   "class-A",
		Date:      future,
		Status:    StatusPresent,
	}

	// Expect Create NOT to be called because validation should fail
	err := svc.MarkAttendance(ctx, "teacher-1", "school-1", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "future")
}

func TestRegression_BulkMixedValidity(t *testing.T) {
	repo := &MockAttRepo{}
	mockUserRepo := new(MockUserRepo)
	svc := NewService(repo, mockUserRepo, nil, nil)
	ctx := context.Background()

	future := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	req := BulkAttendanceRequest{
		ClassID: "class-A",
		Date:    future,
		Statuses: []StudentStatusRequest{
			{StudentID: "s1", Status: StatusPresent},
		},
	}

	err := svc.MarkBulk(ctx, "teacher-1", "school-1", req)
	assert.Error(t, err)
}
