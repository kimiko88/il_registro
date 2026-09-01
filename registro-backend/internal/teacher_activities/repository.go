package teacher_activities

import (
	"database/sql"
)

// Repository definisce le operazioni di persistenza per le attività libere.
type Repository interface {
	Create(act *TeacherFreeActivity) error
	GetByID(id string) (*TeacherFreeActivity, error)
	Update(id string, req UpdateTeacherActivityRequest) (*TeacherFreeActivity, error)
	Delete(id string) error
	GetByTeacher(teacherID, fromDate, toDate string) ([]TeacherFreeActivity, error)
}

type repository struct {
	db *sql.DB
}

// NewRepository restituisce una nuova implementazione del Repository.
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

const selectFields = `
	tfa.id, tfa.teacher_id,
	COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name,
	tfa.date, tfa.start_hour, tfa.duration,
	tfa.activity_type, tfa.description,
	COALESCE(tfa.notes, '') AS notes,
	tfa.created_at, tfa.updated_at
`

const baseFrom = `
	FROM teacher_free_activities tfa
	LEFT JOIN users u ON tfa.teacher_id = u.id
`

func (r *repository) Create(act *TeacherFreeActivity) error {
	if act.ActivityType == "" {
		act.ActivityType = "disponibilita"
	}
	query := `
		WITH inserted AS (
			INSERT INTO teacher_free_activities
				(teacher_id, date, start_hour, duration, activity_type, description, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
			RETURNING id, teacher_id, date, start_hour, duration, activity_type, description,
			          COALESCE(notes, '') AS notes, created_at, updated_at
		)
		SELECT inserted.id, inserted.teacher_id,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name,
		       inserted.date, inserted.start_hour, inserted.duration,
		       inserted.activity_type, inserted.description, inserted.notes,
		       inserted.created_at, inserted.updated_at
		FROM inserted
		LEFT JOIN users u ON inserted.teacher_id = u.id
	`
	return r.db.QueryRow(query,
		act.TeacherID, act.Date, act.StartHour, act.Duration,
		act.ActivityType, act.Description, act.Notes,
	).Scan(
		&act.ID, &act.TeacherID, &act.TeacherName,
		&act.Date, &act.StartHour, &act.Duration,
		&act.ActivityType, &act.Description, &act.Notes,
		&act.CreatedAt, &act.UpdatedAt,
	)
}

func (r *repository) GetByID(id string) (*TeacherFreeActivity, error) {
	query := `SELECT ` + selectFields + baseFrom + `WHERE tfa.id = $1::uuid`
	var act TeacherFreeActivity
	err := r.db.QueryRow(query, id).Scan(
		&act.ID, &act.TeacherID, &act.TeacherName,
		&act.Date, &act.StartHour, &act.Duration,
		&act.ActivityType, &act.Description, &act.Notes,
		&act.CreatedAt, &act.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &act, nil
}

func (r *repository) Update(id string, req UpdateTeacherActivityRequest) (*TeacherFreeActivity, error) {
	query := `
		UPDATE teacher_free_activities
		SET description  = CASE WHEN $2 <> '' THEN $2 ELSE description END,
		    activity_type = CASE WHEN $3 <> '' THEN $3 ELSE activity_type END,
		    notes        = CASE WHEN $4 <> '' THEN $4 ELSE notes END,
		    start_hour   = COALESCE($5, start_hour),
		    duration     = COALESCE($6, duration),
		    date         = CASE WHEN $7 <> '' THEN $7::date ELSE date END,
		    updated_at   = NOW()
		WHERE id = $1::uuid
	`
	var startHour, duration *int
	if req.StartHour != nil {
		startHour = req.StartHour
	}
	if req.Duration != nil {
		duration = req.Duration
	}
	_, err := r.db.Exec(query, id,
		req.Description, req.ActivityType, req.Notes,
		startHour, duration, req.Date,
	)
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *repository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM teacher_free_activities WHERE id = $1::uuid", id)
	return err
}

func (r *repository) GetByTeacher(teacherID, fromDate, toDate string) ([]TeacherFreeActivity, error) {
	query := `
		SELECT ` + selectFields + baseFrom + `
		WHERE tfa.teacher_id = $1::uuid
		  AND ($2 = '' OR tfa.date >= $2::date)
		  AND ($3 = '' OR tfa.date <= $3::date)
		ORDER BY tfa.date DESC, tfa.start_hour ASC
	`
	rows, err := r.db.Query(query, teacherID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []TeacherFreeActivity
	for rows.Next() {
		var act TeacherFreeActivity
		if err := rows.Scan(
			&act.ID, &act.TeacherID, &act.TeacherName,
			&act.Date, &act.StartHour, &act.Duration,
			&act.ActivityType, &act.Description, &act.Notes,
			&act.CreatedAt, &act.UpdatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, act)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
