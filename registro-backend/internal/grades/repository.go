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

	// FindEnrolledSubjects returns the subject IDs a student is enrolled in for a given semester.
	// Used by GetSemesterReport to count subjects with no grades as failed.
	FindEnrolledSubjects(studentID string, semester int) ([]string, error)

	// CreateTest inserts a new class test
	CreateTest(test *ClassTest) error

	// FindTestsByClassAndSubject retrieves class tests
	FindTestsByClassAndSubject(classID string, subjectID string) ([]ClassTest, error)

	// FindUpcomingTestsByClass retrieves all upcoming tests for a class (today or future)
	FindUpcomingTestsByClass(classID string) ([]ClassTest, error)

	// DeleteTest deletes a test (cascade delete will handle grades in DB)
	DeleteTest(id string) error

	// UpdateTest updates a class test's metadata
	UpdateTest(test *ClassTest) error

	// FindGradesByTestID retrieves all grades linked to a class test
	FindGradesByTestID(testID string) ([]Grade, error)

	// FindTestByID retrieves a single class test by its ID
	FindTestByID(id string) (*ClassTest, error)
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
			grade_category, evaluation_type, created_by, test_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, 
			$9, $10, $11, $12, $13,
			$14, $15, $16, $17, NOW(), NOW()
		) RETURNING id`

	err := r.db.QueryRow(query,
		grade.StudentID, grade.SchoolID, grade.SubjectID, grade.TeacherID,
		grade.GradeValue, grade.GradeType, grade.Semester, grade.Date,
		grade.Description, grade.RubricID, grade.Weight, grade.IsPublished, grade.PublishedAt,
		grade.GradeCategory, grade.EvaluationType, grade.CreatedBy, grade.TestID,
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
			grade_category, evaluation_type, created_by, test_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, 
			$9, $10, $11, $12, $13,
			$14, $15, $16, $17, NOW(), NOW()
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
			grade.GradeCategory, grade.EvaluationType, grade.CreatedBy, grade.TestID,
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
		WHERE id = $13::uuid AND deleted_at IS NULL`

	_, err = tx.Exec(updateQuery,
		grade.GradeValue, grade.GradeType, grade.Semester, grade.Date,
		grade.Description, grade.RubricID, grade.Weight, grade.IsPublished, grade.PublishedAt,
		grade.GradeCategory, grade.EvaluationType, grade.ModifiedBy, grade.ID,
	)
	if err != nil {
		return fmt.Errorf("update grade error: %w", err)
	}

	if history != nil {
		historyQuery := `
			INSERT INTO grade_history (
				grade_id, old_value, new_value, 
				old_description, new_description, 
				modified_by, modified_at, reason
			) VALUES ($1::uuid, $2, $3, $4, $5, $6, NOW(), $7) RETURNING id`

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
	query := `UPDATE grades SET deleted_at = NOW(), modified_by = $1 WHERE id = $2::uuid`
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
			       grade_category, evaluation_type, COALESCE(created_by::text, ''), created_at, updated_at, test_id
		FROM grades 
		WHERE id = $1::uuid AND deleted_at IS NULL`

	var g Grade
	err := r.db.QueryRow(query, id).Scan(
		&g.ID, &g.StudentID, &g.SchoolID, &g.SubjectID, &g.TeacherID,
		&g.GradeValue, &g.GradeType, &g.Semester, &g.Date,
		&g.Description, &g.RubricID, &g.Weight, &g.IsPublished, &g.PublishedAt,
		&g.GradeCategory, &g.EvaluationType, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt, &g.TestID,
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
			       grade_category, evaluation_type, COALESCE(created_by::text, ''), created_at, updated_at, test_id
		FROM grades 
		WHERE student_id = $1::uuid AND deleted_at IS NULL
		ORDER BY date DESC`

	return r.scanGrades(query, studentID)
}

