package signatures

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/users"
)

type MockRepo struct{ mock.Mock }

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

type MockDocs struct{ mock.Mock }

func (m *MockDocs) LockDocument(ctx context.Context, id string) error { return m.Called(id).Error(0) }
func (m *MockDocs) GetDocumentContent(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func TestSignDocument_NoMFA_Rejected(t *testing.T) {
	for _, pin := range []string{"1234", "000000", "correct"} {
		mockRepo := new(MockRepo)
		mockDocs := new(MockDocs)
		mockUsers := new(MockUserRepo)
		svc := NewService(mockRepo, mockDocs, mockUsers)
		user := &users.User{ID: "u1", MFAEnabled: false}
		mockUsers.On("GetByID", mock.Anything, "u1").Return(user, nil).Once()
		_, err := svc.SignDocument("u1", SignRequest{DocumentID: "doc-1", Pin: pin})
		assert.Error(t, err, "pin=%s deve essere rifiutato", pin)
		assert.Contains(t, err.Error(), "MFA")
		mockUsers.AssertExpectations(t)
	}
}

func TestSignDocument_RealContentHash(t *testing.T) {
	mockRepo := new(MockRepo)
	mockDocs := new(MockDocs)
	mockUsers := new(MockUserRepo)
	svc := NewService(mockRepo, mockDocs, mockUsers)
	secret, _ := totp.Generate(totp.GenerateOpts{Issuer: "T", AccountName: "t@t.com"})
	totpCode, _ := totp.GenerateCode(secret.Secret(), time.Now())
	const realContent = "Contenuto reale 2025"
	h := sha256.Sum256([]byte(realContent))
	expectedHash := hex.EncodeToString(h[:])
	user := &users.User{ID: "u2", MFAEnabled: true, MFASecret: secret.Secret()}
	mockUsers.On("GetByID", mock.Anything, "u2").Return(user, nil).Once()
	mockDocs.On("GetDocumentContent", mock.Anything, "doc-r").Return(realContent, nil).Once()
	mockRepo.On("Create", mock.AnythingOfType("*signatures.Signature")).Return(nil).Once()
	mockDocs.On("LockDocument", "doc-r").Return(nil).Once()
	sig, err := svc.SignDocument("u2", SignRequest{DocumentID: "doc-r", Pin: totpCode})
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, sig.SignatureHash)
	assert.Len(t, sig.SignatureHash, 64)
	mockUsers.AssertExpectations(t)
	mockDocs.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestSignDocument_EmptyContent_Rejected(t *testing.T) {
	mockRepo := new(MockRepo)
	mockDocs := new(MockDocs)
	mockUsers := new(MockUserRepo)
	svc := NewService(mockRepo, mockDocs, mockUsers)
	secret, _ := totp.Generate(totp.GenerateOpts{Issuer: "T", AccountName: "t@t.com"})
	totpCode, _ := totp.GenerateCode(secret.Secret(), time.Now())
	user := &users.User{ID: "u3", MFAEnabled: true, MFASecret: secret.Secret()}
	mockUsers.On("GetByID", mock.Anything, "u3").Return(user, nil).Once()
	mockDocs.On("GetDocumentContent", mock.Anything, "doc-empty").Return("", nil).Once()
	_, err := svc.SignDocument("u3", SignRequest{DocumentID: "doc-empty", Pin: totpCode})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contenuto")
}

func TestSignDocument_InvalidTOTP_Rejected(t *testing.T) {
	mockRepo := new(MockRepo)
	mockDocs := new(MockDocs)
	mockUsers := new(MockUserRepo)
	svc := NewService(mockRepo, mockDocs, mockUsers)
	user := &users.User{ID: "u5", MFAEnabled: true, MFASecret: "ABCDEFGHIJKLMNOP"}
	mockUsers.On("GetByID", mock.Anything, "u5").Return(user, nil).Once()
	_, err := svc.SignDocument("u5", SignRequest{DocumentID: "doc-1", Pin: "000000"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OTP non valido")
	mockUsers.AssertExpectations(t)
}
