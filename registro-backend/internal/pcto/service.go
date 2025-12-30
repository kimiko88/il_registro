package pcto

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	CreateProject(ctx context.Context, teacherID string, req CreateProjectRequest) error
	GetProjects(ctx context.Context) ([]Project, error)
	AssignStudent(ctx context.Context, projectID, studentID string) error

	LogHours(ctx context.Context, studentID string, req LogHourRequest) error
	GetMyProjects(ctx context.Context, studentID string) ([]Project, error)
	GetMyProjectDetails(ctx context.Context, studentID, projectID string) (*Participation, []HourLog, error)

	CreateCompany(ctx context.Context, c Company) error
	GetCompanies(ctx context.Context) ([]Company, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateProject(ctx context.Context, teacherID string, req CreateProjectRequest) error {
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return errors.New("invalid start_date format (expected YYYY-MM-DD)")
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return errors.New("invalid end_date format (expected YYYY-MM-DD)")
	}

	if end.Before(start) {
		return errors.New("end_date cannot be before start_date")
	}

	p := &Project{
		SchoolID:      "default-school",
		Title:         req.Title,
		Description:   req.Description,
		Type:          req.Type,
		StartDate:     start,
		EndDate:       end,
		TotalHours:    req.TotalHours,
		CompanyID:     req.CompanyID,
		SchoolTutorID: &teacherID,
		CompanyTutor:  req.CompanyTutor,
		CreatedBy:     teacherID,
	}
	return s.repo.CreateProject(ctx, p)
}

func (s *service) GetProjects(ctx context.Context) ([]Project, error) {
	return s.repo.GetProjects(ctx, "default-school")
}

func (s *service) AssignStudent(ctx context.Context, projectID, studentID string) error {
	// Check existing?
	part := &Participation{ProjectID: projectID, StudentID: studentID}
	return s.repo.AssignStudent(ctx, part)
}

func (s *service) LogHours(ctx context.Context, studentID string, req LogHourRequest) error {
	// Verify participation
	part, err := s.repo.GetParticipation(ctx, req.ProjectID, studentID)
	if err != nil {
		return errors.New("student not assigned to project")
	}

	date, _ := time.Parse("2006-01-02", req.Date)

	log := &HourLog{
		ParticipationID: part.ID,
		Date:            date,
		Hours:           req.Hours,
		Activity:        req.Activity,
	}
	return s.repo.LogHours(ctx, log)
}

func (s *service) GetMyProjects(ctx context.Context, studentID string) ([]Project, error) {
	parts, err := s.repo.GetParticipationsByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	var projects []Project
	for _, p := range parts {
		proj, _ := s.repo.GetProjectByID(ctx, p.ProjectID)
		if proj != nil {
			projects = append(projects, *proj)
		}
	}
	return projects, nil
}

func (s *service) GetMyProjectDetails(ctx context.Context, studentID, projectID string) (*Participation, []HourLog, error) {
	part, err := s.repo.GetParticipation(ctx, projectID, studentID)
	if err != nil {
		return nil, nil, err
	}

	logs, err := s.repo.GetHours(ctx, part.ID)
	return part, logs, err
}

func (s *service) CreateCompany(ctx context.Context, c Company) error {
	c.SchoolID = "default-school"
	return s.repo.CreateCompany(ctx, &c)
}

func (s *service) GetCompanies(ctx context.Context) ([]Company, error) {
	return s.repo.GetCompanies(ctx, "default-school")
}
