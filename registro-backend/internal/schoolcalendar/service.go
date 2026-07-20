package schoolcalendar

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Service espone le operazioni sul calendario scolastico.
type Service interface {
	SetSchoolYear(ctx context.Context, actorID, actorRole, schoolID string, req SetYearRequest) (*SchoolYearResponse, error)
	GetSchoolYear(ctx context.Context, schoolID string) (*SchoolYearResponse, error)

	AddNonTeachingDay(ctx context.Context, actorID, actorRole, schoolID string, req AddNonTeachingDayRequest) (*NonTeachingDayResponse, error)
	DeleteNonTeachingDay(ctx context.Context, actorRole, schoolID, id string) error
	ListNonTeachingDays(ctx context.Context, schoolID string) ([]NonTeachingDayResponse, error)

	// CountTeachingDays implementa l'interfaccia attendance.CalendarService.
	CountTeachingDays(ctx context.Context, schoolID string, from, to time.Time) (int, error)
	// GetSchoolYearDates implementa l'interfaccia attendance.CalendarService.
	GetSchoolYearDates(ctx context.Context, schoolID string) (start, end time.Time, err error)

	CreateAcademicPeriod(ctx context.Context, actorRole, schoolID string, req CreateAcademicPeriodRequest) (*AcademicPeriod, error)
	ListAcademicPeriods(ctx context.Context, schoolID string) ([]AcademicPeriod, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func isSecretary(role string) bool {
	return role == "secretary" || role == "admin" || role == "superadmin"
}

func (s *service) SetSchoolYear(ctx context.Context, actorID, actorRole, schoolID string, req SetYearRequest) (*SchoolYearResponse, error) {
	if !isSecretary(actorRole) {
		return nil, errors.New("forbidden: solo la segreteria può impostare l'anno scolastico")
	}
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("start_date non valida: %w", err)
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("end_date non valida: %w", err)
	}
	if !end.After(start) {
		return nil, errors.New("end_date deve essere successiva a start_date")
	}

	settings := &SchoolYearSettings{
		SchoolID:  schoolID,
		YearLabel: req.YearLabel,
		StartDate: start,
		EndDate:   end,
		CreatedBy: actorID,
	}
	if err := s.repo.UpsertYear(settings); err != nil {
		return nil, err
	}
	return s.buildYearResponse(ctx, settings)
}

func (s *service) GetSchoolYear(ctx context.Context, schoolID string) (*SchoolYearResponse, error) {
	settings, err := s.repo.GetYear(schoolID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("anno scolastico non ancora configurato")
		}
		return nil, err
	}
	return s.buildYearResponse(ctx, settings)
}

func (s *service) buildYearResponse(ctx context.Context, settings *SchoolYearSettings) (*SchoolYearResponse, error) {
	teachingDays, err := s.repo.CountTeachingDays(settings.SchoolID, settings.StartDate, settings.EndDate)
	if err != nil {
		teachingDays = 0
	}
	return &SchoolYearResponse{
		ID:           settings.ID,
		SchoolID:     settings.SchoolID,
		YearLabel:    settings.YearLabel,
		StartDate:    settings.StartDate,
		EndDate:      settings.EndDate,
		TeachingDays: teachingDays,
	}, nil
}

func (s *service) AddNonTeachingDay(ctx context.Context, actorID, actorRole, schoolID string, req AddNonTeachingDayRequest) (*NonTeachingDayResponse, error) {
	if !isSecretary(actorRole) {
		return nil, errors.New("forbidden: solo la segreteria può gestire i giorni non didattici")
	}
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("date non valida: %w", err)
	}
	d := &NonTeachingDay{
		SchoolID:  schoolID,
		Date:      date,
		Label:     req.Label,
		CreatedBy: actorID,
	}
	if err := s.repo.AddNonTeachingDay(d); err != nil {
		return nil, err
	}
	return &NonTeachingDayResponse{ID: d.ID, Date: d.Date.Format("2006-01-02"), Label: d.Label}, nil
}

func (s *service) DeleteNonTeachingDay(ctx context.Context, actorRole, schoolID, id string) error {
	if !isSecretary(actorRole) {
		return errors.New("forbidden: solo la segreteria può gestire i giorni non didattici")
	}
	return s.repo.DeleteNonTeachingDay(schoolID, id)
}

func (s *service) ListNonTeachingDays(ctx context.Context, schoolID string) ([]NonTeachingDayResponse, error) {
	days, err := s.repo.ListNonTeachingDays(schoolID)
	if err != nil {
		return nil, err
	}
	var res []NonTeachingDayResponse
	for _, d := range days {
		res = append(res, NonTeachingDayResponse{ID: d.ID, Date: d.Date.Format("2006-01-02"), Label: d.Label})
	}
	if res == nil {
		res = []NonTeachingDayResponse{}
	}
	return res, nil
}

func (s *service) CountTeachingDays(ctx context.Context, schoolID string, from, to time.Time) (int, error) {
	return s.repo.CountTeachingDays(schoolID, from, to)
}

func (s *service) GetSchoolYearDates(ctx context.Context, schoolID string) (start, end time.Time, err error) {
	settings, err := s.repo.GetYear(schoolID)
	if err != nil {
		return
	}
	return settings.StartDate, settings.EndDate, nil
}

func (s *service) CreateAcademicPeriod(ctx context.Context, actorRole, schoolID string, req CreateAcademicPeriodRequest) (*AcademicPeriod, error) {
	if !isSecretary(actorRole) {
		return nil, errors.New("forbidden: solo la segreteria può creare i periodi valutativi")
	}
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("start_date non valida: %w", err)
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("end_date non valida: %w", err)
	}

	p := &AcademicPeriod{
		SchoolID:       schoolID,
		AcademicYearID: req.AcademicYearID,
		Name:           req.Name,
		Code:           req.Code,
		StartDate:      start,
		EndDate:        end,
		IsCurrent:      req.IsCurrent,
	}

	if err := s.repo.CreateAcademicPeriod(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) ListAcademicPeriods(ctx context.Context, schoolID string) ([]AcademicPeriod, error) {
	return s.repo.ListAcademicPeriods(schoolID)
}