func (r *repository) FindByClassAndSubject(classID string, subjectID string, semester int) ([]Grade, error) {
	var query string
	var args []interface{}

	if semester > 0 {
		query = `
			SELECT g.id, g.student_id, g.school_id, g.subject_id, g.teacher_id, 
			       g.grade_value, g.grade_type, g.semester, g.date, 
			       g.description, g.rubric_id, g.weight, g.is_published, g.published_at,
			       g.grade_category, g.evaluation_type, COALESCE(g.created_by::text, ''), g.created_at, g.updated_at, g.test_id
			FROM grades g
			JOIN students s ON g.student_id = s.id
			WHERE s.class_id = $1::uuid AND g.subject_id = $2::uuid AND g.semester = $3 AND g.deleted_at IS NULL
			ORDER BY g.date DESC`
		args = []interface{}{classID, subjectID, semester}
	} else {
		query = `
			SELECT g.id, g.student_id, g.school_id, g.subject_id, g.teacher_id, 
			       g.grade_value, g.grade_type, g.semester, g.date, 
			       g.description, g.rubric_id, g.weight, g.is_published, g.published_at,
			       g.grade_category, g.evaluation_type, COALESCE(g.created_by::text, ''), g.created_at, g.updated_at, g.test_id
			FROM grades g
			JOIN students s ON g.student_id = s.id
			WHERE s.class_id = $1::uuid AND g.subject_id = $2::uuid AND g.deleted_at IS NULL
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
			       g.grade_category, g.evaluation_type, COALESCE(g.created_by::text, ''), g.created_at, g.updated_at, g.test_id
			FROM grades g
			JOIN students s ON g.student_id = s.id
			WHERE s.class_id = $1::uuid AND g.semester = $2 AND g.deleted_at IS NULL
			ORDER BY g.date DESC`
		args = []interface{}{classID, semester}
	} else {
		query = `
			SELECT g.id, g.student_id, g.school_id, g.subject_id, g.teacher_id, 
			       g.grade_value, g.grade_type, g.semester, g.date, 
			       g.description, g.rubric_id, g.weight, g.is_published, g.published_at,
			       g.grade_category, g.evaluation_type, COALESCE(g.created_by::text, ''), g.created_at, g.updated_at, g.test_id
			FROM grades g
			JOIN students s ON g.student_id = s.id
			WHERE s.class_id = $1::uuid AND g.deleted_at IS NULL
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
			       grade_category, evaluation_type, COALESCE(created_by::text, ''), created_at, updated_at, test_id
			FROM grades 
			WHERE subject_id = $1::uuid AND semester = $2 AND deleted_at IS NULL
			ORDER BY date DESC, student_id ASC`
		args = []interface{}{subjectID, semester}
	} else {
		query = `
			SELECT id, student_id, school_id, subject_id, teacher_id, 
			       grade_value, grade_type, semester, date, 
			       description, rubric_id, weight, is_published, published_at,
			       grade_category, evaluation_type, COALESCE(created_by::text, ''), created_at, updated_at, test_id
			FROM grades 
			WHERE subject_id = $1::uuid AND deleted_at IS NULL
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
			       grade_category, evaluation_type, COALESCE(created_by::text, ''), created_at, updated_at, test_id
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
		conditions = append(conditions, fmt.Sprintf("subject_id = $%d::uuid", argIdx))
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

	baseQuery += " ORDER BY date DESC"

	return r.scanGrades(baseQuery, args...)
}

func (r *repository) FindByTeacher(teacherID string) ([]Grade, error) {
	query := `
		SELECT id, student_id, school_id, subject_id, teacher_id, 
			       grade_value, grade_type, semester, date, 
			       description, rubric_id, weight, is_published, published_at,
			       grade_category, evaluation_type, COALESCE(created_by::text, ''), created_at, updated_at, test_id
		FROM grades 
		WHERE teacher_id = $1::uuid AND deleted_at IS NULL
		ORDER BY date DESC`

	return r.scanGrades(query, teacherID)
}

func (r *repository) GetHistory(gradeID string) ([]GradeHistory, error) {
	query := `
		SELECT id, grade_id, old_value, new_value, 
		       old_description, new_description, 
		       modified_by, modified_at, reason
		FROM grade_history
		WHERE grade_id = $1::uuid
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return history, nil
}

// FindEnrolledSubjects returns the distinct subject IDs a student is enrolled in
// for the given semester via their class's class_subjects assignments.
// semester is accepted for API consistency but not used in the query because
// class_subjects does not carry a semester column; the caller filters by semester
// at the grade level.
func (r *repository) FindEnrolledSubjects(studentID string, semester int) ([]string, error) {
	query := `
		SELECT DISTINCT cs.subject_id::text
		FROM class_subjects cs
		JOIN students st ON st.class_id = cs.class_id
		WHERE st.id = $1::uuid`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("FindEnrolledSubjects error: %w", err)
	}
	defer rows.Close()

	var subjects []string
	for rows.Next() {
		var subID string
		if err := rows.Scan(&subID); err != nil {
			return nil, err
		}
		subjects = append(subjects, subID)
	}
	return subjects, rows.Err()
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
			&g.GradeCategory, &g.EvaluationType, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt, &g.TestID,
		); err != nil {
			return nil, err
		}
		results = append(results, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// CreateTest inserts a new class test
func (r *repository) CreateTest(test *ClassTest) error {
	query := `
		INSERT INTO class_tests (
			class_id, subject_id, teacher_id, title, date,
			teacher_notes, parent_notes, evaluation_type, created_at, updated_at
		) VALUES (
			$1::uuid, $2::uuid, $3::uuid, $4, $5,
			$6, $7, $8, NOW(), NOW()
		) RETURNING id`

	err := r.db.QueryRow(query,
		test.ClassID, test.SubjectID, test.TeacherID, test.Title, test.Date,
		test.TeacherNotes, test.ParentNotes, test.EvaluationType,
	).Scan(&test.ID)

	if err != nil {
		return fmt.Errorf("create test error: %w", err)
	}
	return nil
}

// FindTestsByClassAndSubject retrieves class tests
func (r *repository) FindTestsByClassAndSubject(classID string, subjectID string) ([]ClassTest, error) {
	query := `
		SELECT id, class_id, subject_id, teacher_id, title, date,
		       teacher_notes, parent_notes, evaluation_type, created_at, updated_at
		FROM class_tests
		WHERE class_id = $1::uuid AND subject_id = $2::uuid
		ORDER BY date DESC, created_at DESC`

	rows, err := r.db.Query(query, classID, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []ClassTest
	for rows.Next() {
		var t ClassTest
		err := rows.Scan(
			&t.ID, &t.ClassID, &t.SubjectID, &t.TeacherID, &t.Title, &t.Date,
			&t.TeacherNotes, &t.ParentNotes, &t.EvaluationType, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tests = append(tests, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

// FindUpcomingTestsByClass retrieves upcoming tests for a class (today or future)
func (r *repository) FindUpcomingTestsByClass(classID string) ([]ClassTest, error) {
	query := `
		SELECT id, class_id, subject_id, teacher_id, title, date,
		       teacher_notes, parent_notes, evaluation_type, created_at, updated_at
		FROM class_tests
		WHERE class_id = $1::uuid AND date >= CURRENT_DATE
		ORDER BY date ASC`

	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []ClassTest
	for rows.Next() {
		var t ClassTest
		err := rows.Scan(
			&t.ID, &t.ClassID, &t.SubjectID, &t.TeacherID, &t.Title, &t.Date,
			&t.TeacherNotes, &t.ParentNotes, &t.EvaluationType, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tests = append(tests, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

// DeleteTest deletes a test (cascade delete will handle grades in DB)
func (r *repository) DeleteTest(id string) error {
	query := `DELETE FROM class_tests WHERE id = $1::uuid`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete test error: %w", err)
	}
	return nil
}

// UpdateTest updates a class test's metadata
func (r *repository) UpdateTest(test *ClassTest) error {
	query := `
		UPDATE class_tests
		SET title = $1, date = $2, teacher_notes = $3, parent_notes = $4, evaluation_type = $5, updated_at = NOW()
		WHERE id = $6::uuid`

	_, err := r.db.Exec(query,
		test.Title, test.Date, test.TeacherNotes, test.ParentNotes, test.EvaluationType, test.ID,
	)
	if err != nil {
		return fmt.Errorf("update test error: %w", err)
	}
	return nil
}

// FindGradesByTestID retrieves all grades linked to a class test
func (r *repository) FindGradesByTestID(testID string) ([]Grade, error) {
	query := `
		SELECT id, student_id, school_id, subject_id, teacher_id, 
			       grade_value, grade_type, semester, date, 
			       description, rubric_id, weight, is_published, published_at,
			       grade_category, evaluation_type, created_by, created_at, updated_at, test_id
		FROM grades 
		WHERE test_id = $1::uuid AND deleted_at IS NULL`
	return r.scanGrades(query, testID)
}

func (r *repository) FindTestByID(id string) (*ClassTest, error) {
	query := `
		SELECT id, class_id, subject_id, teacher_id, title, date, teacher_notes, parent_notes, evaluation_type, created_at, updated_at
		FROM class_tests
		WHERE id = $1::uuid`
	var t ClassTest
	err := r.db.QueryRow(query, id).Scan(
		&t.ID, &t.ClassID, &t.SubjectID, &t.TeacherID, &t.Title, &t.Date, &t.TeacherNotes, &t.ParentNotes, &t.EvaluationType, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("test not found")
		}
		return nil, fmt.Errorf("find test by id error: %w", err)
	}
	return &t, nil
}
