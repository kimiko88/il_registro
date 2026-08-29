package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"registro-backend/internal/documents"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Documents Repository ─────────────────────────────────────────────

type mockDocsRepo struct {
	docs      map[string]*documents.Document
	contents  map[string]map[int]string // docID -> version -> content
	versions  map[string][]documents.DocumentVersion
	templates map[string]*documents.DocumentTemplate
	sigs      map[string][]documents.DocumentSignature
}

func newMockDocsRepo() *mockDocsRepo {
	return &mockDocsRepo{
		docs:      make(map[string]*documents.Document),
		contents:  make(map[string]map[int]string),
		versions:  make(map[string][]documents.DocumentVersion),
		templates: make(map[string]*documents.DocumentTemplate),
		sigs:      make(map[string][]documents.DocumentSignature),
	}
}

func (m *mockDocsRepo) Create(doc *documents.Document, initialContent string) error {
	doc.ID = "doc-" + strings.ReplaceAll(doc.Title, " ", "_")
	doc.CurrentVersion = 1
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()
	m.docs[doc.ID] = doc

	if m.contents[doc.ID] == nil {
		m.contents[doc.ID] = make(map[int]string)
	}
	m.contents[doc.ID][1] = initialContent
	m.versions[doc.ID] = []documents.DocumentVersion{
		{
			ID:         "v1-" + doc.ID,
			DocumentID: doc.ID,
			VersionNum: 1,
			Content:    initialContent,
			ChangeLog:  "Initial creation",
			CreatedBy:  doc.CreatedBy,
			CreatedAt:  time.Now(),
		},
	}
	return nil
}

func (m *mockDocsRepo) Update(doc *documents.Document, newContent, changeLog string) error {
	existing, ok := m.docs[doc.ID]
	if !ok {
		return errors.New("document not found")
	}
	existing.CurrentVersion++
	existing.UpdatedAt = time.Now()
	if doc.Title != "" {
		existing.Title = doc.Title
	}
	m.contents[doc.ID][existing.CurrentVersion] = newContent
	m.versions[doc.ID] = append(m.versions[doc.ID], documents.DocumentVersion{
		ID:         "v-" + doc.ID,
		DocumentID: doc.ID,
		VersionNum: existing.CurrentVersion,
		Content:    newContent,
		ChangeLog:  changeLog,
		CreatedBy:  doc.CreatedBy,
		CreatedAt:  time.Now(),
	})
	return nil
}

func (m *mockDocsRepo) UpdateStatus(docID string, status documents.DocStatus) error {
	doc, ok := m.docs[docID]
	if !ok {
		return errors.New("document not found")
	}
	doc.Status = status
	doc.UpdatedAt = time.Now()
	return nil
}

func (m *mockDocsRepo) Delete(docID string) error {
	now := time.Now()
	if doc, ok := m.docs[docID]; ok {
		doc.DeletedAt = &now
		return nil
	}
	return errors.New("document not found")
}

func (m *mockDocsRepo) ListAll(schoolID string, docType *documents.DocType) ([]documents.Document, error) {
	var list []documents.Document
	for _, d := range m.docs {
		if d.SchoolID == schoolID && d.DeletedAt == nil {
			if docType == nil || d.Type == *docType {
				list = append(list, *d)
			}
		}
	}
	return list, nil
}

func (m *mockDocsRepo) FindByID(id string) (*documents.Document, error) {
	doc, ok := m.docs[id]
	if !ok || doc.DeletedAt != nil {
		return nil, errors.New("document not found")
	}
	return doc, nil
}

func (m *mockDocsRepo) GetContent(docID string, version int) (string, error) {
	cMap, ok := m.contents[docID]
	if !ok {
		return "", errors.New("document content not found")
	}
	content, ok := cMap[version]
	if !ok {
		return "", errors.New("version not found")
	}
	return content, nil
}

func (m *mockDocsRepo) GetVersions(docID string) ([]documents.DocumentVersion, error) {
	v, ok := m.versions[docID]
	if !ok {
		return []documents.DocumentVersion{}, nil
	}
	return v, nil
}

func (m *mockDocsRepo) FindByClass(classID string) ([]documents.Document, error) {
	var out []documents.Document
	for _, d := range m.docs {
		if d.ClassID != nil && *d.ClassID == classID && d.DeletedAt == nil {
			out = append(out, *d)
		}
	}
	return out, nil
}

