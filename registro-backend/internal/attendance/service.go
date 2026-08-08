package attendance

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"registro-backend/internal/users"
)

// EventBroadcaster defines the interface for real-time notifications
type EventBroadcaster interface {
	BroadcastToUser(userID string, msgType string, payload interface{})
	BroadcastToSchool(schoolID string, msgType string, payload interface{})
}

// CalendarService is used to calculate real school days for absence rate.
type CalendarService interface {
	CountTeachingDays(ctx context.Context, schoolID string, from, to time.Time) (int, error)
	GetSchoolYearDates(ctx context.Context, schoolID string) (start, end time.Time, err error)
}

type Service interface {
	MarkAttendance(ctx context.Context, teacherID, schoolID string, req CreateAttendanceRequest) error
	MarkBulk(ctx context.Context, teacherID, schoolID string, req BulkAttendanceRequest) error
	UpdateAttendance(ctx context.Context, teacherID, schoolID, id string, req UpdateAttendanceRequest) error
	DeleteClassAttendanceHour(ctx context.Context, classID, dateStr string, hour int) error

	GetClassAttendance(ctx context.Context, actorID, actorRole, schoolID, classID string, date string) (*ClassDailyAttendance, error)
	GetStudentAttendance(ctx context.Context, actorID, actorRole, schoolID, studentID string, from, to time.Time) ([]AttendanceResponse, error)

	// Justifications
	RequestJustification(ctx context.Context, parentID string, req JustificationRequest) error
	ProcessJustification(ctx context.Context, teacherID, justificationID string, approve bool) error
	GetPendingJustifications(ctx context.Context, classID string) ([]JustificationResponse, error)
	DeleteJustification(ctx context.Context, actorID string, justificationID string) error

	// Analytics & Summaries
	GetStudentSummary(ctx context.Context, studentID, schoolID string) (*SummaryResponse, error)
	GetSchoolAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error)

	// Parent Access (with guardianship checks)
	GetChildAttendance(ctx context.Context, parentID, studentID string, from, to time.Time) ([]AttendanceResponse, error)
	GetChildSummary(ctx context.Context, parentID, studentID, schoolID string) (*SummaryResponse, error)
	GetChildAttendanceTrends(ctx context.Context, parentID, studentID string) (*TrendsResponse, error)

	// Monthly Breakdown
	GetMonthlyBreakdown(ctx context.Context, studentID, schoolYear string) (*MonthlyBreakdownResponse, error)
	GetChildMonthlyBreakdown(ctx context.Context, parentID, studentID, schoolYear string) (*MonthlyBreakdownResponse, error)

	GetChildUnjustified(ctx context.Context, parentID, studentID string) ([]Attendance, error)
	JustifyChildAbsence(ctx context.Context, parentID, studentID, attendanceID string, req JustifyAbsenceRequest) error
	GetChildAttendanceStats(ctx context.Context, parentID, studentID string) (*AttendanceStats, error)
}

type service struct {
	repo        Repository
	userRepo    users.Repository
	validator   *Validator
	broadcaster EventBroadcaster
	calendar    CalendarService
}

func NewService(repo Repository, uRepo users.Repository, b EventBroadcaster, cal CalendarService) Service {
	if repo == nil {
		panic("attendance.NewService: repo must not be nil")
	}
	if uRepo == nil {
		panic("attendance.NewService: userRepo must not be nil — required for guardianship and access checks")
	}
	return &service{
		repo:        repo,
		userRepo:    uRepo,
		validator:   NewValidator(),
		broadcaster: b,
		calendar:    cal,
	}
}

