package textbooks

import (
	"context"
	"errors"
	"io"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTextbook(ctx context.Context, schoolID string, req CreateTextbookRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errors.New("title is required")
	}
	if req.Price < 0 {
		return errors.New("price cannot be negative")
	}
	t := &Textbook{
		SchoolID:  schoolID,
		Title:     strings.TrimSpace(req.Title),
		Author:    strings.TrimSpace(req.Author),
		Subject:   strings.TrimSpace(req.Subject),
		ISBN:      strings.TrimSpace(req.ISBN),
		Publisher: strings.TrimSpace(req.Publisher),
		Price:     req.Price,
	}
	return s.repo.Create(ctx, t)
}

func (s *Service) UpdateTextbook(ctx context.Context, actorRole, schoolID, id string, req CreateTextbookRequest) error {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if actorRole != "superadmin" && schoolID != "" && t.SchoolID != "" && t.SchoolID != schoolID {
		return errors.New("forbidden: il libro di testo appartiene ad un'altra scuola")
	}

	if strings.TrimSpace(req.Title) == "" {
		return errors.New("title is required")
	}
	if req.Price < 0 {
		return errors.New("price cannot be negative")
	}
	t.Title = strings.TrimSpace(req.Title)
	t.Author = strings.TrimSpace(req.Author)
	t.Subject = strings.TrimSpace(req.Subject)
	t.ISBN = strings.TrimSpace(req.ISBN)
	t.Publisher = strings.TrimSpace(req.Publisher)
	t.Price = req.Price

	return s.repo.Update(ctx, t)
}

func (s *Service) ListTextbooks(ctx context.Context, schoolID string) ([]Textbook, error) {
	return s.repo.List(ctx, schoolID)
}

func (s *Service) DeleteTextbook(ctx context.Context, actorRole, schoolID, id string) error {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if actorRole != "superadmin" && schoolID != "" && t.SchoolID != "" && t.SchoolID != schoolID {
		return errors.New("forbidden: il libro di testo appartiene ad un'altra scuola")
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) AssignToClass(ctx context.Context, classID string, req AssignTextbookRequest) error {
	return s.repo.AssignToClass(ctx, classID, req.SubjectID, req.TextbookID, req.IsOptional)
}

func (s *Service) RemoveFromClass(ctx context.Context, assignmentID string) error {
	return s.repo.RemoveFromClass(ctx, assignmentID)
}

func (s *Service) ListByClass(ctx context.Context, classID string) ([]ClassTextbook, error) {
	return s.repo.ListByClass(ctx, classID)
}

// AIE & Spending Limits methods

func (s *Service) ImportAIECatalog(ctx context.Context, r io.Reader) (int, error) {
	books, err := ParseAIECatalog(r)
	if err != nil {
		return 0, err
	}
	return s.repo.UpsertAIECatalog(ctx, books)
}

func (s *Service) SearchAIECatalog(ctx context.Context, queryStr, subject, schoolOrder string, limit int) ([]AIECatalogBook, error) {
	return s.repo.SearchAIECatalog(ctx, queryStr, subject, schoolOrder, limit)
}

func (s *Service) GetClassSpendingReport(ctx context.Context, schoolID, classID string) (*SpendingReport, error) {
	classYear, schoolOrder, academicYear, _, _, err := s.repo.GetClassInfo(ctx, classID)
	if err != nil {
		return nil, err
	}

	limit, err := s.repo.GetSpendingLimit(ctx, schoolID, classYear, schoolOrder, academicYear)
	if err != nil || limit == nil {
		// Default reference limits according to D.M. 781/2013 if not configured
		limit = &SpendingLimit{
			SchoolID:            schoolID,
			ClassYear:           classYear,
			SchoolOrder:         schoolOrder,
			MaxAmount:           320.00,
			AllowedTolerancePct: 10.00,
			AcademicYear:        academicYear,
		}
	}

	adoptions, err := s.repo.ListClassAdoptions(ctx, classID)
	if err != nil {
		return nil, err
	}

	report := CalculateSpendingReport(adoptions, *limit)
	return &report, nil
}

func (s *Service) UpsertSpendingLimit(ctx context.Context, schoolID string, req SpendingLimit) error {
	req.SchoolID = schoolID
	if req.MaxAmount <= 0 {
		return errors.New("max_amount must be greater than zero")
	}
	if req.AllowedTolerancePct <= 0 {
		req.AllowedTolerancePct = 10.0
	}
	return s.repo.UpsertSpendingLimit(ctx, &req)
}

func (s *Service) AdoptBook(ctx context.Context, item ClassAdoptionItem) error {
	if item.ClassID == "" || item.SubjectID == "" || item.BookID == "" {
		return errors.New("class_id, subject_id, and book_id are required")
	}
	if item.AdoptionType == "" {
		item.AdoptionType = "nuova_adozione"
	}
	return s.repo.SaveClassAdoption(ctx, &item)
}

func (s *Service) DeleteAdoption(ctx context.Context, id string) error {
	return s.repo.DeleteClassAdoption(ctx, id)
}

func (s *Service) ExportClassAIE(ctx context.Context, classID string) (string, error) {
	classYear, _, academicYear, schoolCode, className, err := s.repo.GetClassInfo(ctx, classID)
	if err != nil {
		return "", err
	}
	_ = classYear

	adoptions, err := s.repo.ListClassAdoptions(ctx, classID)
	if err != nil {
		return "", err
	}

	var records []AIEExportRecord
	for _, a := range adoptions {
		records = append(records, AIEExportRecord{
			SchoolCode:   schoolCode,
			AcademicYear: academicYear,
			ClassName:    className,
			SubjectCode:  a.SubjectName,
			ISBN:         a.ISBN,
			Title:        a.BookTitle,
			Authors:      a.Authors,
			Publisher:    a.Publisher,
			Price:        a.Price,
			Volume:       "1",
			AdoptionType: a.AdoptionType,
			AlreadyOwned: a.IsAlreadyOwned,
		})
	}

	return ExportAIEFormat(records)
}
