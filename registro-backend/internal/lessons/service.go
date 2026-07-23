package lessons

import (
	"errors"
	"fmt"
	"time"
)

type Service interface {
	CreateLesson(teacherID string, req CreateLessonRequest) (*LessonResponse, error)
	GetLessonByID(id string) (*LessonResponse, error)
	UpdateLesson(teacherID, role, id string, req UpdateLessonRequest) (*LessonResponse, error)
	DeleteLesson(teacherID, role, id string) error
	GetLessons(classID, subjectID string, date string) ([]LessonResponse, error)
	GetLessonsByGroup(groupID string, date string) ([]LessonResponse, error)
	GetTeacherDiary(teacherID string, fromDate, toDate string) ([]LessonResponse, error)
	CreateHomework(teacherID string, req CreateHomeworkRequest) (*HomeworkResponse, error)
	UpdateHomework(teacherID, role, id string, req UpdateHomeworkRequest) (*HomeworkResponse, error)
	DeleteHomework(teacherID, role, id string) error
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

	activityType := req.ActivityType
	if activityType == "" {
		if req.IsSubstitution {
			activityType = "substitution"
		} else {
			activityType = "standard"
		}
	}

	lesson := &Lesson{
		ClassID:              req.ClassID,
		SubjectID:            req.SubjectID,
		TeacherID:            teacherID,
		Date:                 date,
		Hour:                 req.Hour,
		Duration:             req.Duration,
		Topic:                req.Topic,
		Type:                 req.Type,
		GroupID:              req.GroupID,
		IsSubstitution:       req.IsSubstitution,
		SubstitutedTeacherID: req.SubstitutedTeacherID,
		ActivityType:         activityType,
		Notes:                req.Notes,
	}

	if err := s.repo.CreateLesson(lesson); err != nil {
		return nil, err
	}

	return s.mapLessonResponse(lesson), nil
}

func (s *service) GetLessons(classID, subjectID string, date string) ([]LessonResponse, error) {
	var lessons []Lesson
	var err error

	if subjectID != "" {
		lessons, err = s.repo.GetLessonsByClassAndSubject(classID, subjectID, date)
	} else {
		lessons, err = s.repo.GetLessonsByClass(classID, date)
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

func (s *service) GetLessonsByGroup(groupID string, date string) ([]LessonResponse, error) {
	lessons, err := s.repo.GetLessonsByGroup(groupID, date)
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
		Type:        req.Type,
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
		ID:                     l.ID,
		ClassID:                l.ClassID,
		TeacherID:              l.TeacherID,
		TeacherName:            l.TeacherName,
		SubjectID:              l.SubjectID,
		Date:                   l.Date,
		Hour:                   l.Hour,
		Duration:               l.Duration,
		Topic:                  l.Topic,
		Type:                   l.Type,
		GroupID:                l.GroupID,
		IsSubstitution:         l.IsSubstitution,
		SubstitutedTeacherID:   l.SubstitutedTeacherID,
		SubstitutedTeacherName: l.SubstitutedTeacherName,
		ActivityType:           l.ActivityType,
		Notes:                  l.Notes,
	}
}

func (s *service) mapHomeworkResponse(h *Homework) *HomeworkResponse {
	return &HomeworkResponse{
		ID:          h.ID,
		LessonID:    h.LessonID,
		ClassID:     h.ClassID,
		SubjectID:   h.SubjectID,
		TeacherID:   h.TeacherID,
		TeacherName: h.TeacherName,
		DueDate:     h.DueDate,
		Description: h.Description,
		Type:        h.Type,
	}
}

func (s *service) GetLessonByID(id string) (*LessonResponse, error) {
	l, err := s.repo.GetLessonByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapLessonResponse(l), nil
}

func (s *service) UpdateLesson(teacherID, role, id string, req UpdateLessonRequest) (*LessonResponse, error) {
	existing, err := s.repo.GetLessonByID(id)
	if err != nil {
		return nil, err
	}
	if existing.TeacherID != teacherID && role != "admin" && role != "superadmin" {
		return nil, errors.New("unauthorized: cannot edit another teacher's lesson")
	}

	l, err := s.repo.UpdateLesson(id, req)
	if err != nil {
		return nil, err
	}
	return s.mapLessonResponse(l), nil
}

func (s *service) DeleteLesson(teacherID, role, id string) error {
	existing, err := s.repo.GetLessonByID(id)
	if err != nil {
		return err
	}
	if existing.TeacherID != teacherID && role != "admin" && role != "superadmin" {
		return errors.New("unauthorized: cannot delete another teacher's lesson")
	}

	return s.repo.DeleteLesson(id)
}

func (s *service) UpdateHomework(teacherID, role, id string, req UpdateHomeworkRequest) (*HomeworkResponse, error) {
	existing, err := s.repo.GetHomeworkByID(id)
	if err != nil {
		return nil, err
	}
	if existing.TeacherID != teacherID && role != "admin" && role != "superadmin" {
		return nil, errors.New("unauthorized: cannot edit another teacher's homework")
	}

	h, err := s.repo.UpdateHomework(id, req)
	if err != nil {
		return nil, err
	}
	return s.mapHomeworkResponse(h), nil
}

func (s *service) DeleteHomework(teacherID, role, id string) error {
	existing, err := s.repo.GetHomeworkByID(id)
	if err != nil {
		return err
	}
	if existing.TeacherID != teacherID && role != "admin" && role != "superadmin" {
		return errors.New("unauthorized: cannot delete another teacher's homework")
	}

	return s.repo.DeleteHomework(id)
}

func (s *service) GetTeacherDiary(teacherID string, fromDate, toDate string) ([]LessonResponse, error) {
	lessons, err := s.repo.GetLessonsByTeacher(teacherID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	var res []LessonResponse
	for _, l := range lessons {
		res = append(res, *s.mapLessonResponse(&l))
	}
	return res, nil
}
