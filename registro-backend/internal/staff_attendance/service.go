package staff_attendance

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Service gestisce la logica di business per le presenze del personale
type Service struct {
	repo Repository
}

// NewService crea un nuovo service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetDailySummary restituisce il riepilogo presenze per una data (solo per ruoli autorizzati)
func (s *Service) GetDailySummary(ctx context.Context, actorRole, schoolID, date string) (*DailyStaffSummary, error) {
	if !CanReadAttendance(actorRole) {
		return nil, errors.New("accesso negato: permessi insufficienti per visualizzare le presenze del personale")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return s.repo.GetDailySummary(ctx, schoolID, date)
}

// ListByDate restituisce la lista nominativa del personale per una data
func (s *Service) ListByDate(ctx context.Context, actorRole, schoolID, date string) ([]StaffAttendance, error) {
	if !CanReadAttendance(actorRole) {
		return nil, errors.New("accesso negato: permessi insufficienti")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return s.repo.ListByDate(ctx, schoolID, date)
}

// RecordAttendance registra/aggiorna la presenza di un membro del personale
func (s *Service) RecordAttendance(ctx context.Context, actorRole, actorID, schoolID string, req UpsertStaffAttendanceRequest) (*StaffAttendance, error) {
	if !CanWriteAttendance(actorRole) {
		return nil, errors.New("accesso negato: solo il personale ATA autorizzato può registrare le presenze")
	}
	if req.UserID == "" {
		return nil, errors.New("user_id obbligatorio")
	}
	if req.Date == "" {
		return nil, errors.New("date obbligatorio")
	}
	// Valida il formato data
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		return nil, fmt.Errorf("formato data non valido (atteso YYYY-MM-DD): %w", err)
	}
	// Valida lo status
	validStatuses := map[AttendanceStatus]bool{
		StatusPresent: true, StatusAbsent: true, StatusLate: true,
		StatusMission: true, StatusPermit: true, StatusSickLeave: true, StatusOnStrike: true,
	}
	if !validStatuses[req.Status] {
		return nil, fmt.Errorf("stato non valido: %s", req.Status)
	}

	return s.repo.Upsert(ctx, schoolID, actorID, req)
}

// BulkRecordAttendance registra le presenze di più persone in un'unica operazione
func (s *Service) BulkRecordAttendance(ctx context.Context, actorRole, actorID, schoolID string, req BulkUpsertRequest) ([]StaffAttendance, error) {
	if !CanWriteAttendance(actorRole) {
		return nil, errors.New("accesso negato: solo il personale ATA autorizzato può registrare le presenze")
	}
	if req.Date == "" {
		return nil, errors.New("date obbligatorio")
	}
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		return nil, fmt.Errorf("formato data non valido: %w", err)
	}
	if len(req.Attendances) == 0 {
		return nil, errors.New("nessuna presenza da registrare")
	}
	return s.repo.BulkUpsert(ctx, schoolID, actorID, req)
}

// DeleteAttendance rimuove una presenza (solo admin/superadmin/dsga)
func (s *Service) DeleteAttendance(ctx context.Context, actorRole, schoolID, id string) error {
	allowedRoles := []string{"dsga", "principal", "admin", "superadmin"}
	allowed := false
	for _, r := range allowedRoles {
		if r == actorRole {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("accesso negato: solo DSGA, Dirigente e Admin possono eliminare presenze")
	}
	return s.repo.Delete(ctx, schoolID, id)
}

// RegisterBadgeSwipe gestisce la timbratura da badge hardware
// Questo endpoint è tipicamente chiamato dai terminali badge, senza autenticazione utente
func (s *Service) RegisterBadgeSwipe(ctx context.Context, schoolID string, req BadgeSwipeRequest) (*BadgeSwipe, error) {
	if req.BadgeCode == "" {
		return nil, errors.New("badge_code obbligatorio")
	}
	if req.DeviceID == "" {
		return nil, errors.New("device_id obbligatorio")
	}
	return s.repo.RegisterBadgeSwipe(ctx, schoolID, req)
}

// ProcessBadgeSwipes elabora le timbrature pendenti e aggiorna le presenze
func (s *Service) ProcessBadgeSwipes(ctx context.Context, actorRole, schoolID string) (int, error) {
	allowedRoles := []string{"dsga", "principal", "admin", "superadmin"}
	allowed := false
	for _, r := range allowedRoles {
		if r == actorRole {
			allowed = true
			break
		}
	}
	if !allowed {
		return 0, errors.New("accesso negato")
	}
	return s.repo.ProcessPendingSwipes(ctx, schoolID)
}

// AssignBadge associa un badge a un utente
func (s *Service) AssignBadge(ctx context.Context, actorRole, schoolID, userID, badgeCode, notes string) error {
	allowedRoles := []string{"dsga", "principal", "admin", "superadmin", "secretary"}
	allowed := false
	for _, r := range allowedRoles {
		if r == actorRole {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("accesso negato: solo DSGA, Dirigente o Segreteria possono gestire i badge")
	}
	return s.repo.AssignBadge(ctx, UserBadge{
		SchoolID:  schoolID,
		UserID:    userID,
		BadgeCode: badgeCode,
		Notes:     notes,
		IsActive:  true,
	})
}
