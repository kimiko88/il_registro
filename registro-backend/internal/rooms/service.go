package rooms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden: insufficient permissions")
	ErrNotFound        = errors.New("not found")
	ErrSlotConflict    = errors.New("slot already booked for this room")
	ErrInvalidDate     = errors.New("invalid date format (expected YYYY-MM-DD)")
	ErrPastBookingDate = errors.New("cannot book a room in the past")
	ErrRoomInactive    = errors.New("room is currently inactive")
	ErrNoBookingNeeded = errors.New("this room does not require booking")
)

type Service interface {
	// Buildings
	CreateBuilding(ctx context.Context, schoolID string, req CreateBuildingRequest) (*SchoolBuilding, error)
	GetBuilding(ctx context.Context, id string) (*SchoolBuilding, error)
	ListBuildings(ctx context.Context, schoolID string) ([]SchoolBuilding, error)
	UpdateBuilding(ctx context.Context, schoolID, id string, req UpdateBuildingRequest) (*SchoolBuilding, error)
	DeleteBuilding(ctx context.Context, id string) error

	// Rooms
	CreateRoom(ctx context.Context, schoolID string, req CreateRoomRequest) (*BookableRoom, error)
	GetRoom(ctx context.Context, id string) (*BookableRoom, error)
	ListRooms(ctx context.Context, schoolID string, buildingID *string, roomType *string, activeOnly bool) ([]BookableRoom, error)
	UpdateRoom(ctx context.Context, schoolID, id string, req UpdateRoomRequest) (*BookableRoom, error)
	DeleteRoom(ctx context.Context, id string) error

	// Availability & Bookings
	GetRoomAvailability(ctx context.Context, roomID, fromDate, toDate string) ([]DayAvailability, error)
	CreateBooking(ctx context.Context, schoolID, teacherID string, req CreateBookingRequest) (*RoomBooking, []RoomBooking, error)
	ListBookings(ctx context.Context, filter BookingFilter) ([]RoomBooking, error)
	CancelBooking(ctx context.Context, actorID, actorRole, schoolID, bookingID string, cancelSeries bool) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	if repo == nil {
		panic("rooms.NewService: repo must not be nil")
	}
	return &service{repo: repo}
}

// ----------------- Buildings -----------------

func (s *service) CreateBuilding(ctx context.Context, schoolID string, req CreateBuildingRequest) (*SchoolBuilding, error) {
	if schoolID == "" {
		return nil, errors.New("school_id is required")
	}
	b := &SchoolBuilding{
		SchoolID: schoolID,
		Name:     req.Name,
		Address:  req.Address,
		Notes:    req.Notes,
		IsActive: true,
	}
	return s.repo.CreateBuilding(ctx, b)
}

func (s *service) GetBuilding(ctx context.Context, id string) (*SchoolBuilding, error) {
	return s.repo.GetBuildingByID(ctx, id)
}

func (s *service) ListBuildings(ctx context.Context, schoolID string) ([]SchoolBuilding, error) {
	return s.repo.ListBuildings(ctx, schoolID)
}

func (s *service) UpdateBuilding(ctx context.Context, schoolID, id string, req UpdateBuildingRequest) (*SchoolBuilding, error) {
	b, err := s.repo.GetBuildingByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b.SchoolID != schoolID {
		return nil, ErrForbidden
	}
	if req.Name != nil {
		b.Name = *req.Name
	}
	if req.Address != nil {
		b.Address = req.Address
	}
	if req.Notes != nil {
		b.Notes = req.Notes
	}
	if req.IsActive != nil {
		b.IsActive = *req.IsActive
	}
	return s.repo.UpdateBuilding(ctx, b)
}

func (s *service) DeleteBuilding(ctx context.Context, id string) error {
	return s.repo.DeleteBuilding(ctx, id)
}

// ----------------- Rooms -----------------

func (s *service) CreateRoom(ctx context.Context, schoolID string, req CreateRoomRequest) (*BookableRoom, error) {
	if schoolID == "" {
		return nil, errors.New("school_id is required")
	}
	reqBooking := true
	if req.RequiresBooking != nil {
		reqBooking = *req.RequiresBooking
	}
	roomType := req.RoomType
	if roomType == "" {
		roomType = RoomTypeClassroom
	}
	capacity := req.Capacity
	if capacity <= 0 {
		capacity = 30
	}
	equipJSON, _ := json.Marshal(req.Equipment)
	if len(req.Equipment) == 0 {
		equipJSON = []byte("[]")
	}

	r := &BookableRoom{
		SchoolID:        schoolID,
		BuildingID:      req.BuildingID,
		Name:            req.Name,
		RoomType:        roomType,
		Capacity:        capacity,
		Equipment:       json.RawMessage(equipJSON),
		RequiresBooking: reqBooking,
		IsActive:        true,
		Notes:           req.Notes,
	}
	return s.repo.CreateRoom(ctx, r)
}

