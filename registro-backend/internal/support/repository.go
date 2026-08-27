package support

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	CreateDiaryEntry(ctx context.Context, entry *SupportDiaryEntry) error
	ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]SupportDiaryEntry, error)
	DeleteDiaryEntry(ctx context.Context, id, teacherID string) error

	CreatePeiGoal(ctx context.Context, goal *SupportPeiGoal) error
	UpdatePeiGoalProgress(ctx context.Context, id, progressStatus string) error
	ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]SupportPeiGoal, error)
	DeletePeiGoal(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateDiaryEntry(ctx context.Context, e *SupportDiaryEntry) error {
	query := `
		INSERT INTO support_diaries (
			id, school_id, teacher_id, student_id, class_id, entry_date,
			time_slot, co_teacher_id, activity_type, topic_and_activities,
			student_responses, educator_notes, is_shared_with_family,
			created_at, updated_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4, $5, $6::date,
			$7, NULLIF($8, '')::uuid, $9, $10,
			$11, $12, $13,
			NOW(), NOW()
		) RETURNING id
	`
	coTeacher := ""
	if e.CoTeacherID != nil {
		coTeacher = *e.CoTeacherID
	}
	return r.db.QueryRowContext(ctx, query,
		e.ID, e.SchoolID, e.TeacherID, e.StudentID, e.ClassID, e.EntryDate,
		e.TimeSlot, coTeacher, e.ActivityType, e.TopicAndActivities,
		e.StudentResponses, e.EducatorNotes, e.IsSharedWithFamily,
	).Scan(&e.ID)
}

func (r *repository) ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]SupportDiaryEntry, error) {
	familyFilter := ""
	if isFamily {
		familyFilter = " AND sd.is_shared_with_family = TRUE "
	}

	query := fmt.Sprintf(`
		SELECT sd.id, sd.school_id, sd.teacher_id,
		       TRIM(COALESCE(tu.first_name, '') || ' ' || COALESCE(tu.last_name, '')) as teacher_name,
		       sd.student_id,
		       TRIM(COALESCE(su.first_name, '') || ' ' || COALESCE(su.last_name, '')) as student_name,
		       sd.class_id, cls.name as class_name,
		       sd.entry_date::text, sd.time_slot,
		       sd.co_teacher_id::text,
		       TRIM(COALESCE(cu.first_name, '') || ' ' || COALESCE(cu.last_name, '')) as co_teacher_name,
		       sd.activity_type, sd.topic_and_activities, sd.student_responses,
		       sd.educator_notes, sd.is_shared_with_family, sd.created_at, sd.updated_at
		FROM support_diaries sd
		JOIN teachers t ON sd.teacher_id = t.id
		JOIN users tu ON t.user_id = tu.id
		JOIN students s ON sd.student_id = s.id
		JOIN users su ON s.user_id = su.id
		JOIN classes cls ON sd.class_id = cls.id
		LEFT JOIN teachers ct ON sd.co_teacher_id = ct.id
		LEFT JOIN users cu ON ct.user_id = cu.id
		WHERE sd.school_id = $1
		  %s
		  AND ($2 = '' OR sd.student_id = NULLIF($2, '')::uuid)
		  AND ($3 = '' OR sd.teacher_id = NULLIF($3, '')::uuid)
		  AND ($4 = '' OR sd.class_id = NULLIF($4, '')::uuid)
		ORDER BY sd.entry_date DESC, sd.created_at DESC
	`, familyFilter)

	rows, err := r.db.QueryContext(ctx, query, schoolID, studentID, teacherID, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SupportDiaryEntry
	for rows.Next() {
		var e SupportDiaryEntry
		var coID, coName sql.NullString
		if err := rows.Scan(
			&e.ID, &e.SchoolID, &e.TeacherID, &e.TeacherName,
			&e.StudentID, &e.StudentName,
			&e.ClassID, &e.ClassName,
			&e.EntryDate, &e.TimeSlot,
			&coID, &coName,
			&e.ActivityType, &e.TopicAndActivities, &e.StudentResponses,
			&e.EducatorNotes, &e.IsSharedWithFamily, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if coID.Valid {
			e.CoTeacherID = &coID.String
			e.CoTeacherName = coName.String
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *repository) DeleteDiaryEntry(ctx context.Context, id, teacherID string) error {
	query := `DELETE FROM support_diaries WHERE id = $1 AND ($2 = '' OR teacher_id = NULLIF($2, '')::uuid)`
	_, err := r.db.ExecContext(ctx, query, id, teacherID)
	return err
}

func (r *repository) CreatePeiGoal(ctx context.Context, g *SupportPeiGoal) error {
	query := `
		INSERT INTO support_pei_goals (
			id, school_id, student_id, pei_type, axis, title, description,
			expected_term, progress_status, created_at, updated_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4, $5, $6, $7,
			$8, $9, NOW(), NOW()
		) RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		g.ID, g.SchoolID, g.StudentID, g.PeiType, g.Axis, g.Title, g.Description,
		g.ExpectedTerm, g.ProgressStatus,
	).Scan(&g.ID)
}

func (r *repository) UpdatePeiGoalProgress(ctx context.Context, id, progressStatus string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE support_pei_goals SET progress_status = $1, updated_at = NOW() WHERE id = $2`, progressStatus, id)
	return err
}

func (r *repository) ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]SupportPeiGoal, error) {
	query := `
		SELECT pg.id, pg.school_id, pg.student_id,
		       TRIM(COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '')) as student_name,
		       pg.pei_type, pg.axis, pg.title, pg.description, pg.expected_term, pg.progress_status,
		       pg.created_at, pg.updated_at
		FROM support_pei_goals pg
		JOIN students s ON pg.student_id = s.id
		JOIN users u ON s.user_id = u.id
		WHERE pg.school_id = $1 AND ($2 = '' OR pg.student_id = NULLIF($2, '')::uuid)
		ORDER BY pg.axis ASC, pg.created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SupportPeiGoal
	for rows.Next() {
		var g SupportPeiGoal
		if err := rows.Scan(
			&g.ID, &g.SchoolID, &g.StudentID, &g.StudentName,
			&g.PeiType, &g.Axis, &g.Title, &g.Description, &g.ExpectedTerm, &g.ProgressStatus,
			&g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, g)
	}
	return list, rows.Err()
}

func (r *repository) DeletePeiGoal(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM support_pei_goals WHERE id = $1`, id)
	return err
}
