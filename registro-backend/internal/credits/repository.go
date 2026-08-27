package credits

import (
	"context"
	"database/sql"
)

type Repository interface {
	SaveCredit(ctx context.Context, credit *StudentSchoolCredit) error
	GetCreditByStudentAndYear(ctx context.Context, studentID, academicYear string, gradeLevel int) (*StudentSchoolCredit, error)
	ListCreditsByClass(ctx context.Context, classID, academicYear string) ([]StudentSchoolCredit, error)
	GetStudentCreditSummary(ctx context.Context, studentID string) (*StudentCreditSummary, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) SaveCredit(ctx context.Context, c *StudentSchoolCredit) error {
	query := `
		INSERT INTO student_school_credits (
			id, school_id, student_id, class_id, academic_year, grade_level,
			grade_average, conduct_grade, base_credit_range_min, base_credit_range_max,
			assigned_credit, pcto_hours, has_extracurricular, deliberation_notes, validated_by,
			created_at, updated_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW()
		)
		ON CONFLICT (student_id, academic_year, grade_level)
		DO UPDATE SET
			assigned_credit = EXCLUDED.assigned_credit,
			grade_average = EXCLUDED.grade_average,
			conduct_grade = EXCLUDED.conduct_grade,
			base_credit_range_min = EXCLUDED.base_credit_range_min,
			base_credit_range_max = EXCLUDED.base_credit_range_max,
			pcto_hours = EXCLUDED.pcto_hours,
			has_extracurricular = EXCLUDED.has_extracurricular,
			deliberation_notes = EXCLUDED.deliberation_notes,
			validated_by = EXCLUDED.validated_by,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.SchoolID, c.StudentID, c.ClassID, c.AcademicYear, c.GradeLevel,
		c.GradeAverage, c.ConductGrade, c.BaseCreditRangeMin, c.BaseCreditRangeMax,
		c.AssignedCredit, c.PCTOHours, c.HasExtracurricular, c.DeliberationNotes, c.ValidatedBy,
	)
	return err
}

func (r *repository) GetCreditByStudentAndYear(ctx context.Context, studentID, academicYear string, gradeLevel int) (*StudentSchoolCredit, error) {
	query := `
		SELECT c.id, c.school_id, c.student_id, c.class_id, c.academic_year, c.grade_level,
		       c.grade_average, c.conduct_grade, c.base_credit_range_min, c.base_credit_range_max,
		       c.assigned_credit, c.pcto_hours, c.has_extracurricular, c.deliberation_notes,
		       c.validated_by, c.created_at, c.updated_at
		FROM student_school_credits c
		WHERE c.student_id = $1 AND c.academic_year = $2 AND c.grade_level = $3
	`
	var sc StudentSchoolCredit
	err := r.db.QueryRowContext(ctx, query, studentID, academicYear, gradeLevel).Scan(
		&sc.ID, &sc.SchoolID, &sc.StudentID, &sc.ClassID, &sc.AcademicYear, &sc.GradeLevel,
		&sc.GradeAverage, &sc.ConductGrade, &sc.BaseCreditRangeMin, &sc.BaseCreditRangeMax,
		&sc.AssignedCredit, &sc.PCTOHours, &sc.HasExtracurricular, &sc.DeliberationNotes,
		&sc.ValidatedBy, &sc.CreatedAt, &sc.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func (r *repository) ListCreditsByClass(ctx context.Context, classID, academicYear string) ([]StudentSchoolCredit, error) {
	query := `
		SELECT c.id, c.school_id, c.student_id,
		       TRIM(COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '')) as student_name,
		       c.class_id, cls.name as class_name, c.academic_year, c.grade_level,
		       c.grade_average, c.conduct_grade, c.base_credit_range_min, c.base_credit_range_max,
		       c.assigned_credit, c.pcto_hours, c.has_extracurricular, c.deliberation_notes,
		       c.validated_by, c.created_at, c.updated_at
		FROM student_school_credits c
		JOIN students s ON c.student_id = s.id
		JOIN users u ON s.user_id = u.id
		JOIN classes cls ON c.class_id = cls.id
		WHERE c.class_id = $1 AND ($2 = '' OR c.academic_year = $2)
		ORDER BY u.last_name ASC, u.first_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, classID, academicYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []StudentSchoolCredit
	for rows.Next() {
		var sc StudentSchoolCredit
		if err := rows.Scan(
			&sc.ID, &sc.SchoolID, &sc.StudentID, &sc.StudentName,
			&sc.ClassID, &sc.ClassName, &sc.AcademicYear, &sc.GradeLevel,
			&sc.GradeAverage, &sc.ConductGrade, &sc.BaseCreditRangeMin, &sc.BaseCreditRangeMax,
			&sc.AssignedCredit, &sc.PCTOHours, &sc.HasExtracurricular, &sc.DeliberationNotes,
			&sc.ValidatedBy, &sc.CreatedAt, &sc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, sc)
	}
	return list, rows.Err()
}

func (r *repository) GetStudentCreditSummary(ctx context.Context, studentID string) (*StudentCreditSummary, error) {
	query := `
		SELECT c.id, c.school_id, c.student_id,
		       TRIM(COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '')) as student_name,
		       c.class_id, cls.name as class_name, c.academic_year, c.grade_level,
		       c.grade_average, c.conduct_grade, c.base_credit_range_min, c.base_credit_range_max,
		       c.assigned_credit, c.pcto_hours, c.has_extracurricular, c.deliberation_notes,
		       c.validated_by, c.created_at, c.updated_at
		FROM student_school_credits c
		JOIN students s ON c.student_id = s.id
		JOIN users u ON s.user_id = u.id
		JOIN classes cls ON c.class_id = cls.id
		WHERE c.student_id = $1
		ORDER BY c.grade_level ASC
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := &StudentCreditSummary{
		StudentID:         studentID,
		MaxPossibleCredit: 40,
	}

	for rows.Next() {
		var sc StudentSchoolCredit
		if err := rows.Scan(
			&sc.ID, &sc.SchoolID, &sc.StudentID, &sc.StudentName,
			&sc.ClassID, &sc.ClassName, &sc.AcademicYear, &sc.GradeLevel,
			&sc.GradeAverage, &sc.ConductGrade, &sc.BaseCreditRangeMin, &sc.BaseCreditRangeMax,
			&sc.AssignedCredit, &sc.PCTOHours, &sc.HasExtracurricular, &sc.DeliberationNotes,
			&sc.ValidatedBy, &sc.CreatedAt, &sc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		summary.StudentName = sc.StudentName
		summary.ClassID = sc.ClassID
		summary.ClassName = sc.ClassName
		summary.CreditsByYear = append(summary.CreditsByYear, sc)
		summary.TotalTrienniumCredit += sc.AssignedCredit
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(summary.CreditsByYear) == 0 {
		// Fetch basic student info if no credits entered yet
		_ = r.db.QueryRowContext(ctx, `
			SELECT TRIM(COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '')), COALESCE(s.class_id::text, '')
			FROM students s JOIN users u ON s.user_id = u.id WHERE s.id = $1
		`, studentID).Scan(&summary.StudentName, &summary.ClassID)
	}

	return summary, nil
}
