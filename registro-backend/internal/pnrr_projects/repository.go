package pnrr_projects

import (
	"database/sql"
)

type Repository interface {
	CreateProject(p *PNRRProject) error
	ListProjectsBySchool(schoolID, academicYear string) ([]PNRRProject, error)
	GetProject(id string) (*PNRRProject, error)
	EnrollStudent(ps *PNRRProjectStudent) error
	ListProjectStudents(projectID string) ([]PNRRProjectStudent, error)
	RecordLesson(l *PNRRLessonAttendance, attendedStudentIDs []string) error
	ListLessons(projectID string) ([]PNRRLessonAttendance, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateProject(p *PNRRProject) error {
	query := `
		INSERT INTO pnrr_projects (
			school_id, cup_code, project_title, notice_type, total_hours,
			academic_year, tutor_id, expert_id, target_students_count, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query, p.SchoolID, p.CUPCode, p.ProjectTitle, p.NoticeType, p.TotalHours,
		p.AcademicYear, p.TutorID, p.ExpertID, p.TargetStudentsCount, p.Status,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *repository) ListProjectsBySchool(schoolID, academicYear string) ([]PNRRProject, error) {
	query := `
		SELECT p.id, p.school_id, p.cup_code, p.project_title, p.notice_type, p.total_hours,
		       p.academic_year, p.tutor_id, p.expert_id, p.target_students_count, p.status,
		       p.created_at, p.updated_at,
		       COALESCE(ut.first_name || ' ' || ut.last_name, '') AS tutor_name,
		       COALESCE(ue.first_name || ' ' || ue.last_name, '') AS expert_name,
		       COUNT(ps.id) AS enrolled_students,
		       COALESCE(SUM(CASE WHEN ps.target_75_reached THEN 1 ELSE 0 END), 0) AS target_75_passed
		FROM pnrr_projects p
		LEFT JOIN users ut ON p.tutor_id = ut.id
		LEFT JOIN users ue ON p.expert_id = ue.id
		LEFT JOIN pnrr_project_students ps ON p.id = ps.project_id
		WHERE p.school_id = $1 AND ($2 = '' OR p.academic_year = $2)
		GROUP BY p.id, ut.first_name, ut.last_name, ue.first_name, ue.last_name
		ORDER BY p.created_at DESC
	`
	rows, err := r.db.Query(query, schoolID, academicYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []PNRRProject
	for rows.Next() {
		var p PNRRProject
		var passed int
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.CUPCode, &p.ProjectTitle, &p.NoticeType, &p.TotalHours,
			&p.AcademicYear, &p.TutorID, &p.ExpertID, &p.TargetStudentsCount, &p.Status,
			&p.CreatedAt, &p.UpdatedAt, &p.TutorName, &p.ExpertName,
			&p.EnrolledStudents, &passed,
		); err != nil {
			return nil, err
		}
		p.Target75PassedCount = passed
		if p.EnrolledStudents > 0 {
			p.OverallTargetRate = float64(p.Target75PassedCount) / float64(p.EnrolledStudents) * 100.0
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) GetProject(id string) (*PNRRProject, error) {
	query := `
		SELECT p.id, p.school_id, p.cup_code, p.project_title, p.notice_type, p.total_hours,
		       p.academic_year, p.tutor_id, p.expert_id, p.target_students_count, p.status,
		       p.created_at, p.updated_at,
		       COALESCE(ut.first_name || ' ' || ut.last_name, '') AS tutor_name,
		       COALESCE(ue.first_name || ' ' || ue.last_name, '') AS expert_name,
		       COUNT(ps.id) AS enrolled_students,
		       COALESCE(SUM(CASE WHEN ps.target_75_reached THEN 1 ELSE 0 END), 0) AS target_75_passed
		FROM pnrr_projects p
		LEFT JOIN users ut ON p.tutor_id = ut.id
		LEFT JOIN users ue ON p.expert_id = ue.id
		LEFT JOIN pnrr_project_students ps ON p.id = ps.project_id
		WHERE p.id = $1
		GROUP BY p.id, ut.first_name, ut.last_name, ue.first_name, ue.last_name
	`
	var p PNRRProject
	var passed int
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.SchoolID, &p.CUPCode, &p.ProjectTitle, &p.NoticeType, &p.TotalHours,
		&p.AcademicYear, &p.TutorID, &p.ExpertID, &p.TargetStudentsCount, &p.Status,
		&p.CreatedAt, &p.UpdatedAt, &p.TutorName, &p.ExpertName,
		&p.EnrolledStudents, &passed,
	)
	if err != nil {
		return nil, err
	}
	p.Target75PassedCount = passed
	if p.EnrolledStudents > 0 {
		p.OverallTargetRate = float64(p.Target75PassedCount) / float64(p.EnrolledStudents) * 100.0
	}
	return &p, nil
}

func (r *repository) EnrollStudent(ps *PNRRProjectStudent) error {
	query := `
		INSERT INTO pnrr_project_students (project_id, student_id, enrolled_at)
		VALUES ($1, $2, CURRENT_DATE)
		ON CONFLICT (project_id, student_id) DO NOTHING
		RETURNING id, enrolled_at, created_at, updated_at
	`
	return r.db.QueryRow(query, ps.ProjectID, ps.StudentID).Scan(
		&ps.ID, &ps.EnrolledAt, &ps.CreatedAt, &ps.UpdatedAt,
	)
}

func (r *repository) ListProjectStudents(projectID string) ([]PNRRProjectStudent, error) {
	query := `
		SELECT ps.id, ps.project_id, ps.student_id, TO_CHAR(ps.enrolled_at, 'YYYY-MM-DD'),
		       ps.attended_hours, ps.attendance_percentage, ps.target_75_reached,
		       ps.created_at, ps.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name
		FROM pnrr_project_students ps
		LEFT JOIN users u ON ps.student_id = u.id
		WHERE ps.project_id = $1
		ORDER BY student_name ASC
	`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []PNRRProjectStudent
	for rows.Next() {
		var ps PNRRProjectStudent
		if err := rows.Scan(
			&ps.ID, &ps.ProjectID, &ps.StudentID, &ps.EnrolledAt, &ps.AttendedHours,
			&ps.AttendancePercentage, &ps.Target75Reached, &ps.CreatedAt, &ps.UpdatedAt,
			&ps.StudentFullName,
		); err != nil {
			return nil, err
		}
		list = append(list, ps)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) RecordLesson(l *PNRRLessonAttendance, attendedStudentIDs []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert lesson
	query := `
		INSERT INTO pnrr_lesson_attendances (project_id, lesson_date, hours, topic, signatures_json)
		VALUES ($1, $2, $3, $4, $5::jsonb)
		RETURNING id, created_at
	`
	if l.SignaturesJSON == "" {
		l.SignaturesJSON = "[]"
	}
	if err := tx.QueryRow(query, l.ProjectID, l.LessonDate, l.Hours, l.Topic, l.SignaturesJSON).Scan(&l.ID, &l.CreatedAt); err != nil {
		return err
	}

	// Fetch project total hours for 75% calculation
	var totalHours int
	if err := tx.QueryRow(`SELECT total_hours FROM pnrr_projects WHERE id = $1`, l.ProjectID).Scan(&totalHours); err != nil {
		return err
	}

	// Update attended hours for present students
	for _, studentID := range attendedStudentIDs {
		updateQuery := `
			UPDATE pnrr_project_students
			SET attended_hours = attended_hours + $2,
			    attendance_percentage = ROUND(((attended_hours + $2) / NULLIF($3, 0) * 100)::numeric, 2),
			    target_75_reached = (((attended_hours + $2) / NULLIF($3, 0) * 100) >= 75.0),
			    updated_at = NOW()
			WHERE project_id = $1 AND student_id = $4
		`
		if _, err := tx.Exec(updateQuery, l.ProjectID, l.Hours, totalHours, studentID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *repository) ListLessons(projectID string) ([]PNRRLessonAttendance, error) {
	query := `
		SELECT id, project_id, TO_CHAR(lesson_date, 'YYYY-MM-DD'), hours, topic,
		       signatures_json::text, created_at
		FROM pnrr_lesson_attendances
		WHERE project_id = $1
		ORDER BY lesson_date DESC
	`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []PNRRLessonAttendance
	for rows.Next() {
		var l PNRRLessonAttendance
		if err := rows.Scan(
			&l.ID, &l.ProjectID, &l.LessonDate, &l.Hours, &l.Topic, &l.SignaturesJSON, &l.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
