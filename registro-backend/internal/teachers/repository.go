package teachers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context, schoolID string) ([]Teacher, error)
	Get(ctx context.Context, id string) (*Teacher, error)
	GetByUserID(ctx context.Context, userID string) (*Teacher, error)
	GetBySubject(ctx context.Context, subjectID string) ([]Teacher, error)

	GetSubjects(ctx context.Context, teacherID string) ([]TeacherSubject, error)
	AssignSubject(ctx context.Context, teacherID, subjectID string) error
	RemoveSubject(ctx context.Context, teacherID, subjectID string) error
	GetDashboardStats(ctx context.Context, teacherUserID string) (map[string]interface{}, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context, schoolID string) ([]Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.school_id, t.hiring_date, t.qualification, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		WHERE t.school_id = $1
		ORDER BY u.last_name, u.first_name
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []Teacher
	for rows.Next() {
		var t Teacher
		var hDate sql.NullTime
		var qual sql.NullString
		err := rows.Scan(
			&t.ID, &t.UserID, &t.SchoolID, &hDate, &qual, &t.CreatedAt, &t.UpdatedAt,
			&t.FirstName, &t.LastName, &t.Email,
		)
		if err != nil {
			return nil, err
		}
		if hDate.Valid {
			d := hDate.Time.Format("2006-01-02")
			t.HiringDate = &d
		}
		t.Qualification = qual.String
		teachers = append(teachers, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return teachers, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.school_id, t.hiring_date, t.qualification, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = $1
	`
	var t Teacher
	var hDate sql.NullTime
	var qual sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.UserID, &t.SchoolID, &hDate, &qual, &t.CreatedAt, &t.UpdatedAt,
		&t.FirstName, &t.LastName, &t.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("teacher not found")
		}
		return nil, err
	}
	if hDate.Valid {
		d := hDate.Time.Format("2006-01-02")
		t.HiringDate = &d
	}
	t.Qualification = qual.String
	return &t, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*Teacher, error) {
	query := `SELECT id FROM teachers WHERE user_id = $1::uuid`
	var id string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *PostgresRepository) GetSubjects(ctx context.Context, teacherID string) ([]TeacherSubject, error) {
	query := `
		SELECT ts.id, ts.teacher_id, ts.subject_id, s.name, ts.created_at
		FROM teacher_subjects ts
		JOIN subjects s ON ts.subject_id = s.id
		WHERE ts.teacher_id = $1
		ORDER BY s.name
	`
	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TeacherSubject
	for rows.Next() {
		var ts TeacherSubject
		err := rows.Scan(&ts.ID, &ts.TeacherID, &ts.SubjectID, &ts.SubjectName, &ts.CreatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, ts)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *PostgresRepository) GetBySubject(ctx context.Context, subjectID string) ([]Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.school_id, t.hiring_date, t.qualification, t.created_at, t.updated_at,
		       u.first_name, u.last_name, u.email
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		JOIN teacher_subjects ts ON t.id = ts.teacher_id
		WHERE ts.subject_id = $1
		ORDER BY u.last_name, u.first_name
	`
	rows, err := r.db.QueryContext(ctx, query, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []Teacher
	for rows.Next() {
		var t Teacher
		var hDate sql.NullTime
		var qual sql.NullString
		err := rows.Scan(
			&t.ID, &t.UserID, &t.SchoolID, &hDate, &qual, &t.CreatedAt, &t.UpdatedAt,
			&t.FirstName, &t.LastName, &t.Email,
		)
		if err != nil {
			return nil, err
		}
		if hDate.Valid {
			d := hDate.Time.Format("2006-01-02")
			t.HiringDate = &d
		}
		t.Qualification = qual.String
		teachers = append(teachers, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return teachers, nil
}

func (r *PostgresRepository) AssignSubject(ctx context.Context, teacherID, subjectID string) error {
	id := uuid.New().String()
	query := `INSERT INTO teacher_subjects (id, teacher_id, subject_id, created_at)
	          VALUES ($1, $2, $3, NOW())
			  ON CONFLICT (teacher_id, subject_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, id, teacherID, subjectID)
	return err
}

func (r *PostgresRepository) RemoveSubject(ctx context.Context, teacherID, subjectID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM teacher_subjects WHERE teacher_id = $1::uuid AND subject_id = $2::uuid", teacherID, subjectID)
	return err
}

func (r *PostgresRepository) GetDashboardStats(ctx context.Context, teacherUserID string) (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"classes_count":        int64(0),
		"students_count":       int64(0),
		"lessons_today_count":  int64(0),
		"grades_pending_count": int64(0),
	}

	if teacherUserID == "" {
		return stats, nil
	}
	if _, err := uuid.Parse(teacherUserID); err != nil {
		return stats, nil
	}

	// 1. Classes Count
	var classesCount int64
	classesQuery := `
		SELECT COUNT(DISTINCT class_id) FROM (
			SELECT cs.class_id 
			FROM class_subjects cs
			LEFT JOIN teachers t ON (cs.teacher_id = t.id OR cs.teacher_id = t.user_id)
			WHERE t.user_id = NULLIF($1, '')::uuid OR cs.teacher_id = NULLIF($1, '')::uuid
			UNION
			SELECT id AS class_id
			FROM classes 
			WHERE coordinator_id = NULLIF($1, '')::uuid
		) AS temp`
	if err := r.db.QueryRowContext(ctx, classesQuery, teacherUserID).Scan(&classesCount); err == nil {
		stats["classes_count"] = classesCount
	}

	// 2. Students Count
	var studentsCount int64
	studentsQuery := `
		SELECT COUNT(DISTINCT s.id)
		FROM students s
		WHERE s.class_id IN (
			SELECT cs.class_id 
			FROM class_subjects cs
			LEFT JOIN teachers t ON (cs.teacher_id = t.id OR cs.teacher_id = t.user_id)
			WHERE t.user_id = NULLIF($1, '')::uuid OR cs.teacher_id = NULLIF($1, '')::uuid
			UNION
			SELECT id 
			FROM classes 
			WHERE coordinator_id = NULLIF($1, '')::uuid
		)`
	if err := r.db.QueryRowContext(ctx, studentsQuery, teacherUserID).Scan(&studentsCount); err == nil {
		stats["students_count"] = studentsCount
	}

	// 3. Lessons Today Count
	var lessonsTodayCount int64
	weekday := int(time.Now().Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday
	}
	lessonsQuery := `
		SELECT COUNT(*)
		FROM class_schedules
		WHERE (teacher_id = NULLIF($1, '')::uuid OR teacher_id IN (SELECT id FROM teachers WHERE user_id = NULLIF($1, '')::uuid))
		  AND day_of_week = $2`
	if err := r.db.QueryRowContext(ctx, lessonsQuery, teacherUserID, weekday).Scan(&lessonsTodayCount); err == nil {
		stats["lessons_today_count"] = lessonsTodayCount
	}

	// 4. Grades Pending Count (Voti da inserire)
	var gradesPendingCount int64
	gradesQuery := `
		SELECT COUNT(*)
		FROM class_tests ct
		JOIN students s ON ct.class_id = s.class_id
		LEFT JOIN grades g ON g.test_id = ct.id AND g.student_id = s.id AND g.deleted_at IS NULL
		WHERE (ct.teacher_id = NULLIF($1, '')::uuid OR ct.teacher_id IN (SELECT id FROM teachers WHERE user_id = NULLIF($1, '')::uuid))
		  AND ct.date <= CURRENT_DATE AND g.id IS NULL`
	if err := r.db.QueryRowContext(ctx, gradesQuery, teacherUserID).Scan(&gradesPendingCount); err == nil {
		stats["grades_pending_count"] = gradesPendingCount
	}

	return stats, nil
}