func (s *service) MarkAttendance(ctx context.Context, teacherID, schoolID string, req CreateAttendanceRequest) error {
	if schoolID == "" {
		return fmt.Errorf("school_id mancante nel token")
	}
	if teacherID == "" {
		return fmt.Errorf("forbidden: teacherID mancante")
	}
	isAssigned, err := s.repo.IsTeacherAssignedToClass(ctx, teacherID, req.ClassID)
	if err != nil {
		return fmt.Errorf("errore verifica docente per classe: %w", err)
	}
	if !isAssigned {
		return fmt.Errorf("forbidden: docente non assegnato alla classe")
	}

	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	date, err := time.ParseInLocation("2006-01-02", req.Date, loc)
	if err != nil {
		return fmt.Errorf("data non valida '%s': usa il formato YYYY-MM-DD", req.Date)
	}
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
	if date.After(todayEnd) {
		return fmt.Errorf("impossibile registrare presenze per date future (%s)", req.Date)
	}

	att := &Attendance{
		SchoolID:  schoolID,
		StudentID: req.StudentID,
		ClassID:   req.ClassID,
		Date:      date,
		Hour:      &req.Hour,
		Status:    req.Status,
		Notes:     req.Notes,
	}
	if req.SubjectID != "" {
		subjectIDCopy := req.SubjectID
		att.SubjectID = &subjectIDCopy
	}

	if req.EntryTime != "" {
		att.EntryTime = &req.EntryTime
	}
	if req.ExitTime != "" {
		att.ExitTime = &req.ExitTime
	}

	if err := s.validator.ValidateEntry(att); err != nil {
		return err
	}

	if err := s.repo.Create(att); err != nil {
		return err
	}

	if s.broadcaster != nil {
		s.broadcaster.BroadcastToUser(att.StudentID, "ATTENDANCE_"+string(att.Status), att)
	}

	return nil
}

func (s *service) MarkBulk(ctx context.Context, teacherID, schoolID string, req BulkAttendanceRequest) error {
	if schoolID == "" {
		return fmt.Errorf("school_id mancante nel token")
	}
	if teacherID == "" {
		return fmt.Errorf("forbidden: teacherID mancante")
	}
	if !req.IsSubstitution {
		isAssigned, err := s.repo.IsTeacherAssignedToClass(ctx, teacherID, req.ClassID)
		if err != nil {
			return fmt.Errorf("errore verifica docente per classe: %w", err)
		}
		if !isAssigned {
			return fmt.Errorf("forbidden: docente non autorizzato per la classe")
		}
	} else {
		loc, _ := time.LoadLocation("Europe/Rome")
		dateParsed, _ := time.ParseInLocation("2006-01-02", req.Date, loc)
		isSub, err := s.repo.IsTeacherSubstitute(ctx, teacherID, req.ClassID, dateParsed, req.Hour)
		if err != nil || !isSub {
			isAssigned, _ := s.repo.IsTeacherAssignedToClass(ctx, teacherID, req.ClassID)
			if !isAssigned {
				return fmt.Errorf("forbidden: docente non assegnato e non registrato come supplente per questa classe/ora")
			}
		}
	}

	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	date, err := time.ParseInLocation("2006-01-02", req.Date, loc)
	if err != nil {
		return fmt.Errorf("data non valida '%s': usa il formato YYYY-MM-DD", req.Date)
	}
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
	if date.After(todayEnd) {
		return fmt.Errorf("impossibile registrare presenze per date future (%s)", req.Date)
	}

	if len(req.Statuses) == 0 {
		return fmt.Errorf("nessun record di presenza fornito")
	}
	if len(req.Statuses) > 500 {
		return fmt.Errorf("numero massimo di presenze registrabili in blocco superato (max 500)")
	}

	seen := make(map[string]bool)
	var atts []*Attendance
	for _, r := range req.Statuses {
		key := fmt.Sprintf("%s_%s_%d", r.StudentID, date.Format("2006-01-02"), req.Hour)
		if seen[key] {
			continue
		}
		seen[key] = true

		att := &Attendance{
			SchoolID:  schoolID,
			StudentID: r.StudentID,
			ClassID:   req.ClassID,
			Date:      date,
			Hour:      &req.Hour,
			Status:    r.Status,
			Notes:     r.Notes,
		}
		// Only set SubjectID if a valid non-empty UUID string is provided.
		if req.SubjectID != "" {
			subjectIDCopy := req.SubjectID
			att.SubjectID = &subjectIDCopy
		}
		if r.EntryTime != "" {
			att.EntryTime = &r.EntryTime
		}
		if r.ExitTime != "" {
			att.ExitTime = &r.ExitTime
		}
		if err := s.validator.ValidateEntry(att); err != nil {
			return err
		}
		atts = append(atts, att)
	}

	if err := s.repo.BatchCreate(atts); err != nil {
		return err
	}

	if s.broadcaster != nil {
		for _, att := range atts {
			s.broadcaster.BroadcastToUser(att.StudentID, "ATTENDANCE_"+string(att.Status), att)
		}
	}

	return nil
}

