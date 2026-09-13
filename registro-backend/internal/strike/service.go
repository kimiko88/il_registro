package strike

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrUnauthorized              = errors.New("operazione non autorizzata: permessi insufficienti")
	ErrNotFound                  = errors.New("comunicazione di sciopero non trovata")
	ErrDeclarationDeadlinePassed = errors.New("il termine utile per esprimere o modificare la scelta preventiva è scaduto")
	ErrInvalidIntention          = errors.New("intenzione non valida: deve essere 'participates', 'not_participates' o 'undecided'")
)

func IsAdminOrDSGA(role string) bool {
	switch strings.ToLower(role) {
	case "dsga", "principal", "vice_principal", "admin", "superadmin", "collaboratore_ds":
		return true
	default:
		return false
	}
}

func CanCreateStrikeNotice(role string) bool {
	switch strings.ToLower(role) {
	case "dsga", "principal", "vice_principal", "admin", "superadmin":
		return true
	default:
		return false
	}
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("strike.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func parseDateTime(val string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, val); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("formato data/ora non riconosciuto: %s", val)
}

func (s *Service) CreateNotice(ctx context.Context, actorID, actorRole, schoolID string, req CreateStrikeNoticeRequest) (*StrikeNotice, error) {
	if !CanCreateStrikeNotice(actorRole) {
		return nil, ErrUnauthorized
	}

	if schoolID == "" {
		schoolID = s.repo.ResolveSchoolID(ctx, actorID)
	}
	if schoolID == "" {
		return nil, errors.New("school_id obbligatorio")
	}

	if _, err := time.Parse("2006-01-02", req.StrikeDate); err != nil {
		return nil, fmt.Errorf("formato data sciopero non valido: atteso YYYY-MM-DD")
	}

	deadline, err := parseDateTime(req.DeclarationDeadline)
	if err != nil {
		return nil, fmt.Errorf("data/ora di scadenza dichiarazione non valida: %w", err)
	}

	notice := &StrikeNotice{
		SchoolID:            schoolID,
		Title:               strings.TrimSpace(req.Title),
		ProclaimedBy:        strings.TrimSpace(req.ProclaimedBy),
		StrikeDate:          req.StrikeDate,
		DeclarationDeadline: deadline,
		Content:             strings.TrimSpace(req.Content),
		CreatedBy:           actorID,
		IsPublished:         true,
	}

	if req.PublishToBacheca {
		commID, err := s.repo.CreateBachecaCommunication(ctx, schoolID, actorID, notice.Title, notice.Content, deadline)
		if err == nil && commID != "" {
			notice.CommunicationID = &commID
		}
	}

	if err := s.repo.CreateNotice(ctx, notice); err != nil {
		return nil, err
	}
	return notice, nil
}

func (s *Service) GetNotice(ctx context.Context, id, schoolID, userID string) (*StrikeNotice, error) {
	notice, err := s.repo.GetNoticeByID(ctx, id, schoolID)
	if err != nil {
		return nil, ErrNotFound
	}

	if userID != "" {
		decl, _ := s.repo.GetDeclaration(ctx, id, userID)
		notice.UserDeclaration = decl
	}

	return notice, nil
}

func (s *Service) ListNotices(ctx context.Context, schoolID, userID string) ([]*StrikeNotice, error) {
	list, err := s.repo.ListNotices(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	if userID != "" {
		for _, n := range list {
			decl, _ := s.repo.GetDeclaration(ctx, n.ID, userID)
			n.UserDeclaration = decl
		}
	}

	return list, nil
}

func (s *Service) DeleteNotice(ctx context.Context, id, schoolID, actorRole string) error {
	if !CanCreateStrikeNotice(actorRole) {
		return ErrUnauthorized
	}
	return s.repo.DeleteNotice(ctx, id, schoolID)
}

func (s *Service) SubmitDeclaration(ctx context.Context, noticeID, schoolID, userID, ipAddress string, req SubmitDeclarationRequest) (*StrikeDeclaration, error) {
	if !IsValidIntention(req.Intention) {
		return nil, ErrInvalidIntention
	}

	notice, err := s.repo.GetNoticeByID(ctx, noticeID, schoolID)
	if err != nil {
		return nil, ErrNotFound
	}

	if time.Now().After(notice.DeclarationDeadline) {
		return nil, ErrDeclarationDeadlinePassed
	}

	decl := &StrikeDeclaration{
		StrikeNoticeID: noticeID,
		UserID:         userID,
		Intention:      req.Intention,
		IPAddress:      ipAddress,
		Notes:          strings.TrimSpace(req.Notes),
	}

	if err := s.repo.UpsertDeclaration(ctx, decl); err != nil {
		return nil, err
	}

	return decl, nil
}

func (s *Service) GetNoticeSummary(ctx context.Context, noticeID, schoolID, actorRole string) (*StrikeNoticeSummaryResponse, error) {
	if !IsAdminOrDSGA(actorRole) {
		return nil, ErrUnauthorized
	}
	return s.repo.GetNoticeSummary(ctx, noticeID, schoolID)
}
