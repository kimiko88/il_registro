package recovery

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	CreateCourse(ctx context.Context, c *RecoveryCourse, sessions []RecoveryCourseSession, studentIDs []string) error
	GetCourseByID(ctx context.Context, id string) (*RecoveryCourse, error)
	ListCourses(ctx context.Context, schoolID, academicYear, teacherID string) ([]RecoveryCourse, error)
	UpdateCourseStatus(ctx context.Context, id, status string) error
	UpdateStudentAttendance(ctx context.Context, courseID, studentID string, hours float64, notes string) error

	RecordRecoveryTest(ctx context.Context, test *RecoveryTest) error
	ListRecoveryTests(ctx context.Context, schoolID, classID, studentID string) ([]RecoveryTest, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateCourse(ctx context.Context, c *RecoveryCourse, sessions []RecoveryCourseSession, studentIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryCourse := `
		INSERT INTO recovery_courses (
			id, school_id, subject_id, teacher_id, title, description,
			academic_year, period, total_hours, room, status, created_at, updated_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, NOW(), NOW()
		) RETURNING id
	`
	err = tx.QueryRowContext(ctx, queryCourse,
		c.ID, c.SchoolID, c.SubjectID, c.TeacherID, c.Title, c.Description,
		c.AcademicYear, c.Period, c.TotalHours, c.Room, c.Status,
	).Scan(&c.ID)
	if err != nil {
		return fmt.Errorf("failed to insert recovery course: %w", err)
	}

	for _, s := range sessions {
		querySession := `
			INSERT INTO recovery_course_sessions (course_id, session_date, start_time, end_time, room, topic, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
		`
		if _, err := tx.ExecContext(ctx, querySession, c.ID, s.SessionDate, s.StartTime, s.EndTime, s.Room, s.Topic); err != nil {
			return fmt.Errorf("failed to insert session: %w", err)
		}
	}

	for _, stID := range studentIDs {
		queryStudent := `
			INSERT INTO recovery_course_students (course_id, student_id, attendance_hours, notes, created_at)
			VALUES ($1, $2, 0.0, '', NOW())
			ON CONFLICT (course_id, student_id) DO NOTHING
		`
		if _, err := tx.ExecContext(ctx, queryStudent, c.ID, stID); err != nil {
			return fmt.Errorf("failed to enroll student: %w", err)
		}
	}

	return tx.Commit()
}

func (r *repository) GetCourseByID(ctx context.Context, id string) (*RecoveryCourse, error) {
	query := `
		SELECT rc.id, rc.school_id, rc.subject_id, sub.name as subject_name,
		       rc.teacher_id, TRIM(COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '')) as teacher_name,
		       rc.title, rc.description, rc.academic_year, rc.period, rc.total_hours, rc.room, rc.status,
		       rc.created_at, rc.updated_at
		FROM recovery_courses rc
		JOIN subjects sub ON rc.subject_id = sub.id
		JOIN teachers t ON rc.teacher_id = t.id
		JOIN users u ON t.user_id = u.id
		WHERE rc.id = $1
	`
	var rc RecoveryCourse
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rc.ID, &rc.SchoolID, &rc.SubjectID, &rc.SubjectName,
		&rc.TeacherID, &rc.TeacherName,
		&rc.Title, &rc.Description, &rc.AcademicYear, &rc.Period, &rc.TotalHours, &rc.Room, &rc.Status,
		&rc.CreatedAt, &rc.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Fetch sessions
	sessRows, err := r.db.QueryContext(ctx, `
		SELECT id, course_id, session_date::text, start_time, end_time, room, topic, created_at
		FROM recovery_course_sessions WHERE course_id = $1 ORDER BY session_date ASC, start_time ASC
	`, id)
	if err != nil {
		return nil, err
	}
	defer sessRows.Close()
	for sessRows.Next() {
		var s RecoveryCourseSession
		if err := sessRows.Scan(&s.ID, &s.CourseID, &s.SessionDate, &s.StartTime, &s.EndTime, &s.Room, &s.Topic, &s.CreatedAt); err != nil {
			return nil, err
		}
		rc.Sessions = append(rc.Sessions, s)
	}
	if err := sessRows.Err(); err != nil {
		return nil, err
	}

	// Fetch students
	stRows, err := r.db.QueryContext(ctx, `
		SELECT rcs.id, rcs.course_id, rcs.student_id,
		       TRIM(COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '')) as student_name,
		       cls.name as class_name, rcs.attendance_hours, rcs.notes
		FROM recovery_course_students rcs
		JOIN students s ON rcs.student_id = s.id
		JOIN users u ON s.user_id = u.id
		LEFT JOIN classes cls ON s.class_id = cls.id
		WHERE rcs.course_id = $1
		ORDER BY u.last_name ASC
	`, id)
	if err != nil {
		return nil, err
	}
	defer stRows.Close()
	for stRows.Next() {
		var st RecoveryCourseStudent
		if err := stRows.Scan(&st.ID, &st.CourseID, &st.StudentID, &st.StudentName, &st.ClassName, &st.AttendanceHours, &st.Notes); err != nil {
			return nil, err
		}
		rc.Students = append(rc.Students, st)
	}
	if err := stRows.Err(); err != nil {
		return nil, err
	}

	return &rc, nil

}

func (r *repository) ListCourses(ctx context.Context, schoolID, academicYear, teacherID string) ([]RecoveryCourse, error) {
	query := `
		SELECT rc.id, rc.school_id, rc.subject_id, sub.name as subject_name,
		       rc.teacher_id, TRIM(COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '')) as teacher_name,
		       rc.title, rc.description, rc.academic_year, rc.period, rc.total_hours, rc.room, rc.status,
		       rc.created_at, rc.updated_at
		FROM recovery_courses rc
		JOIN subjects sub ON rc.subject_id = sub.id
		JOIN teachers t ON rc.teacher_id = t.id
		JOIN users u ON t.user_id = u.id
		WHERE rc.school_id = $1
		  AND ($2 = '' OR rc.academic_year = $2)
		  AND ($3 = '' OR rc.teacher_id = NULLIF($3, '')::uuid)
		ORDER BY rc.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, academicYear, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []RecoveryCourse
	for rows.Next() {
		var rc RecoveryCourse
		if err := rows.Scan(
			&rc.ID, &rc.SchoolID, &rc.SubjectID, &rc.SubjectName,
			&rc.TeacherID, &rc.TeacherName,
			&rc.Title, &rc.Description, &rc.AcademicYear, &rc.Period, &rc.TotalHours, &rc.Room, &rc.Status,
			&rc.CreatedAt, &rc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, rc)
	}
	return list, rows.Err()
}

func (r *repository) UpdateCourseStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE recovery_courses SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}

func (r *repository) UpdateStudentAttendance(ctx context.Context, courseID, studentID string, hours float64, notes string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE recovery_course_students 
		SET attendance_hours = $1, notes = $2 
		WHERE course_id = $3 AND student_id = $4
	`, hours, notes, courseID, studentID)
	return err
}

