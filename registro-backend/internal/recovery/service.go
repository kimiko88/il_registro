package recovery

import (
	"context"
	"errors"
)

var (
	ErrCourseNotFound = errors.New("corso di recupero non trovato")
	ErrInvalidGrade   = errors.New("voto della prova di recupero fuori dai limiti (1.0 - 10.0)")
)

type Service interface {
	CreateCourse(ctx context.Context, schoolID string, req CreateCourseRequest) (*RecoveryCourse, error)
	GetCourse(ctx context.Context, id string) (*RecoveryCourse, error)
	ListCourses(ctx context.Context, schoolID, academicYear, teacherID string) ([]RecoveryCourse, error)
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateAttendance(ctx context.Context, courseID, studentID string, hours float64, notes string) error

	RecordTestOutcome(ctx context.Context, schoolID, teacherID string, req RecordTestOutcomeRequest) (*RecoveryTest, error)
	ListTests(ctx context.Context, schoolID, classID, studentID string) ([]RecoveryTest, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateCourse(ctx context.Context, schoolID string, req CreateCourseRequest) (*RecoveryCourse, error) {
	c := &RecoveryCourse{
		SchoolID:     schoolID,
		SubjectID:    req.SubjectID,
		TeacherID:    req.TeacherID,
		Title:        req.Title,
		Description:  req.Description,
		AcademicYear: req.AcademicYear,
		Period:       req.Period,
		TotalHours:   req.TotalHours,
		Room:         req.Room,
		Status:       "scheduled",
	}
	if c.Period == "" {
		c.Period = "summer"
	}
	if c.TotalHours <= 0 {
		c.TotalHours = 10
	}

	if err := s.repo.CreateCourse(ctx, c, req.Sessions, req.StudentIDs); err != nil {
		return nil, err
	}

	return s.repo.GetCourseByID(ctx, c.ID)
}

func (s *service) GetCourse(ctx context.Context, id string) (*RecoveryCourse, error) {
	c, err := s.repo.GetCourseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCourseNotFound
	}
	return c, nil
}

func (s *service) ListCourses(ctx context.Context, schoolID, academicYear, teacherID string) ([]RecoveryCourse, error) {
	return s.repo.ListCourses(ctx, schoolID, academicYear, teacherID)
}

func (s *service) UpdateStatus(ctx context.Context, id, status string) error {
	return s.repo.UpdateCourseStatus(ctx, id, status)
}

func (s *service) UpdateAttendance(ctx context.Context, courseID, studentID string, hours float64, notes string) error {
	return s.repo.UpdateStudentAttendance(ctx, courseID, studentID, hours, notes)
}

func (s *service) RecordTestOutcome(ctx context.Context, schoolID, teacherID string, req RecordTestOutcomeRequest) (*RecoveryTest, error) {
	if req.Grade < 1.0 || req.Grade > 10.0 {
		return nil, ErrInvalidGrade
	}

	outcome := "recuperato"
	deliberation := "Ammesso"
	if req.Grade < 6.0 {
		outcome = "non_recuperato"
		deliberation = "Non Ammesso"
	}
	if req.FinalDeliberation != "" {
		deliberation = req.FinalDeliberation
	}

	test := &RecoveryTest{
		SchoolID:          schoolID,
		DeficiencyID:      req.DeficiencyID,
		StudentID:         req.StudentID,
		SubjectID:         req.SubjectID,
		ClassID:           req.ClassID,
		TeacherID:         teacherID,
		TestDate:          req.TestDate,
		TestType:          req.TestType,
		Grade:             req.Grade,
		Outcome:           outcome,
		FinalDeliberation: deliberation,
		VerbaleNumber:     req.VerbaleNumber,
		Notes:             req.Notes,
	}

	if err := s.repo.RecordRecoveryTest(ctx, test); err != nil {
		return nil, err
	}

	return test, nil
}

func (s *service) ListTests(ctx context.Context, schoolID, classID, studentID string) ([]RecoveryTest, error) {
	return s.repo.ListRecoveryTests(ctx, schoolID, classID, studentID)
}
