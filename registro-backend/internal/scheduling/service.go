package scheduling

import (
	"context"
	"errors"
	"fmt"
	"registro-backend/internal/teachers"
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
	ValidateSchedule(ctx context.Context, slots []Slot) ([]Conflict, error)
	GenerateSchedule(ctx context.Context, req GenerationRequest) ([]Slot, error)
}

type service struct {
	repo        Repository
	teacherRepo teachers.Repository
	validator   *Validator
	generator   *Generator
	notif       *NotificationService
	calendar    *CalendarService
	analytics   *AnalyticsService
}

func NewService(repo Repository, teacherRepo teachers.Repository) Service {
	return &service{
		repo:        repo,
		teacherRepo: teacherRepo,
		validator:   NewValidator(),
		generator:   NewGenerator(),
		notif:       NewNotificationService(),
		calendar:    NewCalendarService(),
		analytics:   NewAnalyticsService(repo),
	}
}

func (s *service) CreateSlots(ctx context.Context, userID string, req CreateSlotRequest) error {
	// 1. Resolve Teacher Profile
	teacher, err := s.teacherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to resolve teacher profile: %w", err)
	}

	start, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return fmt.Errorf("ora inizio non valida '%s': usa HH:MM: %w", req.StartTime, err)
	}
	end, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return fmt.Errorf("ora fine non valida '%s': usa HH:MM: %w", req.EndTime, err)
	}

	var allSlots []ColloquioSlot

	for _, dateStr := range req.Dates {
		firstDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return fmt.Errorf("data non valida '%s': usa YYYY-MM-DD: %w", dateStr, err)
		}
		
		targetDates := []time.Time{firstDate}
		if req.IsRecurring && req.RecurringUntil != "" {
			untilDate, err := time.Parse("2006-01-02", req.RecurringUntil)
			if err != nil {
				return fmt.Errorf("data ricorrenza non valida '%s': usa YYYY-MM-DD: %w", req.RecurringUntil, err)
			}
			current := firstDate.AddDate(0, 0, 7)
			for !current.After(untilDate) {
				targetDates = append(targetDates, current)
				current = current.AddDate(0, 0, 7)
			}
		}

		for _, date := range targetDates {
			if req.Duration > 0 {
				// Split range into slots
				current := start
				for current.Add(time.Duration(req.Duration) * time.Minute).Before(end) || current.Add(time.Duration(req.Duration)*time.Minute).Equal(end) {
					slot := ColloquioSlot{
						TeacherID:   teacher.ID,
						SchoolID:    teacher.SchoolID,
						Date:        date,
						StartTime:   current,
						EndTime:     current.Add(time.Duration(req.Duration) * time.Minute),
						MaxBookings: 1,
						Type:        req.Type,
						Location:    req.Location,
					}
					if req.MaxBookings > 1 {
						slot.MaxBookings = req.MaxBookings
					}
					allSlots = append(allSlots, slot)
					current = current.Add(time.Duration(req.Duration) * time.Minute)
				}
			} else {
				// Single slot for the whole range
				slot := ColloquioSlot{
					TeacherID:   teacher.ID,
					SchoolID:    teacher.SchoolID,
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
				allSlots = append(allSlots, slot)
			}
		}
	}

	if len(allSlots) == 0 {
		return errors.New("no slots generated")
	}

	return s.repo.CreateSlotsBatch(ctx, allSlots)
}

func (s *service) GetMySlots(ctx context.Context, userID string) ([]SlotResponse, error) {
	teacher, err := s.teacherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// next 60 days
	from := time.Now().Truncate(24 * time.Hour)
	to := from.AddDate(0, 0, 60)
	slots, err := s.repo.GetSlots(ctx, teacher.ID, from, to)
	if err != nil {
		return nil, err
	}
	return convertSlots(slots), nil
}

func (s *service) DeleteSlot(ctx context.Context, userID, slotID string) error {
	teacher, err := s.teacherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	slot, err := s.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return err
	}
	if slot.TeacherID != teacher.ID {
		return errors.New("unauthorized")
	}

	if slot.BookingCount > 0 {
		// Cannot delete if there are bookings, mark cancelled
		slot.IsCancelled = true
		return s.repo.UpdateSlot(ctx, slot)
	}

	return s.repo.DeleteSlot(ctx, slotID)
}

func (s *service) GetAvailableSlots(ctx context.Context, teacherID string) ([]SlotResponse, error) {
	schoolID := ""
	if s.teacherRepo != nil {
		teacher, err := s.teacherRepo.GetByUserID(ctx, teacherID)
		if err == nil && teacher != nil {
			schoolID = teacher.SchoolID
		}
	}
	from := time.Now()
	to := from.AddDate(0, 0, 14) // default window
	slots, err := s.repo.GetAvailableSlots(ctx, schoolID, teacherID, from, to)
	if err != nil {
		return nil, err
	}
	return convertSlots(slots), nil
}

