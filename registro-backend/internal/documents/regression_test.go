package documents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRegRepo mocks - renamed to avoid conflict
type MockRegRepo struct {
	mock.Mock
}

func (m *MockRegRepo) Create(doc *Document, content string) error {
	return m.Called(doc, content).Error(0)
}
func (m *MockRegRepo) FindByID(id string) (*Document, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Document), args.Error(1)
}
func (m *MockRegRepo) Update(doc *Document, content, changelog string) error {
	return m.Called(doc, content, changelog).Error(0)
}
func (m *MockRegRepo) UpdateStatus(id string, status DocStatus) error {
	return m.Called(id, status).Error(0)
}
func (m *MockRegRepo) GetTemplate(id string) (*DocumentTemplate, error)         { return nil, nil }
func (m *MockRegRepo) GetContent(id string, version int) (string, error)        { return "content", nil }
func (m *MockRegRepo) GetVersions(id string) ([]DocumentVersion, error)         { return nil, nil }
func (m *MockRegRepo) GetSignatures(id string) ([]DocumentSignature, error)     { return nil, nil }
func (m *MockRegRepo) AddSignature(sig *DocumentSignature) error                { return nil }
func (m *MockRegRepo) GetInbox(schoolID string) ([]Document, error)             { return nil, nil }
func (m *MockRegRepo) GetReviewQueue(schoolID string) ([]Document, error)       { return nil, nil }
func (m *MockRegRepo) CreateTemplate(tpl *DocumentTemplate) error               { return nil }
func (m *MockRegRepo) FindByClass(classID string) ([]Document, error)           { return nil, nil }
func (m *MockRegRepo) FindByStudent(studentID string) ([]Document, error)       { return nil, nil }
func (m *MockRegRepo) GetTemplates(schoolID string) ([]DocumentTemplate, error) { return nil, nil }

func TestVideoRegression_UpdatesOnApprovedDoc(t *testing.T) {
	repo := &MockRegRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	// Scenario: Trying to update an Approved document
	doc := &Document{
		ID:     "doc-123",
		Status: StatusApproved,
		Title:  "Approved Doc",
	}

	repo.On("FindByID", "doc-123").Return(doc, nil)

	err := svc.UpdateDocument(ctx, "user-1", "doc-123", UpdateDocumentRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot edit non-draft document")
}

func TestVideoRegression_IllegalWorkflowTransition(t *testing.T) {
	repo := &MockRegRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	// Scenario: Trying to approve a Draft document directly (skipping Submit)
	doc := &Document{
		ID:     "doc-draft",
		Status: StatusDraft,
	}

	repo.On("FindByID", "doc-draft").Return(doc, nil)

	err := svc.ProcessWorkflow(ctx, "director", "doc-draft", WorkflowActionRequest{Action: "approve_director"})
	assert.Error(t, err)
}
