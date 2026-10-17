package middleschoolexam

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

func (s *Service) GetOrCreateExam(ctx context.Context, schoolID, classID, academicYear, president string) (*MiddleSchoolExam, error) {
	if classID == "" {
		return nil, errors.New("class_id is required")
	}
	if academicYear == "" {
		academicYear = "2026/2027"
	}
	return s.repo.GetOrCreateExamForClass(ctx, schoolID, classID, academicYear, president)
}

func (s *Service) ListCandidates(ctx context.Context, examID string) ([]MiddleSchoolExamCandidate, error) {
	return s.repo.ListCandidates(ctx, examID)
}

func (s *Service) SaveAdmission(ctx context.Context, examID, studentID string, grade int, judgment string, admitted bool) error {
	if grade < 6 || grade > 10 {
		return errors.New("voto di ammissione deve essere compreso tra 6 e 10")
	}
	return s.repo.SaveCandidateAdmission(ctx, examID, studentID, grade, judgment, admitted)
}

func (s *Service) EvaluateCandidate(ctx context.Context, examID, studentID string, grades ExamCandidateGrades, notes string) (*MiddleSchoolExamCandidate, error) {
	candidate, err := s.repo.GetCandidate(ctx, examID, studentID)
	if err != nil {
		return nil, errors.New("candidate not found")
	}

	grades.AdmissionGrade = candidate.AdmissionGrade
	calc := CalculateExamOutcome(grades)

	candidate.GradeItalian = grades.GradeItalian
	candidate.GradeMath = grades.GradeMath
	candidate.GradeEnglish = grades.GradeEnglish
	candidate.GradeSecondLang = grades.GradeSecondLang
	candidate.GradeInterview = grades.GradeInterview
	candidate.ExamMean = calc.ExamMean
	candidate.FinalGrade = calc.FinalGrade
	candidate.HasHonors = calc.HasHonors
	candidate.Outcome = calc.Outcome
	candidate.Notes = notes

	if err := s.repo.SaveCandidateGrades(ctx, candidate); err != nil {
		return nil, err
	}

	return candidate, nil
}

func (s *Service) GenerateDiploma(ctx context.Context, candidateID string) (string, error) {
	diplomaData, err := s.repo.GetCandidateDiplomaData(ctx, candidateID)
	if err != nil {
		return "", err
	}
	return GenerateDiplomaCertificateText(*diplomaData), nil
}

func (s *Service) UpdateStatus(ctx context.Context, examID, status string) error {
	return s.repo.UpdateExamStatus(ctx, examID, status)
}
