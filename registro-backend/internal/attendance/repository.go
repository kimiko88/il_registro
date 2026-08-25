package attendance

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Repository interface {
	Create(att *Attendance) error
	BatchCreate(atts []*Attendance) error
	Update(att *Attendance) error
	DeleteByClassDateHour(schoolID, classID string, date time.Time, hour int) error

	FindByID(id string) (*Attendance, error)
	FindByClassAndDate(classID string, date time.Time) ([]Attendance, error)
	FindByStudent(studentID string, startDate, endDate time.Time) ([]Attendance, error)
	GetStats(studentID string) (*SummaryResponse, error)
	GetStatsBatch(ctx context.Context, studentIDs []string) (map[string]*SummaryResponse, error)
	CountDistinctDays(studentID string) (int, error)
	GetAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error)

	// Justifications
	CreateJustification(j *Justification) error
	UpdateJustification(j *Justification) error
	ProcessJustificationTx(ctx context.Context, j *Justification, teacherID string, approve bool) error
	FindJustificationByID(id string) (*Justification, error)
	FindPendingJustifications(classID, schoolID string) ([]Justification, error)
	FindPendingJustificationsForTeacher(ctx context.Context, teacherID, schoolID string) ([]Justification, error)
	DeleteJustification(id string) error
	DeletePendingJustification(id string) error

	// Teacher & Student Assignment Checks
	IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error)
	IsStudentInClass(ctx context.Context, studentID, classID string) (bool, error)
	AreStudentsInClass(ctx context.Context, studentIDs []string, classID string) (map[string]bool, error)
	IsClassInSchool(ctx context.Context, classID, schoolID string) (bool, error)
	IsTeacherSubstitute(ctx context.Context, teacherID, classID string, date time.Time, hour int) (bool, error)
	HasOverlappingJustification(ctx context.Context, studentID string, startDate, endDate time.Time) (bool, error)

	// Monthly Breakdown
	GetMonthlyBreakdown(ctx context.Context, studentID, schoolYear string) ([]MonthlyBreakdownRow, error)

	FindUnjustifiedByStudent(studentID string) ([]Attendance, error)
	JustifyAbsenceByParent(attendanceID string, studentID string, reason string, notes string) error
	GetStudentAttendanceStats(studentID string) (*AttendanceStats, error)
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
	if len(atts) == 0 {
		return nil
	}
	// Batch processing in chunks of 100 to avoid lock escalation
	chunkSize := 100
	for i := 0; i < len(atts); i += chunkSize {
		end := i + chunkSize
		if end > len(atts) {
			end = len(atts)
		}
		chunk := atts[i:end]

		tx, err := r.db.Begin()
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare(`
			INSERT INTO attendance (
				school_id, student_id, class_id, date, hour, subject_id, status, 
				justified, justified_by, justified_at, notes, entry_time, exit_time, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())
			ON CONFLICT (student_id, class_id, date, hour, school_id) DO UPDATE SET
				status = EXCLUDED.status,
				subject_id = COALESCE(EXCLUDED.subject_id, attendance.subject_id),
				entry_time = EXCLUDED.entry_time,
				exit_time = EXCLUDED.exit_time,
				notes = EXCLUDED.notes,
				updated_at = NOW()
			RETURNING id
		`)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		for _, a := range chunk {
			err := stmt.QueryRow(
				a.SchoolID, a.StudentID, a.ClassID, a.Date, a.Hour, a.SubjectID, a.Status,
				a.Justified, a.JustifiedBy, a.JustifiedAt, a.Notes, a.EntryTime, a.ExitTime,
			).Scan(&a.ID)
			if err != nil {
				stmt.Close()
				_ = tx.Rollback()
				return fmt.Errorf("batch insert error: %w", err)
			}
		}
		stmt.Close()
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Update(a *Attendance) error {
	query := `
		UPDATE attendance SET 
			status=$1, hour=$2, subject_id=$3, notes=$4, justified=$5, justified_by=$6, justified_at=$7, entry_time=$8, exit_time=$9, updated_at=NOW()
		WHERE id=$10::uuid`
	_, err := r.db.Exec(query, a.Status, a.Hour, a.SubjectID, a.Notes, a.Justified, a.JustifiedBy, a.JustifiedAt, a.EntryTime, a.ExitTime, a.ID)
	return err
}

func (r *repository) DeleteByClassDateHour(schoolID, classID string, date time.Time, hour int) error {
	if schoolID != "" {
		_, err := r.db.Exec(
			`DELETE FROM attendance WHERE school_id = $1::uuid AND class_id = $2::uuid AND date = $3 AND hour = $4`,
			schoolID, classID, date, hour,
		)
		return err
	}
	_, err := r.db.Exec(
		`DELETE FROM attendance WHERE class_id = $1::uuid AND date = $2 AND hour = $3`,
		classID, date, hour,
	)
	return err
}

func (r *repository) FindByID(id string) (*Attendance, error) {
	query := `SELECT id, school_id, student_id, class_id, date, hour, subject_id, status, justified, justified_by, justified_at, COALESCE(notes, ''), entry_time, exit_time FROM attendance WHERE id=$1::uuid`
	var a Attendance
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.SchoolID, &a.StudentID, &a.ClassID, &a.Date, &a.Hour, &a.SubjectID, &a.Status,
		&a.Justified, &a.JustifiedBy, &a.JustifiedAt, &a.Notes, &a.EntryTime, &a.ExitTime,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) FindByClassAndDate(classID string, date time.Time) ([]Attendance, error) {
	query := `
		SELECT id, school_id, student_id, class_id, date, hour, subject_id, status, justified, justified_by, justified_at, COALESCE(notes, ''), entry_time, exit_time
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
			&a.ID, &a.SchoolID, &a.StudentID, &a.ClassID, &a.Date, &a.Hour, &a.SubjectID, &a.Status,
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
		SELECT id, school_id, student_id, class_id, date, hour, subject_id, status, justified, justified_by, justified_at, COALESCE(notes, ''), entry_time, exit_time
		FROM attendance 
		WHERE (
			student_id = $1::uuid OR
			student_id IN (SELECT user_id FROM students WHERE id = $1::uuid) OR
			student_id IN (SELECT id FROM students WHERE user_id = $1::uuid)
		) AND date BETWEEN $2 AND $3
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
			&a.ID, &a.SchoolID, &a.StudentID, &a.ClassID, &a.Date, &a.Hour, &a.SubjectID, &a.Status,
			&a.Justified, &a.JustifiedBy, &a.JustifiedAt, &a.Notes, &a.EntryTime, &a.ExitTime,
		); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, nil
}

func (r *repository) GetStats(studentID string) (*SummaryResponse, error) {
	if studentID == "" {
		return &SummaryResponse{}, nil
	}
	query := `
		SELECT 
			COUNT(DISTINCT date) FILTER (WHERE status::text IN ('Absent', 'absent')) as absences,
			COUNT(*) FILTER (WHERE status::text IN ('Late', 'late')) as lates,
			COUNT(*) FILTER (WHERE status::text IN ('LeftEarly', 'left_early')) as early_exits,
			COUNT(*) FILTER (WHERE COALESCE(justified, false) = true OR COALESCE(parent_justified, false) = true) as justified
		FROM attendance
		WHERE (
			student_id = $1::uuid OR
			student_id IN (SELECT user_id FROM students WHERE id = $1::uuid) OR
			student_id IN (SELECT id FROM students WHERE user_id = $1::uuid)
		) AND deleted_at IS NULL`

	var s SummaryResponse
	err := r.db.QueryRow(query, studentID).Scan(&s.TotalAbsences, &s.TotalLates, &s.TotalEarlyExits, &s.JustifiedCount)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) GetStatsBatch(ctx context.Context, studentIDs []string) (map[string]*SummaryResponse, error) {
	result := make(map[string]*SummaryResponse)
	if len(studentIDs) == 0 {
		return result, nil
	}
	query := `
		SELECT 
			student_id::text,
			COUNT(DISTINCT date) FILTER (WHERE status::text IN ('Absent', 'absent')) as absences,
			COUNT(*) FILTER (WHERE status::text IN ('Late', 'late')) as lates,
			COUNT(*) FILTER (WHERE status::text IN ('LeftEarly', 'left_early')) as early_exits,
			COUNT(*) FILTER (WHERE COALESCE(justified, false) = true OR COALESCE(parent_justified, false) = true) as justified
		FROM attendance
		WHERE (
			student_id = ANY($1::uuid[]) OR
			student_id IN (SELECT user_id FROM students WHERE id = ANY($1::uuid[])) OR
			student_id IN (SELECT id FROM students WHERE user_id = ANY($1::uuid[]))
		) AND deleted_at IS NULL
		GROUP BY student_id`

	rows, err := r.db.QueryContext(ctx, query, studentIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sID string
		var s SummaryResponse
		if err := rows.Scan(&sID, &s.TotalAbsences, &s.TotalLates, &s.TotalEarlyExits, &s.JustifiedCount); err == nil {
			result[sID] = &s
		}
	}
	return result, rows.Err()
}

func (r *repository) CountDistinctDays(studentID string) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(DISTINCT date) FROM attendance WHERE student_id=$1::uuid`, studentID).Scan(&count)
	return count, err
}

