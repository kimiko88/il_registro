package extracurricular

import (
	"context"
	"errors"
	"fmt"
	"time"

	"registro-backend/internal/users"
)

var (
	ErrUnauthorized = errors.New("unauthorized action on extracurricular course")
	ErrCourseFull   = errors.New("extracurricular course has reached maximum capacity")
)

type Service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(repo Repository, uRepo ...users.Repository) *Service {
	if repo == nil {
		panic("extracurricular.NewService: repo must not be nil")
	}
	var userRepo users.Repository
	if len(uRepo) > 0 {
		userRepo = uRepo[0]
	}
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) CreateCourse(ctx context.Context, actorRole, teacherID, schoolID string, req CreateCourseRequest) (*Course, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "teacher" {
		return nil, ErrUnauthorized
	}
	if schoolID == "" {
		return nil, fmt.Errorf("school_id required")
	}
	st, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: YYYY-MM-DD")
	}
	et, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format: YYYY-MM-DD")
	}
	if !et.After(st) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}

	c := &Course{
		SchoolID:        schoolID,
		Title:           req.Title,
		Description:     req.Description,
		TeacherID:       teacherID,
		StartDate:       st,
		EndDate:         et,
		MaxParticipants: req.MaxParticipants,
	}

	if err := s.repo.CreateCourse(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListCourses(ctx context.Context, schoolID, studentID string) ([]*Course, error) {
	return s.repo.ListCourses(ctx, schoolID, studentID)
}

func (s *Service) EnrollStudent(ctx context.Context, courseID, studentID string) error {
	course, err := s.repo.GetCourseByID(ctx, courseID)
	if err != nil {
		return err
	}

	if course.MaxParticipants > 0 && course.EnrolledCount >= course.MaxParticipants {
		return ErrCourseFull
	}

	if s.userRepo != nil {
		student, err := s.userRepo.GetByID(ctx, studentID)
		if err != nil || student == nil {
			return fmt.Errorf("student non trovato")
		}
		if student.SchoolID != nil && *student.SchoolID != course.SchoolID {
			return fmt.Errorf("unauthorized: lo studente appartiene ad un'altra scuola")
		}
	}

	return s.repo.EnrollStudent(ctx, courseID, studentID)
}

func (s *Service) ListEnrollments(ctx context.Context, courseID string) ([]*Enrollment, error) {
	return s.repo.ListEnrollments(ctx, courseID)
}

func (s *Service) MarkAttendance(ctx context.Context, teacherID, actorRole string, req MarkAttendanceRequest) error {
	if actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	course, err := s.repo.GetCourseByID(ctx, req.CourseID)
	if err != nil {
		return err
	}
	if actorRole == "teacher" && course.TeacherID != teacherID {
		return errors.New("unauthorized: non sei il docente responsabile di questo corso")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return fmt.Errorf("invalid date format (expected YYYY-MM-DD)")
	}

	att := &AttendanceRecord{
		CourseID:  req.CourseID,
		StudentID: req.StudentID,
		Date:      date,
		Status:    req.Status,
		Hours:     req.Hours,
		SignedBy:  &teacherID,
	}

	return s.repo.MarkAttendance(ctx, att)
}