func (s *service) BookSlot(ctx context.Context, parentID string, req BookSlotRequest) (*BookingResponse, error) {
	parentProfileID, err := s.repo.ResolveParentUserID(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve parent profile: %w", err)
	}

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
		ParentID:  &parentProfileID,
		StudentID: req.StudentID,
		Status:    StatusConfirmed,
	}

	if err := s.repo.CreateBooking(ctx, booking); err != nil {
		return nil, err
	}

	// Hydrate booking response details (parent name/student name)
	hydrated, err := s.repo.GetBooking(ctx, booking.ID)
	if err == nil && hydrated != nil {
		booking = hydrated
		booking.Slot = slot
	}

	s.notif.NotifyBooking(booking, slot, "parent")

	return &BookingResponse{
		ID:       booking.ID,
		Status:   booking.Status,
		BookedAt: booking.BookedAt,
		Notes:    booking.Notes,
		SlotInfo: SlotResponse{
			ID:        slot.ID,
			Date:      slot.Date.Format("2006-01-02"),
			TimeRange: fmt.Sprintf("%s - %s", slot.StartTime.Format("15:04"), slot.EndTime.Format("15:04")),
			Type:      slot.Type,
			Location:  slot.Location,
		},
		ParentName:  booking.ParentName,
		StudentName: booking.StudentName,
	}, nil
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

	// Try to resolve parent profile ID to match parent ownership check
	parentProfileID, _ := s.repo.ResolveParentUserID(ctx, userID)

	isParent := b.ParentID != nil && ((parentProfileID != "" && *b.ParentID == parentProfileID) || *b.ParentID == userID)

	if !isParent {
		// Check if it's the teacher (by userID or resolved teacher profile ID)
		slot, err := s.repo.GetSlotByID(ctx, b.SlotID)
		if err != nil {
			return err
		}
		teacherProfileID := userID
		if s.teacherRepo != nil {
			if t, err := s.teacherRepo.GetByUserID(ctx, userID); err == nil && t != nil {
				teacherProfileID = t.ID
			}
		}
		if slot.TeacherID != userID && slot.TeacherID != teacherProfileID {
			return errors.New("unauthorized: you cannot cancel this booking")
		}
	}

	b.Status = StatusCancelled
	return s.repo.UpdateBooking(ctx, b)
}

func (s *service) UpdateSettings(ctx context.Context, req GeneralScheduleRequest) error {
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end date: %w", err)
	}

	settings := &ColloquioSettings{
		SchoolID:           "default-school", // Temporary default
		BookingWindowDays:  14,
		BookingBufferHours: 24,
		GeneralWindowStart: &start,
		GeneralWindowEnd:   &end,
	}

	return s.repo.UpdateSettings(ctx, settings)
}

func (s *service) GetAnalytics(ctx context.Context) (*AnalyticsResponse, error) {
	return s.analytics.GetStats("default-school"), nil
}

func (s *service) ValidateSchedule(ctx context.Context, slots []Slot) ([]Conflict, error) {
	// Simple wrapper around validator logic
	// In real app, might hydrate teacher names etc.
	return s.validator.ValidateSchedule(slots), nil
}

func (s *service) GenerateSchedule(ctx context.Context, req GenerationRequest) ([]Slot, error) {
	return s.generator.Generate(req)
}

// Helpers
func convertSlots(slots []ColloquioSlot) []SlotResponse {
	var res []SlotResponse
	for _, s := range slots {
		res = append(res, SlotResponse{
			ID: s.ID, Date: s.Date.Format("2006-01-02"),
			TimeRange: fmt.Sprintf("%s - %s", s.StartTime.Format("15:04"), s.EndTime.Format("15:04")),
			Type:      s.Type,
			Available: !s.IsCancelled && s.BookingCount < s.MaxBookings,
			TeacherID: s.TeacherID,
		})
	}
	return res
}

func convertBookings(bookings []ColloquioBooking) []BookingResponse {
	var res []BookingResponse
	for _, b := range bookings {
		br := BookingResponse{
			ID: b.ID, Status: b.Status, Notes: b.Notes, BookedAt: b.BookedAt,
			ParentName:  b.ParentName,
			StudentName: b.StudentName,
		}
		if b.Slot != nil {
			br.SlotInfo = SlotResponse{
				ID:        b.Slot.ID,
				Date:      b.Slot.Date.Format("2006-01-02"),
				TimeRange: fmt.Sprintf("%s - %s", b.Slot.StartTime.Format("15:04"), b.Slot.EndTime.Format("15:04")),
				Type:      b.Slot.Type,
				Location:  b.Slot.Location,
			}
		}
		res = append(res, br)
	}
	return res
}
