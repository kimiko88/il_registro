package primaryeval

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateObjective(ctx context.Context, obj *LearningObjective) (*LearningObjective, error)
	ListObjectives(ctx context.Context, schoolID, subjectID string, classID *string, yearGrade int, academicYear string) ([]LearningObjective, error)
	GetObjective(ctx context.Context, id string) (*LearningObjective, error)
	DeleteObjective(ctx context.Context, id string) error

	SaveEvaluationsBatch(ctx context.Context, schoolID, teacherID string, req *SaveEvaluationsBatchRequest) error
	ListEvaluationsByClassAndSubject(ctx context.Context, classID, subjectID string, semester int) ([]PrimaryEvaluation, error)
	GetStudentEvaluations(ctx context.Context, studentID string, semester int) ([]PrimaryEvaluation, error)
}

type sqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) CreateObjective(ctx context.Context, obj *LearningObjective) (*LearningObjective, error) {
	if obj.ID == "" {
		obj.ID = uuid.New().String()
	}
	now := time.Now()
	obj.CreatedAt = now
	obj.UpdatedAt = now

	query := `
		INSERT INTO primary_learning_objectives (id, school_id, class_id, subject_id, year_grade, title, description, academic_year, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		obj.ID, obj.SchoolID, obj.ClassID, obj.SubjectID, obj.YearGrade, obj.Title, obj.Description, obj.AcademicYear, obj.CreatedAt, obj.UpdatedAt,
	).Scan(&obj.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create primary learning objective: %w", err)
	}
	return obj, nil
}

func (r *sqlRepository) ListObjectives(ctx context.Context, schoolID, subjectID string, classID *string, yearGrade int, academicYear string) ([]LearningObjective, error) {
	query := `
		SELECT o.id, o.school_id, o.class_id, o.subject_id, COALESCE(s.name, ''), o.year_grade, o.title, o.description, o.academic_year, o.created_at, o.updated_at
		FROM primary_learning_objectives o
		LEFT JOIN subjects s ON s.id = o.subject_id
		WHERE o.school_id = $1
	`
	args := []interface{}{schoolID}
	idx := 2

	if subjectID != "" {
		query += fmt.Sprintf(" AND o.subject_id = $%d", idx)
		args = append(args, subjectID)
		idx++
	}
	if yearGrade > 0 {
		query += fmt.Sprintf(" AND o.year_grade = $%d", idx)
		args = append(args, yearGrade)
		idx++
	}
	if academicYear != "" {
		query += fmt.Sprintf(" AND o.academic_year = $%d", idx)
		args = append(args, academicYear)
		idx++
	}
	if classID != nil && *classID != "" {
		query += fmt.Sprintf(" AND (o.class_id = $%d OR o.class_id IS NULL)", idx)
		args = append(args, *classID)
	}

	query += " ORDER BY o.year_grade ASC, o.title ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query primary learning objectives: %w", err)
	}
	defer rows.Close()

	var result []LearningObjective
	for rows.Next() {
		var o LearningObjective
		var classIDNullable sql.NullString
		if err := rows.Scan(
			&o.ID, &o.SchoolID, &classIDNullable, &o.SubjectID, &o.SubjectName,
			&o.YearGrade, &o.Title, &o.Description, &o.AcademicYear, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan primary learning objective: %w", err)
		}
		if classIDNullable.Valid {
			v := classIDNullable.String
			o.ClassID = &v
		}
		result = append(result, o)
	}
	return result, rows.Err()
}

func (r *sqlRepository) GetObjective(ctx context.Context, id string) (*LearningObjective, error) {
	query := `
		SELECT o.id, o.school_id, o.class_id, o.subject_id, COALESCE(s.name, ''), o.year_grade, o.title, o.description, o.academic_year, o.created_at, o.updated_at
		FROM primary_learning_objectives o
		LEFT JOIN subjects s ON s.id = o.subject_id
		WHERE o.id = $1
	`
	var o LearningObjective
	var classIDNullable sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&o.ID, &o.SchoolID, &classIDNullable, &o.SubjectID, &o.SubjectName,
		&o.YearGrade, &o.Title, &o.Description, &o.AcademicYear, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if classIDNullable.Valid {
		v := classIDNullable.String
		o.ClassID = &v
	}
	return &o, nil
}

func (r *sqlRepository) DeleteObjective(ctx context.Context, id string) error {
	query := `DELETE FROM primary_learning_objectives WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *sqlRepository) SaveEvaluationsBatch(ctx context.Context, schoolID, teacherID string, req *SaveEvaluationsBatchRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	evalDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		evalDate = time.Now()
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO primary_evaluations (
			id, school_id, student_id, class_id, subject_id, teacher_id, objective_id,
			level, dimension_autonomy, dimension_continuity, dimension_familiarity, dimension_resources,
			date, semester, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW()
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range req.Evaluations {
		evalID := uuid.New().String()
		autonomy := item.DimensionAutonomy
		if autonomy == "" {
			autonomy = "autonomo"
		}
		continuity := item.DimensionContinuity
		if continuity == "" {
			continuity = "continuo"
		}
		familiarity := item.DimensionFamiliarity
		if familiarity == "" {
			familiarity = "nota"
		}
		resources := item.DimensionResources
		if resources == "" {
			resources = "risorse_proprie"
		}

		_, err := stmt.ExecContext(ctx,
			evalID, schoolID, item.StudentID, req.ClassID, req.SubjectID, teacherID, req.ObjectiveID,
			item.Level, autonomy, continuity, familiarity, resources,
			evalDate, req.Semester, item.Notes,
		)
		if err != nil {
			return fmt.Errorf("failed to insert primary evaluation for student %s: %w", item.StudentID, err)
		}
	}

	return tx.Commit()
}

func (r *sqlRepository) ListEvaluationsByClassAndSubject(ctx context.Context, classID, subjectID string, semester int) ([]PrimaryEvaluation, error) {
	query := `
		SELECT e.id, e.school_id, e.student_id, COALESCE(u.last_name || ' ' || u.first_name, ''),
		       e.class_id, e.subject_id, e.teacher_id, e.objective_id, COALESCE(o.title, ''),
		       e.level, e.dimension_autonomy, e.dimension_continuity, e.dimension_familiarity, e.dimension_resources,
		       e.date, e.semester, e.notes, e.created_at, e.updated_at
		FROM primary_evaluations e
		JOIN students st ON st.id = e.student_id
		JOIN users u ON u.id = st.user_id
		LEFT JOIN primary_learning_objectives o ON o.id = e.objective_id
		WHERE e.class_id = $1 AND e.subject_id = $2
	`
	args := []interface{}{classID, subjectID}
	if semester > 0 {
		query += " AND e.semester = $3"
		args = append(args, semester)
	}
	query += " ORDER BY e.date DESC, u.last_name ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query primary evaluations: %w", err)
	}
	defer rows.Close()

	var result []PrimaryEvaluation
	for rows.Next() {
		var e PrimaryEvaluation
		if err := rows.Scan(
			&e.ID, &e.SchoolID, &e.StudentID, &e.StudentName,
			&e.ClassID, &e.SubjectID, &e.TeacherID, &e.ObjectiveID, &e.ObjectiveTitle,
			&e.Level, &e.DimensionAutonomy, &e.DimensionContinuity, &e.DimensionFamiliarity, &e.DimensionResources,
			&e.Date, &e.Semester, &e.Notes, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan primary evaluation: %w", err)
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *sqlRepository) GetStudentEvaluations(ctx context.Context, studentID string, semester int) ([]PrimaryEvaluation, error) {
	query := `
		SELECT e.id, e.school_id, e.student_id, COALESCE(u.last_name || ' ' || u.first_name, ''),
		       e.class_id, e.subject_id, e.teacher_id, e.objective_id, COALESCE(o.title, ''),
		       e.level, e.dimension_autonomy, e.dimension_continuity, e.dimension_familiarity, e.dimension_resources,
		       e.date, e.semester, e.notes, e.created_at, e.updated_at
		FROM primary_evaluations e
		JOIN students st ON st.id = e.student_id
		JOIN users u ON u.id = st.user_id
		LEFT JOIN primary_learning_objectives o ON o.id = e.objective_id
		WHERE e.student_id = $1
	`
	args := []interface{}{studentID}
	if semester > 0 {
		query += " AND e.semester = $2"
		args = append(args, semester)
	}
	query += " ORDER BY e.date DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query student primary evaluations: %w", err)
	}
	defer rows.Close()

	var result []PrimaryEvaluation
	for rows.Next() {
		var e PrimaryEvaluation
		if err := rows.Scan(
			&e.ID, &e.SchoolID, &e.StudentID, &e.StudentName,
			&e.ClassID, &e.SubjectID, &e.TeacherID, &e.ObjectiveID, &e.ObjectiveTitle,
			&e.Level, &e.DimensionAutonomy, &e.DimensionContinuity, &e.DimensionFamiliarity, &e.DimensionResources,
			&e.Date, &e.Semester, &e.Notes, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan student evaluation: %w", err)
		}
		result = append(result, e)
	}
	return result, rows.Err()
}
