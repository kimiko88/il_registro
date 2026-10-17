package primaryeval

import (
	"context"
	"errors"
	"fmt"
)

type Service interface {
	CreateObjective(ctx context.Context, schoolID string, req *CreateObjectiveRequest) (*LearningObjective, error)
	ListObjectives(ctx context.Context, schoolID, subjectID string, classID *string, yearGrade int, academicYear string) ([]LearningObjective, error)
	DeleteObjective(ctx context.Context, id string) error

	SaveEvaluationsBatch(ctx context.Context, schoolID, teacherID string, req *SaveEvaluationsBatchRequest) error
	ListEvaluations(ctx context.Context, classID, subjectID string, semester int) ([]PrimaryEvaluation, error)
	GetStudentEvaluations(ctx context.Context, studentID string, semester int) ([]PrimaryEvaluation, error)
	GetPrimaryMatrix(ctx context.Context, schoolID, classID, subjectID string, semester int) (*PrimaryMatrixResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

var (
	ErrInvalidLevel     = errors.New("livello non valido: deve essere avanzato, intermedio, base o in_via_di_prima_acquisizione")
	ErrInvalidYearGrade = errors.New("anno di corso non valido: la primaria comprende gli anni da 1 a 5")
	ErrMissingTitle     = errors.New("titolo dell'obiettivo didattico obbligatorio")
)

func IsValidPrimaryLevel(level string) bool {
	switch level {
	case LevelAvanzato, LevelIntermedio, LevelBase, LevelInViaPrimaAcquisiz:
		return true
	default:
		return false
	}
}

func (s *service) CreateObjective(ctx context.Context, schoolID string, req *CreateObjectiveRequest) (*LearningObjective, error) {
	if req.Title == "" {
		return nil, ErrMissingTitle
	}
	if req.YearGrade < 1 || req.YearGrade > 5 {
		return nil, ErrInvalidYearGrade
	}
	if req.AcademicYear == "" {
		req.AcademicYear = "2025/2026"
	}

	obj := &LearningObjective{
		SchoolID:     schoolID,
		ClassID:      req.ClassID,
		SubjectID:    req.SubjectID,
		YearGrade:    req.YearGrade,
		Title:        req.Title,
		Description:  req.Description,
		AcademicYear: req.AcademicYear,
	}

	return s.repo.CreateObjective(ctx, obj)
}

func (s *service) ListObjectives(ctx context.Context, schoolID, subjectID string, classID *string, yearGrade int, academicYear string) ([]LearningObjective, error) {
	return s.repo.ListObjectives(ctx, schoolID, subjectID, classID, yearGrade, academicYear)
}

func (s *service) DeleteObjective(ctx context.Context, id string) error {
	return s.repo.DeleteObjective(ctx, id)
}

func (s *service) SaveEvaluationsBatch(ctx context.Context, schoolID, teacherID string, req *SaveEvaluationsBatchRequest) error {
	if req.ClassID == "" || req.SubjectID == "" || req.ObjectiveID == "" {
		return errors.New("classe, materia e obiettivo didattico sono campi obbligatori")
	}
	if req.Semester != 1 && req.Semester != 2 {
		return errors.New("periodo di valutazione non valido (deve essere 1° o 2° quadrimestre)")
	}

	for _, item := range req.Evaluations {
		if !IsValidPrimaryLevel(item.Level) {
			return fmt.Errorf("%w: ricevuto '%s'", ErrInvalidLevel, item.Level)
		}
	}

	return s.repo.SaveEvaluationsBatch(ctx, schoolID, teacherID, req)
}

func (s *service) ListEvaluations(ctx context.Context, classID, subjectID string, semester int) ([]PrimaryEvaluation, error) {
	return s.repo.ListEvaluationsByClassAndSubject(ctx, classID, subjectID, semester)
}

func (s *service) GetStudentEvaluations(ctx context.Context, studentID string, semester int) ([]PrimaryEvaluation, error) {
	return s.repo.GetStudentEvaluations(ctx, studentID, semester)
}

func (s *service) GetPrimaryMatrix(ctx context.Context, schoolID, classID, subjectID string, semester int) (*PrimaryMatrixResponse, error) {
	objectives, err := s.repo.ListObjectives(ctx, schoolID, subjectID, &classID, 0, "")
	if err != nil {
		return nil, err
	}

	evals, err := s.repo.ListEvaluationsByClassAndSubject(ctx, classID, subjectID, semester)
	if err != nil {
		return nil, err
	}

	// Map students
	studentMap := make(map[string]*PrimaryStudentMatrixRow)
	var studentOrder []string

	for _, e := range evals {
		row, exists := studentMap[e.StudentID]
		if !exists {
			row = &PrimaryStudentMatrixRow{
				StudentID:   e.StudentID,
				StudentName: e.StudentName,
				Evaluations: make(map[string]PrimaryEvaluationCell),
			}
			studentMap[e.StudentID] = row
			studentOrder = append(studentOrder, e.StudentID)
		}

		// Keep most recent evaluation per objective
		if _, hasObj := row.Evaluations[e.ObjectiveID]; !hasObj {
			row.Evaluations[e.ObjectiveID] = PrimaryEvaluationCell{
				Level:                e.Level,
				DimensionAutonomy:    e.DimensionAutonomy,
				DimensionContinuity:  e.DimensionContinuity,
				DimensionFamiliarity: e.DimensionFamiliarity,
				DimensionResources:   e.DimensionResources,
				Notes:                e.Notes,
				Date:                 e.Date,
			}
		}
	}

	studentsList := make([]PrimaryStudentMatrixRow, 0, len(studentOrder))
	for _, id := range studentOrder {
		studentsList = append(studentsList, *studentMap[id])
	}

	return &PrimaryMatrixResponse{
		ClassID:    classID,
		SubjectID:  subjectID,
		Semester:   semester,
		Objectives: objectives,
		Students:   studentsList,
	}, nil
}
