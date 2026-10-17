package parents

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DualSignatureRepository interface {
	CreateAuthorization(ctx context.Context, auth *DualParentalAuthorization) error
	GetAuthorizationByID(ctx context.Context, id string) (*DualParentalAuthorization, error)
	UpdateAuthorization(ctx context.Context, auth *DualParentalAuthorization) error
	ListAuthorizationsForParent(ctx context.Context, parentID string) ([]DualParentalAuthorization, error)
	GetCustodyInfo(ctx context.Context, parentID, studentID string) (*StudentCustodyInfo, error)
	ValidateParentPIN(ctx context.Context, parentID, pin string) (bool, error)
}

type DualSignatureService struct {
	repo           DualSignatureRepository
	mirrorNotifier *MirrorNotifier
}

func NewDualSignatureService(repo DualSignatureRepository) *DualSignatureService {
	return &DualSignatureService{
		repo:           repo,
		mirrorNotifier: NewMirrorNotifier(),
	}
}

func (s *DualSignatureService) CreateAuthorization(ctx context.Context, params CreateDualAuthParams) (*DualParentalAuthorization, error) {
	if params.StudentID == "" || params.DocumentRefID == "" || params.Title == "" || params.Parent1ID == "" {
		return nil, errors.New("missing mandatory fields for dual parental authorization")
	}

	authID := uuid.New().String()
	deadline := params.Deadline
	if deadline.IsZero() {
		deadline = time.Now().Add(14 * 24 * time.Hour) // default 14 days
	}

	record := &DualParentalAuthorization{
		ID:            authID,
		StudentID:     params.StudentID,
		DocumentType:  params.DocumentType,
		DocumentRefID: params.DocumentRefID,
		Title:         params.Title,
		Parent1ID:     params.Parent1ID,
		Parent2ID:     params.Parent2ID,
		Status:        StatusPendingFirst,
		Deadline:      &deadline,
		CreatedAt:     time.Now(),
	}

	if err := s.repo.CreateAuthorization(ctx, record); err != nil {
		return nil, fmt.Errorf("creating authorization: %w", err)
	}

	// Dispatch notification to Parent 1
	_ = s.mirrorNotifier.NotifyNewAuthorization(ctx, record, params.Parent1ID)

	return record, nil
}

func (s *DualSignatureService) SignDocument(ctx context.Context, authID, parentUserID, pin string) (*DualParentalAuthorization, error) {
	record, err := s.repo.GetAuthorizationByID(ctx, authID)
	if err != nil {
		return nil, errors.New("authorization not found")
	}

	if record.Status == StatusCompleted || record.Status == StatusRejected {
		return nil, fmt.Errorf("document already finalized with status: %s", record.Status)
	}

	// Validate PIN
	validPin, err := s.repo.ValidateParentPIN(ctx, parentUserID, pin)
	if err != nil || !validPin {
		return nil, errors.New("invalid signature PIN")
	}

	now := time.Now()
	isParent1 := record.Parent1ID == parentUserID
	isParent2 := record.Parent2ID == parentUserID

	if !isParent1 && !isParent2 {
		return nil, errors.New("unauthorized: user is not a designated parent for this authorization")
	}

	if isParent1 {
		record.Parent1SignedAt = &now
		record.Parent1PinVerified = true

		if record.Parent2ID == "" {
			// Sole custody: single signature is sufficient
			record.Status = StatusCompleted
		} else if record.Parent2SignedAt != nil {
			// Parent 2 already signed
			record.Status = StatusCompleted
		} else {
			// Shared custody: pending second parent
			record.Status = StatusPendingSecond
			// Notify Parent 2 that Parent 1 has signed and their co-signature is required
			_ = s.mirrorNotifier.NotifyPendingSecondSignature(ctx, record, record.Parent2ID, record.Parent1ID)
		}
	} else if isParent2 {
		record.Parent2SignedAt = &now
		record.Parent2PinVerified = true

		if record.Parent1SignedAt != nil {
			// Both signed
			record.Status = StatusCompleted
		} else {
			// Parent 2 signed first
			record.Status = StatusPendingFirst
			_ = s.mirrorNotifier.NotifyPendingSecondSignature(ctx, record, record.Parent1ID, record.Parent2ID)
		}
	}

	if err := s.repo.UpdateAuthorization(ctx, record); err != nil {
		return nil, fmt.Errorf("updating authorization: %w", err)
	}

	return record, nil
}

func (s *DualSignatureService) RejectDocument(ctx context.Context, authID, parentUserID, reason string) (*DualParentalAuthorization, error) {
	record, err := s.repo.GetAuthorizationByID(ctx, authID)
	if err != nil {
		return nil, errors.New("authorization not found")
	}

	if record.Parent1ID != parentUserID && record.Parent2ID != parentUserID {
		return nil, errors.New("unauthorized: user is not a designated parent for this authorization")
	}

	record.Status = StatusRejected
	record.RejectionReason = reason

	if err := s.repo.UpdateAuthorization(ctx, record); err != nil {
		return nil, fmt.Errorf("rejecting authorization: %w", err)
	}

	// Notify other parent of rejection
	otherParent := record.Parent1ID
	if otherParent == parentUserID {
		otherParent = record.Parent2ID
	}
	if otherParent != "" {
		_ = s.mirrorNotifier.NotifyRejection(ctx, record, otherParent, reason)
	}

	return record, nil
}

func (s *DualSignatureService) ListAuthorizations(ctx context.Context, parentUserID string) ([]DualParentalAuthorization, error) {
	return s.repo.ListAuthorizationsForParent(ctx, parentUserID)
}
