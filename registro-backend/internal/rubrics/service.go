package rubrics

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrRubricNotFound = errors.New("rubric not found")
	ErrUnauthorized   = errors.New("unauthorized action on rubric")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("rubrics.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) CreateRubric(ctx context.Context, schoolID, teacherID string, req CreateRubricRequest) (*Rubric, error) {
	if schoolID == "" || teacherID == "" {
		return nil, fmt.Errorf("school_id and teacher_id required")
	}

	rub := &Rubric{
		SchoolID:    schoolID,
		TeacherID:   teacherID,
		SubjectID:   req.SubjectID,
		Title:       req.Title,
		Description: req.Description,
	}

	for _, ci := range req.Criteria {
		crit := Criterion{
			Name:        ci.Name,
			Description: ci.Description,
			MaxScore:    ci.MaxScore,
		}
		if crit.MaxScore <= 0 {
			crit.MaxScore = 10
		}

		for _, li := range ci.Levels {
			crit.Levels = append(crit.Levels, Level{
				Score:       li.Score,
				Label:       li.Label,
				Description: li.Description,
			})
		}
		rub.Criteria = append(rub.Criteria, crit)
	}

	if err := s.repo.CreateRubric(ctx, rub); err != nil {
		return nil, err
	}

	return rub, nil
}

func (s *Service) GetRubric(ctx context.Context, id string) (*Rubric, error) {
	return s.repo.GetRubricByID(ctx, id)
}

func (s *Service) ListRubrics(ctx context.Context, schoolID, teacherID string) ([]*Rubric, error) {
	return s.repo.ListRubrics(ctx, schoolID, teacherID)
}

func (s *Service) UpdateRubric(ctx context.Context, actorID, actorRole, id string, req CreateRubricRequest) (*Rubric, error) {
	rub, err := s.repo.GetRubricByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rub.TeacherID != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return nil, ErrUnauthorized
	}

	rub.Title = req.Title
	rub.Description = req.Description
	rub.SubjectID = req.SubjectID

	if err := s.repo.UpdateRubric(ctx, rub); err != nil {
		return nil, err
	}
	return rub, nil
}

func (s *Service) DeleteRubric(ctx context.Context, actorID, actorRole, id string) error {
	rub, err := s.repo.GetRubricByID(ctx, id)
	if err != nil {
		return err
	}
	if rub.TeacherID != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	return s.repo.DeleteRubric(ctx, id)
}

func (s *Service) AssessStudent(ctx context.Context, teacherID, rubricID string, req CreateAssessmentRequest) (*RubricAssessment, error) {
	var totalScore float64
	for _, cs := range req.Scores {
		totalScore += cs.Score
	}

	var d time.Time
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02", req.Date)
		if err == nil {
			d = parsed
		}
	}
	if d.IsZero() {
		d = time.Now()
	}

	assessment := &RubricAssessment{
		RubricID:   rubricID,
		StudentID:  req.StudentID,
		ClassID:    req.ClassID,
		TeacherID:  teacherID,
		Date:       d,
		Scores:     req.Scores,
		TotalScore: totalScore,
		Notes:      req.Notes,
	}

	if err := s.repo.CreateAssessment(ctx, assessment); err != nil {
		return nil, err
	}
	return assessment, nil
}

func (s *Service) ListAssessmentsByStudent(ctx context.Context, studentID string) ([]*RubricAssessment, error) {
	return s.repo.ListAssessmentsByStudent(ctx, studentID)
}

func (s *Service) ListAssessmentsByClass(ctx context.Context, classID string) ([]*RubricAssessment, error) {
	return s.repo.ListAssessmentsByClass(ctx, classID)
}
