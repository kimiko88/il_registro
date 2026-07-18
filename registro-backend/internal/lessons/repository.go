package lessons

import (
	"database/sql"
)

type Repository interface {
	CreateLesson(lesson *Lesson) error
	GetLessonsByClass(classID string, date string) ([]Lesson, error)
	GetLessonsByClassAndSubject(classID, subjectID string, date string) ([]Lesson, error)
	
	CreateHomework(homework *Homework) error
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
			INSERT INTO class_lessons (class_id, teacher_id, subject_id, date, hour, duration, topic, type, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
			RETURNING id, teacher_id
		)
		SELECT inserted.id, COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name
		FROM inserted
		LEFT JOIN users u ON inserted.teacher_id = u.id
	`
	return r.db.QueryRow(query,
		lesson.ClassID, lesson.TeacherID, lesson.SubjectID, lesson.Date,
		lesson.Hour, lesson.Duration, lesson.Topic, lesson.Type, lesson.Notes,
	).Scan(&lesson.ID, &lesson.TeacherName)
}

func (r *repository) GetLessonsByClass(classID string, date string) ([]Lesson, error) {
	query := `SELECT cl.id, cl.class_id, cl.teacher_id, COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name, cl.subject_id, cl.date, cl.hour, cl.duration, cl.topic, cl.type, cl.notes, cl.created_at, cl.updated_at 
	          FROM class_lessons cl
	          LEFT JOIN users u ON cl.teacher_id = u.id
	          WHERE cl.class_id = $1`
	args := []interface{}{classID}
	if date != "" {
		query += " AND cl.date = $2"
		args = append(args, date)
	}
	query += " ORDER BY cl.date DESC, cl.hour DESC"
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []Lesson
	for rows.Next() {
		var l Lesson
		if err := rows.Scan(&l.ID, &l.ClassID, &l.TeacherID, &l.TeacherName, &l.SubjectID, &l.Date, &l.Hour, &l.Duration, &l.Topic, &l.Type, &l.Notes, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *repository) GetLessonsByClassAndSubject(classID, subjectID string, date string) ([]Lesson, error) {
	query := `SELECT cl.id, cl.class_id, cl.teacher_id, COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name, cl.subject_id, cl.date, cl.hour, cl.duration, cl.topic, cl.type, cl.notes, cl.created_at, cl.updated_at 
	          FROM class_lessons cl
	          LEFT JOIN users u ON cl.teacher_id = u.id
	          WHERE cl.class_id = $1 AND cl.subject_id = $2`
	args := []interface{}{classID, subjectID}
	if date != "" {
		query += " AND cl.date = $3"
		args = append(args, date)
	}
	query += " ORDER BY cl.date DESC, cl.hour DESC"
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []Lesson
	for rows.Next() {
		var l Lesson
		if err := rows.Scan(&l.ID, &l.ClassID, &l.TeacherID, &l.TeacherName, &l.SubjectID, &l.Date, &l.Hour, &l.Duration, &l.Topic, &l.Type, &l.Notes, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *repository) CreateHomework(homework *Homework) error {
	query := `
		WITH inserted AS (
			INSERT INTO class_homeworks (lesson_id, class_id, subject_id, teacher_id, due_date, description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
			RETURNING id, teacher_id
		)
		SELECT inserted.id, COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name
		FROM inserted
		LEFT JOIN users u ON inserted.teacher_id = u.id
	`
	return r.db.QueryRow(query,
		homework.LessonID, homework.ClassID, homework.SubjectID, homework.TeacherID,
		homework.DueDate, homework.Description,
	).Scan(&homework.ID, &homework.TeacherName)
}

func (r *repository) GetHomeworkByClass(classID string) ([]Homework, error) {
	query := `SELECT ch.id, ch.lesson_id, ch.class_id, ch.subject_id, ch.teacher_id, COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name, ch.due_date, ch.description, ch.created_at, ch.updated_at 
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
		if err := rows.Scan(&h.ID, &h.LessonID, &h.ClassID, &h.SubjectID, &h.TeacherID, &h.TeacherName, &h.DueDate, &h.Description, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		homeworks = append(homeworks, h)
	}
	return homeworks, nil
}
