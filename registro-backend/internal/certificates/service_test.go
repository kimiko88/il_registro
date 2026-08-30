package certificates

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"registro-backend/internal/users"
)

// MockRepository mocks certificates.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, cert *Certificate) error {
	args := m.Called(ctx, cert)
	if cert.ID == "" {
		cert.ID = "cert-uuid-123"
	}
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*Certificate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Certificate), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, schoolID, studentID string, certType CertificateType, year string) ([]Certificate, error) {
	args := m.Called(ctx, schoolID, studentID, certType, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Certificate), args.Error(1)
}

func (m *MockRepository) SoftDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) NextProtocolNo(ctx context.Context, schoolID string) (string, error) {
	args := m.Called(ctx, schoolID)
	return args.String(0), args.Error(1)
}

// MockUsersRepository mocks users.Repository
type MockUsersRepository struct {
	mock.Mock
	users.Repository
}

func (m *MockUsersRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func TestGenerateCertificate_Success(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUsersRepository)
	svc := NewService(repo, userRepo)

	ctx := context.Background()
	schoolID := "school-1"
	studentID := "student-1"
	actorID := "admin-1"

	userRepo.On("GetByID", ctx, studentID).Return(&users.User{
		ID:        studentID,
		FirstName: "Mario",
		LastName:  "Rossi",
		SchoolID:  &schoolID,
	}, nil)

	repo.On("NextProtocolNo", ctx, schoolID).Return("PROT-2025-00001", nil)
	repo.On("Create", ctx, mock.AnythingOfType("*certificates.Certificate")).Return(nil)

	req := GenerateCertificateRequest{
		StudentID:    studentID,
		Type:         CertIscrizione,
		AcademicYear: "2024/2025",
		Notes:        "Certificato per uso sportivo",
	}

	cert, pdfBytes, err := svc.GenerateCertificate(ctx, actorID, schoolID, req)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
	assert.NotEmpty(t, pdfBytes)
	assert.Equal(t, "Rossi Mario", cert.StudentName)
	assert.Equal(t, "PROT-2025-00001", cert.ProtocolNo)
	assert.Equal(t, CertIscrizione, cert.Type)

	repo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestGenerateCertificate_StudentSchoolMismatch(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUsersRepository)
	svc := NewService(repo, userRepo)

	ctx := context.Background()
	schoolID := "school-1"
	otherSchoolID := "school-other"
	studentID := "student-1"

	userRepo.On("GetByID", ctx, studentID).Return(&users.User{
		ID:        studentID,
		FirstName: "Mario",
		LastName:  "Rossi",
		SchoolID:  &otherSchoolID,
	}, nil)

	req := GenerateCertificateRequest{
		StudentID: studentID,
		Type:      CertFrequenza,
	}

	cert, pdfBytes, err := svc.GenerateCertificate(ctx, "admin-1", schoolID, req)
	assert.Error(t, err)
	assert.Nil(t, cert)
	assert.Nil(t, pdfBytes)
	assert.Contains(t, err.Error(), "unauthorized: lo studente non appartiene alla scuola")
}

