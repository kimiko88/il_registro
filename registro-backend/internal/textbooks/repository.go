package textbooks

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, t *Textbook) error
	List(ctx context.Context, schoolID string) ([]Textbook, error)
	Delete(ctx context.Context, id string) error
	
	AssignToClass(ctx context.Context, classID string, subjectID string, textbookID string, optional bool) error
	RemoveFromClass(ctx context.Context, assignmentID string) error
	ListByClass(ctx context.Context, classID string) ([]ClassTextbook, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, t *Textbook) error {
	t.ID = uuid.New().String()
	query := `INSERT INTO textbooks (id, school_id, title, author, subject, isbn, publisher, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, t.ID, t.SchoolID, t.Title, t.Author, t.Subject, t.ISBN, t.Publisher, t.Price)
	return err
}

func (r *postgresRepository) List(ctx context.Context, schoolID string) ([]Textbook, error) {
	query := `SELECT id, school_id, title, COALESCE(author, ''), COALESCE(subject, ''), COALESCE(isbn, ''), COALESCE(publisher, ''), COALESCE(price, 0), created_at FROM textbooks WHERE school_id = $1 ORDER BY title`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var res []Textbook
	for rows.Next() {
		var t Textbook
		if err := rows.Scan(&t.ID, &t.SchoolID, &t.Title, &t.Author, &t.Subject, &t.ISBN, &t.Publisher, &t.Price, &t.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM textbooks WHERE id = $1", id)
	return err
}

func (r *postgresRepository) AssignToClass(ctx context.Context, classID string, subjectID string, textbookID string, optional bool) error {
	id := uuid.New().String()
	query := `INSERT INTO class_textbooks (id, class_id, subject_id, textbook_id, is_optional) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, id, classID, subjectID, textbookID, optional)
	return err
}

func (r *postgresRepository) RemoveFromClass(ctx context.Context, assignmentID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM class_textbooks WHERE id = $1", assignmentID)
	return err
}

func (r *postgresRepository) ListByClass(ctx context.Context, classID string) ([]ClassTextbook, error) {
	query := `
		SELECT ct.id, ct.class_id, ct.subject_id, s.name, ct.textbook_id, t.title, t.author, ct.is_optional
		FROM class_textbooks ct
		JOIN textbooks t ON ct.textbook_id = t.id
		JOIN subjects s ON ct.subject_id = s.id
		WHERE ct.class_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var res []ClassTextbook
	for rows.Next() {
		var ct ClassTextbook
		if err := rows.Scan(&ct.ID, &ct.ClassID, &ct.SubjectID, &ct.SubjectName, &ct.TextbookID, &ct.Title, &ct.Author, &ct.IsOptional); err != nil {
			return nil, err
		}
		res = append(res, ct)
	}
	return res, nil
}
