package signatures

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepo
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(ctx context.Context, sig *Signature) error {
	args := m.Called(sig)
	sig.ID = "sig-1"
	sig.SignedAt = time.Now()
	return args.Error(0)
}

func (m *MockRepo) FindByDocumentID(ctx context.Context, docID string) ([]Signature, error) {
	args := m.Called(docID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Signature), args.Error(1)
}

// MockDocs
type MockDocs struct {
	mock.Mock
}

func (m *MockDocs) LockDocument(ctx context.Context, id string) error {
	return m.Called(id).Error(0)
}

func TestService_SignDocument(t *testing.T) {
	mockRepo := new(MockRepo)
	mockDocs := new(MockDocs)
	svc := NewService(mockRepo, mockDocs)

	t.Run("Success", func(t *testing.T) {
		req := SignRequest{DocumentID: "doc-1", Pin: "1234"}

		mockRepo.On("Create", mock.AnythingOfType("*signatures.Signature")).Return(nil)
		mockDocs.On("LockDocument", "doc-1").Return(nil)

		sig, err := svc.SignDocument("user-1", req)
		assert.NoError(t, err)
		assert.NotEmpty(t, sig.SignatureHash)
		assert.Equal(t, "doc-1", sig.DocumentID)
	})

	t.Run("Invalid PIN", func(t *testing.T) {
		req := SignRequest{DocumentID: "doc-1", Pin: "0000"}
		_, err := svc.SignDocument("user-1", req)
		assert.Error(t, err)
		assert.Equal(t, "invalid pin", err.Error())
	})
}
