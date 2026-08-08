package signatures

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"registro-backend/internal/users"
)

// ── Mock feqRepository ────────────────────────────────────────────────────

type MockFEQRepo struct {
	mock.Mock
}

func (m *MockFEQRepo) CreateQualified(_ context.Context, sig *QualifiedSignature) error {
	args := m.Called(sig)
	sig.ID = "qsig-001"
	sig.SignedAt = time.Now().UTC()
	return args.Error(0)
}

func (m *MockFEQRepo) FindQualifiedByID(_ context.Context, id string) (*QualifiedSignature, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*QualifiedSignature), args.Error(1)
}

func (m *MockFEQRepo) FindQualifiedByDocumentID(_ context.Context, docID string) ([]QualifiedSignature, error) {
	args := m.Called(docID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]QualifiedSignature), args.Error(1)
}

func (m *MockFEQRepo) RevokeQualified(_ context.Context, id string, reason string) error {
	args := m.Called(id, reason)
	return args.Error(0)
}

// ── Mock UserLookup ───────────────────────────────────────────────────────

type MockFEQUserRepo struct {
	mock.Mock
}

func (m *MockFEQUserRepo) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

// ── Helper ────────────────────────────────────────────────────────────────

func docBase64() string {
	// "Documento di test per firma FEQ" in base64
	return "RG9jdW1lbnRvIGRpIHRlc3QgcGVyIGZpcm1hIEZFUQ=="
}

// ── Test SignQualified ────────────────────────────────────────────────────

