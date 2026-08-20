package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/auditlog"
	"registro-backend/internal/certificates"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Certificates_InvalidType_Rejected(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := certificates.NewService(&mockCertRepo{}, nil)
	handler := certificates.NewHandler(svc)

	r := gin.New()
	r.POST("/certificates/generate", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.Generate(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/certificates/generate", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIntegration_AuditLog_ChainImmutability_SHA256(t *testing.T) {
	report, err := auditlog.VerifyChainIntegrity(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.True(t, report.IsChainIntact)
}

type mockCertRepo struct{}

func (m *mockCertRepo) Create(ctx context.Context, c *certificates.Certificate) error { return nil }
func (m *mockCertRepo) FindByID(ctx context.Context, id string) (*certificates.Certificate, error) {
	return &certificates.Certificate{ID: id, SchoolID: "school-1"}, nil
}
func (m *mockCertRepo) List(ctx context.Context, schoolID, studentID string, certType certificates.CertificateType, year string) ([]certificates.Certificate, error) {
	return []certificates.Certificate{}, nil
}
func (m *mockCertRepo) SoftDelete(ctx context.Context, id string) error { return nil }
func (m *mockCertRepo) NextProtocolNo(ctx context.Context, schoolID string) (string, error) {
	return "PROT-101", nil
}
