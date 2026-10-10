package family_desk

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

type Repository interface {
	CreateRequest(ctx context.Context, req *FamilyRequest) error
	GetRequest(ctx context.Context, id string) (*FamilyRequest, error)
	ListRequests(ctx context.Context, schoolID, parentID, status string) ([]FamilyRequest, error)
	UpdateRequestStatus(ctx context.Context, id, status, rejectionReason, protocolNumber, reviewerID string) error
	CreatePermanentDelegate(ctx context.Context, delegate *PermanentDelegate) error
	ListPermanentDelegates(ctx context.Context, studentID, schoolID string) ([]PermanentDelegate, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateRequest(ctx context.Context, req *FamilyRequest) error {
	formDataJSON, err := json.Marshal(req.FormData)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO family_requests (
			school_id, student_id, parent_id, request_type, form_data, attachment_urls, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at;
	`
	return r.db.QueryRowContext(ctx, query,
		req.SchoolID, req.StudentID, req.ParentID, req.RequestType, formDataJSON, pq.Array(req.AttachmentURLs), req.Status,
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *postgresRepo) GetRequest(ctx context.Context, id string) (*FamilyRequest, error) {
	query := `
		SELECT id, school_id, student_id, parent_id, request_type, form_data, attachment_urls, status,
		       COALESCE(rejection_reason, ''), COALESCE(protocol_number, ''), reviewed_by, reviewed_at, created_at, updated_at
		FROM family_requests
		WHERE id = $1;
	`
	var req FamilyRequest
	var formDataRaw []byte
	var reviewedBy sql.NullString
	var reviewedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&req.ID, &req.SchoolID, &req.StudentID, &req.ParentID, &req.RequestType, &formDataRaw, pq.Array(&req.AttachmentURLs),
		&req.Status, &req.RejectionReason, &req.ProtocolNumber, &reviewedBy, &reviewedAt, &req.CreatedAt, &req.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(formDataRaw, &req.FormData)
	if reviewedBy.Valid {
		req.ReviewedBy = &reviewedBy.String
	}
	if reviewedAt.Valid {
		req.ReviewedAt = &reviewedAt.Time
	}

	return &req, nil
}

func (r *postgresRepo) ListRequests(ctx context.Context, schoolID, parentID, status string) ([]FamilyRequest, error) {
	query := `
		SELECT id, school_id, student_id, parent_id, request_type, form_data, attachment_urls, status,
		       COALESCE(rejection_reason, ''), COALESCE(protocol_number, ''), reviewed_by, reviewed_at, created_at, updated_at
		FROM family_requests
		WHERE school_id = $1
		  AND ($2 = '' OR parent_id = $2::uuid)
		  AND ($3 = '' OR status = $3)
		ORDER BY created_at DESC;
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, parentID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []FamilyRequest
	for rows.Next() {
		var req FamilyRequest
		var formDataRaw []byte
		var reviewedBy sql.NullString
		var reviewedAt sql.NullTime

		if err := rows.Scan(
			&req.ID, &req.SchoolID, &req.StudentID, &req.ParentID, &req.RequestType, &formDataRaw, pq.Array(&req.AttachmentURLs),
			&req.Status, &req.RejectionReason, &req.ProtocolNumber, &reviewedBy, &reviewedAt, &req.CreatedAt, &req.UpdatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(formDataRaw, &req.FormData)
		if reviewedBy.Valid {
			req.ReviewedBy = &reviewedBy.String
		}
		if reviewedAt.Valid {
			req.ReviewedAt = &reviewedAt.Time
		}
		list = append(list, req)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *postgresRepo) UpdateRequestStatus(ctx context.Context, id, status, rejectionReason, protocolNumber, reviewerID string) error {
	now := time.Now()
	query := `
		UPDATE family_requests
		SET status = $2, rejection_reason = $3, protocol_number = $4, reviewed_by = $5, reviewed_at = $6, updated_at = $6
		WHERE id = $1;
	`
	_, err := r.db.ExecContext(ctx, query, id, status, rejectionReason, protocolNumber, reviewerID, now)
	return err
}

func (r *postgresRepo) CreatePermanentDelegate(ctx context.Context, delegate *PermanentDelegate) error {
	query := `
		INSERT INTO student_permanent_delegates (
			student_id, school_id, first_name, last_name, tax_code, relationship, phone, id_card_details, id_card_file_url, is_valid, approved_request_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at;
	`
	return r.db.QueryRowContext(ctx, query,
		delegate.StudentID, delegate.SchoolID, delegate.FirstName, delegate.LastName, delegate.TaxCode,
		delegate.Relationship, delegate.Phone, delegate.IDCardDetails, delegate.IDCardFileURL, delegate.IsValid, delegate.ApprovedRequestID,
	).Scan(&delegate.ID, &delegate.CreatedAt)
}

func (r *postgresRepo) ListPermanentDelegates(ctx context.Context, studentID, schoolID string) ([]PermanentDelegate, error) {
	query := `
		SELECT id, student_id, school_id, first_name, last_name, tax_code, relationship, phone, id_card_details, COALESCE(id_card_file_url, ''), is_valid, approved_request_id, created_at
		FROM student_permanent_delegates
		WHERE ($1 = '' OR student_id = $1::uuid)
		  AND ($2 = '' OR school_id = $2::uuid)
		  AND is_valid = true
		ORDER BY last_name ASC;
	`
	rows, err := r.db.QueryContext(ctx, query, studentID, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []PermanentDelegate
	for rows.Next() {
		var d PermanentDelegate
		var appID sql.NullString
		if err := rows.Scan(
			&d.ID, &d.StudentID, &d.SchoolID, &d.FirstName, &d.LastName, &d.TaxCode, &d.Relationship, &d.Phone,
			&d.IDCardDetails, &d.IDCardFileURL, &d.IsValid, &appID, &d.CreatedAt,
		); err != nil {
			return nil, err
		}
		if appID.Valid {
			d.ApprovedRequestID = &appID.String
		}
		list = append(list, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
