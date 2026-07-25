package scrutiny

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
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

	// Fetch all grades for class once to avoid N+1 queries (N students * M subjects)
	allClassGrades, _ := s.gradeRepo.FindByClass(classID, semester)

	for _, stu := range allStudents {
		row := StudentScrutinyRow{
			StudentID:   stu.ID,
			StudentName: stu.LastName + " " + stu.FirstName,
			SubjectData: make(map[string]SubjectAverages),
		}

		for _, sub := range subjects {
			var sum float64
			var count int
			for _, g := range allClassGrades {
				if g.StudentID == stu.ID && g.SubjectID == sub.SubjectID && g.IsPublished && g.DeletedAt == nil {
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

		stats, err := s.attRepo.GetStats(stu.ID)
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

type ClassScrutinyOverview struct {
	ClassID             string `json:"class_id"`
	ClassName           string `json:"class_name"`
	Status              string `json:"status"`
	CompletedSubjects   int    `json:"completed_subjects"`
	TotalSubjects       int    `json:"total_subjects"`
	PendingGradesCount int    `json:"pending_grades_count"`
	LastUpdated         string `json:"last_updated"`
}

type ClassScrutinyReport struct {
	ClassID   string                 `json:"class_id"`
	ClassName string                 `json:"class_name"`
	Students  []ClassReportStudentRow `json:"students"`
	Admitted  int                    `json:"admitted"`
	Rejected  int                    `json:"rejected"`
	Suspended int                    `json:"suspended"`
}

type ClassReportStudentRow struct {
	StudentID string            `json:"student_id"`
	Name      string            `json:"name"`
	Grades    map[string]string `json:"grades"`
	Outcome   string            `json:"outcome"`
}

func (s *Service) GetOverview(ctx context.Context) ([]ClassScrutinyOverview, error) {
	classesList, err := s.classRepo.List(ctx, "", "")
	if err != nil {
		return nil, err
	}

	var res []ClassScrutinyOverview
	for _, c := range classesList {
		subjects, _ := s.classRepo.GetClassSubjects(ctx, c.ID)
		records, _ := s.repo.ListRecordsByClass(ctx, c.ID, 2)

		st := "pending"
		lastUpdated := "N/A"
		if len(records) > 0 {
			st = records[0].Status
			if st == "" {
				st = "in_progress"
			}
			lastUpdated = records[0].UpdatedAt.Format("2006-01-02T15:04:05Z")
		}

		res = append(res, ClassScrutinyOverview{
			ClassID:             c.ID,
			ClassName:           c.Name,
			Status:              st,
			CompletedSubjects:   len(subjects),
			TotalSubjects:       len(subjects),
			PendingGradesCount: 0,
			LastUpdated:         lastUpdated,
		})
	}
	return res, nil
}

func (s *Service) GetClassReport(ctx context.Context, classID string) (*ClassScrutinyReport, error) {
	cls, err := s.classRepo.Get(ctx, classID)
	if err != nil {
		return nil, err
	}

	students, err := s.userRepo.GetStudentsByClass(ctx, classID)
	if err != nil {
		return nil, err
	}

	records, _ := s.repo.ListRecordsByClass(ctx, classID, 2)
	recMap := make(map[string]ScrutinyRecord)
	for _, r := range records {
		recMap[r.StudentID] = r
	}

	report := &ClassScrutinyReport{
		ClassID:   classID,
		ClassName: cls.Name,
		Students:  make([]ClassReportStudentRow, 0, len(students)),
	}

	for _, stu := range students {
		outcome := "In corso"
		gradesMap := make(map[string]string)

		if rec, ok := recMap[stu.ID]; ok {
			outcome = rec.FinalDecision
			if outcome == "" {
				outcome = "In corso"
			}
			for _, g := range rec.Grades {
				gradesMap[g.SubjectID] = fmt.Sprintf("%.0f", g.FinalGrade)
			}
		}

		switch outcome {
		case "Ammesso", "Promosso":
			report.Admitted++
		case "Non Ammesso", "Bocciato":
			report.Rejected++
		case "Sospeso", "Giudizio Sospeso":
			report.Suspended++
		}

		report.Students = append(report.Students, ClassReportStudentRow{
			StudentID: stu.ID,
			Name:      stu.LastName + " " + stu.FirstName,
			Grades:    gradesMap,
			Outcome:   outcome,
		})
	}

	return report, nil
}

func (s *Service) FinalizeClass(ctx context.Context, classID string) error {
	return s.repo.UpdateClassScrutinyStatus(ctx, classID, 2, "closed")
}

func (s *Service) ExportAll(ctx context.Context) ([]byte, error) {
	classesList, err := s.classRepo.List(ctx, "", "")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"Classe", "Studente", "Materia", "Voto", "Esito"})

	for _, c := range classesList {
		records, err := s.repo.ListRecordsByClass(ctx, c.ID, 2)
		if err != nil {
			continue
		}
		for _, r := range records {
			studentName := r.StudentID
			if u, err := s.userRepo.GetByID(ctx, r.StudentID); err == nil && u != nil {
				studentName = u.LastName + " " + u.FirstName
			}
			for _, g := range r.Grades {
				_ = w.Write([]string{
					c.Name,
					studentName,
					g.SubjectID,
					fmt.Sprintf("%.0f", g.FinalGrade),
					r.FinalDecision,
				})
			}
		}
	}
	w.Flush()
	return buf.Bytes(), nil
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

	// Check if any record in the scrutiny for this class/semester is already validated or closed
	records, err := s.repo.ListRecordsByClass(ctx, req.ClassID, req.Semester)
	if err == nil {
		for _, r := range records {
			if r.Status == "validated" || r.Status == "closed" {
				return errors.New("forbidden: cannot edit a closed or validated scrutiny")
			}
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
		tID := g.TeacherID
		if tID == "" {
			tID = coordinatorID
		}
		rec.Grades = append(rec.Grades, ScrutinyGrade{
			SubjectID:  g.SubjectID,
			FinalGrade: g.FinalGrade,
			TeacherID:  tID,
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
