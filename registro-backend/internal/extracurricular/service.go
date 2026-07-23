package extracurricular

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized action on extracurricular course")
	ErrCourseFull   = errors.New("extracurricular course has reached maximum capacity")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("extracurricular.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) CreateCourse(ctx context.Context, teacherID, schoolID string, req CreateCourseRequest) (*Course, error) {
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
	if course.EnrolledCount >= course.MaxParticipants {
		return ErrCourseFull
	}
	return s.repo.EnrollStudent(ctx, courseID, studentID)
}

func (s *Service) ListEnrollments(ctx context.Context, courseID string) ([]*Enrollment, error) {
	return s.repo.ListEnrollments(ctx, courseID)
}

func (s *Service) MarkAttendance(ctx context.Context, teacherID, role string, req MarkAttendanceRequest) error {
	if role != "teacher" && role != "admin" && role != "superadmin" {
		return ErrUnauthorized
	}
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: YYYY-MM-DD")
	}

	att := &AttendanceRecord{
		CourseID:  req.CourseID,
		StudentID: req.StudentID,
		Date:      d,
		Status:    req.Status,
		Hours:     req.Hours,
		SignedBy:  &teacherID,
	}

	return s.repo.MarkAttendance(ctx, att)
}
