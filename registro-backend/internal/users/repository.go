package users

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	FindAll() ([]User, error)
	FindByID(id string) (*User, error)
	Create(user *User) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]User, error) {
	query := `
		SELECT id, first_name, last_name, email, role, school_id, is_active, email_verified, created_at, updated_at 
		FROM users 
		WHERE deleted_at IS NULL
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role, &u.SchoolID,
			&u.IsActive, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *repository) FindByID(id string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, role, school_id, is_active, email_verified, created_at, updated_at 
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL
	`
	row := r.db.QueryRow(query, id)

	u := &User{}
	err := row.Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role, &u.SchoolID,
		&u.IsActive, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	return u, nil
}

func (r *repository) Create(user *User) error {
	query := `
		INSERT INTO users (
			first_name, last_name, email, role, school_id, 
			created_at, updated_at, is_active, email_verified
		) 
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), true, false) 
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.FirstName, user.LastName, user.Email, user.Role, user.SchoolID,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}
	return nil
}
