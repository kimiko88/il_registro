package certificates

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"registro-backend/internal/users"
)

type MockCertService struct {
	mock.Mock
}

func (m *MockCertService) GenerateCertificate(ctx context.Context, actorID, schoolID string, req GenerateCertificateRequest) (*Certificate, []byte, error) {
	args := m.Called(ctx, actorID, schoolID, req)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	var data []byte
	if args.Get(1) != nil {
		data = args.Get(1).([]byte)
	}
	return args.Get(0).(*Certificate), data, args.Error(2)
}

func (m *MockCertService) GetCertificateByID(ctx context.Context, id string) (*Certificate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Certificate), args.Error(1)
}

func (m *MockCertService) ListCertificates(ctx context.Context, schoolID, studentID string, certType CertificateType, year string) ([]Certificate, error) {
	args := m.Called(ctx, schoolID, studentID, certType, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Certificate), args.Error(1)
}

func (m *MockCertService) DeleteCertificate(ctx context.Context, id, actorID, actorRole string) error {
	args := m.Called(ctx, id, actorID, actorRole)
	return args.Error(0)
}

func (m *MockCertService) GeneratePDFBytes(ctx context.Context, id string) ([]byte, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCertService) GetUserRepo() users.Repository {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(users.Repository)
}

func setupCertRouter(svc Service, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestCertificatesHandler_List_Permissions(t *testing.T) {
	svc := new(MockCertService)

	// 1. Unauthorized
	rNoAuth := setupCertRouter(svc, "", "", "")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/certificates", nil)
	rNoAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 2. Forbidden (Student cannot list all certificates)
	rStudent := setupCertRouter(svc, "student-1", "student", "school-1")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/certificates", nil)
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 3. Authorized (Secretary)
	expectedList := []Certificate{
		{
			ID:           "cert-1",
			SchoolID:     "school-1",
			StudentID:    "student-1",
			Type:         CertIscrizione,
			IssuedAt:     time.Now(),
			AcademicYear: "2025/2026",
		},
	}
	svc.On("ListCertificates", mock.Anything, "school-1", "", CertificateType(""), "").Return(expectedList, nil).Once()

	rStaff := setupCertRouter(svc, "sec-1", "secretary", "school-1")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/certificates", nil)
	rStaff.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp []Certificate
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "cert-1", resp[0].ID)
}

func TestCertificatesHandler_Generate(t *testing.T) {
	svc := new(MockCertService)
	rStaff := setupCertRouter(svc, "sec-1", "secretary", "school-1")

	genReq := GenerateCertificateRequest{
		StudentID:    "student-100",
		Type:         CertFrequenza,
		AcademicYear: "2025/2026",
		Notes:        "Certificato di regolare frequenza scolastica",
	}

	createdCert := &Certificate{
		ID:           "cert-new-99",
		SchoolID:     "school-1",
		StudentID:    "student-100",
		Type:         CertFrequenza,
		IssuedBy:     "sec-1",
		IssuedAt:     time.Now(),
		AcademicYear: "2025/2026",
		PDFUrl:       "/uploads/school-1/certificates/cert-new-99.pdf",
		ProtocolNo:   "PROT-2026-00042",
	}

	svc.On("GenerateCertificate", mock.Anything, "sec-1", "school-1", genReq).
		Return(createdCert, []byte("%PDF-1.4 mock pdf"), nil).Once()

	body, _ := json.Marshal(genReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/certificates/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rStaff.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "cert-new-99", resp["id"])
	assert.Equal(t, "/uploads/school-1/certificates/cert-new-99.pdf", resp["pdf_url"])
}

func TestCertificatesHandler_DownloadPDF_StudentAccessCheck(t *testing.T) {
	svc := new(MockCertService)

	targetCert := &Certificate{
		ID:        "cert-target",
		SchoolID:  "school-1",
		StudentID: "student-1",
		PDFUrl:    "/certificates/cert-target.pdf",
	}
	svc.On("GetCertificateByID", mock.Anything, "cert-target").Return(targetCert, nil).Twice()

	// 1. Student accessing someone else's certificate -> Forbidden
	rOtherStudent := setupCertRouter(svc, "student-99", "student", "school-1")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/certificates/cert-target/pdf", nil)
	rOtherStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 2. Staff downloading (Admin) -> Allowed to generate/stream PDF
	svc.On("GeneratePDFBytes", mock.Anything, "cert-target").Return([]byte("%PDF-1.4 mock content"), nil).Once()
	rAdmin := setupCertRouter(svc, "admin-1", "admin", "school-1")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/certificates/cert-target/pdf", nil)
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
}
