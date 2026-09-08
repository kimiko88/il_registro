package visitors

import (
	"context"
	"fmt"
)

type Service interface {
	RegisterVisitor(ctx context.Context, schoolID, recordedBy string, req RegisterVisitorRequest) (*Visitor, error)
	RecordVisitorExit(ctx context.Context, id, schoolID, notes string) error
	ListTodayVisitors(ctx context.Context, schoolID, date string) ([]*Visitor, error)

	RecordEarlyExit(ctx context.Context, schoolID, recordedBy string, req RecordEarlyExitRequest) (*EarlyExit, error)
	RecordStudentReturn(ctx context.Context, id, schoolID, notes string) error
	ListTodayEarlyExits(ctx context.Context, schoolID, date string) ([]*EarlyExit, error)

	CreateMaintenanceReport(ctx context.Context, schoolID, reportedBy string, req CreateMaintenanceReportRequest) (*MaintenanceReport, error)
	ListMaintenanceReports(ctx context.Context, schoolID, statusFilter string) ([]*MaintenanceReport, error)
	UpdateMaintenanceStatus(ctx context.Context, id, schoolID string, req UpdateMaintenanceStatusRequest) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	if repo == nil {
		panic("visitors.NewService: repo must not be nil")
	}
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) RegisterVisitor(ctx context.Context, schoolID, recordedBy string, req RegisterVisitorRequest) (*Visitor, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("il nome del visitatore è obbligatorio")
	}
	badge := req.BadgeNumber
	var badgePtr *string
	if badge != "" {
		badgePtr = &badge
	}
	v := &Visitor{
		SchoolID:    schoolID,
		Name:        req.Name,
		DocumentID:  req.DocumentID,
		Purpose:     req.Purpose,
		HostName:    req.HostName,
		BadgeNumber: badgePtr,
		Notes:       req.Notes,
		RecordedBy:  recordedBy,
	}
	if err := s.repo.CreateVisitor(ctx, v); err != nil {
		return nil, fmt.Errorf("errore registrazione visitatore: %w", err)
	}
	return v, nil
}

func (s *serviceImpl) RecordVisitorExit(ctx context.Context, id, schoolID, notes string) error {
	return s.repo.RecordVisitorExit(ctx, id, schoolID, notes)
}

func (s *serviceImpl) ListTodayVisitors(ctx context.Context, schoolID, date string) ([]*Visitor, error) {
	return s.repo.ListTodayVisitors(ctx, schoolID, date)
}

func (s *serviceImpl) RecordEarlyExit(ctx context.Context, schoolID, recordedBy string, req RecordEarlyExitRequest) (*EarlyExit, error) {
	e := &EarlyExit{
		SchoolID:      schoolID,
		StudentID:     req.StudentID,
		DelegateeName: req.DelegateeName,
		DelegateRel:   req.DelegateRel,
		ReasonCode:    req.ReasonCode,
		Notes:         req.Notes,
		RecordedBy:    recordedBy,
	}
	if err := s.repo.CreateEarlyExit(ctx, e); err != nil {
		return nil, fmt.Errorf("errore registrazione uscita anticipata: %w", err)
	}
	return e, nil
}

func (s *serviceImpl) RecordStudentReturn(ctx context.Context, id, schoolID, notes string) error {
	return s.repo.RecordStudentReturn(ctx, id, schoolID, notes)
}

func (s *serviceImpl) ListTodayEarlyExits(ctx context.Context, schoolID, date string) ([]*EarlyExit, error) {
	return s.repo.ListTodayEarlyExits(ctx, schoolID, date)
}

func (s *serviceImpl) CreateMaintenanceReport(ctx context.Context, schoolID, reportedBy string, req CreateMaintenanceReportRequest) (*MaintenanceReport, error) {
	rep := &MaintenanceReport{
		SchoolID:    schoolID,
		Location:    req.Location,
		Category:    req.Category,
		Description: req.Description,
		Priority:    req.Priority,
		ReportedBy:  reportedBy,
	}
	if err := s.repo.CreateMaintenanceReport(ctx, rep); err != nil {
		return nil, fmt.Errorf("errore creazione segnalazione: %w", err)
	}
	return rep, nil
}

func (s *serviceImpl) ListMaintenanceReports(ctx context.Context, schoolID, statusFilter string) ([]*MaintenanceReport, error) {
	return s.repo.ListMaintenanceReports(ctx, schoolID, statusFilter)
}

func (s *serviceImpl) UpdateMaintenanceStatus(ctx context.Context, id, schoolID string, req UpdateMaintenanceStatusRequest) error {
	return s.repo.UpdateMaintenanceStatus(ctx, id, schoolID, req)
}