func (s *service) UpdateAttendance(ctx context.Context, teacherID, schoolID, id string, req UpdateAttendanceRequest) error {
	att, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if att.Justified {
		return fmt.Errorf("impossibile modificare la presenza: il record è già stato giustificato")
	}

	// Ownership check: il docente deve appartenere alla stessa scuola ed essere assegnato alla classe del record.
	if att.SchoolID != schoolID {
		return fmt.Errorf("forbidden: impossibile modificare presenze di un'altra scuola")
	}
	if teacherID == "" {
		return fmt.Errorf("forbidden: teacherID mancante")
	}
	isAssigned, err := s.repo.IsTeacherAssignedToClass(ctx, teacherID, att.ClassID)
	if err != nil || !isAssigned {
		return fmt.Errorf("forbidden: docente non assegnato alla classe del record di presenza")
	}

	if req.Status != nil {
		att.Status = *req.Status
	}
	if req.Notes != nil {
		att.Notes = *req.Notes
	}
	if req.EntryTime != nil {
		att.EntryTime = req.EntryTime
	}
	if req.ExitTime != nil {
		att.ExitTime = req.ExitTime
	}

	return s.repo.Update(att)
}

func (s *service) DeleteClassAttendanceHour(ctx context.Context, classID, dateStr string, hour int) error {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("data non valida '%s': usa il formato YYYY-MM-DD", dateStr)
	}
	return s.repo.DeleteByClassDateHour(classID, date, hour)
}

func (s *service) GetClassAttendance(ctx context.Context, actorID, actorRole, schoolID, classID string, dateStr string) (*ClassDailyAttendance, error) {
	if actorRole == "" {
		return nil, fmt.Errorf("unauthorized: missing actorRole")
	}
	if actorRole == "teacher" {
		isAssigned, err := s.repo.IsTeacherAssignedToClass(ctx, actorID, classID)
		if err != nil || !isAssigned {
			return nil, fmt.Errorf("forbidden: docente non assegnato alla classe")
		}
	} else if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" {
		return nil, fmt.Errorf("forbidden: ruolo non autorizzato alla lettura delle presenze di classe")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("data non valida '%s': usa il formato YYYY-MM-DD", dateStr)
	}
	atts, err := s.repo.FindByClassAndDate(classID, date)
	if err != nil {
		return nil, err
	}

	resp := &ClassDailyAttendance{
		ClassID: classID,
		Date:    dateStr,
		Records: []AttendanceResponse{},
	}

	p, a, l := 0, 0, 0
	for _, att := range atts {
		r := AttendanceResponse{
			ID:          att.ID,
			StudentID:   att.StudentID,
			Date:        att.Date.Format("2006-01-02"),
			Status:      att.Status,
			IsJustified: att.Justified,
			Notes:       att.Notes,
		}
		if att.Hour != nil {
			r.Hour = *att.Hour
		}
		if att.EntryTime != nil && *att.EntryTime != "" {
			r.EntryTime = *att.EntryTime
		}
		if att.ExitTime != nil {
			r.ExitTime = *att.ExitTime
		}

		resp.Records = append(resp.Records, r)

		switch att.Status {
		case StatusPresent:
			p++
		case StatusAbsent:
			a++
		case StatusLate:
			l++
		}
	}
	resp.Summary.Present = p
	resp.Summary.Absent = a
	resp.Summary.Late = l

	return resp, nil
}

func (s *service) GetStudentAttendance(ctx context.Context, actorID, actorRole, schoolID, studentID string, from, to time.Time) ([]AttendanceResponse, error) {
	if actorRole == "student" && actorID != studentID {
		return nil, fmt.Errorf("unauthorized: uno studente può leggere solo le proprie presenze")
	}
	if actorRole == "parent" {
		if s.userRepo == nil {
			return nil, fmt.Errorf("unauthorized: userRepo is missing")
		}
		isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, studentID)
		if err != nil || !isGuardian {
			return nil, fmt.Errorf("unauthorized: not a guardian of this student")
		}
	}

	atts, err := s.repo.FindByStudent(studentID, from, to)
	if err != nil {
		return nil, err
	}

	var resp []AttendanceResponse
	for _, att := range atts {
		r := AttendanceResponse{
			ID:          att.ID,
			StudentID:   att.StudentID,
			Date:        att.Date.Format("2006-01-02"),
			Status:      att.Status,
			IsJustified: att.Justified,
			Notes:       att.Notes,
		}
		if att.EntryTime != nil {
			r.EntryTime = *att.EntryTime
		}
		if att.ExitTime != nil {
			r.ExitTime = *att.ExitTime
		}
		resp = append(resp, r)
	}
	return resp, nil
}

// Justifications

