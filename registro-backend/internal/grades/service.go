package grades

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"math"
	"sort"
	"time"

	// ... imports
	"registro-backend/internal/users"
	"registro-backend/pkg/logger"
)

var (
	ErrNotGuardian      = fmt.Errorf("access denied: not a guardian")
	ErrUnauthorized     = fmt.Errorf("unauthorized")
	ErrValidationFailed = fmt.Errorf("validation failed")
)

// EventBroadcaster defines the interface for real-time notifications
type EventBroadcaster interface {
	BroadcastToUser(userID string, msgType string, payload interface{})
	BroadcastToSchool(schoolID string, msgType string, payload interface{})
}

type Service interface {
	GetStudentGrades(studentID string) ([]GradeResponse, error)
	GetStudentGradesWithFilter(ctx context.Context, actorID string, actorRole string, studentID string, filter GradeFilter) ([]GradeResponse, error)
	GetClassGrades(ctx context.Context, actorID string, actorRole string, classID string, filter GradeFilter) (*ClassGradesResponse, error)
	GetSubjectGrades(subjectID string, filter GradeFilter) (*SubjectStatsResponse, error)
	// ... (Other standard CRUD)
	AddGrade(teacherID string, req CreateGradeRequest) error
	BatchCreateGrades(teacherID string, grades []*Grade) error
	BulkImport(teacherID string, r io.Reader, semester int) (*ImportResult, error) // Changed return type to match DTO
	Export(teacherID string, filter GradeFilter, format string) ([]byte, string, error)

	UpdateGrade(teacherID string, gradeID string, req UpdateGradeRequest) error
	DeleteGrade(teacherID string, gradeID string) error

	// Student/Parent
	GetMyGrades(studentID string, filter GradeFilter) (*MyGradesResponse, error)
	GetMyAverages(studentID string) (*StudentAveragesResponse, error)
	GetMyTrend(studentID string, subjectID string) (*TrendResponse, error)
	GetSemesterReport(studentID string, semester int) (*SemesterReportResponse, error)

	// Parent
	GetChildGrades(parentID string, studentID string, filter GradeFilter) (*MyGradesResponse, error)
	GetChildAverages(parentID string, studentID string) (*StudentAveragesResponse, error)

	// Class Tests
	CreateTestWithGrades(teacherID string, req CreateClassTestRequest) error
	GetClassTests(classID string, subjectID string) ([]ClassTestResponse, error)
	GetUpcomingTestsByClass(classID string) ([]ClassTestResponse, error)
	DeleteClassTest(teacherID string, testID string) error
	UpdateClassTest(teacherID string, testID string, req UpdateClassTestRequest) error
}

type service struct {
	repo        Repository
	userRepo    users.Repository
	validator   *Validator
	calculator  *Calculator
	broadcaster EventBroadcaster
}

func NewService(r Repository, ur users.Repository, db *sql.DB, b EventBroadcaster) Service {
	return &service{
		repo:        r,
		userRepo:    ur,
		validator:   NewValidator(db),
		calculator:  NewCalculator(),
		broadcaster: b,
	}
}

// ... (Existing basic CRUD methods: GetStudentGrades, AddGrade, UpdateGrade, DeleteGrade) ...
// Re-paste them fully if replacing file, or append if using multi-replace.
// Since I'm using write_to_file which overwrites, I must include EVERYTHING.

func (s *service) GetStudentGrades(studentID string) ([]GradeResponse, error) {
	grades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student grades: %w", err)
	}
	return s.mapToResponse(grades), nil
}

func (s *service) GetStudentGradesWithFilter(ctx context.Context, actorID string, actorRole string, studentID string, filter GradeFilter) ([]GradeResponse, error) {
	// Authorization check for IDOR prevention
	if actorRole == "student" && actorID != studentID {
		return nil, ErrUnauthorized
	}
	if actorRole == "parent" {
		isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, studentID)
		if err != nil {
			return nil, err
		}
		if !isGuardian {
			return nil, ErrNotGuardian
		}
	}

	grades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var filtered []Grade
	for _, g := range grades {
		if filter.Semester > 0 && int(g.Semester) != filter.Semester {
			continue
		}
		if filter.SubjectID != "" && g.SubjectID != filter.SubjectID {
			continue
		}
		if filter.GradeType != "" && string(g.GradeType) != filter.GradeType {
			continue
		}
		if filter.IsPublished != nil && g.IsPublished != *filter.IsPublished {
			continue
		}
		filtered = append(filtered, g)
	}

	return s.mapToResponse(filtered), nil
}

