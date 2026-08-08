package scrutiny

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"time"

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
	// Bug 127: admin is NOT dirigenza for scrutiny validation purposes.
	// Only principal and vice_principal can validate (sign off) scrutiny.
	// admin/superadmin can coordinate (start, save, export) but not formally validate.
	isDirigenza := actorRole == "principal" || actorRole == "vice_principal"
	isAdmin := actorRole == "admin" || actorRole == "superadmin"

	cls, err := s.classRepo.Get(ctx, classID)
	if err != nil {
		return false, false, err
	}

	isCoordinator := (cls.CoordinatorID == actorID) || isDirigenza || isAdmin
	return isCoordinator, isDirigenza, nil
}

func (s *Service) GetMatrix(ctx context.Context, actorID, actorRole, classID string, semester int) (*ScrutinyMatrix, error) {
	isCoordinator, isDirigenza, err := s.isDirigenzaOrCoordinator(ctx, actorID, actorRole, classID)
	if err != nil {
		logger.Log.Errorf("scrutiny isDirigenzaOrCoordinator error: %v", err)
		return nil, err
	}

	// 1. Get Existing Records
	records, err := s.repo.ListRecordsByClass(ctx, classID, semester)
	if err != nil {
		logger.Log.Errorf("scrutiny ListRecordsByClass error: %v", err)
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
		logger.Log.Errorf("scrutiny GetClassSubjects error: %v", err)
		return nil, err
	}

	// 3. Get Students in Class
	allStudents, err := s.userRepo.GetStudentsByClass(ctx, classID)
	if err != nil {
		logger.Log.Errorf("scrutiny GetStudentsByClass error: %v", err)
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
	allClassGrades, err := s.gradeRepo.FindByClass(classID, semester)
	if err != nil {
		logger.Log.Errorf("scrutiny FindByClass error: %v", err)
		return nil, fmt.Errorf("failed to load grades for class %s: %w", classID, err)
	}

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
			proposed := 0.0
			if count > 0 {
				avg = sum / float64(count)
				proposed = math.Round(avg)
			}

			row.SubjectData[sub.SubjectID] = SubjectAverages{
				Average:    avg,
				GradeCount: count,
				Proposed:   proposed,
			}
		}

		if rec, ok := recordMap[stu.ID]; ok {
			fullRec := rec
			row.Record = &fullRec
		}

		stats, attErr := s.attRepo.GetStats(stu.ID)
		if attErr != nil {
			logger.Log.Errorf("scrutiny GetStats stu.ID=%s error: %v", stu.ID, attErr)
		}
		if attErr == nil && stats != nil {
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

// GetOverview returns a summary of scrutiny status for all classes in the actor's school.
// Bug 126: requires actorID, actorRole, schoolID; restricts to coordinator/dirigenza/admin;
// filters classes by schoolID to prevent cross-tenant data leaks.
func (s *Service) GetOverview(ctx context.Context, actorID, actorRole, schoolID string) ([]ClassScrutinyOverview, error) {
	if actorRole != "principal" && actorRole != "vice_principal" &&
		actorRole != "admin" && actorRole != "superadmin" && actorRole != "coordinator" && actorRole != "secretary" {
		return nil, errors.New("unauthorized: solo coordinatori, dirigenza, segreteria e admin possono vedere l'overview dello scrutinio")
	}
	if schoolID == "" && actorRole != "superadmin" {
		return nil, errors.New("unauthorized: school_id mancante per utente non superadmin")
	}
	classesList, err := s.classRepo.List(ctx, schoolID, "")
	if err != nil {
		return nil, err
	}

	sem := 2
	now := time.Now()
	if now.Month() >= time.September || now.Month() <= time.January {
		sem = 1
	}

	var res []ClassScrutinyOverview
	for _, c := range classesList {
		subjects, subErr := s.classRepo.GetClassSubjects(ctx, c.ID)
		if subErr != nil {
			logger.Log.Warnf("GetOverview: error fetching class subjects for %s: %v", c.ID, subErr)
		}
		records, recErr := s.repo.ListRecordsByClass(ctx, c.ID, sem)
		if recErr != nil {
			logger.Log.Warnf("GetOverview: error fetching scrutiny records for %s: %v", c.ID, recErr)
		}
		allGrades, grErr := s.gradeRepo.FindByClass(c.ID, sem)
		if grErr != nil {
			logger.Log.Warnf("GetOverview: error fetching grades for class %s: %v", c.ID, grErr)
		}

		gradedSubjects := make(map[string]bool)
		for _, g := range allGrades {
			if g.IsPublished && g.DeletedAt == nil {
				gradedSubjects[g.SubjectID] = true
			}
		}

		completedCount := len(gradedSubjects)
		totalSubjects := len(subjects)
		pendingCount := totalSubjects - completedCount
		if pendingCount < 0 {
			pendingCount = 0
		}

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
			ClassID:            c.ID,
			ClassName:          c.Name,
			Status:             st,
			CompletedSubjects:  completedCount,
			TotalSubjects:      totalSubjects,
			PendingGradesCount: pendingCount,
			LastUpdated:        lastUpdated,
		})
	}
	return res, nil
}

