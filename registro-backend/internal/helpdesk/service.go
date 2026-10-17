package helpdesk

import (
	"context"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSlot(ctx context.Context, slot *HelpDeskSlot) error {
	if slot.TeacherID == "" || slot.SubjectID == "" || slot.SchoolID == "" {
		return errors.New("docente, materia ed istituto sono obbligatori")
	}
	if slot.SlotDate == "" || slot.StartTime == "" || slot.EndTime == "" {
		return errors.New("data, ora inizio e ora fine sono obbligatorie")
	}
	if slot.MaxCapacity <= 0 {
		slot.MaxCapacity = 4
	}
	slot.Status = "open"
	return s.repo.CreateSlot(ctx, slot)
}

func (s *Service) ListSlots(ctx context.Context, schoolID, subjectID, date, status string) ([]HelpDeskSlot, error) {
	if schoolID == "" {
		return nil, errors.New("school_id obbligatorio")
	}
	return s.repo.ListSlots(ctx, schoolID, subjectID, date, status)
}

func (s *Service) BookSlot(ctx context.Context, slotID, studentID, topic string) (*HelpDeskBooking, error) {
	if slotID == "" || studentID == "" {
		return nil, errors.New("slot_id e student_id sono obbligatori")
	}
	booking := &HelpDeskBooking{
		SlotID:           slotID,
		StudentID:        studentID,
		TopicDescription: topic,
	}
	if err := s.repo.BookSlot(ctx, booking); err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *Service) ListBookings(ctx context.Context, slotID string) ([]HelpDeskBooking, error) {
	if slotID == "" {
		return nil, errors.New("slot_id obbligatorio")
	}
	return s.repo.ListBookings(ctx, slotID)
}

func (s *Service) MarkAttendance(ctx context.Context, bookingID, status string) error {
	if bookingID == "" {
		return errors.New("booking_id obbligatorio")
	}
	if status != "attended" && status != "absent" {
		return errors.New("stato presenza deve essere 'attended' o 'absent'")
	}
	return s.repo.MarkAttendance(ctx, bookingID, status)
}

func (s *Service) CompleteSlot(ctx context.Context, slotID string) error {
	if slotID == "" {
		return errors.New("slot_id obbligatorio")
	}
	return s.repo.CompleteSlot(ctx, slotID)
}

func (s *Service) GetFISAccountingReport(ctx context.Context, schoolID string) ([]FISAccountingSummary, error) {
	if schoolID == "" {
		return nil, errors.New("school_id obbligatorio")
	}
	return s.repo.GetFISReport(ctx, schoolID)
}
