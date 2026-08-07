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

var validSearchFilterTypes = map[string]bool{
	"":               true,
	"all":            true,
	"users":          true,
	"students":       true,
	"teachers":       true,
	"classes":        true,
	"communications": true,
	"lessons":        true,
}

func (r *PostgresRepository) GlobalSearch(ctx context.Context, schoolID, q, filterType string) ([]SearchResultItem, error) {
	var results []SearchResultItem
	searchPattern := "%" + q + "%"

	// 1. Search Users (Students / Teachers / Secretary / Admin)
	if filterType == "" || filterType == "all" || filterType == "users" || filterType == "students" || filterType == "teachers" {
		userQuery := `
			SELECT id::text, role, (first_name || ' ' || last_name) AS name, email, COALESCE(fiscal_code, '')
			FROM users
			WHERE (school_id = $1::uuid OR $1 = '')
			  AND (first_name ILIKE $2 OR last_name ILIKE $2 OR email ILIKE $2 OR fiscal_code ILIKE $2 OR role ILIKE $2)
			LIMIT 20
		`
		rows, err := r.db.QueryContext(ctx, userQuery, schoolID, searchPattern)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id, role, name, email, cf string
				if err := rows.Scan(&id, &role, &name, &email, &cf); err == nil {
					roleLabel := "Utente"
					switch role {
					case "student":
						roleLabel = "Studente"
					case "teacher":
						roleLabel = "Docente"
					case "secretary":
						roleLabel = "Personale di Segreteria"
					case "admin":
						roleLabel = "Amministratore"
					case "parent":
						roleLabel = "Genitore"
					}
					results = append(results, SearchResultItem{
						ID:          id,
						Type:        role,
						Title:       name,
						Subtitle:    email,
						Description: fmt.Sprintf("Ruolo: %s | CF: %s", roleLabel, cf),
					})
				}
			}
		}
	}

	// 2. Search Classes
	if filterType == "" || filterType == "all" || filterType == "classes" {
		classQuery := `
			SELECT id::text, COALESCE(section, 'A'), COALESCE(specialization, 'Generale')
			FROM classes
			WHERE (school_id = $1::uuid OR $1 = '')
			  AND (section ILIKE $2 OR specialization ILIKE $2)
			LIMIT 10
		`
		rows, err := r.db.QueryContext(ctx, classQuery, schoolID, searchPattern)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id, section, spec string
				if err := rows.Scan(&id, &section, &spec); err == nil {
					results = append(results, SearchResultItem{
						ID:          id,
						Type:        "class",
						Title:       "Classe " + section,
						Subtitle:    spec,
						Description: "Gestione Classe e Studenti",
					})
				}
			}
		}
	}

	// 3. Search Communications
	if filterType == "" || filterType == "all" || filterType == "communications" {
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
						Subtitle:    fmt.Sprintf("Circolare / Avviso (%s)", ctype),
						Description: body,
					})
				}
			}
		}
	}

	// 4. Search Lessons
	if filterType == "" || filterType == "all" || filterType == "lessons" {
		lessonQuery := `
			SELECT cl.id::text, cl.topic, COALESCE(cl.notes, ''), cl.type
			FROM class_lessons cl
			JOIN classes c ON cl.class_id = c.id
			WHERE (c.school_id = $1::uuid OR $1 = '')
			  AND (cl.topic ILIKE $2 OR cl.notes ILIKE $2)
			LIMIT 20
		`
		rows, err := r.db.QueryContext(ctx, lessonQuery, schoolID, searchPattern)
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

	if results == nil {
		results = []SearchResultItem{}
	}
	return results, nil
}
