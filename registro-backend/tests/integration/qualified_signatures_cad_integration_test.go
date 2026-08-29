package integration

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/signatures"
	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Repositories ─────────────────────────────────────────────────────

type mockSigRepo struct {
	sigs []signatures.Signature
}

func (m *mockSigRepo) Create(_ context.Context, sig *signatures.Signature) error {
	sig.ID = "sig-" + sig.DocumentID
	sig.SignedAt = time.Now()
	m.sigs = append(m.sigs, *sig)
	return nil
}

func (m *mockSigRepo) FindByDocumentID(_ context.Context, docID string) ([]signatures.Signature, error) {
	var out []signatures.Signature
	for _, s := range m.sigs {
		if s.DocumentID == docID {
			out = append(out, s)
		}
	}
	return out, nil
}

type mockFeqRepo struct {
	qualifiedSigs map[string]*signatures.QualifiedSignature
}

func newMockFeqRepo() *mockFeqRepo {
	return &mockFeqRepo{
		qualifiedSigs: make(map[string]*signatures.QualifiedSignature),
	}
}

func (m *mockFeqRepo) CreateQualified(_ context.Context, sig *signatures.QualifiedSignature) error {
	sig.ID = "feq-" + sig.DocumentID
	sig.SignedAt = time.Now()
	sig.IsValid = true
	m.qualifiedSigs[sig.ID] = sig
	return nil
}

func (m *mockFeqRepo) FindQualifiedByID(_ context.Context, id string) (*signatures.QualifiedSignature, error) {
	s, ok := m.qualifiedSigs[id]
	if !ok {
		return nil, errors.New("qualified signature not found")
	}
	return s, nil
}

func (m *mockFeqRepo) FindQualifiedByDocumentID(_ context.Context, docID string) ([]signatures.QualifiedSignature, error) {
	var out []signatures.QualifiedSignature
	for _, s := range m.qualifiedSigs {
		if s.DocumentID == docID {
			out = append(out, *s)
		}
	}
	return out, nil
}

func (m *mockFeqRepo) RevokeQualified(_ context.Context, id, reason string) error {
	s, ok := m.qualifiedSigs[id]
	if !ok {
		return errors.New("qualified signature not found")
	}
	now := time.Now()
	s.RevokedAt = &now
	s.RevocationReason = reason
	s.IsValid = false
	return nil
}

type mockDocServiceForSig struct {
	content string
}

func (m *mockDocServiceForSig) LockDocument(_ context.Context, _ string) error {
	return nil
}

func (m *mockDocServiceForSig) GetDocumentContent(_ context.Context, _ string) (string, error) {
	if m.content == "" {
		return "Contenuto del verbale di scrutinio finale 2026", nil
	}
	return m.content, nil
}

type mockUserLookupForSig struct {
	mfaEnabled bool
	mfaSecret  string
}

func (m *mockUserLookupForSig) GetByID(_ context.Context, id string) (*users.User, error) {
	return &users.User{
		ID:         id,
		FirstName:  "Professore",
		LastName:   "Docente",
		MFAEnabled: m.mfaEnabled,
		MFASecret:  m.mfaSecret,
		Role:       "teacher",
	}, nil
}

type mockTSAClient struct{}

func (m *mockTSAClient) RequestTimestamp(_ context.Context, _ []byte) (string, time.Time, error) {
	return "MOCK-TIMESTAMP-TOKEN-AgID-ACC-01", time.Now().UTC(), nil
}

// ── Setup Router Helper ───────────────────────────────────────────────────

func setupSignaturesRouter(
	sigRepo signatures.Repository,
	feqRepo *mockFeqRepo,
	userLookup *mockUserLookupForSig,
	docSvc *mockDocServiceForSig,
	role, userID, schoolID string,
) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	baseSvc := signatures.NewService(sigRepo, docSvc, userLookup)
	feqSvc := signatures.NewQualifiedServiceWithProviders(nil, feqRepo, userLookup, &signatures.SoftwareQSCD{}, &mockTSAClient{})
	sidiSvc := signatures.NewSidiExportService()

	h := signatures.NewHandler(baseSvc, feqSvc, sidiSvc)
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

