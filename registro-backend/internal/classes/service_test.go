package classes

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, c *Class) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *MockRepository) List(ctx context.Context, schoolID string, academicYear string) ([]Class, error) {
	args := m.Called(ctx, schoolID, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Class), args.Error(1)
}
func (m *MockRepository) Get(ctx context.Context, id string) (*Class, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Class), args.Error(1)
}
func (m *MockRepository) Update(ctx context.Context, c *Class) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepository) ListByTeacher(ctx context.Context, teacherUserID string) ([]Class, error) {
	args := m.Called(ctx, teacherUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Class), args.Error(1)
}
func (m *MockRepository) AssignSubject(ctx context.Context, classID string, subjectID string, teacherID *string, hours float64) error {
	args := m.Called(ctx, classID, subjectID, teacherID, hours)
	return args.Error(0)
}
func (m *MockRepository) UnassignSubject(ctx context.Context, assignmentID string) error {
	args := m.Called(ctx, assignmentID)
	return args.Error(0)
}
func (m *MockRepository) GetClassSubjects(ctx context.Context, classID string) ([]ClassSubject, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ClassSubject), args.Error(1)
}
func (m *MockRepository) GetClassGuardians(ctx context.Context, classID string) ([]GuardianInfo, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]GuardianInfo), args.Error(1)
}
func (m *MockRepository) GetLessonTopics(ctx context.Context, classID string) ([]LessonTopic, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]LessonTopic), args.Error(1)
}
func (m *MockRepository) GetDisciplinaryNotes(ctx context.Context, classID string) ([]DisciplinaryNoteReport, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DisciplinaryNoteReport), args.Error(1)
}


