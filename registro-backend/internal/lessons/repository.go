package lessons

import (
	"database/sql"
)

type Repository interface {
	CreateLesson(lesson *Lesson) error
	GetLessonByID(id string) (*Lesson, error)
	UpdateLesson(id string, req UpdateLessonRequest) (*Lesson, error)
	DeleteLesson(id string) error
	GetLessonsByClass(classID string, date string) ([]Lesson, error)
	GetLessonsByClassAndSubject(classID, subjectID string, date string) ([]Lesson, error)
	GetLessonsByGroup(groupID string, date string) ([]Lesson, error)
	GetLessonsByTeacher(teacherID string, fromDate, toDate string) ([]Lesson, error)

	CreateHomework(homework *Homework) error
	GetHomeworkByID(id string) (*Homework, error)
	UpdateHomework(id string, req UpdateHomeworkRequest) (*Homework, error)
	DeleteHomework(id string) error
	GetHomeworkByClass(classID string) ([]Homework, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateLesson(lesson *Lesson) error {
	query := `
		WITH inserted AS (
			INSERT INTO class_lessons (class_id, teacher_id, subject_id, date, hour, duration, topic, type, group_id, is_substitution, substituted_teacher_id, activity_type, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
			RETURNING id, teacher_id, substituted_teacher_id
		)
		SELECT inserted.id,
		       COALESCE(u1.first_name || ' ' || u1.last_name, '') AS teacher_name,
		       COALESCE(u2.first_name || ' ' || u2.last_name, '') AS substituted_teacher_name
		FROM inserted
		LEFT JOIN users u1 ON inserted.teacher_id = u1.id
		LEFT JOIN users u2 ON inserted.substituted_teacher_id = u2.id
	`
	var groupID, subTeacherID sql.NullString
	if lesson.GroupID != nil && *lesson.GroupID != "" {
		groupID = sql.NullString{String: *lesson.GroupID, Valid: true}
	}
	if lesson.SubstitutedTeacherID != nil && *lesson.SubstitutedTeacherID != "" {
		subTeacherID = sql.NullString{String: *lesson.SubstitutedTeacherID, Valid: true}
	}
	if lesson.ActivityType == "" {
		lesson.ActivityType = "standard"
	}

	return r.db.QueryRow(query,
		lesson.ClassID, lesson.TeacherID, lesson.SubjectID, lesson.Date,
		lesson.Hour, lesson.Duration, lesson.Topic, lesson.Type,
		groupID, lesson.IsSubstitution, subTeacherID, lesson.ActivityType,
	).Scan(&lesson.ID, &lesson.TeacherName, &lesson.SubstitutedTeacherName)
}

func (r *repository) GetLessonsByClass(classID string, date string) ([]Lesson, error) {
	query := `
		SELECT cl.id, cl.class_id, cl.teacher_id, COALESCE(u1.first_name || ' ' || u1.last_name, '') AS teacher_name,
		       cl.subject_id, cl.date, cl.hour, cl.duration, cl.topic, cl.type,
		       cl.group_id, cl.is_substitution, cl.substituted_teacher_id,
		       COALESCE(u2.first_name || ' ' || u2.last_name, '') AS substituted_teacher_name,
		       cl.activity_type, cl.created_at, cl.updated_at
		FROM class_lessons cl
		LEFT JOIN users u1 ON cl.teacher_id = u1.id
		LEFT JOIN users u2 ON cl.substituted_teacher_id = u2.id
		WHERE cl.class_id = $1
	`
	args := []interface{}{classID}
	if date != "" {
		query += " AND cl.date = $2"
		args = append(args, date)
	}
	query += " ORDER BY cl.date DESC, cl.hour DESC"

	return r.scanLessons(query, args...)
}

func (r *repository) GetLessonsByClassAndSubject(classID, subjectID string, date string) ([]Lesson, error) {
	query := `
		SELECT cl.id, cl.class_id, cl.teacher_id, COALESCE(u1.first_name || ' ' || u1.last_name, '') AS teacher_name,
		       cl.subject_id, cl.date, cl.hour, cl.duration, cl.topic, cl.type,
		       cl.group_id, cl.is_substitution, cl.substituted_teacher_id,
		       COALESCE(u2.first_name || ' ' || u2.last_name, '') AS substituted_teacher_name,
		       cl.activity_type, cl.created_at, cl.updated_at
		FROM class_lessons cl
		LEFT JOIN users u1 ON cl.teacher_id = u1.id
		LEFT JOIN users u2 ON cl.substituted_teacher_id = u2.id
		WHERE cl.class_id = $1 AND cl.subject_id = $2
	`
	args := []interface{}{classID, subjectID}
	if date != "" {
		query += " AND cl.date = $3"
		args = append(args, date)
	}
	query += " ORDER BY cl.date DESC, cl.hour DESC"

	return r.scanLessons(query, args...)
}

func (r *repository) GetLessonsByGroup(groupID string, date string) ([]Lesson, error) {
	query := `
		SELECT cl.id, cl.class_id, cl.teacher_id, COALESCE(u1.first_name || ' ' || u1.last_name, '') AS teacher_name,
		       cl.subject_id, cl.date, cl.hour, cl.duration, cl.topic, cl.type,
		       cl.group_id, cl.is_substitution, cl.substituted_teacher_id,
		       COALESCE(u2.first_name || ' ' || u2.last_name, '') AS substituted_teacher_name,
		       cl.activity_type, cl.created_at, cl.updated_at
		FROM class_lessons cl
		LEFT JOIN users u1 ON cl.teacher_id = u1.id
		LEFT JOIN users u2 ON cl.substituted_teacher_id = u2.id
		WHERE cl.group_id = $1
	`
	args := []interface{}{groupID}
	if date != "" {
		query += " AND cl.date = $2"
		args = append(args, date)
	}
	query += " ORDER BY cl.date DESC, cl.hour DESC"

	return r.scanLessons(query, args...)
}

func (r *repository) scanLessons(query string, args ...interface{}) ([]Lesson, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []Lesson
	for rows.Next() {
		var l Lesson
		var groupID, subTeacherID sql.NullString
		if err := rows.Scan(
			&l.ID, &l.ClassID, &l.TeacherID, &l.TeacherName,
			&l.SubjectID, &l.Date, &l.Hour, &l.Duration, &l.Topic, &l.Type,
			&groupID, &l.IsSubstitution, &subTeacherID,
			&l.SubstitutedTeacherName, &l.ActivityType, &l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if groupID.Valid {
			l.GroupID = &groupID.String
		}
		if subTeacherID.Valid {
			l.SubstitutedTeacherID = &subTeacherID.String
		}
		lessons = append(lessons, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lessons, nil
}

func (r *repository) GetLessonsByTeacher(teacherID string, fromDate, toDate string) ([]Lesson, error) {
	query := `
		SELECT cl.id, cl.class_id, cl.teacher_id, COALESCE(u1.first_name || ' ' || u1.last_name, '') AS teacher_name,
		       cl.subject_id, cl.date, cl.hour, cl.duration, cl.topic, cl.type,
		       cl.group_id, cl.is_substitution, cl.substituted_teacher_id,
		       COALESCE(u2.first_name || ' ' || u2.last_name, '') AS substituted_teacher_name,
		       cl.activity_type, cl.created_at, cl.updated_at
		FROM class_lessons cl
		LEFT JOIN users u1 ON cl.teacher_id = u1.id
		LEFT JOIN users u2 ON cl.substituted_teacher_id = u2.id
		WHERE cl.teacher_id = $1::uuid
		  AND ($2 = '' OR cl.date >= $2::date)
		  AND ($3 = '' OR cl.date <= $3::date)
		ORDER BY cl.date DESC, cl.hour ASC
	`
	return r.scanLessons(query, teacherID, fromDate, toDate)
}

func (r *repository) CreateHomework(homework *Homework) error {
	if homework.Type == "" {
		homework.Type = "compito"
	}
	query := `
		WITH inserted AS (
			INSERT INTO class_homeworks (lesson_id, class_id, subject_id, teacher_id, due_date, description, type, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
			RETURNING id, teacher_id
		)
		SELECT inserted.id, COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name
		FROM inserted
		LEFT JOIN users u ON inserted.teacher_id = u.id
	`
	return r.db.QueryRow(query,
		homework.LessonID, homework.ClassID, homework.SubjectID, homework.TeacherID,
		homework.DueDate, homework.Description, homework.Type,
	).Scan(&homework.ID, &homework.TeacherName)
}

func (r *repository) GetHomeworkByClass(classID string) ([]Homework, error) {
	query := `SELECT ch.id, ch.lesson_id, ch.class_id, ch.subject_id, ch.teacher_id, COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name, ch.due_date, ch.description, COALESCE(ch.type, 'compito'), ch.created_at, ch.updated_at 
	          FROM class_homeworks ch
	          LEFT JOIN users u ON ch.teacher_id = u.id
	          WHERE ch.class_id = $1 ORDER BY ch.due_date ASC`
	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var homeworks []Homework
	for rows.Next() {
		var h Homework
		var lessonID sql.NullString
		if err := rows.Scan(&h.ID, &lessonID, &h.ClassID, &h.SubjectID, &h.TeacherID, &h.TeacherName, &h.DueDate, &h.Description, &h.Type, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		if lessonID.Valid {
			h.LessonID = &lessonID.String
		}
		homeworks = append(homeworks, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return homeworks, nil
}

func (r *repository) GetLessonByID(id string) (*Lesson, error) {
	query := `
		SELECT cl.id, cl.class_id, cl.teacher_id, COALESCE(u1.first_name || ' ' || u1.last_name, '') AS teacher_name,
		       cl.subject_id, cl.date, cl.hour, cl.duration, cl.topic, cl.type,
		       cl.group_id, cl.is_substitution, cl.substituted_teacher_id,
		       COALESCE(u2.first_name || ' ' || u2.last_name, '') AS substituted_teacher_name,
		       cl.activity_type, cl.created_at, cl.updated_at
		FROM class_lessons cl
		LEFT JOIN users u1 ON cl.teacher_id = u1.id
		LEFT JOIN users u2 ON cl.substituted_teacher_id = u2.id
		WHERE cl.id = $1::uuid
	`
	lessons, err := r.scanLessons(query, id)
	if err != nil {
		return nil, err
	}
	if len(lessons) == 0 {
		return nil, sql.ErrNoRows
	}
	return &lessons[0], nil
}

func (r *repository) UpdateLesson(id string, req UpdateLessonRequest) (*Lesson, error) {
	query := `
		UPDATE class_lessons
		SET topic = COALESCE(NULLIF($2, ''), topic),
		    type = COALESCE(NULLIF($3, ''), type),
		    activity_type = COALESCE(NULLIF($4, ''), activity_type),
		    updated_at = NOW()
		WHERE id = $1::uuid
	`
	_, err := r.db.Exec(query, id, req.Topic, req.Type, req.ActivityType)
	if err != nil {
		return nil, err
	}
	return r.GetLessonByID(id)
}

func (r *repository) DeleteLesson(id string) error {
	_, err := r.db.Exec("DELETE FROM class_lessons WHERE id = $1::uuid", id)
	return err
}

func (r *repository) UpdateHomework(id string, req UpdateHomeworkRequest) (*Homework, error) {
	query := `
		UPDATE class_homeworks
		SET description = COALESCE(NULLIF($2, ''), description),
		    type = COALESCE(NULLIF($3, ''), type),
		    updated_at = NOW()
		WHERE id = $1::uuid
	`
	_, err := r.db.Exec(query, id, req.Description, req.Type)
	if err != nil {
		return nil, err
	}
	queryGet := `SELECT ch.id, ch.lesson_id, ch.class_id, ch.subject_id, ch.teacher_id, COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name, ch.due_date, ch.description, COALESCE(ch.type, 'compito'), ch.created_at, ch.updated_at 
	             FROM class_homeworks ch
	             LEFT JOIN users u ON ch.teacher_id = u.id
	             WHERE ch.id = $1::uuid`
	var h Homework
	var lessonID sql.NullString
	if err := r.db.QueryRow(queryGet, id).Scan(&h.ID, &lessonID, &h.ClassID, &h.SubjectID, &h.TeacherID, &h.TeacherName, &h.DueDate, &h.Description, &h.Type, &h.CreatedAt, &h.UpdatedAt); err != nil {
		return nil, err
	}
	if lessonID.Valid {
		h.LessonID = &lessonID.String
	}
	return &h, nil
}

func (r *repository) GetHomeworkByID(id string) (*Homework, error) {
	queryGet := `SELECT ch.id, ch.lesson_id, ch.class_id, ch.subject_id, ch.teacher_id, COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name, ch.due_date, ch.description, COALESCE(ch.type, 'compito'), ch.created_at, ch.updated_at 
	             FROM class_homeworks ch
	             LEFT JOIN users u ON ch.teacher_id = u.id
	             WHERE ch.id = $1::uuid`
	var h Homework
	var lessonID sql.NullString
	if err := r.db.QueryRow(queryGet, id).Scan(&h.ID, &lessonID, &h.ClassID, &h.SubjectID, &h.TeacherID, &h.TeacherName, &h.DueDate, &h.Description, &h.Type, &h.CreatedAt, &h.UpdatedAt); err != nil {
		return nil, err
	}
	if lessonID.Valid {
		h.LessonID = &lessonID.String
	}
	return &h, nil
}

func (r *repository) DeleteHomework(id string) error {
	_, err := r.db.Exec("DELETE FROM class_homeworks WHERE id = $1::uuid", id)
	return err
}
