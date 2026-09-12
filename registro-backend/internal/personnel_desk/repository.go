package personnel_desk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, r *DeskRequest) error
	GetByID(ctx context.Context, schoolID, id string) (*DeskRequest, error)
	List(ctx context.Context, schoolID, applicantID, status string) ([]*DeskRequest, error)
	Update(ctx context.Context, r *DeskRequest) error
	Delete(ctx context.Context, schoolID, id, applicantID string) error
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, req *DeskRequest) error {
	req.ID = uuid.New().String()
	now := time.Now()
	req.CreatedAt = now
	req.UpdatedAt = now

	if req.SchoolID == "" {
		_ = r.db.QueryRowContext(ctx, `SELECT COALESCE(school_id::text, '') FROM users WHERE id=$1`, req.ApplicantID).Scan(&req.SchoolID)
	}

	attJSON, err := json.Marshal(req.Attachments)
	if err != nil {
		attJSON = []byte("[]")
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO personnel_desk_requests (
			id, school_id, applicant_id, category, sub_category,
			start_date, end_date, days, hours, description, attachments,
			status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11,
			$12, $13, $14
		)`,
		req.ID, req.SchoolID, req.ApplicantID, req.Category, req.SubCategory,
		req.StartDate, req.EndDate, req.Days, req.Hours, req.Description, string(attJSON),
		req.Status, req.CreatedAt, req.UpdatedAt,
	)
	return err
}

func (r *postgresRepository) GetByID(ctx context.Context, schoolID, id string) (*DeskRequest, error) {
	query := `
		SELECT pdr.id, pdr.school_id, pdr.applicant_id,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS applicant_name,
		       COALESCE(u.role, '') AS applicant_role,
		       pdr.category, COALESCE(pdr.sub_category, ''),
		       pdr.start_date::text, pdr.end_date::text, COALESCE(pdr.days, 0)::float8, COALESCE(pdr.hours, 0)::float8,
		       COALESCE(pdr.description, ''), COALESCE(pdr.attachments::text, '[]'),
		       pdr.status,
		       COALESCE(pdr.aa_note, ''), pdr.aa_reviewed_by, pdr.aa_reviewed_at,
		       COALESCE(pdr.dsga_note, ''), pdr.dsga_signed_by, pdr.dsga_signed_at,
		       pdr.ds_decree_num, COALESCE(pdr.ds_note, ''), pdr.ds_approved_by, pdr.ds_approved_at,
		       pdr.created_at, pdr.updated_at
		FROM personnel_desk_requests pdr
		LEFT JOIN users u ON u.id = pdr.applicant_id
		WHERE pdr.id = $1 AND (pdr.school_id = $2 OR $2 = '')`

	row := r.db.QueryRowContext(ctx, query, id, schoolID)

	req := &DeskRequest{}
	var attStr string
	var aaRevBy, dsgaSignBy, dsDecree, dsAppBy sql.NullString
	var aaRevAt, dsgaSignAt, dsAppAt sql.NullTime

	err := row.Scan(
		&req.ID, &req.SchoolID, &req.ApplicantID,
		&req.ApplicantName, &req.ApplicantRole,
		&req.Category, &req.SubCategory,
		&req.StartDate, &req.EndDate, &req.Days, &req.Hours,
		&req.Description, &attStr,
		&req.Status,
		&req.AANote, &aaRevBy, &aaRevAt,
		&req.DSGANote, &dsgaSignBy, &dsgaSignAt,
		&dsDecree, &req.DSNote, &dsAppBy, &dsAppAt,
		&req.CreatedAt, &req.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("richiesta non trovata")
	}
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal([]byte(attStr), &req.Attachments)
	if aaRevBy.Valid {
		req.AAReviewedBy = &aaRevBy.String
	}
	if aaRevAt.Valid {
		req.AAReviewedAt = &aaRevAt.Time
	}
	if dsgaSignBy.Valid {
		req.DSGASignedBy = &dsgaSignBy.String
	}
	if dsgaSignAt.Valid {
		req.DSGASignedAt = &dsgaSignAt.Time
	}
	if dsDecree.Valid {
		req.DSDecreeNum = &dsDecree.String
	}
	if dsAppBy.Valid {
		req.DSApprovedBy = &dsAppBy.String
	}
	if dsAppAt.Valid {
		req.DSApprovedAt = &dsAppAt.Time
	}

	return req, nil
}

func (r *postgresRepository) List(ctx context.Context, schoolID, applicantID, status string) ([]*DeskRequest, error) {
	query := `
		SELECT pdr.id, pdr.school_id, pdr.applicant_id,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS applicant_name,
		       COALESCE(u.role, '') AS applicant_role,
		       pdr.category, COALESCE(pdr.sub_category, ''),
		       pdr.start_date::text, pdr.end_date::text, COALESCE(pdr.days, 0)::float8, COALESCE(pdr.hours, 0)::float8,
		       COALESCE(pdr.description, ''), COALESCE(pdr.attachments::text, '[]'),
		       pdr.status,
		       COALESCE(pdr.aa_note, ''), pdr.aa_reviewed_by, pdr.aa_reviewed_at,
		       COALESCE(pdr.dsga_note, ''), pdr.dsga_signed_by, pdr.dsga_signed_at,
		       pdr.ds_decree_num, COALESCE(pdr.ds_note, ''), pdr.ds_approved_by, pdr.ds_approved_at,
		       pdr.created_at, pdr.updated_at
		FROM personnel_desk_requests pdr
		LEFT JOIN users u ON u.id = pdr.applicant_id
		WHERE (pdr.school_id = $1 OR $1 = '')`

	args := []interface{}{schoolID}
	idx := 2

	if applicantID != "" {
		query += fmt.Sprintf(" AND pdr.applicant_id = $%d", idx)
		args = append(args, applicantID)
		idx++
	}
	if status != "" {
		query += fmt.Sprintf(" AND pdr.status = $%d", idx)
		args = append(args, status)
		idx++
	}

	query += " ORDER BY pdr.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*DeskRequest
	for rows.Next() {
		req := &DeskRequest{}
		var attStr string
		var aaRevBy, dsgaSignBy, dsDecree, dsAppBy sql.NullString
		var aaRevAt, dsgaSignAt, dsAppAt sql.NullTime

		if err := rows.Scan(
			&req.ID, &req.SchoolID, &req.ApplicantID,
			&req.ApplicantName, &req.ApplicantRole,
			&req.Category, &req.SubCategory,
			&req.StartDate, &req.EndDate, &req.Days, &req.Hours,
			&req.Description, &attStr,
			&req.Status,
			&req.AANote, &aaRevBy, &aaRevAt,
			&req.DSGANote, &dsgaSignBy, &dsgaSignAt,
			&dsDecree, &req.DSNote, &dsAppBy, &dsAppAt,
			&req.CreatedAt, &req.UpdatedAt,
		); err != nil {
			return nil, err
		}

		_ = json.Unmarshal([]byte(attStr), &req.Attachments)
		if aaRevBy.Valid {
			req.AAReviewedBy = &aaRevBy.String
		}
		if aaRevAt.Valid {
			req.AAReviewedAt = &aaRevAt.Time
		}
		if dsgaSignBy.Valid {
			req.DSGASignedBy = &dsgaSignBy.String
		}
		if dsgaSignAt.Valid {
			req.DSGASignedAt = &dsgaSignAt.Time
		}
		if dsDecree.Valid {
			req.DSDecreeNum = &dsDecree.String
		}
		if dsAppBy.Valid {
			req.DSApprovedBy = &dsAppBy.String
		}
		if dsAppAt.Valid {
			req.DSApprovedAt = &dsAppAt.Time
		}

		list = append(list, req)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *postgresRepository) Update(ctx context.Context, req *DeskRequest) error {
	req.UpdatedAt = time.Now()
	attJSON, err := json.Marshal(req.Attachments)
	if err != nil {
		attJSON = []byte("[]")
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE personnel_desk_requests SET
			category = $1,
			sub_category = $2,
			start_date = $3,
			end_date = $4,
			days = $5,
			hours = $6,
			description = $7,
			attachments = $8,
			status = $9,
			aa_note = $10,
			aa_reviewed_by = $11,
			aa_reviewed_at = $12,
			dsga_note = $13,
			dsga_signed_by = $14,
			dsga_signed_at = $15,
			ds_decree_num = $16,
			ds_note = $17,
			ds_approved_by = $18,
			ds_approved_at = $19,
			updated_at = $20
		WHERE id = $21 AND (school_id = $22 OR $22 = '')`,
		req.Category, req.SubCategory,
		req.StartDate, req.EndDate, req.Days, req.Hours,
		req.Description, string(attJSON),
		req.Status,
		req.AANote, req.AAReviewedBy, req.AAReviewedAt,
		req.DSGANote, req.DSGASignedBy, req.DSGASignedAt,
		req.DSDecreeNum, req.DSNote, req.DSApprovedBy, req.DSApprovedAt,
		req.UpdatedAt,
		req.ID, req.SchoolID,
	)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, schoolID, id, applicantID string) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM personnel_desk_requests
		WHERE id = $1 AND (school_id = $2 OR $2 = '') AND applicant_id = $3 AND status IN ('draft', 'submitted')`,
		id, schoolID, applicantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("richiesta non eliminabile o non trovata")
	}
	return nil
}
