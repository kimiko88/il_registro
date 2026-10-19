package documents

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDocumentService struct {
	mock.Mock
}

func (m *MockDocumentService) CreateDocument(ctx context.Context, actorRole, userID, schoolID string, req CreateDocumentRequest) (*DocumentListResponse, error) {
	args := m.Called(ctx, actorRole, userID, schoolID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DocumentListResponse), args.Error(1)
}

func (m *MockDocumentService) GetDocument(ctx context.Context, actorRole, schoolID, id string) (*DocumentDetailResponse, error) {
	args := m.Called(ctx, actorRole, schoolID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DocumentDetailResponse), args.Error(1)
}

func (m *MockDocumentService) UpdateDocument(ctx context.Context, actorRole, schoolID, userID, id string, req UpdateDocumentRequest) error {
	args := m.Called(ctx, actorRole, schoolID, userID, id, req)
	return args.Error(0)
}

func (m *MockDocumentService) DeleteDocument(ctx context.Context, actorRole, schoolID, userID, id string) error {
	args := m.Called(ctx, actorRole, schoolID, userID, id)
	return args.Error(0)
}

func (m *MockDocumentService) AttachFile(ctx context.Context, actorRole, schoolID, docID, fileURL string) error {
	args := m.Called(ctx, actorRole, schoolID, docID, fileURL)
	return args.Error(0)
}

func (m *MockDocumentService) ProcessWorkflow(ctx context.Context, actorRole, schoolID, userID, docID string, req WorkflowActionRequest) error {
	args := m.Called(ctx, actorRole, schoolID, userID, docID, req)
	return args.Error(0)
}

func (m *MockDocumentService) SignDocument(ctx context.Context, actorRole, schoolID, userID, docID string, req SignDocumentRequest) error {
	args := m.Called(ctx, actorRole, schoolID, userID, docID, req)
	return args.Error(0)
}

func (m *MockDocumentService) GetInbox(ctx context.Context, schoolID string) ([]DocumentListResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentListResponse), args.Error(1)
}

func (m *MockDocumentService) GetReviewQueue(ctx context.Context, schoolID string) ([]DocumentListResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentListResponse), args.Error(1)
}

func (m *MockDocumentService) GetMyDocuments(ctx context.Context, schoolID, userID string) ([]DocumentListResponse, error) {
	args := m.Called(ctx, schoolID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentListResponse), args.Error(1)
}

func (m *MockDocumentService) ListDocuments(ctx context.Context, actorRole, schoolID string, docType *DocType) ([]DocumentListResponse, error) {
	args := m.Called(ctx, actorRole, schoolID, docType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentListResponse), args.Error(1)
}

func (m *MockDocumentService) CreateTemplate(ctx context.Context, schoolID string, req TemplateRequest) error {
	args := m.Called(ctx, schoolID, req)
	return args.Error(0)
}

func (m *MockDocumentService) ListTemplates(ctx context.Context, schoolID string) ([]DocumentTemplate, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentTemplate), args.Error(1)
}

func (m *MockDocumentService) UpdateTemplate(ctx context.Context, schoolID, id string, req TemplateRequest) error {
	args := m.Called(ctx, schoolID, id, req)
	return args.Error(0)
}

func (m *MockDocumentService) DeleteTemplate(ctx context.Context, schoolID, id string) error {
	args := m.Called(ctx, schoolID, id)
	return args.Error(0)
}

func (m *MockDocumentService) ExportDocument(ctx context.Context, actorRole, schoolID, id, format string) ([]byte, string, error) {
	args := m.Called(ctx, actorRole, schoolID, id, format)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).([]byte), args.String(1), args.Error(2)
}

func (m *MockDocumentService) GetDocumentContent(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *MockDocumentService) LockDocument(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentService) GetDocumentVersions(ctx context.Context, actorRole, schoolID, docID string) ([]DocumentVersion, error) {
	args := m.Called(ctx, actorRole, schoolID, docID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentVersion), args.Error(1)
}

type MockStorageUploader struct {
	mock.Mock
}

func (u *MockStorageUploader) UploadFile(ctx context.Context, storagePath, mimeType string, r io.Reader) (string, error) {
	args := u.Called(ctx, storagePath, mimeType, r)
	return args.String(0), args.Error(1)
}

func (u *MockStorageUploader) DeleteFile(ctx context.Context, storagePath string) error {
	args := u.Called(ctx, storagePath)
	return args.Error(0)
}

