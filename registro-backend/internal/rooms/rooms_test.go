package rooms

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockRoomsRepo struct {
	buildings []SchoolBuilding
	rooms     []BookableRoom
	bookings  []RoomBooking
}

func newMockRoomsRepo() *mockRoomsRepo {
	return &mockRoomsRepo{
		buildings: make([]SchoolBuilding, 0),
		rooms:     make([]BookableRoom, 0),
		bookings:  make([]RoomBooking, 0),
	}
}

func (m *mockRoomsRepo) CreateBuilding(ctx context.Context, b *SchoolBuilding) (*SchoolBuilding, error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	m.buildings = append(m.buildings, *b)
	return b, nil
}

func (m *mockRoomsRepo) GetBuildingByID(ctx context.Context, id string) (*SchoolBuilding, error) {
	for _, b := range m.buildings {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockRoomsRepo) ListBuildings(ctx context.Context, schoolID string) ([]SchoolBuilding, error) {
	var res []SchoolBuilding
	for _, b := range m.buildings {
		if b.SchoolID == schoolID {
			res = append(res, b)
		}
	}
	return res, nil
}

func (m *mockRoomsRepo) UpdateBuilding(ctx context.Context, b *SchoolBuilding) (*SchoolBuilding, error) {
	for i, existing := range m.buildings {
		if existing.ID == b.ID {
			m.buildings[i] = *b
			return b, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockRoomsRepo) DeleteBuilding(ctx context.Context, id string) error {
	var filtered []SchoolBuilding
	for _, b := range m.buildings {
		if b.ID != id {
			filtered = append(filtered, b)
		}
	}
	m.buildings = filtered
	return nil
}

func (m *mockRoomsRepo) CreateRoom(ctx context.Context, r *BookableRoom) (*BookableRoom, error) {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	m.rooms = append(m.rooms, *r)
	return r, nil
}

func (m *mockRoomsRepo) GetRoomByID(ctx context.Context, id string) (*BookableRoom, error) {
	for _, r := range m.rooms {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockRoomsRepo) ListRooms(ctx context.Context, schoolID string, buildingID *string, roomType *string, activeOnly bool) ([]BookableRoom, error) {
	var res []BookableRoom
	for _, r := range m.rooms {
		if r.SchoolID != schoolID {
			continue
		}
		if buildingID != nil && (r.BuildingID == nil || *r.BuildingID != *buildingID) {
			continue
		}
		if roomType != nil && r.RoomType != *roomType {
			continue
		}
		if activeOnly && !r.IsActive {
			continue
		}
		res = append(res, r)
	}
	return res, nil
}

func (m *mockRoomsRepo) UpdateRoom(ctx context.Context, r *BookableRoom) (*BookableRoom, error) {
	for i, existing := range m.rooms {
		if existing.ID == r.ID {
			m.rooms[i] = *r
			return r, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockRoomsRepo) DeleteRoom(ctx context.Context, id string) error {
	var filtered []BookableRoom
	for _, r := range m.rooms {
		if r.ID != id {
			filtered = append(filtered, r)
		}
	}
	m.rooms = filtered
	return nil
}

func (m *mockRoomsRepo) CheckSlotConflict(ctx context.Context, roomID string, date string, hourIndex int, excludeBookingID *string) (bool, error) {
	for _, b := range m.bookings {
		if b.RoomID == roomID && b.BookingDate == date && b.HourIndex == hourIndex && b.Status == BookingStatusConfirmed {
			if excludeBookingID == nil || b.ID != *excludeBookingID {
				return true, nil
			}
		}
	}
	return false, nil
}

func (m *mockRoomsRepo) CreateBooking(ctx context.Context, b *RoomBooking) (*RoomBooking, error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	m.bookings = append(m.bookings, *b)
	return b, nil
}

func (m *mockRoomsRepo) CreateBookingBatch(ctx context.Context, bookings []RoomBooking) ([]RoomBooking, error) {
	for i := range bookings {
		if bookings[i].ID == "" {
			bookings[i].ID = uuid.New().String()
		}
		m.bookings = append(m.bookings, bookings[i])
	}
	return bookings, nil
}

func (m *mockRoomsRepo) GetBookingByID(ctx context.Context, id string) (*RoomBooking, error) {
	for _, b := range m.bookings {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockRoomsRepo) ListBookings(ctx context.Context, filter BookingFilter) ([]RoomBooking, error) {
	var res []RoomBooking
	for _, b := range m.bookings {
		if b.SchoolID != filter.SchoolID {
			continue
		}
		if filter.RoomID != nil && b.RoomID != *filter.RoomID {
			continue
		}
		if filter.TeacherID != nil && b.TeacherID != *filter.TeacherID {
			continue
		}
		if filter.Status != nil && b.Status != *filter.Status {
			continue
		}
		res = append(res, b)
	}
	return res, nil
}

func (m *mockRoomsRepo) CancelBooking(ctx context.Context, id string) error {
	for i, b := range m.bookings {
		if b.ID == id {
			m.bookings[i].Status = BookingStatusCancelled
			return nil
		}
	}
	return ErrNotFound
}

func (m *mockRoomsRepo) CancelRecurringSeries(ctx context.Context, parentBookingID string) error {
	for i, b := range m.bookings {
		if b.ID == parentBookingID || (b.ParentBookingID != nil && *b.ParentBookingID == parentBookingID) {
			m.bookings[i].Status = BookingStatusCancelled
		}
	}
	return nil
}

func (m *mockRoomsRepo) GetRoomBookingsForDateRange(ctx context.Context, roomID string, fromDate, toDate string) ([]RoomBooking, error) {
	var res []RoomBooking
	for _, b := range m.bookings {
		if b.RoomID == roomID && b.BookingDate >= fromDate && b.BookingDate <= toDate && b.Status == BookingStatusConfirmed {
			res = append(res, b)
		}
	}
	return res, nil
}

func TestRoomsService(t *testing.T) {
	repo := newMockRoomsRepo()
	svc := NewService(repo)
	ctx := context.Background()

	schoolID := "school-1"
	teacherID := "teacher-1"

	// 1. Create Building (Plesso)
	b, err := svc.CreateBuilding(ctx, schoolID, CreateBuildingRequest{
		Name: "Sede Centrale",
	})
	if err != nil {
		t.Fatalf("unexpected error creating building: %v", err)
	}
	if b.Name != "Sede Centrale" {
		t.Errorf("expected Sede Centrale, got %s", b.Name)
	}

	// 2. Create Room in Building
	r, err := svc.CreateRoom(ctx, schoolID, CreateRoomRequest{
		BuildingID: &b.ID,
		Name:       "Lab Informatica 1",
		RoomType:   RoomTypeLabInformatica,
		Capacity:   25,
	})
	if err != nil {
		t.Fatalf("unexpected error creating room: %v", err)
	}
	if r.RoomType != RoomTypeLabInformatica {
		t.Errorf("expected %s, got %s", RoomTypeLabInformatica, r.RoomType)
	}

	// 3. Create Spot Booking
	req := CreateBookingRequest{
		RoomID:      r.ID,
		BookingDate: "2026-10-05",
		HourIndex:   2,
	}
	booking, _, err := svc.CreateBooking(ctx, schoolID, teacherID, req)
	if err != nil {
		t.Fatalf("unexpected error booking room: %v", err)
	}
	if booking.HourIndex != 2 {
		t.Errorf("expected hour 2, got %d", booking.HourIndex)
	}

	// 4. Test Slot Conflict (same room, same day, same hour)
	_, _, err = svc.CreateBooking(ctx, schoolID, "teacher-2", req)
	if err == nil {
		t.Fatalf("expected error on duplicate slot booking, got nil")
	}

	// 5. Test Recurring Booking (Weekly for 3 weeks)
	recReq := CreateBookingRequest{
		RoomID:            r.ID,
		BookingDate:       "2026-10-06",
		HourIndex:         3,
		IsRecurring:       true,
		RecurrencePattern: func(s string) *string { return &s }(RecurrenceWeekly),
		RecurringUntil:    func(s string) *string { return &s }("2026-10-20"), // 3 occurrences: 06, 13, 20
	}
	parent, occurrences, err := svc.CreateBooking(ctx, schoolID, teacherID, recReq)
	if err != nil {
		t.Fatalf("unexpected error creating recurring booking: %v", err)
	}
	if len(occurrences) != 3 {
		t.Fatalf("expected 3 occurrences, got %d", len(occurrences))
	}
	if parent.ParentBookingID != nil {
		t.Errorf("parent booking should have nil ParentBookingID")
	}
	if occurrences[1].ParentBookingID == nil || *occurrences[1].ParentBookingID != parent.ID {
		t.Errorf("child booking parent ID mismatch")
	}

	// 6. Test Cancel Series
	err = svc.CancelBooking(ctx, teacherID, "teacher", schoolID, occurrences[1].ID, true)
	if err != nil {
		t.Fatalf("unexpected error cancelling series: %v", err)
	}

	// Check that all occurrences are cancelled
	for _, occ := range occurrences {
		found, err := repo.GetBookingByID(ctx, occ.ID)
		if err != nil {
			t.Fatalf("error fetching booking: %v", err)
		}
		if found.Status != BookingStatusCancelled {
			t.Errorf("expected booking %s to be cancelled, got %s", occ.ID, found.Status)
		}
	}
}

func TestRoomsService_ValidationErrors(t *testing.T) {
	repo := newMockRoomsRepo()
	svc := NewService(repo)
	ctx := context.Background()

	schoolID := "school-1"
	teacherID := "teacher-1"

	// 1. Invalid Hour Index (< 1)
	_, _, err := svc.CreateBooking(ctx, schoolID, teacherID, CreateBookingRequest{
		RoomID:      "room-1",
		BookingDate: "2026-10-05",
		HourIndex:   0,
	})
	if err == nil {
		t.Errorf("expected error for hour index 0, got nil")
	}

	// 2. Invalid Date format
	_, _, err = svc.CreateBooking(ctx, schoolID, teacherID, CreateBookingRequest{
		RoomID:      "room-1",
		BookingDate: "not-a-date",
		HourIndex:   1,
	})
	if err == nil {
		t.Errorf("expected error for invalid date, got nil")
	}

	// 3. Recurring booking without recurring_until
	_, _, err = svc.CreateBooking(ctx, schoolID, teacherID, CreateBookingRequest{
		RoomID:      "room-1",
		BookingDate: "2026-10-05",
		HourIndex:   1,
		IsRecurring: true,
	})
	if err == nil {
		t.Errorf("expected error for recurring booking missing recurring_until, got nil")
	}

	// 4. Cancel booking unauthorized
	// First create a room and booking by teacher-1
	r, err := repo.CreateRoom(ctx, &BookableRoom{
		SchoolID: schoolID,
		Name:     "Aula 10",
		RoomType: RoomTypeClassroom,
		Capacity: 20,
		IsActive: true,
	})
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	b, _, err := svc.CreateBooking(ctx, schoolID, teacherID, CreateBookingRequest{
		RoomID:      r.ID,
		BookingDate: "2026-10-12",
		HourIndex:   3,
	})
	if err != nil {
		t.Fatalf("failed to create booking: %v", err)
	}

	// Another teacher tries to cancel it without secretary/admin role
	err = svc.CancelBooking(ctx, "teacher-2", "teacher", schoolID, b.ID, false)
	if err == nil {
		t.Errorf("expected error when teacher-2 attempts to cancel teacher-1 booking, got nil")
	}
}
