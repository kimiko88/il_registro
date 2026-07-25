package pcto

import (
	"context"
	"errors"
	"time"

	"registro-backend/internal/permissions"
)

type Service interface {
	CreateProject(ctx context.Context, schoolID, actorRole, teacherID string, req CreateProjectRequest) error
	GetProjects(ctx context.Context, schoolID string) ([]Project, error)
	UpdateProject(ctx context.Context, actorID, actorRole, id string, req CreateProjectRequest) error
	DeleteProject(ctx context.Context, actorID, actorRole, id string) error
	AssignStudent(ctx context.Context, actorRole, projectID, studentID string) error

	LogHours(ctx context.Context, studentID string, req LogHourRequest) error
	GetMyProjects(ctx context.Context, studentID string) ([]Project, error)
	GetMyProjectDetails(ctx context.Context, studentID, projectID string) (*Participation, []HourLog, error)

	CreateCompany(ctx context.Context, schoolID, actorRole string, c Company) error
	GetCompanies(ctx context.Context, schoolID string) ([]Company, error)
	GetStats(ctx context.Context, schoolID string) (*PCTOStats, error)
	ApproveHours(ctx context.Context, actorID, actorRole, logID string, approved bool) error
}

type service struct {
	repo        Repository
	permManager *permissions.Manager
}

func NewService(repo Repository) Service {
	return &service{
		repo:        repo,
		permManager: permissions.NewManager(),
	}
}

func (s *service) CreateProject(ctx context.Context, schoolID, actorRole, teacherID string, req CreateProjectRequest) error {
	if !s.permManager.HasPermission(actorRole, permissions.PCTOCreate) {
		return errors.New("unauthorized")
	}
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
		SchoolID:      schoolID,
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

func (s *service) GetProjects(ctx context.Context, schoolID string) ([]Project, error) {
	return s.repo.GetProjects(ctx, schoolID)
}

func (s *service) AssignStudent(ctx context.Context, actorRole, projectID, studentID string) error {
	if !s.permManager.HasPermission(actorRole, permissions.PCTOUpdate) {
		return errors.New("unauthorized")
	}
	existing, err := s.repo.GetParticipation(ctx, projectID, studentID)
	if err == nil && existing != nil {
		return errors.New("student is already assigned to this project")
	}
	part := &Participation{ProjectID: projectID, StudentID: studentID}
	return s.repo.AssignStudent(ctx, part)
}

func (s *service) LogHours(ctx context.Context, studentID string, req LogHourRequest) error {
	// Verify participation
	part, err := s.repo.GetParticipation(ctx, req.ProjectID, studentID)
	if err != nil {
		return errors.New("student not assigned to project")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return errors.New("invalid date format (expected YYYY-MM-DD)")
	}

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

func (s *service) CreateCompany(ctx context.Context, schoolID, actorRole string, c Company) error {
	if !s.permManager.HasPermission(actorRole, permissions.PCTOCreate) {
		return errors.New("unauthorized")
	}
	c.SchoolID = schoolID
	return s.repo.CreateCompany(ctx, &c)
}

func (s *service) UpdateProject(ctx context.Context, actorID, actorRole, id string, req CreateProjectRequest) error {
	if !s.permManager.HasPermission(actorRole, permissions.PCTOUpdate) {
		return errors.New("unauthorized")
	}
	p, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return err
	}
	if p.CreatedBy != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return errors.New("unauthorized: not the creator of this PCTO project")
	}
	p.Title = req.Title
	p.Description = req.Description
	p.TotalHours = req.TotalHours
	return s.repo.UpdateProject(ctx, p)
}

func (s *service) DeleteProject(ctx context.Context, actorID, actorRole, id string) error {
	if !s.permManager.HasPermission(actorRole, permissions.PCTODelete) {
		return errors.New("unauthorized")
	}
	p, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return err
	}
	if p.CreatedBy != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return errors.New("unauthorized: not the creator of this PCTO project")
	}
	return s.repo.DeleteProject(ctx, id)
}

func (s *service) GetCompanies(ctx context.Context, schoolID string) ([]Company, error) {
	return s.repo.GetCompanies(ctx, schoolID)
}

func (s *service) GetStats(ctx context.Context, schoolID string) (*PCTOStats, error) {
	return s.repo.GetStats(ctx, schoolID)
}

func (s *service) ApproveHours(ctx context.Context, actorID, actorRole, logID string, approved bool) error {
	if !s.permManager.HasPermission(actorRole, permissions.PCTOUpdate) && actorRole != "teacher" && actorRole != "tutor" && actorRole != "admin" && actorRole != "superadmin" {
		return errors.New("unauthorized")
	}
	status := "approved"
	if !approved {
		status = "rejected"
	}
	return s.repo.UpdateHourLogStatus(ctx, logID, status)
}
