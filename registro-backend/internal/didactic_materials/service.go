package didactic_materials

import (
	"context"
	"errors"

	"registro-backend/internal/users"
)

type Service interface {
	CreateMaterial(teacherID string, schoolID string, req CreateMaterialRequest) (*MaterialResponse, error)
	GetMaterialsByClass(ctx context.Context, userID, role, classID string) ([]MaterialResponse, error)
	DeleteMaterial(id string, teacherID string) error
}

type service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(r Repository, uRepo ...users.Repository) Service {
	var userRepo users.Repository
	if len(uRepo) > 0 {
		userRepo = uRepo[0]
	}
	return &service{repo: r, userRepo: userRepo}
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

func (s *service) GetMaterialsByClass(ctx context.Context, userID, role, classID string) ([]MaterialResponse, error) {
	// Role and Class assignment checks
	if role == "student" {
		if s.userRepo != nil {
			student, err := s.userRepo.GetByID(ctx, userID)
			if err != nil {
				return nil, err
			}
			if student.ClassID == nil || *student.ClassID != classID {
				return nil, errors.New("unauthorized: student does not belong to this class")
			}
		}
	} else if role == "parent" {
		if s.userRepo != nil {
			children, err := s.userRepo.GetChildren(ctx, userID)
			if err != nil {
				return nil, err
			}
			authorized := false
			for _, child := range children {
				if child.ClassID == classID {
					authorized = true
					break
				}
			}
			if !authorized {
				return nil, errors.New("unauthorized: parent does not have any children in this class")
			}
		}
	} else if role != "teacher" && role != "admin" && role != "superadmin" {
		return nil, errors.New("unauthorized: invalid role")
	}

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