func (s *service) RequestJustification(ctx context.Context, parentID string, req JustificationRequest) error {
	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, req.StudentID)
	if err != nil || !isGuardian {
		return fmt.Errorf("unauthorized: parent is not a guardian of this student")
	}

	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return fmt.Errorf("start_date non valida '%s': usa il formato YYYY-MM-DD", req.StartDate)
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return fmt.Errorf("end_date non valida '%s': usa il formato YYYY-MM-DD", req.EndDate)
	}

	if end.Before(start) {
		return fmt.Errorf("end_date (%s) non può essere precedente a start_date (%s)", req.EndDate, req.StartDate)
	}

	hasOverlap, err := s.repo.HasOverlappingJustification(ctx, req.StudentID, start, end)
	if err != nil {
		return fmt.Errorf("failed to check existing justifications: %w", err)
	}
	if hasOverlap {
		return fmt.Errorf("esiste già una richiesta di giustificazione per questo periodo")
	}

	j := &Justification{
		StudentID: req.StudentID,
		ParentID:  parentID,
		StartDate: start,
		EndDate:   end,
		Reason:    req.Reason,
		Status:    JustificationPending,
	}

	if err := s.validator.ValidateJustification(j); err != nil {
		return err
	}

	return s.repo.CreateJustification(j)
}

func (s *service) ProcessJustification(ctx context.Context, teacherID, justificationID string, approve bool) error {
	j, err := s.repo.FindJustificationByID(justificationID)
	if err != nil {
		return err
	}

	if j.Status != JustificationPending {
		return fmt.Errorf("la giustifica %s è già stata elaborata (stato attuale: %s)", justificationID, j.Status)
	}

	if s.userRepo != nil {
		studentUser, err := s.userRepo.GetByID(ctx, j.StudentID)
		if err != nil || studentUser == nil {
			return fmt.Errorf("student not found or error fetching student: %w", err)
		}
		if studentUser.ClassID == nil || *studentUser.ClassID == "" {
			return fmt.Errorf("impossibile processare giustifica: studente non assegnato a nessuna classe")
		}
		isAssigned, err := s.repo.IsTeacherAssignedToClass(ctx, teacherID, *studentUser.ClassID)
		if err != nil || !isAssigned {
			return fmt.Errorf("forbidden: docente non assegnato alla classe dello studente")
		}
	}

	if err := s.repo.ProcessJustificationTx(ctx, j, teacherID, approve); err != nil {
		return err
	}

	if approve {
		j.Status = JustificationApproved
	} else {
		j.Status = JustificationRejected
	}

	if s.broadcaster != nil {
		payload := map[string]interface{}{
			"id":     j.ID,
			"status": string(j.Status),
			"reason": j.Reason,
		}
		if j.ParentID != "" {
			s.broadcaster.BroadcastToUser(j.ParentID, "justification_processed", payload)
		}
		if j.StudentID != "" {
			s.broadcaster.BroadcastToUser(j.StudentID, "justification_processed", payload)
		}
	}

	return nil
}

func (s *service) GetPendingJustifications(ctx context.Context, classID string) ([]JustificationResponse, error) {
	js, err := s.repo.FindPendingJustifications(classID)
	if err != nil {
		return nil, err
	}

	var resp []JustificationResponse
	for _, j := range js {
		resp = append(resp, JustificationResponse{
			ID:          j.ID,
			StudentID:   j.StudentID,
			StudentName: j.StudentName,
			Date:        j.StartDate.Format("2006-01-02"),
			Status:      string(j.Status),
			Reason:      j.Reason,
			DateRange:   fmt.Sprintf("%s - %s", j.StartDate.Format("2006-01-02"), j.EndDate.Format("2006-01-02")),
		})
	}
	return resp, nil
}

