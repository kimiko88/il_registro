package religion_alternative

import (
	"context"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RecordOption(ctx context.Context, opt *StudentReligionOption) error {
	if opt.StudentID == "" {
		return errors.New("student_id obbligatorio")
	}
	if opt.SchoolID == "" {
		return errors.New("school_id obbligatorio")
	}
	if opt.AcademicYear == "" {
		opt.AcademicYear = "2026/2027"
	}
	if !IsValidOptionType(opt.OptionType) {
		return errors.New("opzione non valida secondo le tipologie ministeriali")
	}
	return s.repo.SaveOption(ctx, opt)
}

func (s *Service) GetOption(ctx context.Context, studentID, academicYear string) (*StudentReligionOption, error) {
	if studentID == "" {
		return nil, errors.New("student_id obbligatorio")
	}
	if academicYear == "" {
		academicYear = "2026/2027"
	}
	return s.repo.GetOption(ctx, studentID, academicYear)
}

func (s *Service) ListOptions(ctx context.Context, schoolID, academicYear string) ([]StudentReligionOption, error) {
	if schoolID == "" {
		return nil, errors.New("school_id obbligatorio")
	}
	if academicYear == "" {
		academicYear = "2026/2027"
	}
	return s.repo.ListOptionsBySchool(ctx, schoolID, academicYear)
}

func (s *Service) RecordEvaluation(ctx context.Context, eval *AlternativeEvaluation) error {
	if eval.StudentID == "" {
		return errors.New("student_id obbligatorio")
	}
	if eval.SchoolID == "" {
		return errors.New("school_id obbligatorio")
	}
	if eval.Period == "" {
		return errors.New("periodo di valutazione obbligatorio (es. q1, q2, finale)")
	}
	if eval.SubjectKind != "irc" && eval.SubjectKind != "materia_alternativa" {
		return errors.New("subject_kind deve essere 'irc' oppure 'materia_alternativa'")
	}
	if !IsValidJudgmentLevel(eval.JudgmentLevel) {
		return errors.New("livello di giudizio non valido (attesi: ottimo, distinto, buono, sufficiente, non_sufficiente)")
	}
	return s.repo.SaveEvaluation(ctx, eval)
}

func (s *Service) ListEvaluations(ctx context.Context, schoolID, period, subjectKind string) ([]AlternativeEvaluation, error) {
	if schoolID == "" {
		return nil, errors.New("school_id obbligatorio")
	}
	return s.repo.ListEvaluations(ctx, schoolID, period, subjectKind)
}
