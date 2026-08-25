package groups

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockGroupsRepoForSec struct {
	groups map[string]*Group
}

func (m *mockGroupsRepoForSec) Create(ctx context.Context, g *Group) error {
	g.ID = "g-1"
	m.groups[g.ID] = g
	return nil
}
func (m *mockGroupsRepoForSec) Update(ctx context.Context, g *Group) error {
	m.groups[g.ID] = g
	return nil
}
func (m *mockGroupsRepoForSec) Delete(ctx context.Context, id string) error {
	delete(m.groups, id)
	return nil
}
func (m *mockGroupsRepoForSec) GetByID(ctx context.Context, id string) (*Group, error) {
	if g, ok := m.groups[id]; ok {
		return g, nil
	}
	return nil, ErrGroupNotFound
}
func (m *mockGroupsRepoForSec) ListBySchool(ctx context.Context, schoolID string) ([]Group, error) {
	var res []Group
	for _, g := range m.groups {
		if g.SchoolID == schoolID {
			res = append(res, *g)
		}
	}
	return res, nil
}
func (m *mockGroupsRepoForSec) ListByTeacher(ctx context.Context, teacherID string) ([]Group, error) {
	return nil, nil
}
func (m *mockGroupsRepoForSec) ListByStudent(ctx context.Context, studentID string) ([]Group, error) {
	return nil, nil
}
func (m *mockGroupsRepoForSec) AddStudents(ctx context.Context, groupID string, studentIDs []string) error {
	return nil
}
func (m *mockGroupsRepoForSec) RemoveStudent(ctx context.Context, groupID, studentID string) error {
	return nil
}
func (m *mockGroupsRepoForSec) GetStudentsInGroup(ctx context.Context, groupID string) ([]GroupStudentInfo, error) {
	return nil, nil
}

func TestGroups_CrossTenant_Update_Blocked(t *testing.T) {
	repo := &mockGroupsRepoForSec{
		groups: map[string]*Group{
			"g-1": {ID: "g-1", SchoolID: "school-1", Name: "Gruppo A"},
		},
	}
	svc := NewService(repo)

	// Admin of school-2 attempts to update group in school-1
	_, err := svc.UpdateGroup(context.Background(), "admin", "school-2", "g-1", UpdateGroupRequest{
		Name: "Hacked Name",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestGroups_CrossTenant_Delete_Blocked(t *testing.T) {
	repo := &mockGroupsRepoForSec{
		groups: map[string]*Group{
			"g-1": {ID: "g-1", SchoolID: "school-1", Name: "Gruppo A"},
		},
	}
	svc := NewService(repo)

	// Admin of school-2 attempts to delete group in school-1
	err := svc.DeleteGroup(context.Background(), "admin", "school-2", "g-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}