// DeleteJustification elimina fisicamente la giustifica previa verifica dell'autorizzazione dell'actorID.
func (s *service) DeleteJustification(ctx context.Context, actorID string, justificationID string) error {
	if actorID == "" {
		return errors.New("unauthorized: actorID non fornito")
	}

	j, err := s.repo.FindJustificationByID(justificationID)
	if err != nil {
		return fmt.Errorf("giustifica non trovata: %w", err)
	}

	if j.Status != JustificationPending {
		return fmt.Errorf("impossibile eliminare una giustifica già elaborata (stato attuale: %s)", j.Status)
	}

	isGuardian := false
	if s.userRepo != nil && j.StudentID != "" {
		g, _ := s.userRepo.IsGuardian(ctx, actorID, j.StudentID)
		isGuardian = g
	}

	if actorID != j.ParentID && actorID != j.StudentID && !isGuardian {
		if s.userRepo == nil {
			return errors.New("unauthorized: userRepo non configurato per il controllo permessi")
		}
		actorUser, err := s.userRepo.GetByID(ctx, actorID)
		if err != nil {
			return fmt.Errorf("unauthorized: impossibile verificare permessi dell'utente: %w", err)
		}
		if actorUser.Role == "teacher" {
			return fmt.Errorf("forbidden: i docenti non possono eliminare le giustifiche create dai genitori")
		}
		studentUser, err := s.userRepo.GetByID(ctx, j.StudentID)
		if err != nil {
			return fmt.Errorf("unauthorized: impossibile verificare lo studente della giustifica: %w", err)
		}
		if actorUser.Role != "superadmin" {
			// Bug 104: fix nil-SchoolID bypass — verifica esplicita dei puntatori prima di dereferenziare
			if actorUser.SchoolID == nil || studentUser.SchoolID == nil {
				// Se uno dei due SchoolID mancano non possiamo garantire la co-appartenenza
				return fmt.Errorf("forbidden: impossibile verificare appartenenza scolastica (schoolID mancante)")
			}
			if *actorUser.SchoolID != *studentUser.SchoolID {
				return fmt.Errorf("forbidden: impossibile eliminare giustifiche di un'altra scuola")
			}
		}
	}

	return s.repo.DeleteJustification(justificationID)
}

func (s *service) GetStudentSummary(ctx context.Context, studentID, schoolID string) (*SummaryResponse, error) {
	stats, err := s.repo.GetStats(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	now := time.Now()
	startYear := now.Year()
	if now.Month() < time.September {
		startYear--
	}
	startOfSchoolYear := time.Date(startYear, time.September, 1, 0, 0, 0, 0, now.Location())
	elapsedDays := int(now.Sub(startOfSchoolYear).Hours() / 24)
	if elapsedDays <= 0 {
		elapsedDays = 1
	}

	totalDays := elapsedDays
	if s.calendar != nil && schoolID != "" {
		count, err := s.calendar.CountTeachingDays(ctx, schoolID, startOfSchoolYear, now)
		if err == nil && count > 0 {
			totalDays = count
		}
	}

	if totalDays > 0 {
		stats.AbsenceRate = (float64(stats.TotalAbsences) / float64(totalDays)) * 100
	} else {
		stats.AbsenceRate = 0
	}

	if stats.AbsenceRate > 25.0 {
		stats.RiskLevel = "Critical"
	} else if stats.AbsenceRate > 15.0 {
		stats.RiskLevel = "Warning"
	} else {
		stats.RiskLevel = "Normal"
	}

	return stats, nil
}

func (s *service) GetSchoolAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error) {
	return s.repo.GetAnalytics(ctx, schoolID)
}

func (s *service) GetChildAttendance(ctx context.Context, parentID, studentID string, from, to time.Time) ([]AttendanceResponse, error) {
	if s.userRepo == nil {
		return nil, fmt.Errorf("unauthorized: userRepo is nil")
	}
	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
	if err != nil {
		return nil, err
	}
	if !isGuardian {
		return nil, fmt.Errorf("unauthorized: not a guardian of this student")
	}
	studentUser, err := s.userRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("impossibile recuperare lo studente: %w", err)
	}
	if studentUser == nil || studentUser.SchoolID == nil {
		return nil, fmt.Errorf("lo studente %s non ha una scuola associata", studentID)
	}
	schoolID := *studentUser.SchoolID
	return s.GetStudentAttendance(ctx, parentID, "parent", schoolID, studentID, from, to)
}

func (s *service) GetChildSummary(ctx context.Context, parentID, studentID, schoolID string) (*SummaryResponse, error) {
	if s.userRepo == nil {
		return nil, fmt.Errorf("unauthorized: userRepo is nil")
	}
	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
	if err != nil {
		return nil, err
	}
	if !isGuardian {
		return nil, fmt.Errorf("unauthorized: not a guardian of this student")
	}
	return s.GetStudentSummary(ctx, studentID, schoolID)
}

