package groups

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository interface {
	Create(ctx context.Context, g *Group) error
	Update(ctx context.Context, g *Group) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Group, error)
	ListBySchool(ctx context.Context, schoolID string) ([]Group, error)
	ListByTeacher(ctx context.Context, teacherID string) ([]Group, error)
	ListByStudent(ctx context.Context, studentID string) ([]Group, error)
	AddStudents(ctx context.Context, groupID string, studentIDs []string) error
	RemoveStudent(ctx context.Context, groupID, studentID string) error
	GetStudentsInGroup(ctx context.Context, groupID string) ([]GroupStudentInfo, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, g *Group) error {
	query := `
		INSERT INTO groups (id, school_id, name, subject_id, teacher_id, academic_year, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	if g.ID == "" {
		g.ID = fmt.Sprintf("group-%d", time.Now().UnixNano())
	}
	now := time.Now()
	g.CreatedAt = now
	g.UpdatedAt = now

	var subjectID, teacherID sql.NullString
	if g.SubjectID != nil && *g.SubjectID != "" {
		subjectID = sql.NullString{String: *g.SubjectID, Valid: true}
	}
	if g.TeacherID != nil && *g.TeacherID != "" {
		teacherID = sql.NullString{String: *g.TeacherID, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query, g.ID, g.SchoolID, g.Name, subjectID, teacherID, g.AcademicYear, g.Description, g.CreatedAt, g.UpdatedAt)
	return err
}

func (r *repository) Update(ctx context.Context, g *Group) error {
	query := `
		UPDATE groups
		SET name = $1, subject_id = $2, teacher_id = $3, description = $4, updated_at = $5
		WHERE id = $6
	`
	var subjectID, teacherID sql.NullString
	if g.SubjectID != nil && *g.SubjectID != "" {
		subjectID = sql.NullString{String: *g.SubjectID, Valid: true}
	}
	if g.TeacherID != nil && *g.TeacherID != "" {
		teacherID = sql.NullString{String: *g.TeacherID, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query, g.Name, subjectID, teacherID, g.Description, time.Now(), g.ID)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM groups WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) GetByID(ctx context.Context, id string) (*Group, error) {
	query := `
		SELECT g.id, g.school_id, g.name, g.subject_id, COALESCE(s.name, ''), g.teacher_id,
		       COALESCE(u.first_name || ' ' || u.last_name, ''), g.academic_year, COALESCE(g.description, ''),
		       g.created_at, g.updated_at
		FROM groups g
		LEFT JOIN subjects s ON s.id = g.subject_id
		LEFT JOIN users u ON u.id = g.teacher_id
		WHERE g.id = $1
	`
	var g Group
	var subjectID, teacherID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&g.ID, &g.SchoolID, &g.Name, &subjectID, &g.SubjectName,
		&teacherID, &g.TeacherName, &g.AcademicYear, &g.Description,
		&g.CreatedAt, &g.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if subjectID.Valid {
		g.SubjectID = &subjectID.String
	}
	if teacherID.Valid {
		g.TeacherID = &teacherID.String
	}

	students, err := r.GetStudentsInGroup(ctx, id)
	if err == nil {
		g.Students = students
		g.StudentCount = len(students)
	}

	return &g, nil
}

func (r *repository) ListBySchool(ctx context.Context, schoolID string) ([]Group, error) {
	query := `
		SELECT g.id, g.school_id, g.name, g.subject_id, COALESCE(s.name, ''), g.teacher_id,
		       COALESCE(u.first_name || ' ' || u.last_name, ''), g.academic_year, COALESCE(g.description, ''),
		       (SELECT COUNT(*) FROM group_students gs WHERE gs.group_id = g.id) AS student_count,
		       g.created_at, g.updated_at
		FROM groups g
		LEFT JOIN subjects s ON s.id = g.subject_id
		LEFT JOIN users u ON u.id = g.teacher_id
		WHERE g.school_id = $1
		ORDER BY g.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Group
	for rows.Next() {
		var g Group
		var subjectID, teacherID sql.NullString
		if err := rows.Scan(
			&g.ID, &g.SchoolID, &g.Name, &subjectID, &g.SubjectName,
			&teacherID, &g.TeacherName, &g.AcademicYear, &g.Description,
			&g.StudentCount, &g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if subjectID.Valid {
			g.SubjectID = &subjectID.String
		}
		if teacherID.Valid {
			g.TeacherID = &teacherID.String
		}
		result = append(result, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) ListByTeacher(ctx context.Context, teacherID string) ([]Group, error) {
	query := `
		SELECT g.id, g.school_id, g.name, g.subject_id, COALESCE(s.name, ''), g.teacher_id,
		       COALESCE(u.first_name || ' ' || u.last_name, ''), g.academic_year, COALESCE(g.description, ''),
		       (SELECT COUNT(*) FROM group_students gs WHERE gs.group_id = g.id) AS student_count,
		       g.created_at, g.updated_at
		FROM groups g
		LEFT JOIN subjects s ON s.id = g.subject_id
		LEFT JOIN users u ON u.id = g.teacher_id
		WHERE g.teacher_id = $1
		ORDER BY g.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Group
	for rows.Next() {
		var g Group
		var subjectID, teacherID sql.NullString
		if err := rows.Scan(
			&g.ID, &g.SchoolID, &g.Name, &subjectID, &g.SubjectName,
			&teacherID, &g.TeacherName, &g.AcademicYear, &g.Description,
			&g.StudentCount, &g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if subjectID.Valid {
			g.SubjectID = &subjectID.String
		}
		if teacherID.Valid {
			g.TeacherID = &teacherID.String
		}
		result = append(result, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) ListByStudent(ctx context.Context, studentID string) ([]Group, error) {
	query := `
		SELECT g.id, g.school_id, g.name, g.subject_id, COALESCE(s.name, ''), g.teacher_id,
		       COALESCE(u.first_name || ' ' || u.last_name, ''), g.academic_year, COALESCE(g.description, ''),
		       (SELECT COUNT(*) FROM group_students gs WHERE gs.group_id = g.id) AS student_count,
		       g.created_at, g.updated_at
		FROM groups g
		JOIN group_students gs ON gs.group_id = g.id
		LEFT JOIN subjects s ON s.id = g.subject_id
		LEFT JOIN users u ON u.id = g.teacher_id
		WHERE gs.student_id = $1
		ORDER BY g.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Group
	for rows.Next() {
		var g Group
		var subjectID, teacherID sql.NullString
		if err := rows.Scan(
			&g.ID, &g.SchoolID, &g.Name, &subjectID, &g.SubjectName,
			&teacherID, &g.TeacherName, &g.AcademicYear, &g.Description,
			&g.StudentCount, &g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if subjectID.Valid {
			g.SubjectID = &subjectID.String
		}
		if teacherID.Valid {
			g.TeacherID = &teacherID.String
		}
		result = append(result, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) AddStudents(ctx context.Context, groupID string, studentIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO group_students (group_id, student_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	for _, sid := range studentIDs {
		if sid == "" {
			continue
		}
		if _, err := stmt.ExecContext(ctx, groupID, sid); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *repository) RemoveStudent(ctx context.Context, groupID, studentID string) error {
	query := `DELETE FROM group_students WHERE group_id = $1 AND student_id = $2`
	_, err := r.db.ExecContext(ctx, query, groupID, studentID)
	return err
}

func (r *repository) GetStudentsInGroup(ctx context.Context, groupID string) ([]GroupStudentInfo, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, COALESCE(c.name, '')
		FROM group_students gs
		JOIN users u ON u.id = gs.student_id
		LEFT JOIN classes c ON c.id = u.class_id
		WHERE gs.group_id = $1
		ORDER BY u.last_name, u.first_name
	`
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []GroupStudentInfo
	for rows.Next() {
		var s GroupStudentInfo
		if err := rows.Scan(&s.StudentID, &s.FirstName, &s.LastName, &s.ClassName); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}
