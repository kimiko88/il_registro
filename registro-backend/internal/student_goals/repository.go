package student_goals

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, goal *StudentGoal) error
	GetByID(ctx context.Context, id string) (*StudentGoal, error)
	ListByStudent(ctx context.Context, studentID string) ([]*StudentGoal, error)
	UpdateStatus(ctx context.Context, id string, status GoalStatus) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, g *StudentGoal) error {
	g.ID = uuid.New().String()
	g.CreatedAt = time.Now()
	if g.Status == "" {
		g.Status = StatusPending
	}

	query := `
		INSERT INTO student_goals (id, student_id, teacher_id, title, description, badge_name, badge_icon, category, status, points, due_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		g.ID, g.StudentID, g.TeacherID, g.Title, g.Description, g.BadgeName, g.BadgeIcon,
		g.Category, g.Status, g.Points, g.DueDate, g.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*StudentGoal, error) {
	query := `
		SELECT id, student_id, teacher_id, title, description, COALESCE(badge_name, ''), COALESCE(badge_icon, ''),
		       category, status, points, due_date, created_at, completed_at
		FROM student_goals WHERE id = $1::uuid
	`
	g := &StudentGoal{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&g.ID, &g.StudentID, &g.TeacherID, &g.Title, &g.Description, &g.BadgeName, &g.BadgeIcon,
		&g.Category, &g.Status, &g.Points, &g.DueDate, &g.CreatedAt, &g.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (r *PostgresRepository) ListByStudent(ctx context.Context, studentID string) ([]*StudentGoal, error) {
	query := `
		SELECT id, student_id, teacher_id, title, description, COALESCE(badge_name, ''), COALESCE(badge_icon, ''),
		       category, status, points, due_date, created_at, completed_at
		FROM student_goals 
		WHERE (
			student_id = $1::uuid OR
			student_id IN (SELECT user_id FROM students WHERE id = $1::uuid) OR
			student_id IN (SELECT id FROM students WHERE user_id = $1::uuid)
		)
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var list []*StudentGoal
	for rows.Next() {
		g := &StudentGoal{}
		if err := rows.Scan(
			&g.ID, &g.StudentID, &g.TeacherID, &g.Title, &g.Description, &g.BadgeName, &g.BadgeIcon,
			&g.Category, &g.Status, &g.Points, &g.DueDate, &g.CreatedAt, &g.CompletedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, status GoalStatus) error {
	var completedAt *time.Time
	if status == StatusCompleted {
		now := time.Now()
		completedAt = &now
	}
	query := `UPDATE student_goals SET status = $2, completed_at = $3 WHERE id = $1::uuid`
	_, err := r.db.ExecContext(ctx, query, id, status, completedAt)
	return err
}
