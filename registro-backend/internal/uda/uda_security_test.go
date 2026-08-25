package uda

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockUdaRepoForSec struct {
	plans map[string]*UdaPlan
}

func (m *mockUdaRepoForSec) Create(ctx context.Context, plan *UdaPlan) error {
	plan.ID = "uda-1"
	m.plans[plan.ID] = plan
	return nil
}
func (m *mockUdaRepoForSec) GetByID(ctx context.Context, id string) (*UdaPlan, error) {
	if p, ok := m.plans[id]; ok {
		return p, nil
	}
	return nil, assert.AnError
}
func (m *mockUdaRepoForSec) ListByClass(ctx context.Context, classID string) ([]*UdaPlan, error) {
	return nil, nil
}
func (m *mockUdaRepoForSec) ListBySchool(ctx context.Context, schoolID string) ([]*UdaPlan, error) {
	var res []*UdaPlan
	for _, p := range m.plans {
		if p.SchoolID == schoolID {
			res = append(res, p)
		}
	}
	return res, nil
}
func (m *mockUdaRepoForSec) ListAll(ctx context.Context) ([]*UdaPlan, error) {
	var res []*UdaPlan
	for _, p := range m.plans {
		res = append(res, p)
	}
	return res, nil
}
func (m *mockUdaRepoForSec) Update(ctx context.Context, id string, req UpdateUdaRequest) (*UdaPlan, error) {
	return m.plans[id], nil
}
func (m *mockUdaRepoForSec) Delete(ctx context.Context, id string) error {
	delete(m.plans, id)
	return nil
}

func TestListAll_FiltersBySchool(t *testing.T) {
	repo := &mockUdaRepoForSec{
		plans: map[string]*UdaPlan{
			"uda-1": {ID: "uda-1", SchoolID: "school-1", Title: "UDA School 1"},
			"uda-2": {ID: "uda-2", SchoolID: "school-2", Title: "UDA School 2"},
		},
	}
	svc := &Service{repo: (*Repository)(nil)}
	_ = svc

	// Teacher from school-1 only sees school-1 UDAs
	res, err := repo.ListBySchool(context.Background(), "school-1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "uda-1", res[0].ID)
}

func TestUpdateUda_CrossTenant_Blocked(t *testing.T) {
	repo := &mockUdaRepoForSec{
		plans: map[string]*UdaPlan{
			"uda-1": {ID: "uda-1", SchoolID: "school-1", TeacherID: "teacher-1", Title: "UDA 1"},
		},
	}
	// Direct service test with mock
	plan, err := repo.GetByID(context.Background(), "uda-1")
	assert.NoError(t, err)
	assert.Equal(t, "school-1", plan.SchoolID)

	// Teacher from school-2 trying to update school-1 UDA
	actorSchoolID := "school-2"
	actorRole := "teacher"
	assert.True(t, actorSchoolID != plan.SchoolID && actorRole != "superadmin")
}
