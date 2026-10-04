package pcto_tutor

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Repository interface {
	GetTutorByToken(ctx context.Context, token string) (*CompanyTutor, error)
	GetTutorByID(ctx context.Context, id string) (*CompanyTutor, error)
	CreateTutor(ctx context.Context, tutor *CompanyTutor) error
	AssignStudent(ctx context.Context, tutorID, projectID, studentID string) error
	ListAssignedStudents(ctx context.Context, tutorID string) ([]AssignedStudentInfo, error)
	SaveTimesheetVerification(ctx context.Context, entry *TimesheetVerification) error
	SaveCompanyEvaluation(ctx context.Context, eval *CompanyEvaluation) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetTutorByToken(ctx context.Context, token string) (*CompanyTutor, error) {
	query := `
		SELECT id, school_id, company_name, tutor_first_name, tutor_last_name, email,
		       COALESCE(phone, ''), access_token, token_expires_at, is_active, created_at, updated_at
		FROM pcto_company_tutors
		WHERE access_token = $1 AND is_active = true
	`
	var t CompanyTutor
	var expiresAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&t.ID, &t.SchoolID, &t.CompanyName, &t.TutorFirstName, &t.TutorLastName, &t.Email,
		&t.Phone, &t.AccessToken, &expiresAt, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if expiresAt.Valid {
		t.TokenExpiresAt = &expiresAt.Time
		if t.TokenExpiresAt.Before(time.Now()) {
			return nil, errors.New("il link di accesso è scaduto")
		}
	}
	return &t, nil
}

func (r *postgresRepo) GetTutorByID(ctx context.Context, id string) (*CompanyTutor, error) {
	query := `
		SELECT id, school_id, company_name, tutor_first_name, tutor_last_name, email,
		       COALESCE(phone, ''), COALESCE(access_token, ''), token_expires_at, is_active, created_at, updated_at
		FROM pcto_company_tutors
		WHERE id = $1
	`
	var t CompanyTutor
	var expiresAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.SchoolID, &t.CompanyName, &t.TutorFirstName, &t.TutorLastName, &t.Email,
		&t.Phone, &t.AccessToken, &expiresAt, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if expiresAt.Valid {
		t.TokenExpiresAt = &expiresAt.Time
	}
	return &t, nil
}

func (r *postgresRepo) CreateTutor(ctx context.Context, tutor *CompanyTutor) error {
	query := `
		INSERT INTO pcto_company_tutors (
			school_id, company_name, tutor_first_name, tutor_last_name, email, phone, access_token, token_expires_at, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at;
	`
	return r.db.QueryRowContext(ctx, query,
		tutor.SchoolID, tutor.CompanyName, tutor.TutorFirstName, tutor.TutorLastName,
		tutor.Email, tutor.Phone, tutor.AccessToken, tutor.TokenExpiresAt, tutor.IsActive,
	).Scan(&tutor.ID, &tutor.CreatedAt, &tutor.UpdatedAt)
}

func (r *postgresRepo) AssignStudent(ctx context.Context, tutorID, projectID, studentID string) error {
	query := `
		INSERT INTO pcto_tutor_assignments (tutor_id, project_id, student_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (tutor_id, project_id, student_id) DO NOTHING;
	`
	_, err := r.db.ExecContext(ctx, query, tutorID, projectID, studentID)
	return err
}

func (r *postgresRepo) ListAssignedStudents(ctx context.Context, tutorID string) ([]AssignedStudentInfo, error) {
	query := `
		SELECT 
			a.student_id,
			COALESCE(u.first_name || ' ' || u.last_name, 'Studente') AS student_name,
			a.project_id,
			COALESCE(p.title, 'Progetto PCTO') AS project_title,
			COALESCE(p.total_hours, 0) AS total_hours,
			COALESCE(part.hours_completed, 0) AS completed_hours,
			EXISTS(
				SELECT 1 FROM pcto_company_evaluations e 
				WHERE e.tutor_id = a.tutor_id AND e.student_id = a.student_id AND e.project_id = a.project_id
			) AS is_evaluated
		FROM pcto_tutor_assignments a
		LEFT JOIN users u ON u.id = a.student_id
		LEFT JOIN pcto_projects p ON p.id = a.project_id
		LEFT JOIN pcto_participations part ON part.project_id = a.project_id AND part.student_id = a.student_id
		WHERE a.tutor_id = $1
		ORDER BY student_name ASC;
	`
	rows, err := r.db.QueryContext(ctx, query, tutorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AssignedStudentInfo
	for rows.Next() {
		var s AssignedStudentInfo
		if err := rows.Scan(&s.StudentID, &s.StudentName, &s.ProjectID, &s.ProjectTitle, &s.TotalHours, &s.CompletedHours, &s.IsEvaluated); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *postgresRepo) SaveTimesheetVerification(ctx context.Context, entry *TimesheetVerification) error {
	query := `
		INSERT INTO pcto_timesheet_verifications (
			project_id, student_id, activity_date, hours_declared, hours_approved, tutor_id, tutor_notes
		) VALUES ($1, $2, $3::date, $4, $5, $6, $7)
		RETURNING id, signed_at;
	`
	return r.db.QueryRowContext(ctx, query,
		entry.ProjectID, entry.StudentID, entry.ActivityDate, entry.HoursDeclared, entry.HoursApproved, entry.TutorID, entry.TutorNotes,
	).Scan(&entry.ID, &entry.SignedAt)
}

func (r *postgresRepo) SaveCompanyEvaluation(ctx context.Context, eval *CompanyEvaluation) error {
	query := `
		INSERT INTO pcto_company_evaluations (
			tutor_id, student_id, project_id, reliability_level, technical_skills, teamwork_skills, final_feedback
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tutor_id, student_id, project_id)
		DO UPDATE SET
			reliability_level = EXCLUDED.reliability_level,
			technical_skills = EXCLUDED.technical_skills,
			teamwork_skills = EXCLUDED.teamwork_skills,
			final_feedback = EXCLUDED.final_feedback,
			submitted_at = CURRENT_TIMESTAMP
		RETURNING id, submitted_at;
	`
	return r.db.QueryRowContext(ctx, query,
		eval.TutorID, eval.StudentID, eval.ProjectID, eval.ReliabilityLevel, eval.TechnicalSkills, eval.TeamworkSkills, eval.FinalFeedback,
	).Scan(&eval.ID, &eval.SubmittedAt)
}