func (s *service) GetChildAttendanceTrends(ctx context.Context, parentID, studentID string) (*TrendsResponse, error) {
	if s.userRepo == nil {
		return nil, fmt.Errorf("unauthorized: userRepo is nil")
	}
	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
	if err != nil {
		return nil, err
	}
	if !isGuardian {
		return nil, fmt.Errorf("unauthorized: not a guardian of this student")
	}

	now := time.Now()
	startYear := now.Year()
	if now.Month() < time.September {
		startYear--
	}
	from := time.Date(startYear, time.September, 1, 0, 0, 0, 0, now.Location())
	to := now
	atts, err := s.repo.FindByStudent(studentID, from, to)
	if err != nil {
		return nil, err
	}

	monthlyMap := make(map[string]*MonthlyTrend)
	monthlyTotalMap := make(map[string]int)
	monthlyPresentMap := make(map[string]int)
	var monthKeys []string

	for _, a := range atts {
		mKey := a.Date.Format("2006-01")
		tr, exists := monthlyMap[mKey]
		if !exists {
			tr = &MonthlyTrend{Month: mKey}
			monthlyMap[mKey] = tr
			monthKeys = append(monthKeys, mKey)
		}
		monthlyTotalMap[mKey]++
		if a.Status == StatusPresent || a.Status == StatusLate || a.Status == StatusEarlyExit {
			monthlyPresentMap[mKey]++
		}
		switch a.Status {
		case StatusAbsent:
			tr.Absences++
		case StatusLate:
			tr.Lates++
		case StatusEarlyExit:
			tr.EarlyExits++
		}
	}

	sort.Strings(monthKeys)

	resp := &TrendsResponse{
		StudentID: studentID,
		Trends:    []MonthlyTrend{},
	}
	for _, k := range monthKeys {
		tr := monthlyMap[k]
		totalEntries := monthlyTotalMap[k]
		presentCount := monthlyPresentMap[k]
		rate := 100.0
		if totalEntries > 0 {
			rate = (float64(presentCount) / float64(totalEntries)) * 100.0
		}
		tr.PresenceRate = rate
		resp.Trends = append(resp.Trends, *tr)
	}

	return resp, nil
}

// GetMonthlyBreakdown returns a detailed per-month attendance breakdown for a student.
func (s *service) GetMonthlyBreakdown(ctx context.Context, studentID, schoolYear string) (*MonthlyBreakdownResponse, error) {
	if schoolYear == "" {
		now := time.Now()
		if now.Month() >= 9 {
			schoolYear = fmt.Sprintf("%d-%d", now.Year(), now.Year()+1)
		} else {
			schoolYear = fmt.Sprintf("%d-%d", now.Year()-1, now.Year())
		}
	}
	months, err := s.repo.GetMonthlyBreakdown(ctx, studentID, schoolYear)
	if err != nil {
		return nil, fmt.Errorf("GetMonthlyBreakdown: %w", err)
	}
	if months == nil {
		months = []MonthlyBreakdownRow{}
	}
	return &MonthlyBreakdownResponse{
		StudentID:  studentID,
		SchoolYear: schoolYear,
		Months:     months,
	}, nil
}

func (s *service) GetChildUnjustified(ctx context.Context, parentID, studentID string) ([]Attendance, error) {
	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
	if err != nil {
		return nil, err
	}
	if !isGuardian {
		return nil, fmt.Errorf("parent is not a guardian of student")
	}
	return s.repo.FindUnjustifiedByStudent(studentID)
}

func (s *service) JustifyChildAbsence(ctx context.Context, parentID, studentID, attendanceID string, req JustifyAbsenceRequest) error {
	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
	if err != nil {
		return err
	}
	if !isGuardian {
		return fmt.Errorf("parent is not a guardian of student")
	}
	return s.repo.JustifyAbsenceByParent(attendanceID, req.Reason, req.Notes)
}

func (s *service) GetChildAttendanceStats(ctx context.Context, parentID, studentID string) (*AttendanceStats, error) {
	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
	if err != nil {
		return nil, err
	}
	if !isGuardian {
		return nil, fmt.Errorf("parent is not a guardian of student")
	}
	return s.repo.GetStudentAttendanceStats(studentID)
}

// GetChildMonthlyBreakdown is the parent-facing version with guardianship check.
func (s *service) GetChildMonthlyBreakdown(ctx context.Context, parentID, studentID, schoolYear string) (*MonthlyBreakdownResponse, error) {
	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
	if err != nil || !isGuardian {
		return nil, fmt.Errorf("access denied: not a guardian of this student")
	}
	return s.GetMonthlyBreakdown(ctx, studentID, schoolYear)
}
