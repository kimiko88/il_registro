package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/documents"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockDocRepoForBES struct {
	docs map[string]*documents.Document
}

func (m *mockDocRepoForBES) Create(doc *documents.Document, initialContent string) error {
	if doc.ID == "" {
		doc.ID = "doc-pdp-1"
	}
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()
	m.docs[doc.ID] = doc
	return nil
}

func (m *mockDocRepoForBES) Update(doc *documents.Document, newContent, changeLog string) error {
	m.docs[doc.ID] = doc
	return nil
}

func (m *mockDocRepoForBES) UpdateStatus(docID string, status documents.DocStatus) error {
	if doc, ok := m.docs[docID]; ok {
		doc.Status = status
	}
	return nil
}

func (m *mockDocRepoForBES) Delete(docID string) error {
	delete(m.docs, docID)
	return nil
}

func (m *mockDocRepoForBES) ListAll(schoolID string, docType *documents.DocType) ([]documents.Document, error) {
	res := make([]documents.Document, 0)
	for _, d := range m.docs {
		res = append(res, *d)
	}
	return res, nil
}

func (m *mockDocRepoForBES) FindByID(id string) (*documents.Document, error) {
	d, ok := m.docs[id]
	if !ok {
		return nil, errors.New("document not found")
	}
	return d, nil
}

func (m *mockDocRepoForBES) GetContent(docID string, version int) (string, error) {
	return "{\"measures\": [\"Tempo aggiuntivo 30%\", \"Calcolatrice\"]}", nil
}

func (m *mockDocRepoForBES) GetVersions(docID string) ([]documents.DocumentVersion, error) {
	return []documents.DocumentVersion{}, nil
}

func (m *mockDocRepoForBES) FindByClass(classID string) ([]documents.Document, error) {
	return []documents.Document{}, nil
}

func (m *mockDocRepoForBES) FindByStudent(studentID string) ([]documents.Document, error) {
	return []documents.Document{}, nil
}

func (m *mockDocRepoForBES) GetInbox(schoolID string) ([]documents.Document, error) {
	return []documents.Document{}, nil
}

func (m *mockDocRepoForBES) GetReviewQueue(schoolID string) ([]documents.Document, error) {
	return []documents.Document{}, nil
}

func (m *mockDocRepoForBES) AddSignature(sig *documents.DocumentSignature) error {
	return nil
}

func (m *mockDocRepoForBES) GetSignatures(docID string) ([]documents.DocumentSignature, error) {
	return []documents.DocumentSignature{}, nil
}

func (m *mockDocRepoForBES) GetTemplates(schoolID string) ([]documents.DocumentTemplate, error) {
	return []documents.DocumentTemplate{}, nil
}

func (m *mockDocRepoForBES) GetTemplate(id string) (*documents.DocumentTemplate, error) {
	return nil, nil
}

func (m *mockDocRepoForBES) CreateTemplate(tpl *documents.DocumentTemplate) error {
	return nil
}

func (m *mockDocRepoForBES) UpdateTemplate(tpl *documents.DocumentTemplate) error {
	return nil
}

func (m *mockDocRepoForBES) DeleteTemplate(id string) error {
	return nil
}

func TestIntegration_BES_PDP_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockDocRepoForBES{docs: make(map[string]*documents.Document)}
	svc := documents.NewService(repo, nil)
	handler := documents.NewHandler(svc)

	r := gin.New()
	r.POST("/documents", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.CreateDocument(c)
	})
	r.POST("/documents/:id/workflow", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.ProcessWorkflow(c)
	})

	// 1. Create PDP draft document for BES/DSA student
	studID := "00000000-0000-0000-0000-000000000001"
	classID := "00000000-0000-0000-0000-000000000002"
	createReq := documents.CreateDocumentRequest{
		Title:     "Piano Didattico Personalizzato (PDP) - Mario Rossi",
		Type:      documents.TypePDP,
		StudentID: &studID,
		ClassID:   &classID,
		Content:   "{\"measures\": [\"Tempo aggiuntivo 30%\", \"Calcolatrice\"]}",
	}
	body, _ := json.Marshal(createReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	t.Logf("CreateDoc response status %d, body: %s", w1.Code, w1.Body.String())
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Submit PDP document for review
	w2 := httptest.NewRecorder()
	wfReq := documents.WorkflowActionRequest{
		Action: "submit",
	}
	body2, _ := json.Marshal(wfReq)
	req2, _ := http.NewRequest("POST", "/documents/doc-pdp-1/workflow", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
