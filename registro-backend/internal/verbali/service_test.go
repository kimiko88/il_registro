package verbali

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
func (m *MockRepository) ClassBelongsToSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	args := m.Called(ctx, classID, schoolID)
	return args.Bool(0), args.Error(1)
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

	mockRepo.On("ClassBelongsToSchool", mock.Anything, "class-1", "school-1").Return(true, nil).Once()
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

func TestCreateMeeting_CrossTenantBlocked(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	reqMeeting := CreateMeetingRequest{
		ClassID:   "class-other-school",
		Title:     "Consiglio di Classe",
		Date:      "2025-11-10",
		StartTime: "16:00",
		EndTime:   "17:30",
	}

	mockRepo.On("ClassBelongsToSchool", mock.Anything, "class-other-school", "school-1").Return(false, nil).Once()

	_, err := svc.CreateMeeting(context.Background(), "t-1", "school-1", reqMeeting)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong to school")
	mockRepo.AssertExpectations(t)
}

func TestSignVerbale(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	secID := "u-1"
	verbale := &MeetingVerbale{ID: "v-1", SecretaryID: &secID, IsPublished: true}
	mockRepo.On("GetVerbaleByID", mock.Anything, "v-1", "u-1").Return(verbale, nil).Once()
	mockRepo.On("SignVerbale", mock.Anything, "v-1", "u-1", "10.0.0.1").Return(nil).Once()

	err := svc.SignVerbale(ctx, "v-1", "u-1", "10.0.0.1")
	assert.NoError(t, err)

	// Unpublished error
	unpub := &MeetingVerbale{ID: "v-2", SecretaryID: &secID, IsPublished: false}
	mockRepo.On("GetVerbaleByID", mock.Anything, "v-2", "u-1").Return(unpub, nil).Once()
	err = svc.SignVerbale(ctx, "v-2", "u-1", "10.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot sign an unpublished verbale")

	// Unauthorized signer error
	mockRepo.On("GetVerbaleByID", mock.Anything, "v-1", "other-user").Return(verbale, nil).Once()
	err = svc.SignVerbale(ctx, "v-1", "other-user", "10.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized: non sei il segretario o il presidente")

	// Not found error
	mockRepo.On("GetVerbaleByID", mock.Anything, "v-missing", "u-1").Return(nil, errors.New("not found")).Once()
	err = svc.SignVerbale(ctx, "v-missing", "u-1", "10.0.0.1")
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestService_AdditionalMethods(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	// 1. CreateVerbale unauthorized role
	_, err := svc.CreateVerbale(ctx, "u-1", "student", CreateVerbaleRequest{MeetingID: "m-1"})
	assert.ErrorIs(t, err, ErrUnauthorized)

	// 2. CreateVerbale meeting not found
	mockRepo.On("GetMeetingByID", ctx, "m-missing").Return(nil, errors.New("not found")).Once()
	_, err = svc.CreateVerbale(ctx, "u-1", "teacher", CreateVerbaleRequest{MeetingID: "m-missing"})
	assert.ErrorIs(t, err, ErrNotFound)

	// 3. ListMeetings
	mockRepo.On("ListMeetings", ctx, "school-1", "c-1").Return([]*CouncilMeeting{{ID: "m-1"}}, nil).Once()
	meetings, err := svc.ListMeetings(ctx, "school-1", "c-1")
	assert.NoError(t, err)
	assert.Len(t, meetings, 1)

	// 4. GetVerbale
	mockRepo.On("GetVerbaleByID", ctx, "v-1", "u-1").Return(&MeetingVerbale{ID: "v-1"}, nil).Once()
	v, err := svc.GetVerbale(ctx, "v-1", "u-1")
	assert.NoError(t, err)
	assert.Equal(t, "v-1", v.ID)

	// 5. ListVerbali
	mockRepo.On("ListVerbali", ctx, "m-1", "u-1").Return([]*MeetingVerbale{{ID: "v-1"}}, nil).Once()
	verbali, err := svc.ListVerbali(ctx, "m-1", "u-1")
	assert.NoError(t, err)
	assert.Len(t, verbali, 1)

	// 6. GetSignatures
	mockRepo.On("GetSignatures", ctx, "v-1").Return([]VerbaleSignature{{ID: "sig-1"}}, nil).Once()
	sigs, err := svc.GetSignatures(ctx, "v-1")
	assert.NoError(t, err)
	assert.Len(t, sigs, 1)

	// 7. GetMeeting
	mockRepo.On("GetMeetingByID", ctx, "m-1").Return(&CouncilMeeting{ID: "m-1"}, nil).Once()
	m, err := svc.GetMeeting(ctx, "m-1")
	assert.NoError(t, err)
	assert.Equal(t, "m-1", m.ID)

	mockRepo.AssertExpectations(t)
}
