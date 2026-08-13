package documents

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(d *Document, c string) error {
	args := m.Called(d, c)
	d.ID = "new-id"
	return args.Error(0)
}
func (m *MockRepo) FindByID(id string) (*Document, error) {
	args := m.Called(id)
	return args.Get(0).(*Document), args.Error(1)
}
func (m *MockRepo) Update(d *Document, c, l string) error {
	args := m.Called(d, c, l)
	return args.Error(0)
}
func (m *MockRepo) UpdateStatus(id string, s DocStatus) error {
	args := m.Called(id, s)
	return args.Error(0)
}
func (m *MockRepo) GetContent(id string, v int) (string, error) {
	args := m.Called(id, v)
	return args.String(0), args.Error(1)
}
func (m *MockRepo) GetVersions(id string) ([]DocumentVersion, error) {
	args := m.Called(id)
	return args.Get(0).([]DocumentVersion), args.Error(1)
}
func (m *MockRepo) AddSignature(s *DocumentSignature) error {
	args := m.Called(s)
	return args.Error(0)
}
func (m *MockRepo) GetSignatures(id string) ([]DocumentSignature, error) {
	args := m.Called(id)
	return args.Get(0).([]DocumentSignature), args.Error(1)
}

// Stubs for others
func (m *MockRepo) GetTemplate(id string) (*DocumentTemplate, error)   { return nil, nil }
func (m *MockRepo) GetTemplates(id string) ([]DocumentTemplate, error) { return nil, nil }
func (m *MockRepo) CreateTemplate(t *DocumentTemplate) error           { return nil }
func (m *MockRepo) UpdateTemplate(t *DocumentTemplate) error           { return nil }
func (m *MockRepo) DeleteTemplate(id string) error                     { return nil }
func (m *MockRepo) FindByClass(c string) ([]Document, error)           { return nil, nil }
func (m *MockRepo) FindByStudent(s string) ([]Document, error)         { return nil, nil }
func (m *MockRepo) GetInbox(s string) ([]Document, error) {
	args := m.Called(s)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Document), args.Error(1)
}
func (m *MockRepo) GetReviewQueue(s string) ([]Document, error) { return nil, nil }
func (m *MockRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockRepo) ListAll(s string, t *DocType) ([]Document, error) {
	args := m.Called(s, t)
	return args.Get(0).([]Document), args.Error(1)
}

// --- Tests ---

func TestService_CreateDocument(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewService(mockRepo)

	req := CreateDocumentRequest{
		Title:   "Test PDP",
		Type:    TypePDP,
		Content: "Some content",
	}

	mockRepo.On("Create", mock.Anything, "Some content").Return(nil)

	res, err := svc.CreateDocument(context.Background(), "secretary", "user1", "school1", req)
	assert.NoError(t, err)
	if err == nil {
		assert.NotNil(t, res)
	}

	mockRepo.AssertExpectations(t)
}

func TestWorkflow_Transitions(t *testing.T) {
	w := NewWorkflowEngine()
	assert.NoError(t, w.CanTransition(StatusDraft, StatusSubmitted))
	assert.Error(t, w.CanTransition(StatusDraft, StatusSigned))
}

func TestService_SignDocument(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewService(mockRepo)

	docID := "doc-1"
	doc := &Document{ID: docID, SchoolID: "school-1", Status: StatusApproved, CurrentVersion: 1}

	mockRepo.On("FindByID", docID).Return(doc, nil)
	mockRepo.On("GetContent", docID, 1).Return("sample document content", nil)
	mockRepo.On("GetVersions", docID).Return([]DocumentVersion{{ID: "v1"}}, nil)
	mockRepo.On("AddSignature", mock.Anything).Return(nil)

	// Generate a valid RSA key and self-signed certificate for real crypto verification test
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	tmpl := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "Test Signer"},
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	assert.NoError(t, err)

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})

	digest := sha256.Sum256([]byte("sample document content"))
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, digest[:])
	assert.NoError(t, err)

	req := SignDocumentRequest{CertificateData: string(pemCert), SignatureData: base64.StdEncoding.EncodeToString(sig)}
	err = svc.SignDocument(context.Background(), "admin", "school-1", "user1", docID, req)
	assert.NoError(t, err)
}

func TestService_GetInbox(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewService(mockRepo)

	docs := []Document{{ID: "doc1", Title: "Inbox Doc", Status: StatusSubmitted}}
	mockRepo.On("GetInbox", "school1").Return(docs, nil)

	res, err := svc.GetInbox(context.Background(), "school1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "Inbox Doc", res[0].Title)
	mockRepo.AssertExpectations(t)
}
