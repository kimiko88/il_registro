package rubrics

import (
	"context"
	"errors"
	"fmt"
	"time"

	"registro-backend/internal/users"
)

var (
	ErrRubricNotFound = errors.New("rubric not found")
	ErrUnauthorized   = errors.New("unauthorized action on rubric")
)

type Service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(repo Repository, uRepo ...users.Repository) *Service {
	if repo == nil {
		panic("rubrics.NewService: repo must not be nil")
	}
	s := &Service{repo: repo}
	if len(uRepo) > 0 && uRepo[0] != nil {
		s.userRepo = uRepo[0]
	}
	return s
}

func (s *Service) GetUserRepo() users.Repository {
	return s.userRepo
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

// UpdateRubric updates a rubric, verifying ownership or admin school membership.
// Bug 138: admin can only update rubrics belonging to their own school.
func (s *Service) UpdateRubric(ctx context.Context, actorID, actorRole, actorSchoolID, id string, req CreateRubricRequest) (*Rubric, error) {
	rub, err := s.repo.GetRubricByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rub.TeacherID != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return nil, ErrUnauthorized
	}
	// Bug 138: admin can only edit rubrics within their own school
	if actorRole == "admin" && actorSchoolID != "" && rub.SchoolID != actorSchoolID {
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

// DeleteRubric deletes a rubric, verifying ownership or admin school membership.
// Bug 138: admin can only delete rubrics belonging to their own school.
func (s *Service) DeleteRubric(ctx context.Context, actorID, actorRole, actorSchoolID, id string) error {
	rub, err := s.repo.GetRubricByID(ctx, id)
	if err != nil {
		return err
	}
	if rub.TeacherID != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	// Bug 138: admin can only delete rubrics within their own school
	if actorRole == "admin" && actorSchoolID != "" && rub.SchoolID != actorSchoolID {
		return ErrUnauthorized
	}
	return s.repo.DeleteRubric(ctx, id)
}

// AssessStudent creates a rubric assessment for a student.
// Bug 137: verifying teacher-class assignment requires a classRepo or teacherRepo dependency.
// TODO: add IsTeacherAssignedToClass(ctx, teacherID, classID) to the Repository and call it here.
func (s *Service) AssessStudent(ctx context.Context, teacherID, rubricID string, req CreateAssessmentRequest) (*RubricAssessment, error) {
	var totalScore float64
	for _, cs := range req.Scores {
		if cs.Score < 0 {
			return nil, errors.New("score cannot be negative")
		}
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