func (r *repository) GetAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error) {
	var res AnalyticsResponse

	query := `
		SELECT
			ROUND(
				100.0 * COUNT(*) FILTER (WHERE status = 'Absent') /
				NULLIF(COUNT(*), 0),
			2) AS avg_absence_rate
		FROM attendance
		WHERE school_id = $1`
	err := r.db.QueryRowContext(ctx, query, schoolID).Scan(&res.AverageAbsenceRate)
	if err != nil {
		return nil, err
	}

	topQuery := `
		SELECT u.first_name || ' ' || u.last_name
		FROM attendance a
		JOIN users u ON u.id = a.student_id::uuid
		WHERE a.school_id = $1 AND a.status = 'Absent'
		GROUP BY a.student_id, u.first_name, u.last_name
		ORDER BY COUNT(*) DESC
		LIMIT 5`
	rows, err := r.db.QueryContext(ctx, topQuery, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		res.TopAbsentees = append(res.TopAbsentees, name)
	}
	if res.TopAbsentees == nil {
		res.TopAbsentees = []string{}
	}
	return &res, nil
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

func (r *repository) ProcessJustificationTx(ctx context.Context, j *Justification, teacherID string, approve bool) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Lock the justification row FOR UPDATE inside the transaction to prevent race conditions
	var currentStatus JustificationStatus
	if err := tx.QueryRowContext(ctx, `SELECT status FROM justifications WHERE id=$1::uuid FOR UPDATE`, j.ID).Scan(&currentStatus); err != nil {
		return fmt.Errorf("la giustifica %s non esiste: %w", j.ID, err)
	}
	if currentStatus != JustificationPending {
		return fmt.Errorf("la giustifica %s è già stata elaborata (stato attuale: %s)", j.ID, currentStatus)
	}

	now := time.Now()
	if approve {
		j.Status = JustificationApproved
		j.ApprovedBy = &teacherID
		j.ApprovedAt = &now

		queryJust := `UPDATE justifications SET status=$1, approved_by=$2::uuid, approved_at=$3, updated_at=NOW() WHERE id=$4::uuid`
		if _, err := tx.ExecContext(ctx, queryJust, j.Status, j.ApprovedBy, j.ApprovedAt, j.ID); err != nil {
			return fmt.Errorf("failed to update justification: %w", err)
		}

		queryAtt := `
			UPDATE attendance 
			SET justified = true, justified_by = $1::uuid, justified_at = $2, updated_at = NOW() 
			WHERE student_id = $3::uuid AND date >= $4 AND date <= $5 AND status = 'Absent'
		`
		if _, err := tx.ExecContext(ctx, queryAtt, teacherID, now, j.StudentID, j.StartDate, j.EndDate); err != nil {
			return fmt.Errorf("failed to update attendance records: %w", err)
		}
	} else {
		j.Status = JustificationRejected
		queryJust := `UPDATE justifications SET status=$1, updated_at=NOW() WHERE id=$2::uuid`
		if _, err := tx.ExecContext(ctx, queryJust, j.Status, j.ID); err != nil {
			return fmt.Errorf("failed to update justification: %w", err)
		}
	}

	return tx.Commit()
}

