package grades

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"math"
	"sort"
	"time"

	"github.com/go-pdf/fpdf"

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
	// GetStudentGradesPaged returns a paginated response for list-grades endpoints.
	GetStudentGradesPaged(ctx context.Context, actorID string, actorRole string, studentID string, filter GradeFilter) (*PaginatedGradesResponse, error)
	GetClassGrades(ctx context.Context, actorID string, actorRole string, classID string, filter GradeFilter) (*ClassGradesResponse, error)
	GetSubjectGrades(ctx context.Context, actorID string, actorRole string, subjectID string, filter GradeFilter) (*SubjectStatsResponse, error)
	AddGrade(teacherID string, req CreateGradeRequest) (*GradeResponse, error)
	BatchCreateGrades(teacherID string, grades []*Grade) error
	BulkImport(teacherID string, r io.Reader, semester int) (*ImportResult, error)
	Export(teacherID string, filter GradeFilter, format string) ([]byte, string, error)
	UpdateGrade(teacherID string, gradeID string, req UpdateGradeRequest) (*GradeResponse, error)
	DeleteGrade(teacherID string, gradeID string) error

	// Student/Parent
	GetMyGrades(studentID string, filter GradeFilter) (*MyGradesResponse, error)
	GetMyAverages(studentID string) (*StudentAveragesResponse, error)
	GetMyTrend(studentID string, subjectID string) (*TrendResponse, error)
	GetSemesterReport(studentID string, semester int) (*SemesterReportResponse, error)
	GenerateSemesterReportPDF(studentID string, semester int) ([]byte, error)

	// Parent
	GetChildGrades(parentID string, studentID string, filter GradeFilter) (*MyGradesResponse, error)
	GetChildAverages(parentID string, studentID string) (*StudentAveragesResponse, error)
	GetChildSemesterReport(ctx context.Context, parentID, studentID string, semester int) (*SemesterReportResponse, error)

	// Class Tests
	CreateTestWithGrades(teacherID string, req CreateClassTestRequest) error
	GetClassTests(classID string, subjectID string) ([]ClassTestResponse, error)
	GetUpcomingTestsByClass(classID string) ([]ClassTestResponse, error)
	DeleteClassTest(teacherID string, testID string) error
	UpdateClassTest(teacherID string, testID string, req UpdateClassTestRequest) error

	// Weight Config
	GetWeightConfigs(schoolID, subjectID, classID string) ([]GradeWeightConfig, error)
	UpsertWeightConfig(actorID, schoolID string, req UpsertWeightConfigRequest) (*GradeWeightConfig, error)
	DeleteWeightConfig(actorID, configID string) error
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

func (s *service) GetStudentGrades(studentID string) ([]GradeResponse, error) {
	grades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student grades: %w", err)
	}
	return s.mapToResponse(grades), nil
}

