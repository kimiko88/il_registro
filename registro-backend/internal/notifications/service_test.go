package notifications

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveToken(ctx context.Context, token *PushToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
func (m *MockRepository) GetUserTokens(ctx context.Context, userID string) ([]PushToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]PushToken), args.Error(1)
}
func (m *MockRepository) DeleteToken(ctx context.Context, userID, deviceToken string) error {
	args := m.Called(ctx, userID, deviceToken)
	return args.Error(0)
}
func (m *MockRepository) CreateDBNotification(ctx context.Context, n *DBNotification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}
func (m *MockRepository) ListDBNotifications(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]DBNotification, error) {
	args := m.Called(ctx, userID, unreadOnly, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DBNotification), args.Error(1)
}
func (m *MockRepository) MarkAsRead(ctx context.Context, userID, notificationID string) error {
	args := m.Called(ctx, userID, notificationID)
	return args.Error(0)
}
func (m *MockRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestRegisterAndSendPush(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	mockRepo.On("SaveToken", mock.Anything, mock.MatchedBy(func(t *PushToken) bool {
		return t.UserID == "user-1" && t.DeviceToken == "token-123"
	})).Return(nil).Once()

	err := svc.RegisterToken(context.Background(), "user-1", RegisterTokenRequest{
		DeviceToken: "token-123",
		Platform:    "android",
	})
	assert.NoError(t, err)

	mockRepo.On("GetUserTokens", mock.Anything, "user-1").Return([]PushToken{
		{ID: "1", UserID: "user-1", DeviceToken: "token-123", Platform: "android"},
	}, nil).Once()

	count, err := svc.SendPushNotification(context.Background(), SendNotificationRequest{
		UserID: "user-1",
		Title:  "Nuovo Voto",
		Body:   "Hai ricevuto 8 in Matematica",
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	mockRepo.AssertExpectations(t)
}

func TestDBNotification(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	mockRepo.On("CreateDBNotification", mock.Anything, mock.MatchedBy(func(n *DBNotification) bool {
		return n.UserID == "user-1" && n.Title == "Nuova Assenza"
	})).Return(nil).Once()

	n, err := svc.CreateInAppNotification(context.Background(), "user-1", "Nuova Assenza", "Assenza registrata il 15-11", "absence", nil)
	assert.NoError(t, err)
	assert.NotNil(t, n)

	mockRepo.On("ListDBNotifications", mock.Anything, "user-1", true, 50, 0).Return([]DBNotification{*n}, nil).Once()
	list, err := svc.ListDBNotifications(context.Background(), "user-1", true, 50, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	mockRepo.On("MarkAllAsRead", mock.Anything, "user-1").Return(nil).Once()
	err = svc.MarkAllAsRead(context.Background(), "user-1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
