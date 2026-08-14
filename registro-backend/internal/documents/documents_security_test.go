package documents

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockServiceForSecTest struct {
	mock.Mock
}

func (m *mockServiceForSecTest) CreateDocument(ctx context.Context, actorRole, userID, schoolID string, req CreateDocumentRequest) (*DocumentListResponse, error) {
	args := m.Called(ctx, actorRole, userID, schoolID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DocumentListResponse), args.Error(1)
}
func (m *mockServiceForSecTest) GetDocument(ctx context.Context, actorRole, schoolID, id string) (*DocumentDetailResponse, error) {
	args := m.Called(ctx, actorRole, schoolID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DocumentDetailResponse), args.Error(1)
}
func (m *mockServiceForSecTest) UpdateDocument(ctx context.Context, actorRole, schoolID, userID, id string, req UpdateDocumentRequest) error {
	return m.Called(ctx, actorRole, schoolID, userID, id, req).Error(0)
}
func (m *mockServiceForSecTest) DeleteDocument(ctx context.Context, actorRole, schoolID, userID, id string) error {
	return m.Called(ctx, actorRole, schoolID, userID, id).Error(0)
}
func (m *mockServiceForSecTest) AttachFile(ctx context.Context, actorRole, schoolID, docID, fileURL string) error {
	return m.Called(ctx, actorRole, schoolID, docID, fileURL).Error(0)
}
func (m *mockServiceForSecTest) ProcessWorkflow(ctx context.Context, actorRole, schoolID, userID, docID string, req WorkflowActionRequest) error {
	return m.Called(ctx, actorRole, schoolID, userID, docID, req).Error(0)
}
func (m *mockServiceForSecTest) SignDocument(ctx context.Context, actorRole, schoolID, userID, docID string, req SignDocumentRequest) error {
	return m.Called(ctx, actorRole, schoolID, userID, docID, req).Error(0)
}
func (m *mockServiceForSecTest) GetInbox(ctx context.Context, schoolID string) ([]DocumentListResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentListResponse), args.Error(1)
}
func (m *mockServiceForSecTest) GetReviewQueue(ctx context.Context, schoolID string) ([]DocumentListResponse, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentListResponse), args.Error(1)
}
func (m *mockServiceForSecTest) GetMyDocuments(ctx context.Context, schoolID, userID string) ([]DocumentListResponse, error) {
	args := m.Called(ctx, schoolID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentListResponse), args.Error(1)
}
func (m *mockServiceForSecTest) ListDocuments(ctx context.Context, actorRole, schoolID string, docType *DocType) ([]DocumentListResponse, error) {
	args := m.Called(ctx, actorRole, schoolID, docType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentListResponse), args.Error(1)
}
func (m *mockServiceForSecTest) CreateTemplate(ctx context.Context, schoolID string, req TemplateRequest) error {
	return m.Called(ctx, schoolID, req).Error(0)
}
func (m *mockServiceForSecTest) ListTemplates(ctx context.Context, schoolID string) ([]DocumentTemplate, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentTemplate), args.Error(1)
}
func (m *mockServiceForSecTest) UpdateTemplate(ctx context.Context, schoolID, id string, req TemplateRequest) error {
	return m.Called(ctx, schoolID, id, req).Error(0)
}
func (m *mockServiceForSecTest) DeleteTemplate(ctx context.Context, schoolID, id string) error {
	return m.Called(ctx, schoolID, id).Error(0)
}
func (m *mockServiceForSecTest) ExportDocument(ctx context.Context, actorRole, schoolID, id, format string) ([]byte, string, error) {
	args := m.Called(ctx, actorRole, schoolID, id, format)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).([]byte), args.String(1), args.Error(2)
}
func (m *mockServiceForSecTest) LockDocument(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockServiceForSecTest) GetDocumentVersions(ctx context.Context, actorRole, schoolID, docID string) ([]DocumentVersion, error) {
	args := m.Called(ctx, actorRole, schoolID, docID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentVersion), args.Error(1)
}

func TestGetInbox_TeacherForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForSecTest)
	h := NewHandler(mockSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
	})
	h.RegisterRoutes(r.Group("/api/v1"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/documents/inbox", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestSignDocument_PrincipalAllowed_TeacherForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForSecTest)
	h := NewHandler(mockSvc)

	// 1. Teacher -> Forbidden
	r1 := gin.New()
	r1.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
	})
	h.RegisterRoutes(r1.Group("/api/v1"))

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/documents/doc-1/sign", strings.NewReader(`{}`))
	req1.Header.Set("Content-Type", "application/json")
	r1.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusForbidden, w1.Code)

	// 2. Vice Principal -> Authorized (proceeds to service call)
	r2 := gin.New()
	r2.Use(func(c *gin.Context) {
		c.Set("user_id", "vp-1")
		c.Set("role", "vice_principal")
		c.Set("school_id", "school-1")
	})
	h.RegisterRoutes(r2.Group("/api/v1"))

	mockSvc.On("SignDocument", mock.Anything, "vice_principal", "school-1", "vp-1", "doc-1", mock.Anything).
		Return(nil).Once()

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/documents/doc-1/sign", strings.NewReader(`{"certificate_data":"valid-cert","signature_data":"valid-sig"}`))
	req2.Header.Set("Content-Type", "application/json")
	r2.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
