package colloqui

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUnauthorized      = errors.New("unauthorized action on colloquio")
	ErrInvalidDate       = errors.New("invalid date or time format")
	ErrPastDate          = errors.New("la data dello slot non può essere nel passato")
	ErrOverlappingSlot   = errors.New("esiste già uno slot sovrapposto per questo docente in questa fascia oraria")
	ErrNotGuardian       = errors.New("il genitore non è tutore legale dello studente indicato")
)

type Service interface {
	CreateSlot(ctx context.Context, teacherUserID, schoolID string, req CreateSlotRequest) (*ColloquioSlot, error)
	ListSlots(ctx context.Context, schoolID, teacherID string, from, to time.Time, availableOnly bool) ([]*ColloquioSlot, error)
	CancelSlot(ctx context.Context, actorID, actorRole, actorSchoolID, slotID string) error
	BookSlot(ctx context.Context, parentUserID, parentSchoolID string, req CreateBookingRequest) (*ColloquioBooking, error)
	ListMyBookings(ctx context.Context, userID, schoolID string) ([]*ColloquioBooking, error)
	ListSlotBookings(ctx context.Context, actorID, actorRole, slotID string) ([]*ColloquioBooking, error)
	UpdateBookingStatus(ctx context.Context, actorID, actorRole, bookingID string, req UpdateBookingStatusRequest) error
	PatchSlot(ctx context.Context, actorID, actorRole, actorSchoolID, slotID, startTime, endTime string) error
	GetBookingByID(ctx context.Context, actorID, actorRole, bookingID string) (*ColloquioBooking, error)
	CreateAssembly(ctx context.Context, actorRole, teacherID, schoolID string, req CreateAssemblyRequest) (*ColloquioSlot, error)
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	if repo == nil {
		panic("colloqui.NewService: repo must not be nil")
	}
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) CreateSlot(ctx context.Context, teacherUserID, schoolID string, req CreateSlotRequest) (*ColloquioSlot, error) {
	if schoolID == "" {
		return nil, fmt.Errorf("school_id required")
	}
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDate
	}

	// Bug 146: confronto data con timezone locale della scuola/server
	todayStr := time.Now().In(time.Local).Format("2006-01-02")
	if req.Date < todayStr {
		return nil, ErrPastDate
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

	// Bug 98: verifica slot sovrapposti per lo stesso docente
	if req.StartTime != "" && req.EndTime != "" {
		overlaps, err := s.repo.ExistsOverlappingSlot(ctx, teacherProfileID, req.Date, req.StartTime, req.EndTime)
		if err == nil && overlaps {
			return nil, ErrOverlappingSlot
		}
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

func (s *serviceImpl) ListSlots(ctx context.Context, schoolID, teacherID string, from, to time.Time, availableOnly bool) ([]*ColloquioSlot, error) {
	if from.IsZero() {
		from = time.Now().Truncate(24 * time.Hour)
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

func (s *serviceImpl) CancelSlot(ctx context.Context, actorID, actorRole, actorSchoolID, slotID string) error {
	slot, err := s.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return err
	}

	if actorRole == "admin" && actorSchoolID != "" && slot.SchoolID != actorSchoolID {
		return ErrUnauthorized
	}

	teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, actorID)
	if slot.TeacherID != actorID && slot.TeacherID != teacherProfileID && actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}

	return s.repo.CancelSlot(ctx, slotID)
}

func (s *serviceImpl) BookSlot(ctx context.Context, parentUserID, parentSchoolID string, req CreateBookingRequest) (*ColloquioBooking, error) {
	slot, err := s.repo.GetSlotByID(ctx, req.SlotID)
	if err != nil {
		return nil, fmt.Errorf("slot non trovato: %w", err)
	}

	if parentSchoolID != "" && slot.SchoolID != "" && slot.SchoolID != parentSchoolID {
		return nil, errors.New("unauthorized: non puoi prenotare uno slot di un'altra scuola")
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
		// Bug 97/129: verifica guardianship prima di associare lo studente
		isGuardian, err := s.repo.IsGuardian(ctx, parentUserID, *req.StudentID)
		if err != nil {
			return nil, fmt.Errorf("errore verifica tutela: %w", err)
		}
		if !isGuardian {
			return nil, ErrNotGuardian
		}

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

// ListMyBookings restituisce le prenotazioni dell'utente filtrate per scuola (Bug 100).
func (s *serviceImpl) ListMyBookings(ctx context.Context, userID, schoolID string) ([]*ColloquioBooking, error) {
	return s.repo.ListUserBookings(ctx, userID, schoolID)
}

func (s *serviceImpl) ListSlotBookings(ctx context.Context, actorID, actorRole, slotID string) ([]*ColloquioBooking, error) {
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

func (s *serviceImpl) UpdateBookingStatus(ctx context.Context, actorID, actorRole, bookingID string, req UpdateBookingStatusRequest) error {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}

	parentProfileID, _ := s.repo.GetParentProfileID(ctx, actorID)
	isOwner := (booking.ParentID != nil && (*booking.ParentID == actorID || *booking.ParentID == parentProfileID))

	if isOwner && req.Status != "cancelled" && req.Status != StatusCancelled {
		return errors.New("unauthorized: i genitori possono solo cancellare le proprie prenotazioni")
	}

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

func (s *serviceImpl) PatchSlot(ctx context.Context, actorID, actorRole, actorSchoolID, slotID, startTime, endTime string) error {
	slot, err := s.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return err
	}
	if slot.BookingCount > 0 {
		return errors.New("impossibile modificare l'orario di uno slot con prenotazioni attive")
	}
	if actorRole == "admin" && actorSchoolID != "" && slot.SchoolID != actorSchoolID {
		return ErrUnauthorized
	}
	teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, actorID)
	if slot.TeacherID != actorID && slot.TeacherID != teacherProfileID && actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	return s.repo.PatchSlot(ctx, slotID, startTime, endTime)
}

func (s *serviceImpl) GetBookingByID(ctx context.Context, actorID, actorRole, bookingID string) (*ColloquioBooking, error) {
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

func (s *serviceImpl) CreateAssembly(ctx context.Context, actorRole, teacherID, schoolID string, req CreateAssemblyRequest) (*ColloquioSlot, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "vice_principal" {
		return nil, ErrUnauthorized
	}
	maxBookings := req.MaxBookings
	if maxBookings <= 0 {
		maxBookings = 100
	}
	return s.CreateSlot(ctx, teacherID, schoolID, CreateSlotRequest{
		Date:        req.Date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		MaxBookings: maxBookings,
		Type:        TypeAssembly,
		Location:    req.Location,
	})
}

