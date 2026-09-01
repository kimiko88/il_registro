package pcto

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock
type MockRepo struct{ mock.Mock }

func (m *MockRepo) CreateProject(ctx context.Context, p *Project) error          { p.ID = "p1"; return nil }
func (m *MockRepo) GetProjects(ctx context.Context, s string) ([]Project, error) { return nil, nil }
func (m *MockRepo) CreateCompany(ctx context.Context, c *Company) error          { return nil }
func (m *MockRepo) AssignStudent(ctx context.Context, p *Participation) error    { return nil }
func (m *MockRepo) LogHours(ctx context.Context, h *HourLog) error               { return nil }
func (m *MockRepo) GetParticipation(ctx context.Context, p, s string) (*Participation, error) {
	return &Participation{ID: "part1"}, nil
}

// Stubs
func (m *MockRepo) GetProjectByID(ctx context.Context, id string) (*Project, error) {
	return &Project{}, nil
}
func (m *MockRepo) UpdateProject(ctx context.Context, p *Project) error { return nil }
func (m *MockRepo) GetParticipationsByProject(ctx context.Context, id string) ([]Participation, error) {
	return nil, nil
}
func (m *MockRepo) GetParticipationsByStudent(ctx context.Context, id string) ([]Participation, error) {
	return []Participation{{ProjectID: "p1"}}, nil
}
func (m *MockRepo) GetHours(ctx context.Context, id string) ([]HourLog, error)    { return nil, nil }
func (m *MockRepo) VerifyHours(ctx context.Context, h, t string) error            { return nil }
func (m *MockRepo) GetCompanies(ctx context.Context, s string) ([]Company, error) { return nil, nil }
func (m *MockRepo) DeleteProject(ctx context.Context, id string) error            { return nil }
func (m *MockRepo) GetHourLogByID(ctx context.Context, id string) (*HourLog, error) {
	return &HourLog{ID: id, ParticipationID: "part1"}, nil
}
func (m *MockRepo) GetParticipationByID(ctx context.Context, id string) (*Participation, error) {
	return &Participation{ID: id, ProjectID: "p1"}, nil
}
func (m *MockRepo) GetStats(ctx context.Context, schoolID string) (*PCTOStats, error) {
	return nil, nil
}
func (m *MockRepo) UpdateHourLogStatus(ctx context.Context, logID, status string) error { return nil }

func TestService_CreateProject(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	req := CreateProjectRequest{Title: "Internship", Type: "External", StartDate: "2025-01-01", EndDate: "2025-02-01"}
	err := svc.CreateProject(context.Background(), "school1", "secretary", "t1", req)
	assert.NoError(t, err)
}

func TestService_LogHours(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	// Valid logging
	err := svc.LogHours(context.Background(), "s1", LogHourRequest{ProjectID: "p1", Date: "2025-01-05", Hours: 4})
	assert.NoError(t, err)

	// Zero or negative hours rejected
	err = svc.LogHours(context.Background(), "s1", LogHourRequest{ProjectID: "p1", Date: "2025-01-05", Hours: 0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hours must be greater than 0")

	err = svc.LogHours(context.Background(), "s1", LogHourRequest{ProjectID: "p1", Date: "2025-01-05", Hours: -2})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hours must be greater than 0")

	// Invalid date format rejected
	err = svc.LogHours(context.Background(), "s1", LogHourRequest{ProjectID: "p1", Date: "invalid-date", Hours: 3})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid date format")
}

func TestService_CreateProject_Validation(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)
	ctx := context.Background()

	// End before start rejected
	err := svc.CreateProject(ctx, "school1", "secretary", "t1", CreateProjectRequest{
		Title:     "Project",
		Type:      "Internal",
		StartDate: "2025-06-01",
		EndDate:   "2025-05-01",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "end_date cannot be before start_date")

	// Invalid date format rejected
	err = svc.CreateProject(ctx, "school1", "secretary", "t1", CreateProjectRequest{
		Title:     "Project",
		Type:      "Internal",
		StartDate: "invalid",
		EndDate:   "2025-05-01",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid start_date format")
}

func TestService_UpdateAndDelete_CreatorCheck(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)
	ctx := context.Background()

	// Mock project created by t1
	tutor := "t1"
	repo.On("GetProjectByID", ctx, "p1").Return(&Project{ID: "p1", CreatedBy: "t1", SchoolTutorID: &tutor}, nil)

	// Secretary s2 (who has PCTOUpdate permission) is not the creator, so update is blocked
	err := svc.UpdateProject(ctx, "s2", "secretary", "p1", CreateProjectRequest{Title: "Updated"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized: not the creator")

	// Secretary s2 is not creator, so delete is blocked
	err = svc.DeleteProject(ctx, "s2", "secretary", "p1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized: not the creator")

	// Admin is authorized even if not the creator
	err = svc.UpdateProject(ctx, "admin1", "admin", "p1", CreateProjectRequest{Title: "Updated By Admin"})
	assert.NoError(t, err)
}
