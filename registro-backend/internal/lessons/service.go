package lessons

import (
	"database/sql"
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
	GetHomeworks(classID string, fromDate ...string) ([]HomeworkResponse, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	if r == nil {
		panic("lessons.NewService: repo must not be nil")
	}
	return &service{repo: r}
}

func (s *service) CreateLesson(teacherID string, req CreateLessonRequest) (*LessonResponse, error) {
	if req.IsCoTeaching && req.IsSubstitution {
		return nil, errors.New("non è possibile contrassegnare una lezione sia come Compresenza che come Sostituzione")
	}
	if req.ClassID == "" {
		return nil, errors.New("class_id is required")
	}
	isAssigned, err := s.repo.IsTeacherAssignedToClass(teacherID, req.ClassID)
	if err != nil {
		return nil, fmt.Errorf("failed to check teacher assignment: %w", err)
	}
	if !isAssigned {
		if req.IsSubstitution {
			hasSub, subErr := s.repo.HasApprovedSubstitution(teacherID, req.ClassID, req.Date, req.Hour)
			if subErr != nil {
				return nil, fmt.Errorf("failed to check substitution: %w", subErr)
			}
			if !hasSub {
				return nil, errors.New("forbidden: docente non assegnato alla classe e nessuna sostituzione approvata trovata")
			}
		} else {
			return nil, errors.New("forbidden: docente non assegnato alla classe")
		}
	} else if req.IsSubstitution {
		hasSub, subErr := s.repo.HasApprovedSubstitution(teacherID, req.ClassID, req.Date, req.Hour)
		if subErr != nil {
			return nil, fmt.Errorf("failed to check substitution: %w", subErr)
		}
		if !hasSub {
			return nil, errors.New("forbidden: nessuna sostituzione approvata trovata per la data e ora indicate")
		}
	}

	if req.Duration <= 0 {
		return nil, errors.New("duration must be greater than 0")
	}
	if req.Hour <= 0 {
		return nil, errors.New("hour must be greater than 0")
	}

	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		loc = time.Local
	}

	date, err := time.ParseInLocation("2006-01-02", req.Date, loc)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	now := time.Now().In(loc)
	minDate := now.AddDate(-2, 0, 0)
	maxDate := now.AddDate(1, 0, 0)
	if date.Before(minDate) || date.After(maxDate) {
		return nil, fmt.Errorf("la data della lezione (%s) è fuori dall'intervallo consentito", req.Date)
	}

	// Overlap / Compresenza check: verify no conflicting lesson exists for the same class, hour, and duration on that date
	existing, err := s.repo.GetLessonsByClass(req.ClassID, req.Date)
	if err == nil {
		reqDuration := req.Duration
		if reqDuration <= 0 {
			reqDuration = 1
		}
		newStart, newEnd := req.Hour, req.Hour+reqDuration

		for _, l := range existing {
			lDur := l.Duration
			if lDur <= 0 {
				lDur = 1
			}
			exStart, exEnd := l.Hour, l.Hour+lDur

			if newStart < exEnd && exStart < newEnd {
				if !l.IsCoTeaching || !req.IsCoTeaching {
					return nil, fmt.Errorf("impossibile inserire più lezioni nella stessa ora (%dª ora) per questa classe a meno che non sia spuntata la voce 'Compresenza'", req.Hour)
				}
			}
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
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
		IsCoTeaching:         req.IsCoTeaching,
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

	res := []LessonResponse{}
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

	res := []LessonResponse{}
	for _, l := range lessons {
		res = append(res, *s.mapLessonResponse(&l))
	}
	return res, nil
}

func (s *service) CreateHomework(teacherID string, req CreateHomeworkRequest) (*HomeworkResponse, error) {
	if req.ClassID == "" {
		return nil, errors.New("class_id is required")
	}
	isAssigned, err := s.repo.IsTeacherAssignedToClass(teacherID, req.ClassID)
	if err != nil {
		return nil, fmt.Errorf("failed to check teacher assignment: %w", err)
	}
	if !isAssigned {
		return nil, errors.New("forbidden: docente non assegnato alla classe")
	}

	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	dueDate, err := time.ParseInLocation("2006-01-02", req.DueDate, loc)
	if err != nil {
		return nil, fmt.Errorf("invalid due_date format: %w", err)
	}
	if !dueDate.After(today) {
		return nil, errors.New("due_date non può essere nel passato")
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

func (s *service) GetHomeworks(classID string, fromDate ...string) ([]HomeworkResponse, error) {
	filterDate := ""
	if len(fromDate) > 0 && fromDate[0] != "" {
		filterDate = fromDate[0]
	} else {
		filterDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	homeworks, err := s.repo.GetHomeworkByClass(classID, filterDate)
	if err != nil {
		return nil, err
	}

	res := []HomeworkResponse{}
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
		IsCoTeaching:           l.IsCoTeaching,
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

func canManageLessonOrHomework(role, teacherID, existingTeacherID string) error {
	if role != "teacher" && role != "coordinator" && role != "admin" && role != "superadmin" && role != "secretary" && role != "principal" && role != "vice_principal" {
		return errors.New("unauthorized: insufficient permissions")
	}
	if existingTeacherID != "" && teacherID != existingTeacherID && role != "admin" && role != "superadmin" && role != "secretary" && role != "principal" && role != "vice_principal" {
		return errors.New("unauthorized: cannot modify another teacher's record")
	}
	return nil
}

func (s *service) UpdateLesson(teacherID, role, id string, req UpdateLessonRequest) (*LessonResponse, error) {
	existing, err := s.repo.GetLessonByID(id)
	if err != nil {
		return nil, err
	}
	if err := canManageLessonOrHomework(role, teacherID, existing.TeacherID); err != nil {
		return nil, err
	}

	targetHour := existing.Hour
	if req.Hour != nil && *req.Hour > 0 {
		targetHour = *req.Hour
	}
	targetDuration := existing.Duration
	if req.Duration != nil && *req.Duration > 0 {
		targetDuration = *req.Duration
	}
	targetCoTeaching := existing.IsCoTeaching
	if req.IsCoTeaching != nil {
		targetCoTeaching = *req.IsCoTeaching
	}
	targetSub := existing.IsSubstitution
	if req.IsSubstitution != nil {
		targetSub = *req.IsSubstitution
	}
	if targetCoTeaching && targetSub {
		return nil, errors.New("non è possibile contrassegnare una lezione sia come Compresenza che come Sostituzione")
	}

	dateStr := existing.Date.Format("2006-01-02")
	classLessons, err := s.repo.GetLessonsByClass(existing.ClassID, dateStr)
	if err == nil {
		newStart, newEnd := targetHour, targetHour+targetDuration
		for _, l := range classLessons {
			if l.ID == id {
				continue
			}
			lDur := l.Duration
			if lDur <= 0 {
				lDur = 1
			}
			exStart, exEnd := l.Hour, l.Hour+lDur
			if newStart < exEnd && exStart < newEnd {
				if !l.IsCoTeaching || !targetCoTeaching {
					return nil, fmt.Errorf("impossibile registrare più lezioni nella stessa ora (%dª ora) per questa classe senza la spunta 'Compresenza'", targetHour)
				}
			}
		}
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
	if err := canManageLessonOrHomework(role, teacherID, existing.TeacherID); err != nil {
		return err
	}

	return s.repo.DeleteLesson(id)
}

func (s *service) UpdateHomework(teacherID, role, id string, req UpdateHomeworkRequest) (*HomeworkResponse, error) {
	existing, err := s.repo.GetHomeworkByID(id)
	if err != nil {
		return nil, err
	}
	if err := canManageLessonOrHomework(role, teacherID, existing.TeacherID); err != nil {
		return nil, err
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
	if err := canManageLessonOrHomework(role, teacherID, existing.TeacherID); err != nil {
		return err
	}

	return s.repo.DeleteHomework(id)
}

func (s *service) GetTeacherDiary(teacherID string, fromDate, toDate string) ([]LessonResponse, error) {
	if fromDate == "" && toDate == "" {
		loc, err := time.LoadLocation("Europe/Rome")
		if err != nil {
			loc = time.Local
		}
		now := time.Now().In(loc)
		fromDate = now.AddDate(0, 0, -30).Format("2006-01-02")
		toDate = now.AddDate(0, 0, 30).Format("2006-01-02")
	} else if fromDate != "" && toDate != "" {
		from, err1 := time.Parse("2006-01-02", fromDate)
		to, err2 := time.Parse("2006-01-02", toDate)
		if err1 != nil || err2 != nil {
			return nil, errors.New("invalid date format: expected YYYY-MM-DD")
		}
		if from.After(to) {
			return nil, errors.New("fromDate cannot be after toDate")
		}
		if to.Sub(from) > 90*24*time.Hour {
			return nil, errors.New("date range cannot exceed 90 days")
		}
	}
	lessons, err := s.repo.GetLessonsByTeacher(teacherID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	res := []LessonResponse{}
	for _, l := range lessons {
		res = append(res, *s.mapLessonResponse(&l))
	}
	return res, nil
}