func (m *mockDocsRepo) FindByStudent(studentID string) ([]documents.Document, error) {
	var out []documents.Document
	for _, d := range m.docs {
		if d.StudentID != nil && *d.StudentID == studentID && d.DeletedAt == nil {
			out = append(out, *d)
		}
	}
	return out, nil
}

func (m *mockDocsRepo) GetInbox(schoolID string) ([]documents.Document, error) {
	var out []documents.Document
	for _, d := range m.docs {
		if d.SchoolID == schoolID && d.Status == documents.StatusSubmitted && d.DeletedAt == nil {
			out = append(out, *d)
		}
	}
	return out, nil
}

func (m *mockDocsRepo) GetReviewQueue(schoolID string) ([]documents.Document, error) {
	var out []documents.Document
	for _, d := range m.docs {
		if d.SchoolID == schoolID && d.Status == documents.StatusReview && d.DeletedAt == nil {
			out = append(out, *d)
		}
	}
	return out, nil
}

func (m *mockDocsRepo) AddSignature(sig *documents.DocumentSignature) error {
	m.sigs[sig.DocumentID] = append(m.sigs[sig.DocumentID], *sig)
	return nil
}

func (m *mockDocsRepo) GetSignatures(docID string) ([]documents.DocumentSignature, error) {
	return m.sigs[docID], nil
}

func (m *mockDocsRepo) GetTemplates(schoolID string) ([]documents.DocumentTemplate, error) {
	var out []documents.DocumentTemplate
	for _, t := range m.templates {
		if t.SchoolID == schoolID && t.IsActive {
			out = append(out, *t)
		}
	}
	return out, nil
}

func (m *mockDocsRepo) GetTemplate(id string) (*documents.DocumentTemplate, error) {
	t, ok := m.templates[id]
	if !ok {
		return nil, errors.New("template not found")
	}
	return t, nil
}

func (m *mockDocsRepo) CreateTemplate(tpl *documents.DocumentTemplate) error {
	tpl.ID = "tpl-" + strings.ReplaceAll(tpl.Name, " ", "_")
	tpl.IsActive = true
	m.templates[tpl.ID] = tpl
	return nil
}

func (m *mockDocsRepo) UpdateTemplate(tpl *documents.DocumentTemplate) error {
	m.templates[tpl.ID] = tpl
	return nil
}

func (m *mockDocsRepo) DeleteTemplate(id string) error {
	delete(m.templates, id)
	return nil
}

// ── Setup Router Helper ───────────────────────────────────────────────────

