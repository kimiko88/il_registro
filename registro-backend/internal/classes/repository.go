package classes

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, class *Class) error
	List(ctx context.Context, schoolID string, academicYear string) ([]Class, error)
	Get(ctx context.Context, id string) (*Class, error)
	Update(ctx context.Context, class *Class) error
	Delete(ctx context.Context, id string) error
	ListByTeacher(ctx context.Context, teacherUserID string, schoolYear string) ([]Class, error)
	AssignSubject(ctx context.Context, classID string, subjectID string, teacherID *string, hours float64) error
	UnassignSubject(ctx context.Context, assignmentID string) error
	GetClassSubjects(ctx context.Context, classID string) ([]ClassSubject, error)
	GetClassGuardians(ctx context.Context, classID string) ([]GuardianInfo, error)
	GetLessonTopics(ctx context.Context, classID string) ([]LessonTopic, error)
	GetDisciplinaryNotes(ctx context.Context, classID string) ([]DisciplinaryNoteReport, error)
	BulkMigrateStudents(ctx context.Context, migrations []StudentMigrationItem) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// ... existing ListByTeacher ...

func (r *PostgresRepository) AssignSubject(ctx context.Context, classID string, subjectID string, teacherID *string, hours float64) error {
	// Verify teacher exists if provided
	if teacherID != nil {
		query := `SELECT id FROM teachers WHERE id = $1`
		var tid string
		err := r.db.QueryRowContext(ctx, query, *teacherID).Scan(&tid)
		if err != nil {
			return errors.New("teacher profile not found")
		}
	}

	// Insert
	id := uuid.New().String()
	// Constraint is (class_id, subject_id, teacher_id).
	query := `INSERT INTO class_subjects (id, class_id, subject_id, teacher_id, hours_per_week)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, id, classID, subjectID, teacherID, hours)
	return err
}

func (r *PostgresRepository) UnassignSubject(ctx context.Context, assignmentID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM class_subjects WHERE id = $1", assignmentID)
	return err
}

func (r *PostgresRepository) GetClassSubjects(ctx context.Context, classID string) ([]ClassSubject, error) {
	if classID == "" {
		return []ClassSubject{}, nil
	}
	query := `
		SELECT cs.id, cs.class_id, cs.subject_id, s.name, COALESCE(cs.teacher_id::text, ''), COALESCE(u.last_name, ''), COALESCE(u.first_name, ''), COALESCE(cs.hours_per_week, 0)
		FROM class_subjects cs
		JOIN subjects s ON cs.subject_id::text = s.id::text
		LEFT JOIN teachers t ON (NULLIF(cs.teacher_id::text, '') = t.id::text OR NULLIF(cs.teacher_id::text, '') = t.user_id::text)
		LEFT JOIN users u ON t.user_id::text = u.id::text
		WHERE cs.class_id::text = $1
		ORDER BY s.name
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ClassSubject
	for rows.Next() {
		var cs ClassSubject
		var tLen, tFirst sql.NullString
		err := rows.Scan(
			&cs.ID, &cs.ClassID, &cs.SubjectID, &cs.SubjectName, &cs.TeacherID, &tLen, &tFirst, &cs.HoursPerWeek,
		)
		if err != nil {
			return nil, err
		}
		if tLen.Valid {
			cs.TeacherName = tLen.String + " " + tFirst.String
		}
		results = append(results, cs)
	}
	return results, nil
}

// ... existing methods ...

