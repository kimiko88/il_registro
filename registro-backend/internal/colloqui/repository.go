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
	PatchSlot(ctx context.Context, slotID string, startTime, endTime string) error
	CancelSlot(ctx context.Context, slotID string) error

	CreateBooking(ctx context.Context, booking *ColloquioBooking) error
	GetBookingByID(ctx context.Context, id string) (*ColloquioBooking, error)
	ListUserBookings(ctx context.Context, userID, schoolID string) ([]*ColloquioBooking, error)
	ListSlotBookings(ctx context.Context, slotID string) ([]*ColloquioBooking, error)
	UpdateBookingStatus(ctx context.Context, bookingID string, status BookingStatus, changedBy string, reason string) error

	GetTeacherProfileID(ctx context.Context, userID string) (string, error)
	GetParentProfileID(ctx context.Context, userID string) (string, error)
	GetStudentProfileID(ctx context.Context, userID string) (string, error)

	// Bug 97/129: verifica che il genitore sia tutore legale dello studente
	IsGuardian(ctx context.Context, parentUserID, studentUserID string) (bool, error)
	// Bug 98: verifica slot sovrapposti per lo stesso docente
	ExistsOverlappingSlot(ctx context.Context, teacherID string, date, startTime, endTime string) (bool, error)

	// General Meetings & Virtual Queue
	CreateGeneralMeeting(ctx context.Context, m *GeneralParentMeeting, teacherIDs []string) error
	ListGeneralMeetings(ctx context.Context, schoolID string) ([]GeneralParentMeeting, error)
	GetGeneralMeeting(ctx context.Context, id string) (*GeneralParentMeeting, error)
	BookQueueTicket(ctx context.Context, t *GeneralMeetingQueueTicket) error
	ListQueueTickets(ctx context.Context, meetingID, teacherID, parentID string) ([]GeneralMeetingQueueTicket, error)
	UpdateTicketStatus(ctx context.Context, id, status, notes string) error
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
		WHERE ($1 = '' OR s.school_id = NULLIF($1, '')::uuid)
		  AND s.is_cancelled = false
		  AND ($2 = '' OR s.teacher_id = NULLIF($2, '')::uuid OR t.user_id = NULLIF($2, '')::uuid)
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

func (r *PostgresRepository) PatchSlot(ctx context.Context, slotID string, startTime, endTime string) error {
	query := `
		UPDATE colloquio_slots
		SET start_time = COALESCE(NULLIF($2, ''), start_time),
		    end_time = COALESCE(NULLIF($3, ''), end_time)
		WHERE id = $1::uuid AND booking_count = 0 AND NOT EXISTS (SELECT 1 FROM colloquio_bookings WHERE slot_id = $1::uuid AND status != 'Cancelled')
	`
	res, err := r.db.ExecContext(ctx, query, slotID, startTime, endTime)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return errors.New("impossibile modificare l'orario di uno slot con prenotazioni attive")
	}
	return nil
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

func (r *PostgresRepository) ListUserBookings(ctx context.Context, userID, schoolID string) ([]*ColloquioBooking, error) {
	query := `
		SELECT b.id, b.slot_id, b.parent_id, b.student_id, b.status, COALESCE(b.notes, ''), b.booked_at,
		       COALESCE(pu.first_name || ' ' || pu.last_name, '') AS parent_name,
		       COALESCE(su.first_name || ' ' || su.last_name, '') AS student_name
		FROM colloquio_bookings b
		JOIN colloquio_slots sl ON b.slot_id = sl.id
		LEFT JOIN parents p ON b.parent_id = p.id
		LEFT JOIN users pu ON p.user_id = pu.id OR b.parent_id = pu.id
		LEFT JOIN students s ON b.student_id = s.id
		LEFT JOIN users su ON s.user_id = su.id OR b.student_id = su.id
		WHERE (b.parent_id = $1::uuid OR p.user_id = $1::uuid OR b.student_id = $1::uuid OR s.user_id = $1::uuid)
		  AND ($2 = '' OR sl.school_id = $2::uuid)
		ORDER BY b.booked_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, schoolID)
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
	if _, err := tx.ExecContext(ctx, historyQuery, bookingID, string(status), changedByUUID, reason); err != nil {
		return err
	}

	return tx.Commit()
}



// IsGuardian verifica che parentUserID sia tutore legale di studentUserID tramite la
// tabella parent_students (Bug 97/129).
func (r *PostgresRepository) IsGuardian(ctx context.Context, parentUserID, studentUserID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM student_parents sp
			LEFT JOIN parents p ON sp.parent_id = p.id
			LEFT JOIN students s ON sp.student_id = s.id
			WHERE (sp.parent_id = $1::uuid OR p.user_id = $1::uuid OR p.id = $1::uuid)
			  AND (sp.student_id = $2::uuid OR s.user_id = $2::uuid OR s.id = $2::uuid)
		)
	`
	var ok bool
	err := r.db.QueryRowContext(ctx, query, parentUserID, studentUserID).Scan(&ok)
	return ok, err
}