func (r *repository) FindJustificationByID(id string) (*Justification, error) {
	var j Justification
	query := `SELECT id, student_id, parent_id, start_date, end_date, reason, status, approved_by, approved_at FROM justifications WHERE id=$1::uuid`
	err := r.db.QueryRow(query, id).Scan(&j.ID, &j.StudentID, &j.ParentID, &j.StartDate, &j.EndDate, &j.Reason, &j.Status, &j.ApprovedBy, &j.ApprovedAt)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *repository) FindPendingJustifications(classID, schoolID string) ([]Justification, error) {
	query := `
		SELECT j.id, j.student_id, COALESCE(u.first_name || ' ' || u.last_name, 'Studente') AS student_name, j.start_date, j.end_date, j.reason, j.status 
		FROM justifications j
		JOIN users u ON j.student_id = u.id::uuid
		LEFT JOIN students s ON s.user_id = u.id::uuid
		WHERE ($1::text = '' OR s.class_id::text = $1::text)
		  AND ($2::text = '' OR u.school_id::text = $2::text OR s.school_id::text = $2::text)
		  AND j.status = 'pending'`

	rows, err := r.db.Query(query, classID, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Justification
	for rows.Next() {
		var j Justification
		if err := rows.Scan(&j.ID, &j.StudentID, &j.StudentName, &j.StartDate, &j.EndDate, &j.Reason, &j.Status); err != nil {
			return nil, err
		}
		res = append(res, j)
	}
	return res, nil
}

func (r *repository) FindPendingJustificationsForTeacher(ctx context.Context, teacherID, schoolID string) ([]Justification, error) {
	query := `
		SELECT j.id, j.student_id, COALESCE(u.first_name || ' ' || u.last_name, 'Studente') AS student_name, j.start_date, j.end_date, j.reason, j.status 
		FROM justifications j
		JOIN users u ON j.student_id = u.id::uuid
		JOIN students s ON s.user_id = u.id::uuid
		WHERE j.status = 'pending'
		  AND ($2::text = '' OR u.school_id::text = $2::text OR s.school_id::text = $2::text)
		  AND (
			EXISTS (
				SELECT 1 FROM classes c
				LEFT JOIN class_subjects cs ON c.id::text = cs.class_id::text
				LEFT JOIN teachers t ON (cs.teacher_id::text = t.id::text OR cs.teacher_id::text = t.user_id::text)
				WHERE c.id::text = s.class_id::text AND (
					t.user_id::text = $1::text OR
					t.id::text = $1::text OR
					cs.teacher_id::text = $1::text OR
					c.coordinator_id::text = $1::text
				)
			)
		  )
		ORDER BY j.start_date DESC`

	rows, err := r.db.QueryContext(ctx, query, teacherID, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Justification
	for rows.Next() {
		var j Justification
		if err := rows.Scan(&j.ID, &j.StudentID, &j.StudentName, &j.StartDate, &j.EndDate, &j.Reason, &j.Status); err != nil {
			return nil, err
		}
		res = append(res, j)
	}
	return res, nil
}

func (r *repository) DeleteJustification(id string) error {
	_, err := r.db.Exec(`DELETE FROM justifications WHERE id=$1::uuid`, id)
	return err
}

func (r *repository) DeletePendingJustification(id string) error {
	res, err := r.db.Exec(`DELETE FROM justifications WHERE id=$1::uuid AND status = 'pending'`, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("impossibile eliminare la giustifica: non trovata o già elaborata")
	}
	return nil
}

func (r *repository) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM classes c
			LEFT JOIN class_subjects cs ON c.id::text = cs.class_id::text
			LEFT JOIN teachers t ON (cs.teacher_id::text = t.id::text OR cs.teacher_id::text = t.user_id::text)
			LEFT JOIN users u ON u.id::text = $2::text
			WHERE c.id::text = $1::text AND (
				t.user_id::text = $2::text OR
				t.id::text = $2::text OR
				cs.teacher_id::text = $2::text OR
				c.coordinator_id::text = $2::text OR
				(u.role IN ('admin', 'superadmin', 'secretary') AND c.school_id = u.school_id) OR
				(COALESCE(u.is_staff, false) = true AND c.school_id = u.school_id)
			)
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, classID, teacherID).Scan(&exists)
	return exists, err
}

