package rooms

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	// Buildings
	CreateBuilding(ctx context.Context, b *SchoolBuilding) (*SchoolBuilding, error)
	GetBuildingByID(ctx context.Context, id string) (*SchoolBuilding, error)
	ListBuildings(ctx context.Context, schoolID string) ([]SchoolBuilding, error)
	UpdateBuilding(ctx context.Context, b *SchoolBuilding) (*SchoolBuilding, error)
	DeleteBuilding(ctx context.Context, id string) error

	// Rooms
	CreateRoom(ctx context.Context, r *BookableRoom) (*BookableRoom, error)
	GetRoomByID(ctx context.Context, id string) (*BookableRoom, error)
	ListRooms(ctx context.Context, schoolID string, buildingID *string, roomType *string, activeOnly bool) ([]BookableRoom, error)
	UpdateRoom(ctx context.Context, r *BookableRoom) (*BookableRoom, error)
	DeleteRoom(ctx context.Context, id string) error

	// Bookings
	CreateBooking(ctx context.Context, b *RoomBooking) (*RoomBooking, error)
	CreateBookingBatch(ctx context.Context, bookings []RoomBooking) ([]RoomBooking, error)
	GetBookingByID(ctx context.Context, id string) (*RoomBooking, error)
	ListBookings(ctx context.Context, filter BookingFilter) ([]RoomBooking, error)
	CancelBooking(ctx context.Context, id string) error
	CancelRecurringSeries(ctx context.Context, parentBookingID string) error
	CheckSlotConflict(ctx context.Context, roomID string, date string, hourIndex int, excludeBookingID *string) (bool, error)
	GetRoomBookingsForDateRange(ctx context.Context, roomID string, fromDate, toDate string) ([]RoomBooking, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// ----------------- Buildings -----------------

func (r *PostgresRepository) CreateBuilding(ctx context.Context, b *SchoolBuilding) (*SchoolBuilding, error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	query := `
		INSERT INTO school_buildings (id, school_id, name, address, notes, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, b.ID, b.SchoolID, b.Name, b.Address, b.Notes, b.IsActive).
		Scan(&b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create building: %w", err)
	}
	return b, nil
}

func (r *PostgresRepository) GetBuildingByID(ctx context.Context, id string) (*SchoolBuilding, error) {
	query := `
		SELECT id, school_id, name, address, notes, is_active, created_at, updated_at
		FROM school_buildings
		WHERE id = $1
	`
	var b SchoolBuilding
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.SchoolID, &b.Name, &b.Address, &b.Notes, &b.IsActive, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *PostgresRepository) ListBuildings(ctx context.Context, schoolID string) ([]SchoolBuilding, error) {
	query := `
		SELECT id, school_id, name, address, notes, is_active, created_at, updated_at
		FROM school_buildings
		WHERE school_id = $1
		ORDER BY name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to list buildings: %w", err)
	}
	defer rows.Close()

	var buildings []SchoolBuilding
	for rows.Next() {
		var b SchoolBuilding
		if err := rows.Scan(&b.ID, &b.SchoolID, &b.Name, &b.Address, &b.Notes, &b.IsActive, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		buildings = append(buildings, b)
	}
	return buildings, rows.Err()
}

func (r *PostgresRepository) UpdateBuilding(ctx context.Context, b *SchoolBuilding) (*SchoolBuilding, error) {
	query := `
		UPDATE school_buildings
		SET name = $1, address = $2, notes = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5 AND school_id = $6
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(ctx, query, b.Name, b.Address, b.Notes, b.IsActive, b.ID, b.SchoolID).
		Scan(&b.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update building: %w", err)
	}
	return b, nil
}

func (r *PostgresRepository) DeleteBuilding(ctx context.Context, id string) error {
	query := `DELETE FROM school_buildings WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ----------------- Rooms -----------------

func (r *PostgresRepository) CreateRoom(ctx context.Context, room *BookableRoom) (*BookableRoom, error) {
	if room.ID == "" {
		room.ID = uuid.New().String()
	}
	if len(room.Equipment) == 0 {
		room.Equipment = json.RawMessage("[]")
	}
	query := `
		INSERT INTO bookable_rooms (id, school_id, building_id, name, room_type, capacity, equipment, requires_booking, is_active, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		room.ID, room.SchoolID, room.BuildingID, room.Name, room.RoomType, room.Capacity,
		string(room.Equipment), room.RequiresBooking, room.IsActive, room.Notes,
	).Scan(&room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}
	return room, nil
}

func (r *PostgresRepository) GetRoomByID(ctx context.Context, id string) (*BookableRoom, error) {
	query := `
		SELECT r.id, r.school_id, r.building_id, COALESCE(b.name, ''), r.name, r.room_type,
		       r.capacity, r.equipment, r.requires_booking, r.is_active, r.notes, r.created_at, r.updated_at
		FROM bookable_rooms r
		LEFT JOIN school_buildings b ON r.building_id = b.id
		WHERE r.id = $1
	`
	var rm BookableRoom
	var equipStr string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rm.ID, &rm.SchoolID, &rm.BuildingID, &rm.BuildingName, &rm.Name, &rm.RoomType,
		&rm.Capacity, &equipStr, &rm.RequiresBooking, &rm.IsActive, &rm.Notes, &rm.CreatedAt, &rm.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	rm.Equipment = json.RawMessage(equipStr)
	return &rm, nil
}

func (r *PostgresRepository) ListRooms(ctx context.Context, schoolID string, buildingID *string, roomType *string, activeOnly bool) ([]BookableRoom, error) {
	var conditions []string
	var args []interface{}
	idx := 1

	conditions = append(conditions, fmt.Sprintf("r.school_id = $%d", idx))
	args = append(args, schoolID)
	idx++

	if buildingID != nil && *buildingID != "" {
		conditions = append(conditions, fmt.Sprintf("r.building_id = $%d", idx))
		args = append(args, *buildingID)
		idx++
	}

	if roomType != nil && *roomType != "" {
		conditions = append(conditions, fmt.Sprintf("r.room_type = $%d", idx))
		args = append(args, *roomType)
	}

	if activeOnly {
		conditions = append(conditions, "r.is_active = TRUE")
	}

	query := fmt.Sprintf(`
		SELECT r.id, r.school_id, r.building_id, COALESCE(b.name, ''), r.name, r.room_type,
		       r.capacity, r.equipment, r.requires_booking, r.is_active, r.notes, r.created_at, r.updated_at
		FROM bookable_rooms r
		LEFT JOIN school_buildings b ON r.building_id = b.id
		WHERE %s
		ORDER BY b.name NULLS FIRST, r.name ASC
	`, strings.Join(conditions, " AND "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list rooms: %w", err)
	}
	defer rows.Close()

	var rooms []BookableRoom
	for rows.Next() {
		var rm BookableRoom
		var equipStr string
		if err := rows.Scan(
			&rm.ID, &rm.SchoolID, &rm.BuildingID, &rm.BuildingName, &rm.Name, &rm.RoomType,
			&rm.Capacity, &equipStr, &rm.RequiresBooking, &rm.IsActive, &rm.Notes, &rm.CreatedAt, &rm.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rm.Equipment = json.RawMessage(equipStr)
		rooms = append(rooms, rm)
	}
	return rooms, rows.Err()
}

func (r *PostgresRepository) UpdateRoom(ctx context.Context, rm *BookableRoom) (*BookableRoom, error) {
	if len(rm.Equipment) == 0 {
		rm.Equipment = json.RawMessage("[]")
	}
	query := `
		UPDATE bookable_rooms
		SET building_id = $1, name = $2, room_type = $3, capacity = $4, equipment = $5,
		    requires_booking = $6, is_active = $7, notes = $8, updated_at = NOW()
		WHERE id = $9 AND school_id = $10
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		rm.BuildingID, rm.Name, rm.RoomType, rm.Capacity, string(rm.Equipment),
		rm.RequiresBooking, rm.IsActive, rm.Notes, rm.ID, rm.SchoolID,
	).Scan(&rm.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update room: %w", err)
	}
	return rm, nil
}

func (r *PostgresRepository) DeleteRoom(ctx context.Context, id string) error {
	query := `DELETE FROM bookable_rooms WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ----------------- Bookings -----------------

func (r *PostgresRepository) CheckSlotConflict(ctx context.Context, roomID string, date string, hourIndex int, excludeBookingID *string) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM room_bookings
		WHERE room_id = $1
		  AND booking_date = $2
		  AND hour_index = $3
		  AND status = 'confirmed'
	`
	var count int
	var err error
	if excludeBookingID != nil && *excludeBookingID != "" {
		query += " AND id != $4"
		err = r.db.QueryRowContext(ctx, query, roomID, date, hourIndex, *excludeBookingID).Scan(&count)
	} else {
		err = r.db.QueryRowContext(ctx, query, roomID, date, hourIndex).Scan(&count)
	}
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresRepository) CreateBooking(ctx context.Context, b *RoomBooking) (*RoomBooking, error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	if b.Status == "" {
		b.Status = BookingStatusConfirmed
	}
	query := `
		INSERT INTO room_bookings (
			id, school_id, room_id, teacher_id, class_id, subject_id,
			booking_date, hour_index, status, notes,
			is_recurring, recurrence_pattern, recurring_until, parent_booking_id,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14,
			NOW(), NOW()
		)
		RETURNING created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		b.ID, b.SchoolID, b.RoomID, b.TeacherID, b.ClassID, b.SubjectID,
		b.BookingDate, b.HourIndex, b.Status, b.Notes,
		b.IsRecurring, b.RecurrencePattern, b.RecurringUntil, b.ParentBookingID,
	).Scan(&b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}
	return b, nil
}

func (r *PostgresRepository) CreateBookingBatch(ctx context.Context, bookings []RoomBooking) ([]RoomBooking, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO room_bookings (
			id, school_id, room_id, teacher_id, class_id, subject_id,
			booking_date, hour_index, status, notes,
			is_recurring, recurrence_pattern, recurring_until, parent_booking_id,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14,
			NOW(), NOW()
		)
		RETURNING created_at, updated_at
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	now := time.Now()
	for i := range bookings {
		b := &bookings[i]
		if b.ID == "" {
			b.ID = uuid.New().String()
		}
		if b.Status == "" {
			b.Status = BookingStatusConfirmed
		}
		err := stmt.QueryRowContext(ctx,
			b.ID, b.SchoolID, b.RoomID, b.TeacherID, b.ClassID, b.SubjectID,
			b.BookingDate, b.HourIndex, b.Status, b.Notes,
			b.IsRecurring, b.RecurrencePattern, b.RecurringUntil, b.ParentBookingID,
		).Scan(&b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to create booking batch item (date %s hour %d): %w", b.BookingDate, b.HourIndex, err)
		}
		b.CreatedAt = now
		b.UpdatedAt = now
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *PostgresRepository) GetBookingByID(ctx context.Context, id string) (*RoomBooking, error) {
	query := `
		SELECT b.id, b.school_id, b.room_id, r.name, r.building_id, COALESCE(sb.name, ''),
		       b.teacher_id, COALESCE(u.first_name || ' ' || u.last_name, ''),
		       b.class_id, COALESCE(c.name, ''),
		       b.subject_id, COALESCE(s.name, ''),
		       b.booking_date::text, b.hour_index, b.status, b.notes,
		       b.is_recurring, b.recurrence_pattern, b.recurring_until::text, b.parent_booking_id,
		       b.created_at, b.updated_at
		FROM room_bookings b
		JOIN bookable_rooms r ON b.room_id = r.id
		LEFT JOIN school_buildings sb ON r.building_id = sb.id
		LEFT JOIN users u ON b.teacher_id = u.id
		LEFT JOIN classes c ON b.class_id = c.id
		LEFT JOIN subjects s ON b.subject_id = s.id
		WHERE b.id = $1
	`
	var b RoomBooking
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.SchoolID, &b.RoomID, &b.RoomName, &b.BuildingID, &b.BuildingName,
		&b.TeacherID, &b.TeacherName,
		&b.ClassID, &b.ClassName,
		&b.SubjectID, &b.SubjectName,
		&b.BookingDate, &b.HourIndex, &b.Status, &b.Notes,
		&b.IsRecurring, &b.RecurrencePattern, &b.RecurringUntil, &b.ParentBookingID,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *PostgresRepository) ListBookings(ctx context.Context, filter BookingFilter) ([]RoomBooking, error) {
	var conditions []string
	var args []interface{}
	idx := 1

	conditions = append(conditions, fmt.Sprintf("b.school_id = $%d", idx))
	args = append(args, filter.SchoolID)
	idx++

	if filter.RoomID != nil && *filter.RoomID != "" {
		conditions = append(conditions, fmt.Sprintf("b.room_id = $%d", idx))
		args = append(args, *filter.RoomID)
		idx++
	}

	if filter.BuildingID != nil && *filter.BuildingID != "" {
		conditions = append(conditions, fmt.Sprintf("r.building_id = $%d", idx))
		args = append(args, *filter.BuildingID)
		idx++
	}

	if filter.TeacherID != nil && *filter.TeacherID != "" {
		conditions = append(conditions, fmt.Sprintf("b.teacher_id = $%d", idx))
		args = append(args, *filter.TeacherID)
		idx++
	}

	if filter.ClassID != nil && *filter.ClassID != "" {
		conditions = append(conditions, fmt.Sprintf("b.class_id = $%d", idx))
		args = append(args, *filter.ClassID)
		idx++
	}

	if filter.FromDate != nil && *filter.FromDate != "" {
		conditions = append(conditions, fmt.Sprintf("b.booking_date >= $%d", idx))
		args = append(args, *filter.FromDate)
		idx++
	}

	if filter.ToDate != nil && *filter.ToDate != "" {
		conditions = append(conditions, fmt.Sprintf("b.booking_date <= $%d", idx))
		args = append(args, *filter.ToDate)
		idx++
	}

	if filter.Status != nil && *filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("b.status = $%d", idx))
		args = append(args, *filter.Status)
		idx++
	}

	if filter.ParentBookingID != nil && *filter.ParentBookingID != "" {
		conditions = append(conditions, fmt.Sprintf("(b.parent_booking_id = $%d OR b.id = $%d)", idx, idx))
		args = append(args, *filter.ParentBookingID)
	}

	query := fmt.Sprintf(`
		SELECT b.id, b.school_id, b.room_id, r.name, r.building_id, COALESCE(sb.name, ''),
		       b.teacher_id, COALESCE(u.first_name || ' ' || u.last_name, ''),
		       b.class_id, COALESCE(c.name, ''),
		       b.subject_id, COALESCE(s.name, ''),
		       b.booking_date::text, b.hour_index, b.status, b.notes,
		       b.is_recurring, b.recurrence_pattern, b.recurring_until::text, b.parent_booking_id,
		       b.created_at, b.updated_at
		FROM room_bookings b
		JOIN bookable_rooms r ON b.room_id = r.id
		LEFT JOIN school_buildings sb ON r.building_id = sb.id
		LEFT JOIN users u ON b.teacher_id = u.id
		LEFT JOIN classes c ON b.class_id = c.id
		LEFT JOIN subjects s ON b.subject_id = s.id
		WHERE %s
		ORDER BY b.booking_date ASC, b.hour_index ASC
	`, strings.Join(conditions, " AND "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list bookings: %w", err)
	}
	defer rows.Close()

	var bookings []RoomBooking
	for rows.Next() {
		var b RoomBooking
		if err := rows.Scan(
			&b.ID, &b.SchoolID, &b.RoomID, &b.RoomName, &b.BuildingID, &b.BuildingName,
			&b.TeacherID, &b.TeacherName,
			&b.ClassID, &b.ClassName,
			&b.SubjectID, &b.SubjectName,
			&b.BookingDate, &b.HourIndex, &b.Status, &b.Notes,
			&b.IsRecurring, &b.RecurrencePattern, &b.RecurringUntil, &b.ParentBookingID,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}

func (r *PostgresRepository) CancelBooking(ctx context.Context, id string) error {
	query := `
		UPDATE room_bookings
		SET status = 'cancelled', updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresRepository) CancelRecurringSeries(ctx context.Context, parentBookingID string) error {
	query := `
		UPDATE room_bookings
		SET status = 'cancelled', updated_at = NOW()
		WHERE parent_booking_id = $1 OR id = $1
	`
	_, err := r.db.ExecContext(ctx, query, parentBookingID)
	return err
}

func (r *PostgresRepository) GetRoomBookingsForDateRange(ctx context.Context, roomID string, fromDate, toDate string) ([]RoomBooking, error) {
	query := `
		SELECT b.id, b.school_id, b.room_id, r.name, r.building_id, COALESCE(sb.name, ''),
		       b.teacher_id, COALESCE(u.first_name || ' ' || u.last_name, ''),
		       b.class_id, COALESCE(c.name, ''),
		       b.subject_id, COALESCE(s.name, ''),
		       b.booking_date::text, b.hour_index, b.status, b.notes,
		       b.is_recurring, b.recurrence_pattern, b.recurring_until::text, b.parent_booking_id,
		       b.created_at, b.updated_at
		FROM room_bookings b
		JOIN bookable_rooms r ON b.room_id = r.id
		LEFT JOIN school_buildings sb ON r.building_id = sb.id
		LEFT JOIN users u ON b.teacher_id = u.id
		LEFT JOIN classes c ON b.class_id = c.id
		LEFT JOIN subjects s ON b.subject_id = s.id
		WHERE b.room_id = $1
		  AND b.booking_date >= $2
		  AND b.booking_date <= $3
		  AND b.status = 'confirmed'
		ORDER BY b.booking_date ASC, b.hour_index ASC
	`
	rows, err := r.db.QueryContext(ctx, query, roomID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []RoomBooking
	for rows.Next() {
		var b RoomBooking
		if err := rows.Scan(
			&b.ID, &b.SchoolID, &b.RoomID, &b.RoomName, &b.BuildingID, &b.BuildingName,
			&b.TeacherID, &b.TeacherName,
			&b.ClassID, &b.ClassName,
			&b.SubjectID, &b.SubjectName,
			&b.BookingDate, &b.HourIndex, &b.Status, &b.Notes,
			&b.IsRecurring, &b.RecurrencePattern, &b.RecurringUntil, &b.ParentBookingID,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}