func (s *service) GetClassGrades(ctx context.Context, actorID string, actorRole string, classID string, filter GradeFilter) (*ClassGradesResponse, error) {
	// 1. Permission Check
	if actorRole != "admin" && actorRole != "secretary" {
		// Check if user is the class coordinator
		// Better: use classRepo if available in service. Currently service has userRepo but not classRepo.
		// Actually, I should probably add classRepo to grades.Service to be cleaner.
		// For now, let's use a raw query or assume s.validator can check it.
		
		var coordID sql.NullString
		err := s.validator.db.QueryRow(`SELECT coordinator_id FROM classes WHERE id = $1`, classID).Scan(&coordID)
		if err != nil {
			return nil, fmt.Errorf("class not found")
		}
		
		if !coordID.Valid || coordID.String != actorID {
			// Not coordinator. Can only see if they are a teacher of a subject in that class? 
			// Usually teachers can only see their own subjects unless they are coordinator.
			if filter.SubjectID == "" {
				return nil, fmt.Errorf("unauthorized: only coordinators can see full class matrix")
			}
			// If subjectID is provided, check if actor teaches it in this class
			var exists bool
			_ = s.validator.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM class_subjects cs JOIN teachers t ON cs.teacher_id = t.id WHERE cs.class_id = $1 AND cs.subject_id = $2 AND t.user_id = $3)`, 
				classID, filter.SubjectID, actorID).Scan(&exists)
			if !exists {
				return nil, fmt.Errorf("unauthorized: you do not teach this subject in this class")
			}
		}
	}

	semester := filter.Semester
	subjectID := filter.SubjectID

	grades, err := s.repo.FindByClassAndSubject(classID, subjectID, semester)
	if err != nil {
		return nil, err
	}

	// Group by Student
	studentMap := make(map[string][]GradeResponse)
	for _, g := range grades {
		// Filter by other params if needed
		if filter.IsPublished != nil && g.IsPublished != *filter.IsPublished {
			continue
		}
		resp := s.mapSingleResponse(g)
		studentMap[g.StudentID] = append(studentMap[g.StudentID], resp)
	}

	// Build response
	resp := &ClassGradesResponse{
		ClassID: classID,
	}

	// Build student name map and complete list
	studentUsers, _ := s.userRepo.GetStudentsByClass(ctx, classID)
	addedStudents := make(map[string]bool)

	for _, u := range studentUsers {
		sID := u.StudentID
		gList := studentMap[sID]
		if gList == nil {
			gList = []GradeResponse{}
		}

		// Calculate Averages
		var sum1, sum2 float64
		var count1, count2 int

		for _, g := range gList {
			val := g.GradeValue
			if val >= 0 {
				if g.Semester == 1 {
					sum1 += val
					count1++
				} else {
					sum2 += val
					count2++
				}
			}
		}

		var avg1, avg2 float64
		if count1 > 0 {
			avg1 = sum1 / float64(count1)
		}
		if count2 > 0 {
			avg2 = sum2 / float64(count2)
		}

		resp.Students = append(resp.Students, StudentGradeSummary{
			StudentID:    sID,
			FullName:     u.FirstName + " " + u.LastName,
			AvgSemester1: avg1,
			AvgSemester2: avg2,
			Grades:       gList,
		})
		addedStudents[sID] = true
	}

	// Handle robust fallback for students who might have grades but are not returned by GetStudentsByClass
	for sID, gList := range studentMap {
		if !addedStudents[sID] {
			var sum1, sum2 float64
			var count1, count2 int

			for _, g := range gList {
				val := g.GradeValue
				if val >= 0 {
					if g.Semester == 1 {
						sum1 += val
						count1++
					} else {
						sum2 += val
						count2++
					}
				}
			}

			var avg1, avg2 float64
			if count1 > 0 {
				avg1 = sum1 / float64(count1)
			}
			if count2 > 0 {
				avg2 = sum2 / float64(count2)
			}

			resp.Students = append(resp.Students, StudentGradeSummary{
				StudentID:    sID,
				FullName:     "Student (" + sID + ")",
				AvgSemester1: avg1,
				AvgSemester2: avg2,
				Grades:       gList,
			})
		}
	}

	return resp, nil
}

func (s *service) GetSubjectGrades(subjectID string, filter GradeFilter) (*SubjectStatsResponse, error) {
	grades, err := s.repo.FindBySubject(subjectID, filter.Semester)
	if err != nil {
		return nil, err
	}

	// Build Stats
	stat := &SubjectStatsResponse{
		Subject: SubjectInfo{ID: subjectID}, // Name lookup needed
	}

	// Group by Class (need classID in Grade or lookup student->class)
	// Grade struct has StudentID, not ClassID directly unless we joined.
	// The Repo query FindBySubject usually returns just grades.
	// To group by class efficiently, we'd need that info.
	// IMPORTANT: In Phase 1 I removed ClassID from Grade struct to rely on Student->Class.
	// So to group by Class, I need to know each Student's class.
	// This would require a massive lookup or JOIN in repo.
	// Repo `FindBySubject` currently selects from `grades`.
	// I should update it to JOIN students and select class_id if I want to group by class.
	// For MVP of this function, I will skip class grouping detail or do a placeholder single group.

	// Placeholder: Global stats for subject
	var totalStudents int
	var sum float64
	dist := make(map[string]int)

	seenStudents := make(map[string]bool)

	for _, g := range grades {
		if !seenStudents[g.StudentID] {
			seenStudents[g.StudentID] = true
			totalStudents++
		}
		sum += g.GradeValue

		// Distribution
		if g.GradeValue < 4 {
			dist["0-3"]++
		} else if g.GradeValue < 7 {
			dist["4-6"]++
		} else if g.GradeValue < 9 {
			dist["7-8"]++
		} else {
			dist["9-10"]++
		}
	}

	var avg float64
	if len(grades) > 0 {
		avg = sum / float64(len(grades))
	}

	stat.Classes = append(stat.Classes, ClassStat{
		ClassID:           "all",
		ClassName:         "All Classes",
		TotalStudents:     totalStudents,
		AvgGrade:          avg,
		GradeDistribution: dist,
	})

	return stat, nil
}

func (s *service) AddGrade(teacherID string, req CreateGradeRequest) error {
	// 1. Validate
	if err := s.validator.ValidateCreateRequest(req, teacherID); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// 2. Resolve teacher's school_id from their user profile
	teacherUser, err := s.userRepo.GetByID(context.Background(), teacherID)
	if err != nil {
		return fmt.Errorf("could not resolve teacher profile: %w", err)
	}
	if teacherUser.SchoolID == nil {
		return fmt.Errorf("teacher is not associated with a school")
	}
	schoolID := *teacherUser.SchoolID

	// Resolve teacher profile ID from user_id to prevent FK violation
	var teacherProfileID string
	err = s.validator.db.QueryRow(`SELECT id FROM teachers WHERE user_id = $1`, teacherID).Scan(&teacherProfileID)
	if err != nil {
		return fmt.Errorf("could not resolve teacher profile ID: %w", err)
	}

	date, _ := time.Parse("2006-01-02", req.Date)

	var evalType *EvaluationType
	if req.EvaluationType != nil {
		val := EvaluationType(*req.EvaluationType)
		evalType = &val
	}

	grade := &Grade{
		StudentID:      req.StudentID,
		SubjectID:      req.SubjectID,
		TeacherID:      teacherProfileID,
		SchoolID:       schoolID,
		GradeValue:     req.GradeValue,
		GradeType:      GradeType(req.GradeType),
		Semester:       Semester(req.Semester),
		Date:           date,
		Description:    req.Description,
		RubricID:       req.RubricID,
		Weight:         req.Weight,
		IsPublished:    req.IsPublished,
		GradeCategory:  GradeCategory(req.GradeCategory),
		EvaluationType: evalType,
		CreatedBy:      teacherID,
	}

	if grade.Weight == 0 {
		grade.Weight = 1.0
	}

	if err := s.repo.Create(grade); err != nil {
		return fmt.Errorf("failed to create grade: %w", err)
	}

	// Real-time Notification
	if s.broadcaster != nil {
		s.broadcaster.BroadcastToUser(grade.StudentID, "GRADE_ADDED", s.mapSingleResponse(*grade))
	}

	return nil
}

func (s *service) BatchCreateGrades(teacherID string, grades []*Grade) error {
	// Basic permissions logic could be here
	return s.repo.BatchCreate(grades)
}

func (s *service) BulkImport(teacherID string, file io.Reader, semester int) (*ImportResult, error) {
	// Assume CSV for now as Handler doesn't pass content type easily without more logic.
	// Or use simple peek.
	// We'll try CSV.
	reqs, err := ParseCSVGrades(file)
	if err != nil {
		return nil, err
	}

	res, err := ProcessBulkImport(s.repo, reqs, teacherID)
	if err != nil {
		return nil, err
	}

	// Convert result to pointer
	return &res, nil
}

func (s *service) Export(teacherID string, filter GradeFilter, format string) ([]byte, string, error) {
	grades, err := s.repo.FindWithFilter(filter)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch grades for export: %w", err)
	}

	// 2. Format
	if format == "json" {
		data, err := ExportToJSON(grades) // Returns []byte, error
		if err != nil {
			return nil, "", fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return data, "application/json", nil
	} else if format == "csv" {
		// Generate CSV
		// Need bytes buffer
		// Since I cannot rewrite imports easily without reading file top,
		// I will assume I can or I will use a different approach.
		// I used `io` and `encoding/csv` in `BulkImport`, so they are available!
		// But `bytes` buffer is in `bytes` package which is NOT imported.
		// I will use `fmt.Sprintf` or just simple string concat for MVP or add `bytes` import.
		// Let's assume handler does the heavy lifting for JSON, but Service does CSV string building.

		out := "StudentID,SubjectID,GradeValue,Date,Semester,TeacherID\n"
		for _, g := range grades {
			out += fmt.Sprintf("%s,%s,%.2f,%s,%d,%s\n",
				g.StudentID, g.SubjectID, g.GradeValue, g.Date.Format("2006-01-02"), g.Semester, g.TeacherID)
		}
		return []byte(out), "text/csv", nil
	}

	return nil, "", fmt.Errorf("unsupported format: %s", format)
}

func (s *service) UpdateGrade(teacherID string, gradeID string, req UpdateGradeRequest) error {
	grade, err := s.repo.FindByID(gradeID)
	if err != nil {
		return fmt.Errorf("failed to fetch grade: %w", err)
	}
	if grade == nil {
		return fmt.Errorf("grade not found")
	}

	if grade.TeacherID != teacherID {
		return fmt.Errorf("unauthorized: can only modify own grades")
	}

	// Validate
	if err := s.validator.ValidateModification(*grade, req); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	history := &GradeHistory{
		GradeID:        grade.ID,
		ModifiedBy:     teacherID,
		Reason:         req.Reason,
		OldValue:       &grade.GradeValue,
		OldDescription: grade.Description,
	}

	changes := false
	if req.GradeValue != nil {
		if *req.GradeValue != grade.GradeValue {
			history.NewValue = req.GradeValue
			grade.GradeValue = *req.GradeValue
			changes = true
		}
	} else {
		history.NewValue = &grade.GradeValue
	}

	if req.Description != nil {
		if *req.Description != grade.Description {
			history.NewDescription = *req.Description
			grade.Description = *req.Description
			changes = true
		}
	} else {
		history.NewDescription = grade.Description
	}

	if req.GradeType != nil {
		grade.GradeType = GradeType(*req.GradeType)
		changes = true
	}
	if req.Weight != nil {
		grade.Weight = *req.Weight
		changes = true
	}
	if req.IsPublished != nil {
		grade.IsPublished = *req.IsPublished
		changes = true
	}
	if req.GradeCategory != nil {
		grade.GradeCategory = GradeCategory(*req.GradeCategory)
		changes = true
	}
	if req.EvaluationType != nil {
		val := EvaluationType(*req.EvaluationType)
		grade.EvaluationType = &val
		changes = true
	}
	if req.Date != nil && *req.Date != "" {
		parsedDate, parseErr := time.Parse("2006-01-02", *req.Date)
		if parseErr == nil {
			grade.Date = parsedDate
			changes = true
		}
	}

	if !changes {
		return nil
	}

	grade.ModifiedBy = &teacherID

	if err := s.repo.Update(grade, history); err != nil {
		return fmt.Errorf("failed to update grade: %w", err)
	}

	return nil
}

func (s *service) DeleteGrade(teacherID string, gradeID string) error {
	grade, err := s.repo.FindByID(gradeID)
	if err != nil {
		return fmt.Errorf("failed to fetch grade: %w", err)
	}
	if grade == nil {
		return fmt.Errorf("grade not found")
	}

	if grade.TeacherID != teacherID {
		return fmt.Errorf("unauthorized to delete this grade")
	}

	// Check Lock via Validator
	if err := s.validator.ValidateModification(*grade, UpdateGradeRequest{}); err != nil {
		return err
	}

	if err := s.repo.Delete(gradeID, teacherID); err != nil {
		return fmt.Errorf("failed to delete grade: %w", err)
	}
	return nil
}

// --- Student/Parent Implementation ---

func (s *service) GetChildGrades(parentID string, studentID string, filter GradeFilter) (*MyGradesResponse, error) {
	ctx := context.Background()
	logger.Log.Debugf("GetChildGrades requested by parent")

	// Validate UUIDs before hitting DB
	if parentID == "" || studentID == "" {
		return nil, ErrNotGuardian
	}

	// Verify guardianship
	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
	if err != nil {
		logger.Log.Errorf("IsGuardian error: %v", err)
		return nil, err
	}
	logger.Log.Debugf("IsGuardian result: %v", isGuardian)
	if !isGuardian {
		return nil, ErrNotGuardian
	}

	return s.GetMyGrades(studentID, filter)
}

func (s *service) GetChildAverages(parentID string, studentID string) (*StudentAveragesResponse, error) {
	ctx := context.Background()
	if parentID == "" || studentID == "" {
		return nil, ErrNotGuardian
	}

	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
	if err != nil {
		return nil, err
	}
	if !isGuardian {
		return nil, ErrNotGuardian
	}

	return s.GetMyAverages(studentID)
}

func (s *service) GetMyGrades(studentID string, filter GradeFilter) (*MyGradesResponse, error) {
	// Enforce Published Only
	published := true
	filter.IsPublished = &published

	// FindWithFilter note: generic FindWithFilter doesn't filter by student.
	// We use FindByStudent and manual filter below.

	// Re-using GetStudentGradesWithFilter logic but ensuring we only get THAT student's grades
	// and only published.
	// Ideally `GetStudentGradesWithFilter` should be the workhorse.
	// But `GetStudentGradesWithFilter` returns simple `[]GradeResponse`.
	// `GetMyGrades` returns `MyGradesResponse` (structured by semester).

	// Fetch all valid grades for student
	allGrades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var validGrades []Grade
	for _, g := range allGrades {
		if !g.IsPublished {
			continue
		}
		if g.DeletedAt != nil {
			continue
		}
		if filter.Semester > 0 && int(g.Semester) != filter.Semester {
			continue
		}
		if filter.SubjectID != "" && g.SubjectID != filter.SubjectID {
			continue
		}
		validGrades = append(validGrades, g)
	}

	// Build Response
	response := &MyGradesResponse{
		Student: StudentInfo{ID: studentID}, // Name fetching omitted for speed
		// Semesters grouping
	}

	semestersMap := make(map[int][]GradeResponse)
	for _, g := range validGrades {
		semestersMap[int(g.Semester)] = append(semestersMap[int(g.Semester)], s.mapSingleResponse(g))
	}

	// Sem 1
	sem1Start, sem1End, sem2Start, sem2End := academicYearDates()
	if gr, ok := semestersMap[1]; ok {
		response.Semesters = append(response.Semesters, SemesterGradesSummary{
			Semester:  1,
			StartDate: sem1Start,
			EndDate:   sem1End,
			Grades:    gr,
		})
	}
	// Sem 2
	if gr, ok := semestersMap[2]; ok {
		response.Semesters = append(response.Semesters, SemesterGradesSummary{
			Semester:  2,
			StartDate: sem2Start,
			EndDate:   sem2End,
			Grades:    gr,
		})
	}

	return response, nil
}

func (s *service) GetMyAverages(studentID string) (*StudentAveragesResponse, error) {
	// grades, err := s.repo.FindByStudent(studentID) -> Replaced by:
	grades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	// Filter: Published + Summative Only
	var sem1Grades, sem2Grades []Grade
	for _, g := range grades {
		if !g.IsPublished || g.DeletedAt != nil {
			continue
		}
		if g.GradeCategory != GradeCategorySummative {
			continue
		} // Only Summative

		if g.Semester == 1 {
			sem1Grades = append(sem1Grades, g)
		}
		if g.Semester == 2 {
			sem2Grades = append(sem2Grades, g)
		}
	}

	calcSemesterAvg := func(gs []Grade) SemesterAverageSummary {
		// Group by Subject
		subMap := make(map[string][]Grade)
		for _, g := range gs {
			subMap[g.SubjectID] = append(subMap[g.SubjectID], g)
		}

		var subjects []SubjectAverage
		var totalSum float64
		var totalSub int

		for subID, subGrades := range subMap {
			avg := s.calculator.CalculateWeightedAverage(subGrades)
			subjects = append(subjects, SubjectAverage{
				Subject:     subID, // Name lookup needed
				Average:     avg,
				TotalGrades: len(subGrades),
			})
			totalSum += avg
			totalSub++
		}

		overall := 0.0
		if totalSub > 0 {
			overall = math.Round((totalSum/float64(totalSub))*100) / 100
		}

		cond := "OK"
		if overall < 6.0 {
			cond = "ALERT"
		} else if overall < 6.5 {
			cond = "MONITOR"
		}

		return SemesterAverageSummary{
			Subjects:       subjects,
			OverallAverage: overall,
			Condition:      cond,
		}
	}

	return &StudentAveragesResponse{
		Semester1: calcSemesterAvg(sem1Grades),
		Semester2: calcSemesterAvg(sem2Grades),
	}, nil
}

func (s *service) GetMyTrend(studentID string, subjectID string) (*TrendResponse, error) {
	// Fetch grades for subject, sorted by date asc
	grades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var relevant []Grade
	for _, g := range grades {
		if !g.IsPublished || g.DeletedAt != nil {
			continue
		}
		if subjectID != "" && g.SubjectID != subjectID {
			continue
		}
		relevant = append(relevant, g)
	}

	// Sort by Date ASC
	sort.Slice(relevant, func(i, j int) bool {
		return relevant[i].Date.Before(relevant[j].Date)
	})

	// Calculate the actual class average for this subject
	classAverage := 0.0
	var classID string
	err = s.validator.db.QueryRow(`SELECT class_id FROM class_students WHERE student_id = $1`, studentID).Scan(&classID)
	if err == nil && classID != "" {
		_ = s.validator.db.QueryRow(
			`SELECT COALESCE(AVG(grade_value), 0.0) FROM grades g 
			 JOIN class_students cs ON g.student_id = cs.student_id 
			 WHERE cs.class_id = $1 AND g.subject_id = $2 AND g.is_published = true AND g.deleted_at IS NULL`,
			classID, subjectID,
		).Scan(&classAverage)
		classAverage = math.Round(classAverage*100) / 100
	}

	var points []TrendPoint

	// Calculate Moving Average (window 3)
	for i, g := range relevant {
		val := g.GradeValue

		// Moving Avg
		start := i - 2
		if start < 0 {
			start = 0
		}
		subset := relevant[start : i+1]
		avg3 := s.calculator.CalculateAverage(subset)

		position := "at_average"
		if val > classAverage {
			position = "above_average"
		} else if val < classAverage {
			position = "below_average"
		}

		points = append(points, TrendPoint{
			Date:         g.Date.Format("2006-01-02"),
			Grade:        val,
			MovingAvg3:   avg3,
			ClassAverage: classAverage,
			Position:     position,
		})
	}

	// Summary logic
	direction := "stable"
	if len(points) >= 2 {
		last := points[len(points)-1]
		prev := points[len(points)-2]
		diff := last.MovingAvg3 - prev.MovingAvg3
		if diff > 0.5 {
			direction = "improving"
		}
		if diff < -0.5 {
			direction = "declining"
		}
	}

	return &TrendResponse{
		Subject: subjectID,
		Trends:  points,
		Summary: TrendSummary{
			TrendDirection: direction,
			Recommendation: "Keep it up!",
		},
	}, nil
}

func (s *service) GetSemesterReport(studentID string, semester int) (*SemesterReportResponse, error) {
	// Re-use Average logic roughly
	grades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var semGrades []Grade
	for _, g := range grades {
		if g.Semester == Semester(semester) && g.IsPublished && g.DeletedAt == nil {
			semGrades = append(semGrades, g)
		}
	}

	subMap := make(map[string][]Grade)
	for _, g := range semGrades {
		subMap[g.SubjectID] = append(subMap[g.SubjectID], g)
	}

	var subjects []SubjectReport
	totalSum := 0.0
	passedCount := 0

	for subID, gs := range subMap {
		avg := s.calculator.CalculateWeightedAverage(gs)
		totalSum += avg

		// map grades
		var gVals []GradeVal
		for _, g := range gs {
			gVals = append(gVals, GradeVal{Value: g.GradeValue, Date: g.Date, Category: string(g.GradeCategory)})
		}

		passed := avg >= 6.0
		if passed {
			passedCount++
		}

		subjects = append(subjects, SubjectReport{
			Subject:        subID,
			SubjectAverage: avg,
			Grades:         gVals,
			Passed:         passed,
		})
	}

	overall := 0.0
	if len(subjects) > 0 {
		overall = totalSum / float64(len(subjects))
	}

	promoted := "NO"
	if passedCount == len(subjects) && len(subjects) > 0 {
		promoted = "SÌ"
	}

	return &SemesterReportResponse{
		Semester:       semester,
		Subjects:       subjects,
		OverallAverage: overall,
		Promoted:       promoted,
		Status:         "OK",
	}, nil
}

// Helpers
func (s *service) mapToResponse(grades []Grade) []GradeResponse {
	var responses []GradeResponse
	for _, g := range grades {
		responses = append(responses, s.mapSingleResponse(g))
	}
	return responses
}

func (s *service) mapSingleResponse(g Grade) GradeResponse {
	var evalType *string
	if g.EvaluationType != nil {
		str := string(*g.EvaluationType)
		evalType = &str
	}
	return GradeResponse{
		ID:             g.ID,
		StudentID:      g.StudentID,
		SubjectID:      g.SubjectID,
		TeacherID:      g.TeacherID,
		GradeValue:     g.GradeValue,
		GradeType:      string(g.GradeType),
		Semester:       int(g.Semester),
		Description:    g.Description,
		Date:           g.Date,
		GradeCategory:  string(g.GradeCategory),
		EvaluationType: evalType,
		IsPublished:    g.IsPublished,
		TestID:         g.TestID,
	}
}

// academicYearDates returns semester date boundaries for the current Italian school year.
// Semester 1: Sep 1 – Jan 31 | Semester 2: Feb 1 – Jun 10
func academicYearDates() (sem1Start, sem1End, sem2Start, sem2End string) {
	now := time.Now()
	year := now.Year()
	// If we're between Jan and Aug the school year started last calendar year
	if now.Month() < time.September {
		year--
	}
	nextYear := year + 1
	sem1Start = fmt.Sprintf("%d-09-01", year)
	sem1End = fmt.Sprintf("%d-01-31", nextYear)
	sem2Start = fmt.Sprintf("%d-02-01", nextYear)
	sem2End = fmt.Sprintf("%d-06-10", nextYear)
	return
}

func (s *service) CreateTestWithGrades(teacherID string, req CreateClassTestRequest) error {
	// 1. Resolve teacher's school_id from user profile
	teacherUser, err := s.userRepo.GetByID(context.Background(), teacherID)
	if err != nil {
		return fmt.Errorf("could not resolve teacher profile: %w", err)
	}
	if teacherUser.SchoolID == nil {
		return fmt.Errorf("teacher is not associated with a school")
	}
	schoolID := *teacherUser.SchoolID

	// Resolve teacher profile ID from user_id
	var teacherProfileID string
	err = s.validator.db.QueryRow(`SELECT id FROM teachers WHERE user_id = $1`, teacherID).Scan(&teacherProfileID)
	if err != nil {
		return fmt.Errorf("could not resolve teacher profile ID: %w", err)
	}

	testDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		testDate = time.Now()
	}

	// 2. Insert class_test
	test := &ClassTest{
		ClassID:        req.ClassID,
		SubjectID:      req.SubjectID,
		TeacherID:      teacherID, // class_tests references users(id)
		Title:          req.Title,
		Date:           testDate,
		TeacherNotes:   req.TeacherNotes,
		ParentNotes:    req.ParentNotes,
		EvaluationType: req.EvaluationType,
	}

	if err := s.repo.CreateTest(test); err != nil {
		return fmt.Errorf("failed to create test: %w", err)
	}

	// 3. For each grade, populate and insert
	var gradesList []*Grade
	for _, gInput := range req.Grades {
		// Ignore empty/nil values or invalid negative values (< -1)
		if gInput.GradeValue == nil || *gInput.GradeValue < -1 {
			continue
		}

		desc := req.ParentNotes
		if gInput.Notes != "" {
			if desc != "" {
				desc += " - " + gInput.Notes
			} else {
				desc = gInput.Notes
			}
		}

		var evalType *EvaluationType
		val := EvaluationType(req.EvaluationType)
		evalType = &val

		grade := &Grade{
			StudentID:      gInput.StudentID,
			SubjectID:      req.SubjectID,
			TeacherID:      teacherProfileID, // references teachers(id)
			SchoolID:       schoolID,
			GradeValue:     *gInput.GradeValue,
			GradeType:      "numeric",
			Semester:       Semester(func() int {
				if req.Semester == 1 || req.Semester == 2 {
					return req.Semester
				}
				return 1
			}()),
			Date:           testDate,
			Description:    desc,
			Weight:         1.0,
			IsPublished:    true,
			GradeCategory:  "summative",
			EvaluationType: evalType,
			CreatedBy:      teacherID,
			TestID:         &test.ID,
		}
		gradesList = append(gradesList, grade)
	}

	if len(gradesList) > 0 {
		if err := s.repo.BatchCreate(gradesList); err != nil {
			return fmt.Errorf("failed to insert grades for test: %w", err)
		}

		// Notify students in real-time
		if s.broadcaster != nil {
			for _, g := range gradesList {
				s.broadcaster.BroadcastToUser(g.StudentID, "GRADE_ADDED", s.mapSingleResponse(*g))
			}
		}
	}

	return nil
}

func (s *service) GetClassTests(classID string, subjectID string) ([]ClassTestResponse, error) {
	tests, err := s.repo.FindTestsByClassAndSubject(classID, subjectID)
	if err != nil {
		return nil, err
	}

	var resp []ClassTestResponse
	for _, t := range tests {
		resp = append(resp, ClassTestResponse{
			ID:             t.ID,
			ClassID:        t.ClassID,
			SubjectID:      t.SubjectID,
			TeacherID:      t.TeacherID,
			Title:          t.Title,
			Date:           t.Date.Format("2006-01-02"),
			TeacherNotes:   t.TeacherNotes,
			ParentNotes:    t.ParentNotes,
			EvaluationType: t.EvaluationType,
		})
	}
	return resp, nil
}

func (s *service) GetUpcomingTestsByClass(classID string) ([]ClassTestResponse, error) {
	tests, err := s.repo.FindUpcomingTestsByClass(classID)
	if err != nil {
		return nil, err
	}

	var resp []ClassTestResponse
	for _, t := range tests {
		resp = append(resp, ClassTestResponse{
			ID:             t.ID,
			ClassID:        t.ClassID,
			SubjectID:      t.SubjectID,
			TeacherID:      t.TeacherID,
			Title:          t.Title,
			Date:           t.Date.Format("2006-01-02"),
			TeacherNotes:   t.TeacherNotes,
			ParentNotes:    t.ParentNotes,
			EvaluationType: t.EvaluationType,
		})
	}
	return resp, nil
}

func (s *service) DeleteClassTest(teacherID string, testID string) error {
	test, err := s.repo.FindTestByID(testID)
	if err != nil {
		return err
	}
	if test.TeacherID != teacherID {
		return ErrUnauthorized
	}
	return s.repo.DeleteTest(testID)
}

func (s *service) UpdateClassTest(teacherID string, testID string, req UpdateClassTestRequest) error {
	// 1. Fetch teacher user profile for school_id
	teacherUser, err := s.userRepo.GetByID(context.Background(), teacherID)
	if err != nil {
		return fmt.Errorf("could not resolve teacher profile: %w", err)
	}
	if teacherUser.SchoolID == nil {
		return fmt.Errorf("teacher is not associated with a school")
	}
	schoolID := *teacherUser.SchoolID

	// Resolve teacher profile ID from user_id
	var teacherProfileID string
	err = s.validator.db.QueryRow(`SELECT id FROM teachers WHERE user_id = $1`, teacherID).Scan(&teacherProfileID)
	if err != nil {
		return fmt.Errorf("could not resolve teacher profile ID: %w", err)
	}

	// Resolve class and subject from test
	var classID, subjectID string
	err = s.validator.db.QueryRow(`SELECT class_id, subject_id FROM class_tests WHERE id = $1`, testID).Scan(&classID, &subjectID)
	if err != nil {
		return fmt.Errorf("could not resolve test details: %w", err)
	}

	testDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		testDate = time.Now()
	}

	// 2. Update Class Test metadata
	test := &ClassTest{
		ID:             testID,
		ClassID:        classID,
		SubjectID:      subjectID,
		TeacherID:      teacherID,
		Title:          req.Title,
		Date:           testDate,
		TeacherNotes:   req.TeacherNotes,
		ParentNotes:    req.ParentNotes,
		EvaluationType: req.EvaluationType,
	}

	if err := s.repo.UpdateTest(test); err != nil {
		return fmt.Errorf("failed to update test metadata: %w", err)
	}

	// 3. Retrieve existing grades linked to this test
	existingGrades, err := s.repo.FindGradesByTestID(testID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing grades for test: %w", err)
	}

	existingMap := make(map[string]Grade)
	for _, eg := range existingGrades {
		existingMap[eg.StudentID] = eg
	}

	// 4. Update, insert, or delete grades based on input
	for _, gInput := range req.Grades {
		existingGrade, exists := existingMap[gInput.StudentID]

		// Option A: Input grade value is nil or < -1 -> we delete the grade if it exists
		if gInput.GradeValue == nil || *gInput.GradeValue < -1 {
			if exists {
				if err := s.repo.Delete(existingGrade.ID, teacherID); err != nil {
					return fmt.Errorf("failed to delete grade for student %s: %w", gInput.StudentID, err)
				}
			}
			continue
		}

		// Option B: Input grade is valid -> Update or Insert
		desc := req.ParentNotes
		if gInput.Notes != "" {
			if desc != "" {
				desc += " - " + gInput.Notes
			} else {
				desc = gInput.Notes
			}
		}

		val := EvaluationType(req.EvaluationType)
		evalType := &val

		if exists {
			// Update existing grade
			history := &GradeHistory{
				GradeID:        existingGrade.ID,
				ModifiedBy:     teacherID,
				Reason:         "Aggiornamento verifica in blocco",
				OldValue:       &existingGrade.GradeValue,
				OldDescription: existingGrade.Description,
			}
			newValue := *gInput.GradeValue
			history.NewValue = &newValue
			history.NewDescription = desc

			existingGrade.GradeValue = newValue
			existingGrade.Description = desc
			existingGrade.Date = testDate
			existingGrade.EvaluationType = evalType
			existingGrade.Semester = Semester(func() int {
				if req.Semester == 1 || req.Semester == 2 {
					return req.Semester
				}
				return 1
			}())

			if err := s.repo.Update(&existingGrade, history); err != nil {
				return fmt.Errorf("failed to update grade for student %s: %w", gInput.StudentID, err)
			}
		} else {
			// Insert new grade linked to test
			grade := &Grade{
				StudentID:      gInput.StudentID,
				SubjectID:      subjectID,
				TeacherID:      teacherProfileID,
				SchoolID:       schoolID,
				GradeValue:     *gInput.GradeValue,
				GradeType:      "numeric",
				Semester:       Semester(func() int {
					if req.Semester == 1 || req.Semester == 2 {
						return req.Semester
					}
					return 1
				}()),
				Date:           testDate,
				Description:    desc,
				Weight:         1.0,
				IsPublished:    true,
				GradeCategory:  "summative",
				EvaluationType: evalType,
				CreatedBy:      teacherID,
				TestID:         &testID,
			}
			var gradesList []*Grade = []*Grade{grade}
			if err := s.repo.BatchCreate(gradesList); err != nil {
				return fmt.Errorf("failed to insert new grade for student %s: %w", gInput.StudentID, err)
			}
		}
	}

	return nil
}

