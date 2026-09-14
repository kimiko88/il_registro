package personnel_desk

import (
	"context"
	"fmt"
	"time"
)

type Service interface {
	CreateRequest(ctx context.Context, schoolID, applicantID string, input CreateDeskRequestInput) (*DeskRequest, error)
	GetRequest(ctx context.Context, schoolID, id, currentUserID, currentUserRole string) (*DeskRequest, error)
	ListRequests(ctx context.Context, schoolID, currentUserID, currentUserRole, filterStatus string) ([]*DeskRequest, error)
	SubmitRequest(ctx context.Context, schoolID, id, currentUserID string) (*DeskRequest, error)
	AAReview(ctx context.Context, schoolID, id, reviewerID string, input AAReviewInput) (*DeskRequest, error)
	DSGASign(ctx context.Context, schoolID, id, signerID string, input DSGASignInput) (*DeskRequest, error)
	DSApprove(ctx context.Context, schoolID, id, approverID string, input DSApproveInput) (*DeskRequest, error)
	DeleteRequest(ctx context.Context, schoolID, id, currentUserID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func isReviewerRole(role string) bool {
	switch role {
	case "assistente_amministrativo", "dsga", "principal", "vice_principal", "secretary", "admin", "superadmin", "collaboratore_ds":
		return true
	default:
		return false
	}
}

func (s *service) CreateRequest(ctx context.Context, schoolID, applicantID string, input CreateDeskRequestInput) (*DeskRequest, error) {
	initialStatus := StatusDraft
	if input.SubmitNow {
		initialStatus = StatusSubmitted
	}

	req := &DeskRequest{
		SchoolID:    schoolID,
		ApplicantID: applicantID,
		Category:    input.Category,
		SubCategory: input.SubCategory,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Days:        input.Days,
		Hours:       input.Hours,
		Description: input.Description,
		Attachments: input.Attachments,
		Status:      initialStatus,
	}

	if err := s.repo.Create(ctx, req); err != nil {
		return nil, fmt.Errorf("errore salvataggio richiesta: %w", err)
	}

	return s.repo.GetByID(ctx, schoolID, req.ID)
}

func (s *service) GetRequest(ctx context.Context, schoolID, id, currentUserID, currentUserRole string) (*DeskRequest, error) {
	req, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	if !isReviewerRole(currentUserRole) && req.ApplicantID != currentUserID {
		return nil, fmt.Errorf("accesso non autorizzato a questa richiesta")
	}

	return req, nil
}

func (s *service) ListRequests(ctx context.Context, schoolID, currentUserID, currentUserRole, filterStatus string) ([]*DeskRequest, error) {
	applicantFilter := ""
	if !isReviewerRole(currentUserRole) {
		applicantFilter = currentUserID
	}

	return s.repo.List(ctx, schoolID, applicantFilter, filterStatus)
}

func (s *service) SubmitRequest(ctx context.Context, schoolID, id, currentUserID string) (*DeskRequest, error) {
	req, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	if req.ApplicantID != currentUserID {
		return nil, fmt.Errorf("solo il richiedente può inviare la richiesta")
	}
	if req.Status != StatusDraft {
		return nil, fmt.Errorf("la richiesta è già stata inviata")
	}

	req.Status = StatusSubmitted
	if err := s.repo.Update(ctx, req); err != nil {
		return nil, fmt.Errorf("errore invio richiesta: %w", err)
	}

	return s.repo.GetByID(ctx, schoolID, id)
}

func (s *service) AAReview(ctx context.Context, schoolID, id, reviewerID string, input AAReviewInput) (*DeskRequest, error) {
	req, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	if req.Status != StatusSubmitted && req.Status != StatusAAReview {
		return nil, fmt.Errorf("la richiesta non si trova nello stato di istruttoria")
	}

	now := time.Now()
	req.AANote = input.Note
	req.AAReviewedBy = &reviewerID
	req.AAReviewedAt = &now

	if input.Approve {
		req.Status = StatusDSGAReview
	} else {
		req.Status = StatusRejected
	}

	if err := s.repo.Update(ctx, req); err != nil {
		return nil, fmt.Errorf("errore aggiornamento istruttoria AA: %w", err)
	}

	return s.repo.GetByID(ctx, schoolID, id)
}

func (s *service) DSGASign(ctx context.Context, schoolID, id, signerID string, input DSGASignInput) (*DeskRequest, error) {
	req, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	if req.Status != StatusDSGAReview {
		return nil, fmt.Errorf("la richiesta non è in attesa di visto DSGA")
	}

	now := time.Now()
	req.DSGANote = input.Note
	req.DSGASignedBy = &signerID
	req.DSGASignedAt = &now

	if input.Approve {
		req.Status = StatusDSReview
	} else {
		req.Status = StatusRejected
	}

	if err := s.repo.Update(ctx, req); err != nil {
		return nil, fmt.Errorf("errore firma visto DSGA: %w", err)
	}

	return s.repo.GetByID(ctx, schoolID, id)
}

func (s *service) DSApprove(ctx context.Context, schoolID, id, approverID string, input DSApproveInput) (*DeskRequest, error) {
	req, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	if req.Status != StatusDSReview {
		return nil, fmt.Errorf("la richiesta non è in attesa di approvazione DS")
	}

	now := time.Now()
	if input.DecreeNum != "" {
		req.DSDecreeNum = &input.DecreeNum
	}
	req.DSNote = input.Note
	req.DSApprovedBy = &approverID
	req.DSApprovedAt = &now

	if input.Approve {
		req.Status = StatusApproved
	} else {
		req.Status = StatusRejected
	}

	if err := s.repo.Update(ctx, req); err != nil {
		return nil, fmt.Errorf("errore approvazione DS: %w", err)
	}

	return s.repo.GetByID(ctx, schoolID, id)
}

func (s *service) DeleteRequest(ctx context.Context, schoolID, id, currentUserID string) error {
	return s.repo.Delete(ctx, schoolID, id, currentUserID)
}
