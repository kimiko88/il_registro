package orientamento

import (
	"context"
	"database/sql"
)

type Repository interface {
	CreateEvent(ctx context.Context, e *Event) error
	GetEvents(ctx context.Context, schoolID string) ([]Event, error)
	RegisterStudent(ctx context.Context, p *Participation) error
	GetParticipations(ctx context.Context, studentID string) ([]Participation, error)
	MarkAttendance(ctx context.Context, eventID, studentID string, attended bool) error
	SavePreference(ctx context.Context, p *StudentPreference) error
	GetPreference(ctx context.Context, studentID string) (*StudentPreference, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateEvent(ctx context.Context, e *Event) error {
	query := `
		INSERT INTO orientamento_events (school_id, title, description, category, date, end_date, location, hours, max_attendees, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	return r.db.QueryRowContext(ctx, query,
		e.SchoolID, e.Title, e.Description, e.Category, e.Date, e.EndDate, e.Location, e.Hours, e.MaxAttendees, e.CreatedBy).Scan(&e.ID)
}

func (r *repository) GetEvents(ctx context.Context, schoolID string) ([]Event, error) {
	var rows *sql.Rows
	var err error
	if schoolID != "" {
		rows, err = r.db.QueryContext(ctx, `SELECT id, school_id, title, description, category, date, end_date, location, hours, max_attendees, created_by FROM orientamento_events WHERE school_id=$1 ORDER BY date DESC`, schoolID)
	} else {
		rows, err = r.db.QueryContext(ctx, `SELECT id, school_id, title, description, category, date, end_date, location, hours, max_attendees, created_by FROM orientamento_events ORDER BY date DESC`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.SchoolID, &e.Title, &e.Description, &e.Category, &e.Date, &e.EndDate, &e.Location, &e.Hours, &e.MaxAttendees, &e.CreatedBy); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *repository) RegisterStudent(ctx context.Context, p *Participation) error {
	query := `
		INSERT INTO orientamento_participations (event_id, student_id, status)
		VALUES ($1, COALESCE((SELECT id FROM students WHERE user_id = $2::uuid OR id = $2::uuid LIMIT 1), $2::uuid), 'Registered')
		RETURNING id`
	return r.db.QueryRowContext(ctx, query, p.EventID, p.StudentID).Scan(&p.ID)
}

func (r *repository) GetParticipations(ctx context.Context, studentID string) ([]Participation, error) {
	query := `
		SELECT p.id, p.event_id, p.student_id, p.status, p.attended, p.registered_at,
		       e.id, e.school_id, e.title, e.description, e.category, e.date, e.end_date, e.location, e.hours, e.max_attendees, e.created_by
		FROM orientamento_participations p
		LEFT JOIN orientamento_events e ON p.event_id = e.id
		WHERE (
			p.student_id = $1::uuid OR
			p.student_id IN (SELECT id FROM students WHERE user_id = $1::uuid OR id = $1::uuid) OR
			p.student_id IN (SELECT user_id FROM students WHERE id = $1::uuid OR user_id = $1::uuid)
		)
		ORDER BY e.date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var parts []Participation
	for rows.Next() {
		var p Participation
		var e Event
		var eID, eSchoolID, eTitle, eDesc, eCat, eLoc, eCreatedBy sql.NullString
		var eDate, eEndDate sql.NullTime
		var eHours sql.NullFloat64
		var eMaxAttendees sql.NullInt64

		err := rows.Scan(
			&p.ID, &p.EventID, &p.StudentID, &p.Status, &p.Attended, &p.RegisteredAt,
			&eID, &eSchoolID, &eTitle, &eDesc, &eCat, &eDate, &eEndDate, &eLoc, &eHours, &eMaxAttendees, &eCreatedBy,
		)
		if err != nil {
			return nil, err
		}
		if eID.Valid {
			e.ID = eID.String
			e.SchoolID = eSchoolID.String
			e.Title = eTitle.String
			e.Description = eDesc.String
			e.Category = eCat.String
			if eDate.Valid {
				e.Date = eDate.Time
			}
			if eEndDate.Valid {
				e.EndDate = eEndDate.Time
			}
			e.Location = eLoc.String
			e.Hours = eHours.Float64
			e.MaxAttendees = int(eMaxAttendees.Int64)
			e.CreatedBy = eCreatedBy.String
			p.Event = &e
		}
		parts = append(parts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return parts, nil
}


func (r *repository) MarkAttendance(ctx context.Context, eventID, studentID string, attended bool) error {
	status := "NoShow"
	if attended {
		status = "Attended"
	}
	query := `
		UPDATE orientamento_participations 
		SET attended = $1, status = $2 
		WHERE event_id = $3::uuid AND (
			student_id = $4::uuid OR 
			student_id IN (SELECT id FROM students WHERE user_id = $4::uuid OR id = $4::uuid) OR
			student_id IN (SELECT user_id FROM students WHERE id = $4::uuid OR user_id = $4::uuid)
		)
	`
	_, err := r.db.ExecContext(ctx, query, attended, status, eventID, studentID)
	return err
}


func (r *repository) SavePreference(ctx context.Context, p *StudentPreference) error {
	query := `
		INSERT INTO orientamento_preferences (student_id, preferred_track, target_field, notes, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (student_id) DO UPDATE SET
			preferred_track = EXCLUDED.preferred_track,
			target_field = EXCLUDED.target_field,
			notes = EXCLUDED.notes,
			updated_at = NOW()`
	_, err := r.db.ExecContext(ctx, query, p.StudentID, p.PreferredTrack, p.TargetField, p.Notes)
	return err
}

func (r *repository) GetPreference(ctx context.Context, studentID string) (*StudentPreference, error) {
	var p StudentPreference
	query := `SELECT id, student_id, preferred_track, target_field, notes, updated_at FROM orientamento_preferences WHERE student_id = $1`
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(&p.ID, &p.StudentID, &p.PreferredTrack, &p.TargetField, &p.Notes, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &StudentPreference{StudentID: studentID}, nil
		}
		return nil, err
	}
	return &p, nil
}
