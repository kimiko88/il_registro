package homework_submissions

import (
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	UpsertSubmission(s *HomeworkSubmission) error
	GetSubmission(id string) (*HomeworkSubmission, error)
	GetSubmissionByHomeworkAndStudent(homeworkID, studentID string) (*HomeworkSubmission, error)
	ListSubmissionsByHomework(homeworkID string) ([]HomeworkSubmission, error)
	ListSubmissionsByStudent(studentID string) ([]HomeworkSubmission, error)
	SaveGrading(id string, req GradeSubmissionRequest, teacherID string) error
	SetClassroomSync(id, classroomSubmissionID string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) UpsertSubmission(s *HomeworkSubmission) error {
	query := `
		INSERT INTO homework_submissions (
			homework_id, student_id, file_url, file_type, submitted_at, status
		) VALUES ($1, $2, $3, $4, NOW(), 'consegnato')
		ON CONFLICT (homework_id, student_id)
		DO UPDATE SET
			file_url = EXCLUDED.file_url,
			file_type = EXCLUDED.file_type,
			submitted_at = NOW(),
			status = 'consegnato',
			updated_at = NOW()
		RETURNING id, submitted_at, created_at, updated_at
	`
	return r.db.QueryRow(
		query, s.HomeworkID, s.StudentID, s.FileURL, s.FileType,
	).Scan(&s.ID, &s.SubmittedAt, &s.CreatedAt, &s.UpdatedAt)
}

func (r *repository) GetSubmission(id string) (*HomeworkSubmission, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, sql.ErrNoRows
	}
	query := `
		SELECT s.id, s.homework_id, s.student_id, s.file_url, s.file_type, s.submitted_at,
		       COALESCE(s.teacher_annotated_file_url, ''), s.annotations_json::text,
		       s.teacher_grade, COALESCE(s.teacher_feedback, ''), s.status, s.graded_at, s.graded_by,
		       COALESCE(s.google_classroom_submission_id, ''), s.synced_with_classroom,
		       s.created_at, s.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name,
		       COALESCE(h.description, '') AS homework_title
		FROM homework_submissions s
		LEFT JOIN students st ON s.student_id = st.id
		LEFT JOIN users u ON (s.student_id = u.id OR st.user_id = u.id)
		LEFT JOIN class_homeworks h ON s.homework_id = h.id
		WHERE s.id = $1
	`
	var sub HomeworkSubmission
	err := r.db.QueryRow(query, id).Scan(
		&sub.ID, &sub.HomeworkID, &sub.StudentID, &sub.FileURL, &sub.FileType, &sub.SubmittedAt,
		&sub.TeacherAnnotatedFileURL, &sub.AnnotationsJSON, &sub.TeacherGrade, &sub.TeacherFeedback,
		&sub.Status, &sub.GradedAt, &sub.GradedBy, &sub.GoogleClassroomSubmissionID,
		&sub.SyncedWithClassroom, &sub.CreatedAt, &sub.UpdatedAt, &sub.StudentFullName, &sub.HomeworkTitle,
	)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *repository) GetSubmissionByHomeworkAndStudent(homeworkID, studentID string) (*HomeworkSubmission, error) {
	if _, err := uuid.Parse(homeworkID); err != nil {
		return nil, sql.ErrNoRows
	}
	query := `
		SELECT s.id, s.homework_id, s.student_id, s.file_url, s.file_type, s.submitted_at,
		       COALESCE(s.teacher_annotated_file_url, ''), s.annotations_json::text,
		       s.teacher_grade, COALESCE(s.teacher_feedback, ''), s.status, s.graded_at, s.graded_by,
		       COALESCE(s.google_classroom_submission_id, ''), s.synced_with_classroom,
		       s.created_at, s.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name,
		       COALESCE(h.description, '') AS homework_title
		FROM homework_submissions s
		LEFT JOIN students st ON s.student_id = st.id
		LEFT JOIN users u ON (s.student_id = u.id OR st.user_id = u.id)
		LEFT JOIN class_homeworks h ON s.homework_id = h.id
		WHERE s.homework_id = $1 AND (s.student_id = $2 OR st.user_id::text = $2)
	`
	var sub HomeworkSubmission
	err := r.db.QueryRow(query, homeworkID, studentID).Scan(
		&sub.ID, &sub.HomeworkID, &sub.StudentID, &sub.FileURL, &sub.FileType, &sub.SubmittedAt,
		&sub.TeacherAnnotatedFileURL, &sub.AnnotationsJSON, &sub.TeacherGrade, &sub.TeacherFeedback,
		&sub.Status, &sub.GradedAt, &sub.GradedBy, &sub.GoogleClassroomSubmissionID,
		&sub.SyncedWithClassroom, &sub.CreatedAt, &sub.UpdatedAt, &sub.StudentFullName, &sub.HomeworkTitle,
	)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *repository) ListSubmissionsByHomework(homeworkID string) ([]HomeworkSubmission, error) {
	if _, err := uuid.Parse(homeworkID); err != nil {
		return []HomeworkSubmission{}, nil
	}
	query := `
		SELECT s.id, s.homework_id, s.student_id, s.file_url, s.file_type, s.submitted_at,
		       COALESCE(s.teacher_annotated_file_url, ''), s.annotations_json::text,
		       s.teacher_grade, COALESCE(s.teacher_feedback, ''), s.status, s.graded_at, s.graded_by,
		       COALESCE(s.google_classroom_submission_id, ''), s.synced_with_classroom,
		       s.created_at, s.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name,
		       COALESCE(h.description, '') AS homework_title
		FROM homework_submissions s
		LEFT JOIN students st ON s.student_id = st.id
		LEFT JOIN users u ON (s.student_id = u.id OR st.user_id = u.id)
		LEFT JOIN class_homeworks h ON s.homework_id = h.id
		WHERE s.homework_id = $1
		ORDER BY s.submitted_at DESC
	`
	rows, err := r.db.Query(query, homeworkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []HomeworkSubmission
	for rows.Next() {
		var sub HomeworkSubmission
		if err := rows.Scan(
			&sub.ID, &sub.HomeworkID, &sub.StudentID, &sub.FileURL, &sub.FileType, &sub.SubmittedAt,
			&sub.TeacherAnnotatedFileURL, &sub.AnnotationsJSON, &sub.TeacherGrade, &sub.TeacherFeedback,
			&sub.Status, &sub.GradedAt, &sub.GradedBy, &sub.GoogleClassroomSubmissionID,
			&sub.SyncedWithClassroom, &sub.CreatedAt, &sub.UpdatedAt, &sub.StudentFullName, &sub.HomeworkTitle,
		); err != nil {
			return nil, err
		}
		list = append(list, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if list == nil {
		list = []HomeworkSubmission{}
	}
	return list, nil
}

func (r *repository) ListSubmissionsByStudent(studentID string) ([]HomeworkSubmission, error) {
	query := `
		SELECT s.id, s.homework_id, s.student_id, s.file_url, s.file_type, s.submitted_at,
		       COALESCE(s.teacher_annotated_file_url, ''), s.annotations_json::text,
		       s.teacher_grade, COALESCE(s.teacher_feedback, ''), s.status, s.graded_at, s.graded_by,
		       COALESCE(s.google_classroom_submission_id, ''), s.synced_with_classroom,
		       s.created_at, s.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name,
		       COALESCE(h.description, '') AS homework_title
		FROM homework_submissions s
		LEFT JOIN students st ON s.student_id = st.id
		LEFT JOIN users u ON (s.student_id = u.id OR st.user_id = u.id)
		LEFT JOIN class_homeworks h ON s.homework_id = h.id
		WHERE s.student_id = $1 OR st.user_id::text = $1
		ORDER BY s.submitted_at DESC
	`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []HomeworkSubmission
	for rows.Next() {
		var sub HomeworkSubmission
		if err := rows.Scan(
			&sub.ID, &sub.HomeworkID, &sub.StudentID, &sub.FileURL, &sub.FileType, &sub.SubmittedAt,
			&sub.TeacherAnnotatedFileURL, &sub.AnnotationsJSON, &sub.TeacherGrade, &sub.TeacherFeedback,
			&sub.Status, &sub.GradedAt, &sub.GradedBy, &sub.GoogleClassroomSubmissionID,
			&sub.SyncedWithClassroom, &sub.CreatedAt, &sub.UpdatedAt, &sub.StudentFullName, &sub.HomeworkTitle,
		); err != nil {
			return nil, err
		}
		list = append(list, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if list == nil {
		list = []HomeworkSubmission{}
	}
	return list, nil
}

func (r *repository) SaveGrading(id string, req GradeSubmissionRequest, teacherID string) error {
	annJSON := req.AnnotationsJSON
	if annJSON == "" {
		annJSON = "[]"
	}
	query := `
		UPDATE homework_submissions
		SET teacher_annotated_file_url = $2,
		    annotations_json = $3::jsonb,
		    teacher_grade = $4,
		    teacher_feedback = $5,
		    status = 'restituito',
		    graded_at = NOW(),
		    graded_by = $6,
		    synced_with_classroom = $7,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, req.AnnotatedURL, annJSON, req.TeacherGrade, req.TeacherFeedback, teacherID, req.SyncClassroom)
	return err
}

func (r *repository) SetClassroomSync(id, classroomSubmissionID string) error {
	query := `
		UPDATE homework_submissions
		SET google_classroom_submission_id = $2,
		    synced_with_classroom = TRUE,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, classroomSubmissionID)
	return err
}
