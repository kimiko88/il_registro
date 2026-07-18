package scrutiny

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	SaveRecord(ctx context.Context, record *ScrutinyRecord) error
	GetRecord(ctx context.Context, studentID, classID string, semester int) (*ScrutinyRecord, error)
	ListRecordsByClass(ctx context.Context, classID string, semester int) ([]ScrutinyRecord, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) SaveRecord(ctx context.Context, rec *ScrutinyRecord) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// 1. Upsert Scrutiny Record
	if rec.ID == "" {
		rec.ID = uuid.New().String()
	}
	rec.UpdatedAt = time.Now()

	query := `
		INSERT INTO scrutiny_records (id, student_id, class_id, semester, conduct_grade, final_decision, notes, coordinator_id, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (student_id, class_id, semester) 
		DO UPDATE SET conduct_grade = EXCLUDED.conduct_grade, final_decision = EXCLUDED.final_decision, notes = EXCLUDED.notes, updated_at = EXCLUDED.updated_at
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, query, rec.ID, rec.StudentID, rec.ClassID, rec.Semester, rec.ConductGrade, rec.FinalDecision, rec.Notes, rec.CoordinatorID, rec.UpdatedAt).Scan(&rec.ID)
	if err != nil {
		return err
	}

	// 2. Sync Grades
	_, err = tx.ExecContext(ctx, "DELETE FROM scrutiny_grades WHERE scrutiny_record_id = $1", rec.ID)
	if err != nil {
		return err
	}

	for _, g := range rec.Grades {
		g.ID = uuid.New().String()
		_, err = tx.ExecContext(ctx, "INSERT INTO scrutiny_grades (id, scrutiny_record_id, subject_id, final_grade, teacher_id) VALUES ($1, $2, $3, $4, $5)",
			g.ID, rec.ID, g.SubjectID, g.FinalGrade, g.TeacherID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *postgresRepository) GetRecord(ctx context.Context, studentID, classID string, semester int) (*ScrutinyRecord, error) {
	query := `SELECT id, student_id, class_id, semester, conduct_grade, final_decision, notes, coordinator_id, created_at, updated_at 
	          FROM scrutiny_records WHERE student_id = $1 AND class_id = $2 AND semester = $3`
	
	rec := &ScrutinyRecord{}
	err := r.db.QueryRowContext(ctx, query, studentID, classID, semester).Scan(
		&rec.ID, &rec.StudentID, &rec.ClassID, &rec.Semester, &rec.ConductGrade, &rec.FinalDecision, &rec.Notes, &rec.CoordinatorID, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Load grades
	gRows, err := r.db.QueryContext(ctx, "SELECT id, scrutiny_record_id, subject_id, final_grade, teacher_id FROM scrutiny_grades WHERE scrutiny_record_id = $1", rec.ID)
	if err != nil {
		return nil, err
	}
	defer gRows.Close()

	for gRows.Next() {
		var g ScrutinyGrade
		if err := gRows.Scan(&g.ID, &g.ScrutinyRecordID, &g.SubjectID, &g.FinalGrade, &g.TeacherID); err != nil {
			return nil, err
		}
		rec.Grades = append(rec.Grades, g)
	}

	return rec, nil
}

func (r *postgresRepository) ListRecordsByClass(ctx context.Context, classID string, semester int) ([]ScrutinyRecord, error) {
	query := `SELECT id, student_id, class_id, semester, conduct_grade, final_decision, notes, coordinator_id, created_at, updated_at 
	          FROM scrutiny_records WHERE class_id = $1 AND semester = $2`
	
	rows, err := r.db.QueryContext(ctx, query, classID, semester)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []ScrutinyRecord
	for rows.Next() {
		var rec ScrutinyRecord
		if err := rows.Scan(&rec.ID, &rec.StudentID, &rec.ClassID, &rec.Semester, &rec.ConductGrade, &rec.FinalDecision, &rec.Notes, &rec.CoordinatorID, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, rec)
	}
	return res, nil
}
