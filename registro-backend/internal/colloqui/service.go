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

	teacherProfileID, err := s.repo.GetTeacherProfileID(ctx, teacherUserID)
	if err != nil {
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
	if schoolID == "" {
		return nil, fmt.Errorf("school_id required")
	}
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
	parentProfileID, err := s.repo.GetParentProfileID(ctx, parentUserID)
	if err != nil {
		parentProfileID = parentUserID
	}

	booking := &ColloquioBooking{
		SlotID:   req.SlotID,
		ParentID: &parentProfileID,
		Notes:    req.Notes,
	}
	if req.StudentID != nil && *req.StudentID != "" {
		stdProfileID, err := s.repo.GetStudentProfileID(ctx, *req.StudentID)
		if err != nil {
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

	if !isOwner && actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}

	return s.repo.UpdateBookingStatus(ctx, bookingID, req.Status, actorID, req.Reason)
}
