package colloqui

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized action on colloquio")
	ErrInvalidDate  = errors.New("invalid date or time format")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("colloqui.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) CreateSlot(ctx context.Context, teacherUserID, schoolID string, req CreateSlotRequest) (*ColloquioSlot, error) {
	if schoolID == "" {
		return nil, fmt.Errorf("school_id required")
	}
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDate
	}
	if req.StartTime != "" && req.EndTime != "" {
		st, err1 := time.Parse("15:04", req.StartTime)
		et, err2 := time.Parse("15:04", req.EndTime)
		if err1 == nil && err2 == nil && !et.After(st) {
			return nil, errors.New("l'orario di fine deve essere successivo all'orario di inizio")
		}
	}

	teacherProfileID, err := s.repo.GetTeacherProfileID(ctx, teacherUserID)
	if err != nil || teacherProfileID == "" {
		teacherProfileID = teacherUserID
	}

	slot := &ColloquioSlot{
		TeacherID:   teacherProfileID,
		SchoolID:    schoolID,
		Date:        d,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		MaxBookings: req.MaxBookings,
		Type:        req.Type,
		Location:    req.Location,
	}

	if err := s.repo.CreateSlot(ctx, slot); err != nil {
		return nil, err
	}
	return slot, nil
}

func (s *Service) ListSlots(ctx context.Context, schoolID, teacherID string, from, to time.Time, availableOnly bool) ([]*ColloquioSlot, error) {
	if from.IsZero() {
		from = time.Now().AddDate(0, 0, -1)
	}
	if to.IsZero() {
		to = from.AddDate(0, 2, 0)
	}

	filter := SlotFilter{
		TeacherID: teacherID,
		SchoolID:  schoolID,
		From:      from,
		To:        to,
		Available: availableOnly,
	}

	return s.repo.ListSlots(ctx, filter)
}

func (s *Service) CancelSlot(ctx context.Context, actorID, actorRole, slotID string) error {
	slot, err := s.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return err
	}

	teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, actorID)
	if slot.TeacherID != actorID && slot.TeacherID != teacherProfileID && actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}

	return s.repo.CancelSlot(ctx, slotID)
}

func (s *Service) BookSlot(ctx context.Context, parentUserID string, req CreateBookingRequest) (*ColloquioBooking, error) {
	slot, err := s.repo.GetSlotByID(ctx, req.SlotID)
	if err != nil {
		return nil, fmt.Errorf("slot non trovato: %w", err)
	}

	if slot.IsCancelled {
		return nil, errors.New("slot non disponibile per la prenotazione (cancellato)")
	}
	if slot.BookingCount >= slot.MaxBookings {
		return nil, errors.New("slot esaurito: capienza massima raggiunta")
	}

	parentProfileID, err := s.repo.GetParentProfileID(ctx, parentUserID)
	if err != nil || parentProfileID == "" {
		parentProfileID = parentUserID
	}

	booking := &ColloquioBooking{
		SlotID:   req.SlotID,
		ParentID: &parentProfileID,
		Notes:    req.Notes,
	}
	if req.StudentID != nil && *req.StudentID != "" {
		stdProfileID, err := s.repo.GetStudentProfileID(ctx, *req.StudentID)
		if err != nil || stdProfileID == "" {
			stdProfileID = *req.StudentID
		}
		booking.StudentID = &stdProfileID
	}

	if err := s.repo.CreateBooking(ctx, booking); err != nil {
		return nil, err
	}

	return s.repo.GetBookingByID(ctx, booking.ID)
}

func (s *Service) ListMyBookings(ctx context.Context, userID string) ([]*ColloquioBooking, error) {
	return s.repo.ListUserBookings(ctx, userID)
}

func (s *Service) ListSlotBookings(ctx context.Context, actorID, actorRole, slotID string) ([]*ColloquioBooking, error) {
	slot, err := s.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return nil, err
	}

	teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, actorID)
	if slot.TeacherID != actorID && slot.TeacherID != teacherProfileID && actorRole != "admin" && actorRole != "superadmin" {
		return nil, ErrUnauthorized
	}

	return s.repo.ListSlotBookings(ctx, slotID)
}

func (s *Service) UpdateBookingStatus(ctx context.Context, actorID, actorRole, bookingID string, req UpdateBookingStatusRequest) error {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}

	parentProfileID, _ := s.repo.GetParentProfileID(ctx, actorID)
	isOwner := (booking.ParentID != nil && (*booking.ParentID == actorID || *booking.ParentID == parentProfileID))

	if !isOwner {
		if actorRole != "admin" && actorRole != "superadmin" {
			if actorRole != "teacher" {
				return ErrUnauthorized
			}
			// For teachers: verify ownership of the underlying slot
			slot, err := s.repo.GetSlotByID(ctx, booking.SlotID)
			if err != nil {
				return err
			}
			teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, actorID)
			if slot.TeacherID != actorID && slot.TeacherID != teacherProfileID {
				return ErrUnauthorized
			}
		}
	}

	return s.repo.UpdateBookingStatus(ctx, bookingID, req.Status, actorID, req.Reason)
}

func (s *Service) PatchSlot(ctx context.Context, actorID, actorRole, slotID, startTime, endTime string) error {
	slot, err := s.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return err
	}
	if slot.BookingCount > 0 {
		return errors.New("impossibile modificare l'orario di uno slot con prenotazioni attive")
	}
	teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, actorID)
	if slot.TeacherID != actorID && slot.TeacherID != teacherProfileID && actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	return s.repo.PatchSlot(ctx, slotID, startTime, endTime)
}

func (s *Service) GetBookingByID(ctx context.Context, actorID, actorRole, bookingID string) (*ColloquioBooking, error) {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if actorRole == "admin" || actorRole == "superadmin" {
		return booking, nil
	}

	parentProfileID, _ := s.repo.GetParentProfileID(ctx, actorID)
	if booking.ParentID != nil && (*booking.ParentID == actorID || *booking.ParentID == parentProfileID) {
		return booking, nil
	}

	slot, err := s.repo.GetSlotByID(ctx, booking.SlotID)
	if err == nil {
		teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, actorID)
		if slot.TeacherID == actorID || slot.TeacherID == teacherProfileID {
			return booking, nil
		}
	}

	return nil, ErrUnauthorized
}

func (s *Service) CreateAssembly(ctx context.Context, teacherID, schoolID string, req CreateAssemblyRequest) (*ColloquioSlot, error) {
	return s.CreateSlot(ctx, teacherID, schoolID, CreateSlotRequest{
		Date:        req.Date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		MaxBookings: 100,
		Type:        TypeAssembly,
		Location:    req.Location,
	})
}
