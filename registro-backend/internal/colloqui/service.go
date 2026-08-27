package colloqui

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUnauthorized    = errors.New("unauthorized action on colloquio")
	ErrInvalidDate     = errors.New("invalid date or time format")
	ErrPastDate        = errors.New("la data dello slot non può essere nel passato")
	ErrOverlappingSlot = errors.New("esiste già uno slot sovrapposto per questo docente in questa fascia oraria")
	ErrNotGuardian     = errors.New("il genitore non è tutore legale dello studente indicato")
	ErrBookingNotFound = errors.New("booking not found")
	ErrSlotNotFound    = errors.New("slot not found")
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
	CreateGeneralMeeting(ctx context.Context, schoolID string, req CreateGeneralParentMeetingRequest) (*GeneralParentMeeting, error)
	ListGeneralMeetings(ctx context.Context, schoolID string) ([]GeneralParentMeeting, error)
	GetGeneralMeeting(ctx context.Context, id string) (*GeneralParentMeeting, error)
	BookQueueTicket(ctx context.Context, parentID string, req BookQueueTicketRequest) (*GeneralMeetingQueueTicket, error)
	ListQueueTickets(ctx context.Context, meetingID, teacherID, parentID string) ([]GeneralMeetingQueueTicket, error)
	UpdateTicketStatus(ctx context.Context, id, status, notes string) error
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
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDate
	}

	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		loc = time.Local
	}
	todayStr := time.Now().In(loc).Format("2006-01-02")
	if req.Date < todayStr {
		return nil, ErrPastDate
	}

	// Fetch teacher profile
	teacherProfileID, err := s.repo.GetTeacherProfileID(ctx, teacherUserID)
	if err != nil {
		return nil, err
	}
	if teacherProfileID == "" {
		teacherProfileID = teacherUserID
	}

	// Overlapping check
	overlap, err := s.repo.ExistsOverlappingSlot(ctx, teacherProfileID, req.Date, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, ErrOverlappingSlot
	}

	maxBookings := req.MaxBookings
	if maxBookings <= 0 {
		maxBookings = 1
	}

	slotType := req.Type
	if slotType == "" {
		slotType = TypeIndividual
	}

	slot := &ColloquioSlot{
		TeacherID:   teacherProfileID,
		SchoolID:    schoolID,
		Date:        d,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		MaxBookings: maxBookings,
		Type:        slotType,
		Location:    req.Location,
	}

	if err := s.repo.CreateSlot(ctx, slot); err != nil {
		return nil, err
	}
	return slot, nil
}

func (s *serviceImpl) ListSlots(ctx context.Context, schoolID, teacherID string, from, to time.Time, availableOnly bool) ([]*ColloquioSlot, error) {
	if schoolID == "" && teacherID == "" {
		return nil, errors.New("parametro school_id o teacher_id obbligatorio per la ricerca degli slot")
	}

	if from.IsZero() {
		from = time.Now().Truncate(24 * time.Hour)
	}
	if to.IsZero() {
		to = from.AddDate(0, 2, 0)
	}

	return s.repo.ListSlots(ctx, SlotFilter{
		SchoolID:  schoolID,
		TeacherID: teacherID,
		From:      from,
		To:        to,
		Available: availableOnly,
	})
}

func (s *serviceImpl) CancelSlot(ctx context.Context, actorID, actorRole, actorSchoolID, slotID string) error {
	slot, err := s.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return err
	}
	if slot == nil {
		return ErrSlotNotFound
	}

	if actorRole == "admin" && (actorSchoolID == "" || slot.SchoolID != actorSchoolID) {
		return ErrUnauthorized
	}

	if actorRole != "admin" && actorRole != "superadmin" {
		teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, actorID)
		if slot.TeacherID != actorID && slot.TeacherID != teacherProfileID {
			return ErrUnauthorized
		}
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

func (s *serviceImpl) ListMyBookings(ctx context.Context, userID, schoolID string) ([]*ColloquioBooking, error) {
	parentProfileID, _ := s.repo.GetParentProfileID(ctx, userID)
	if parentProfileID == "" {
		parentProfileID = userID
	}
	return s.repo.ListUserBookings(ctx, parentProfileID, schoolID)
}

