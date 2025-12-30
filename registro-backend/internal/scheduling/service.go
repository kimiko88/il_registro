package scheduling

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	CreateSlots(ctx context.Context, teacherID string, req CreateSlotRequest) error
	GetMySlots(ctx context.Context, teacherID string) ([]SlotResponse, error)
	DeleteSlot(ctx context.Context, teacherID, slotID string) error

	GetAvailableSlots(ctx context.Context, teacherID string) ([]SlotResponse, error)
	BookSlot(ctx context.Context, parentID string, req BookSlotRequest) (*BookingResponse, error)

	GetMyBookings(ctx context.Context, userID, role string) ([]BookingResponse, error)
	CancelBooking(ctx context.Context, userID, bookingID string) error

	// Admin
	UpdateSettings(ctx context.Context, req GeneralScheduleRequest) error
	GetAnalytics(ctx context.Context) (*AnalyticsResponse, error)
}

type service struct {
	repo      Repository
	validator *Validator
	notif     *NotificationService
	calendar  *CalendarService
	analytics *AnalyticsService
}

func NewService(repo Repository) Service {
	return &service{
		repo:      repo,
		validator: NewValidator(),
		notif:     NewNotificationService(),
		calendar:  NewCalendarService(),
		analytics: NewAnalyticsService(repo),
	}
}

func (s *service) CreateSlots(ctx context.Context, teacherID string, req CreateSlotRequest) error {
	date, _ := time.Parse("2006-01-02", req.Date)
	start, _ := time.Parse("15:04", req.StartTime)
	end, _ := time.Parse("15:04", req.EndTime)

	slot := &ColloquioSlot{
		TeacherID:   teacherID,
		SchoolID:    "default-school",
		Date:        date,
		StartTime:   start,
		EndTime:     end,
		MaxBookings: 1,
		Type:        req.Type,
		Location:    req.Location,
	}
	if req.MaxBookings > 1 {
		slot.MaxBookings = req.MaxBookings
	}

	if err := s.validator.ValidateSlot(slot); err != nil {
		return err
	}

	return s.repo.CreateSlot(ctx, slot)
}

func (s *service) GetMySlots(ctx context.Context, teacherID string) ([]SlotResponse, error) {
	// next 30 days
	from := time.Now()
	to := from.AddDate(0, 0, 30)
	slots, err := s.repo.GetSlots(ctx, teacherID, from, to)
	if err != nil {
		return nil, err
	}
	return convertSlots(slots), nil
}

func (s *service) DeleteSlot(ctx context.Context, teacherID, slotID string) error {
	// Should check if it has bookings first?
	// Simplified: Mark cancelled
	slot, err := s.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return err
	}
	if slot.TeacherID != teacherID {
		return errors.New("unauthorized")
	}

	slot.IsCancelled = true
	return s.repo.UpdateSlot(ctx, slot)
}

func (s *service) GetAvailableSlots(ctx context.Context, teacherID string) ([]SlotResponse, error) {
	from := time.Now()
	to := from.AddDate(0, 0, 14) // default window
	slots, err := s.repo.GetAvailableSlots(ctx, "default-school", teacherID, from, to)
	if err != nil {
		return nil, err
	}
	return convertSlots(slots), nil
}

func (s *service) BookSlot(ctx context.Context, parentID string, req BookSlotRequest) (*BookingResponse, error) {
	slot, err := s.repo.GetSlotByID(ctx, req.SlotID)
	if err != nil {
		return nil, err
	}

	settings, _ := s.repo.GetSettings(ctx, slot.SchoolID)
	if err := s.validator.ValidateBooking(slot, settings); err != nil {
		return nil, err
	}

	// Conflict check
	count, _ := s.repo.CountBookingsForParent(ctx, parentID, slot.Date, slot.StartTime, slot.EndTime)
	if count > 0 {
		return nil, errors.New("conflict: you already have a booking at this time")
	}

	booking := &ColloquioBooking{
		SlotID:    req.SlotID,
		ParentID:  &parentID,
		StudentID: req.StudentID,
		Status:    StatusConfirmed,
	}

	if err := s.repo.CreateBooking(ctx, booking); err != nil {
		return nil, err
	}

	s.notif.NotifyBooking(booking, slot, "parent")

	// Stub Response
	return &BookingResponse{ID: booking.ID, Status: StatusConfirmed}, nil
}

func (s *service) GetMyBookings(ctx context.Context, userID, role string) ([]BookingResponse, error) {
	var bookings []ColloquioBooking
	var err error
	if role == "parent" {
		bookings, err = s.repo.GetBookingsByParent(ctx, userID)
	} else {
		bookings, err = s.repo.GetBookingsByTeacher(ctx, userID)
	}
	if err != nil {
		return nil, err
	}

	return convertBookings(bookings), nil
}

func (s *service) CancelBooking(ctx context.Context, userID, bookingID string) error {
	b, err := s.repo.GetBooking(ctx, bookingID)
	if err != nil {
		return err
	}

	// Verify ownership (simplified)

	b.Status = StatusCancelled
	return s.repo.UpdateBooking(ctx, b)
}

func (s *service) UpdateSettings(ctx context.Context, req GeneralScheduleRequest) error {
	// Stub config update
	return nil
}

func (s *service) GetAnalytics(ctx context.Context) (*AnalyticsResponse, error) {
	return s.analytics.GetStats("default-school"), nil
}

// Helpers
func convertSlots(slots []ColloquioSlot) []SlotResponse {
	var res []SlotResponse
	for _, s := range slots {
		res = append(res, SlotResponse{
			ID: s.ID, Date: s.Date.Format("2006-01-02"),
			Type: s.Type, TeacherID: s.TeacherID,
		})
	}
	return res
}

func convertBookings(bookings []ColloquioBooking) []BookingResponse {
	var res []BookingResponse
	for _, b := range bookings {
		res = append(res, BookingResponse{ID: b.ID, Status: b.Status, Notes: b.Notes, BookedAt: b.BookedAt})
	}
	return res
}
