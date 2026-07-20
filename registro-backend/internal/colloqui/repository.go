package colloqui

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSlotFull      = errors.New("colloquio slot is fully booked")
	ErrAlreadyBooked = errors.New("parent has already booked this slot")
	ErrSlotCancelled = errors.New("colloquio slot has been cancelled")
)

type Repository interface {
	CreateSlot(ctx context.Context, slot *ColloquioSlot) error
	GetSlotByID(ctx context.Context, id string) (*ColloquioSlot, error)
	ListSlots(ctx context.Context, filter SlotFilter) ([]*ColloquioSlot, error)
	CancelSlot(ctx context.Context, slotID string) error

	CreateBooking(ctx context.Context, booking *ColloquioBooking) error
	GetBookingByID(ctx context.Context, id string) (*ColloquioBooking, error)
	ListUserBookings(ctx context.Context, userID string) ([]*ColloquioBooking, error)
	ListSlotBookings(ctx context.Context, slotID string) ([]*ColloquioBooking, error)
	UpdateBookingStatus(ctx context.Context, bookingID string, status BookingStatus, changedBy string, reason string) error

	GetTeacherProfileID(ctx context.Context, userID string) (string, error)
	GetParentProfileID(ctx context.Context, userID string) (string, error)
	GetStudentProfileID(ctx context.Context, userID string) (string, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	query := `SELECT id FROM teachers WHERE user_id = $1::uuid`
	var id string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if err == sql.ErrNoRows {
		return userID, nil // Fallback to userID if no dedicated profile
	}
	return id, err
}

func (r *PostgresRepository) GetParentProfileID(ctx context.Context, userID string) (string, error) {
	query := `SELECT id FROM parents WHERE user_id = $1::uuid`
	var id string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if err == sql.ErrNoRows {
		return userID, nil // Fallback
	}
	return id, err
}

func (r *PostgresRepository) GetStudentProfileID(ctx context.Context, userID string) (string, error) {
	query := `SELECT id FROM students WHERE user_id = $1::uuid`
	var id string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if err == sql.ErrNoRows {
		return userID, nil
	}
	return id, err
}

func (r *PostgresRepository) CreateSlot(ctx context.Context, slot *ColloquioSlot) error {
	slot.ID = uuid.New().String()
	slot.CreatedAt = time.Now()
	if slot.MaxBookings <= 0 {
		slot.MaxBookings = 1
	}
	if slot.Type == "" {
		slot.Type = TypeIndividual
	}

	query := `
		INSERT INTO colloquio_slots (id, teacher_id, school_id, date, start_time, end_time, max_bookings, booking_count, type, location, is_cancelled, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 0, $8, $9, false, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		slot.ID, slot.TeacherID, slot.SchoolID, slot.Date, slot.StartTime, slot.EndTime,
		slot.MaxBookings, slot.Type, slot.Location, slot.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetSlotByID(ctx context.Context, id string) (*ColloquioSlot, error) {
	query := `
		SELECT s.id, s.teacher_id, s.school_id, s.date, s.start_time, s.end_time,
		       s.max_bookings, s.booking_count, s.type, COALESCE(s.location, ''), s.is_cancelled, s.created_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name
		FROM colloquio_slots s
		LEFT JOIN teachers t ON s.teacher_id = t.id
		LEFT JOIN users u ON t.user_id = u.id OR s.teacher_id = u.id
		WHERE s.id = $1::uuid
	`
	s := &ColloquioSlot{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.TeacherID, &s.SchoolID, &s.Date, &s.StartTime, &s.EndTime,
		&s.MaxBookings, &s.BookingCount, &s.Type, &s.Location, &s.IsCancelled, &s.CreatedAt,
		&s.TeacherName,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *PostgresRepository) ListSlots(ctx context.Context, filter SlotFilter) ([]*ColloquioSlot, error) {
	query := `
		SELECT s.id, s.teacher_id, s.school_id, s.date, s.start_time, s.end_time,
		       s.max_bookings, s.booking_count, s.type, COALESCE(s.location, ''), s.is_cancelled, s.created_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name
		FROM colloquio_slots s
		LEFT JOIN teachers t ON s.teacher_id = t.id
		LEFT JOIN users u ON t.user_id = u.id OR s.teacher_id = u.id
		WHERE s.school_id = $1::uuid
		  AND s.is_cancelled = false
		  AND ($2 = '' OR s.teacher_id = $2::uuid OR t.user_id = $2::uuid)
		  AND ($3 = false OR s.booking_count < s.max_bookings)
		  AND s.date >= $4 AND s.date <= $5
		ORDER BY s.date ASC, s.start_time ASC
	`
	rows, err := r.db.QueryContext(ctx, query, filter.SchoolID, filter.TeacherID, filter.Available, filter.From, filter.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*ColloquioSlot
	for rows.Next() {
		s := &ColloquioSlot{}
		if err := rows.Scan(
			&s.ID, &s.TeacherID, &s.SchoolID, &s.Date, &s.StartTime, &s.EndTime,
			&s.MaxBookings, &s.BookingCount, &s.Type, &s.Location, &s.IsCancelled, &s.CreatedAt,
			&s.TeacherName,
		); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, rows.Err()
}

func (r *PostgresRepository) CancelSlot(ctx context.Context, slotID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `UPDATE colloquio_slots SET is_cancelled = true WHERE id = $1::uuid`, slotID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE colloquio_bookings SET status = 'Cancelled' WHERE slot_id = $1::uuid`, slotID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) CreateBooking(ctx context.Context, b *ColloquioBooking) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var isCancelled bool
	var maxB, countB int
	err = tx.QueryRowContext(ctx, `SELECT max_bookings, booking_count, is_cancelled FROM colloquio_slots WHERE id = $1::uuid FOR UPDATE`, b.SlotID).Scan(&maxB, &countB, &isCancelled)
	if err != nil {
		return fmt.Errorf("colloquio slot not found: %w", err)
	}
	if isCancelled {
		return ErrSlotCancelled
	}
	if countB >= maxB {
		return ErrSlotFull
	}

	b.ID = uuid.New().String()
	b.BookedAt = time.Now()
	b.Status = StatusConfirmed

	query := `
		INSERT INTO colloquio_bookings (id, slot_id, parent_id, student_id, status, notes, booked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.ExecContext(ctx, query, b.ID, b.SlotID, b.ParentID, b.StudentID, b.Status, b.Notes, b.BookedAt)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	_, err = tx.ExecContext(ctx, `UPDATE colloquio_slots SET booking_count = booking_count + 1 WHERE id = $1::uuid`, b.SlotID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetBookingByID(ctx context.Context, id string) (*ColloquioBooking, error) {
	query := `
		SELECT b.id, b.slot_id, b.parent_id, b.student_id, b.status, COALESCE(b.notes, ''), b.booked_at,
		       COALESCE(pu.first_name || ' ' || pu.last_name, '') AS parent_name,
		       COALESCE(su.first_name || ' ' || su.last_name, '') AS student_name
		FROM colloquio_bookings b
		LEFT JOIN parents p ON b.parent_id = p.id
		LEFT JOIN users pu ON p.user_id = pu.id OR b.parent_id = pu.id
		LEFT JOIN students s ON b.student_id = s.id
		LEFT JOIN users su ON s.user_id = su.id OR b.student_id = su.id
		WHERE b.id = $1::uuid
	`
	b := &ColloquioBooking{}
	var parentID, studentID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.SlotID, &parentID, &studentID, &b.Status, &b.Notes, &b.BookedAt,
		&b.ParentName, &b.StudentName,
	)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		b.ParentID = &parentID.String
	}
	if studentID.Valid {
		b.StudentID = &studentID.String
	}

	slot, err := r.GetSlotByID(ctx, b.SlotID)
	if err == nil {
		b.Slot = slot
	}
	return b, nil
}

func (r *PostgresRepository) ListUserBookings(ctx context.Context, userID string) ([]*ColloquioBooking, error) {
	query := `
		SELECT b.id, b.slot_id, b.parent_id, b.student_id, b.status, COALESCE(b.notes, ''), b.booked_at,
		       COALESCE(pu.first_name || ' ' || pu.last_name, '') AS parent_name,
		       COALESCE(su.first_name || ' ' || su.last_name, '') AS student_name
		FROM colloquio_bookings b
		LEFT JOIN parents p ON b.parent_id = p.id
		LEFT JOIN users pu ON p.user_id = pu.id OR b.parent_id = pu.id
		LEFT JOIN students s ON b.student_id = s.id
		LEFT JOIN users su ON s.user_id = su.id OR b.student_id = su.id
		WHERE b.parent_id = $1::uuid OR p.user_id = $1::uuid OR b.student_id = $1::uuid OR s.user_id = $1::uuid
		ORDER BY b.booked_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*ColloquioBooking
	for rows.Next() {
		b := &ColloquioBooking{}
		var parentID, studentID sql.NullString
		if err := rows.Scan(
			&b.ID, &b.SlotID, &parentID, &studentID, &b.Status, &b.Notes, &b.BookedAt,
			&b.ParentName, &b.StudentName,
		); err != nil {
			return nil, err
		}
		if parentID.Valid {
			b.ParentID = &parentID.String
		}
		if studentID.Valid {
			b.StudentID = &studentID.String
		}

		slot, err := r.GetSlotByID(ctx, b.SlotID)
		if err == nil {
			b.Slot = slot
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}

func (r *PostgresRepository) ListSlotBookings(ctx context.Context, slotID string) ([]*ColloquioBooking, error) {
	query := `
		SELECT b.id, b.slot_id, b.parent_id, b.student_id, b.status, COALESCE(b.notes, ''), b.booked_at,
		       COALESCE(pu.first_name || ' ' || pu.last_name, '') AS parent_name,
		       COALESCE(su.first_name || ' ' || su.last_name, '') AS student_name
		FROM colloquio_bookings b
		LEFT JOIN parents p ON b.parent_id = p.id
		LEFT JOIN users pu ON p.user_id = pu.id OR b.parent_id = pu.id
		LEFT JOIN students s ON b.student_id = s.id
		LEFT JOIN users su ON s.user_id = su.id OR b.student_id = su.id
		WHERE b.slot_id = $1::uuid
		ORDER BY b.booked_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, slotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*ColloquioBooking
	for rows.Next() {
		b := &ColloquioBooking{}
		var parentID, studentID sql.NullString
		if err := rows.Scan(
			&b.ID, &b.SlotID, &parentID, &studentID, &b.Status, &b.Notes, &b.BookedAt,
			&b.ParentName, &b.StudentName,
		); err != nil {
			return nil, err
		}
		if parentID.Valid {
			b.ParentID = &parentID.String
		}
		if studentID.Valid {
			b.StudentID = &studentID.String
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}

func (r *PostgresRepository) UpdateBookingStatus(ctx context.Context, bookingID string, status BookingStatus, changedBy string, reason string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var slotID string
	var prevStatus string
	err = tx.QueryRowContext(ctx, `SELECT slot_id, status FROM colloquio_bookings WHERE id = $1::uuid FOR UPDATE`, bookingID).Scan(&slotID, &prevStatus)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE colloquio_bookings SET status = $1 WHERE id = $2::uuid`, status, bookingID)
	if err != nil {
		return err
	}

	// Decrement slot booking_count if cancelled
	if prevStatus != string(StatusCancelled) && status == StatusCancelled {
		_, err = tx.ExecContext(ctx, `UPDATE colloquio_slots SET booking_count = GREATEST(0, booking_count - 1) WHERE id = $1::uuid`, slotID)
		if err != nil {
			return err
		}
	}

	// Log history
	var changedByUUID sql.NullString
	if changedBy != "" {
		changedByUUID = sql.NullString{String: changedBy, Valid: true}
	}
	historyQuery := `
		INSERT INTO colloquio_history (booking_id, new_status, changed_by, changed_at, reason)
		VALUES ($1::uuid, $2, $3::uuid, NOW(), $4)
	`
	_, _ = tx.ExecContext(ctx, historyQuery, bookingID, string(status), changedByUUID, reason)

	return tx.Commit()
}