// ExistsOverlappingSlot controlla se esiste già uno slot per il docente in quella data e fascia oraria (Bug 98).
func (r *PostgresRepository) ExistsOverlappingSlot(ctx context.Context, teacherID, date, startTime, endTime string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM colloquio_slots
			WHERE teacher_id = $1::uuid
			  AND date = $2::date
			  AND is_cancelled = false
			  AND start_time < $4
			  AND end_time > $3
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, teacherID, date, startTime, endTime).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) CreateGeneralMeeting(ctx context.Context, m *GeneralParentMeeting, teacherIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO general_parent_meetings (
			id, school_id, title, event_date, start_time, end_time,
			slot_duration_minutes, location_type, status, created_at, updated_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4::date, $5, $6,
			$7, $8, $9, NOW(), NOW()
		) RETURNING id
	`
	if m.SlotDurationMinutes <= 0 {
		m.SlotDurationMinutes = 7
	}
	if m.LocationType == "" {
		m.LocationType = "in_presenza"
	}
	if m.Status == "" {
		m.Status = "open_for_booking"
	}

	err = tx.QueryRowContext(ctx, query,
		m.ID, m.SchoolID, m.Title, m.EventDate, m.StartTime, m.EndTime,
		m.SlotDurationMinutes, m.LocationType, m.Status,
	).Scan(&m.ID)
	if err != nil {
		return err
	}

	for _, tID := range teacherIDs {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO general_meeting_teacher_slots (
				meeting_id, teacher_id, room_or_table, max_bookings, created_at
			) VALUES ($1, $2, 'Aula Magna', 30, NOW())
			ON CONFLICT (meeting_id, teacher_id) DO NOTHING
		`, m.ID, tID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) ListGeneralMeetings(ctx context.Context, schoolID string) ([]GeneralParentMeeting, error) {
	query := `
		SELECT id, school_id, title, event_date::text, start_time, end_time,
		       slot_duration_minutes, location_type, status, created_at, updated_at
		FROM general_parent_meetings
		WHERE school_id = $1
		ORDER BY event_date DESC, start_time ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []GeneralParentMeeting
	for rows.Next() {
		var m GeneralParentMeeting
		if err := rows.Scan(
			&m.ID, &m.SchoolID, &m.Title, &m.EventDate, &m.StartTime, &m.EndTime,
			&m.SlotDurationMinutes, &m.LocationType, &m.Status, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) GetGeneralMeeting(ctx context.Context, id string) (*GeneralParentMeeting, error) {
	query := `
		SELECT id, school_id, title, event_date::text, start_time, end_time,
		       slot_duration_minutes, location_type, status, created_at, updated_at
		FROM general_parent_meetings
		WHERE id = $1
	`
	var m GeneralParentMeeting
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.SchoolID, &m.Title, &m.EventDate, &m.StartTime, &m.EndTime,
		&m.SlotDurationMinutes, &m.LocationType, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Fetch teacher slots
	slotRows, err := r.db.QueryContext(ctx, `
		SELECT gmts.id, gmts.meeting_id, gmts.teacher_id,
		       TRIM(COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '')) as teacher_name,
		       COALESCE(sub.name, '') as subject_name,
		       gmts.room_or_table, gmts.meet_url, gmts.max_bookings,
		       (SELECT COUNT(*) FROM general_meeting_queue_tickets WHERE meeting_id = gmts.meeting_id AND teacher_id = gmts.teacher_id AND status != 'annullato') as booked_count
		FROM general_meeting_teacher_slots gmts
		JOIN teachers t ON gmts.teacher_id = t.id
		JOIN users u ON t.user_id = u.id
		LEFT JOIN class_subjects cs ON cs.teacher_id = t.id
		LEFT JOIN subjects sub ON cs.subject_id = sub.id
		WHERE gmts.meeting_id = $1
		ORDER BY u.last_name ASC
	`, id)
	if err == nil {
		defer slotRows.Close()
		for slotRows.Next() {
			var s GeneralMeetingTeacherSlot
			if err := slotRows.Scan(&s.ID, &s.MeetingID, &s.TeacherID, &s.TeacherName, &s.SubjectName, &s.RoomOrTable, &s.MeetURL, &s.MaxBookings, &s.BookedCount); err != nil {
				return nil, err
			}
			m.TeacherSlots = append(m.TeacherSlots, s)
		}
		if err := slotRows.Err(); err != nil {
			return nil, err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}


	return &m, nil
}

func (r *PostgresRepository) BookQueueTicket(ctx context.Context, t *GeneralMeetingQueueTicket) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get next progressive ticket number for this teacher & meeting
	var maxTicket sql.NullInt64
	_ = tx.QueryRowContext(ctx, `
		SELECT MAX(ticket_number) FROM general_meeting_queue_tickets 
		WHERE meeting_id = $1 AND teacher_id = $2
	`, t.MeetingID, t.TeacherID).Scan(&maxTicket)

	nextTicket := 1
	if maxTicket.Valid {
		nextTicket = int(maxTicket.Int64) + 1
	}
	t.TicketNumber = nextTicket

	// Fetch meeting slot duration and start time to calculate scheduled time
	var startTime string
	var slotDuration int
	_ = tx.QueryRowContext(ctx, `SELECT start_time, slot_duration_minutes FROM general_parent_meetings WHERE id = $1`, t.MeetingID).Scan(&startTime, &slotDuration)
	if slotDuration <= 0 {
		slotDuration = 7
	}

	// Calculate scheduled time (e.g. 15:00 + (ticket-1)*7 min)
	parsedStart, err := time.Parse("15:04", startTime)
	if err == nil {
		scheduled := parsedStart.Add(time.Duration((nextTicket-1)*slotDuration) * time.Minute)
		t.ScheduledTime = scheduled.Format("15:04")
	} else {
		t.ScheduledTime = startTime
	}

	query := `
		INSERT INTO general_meeting_queue_tickets (
			id, meeting_id, teacher_id, parent_id, student_id,
			ticket_number, scheduled_time, status, notes, created_at, updated_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4, $5,
			$6, $7, 'prenotato', $8, NOW(), NOW()
		) RETURNING id
	`
	err = tx.QueryRowContext(ctx, query,
		t.ID, t.MeetingID, t.TeacherID, t.ParentID, t.StudentID,
		t.TicketNumber, t.ScheduledTime, t.Notes,
	).Scan(&t.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) ListQueueTickets(ctx context.Context, meetingID, teacherID, parentID string) ([]GeneralMeetingQueueTicket, error) {
	query := `
		SELECT gt.id, gt.meeting_id, gpm.title as meeting_title,
		       gt.teacher_id, TRIM(COALESCE(tu.first_name, '') || ' ' || COALESCE(tu.last_name, '')) as teacher_name,
		       gt.parent_id, TRIM(COALESCE(pu.first_name, '') || ' ' || COALESCE(pu.last_name, '')) as parent_name,
		       gt.student_id, TRIM(COALESCE(su.first_name, '') || ' ' || COALESCE(su.last_name, '')) as student_name,
		       gt.ticket_number, gt.scheduled_time, gt.status, gt.notes,
		       COALESCE(gmts.room_or_table, 'Aula Magna') as room_or_table,
		       gt.called_at, gt.completed_at, gt.created_at, gt.updated_at
		FROM general_meeting_queue_tickets gt
		JOIN general_parent_meetings gpm ON gt.meeting_id = gpm.id
		JOIN teachers t ON gt.teacher_id = t.id
		JOIN users tu ON t.user_id = tu.id
		JOIN parents p ON gt.parent_id = p.id
		JOIN users pu ON p.user_id = pu.id
		JOIN students s ON gt.student_id = s.id
		JOIN users su ON s.user_id = su.id
		LEFT JOIN general_meeting_teacher_slots gmts ON gmts.meeting_id = gt.meeting_id AND gmts.teacher_id = gt.teacher_id
		WHERE gt.meeting_id = $1
		  AND ($2 = '' OR gt.teacher_id = NULLIF($2, '')::uuid)
		  AND ($3 = '' OR gt.parent_id = NULLIF($3, '')::uuid)
		ORDER BY gt.ticket_number ASC
	`
	rows, err := r.db.QueryContext(ctx, query, meetingID, teacherID, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []GeneralMeetingQueueTicket
	for rows.Next() {
		var t GeneralMeetingQueueTicket
		if err := rows.Scan(
			&t.ID, &t.MeetingID, &t.MeetingTitle,
			&t.TeacherID, &t.TeacherName,
			&t.ParentID, &t.ParentName,
			&t.StudentID, &t.StudentName,
			&t.TicketNumber, &t.ScheduledTime, &t.Status, &t.Notes,
			&t.RoomOrTable, &t.CalledAt, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) UpdateTicketStatus(ctx context.Context, id, status, notes string) error {
	query := `
		UPDATE general_meeting_queue_tickets
		SET status = $1, notes = COALESCE(NULLIF($2, ''), notes),
		    called_at = CASE WHEN $1 = 'chiamato' AND called_at IS NULL THEN NOW() ELSE called_at END,
		    completed_at = CASE WHEN $1 IN ('concluso', 'assente') AND completed_at IS NULL THEN NOW() ELSE completed_at END,
		    updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, status, notes, id)
	return err
}
