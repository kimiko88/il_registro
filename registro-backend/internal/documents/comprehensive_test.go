package documents

import (
	"context"
	"testing"

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
func (m *MockRepo) FindByClass(c string) ([]Document, error)           { return nil, nil }
func (m *MockRepo) FindByStudent(s string) ([]Document, error)         { return nil, nil }
func (m *MockRepo) GetInbox(s string) ([]Document, error)              { return nil, nil }
func (m *MockRepo) GetReviewQueue(s string) ([]Document, error)        { return nil, nil }

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

	res, err := svc.CreateDocument(context.Background(), "user1", req)
	assert.NoError(t, err)
	assert.Equal(t, "new-id", res.ID)

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
	doc := &Document{ID: docID, Status: StatusApproved, CurrentVersion: 1}

	mockRepo.On("FindByID", docID).Return(doc, nil)
	mockRepo.On("GetVersions", docID).Return([]DocumentVersion{{ID: "v1"}}, nil)
	mockRepo.On("AddSignature", mock.Anything).Return(nil)

	req := SignDocumentRequest{CertificateData: "cert", SignatureData: "sig"}
	err := svc.SignDocument(context.Background(), "user1", docID, req)
	assert.NoError(t, err)
}
