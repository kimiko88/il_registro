package middleschoolexam

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	GetOrCreateExamForClass(ctx context.Context, schoolID, classID, academicYear, president string) (*MiddleSchoolExam, error)
	GetExamByID(ctx context.Context, id string) (*MiddleSchoolExam, error)
	UpdateExamStatus(ctx context.Context, id, status string) error
	ListCandidates(ctx context.Context, examID string) ([]MiddleSchoolExamCandidate, error)
	SaveCandidateAdmission(ctx context.Context, examID, studentID string, grade int, judgment string, admitted bool) error
	SaveCandidateGrades(ctx context.Context, candidate *MiddleSchoolExamCandidate) error
	GetCandidate(ctx context.Context, examID, studentID string) (*MiddleSchoolExamCandidate, error)
	GetCandidateDiplomaData(ctx context.Context, candidateID string) (*DiplomaData, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetOrCreateExamForClass(ctx context.Context, schoolID, classID, academicYear, president string) (*MiddleSchoolExam, error) {
	if president == "" {
		president = "Dirigente Scolastico"
	}
	querySelect := `
		SELECT id, school_id, class_id, academic_year, subcommission_number, president_name, status, created_at
		FROM middle_school_exams
		WHERE class_id = $1 AND academic_year = $2
	`
	var e MiddleSchoolExam
	err := r.db.QueryRowContext(ctx, querySelect, classID, academicYear).Scan(
		&e.ID, &e.SchoolID, &e.ClassID, &e.AcademicYear, &e.SubcommissionNumber, &e.PresidentName, &e.Status, &e.CreatedAt,
	)
	if err == nil {
		return &e, nil
	}

	e.ID = uuid.New().String()
	e.SchoolID = schoolID
	e.ClassID = classID
	e.AcademicYear = academicYear
	e.SubcommissionNumber = 1
	e.PresidentName = president
	e.Status = "admission"
	e.CreatedAt = time.Now()

	queryInsert := `
		INSERT INTO middle_school_exams (id, school_id, class_id, academic_year, subcommission_number, president_name, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (class_id, academic_year) DO UPDATE SET president_name = EXCLUDED.president_name
		RETURNING id
	`
	_ = r.db.QueryRowContext(ctx, queryInsert, e.ID, e.SchoolID, e.ClassID, e.AcademicYear, e.SubcommissionNumber, e.PresidentName, e.Status, e.CreatedAt).Scan(&e.ID)

	return &e, nil
}

func (r *postgresRepository) GetExamByID(ctx context.Context, id string) (*MiddleSchoolExam, error) {
	query := `
		SELECT id, school_id, class_id, academic_year, subcommission_number, president_name, status, created_at
		FROM middle_school_exams
		WHERE id = $1
	`
	var e MiddleSchoolExam
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.SchoolID, &e.ClassID, &e.AcademicYear, &e.SubcommissionNumber, &e.PresidentName, &e.Status, &e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *postgresRepository) UpdateExamStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE middle_school_exams SET status = $1 WHERE id = $2", status, id)
	return err
}

