package timetables

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrString(s string) *string {
	return &s
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByClass(ctx context.Context, classID string) ([]ClassSchedule, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ClassSchedule), args.Error(1)
}

func (m *MockRepository) GetStudentClassID(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) GetParentStudentClassID(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) GetTeacherSchedule(ctx context.Context, userID string) ([]ClassSchedule, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ClassSchedule), args.Error(1)
}

func (m *MockRepository) GetClassSchoolID(ctx context.Context, classID string) (string, error) {
	args := m.Called(ctx, classID)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, classID string, entries []ScheduleEntry) error {
	args := m.Called(ctx, classID, entries)
	return args.Error(0)
}

func (m *MockRepository) UpdateTeacher(ctx context.Context, teacherID string, entries []TeacherScheduleEntry) error {
	args := m.Called(ctx, teacherID, entries)
	return args.Error(0)
}

func TestService_GetByClass(t *testing.T) {
	t.Run("Missing context", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		_, err := svc.GetByClass(context.Background(), "", "", "school-1", "class-1")
		assert.Error(t, err)
	})

	t.Run("Missing class_id", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		_, err := svc.GetByClass(context.Background(), "user-1", "admin", "school-1", "")
		assert.Error(t, err)
	})

	t.Run("Student viewing own class -> Success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		mockRepo.On("GetStudentClassID", mock.Anything, "st-1").Return("class-1", nil)
		mockRepo.On("GetByClass", mock.Anything, "class-1").Return([]ClassSchedule{
			{ClassID: "class-1", DayOfWeek: 1, HourIndex: 1, SubjectName: "Matematica"},
		}, nil)

		res, err := svc.GetByClass(context.Background(), "st-1", "student", "school-1", "class-1")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "Matematica", res[0].SubjectName)
	})

	t.Run("Student viewing foreign class -> Forbidden", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		mockRepo.On("GetStudentClassID", mock.Anything, "st-1").Return("class-1", nil)

		_, err := svc.GetByClass(context.Background(), "st-1", "student", "school-1", "class-other")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "forbidden")
	})

	t.Run("Parent viewing own child's class -> Success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		mockRepo.On("GetParentStudentClassID", mock.Anything, "parent-1").Return("class-1", nil)
		mockRepo.On("GetByClass", mock.Anything, "class-1").Return([]ClassSchedule{
			{ClassID: "class-1", DayOfWeek: 1, HourIndex: 2, SubjectName: "Italiano"},
		}, nil)

		res, err := svc.GetByClass(context.Background(), "parent-1", "parent", "school-1", "class-1")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestService_GetMySchedule(t *testing.T) {
	t.Run("Teacher schedule", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		mockRepo.On("GetTeacherSchedule", mock.Anything, "t-1").Return([]ClassSchedule{
			{TeacherID: ptrString("t-1"), DayOfWeek: 1, HourIndex: 1, ClassName: "1A"},
		}, nil)

		res, err := svc.GetMySchedule(context.Background(), "t-1", "teacher")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "1A", res[0].ClassName)
	})

	t.Run("Unauthorized actor", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		_, err := svc.GetMySchedule(context.Background(), "", "teacher")
		assert.Error(t, err)
	})
}

func TestService_GetByTeacher(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		mockRepo.On("GetTeacherSchedule", mock.Anything, "t-100").Return([]ClassSchedule{
			{TeacherID: ptrString("t-100"), SubjectName: "Fisica"},
		}, nil)

		res, err := svc.GetByTeacher(context.Background(), "admin-1", "admin", "school-1", "t-100")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("Missing teacher ID", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		_, err := svc.GetByTeacher(context.Background(), "admin-1", "admin", "school-1", "")
		assert.Error(t, err)
	})
}

func TestService_UpdateClass(t *testing.T) {
	t.Run("Secretary updates class timetable -> Allowed", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		entries := []ScheduleEntry{
			{DayOfWeek: 1, HourIndex: 1, SubjectID: "sub-1", TeacherID: ptrString("t-1")},
		}
		mockRepo.On("GetClassSchoolID", mock.Anything, "class-1").Return("school-1", nil)
		mockRepo.On("Update", mock.Anything, "class-1", entries).Return(nil)

		err := svc.Update(context.Background(), "sec-1", "secretary", "school-1", "class-1", entries)
		assert.NoError(t, err)
	})

	t.Run("Student updates class timetable -> Forbidden", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		err := svc.Update(context.Background(), "st-1", "student", "school-1", "class-1", []ScheduleEntry{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "forbidden")
	})
}

func TestService_UpdateTeacher(t *testing.T) {
	t.Run("Admin updates teacher timetable -> Allowed", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		entries := []TeacherScheduleEntry{
			{DayOfWeek: 1, HourIndex: 1, ClassID: "c-1", SubjectID: "sub-1"},
		}
		mockRepo.On("UpdateTeacher", mock.Anything, "t-1", entries).Return(nil)

		err := svc.UpdateTeacher(context.Background(), "admin-1", "admin", "school-1", "t-1", entries)
		assert.NoError(t, err)
	})

	t.Run("Teacher updates teacher timetable -> Forbidden", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)
		err := svc.UpdateTeacher(context.Background(), "t-1", "teacher", "school-1", "t-1", []TeacherScheduleEntry{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "forbidden")
	})
}
