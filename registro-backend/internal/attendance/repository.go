package attendance

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository interface {
	Create(att *Attendance) error
	BatchCreate(atts []*Attendance) error
	Update(att *Attendance) error

	FindByID(id string) (*Attendance, error)
	FindByClassAndDate(classID string, date time.Time) ([]Attendance, error)
	FindByStudent(studentID string, startDate, endDate time.Time) ([]Attendance, error)
	GetStats(studentID string) (*SummaryResponse, error)

	// Justifications
	CreateJustification(j *Justification) error
	UpdateJustification(j *Justification) error
	FindJustificationByID(id string) (*Justification, error)
	FindPendingJustifications(classID string) ([]Justification, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// --- Attendance ---

func (r *repository) Create(a *Attendance) error {
	query := `
		INSERT INTO attendance (
			school_id, student_id, class_id, date, hour, subject_id, status, 
			justified, justified_by, justified_at, notes, entry_time, exit_time, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW()
		) RETURNING id`

	return r.db.QueryRow(query,
		a.SchoolID, a.StudentID, a.ClassID, a.Date, a.Hour, a.SubjectID, a.Status,
		a.Justified, a.JustifiedBy, a.JustifiedAt, a.Notes, a.EntryTime, a.ExitTime,
	).Scan(&a.ID)
}

func (r *repository) BatchCreate(atts []*Attendance) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.Prepare(`
		INSERT INTO attendance (
			school_id, student_id, class_id, date, hour, subject_id, status, 
			justified, justified_by, justified_at, notes, entry_time, exit_time, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
		RETURNING id
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, a := range atts {
		err := stmt.QueryRow(
			a.SchoolID, a.StudentID, a.ClassID, a.Date, a.Hour, a.SubjectID, a.Status,
			a.Justified, a.JustifiedBy, a.JustifiedAt, a.Notes, a.EntryTime, a.ExitTime,
		).Scan(&a.ID)
		if err != nil {
			return fmt.Errorf("batch insert error: %w", err)
		}
	}
	return tx.Commit()
}

func (r *repository) Update(a *Attendance) error {
	query := `
		UPDATE attendance SET 
			status=$1, hour=$2, subject_id=$3, notes=$4, justified=$5, justified_by=$6, justified_at=$7, entry_time=$8, exit_time=$9, updated_at=NOW()
		WHERE id=$10::uuid`
	_, err := r.db.Exec(query, a.Status, a.Hour, a.SubjectID, a.Notes, a.Justified, a.JustifiedBy, a.JustifiedAt, a.EntryTime, a.ExitTime, a.ID)
	return err
}

func (r *repository) FindByID(id string) (*Attendance, error) {
	query := `SELECT id, student_id, class_id, date, hour, subject_id, status, justified, justified_by, justified_at, COALESCE(notes, ''), entry_time, exit_time FROM attendance WHERE id=$1::uuid`
	var a Attendance
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.StudentID, &a.ClassID, &a.Date, &a.Hour, &a.SubjectID, &a.Status,
		&a.Justified, &a.JustifiedBy, &a.JustifiedAt, &a.Notes, &a.EntryTime, &a.ExitTime,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) FindByClassAndDate(classID string, date time.Time) ([]Attendance, error) {
	query := `
		SELECT id, student_id, class_id, date, hour, subject_id, status, justified, justified_by, justified_at, COALESCE(notes, ''), entry_time, exit_time
		FROM attendance 
		WHERE class_id=$1::uuid AND date=$2`

	rows, err := r.db.Query(query, classID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Attendance
	for rows.Next() {
		var a Attendance
		if err := rows.Scan(
			&a.ID, &a.StudentID, &a.ClassID, &a.Date, &a.Hour, &a.SubjectID, &a.Status, 
			&a.Justified, &a.JustifiedBy, &a.JustifiedAt, &a.Notes, &a.EntryTime, &a.ExitTime,
		); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, nil
}

func (r *repository) FindByStudent(studentID string, startDate, endDate time.Time) ([]Attendance, error) {
	query := `
		SELECT id, student_id, class_id, date, hour, subject_id, status, justified, justified_by, justified_at, COALESCE(notes, ''), entry_time, exit_time
		FROM attendance 
		WHERE student_id=$1::uuid AND date BETWEEN $2 AND $3
		ORDER BY date DESC`

	rows, err := r.db.Query(query, studentID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Attendance
	for rows.Next() {
		var a Attendance
		if err := rows.Scan(
			&a.ID, &a.StudentID, &a.ClassID, &a.Date, &a.Hour, &a.SubjectID, &a.Status, 
			&a.Justified, &a.JustifiedBy, &a.JustifiedAt, &a.Notes, &a.EntryTime, &a.ExitTime,
		); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, nil
}

func (r *repository) GetStats(studentID string) (*SummaryResponse, error) {
	// Simple aggregation
	query := `
		SELECT 
			COUNT(*) FILTER (WHERE status = 'Absent') as absences,
			COUNT(*) FILTER (WHERE status = 'Late') as lates,
			COUNT(*) FILTER (WHERE status = 'LeftEarly') as early_exits,
			COUNT(*) FILTER (WHERE justified = true) as justified
		FROM attendance
		WHERE student_id = $1::uuid`

	var s SummaryResponse
	err := r.db.QueryRow(query, studentID).Scan(&s.TotalAbsences, &s.TotalLates, &s.TotalEarlyExits, &s.JustifiedCount)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// --- Justifications ---

func (r *repository) CreateJustification(j *Justification) error {
	query := `
		INSERT INTO justifications (student_id, parent_id, start_date, end_date, reason, status, created_at, updated_at)
		VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6, NOW(), NOW()) RETURNING id`
	return r.db.QueryRow(query, j.StudentID, j.ParentID, j.StartDate, j.EndDate, j.Reason, j.Status).Scan(&j.ID)
}

func (r *repository) UpdateJustification(j *Justification) error {
	query := `UPDATE justifications SET status=$1, approved_by=$2::uuid, approved_at=$3, updated_at=NOW() WHERE id=$4::uuid`
	_, err := r.db.Exec(query, j.Status, j.ApprovedBy, j.ApprovedAt, j.ID)
	return err
}

func (r *repository) FindJustificationByID(id string) (*Justification, error) {
	var j Justification
	query := `SELECT id, student_id, parent_id, start_date, end_date, reason, status, approved_by FROM justifications WHERE id=$1::uuid`
	err := r.db.QueryRow(query, id).Scan(&j.ID, &j.StudentID, &j.ParentID, &j.StartDate, &j.EndDate, &j.Reason, &j.Status, &j.ApprovedBy)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *repository) FindPendingJustifications(classID string) ([]Justification, error) {
	// Join with students to filter by class
	query := `
		SELECT j.id, j.student_id, j.start_date, j.end_date, j.reason, j.status 
		FROM justifications j
		JOIN students s ON j.student_id = s.id
		WHERE s.class_id = $1::uuid AND j.status = 'pending'`

	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Justification
	for rows.Next() {
		var j Justification
		if err := rows.Scan(&j.ID, &j.StudentID, &j.StartDate, &j.EndDate, &j.Reason, &j.Status); err != nil {
			return nil, err
		}
		res = append(res, j)
	}
	return res, nil
}