func TestDeleteCertificate_Authorization(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo, nil)
	ctx := context.Background()

	// Teacher cannot delete certificate
	err := svc.DeleteCertificate(ctx, "cert-1", "teacher-1", "teacher")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")

	// Admin can delete certificate
	repo.On("SoftDelete", ctx, "cert-1").Return(nil)
	err = svc.DeleteCertificate(ctx, "cert-1", "admin-1", "admin")
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestGenerateCertificate_InvalidType(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	_, _, err := svc.GenerateCertificate(context.Background(), "admin-1", "school-1", GenerateCertificateRequest{
		StudentID: "student-1",
		Type:      CertificateType("invalid_type_xyz"),
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid certificate type")
}

func TestGetAndListCertificates(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo, nil)
	ctx := context.Background()

	expectedCert := &Certificate{
		ID:         "cert-100",
		SchoolID:   "school-1",
		StudentID:  "student-1",
		Type:       CertPromozione,
		ProtocolNo: "PROT-100",
	}

	repo.On("FindByID", ctx, "cert-100").Return(expectedCert, nil)
	cert, err := svc.GetCertificateByID(ctx, "cert-100")
	assert.NoError(t, err)
	assert.Equal(t, "cert-100", cert.ID)

	repo.On("List", ctx, "school-1", "student-1", CertPromozione, "2024/2025").Return([]Certificate{*expectedCert}, nil)
	list, err := svc.ListCertificates(ctx, "school-1", "student-1", CertPromozione, "2024/2025")
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	repo.AssertExpectations(t)
}

func TestGeneratePDFBytes(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUsersRepository)
	svc := NewService(repo, userRepo)
	ctx := context.Background()

	cert := &Certificate{
		ID:          "cert-pdf",
		SchoolID:    "school-1",
		StudentID:   "student-pdf",
		StudentName: "[Studente]",
		Type:        CertBehavior,
		IssuedAt:    time.Now(),
	}

	repo.On("FindByID", ctx, "cert-pdf").Return(cert, nil)
	schoolID := "school-1"
	userRepo.On("GetByID", ctx, "student-pdf").Return(&users.User{
		ID:        "student-pdf",
		FirstName: "Anna",
		LastName:  "Verdi",
		SchoolID:  &schoolID,
	}, nil)

	pdf, err := svc.GeneratePDFBytes(ctx, "cert-pdf")
	assert.NoError(t, err)
	assert.NotEmpty(t, pdf)
}

func TestGeneratePDFBytes_NotFound(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo, nil)
	ctx := context.Background()

	repo.On("FindByID", ctx, "nonexistent").Return(nil, errors.New("not found"))

	pdf, err := svc.GeneratePDFBytes(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, pdf)
}

func TestDeleteCertificate_StaffRoles(t *testing.T) {
	ctx := context.Background()

	// Superadmin is authorized
	repo := new(MockRepository)
	svc := NewService(repo, nil)
	repo.On("SoftDelete", ctx, "cert-10").Return(nil).Once()
	err := svc.DeleteCertificate(ctx, "cert-10", "staff-1", "superadmin")
	assert.NoError(t, err)
	repo.AssertExpectations(t)

	// Secretary, Teacher, Student are unauthorized
	for _, role := range []string{"secretary", "teacher", "student", "parent"} {
		err := svc.DeleteCertificate(ctx, "cert-10", "user-1", role)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	}
}

func TestGenerateCertificate_Errors(t *testing.T) {
	ctx := context.Background()
	schoolID := "school-1"
	studentID := "student-1"

	t.Run("UserLookupError", func(t *testing.T) {
		repo := new(MockRepository)
		userRepo := new(MockUsersRepository)
		svc := NewService(repo, userRepo)

		userRepo.On("GetByID", ctx, studentID).Return(nil, errors.New("user db error")).Once()

		_, _, err := svc.GenerateCertificate(ctx, "admin-1", schoolID, GenerateCertificateRequest{
			StudentID: studentID,
			Type:      CertIscrizione,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user db error")
	})

	t.Run("NextProtocolFallbackAndCreateError", func(t *testing.T) {
		repo := new(MockRepository)
		userRepo := new(MockUsersRepository)
		svc := NewService(repo, userRepo)

		userRepo.On("GetByID", ctx, studentID).Return(&users.User{
			ID:        studentID,
			FirstName: "Mario",
			LastName:  "Rossi",
			SchoolID:  &schoolID,
		}, nil).Once()
		repo.On("NextProtocolNo", ctx, schoolID).Return("", errors.New("protocol seq error")).Once()
		repo.On("Create", ctx, mock.AnythingOfType("*certificates.Certificate")).Return(errors.New("db insert error")).Once()

		_, _, err := svc.GenerateCertificate(ctx, "admin-1", schoolID, GenerateCertificateRequest{
			StudentID: studentID,
			Type:      CertIscrizione,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create certificate record")
	})
}


