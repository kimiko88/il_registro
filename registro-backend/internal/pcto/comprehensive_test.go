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

	err := svc.LogHours(context.Background(), "s1", LogHourRequest{ProjectID: "p1", Date: "2025-01-05", Hours: 4})
	assert.NoError(t, err)
}
