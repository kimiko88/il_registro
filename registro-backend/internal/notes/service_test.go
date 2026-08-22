package notes

import (
	"context"
	"errors"
	"testing"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, note *StudentNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockRepository) Get(ctx context.Context, id string) (*StudentNote, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*StudentNote), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, note *StudentNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) DeleteWithReason(ctx context.Context, id string, reason string) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "DeleteWithReason" {
			args := m.Called(ctx, id, reason)
			return args.Error(0)
		}
	}
	return m.Delete(ctx, id)
}

func (m *MockRepository) MarkAsViewedByParent(ctx context.Context, id string) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "MarkAsViewedByParent" {
			args := m.Called(ctx, id)
			return args.Error(0)
		}
	}
	return nil
}

func (m *MockRepository) MarkManyAsViewedByParent(ctx context.Context, ids []string) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "MarkManyAsViewedByParent" {
			args := m.Called(ctx, ids)
			return args.Error(0)
		}
	}
	return nil
}

func (m *MockRepository) List(ctx context.Context, filter NoteFilter) ([]StudentNote, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]StudentNote), args.Error(1)
}

func (m *MockRepository) ApproveNote(ctx context.Context, id string, approverID string) error {
	args := m.Called(ctx, id, approverID)
	return args.Error(0)
}

func (m *MockRepository) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "IsTeacherAssignedToClass" {
			args := m.Called(ctx, teacherID, classID)
			return args.Bool(0), args.Error(1)
		}
	}
	return true, nil
}

type MockUserRepo struct {
	mock.Mock
	users.Repository
}

func (m *MockUserRepo) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

