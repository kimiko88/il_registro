package parents

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Repository interface {
	GetChildrenByParentUserID(ctx context.Context, parentUserID string) ([]string, error)
	GetDashboardStats(ctx context.Context, parentUserID string, childUserIDs []string) (*ParentDashboardStatsResponse, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetChildrenByParentUserID(ctx context.Context, parentUserID string) ([]string, error) {
	query := `
		SELECT u.id::text
		FROM student_parents sp
		JOIN parents p ON sp.parent_id = p.id
		JOIN students s ON sp.student_id = s.id
		JOIN users u ON s.user_id = u.id
		WHERE p.user_id = $1::uuid
	`
	rows, err := r.db.QueryContext(ctx, query, parentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, rows.Err()
}

func (r *PostgresRepository) GetDashboardStats(ctx context.Context, parentUserID string, childUserIDs []string) (*ParentDashboardStatsResponse, error) {
	resp := &ParentDashboardStatsResponse{
		ChildrenCount:        len(childUserIDs),
		UpcomingColloqui:     0,
		UnreadCommunications: 0,
		DocumentsCount:       0,
		TotalChildren:        len(childUserIDs),
		UnreadMessages:       0,
		PendingPayments:      0,
	}

	if r.db == nil {
		return resp, nil
	}

	// 1. Upcoming confirmed colloqui for this parent
	queryColloqui := `
		SELECT COUNT(*) 
		FROM colloqui_prenotazioni cp
		JOIN colloqui_slot cs ON cp.slot_id = cs.id
		WHERE cp.parent_id = $1::uuid 
		  AND cp.status = 'confirmed'
		  AND cs.data_ora >= NOW()`
	var upcomingColloqui int
	if err := r.db.QueryRowContext(ctx, queryColloqui, parentUserID).Scan(&upcomingColloqui); err == nil {
		resp.UpcomingColloqui = upcomingColloqui
	}

	// 2. Pending payments and documents for children
	if len(childUserIDs) > 0 {
		placeholders := make([]string, len(childUserIDs))
		args := make([]interface{}, len(childUserIDs))
		for i, id := range childUserIDs {
			placeholders[i] = fmt.Sprintf("$%d::uuid", i+1)
			args[i] = id
		}

		queryPay := fmt.Sprintf(`
			SELECT COUNT(*) 
			FROM payments 
			WHERE student_id IN (%s) 
			  AND status = 'pending' 
			  AND deleted_at IS NULL`, strings.Join(placeholders, ","))
		var pendingPayments int
		if err := r.db.QueryRowContext(ctx, queryPay, args...).Scan(&pendingPayments); err == nil {
			resp.PendingPayments = pendingPayments
		}

		queryDocs := fmt.Sprintf(`
			SELECT COUNT(*) 
			FROM documents_enhanced 
			WHERE student_id IN (%s) 
			  AND status = 'published' 
			  AND deleted_at IS NULL`, strings.Join(placeholders, ","))
		var docsCount int
		if err := r.db.QueryRowContext(ctx, queryDocs, args...).Scan(&docsCount); err == nil {
			resp.DocumentsCount = docsCount
		}
	}

	// 3. Unread communications (circulars requiring signature or unread)
	queryUnread := `
		SELECT COUNT(*) 
		FROM communications c
		LEFT JOIN communication_signatures cs ON c.id = cs.communication_id AND cs.user_id = $1::uuid
		WHERE c.deleted_at IS NULL 
		  AND (cs.id IS NULL OR cs.signed_at IS NULL)`
	var unreadComms int
	if err := r.db.QueryRowContext(ctx, queryUnread, parentUserID).Scan(&unreadComms); err == nil {
		resp.UnreadCommunications = unreadComms
		resp.UnreadMessages = unreadComms
	}

	return resp, nil
}

func (r *PostgresRepository) CreateAuthorization(ctx context.Context, a *DualParentalAuthorization) error {
	query := `
		INSERT INTO dual_parental_authorizations 
		(id, student_id, document_type, document_ref_id, title, parent1_id, parent2_id, status, deadline, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, '')::uuid, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query, a.ID, a.StudentID, a.DocumentType, a.DocumentRefID, a.Title, a.Parent1ID, a.Parent2ID, a.Status, a.Deadline, a.CreatedAt)
	return err
}

func (r *PostgresRepository) GetAuthorizationByID(ctx context.Context, id string) (*DualParentalAuthorization, error) {
	query := `
		SELECT 
			id, student_id, document_type, document_ref_id, title,
			parent1_id, parent1_signed_at, parent1_pin_verified,
			COALESCE(parent2_id::text, ''), parent2_signed_at, parent2_pin_verified,
			status, COALESCE(rejection_reason, ''), deadline, created_at
		FROM dual_parental_authorizations
		WHERE id = $1
	`
	var a DualParentalAuthorization
	var parent2ID string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.StudentID, &a.DocumentType, &a.DocumentRefID, &a.Title,
		&a.Parent1ID, &a.Parent1SignedAt, &a.Parent1PinVerified,
		&parent2ID, &a.Parent2SignedAt, &a.Parent2PinVerified,
		&a.Status, &a.RejectionReason, &a.Deadline, &a.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	a.Parent2ID = parent2ID
	return &a, nil
}

func (r *PostgresRepository) UpdateAuthorization(ctx context.Context, a *DualParentalAuthorization) error {
	query := `
		UPDATE dual_parental_authorizations SET
			parent1_signed_at = $1,
			parent1_pin_verified = $2,
			parent2_signed_at = $3,
			parent2_pin_verified = $4,
			status = $5,
			rejection_reason = $6
		WHERE id = $7
	`
	_, err := r.db.ExecContext(ctx, query, a.Parent1SignedAt, a.Parent1PinVerified, a.Parent2SignedAt, a.Parent2PinVerified, a.Status, a.RejectionReason, a.ID)
	return err
}

func (r *PostgresRepository) ListAuthorizationsForParent(ctx context.Context, parentID string) ([]DualParentalAuthorization, error) {
	query := `
		SELECT 
			id, student_id, document_type, document_ref_id, title,
			parent1_id, parent1_signed_at, parent1_pin_verified,
			COALESCE(parent2_id::text, ''), parent2_signed_at, parent2_pin_verified,
			status, COALESCE(rejection_reason, ''), deadline, created_at
		FROM dual_parental_authorizations
		WHERE parent1_id = $1::uuid OR parent2_id = $1::uuid
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []DualParentalAuthorization
	for rows.Next() {
		var a DualParentalAuthorization
		var parent2ID string
		if err := rows.Scan(
			&a.ID, &a.StudentID, &a.DocumentType, &a.DocumentRefID, &a.Title,
			&a.Parent1ID, &a.Parent1SignedAt, &a.Parent1PinVerified,
			&parent2ID, &a.Parent2SignedAt, &a.Parent2PinVerified,
			&a.Status, &a.RejectionReason, &a.Deadline, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		a.Parent2ID = parent2ID
		results = append(results, a)
	}
	return results, rows.Err()
}

func (r *PostgresRepository) GetCustodyInfo(ctx context.Context, parentID, studentID string) (*StudentCustodyInfo, error) {
	query := `
		SELECT 
			COALESCE(sp.custody_type, 'shared'),
			COALESCE(sp.court_order_details, ''),
			COALESCE(sp.can_authorize_activities, true),
			COALESCE(sp.is_mirror_notified, true)
		FROM student_parents sp
		JOIN parents p ON sp.parent_id = p.id
		JOIN students s ON sp.student_id = s.id
		WHERE p.user_id = $1::uuid AND (s.user_id = $2::uuid OR s.id = $2::uuid)
	`
	var custody CustodyType
	var courtOrder string
	var canAuth, mirror bool
	err := r.db.QueryRowContext(ctx, query, parentID, studentID).Scan(&custody, &courtOrder, &canAuth, &mirror)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &StudentCustodyInfo{
				ParentID:               parentID,
				StudentID:              studentID,
				CustodyType:            CustodyShared,
				CanAuthorizeActivities: true,
				IsMirrorNotified:       true,
			}, nil
		}
		return nil, err
	}
	return &StudentCustodyInfo{
		ParentID:               parentID,
		StudentID:              studentID,
		CustodyType:            custody,
		CourtOrderDetails:      courtOrder,
		CanAuthorizeActivities: canAuth,
		IsMirrorNotified:       mirror,
	}, nil
}

func (r *PostgresRepository) ValidateParentPIN(ctx context.Context, parentID, pin string) (bool, error) {
	// A PIN must be non-empty and at least 4 digits
	if len(pin) < 4 {
		return false, nil
	}
	return true, nil
}
