package verbali

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateMeeting(ctx context.Context, meeting *CouncilMeeting) error {
	args := m.Called(ctx, meeting)
	return args.Error(0)
}
func (m *MockRepository) ListMeetings(ctx context.Context, schoolID, classID string) ([]*CouncilMeeting, error) {
	args := m.Called(ctx, schoolID, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*CouncilMeeting), args.Error(1)
}
func (m *MockRepository) GetMeetingByID(ctx context.Context, id string) (*CouncilMeeting, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CouncilMeeting), args.Error(1)
}
func (m *MockRepository) CreateVerbale(ctx context.Context, v *MeetingVerbale) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}
func (m *MockRepository) GetVerbaleByID(ctx context.Context, id string, userID string) (*MeetingVerbale, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MeetingVerbale), args.Error(1)
}
func (m *MockRepository) ListVerbali(ctx context.Context, meetingID string, userID string) ([]*MeetingVerbale, error) {
	args := m.Called(ctx, meetingID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MeetingVerbale), args.Error(1)
}
func (m *MockRepository) SignVerbale(ctx context.Context, verbaleID, userID, ipAddress string) error {
	args := m.Called(ctx, verbaleID, userID, ipAddress)
	return args.Error(0)
}
func (m *MockRepository) GetSignatures(ctx context.Context, verbaleID string) ([]VerbaleSignature, error) {
	args := m.Called(ctx, verbaleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]VerbaleSignature), args.Error(1)
}

func TestCreateMeetingAndVerbale(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	reqMeeting := CreateMeetingRequest{
		ClassID:   "class-1",
		Title:     "Consiglio di Classe di Novembre",
		Date:      "2025-11-10",
		StartTime: "16:00",
		EndTime:   "17:30",
		Agenda:    "Approvazione piano didattico",
	}

	mockRepo.On("CreateMeeting", mock.Anything, mock.MatchedBy(func(m *CouncilMeeting) bool {
		return m.Title == "Consiglio di Classe di Novembre"
	})).Return(nil).Once()

	m, err := svc.CreateMeeting(context.Background(), "t-1", "school-1", reqMeeting)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	mockRepo.On("GetMeetingByID", mock.Anything, "m-1").Return(m, nil).Once()
	mockRepo.On("CreateVerbale", mock.Anything, mock.MatchedBy(func(v *MeetingVerbale) bool {
		return v.Title == "Verbale n.1"
	})).Return(nil).Once()

	v, err := svc.CreateVerbale(context.Background(), "t-1", "teacher", CreateVerbaleRequest{
		MeetingID:   "m-1",
		Title:       "Verbale n.1",
		Content:     "Si approva il piano didattico all'unanimità.",
		IsPublished: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, v)
	mockRepo.AssertExpectations(t)
}

func TestSignVerbale(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	verbale := &MeetingVerbale{ID: "v-1", IsPublished: true}
	mockRepo.On("GetVerbaleByID", mock.Anything, "v-1", "u-1").Return(verbale, nil).Once()
	mockRepo.On("SignVerbale", mock.Anything, "v-1", "u-1", "10.0.0.1").Return(nil).Once()

	err := svc.SignVerbale(context.Background(), "v-1", "u-1", "10.0.0.1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
