package grades

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository interface {
	// Create inserts a new grade
	Create(grade *Grade) error

	// BatchCreate inserts multiple grades in a transaction
	BatchCreate(grades []*Grade) error

	// Update modifies an existing grade and logs the change to history
	Update(grade *Grade, history *GradeHistory) error

	// Delete performs a soft delete
	Delete(id string, deletedBy string) error

	// FindByID retrieves a single grade
	FindByID(id string) (*Grade, error)

	// FindByStudent retrieves all grades for a student
	FindByStudent(studentID string) ([]Grade, error)

	// FindByClassAndSubject retrieves grades for a specific class context
	FindByClassAndSubject(classID string, subjectID string, semester int) ([]Grade, error)

	// FindByClass retrieves all grades for a class (across all subjects)
	FindByClass(classID string, semester int) ([]Grade, error)

	// FindBySubject retrieves all grades for a subject (across classes if needed, or filtered)
	FindBySubject(subjectID string, semester int) ([]Grade, error)

	// FindWithFilter generic filter for export/advanced search
	FindWithFilter(filter GradeFilter) ([]Grade, error)

	// FindByTeacher retrieves grades assigned by a teacher (optional utility)
	FindByTeacher(teacherID string) ([]Grade, error)

	// GetHistory retrieves the modification history of a grade
	GetHistory(gradeID string) ([]GradeHistory, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(grade *Grade) error {
	query := `
		INSERT INTO grades (
			student_id, school_id, subject_id, teacher_id,
			grade_value, grade_type, semester, date, 
			description, rubric_id, weight, is_published, published_at,
			grade_category, evaluation_type, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, 
			$9, $10, $11, $12, $13,
			$14, $15, $16, NOW(), NOW()
		) RETURNING id`

	err := r.db.QueryRow(query,
		grade.StudentID, grade.SchoolID, grade.SubjectID, grade.TeacherID,
		grade.GradeValue, grade.GradeType, grade.Semester, grade.Date,
		grade.Description, grade.RubricID, grade.Weight, grade.IsPublished, grade.PublishedAt,
		grade.GradeCategory, grade.EvaluationType, grade.CreatedBy,
	).Scan(&grade.ID)

	if err != nil {
		return fmt.Errorf("create grade error: %w", err)
	}
	return nil
}

func (r *repository) BatchCreate(grades []*Grade) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("batch create begin tx error: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	query := `
		INSERT INTO grades (
			student_id, school_id, subject_id, teacher_id,
			grade_value, grade_type, semester, date, 
			description, rubric_id, weight, is_published, published_at,
			grade_category, evaluation_type, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, 
			$9, $10, $11, $12, $13,
			$14, $15, $16, NOW(), NOW()
		) RETURNING id`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare batch statement error: %w", err)
	}
	defer stmt.Close()

	for _, grade := range grades {
		err := stmt.QueryRow(
			grade.StudentID, grade.SchoolID, grade.SubjectID, grade.TeacherID,
			grade.GradeValue, grade.GradeType, grade.Semester, grade.Date,
			grade.Description, grade.RubricID, grade.Weight, grade.IsPublished, grade.PublishedAt,
			grade.GradeCategory, grade.EvaluationType, grade.CreatedBy,
		).Scan(&grade.ID)

		if err != nil {
			return fmt.Errorf("batch insert error for student %s: %w", grade.StudentID, err)
		}
	}

	return tx.Commit()
}

func (r *repository) Update(grade *Grade, history *GradeHistory) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// 1. Update Grade
	updateQuery := `
		UPDATE grades SET
			grade_value = $1,
			grade_type = $2,
			semester = $3,
			date = $4,
			description = $5,
			rubric_id = $6,
			weight = $7,
			is_published = $8,
			published_at = $9,
			grade_category = $10,
			evaluation_type = $11,
			modified_by = $12,
			updated_at = NOW()
		WHERE id = $13 AND deleted_at IS NULL`

	_, err = tx.Exec(updateQuery,
		grade.GradeValue, grade.GradeType, grade.Semester, grade.Date,
		grade.Description, grade.RubricID, grade.Weight, grade.IsPublished, grade.PublishedAt,
		grade.GradeCategory, grade.EvaluationType, grade.ModifiedBy, grade.ID,
	)
	if err != nil {
		return fmt.Errorf("update grade error: %w", err)
	}

	// 2. Insert History if provided
	if history != nil {
		historyQuery := `
			INSERT INTO grade_history (
				grade_id, old_value, new_value, 
				old_description, new_description, 
				modified_by, modified_at, reason
			) VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7) RETURNING id`

		err = tx.QueryRow(historyQuery,
			grade.ID, history.OldValue, history.NewValue,
			history.OldDescription, history.NewDescription,
			history.ModifiedBy, history.Reason,
		).Scan(&history.ID)

		if err != nil {
			return fmt.Errorf("insert history error: %w", err)
		}
	}

	return tx.Commit()
}

func (r *repository) Delete(id string, deletedBy string) error {
	query := `UPDATE grades SET deleted_at = NOW(), modified_by = $1 WHERE id = $2`
	_, err := r.db.Exec(query, deletedBy, id)
	if err != nil {
		return fmt.Errorf("delete grade error: %w", err)
	}
	return nil
}

func (r *repository) FindByID(id string) (*Grade, error) {
	query := `
		SELECT id, student_id, school_id, subject_id, teacher_id, 
			       grade_value, grade_type, semester, date, 
			       description, rubric_id, weight, is_published, published_at,
			       grade_category, evaluation_type, created_by, created_at, updated_at 
		FROM grades 
		WHERE id = $1 AND deleted_at IS NULL`

	var g Grade
	err := r.db.QueryRow(query, id).Scan(
		&g.ID, &g.StudentID, &g.SchoolID, &g.SubjectID, &g.TeacherID,
		&g.GradeValue, &g.GradeType, &g.Semester, &g.Date,
		&g.Description, &g.RubricID, &g.Weight, &g.IsPublished, &g.PublishedAt,
		&g.GradeCategory, &g.EvaluationType, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find grade error: %w", err)
	}
	return &g, nil
}

func (r *repository) FindByStudent(studentID string) ([]Grade, error) {
	query := `
		SELECT id, student_id, school_id, subject_id, teacher_id, 
			       grade_value, grade_type, semester, date, 
			       description, rubric_id, weight, is_published, published_at,
			       grade_category, evaluation_type, created_by, created_at, updated_at 
		FROM grades 
		WHERE student_id = $1 AND deleted_at IS NULL
		ORDER BY date DESC`

	return r.scanGrades(query, studentID)
}

func (r *repository) FindByClassAndSubject(classID string, subjectID string, semester int) ([]Grade, error) {
	// This requires joining with students table to filter by class_id
	var query string
	var args []interface{}

	if semester > 0 {
		query = `
			SELECT g.id, g.student_id, g.school_id, g.subject_id, g.teacher_id, 
			       g.grade_value, g.grade_type, g.semester, g.date, 
			       g.description, g.rubric_id, g.weight, g.is_published, g.published_at,
			       g.grade_category, g.evaluation_type, g.created_by, g.created_at, g.updated_at 
			FROM grades g
			JOIN students s ON g.student_id = s.id
			WHERE s.class_id = $1 AND g.subject_id = $2 AND g.semester = $3 AND g.deleted_at IS NULL
			ORDER BY g.date DESC`
		args = []interface{}{classID, subjectID, semester}
	} else {
		// All semesters
		query = `
			SELECT g.id, g.student_id, g.school_id, g.subject_id, g.teacher_id, 
			       g.grade_value, g.grade_type, g.semester, g.date, 
			       g.description, g.rubric_id, g.weight, g.is_published, g.published_at,
			       g.grade_category, g.evaluation_type, g.created_by, g.created_at, g.updated_at 
			FROM grades g
			JOIN students s ON g.student_id = s.id
			WHERE s.class_id = $1 AND g.subject_id = $2 AND g.deleted_at IS NULL
			ORDER BY g.date DESC`
		args = []interface{}{classID, subjectID}
	}

	return r.scanGrades(query, args...)
}

func (r *repository) FindByClass(classID string, semester int) ([]Grade, error) {
	var query string
	var args []interface{}

	if semester > 0 {
		query = `
			SELECT g.id, g.student_id, g.school_id, g.subject_id, g.teacher_id, 
			       g.grade_value, g.grade_type, g.semester, g.date, 
			       g.description, g.rubric_id, g.weight, g.is_published, g.published_at,
			       g.grade_category, g.evaluation_type, g.created_by, g.created_at, g.updated_at 
			FROM grades g
			JOIN students s ON g.student_id = s.id
			WHERE s.class_id = $1 AND g.semester = $2 AND g.deleted_at IS NULL
			ORDER BY g.date DESC`
		args = []interface{}{classID, semester}
	} else {
		query = `
			SELECT g.id, g.student_id, g.school_id, g.subject_id, g.teacher_id, 
			       g.grade_value, g.grade_type, g.semester, g.date, 
			       g.description, g.rubric_id, g.weight, g.is_published, g.published_at,
			       g.grade_category, g.evaluation_type, g.created_by, g.created_at, g.updated_at 
			FROM grades g
			JOIN students s ON g.student_id = s.id
			WHERE s.class_id = $1 AND g.deleted_at IS NULL
			ORDER BY g.date DESC`
		args = []interface{}{classID}
	}

	return r.scanGrades(query, args...)
}

func (r *repository) FindBySubject(subjectID string, semester int) ([]Grade, error) {
	var query string
	var args []interface{}

	if semester > 0 {
		query = `
			SELECT id, student_id, school_id, subject_id, teacher_id, 
			       grade_value, grade_type, semester, date, 
			       description, rubric_id, weight, is_published, published_at,
			       grade_category, evaluation_type, created_by, created_at, updated_at 
			FROM grades 
			WHERE subject_id = $1 AND semester = $2 AND deleted_at IS NULL
			ORDER BY date DESC, student_id ASC`
		args = []interface{}{subjectID, semester}
	} else {
		query = `
			SELECT id, student_id, school_id, subject_id, teacher_id, 
			       grade_value, grade_type, semester, date, 
			       description, rubric_id, weight, is_published, published_at,
			       grade_category, evaluation_type, created_by, created_at, updated_at 
			FROM grades 
			WHERE subject_id = $1 AND deleted_at IS NULL
			ORDER BY date DESC, student_id ASC`
		args = []interface{}{subjectID}
	}

	return r.scanGrades(query, args...)
}

func (r *repository) FindWithFilter(filter GradeFilter) ([]Grade, error) {
	baseQuery := `
		SELECT id, student_id, school_id, subject_id, teacher_id, 
			       grade_value, grade_type, semester, date, 
			       description, rubric_id, weight, is_published, published_at,
			       grade_category, evaluation_type, created_by, created_at, updated_at 
		FROM grades 
		WHERE deleted_at IS NULL`

	var args []interface{}
	var conditions []string
	argIdx := 1

	if filter.Semester > 0 {
		conditions = append(conditions, fmt.Sprintf("semester = $%d", argIdx))
		args = append(args, filter.Semester)
		argIdx++
	}
	if filter.SubjectID != "" {
		conditions = append(conditions, fmt.Sprintf("subject_id = $%d", argIdx))
		args = append(args, filter.SubjectID)
		argIdx++
	}
	if filter.GradeType != "" {
		conditions = append(conditions, fmt.Sprintf("grade_type = $%d", argIdx))
		args = append(args, filter.GradeType)
		argIdx++
	}
	if filter.IsPublished != nil {
		conditions = append(conditions, fmt.Sprintf("is_published = $%d", argIdx))
		args = append(args, *filter.IsPublished)
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY date DESC" // Default sorting

	return r.scanGrades(baseQuery, args...)
}

func (r *repository) FindByTeacher(teacherID string) ([]Grade, error) {
	query := `
		SELECT id, student_id, school_id, subject_id, teacher_id, 
			       grade_value, grade_type, semester, date, 
			       description, rubric_id, weight, is_published, published_at,
			       grade_category, evaluation_type, created_by, created_at, updated_at 
		FROM grades 
		WHERE teacher_id = $1 AND deleted_at IS NULL
		ORDER BY date DESC`

	return r.scanGrades(query, teacherID)
}

func (r *repository) GetHistory(gradeID string) ([]GradeHistory, error) {
	query := `
		SELECT id, grade_id, old_value, new_value, 
		       old_description, new_description, 
		       modified_by, modified_at, reason
		FROM grade_history
		WHERE grade_id = $1
		ORDER BY modified_at DESC`

	rows, err := r.db.Query(query, gradeID)
	if err != nil {
		return nil, fmt.Errorf("query history error: %w", err)
	}
	defer rows.Close()

	var history []GradeHistory
	for rows.Next() {
		var h GradeHistory
		if err := rows.Scan(
			&h.ID, &h.GradeID, &h.OldValue, &h.NewValue,
			&h.OldDescription, &h.NewDescription,
			&h.ModifiedBy, &h.ModifiedAt, &h.Reason,
		); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

// Helper to scan rows into Grade slice
func (r *repository) scanGrades(query string, args ...interface{}) ([]Grade, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *repository) scanRows(rows *sql.Rows) ([]Grade, error) {
	var results []Grade
	for rows.Next() {
		var g Grade
		if err := rows.Scan(
			&g.ID, &g.StudentID, &g.SchoolID, &g.SubjectID, &g.TeacherID,
			&g.GradeValue, &g.GradeType, &g.Semester, &g.Date,
			&g.Description, &g.RubricID, &g.Weight, &g.IsPublished, &g.PublishedAt,
			&g.GradeCategory, &g.EvaluationType, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, g)
	}
	return results, nil
}
