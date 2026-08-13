package scheduling

import (
	"context"
	"errors"
	"fmt"
	"registro-backend/internal/teachers"
	"time"
)

type NotificationSender interface {
	NotifyBooking(booking *ColloquioBooking, slot *ColloquioSlot, role string)
}

type CalendarExporter interface {
	GenerateICS(slots []ColloquioSlot, bookings []ColloquioBooking) string
}

type AnalyticsTracker interface {
	GetStats(ctx context.Context, schoolID string) *AnalyticsResponse
}

type Service interface {
	CreateSlots(ctx context.Context, teacherID string, req CreateSlotRequest) error
	GetMySlots(ctx context.Context, teacherID string) ([]SlotResponse, error)
	DeleteSlot(ctx context.Context, teacherID, slotID string) error

	GetAvailableSlots(ctx context.Context, teacherID string) ([]SlotResponse, error)
	BookSlot(ctx context.Context, parentID string, req BookSlotRequest) (*BookingResponse, error)

	GetMyBookings(ctx context.Context, userID, role string) ([]BookingResponse, error)
	CancelBooking(ctx context.Context, userID, bookingID string) error

	// Admin
	UpdateSettings(ctx context.Context, schoolID string, req GeneralScheduleRequest) error
	GetAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error)
	ValidateSchedule(ctx context.Context, slots []Slot) ([]Conflict, error)
	GenerateSchedule(ctx context.Context, req GenerationRequest) ([]Slot, error)
}

type service struct {
	repo        Repository
	teacherRepo teachers.Repository
	validator   *Validator
	generator   *Generator
	notif       NotificationSender
	calendar    CalendarExporter
	analytics   AnalyticsTracker
}

func NewService(repo Repository, teacherRepo teachers.Repository, notif NotificationSender, calendar CalendarExporter, analytics AnalyticsTracker) Service {
	var n NotificationSender = notif
	if n == nil {
		n = NewNotificationService()
	}
	var c CalendarExporter = calendar
	if c == nil {
		c = NewCalendarService()
	}
	var a AnalyticsTracker = analytics
	if a == nil {
		a = NewAnalyticsService(repo)
	}
	return &service{
		repo:        repo,
		teacherRepo: teacherRepo,
		validator:   NewValidator(),
		generator:   NewGenerator(),
		notif:       n,
		calendar:    c,
		analytics:   a,
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
		today := time.Now().Truncate(24 * time.Hour)
		if firstDate.Before(today) {
			return fmt.Errorf("la data dello slot non può essere nel passato (%s)", dateStr)
		}

		targetDates := []time.Time{firstDate}
		if req.IsRecurring && req.RecurringUntil != "" {
			untilDate, err := time.Parse("2006-01-02", req.RecurringUntil)
			if err != nil {
				return fmt.Errorf("data ricorrenza non valida '%s': usa YYYY-MM-DD: %w", req.RecurringUntil, err)
			}
			maxUntil := firstDate.AddDate(1, 0, 0)
			if untilDate.After(maxUntil) {
				untilDate = maxUntil
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
				for current.Add(time.Duration(req.Duration)*time.Minute).Before(end) || current.Add(time.Duration(req.Duration)*time.Minute).Equal(end) {
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

	// Overlap check: verify new slots do not overlap with existing slots for teacher
	if len(allSlots) > 0 {
		minDate := allSlots[0].Date
		maxDate := allSlots[0].Date
		for _, sl := range allSlots {
			if sl.Date.Before(minDate) {
				minDate = sl.Date
			}
			if sl.Date.After(maxDate) {
				maxDate = sl.Date
			}
		}
		existingSlots, _ := s.repo.GetSlots(ctx, teacher.ID, minDate, maxDate)
		for _, newSlot := range allSlots {
			for _, ex := range existingSlots {
				if ex.Date.Equal(newSlot.Date) && !ex.IsCancelled {
					if newSlot.StartTime.Before(ex.EndTime) && newSlot.EndTime.After(ex.StartTime) {
						return fmt.Errorf("sovrapposizione oraria rilevata per il giorno %s (%s - %s)",
							newSlot.Date.Format("2006-01-02"),
							newSlot.StartTime.Format("15:04"),
							newSlot.EndTime.Format("15:04"),
						)
					}
				}
			}
		}
	}

	return s.repo.CreateSlotsBatch(ctx, allSlots)
}

func (s *service) GetMySlots(ctx context.Context, userID string) ([]SlotResponse, error) {
	teacher, err := s.teacherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// next 60 days (timezone-safe local midnight)
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
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
		if err != nil || teacher == nil {
			return nil, fmt.Errorf("could not resolve teacher profile for slot lookup: %w", err)
		}
		schoolID = teacher.SchoolID
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

	if req.StudentID != nil && *req.StudentID != "" {
		isGuardian, err := s.repo.IsGuardian(ctx, parentID, *req.StudentID)
		if err != nil || !isGuardian {
			return nil, errors.New("unauthorized: parent is not a guardian of this student")
		}
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
	booking, err := s.repo.GetBooking(ctx, bookingID)
	if err != nil {
		return err
	}

	parentProfileID, _ := s.repo.ResolveParentUserID(ctx, userID)
	isParent := booking.ParentID != nil && *booking.ParentID != "" && (*booking.ParentID == userID || *booking.ParentID == parentProfileID)
	isTeacher := false
	if booking.Slot != nil && booking.Slot.TeacherID != "" {
		teacher, err := s.teacherRepo.GetByUserID(ctx, userID)
		if err == nil && teacher != nil && teacher.ID == booking.Slot.TeacherID {
			isTeacher = true
		}
	}

	if !isParent && !isTeacher {
		return errors.New("unauthorized: cannot cancel this booking")
	}

	booking.Status = StatusCancelled
	return s.repo.UpdateBooking(ctx, booking)
}

func (s *service) UpdateSettings(ctx context.Context, schoolID string, req GeneralScheduleRequest) error {
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end date: %w", err)
	}

	bookingWindow := req.BookingWindowDays
	if bookingWindow <= 0 {
		bookingWindow = 14
	}
	bookingBuffer := req.BookingBufferHours
	if bookingBuffer <= 0 {
		bookingBuffer = 24
	}

	settings := &ColloquioSettings{
		SchoolID:           schoolID,
		BookingWindowDays:  bookingWindow,
		BookingBufferHours: bookingBuffer,
		GeneralWindowStart: &start,
		GeneralWindowEnd:   &end,
	}

	return s.repo.UpdateSettings(ctx, settings)
}

func (s *service) GetAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error) {
	if schoolID == "" {
		return nil, fmt.Errorf("school_id is required")
	}
	return s.analytics.GetStats(ctx, schoolID), nil
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
