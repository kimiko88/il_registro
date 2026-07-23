package tenants

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateTenant(ctx context.Context, t *Tenant) error
	GetTenantByID(ctx context.Context, id string) (*Tenant, error)
	ListTenants(ctx context.Context) ([]*Tenant, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateTenant(ctx context.Context, t *Tenant) error {
	t.ID = uuid.New().String()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	t.Status = "active"

	query := `
		INSERT INTO schools (id, name, code, is_active, created_at, updated_at)
		VALUES ($1::uuid, $2, $3, true, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, t.ID, t.Name, t.Code, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetTenantByID(ctx context.Context, id string) (*Tenant, error) {
	query := `SELECT id, name, COALESCE(code, ''), created_at FROM schools WHERE id = $1::uuid`
	t := &Tenant{
		Quota: TenantQuota{MaxStudents: 1000, MaxTeachers: 100, MaxStorageMB: 10240},
		Status: "active",
	}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Name, &t.Code, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *PostgresRepository) ListTenants(ctx context.Context) ([]*Tenant, error) {
	query := `SELECT id, name, COALESCE(code, ''), created_at FROM schools ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Tenant
	for rows.Next() {
		t := &Tenant{
			Quota: TenantQuota{MaxStudents: 1000, MaxTeachers: 100, MaxStorageMB: 10240},
			Status: "active",
		}
		if err := rows.Scan(&t.ID, &t.Name, &t.Code, &t.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}
