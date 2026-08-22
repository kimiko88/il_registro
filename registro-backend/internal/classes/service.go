package classes

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateClass(ctx context.Context, schoolID string, req CreateClassRequest) (*Class, error) {
	if schoolID == "" {
		schoolID = req.SchoolID
	}
	if schoolID == "" {
		return nil, fmt.Errorf("school_id is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("class name is required")
	}
	if req.AcademicYear == "" {
		return nil, fmt.Errorf("academic year is required")
	}

	c := &Class{
		SchoolID:      schoolID,
		Name:          req.Name,
		Section:       req.Section,
		Articolazione: req.Articolazione,
		AcademicYear:  req.AcademicYear,
		CoordinatorID: req.CoordinatorID,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListClasses(ctx context.Context, schoolID string, academicYear string) ([]Class, error) {
	return s.repo.List(ctx, schoolID, academicYear)
}

func (s *Service) GetTeacherClasses(ctx context.Context, teacherID string, schoolYear string) ([]Class, error) {
	return s.repo.ListByTeacher(ctx, teacherID, schoolYear)
}

func (s *Service) checkClassSchool(ctx context.Context, schoolID, classID string) error {
	if schoolID == "" {
		return nil
	}
	c, err := s.repo.Get(ctx, classID)
	if err != nil {
		return err
	}
	if c.SchoolID != schoolID {
		return fmt.Errorf("forbidden: class belongs to another school")
	}
	return nil
}

func (s *Service) GetClass(ctx context.Context, schoolID, id string) (*Class, error) {
	if err := s.checkClassSchool(ctx, schoolID, id); err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, id)
}

func (s *Service) UpdateClass(ctx context.Context, schoolID, id string, req CreateClassRequest) (*Class, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.SchoolID != schoolID {
		return nil, fmt.Errorf("forbidden: class belongs to another school")
	}

	if req.Name == "" {
		return nil, fmt.Errorf("class name is required")
	}

	c.Name = req.Name
	c.Section = req.Section
	c.Articolazione = req.Articolazione
	c.AcademicYear = req.AcademicYear
	c.CoordinatorID = req.CoordinatorID

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteClass(ctx context.Context, schoolID, id string) error {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if c.SchoolID != schoolID {
		return fmt.Errorf("forbidden: class belongs to another school")
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) AssignSubject(ctx context.Context, schoolID, classID string, req AssignSubjectRequest) error {
	c, err := s.repo.Get(ctx, classID)
	if err != nil {
		return err
	}
	if c.SchoolID != schoolID {
		return fmt.Errorf("forbidden: class belongs to another school")
	}
	return s.repo.AssignSubject(ctx, classID, req.SubjectID, req.TeacherID, req.HoursPerWeek)
}

func (s *Service) RemoveSubject(ctx context.Context, assignmentID string, schoolAndClassIDs ...string) error {
	if len(schoolAndClassIDs) >= 2 && schoolAndClassIDs[0] != "" && schoolAndClassIDs[1] != "" {
		if err := s.checkClassSchool(ctx, schoolAndClassIDs[0], schoolAndClassIDs[1]); err != nil {
			return err
		}
	}
	return s.repo.UnassignSubject(ctx, assignmentID)
}

func (s *Service) GetClassSubjects(ctx context.Context, schoolID, classID string) ([]ClassSubject, error) {
	if err := s.checkClassSchool(ctx, schoolID, classID); err != nil {
		return nil, err
	}
	return s.repo.GetClassSubjects(ctx, classID)
}

func (s *Service) GetClassGuardians(ctx context.Context, schoolID, classID string) ([]GuardianInfo, error) {
	if err := s.checkClassSchool(ctx, schoolID, classID); err != nil {
		return nil, err
	}
	return s.repo.GetClassGuardians(ctx, classID)
}

func (s *Service) GetLessonTopics(ctx context.Context, schoolID, classID string) ([]LessonTopic, error) {
	if err := s.checkClassSchool(ctx, schoolID, classID); err != nil {
		return nil, err
	}
	return s.repo.GetLessonTopics(ctx, classID)
}

func (s *Service) GetDisciplinaryNotes(ctx context.Context, schoolID, classID string) ([]DisciplinaryNoteReport, error) {
	if err := s.checkClassSchool(ctx, schoolID, classID); err != nil {
		return nil, err
	}
	return s.repo.GetDisciplinaryNotes(ctx, classID)
}

func (s *Service) BulkMigrateStudents(ctx context.Context, schoolID string, req BulkStudentMigrationRequest) error {
	if req.SourceAcademicYear == "" || req.TargetAcademicYear == "" {
		return fmt.Errorf("source and target academic years are required")
	}
	if len(req.Migrations) == 0 {
		return fmt.Errorf("no student migrations provided")
	}
	return s.repo.BulkMigrateStudents(ctx, req.Migrations)
}
