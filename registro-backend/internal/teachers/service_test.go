package teachers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) List(ctx context.Context, schoolID string) ([]Teacher, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Teacher), args.Error(1)
}

func (m *MockRepository) Get(ctx context.Context, id string) (*Teacher, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Teacher), args.Error(1)
}

func (m *MockRepository) GetByUserID(ctx context.Context, userID string) (*Teacher, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Teacher), args.Error(1)
}

func (m *MockRepository) GetSubjects(ctx context.Context, teacherID string) ([]TeacherSubject, error) {
	args := m.Called(ctx, teacherID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]TeacherSubject), args.Error(1)
}

func (m *MockRepository) GetBySubject(ctx context.Context, subjectID string) ([]Teacher, error) {
	args := m.Called(ctx, subjectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Teacher), args.Error(1)
}

func (m *MockRepository) AssignSubject(ctx context.Context, teacherID, subjectID string) error {
	args := m.Called(ctx, teacherID, subjectID)
	return args.Error(0)
}

func (m *MockRepository) RemoveSubject(ctx context.Context, teacherID, subjectID string) error {
	args := m.Called(ctx, teacherID, subjectID)
	return args.Error(0)
}

func (m *MockRepository) GetDashboardStats(ctx context.Context, teacherUserID string) (map[string]interface{}, error) {
	args := m.Called(ctx, teacherUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestService_ListTeachers(t *testing.T) {
	tests := []struct {
		name     string
		schoolID string
		mockFn   func(*MockRepository)
		wantErr  bool
		wantLen  int
	}{
		{
			name:     "successful list",
			schoolID: "school-123",
			mockFn: func(m *MockRepository) {
				teachers := []Teacher{
					{ID: "t1", UserID: "u1", SchoolID: "school-123", FirstName: "John", LastName: "Doe"},
					{ID: "t2", UserID: "u2", SchoolID: "school-123", FirstName: "Jane", LastName: "Smith"},
				}
				m.On("List", mock.Anything, "school-123").Return(teachers, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:     "empty list",
			schoolID: "school-456",
			mockFn: func(m *MockRepository) {
				m.On("List", mock.Anything, "school-456").Return([]Teacher{}, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:     "repository error",
			schoolID: "school-789",
			mockFn: func(m *MockRepository) {
				m.On("List", mock.Anything, "school-789").Return(nil, errors.New("database error"))
			},
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo)
			teachers, err := service.ListTeachers(context.Background(), tt.schoolID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, teachers, tt.wantLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetTeacher(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful get",
			id:   "teacher-123",
			mockFn: func(m *MockRepository) {
				teacher := &Teacher{ID: "teacher-123", UserID: "user-1", FirstName: "John"}
				m.On("Get", mock.Anything, "teacher-123").Return(teacher, nil)
			},
			wantErr: false,
		},
		{
			name: "teacher not found",
			id:   "teacher-999",
			mockFn: func(m *MockRepository) {
				m.On("Get", mock.Anything, "teacher-999").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo)
			teacher, err := service.GetTeacher(context.Background(), "", tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, teacher)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, teacher)
				assert.Equal(t, tt.id, teacher.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetTeacherSubjects(t *testing.T) {
	tests := []struct {
		name      string
		teacherID string
		mockFn    func(*MockRepository)
		wantErr   bool
		wantLen   int
	}{
		{
			name:      "successful get subjects",
			teacherID: "teacher-123",
			mockFn: func(m *MockRepository) {
				subjects := []TeacherSubject{
					{ID: "ts1", TeacherID: "teacher-123", SubjectID: "sub-1", SubjectName: "Math"},
					{ID: "ts2", TeacherID: "teacher-123", SubjectID: "sub-2", SubjectName: "Physics"},
				}
				m.On("GetSubjects", mock.Anything, "teacher-123").Return(subjects, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:      "teacher has no subjects",
			teacherID: "teacher-456",
			mockFn: func(m *MockRepository) {
				m.On("GetSubjects", mock.Anything, "teacher-456").Return([]TeacherSubject{}, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:      "repository error",
			teacherID: "teacher-789",
			mockFn: func(m *MockRepository) {
				m.On("GetSubjects", mock.Anything, "teacher-789").Return(nil, errors.New("database error"))
			},
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo)
			subjects, err := service.GetTeacherSubjects(context.Background(), tt.teacherID)

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

func TestService_GetTeachersBySubject(t *testing.T) {
	tests := []struct {
		name      string
		subjectID string
		mockFn    func(*MockRepository)
		wantErr   bool
		wantLen   int
	}{
		{
			name:      "successful get teachers",
			subjectID: "subject-123",
			mockFn: func(m *MockRepository) {
				teachers := []Teacher{
					{ID: "t1", UserID: "u1", FirstName: "John"},
					{ID: "t2", UserID: "u2", FirstName: "Jane"},
				}
				m.On("GetBySubject", mock.Anything, "subject-123").Return(teachers, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:      "no teachers for subject",
			subjectID: "subject-456",
			mockFn: func(m *MockRepository) {
				m.On("GetBySubject", mock.Anything, "subject-456").Return([]Teacher{}, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:      "repository error",
			subjectID: "subject-789",
			mockFn: func(m *MockRepository) {
				m.On("GetBySubject", mock.Anything, "subject-789").Return(nil, errors.New("database error"))
			},
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo)
			teachers, err := service.GetTeachersBySubject(context.Background(), tt.subjectID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, teachers, tt.wantLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_AssignSubject(t *testing.T) {
	tests := []struct {
		name      string
		teacherID string
		subjectID string
		mockFn    func(*MockRepository)
		wantErr   bool
	}{
		{
			name:      "successful assignment",
			teacherID: "teacher-123",
			subjectID: "subject-456",
			mockFn: func(m *MockRepository) {
				m.On("AssignSubject", mock.Anything, "teacher-123", "subject-456").Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "assignment error",
			teacherID: "teacher-123",
			subjectID: "subject-999",
			mockFn: func(m *MockRepository) {
				m.On("AssignSubject", mock.Anything, "teacher-123", "subject-999").Return(errors.New("assignment failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo)
			err := service.AssignSubject(context.Background(), "admin", "", tt.teacherID, tt.subjectID)

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
		name      string
		teacherID string
		subjectID string
		mockFn    func(*MockRepository)
		wantErr   bool
	}{
		{
			name:      "successful removal",
			teacherID: "teacher-123",
			subjectID: "subject-456",
			mockFn: func(m *MockRepository) {
				m.On("RemoveSubject", mock.Anything, "teacher-123", "subject-456").Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "removal error",
			teacherID: "teacher-123",
			subjectID: "subject-999",
			mockFn: func(m *MockRepository) {
				m.On("RemoveSubject", mock.Anything, "teacher-123", "subject-999").Return(errors.New("removal failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo)
			err := service.RemoveSubject(context.Background(), "admin", "", tt.teacherID, tt.subjectID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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

func TestService_GetDashboardStats(t *testing.T) {
	ctx := context.Background()

	// 1. Unauthorized when accessing another teacher's dashboard without admin/secretary role
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	_, err := svc.GetDashboardStats(ctx, "teacher-1", "teacher-2", "teacher")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")

	// 2. Authorized when accessing own dashboard
	mockRepo.On("GetDashboardStats", ctx, "teacher-1").Return(map[string]interface{}{"ok": true}, nil).Once()
	stats, err := svc.GetDashboardStats(ctx, "teacher-1", "teacher-1", "teacher")
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// 3. Authorized when admin accesses another teacher's dashboard
	mockRepo.On("GetDashboardStats", ctx, "teacher-2").Return(map[string]interface{}{"ok": true}, nil).Once()
	stats, err = svc.GetDashboardStats(ctx, "admin-1", "teacher-2", "admin")
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	mockRepo.AssertExpectations(t)
}

func TestService_AssignRemoveSubject_Permissions(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	// Unauthorized role for Assign
	err := svc.AssignSubject(ctx, "student", "school-1", "t-1", "s-1")
	assert.Error(t, err)

	// Unauthorized role for Remove
	err = svc.RemoveSubject(ctx, "parent", "school-1", "t-1", "s-1")
	assert.Error(t, err)

	// Cross school for Assign
	mockRepo.On("Get", ctx, "t-1").Return(&Teacher{ID: "t-1", SchoolID: "school-A"}, nil).Once()
	err = svc.AssignSubject(ctx, "admin", "school-B", "t-1", "s-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "belongs to another school")

	// Cross school for Remove
	mockRepo.On("Get", ctx, "t-1").Return(&Teacher{ID: "t-1", SchoolID: "school-A"}, nil).Once()
	err = svc.RemoveSubject(ctx, "admin", "school-B", "t-1", "s-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "belongs to another school")

	mockRepo.AssertExpectations(t)
}
