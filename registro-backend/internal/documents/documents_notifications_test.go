package documents

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttachFile_CrossSchool_Denied(t *testing.T) {
	mockRepo := &mockDocRepoForTesting{
		doc: &Document{
			ID:       "doc-1",
			SchoolID: "school-A",
		},
	}

	svc := NewService(mockRepo)

	err := svc.AttachFile(context.Background(), "admin", "school-B", "doc-1", "https://storage.example.com/file.pdf")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot access documents of another school")
}

type mockDocRepoForTesting struct {
	Repository
	doc *Document
}

func (m *mockDocRepoForTesting) FindByID(id string) (*Document, error) {
	if m.doc != nil && m.doc.ID == id {
		return m.doc, nil
	}
	return nil, errors.New("document not found")
}
