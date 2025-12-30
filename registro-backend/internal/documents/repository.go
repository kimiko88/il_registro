package documents

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	FindAll() ([]Document, error)
	Create(doc *Document) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Document, error) {
	query := `SELECT id, title, description, file_path, uploaded_by, type, date, created_at, updated_at FROM documents WHERE deleted_at IS NULL`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var docs []Document
	for rows.Next() {
		var d Document
		if err := rows.Scan(&d.ID, &d.Title, &d.Description, &d.FilePath, &d.UploadedBy, &d.Type, &d.Date, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

func (r *repository) Create(doc *Document) error {
	query := `INSERT INTO documents (title, description, file_path, uploaded_by, type, date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id`

	err := r.db.QueryRow(query, doc.Title, doc.Description, doc.FilePath, doc.UploadedBy, doc.Type, doc.Date).Scan(&doc.ID)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}
	return nil
}