func (s *Service) GetClassReport(ctx context.Context, actorID, actorRole, classID string, semester int) (*ClassScrutinyReport, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "vice_principal" && actorRole != "secretary" && actorRole != "teacher" {
		return nil, errors.New("unauthorized: insufficient permissions to view class scrutiny report")
	}
	if actorRole == "teacher" {
		isCoordinator, _, err := s.isDirigenzaOrCoordinator(ctx, actorID, actorRole, classID)
		if err != nil || !isCoordinator {
			return nil, errors.New("unauthorized: solo il coordinatore di classe o la dirigenza possono accedere al report di scrutinio")
		}
	}
	if semester <= 0 {
		semester = 1
	}
	cls, err := s.classRepo.Get(ctx, classID)
	if err != nil {
		return nil, err
	}

	students, err := s.userRepo.GetStudentsByClass(ctx, classID)
	if err != nil {
		return nil, err
	}

	records, err := s.repo.ListRecordsByClass(ctx, classID, semester)
	if err != nil {
		records = []ScrutinyRecord{}
	}
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

func (s *Service) FinalizeClass(ctx context.Context, actorID, actorRole, classID string, semester int) error {
	isCoordinator, isDirigenza, err := s.isDirigenzaOrCoordinator(ctx, actorID, actorRole, classID)
	if err != nil {
		return err
	}
	if !isCoordinator && !isDirigenza {
		return ErrUnauthorizedScrutiny
	}
	if semester <= 0 {
		semester = 2
	}
	return s.repo.UpdateClassScrutinyStatus(ctx, classID, semester, "closed")
}

// ExportAll exports scrutiny data for all classes in the actor's school as CSV.
func (s *Service) ExportAll(ctx context.Context, actorID, actorRole, schoolID string, semester int) ([]byte, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "vice_principal" && actorRole != "secretary" {
		return nil, errors.New("unauthorized: solo dirigenza e admin possono esportare tutti gli scrutini")
	}
	if semester <= 0 {
		semester = 2
	}
	classesList, err := s.classRepo.List(ctx, schoolID, "")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"Classe", "Studente", "Materia", "Voto", "Esito"})

	for _, c := range classesList {
		subs, _ := s.classRepo.GetClassSubjects(ctx, c.ID)
		subNameMap := make(map[string]string)
		for _, sb := range subs {
			subNameMap[sb.SubjectID] = sb.SubjectName
		}
		students, _ := s.userRepo.GetStudentsByClass(ctx, c.ID)
		studentMap := make(map[string]string)
		for _, st := range students {
			studentMap[st.ID] = st.LastName + " " + st.FirstName
		}

		records, err := s.repo.ListRecordsByClass(ctx, c.ID, semester)
		if err != nil {
			continue
		}
		for _, r := range records {
			studentName := studentMap[r.StudentID]
			if studentName == "" {
				studentName = r.StudentID
			}
			for _, g := range r.Grades {
				subName := subNameMap[g.SubjectID]
				if subName == "" {
					subName = g.SubjectID
				}
				_ = w.Write([]string{
					c.Name,
					studentName,
					subName,
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
			if clsSubs, err := s.classRepo.GetClassSubjects(ctx, req.ClassID); err == nil {
				for _, cs := range clsSubs {
					if cs.SubjectID == g.SubjectID && cs.TeacherID != nil && *cs.TeacherID != "" {
						tID = *cs.TeacherID
						break
					}
				}
			}
			if tID == "" {
				tID = coordinatorID
			}
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
