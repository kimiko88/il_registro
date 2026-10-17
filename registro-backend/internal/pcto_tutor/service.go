package pcto_tutor

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AuthenticateWithToken(ctx context.Context, token string) (*CompanyTutor, error) {
	if token == "" {
		return nil, errors.New("token di accesso obbligatorio")
	}
	return s.repo.GetTutorByToken(ctx, token)
}

func (s *Service) RegisterTutor(ctx context.Context, tutor *CompanyTutor) error {
	if tutor.Email == "" || tutor.CompanyName == "" || tutor.TutorLastName == "" {
		return errors.New("email, azienda e cognome tutor sono obbligatori")
	}

	token, expiresAt, err := GenerateMagicLinkToken(14 * 24 * time.Hour) // 14 days valid
	if err != nil {
		return err
	}
	tutor.AccessToken = token
	tutor.TokenExpiresAt = &expiresAt
	tutor.IsActive = true

	return s.repo.CreateTutor(ctx, tutor)
}

func (s *Service) AssignStudent(ctx context.Context, tutorID, projectID, studentID string) error {
	if tutorID == "" || projectID == "" || studentID == "" {
		return errors.New("tutor_id, project_id e student_id sono obbligatori")
	}
	return s.repo.AssignStudent(ctx, tutorID, projectID, studentID)
}

func (s *Service) GetAssignedStudents(ctx context.Context, tutorID string) ([]AssignedStudentInfo, error) {
	if tutorID == "" {
		return nil, errors.New("tutor_id obbligatorio")
	}
	return s.repo.ListAssignedStudents(ctx, tutorID)
}

func (s *Service) VerifyTimesheet(ctx context.Context, entry *TimesheetVerification) error {
	if entry.TutorID == "" || entry.StudentID == "" || entry.ProjectID == "" {
		return errors.New("tutor_id, student_id e project_id obbligatori")
	}
	if err := ValidateTimesheetEntry(entry); err != nil {
		return err
	}
	return s.repo.SaveTimesheetVerification(ctx, entry)
}

func (s *Service) SubmitEvaluation(ctx context.Context, eval *CompanyEvaluation) error {
	if eval.TutorID == "" || eval.StudentID == "" || eval.ProjectID == "" {
		return errors.New("tutor_id, student_id e project_id obbligatori")
	}
	if err := ValidateEvaluation(eval); err != nil {
		return err
	}
	return s.repo.SaveCompanyEvaluation(ctx, eval)
}
