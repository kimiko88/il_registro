package pcto

import (
	"context"
	"database/sql"
	"strings"
)

type Repository interface {
	CreateProject(ctx context.Context, p *Project) error
	GetProjects(ctx context.Context, schoolID string) ([]Project, error)
	GetProjectByID(ctx context.Context, id string) (*Project, error)
	UpdateProject(ctx context.Context, p *Project) error
	DeleteProject(ctx context.Context, id string) error

	AssignStudent(ctx context.Context, participation *Participation) error
	GetParticipationsByProject(ctx context.Context, projectID string) ([]Participation, error)
	GetParticipationsByStudent(ctx context.Context, studentID string) ([]Participation, error)
	GetParticipation(ctx context.Context, projectID, studentID string) (*Participation, error)

	LogHours(ctx context.Context, h *HourLog) error
	GetHours(ctx context.Context, participationID string) ([]HourLog, error)
	GetHourLogByID(ctx context.Context, id string) (*HourLog, error)
	GetParticipationByID(ctx context.Context, id string) (*Participation, error)
	VerifyHours(ctx context.Context, hourID, teacherID string) error
	UpdateHourLogStatus(ctx context.Context, logID, status string) error

	CreateCompany(ctx context.Context, c *Company) error
	GetCompanies(ctx context.Context, schoolID string) ([]Company, error)
	GetStats(ctx context.Context, schoolID string) (*PCTOStats, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateProject(ctx context.Context, p *Project) error {
	query := `
		INSERT INTO pcto_projects (school_id, title, description, type, start_date, end_date, total_hours, company_id, school_tutor_id, company_tutor_name, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	return r.db.QueryRowContext(ctx, query,
		p.SchoolID, p.Title, p.Description, p.Type, p.StartDate, p.EndDate, p.TotalHours, p.CompanyID, p.SchoolTutorID, p.CompanyTutor, p.CreatedBy,
	).Scan(&p.ID)
}

func (r *repository) GetProjects(ctx context.Context, schoolID string) ([]Project, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, school_id, title, description, type, start_date, end_date, total_hours, company_id, school_tutor_id, company_tutor_name, created_by FROM pcto_projects WHERE school_id=$1 ORDER BY start_date DESC`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.SchoolID, &p.Title, &p.Description, &p.Type, &p.StartDate, &p.EndDate, &p.TotalHours, &p.CompanyID, &p.SchoolTutorID, &p.CompanyTutor, &p.CreatedBy); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *repository) GetProjectByID(ctx context.Context, id string) (*Project, error) {
	var p Project
	var compName sql.NullString
	query := `
		SELECT p.id, p.school_id, p.title, p.description, p.type, p.start_date, p.end_date, p.total_hours, p.company_id, p.school_tutor_id, p.company_tutor_name, p.created_by,
		       COALESCE(c.name, '')
		FROM pcto_projects p
		LEFT JOIN pcto_companies c ON p.company_id = c.id
		WHERE p.id=$1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SchoolID, &p.Title, &p.Description, &p.Type, &p.StartDate, &p.EndDate, &p.TotalHours, &p.CompanyID, &p.SchoolTutorID, &p.CompanyTutor, &p.CreatedBy,
		&compName,
	)
	if err != nil {
		return nil, err
	}
	if compName.Valid {
		p.CompanyName = compName.String
	}
	return &p, nil
}

func (r *repository) UpdateProject(ctx context.Context, p *Project) error {
	_, err := r.db.ExecContext(ctx, `UPDATE pcto_projects SET title=$1, description=$2, total_hours=$3 WHERE id=$4`, p.Title, p.Description, p.TotalHours, p.ID)
	return err
}

func (r *repository) DeleteProject(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM pcto_projects WHERE id=$1`, id)
	return err
}

func (r *repository) AssignStudent(ctx context.Context, p *Participation) error {
	query := `INSERT INTO pcto_participations (project_id, student_id, status) VALUES ($1, $2, 'Active') RETURNING id`
	return r.db.QueryRowContext(ctx, query, p.ProjectID, p.StudentID).Scan(&p.ID)
}

