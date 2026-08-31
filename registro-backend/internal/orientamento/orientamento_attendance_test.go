package orientamento_test

import (
	"context"
	"testing"
	"time"

	"registro-backend/internal/orientamento"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockOrientRepo struct {
	mock.Mock
}

func (m *mockOrientRepo) CreateEvent(ctx context.Context, e *orientamento.Event) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *mockOrientRepo) GetEvents(ctx context.Context, schoolID string) ([]orientamento.Event, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).([]orientamento.Event), args.Error(1)
}

func (m *mockOrientRepo) RegisterStudent(ctx context.Context, p *orientamento.Participation) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockOrientRepo) GetParticipations(ctx context.Context, studentID string) ([]orientamento.Participation, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]orientamento.Participation), args.Error(1)
}

func (m *mockOrientRepo) MarkAttendance(ctx context.Context, eventID, studentID string, attended bool) error {
	args := m.Called(ctx, eventID, studentID, attended)
	return args.Error(0)
}

func (m *mockOrientRepo) SavePreference(ctx context.Context, p *orientamento.StudentPreference) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockOrientRepo) GetPreference(ctx context.Context, studentID string) (*orientamento.StudentPreference, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orientamento.StudentPreference), args.Error(1)
}

func (m *mockOrientRepo) SaveCapolavoro(ctx context.Context, c *orientamento.Capolavoro) error {
	return nil
}
func (m *mockOrientRepo) GetCapolavori(ctx context.Context, studentID string) ([]orientamento.Capolavoro, error) {
	return nil, nil
}
func (m *mockOrientRepo) GetCurriculumStudente(ctx context.Context, studentID string) (*orientamento.CurriculumStudenteSummary, error) {
	return nil, nil
}

func TestOrientamento_MarkAttendance(t *testing.T) {
	mockRepo := new(mockOrientRepo)
	svc := orientamento.NewService(mockRepo)
	ctx := context.Background()

	mockRepo.On("MarkAttendance", ctx, "event-123", "student-456", true).Return(nil).Once()

	err := svc.MarkAttendance(ctx, "event-123", "student-456")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestOrientamento_PreventDuplicateRegistration(t *testing.T) {
	mockRepo := new(mockOrientRepo)
	svc := orientamento.NewService(mockRepo)
	ctx := context.Background()

	existing := []orientamento.Participation{
		{
			ID:           "part-1",
			EventID:      "event-123",
			StudentID:    "student-456",
			Status:       "Registered",
			RegisteredAt: time.Now(),
		},
	}

	mockRepo.On("GetParticipations", ctx, "student-456").Return(existing, nil).Once()

	err := svc.RegisterStudent(ctx, "student-456", "event-123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "già iscritto")
	mockRepo.AssertExpectations(t)
}
