package credits

import (
	"context"
	"errors"
	"fmt"

	"registro-backend/internal/users"
)

var (
	ErrInvalidGradeLevel = errors.New("livello classe non valido per crediti scolastici (richiesto 3, 4 o 5)")
	ErrCreditOutOfRange  = errors.New("credito assegnato fuori dalla fascia ministeriale consentita")
)

type Service interface {
	CalculateSuggestedCredit(gradeLevel int, average float64, conductGrade int, pctoHours int, hasExtracurricular bool) CreditCalculationResult
	AssignCredit(ctx context.Context, actorID, schoolID string, req AssignCreditRequest) (*StudentSchoolCredit, error)
	ListClassCredits(ctx context.Context, classID, academicYear string) ([]StudentSchoolCredit, error)
	GetStudentSummary(ctx context.Context, studentID string) (*StudentCreditSummary, error)
	GetUserRepo() users.Repository
}

type service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(repo Repository, uRepo ...users.Repository) Service {
	svc := &service{repo: repo}
	if len(uRepo) > 0 && uRepo[0] != nil {
		svc.userRepo = uRepo[0]
	}
	return svc
}

func (s *service) GetUserRepo() users.Repository {
	return s.userRepo
}


func (s *service) CalculateSuggestedCredit(gradeLevel int, average float64, conductGrade int, pctoHours int, hasExtracurricular bool) CreditCalculationResult {
	return CalculateCreditRange(gradeLevel, average, conductGrade, pctoHours, hasExtracurricular)
}

func (s *service) AssignCredit(ctx context.Context, actorID, schoolID string, req AssignCreditRequest) (*StudentSchoolCredit, error) {
	if req.GradeLevel < 3 || req.GradeLevel > 5 {
		return nil, ErrInvalidGradeLevel
	}

	calc := CalculateCreditRange(req.GradeLevel, req.GradeAverage, req.ConductGrade, req.PCTOHours, req.HasExtracurricular)

	// Validate that assigned credit is within range or justified
	if req.AssignedCredit < calc.BaseCreditRangeMin || req.AssignedCredit > calc.BaseCreditRangeMax {
		if req.DeliberationNotes == "" {
			return nil, fmt.Errorf("%w (%d - %d), motivazione delibera obbligatoria per eccezioni", ErrCreditOutOfRange, calc.BaseCreditRangeMin, calc.BaseCreditRangeMax)
		}
	}

	credit := &StudentSchoolCredit{
		SchoolID:           schoolID,
		StudentID:          req.StudentID,
		ClassID:            req.ClassID,
		AcademicYear:       req.AcademicYear,
		GradeLevel:         req.GradeLevel,
		GradeAverage:       req.GradeAverage,
		ConductGrade:       req.ConductGrade,
		BaseCreditRangeMin: calc.BaseCreditRangeMin,
		BaseCreditRangeMax: calc.BaseCreditRangeMax,
		AssignedCredit:     req.AssignedCredit,
		PCTOHours:          req.PCTOHours,
		HasExtracurricular: req.HasExtracurricular,
		DeliberationNotes:  req.DeliberationNotes,
		ValidatedBy:        &actorID,
	}

	if err := s.repo.SaveCredit(ctx, credit); err != nil {
		return nil, err
	}

	return credit, nil
}

func (s *service) ListClassCredits(ctx context.Context, classID, academicYear string) ([]StudentSchoolCredit, error) {
	return s.repo.ListCreditsByClass(ctx, classID, academicYear)
}

func (s *service) GetStudentSummary(ctx context.Context, studentID string) (*StudentCreditSummary, error) {
	return s.repo.GetStudentCreditSummary(ctx, studentID)
}
