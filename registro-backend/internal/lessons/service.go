package lessons

import (
	"fmt"
	"time"
)

type Service interface {
	CreateLesson(teacherID string, req CreateLessonRequest) (*LessonResponse, error)
	GetLessons(classID, subjectID string) ([]LessonResponse, error)
	CreateHomework(teacherID string, req CreateHomeworkRequest) (*HomeworkResponse, error)
	GetHomeworks(classID string) ([]HomeworkResponse, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateLesson(teacherID string, req CreateLessonRequest) (*LessonResponse, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	lesson := &Lesson{
		ClassID:   req.ClassID,
		SubjectID: req.SubjectID,
		TeacherID: teacherID,
		Date:      date,
		Topic:     req.Topic,
		Type:      req.Type,
		Notes:     req.Notes,
	}

	if err := s.repo.CreateLesson(lesson); err != nil {
		return nil, err
	}

	return s.mapLessonResponse(lesson), nil
}

func (s *service) GetLessons(classID, subjectID string) ([]LessonResponse, error) {
	var lessons []Lesson
	var err error

	if subjectID != "" {
		lessons, err = s.repo.GetLessonsByClassAndSubject(classID, subjectID)
	} else {
		lessons, err = s.repo.GetLessonsByClass(classID)
	}

	if err != nil {
		return nil, err
	}

	var res []LessonResponse
	for _, l := range lessons {
		res = append(res, *s.mapLessonResponse(&l))
	}
	return res, nil
}

func (s *service) CreateHomework(teacherID string, req CreateHomeworkRequest) (*HomeworkResponse, error) {
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid due_date format: %w", err)
	}

	hw := &Homework{
		ClassID:     req.ClassID,
		SubjectID:   req.SubjectID,
		TeacherID:   teacherID,
		LessonID:    req.LessonID,
		DueDate:     dueDate,
		Description: req.Description,
	}

	if err := s.repo.CreateHomework(hw); err != nil {
		return nil, err
	}

	return s.mapHomeworkResponse(hw), nil
}

func (s *service) GetHomeworks(classID string) ([]HomeworkResponse, error) {
	homeworks, err := s.repo.GetHomeworkByClass(classID)
	if err != nil {
		return nil, err
	}

	var res []HomeworkResponse
	for _, h := range homeworks {
		res = append(res, *s.mapHomeworkResponse(&h))
	}
	return res, nil
}

func (s *service) mapLessonResponse(l *Lesson) *LessonResponse {
	return &LessonResponse{
		ID:        l.ID,
		ClassID:   l.ClassID,
		TeacherID: l.TeacherID,
		SubjectID: l.SubjectID,
		Date:      l.Date,
		Topic:     l.Topic,
		Type:      l.Type,
		Notes:     l.Notes,
	}
}

func (s *service) mapHomeworkResponse(h *Homework) *HomeworkResponse {
	return &HomeworkResponse{
		ID:          h.ID,
		LessonID:    h.LessonID,
		ClassID:     h.ClassID,
		SubjectID:   h.SubjectID,
		TeacherID:   h.TeacherID,
		DueDate:     h.DueDate,
		Description: h.Description,
	}
}
