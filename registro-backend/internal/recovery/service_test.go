package recovery

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRecoveryRepo struct {
	mock.Mock
}

func (m *MockRecoveryRepo) CreateCourse(ctx context.Context, c *RecoveryCourse, sessions []RecoveryCourseSession, studentIDs []string) error {
	args := m.Called(ctx, c, sessions, studentIDs)
	c.ID = "course-123"
	return args.Error(0)
}

func (m *MockRecoveryRepo) GetCourseByID(ctx context.Context, id string) (*RecoveryCourse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RecoveryCourse), args.Error(1)
}

func (m *MockRecoveryRepo) ListCourses(ctx context.Context, schoolID, academicYear, teacherID string) ([]RecoveryCourse, error) {
	args := m.Called(ctx, schoolID, academicYear, teacherID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]RecoveryCourse), args.Error(1)
}

func (m *MockRecoveryRepo) UpdateCourseStatus(ctx context.Context, id, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockRecoveryRepo) UpdateStudentAttendance(ctx context.Context, courseID, studentID string, hours float64, notes string) error {
	args := m.Called(ctx, courseID, studentID, hours, notes)
	return args.Error(0)
}

func (m *MockRecoveryRepo) RecordRecoveryTest(ctx context.Context, test *RecoveryTest) error {
	args := m.Called(ctx, test)
	test.ID = "test-123"
	return args.Error(0)
}

func (m *MockRecoveryRepo) ListRecoveryTests(ctx context.Context, schoolID, classID, studentID string) ([]RecoveryTest, error) {
	args := m.Called(ctx, schoolID, classID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]RecoveryTest), args.Error(1)
}

func TestRecoveryService_CreateCourse(t *testing.T) {
	repo := new(MockRecoveryRepo)
	svc := NewService(repo)

	req := CreateCourseRequest{
		SubjectID:    "sub-1",
		TeacherID:    "t-1",
		Title:        "Recupero Matematica Estivo",
		AcademicYear: "2024/2025",
		Period:       "summer",
		TotalHours:   12,
		Sessions: []RecoveryCourseSession{
			{SessionDate: "2025-07-01", StartTime: "09:00", EndTime: "11:00"},
		},
		StudentIDs: []string{"s-1", "s-2"},
	}

	repo.On("CreateCourse", mock.Anything, mock.Anything, req.Sessions, req.StudentIDs).Return(nil).Once()
	repo.On("GetCourseByID", mock.Anything, "course-123").Return(&RecoveryCourse{
		ID:    "course-123",
		Title: "Recupero Matematica Estivo",
	}, nil).Once()

	c, err := svc.CreateCourse(context.Background(), "school-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "course-123", c.ID)
	repo.AssertExpectations(t)
}

func TestRecoveryService_RecordTestOutcome(t *testing.T) {
	repo := new(MockRecoveryRepo)
	svc := NewService(repo)

	// Test case 1: Passed test (grade 6.5) -> recuperato / Ammesso
	defID := "def-1"
	reqPassed := RecordTestOutcomeRequest{
		DeficiencyID:  &defID,
		StudentID:     "s-1",
		SubjectID:     "sub-1",
		ClassID:       "c-1",
		TestDate:      "2025-09-02",
		TestType:      "written",
		Grade:         6.5,
		VerbaleNumber: "VERB-01/2025",
	}

	repo.On("RecordRecoveryTest", mock.Anything, mock.MatchedBy(func(test *RecoveryTest) bool {
		return test.Grade == 6.5 && test.Outcome == "recuperato" && test.FinalDeliberation == "Ammesso"
	})).Return(nil).Once()

	testRes, err := svc.RecordTestOutcome(context.Background(), "school-1", "teacher-1", reqPassed)
	assert.NoError(t, err)
	assert.Equal(t, "recuperato", testRes.Outcome)
	assert.Equal(t, "Ammesso", testRes.FinalDeliberation)

	// Test case 2: Failed test (grade 4.5) -> non_recuperato / Non Ammesso
	reqFailed := RecordTestOutcomeRequest{
		DeficiencyID:  &defID,
		StudentID:     "s-2",
		SubjectID:     "sub-1",
		ClassID:       "c-1",
		TestDate:      "2025-09-02",
		TestType:      "oral",
		Grade:         4.5,
		VerbaleNumber: "VERB-01/2025",
	}

	repo.On("RecordRecoveryTest", mock.Anything, mock.MatchedBy(func(test *RecoveryTest) bool {
		return test.Grade == 4.5 && test.Outcome == "non_recuperato" && test.FinalDeliberation == "Non Ammesso"
	})).Return(nil).Once()

	testRes2, err := svc.RecordTestOutcome(context.Background(), "school-1", "teacher-1", reqFailed)
	assert.NoError(t, err)
	assert.Equal(t, "non_recuperato", testRes2.Outcome)
	assert.Equal(t, "Non Ammesso", testRes2.FinalDeliberation)

	// Test case 3: Invalid grade (< 1.0 or > 10.0)
	reqInvalid := RecordTestOutcomeRequest{
		StudentID: "s-1",
		SubjectID: "sub-1",
		ClassID:   "c-1",
		TestDate:  "2025-09-02",
		TestType:  "written",
		Grade:     11.0,
	}
	_, errInvalid := svc.RecordTestOutcome(context.Background(), "school-1", "teacher-1", reqInvalid)
	assert.ErrorIs(t, errInvalid, ErrInvalidGrade)
}

func TestService_RemainingMethods(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRecoveryRepo)
	svc := NewService(repo)

	// 1. GetCourse - Found
	repo.On("GetCourseByID", ctx, "c-1").Return(&RecoveryCourse{ID: "c-1", Title: "Recupero"}, nil).Once()
	c, err := svc.GetCourse(ctx, "c-1")
	assert.NoError(t, err)
	assert.Equal(t, "c-1", c.ID)

	// 2. GetCourse - Not found
	repo.On("GetCourseByID", ctx, "c-404").Return(nil, nil).Once()
	_, err = svc.GetCourse(ctx, "c-404")
	assert.ErrorIs(t, err, ErrCourseNotFound)

	// 3. ListCourses
	repo.On("ListCourses", ctx, "s-1", "2026", "t-1").Return([]RecoveryCourse{{ID: "c-1"}}, nil).Once()
	list, err := svc.ListCourses(ctx, "s-1", "2026", "t-1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// 4. UpdateStatus
	repo.On("UpdateCourseStatus", ctx, "c-1", "completed").Return(nil).Once()
	err = svc.UpdateStatus(ctx, "c-1", "completed")
	assert.NoError(t, err)

	// 5. UpdateAttendance
	repo.On("UpdateStudentAttendance", ctx, "c-1", "s-1", 5.0, "presente").Return(nil).Once()
	err = svc.UpdateAttendance(ctx, "c-1", "s-1", 5.0, "presente")
	assert.NoError(t, err)

	// 6. ListTests
	repo.On("ListRecoveryTests", ctx, "s-1", "class-1", "student-1").Return([]RecoveryTest{{ID: "t-1"}}, nil).Once()
	tests, err := svc.ListTests(ctx, "s-1", "class-1", "student-1")
	assert.NoError(t, err)
	assert.Len(t, tests, 1)

	repo.AssertExpectations(t)
}