func TestService_CreateClass(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateClassRequest
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful creation",
			req:  CreateClassRequest{Name: "1A", Section: "A", AcademicYear: "2024-2025"},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(c *Class) bool {
					return c.Name == "1A" && c.Section == "A"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			req:  CreateClassRequest{Name: "2B", Section: "B", AcademicYear: "2024-2025"},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*classes.Class")).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)
			service := NewService(mockRepo)

			class, err := service.CreateClass(context.Background(), "school-1", tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, class)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, class)
				assert.Equal(t, tt.req.Name, class.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_ListClasses(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(*MockRepository)
		wantErr bool
		wantLen int
	}{
		{
			name: "successful list",
			mockFn: func(m *MockRepository) {
				classes := []Class{{Name: "1A", ID: "c1"}, {Name: "2B", ID: "c2"}}
				m.On("List", mock.Anything, "school-1", "2024-2025").Return(classes, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "empty list",
			mockFn: func(m *MockRepository) {
				m.On("List", mock.Anything, "school-1", "2024-2025").Return([]Class{}, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "repository error",
			mockFn: func(m *MockRepository) {
				m.On("List", mock.Anything, "school-1", "2024-2025").Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)
			service := NewService(mockRepo)

			list, err := service.ListClasses(context.Background(), "school-1", "2024-2025")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, list, tt.wantLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetClass(t *testing.T) {
	tests := []struct {
		name    string
		classID string
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name:    "successful get",
			classID: "class-123",
			mockFn: func(m *MockRepository) {
				m.On("Get", mock.Anything, "class-123").Return(&Class{ID: "class-123", Name: "1A"}, nil)
			},
			wantErr: false,
		},
		{
			name:    "class not found",
			classID: "class-999",
			mockFn: func(m *MockRepository) {
				m.On("Get", mock.Anything, "class-999").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)
			service := NewService(mockRepo)

			class, err := service.GetClass(context.Background(), tt.classID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, class)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, class)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetTeacherClasses(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(*MockRepository)
		wantErr bool
		wantLen int
	}{
		{
			name: "successful get teacher classes",
			mockFn: func(m *MockRepository) {
				classes := []Class{{ID: "c1", Name: "1A"}, {ID: "c2", Name: "2B"}}
				m.On("ListByTeacher", mock.Anything, "teacher-123").Return(classes, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "repository error",
			mockFn: func(m *MockRepository) {
				m.On("ListByTeacher", mock.Anything, "teacher-123").Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)
			service := NewService(mockRepo)

			classes, err := service.GetTeacherClasses(context.Background(), "teacher-123")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, classes, tt.wantLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_UpdateClass(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful update",
			mockFn: func(m *MockRepository) {
				existing := &Class{ID: "class-123", SchoolID: "school-1", Name: "1A"}
				m.On("Get", mock.Anything, "class-123").Return(existing, nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("*classes.Class")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "class not found",
			mockFn: func(m *MockRepository) {
				m.On("Get", mock.Anything, "class-999").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "update error",
			mockFn: func(m *MockRepository) {
				existing := &Class{ID: "class-123", SchoolID: "school-1", Name: "1A"}
				m.On("Get", mock.Anything, "class-123").Return(existing, nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("*classes.Class")).Return(errors.New("update failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)
			service := NewService(mockRepo)

			req := CreateClassRequest{Name: "1A Updated", Section: "A", AcademicYear: "2024-2025"}
			classID := "class-123"
			if tt.name == "class not found" {
				classID = "class-999"
			}
			class, err := service.UpdateClass(context.Background(), "school-1", classID, req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, class)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_DeleteClass(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful delete",
			mockFn: func(m *MockRepository) {
				m.On("Get", mock.Anything, "class-123").Return(&Class{ID: "class-123", SchoolID: "school-1"}, nil)
				m.On("Delete", mock.Anything, "class-123").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			mockFn: func(m *MockRepository) {
				m.On("Get", mock.Anything, "class-123").Return(&Class{ID: "class-123", SchoolID: "school-1"}, nil)
				m.On("Delete", mock.Anything, "class-123").Return(errors.New("delete failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)
			service := NewService(mockRepo)

			err := service.DeleteClass(context.Background(), "school-1", "class-123")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_AssignSubject(t *testing.T) {
	teacherID := "teacher-123"

	tests := []struct {
		name    string
		req     AssignSubjectRequest
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful assignment with teacher",
			req: AssignSubjectRequest{
				SubjectID:    "subject-456",
				TeacherID:    &teacherID,
				HoursPerWeek: 4.0,
			},
			mockFn: func(m *MockRepository) {
				m.On("Get", mock.Anything, "class-123").Return(&Class{ID: "class-123", SchoolID: "school-1"}, nil)
				m.On("AssignSubject", mock.Anything, "class-123", "subject-456", &teacherID, 4.0).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "successful assignment without teacher",
			req: AssignSubjectRequest{
				SubjectID:    "subject-789",
				TeacherID:    nil,
				HoursPerWeek: 3.0,
			},
			mockFn: func(m *MockRepository) {
				var nilTeacher *string
				m.On("Get", mock.Anything, "class-123").Return(&Class{ID: "class-123", SchoolID: "school-1"}, nil)
				m.On("AssignSubject", mock.Anything, "class-123", "subject-789", nilTeacher, 3.0).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			req: AssignSubjectRequest{
				SubjectID:    "subject-456",
				HoursPerWeek: 4.0,
			},
			mockFn: func(m *MockRepository) {
				var nilTeacher *string
				m.On("Get", mock.Anything, "class-123").Return(&Class{ID: "class-123", SchoolID: "school-1"}, nil)
				m.On("AssignSubject", mock.Anything, "class-123", "subject-456", nilTeacher, 4.0).Return(errors.New("assignment failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)
			service := NewService(mockRepo)

			err := service.AssignSubject(context.Background(), "school-1", "class-123", tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_RemoveSubject(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful removal",
			mockFn: func(m *MockRepository) {
				m.On("UnassignSubject", mock.Anything, "assignment-123").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			mockFn: func(m *MockRepository) {
				m.On("UnassignSubject", mock.Anything, "assignment-123").Return(errors.New("removal failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)
			service := NewService(mockRepo)

			err := service.RemoveSubject(context.Background(), "assignment-123")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetClassSubjects(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(*MockRepository)
		wantErr bool
		wantLen int
	}{
		{
			name: "successful get subjects",
			mockFn: func(m *MockRepository) {
				subjects := []ClassSubject{
					{SubjectID: "s1", SubjectName: "Math"},
					{SubjectID: "s2", SubjectName: "English"},
				}
				m.On("GetClassSubjects", mock.Anything, "class-123").Return(subjects, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "empty list",
			mockFn: func(m *MockRepository) {
				m.On("GetClassSubjects", mock.Anything, "class-123").Return([]ClassSubject{}, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "repository error",
			mockFn: func(m *MockRepository) {
				m.On("GetClassSubjects", mock.Anything, "class-123").Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)
			service := NewService(mockRepo)

			subjects, err := service.GetClassSubjects(context.Background(), "class-123")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, subjects, tt.wantLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewService(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}

func TestService_GetLessonTopics(t *testing.T) {
	mockRepo := new(MockRepository)
	expected := []LessonTopic{
		{ID: "lt-1", Topic: "Algebra"},
	}
	mockRepo.On("GetLessonTopics", mock.Anything, "class-1").Return(expected, nil)

	service := NewService(mockRepo)
	topics, err := service.GetLessonTopics(context.Background(), "class-1")

	assert.NoError(t, err)
	assert.Equal(t, expected, topics)
	mockRepo.AssertExpectations(t)
}

func TestService_GetDisciplinaryNotes(t *testing.T) {
	mockRepo := new(MockRepository)
	expected := []DisciplinaryNoteReport{
		{ID: "dn-1", Description: "Disturbo in classe"},
	}
	mockRepo.On("GetDisciplinaryNotes", mock.Anything, "class-1").Return(expected, nil)

	service := NewService(mockRepo)
	notes, err := service.GetDisciplinaryNotes(context.Background(), "class-1")

	assert.NoError(t, err)
	assert.Equal(t, expected, notes)
	mockRepo.AssertExpectations(t)
}

