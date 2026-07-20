package scrutiny

import (
	"context"
	"errors"
	"math"
	"registro-backend/internal/attendance"
	"registro-backend/internal/classes"
	"registro-backend/internal/grades"
	"registro-backend/internal/users"
	"registro-backend/pkg/logger"
)

var (
	ErrScrutinyNotValidated = errors.New("scrutiny is not yet validated by the principal")
	ErrUnauthorizedScrutiny = errors.New("unauthorized: only class coordinator or dirigenza can perform this operation")
)

type Service struct {
	repo       Repository
	gradeRepo  grades.Repository
	classRepo  classes.Repository
	userRepo   users.Repository
	attRepo    attendance.Repository
}

func NewService(repo Repository, gr grades.Repository, cr classes.Repository, ur users.Repository, ar attendance.Repository) *Service {
	return &Service{
		repo:      repo,
		gradeRepo: gr,
		classRepo: cr,
		userRepo:  ur,
		attRepo:   ar,
	}
}

func (s *Service) isDirigenzaOrCoordinator(ctx context.Context, actorID, actorRole, classID string) (bool, bool, error) {
	isDirigenza := actorRole == "principal" || actorRole == "vice_principal" || actorRole == "admin" || actorRole == "superadmin"
	
	cls, err := s.classRepo.Get(ctx, classID)
	if err != nil {
		return false, false, err
	}

	isCoordinator := (cls.CoordinatorID == actorID) || isDirigenza
	return isCoordinator, isDirigenza, nil
}

func (s *Service) GetMatrix(ctx context.Context, actorID, actorRole, classID string, semester int) (*ScrutinyMatrix, error) {
	isCoordinator, isDirigenza, err := s.isDirigenzaOrCoordinator(ctx, actorID, actorRole, classID)
	if err != nil {
		return nil, err
	}

	// 1. Get Existing Records
	records, err := s.repo.ListRecordsByClass(ctx, classID, semester)
	if err != nil {
		return nil, err
	}

	// If caller is NOT coordinator or dirigenza (e.g. parent, student, or subject teacher),
	// check if scrutiny has been validated.
	if !isCoordinator && !isDirigenza {
		isValidated := false
		for _, r := range records {
			if r.Status == "validated" {
				isValidated = true
				break
			}
		}
		if !isValidated {
			return nil, ErrScrutinyNotValidated
		}
	}

	// 2. Get Class Subjects
	subjects, err := s.classRepo.GetClassSubjects(ctx, classID)
	if err != nil {
		return nil, err
	}

	// 3. Get Students in Class
	allStudents, err := s.userRepo.GetStudentsByClass(ctx, classID)
	if err != nil {
		return nil, err
	}
	logger.Log.Debugf("Found %d students for class %s", len(allStudents), classID)

	recordMap := make(map[string]ScrutinyRecord)
	for _, r := range records {
		recordMap[r.StudentID] = r
	}

	matrix := &ScrutinyMatrix{
		ClassID:  classID,
		Semester: semester,
	}

	for _, sub := range subjects {
		matrix.Subjects = append(matrix.Subjects, SubjectInfo{ID: sub.SubjectID, Name: sub.SubjectName})
	}

	for _, stu := range allStudents {
		row := StudentScrutinyRow{
			StudentID:   stu.ID,
			StudentName: stu.LastName + " " + stu.FirstName,
			SubjectData: make(map[string]SubjectAverages),
		}

		for _, sub := range subjects {
			gradesList, err := s.gradeRepo.FindByClassAndSubject(classID, sub.SubjectID, semester)
			if err != nil {
				continue
			}

			var sum float64
			var count int
			for _, g := range gradesList {
				if g.StudentID == stu.StudentID && g.IsPublished && g.DeletedAt == nil {
					sum += g.GradeValue
					count++
				}
			}

			avg := 0.0
			if count > 0 {
				avg = sum / float64(count)
			}

			row.SubjectData[sub.SubjectID] = SubjectAverages{
				Average:    avg,
				GradeCount: count,
				Proposed:   math.Round(avg),
			}
		}

		if _, ok := recordMap[stu.ID]; ok {
			full, _ := s.repo.GetRecord(ctx, stu.ID, classID, semester)
			row.Record = full
		}

		stats, err := s.attRepo.GetStats(stu.StudentID)
		if err == nil && stats != nil {
			row.AttendanceStats = AttendanceSummary{
				Absences:   stats.TotalAbsences,
				Lates:      stats.TotalLates,
				EarlyExits: stats.TotalEarlyExits,
			}
		}

		matrix.Students = append(matrix.Students, row)
	}

	return matrix, nil
}