func (s *service) GetStudentGradesWithFilter(ctx context.Context, actorID string, actorRole string, studentID string, filter GradeFilter) ([]GradeResponse, error) {
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

// GetStudentGradesPaged applies the same ownership checks as GetStudentGradesWithFilter
// then delegates to the paginated repository method.
func (s *service) GetStudentGradesPaged(ctx context.Context, actorID string, actorRole string, studentID string, filter GradeFilter) (*PaginatedGradesResponse, error) {
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

	// Inject studentID as a filter — FindWithFilterPaginated honours the
	// existing conditions; we add it via a separate lookup to keep the
	// generic repo method clean.
	// For the paginated path we call FindByStudent (already sorted) and
	// apply in-process pagination so that no schema changes are required.
	grades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	// Apply in-memory filters (same logic as GetStudentGradesWithFilter)
	var filtered []Grade
	for _, g := range grades {
		if g.DeletedAt != nil {
			continue
		}
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

	total := len(filtered)
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	if offset > total {
		offset = total
	}
	end := offset + pageSize
	if end > total {
		end = total
	}

	totalPages := (total + pageSize - 1) / pageSize
	if total == 0 {
		totalPages = 0
	}

	slice := filtered[offset:end]
	return &PaginatedGradesResponse{
		Data:       s.mapToResponse(slice),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *service) GetClassGrades(ctx context.Context, actorID string, actorRole string, classID string, filter GradeFilter) (*ClassGradesResponse, error) {
	// Permission check
	if actorRole != "admin" && actorRole != "secretary" {
		var coordID sql.NullString
		if err := s.validator.db.QueryRow(
			`SELECT coordinator_id FROM classes WHERE id = $1`, classID,
		).Scan(&coordID); err != nil {
			return nil, fmt.Errorf("class not found")
		}

		if !coordID.Valid || coordID.String != actorID {
			if filter.SubjectID == "" {
				return nil, fmt.Errorf("unauthorized: only coordinators can see full class matrix")
			}
			// If a specific subject is requested, verify the actor teaches it.
			// Propagate any DB error instead of silently treating it as "not found".
			var exists bool
			if err := s.validator.db.QueryRow(
				`SELECT EXISTS(
					SELECT 1 FROM class_subjects cs
					JOIN teachers t ON cs.teacher_id = t.id
					WHERE cs.class_id = $1 AND cs.subject_id = $2 AND t.user_id = $3
				)`,
				classID, filter.SubjectID, actorID,
			).Scan(&exists); err != nil {
				return nil, fmt.Errorf("authorization check failed: %w", err)
			}
			if !exists {
				return nil, fmt.Errorf("unauthorized: you do not teach this subject in this class")
			}
		}
	}

	grades, err := s.repo.FindByClassAndSubject(classID, filter.SubjectID, filter.Semester)
	if err != nil {
		return nil, err
	}

	// Group by student
	studentMap := make(map[string][]GradeResponse)
	for _, g := range grades {
		if filter.IsPublished != nil && g.IsPublished != *filter.IsPublished {
			continue
		}
		studentMap[g.StudentID] = append(studentMap[g.StudentID], s.mapSingleResponse(g))
	}

	resp := &ClassGradesResponse{ClassID: classID}
	studentUsers, _ := s.userRepo.GetStudentsByClass(ctx, classID)
	addedStudents := make(map[string]bool)

	for _, u := range studentUsers {
		sID := u.StudentID
		gList := studentMap[sID]
		if gList == nil {
			gList = []GradeResponse{}
		}
		avg1, avg2 := calcSemesterAverages(gList)
		resp.Students = append(resp.Students, StudentGradeSummary{
			StudentID:    sID,
			FullName:     u.FirstName + " " + u.LastName,
			AvgSemester1: avg1,
			AvgSemester2: avg2,
			Grades:       gList,
		})
		addedStudents[sID] = true
	}

	// Fallback for students who have grades but are not in GetStudentsByClass result
	for sID, gList := range studentMap {
		if addedStudents[sID] {
			continue
		}
		avg1, avg2 := calcSemesterAverages(gList)
		resp.Students = append(resp.Students, StudentGradeSummary{
			StudentID:    sID,
			FullName:     "Student (" + sID + ")",
			AvgSemester1: avg1,
			AvgSemester2: avg2,
			Grades:       gList,
		})
	}

	return resp, nil
}

// calcSemesterAverages computes per-semester unweighted averages from a grade list.
func calcSemesterAverages(grades []GradeResponse) (avg1, avg2 float64) {
	var sum1, sum2 float64
	var count1, count2 int
	for _, g := range grades {
		if g.GradeCategory != "" && g.GradeCategory != string(GradeCategorySummative) {
			continue
		}
		if g.GradeValue >= 0 {
			if g.Semester == 1 {
				sum1 += g.GradeValue
				count1++
			} else {
				sum2 += g.GradeValue
				count2++
			}
		}
	}
	if count1 > 0 {
		avg1 = sum1 / float64(count1)
	}
	if count2 > 0 {
		avg2 = sum2 / float64(count2)
	}
	return
}

func (s *service) GetSubjectGrades(ctx context.Context, actorID string, actorRole string, subjectID string, filter GradeFilter) (*SubjectStatsResponse, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" {
		if actorRole != "teacher" {
			return nil, ErrUnauthorized
		}
		var exists bool
		if err := s.validator.db.QueryRowContext(
			ctx,
			`SELECT EXISTS(
				SELECT 1 FROM class_subjects cs
				JOIN teachers t ON cs.teacher_id = t.id
				WHERE cs.subject_id = $1 AND t.user_id = $2
			)`,
			subjectID, actorID,
		).Scan(&exists); err != nil {
			return nil, fmt.Errorf("authorization check failed: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("%w: you are not assigned to teach this subject", ErrUnauthorized)
		}
	}

	grades, err := s.repo.FindBySubject(subjectID, filter.Semester)
	if err != nil {
		return nil, err
	}

	stat := &SubjectStatsResponse{
		Subject: SubjectInfo{ID: subjectID},
	}

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
		switch {
		case g.GradeValue < 4:
			dist["0-3"]++
		case g.GradeValue < 7:
			dist["4-6"]++
		case g.GradeValue < 9:
			dist["7-8"]++
		default:
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

func (s *service) AddGrade(teacherID string, req CreateGradeRequest) (*GradeResponse, error) {
	if err := s.validator.ValidateCreateRequest(req, teacherID); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	teacherUser, err := s.userRepo.GetByID(context.Background(), teacherID)
	if err != nil {
		return nil, fmt.Errorf("could not resolve teacher profile: %w", err)
	}
	if teacherUser.SchoolID == nil {
		return nil, fmt.Errorf("teacher is not associated with a school")
	}
	schoolID := *teacherUser.SchoolID

	var teacherProfileID string
	if err := s.validator.db.QueryRow(
		`SELECT id FROM teachers WHERE user_id = $1`, teacherID,
	).Scan(&teacherProfileID); err != nil {
		return nil, fmt.Errorf("could not resolve teacher profile ID: %w", err)
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("formato data non valido: %w", err)
	}

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
		return nil, fmt.Errorf("failed to create grade: %w", err)
	}

	resp := s.mapSingleResponse(*grade)
	if s.broadcaster != nil {
		s.broadcaster.BroadcastToUser(grade.StudentID, "GRADE_ADDED", resp)
	}

	return &resp, nil
}

func (s *service) BatchCreateGrades(teacherID string, grades []*Grade) error {
	return s.repo.BatchCreate(grades)
}

func (s *service) BulkImport(teacherID string, file io.Reader, semester int) (*ImportResult, error) {
	// ParseCSVGrades now receives the caller-chosen semester so that grades
	// are assigned to the correct period instead of a CSV-embedded or hardcoded default.
	reqs, err := ParseCSVGrades(file, semester)
	if err != nil {
		return nil, err
	}

	res, err := ProcessBulkImport(s.repo, reqs, teacherID)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *service) Export(teacherID string, filter GradeFilter, format string) ([]byte, string, error) {
	grades, err := s.repo.FindWithFilter(filter)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch grades for export: %w", err)
	}

	if format == "json" {
		data, err := ExportToJSON(grades)
		if err != nil {
			return nil, "", fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return data, "application/json", nil
	}

	if format == "csv" {
		out := "StudentID,SubjectID,GradeValue,Date,Semester,TeacherID\n"
		for _, g := range grades {
			out += fmt.Sprintf("%s,%s,%.2f,%s,%d,%s\n",
				g.StudentID, g.SubjectID, g.GradeValue,
				g.Date.Format("2006-01-02"), g.Semester, g.TeacherID)
		}
		return []byte(out), "text/csv", nil
	}

	return nil, "", fmt.Errorf("unsupported format: %s", format)
}

func (s *service) UpdateGrade(teacherID string, gradeID string, req UpdateGradeRequest) (*GradeResponse, error) {
	grade, err := s.repo.FindByID(gradeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch grade: %w", err)
	}
	if grade == nil {
		return nil, fmt.Errorf("grade not found")
	}

	// Resolve the teacher's profile ID (teachers.id) from the user ID (users.id)
	// so we compare the same domain — grade.TeacherID references teachers.id.
	var teacherProfileID string
	if err := s.validator.db.QueryRow(
		`SELECT id FROM teachers WHERE user_id = $1`, teacherID,
	).Scan(&teacherProfileID); err != nil {
		return nil, fmt.Errorf("could not resolve teacher profile: %w", err)
	}

	if grade.TeacherID != teacherProfileID {
		return nil, fmt.Errorf("unauthorized: can only modify own grades")
	}

	if err := s.validator.ValidateModification(*grade, req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
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
		if parseErr != nil {
			return nil, fmt.Errorf("invalid date format: %w", parseErr)
		}
		grade.Date = parsedDate
		changes = true
	}

	if !changes {
		resp := s.mapSingleResponse(*grade)
		return &resp, nil
	}

	grade.ModifiedBy = &teacherID

	if err := s.repo.Update(grade, history); err != nil {
		return nil, fmt.Errorf("failed to update grade: %w", err)
	}

	resp := s.mapSingleResponse(*grade)
	return &resp, nil
}

func (s *service) DeleteGrade(teacherID string, gradeID string) error {
	grade, err := s.repo.FindByID(gradeID)
	if err != nil {
		return fmt.Errorf("failed to fetch grade: %w", err)
	}
	if grade == nil {
		return fmt.Errorf("grade not found")
	}

	// Resolve teacher profile ID for correct ownership comparison.
	var teacherProfileID string
	if err := s.validator.db.QueryRow(
		`SELECT id FROM teachers WHERE user_id = $1`, teacherID,
	).Scan(&teacherProfileID); err != nil {
		return fmt.Errorf("could not resolve teacher profile: %w", err)
	}

	if grade.TeacherID != teacherProfileID {
		return fmt.Errorf("unauthorized to delete this grade")
	}

	if err := s.validator.ValidateModification(*grade, UpdateGradeRequest{}); err != nil {
		return err
	}

	if err := s.repo.Delete(gradeID, teacherID); err != nil {
		return fmt.Errorf("failed to delete grade: %w", err)
	}
	return nil
}

// --- Student / Parent ---

func (s *service) GetChildGrades(parentID string, studentID string, filter GradeFilter) (*MyGradesResponse, error) {
	ctx := context.Background()
	logger.Log.Debugf("GetChildGrades requested by parent")

	if parentID == "" || studentID == "" {
		return nil, ErrNotGuardian
	}

	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
	if err != nil {
		logger.Log.Errorf("IsGuardian error: %v", err)
		return nil, err
	}
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
	published := true
	filter.IsPublished = &published

	allGrades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var validGrades []Grade
	for _, g := range allGrades {
		if !g.IsPublished || g.DeletedAt != nil {
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

	response := &MyGradesResponse{
		Student: StudentInfo{ID: studentID},
	}

	semestersMap := make(map[int][]GradeResponse)
	for _, g := range validGrades {
		semestersMap[int(g.Semester)] = append(semestersMap[int(g.Semester)], s.mapSingleResponse(g))
	}

	sem1Start, sem1End, sem2Start, sem2End := academicYearDates()
	if gr, ok := semestersMap[1]; ok {
		response.Semesters = append(response.Semesters, SemesterGradesSummary{
			Semester:  1,
			StartDate: sem1Start,
			EndDate:   sem1End,
			Grades:    gr,
		})
	}
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
	grades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var sem1Grades, sem2Grades []Grade
	for _, g := range grades {
		if !g.IsPublished || g.DeletedAt != nil {
			continue
		}
		if g.GradeCategory != GradeCategorySummative {
			continue
		}
		if g.Semester == 1 {
			sem1Grades = append(sem1Grades, g)
		}
		if g.Semester == 2 {
			sem2Grades = append(sem2Grades, g)
		}
	}

	calcSemesterAvg := func(gs []Grade) SemesterAverageSummary {
		subMap := make(map[string][]Grade)
		for _, g := range gs {
			subMap[g.SubjectID] = append(subMap[g.SubjectID], g)
		}

		var subjects []SubjectAverage
		var totalSum float64
		var totalSub int

		for subID, subGrades := range subMap {
			avg := s.calculator.CalculateAverage(subGrades)
			weightedAvg := s.calculator.CalculateWeightedAverage(subGrades)
			subjects = append(subjects, SubjectAverage{
				Subject:         subID,
				Average:         avg,
				WeightedAverage: weightedAvg,
				TotalGrades:     len(subGrades),
			})
			totalSum += weightedAvg
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

	sort.Slice(relevant, func(i, j int) bool {
		return relevant[i].Date.Before(relevant[j].Date)
	})

	classAverage := 0.0
	var classID string
	if err := s.validator.db.QueryRow(
		`SELECT class_id FROM class_students WHERE student_id = $1`, studentID,
	).Scan(&classID); err == nil && classID != "" {
		_ = s.validator.db.QueryRow(
			`SELECT COALESCE(AVG(grade_value), 0.0)
			 FROM grades g
			 JOIN class_students cs ON g.student_id = cs.student_id
			 WHERE cs.class_id = $1 AND g.subject_id = $2
			   AND g.is_published = true AND g.deleted_at IS NULL`,
			classID, subjectID,
		).Scan(&classAverage)
		classAverage = math.Round(classAverage*100) / 100
	}

	var points []TrendPoint
	for i, g := range relevant {
		start := i - 2
		if start < 0 {
			start = 0
		}
		subset := relevant[start : i+1]
		movingAvg := s.calculator.CalculateAverage(subset)

		// Compare the smoothed moving average against the class average,
		// not the raw grade value which may be a noisy outlier.
		position := "at_average"
		if movingAvg > classAverage {
			position = "above_average"
		} else if movingAvg < classAverage {
			position = "below_average"
		}

		points = append(points, TrendPoint{
			Date:         g.Date.Format("2006-01-02"),
			Grade:        g.GradeValue,
			MovingAvg3:   movingAvg,
			ClassAverage: classAverage,
			Position:     position,
		})
	}

	direction := "stable"
	if len(points) >= 2 {
		last := points[len(points)-1]
		prev := points[len(points)-2]
		diff := last.MovingAvg3 - prev.MovingAvg3
		if diff > 0.5 {
			direction = "improving"
		} else if diff < -0.5 {
			direction = "declining"
		}
	}

	// Context-aware recommendation instead of a hardcoded English string.
	var recommendation string
	switch direction {
	case "improving":
		recommendation = "Ottimo progresso! Continua così."
	case "declining":
		recommendation = "Il rendimento è in calo. Si consiglia di rivedere gli argomenti recenti."
	default:
		recommendation = "Rendimento stabile. Mantieni la costanza."
	}

	return &TrendResponse{
		Subject: subjectID,
		Trends:  points,
		Summary: TrendSummary{
			TrendDirection: direction,
			Recommendation: recommendation,
		},
	}, nil
}

func (s *service) GetSemesterReport(studentID string, semester int) (*SemesterReportResponse, error) {
	grades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var studentName, className, schoolYear string
	if s.validator != nil && s.validator.db != nil {
		_ = s.validator.db.QueryRow(
			`SELECT u.first_name || ' ' || u.last_name, COALESCE(c.name, 'N/D')
			 FROM users u
			 LEFT JOIN class_students cs ON u.id = cs.student_id
			 LEFT JOIN classes c ON cs.class_id = c.id
			 WHERE u.id = $1`, studentID,
		).Scan(&studentName, &className)
	}
	if studentName == "" {
		studentName = "Studente " + studentID
	}
	if className == "" {
		className = "Classe N/D"
	}
	schoolYear = currentSchoolYear()

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

	enrolledSubjects, enrollErr := s.repo.FindEnrolledSubjects(studentID, semester)

	var subjects []SubjectReport
	totalSum := 0.0
	passedCount := 0

	processedSubjects := make(map[string]bool)

	for subID, gs := range subMap {
		avg := s.calculator.CalculateWeightedAverage(gs)
		totalSum += avg

		var gVals []GradeVal
		for _, g := range gs {
			gVals = append(gVals, GradeVal{
				Value:    g.GradeValue,
				Date:     g.Date,
				Category: string(g.GradeCategory),
			})
		}

		passed := avg >= 6.0
		if passed {
			passedCount++
		}

		finalGrade := math.Round(avg)
		if finalGrade < 1 {
			finalGrade = 1
		}

		subjects = append(subjects, SubjectReport{
			Subject:        subID,
			SubjectID:      subID,
			Teacher:        "Docente",
			FinalGrade:     finalGrade,
			SubjectAverage: math.Round(avg*100) / 100,
			GradeCount:     len(gs),
			AbsenceDays:    0,
			Notes:          "",
			Grades:         gVals,
			Passed:         passed,
		})
		processedSubjects[subID] = true
	}

	if enrollErr == nil {
		for _, subID := range enrolledSubjects {
			if processedSubjects[subID] {
				continue
			}
			subjects = append(subjects, SubjectReport{
				Subject:        subID,
				SubjectID:      subID,
				Teacher:        "Docente",
				FinalGrade:     0,
				SubjectAverage: 0,
				GradeCount:     0,
				AbsenceDays:    0,
				Notes:          "",
				Grades:         []GradeVal{},
				Passed:         false,
			})
		}
	}

	overall := 0.0
	if len(subjects) > 0 {
		overall = math.Round((totalSum/float64(len(subjects)))*100) / 100
	}

	totalSubjectCount := len(subjects)
	promoted := "NO"
	if enrollErr == nil && len(enrolledSubjects) > 0 && totalSubjectCount > 0 && passedCount == len(enrolledSubjects) && passedCount == totalSubjectCount {
		promoted = "SÌ"
	}

	return &SemesterReportResponse{
		Semester:         semester,
		StudentName:      studentName,
		ClassName:        className,
		SchoolYear:       schoolYear,
		BehaviorGrade:    8.0,
		ScholasticCredit: 8.0,
		OverallAverage:   overall,
		TotalAbsenceDays: 0,
		Subjects:         subjects,
		Promoted:         promoted,
		Status:           "OK",
		LastUpdate:       time.Now(),
	}, nil
}

func (s *service) GetChildSemesterReport(ctx context.Context, parentID, studentID string, semester int) (*SemesterReportResponse, error) {
	if s.userRepo != nil {
		isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
		if err != nil {
			return nil, err
		}
		if !isGuardian {
			return nil, ErrNotGuardian
		}
	}
	return s.GetSemesterReport(studentID, semester)
}

func (s *service) GenerateSemesterReportPDF(studentID string, semester int) ([]byte, error) {
	report, err := s.GetSemesterReport(studentID, semester)
	if err != nil {
		return nil, err
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.CellFormat(190, 10, "REGISTRO ELETTRONICO - SCHEDA DI VALUTAZIONE", "0", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	semText := fmt.Sprintf("%d° Quadrimestre", semester)
	pdf.CellFormat(190, 8, fmt.Sprintf("Valutazione Finale - %s - A.S. %s", semText, report.SchoolYear), "0", 1, "C", false, 0, "")
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(95, 7, fmt.Sprintf("Studente: %s", report.StudentName), "1", 0, "L", false, 0, "")
	pdf.CellFormat(95, 7, fmt.Sprintf("Classe: %s", report.ClassName), "1", 1, "L", false, 0, "")
	pdf.CellFormat(95, 7, fmt.Sprintf("Media Generale: %.2f", report.OverallAverage), "1", 0, "L", false, 0, "")
	pdf.CellFormat(95, 7, fmt.Sprintf("Voto Comportamento: %.0f", report.BehaviorGrade), "1", 1, "L", false, 0, "")
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(60, 8, "Materia", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 8, "Docente", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "Voti (N°)", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Media", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Voto Finale", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 10)
	for _, sub := range report.Subjects {
		subName := sub.Subject
		if len(subName) > 25 {
			subName = subName[:22] + "..."
		}
		teacherName := sub.Teacher
		if len(teacherName) > 18 {
			teacherName = teacherName[:15] + "..."
		}

		pdf.CellFormat(60, 7, subName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 7, teacherName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%d", sub.GradeCount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", sub.SubjectAverage), "1", 0, "C", false, 0, "")

		finalGradeStr := fmt.Sprintf("%.0f", sub.FinalGrade)
		if sub.FinalGrade == 0 {
			finalGradeStr = fmt.Sprintf("%.1f", sub.SubjectAverage)
		}
		pdf.CellFormat(30, 7, finalGradeStr, "1", 1, "C", false, 0, "")
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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

func currentSchoolYear() string {
	now := time.Now()
	year := now.Year()
	if now.Month() < time.September {
		year--
	}
	return fmt.Sprintf("%d/%d", year, year+1)
}

func (s *service) CreateTestWithGrades(teacherID string, req CreateClassTestRequest) error {
	teacherUser, err := s.userRepo.GetByID(context.Background(), teacherID)
	if err != nil {
		return fmt.Errorf("could not resolve teacher profile: %w", err)
	}
	if teacherUser.SchoolID == nil {
		return fmt.Errorf("teacher is not associated with a school")
	}
	schoolID := *teacherUser.SchoolID

	var teacherProfileID string
	if err := s.validator.db.QueryRow(
		`SELECT id FROM teachers WHERE user_id = $1`, teacherID,
	).Scan(&teacherProfileID); err != nil {
		return fmt.Errorf("could not resolve teacher profile ID: %w", err)
	}

	testDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		testDate = time.Now()
	}

	test := &ClassTest{
		ClassID:        req.ClassID,
		SubjectID:      req.SubjectID,
		TeacherID:      teacherID,
		Title:          req.Title,
		Date:           testDate,
		TeacherNotes:   req.TeacherNotes,
		ParentNotes:    req.ParentNotes,
		EvaluationType: req.EvaluationType,
	}

	if err := s.repo.CreateTest(test); err != nil {
		return fmt.Errorf("failed to create test: %w", err)
	}

	validSemester := func(sem int) int {
		if sem == 1 || sem == 2 {
			return sem
		}
		return 1
	}

	var gradesList []*Grade
	for _, gInput := range req.Grades {
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

		val := EvaluationType(req.EvaluationType)
		evalType := &val

		grade := &Grade{
			StudentID:      gInput.StudentID,
			SubjectID:      req.SubjectID,
			TeacherID:      teacherProfileID,
			SchoolID:       schoolID,
			GradeValue:     *gInput.GradeValue,
			GradeType:      "numeric",
			Semester:       Semester(validSemester(req.Semester)),
			Date:           testDate,
			Description:    desc,
			Weight:         1.0,
			IsPublished:    false, // drafts: teacher must explicitly publish
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
	var teacherProfileID string
	if s.validator != nil && s.validator.db != nil {
		_ = s.validator.db.QueryRow(
			`SELECT id FROM teachers WHERE user_id = $1`, teacherID,
		).Scan(&teacherProfileID)
	}
	if test.TeacherID != teacherID && test.TeacherID != teacherProfileID {
		return ErrUnauthorized
	}
	return s.repo.DeleteTest(testID)
}

func (s *service) UpdateClassTest(teacherID string, testID string, req UpdateClassTestRequest) error {
	teacherUser, err := s.userRepo.GetByID(context.Background(), teacherID)
	if err != nil {
		return fmt.Errorf("could not resolve teacher profile: %w", err)
	}
	if teacherUser.SchoolID == nil {
		return fmt.Errorf("teacher is not associated with a school")
	}
	schoolID := *teacherUser.SchoolID

	var teacherProfileID string
	if err := s.validator.db.QueryRow(
		`SELECT id FROM teachers WHERE user_id = $1`, teacherID,
	).Scan(&teacherProfileID); err != nil {
		return fmt.Errorf("could not resolve teacher profile ID: %w", err)
	}

	var classID, subjectID string
	if err := s.validator.db.QueryRow(
		`SELECT class_id, subject_id FROM class_tests WHERE id = $1`, testID,
	).Scan(&classID, &subjectID); err != nil {
		return fmt.Errorf("could not resolve test details: %w", err)
	}

	testDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return fmt.Errorf("invalid test date format '%s': %w", req.Date, err)
	}

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

	existingGrades, err := s.repo.FindGradesByTestID(testID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing grades for test: %w", err)
	}

	existingMap := make(map[string]Grade)
	for _, eg := range existingGrades {
		existingMap[eg.StudentID] = eg
	}

	validSemester := func(sem int) int {
		if sem == 1 || sem == 2 {
			return sem
		}
		return 1
	}

	for _, gInput := range req.Grades {
		existingGrade, exists := existingMap[gInput.StudentID]

		if gInput.GradeValue == nil || *gInput.GradeValue < -1 {
			if exists {
				if err := s.repo.Delete(existingGrade.ID, teacherID); err != nil {
					return fmt.Errorf("failed to delete grade for student %s: %w", gInput.StudentID, err)
				}
			}
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

		val := EvaluationType(req.EvaluationType)
		evalType := &val

		if exists {
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
			existingGrade.Semester = Semester(validSemester(req.Semester))

			if err := s.repo.Update(&existingGrade, history); err != nil {
				return fmt.Errorf("failed to update grade for student %s: %w", gInput.StudentID, err)
			}
		} else {
			grade := &Grade{
				StudentID:      gInput.StudentID,
				SubjectID:      subjectID,
				TeacherID:      teacherProfileID,
				SchoolID:       schoolID,
				GradeValue:     *gInput.GradeValue,
				GradeType:      "numeric",
				Semester:       Semester(validSemester(req.Semester)),
				Date:           testDate,
				Description:    desc,
				Weight:         1.0,
				IsPublished:    false,
				GradeCategory:  "summative",
				EvaluationType: evalType,
				CreatedBy:      teacherID,
				TestID:         &testID,
			}
			if err := s.repo.BatchCreate([]*Grade{grade}); err != nil {
				return fmt.Errorf("failed to insert new grade for student %s: %w", gInput.StudentID, err)
			}
		}
	}

	return nil
}

// --- Grade Weight Config ---

// GetWeightConfigs returns the weight configurations for a given school/subject/class.
func (s *service) GetWeightConfigs(schoolID, subjectID, classID string) ([]GradeWeightConfig, error) {
	configs, err := s.repo.GetWeightConfigs(schoolID, subjectID, classID)
	if err != nil {
		return nil, fmt.Errorf("GetWeightConfigs: %w", err)
	}
	if configs == nil {
		configs = []GradeWeightConfig{}
	}
	return configs, nil
}

// UpsertWeightConfig creates or updates a weight configuration entry.
func (s *service) UpsertWeightConfig(actorID, schoolID string, req UpsertWeightConfigRequest) (*GradeWeightConfig, error) {
	if schoolID == "" {
		return nil, fmt.Errorf("school_id is required")
	}
	validCategories := map[string]bool{"formative": true, "summative": true, "practical": true}
	if !validCategories[req.GradeCategory] {
		return nil, fmt.Errorf("invalid grade_category: must be formative, summative, or practical")
	}
	if req.EvaluationType != nil {
		validTypes := map[string]bool{"Written": true, "Oral": true, "Practical": true}
		if !validTypes[*req.EvaluationType] {
			return nil, fmt.Errorf("invalid evaluation_type: must be Written, Oral, or Practical")
		}
	}
	cfg := &GradeWeightConfig{
		SchoolID:       schoolID,
		SubjectID:      req.SubjectID,
		ClassID:        req.ClassID,
		GradeCategory:  req.GradeCategory,
		EvaluationType: req.EvaluationType,
		Weight:         req.Weight,
		CreatedBy:      actorID,
	}
	return s.repo.UpsertWeightConfig(cfg)
}

// DeleteWeightConfig removes a weight configuration entry.
func (s *service) DeleteWeightConfig(actorID, configID string) error {
	if configID == "" {
		return fmt.Errorf("config id required")
	}
	return s.repo.DeleteWeightConfig(configID)
}

// effectiveWeight returns the resolved weight for a grade, applying configured defaults.
// Priority: grade.Weight (if != 1.0) > matching config > 1.0 default.
func effectiveWeight(g Grade, configs []GradeWeightConfig) float64 {
	if g.Weight != 1.0 {
		return g.Weight
	}
	var bestWeight *float64
	bestScore := -1
	for _, c := range configs {
		if string(g.GradeCategory) != c.GradeCategory {
			continue
		}
		if c.EvaluationType != nil && g.EvaluationType != nil && string(*g.EvaluationType) != *c.EvaluationType {
			continue
		}
		score := 0
		if c.SubjectID != nil && *c.SubjectID == g.SubjectID {
			score += 2
		}
		if c.ClassID != nil {
			score++
		}
		if score > bestScore {
			bestScore = score
			w := c.Weight
			bestWeight = &w
		}
	}
	if bestWeight != nil {
		return *bestWeight
	}
	return 1.0
}