func (s *service) GetRoom(ctx context.Context, id string) (*BookableRoom, error) {
	return s.repo.GetRoomByID(ctx, id)
}

func (s *service) ListRooms(ctx context.Context, schoolID string, buildingID *string, roomType *string, activeOnly bool) ([]BookableRoom, error) {
	return s.repo.ListRooms(ctx, schoolID, buildingID, roomType, activeOnly)
}

func (s *service) UpdateRoom(ctx context.Context, schoolID, id string, req UpdateRoomRequest) (*BookableRoom, error) {
	r, err := s.repo.GetRoomByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if r.SchoolID != schoolID {
		return nil, ErrForbidden
	}
	if req.BuildingID != nil {
		r.BuildingID = req.BuildingID
	}
	if req.Name != nil {
		r.Name = *req.Name
	}
	if req.RoomType != nil {
		r.RoomType = *req.RoomType
	}
	if req.Capacity != nil {
		r.Capacity = *req.Capacity
	}
	if req.Equipment != nil {
		equipJSON, _ := json.Marshal(req.Equipment)
		r.Equipment = json.RawMessage(equipJSON)
	}
	if req.RequiresBooking != nil {
		r.RequiresBooking = *req.RequiresBooking
	}
	if req.IsActive != nil {
		r.IsActive = *req.IsActive
	}
	if req.Notes != nil {
		r.Notes = req.Notes
	}
	return s.repo.UpdateRoom(ctx, r)
}

func (s *service) DeleteRoom(ctx context.Context, id string) error {
	return s.repo.DeleteRoom(ctx, id)
}

// ----------------- Bookings & Availability -----------------

func (s *service) GetRoomAvailability(ctx context.Context, roomID, fromDate, toDate string) ([]DayAvailability, error) {
	start, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	end, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	if end.Before(start) {
		end = start
	}

	bookings, err := s.repo.GetRoomBookingsForDateRange(ctx, roomID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	// Map bookings by date -> hour_index
	bookingMap := make(map[string]map[int]RoomBooking)
	for _, b := range bookings {
		if bookingMap[b.BookingDate] == nil {
			bookingMap[b.BookingDate] = make(map[int]RoomBooking)
		}
		bookingMap[b.BookingDate][b.HourIndex] = b
	}

	var result []DayAvailability
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		daySlots := make([]RoomAvailabilitySlot, 0, 12)
		for h := 1; h <= 10; h++ {
			slot := RoomAvailabilitySlot{
				HourIndex:   h,
				IsAvailable: true,
			}
			if b, exists := bookingMap[dateStr][h]; exists {
				slot.IsAvailable = false
				cpy := b
				slot.Booking = &cpy
			}
			daySlots = append(daySlots, slot)
		}
		result = append(result, DayAvailability{
			Date:  dateStr,
			Slots: daySlots,
		})
	}

	return result, nil
}