func TestFEQService_SignQualified_FES_NoMFA(t *testing.T) {
	repo := new(MockFEQRepo)
	userRepo := new(MockFEQUserRepo)
	svc := NewQualifiedService(nil, repo, userRepo)

	user := &users.User{ID: "u-1", FirstName: "Mario", LastName: "Rossi", MFAEnabled: false}
	userRepo.On("GetByID", mock.Anything, "u-1").Return(user, nil).Once()
	repo.On("CreateQualified", mock.AnythingOfType("*signatures.QualifiedSignature")).Return(nil).Once()

	req := QualifiedSignRequest{
		DocumentID:      "doc-fes-001",
		DocumentContent: docBase64(),
		Level:           SignatureLevelFES,
		Pin:             "qualsiasi", // FES non richiede TOTP
	}

	sig, err := svc.SignQualified(context.Background(), "u-1", "school-001", req)
	assert.NoError(t, err)
	assert.Equal(t, SignatureLevelFES, sig.Level)
	assert.NotEmpty(t, sig.DocumentHash)
	assert.NotEmpty(t, sig.SignatureValue)
	assert.NotEmpty(t, sig.TimestampToken)
	assert.True(t, sig.IsValid)
	repo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestFEQService_SignQualified_FEQ_RequiresMFA(t *testing.T) {
	repo := new(MockFEQRepo)
	userRepo := new(MockFEQUserRepo)
	svc := NewQualifiedService(nil, repo, userRepo)

	user := &users.User{ID: "u-2", FirstName: "Luigi", LastName: "Verdi", MFAEnabled: false}
	userRepo.On("GetByID", mock.Anything, "u-2").Return(user, nil).Once()

	req := QualifiedSignRequest{
		DocumentID:      "doc-feq-001",
		DocumentContent: docBase64(),
		Level:           SignatureLevelFEQ,
		Pin:             "123456",
	}

	_, err := svc.SignQualified(context.Background(), "u-2", "school-001", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MFA TOTP abilitato")
	userRepo.AssertExpectations(t)
}

func TestFEQService_SignQualified_FEQ_InvalidTOTP(t *testing.T) {
	repo := new(MockFEQRepo)
	userRepo := new(MockFEQUserRepo)
	svc := NewQualifiedService(nil, repo, userRepo)

	user := &users.User{ID: "u-3", FirstName: "Anna", LastName: "Bianchi", MFAEnabled: true, MFASecret: "JBSWY3DPEHPK3PXP"}
	userRepo.On("GetByID", mock.Anything, "u-3").Return(user, nil).Once()

	req := QualifiedSignRequest{
		DocumentID:      "doc-feq-002",
		DocumentContent: docBase64(),
		Level:           SignatureLevelFEQ,
		Pin:             "000000", // OTP errato
	}

	_, err := svc.SignQualified(context.Background(), "u-3", "school-001", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OTP non valido")
	userRepo.AssertExpectations(t)
}

func TestFEQService_VerifyQualified_Valid(t *testing.T) {
	repo := new(MockFEQRepo)
	userRepo := new(MockFEQUserRepo)
	svc := NewQualifiedService(nil, repo, userRepo)

	qscd := &SoftwareQSCD{}
	docHash := sha256.Sum256([]byte("Documento di test per firma FEQ"))
	docHashHex := hex.EncodeToString(docHash[:])
	sigBytes, pubKeyPEM, certDN, certSerial, err := qscd.Sign(context.Background(), docHash[:])
	require.NoError(t, err)
	sigHex := hex.EncodeToString(sigBytes)

	now := time.Now().UTC()
	qSig := &QualifiedSignature{
		ID:                "qsig-verify-1",
		DocumentID:        "doc-001",
		SignerID:          "u-1",
		Level:             SignatureLevelFEQ,
		DocumentHash:      docHashHex,
		SignatureValue:    sigHex,
		PublicKeyPEM:      string(pubKeyPEM),
		CertificateSerial: certSerial,
		CertificateDN:     certDN,
		TimestampToken:    "cafebabe1234",
		SignedAt:          now,
		TimestampAt:       now,
		IsValid:           true,
		RevokedAt:         nil,
	}
	repo.On("FindQualifiedByID", "qsig-verify-1").Return(qSig, nil).Once()

	result, err := svc.VerifyQualified(context.Background(), "qsig-verify-1")
	assert.NoError(t, err)
	assert.True(t, result.IsValid)
	assert.Equal(t, SignatureLevelFEQ, result.Level)
	assert.Contains(t, result.Message, "valida")
	repo.AssertExpectations(t)
}

func TestFEQService_VerifyQualified_Revoked(t *testing.T) {
	repo := new(MockFEQRepo)
	userRepo := new(MockFEQUserRepo)
	svc := NewQualifiedService(nil, repo, userRepo)

	now := time.Now().UTC()
	revokedAt := now.Add(-24 * time.Hour)
	qSig := &QualifiedSignature{
		ID:             "qsig-revoked-1",
		DocumentID:     "doc-002",
		SignerID:       "u-1",
		Level:          SignatureLevelFEQ,
		DocumentHash:   "abc123",
		SignatureValue: "sig123",
		TimestampToken: "ts123",
		SignedAt:       now.Add(-48 * time.Hour),
		TimestampAt:    now.Add(-48 * time.Hour),
		IsValid:        false,
		RevokedAt:      &revokedAt,
	}
	repo.On("FindQualifiedByID", "qsig-revoked-1").Return(qSig, nil).Once()

	result, err := svc.VerifyQualified(context.Background(), "qsig-revoked-1")
	assert.NoError(t, err)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Message, "revocata")
	repo.AssertExpectations(t)
}

// ── Test SIDI Export ──────────────────────────────────────────────────────

func TestSidiExport_GeneratePackage_Scrutini(t *testing.T) {
	svc := NewSidiExportService()
	zipBytes, record, err := svc.GenerateSidiPackage(
		context.Background(),
		"BOIT00100A", "Istituto Tecnico Informatica", "2025/2026", "SCRUTINI",
	)
	assert.NoError(t, err)
	assert.NotNil(t, zipBytes)
	assert.Greater(t, len(zipBytes), 100)
	assert.Equal(t, "BOIT00100A", record.CodiceMeccanografico)
	assert.Equal(t, "SCRUTINI", record.TipologiaExport)
	assert.Equal(t, "PRONTO", record.StatoTrasmissione)
	assert.NotEmpty(t, record.HashIntegrità)
}

func TestSidiExport_GeneratePackage_InvalidTipologia(t *testing.T) {
	svc := NewSidiExportService()
	_, _, err := svc.GenerateSidiPackage(
		context.Background(),
		"BOIT00100A", "Scuola", "2025/2026", "INVALIDO",
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tipologia non valida")
}

func TestCadPreservation_GeneratePackage(t *testing.T) {
	zipBytes, err := GenerateCadPreservationPackage(context.Background(), "BOIT00100A", "2025/2026")
	assert.NoError(t, err)
	assert.NotNil(t, zipBytes)
	assert.Greater(t, len(zipBytes), 200)
}

func TestCadPreservation_MissingSchoolID(t *testing.T) {
	_, err := GenerateCadPreservationPackage(context.Background(), "", "2025/2026")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "schoolID")
}