func setupDocumentsRouter(repo *mockDocsRepo, role, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := documents.NewService(repo)
	h := documents.NewHandler(svc)
	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Set("school_id", schoolID)
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

// ── Tests ─────────────────────────────────────────────────────────────────

// DOC01 — Admin can create a new document
func TestDocuments_CreateDocument_Admin(t *testing.T) {
	repo := newMockDocsRepo()
	r := setupDocumentsRouter(repo, "admin", "admin-1", "school-1")

	body, _ := json.Marshal(documents.CreateDocumentRequest{
		Title:   "Piano Didattico Personalizzato",
		Type:    documents.TypePDP,
		Content: "<p>Contenuto PDP per studente BES</p>",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/documents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var resp documents.DocumentListResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, "Piano Didattico Personalizzato", resp.Title)
	assert.Equal(t, documents.StatusDraft, resp.Status)
}

// DOC02 — Student is forbidden from creating documents
func TestDocuments_CreateDocument_ForbiddenForStudent(t *testing.T) {
	repo := newMockDocsRepo()
	r := setupDocumentsRouter(repo, "student", "student-1", "school-1")

	body, _ := json.Marshal(documents.CreateDocumentRequest{
		Title:   "Documento Non Autorizzato",
		Type:    documents.TypeGeneric,
		Content: "test",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/documents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// DOC03 — List documents for school
func TestDocuments_ListDocuments(t *testing.T) {
	repo := newMockDocsRepo()
	_ = repo.Create(&documents.Document{
		SchoolID:  "school-1",
		Title:     "Doc 1",
		Type:      documents.TypePDP,
		Status:    documents.StatusDraft,
		CreatedBy: "admin-1",
	}, "Content 1")
	_ = repo.Create(&documents.Document{
		SchoolID:  "school-1",
		Title:     "Doc 2",
		Type:      documents.TypePCTO,
		Status:    documents.StatusDraft,
		CreatedBy: "admin-1",
	}, "Content 2")

	r := setupDocumentsRouter(repo, "admin", "admin-1", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp []documents.DocumentListResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
}

// DOC04 — Get document details including content and versions
func TestDocuments_GetDocument_Detail(t *testing.T) {
	repo := newMockDocsRepo()
	doc := &documents.Document{
		SchoolID:  "school-1",
		Title:     "Relazione Finale 5A",
		Type:      documents.TypeMay15,
		Status:    documents.StatusDraft,
		CreatedBy: "admin-1",
	}
	_ = repo.Create(doc, "<h1>Relazione 15 Maggio</h1>")

	r := setupDocumentsRouter(repo, "admin", "admin-1", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents/"+doc.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp documents.DocumentDetailResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "Relazione Finale 5A", resp.Title)
	assert.Contains(t, resp.Content, "Relazione 15 Maggio")
	assert.Equal(t, 1, resp.CurrentVersion)
}

// DOC05 — Update document content creates a new version
func TestDocuments_UpdateDocument_CreatesNewVersion(t *testing.T) {
	repo := newMockDocsRepo()
	doc := &documents.Document{
		SchoolID:  "school-1",
		Title:     "Bozza Documento",
		Type:      documents.TypeGeneric,
		Status:    documents.StatusDraft,
		CreatedBy: "admin-1",
	}
	_ = repo.Create(doc, "Versione 1")

	r := setupDocumentsRouter(repo, "admin", "admin-1", "school-1")
	newTitle := "Bozza Aggiornata"
	newContent := "Versione 2 con modifiche"
	body, _ := json.Marshal(documents.UpdateDocumentRequest{
		Title:     &newTitle,
		Content:   &newContent,
		ChangeLog: "Aggiunta sezione obiettivi",
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/documents/"+doc.ID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	// Verify version incremented
	assert.Equal(t, 2, doc.CurrentVersion)
}

// DOC06 — Workflow transitions: submit document for secretary review
func TestDocuments_Workflow_SubmitToInbox(t *testing.T) {
	repo := newMockDocsRepo()
	doc := &documents.Document{
		SchoolID:  "school-1",
		Title:     "PDP da inviare",
		Type:      documents.TypePDP,
		Status:    documents.StatusDraft,
		CreatedBy: "secretary-1",
	}
	_ = repo.Create(doc, "Contenuto")

	r := setupDocumentsRouter(repo, "secretary", "secretary-1", "school-1")
	body, _ := json.Marshal(documents.WorkflowActionRequest{
		Action: "submit",
		Reason: "Pronto per revisione",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/documents/"+doc.ID+"/workflow", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, documents.StatusSubmitted, doc.Status)
}

// DOC07 — Secretary retrieves inbox containing submitted documents
func TestDocuments_Secretary_GetInbox(t *testing.T) {
	repo := newMockDocsRepo()
	doc := &documents.Document{
		SchoolID:  "school-1",
		Title:     "PDP in attesa",
		Type:      documents.TypePDP,
		Status:    documents.StatusSubmitted,
		CreatedBy: "teacher-1",
	}
	_ = repo.Create(doc, "Contenuto")

	r := setupDocumentsRouter(repo, "secretary", "sec-1", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents/inbox", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp []documents.DocumentListResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "PDP in attesa", resp[0].Title)
}

// DOC08 — Document template management (Create, List, Delete)
func TestDocuments_Templates_Lifecycle(t *testing.T) {
	repo := newMockDocsRepo()
	r := setupDocumentsRouter(repo, "admin", "admin-1", "school-1")

	// 1. Create Template
	createBody, _ := json.Marshal(documents.TemplateRequest{
		Name:    "Modello PDP 2026",
		Type:    documents.TypePDP,
		Content: "<h1>PDP {{student_name}}</h1>",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/documents/template", bytes.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// 2. List Templates
	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/documents/template", nil)
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)
	require.Equal(t, http.StatusOK, wList.Code)
	var listResp []documents.DocumentTemplate
	err := json.Unmarshal(wList.Body.Bytes(), &listResp)
	require.NoError(t, err)
	assert.Len(t, listResp, 1)
	assert.Equal(t, "Modello PDP 2026", listResp[0].Name)

	// 3. Delete Template
	reqDel := httptest.NewRequest(http.MethodDelete, "/api/v1/documents/template/"+listResp[0].ID, nil)
	wDel := httptest.NewRecorder()
	r.ServeHTTP(wDel, reqDel)
	assert.Equal(t, http.StatusOK, wDel.Code)
}

// DOC09 — Unauthorized access without user_id returns 401
func TestDocuments_Unauthorized_NoUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockDocsRepo()
	r := gin.New()
	svc := documents.NewService(repo)
	h := documents.NewHandler(svc)
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