func (s *serviceImpl) ListSlotBookings(ctx context.Context, actorID, actorRole, slotID string) ([]*ColloquioBooking, error) {
	slot, err := s.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return nil, err
	}
	if slot == nil {
		return nil, ErrSlotNotFound
	}

	if actorRole != "admin" && actorRole != "superadmin" {
		teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, actorID)
		if slot.TeacherID != actorID && slot.TeacherID != teacherProfileID {
			return nil, ErrUnauthorized
		}
	}

	return s.repo.ListSlotBookings(ctx, slotID)
}

func (s *serviceImpl) UpdateBookingStatus(ctx context.Context, actorID, actorRole, bookingID string, req UpdateBookingStatusRequest) error {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if booking == nil {
		return ErrBookingNotFound
	}

	if actorRole != "admin" && actorRole != "superadmin" {
		teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, actorID)
		parentProfileID, _ := s.repo.GetParentProfileID(ctx, actorID)

		slot, _ := s.repo.GetSlotByID(ctx, booking.SlotID)
		isTeacher := slot != nil && (slot.TeacherID == actorID || slot.TeacherID == teacherProfileID)
		isParent := booking.ParentID != nil && (*booking.ParentID == actorID || *booking.ParentID == parentProfileID)

		if !isTeacher && !isParent {
			return ErrUnauthorized
		}
	}

	return s.repo.UpdateBookingStatus(ctx, bookingID, req.Status, actorID, req.Reason)
}

func (s *serviceImpl) PatchSlot(ctx context.Context, actorID, actorRole, actorSchoolID, slotID, startTime, endTime string) error {
	slot, err := s.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return err
	}
	if slot == nil {
		return ErrSlotNotFound
	}

	if actorRole != "admin" && actorRole != "superadmin" {
		teacherProfileID, _ := s.repo.GetTeacherProfileID(ctx, actorID)
		if slot.TeacherID != actorID && slot.TeacherID != teacherProfileID {
			return ErrUnauthorized
		}
	}

	return s.repo.PatchSlot(ctx, slotID, startTime, endTime)
}

func (s *serviceImpl) GetBookingByID(ctx context.Context, actorID, actorRole, bookingID string) (*ColloquioBooking, error) {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, ErrBookingNotFound
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

func (s *serviceImpl) CreateGeneralMeeting(ctx context.Context, schoolID string, req CreateGeneralParentMeetingRequest) (*GeneralParentMeeting, error) {
	m := &GeneralParentMeeting{
		SchoolID:            schoolID,
		Title:               req.Title,
		EventDate:           req.EventDate,
		StartTime:           req.StartTime,
		EndTime:             req.EndTime,
		SlotDurationMinutes: req.SlotDurationMinutes,
		LocationType:        req.LocationType,
		Status:              "open_for_booking",
	}
	if m.SlotDurationMinutes <= 0 {
		m.SlotDurationMinutes = 7
	}
	if err := s.repo.CreateGeneralMeeting(ctx, m, req.TeacherIDs); err != nil {
		return nil, err
	}
	return s.repo.GetGeneralMeeting(ctx, m.ID)
}

func (s *serviceImpl) ListGeneralMeetings(ctx context.Context, schoolID string) ([]GeneralParentMeeting, error) {
	return s.repo.ListGeneralMeetings(ctx, schoolID)
}

func (s *serviceImpl) GetGeneralMeeting(ctx context.Context, id string) (*GeneralParentMeeting, error) {
	m, err := s.repo.GetGeneralMeeting(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrSlotNotFound
	}
	return m, nil
}

func (s *serviceImpl) BookQueueTicket(ctx context.Context, parentID string, req BookQueueTicketRequest) (*GeneralMeetingQueueTicket, error) {
	t := &GeneralMeetingQueueTicket{
		MeetingID: req.MeetingID,
		TeacherID: req.TeacherID,
		ParentID:  parentID,
		StudentID: req.StudentID,
		Notes:     req.Notes,
	}
	if err := s.repo.BookQueueTicket(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *serviceImpl) ListQueueTickets(ctx context.Context, meetingID, teacherID, parentID string) ([]GeneralMeetingQueueTicket, error) {
	return s.repo.ListQueueTickets(ctx, meetingID, teacherID, parentID)
}

func (s *serviceImpl) UpdateTicketStatus(ctx context.Context, id, status, notes string) error {
	return s.repo.UpdateTicketStatus(ctx, id, status, notes)
}
