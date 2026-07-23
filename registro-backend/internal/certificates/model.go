package certificates

import (
	"time"
)

type CertificateType string

const (
	CertIscrizione CertificateType = "iscrizione"
	CertFrequenza  CertificateType = "frequenza"
	CertPromozione CertificateType = "promozione"
	CertBehavior   CertificateType = "condotta"
)

type Certificate struct {
	ID           string          `json:"id" db:"id"`
	SchoolID     string          `json:"school_id" db:"school_id"`
	StudentID    string          `json:"student_id" db:"student_id"`
	StudentName  string          `json:"student_name,omitempty" db:"student_name"`
	ClassName    string          `json:"class_name,omitempty" db:"class_name"`
	Type         CertificateType `json:"type" db:"type"`
	IssuedBy     string          `json:"issued_by" db:"issued_by"`
	IssuedByName string          `json:"issued_by_name,omitempty" db:"issued_by_name"`
	IssuedAt     time.Time       `json:"issued_at" db:"issued_at"`
	AcademicYear string          `json:"academic_year" db:"academic_year"`
	Notes        string          `json:"notes" db:"notes"`
	PDFUrl       string          `json:"pdf_url" db:"pdf_url"`
	ProtocolNo   string          `json:"protocol_no" db:"protocol_no"`
	IsDeleted    bool            `json:"is_deleted" db:"is_deleted"`
}

type GenerateCertificateRequest struct {
	StudentID    string          `json:"student_id" binding:"required"`
	Type         CertificateType `json:"type" binding:"required"`
	AcademicYear string          `json:"academic_year"`
	Notes        string          `json:"notes"`
}
