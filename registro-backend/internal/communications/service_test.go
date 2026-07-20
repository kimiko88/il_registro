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

func (m *MockRepository) ListBacheca(ctx context.Context, schoolID, userID string) ([]*Message, error) {
	args := m.Called(ctx, schoolID, userID)
	return args.Get(0).([]*Message), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) Sign(ctx context.Context, communicationID string, userID string) error {
	args := m.Called(ctx, communicationID, userID)
	return args.Error(0)
}

func (m *MockRepository) SignWithIP(ctx context.Context, communicationID string, userID string, ipAddress string) error {
	args := m.Called(ctx, communicationID, userID, ipAddress)
	return args.Error(0)
}

func (m *MockRepository) GetSignatures(ctx context.Context, communicationID string) ([]string, error) {
	args := m.Called(ctx, communicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRepository) GetSignatureReport(ctx context.Context, communicationID string) (*SignatureReportResponse, error) {
	args := m.Called(ctx, communicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SignatureReportResponse), args.Error(1)
}

func (m *MockRepository) Get(ctx context.Context, id string) (*Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Message), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, id string, subject, body string) error {
	args := m.Called(ctx, id, subject, body)
	return args.Error(0)
}

func (m *MockRepository) MarkAsRead(ctx context.Context, communicationID, userID, ipAddress string) error {
	args := m.Called(ctx, communicationID, userID, ipAddress)
	return args.Error(0)
}

func (m *MockRepository) GetUnreadUsers(ctx context.Context, communicationID string) ([]string, error) {
	args := m.Called(ctx, communicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func TestService_SendMessage(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := CreateMessageRequest{
		Recipients:        []string{"user-2"},
		Subject:           "Circolare N.12",
		Body:              "Chiusura scuola per festività",
		Type:              "circular",
		RequiresSignature: true,
	}
	senderID := "user-1"

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(m *Message) bool {
		return m.SenderID == senderID && m.Subject == "Circolare N.12" && m.RequiresSignature == true
	})).Return(nil)

	msg, err := svc.SendMessage(context.Background(), senderID, req)
	assert.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, senderID, msg.SenderID)
	assert.True(t, msg.RequiresSignature)
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

func TestService_SignMessageWithIP(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	mockRepo.On("SignWithIP", mock.Anything, "msg-1", "user-1", "192.168.1.50").Return(nil)

	err := svc.SignMessageWithIP(context.Background(), "msg-1", "user-1", "192.168.1.50")
	assert.NoError(t, err)
}

func TestService_SignatureReport(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	expectedReport := &SignatureReportResponse{
		CommunicationID: "msg-1",
		Subject:         "Circolare 10",
		TotalRecipients: 30,
		SignedCount:     25,
		PendingCount:    5,
	}

	mockRepo.On("GetSignatureReport", mock.Anything, "msg-1").Return(expectedReport, nil).Once()

	report, err := svc.GetSignatureReport(context.Background(), "admin", "msg-1")
	assert.NoError(t, err)
	assert.Equal(t, 25, report.SignedCount)
	assert.Equal(t, 5, report.PendingCount)
}
