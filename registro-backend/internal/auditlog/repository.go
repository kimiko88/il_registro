package auditlog

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	Insert(ctx context.Context, event *AuditEvent) error
	InsertAsync(event AuditEvent)
	List(ctx context.Context, p FilterParams) ([]AuditEvent, int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	r := &repository{db: db}
	r.ensureTable()
	return r
}

func (r *repository) ensureTable() {
	query := `
		CREATE TABLE IF NOT EXISTS audit_logs (
			id UUID PRIMARY KEY,
			school_id TEXT,
			actor_id TEXT,
			actor_role TEXT,
			actor_name TEXT,
			action TEXT,
			entity_type TEXT,
			entity_id TEXT,
			details TEXT,
			ip_address TEXT,
			user_agent TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_audit_logs_school_actor ON audit_logs(school_id, actor_id);
		CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
		CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
	`
	if r.db != nil {
		_, err := r.db.Exec(query)
		if err != nil {
			logrus.Warnf("audit_logs table init check warning: %v", err)
		}
	}
}

func (r *repository) Insert(ctx context.Context, event *AuditEvent) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO audit_logs (id, school_id, actor_id, actor_role, actor_name, action, entity_type, entity_id, details, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		event.ID, event.SchoolID, event.ActorID, event.ActorRole, event.ActorName,
		event.Action, event.EntityType, event.EntityID, event.Details, event.IPAddress, event.UserAgent, event.CreatedAt,
	)
	return err
}

func (r *repository) InsertAsync(event AuditEvent) {
	go func(ev AuditEvent) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := r.Insert(ctx, &ev); err != nil {
			logrus.Errorf("Async audit log insert error: %v", err)
		}
	}(event)
}

func (r *repository) List(ctx context.Context, p FilterParams) ([]AuditEvent, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	argID := 1

	if p.SchoolID != "" {
		where += fmt.Sprintf(" AND school_id = $%d", argID)
		args = append(args, p.SchoolID)
		argID++
	}
	if p.ActorID != "" {
		where += fmt.Sprintf(" AND (actor_id ILIKE $%d OR actor_name ILIKE $%d)", argID, argID)
		args = append(args, "%"+p.ActorID+"%")
		argID++
	}
	if p.Action != "" {
		where += fmt.Sprintf(" AND action ILIKE $%d", argID)
		args = append(args, "%"+p.Action+"%")
		argID++
	}
	if p.EntityType != "" {
		where += fmt.Sprintf(" AND entity_type ILIKE $%d", argID)
		args = append(args, "%"+p.EntityType+"%")
		argID++
	}
	if p.IPAddress != "" {
		where += fmt.Sprintf(" AND ip_address ILIKE $%d", argID)
		args = append(args, "%"+p.IPAddress+"%")
		argID++
	}
	if p.From != "" {
		where += fmt.Sprintf(" AND created_at >= $%d::timestamp", argID)
		args = append(args, p.From)
		argID++
	}
	if p.To != "" {
		where += fmt.Sprintf(" AND created_at <= $%d::timestamp", argID)
		args = append(args, p.To+" 23:59:59")
		argID++
	}

	countQuery := "SELECT COUNT(*) FROM audit_logs " + where
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		total = 0
	}

	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Page <= 0 {
		p.Page = 1
	}
	offset := (p.Page - 1) * p.Limit

	listQuery := fmt.Sprintf("SELECT id, school_id, actor_id, actor_role, actor_name, action, entity_type, entity_id, details, ip_address, user_agent, created_at FROM audit_logs %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", where, argID, argID+1)
	args = append(args, p.Limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []AuditEvent
	for rows.Next() {
		var ev AuditEvent
		err := rows.Scan(
			&ev.ID, &ev.SchoolID, &ev.ActorID, &ev.ActorRole, &ev.ActorName,
			&ev.Action, &ev.EntityType, &ev.EntityID, &ev.Details, &ev.IPAddress, &ev.UserAgent, &ev.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
