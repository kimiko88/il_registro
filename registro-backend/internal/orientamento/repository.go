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
	rows, err := r.db.QueryContext(ctx, `SELECT id, school_id, title, description, category, date, end_date, location, hours, max_attendees, created_by FROM orientamento_events WHERE school_id=$1 ORDER BY date DESC`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var e Event
		_ = rows.Scan(&e.ID, &e.SchoolID, &e.Title, &e.Description, &e.Category, &e.Date, &e.EndDate, &e.Location, &e.Hours, &e.MaxAttendees, &e.CreatedBy)
		events = append(events, e)
	}
	return events, nil
}

func (r *repository) RegisterStudent(ctx context.Context, p *Participation) error {
	return r.db.QueryRowContext(ctx, `INSERT INTO orientamento_participations (event_id, student_id, status) VALUES ($1, $2, 'Registered') RETURNING id`,
		p.EventID, p.StudentID).Scan(&p.ID)
}

func (r *repository) GetParticipations(ctx context.Context, studentID string) ([]Participation, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, event_id, student_id, status, attended, registered_at FROM orientamento_participations WHERE student_id=$1`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var parts []Participation
	for rows.Next() {
		var p Participation
		_ = rows.Scan(&p.ID, &p.EventID, &p.StudentID, &p.Status, &p.Attended, &p.RegisteredAt)
		parts = append(parts, p)
	}
	return parts, nil
}

func (r *repository) MarkAttendance(ctx context.Context, eventID, studentID string, attended bool) error {
	status := "NoShow"
	if attended {
		status = "Attended"
	}
	_, err := r.db.ExecContext(ctx, `UPDATE orientamento_participations SET attended=$1, status=$2 WHERE event_id=$3 AND student_id=$4`, attended, status, eventID, studentID)
	return err
}
