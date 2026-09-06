package teachers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	GetPersonalRegisterData(ctx context.Context, teacherID, classID, subjectID string) (*TeacherRegisterData, error)
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
	defer func() { _ = rows.Close() }()

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
	defer func() { _ = rows.Close() }()

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
	defer func() { _ = rows.Close() }()

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

func (r *PostgresRepository) GetPersonalRegisterData(ctx context.Context, teacherID, classID, subjectID string) (*TeacherRegisterData, error) {
	data := &TeacherRegisterData{
		AcademicYear: "2023/2024",
		Students:     make([]TeacherRegisterStudent, 0),
		Lessons:      make([]TeacherRegisterLesson, 0),
	}

	// 1. Teacher & School details
	teacherQuery := `
		SELECT COALESCE(u.first_name || ' ' || u.last_name, 'Docente'), COALESCE(s.name, 'Istituto Comprensivo Statale')
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		LEFT JOIN schools s ON t.school_id = s.id
		WHERE t.id = NULLIF($1, '')::uuid OR t.user_id = NULLIF($1, '')::uuid
		LIMIT 1`
	_ = r.db.QueryRowContext(ctx, teacherQuery, teacherID).Scan(&data.TeacherName, &data.SchoolName)
	if data.TeacherName == "" {
		data.TeacherName = "Docente"
	}
	if data.SchoolName == "" {
		data.SchoolName = "Istituto Scolastico"
	}

	// 2. Class & Subject determination
	if classID == "" || subjectID == "" {
		assignQuery := `
			SELECT cs.class_id::text, cs.subject_id::text,
			       COALESCE(al.name || ' ', '') || c.section, COALESCE(sub.name, 'Materia')
			FROM class_subjects cs
			JOIN classes c ON cs.class_id = c.id
			JOIN subjects sub ON cs.subject_id = sub.id
			LEFT JOIN academic_levels al ON c.level_id = al.id
			WHERE cs.teacher_id = NULLIF($1, '')::uuid OR cs.teacher_id IN (SELECT id FROM teachers WHERE user_id = NULLIF($1, '')::uuid)
			LIMIT 1`
		_ = r.db.QueryRowContext(ctx, assignQuery, teacherID).Scan(&classID, &subjectID, &data.ClassName, &data.SubjectName)
	} else {
		classQuery := `SELECT COALESCE(al.name || ' ', '') || c.section FROM classes c LEFT JOIN academic_levels al ON c.level_id = al.id WHERE c.id = NULLIF($1, '')::uuid`
		_ = r.db.QueryRowContext(ctx, classQuery, classID).Scan(&data.ClassName)
		subQuery := `SELECT name FROM subjects WHERE id = NULLIF($1, '')::uuid`
		_ = r.db.QueryRowContext(ctx, subQuery, subjectID).Scan(&data.SubjectName)
	}

	if data.ClassName == "" {
		data.ClassName = "Classe Non Specificata"
	}
	if data.SubjectName == "" {
		data.SubjectName = "Materia Generale"
	}

	// 3. Students list
	studentMap := make(map[string]*TeacherRegisterStudent)
	studentOrder := make([]string, 0)
	if classID != "" {
		stQuery := `
			SELECT s.id::text, COALESCE(u.last_name || ' ' || u.first_name, 'Studente')
			FROM students s
			JOIN users u ON s.user_id = u.id
			WHERE s.class_id = NULLIF($1, '')::uuid AND s.deleted_at IS NULL AND u.deleted_at IS NULL
			ORDER BY u.last_name, u.first_name`
		if rows, err := r.db.QueryContext(ctx, stQuery, classID); err == nil {
			defer rows.Close()
			for rows.Next() {
				var sid, sname string
				if err := rows.Scan(&sid, &sname); err == nil {
					st := &TeacherRegisterStudent{
						ID:          sid,
						Name:        sname,
						Q1Written:   "-",
						Q1Oral:      "-",
						Q1Practical: "-",
						Q1Avg:       "-",
						Q2Written:   "-",
						Q2Oral:      "-",
						Q2Practical: "-",
						Q2Avg:       "-",
						FinalAvg:    "-",
						Absences:    0,
					}
					studentMap[sid] = st
					studentOrder = append(studentOrder, sid)
				}
			}
		}

		// 4. Grades aggregation
		gradesQuery := `
			SELECT student_id::text, grade_value, grade_type, semester
			FROM grades
			WHERE class_id = NULLIF($1, '')::uuid 
			  AND (subject_id = NULLIF($2, '')::uuid OR $2 = '')
			  AND deleted_at IS NULL`
		if rows, err := r.db.QueryContext(ctx, gradesQuery, classID, subjectID); err == nil {
			defer rows.Close()
			type studentGrades struct {
				q1Written, q1Oral, q1Prac []float64
				q2Written, q2Oral, q2Prac []float64
			}
			stGradesMap := make(map[string]*studentGrades)

			for rows.Next() {
				var sid string
				var val float64
				var gType string
				var sem int
				if err := rows.Scan(&sid, &val, &gType, &sem); err == nil {
					if _, ok := stGradesMap[sid]; !ok {
						stGradesMap[sid] = &studentGrades{}
					}
					sg := stGradesMap[sid]
					if sem == 1 {
						switch gType {
						case "Written":
							sg.q1Written = append(sg.q1Written, val)
						case "Oral":
							sg.q1Oral = append(sg.q1Oral, val)
						case "Practical":
							sg.q1Prac = append(sg.q1Prac, val)
						default:
							sg.q1Oral = append(sg.q1Oral, val)
						}
					} else {
						switch gType {
						case "Written":
							sg.q2Written = append(sg.q2Written, val)
						case "Oral":
							sg.q2Oral = append(sg.q2Oral, val)
						case "Practical":
							sg.q2Prac = append(sg.q2Prac, val)
						default:
							sg.q2Oral = append(sg.q2Oral, val)
						}
					}
				}
			}

			formatAvg := func(vals []float64) (string, float64, bool) {
				if len(vals) == 0 {
					return "-", 0, false
				}
				sum := 0.0
				for _, v := range vals {
					sum += v
				}
				avg := sum / float64(len(vals))
				return fmt.Sprintf("%.1f", avg), avg, true
			}

			for sid, sg := range stGradesMap {
				if st, ok := studentMap[sid]; ok {
					st.Q1Written, _, _ = formatAvg(sg.q1Written)
					st.Q1Oral, _, _ = formatAvg(sg.q1Oral)
					st.Q1Practical, _, _ = formatAvg(sg.q1Prac)
					allQ1 := append(append(sg.q1Written, sg.q1Oral...), sg.q1Prac...)
					_, q1AvgVal, hasQ1 := formatAvg(allQ1)
					if hasQ1 {
						st.Q1Avg = fmt.Sprintf("%.2f", q1AvgVal)
					}

					st.Q2Written, _, _ = formatAvg(sg.q2Written)
					st.Q2Oral, _, _ = formatAvg(sg.q2Oral)
					st.Q2Practical, _, _ = formatAvg(sg.q2Prac)
					allQ2 := append(append(sg.q2Written, sg.q2Oral...), sg.q2Prac...)
					_, q2AvgVal, hasQ2 := formatAvg(allQ2)
					if hasQ2 {
						st.Q2Avg = fmt.Sprintf("%.2f", q2AvgVal)
					}

					allGrades := append(allQ1, allQ2...)
					_, finalAvgVal, hasFinal := formatAvg(allGrades)
					if hasFinal {
						st.FinalAvg = fmt.Sprintf("%.2f", finalAvgVal)
					}
				}
			}
		}

		// 5. Subject Absences count
		absQuery := `
			SELECT student_id::text, COUNT(*)
			FROM attendance
			WHERE class_id = NULLIF($1, '')::uuid
			  AND (subject_id = NULLIF($2, '')::uuid OR $2 = '' OR subject_id IS NULL)
			  AND status = 'Absent'
			GROUP BY student_id`
		if rows, err := r.db.QueryContext(ctx, absQuery, classID, subjectID); err == nil {
			defer rows.Close()
			for rows.Next() {
				var sid string
				var count int
				if err := rows.Scan(&sid, &count); err == nil {
					if st, ok := studentMap[sid]; ok {
						st.Absences = count
					}
				}
			}
		}

		// 6. Signed lessons
		lessonsQuery := `
			SELECT TO_CHAR(l.date, 'DD/MM/YYYY'), COALESCE(l.hour, 1), l.topic, l.type,
			       COALESCE(u.first_name || ' ' || u.last_name, '')
			FROM class_lessons l
			JOIN users u ON l.teacher_id = u.id
			WHERE l.class_id = NULLIF($1, '')::uuid
			  AND (l.subject_id = NULLIF($2, '')::uuid OR $2 = '')
			ORDER BY l.date, l.hour LIMIT 80`
		if rows, err := r.db.QueryContext(ctx, lessonsQuery, classID, subjectID); err == nil {
			defer rows.Close()
			for rows.Next() {
				var l TeacherRegisterLesson
				if err := rows.Scan(&l.Date, &l.Hour, &l.Topic, &l.Type, &l.SignedBy); err == nil {
					data.Lessons = append(data.Lessons, l)
				}
			}
		}
	}

	for _, sid := range studentOrder {
		data.Students = append(data.Students, *studentMap[sid])
	}

	return data, nil
}
