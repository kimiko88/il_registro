package signatures

import (
	"context"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/users"
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

// MockUserRepo
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func TestService_SignDocument(t *testing.T) {
	mockRepo := new(MockRepo)
	mockDocs := new(MockDocs)
	mockUsers := new(MockUserRepo)
	svc := NewService(mockRepo, mockDocs, mockUsers)

	t.Run("Success without MFA (legacy PIN)", func(t *testing.T) {
		req := SignRequest{DocumentID: "doc-1", Pin: "1234"}
		user := &users.User{ID: "user-1", MFAEnabled: false}

		mockUsers.On("GetByID", mock.Anything, "user-1").Return(user, nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("*signatures.Signature")).Return(nil).Once()
		mockDocs.On("LockDocument", "doc-1").Return(nil).Once()

		sig, err := svc.SignDocument("user-1", req)
		assert.NoError(t, err)
		assert.NotEmpty(t, sig.SignatureHash)
		assert.Equal(t, "doc-1", sig.DocumentID)
		mockUsers.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockDocs.AssertExpectations(t)
	})

	t.Run("Success with MFA (TOTP)", func(t *testing.T) {
		secret, _ := totp.Generate(totp.GenerateOpts{
			Issuer:      "Test",
			AccountName: "test@example.com",
		})
		totpCode, _ := totp.GenerateCode(secret.Secret(), time.Now())
		req := SignRequest{DocumentID: "doc-2", Pin: totpCode}
		user := &users.User{ID: "user-2", MFAEnabled: true, MFASecret: secret.Secret()}

		mockUsers.On("GetByID", mock.Anything, "user-2").Return(user, nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("*signatures.Signature")).Return(nil).Once()
		mockDocs.On("LockDocument", "doc-2").Return(nil).Once()

		sig, err := svc.SignDocument("user-2", req)
		assert.NoError(t, err)
		assert.NotEmpty(t, sig.SignatureHash)
		mockUsers.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockDocs.AssertExpectations(t)
	})

	t.Run("Invalid legacy PIN", func(t *testing.T) {
		req := SignRequest{DocumentID: "doc-1", Pin: "0000"}
		user := &users.User{ID: "user-1", MFAEnabled: false}

		mockUsers.On("GetByID", mock.Anything, "user-1").Return(user, nil).Once()
		_, err := svc.SignDocument("user-1", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PIN non valido")
		mockUsers.AssertExpectations(t)
	})

	t.Run("Invalid TOTP", func(t *testing.T) {
		req := SignRequest{DocumentID: "doc-1", Pin: "000000"}
		user := &users.User{ID: "user-3", MFAEnabled: true, MFASecret: "ABCD"}

		mockUsers.On("GetByID", mock.Anything, "user-3").Return(user, nil).Once()
		_, err := svc.SignDocument("user-3", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "OTP non valido")
		mockUsers.AssertExpectations(t)
	})
}
