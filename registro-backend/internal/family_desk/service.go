package family_desk

import (
	"context"
	"errors"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SubmitRequest(ctx context.Context, req *FamilyRequest) error {
	if req.StudentID == "" || req.ParentID == "" || req.SchoolID == "" {
		return errors.New("student_id, parent_id e school_id sono obbligatori")
	}
	if !IsValidRequestType(req.RequestType) {
		return fmt.Errorf("tipologia di istanza non valida: %s", req.RequestType)
	}
	if err := ValidateRequestFormData(req); err != nil {
		return err
	}

	req.Status = StatusSubmitted
	return s.repo.CreateRequest(ctx, req)
}

func (s *Service) ListRequests(ctx context.Context, schoolID, parentID, status string) ([]FamilyRequest, error) {
	if schoolID == "" {
		return nil, errors.New("school_id obbligatorio")
	}
	return s.repo.ListRequests(ctx, schoolID, parentID, status)
}

func (s *Service) GetRequest(ctx context.Context, id string) (*FamilyRequest, error) {
	if id == "" {
		return nil, errors.New("id istanza obbligatorio")
	}
	return s.repo.GetRequest(ctx, id)
}

func (s *Service) ReviewRequest(ctx context.Context, id, status, rejectionReason, protocolNumber, reviewerID string) error {
	if id == "" {
		return errors.New("id istanza obbligatorio")
	}
	if status != StatusInIstruttoria && status != StatusApproved && status != StatusRejected {
		return errors.New("stato revisione non valido")
	}

	req, err := s.repo.GetRequest(ctx, id)
	if err != nil {
		return errors.New("istanza non trovata")
	}

	// A permanent delegate and the approved request must be committed atomically.
	if status == StatusApproved && req.RequestType == TypeDelegaRitiro {
		delegate, err := BuildDelegateFromRequest(req)
		if err != nil {
			return err
		}
		atomicRepo, ok := s.repo.(interface {
			ApproveRequestWithDelegate(context.Context, string, string, string, string, string, *PermanentDelegate) error
		})
		if !ok {
			return errors.New("repository does not support atomic approval and delegate creation")
		}
		return atomicRepo.ApproveRequestWithDelegate(ctx, id, status, rejectionReason, protocolNumber, reviewerID, delegate)
	}
	return s.repo.UpdateRequestStatus(ctx, id, status, rejectionReason, protocolNumber, reviewerID)
}

func (s *Service) ListPermanentDelegates(ctx context.Context, studentID, schoolID string) ([]PermanentDelegate, error) {
	return s.repo.ListPermanentDelegates(ctx, studentID, schoolID)
}
