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

func (m *MockAttRepo) Create(a *Attendance) error              { return m.Called(a).Error(0) }
func (m *MockAttRepo) BatchCreate(atts []*Attendance) error    { return m.Called(atts).Error(0) }
func (m *MockAttRepo) FindByID(id string) (*Attendance, error) { return nil, nil }
func (m *MockAttRepo) Update(a *Attendance) error              { return nil }
func (m *MockAttRepo) FindByClassAndDate(classID string, date time.Time) ([]Attendance, error) {
	return nil, nil
}
func (m *MockAttRepo) FindByStudent(studentID string, start, end time.Time) ([]Attendance, error) {
	return nil, nil
}
func (m *MockAttRepo) CreateJustification(j *Justification) error              { return nil }
func (m *MockAttRepo) FindJustificationByID(id string) (*Justification, error) { return nil, nil }
func (m *MockAttRepo) UpdateJustification(j *Justification) error              { return nil }
func (m *MockAttRepo) FindPendingJustifications(classID string) ([]Justification, error) {
	return nil, nil
}
func (m *MockAttRepo) GetStats(studentID string) (*SummaryResponse, error) { return nil, nil }

func TestRegression_FutureAttendance(t *testing.T) {
	repo := &MockAttRepo{}
	svc := NewService(repo, nil)
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
	err := svc.MarkAttendance(ctx, "teacher-1", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "future")
}

func TestRegression_BulkMixedValidity(t *testing.T) {
	repo := &MockAttRepo{}
	svc := NewService(repo, nil)
	ctx := context.Background()

	future := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	req := BulkAttendanceRequest{
		ClassID: "class-A",
		Date:    future,
		Statuses: []CreateAttendanceRequest{
			{StudentID: "s1", Status: StatusPresent},
		},
	}

	err := svc.MarkBulk(ctx, "teacher-1", req)
	assert.Error(t, err)
}
