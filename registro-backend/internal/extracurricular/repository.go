package extracurricular

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateCourse(ctx context.Context, c *Course) error
	GetCourseByID(ctx context.Context, id string) (*Course, error)
	ListCourses(ctx context.Context, schoolID, studentID string) ([]*Course, error)

	EnrollStudent(ctx context.Context, courseID, studentID string) error
	ListEnrollments(ctx context.Context, courseID string) ([]*Enrollment, error)

	MarkAttendance(ctx context.Context, att *AttendanceRecord) error
	ListAttendance(ctx context.Context, courseID string, date time.Time) ([]*AttendanceRecord, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateCourse(ctx context.Context, c *Course) error {
	c.ID = uuid.New().String()
	c.CreatedAt = time.Now()
	if c.MaxParticipants <= 0 {
		c.MaxParticipants = 30
	}

	query := `
		INSERT INTO extracurricular_courses (id, school_id, title, description, teacher_id, start_date, end_date, max_participants, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.SchoolID, c.Title, c.Description, c.TeacherID, c.StartDate, c.EndDate, c.MaxParticipants, c.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetCourseByID(ctx context.Context, id string) (*Course, error) {
	query := `
		SELECT c.id, c.school_id, c.title, COALESCE(c.description, ''), c.teacher_id, c.start_date, c.end_date, c.max_participants, c.created_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name,
		       (SELECT COUNT(*) FROM extracurricular_enrollments e WHERE e.course_id = c.id) AS enrolled_count
		FROM extracurricular_courses c
		LEFT JOIN users u ON c.teacher_id = u.id
		WHERE c.id = $1::uuid
	`
	c := &Course{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.SchoolID, &c.Title, &c.Description, &c.TeacherID, &c.StartDate, &c.EndDate, &c.MaxParticipants, &c.CreatedAt,
		&c.TeacherName, &c.EnrolledCount,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *PostgresRepository) ListCourses(ctx context.Context, schoolID, studentID string) ([]*Course, error) {
	stdUUID := studentID
	if stdUUID == "" {
		stdUUID = "00000000-0000-0000-0000-000000000000"
	}
	query := `
		SELECT c.id, c.school_id, c.title, COALESCE(c.description, ''), c.teacher_id, c.start_date, c.end_date, c.max_participants, c.created_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name,
		       (SELECT COUNT(*) FROM extracurricular_enrollments e WHERE e.course_id = c.id) AS enrolled_count,
		       EXISTS(SELECT 1 FROM extracurricular_enrollments e WHERE e.course_id = c.id AND e.student_id = $2::uuid) AS is_student_enrolled
		FROM extracurricular_courses c
		LEFT JOIN users u ON c.teacher_id = u.id
		WHERE c.school_id = $1::uuid
		ORDER BY c.start_date ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, stdUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*Course
	for rows.Next() {
		c := &Course{}
		if err := rows.Scan(
			&c.ID, &c.SchoolID, &c.Title, &c.Description, &c.TeacherID, &c.StartDate, &c.EndDate, &c.MaxParticipants, &c.CreatedAt,
			&c.TeacherName, &c.EnrolledCount, &c.IsStudentEnrolled,
		); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}

func (r *PostgresRepository) EnrollStudent(ctx context.Context, courseID, studentID string) error {
	query := `
		INSERT INTO extracurricular_enrollments (id, course_id, student_id, enrolled_at)
		VALUES ($1, $2::uuid, $3::uuid, NOW())
		ON CONFLICT (course_id, student_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, uuid.New().String(), courseID, studentID)
	return err
}

func (r *PostgresRepository) ListEnrollments(ctx context.Context, courseID string) ([]*Enrollment, error) {
	query := `
		SELECT e.id, e.course_id, e.student_id, e.enrolled_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name
		FROM extracurricular_enrollments e
		JOIN users u ON e.student_id = u.id
		WHERE e.course_id = $1::uuid
		ORDER BY e.enrolled_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Enrollment
	for rows.Next() {
		e := &Enrollment{}
		if err := rows.Scan(&e.ID, &e.CourseID, &e.StudentID, &e.EnrolledAt, &e.StudentName); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) MarkAttendance(ctx context.Context, att *AttendanceRecord) error {
	att.ID = uuid.New().String()
	att.SignedAt = time.Now()

	query := `
		INSERT INTO extracurricular_attendance (id, course_id, student_id, date, status, hours, signed_by, signed_at)
		VALUES ($1, $2::uuid, $3::uuid, $4, $5, $6, $7::uuid, $8)
		ON CONFLICT (course_id, student_id, date) DO UPDATE
		SET status = EXCLUDED.status, hours = EXCLUDED.hours, signed_by = EXCLUDED.signed_by, signed_at = EXCLUDED.signed_at
	`
	var signedByUUID interface{} = nil
	if att.SignedBy != nil && *att.SignedBy != "" {
		signedByUUID = *att.SignedBy
	}
	_, err := r.db.ExecContext(ctx, query,
		att.ID, att.CourseID, att.StudentID, att.Date, att.Status, att.Hours, signedByUUID, att.SignedAt,
	)
	return err
}

func (r *PostgresRepository) ListAttendance(ctx context.Context, courseID string, date time.Time) ([]*AttendanceRecord, error) {
	query := `
		SELECT a.id, a.course_id, a.student_id, a.date, a.status, a.hours, a.signed_by, a.signed_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name
		FROM extracurricular_attendance a
		JOIN users u ON a.student_id = u.id
		WHERE a.course_id = $1::uuid AND a.date = $2
		ORDER BY u.last_name ASC, u.first_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, courseID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*AttendanceRecord
	for rows.Next() {
		a := &AttendanceRecord{}
		var signedBy sql.NullString
		if err := rows.Scan(&a.ID, &a.CourseID, &a.StudentID, &a.Date, &a.Status, &a.Hours, &signedBy, &a.SignedAt, &a.StudentName); err != nil {
			return nil, err
		}
		if signedBy.Valid {
			a.SignedBy = &signedBy.String
		}
		list = append(list, a)
	}
	return list, rows.Err()
}