func (r *postgresRepository) ListCandidates(ctx context.Context, examID string) ([]MiddleSchoolExamCandidate, error) {
	query := `
		SELECT 
			c.id, c.exam_id, c.student_id, 
			COALESCE(u.last_name || ' ' || u.first_name, ''),
			COALESCE(c.admission_grade, 6),
			COALESCE(c.admission_judgment, ''),
			c.is_admitted,
			COALESCE(c.grade_italian, 0),
			COALESCE(c.grade_math, 0),
			COALESCE(c.grade_english, 0),
			COALESCE(c.grade_second_lang, 0),
			COALESCE(c.grade_interview, 0),
			COALESCE(c.exam_mean, 0),
			COALESCE(c.final_grade, 0),
			COALESCE(c.has_honors, false),
			COALESCE(c.outcome, ''),
			c.deliberated_at,
			COALESCE(c.notes, ''),
			c.created_at
		FROM middle_school_exam_candidates c
		JOIN users u ON c.student_id = u.id
		WHERE c.exam_id = $1
		ORDER BY u.last_name ASC, u.first_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, examID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []MiddleSchoolExamCandidate
	for rows.Next() {
		var c MiddleSchoolExamCandidate
		if err := rows.Scan(
			&c.ID, &c.ExamID, &c.StudentID, &c.StudentName,
			&c.AdmissionGrade, &c.AdmissionJudgment, &c.IsAdmitted,
			&c.GradeItalian, &c.GradeMath, &c.GradeEnglish, &c.GradeSecondLang, &c.GradeInterview,
			&c.ExamMean, &c.FinalGrade, &c.HasHonors, &c.Outcome,
			&c.DeliberatedAt, &c.Notes, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *postgresRepository) SaveCandidateAdmission(ctx context.Context, examID, studentID string, grade int, judgment string, admitted bool) error {
	id := uuid.New().String()
	query := `
		INSERT INTO middle_school_exam_candidates (id, exam_id, student_id, admission_grade, admission_judgment, is_admitted)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (exam_id, student_id) DO UPDATE SET
			admission_grade = EXCLUDED.admission_grade,
			admission_judgment = EXCLUDED.admission_judgment,
			is_admitted = EXCLUDED.is_admitted
	`
	_, err := r.db.ExecContext(ctx, query, id, examID, studentID, grade, judgment, admitted)
	return err
}

func (r *postgresRepository) SaveCandidateGrades(ctx context.Context, c *MiddleSchoolExamCandidate) error {
	now := time.Now()
	query := `
		UPDATE middle_school_exam_candidates SET
			grade_italian = $1,
			grade_math = $2,
			grade_english = $3,
			grade_second_lang = $4,
			grade_interview = $5,
			exam_mean = $6,
			final_grade = $7,
			has_honors = $8,
			outcome = $9,
			deliberated_at = $10,
			notes = $11
		WHERE exam_id = $12 AND student_id = $13
	`
	_, err := r.db.ExecContext(ctx, query,
		c.GradeItalian, c.GradeMath, c.GradeEnglish, c.GradeSecondLang, c.GradeInterview,
		c.ExamMean, c.FinalGrade, c.HasHonors, c.Outcome, now, c.Notes,
		c.ExamID, c.StudentID,
	)
	return err
}

func (r *postgresRepository) GetCandidate(ctx context.Context, examID, studentID string) (*MiddleSchoolExamCandidate, error) {
	query := `
		SELECT 
			c.id, c.exam_id, c.student_id, 
			COALESCE(u.last_name || ' ' || u.first_name, ''),
			COALESCE(c.admission_grade, 6),
			COALESCE(c.admission_judgment, ''),
			c.is_admitted,
			COALESCE(c.grade_italian, 0),
			COALESCE(c.grade_math, 0),
			COALESCE(c.grade_english, 0),
			COALESCE(c.grade_second_lang, 0),
			COALESCE(c.grade_interview, 0),
			COALESCE(c.exam_mean, 0),
			COALESCE(c.final_grade, 0),
			COALESCE(c.has_honors, false),
			COALESCE(c.outcome, ''),
			c.deliberated_at,
			COALESCE(c.notes, ''),
			c.created_at
		FROM middle_school_exam_candidates c
		JOIN users u ON c.student_id = u.id
		WHERE c.exam_id = $1 AND c.student_id = $2
	`
	var c MiddleSchoolExamCandidate
	err := r.db.QueryRowContext(ctx, query, examID, studentID).Scan(
		&c.ID, &c.ExamID, &c.StudentID, &c.StudentName,
		&c.AdmissionGrade, &c.AdmissionJudgment, &c.IsAdmitted,
		&c.GradeItalian, &c.GradeMath, &c.GradeEnglish, &c.GradeSecondLang, &c.GradeInterview,
		&c.ExamMean, &c.FinalGrade, &c.HasHonors, &c.Outcome,
		&c.DeliberatedAt, &c.Notes, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *postgresRepository) GetCandidateDiplomaData(ctx context.Context, candidateID string) (*DiplomaData, error) {
	query := `
		SELECT 
			COALESCE(s.name, 'Scuola Secondaria di I Grado'),
			COALESCE(s.code, 'SCUOLA01'),
			e.academic_year,
			COALESCE(u.first_name || ' ' || u.last_name, ''),
			to_char(COALESCE(u.birth_date, '2012-01-01'::date), 'DD/MM/YYYY'),
			COALESCE(u.birth_place, 'Roma'),
			COALESCE(c.final_grade, 6),
			COALESCE(c.has_honors, false),
			e.president_name
		FROM middle_school_exam_candidates c
		JOIN middle_school_exams e ON c.exam_id = e.id
		JOIN schools s ON e.school_id = s.id
		JOIN users u ON c.student_id = u.id
		WHERE c.id = $1
	`
	var d DiplomaData
	err := r.db.QueryRowContext(ctx, query, candidateID).Scan(
		&d.SchoolName, &d.SchoolCode, &d.AcademicYear, &d.StudentFullName,
		&d.BirthDate, &d.BirthPlace, &d.FinalGrade, &d.HasHonors, &d.PresidentName,
	)
	if err != nil {
		return nil, err
	}
	d.IssueDate = time.Now()
	return &d, nil
}
