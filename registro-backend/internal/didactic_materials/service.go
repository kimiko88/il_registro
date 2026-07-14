package didactic_materials

import "errors"

type Service interface {
	CreateMaterial(teacherID string, schoolID string, req CreateMaterialRequest) (*MaterialResponse, error)
	GetMaterialsByClass(classID string) ([]MaterialResponse, error)
	DeleteMaterial(id string, teacherID string) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateMaterial(teacherID string, schoolID string, req CreateMaterialRequest) (*MaterialResponse, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	m := &DidacticMaterial{
		SchoolID:      schoolID,
		ClassID:       req.ClassID,
		SubjectID:     req.SubjectID,
		TeacherID:     teacherID,
		Title:         req.Title,
		Description:   req.Description,
		AttachmentURL: req.AttachmentURL,
	}

	if err := s.repo.Create(m); err != nil {
		return nil, err
	}

	return s.mapResponse(m), nil
}

func (s *service) GetMaterialsByClass(classID string) ([]MaterialResponse, error) {
	list, err := s.repo.GetByClass(classID)
	if err != nil {
		return nil, err
	}

	var res []MaterialResponse
	for _, m := range list {
		res = append(res, *s.mapResponse(&m))
	}
	return res, nil
}

func (s *service) DeleteMaterial(id string, teacherID string) error {
	return s.repo.Delete(id, teacherID)
}

func (s *service) mapResponse(m *DidacticMaterial) *MaterialResponse {
	return &MaterialResponse{
		ID:            m.ID,
		ClassID:       m.ClassID,
		SubjectID:     m.SubjectID,
		TeacherID:     m.TeacherID,
		TeacherName:   m.TeacherName,
		Title:         m.Title,
		Description:   m.Description,
		AttachmentURL: m.AttachmentURL,
		CreatedAt:     m.CreatedAt,
	}
}
