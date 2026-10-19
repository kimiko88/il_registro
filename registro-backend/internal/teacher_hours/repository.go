package teacher_hours

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository interface {
	GetOrCreateObligation(schoolID, teacherID, academicYear string) (*TeacherHourObligation, error)
	GetObligation(obligationID string) (*TeacherHourObligation, error)
	ListObligationsBySchool(schoolID, academicYear string) ([]TeacherHourObligation, error)
	LogHourActivity(log *TeacherHourLog) error
	ListLogsByObligation(obligationID string) ([]TeacherHourLog, error)
	CreateShortLeave(leave *TeacherShortLeave) error
	ListShortLeaves(schoolID, teacherID string) ([]TeacherShortLeave, error)
	UpdateShortLeaveStatus(id, status string, notes string) error
	RecordShortLeaveRecovery(id string, recoveredHours float64, recoveredDate string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) resolveSchoolID(schoolID, teacherID string) string {
	if schoolID != "" {
		return schoolID
	}
	if teacherID != "" {
		var userSchoolID sql.NullString
		_ = r.db.QueryRow(`SELECT school_id::text FROM users WHERE id::text = $1 AND school_id IS NOT NULL`, teacherID).Scan(&userSchoolID)
		if userSchoolID.Valid && userSchoolID.String != "" {
			return userSchoolID.String
		}
	}
	var defaultSchoolID string
	_ = r.db.QueryRow(`SELECT id::text FROM schools ORDER BY created_at ASC LIMIT 1`).Scan(&defaultSchoolID)
	return defaultSchoolID
}

func (r *repository) GetOrCreateObligation(schoolID, teacherID, academicYear string) (*TeacherHourObligation, error) {
	if academicYear == "" {
		academicYear = "2025/2026"
	}
	schoolID = r.resolveSchoolID(schoolID, teacherID)
	query := `
		INSERT INTO teacher_hour_obligations (school_id, teacher_id, academic_year, collegio_planned_hours, consigli_planned_hours)
		VALUES ($1, $2, $3, 40.00, 40.00)
		ON CONFLICT (school_id, teacher_id, academic_year)
		DO UPDATE SET updated_at = NOW()
		RETURNING id, school_id, teacher_id, academic_year, collegio_planned_hours, consigli_planned_hours,
		          eccedenti_hours, disposizione_hours, banca_ore_balance, created_at, updated_at
	`
	var o TeacherHourObligation
	err := r.db.QueryRow(query, schoolID, teacherID, academicYear).Scan(
		&o.ID, &o.SchoolID, &o.TeacherID, &o.AcademicYear, &o.CollegioPlannedHours, &o.ConsigliPlannedHours,
		&o.EccedentiHours, &o.DisposizioneHours, &o.BancaOreBalance, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get or create obligation: %w", err)
	}

	// Calculate done hours
	r.enrichObligationWithLogs(&o)
	return &o, nil
}

func (r *repository) GetObligation(obligationID string) (*TeacherHourObligation, error) {
	query := `
		SELECT id, school_id, teacher_id, academic_year, collegio_planned_hours, consigli_planned_hours,
		       eccedenti_hours, disposizione_hours, banca_ore_balance, created_at, updated_at
		FROM teacher_hour_obligations
		WHERE id = $1
	`
	var o TeacherHourObligation
	err := r.db.QueryRow(query, obligationID).Scan(
		&o.ID, &o.SchoolID, &o.TeacherID, &o.AcademicYear, &o.CollegioPlannedHours, &o.ConsigliPlannedHours,
		&o.EccedentiHours, &o.DisposizioneHours, &o.BancaOreBalance, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	r.enrichObligationWithLogs(&o)
	return &o, nil
}

func (r *repository) ListObligationsBySchool(schoolID, academicYear string) ([]TeacherHourObligation, error) {
	query := `
		SELECT id, school_id, teacher_id, academic_year, collegio_planned_hours, consigli_planned_hours,
		       eccedenti_hours, disposizione_hours, banca_ore_balance, created_at, updated_at
		FROM teacher_hour_obligations
		WHERE (school_id::text = $1 OR $1 = '') AND ($2 = '' OR academic_year = $2)
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, schoolID, academicYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TeacherHourObligation
	for rows.Next() {
		var o TeacherHourObligation
		if err := rows.Scan(
			&o.ID, &o.SchoolID, &o.TeacherID, &o.AcademicYear, &o.CollegioPlannedHours, &o.ConsigliPlannedHours,
			&o.EccedentiHours, &o.DisposizioneHours, &o.BancaOreBalance, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		r.enrichObligationWithLogs(&o)
		result = append(result, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) enrichObligationWithLogs(o *TeacherHourObligation) {
	query := `
		SELECT category, COALESCE(SUM(hours), 0)
		FROM teacher_hour_logs
		WHERE obligation_id = $1
		GROUP BY category
	`
	rows, err := r.db.Query(query, o.ID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cat string
		var sum float64
		if err := rows.Scan(&cat, &sum); err == nil {
			switch cat {
			case "collegio":
				o.CollegioDoneHours = sum
			case "consigli":
				o.ConsigliDoneHours = sum
			}
		}
	}
	if err := rows.Err(); err != nil {
		return
	}
	o.CollegioRemainingHours = o.CollegioPlannedHours - o.CollegioDoneHours
	o.ConsigliRemainingHours = o.ConsigliPlannedHours - o.ConsigliDoneHours
	o.HasOvertime = (o.CollegioDoneHours > o.CollegioPlannedHours) || (o.ConsigliDoneHours > o.ConsigliPlannedHours)
	o.HasDeficit = (o.CollegioRemainingHours > 0) || (o.ConsigliRemainingHours > 0)
}

func (r *repository) LogHourActivity(log *TeacherHourLog) error {
	query := `
		INSERT INTO teacher_hour_logs (obligation_id, teacher_id, category, hours, activity_date, meeting_id, description, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query, log.ObligationID, log.TeacherID, log.Category, log.Hours, log.ActivityDate, log.MeetingID, log.Description, log.Status,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *repository) ListLogsByObligation(obligationID string) ([]TeacherHourLog, error) {
	query := `
		SELECT id, obligation_id, teacher_id, category, hours, TO_CHAR(activity_date, 'YYYY-MM-DD'), meeting_id, description, status, created_at
		FROM teacher_hour_logs
		WHERE obligation_id = $1
		ORDER BY activity_date DESC, created_at DESC
	`
	rows, err := r.db.Query(query, obligationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []TeacherHourLog
	for rows.Next() {
		var l TeacherHourLog
		if err := rows.Scan(&l.ID, &l.ObligationID, &l.TeacherID, &l.Category, &l.Hours, &l.ActivityDate, &l.MeetingID, &l.Description, &l.Status, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *repository) CreateShortLeave(leave *TeacherShortLeave) error {
	leave.SchoolID = r.resolveSchoolID(leave.SchoolID, leave.TeacherID)
	reqDate, err := time.Parse("2006-01-02", leave.RequestDate)
	if err != nil {
		reqDate = time.Now()
	}
	// CCNL: recovery within 2 months (60 days)
	deadline := reqDate.AddDate(0, 2, 0).Format("2006-01-02")
	leave.RecoveryDeadline = deadline

	query := `
		INSERT INTO teacher_short_leaves (school_id, teacher_id, request_date, hours, reason, status, recovery_deadline, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, recovery_deadline, created_at, updated_at
	`
	return r.db.QueryRow(
		query, leave.SchoolID, leave.TeacherID, leave.RequestDate, leave.Hours, leave.Reason, leave.Status, leave.RecoveryDeadline, leave.Notes,
	).Scan(&leave.ID, &leave.RecoveryDeadline, &leave.CreatedAt, &leave.UpdatedAt)
}

func (r *repository) ListShortLeaves(schoolID, teacherID string) ([]TeacherShortLeave, error) {
	schoolID = r.resolveSchoolID(schoolID, teacherID)
	query := `
		SELECT id, school_id, teacher_id, TO_CHAR(request_date, 'YYYY-MM-DD'), hours, reason, status,
		       TO_CHAR(recovery_deadline, 'YYYY-MM-DD'), recovered_hours,
		       CASE WHEN recovered_date IS NOT NULL THEN TO_CHAR(recovered_date, 'YYYY-MM-DD') ELSE NULL END,
		       COALESCE(notes, ''), created_at, updated_at
		FROM teacher_short_leaves
		WHERE (school_id::text = $1 OR $1 = '') AND ($2 = '' OR teacher_id::text = $2)
		ORDER BY request_date DESC, created_at DESC
	`
	rows, err := r.db.Query(query, schoolID, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaves []TeacherShortLeave
	today := time.Now().Format("2006-01-02")
	for rows.Next() {
		var l TeacherShortLeave
		if err := rows.Scan(
			&l.ID, &l.SchoolID, &l.TeacherID, &l.RequestDate, &l.Hours, &l.Reason, &l.Status,
			&l.RecoveryDeadline, &l.RecoveredHours, &l.RecoveredDate, &l.Notes, &l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if l.Status != "recovered" && l.RecoveryDeadline < today {
			l.IsExpired = true
		}
		leaves = append(leaves, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return leaves, nil
}

func (r *repository) UpdateShortLeaveStatus(id, status string, notes string) error {
	query := `
		UPDATE teacher_short_leaves
		SET status = $2, notes = CASE WHEN $3 <> '' THEN $3 ELSE notes END, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, status, notes)
	return err
}

func (r *repository) RecordShortLeaveRecovery(id string, recoveredHours float64, recoveredDate string) error {
	query := `
		UPDATE teacher_short_leaves
		SET recovered_hours = recovered_hours + $2,
		    recovered_date = $3,
		    status = CASE WHEN (recovered_hours + $2) >= hours THEN 'recovered' ELSE status END,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, recoveredHours, recoveredDate)
	return err
}
