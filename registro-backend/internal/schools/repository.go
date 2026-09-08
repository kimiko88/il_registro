package schools

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostgresRepository handles database operations for schools
type PostgresRepository struct {
	db *sql.DB
}

// NewRepository creates a new schools repository
func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create inserts a new school
func (r *PostgresRepository) Create(ctx context.Context, school *School) error {
	school.ID = uuid.New().String()
	school.CreatedAt = time.Now()
	school.UpdatedAt = time.Now()

	query := `
		INSERT INTO schools (id, name, code, address, city, phone, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		school.ID, school.Name, school.Code, school.Address, school.City,
		school.Phone, school.Email, school.CreatedAt, school.UpdatedAt,
	)
	return err
}

// GetByID retrieves a school by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*School, error) {
	query := `SELECT id, name, COALESCE(code, ''), address, COALESCE(city, ''), phone, email, created_at, updated_at FROM schools WHERE id = $1 AND deleted_at IS NULL`
	school := &School{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&school.ID, &school.Name, &school.Code, &school.Address, &school.City,
		&school.Phone, &school.Email, &school.CreatedAt, &school.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return school, err
}

// List retrieves schools with pagination
func (r *PostgresRepository) List(ctx context.Context, params *ListParams) ([]*School, int, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 {
		pageSize = params.Limit
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (params.Page - 1) * pageSize

	// Count query
	countQuery := `SELECT COUNT(*) FROM schools WHERE deleted_at IS NULL`
	args := []interface{}{}
	argIndex := 1

	if params.Search != "" {
		countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR code ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+params.Search+"%")
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// List query
	listQuery := `SELECT id, name, COALESCE(code, ''), COALESCE(address, ''), COALESCE(city, ''), COALESCE(phone, ''), COALESCE(email, ''), created_at, updated_at FROM schools WHERE deleted_at IS NULL`
	if params.Search != "" {
		listQuery += " AND (name ILIKE $1 OR code ILIKE $1)"
	}
	listQuery += fmt.Sprintf(" ORDER BY name LIMIT %d OFFSET %d", pageSize, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var schools []*School
	for rows.Next() {
		s := &School{}
		if err := rows.Scan(
			&s.ID, &s.Name, &s.Code, &s.Address, &s.City,
			&s.Phone, &s.Email, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		schools = append(schools, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return schools, total, nil
}

// Update updates a school
func (r *PostgresRepository) Update(ctx context.Context, id string, req *UpdateSchoolRequest) error {
	query := `UPDATE schools SET updated_at = $1`
	args := []interface{}{time.Now()}
	argIndex := 2

	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argIndex)
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Code != nil {
		query += fmt.Sprintf(", code = $%d", argIndex)
		args = append(args, *req.Code)
		argIndex++
	}
	if req.Address != nil {
		query += fmt.Sprintf(", address = $%d", argIndex)
		args = append(args, *req.Address)
		argIndex++
	}
	if req.City != nil {
		query += fmt.Sprintf(", city = $%d", argIndex)
		args = append(args, *req.City)
		argIndex++
	}
	if req.Phone != nil {
		query += fmt.Sprintf(", phone = $%d", argIndex)
		args = append(args, *req.Phone)
		argIndex++
	}
	if req.Email != nil {
		query += fmt.Sprintf(", email = $%d", argIndex)
		args = append(args, *req.Email)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", argIndex)
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("school not found")
	}
	return nil
}

// Delete soft-deletes a school
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE schools SET deleted_at = $1, updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