func TestService_CreateNote(t *testing.T) {
	subjectID := "subject-123"

	tests := []struct {
		name    string
		req     CreateNoteRequest
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name: "successful creation - generic note",
			req: CreateNoteRequest{
				StudentID: "student-123",
				ClassID:   "class-456",
				SubjectID: &subjectID,
				Type:      NoteTypeGeneric,
				Note:      "Good participation",
				Date:      "2024-01-15",
			},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(n *StudentNote) bool {
					return n.StudentID == "student-123" && n.Type == NoteTypeGeneric
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "successful creation - behavior note",
			req: CreateNoteRequest{
				StudentID: "student-123",
				ClassID:   "class-456",
				Type:      NoteTypeBehavior,
				Note:      "Disrupting class",
				Date:      "2024-01-15",
			},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(n *StudentNote) bool {
					return n.Type == NoteTypeBehavior && n.SubjectID == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "successful creation - homework note",
			req: CreateNoteRequest{
				StudentID: "student-123",
				ClassID:   "class-456",
				SubjectID: &subjectID,
				Type:      NoteTypeHomework,
				Note:      "Missing homework",
				Date:      "2024-01-15",
			},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*notes.StudentNote")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			req: CreateNoteRequest{
				StudentID: "student-123",
				ClassID:   "class-456",
				Type:      NoteTypeGeneric,
				Note:      "Test note",
				Date:      "2024-01-15",
			},
			mockFn: func(m *MockRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*notes.StudentNote")).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo, new(MockUserRepo))
			note, err := service.CreateNote(context.Background(), "teacher-789", "school-123", tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, note)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, note)
				assert.Equal(t, tt.req.StudentID, note.StudentID)
				assert.Equal(t, tt.req.Type, note.Type)
				assert.Equal(t, "teacher-789", note.TeacherID)
				assert.Equal(t, "school-123", note.SchoolID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_UpdateNote(t *testing.T) {
	tests := []struct {
		name      string
		teacherID string
		noteID    string
		req       UpdateNoteRequest
		mockFn    func(*MockRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "successful update",
			teacherID: "teacher-123",
			noteID:    "note-456",
			req: UpdateNoteRequest{
				Type: NoteTypeBehavior,
				Note: "Updated note",
				Date: "2024-01-16",
			},
			mockFn: func(m *MockRepository) {
				existingNote := &StudentNote{
					ID:        "note-456",
					TeacherID: "teacher-123",
					Type:      NoteTypeGeneric,
					Note:      "Original note",
				}
				m.On("Get", mock.Anything, "note-456").Return(existingNote, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(n *StudentNote) bool {
					return n.Type == NoteTypeBehavior && n.Note == "Updated note"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "unauthorized - different teacher",
			teacherID: "teacher-123",
			noteID:    "note-456",
			req: UpdateNoteRequest{
				Note: "Updated note",
			},
			mockFn: func(m *MockRepository) {
				existingNote := &StudentNote{
					ID:        "note-456",
					TeacherID: "teacher-999", // Different teacher
					Note:      "Original note",
				}
				m.On("Get", mock.Anything, "note-456").Return(existingNote, nil)
			},
			wantErr: true,
			errMsg:  "unauthorized",
		},
		{
			name:      "note not found",
			teacherID: "teacher-123",
			noteID:    "note-999",
			req: UpdateNoteRequest{
				Note: "Updated note",
			},
			mockFn: func(m *MockRepository) {
				m.On("Get", mock.Anything, "note-999").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:      "update error",
			teacherID: "teacher-123",
			noteID:    "note-456",
			req: UpdateNoteRequest{
				Note: "Updated note",
			},
			mockFn: func(m *MockRepository) {
				existingNote := &StudentNote{
					ID:        "note-456",
					TeacherID: "teacher-123",
					Note:      "Original note",
				}
				m.On("Get", mock.Anything, "note-456").Return(existingNote, nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("*notes.StudentNote")).Return(errors.New("update failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo, new(MockUserRepo))
			note, err := service.UpdateNote(context.Background(), tt.teacherID, "teacher", tt.noteID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, note)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_DeleteNote(t *testing.T) {
	tests := []struct {
		name      string
		teacherID string
		noteID    string
		mockFn    func(*MockRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "successful delete",
			teacherID: "teacher-123",
			noteID:    "note-456",
			mockFn: func(m *MockRepository) {
				existingNote := &StudentNote{
					ID:        "note-456",
					TeacherID: "teacher-123",
				}
				m.On("Get", mock.Anything, "note-456").Return(existingNote, nil)
				m.On("Delete", mock.Anything, "note-456").Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "unauthorized - different teacher",
			teacherID: "teacher-123",
			noteID:    "note-456",
			mockFn: func(m *MockRepository) {
				existingNote := &StudentNote{
					ID:        "note-456",
					TeacherID: "teacher-999",
				}
				m.On("Get", mock.Anything, "note-456").Return(existingNote, nil)
			},
			wantErr: true,
			errMsg:  "unauthorized",
		},
		{
			name:      "note not found",
			teacherID: "teacher-123",
			noteID:    "note-999",
			mockFn: func(m *MockRepository) {
				m.On("Get", mock.Anything, "note-999").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo, new(MockUserRepo))
			err := service.DeleteNote(context.Background(), tt.teacherID, "teacher", tt.noteID, "")

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_ListNotes(t *testing.T) {
	tests := []struct {
		name    string
		filter  NoteFilter
		mockFn  func(*MockRepository)
		wantErr bool
		wantLen int
	}{
		{
			name: "successful list by student",
			filter: NoteFilter{
				StudentID: "student-123",
			},
			mockFn: func(m *MockRepository) {
				notes := []StudentNote{
					{ID: "n1", StudentID: "student-123", Type: NoteTypeGeneric},
					{ID: "n2", StudentID: "student-123", Type: NoteTypeBehavior},
				}
				m.On("List", mock.Anything, mock.MatchedBy(func(f NoteFilter) bool {
					return f.StudentID == "student-123"
				})).Return(notes, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "successful list by class",
			filter: NoteFilter{
				ClassID: "class-456",
			},
			mockFn: func(m *MockRepository) {
				notes := []StudentNote{
					{ID: "n1", ClassID: "class-456"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("notes.NoteFilter")).Return(notes, nil)
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "empty list",
			filter: NoteFilter{
				StudentID: "student-999",
			},
			mockFn: func(m *MockRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("notes.NoteFilter")).Return([]StudentNote{}, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "repository error",
			filter: NoteFilter{
				StudentID: "student-123",
			},
			mockFn: func(m *MockRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("notes.NoteFilter")).Return(nil, errors.New("database error"))
			},
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockFn(mockRepo)

			service := NewService(mockRepo, new(MockUserRepo))
			notes, err := service.ListNotes(context.Background(), tt.filter)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, notes, tt.wantLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewService(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUserRepo := new(MockUserRepo)
	service := NewService(mockRepo, mockUserRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}
