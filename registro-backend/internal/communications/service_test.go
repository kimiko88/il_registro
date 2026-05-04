package communications

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, msg *Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, userID string) ([]*Message, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*Message), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_SendMessage(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := CreateMessageRequest{
		Recipients: []string{"user-2"},
		Subject:    "Hello",
		Body:       "World",
		Type:       "email",
	}
	senderID := "user-1"

	// Expect creation
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(m *Message) bool {
		return m.SenderID == senderID && m.Subject == "Hello"
	})).Return(nil)

	msg, err := svc.SendMessage(context.Background(), senderID, req)
	assert.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, senderID, msg.SenderID)
}

func TestService_ListMessages(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	userID := "user-1"
	expected := []*Message{
		{ID: "msg-1", SenderID: "user-2", Subject: "Hi"},
	}

	mockRepo.On("List", mock.Anything, userID).Return(expected, nil)

	msgs, err := svc.ListMessages(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
	assert.Equal(t, "Hi", msgs[0].Subject)
}
