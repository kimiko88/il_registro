package search

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	GlobalSearch(ctx context.Context, schoolID, query, filterType string) ([]SearchResultItem, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GlobalSearch(ctx context.Context, schoolID, q, filterType string) ([]SearchResultItem, error) {
	var results []SearchResultItem
	searchPattern := "%" + q + "%"

	// 1. Search Users (Students / Teachers)
	if filterType == "" || filterType == "users" || filterType == "students" || filterType == "teachers" {
		userQuery := `
			SELECT id::text, role, (first_name || ' ' || last_name) AS name, email, COALESCE(fiscal_code, '')
			FROM users
			WHERE (school_id = $1::uuid OR $1 = '')
			  AND (first_name ILIKE $2 OR last_name ILIKE $2 OR email ILIKE $2 OR fiscal_code ILIKE $2)
			LIMIT 20
		`
		rows, err := r.db.QueryContext(ctx, userQuery, schoolID, searchPattern)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id, role, name, email, cf string
				if err := rows.Scan(&id, &role, &name, &email, &cf); err == nil {
					results = append(results, SearchResultItem{
						ID:          id,
						Type:        role,
						Title:       name,
						Subtitle:    email,
						Description: fmt.Sprintf("Ruolo: %s | Codice Fiscale: %s", role, cf),
					})
				}
			}
		}
	}

	// 2. Search Communications
	if filterType == "" || filterType == "communications" {
		commQuery := `
			SELECT id::text, subject, type, body
			FROM communications
			WHERE (school_id = $1::uuid OR $1 = '')
			  AND (subject ILIKE $2 OR body ILIKE $2)
			LIMIT 20
		`
		rows, err := r.db.QueryContext(ctx, commQuery, schoolID, searchPattern)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id, subject, ctype, body string
				if err := rows.Scan(&id, &subject, &ctype, &body); err == nil {
					if len(body) > 100 {
						body = body[:100] + "..."
					}
					results = append(results, SearchResultItem{
						ID:          id,
						Type:        "communication",
						Title:       subject,
						Subtitle:    fmt.Sprintf("Tipo: %s", ctype),
						Description: body,
					})
				}
			}
		}
	}

	// 3. Search Lessons
	if filterType == "" || filterType == "lessons" {
		lessonQuery := `
			SELECT id::text, topic, COALESCE(notes, ''), type
			FROM class_lessons
			WHERE topic ILIKE $1 OR notes ILIKE $1
			LIMIT 20
		`
		rows, err := r.db.QueryContext(ctx, lessonQuery, searchPattern)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id, topic, notes, ltype string
				if err := rows.Scan(&id, &topic, &notes, &ltype); err == nil {
					results = append(results, SearchResultItem{
						ID:          id,
						Type:        "lesson",
						Title:       topic,
						Subtitle:    fmt.Sprintf("Tipo Lezione: %s", ltype),
						Description: notes,
					})
				}
			}
		}
	}

	return results, nil
}