func (s *Service) StartScrutiny(ctx context.Context, actorID, actorRole, classID string, semester int) error {
	isCoordinator, isDirigenza, err := s.isDirigenzaOrCoordinator(ctx, actorID, actorRole, classID)
	if err != nil {
		return err
	}
	if !isCoordinator && !isDirigenza {
		return ErrUnauthorizedScrutiny
	}

	return s.repo.UpdateClassScrutinyStatus(ctx, classID, semester, "in_progress")
}

func (s *Service) ValidateScrutiny(ctx context.Context, actorID, actorRole, classID string, semester int) error {
	_, isDirigenza, err := s.isDirigenzaOrCoordinator(ctx, actorID, actorRole, classID)
	if err != nil {
		return err
	}
	if !isDirigenza {
		return errors.New("unauthorized: only dirigenza (principal / vice_principal) can validate scrutiny")
	}

	return s.repo.ValidateClassScrutiny(ctx, classID, semester, actorID)
}

func (s *Service) CloseScrutiny(ctx context.Context, actorID, actorRole, classID string, semester int) error {
	isCoordinator, isDirigenza, err := s.isDirigenzaOrCoordinator(ctx, actorID, actorRole, classID)
	if err != nil {
		return err
	}
	if !isCoordinator && !isDirigenza {
		return ErrUnauthorizedScrutiny
	}

	return s.repo.UpdateClassScrutinyStatus(ctx, classID, semester, "closed")
}

func (s *Service) SaveScrutiny(ctx context.Context, coordinatorID, actorRole string, req SaveScrutinyRequest) error {
	isCoordinator, isDirigenza, err := s.isDirigenzaOrCoordinator(ctx, coordinatorID, actorRole, req.ClassID)
	if err != nil {
		return err
	}
	if !isCoordinator && !isDirigenza {
		return ErrUnauthorizedScrutiny
	}

	// Check if scrutiny is already validated or closed
	records, err := s.repo.ListRecordsByClass(ctx, req.ClassID, req.Semester)
	if err == nil && len(records) > 0 {
		if records[0].Status == "validated" || records[0].Status == "closed" {
			return errors.New("forbidden: cannot edit a closed or validated scrutiny")
		}
	}

	rec := &ScrutinyRecord{
		StudentID:     req.StudentID,
		ClassID:       req.ClassID,
		Semester:      req.Semester,
		ConductGrade:  req.ConductGrade,
		FinalDecision: req.FinalDecision,
		Notes:         req.Notes,
		CoordinatorID: coordinatorID,
		Status:        "in_progress",
	}

	for _, g := range req.Grades {
		rec.Grades = append(rec.Grades, ScrutinyGrade{
			SubjectID:  g.SubjectID,
			FinalGrade: g.FinalGrade,
			TeacherID:  coordinatorID,
		})
	}

	return s.repo.SaveRecord(ctx, rec)
}

func (s *Service) ExportPagellaPDF(ctx context.Context, actorID, actorRole, classID, studentID string, semester int) ([]byte, error) {
	matrix, err := s.GetMatrix(ctx, actorID, actorRole, classID, semester)
	if err != nil {
		return nil, err
	}
	return GeneratePagellaPDF(matrix, studentID)
}
