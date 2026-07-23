package certificates

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, cert *Certificate) error
	FindByID(ctx context.Context, id string) (*Certificate, error)
	List(ctx context.Context, schoolID, studentID string, certType CertificateType, year string) ([]Certificate, error)
	SoftDelete(ctx context.Context, id string) error
	NextProtocolNo(ctx context.Context, schoolID string) (string, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, cert *Certificate) error {
	if cert.ID == "" {
		cert.ID = uuid.New().String()
	}
	if cert.IssuedAt.IsZero() {
		cert.IssuedAt = time.Now()
	}

	query := `
		INSERT INTO certificates (id, school_id, student_id, type, issued_by, issued_at, academic_year, notes, pdf_url, protocol_no, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, false)
	`
	_, err := r.db.ExecContext(ctx, query,
		cert.ID, cert.SchoolID, cert.StudentID, cert.Type, cert.IssuedBy, cert.IssuedAt, cert.AcademicYear, cert.Notes, cert.PDFUrl, cert.ProtocolNo,
	)
	return err
}

func (r *repository) FindByID(ctx context.Context, id string) (*Certificate, error) {
	query := `
		SELECT c.id, c.school_id, c.student_id, COALESCE(u.last_name || ' ' || u.first_name, '') as student_name,
		       COALESCE(cl.name || cl.section, '') as class_name, c.type, c.issued_by,
		       COALESCE(ib.last_name || ' ' || ib.first_name, '') as issued_by_name,
		       c.issued_at, c.academic_year, c.notes, c.pdf_url, c.protocol_no, c.is_deleted
		FROM certificates c
		LEFT JOIN users u ON c.student_id = u.id
		LEFT JOIN students s ON s.user_id = u.id
		LEFT JOIN users ib ON c.issued_by = ib.id
		LEFT JOIN classes cl ON s.class_id = cl.id
		WHERE c.id = $1 AND c.is_deleted = false
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var c Certificate
	err := row.Scan(
		&c.ID, &c.SchoolID, &c.StudentID, &c.StudentName, &c.ClassName, &c.Type,
		&c.IssuedBy, &c.IssuedByName, &c.IssuedAt, &c.AcademicYear, &c.Notes, &c.PDFUrl, &c.ProtocolNo, &c.IsDeleted,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) List(ctx context.Context, schoolID, studentID string, certType CertificateType, year string) ([]Certificate, error) {
	query := `
		SELECT c.id, c.school_id, c.student_id, COALESCE(u.last_name || ' ' || u.first_name, '') as student_name,
		       COALESCE(cl.name || cl.section, '') as class_name, c.type, c.issued_by,
		       COALESCE(ib.last_name || ' ' || ib.first_name, '') as issued_by_name,
		       c.issued_at, c.academic_year, c.notes, c.pdf_url, c.protocol_no, c.is_deleted
		FROM certificates c
		LEFT JOIN users u ON c.student_id = u.id
		LEFT JOIN students s ON s.user_id = u.id
		LEFT JOIN users ib ON c.issued_by = ib.id
		LEFT JOIN classes cl ON s.class_id = cl.id
		WHERE c.is_deleted = false
	`
	args := []interface{}{}
	argID := 1

	if schoolID != "" {
		query += fmt.Sprintf(" AND c.school_id = $%d", argID)
		args = append(args, schoolID)
		argID++
	}
	if studentID != "" {
		query += fmt.Sprintf(" AND c.student_id = $%d", argID)
		args = append(args, studentID)
		argID++
	}
	if certType != "" {
		query += fmt.Sprintf(" AND c.type = $%d", argID)
		args = append(args, certType)
		argID++
	}
	if year != "" {
		query += fmt.Sprintf(" AND c.academic_year = $%d", argID)
		args = append(args, year)
		argID++
	}

	query += " ORDER BY c.issued_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Certificate
	for rows.Next() {
		var c Certificate
		err := rows.Scan(
			&c.ID, &c.SchoolID, &c.StudentID, &c.StudentName, &c.ClassName, &c.Type,
			&c.IssuedBy, &c.IssuedByName, &c.IssuedAt, &c.AcademicYear, &c.Notes, &c.PDFUrl, &c.ProtocolNo, &c.IsDeleted,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *repository) SoftDelete(ctx context.Context, id string) error {
	query := `UPDATE certificates SET is_deleted = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) NextProtocolNo(ctx context.Context, schoolID string) (string, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM certificates WHERE school_id = $1`, schoolID).Scan(&count)
	if err != nil {
		count = 0
	}
	year := time.Now().Year()
	return fmt.Sprintf("CERT-%d-%05d", year, count+1), nil
}
