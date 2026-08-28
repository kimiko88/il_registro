package payments

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"registro-backend/internal/users"
)

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrPaymentNotFound = errors.New("pagamento non trovato")
)

type Service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(repo Repository, uRepo ...users.Repository) *Service {
	if repo == nil {
		panic("payments.NewService: repo must not be nil")
	}
	svc := &Service{repo: repo}
	if len(uRepo) > 0 && uRepo[0] != nil {
		svc.userRepo = uRepo[0]
	}
	return svc
}

func (s *Service) List(ctx context.Context, actorID, actorRole, schoolID, targetStudentID string) (*PaymentSummaryResponse, error) {
	if actorID == "" {
		return nil, ErrUnauthorized
	}

	var list []*SchoolPayment
	var err error

	switch actorRole {
	case "parent":
		if targetStudentID != "" {
			if s.userRepo != nil {
				isG, errG := s.userRepo.IsGuardian(ctx, actorID, targetStudentID)
				if errG != nil || !isG {
					return nil, ErrUnauthorized
				}
			}
			list, err = s.repo.ListByStudent(ctx, schoolID, targetStudentID)
		} else {
			list, err = s.repo.ListByParent(ctx, schoolID, actorID)
		}
	case "student":
		list, err = s.repo.ListByStudent(ctx, schoolID, actorID)
	case "admin", "superadmin", "secretary", "principal", "vice_principal":
		if targetStudentID != "" {
			list, err = s.repo.ListByStudent(ctx, schoolID, targetStudentID)
		} else {
			list, err = s.repo.ListBySchool(ctx, schoolID)
		}
	default:
		return nil, ErrUnauthorized
	}

	if err != nil {
		return nil, err
	}

	summary := &PaymentSummaryResponse{
		Payments: list,
	}

	for _, p := range list {
		switch p.Status {
		case "pending":
			summary.PendingTotal += p.Amount
			summary.PendingCount++
		case "paid":
			summary.PaidCount++
		}
	}

	return summary, nil
}

func (s *Service) GetByID(ctx context.Context, actorID, actorRole, schoolID, id string) (*SchoolPayment, error) {
	if actorID == "" {
		return nil, ErrUnauthorized
	}
	p, err := s.repo.GetByID(ctx, id)
	if err != nil || p == nil {
		return nil, ErrPaymentNotFound
	}
	if actorRole != "superadmin" && schoolID != "" && p.SchoolID != "" && p.SchoolID != schoolID {
		return nil, ErrUnauthorized
	}
	if actorRole == "student" {
		if p.StudentID != actorID {
			return nil, ErrUnauthorized
		}
	} else if actorRole == "parent" {
		if s.userRepo != nil && p.StudentID != "" {
			isG, errG := s.userRepo.IsGuardian(ctx, actorID, p.StudentID)
			if errG != nil || !isG {
				return nil, ErrUnauthorized
			}
		}
	}
	return p, nil
}

func (s *Service) Pay(ctx context.Context, actorID, actorRole, id, method string) (*SchoolPayment, error) {
	if actorID == "" {
		return nil, ErrUnauthorized
	}
	if actorRole != "parent" && actorRole != "student" && actorRole != "admin" && actorRole != "superadmin" {
		return nil, ErrUnauthorized
	}

	p, err := s.repo.GetByID(ctx, id)
	if err != nil || p == nil {
		return nil, ErrPaymentNotFound
	}
	if actorRole == "student" {
		if p.StudentID != actorID {
			return nil, ErrUnauthorized
		}
	} else if actorRole == "parent" {
		if s.userRepo != nil && p.StudentID != "" {
			isG, errG := s.userRepo.IsGuardian(ctx, actorID, p.StudentID)
			if errG != nil || !isG {
				return nil, ErrUnauthorized
			}
		}
	}
	if p.Status == "paid" {
		return nil, fmt.Errorf("questo contributo o tassa è già stato pagato")
	}

	if method == "" {
		method = "PagoPA"
	}

	randBytes := make([]byte, 6)
	_, _ = rand.Read(randBytes)
	txID := fmt.Sprintf("PAGOPA-%s", hex.EncodeToString(randBytes))
	receiptNum := fmt.Sprintf("REC-%d-%s", time.Now().Year(), hex.EncodeToString(randBytes[:4]))

	if err := s.repo.Pay(ctx, id, actorID, method, txID, receiptNum); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, actorID, actorRole, schoolID string, req CreatePaymentRequest) (*SchoolPayment, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" {
		return nil, ErrUnauthorized
	}
	if schoolID == "" {
		return nil, fmt.Errorf("school_id is required")
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		dueDate, err = time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid due_date format, expected YYYY-MM-DD")
		}
	}

	p := &SchoolPayment{
		SchoolID:    schoolID,
		StudentID:   req.StudentID,
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
		DueDate:     dueDate,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
