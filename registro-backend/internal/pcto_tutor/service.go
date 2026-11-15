package pcto_tutor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func hashTutorToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *Service) AuthenticateWithToken(ctx context.Context, token string) (*CompanyTutor, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("token di accesso obbligatorio")
	}
	return s.repo.GetTutorByToken(ctx, hashTutorToken(token))
}

func (s *Service) RegisterTutor(ctx context.Context, tutor *CompanyTutor) error {
	if tutor.Email == "" || tutor.CompanyName == "" || tutor.TutorLastName == "" {
		return errors.New("email, azienda e cognome tutor sono obbligatori")
	}

	token, expiresAt, err := GenerateMagicLinkToken(14 * 24 * time.Hour)
	if err != nil {
		return err
	}
	tutor.TokenExpiresAt = &expiresAt
	tutor.IsActive = true
	tutor.AccessToken = hashTutorToken(token)
	err = s.repo.CreateTutor(ctx, tutor)
	// The raw token is returned once to the authorized registrar; only its
	// SHA-256 digest is persisted and subsequent reads never return it.
	tutor.AccessToken = token
	return err
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
