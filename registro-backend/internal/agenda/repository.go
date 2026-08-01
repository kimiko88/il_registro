package agenda

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, item *AgendaItem) error
	GetByID(ctx context.Context, id string) (*AgendaItem, error)
	Update(ctx context.Context, item *AgendaItem) error
	Delete(ctx context.Context, id string) error
	ListCalendar(ctx context.Context, schoolID string, filter CalendarFilter) ([]*AgendaItem, error)
	SetCompletion(ctx context.Context, itemID, studentID string, completed bool) error
	IsStudentInClass(ctx context.Context, studentID, classID string) (bool, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, item *AgendaItem) error {
	item.ID = uuid.New().String()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	query := `
		INSERT INTO agenda_items (id, school_id, class_id, subject_id, teacher_id, title, description, type, all_day, start_time, end_time, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.SchoolID, item.ClassID, item.SubjectID, item.TeacherID,
		item.Title, item.Description, item.Type, item.AllDay, item.StartTime, item.EndTime, item.Date, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*AgendaItem, error) {
	query := `
		SELECT a.id, a.school_id, a.class_id, a.subject_id, a.teacher_id, a.title, COALESCE(a.description, ''),
		       a.type, COALESCE(a.all_day, false), COALESCE(a.start_time, '09:00'), COALESCE(a.end_time, '10:00'),
		       a.date, a.created_at, a.updated_at,
		       COALESCE(s.name, '') AS subject_name,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name
		FROM agenda_items a
		LEFT JOIN subjects s ON a.subject_id = s.id
		LEFT JOIN users u ON a.teacher_id = u.id
		WHERE a.id = $1::uuid
	`
	item := &AgendaItem{}
	var subjID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.SchoolID, &item.ClassID, &subjID, &item.TeacherID, &item.Title, &item.Description,
		&item.Type, &item.AllDay, &item.StartTime, &item.EndTime, &item.Date, &item.CreatedAt, &item.UpdatedAt, &item.SubjectName, &item.TeacherName,
	)
	if err != nil {
		return nil, err
	}
	if subjID.Valid {
		item.SubjectID = &subjID.String
	}
	return item, nil
}

func (r *PostgresRepository) Update(ctx context.Context, item *AgendaItem) error {
	item.UpdatedAt = time.Now()
	query := `
		UPDATE agenda_items 
		SET title = $1, description = $2, type = $3, all_day = $4, start_time = $5, end_time = $6, date = $7, updated_at = $8
		WHERE id = $9::uuid
	`
	_, err := r.db.ExecContext(ctx, query, item.Title, item.Description, item.Type, item.AllDay, item.StartTime, item.EndTime, item.Date, item.UpdatedAt, item.ID)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM agenda_items WHERE id = $1::uuid", id)
	return err
}

func (r *PostgresRepository) ListCalendar(ctx context.Context, schoolID string, filter CalendarFilter) ([]*AgendaItem, error) {
	query := `
		SELECT a.id, a.school_id, a.class_id, a.subject_id, a.teacher_id, a.title, COALESCE(a.description, ''),
		       a.type, COALESCE(a.all_day, false), COALESCE(a.start_time, '09:00'), COALESCE(a.end_time, '10:00'),
		       a.date, a.created_at, a.updated_at,
		       COALESCE(s.name, '') AS subject_name,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name,
		       CASE WHEN $2 = '' THEN false ELSE EXISTS(SELECT 1 FROM student_agenda_completions sac WHERE sac.agenda_item_id = a.id AND sac.student_id = NULLIF($2, '')::uuid) END AS is_completed
		FROM agenda_items a
		LEFT JOIN subjects s ON a.subject_id = s.id
		LEFT JOIN users u ON a.teacher_id = u.id
		WHERE ($1 = '' OR a.school_id = NULLIF($1, '')::uuid)
		  AND ($3 = '' OR a.class_id = NULLIF($3, '')::uuid)
		  AND ($4 = '' OR a.subject_id = NULLIF($4, '')::uuid)
		  AND ($5 = '' OR a.type = $5)
		  AND a.date >= $6 AND a.date <= $7
		ORDER BY a.date ASC, a.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, schoolID, filter.StudentID, filter.ClassID, filter.SubjectID, filter.Type, filter.From, filter.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*AgendaItem
	for rows.Next() {
		item := &AgendaItem{}
		var subjID sql.NullString
		if err := rows.Scan(
			&item.ID, &item.SchoolID, &item.ClassID, &subjID, &item.TeacherID, &item.Title, &item.Description,
			&item.Type, &item.AllDay, &item.StartTime, &item.EndTime, &item.Date, &item.CreatedAt, &item.UpdatedAt, &item.SubjectName, &item.TeacherName, &item.IsCompleted,
		); err != nil {
			return nil, err
		}
		if subjID.Valid {
			item.SubjectID = &subjID.String
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) SetCompletion(ctx context.Context, itemID, studentID string, completed bool) error {
	if completed {
		query := `
			INSERT INTO student_agenda_completions (agenda_item_id, student_id, completed_at)
			VALUES ($1::uuid, $2::uuid, NOW())
			ON CONFLICT (agenda_item_id, student_id) DO NOTHING
		`
		_, err := r.db.ExecContext(ctx, query, itemID, studentID)
		return err
	}
	query := `DELETE FROM student_agenda_completions WHERE agenda_item_id = $1::uuid AND student_id = $2::uuid`
	_, err := r.db.ExecContext(ctx, query, itemID, studentID)
	return err
}

func (r *PostgresRepository) IsStudentInClass(ctx context.Context, studentID, classID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM users u
			LEFT JOIN students s ON u.id = s.user_id
			WHERE u.id = $1::uuid AND (u.class_id = $2::uuid OR s.class_id = $2::uuid)
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, studentID, classID).Scan(&exists)
	return exists, err
}
