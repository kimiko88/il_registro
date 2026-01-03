package subjects

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

func (m *MockRepository) Create(ctx context.Context, subject *Subject) error {
	args := m.Called(ctx, subject)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, schoolID string) ([]Subject, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Subject), args.Error(1)
}

func (m *MockRepository) Get(ctx context.Context, id string) (*Subject, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Subject), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, subject *Subject) error {
	args := m.Called(ctx, subject)
	return args.Error(0)
}

func TestService_CreateSubject(t *testing.T) {
	mandatory := true
	optional := false

	tests := []struct {
		name    string
		req     CreateSubjectRequest
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful creation with mandatory subject",
			req: CreateSubjectRequest{
				Name:        "Mathematics",
				Code:        "MATH101",
				Description: "Basic mathematics",
				IsMandatory: &mandatory,
			},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(s *Subject) bool {
					return s.Name == "Mathematics" && s.IsMandatory == true
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "successful creation with optional subject",
			req: CreateSubjectRequest{
				Name:        "Art",
				Code:        "ART101",
				Description: "Art class",
				IsMandatory: &optional,
			},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(s *Subject) bool {
					return s.Name == "Art" && s.IsMandatory == false
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "default to mandatory when not specified",
			req: CreateSubjectRequest{
				Name:        "Science",
				Code:        "SCI101",
				Description: "Science class",
				IsMandatory: nil,
			},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(s *Subject) bool {
					return s.Name == "Science" && s.IsMandatory == true
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			req: CreateSubjectRequest{
				Name: "English",
				Code: "ENG101",
			},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*subjects.Subject")).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo)
			subject, err := service.CreateSubject(context.Background(), "school-123", tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, subject)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, subject)
				assert.Equal(t, tt.req.Name, subject.Name)
				assert.Equal(t, tt.req.Code, subject.Code)
				assert.Equal(t, "school-123", subject.SchoolID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_ListSubjects(t *testing.T) {
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
				subjects := []Subject{
					{ID: "1", Name: "Math", SchoolID: "school-123"},
					{ID: "2", Name: "Science", SchoolID: "school-123"},
				}
				m.On("List", mock.Anything, "school-123").Return(subjects, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:     "empty list",
			schoolID: "school-456",
			mockFn: func(m *MockRepository) {
				m.On("List", mock.Anything, "school-456").Return([]Subject{}, nil)
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
			subjects, err := service.ListSubjects(context.Background(), tt.schoolID)

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

func TestService_GetSubject(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful get",
			id:   "subject-123",
			mockFn: func(m *MockRepository) {
				subject := &Subject{ID: "subject-123", Name: "Math"}
				m.On("Get", mock.Anything, "subject-123").Return(subject, nil)
			},
			wantErr: false,
		},
		{
			name: "subject not found",
			id:   "subject-999",
			mockFn: func(m *MockRepository) {
				m.On("Get", mock.Anything, "subject-999").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo)
			subject, err := service.GetSubject(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, subject)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, subject)
				assert.Equal(t, tt.id, subject.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_UpdateSubject(t *testing.T) {
	mandatory := false

	tests := []struct {
		name    string
		id      string
		req     CreateSubjectRequest
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful update",
			id:   "subject-123",
			req: CreateSubjectRequest{
				Name:        "Updated Math",
				Code:        "MATH102",
				Description: "Advanced math",
				IsMandatory: &mandatory,
			},
			mockFn: func(m *MockRepository) {
				existing := &Subject{ID: "subject-123", Name: "Math", IsMandatory: true}
				m.On("Get", mock.Anything, "subject-123").Return(existing, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(s *Subject) bool {
					return s.Name == "Updated Math" && s.IsMandatory == false
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "subject not found",
			id:   "subject-999",
			req: CreateSubjectRequest{
				Name: "Test",
			},
			mockFn: func(m *MockRepository) {
				m.On("Get", mock.Anything, "subject-999").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "update error",
			id:   "subject-123",
			req: CreateSubjectRequest{
				Name: "Updated Math",
			},
			mockFn: func(m *MockRepository) {
				existing := &Subject{ID: "subject-123", Name: "Math"}
				m.On("Get", mock.Anything, "subject-123").Return(existing, nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("*subjects.Subject")).Return(errors.New("update failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo)
			subject, err := service.UpdateSubject(context.Background(), tt.id, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, subject)
				assert.Equal(t, tt.req.Name, subject.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_DeleteSubject(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful delete",
			id:   "subject-123",
			mockFn: func(m *MockRepository) {
				m.On("Delete", mock.Anything, "subject-123").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "delete error",
			id:   "subject-999",
			mockFn: func(m *MockRepository) {
				m.On("Delete", mock.Anything, "subject-999").Return(errors.New("delete failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo)
			err := service.DeleteSubject(context.Background(), tt.id)

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
