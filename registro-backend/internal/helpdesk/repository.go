package helpdesk

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Repository interface {
	CreateSlot(ctx context.Context, slot *HelpDeskSlot) error
	GetSlot(ctx context.Context, id string) (*HelpDeskSlot, error)
	ListSlots(ctx context.Context, schoolID, subjectID, date, status string) ([]HelpDeskSlot, error)
	BookSlot(ctx context.Context, booking *HelpDeskBooking) error
	ListBookings(ctx context.Context, slotID string) ([]HelpDeskBooking, error)
	MarkAttendance(ctx context.Context, bookingID, status string) error
	CompleteSlot(ctx context.Context, slotID string) error
	GetFISReport(ctx context.Context, schoolID string) ([]FISAccountingSummary, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateSlot(ctx context.Context, slot *HelpDeskSlot) error {
	query := `
		INSERT INTO help_desk_slots (
			school_id, teacher_id, subject_id, slot_date, start_time, end_time, room_id, max_capacity, status
		) VALUES ($1, $2, $3, $4::date, $5::time, $6::time, $7, $8, $9)
		RETURNING id, created_at;
	`
	return r.db.QueryRowContext(ctx, query,
		slot.SchoolID, slot.TeacherID, slot.SubjectID, slot.SlotDate, slot.StartTime, slot.EndTime, slot.RoomID, slot.MaxCapacity, slot.Status,
	).Scan(&slot.ID, &slot.CreatedAt)
}

func (r *postgresRepo) GetSlot(ctx context.Context, id string) (*HelpDeskSlot, error) {
	query := `
		SELECT s.id, s.school_id, s.teacher_id, s.subject_id, s.slot_date::text, s.start_time::text, s.end_time::text,
		       s.room_id, s.max_capacity, s.status, s.created_at,
		       COALESCE(u.first_name || ' ' || u.last_name, 'Docente') as teacher_name,
		       COALESCE(sub.name, 'Materia') as subject_name,
		       (SELECT COUNT(*) FROM help_desk_bookings b WHERE b.slot_id = s.id AND b.status != 'cancelled') as bookings_count
		FROM help_desk_slots s
		LEFT JOIN users u ON u.id = s.teacher_id
		LEFT JOIN subjects sub ON sub.id = s.subject_id
		WHERE s.id = $1;
	`
	var s HelpDeskSlot
	var roomID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.SchoolID, &s.TeacherID, &s.SubjectID, &s.SlotDate, &s.StartTime, &s.EndTime,
		&roomID, &s.MaxCapacity, &s.Status, &s.CreatedAt, &s.TeacherName, &s.SubjectName, &s.BookingsCount,
	)
	if err != nil {
		return nil, err
	}
	if roomID.Valid {
		s.RoomID = &roomID.String
	}
	return &s, nil
}

func (r *postgresRepo) ListSlots(ctx context.Context, schoolID, subjectID, date, status string) ([]HelpDeskSlot, error) {
	query := `
		SELECT s.id, s.school_id, s.teacher_id, s.subject_id, s.slot_date::text, s.start_time::text, s.end_time::text,
		       s.room_id, s.max_capacity, s.status, s.created_at,
		       COALESCE(u.first_name || ' ' || u.last_name, 'Docente') as teacher_name,
		       COALESCE(sub.name, 'Materia') as subject_name,
		       (SELECT COUNT(*) FROM help_desk_bookings b WHERE b.slot_id = s.id AND b.status != 'cancelled') as bookings_count
		FROM help_desk_slots s
		LEFT JOIN users u ON u.id = s.teacher_id
		LEFT JOIN subjects sub ON sub.id = s.subject_id
		WHERE s.school_id = $1
		  AND ($2 = '' OR s.subject_id = $2::uuid)
		  AND ($3 = '' OR s.slot_date = $3::date)
		  AND ($4 = '' OR s.status = $4)
		ORDER BY s.slot_date ASC, s.start_time ASC;
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, subjectID, date, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []HelpDeskSlot
	for rows.Next() {
		var s HelpDeskSlot
		var roomID sql.NullString
		if err := rows.Scan(
			&s.ID, &s.SchoolID, &s.TeacherID, &s.SubjectID, &s.SlotDate, &s.StartTime, &s.EndTime,
			&roomID, &s.MaxCapacity, &s.Status, &s.CreatedAt, &s.TeacherName, &s.SubjectName, &s.BookingsCount,
		); err != nil {
			return nil, err
		}
		if roomID.Valid {
			s.RoomID = &roomID.String
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *postgresRepo) BookSlot(ctx context.Context, booking *HelpDeskBooking) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Lock slot row and check capacity
	var maxCapacity int
	var currentCount int
	err = tx.QueryRowContext(ctx, "SELECT max_capacity FROM help_desk_slots WHERE id = $1 FOR UPDATE", booking.SlotID).Scan(&maxCapacity)
	if err != nil {
		return errors.New("slot help desk non trovato")
	}

	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM help_desk_bookings WHERE slot_id = $1 AND status != 'cancelled'", booking.SlotID).Scan(&currentCount)
	if err != nil {
		return err
	}

	if currentCount >= maxCapacity {
		return errors.New("lo slot è al completo (raggiunta capienza massima)")
	}

	// 2. Insert booking
	query := `
		INSERT INTO help_desk_bookings (slot_id, student_id, topic_description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, booked_at;
	`
	if err := tx.QueryRowContext(ctx, query, booking.SlotID, booking.StudentID, booking.TopicDescription, "booked").Scan(&booking.ID, &booking.BookedAt); err != nil {
		return errors.New("prenotazione già esistente per questo studente")
	}

	// 3. If full, update status to fully_booked
	if currentCount+1 >= maxCapacity {
		_, _ = tx.ExecContext(ctx, "UPDATE help_desk_slots SET status = 'fully_booked' WHERE id = $1", booking.SlotID)
	}

	return tx.Commit()
}

func (r *postgresRepo) ListBookings(ctx context.Context, slotID string) ([]HelpDeskBooking, error) {
	query := `
		SELECT b.id, b.slot_id, b.student_id, b.topic_description, b.status, b.booked_at, b.attended_at,
		       COALESCE(u.first_name || ' ' || u.last_name, 'Studente') as student_name
		FROM help_desk_bookings b
		LEFT JOIN users u ON u.id = b.student_id
		WHERE b.slot_id = $1
		ORDER BY b.booked_at ASC;
	`
	rows, err := r.db.QueryContext(ctx, query, slotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []HelpDeskBooking
	for rows.Next() {
		var b HelpDeskBooking
		var attendedAt sql.NullTime
		if err := rows.Scan(&b.ID, &b.SlotID, &b.StudentID, &b.TopicDescription, &b.Status, &b.BookedAt, &attendedAt, &b.StudentName); err != nil {
			return nil, err
		}
		if attendedAt.Valid {
			b.AttendedAt = &attendedAt.Time
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *postgresRepo) MarkAttendance(ctx context.Context, bookingID, status string) error {
	now := time.Now()
	query := `
		UPDATE help_desk_bookings
		SET status = $2, attended_at = $3
		WHERE id = $1;
	`
	_, err := r.db.ExecContext(ctx, query, bookingID, status, now)
	return err
}

func (r *postgresRepo) CompleteSlot(ctx context.Context, slotID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE help_desk_slots SET status = 'completed' WHERE id = $1", slotID)
	return err
}

func (r *postgresRepo) GetFISReport(ctx context.Context, schoolID string) ([]FISAccountingSummary, error) {
	query := `
		SELECT s.teacher_id,
		       COALESCE(u.first_name || ' ' || u.last_name, 'Docente') as teacher_name,
		       COUNT(DISTINCT s.id) as completed_slots,
		       COUNT(b.id) FILTER (WHERE b.status = 'attended') as attended_count
		FROM help_desk_slots s
		LEFT JOIN users u ON u.id = s.teacher_id
		LEFT JOIN help_desk_bookings b ON b.slot_id = s.id
		WHERE s.school_id = $1 AND s.status = 'completed'
		GROUP BY s.teacher_id, u.first_name, u.last_name;
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []FISAccountingSummary
	for rows.Next() {
		var f FISAccountingSummary
		if err := rows.Scan(&f.TeacherID, &f.TeacherName, &f.CompletedSlots, &f.AttendedCount); err != nil {
			return nil, err
		}
		// Default 1.5h per slot completed
		f.TotalHours = float64(f.CompletedSlots) * 1.5
		list = append(list, f)
	}
	return list, nil
}