func (r *PostgresRepository) ListByTeacher(ctx context.Context, teacherUserID string, schoolYear string) ([]Class, error) {
	if teacherUserID == "" {
		return []Class{}, nil
	}
	if _, err := uuid.Parse(teacherUserID); err != nil {
		return []Class{}, nil
	}
	// Combine classes where user is coordinator OR assigned as teacher (via class_subjects)
	var query string
	var args []interface{}
	args = append(args, teacherUserID)

	if schoolYear != "" {
		query = `
			SELECT DISTINCT c.id, c.school_id, c.name, c.section, c.articolazione, c.academic_year, c.coordinator_id,
			       (SELECT COUNT(*) FROM students s WHERE s.class_id = c.id) AS students_count,
			       c.created_at, c.updated_at
			FROM classes c
			LEFT JOIN class_subjects cs ON c.id::text = cs.class_id::text
			LEFT JOIN teachers t ON (NULLIF(cs.teacher_id::text, '') = t.id::text OR NULLIF(cs.teacher_id::text, '') = t.user_id::text)
			WHERE (t.user_id::text = $1 OR cs.teacher_id::text = $1 OR c.coordinator_id::text = $1)
			  AND (c.academic_year IS NULL OR c.academic_year = '' OR REPLACE(c.academic_year, '-', '/') = REPLACE($2, '-', '/'))
			ORDER BY c.name
		`
		args = append(args, schoolYear)
	} else {
		query = `
			SELECT DISTINCT c.id, c.school_id, c.name, c.section, c.articolazione, c.academic_year, c.coordinator_id,
			       (SELECT COUNT(*) FROM students s WHERE s.class_id = c.id) AS students_count,
			       c.created_at, c.updated_at
			FROM classes c
			LEFT JOIN class_subjects cs ON c.id::text = cs.class_id::text
			LEFT JOIN teachers t ON (NULLIF(cs.teacher_id::text, '') = t.id::text OR NULLIF(cs.teacher_id::text, '') = t.user_id::text)
			WHERE (t.user_id::text = $1 OR cs.teacher_id::text = $1 OR c.coordinator_id::text = $1)
			ORDER BY c.name
		`
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var c Class
		var sec, art, coord sql.NullString
		if err := rows.Scan(&c.ID, &c.SchoolID, &c.Name, &sec, &art, &c.AcademicYear, &coord, &c.StudentsCount, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.Section = sec.String
		c.Articolazione = art.String
		c.CoordinatorID = coord.String
		classes = append(classes, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if classes == nil {
		classes = []Class{}
	}
	return classes, nil
}

func (r *PostgresRepository) Create(ctx context.Context, c *Class) error {
	c.ID = uuid.New().String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	query := `INSERT INTO classes (id, school_id, name, section, articolazione, location, academic_year, coordinator_id, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.SchoolID, c.Name, c.Section, c.Articolazione, c.Location, c.AcademicYear,
		sql.NullString{String: c.CoordinatorID, Valid: c.CoordinatorID != ""},
		c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *PostgresRepository) List(ctx context.Context, schoolID string, academicYear string) ([]Class, error) {
	if schoolID != "" {
		if _, err := uuid.Parse(schoolID); err != nil {
			return []Class{}, nil
		}
	}
	query := `SELECT c.id, c.school_id, c.name, c.section, c.articolazione, COALESCE(c.location, ''), c.academic_year, c.coordinator_id,
	                 (SELECT COUNT(*) FROM students s WHERE s.class_id = c.id) AS students_count,
	                 c.created_at, c.updated_at 
	          FROM classes c WHERE ($1 = '' OR c.school_id = NULLIF($1, '')::uuid)`

	args := []interface{}{schoolID}
	if academicYear != "" {
		query += " AND c.academic_year = $2"
		args = append(args, academicYear)
	}
	query += " ORDER BY c.name"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var c Class
		var sec, art, loc, coord sql.NullString
		if err := rows.Scan(&c.ID, &c.SchoolID, &c.Name, &sec, &art, &loc, &c.AcademicYear, &coord, &c.StudentsCount, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.Section = sec.String
		c.Articolazione = art.String
		c.Location = loc.String
		c.CoordinatorID = coord.String
		classes = append(classes, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if classes == nil {
		classes = []Class{}
	}
	return classes, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*Class, error) {
	query := `SELECT c.id, c.school_id, c.name, c.section, c.articolazione, COALESCE(c.location, ''), c.academic_year, c.coordinator_id,
	                 (SELECT COUNT(*) FROM students s WHERE s.class_id = c.id) AS students_count,
	                 c.created_at, c.updated_at 
	          FROM classes c WHERE c.id = $1`
	var c Class
	var sec, art, loc, coord sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.SchoolID, &c.Name, &sec, &art, &loc, &c.AcademicYear, &coord, &c.StudentsCount, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	c.Section = sec.String
	c.Articolazione = art.String
	c.Location = loc.String
	c.CoordinatorID = coord.String
	return &c, nil
}

func (r *PostgresRepository) Update(ctx context.Context, c *Class) error {
	c.UpdatedAt = time.Now()
	query := `UPDATE classes 
			  SET name = $1, section = $2, articolazione = $3, location = $4, academic_year = $5, coordinator_id = $6, updated_at = $7
			  WHERE id = $8`
	res, err := r.db.ExecContext(ctx, query,
		c.Name, c.Section, c.Articolazione, c.Location, c.AcademicYear,
		sql.NullString{String: c.CoordinatorID, Valid: c.CoordinatorID != ""},
		c.UpdatedAt, c.ID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("class not found")
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM classes WHERE id = $1", id)
	return err
}

func (r *PostgresRepository) GetClassGuardians(ctx context.Context, classID string) ([]GuardianInfo, error) {
	if _, err := uuid.Parse(classID); err != nil {
		return []GuardianInfo{}, nil
	}
	query := `
		SELECT DISTINCT u.id, u.first_name, u.last_name, u.email, COALESCE(u.phone_number, ''),
		       su.id, su.first_name || ' ' || su.last_name AS student_name, su.last_name AS s_last, su.first_name AS s_first
		FROM students st
		JOIN users su ON su.id = st.user_id
		JOIN student_parents sp ON (sp.student_id = st.id OR sp.student_id = su.id)
		JOIN parents p ON (p.id = sp.parent_id OR p.user_id = sp.parent_id)
		JOIN users u ON u.id = p.user_id
		WHERE st.class_id = $1::uuid
		ORDER BY su.last_name, su.first_name, u.last_name
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []GuardianInfo
	for rows.Next() {
		var g GuardianInfo
		var sLast, sFirst string
		if err := rows.Scan(&g.GuardianID, &g.FirstName, &g.LastName, &g.Email, &g.Phone, &g.StudentID, &g.StudentName, &sLast, &sFirst); err != nil {
			return nil, err
		}
		result = append(result, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PostgresRepository) GetLessonTopics(ctx context.Context, classID string) ([]LessonTopic, error) {
	if _, err := uuid.Parse(classID); err != nil {
		return []LessonTopic{}, nil
	}
	query := `
		SELECT cl.id::text, cl.date, COALESCE(s.name, ''), COALESCE(cl.topic, ''), COALESCE(u.first_name, ''), COALESCE(u.last_name, '')
		FROM class_lessons cl
		LEFT JOIN subjects s ON cl.subject_id = s.id
		LEFT JOIN users u ON cl.teacher_id = u.id
		WHERE cl.class_id = $1::uuid
		ORDER BY cl.date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []LessonTopic
	for rows.Next() {
		var lt LessonTopic
		if err := rows.Scan(&lt.ID, &lt.Date, &lt.SubjectName, &lt.Topic, &lt.TeacherFirstName, &lt.TeacherLastName); err != nil {
			return nil, err
		}
		topics = append(topics, lt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if topics == nil {
		topics = []LessonTopic{}
	}
	return topics, nil
}

func (r *PostgresRepository) GetDisciplinaryNotes(ctx context.Context, classID string) ([]DisciplinaryNoteReport, error) {
	if _, err := uuid.Parse(classID); err != nil {
		return []DisciplinaryNoteReport{}, nil
	}
	query := `
		SELECT sn.id::text, sn.date, COALESCE(su.first_name, ''), COALESCE(su.last_name, ''), sn.type::text, COALESCE(sn.note, ''), COALESCE(tu.first_name, ''), COALESCE(tu.last_name, '')
		FROM student_notes sn
		LEFT JOIN users su ON sn.student_id = su.id
		LEFT JOIN users tu ON sn.teacher_id = tu.id
		WHERE sn.class_id = $1::uuid
		ORDER BY sn.date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []DisciplinaryNoteReport
	for rows.Next() {
		var dnr DisciplinaryNoteReport
		if err := rows.Scan(&dnr.ID, &dnr.Date, &dnr.StudentFirstName, &dnr.StudentLastName, &dnr.NoteType, &dnr.Description, &dnr.TeacherFirstName, &dnr.TeacherLastName); err != nil {
			return nil, err
		}
		notes = append(notes, dnr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if notes == nil {
		notes = []DisciplinaryNoteReport{}
	}
	return notes, nil
}

func (r *PostgresRepository) BulkMigrateStudents(ctx context.Context, migrations []StudentMigrationItem) error {
	if len(migrations) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmtUpdateClass, err := tx.PrepareContext(ctx, `
		UPDATE students
		SET class_id = $1::uuid, updated_at = NOW()
		WHERE id = $2::uuid OR user_id = $2::uuid
	`)
	if err != nil {
		return err
	}
	defer stmtUpdateClass.Close()

	stmtUnassignClass, err := tx.PrepareContext(ctx, `
		UPDATE students
		SET class_id = NULL, updated_at = NOW()
		WHERE id = $1::uuid OR user_id = $1::uuid
	`)
	if err != nil {
		return err
	}
	defer stmtUnassignClass.Close()

	for _, item := range migrations {
		switch item.Action {
		case "promoted", "repeater":
			if item.TargetClassID != "" {
				if _, err := stmtUpdateClass.ExecContext(ctx, item.TargetClassID, item.StudentID); err != nil {
					return err
				}
			}
		case "graduated", "left":
			if _, err := stmtUnassignClass.ExecContext(ctx, item.StudentID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
