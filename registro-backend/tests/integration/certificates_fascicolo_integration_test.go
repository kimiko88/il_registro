package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/certificates"
	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCertificatesRepo struct {
	certs   map[string]*certificates.Certificate
	protoNo int
}

func newMockCertificatesRepo() *mockCertificatesRepo {
	return &mockCertificatesRepo{
		certs:   make(map[string]*certificates.Certificate),
		protoNo: 100,
	}
}

func (m *mockCertificatesRepo) Create(ctx context.Context, c *certificates.Certificate) error {
	if c.ID == "" {
		c.ID = "cert-" + time.Now().Format("150405.000000")
	}
	c.IssuedAt = time.Now()
	m.certs[c.ID] = c
	return nil
}

func (m *mockCertificatesRepo) FindByID(ctx context.Context, id string) (*certificates.Certificate, error) {
	c, ok := m.certs[id]
	if !ok {
		return nil, assert.AnError
	}
	return c, nil
}

func (m *mockCertificatesRepo) List(ctx context.Context, schoolID, studentID string, certType certificates.CertificateType, year string) ([]certificates.Certificate, error) {
	var res []certificates.Certificate
	for _, c := range m.certs {
		if schoolID != "" && c.SchoolID != schoolID {
			continue
		}
		if studentID != "" && c.StudentID != studentID {
			continue
		}
		res = append(res, *c)
	}
	return res, nil
}

func (m *mockCertificatesRepo) SoftDelete(ctx context.Context, id string) error {
	delete(m.certs, id)
	return nil
}

func (m *mockCertificatesRepo) NextProtocolNo(ctx context.Context, schoolID string) (string, error) {
	m.protoNo++
	return "PROT-2026-00123", nil
}

type mockUserRepoForCerts struct {
	users.Repository
}

func (m *mockUserRepoForCerts) GetByID(ctx context.Context, id string) (*users.User, error) {
	schoolID := "school-1"
	return &users.User{
		ID:        id,
		SchoolID:  &schoolID,
		FirstName: "Mario",
		LastName:  "Rossi",
		Role:      "student",
	}, nil
}

func setupCertificatesTestRouter(repo certificates.Repository, uRepo users.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := certificates.NewService(repo, uRepo)
	handler := certificates.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = "secretary"
		}
		c.Set("role", role)
		c.Set("user_id", "secr-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_Certificates_Fascicolo_Workflow(t *testing.T) {
	repo := newMockCertificatesRepo()
	uRepo := &mockUserRepoForCerts{}
	r := setupCertificatesTestRouter(repo, uRepo)

	// 1. Generate an Attendance Certificate (Certificato di Frequenza)
	genReq := certificates.GenerateCertificateRequest{
		StudentID:    "student-1",
		Type:         certificates.CertFrequenza,
		AcademicYear: "2025/2026",
		Notes:        "Rilasciato per usi consentiti dalla legge (detrazioni e trasporti)",
	}
	body, _ := json.Marshal(genReq)
	req, _ := http.NewRequest("POST", "/api/v1/certificates/generate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var certResp struct {
		ID     string                    `json:"id"`
		PDFUrl string                    `json:"pdf_url"`
		Cert   *certificates.Certificate `json:"cert"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &certResp)
	require.NoError(t, err)
	require.NotNil(t, certResp.Cert)
	assert.NotEmpty(t, certResp.Cert.ID)
	assert.Equal(t, certificates.CertFrequenza, certResp.Cert.Type)
	assert.NotEmpty(t, certResp.Cert.ProtocolNo)

	// 2. List certificates in student dossier (Fascicolo dello studente)
	listReq, _ := http.NewRequest("GET", "/api/v1/certificates?student_id=student-1", nil)
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, listReq)
	require.Equal(t, http.StatusOK, wList.Code)

	var list []certificates.Certificate
	err = json.Unmarshal(wList.Body.Bytes(), &list)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// 3. Download PDF generated certificate
	pdfReq, _ := http.NewRequest("GET", "/api/v1/certificates/"+certResp.Cert.ID+"/pdf", nil)
	wPdf := httptest.NewRecorder()
	r.ServeHTTP(wPdf, pdfReq)
	require.Equal(t, http.StatusOK, wPdf.Code)
	assert.Equal(t, "application/pdf", wPdf.Header().Get("Content-Type"))
	assert.NotEmpty(t, wPdf.Body.Bytes())
}
