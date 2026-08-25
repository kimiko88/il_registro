package certificates

import (
	"context"
	"fmt"
	"time"

	"registro-backend/internal/users"
)

type Service interface {
	GenerateCertificate(ctx context.Context, actorID, schoolID string, req GenerateCertificateRequest) (*Certificate, []byte, error)
	GetCertificateByID(ctx context.Context, id string) (*Certificate, error)
	ListCertificates(ctx context.Context, schoolID, studentID string, certType CertificateType, year string) ([]Certificate, error)
	DeleteCertificate(ctx context.Context, id, actorID, actorRole string) error
	GeneratePDFBytes(ctx context.Context, id string) ([]byte, error)
}

type service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(repo Repository, userRepo users.Repository) Service {
	return &service{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *service) GenerateCertificate(ctx context.Context, actorID, schoolID string, req GenerateCertificateRequest) (*Certificate, []byte, error) {
	if req.StudentID == "" {
		return nil, nil, fmt.Errorf("student_id required")
	}
	if !IsValidCertificateType(req.Type) {
		return nil, nil, fmt.Errorf("invalid certificate type: %s", req.Type)
	}
	if req.AcademicYear == "" {
		req.AcademicYear = currentAcademicYear()
	}

	var studentName string
	if s.userRepo != nil {
		std, err := s.userRepo.GetByID(ctx, req.StudentID)
		if err != nil {
			return nil, nil, fmt.Errorf("studente non trovato: %w", err)
		}
		// Bug 118: verify student belongs to the issuing school
		if std.SchoolID == nil || *std.SchoolID != schoolID {
			return nil, nil, fmt.Errorf("unauthorized: lo studente non appartiene alla scuola del certificatore")
		}
		studentName = fmt.Sprintf("%s %s", std.LastName, std.FirstName)
	}

	protoNo, err := s.repo.NextProtocolNo(ctx, schoolID)
	if err != nil {
		protoNo = fmt.Sprintf("CERT-%d-00001", time.Now().Year())
	}

	cert := &Certificate{
		SchoolID:     schoolID,
		StudentID:    req.StudentID,
		StudentName:  studentName,
		Type:         req.Type,
		IssuedBy:     actorID,
		IssuedAt:     time.Now(),
		AcademicYear: req.AcademicYear,
		Notes:        req.Notes,
		ProtocolNo:   protoNo,
		PDFUrl:       "/api/v1/certificates/download/temp",
	}

	pdfBytes, err := GenerateCertificatePDF(cert, "Istituto Scolastico il_registro")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	err = s.repo.Create(ctx, cert)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate record: %w", err)
	}

	cert.PDFUrl = fmt.Sprintf("/api/v1/certificates/%s/pdf", cert.ID)

	return cert, pdfBytes, nil
}

func (s *service) GetCertificateByID(ctx context.Context, id string) (*Certificate, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) ListCertificates(ctx context.Context, schoolID, studentID string, certType CertificateType, year string) ([]Certificate, error) {
	return s.repo.List(ctx, schoolID, studentID, certType, year)
}

// DeleteCertificate soft-deletes a certificate.
// Bug 119: restricted to admin and superadmin roles only.
func (s *service) DeleteCertificate(ctx context.Context, id, actorID, actorRole string) error {
	if actorRole != "admin" && actorRole != "superadmin" {
		return fmt.Errorf("unauthorized: solo admin e superadmin possono eliminare i certificati")
	}
	return s.repo.SoftDelete(ctx, id)
}

func (s *service) GeneratePDFBytes(ctx context.Context, id string) ([]byte, error) {
	cert, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if (cert.StudentName == "" || cert.StudentName == "[Studente]") && s.userRepo != nil {
		std, err := s.userRepo.GetByID(ctx, cert.StudentID)
		if err == nil && std != nil {
			cert.StudentName = fmt.Sprintf("%s %s", std.LastName, std.FirstName)
		}
	}
	return GenerateCertificatePDF(cert, "Istituto Scolastico il_registro")
}

// currentAcademicYear returns the school year label for today's date.
// Months September-December belong to the year starting in that calendar year;
// months January-August belong to the year that started in the previous calendar year.
func currentAcademicYear() string {
	now := time.Now()
	year := now.Year()
	if now.Month() < time.September {
		year--
	}
	return fmt.Sprintf("%d/%d", year, year+1)
}
