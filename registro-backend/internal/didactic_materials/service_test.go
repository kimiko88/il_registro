package didactic_materials

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	materials []DidacticMaterial
	errCreate error
	errGet    error
	errDelete error
}

func (m *mockRepository) Create(mat *DidacticMaterial) error {
	if m.errCreate != nil {
		return m.errCreate
	}
	mat.ID = "test-id"
	mat.TeacherName = "John Doe"
	m.materials = append(m.materials, *mat)
	return nil
}

func (m *mockRepository) GetByClass(classID string) ([]DidacticMaterial, error) {
	if m.errGet != nil {
		return nil, m.errGet
	}
	var res []DidacticMaterial
	for _, mat := range m.materials {
		if mat.ClassID == classID {
			res = append(res, mat)
		}
	}
	return res, nil
}

func (m *mockRepository) Delete(id string, teacherID string) error {
	if m.errDelete != nil {
		return m.errDelete
	}
	var keep []DidacticMaterial
	for _, mat := range m.materials {
		if mat.ID == id && mat.TeacherID == teacherID {
			continue
		}
		keep = append(keep, mat)
	}
	m.materials = keep
	return nil
}

func TestCreateMaterial(t *testing.T) {
	repo := &mockRepository{}
	s := NewService(repo, nil)

	// Test validation error
	req := CreateMaterialRequest{
		ClassID:   "class-1",
		SubjectID: "math",
		Title:     "",
	}
	resp, err := s.CreateMaterial(context.Background(), "teacher-1", "school-1", req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "title is required")

	// Test successful creation
	req.Title = "Math Slide"
	req.Description = "Class slides"
	resp, err = s.CreateMaterial(context.Background(), "teacher-1", "school-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-id", resp.ID)
	assert.Equal(t, "John Doe", resp.TeacherName)
	assert.Equal(t, "Math Slide", resp.Title)

	// Test repository error
	repo.errCreate = errors.New("db error")
	resp, err = s.CreateMaterial(context.Background(), "teacher-1", "school-1", req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "db error")
}

func TestGetMaterialsByClass(t *testing.T) {
	repo := &mockRepository{
		materials: []DidacticMaterial{
			{ID: "id-1", ClassID: "class-1", Title: "Math"},
			{ID: "id-2", ClassID: "class-2", Title: "Physics"},
		},
	}
	s := NewService(repo, nil)

	res, err := s.GetMaterialsByClass(context.Background(), "user-1", "teacher", "class-1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "Math", res[0].Title)

	// Test repository error
	repo.errGet = errors.New("get error")
	res, err = s.GetMaterialsByClass(context.Background(), "user-1", "teacher", "class-1")
	assert.Nil(t, res)
	assert.EqualError(t, err, "get error")
}

func TestDeleteMaterial(t *testing.T) {
	repo := &mockRepository{
		materials: []DidacticMaterial{
			{ID: "id-1", TeacherID: "teacher-1", Title: "Math"},
		},
	}
	s := NewService(repo, nil)

	err := s.DeleteMaterial(context.Background(), "id-1", "teacher-1")
	assert.NoError(t, err)
	assert.Len(t, repo.materials, 0)
}
