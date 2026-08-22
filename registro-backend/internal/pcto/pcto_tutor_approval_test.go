package pcto

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type tutorTestMockRepo struct {
	MockRepo
	hourLog       *HourLog
	participation *Participation
	project       *Project
	statusUpdated string
}

func (m *tutorTestMockRepo) GetHourLogByID(ctx context.Context, id string) (*HourLog, error) {
	return m.hourLog, nil
}

func (m *tutorTestMockRepo) GetParticipationByID(ctx context.Context, id string) (*Participation, error) {
	return m.participation, nil
}

func (m *tutorTestMockRepo) GetProjectByID(ctx context.Context, id string) (*Project, error) {
	return m.project, nil
}

func (m *tutorTestMockRepo) UpdateHourLogStatus(ctx context.Context, logID, status string) error {
	m.statusUpdated = status
	return nil
}

func TestPCTO_ApproveHours_TutorOwnership(t *testing.T) {
	tutorID := "teacher-tutor-1"
	otherTeacherID := "teacher-other-2"

	repo := &tutorTestMockRepo{
		hourLog:       &HourLog{ID: "log-1", ParticipationID: "part-1"},
		participation: &Participation{ID: "part-1", ProjectID: "proj-1"},
		project:       &Project{ID: "proj-1", SchoolTutorID: &tutorID, CreatedBy: tutorID},
	}
	svc := NewService(repo)

	// 1. Other teacher cannot approve
	err := svc.ApproveHours(context.Background(), otherTeacherID, "teacher", "log-1", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non sei il tutor scolastico")

	// 2. Assigned tutor can approve
	err = svc.ApproveHours(context.Background(), tutorID, "teacher", "log-1", true)
	assert.NoError(t, err)
	assert.Equal(t, "approved", repo.statusUpdated)

	// 3. Admin can approve
	err = svc.ApproveHours(context.Background(), "admin-1", "admin", "log-1", false)
	assert.NoError(t, err)
	assert.Equal(t, "rejected", repo.statusUpdated)
}