func (r *repository) GetParticipationsByProject(ctx context.Context, projectID string) ([]Participation, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, project_id, student_id, status, hours_completed FROM pcto_participations WHERE project_id=$1`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var parts []Participation
	for rows.Next() {
		var p Participation
		if err := rows.Scan(&p.ID, &p.ProjectID, &p.StudentID, &p.Status, &p.HoursCompleted); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return parts, nil
}

func (r *repository) GetParticipationsByStudent(ctx context.Context, studentID string) ([]Participation, error) {
	query := `
		SELECT id, project_id, student_id, status, hours_completed 
		FROM pcto_participations 
		WHERE (
			student_id = $1::uuid OR
			student_id IN (SELECT id FROM students WHERE user_id = $1::uuid OR id = $1::uuid) OR
			student_id IN (SELECT user_id FROM students WHERE id = $1::uuid OR user_id = $1::uuid)
		)
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var parts []Participation
	for rows.Next() {
		var p Participation
		if err := rows.Scan(&p.ID, &p.ProjectID, &p.StudentID, &p.Status, &p.HoursCompleted); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return parts, nil
}

func (r *repository) GetParticipation(ctx context.Context, projectID, studentID string) (*Participation, error) {
	var p Participation
	query := `
		SELECT id, project_id, student_id, status, hours_completed 
		FROM pcto_participations 
		WHERE project_id = $1::uuid 
		  AND (
			student_id = $2::uuid OR
			student_id IN (SELECT id FROM students WHERE user_id = $2::uuid OR id = $2::uuid) OR
			student_id IN (SELECT user_id FROM students WHERE id = $2::uuid OR user_id = $2::uuid)
		  )
	`
	err := r.db.QueryRowContext(ctx, query, projectID, studentID).Scan(&p.ID, &p.ProjectID, &p.StudentID, &p.Status, &p.HoursCompleted)
	return &p, err
}

func (r *repository) LogHours(ctx context.Context, h *HourLog) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Insert Log
	err = tx.QueryRowContext(ctx, `INSERT INTO pcto_hours (participation_id, date, hours, activity_description) VALUES ($1, $2, $3, $4) RETURNING id`,
		h.ParticipationID, h.Date, h.Hours, h.Activity).Scan(&h.ID)
	if err != nil {
		return err
	}

	// Update Accumulator
	_, err = tx.ExecContext(ctx, `UPDATE pcto_participations SET hours_completed = hours_completed + $1 WHERE id=$2`, h.Hours, h.ParticipationID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) GetHours(ctx context.Context, participationID string) ([]HourLog, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, participation_id, date, hours, activity_description, verified, verified_by FROM pcto_hours WHERE participation_id=$1 ORDER BY date DESC`, participationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []HourLog
	for rows.Next() {
		var h HourLog
		var verifiedBy sql.NullString
		if err := rows.Scan(&h.ID, &h.ParticipationID, &h.Date, &h.Hours, &h.Activity, &h.Verified, &verifiedBy); err != nil {
			return nil, err
		}
		if verifiedBy.Valid {
			val := verifiedBy.String
			h.VerifiedBy = &val
		}
		logs = append(logs, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logs, nil
}


func (r *repository) GetHourLogByID(ctx context.Context, id string) (*HourLog, error) {
	var h HourLog
	var verifiedBy sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT id, participation_id, date, hours, activity_description, verified, verified_by FROM pcto_hours WHERE id=$1`, id).Scan(&h.ID, &h.ParticipationID, &h.Date, &h.Hours, &h.Activity, &h.Verified, &verifiedBy)
	if err != nil {
		return nil, err
	}
	if verifiedBy.Valid {
		val := verifiedBy.String
		h.VerifiedBy = &val
	}
	return &h, nil
}

func (r *repository) GetParticipationByID(ctx context.Context, id string) (*Participation, error) {
	var p Participation
	err := r.db.QueryRowContext(ctx, `SELECT id, project_id, student_id, status, hours_completed FROM pcto_participations WHERE id=$1`, id).Scan(&p.ID, &p.ProjectID, &p.StudentID, &p.Status, &p.HoursCompleted)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) VerifyHours(ctx context.Context, hourID, teacherID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE pcto_hours SET verified=TRUE, verified_by=$1, verified_at=NOW() WHERE id=$2`, teacherID, hourID)
	return err
}

func (r *repository) CreateCompany(ctx context.Context, c *Company) error {
	if c.ContactPerson == "" && (c.ContactPersonFirstName != "" || c.ContactPersonLastName != "") {
		c.ContactPerson = strings.TrimSpace(c.ContactPersonFirstName + " " + c.ContactPersonLastName)
	}
	query := `
		INSERT INTO pcto_companies (school_id, name, vat_number, address, contact_person, contact_person_first_name, contact_person_last_name, contact_person_phone, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		c.SchoolID, c.Name, c.VatNumber, c.Address, c.ContactPerson, c.ContactPersonFirstName, c.ContactPersonLastName, c.ContactPersonPhone, c.Email).Scan(&c.ID)
}

func (r *repository) GetCompanies(ctx context.Context, schoolID string) ([]Company, error) {
	query := `
		SELECT id, school_id, name, vat_number, address, contact_person, COALESCE(contact_person_first_name, ''), COALESCE(contact_person_last_name, ''), COALESCE(contact_person_phone, ''), email
		FROM pcto_companies
		WHERE school_id=$1
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comps []Company
	for rows.Next() {
		var c Company
		err := rows.Scan(
			&c.ID, &c.SchoolID, &c.Name, &c.VatNumber, &c.Address, &c.ContactPerson, &c.ContactPersonFirstName, &c.ContactPersonLastName, &c.ContactPersonPhone, &c.Email,
		)
		if err != nil {
			return nil, err
		}
		comps = append(comps, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comps, nil
}


func (r *repository) GetStats(ctx context.Context, schoolID string) (*PCTOStats, error) {
	stats := &PCTOStats{}

	// Total Projects
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM pcto_projects WHERE school_id = $1", schoolID).Scan(&stats.TotalProjects)
	if err != nil {
		return nil, err
	}

	// Total Students Involved
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(DISTINCT student_id) FROM pcto_participations WHERE project_id IN (SELECT id FROM pcto_projects WHERE school_id = $1)", schoolID).Scan(&stats.TotalStudents)
	if err != nil {
		return nil, err
	}

	// Total Hours Recorded
	err = r.db.QueryRowContext(ctx, "SELECT COALESCE(SUM(hours), 0) FROM pcto_hours WHERE participation_id IN (SELECT id FROM pcto_participations WHERE project_id IN (SELECT id FROM pcto_projects WHERE school_id = $1))", schoolID).Scan(&stats.TotalHours)
	if err != nil {
		return nil, err
	}

	// Active Companies
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM pcto_companies WHERE school_id = $1", schoolID).Scan(&stats.ActiveCompanies)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *repository) UpdateHourLogStatus(ctx context.Context, logID, status string) error {
	query := `UPDATE pcto_hours SET status = $1, verified = ($1 = 'approved') WHERE id = $2::uuid`
	_, err := r.db.ExecContext(ctx, query, status, logID)
	return err
}
