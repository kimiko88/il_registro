package may15

import (
	"context"
	"fmt"
)

type Service interface {
	GetDocument(ctx context.Context, classID, academicYear string) (*ClassMay15Document, error)
	SaveDocument(ctx context.Context, schoolID, classID string, req SaveMay15Request) (*ClassMay15Document, error)
	PublishDocument(ctx context.Context, classID, academicYear string) (*ClassMay15Document, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetDocument(ctx context.Context, classID, academicYear string) (*ClassMay15Document, error) {
	if classID == "" {
		return nil, fmt.Errorf("class_id is required")
	}
	if academicYear == "" {
		academicYear = "2025/2026"
	}
	doc, err := s.repo.GetByClassAndYear(ctx, classID, academicYear)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		// Return empty draft skeleton
		return &ClassMay15Document{
			ClassID:      classID,
			AcademicYear: academicYear,
			Status:       StatusBozza,
		}, nil
	}
	return doc, nil
}

func (s *service) SaveDocument(ctx context.Context, schoolID, classID string, req SaveMay15Request) (*ClassMay15Document, error) {
	if classID == "" {
		return nil, fmt.Errorf("class_id is required")
	}
	if req.AcademicYear == "" {
		req.AcademicYear = "2025/2026"
	}
	if req.Status == "" {
		req.Status = StatusBozza
	}

	doc := &ClassMay15Document{
		SchoolID:           schoolID,
		ClassID:            classID,
		AcademicYear:       req.AcademicYear,
		Status:             req.Status,
		ClassPresentation:  req.ClassPresentation,
		TeachingContinuity: req.TeachingContinuity,
		PCTOPathways:       req.PCTOPathways,
		ExamSimulations:    req.ExamSimulations,
		EvaluationRubrics:  req.EvaluationRubrics,
		CLILModules:        req.CLILModules,
	}

	if err := s.repo.Upsert(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *service) PublishDocument(ctx context.Context, classID, academicYear string) (*ClassMay15Document, error) {
	if classID == "" {
		return nil, fmt.Errorf("class_id is required")
	}
	if academicYear == "" {
		academicYear = "2025/2026"
	}
	return s.repo.Publish(ctx, classID, academicYear)
}
