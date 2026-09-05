package parents

import (
	"context"
	"database/sql"
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