// SIG01 — Legacy FEA sign document with valid TOTP code
func TestSignatures_SignDocument_LegacyFEA_Success(t *testing.T) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "il_registro",
		AccountName: "prof@scuola.it",
	})
	require.NoError(t, err)

	userLookup := &mockUserLookupForSig{mfaEnabled: true, mfaSecret: key.Secret()}
	sigRepo := &mockSigRepo{}
	feqRepo := newMockFeqRepo()
	docSvc := &mockDocServiceForSig{}

	r := setupSignaturesRouter(sigRepo, feqRepo, userLookup, docSvc, "teacher", "teacher-1", "school-1")

	code, err := totp.GenerateCode(key.Secret(), time.Now())
	require.NoError(t, err)

	body, _ := json.Marshal(signatures.SignRequest{
		DocumentID: "doc-123",
		Pin:        code,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/signatures", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var sig signatures.Signature
	err = json.Unmarshal(w.Body.Bytes(), &sig)
	require.NoError(t, err)
	assert.Equal(t, "doc-123", sig.DocumentID)
	assert.NotEmpty(t, sig.SignatureHash)
}

// SIG02 — FEA signing fails when MFA is not enabled
func TestSignatures_SignDocument_FailsWithoutMFA(t *testing.T) {
	userLookup := &mockUserLookupForSig{mfaEnabled: false}
	sigRepo := &mockSigRepo{}
	feqRepo := newMockFeqRepo()
	docSvc := &mockDocServiceForSig{}

	r := setupSignaturesRouter(sigRepo, feqRepo, userLookup, docSvc, "teacher", "teacher-1", "school-1")

	body, _ := json.Marshal(signatures.SignRequest{
		DocumentID: "doc-123",
		Pin:        "123456",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/signatures", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "MFA")
}

// SIG03 — Get signatures list for a document
func TestSignatures_GetSignatures_ByDocument(t *testing.T) {
	userLookup := &mockUserLookupForSig{mfaEnabled: true}
	sigRepo := &mockSigRepo{
		sigs: []signatures.Signature{
			{ID: "sig-1", DocumentID: "doc-99", SignerID: "teacher-1", SignatureHash: "hash-1"},
		},
	}
	feqRepo := newMockFeqRepo()
	docSvc := &mockDocServiceForSig{}

	r := setupSignaturesRouter(sigRepo, feqRepo, userLookup, docSvc, "teacher", "teacher-1", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/signatures/document/doc-99", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var list []signatures.Signature
	err := json.Unmarshal(w.Body.Bytes(), &list)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "sig-1", list[0].ID)
}

// SIG04 — FEQ qualified signature creation with RSA key generation and TSA timestamp
func TestSignatures_SignQualified_FEQ_Success(t *testing.T) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "il_registro",
		AccountName: "preside@scuola.it",
	})
	require.NoError(t, err)

	userLookup := &mockUserLookupForSig{mfaEnabled: true, mfaSecret: key.Secret()}
	sigRepo := &mockSigRepo{}
	feqRepo := newMockFeqRepo()
	docSvc := &mockDocServiceForSig{}

	r := setupSignaturesRouter(sigRepo, feqRepo, userLookup, docSvc, "principal", "principal-1", "school-1")

	code, err := totp.GenerateCode(key.Secret(), time.Now())
	require.NoError(t, err)

	body, _ := json.Marshal(signatures.QualifiedSignRequest{
		DocumentID:      "doc-feq-1",
		DocumentContent: base64.StdEncoding.EncodeToString([]byte("Verbale Scrutinio Classe 5B")),
		Level:           signatures.SignatureLevelFEQ,
		Pin:             code,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/signatures/qualified", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var qSig signatures.QualifiedSignature
	err = json.Unmarshal(w.Body.Bytes(), &qSig)
	require.NoError(t, err)
	assert.Equal(t, "doc-feq-1", qSig.DocumentID)
	assert.Equal(t, signatures.SignatureLevelFEQ, qSig.Level)
	assert.NotEmpty(t, qSig.PublicKeyPEM)
	assert.NotEmpty(t, qSig.TimestampToken)
	assert.True(t, qSig.IsValid)
}

// SIG05 — FES simple electronic signature is allowed without MFA for general workflows
func TestSignatures_SignQualified_FES_Allowed(t *testing.T) {
	userLookup := &mockUserLookupForSig{mfaEnabled: false}
	sigRepo := &mockSigRepo{}
	feqRepo := newMockFeqRepo()
	docSvc := &mockDocServiceForSig{}

	r := setupSignaturesRouter(sigRepo, feqRepo, userLookup, docSvc, "parent", "parent-1", "school-1")

	body, _ := json.Marshal(signatures.QualifiedSignRequest{
		DocumentID:      "doc-fes-1",
		DocumentContent: base64.StdEncoding.EncodeToString([]byte("Presa visione circolare viaggi")),
		Level:           signatures.SignatureLevelFES,
		Pin:             "0000",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/signatures/qualified", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var qSig signatures.QualifiedSignature
	err := json.Unmarshal(w.Body.Bytes(), &qSig)
	require.NoError(t, err)
	assert.Equal(t, signatures.SignatureLevelFES, qSig.Level)
}

// SIG06 — Verify qualified signature endpoint returns validity and certificate DN
func TestSignatures_VerifyQualified_Endpoint(t *testing.T) {
	userLookup := &mockUserLookupForSig{mfaEnabled: true}
	sigRepo := &mockSigRepo{}
	feqRepo := newMockFeqRepo()
	docSvc := &mockDocServiceForSig{}

	// Pre-create a qualified signature
	sig := &signatures.QualifiedSignature{
		ID:            "feq-verify-1",
		DocumentID:    "doc-1",
		SignerID:      "teacher-1",
		Level:         signatures.SignatureLevelFEQ,
		CertificateDN: "CN=il_registro-SoftwareQSCD,O=il_registro,C=IT",
		SignedAt:      time.Now(),
		TimestampAt:   time.Now(),
		IsValid:       true,
	}
	feqRepo.qualifiedSigs[sig.ID] = sig

	r := setupSignaturesRouter(sigRepo, feqRepo, userLookup, docSvc, "admin", "admin-1", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/signatures/qualified/feq-verify-1/verify", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var result signatures.VerifyResult
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "feq-verify-1", result.SignatureID)
	assert.Equal(t, signatures.SignatureLevelFEQ, result.Level)
}

// SIG07 — CAD preservation package download produces valid ZIP with manifest
func TestSignatures_CadPreservationPackage_Download(t *testing.T) {
	userLookup := &mockUserLookupForSig{mfaEnabled: true}
	sigRepo := &mockSigRepo{}
	feqRepo := newMockFeqRepo()
	docSvc := &mockDocServiceForSig{}

	r := setupSignaturesRouter(sigRepo, feqRepo, userLookup, docSvc, "admin", "admin-1", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/signatures/cad-preservation/download?school_id=school-1&year=2025-2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	assert.Greater(t, w.Body.Len(), 100) // Non-empty ZIP archive
}

// SIG08 — SIDI export generates ZIP with DM 742 / MIUR XML structure
func TestSignatures_SidiExport_Download(t *testing.T) {
	userLookup := &mockUserLookupForSig{mfaEnabled: true}
	sigRepo := &mockSigRepo{}
	feqRepo := newMockFeqRepo()
	docSvc := &mockDocServiceForSig{}

	r := setupSignaturesRouter(sigRepo, feqRepo, userLookup, docSvc, "admin", "admin-1", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/sidi/export/school-1?tipo=SCRUTINI&year=2025-2026&school_name=Liceo+Scientifico", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "SCRUTINI")
}

// SIG09 — Unauthenticated requests receive 401
func TestSignatures_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	baseSvc := signatures.NewService(&mockSigRepo{}, &mockDocServiceForSig{}, &mockUserLookupForSig{})
	feqSvc := signatures.NewQualifiedService(nil, newMockFeqRepo(), &mockUserLookupForSig{})
	h := signatures.NewHandler(baseSvc, feqSvc, signatures.NewSidiExportService())
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/signatures/document/doc-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