func (s *service) CreateBooking(ctx context.Context, schoolID, teacherID string, req CreateBookingRequest) (*RoomBooking, []RoomBooking, error) {
	room, err := s.repo.GetRoomByID(ctx, req.RoomID)
	if err != nil {
		return nil, nil, fmt.Errorf("room not found: %w", err)
	}
	if !room.IsActive {
		return nil, nil, ErrRoomInactive
	}
	if room.SchoolID != schoolID {
		return nil, nil, ErrForbidden
	}

	startDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		return nil, nil, ErrInvalidDate
	}

	// Check spot booking or single instance
	if !req.IsRecurring {
		hasConflict, err := s.repo.CheckSlotConflict(ctx, req.RoomID, req.BookingDate, req.HourIndex, nil)
		if err != nil {
			return nil, nil, err
		}
		if hasConflict {
			return nil, nil, fmt.Errorf("%w for date %s at hour %d", ErrSlotConflict, req.BookingDate, req.HourIndex)
		}

		b := &RoomBooking{
			SchoolID:    schoolID,
			RoomID:      req.RoomID,
			TeacherID:   teacherID,
			ClassID:     req.ClassID,
			SubjectID:   req.SubjectID,
			BookingDate: req.BookingDate,
			HourIndex:   req.HourIndex,
			Status:      BookingStatusConfirmed,
			Notes:       req.Notes,
			IsRecurring: false,
		}
		created, err := s.repo.CreateBooking(ctx, b)
		if err != nil {
			return nil, nil, err
		}
		return created, nil, nil
	}

	// Recurring booking logic
	if req.RecurringUntil == nil || *req.RecurringUntil == "" {
		return nil, nil, errors.New("recurring_until date is required for recurring bookings")
	}
	endDate, err := time.Parse("2006-01-02", *req.RecurringUntil)
	if err != nil {
		return nil, nil, ErrInvalidDate
	}
	if endDate.Before(startDate) {
		return nil, nil, errors.New("recurring_until cannot be earlier than booking_date")
	}

	stepDays := 7
	pattern := RecurrenceWeekly
	if req.RecurrencePattern != nil && *req.RecurrencePattern == RecurrenceBiweekly {
		stepDays = 14
		pattern = RecurrenceBiweekly
	}

	// Generate target dates
	var targetDates []string
	for cur := startDate; !cur.After(endDate); cur = cur.AddDate(0, 0, stepDays) {
		targetDates = append(targetDates, cur.Format("2006-01-02"))
	}

	if len(targetDates) == 0 {
		return nil, nil, errors.New("no occurrence dates generated for recurring pattern")
	}

	// Validate slot conflicts for all target dates before creating
	for _, dateStr := range targetDates {
		conflict, err := s.repo.CheckSlotConflict(ctx, req.RoomID, dateStr, req.HourIndex, nil)
		if err != nil {
			return nil, nil, err
		}
		if conflict {
			return nil, nil, fmt.Errorf("conflitto per l'aula %s alla data %s ora %d: slot già occupato", room.Name, dateStr, req.HourIndex)
		}
	}

	parentID := uuid.New().String()
	patternStr := pattern
	untilStr := *req.RecurringUntil

	// Create parent booking (first occurrence)
	parentBooking := RoomBooking{
		ID:                parentID,
		SchoolID:          schoolID,
		RoomID:            req.RoomID,
		TeacherID:         teacherID,
		ClassID:           req.ClassID,
		SubjectID:         req.SubjectID,
		BookingDate:       targetDates[0],
		HourIndex:         req.HourIndex,
		Status:            BookingStatusConfirmed,
		Notes:             req.Notes,
		IsRecurring:       true,
		RecurrencePattern: &patternStr,
		RecurringUntil:    &untilStr,
		ParentBookingID:   nil,
	}

	var allBookings []RoomBooking
	allBookings = append(allBookings, parentBooking)

	// Create child occurrences
	for i := 1; i < len(targetDates); i++ {
		child := RoomBooking{
			ID:                uuid.New().String(),
			SchoolID:          schoolID,
			RoomID:            req.RoomID,
			TeacherID:         teacherID,
			ClassID:           req.ClassID,
			SubjectID:         req.SubjectID,
			BookingDate:       targetDates[i],
			HourIndex:         req.HourIndex,
			Status:            BookingStatusConfirmed,
			Notes:             req.Notes,
			IsRecurring:       true,
			RecurrencePattern: &patternStr,
			RecurringUntil:    &untilStr,
			ParentBookingID:   &parentID,
		}
		allBookings = append(allBookings, child)
	}

	saved, err := s.repo.CreateBookingBatch(ctx, allBookings)
	if err != nil {
		return nil, nil, err
	}

	return &saved[0], saved, nil
}

func (s *service) ListBookings(ctx context.Context, filter BookingFilter) ([]RoomBooking, error) {
	return s.repo.ListBookings(ctx, filter)
}

func (s *service) CancelBooking(ctx context.Context, actorID, actorRole, schoolID, bookingID string, cancelSeries bool) error {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if booking.SchoolID != schoolID {
		return ErrForbidden
	}

	// Permissions check: teacher can only cancel own bookings; secretary, admin, principal can cancel any
	isManager := actorRole == "secretary" || actorRole == "admin" || actorRole == "superadmin" ||
		actorRole == "principal" || actorRole == "vice_principal" || actorRole == "collaboratore_ds"

	if !isManager && booking.TeacherID != actorID {
		return ErrForbidden
	}

	if cancelSeries && booking.IsRecurring {
		parentID := booking.ID
		if booking.ParentBookingID != nil && *booking.ParentBookingID != "" {
			parentID = *booking.ParentBookingID
		}
		return s.repo.CancelRecurringSeries(ctx, parentID)
	}

	return s.repo.CancelBooking(ctx, bookingID)
}
