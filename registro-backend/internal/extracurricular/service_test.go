package extracurricular

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateCourse(ctx context.Context, c *Course) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *MockRepository) GetCourseByID(ctx context.Context, id string) (*Course, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Course), args.Error(1)
}
func (m *MockRepository) ListCourses(ctx context.Context, schoolID, studentID string) ([]*Course, error) {
	args := m.Called(ctx, schoolID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Course), args.Error(1)
}
func (m *MockRepository) EnrollStudent(ctx context.Context, courseID, studentID string) error {
	args := m.Called(ctx, courseID, studentID)
	return args.Error(0)
}
func (m *MockRepository) ListEnrollments(ctx context.Context, courseID string) ([]*Enrollment, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Enrollment), args.Error(1)
}
func (m *MockRepository) MarkAttendance(ctx context.Context, att *AttendanceRecord) error {
	args := m.Called(ctx, att)
	return args.Error(0)
}
func (m *MockRepository) ListAttendance(ctx context.Context, courseID string, date time.Time) ([]*AttendanceRecord, error) {
	args := m.Called(ctx, courseID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*AttendanceRecord), args.Error(1)
}

func TestCreateAndEnrollCourse(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := CreateCourseRequest{
		Title:           "Laboratorio di Robottica",
		Description:     "Corso pratico Arduino e C++",
		StartDate:       "2025-11-01",
		EndDate:         "2025-12-20",
		MaxParticipants: 15,
	}

	mockRepo.On("CreateCourse", mock.Anything, mock.MatchedBy(func(c *Course) bool {
		return c.Title == "Laboratorio di Robottica"
	})).Return(nil).Once()

	c, err := svc.CreateCourse(context.Background(), "teacher", "t-1", "school-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	mockRepo.On("GetCourseByID", mock.Anything, "course-1").Return(&Course{
		ID:              "course-1",
		MaxParticipants: 15,
		EnrolledCount:   5,
	}, nil).Once()
	mockRepo.On("EnrollStudent", mock.Anything, "course-1", "student-1").Return(nil).Once()

	err = svc.EnrollStudent(context.Background(), "course-1", "student-1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