func (r *repository) RecordRecoveryTest(ctx context.Context, t *RecoveryTest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO recovery_tests (
			id, school_id, deficiency_id, student_id, subject_id, class_id,
			teacher_id, test_date, test_type, grade, outcome, final_deliberation,
			verbale_number, notes, created_at, updated_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, NULLIF($3, '')::uuid, $4, $5, $6,
			$7, $8, $9, $10, $11, $12,
			$13, $14, NOW(), NOW()
		) RETURNING id
	`
	err = tx.QueryRowContext(ctx, query,
		t.ID, t.SchoolID, t.DeficiencyID, t.StudentID, t.SubjectID, t.ClassID,
		t.TeacherID, t.TestDate, t.TestType, t.Grade, t.Outcome, t.FinalDeliberation,
		t.VerbaleNumber, t.Notes,
	).Scan(&t.ID)
	if err != nil {
		return fmt.Errorf("failed to record recovery test: %w", err)
	}

	// Update corresponding deficiency status if linked
	if t.DeficiencyID != nil && *t.DeficiencyID != "" {
		defStatus := "recuperato"
		if t.Outcome == "non_recuperato" {
			defStatus = "non_recuperato"
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE student_deficiencies
			SET status = $1, recovery_grade = $2, recovery_date = $3::date, updated_at = NOW()
			WHERE id = $4
		`, defStatus, t.Grade, t.TestDate, *t.DeficiencyID); err != nil {
			return fmt.Errorf("failed to update deficiency status: %w", err)
		}
	}

	return tx.Commit()

}

func (r *repository) ListRecoveryTests(ctx context.Context, schoolID, classID, studentID string) ([]RecoveryTest, error) {
	query := `
		SELECT rt.id, rt.school_id, rt.deficiency_id::text, rt.student_id,
		       TRIM(COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '')) as student_name,
		       rt.subject_id, sub.name as subject_name,
		       rt.class_id, cls.name as class_name,
		       rt.teacher_id, TRIM(COALESCE(tu.first_name, '') || ' ' || COALESCE(tu.last_name, '')) as teacher_name,
		       rt.test_date::text, rt.test_type, rt.grade, rt.outcome, rt.final_deliberation,
		       rt.verbale_number, rt.notes, rt.created_at, rt.updated_at
		FROM recovery_tests rt
		JOIN students s ON rt.student_id = s.id
		JOIN users u ON s.user_id = u.id
		JOIN subjects sub ON rt.subject_id = sub.id
		JOIN classes cls ON rt.class_id = cls.id
		JOIN teachers t ON rt.teacher_id = t.id
		JOIN users tu ON t.user_id = tu.id
		WHERE rt.school_id = $1
		  AND ($2 = '' OR rt.class_id = NULLIF($2, '')::uuid)
		  AND ($3 = '' OR rt.student_id = NULLIF($3, '')::uuid)
		ORDER BY rt.test_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, classID, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []RecoveryTest
	for rows.Next() {
		var rt RecoveryTest
		var defID sql.NullString
		if err := rows.Scan(
			&rt.ID, &rt.SchoolID, &defID, &rt.StudentID, &rt.StudentName,
			&rt.SubjectID, &rt.SubjectName,
			&rt.ClassID, &rt.ClassName,
			&rt.TeacherID, &rt.TeacherName,
			&rt.TestDate, &rt.TestType, &rt.Grade, &rt.Outcome, &rt.FinalDeliberation,
			&rt.VerbaleNumber, &rt.Notes, &rt.CreatedAt, &rt.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if defID.Valid {
			rt.DeficiencyID = &defID.String
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}