// GetMonthlyBreakdown returns per-month attendance counts for a student across a school year.
// schoolYear format: "2024-2025" → from 2024-09-01 to 2025-08-31.
func (r *repository) GetMonthlyBreakdown(ctx context.Context, studentID, schoolYear string) ([]MonthlyBreakdownRow, error) {
	// Parse school year (e.g. "2024-2025")
	var startYear int
	if _, err := fmt.Sscanf(schoolYear, "%d-%*d", &startYear); err != nil {
		// Default to current school year
		now := time.Now()
		if now.Month() >= 9 {
			startYear = now.Year()
		} else {
			startYear = now.Year() - 1
		}
	}

	query := `
		SELECT
			to_char(date, 'YYYY-MM') AS month,
			COUNT(*) FILTER (WHERE status = 'Absent')              AS absences,
			COUNT(*) FILTER (WHERE status = 'Late')                AS lates,
			COUNT(*) FILTER (WHERE status = 'LeftEarly')           AS early_exits,
			COUNT(*) FILTER (WHERE status = 'Absent' AND justified = true) AS justified_absences,
			COUNT(DISTINCT date) AS total_school_days
		FROM attendance
		WHERE student_id = $1::uuid
		  AND date >= make_date($2, 9, 1)
		  AND date < make_date($2 + 1, 9, 1)
		GROUP BY month
		ORDER BY month
	`

	rows, err := r.db.QueryContext(ctx, query, studentID, startYear)
	if err != nil {
		return nil, fmt.Errorf("GetMonthlyBreakdown query: %w", err)
	}
	defer rows.Close()

	var result []MonthlyBreakdownRow
	for rows.Next() {
		var row MonthlyBreakdownRow
		if err := rows.Scan(
			&row.Month,
			&row.Absences,
			&row.Lates,
			&row.EarlyExits,
			&row.JustifiedAbsences,
			&row.TotalSchoolDays,
		); err != nil {
			return nil, fmt.Errorf("GetMonthlyBreakdown scan: %w", err)
		}
		// Compute presence rate: days present / total * 100
		// Compute presence rate: (total - absences) / total * 100
		present := row.TotalSchoolDays - row.Absences
		if present < 0 {
			present = 0
		}
		if row.TotalSchoolDays > 0 {
			row.PresenceRate = float64(present) / float64(row.TotalSchoolDays) * 100
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) FindUnjustifiedByStudent(studentID string) ([]Attendance, error) {
	query := `
		SELECT id, school_id, student_id, class_id, date, hour, subject_id, status, justified, justified_by, justified_at, COALESCE(notes, ''), entry_time, exit_time
		FROM attendance
		WHERE (
			student_id = $1::uuid OR
			student_id IN (SELECT user_id FROM students WHERE id = $1::uuid) OR
			student_id IN (SELECT id FROM students WHERE user_id = $1::uuid)
		) AND status IN ('Absent', 'Late', 'LeftEarly')
		  AND COALESCE(justified, false) = false AND COALESCE(parent_justified, false) = false
		ORDER BY date DESC
	`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Attendance
	for rows.Next() {
		var a Attendance
		if err := rows.Scan(
			&a.ID, &a.SchoolID, &a.StudentID, &a.ClassID, &a.Date, &a.Hour, &a.SubjectID, &a.Status,
			&a.Justified, &a.JustifiedBy, &a.JustifiedAt, &a.Notes, &a.EntryTime, &a.ExitTime,
		); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	if result == nil {
		result = []Attendance{}
	}
	return result, nil
}

func (r *repository) JustifyAbsenceByParent(attendanceID string, studentID string, reason string, notes string) error {
	query := `
		UPDATE attendance
		SET parent_justified = true, parent_justified_at = NOW(),
		    justification_reason = $1, notes = COALESCE($2, notes), updated_at = NOW()
		WHERE id = $3::uuid 
		  AND (
			student_id = $4::uuid OR
			student_id IN (SELECT user_id FROM students WHERE id = $4::uuid) OR
			student_id IN (SELECT id FROM students WHERE user_id = $4::uuid)
		  )
		  AND COALESCE(parent_justified, false) = false
	`
	res, err := r.db.Exec(query, reason, notes, attendanceID, studentID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("forbidden: l'assenza non appartiene allo studente o è già stata giustificata")
	}
	return nil
}

func (r *repository) GetStudentAttendanceStats(studentID string) (*AttendanceStats, error) {
	query := `
		SELECT
			COUNT(*) AS total_records,
			COUNT(*) FILTER (WHERE status = 'Present') AS present,
			COUNT(*) FILTER (WHERE status = 'Absent') AS absent,
			COUNT(*) FILTER (WHERE status = 'Late') AS lates,
			COUNT(*) FILTER (WHERE status = 'LeftEarly') AS early_exits,
			COUNT(*) FILTER (WHERE (COALESCE(justified, false) = true OR COALESCE(parent_justified, false) = true)) AS justified,
			COUNT(*) FILTER (WHERE status IN ('Absent', 'Late', 'LeftEarly') AND COALESCE(justified, false) = false AND COALESCE(parent_justified, false) = false) AS unjustified
		FROM attendance
		WHERE (
			student_id = $1::uuid OR
			student_id IN (SELECT user_id FROM students WHERE id = $1::uuid) OR
			student_id IN (SELECT id FROM students WHERE user_id = $1::uuid)
		)
	`
	stats := &AttendanceStats{}
	err := r.db.QueryRow(query, studentID).Scan(
		&stats.TotalSchoolDays,
		&stats.DaysPresent,
		&stats.DaysAbsent,
		&stats.LateArrivals,
		&stats.EarlyExits,
		&stats.Justified,
		&stats.Unjustified,
	)
	if err != nil {
		return nil, err
	}

	if stats.TotalSchoolDays == 0 {
		stats.AbsencePercentage = 0
		stats.MonthlyBreakdown = []MonthlyAttendance{}
		return stats, nil
	}

	stats.AbsencePercentage = (float64(stats.DaysAbsent) / float64(stats.TotalSchoolDays)) * 100

	mQuery := `
		SELECT
			to_char(date, 'Mon') AS m_name,
			COUNT(*) FILTER (WHERE status = 'Present') AS p_cnt,
			COUNT(*) FILTER (WHERE status = 'Absent') AS a_cnt
		FROM attendance
		WHERE (
			student_id = $1::uuid OR
			student_id IN (SELECT user_id FROM students WHERE id = $1::uuid) OR
			student_id IN (SELECT id FROM students WHERE user_id = $1::uuid)
		)
		GROUP BY to_char(date, 'YYYY-MM'), m_name
		ORDER BY to_char(date, 'YYYY-MM')
	`
	rows, err := r.db.Query(mQuery, studentID)
	if err == nil {
		for rows.Next() {
			var ma MonthlyAttendance
			if err := rows.Scan(&ma.Month, &ma.Present, &ma.Absent); err == nil {
				stats.MonthlyBreakdown = append(stats.MonthlyBreakdown, ma)
			}
		}
		rows.Close()
	}

	if stats.MonthlyBreakdown == nil {
		stats.MonthlyBreakdown = []MonthlyAttendance{}
	}

	return stats, nil
}

func (r *repository) IsTeacherSubstitute(ctx context.Context, teacherID, classID string, date time.Time, hour int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM substitutions
			WHERE (teacher_id = $1::uuid OR substitute_teacher_id = $1::uuid)
			  AND class_id = $2::uuid
			  AND date = $3::date
			  AND (hour = $4 OR hour IS NULL)
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, teacherID, classID, date, hour).Scan(&exists)
	if err != nil {
		return false, nil
	}
	return exists, nil
}

func (r *repository) HasOverlappingJustification(ctx context.Context, studentID string, startDate, endDate time.Time) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM justifications
			WHERE student_id = $1::uuid
			  AND LOWER(status::text) IN ('pending', 'approved')
			  AND start_date <= $3::date
			  AND end_date >= $2::date
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, studentID, startDate, endDate).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *repository) IsStudentInClass(ctx context.Context, studentID, classID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM class_students
			WHERE student_id::text = $1::text AND class_id::text = $2::text
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, studentID, classID).Scan(&exists)
	if err != nil {
		return false, nil
	}
	return exists, nil
}

func (r *repository) AreStudentsInClass(ctx context.Context, studentIDs []string, classID string) (map[string]bool, error) {
	res := make(map[string]bool)
	if len(studentIDs) == 0 || classID == "" {
		return res, nil
	}
	query := `
		SELECT student_id::text FROM class_students
		WHERE class_id::text = $1::text AND student_id::text = ANY($2)
	`
	rows, err := r.db.QueryContext(ctx, query, classID, pq.Array(studentIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sID string
		if err := rows.Scan(&sID); err == nil {
			res[sID] = true
		}
	}
	return res, rows.Err()
}

func (r *repository) IsClassInSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	if classID == "" || schoolID == "" {
		return false, nil
	}
	query := `
		SELECT EXISTS (
			SELECT 1 FROM classes
			WHERE id::text = $1::text AND school_id::text = $2::text
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, classID, schoolID).Scan(&exists)
	if err != nil {
		return false, nil
	}
	return exists, nil
}
