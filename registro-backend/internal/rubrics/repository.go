package rubrics

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateRubric(ctx context.Context, r *Rubric) error
	GetRubricByID(ctx context.Context, id string) (*Rubric, error)
	ListRubrics(ctx context.Context, schoolID, teacherID string) ([]*Rubric, error)
	UpdateRubric(ctx context.Context, r *Rubric) error
	DeleteRubric(ctx context.Context, id string) error
	CreateAssessment(ctx context.Context, a *RubricAssessment) error
	ListAssessmentsByStudent(ctx context.Context, studentID string) ([]*RubricAssessment, error)
	ListAssessmentsByClass(ctx context.Context, classID string) ([]*RubricAssessment, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateRubric(ctx context.Context, rubric *Rubric) error {
	if rubric.ID == "" {
		rubric.ID = uuid.New().String()
	}
	rubric.CreatedAt = time.Now()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO rubrics (id, school_id, teacher_id, subject_id, title, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, rubric.ID, rubric.SchoolID, rubric.TeacherID, rubric.SubjectID, rubric.Title, rubric.Description, rubric.CreatedAt)
	if err != nil {
		return err
	}

	for i := range rubric.Criteria {
		c := &rubric.Criteria[i]
		if c.ID == "" {
			c.ID = uuid.New().String()
		}
		c.RubricID = rubric.ID

		_, err = tx.ExecContext(ctx, `
			INSERT INTO rubric_criteria (id, rubric_id, name, description, max_score)
			VALUES ($1, $2, $3, $4, $5)
		`, c.ID, c.RubricID, c.Name, c.Description, c.MaxScore)
		if err != nil {
			return err
		}

		for j := range c.Levels {
			l := &c.Levels[j]
			if l.ID == "" {
				l.ID = uuid.New().String()
			}
			l.CriterionID = c.ID

			_, err = tx.ExecContext(ctx, `
				INSERT INTO rubric_levels (id, criterion_id, score, label, description)
				VALUES ($1, $2, $3, $4, $5)
			`, l.ID, l.CriterionID, l.Score, l.Label, l.Description)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetRubricByID(ctx context.Context, id string) (*Rubric, error) {
	rub := &Rubric{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, school_id, teacher_id, subject_id, title, COALESCE(description, ''), created_at
		FROM rubrics WHERE id = $1::uuid
	`, id).Scan(&rub.ID, &rub.SchoolID, &rub.TeacherID, &rub.SubjectID, &rub.Title, &rub.Description, &rub.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Fetch Criteria
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, rubric_id, name, COALESCE(description, ''), max_score
		FROM rubric_criteria WHERE rubric_id = $1::uuid ORDER BY name
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Criterion
		if err := rows.Scan(&c.ID, &c.RubricID, &c.Name, &c.Description, &c.MaxScore); err != nil {
			return nil, err
		}

		// Fetch Levels for this criterion
		lRows, err := r.db.QueryContext(ctx, `
			SELECT id, criterion_id, score, label, COALESCE(description, '')
			FROM rubric_levels WHERE criterion_id = $1::uuid ORDER BY score DESC
		`, c.ID)
		if err == nil {
			for lRows.Next() {
				var l Level
				if err := lRows.Scan(&l.ID, &l.CriterionID, &l.Score, &l.Label, &l.Description); err == nil {
					c.Levels = append(c.Levels, l)
				}
			}
			if err := lRows.Err(); err != nil {
				lRows.Close()
				return nil, err
			}
			lRows.Close()
		}

		rub.Criteria = append(rub.Criteria, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rub, nil
}


func (r *PostgresRepository) ListRubrics(ctx context.Context, schoolID, teacherID string) ([]*Rubric, error) {
	query := `
		SELECT id, school_id, teacher_id, subject_id, title, COALESCE(description, ''), created_at
		FROM rubrics WHERE school_id = $1
	`
	args := []interface{}{schoolID}
	if teacherID != "" {
		query += " AND teacher_id = $2"
		args = append(args, teacherID)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Rubric
	for rows.Next() {
		rub := &Rubric{}
		if err := rows.Scan(&rub.ID, &rub.SchoolID, &rub.TeacherID, &rub.SubjectID, &rub.Title, &rub.Description, &rub.CreatedAt); err != nil {
			return nil, err
		}

		// Fetch Criteria
		cRows, cErr := r.db.QueryContext(ctx, `
			SELECT id, rubric_id, name, COALESCE(description, ''), max_score
			FROM rubric_criteria WHERE rubric_id = $1::uuid ORDER BY name
		`, rub.ID)
		if cErr == nil {
			for cRows.Next() {
				var c Criterion
				if err := cRows.Scan(&c.ID, &c.RubricID, &c.Name, &c.Description, &c.MaxScore); err == nil {
					rub.Criteria = append(rub.Criteria, c)
				}
			}
			cRows.Close()
		}
		if rub.Criteria == nil {
			rub.Criteria = []Criterion{}
		}

		list = append(list, rub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *PostgresRepository) UpdateRubric(ctx context.Context, rubric *Rubric) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		UPDATE rubrics SET title = $1, description = $2, subject_id = $3 WHERE id = $4::uuid
	`, rubric.Title, rubric.Description, rubric.SubjectID, rubric.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) DeleteRubric(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM rubrics WHERE id = $1::uuid`, id)
	return err
}

func (r *PostgresRepository) CreateAssessment(ctx context.Context, a *RubricAssessment) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	a.CreatedAt = time.Now()
	if a.Date.IsZero() {
		a.Date = time.Now()
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO rubric_assessments (id, rubric_id, student_id, class_id, teacher_id, date, total_score, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, a.ID, a.RubricID, a.StudentID, a.ClassID, a.TeacherID, a.Date, a.TotalScore, a.Notes, a.CreatedAt)
	if err != nil {
		return err
	}

	for _, cs := range a.Scores {
		var levelID sql.NullString
		if cs.LevelID != "" {
			levelID = sql.NullString{String: cs.LevelID, Valid: true}
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO rubric_criterion_scores (id, assessment_id, criterion_id, level_id, score)
			VALUES ($1, $2, $3, $4, $5)
		`, uuid.New().String(), a.ID, cs.CriterionID, levelID, cs.Score)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) ListAssessmentsByStudent(ctx context.Context, studentID string) ([]*RubricAssessment, error) {
	query := `
		SELECT a.id, a.rubric_id, a.student_id, a.class_id, a.teacher_id, a.date, a.total_score, COALESCE(a.notes, ''), a.created_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name,
		       COALESCE(r.title, '') AS rubric_title
		FROM rubric_assessments a
		LEFT JOIN users u ON a.student_id = u.id
		LEFT JOIN rubrics r ON a.rubric_id = r.id
		WHERE a.student_id = $1
		ORDER BY a.date DESC, a.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*RubricAssessment
	for rows.Next() {
		a := &RubricAssessment{}
		if err := rows.Scan(&a.ID, &a.RubricID, &a.StudentID, &a.ClassID, &a.TeacherID, &a.Date, &a.TotalScore, &a.Notes, &a.CreatedAt, &a.StudentName, &a.RubricTitle); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PostgresRepository) ListAssessmentsByClass(ctx context.Context, classID string) ([]*RubricAssessment, error) {
	query := `
		SELECT a.id, a.rubric_id, a.student_id, a.class_id, a.teacher_id, a.date, a.total_score, COALESCE(a.notes, ''), a.created_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name,
		       COALESCE(r.title, '') AS rubric_title
		FROM rubric_assessments a
		LEFT JOIN users u ON a.student_id = u.id
		LEFT JOIN rubrics r ON a.rubric_id = r.id
		WHERE a.class_id = $1
		ORDER BY a.date DESC, a.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*RubricAssessment
	for rows.Next() {
		a := &RubricAssessment{}
		if err := rows.Scan(&a.ID, &a.RubricID, &a.StudentID, &a.ClassID, &a.TeacherID, &a.Date, &a.TotalScore, &a.Notes, &a.CreatedAt, &a.StudentName, &a.RubricTitle); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

