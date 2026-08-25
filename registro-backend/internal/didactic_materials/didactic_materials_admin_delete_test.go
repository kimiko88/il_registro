package didactic_materials

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepoForAdminDelete struct {
	materials map[string]*DidacticMaterial
}

func (m *mockRepoForAdminDelete) Create(dm *DidacticMaterial) error {
	m.materials[dm.ID] = dm
	return nil
}
func (m *mockRepoForAdminDelete) GetByClass(classID string) ([]DidacticMaterial, error) {
	return nil, nil
}
func (m *mockRepoForAdminDelete) GetByID(id string) (*DidacticMaterial, error) {
	if dm, ok := m.materials[id]; ok {
		return dm, nil
	}
	return nil, assert.AnError
}
func (m *mockRepoForAdminDelete) Delete(id, teacherID string) error {
	delete(m.materials, id)
	return nil
}
func (m *mockRepoForAdminDelete) DeleteByID(id string) error {
	delete(m.materials, id)
	return nil
}

func TestDeleteMaterial_AdminCanDeleteOtherTeacherMaterial(t *testing.T) {
	repo := &mockRepoForAdminDelete{
		materials: map[string]*DidacticMaterial{
			"mat-1": {ID: "mat-1", SchoolID: "school-1", TeacherID: "teacher-1", Title: "Lesson Material"},
		},
	}
	svc := NewService(repo, nil)

	// Admin of the same school can delete material
	err := svc.DeleteMaterial(context.Background(), "admin-1", "admin", "school-1", "mat-1")
	assert.NoError(t, err)
	assert.Len(t, repo.materials, 0)
}

func TestDeleteMaterial_DifferentTeacherCannotDelete(t *testing.T) {
	repo := &mockRepoForAdminDelete{
		materials: map[string]*DidacticMaterial{
			"mat-1": {ID: "mat-1", SchoolID: "school-1", TeacherID: "teacher-1", Title: "Lesson Material"},
		},
	}
	svc := NewService(repo, nil)

	// Teacher-2 from the same school CANNOT delete teacher-1's material
	err := svc.DeleteMaterial(context.Background(), "teacher-2", "teacher", "school-1", "mat-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	assert.Len(t, repo.materials, 1)
}

func TestDeleteMaterial_AdminDifferentSchoolCannotDelete(t *testing.T) {
	repo := &mockRepoForAdminDelete{
		materials: map[string]*DidacticMaterial{
			"mat-1": {ID: "mat-1", SchoolID: "school-1", TeacherID: "teacher-1", Title: "Lesson Material"},
		},
	}
	svc := NewService(repo, nil)

	// Admin from school-2 CANNOT delete school-1's material
	err := svc.DeleteMaterial(context.Background(), "admin-2", "admin", "school-2", "mat-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	assert.Len(t, repo.materials, 1)
}
