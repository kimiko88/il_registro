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
	repo      Repository
	gradeRepo grades.Repository
	classRepo classes.Repository
	userRepo  users.Repository
	attRepo   attendance.Repository
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
	if actorRole == "student" || actorRole == "parent" {
		return nil, errors.New("forbidden: studenti e genitori non possono accedere alla matrice generale di classe dello scrutinio")
	}

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

	// If caller is NOT coordinator or dirigenza (e.g. subject teacher),
	// check if scrutiny has been validated.
	if !isCoordinator && !isDirigenza {
		if actorRole == "teacher" {
			clsSubs, err := s.classRepo.GetClassSubjects(ctx, classID)
			if err != nil {
				return nil, err
			}
			isTeacherAssigned := false
			for _, cs := range clsSubs {
				if cs.TeacherID != nil && *cs.TeacherID == actorID {
					isTeacherAssigned = true
					break
				}
			}
			if !isTeacherAssigned {
				return nil, errors.New("forbidden: docente non appartenente al consiglio di classe")
			}
		}

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

	return s.buildMatrix(ctx, classID, semester, records)
}

func (s *Service) buildMatrix(ctx context.Context, classID string, semester int, records []ScrutinyRecord) (*ScrutinyMatrix, error) {
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
	allClassGrades, err := s.gradeRepo.FindByClass(ctx, classID, semester)
	if err != nil {
		logger.Log.Errorf("scrutiny FindByClass error: %v", err)
		return nil, fmt.Errorf("failed to load grades for class %s: %w", classID, err)
	}

	// Pre-indice voti: map[studentID][subjectID] → []Grade
	// Riduce la complessità da O(S×G×M) a O(G + S×M), eliminando il loop triplo.
	type gradeEntry struct {
		sum   float64
		count int
	}
	gradeIndex := make(map[string]map[string]*gradeEntry)
	for _, g := range allClassGrades {
		if !g.IsPublished || g.DeletedAt != nil {
			continue
		}
		if _, ok := gradeIndex[g.StudentID]; !ok {
			gradeIndex[g.StudentID] = make(map[string]*gradeEntry)
		}
		e := gradeIndex[g.StudentID][g.SubjectID]
		if e == nil {
			e = &gradeEntry{}
			gradeIndex[g.StudentID][g.SubjectID] = e
		}
		e.sum += g.GradeValue
		e.count++
	}

	studentIDs := make([]string, len(allStudents))
	for i, stu := range allStudents {
		studentIDs[i] = stu.ID
	}
	statsMap, _ := s.attRepo.GetStatsBatch(ctx, studentIDs)

	for _, stu := range allStudents {
		row := StudentScrutinyRow{
			StudentID:   stu.ID,
			StudentName: stu.LastName + " " + stu.FirstName,
			SubjectData: make(map[string]SubjectAverages),
		}

		for _, sub := range subjects {
			avg := 0.0
			proposed := 0.0
			count := 0
			if smap, ok := gradeIndex[stu.ID]; ok {
				if e, ok := smap[sub.SubjectID]; ok && e.count > 0 {
					count = e.count
					avg = e.sum / float64(e.count)
					proposed = math.Min(10, math.Max(1, math.Round(avg)))
				}
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

		if statsMap != nil {
			if stats, ok := statsMap[stu.ID]; ok && stats != nil {
				row.AttendanceStats = AttendanceSummary{
					Absences:   stats.TotalAbsences,
					Lates:      stats.TotalLates,
					EarlyExits: stats.TotalEarlyExits,
				}
			}
		}

		matrix.Students = append(matrix.Students, row)
	}

	return matrix, nil
}

type ClassScrutinyOverview struct {
	ClassID            string `json:"class_id"`
	ClassName          string `json:"class_name"`
	Status             string `json:"status"`
	CompletedSubjects  int    `json:"completed_subjects"`
	TotalSubjects      int    `json:"total_subjects"`
	PendingGradesCount int    `json:"pending_grades_count"`
	LastUpdated        string `json:"last_updated"`
}

type ClassScrutinyReport struct {
	ClassID   string                  `json:"class_id"`
	ClassName string                  `json:"class_name"`
	Students  []ClassReportStudentRow `json:"students"`
	Admitted  int                     `json:"admitted"`
	Rejected  int                     `json:"rejected"`
	Suspended int                     `json:"suspended"`
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
func (s *Service) GetOverview(ctx context.Context, actorID, actorRole, schoolID string, semester ...int) ([]ClassScrutinyOverview, error) {
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

	sem := 0
	if len(semester) > 0 && semester[0] > 0 {
		sem = semester[0]
	}
	if sem == 0 {
		sem = 2
		loc, err := time.LoadLocation("Europe/Rome")
		if err != nil {
			loc = time.Local
		}
		now := time.Now().In(loc)
		if now.Month() >= time.September || now.Month() <= time.January {
			sem = 1
		}
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
		allGrades, grErr := s.gradeRepo.FindByClass(ctx, c.ID, sem)
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
	cls, err := s.classRepo.Get(ctx, classID)
	if err != nil {
		return nil, err
	}
	if actorRole == "teacher" {
		isCoordinator, _, err := s.isDirigenzaOrCoordinator(ctx, actorID, actorRole, classID)
		if err != nil || !isCoordinator {
			return nil, errors.New("unauthorized: solo il coordinatore di classe o la dirigenza possono accedere al report di scrutinio")
		}
		if s.userRepo != nil {
			if user, err := s.userRepo.GetByID(ctx, actorID); err == nil && user != nil && user.SchoolID != nil && *user.SchoolID != "" && cls.SchoolID != "" && cls.SchoolID != *user.SchoolID {
				return nil, errors.New("unauthorized: classe appartenente a un'altra scuola")
			}
		}
	}
	if semester <= 0 {
		semester = 1
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

func (s *Service) SaveScrutiny(ctx context.Context, coordinatorID, actorRole, actorSchoolID string, req SaveScrutinyRequest) error {
	isCoordinator, isDirigenza, err := s.isDirigenzaOrCoordinator(ctx, coordinatorID, actorRole, req.ClassID)
	if err != nil {
		return err
	}
	if !isCoordinator && !isDirigenza {
		return ErrUnauthorizedScrutiny
	}

	// Fix cross-tenant: verifica che la classe appartenga alla stessa scuola dell'attore.
	// isDirigenzaOrCoordinator ha già caricato la classe — la recuperiamo per il check.
	if actorRole != "superadmin" && actorSchoolID != "" {
		cls, clsErr := s.classRepo.Get(ctx, req.ClassID)
		if clsErr != nil {
			return clsErr
		}
		if cls.SchoolID != actorSchoolID {
			return errors.New("forbidden: la classe non appartiene alla tua scuola")
		}
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

	if req.ConductGrade != 0 && (req.ConductGrade < 1 || req.ConductGrade > 10) {
		return fmt.Errorf("voto di condotta non valido (%d): deve essere compreso tra 1 e 10", req.ConductGrade)
	}

	canonicalDecision := req.FinalDecision
	switch req.FinalDecision {
	case "Promosso", "promosso", "ammesso", "Ammesso":
		canonicalDecision = "Ammesso"
	case "Bocciato", "bocciato", "non ammesso", "Non Ammesso", "Non ammesso":
		canonicalDecision = "Non Ammesso"
	case "Giudizio Sospeso", "giudizio sospeso", "sospeso", "Sospeso":
		canonicalDecision = "Sospeso"
	}

	rec := &ScrutinyRecord{
		StudentID:     req.StudentID,
		ClassID:       req.ClassID,
		Semester:      req.Semester,
		ConductGrade:  req.ConductGrade,
		FinalDecision: canonicalDecision,
		Notes:         req.Notes,
		CoordinatorID: coordinatorID,
		Status:        "in_progress",
	}

	var clsSubs []classes.ClassSubject
	var clsSubsLoaded bool

	for _, g := range req.Grades {
		tID := g.TeacherID
		if tID == "" {
			if !clsSubsLoaded {
				if fetched, err := s.classRepo.GetClassSubjects(ctx, req.ClassID); err == nil {
					clsSubs = fetched
				}
				clsSubsLoaded = true
			}
			for _, cs := range clsSubs {
				if cs.SubjectID == g.SubjectID && cs.TeacherID != nil && *cs.TeacherID != "" {
					tID = *cs.TeacherID
					break
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
	if actorRole == "student" {
		if actorID != studentID {
			return nil, errors.New("forbidden: lo studente può scaricare solo la propria pagella")
		}
	} else if actorRole == "parent" {
		isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, studentID)
		if err != nil || !isGuardian {
			return nil, errors.New("forbidden: genitore non autorizzato per questo studente")
		}
	} else if actorRole == "teacher" {
		isCoordinator, isDirigenza, _ := s.isDirigenzaOrCoordinator(ctx, actorID, actorRole, classID)
		if !isCoordinator && !isDirigenza {
			clsSubs, err := s.classRepo.GetClassSubjects(ctx, classID)
			if err != nil {
				return nil, err
			}
			isTeacherAssigned := false
			for _, cs := range clsSubs {
				if cs.TeacherID != nil && *cs.TeacherID == actorID {
					isTeacherAssigned = true
					break
				}
			}
			if !isTeacherAssigned {
				return nil, errors.New("forbidden: docente non appartenente al consiglio di classe")
			}
		}
	} else if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "vice_principal" && actorRole != "secretary" {
		return nil, errors.New("forbidden: ruolo non autorizzato all'esportazione pagella")
	}

	records, err := s.repo.ListRecordsByClass(ctx, classID, semester)
	if err != nil {
		return nil, err
	}

	if actorRole == "student" || actorRole == "parent" {
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

	matrix, err := s.buildMatrix(ctx, classID, semester, records)
	if err != nil {
		return nil, err
	}

	return GeneratePagellaPDF(matrix, studentID)
}

func (s *Service) ExportClassPagelleZIP(ctx context.Context, actorID, actorRole, classID string, semester int) ([]byte, error) {
	if actorRole == "teacher" {
		isCoordinator, isDirigenza, _ := s.isDirigenzaOrCoordinator(ctx, actorID, actorRole, classID)
		if !isCoordinator && !isDirigenza {
			return nil, errors.New("forbidden: solo il coordinatore di classe o la dirigenza può scaricare tutte le pagelle")
		}
	} else if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "vice_principal" && actorRole != "secretary" {
		return nil, errors.New("forbidden: ruolo non autorizzato all'esportazione pagelle di classe")
	}

	records, err := s.repo.ListRecordsByClass(ctx, classID, semester)
	if err != nil {
		return nil, err
	}

	matrix, err := s.buildMatrix(ctx, classID, semester, records)
	if err != nil {
		return nil, err
	}

	return GenerateClassPagelleZIP(matrix)
}

func (s *Service) SaveDeficiency(ctx context.Context, actorID, actorRole string, req *SaveDeficiencyRequest) error {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "teacher" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" {
		return errors.New("forbidden: non hai i permessi per inserire o modificare carenze")
	}
	if req.StudentID == "" || req.ClassID == "" || req.SubjectID == "" {
		return errors.New("student_id, class_id, and subject_id are required")
	}

	def := &StudentDeficiency{
		ID:            req.ID,
		SchoolID:      req.SchoolID,
		StudentID:     req.StudentID,
		ClassID:       req.ClassID,
		SubjectID:     req.SubjectID,
		Semester:      req.Semester,
		PeriodType:    req.PeriodType,
		Topics:        req.Topics,
		RecoveryMode:  req.RecoveryMode,
		Status:        req.Status,
		RecoveryGrade: req.RecoveryGrade,
		Notes:         req.Notes,
	}
	if req.RecoveryDate != nil && *req.RecoveryDate != "" {
		if t, err := time.Parse("2006-01-02", *req.RecoveryDate); err == nil {
			def.RecoveryDate = &t
		}
	}
	return s.repo.SaveDeficiency(ctx, def)
}

func (s *Service) GetStudentDeficiencies(ctx context.Context, actorID, actorRole, studentID string) ([]StudentDeficiency, error) {
	if studentID == "" {
		return []StudentDeficiency{}, nil
	}
	if actorRole == "student" && actorID != studentID {
		return nil, errors.New("forbidden: student can only access own deficiencies")
	}
	if actorRole == "parent" {
		isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, studentID)
		if err != nil || !isGuardian {
			return nil, errors.New("forbidden: parent is not a guardian of this student")
		}
	}
	return s.repo.GetDeficienciesByStudent(ctx, studentID)
}

func (s *Service) GetClassDeficiencies(ctx context.Context, actorID, actorRole, classID string, semester int) ([]StudentDeficiency, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "teacher" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" {
		return nil, errors.New("forbidden: non hai i permessi per accedere alle carenze della classe")
	}
	if classID == "" {
		return []StudentDeficiency{}, nil
	}
	return s.repo.GetDeficienciesByClass(ctx, classID, semester)
}

func (s *Service) SaveDeferredScrutiny(ctx context.Context, actorID, actorRole string, req *SaveDeferredScrutinyRequest) error {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "teacher" && actorRole != "principal" && actorRole != "vice_principal" {
		return errors.New("forbidden: non hai i permessi per gestire lo scrutinio differito")
	}
	if req.StudentID == "" || req.ClassID == "" {
		return errors.New("student_id and class_id are required")
	}
	return s.repo.SaveDeferredScrutiny(ctx, req)
}