func setupTestDocumentRouter(svc Service, upl *MockStorageUploader, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		if role != "" {
			c.Set("role", role)
		}
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	var h *Handler
	if upl != nil {
		h = NewHandler(svc, upl)
	} else {
		h = NewHandler(svc)
	}
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestDocumentsHandler_AllEndpoints(t *testing.T) {
	svc := new(MockDocumentService)
	upl := new(MockStorageUploader)
	rTeacher := setupTestDocumentRouter(svc, upl, "teacher-1", "teacher", "sch-1")

	// 1. POST /documents - Create Document
	svc.On("CreateDocument", mock.Anything, "teacher", "teacher-1", "sch-1", mock.AnythingOfType("documents.CreateDocumentRequest")).
		Return(&DocumentListResponse{ID: "doc-1", Title: "Verbale Consiglio"}, nil).Once()

	createPayload := `{"title":"Verbale Consiglio","type":"report","content":"Contenuto verbale"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/documents", bytes.NewBufferString(createPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. GET /documents - List Documents
	svc.On("ListDocuments", mock.Anything, "teacher", "sch-1", (*DocType)(nil)).
		Return([]DocumentListResponse{{ID: "doc-1", Title: "Verbale Consiglio"}}, nil).Once()

	req, _ = http.NewRequest(http.MethodGet, "/api/v1/documents", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. GET /documents/:id - Get Document
	now := time.Now()
	svc.On("GetDocument", mock.Anything, "teacher", "sch-1", "doc-1").
		Return(&DocumentDetailResponse{
			DocumentListResponse: DocumentListResponse{ID: "doc-1", Title: "Verbale Consiglio", Status: StatusDraft, UpdatedAt: now},
			Content:              "Contenuto verbale",
			CurrentVersion:       1,
		}, nil).Once()

	req, _ = http.NewRequest(http.MethodGet, "/api/v1/documents/doc-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. GET /documents/:id/versions
	svc.On("GetDocumentVersions", mock.Anything, "teacher", "sch-1", "doc-1").
		Return([]DocumentVersion{{ID: "v1", VersionNum: 1, DocumentID: "doc-1"}}, nil).Once()

	req, _ = http.NewRequest(http.MethodGet, "/api/v1/documents/doc-1/versions", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. PATCH /documents/:id - Update Document
	svc.On("UpdateDocument", mock.Anything, "teacher", "sch-1", "teacher-1", "doc-1", mock.AnythingOfType("documents.UpdateDocumentRequest")).
		Return(nil).Once()

	updatePayload := `{"title":"Verbale Aggiornato","change_log":"Corretto errore"}`
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/documents/doc-1", bytes.NewBufferString(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. POST /documents/:id/workflow - Process workflow
	svc.On("ProcessWorkflow", mock.Anything, "teacher", "sch-1", "teacher-1", "doc-1", WorkflowActionRequest{Action: "submit"}).
		Return(nil).Once()

	wfPayload := `{"action":"submit"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/documents/doc-1/workflow", bytes.NewBufferString(wfPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. GET /documents/:id/export - Export PDF
	svc.On("ExportDocument", mock.Anything, "teacher", "sch-1", "doc-1", "pdf").
		Return([]byte("%PDF-1.4 test"), "application/pdf", nil).Once()

	req, _ = http.NewRequest(http.MethodGet, "/api/v1/documents/doc-1/export?format=pdf", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))

	// 8. DELETE /documents/:id - Delete Document
	svc.On("DeleteDocument", mock.Anything, "teacher", "sch-1", "teacher-1", "doc-1").
		Return(nil).Once()

	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/documents/doc-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 9. Secretary & Director endpoints
	rSecretary := setupTestDocumentRouter(svc, upl, "sec-1", "secretary", "sch-1")
	svc.On("GetInbox", mock.Anything, "sch-1").
		Return([]DocumentListResponse{{ID: "doc-1", Status: StatusSubmitted}}, nil).Once()

	req, _ = http.NewRequest(http.MethodGet, "/api/v1/documents/inbox", nil)
	w = httptest.NewRecorder()
	rSecretary.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	rPrincipal := setupTestDocumentRouter(svc, upl, "p-1", "principal", "sch-1")
	svc.On("GetReviewQueue", mock.Anything, "sch-1").
		Return([]DocumentListResponse{{ID: "doc-1", Status: StatusReview}}, nil).Once()

	req, _ = http.NewRequest(http.MethodGet, "/api/v1/documents/review-queue", nil)
	w = httptest.NewRecorder()
	rPrincipal.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 10. POST /documents/:id/sign - Sign document
	svc.On("SignDocument", mock.Anything, "principal", "sch-1", "p-1", "doc-1", mock.AnythingOfType("documents.SignDocumentRequest")).
		Return(nil).Once()

	signPayload := `{"signature_data":"base64sig==","certificate_data":"cert=="}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/documents/doc-1/sign", bytes.NewBufferString(signPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rPrincipal.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 11. Templates endpoints
	rAdmin := setupTestDocumentRouter(svc, upl, "adm-1", "admin", "sch-1")

	// CreateTemplate
	svc.On("CreateTemplate", mock.Anything, "sch-1", mock.AnythingOfType("documents.TemplateRequest")).
		Return(nil).Once()
	tplPayload := `{"name":"Modello PDP","type":"pdp","content":"Struttura modello"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/documents/template", bytes.NewBufferString(tplPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// ListTemplates
	svc.On("ListTemplates", mock.Anything, "sch-1").
		Return([]DocumentTemplate{{ID: "tpl-1", Name: "Modello PDP"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/documents/template", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// UpdateTemplate
	svc.On("UpdateTemplate", mock.Anything, "sch-1", "tpl-1", mock.AnythingOfType("documents.TemplateRequest")).
		Return(nil).Once()
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/documents/template/tpl-1", bytes.NewBufferString(tplPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// DeleteTemplate
	svc.On("DeleteTemplate", mock.Anything, "sch-1", "tpl-1").
		Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/documents/template/tpl-1", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 12. UploadFile endpoint
	// Test file with invalid extension
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "malicious.exe")
	_, _ = part.Write([]byte("MZ executable header dummy"))
	_ = writer.Close()

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/documents/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusUnprocessableEntity)

	// Test missing file
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/documents/upload", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=something")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 13. Unauthorized error handling checks
	svc.On("GetDocument", mock.Anything, "teacher", "sch-1", "unauthorized-doc").
		Return(nil, errors.New("unauthorized")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/documents/unauthorized-doc", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	svc.AssertExpectations(t)
}
