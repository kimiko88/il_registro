package grades

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"time"

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
	BroadcastToUser(userID string, schoolID string, msgType string, payload interface{})
	BroadcastToSchool(schoolID string, msgType string, payload interface{})
}

type Service interface {
	GetStudentGrades(ctx context.Context, actorID string, actorRole string, studentID string) ([]GradeResponse, error)
	GetStudentGradesWithFilter(ctx context.Context, actorID string, actorRole string, studentID string, filter GradeFilter) ([]GradeResponse, error)
	// GetStudentGradesPaged returns a paginated response for list-grades endpoints.
	GetStudentGradesPaged(ctx context.Context, actorID string, actorRole string, studentID string, filter GradeFilter) (*PaginatedGradesResponse, error)
	GetClassGrades(ctx context.Context, actorID string, actorRole string, classID string, filter GradeFilter) (*ClassGradesResponse, error)
	GetSubjectGrades(ctx context.Context, actorID string, actorRole string, subjectID string, filter GradeFilter) (*SubjectStatsResponse, error)
	AddGrade(ctx context.Context, teacherID string, req CreateGradeRequest) (*GradeResponse, error)
	BatchCreateGrades(teacherID, actorRole, schoolID string, grades []*Grade) error
	BulkImport(teacherID, schoolID string, r io.Reader, semester int) (*ImportResult, error)
	Export(teacherID, schoolID string, filter GradeFilter, format string) ([]byte, string, error)
	UpdateGrade(ctx context.Context, teacherID string, gradeID string, req UpdateGradeRequest) (*GradeResponse, error)
	DeleteGrade(ctx context.Context, teacherID string, gradeID string) error

	// Student/Parent
	GetMyGrades(ctx context.Context, studentID string, filter GradeFilter) (*MyGradesResponse, error)
	GetMyAverages(ctx context.Context, studentID string) (*StudentAveragesResponse, error)
	GetMyTrend(ctx context.Context, actorID string, actorRole string, studentID string, subjectID string) (*TrendResponse, error)
	GetSemesterReport(ctx context.Context, actorID string, actorRole string, studentID string, semester int) (*SemesterReportResponse, error)
	GenerateSemesterReportPDF(ctx context.Context, actorID, actorRole, studentID string, semester int) ([]byte, error)

	// Parent
	ValidateParentGuardian(ctx context.Context, parentID, studentID string) error
	GetChildGrades(ctx context.Context, parentID string, studentID string, filter GradeFilter) (*MyGradesResponse, error)
	GetChildAverages(ctx context.Context, parentID string, studentID string) (*StudentAveragesResponse, error)
	GetChildSemesterReport(ctx context.Context, parentID, studentID string, semester int) (*SemesterReportResponse, error)

	// Class Tests
	CreateTestWithGrades(teacherID string, req CreateClassTestRequest) (*ClassTest, error)
	GetClassTests(ctx context.Context, actorID string, actorRole string, classID string, subjectID string) ([]ClassTestResponse, error)
	GetUpcomingTestsByClass(ctx context.Context, actorID string, actorRole string, classID string) ([]ClassTestResponse, error)
	DeleteClassTest(teacherID string, testID string) error
	UpdateClassTest(teacherID string, testID string, req UpdateClassTestRequest) error

	// Weight Config
	GetWeightConfigs(schoolID, subjectID, classID string) ([]GradeWeightConfig, error)
	UpsertWeightConfig(actorID, actorRole, schoolID string, req UpsertWeightConfigRequest) (*GradeWeightConfig, error)
	DeleteWeightConfig(actorID, actorRole, schoolID, configID string) error
}

type service struct {
	repo        Repository
	userRepo    users.Repository
	validator   *Validator
	calculator  *Calculator
	broadcaster EventBroadcaster
	pdfExporter ReportCardPDFExporter
}

func NewService(r Repository, ur users.Repository, db *sql.DB, b EventBroadcaster) Service {
	return &service{
		repo:        r,
		userRepo:    ur,
		validator:   NewValidator(db),
		calculator:  NewCalculator(),
		broadcaster: b,
		pdfExporter: NewReportCardPDFExporter(),
	}
}

func (s *service) GetStudentGrades(ctx context.Context, actorID, actorRole, studentID string) ([]GradeResponse, error) {
	return s.GetStudentGradesWithFilter(ctx, actorID, actorRole, studentID, GradeFilter{})
}

func (s *service) checkGradeAccessPermissions(ctx context.Context, actorID string, actorRole string, studentID string) error {
	if actorRole == "" {
		return ErrUnauthorized
	}
	if actorRole == "student" && actorID != studentID {
		return ErrUnauthorized
	}
	if actorRole == "parent" {
		isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, studentID)
		if err != nil {
			return err
		}
		if !isGuardian {
			return ErrNotGuardian
		}
	}
	if (actorRole == "teacher" || actorRole == "coordinator") && s.validator != nil {
		isAssigned, err := s.validator.IsTeacherAssignedToStudent(ctx, actorID, studentID)
		if err != nil {
			return err
		}
		if !isAssigned {
			return ErrUnauthorized
		}
	}
	return nil
}

func (s *service) GetStudentGradesWithFilter(ctx context.Context, actorID string, actorRole string, studentID string, filter GradeFilter) ([]GradeResponse, error) {
	if err := s.checkGradeAccessPermissions(ctx, actorID, actorRole, studentID); err != nil {
		return nil, err
	}

	if filter.Semester == 0 && filter.SubjectID == "" && filter.GradeType == "" && filter.IsPublished == nil {
		grades, err := s.repo.FindByStudent(studentID)
		if err != nil {
			return nil, err
		}
		return s.mapToResponse(grades), nil
	}

	filter.StudentID = studentID
	grades, err := s.repo.FindWithFilter(filter)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(grades), nil
}

// GetStudentGradesPaged applies the same ownership checks as GetStudentGradesWithFilter
// then delegates to the paginated repository method.
func (s *service) GetStudentGradesPaged(ctx context.Context, actorID string, actorRole string, studentID string, filter GradeFilter) (*PaginatedGradesResponse, error) {
	if err := s.checkGradeAccessPermissions(ctx, actorID, actorRole, studentID); err != nil {
		return nil, err
	}

	filter.StudentID = studentID
	grades, total, err := s.repo.FindWithFilterPaginated(filter)
	if err != nil {
		return nil, err
	}

	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}

	return &PaginatedGradesResponse{
		Data:       s.mapToResponse(grades),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *service) GetClassGrades(ctx context.Context, actorID string, actorRole string, classID string, filter GradeFilter) (*ClassGradesResponse, error) {
	// Permission check
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" {
		isCoord, err := s.validator.IsClassCoordinator(actorID, classID)
		if err != nil {
			return nil, fmt.Errorf("class authorization check failed: %w", err)
		}

		// FIX: a teacher who IS coordinator gets the full class view; any other
		// teacher (including a coordinator of a different class, i.e. isCoord==false)
		// must supply a SubjectID and be assigned to it.
		if !isCoord {
			if actorRole != "teacher" && actorRole != "coordinator" {
				return nil, fmt.Errorf("%w: solo il coordinatore di classe o la dirigenza possono accedere al quadro completo della classe", ErrUnauthorized)
			}
			if filter.SubjectID == "" {
				return nil, fmt.Errorf("%w: solo il coordinatore di classe o la dirigenza possono accedere al quadro completo della classe", ErrUnauthorized)
			}
			assigned, err := s.validator.IsTeacherAssignedToSubject(actorID, filter.SubjectID, classID)
			if err != nil {
				return nil, fmt.Errorf("authorization check failed: %w", err)
			}
			if !assigned {
				return nil, fmt.Errorf("%w: non sei assegnato a questa materia per la classe indicata", ErrUnauthorized)
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
	studentUsers, err := s.userRepo.GetStudentsByClass(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch class students: %w", err)
	}
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
	}

	sort.Slice(resp.Students, func(i, j int) bool {
		return resp.Students[i].FullName < resp.Students[j].FullName
	})

	return resp, nil
}

// calcSemesterAverages computes per-semester weighted averages from a grade list.
func calcSemesterAverages(grades []GradeResponse) (avg1, avg2 float64) {
	var weightedSum1, totalWeight1 float64
	var weightedSum2, totalWeight2 float64
	for _, g := range grades {
		// Skip non-summative categories (including empty category which was
		// previously let through and polluted averages with formative grades).
		if g.GradeCategory != string(GradeCategorySummative) {
			continue
		}
		if g.GradeValue > 0 {
			w := g.Weight
			if w <= 0 {
				w = 1.0
			}
			if g.Semester == 1 {
				weightedSum1 += g.GradeValue * w
				totalWeight1 += w
			} else {
				weightedSum2 += g.GradeValue * w
				totalWeight2 += w
			}
		}
	}
	if totalWeight1 > 0 {
		avg1 = math.Round((weightedSum1/totalWeight1)*100) / 100
	}
	if totalWeight2 > 0 {
		avg2 = math.Round((weightedSum2/totalWeight2)*100) / 100
	}
	return
}

func (s *service) GetSubjectGrades(ctx context.Context, actorID string, actorRole string, subjectID string, filter GradeFilter) (*SubjectStatsResponse, error) {
	// FIX: principal and vice_principal added to the bypass whitelist, consistent
	// with the policy already in place for GetClassGrades.
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" &&
		actorRole != "principal" && actorRole != "vice_principal" {
		if actorRole != "teacher" {
			return nil, ErrUnauthorized
		}
		// FIX: use validator method instead of inline SQL to avoid nil-panic when s.validator is nil
		// and to respect the repository pattern.
		if s.validator == nil {
			return nil, fmt.Errorf("%w: authorization service unavailable", ErrUnauthorized)
		}
		assigned, err := s.validator.IsTeacherAssignedToSubjectBySubjectID(ctx, actorID, subjectID)
		if err != nil {
			return nil, fmt.Errorf("authorization check failed: %w", err)
		}
		if !assigned {
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

	var validGrades []Grade
	for _, g := range grades {
		if g.DeletedAt != nil || !g.IsPublished {
			continue
		}
		val := g.GradeValue
		if val == 0 && g.GradeType == GradeTypeJudgment {
			val = s.calculator.ConvertJudgmentToValue(g.Description)
		}
		if val < 1.0 || val > 10.0 {
			continue
		}
		validGrades = append(validGrades, g)
		if !seenStudents[g.StudentID] {
			seenStudents[g.StudentID] = true
			totalStudents++
		}
		sum += val
		switch {
		case val < 4:
			dist["1-3"]++
		case val < 7:
			dist["4-6"]++
		case val < 9:
			dist["7-8"]++
		default:
			dist["9-10"]++
		}
	}

	var avg float64
	if len(validGrades) > 0 {
		avg = sum / float64(len(validGrades))
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

func (s *service) AddGrade(ctx context.Context, teacherID string, req CreateGradeRequest) (*GradeResponse, error) {
	if err := s.validator.ValidateCreateRequest(req, teacherID); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	teacherUser, err := s.userRepo.GetByID(ctx, teacherID)
	if err != nil {
		return nil, fmt.Errorf("could not resolve teacher profile: %w", err)
	}
	if teacherUser.SchoolID == nil {
		return nil, fmt.Errorf("teacher is not associated with a school")
	}
	schoolID := *teacherUser.SchoolID

	teacherProfileID, err := s.resolveTeacherProfileID(ctx, teacherID)
	if err != nil {
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
		s.broadcaster.BroadcastToUser(grade.StudentID, schoolID, "GRADE_ADDED", resp)
	}

	return &resp, nil
}

func (s *service) BatchCreateGrades(teacherID, actorRole, schoolID string, grades []*Grade) error {
	if actorRole != "teacher" && actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	if len(grades) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	var deduped []*Grade
	for _, g := range grades {
		if g == nil {
			continue
		}
		if actorRole != "superadmin" && schoolID != "" && g.SchoolID != "" && g.SchoolID != schoolID {
			return fmt.Errorf("%w: cannot create grades for another school", ErrUnauthorized)
		}
		evalTypeStr := ""
		if g.EvaluationType != nil {
			evalTypeStr = string(*g.EvaluationType)
		}
		testIDStr := ""
		if g.TestID != nil {
			testIDStr = *g.TestID
		}
		key := fmt.Sprintf("%s_%s_%s_%s_%s_%s_%.2f_%s", g.StudentID, g.SubjectID, testIDStr, g.Date.Format("2006-01-02"), g.GradeType, evalTypeStr, g.GradeValue, g.Description)
		if g.ID != "" {
			key = g.ID + "_" + key
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		deduped = append(deduped, g)
	}
	return s.repo.BatchCreate(deduped)
}

func (s *service) BulkImport(teacherID, schoolID string, file io.Reader, semester int) (*ImportResult, error) {
	// ParseCSVGrades now receives the caller-chosen semester so that grades
	// are assigned to the correct period instead of a CSV-embedded or hardcoded default.
	reqs, err := ParseCSVGrades(file, semester)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	teacherUser, err := s.userRepo.GetByID(ctx, teacherID)
	if err != nil || teacherUser == nil {
		return nil, fmt.Errorf("teacher user not found: %w", err)
	}
	if teacherUser.SchoolID != nil && *teacherUser.SchoolID != "" {
		if schoolID != "" && *teacherUser.SchoolID != schoolID {
			return nil, fmt.Errorf("school ID mismatch with teacher record")
		}
		schoolID = *teacherUser.SchoolID
	}

	teacherProfileID, err := s.resolveTeacherProfileID(ctx, teacherID)
	if err != nil {
		return nil, fmt.Errorf("teacher profile ID is required for bulk import: %w", err)
	}

	// Validate teacher ownership & assignment for each row when teacher is importing
	if teacherUser.Role != "admin" && teacherUser.Role != "superadmin" && s.validator != nil && s.validator.db != nil {
		for idx, req := range reqs {
			assigned, err := s.validator.IsTeacherAssignedToSubjectBySubjectID(ctx, teacherID, req.SubjectID)
			if err != nil || !assigned {
				return nil, fmt.Errorf("row %d: %w: teacher is not assigned to teach subject %s", idx+1, ErrUnauthorized, req.SubjectID)
			}
			assignedStudent, err := s.validator.IsTeacherAssignedToStudent(ctx, teacherID, req.StudentID)
			if err != nil || !assignedStudent {
				return nil, fmt.Errorf("row %d: %w: teacher is not assigned to student %s", idx+1, ErrUnauthorized, req.StudentID)
			}
		}
	}

	res, err := ProcessBulkImport(s.repo, reqs, teacherID, teacherProfileID, schoolID)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *service) Export(teacherID, schoolID string, filter GradeFilter, format string) ([]byte, string, error) {
	if teacherID == "" {
		return nil, "", ErrUnauthorized
	}
	// Force the filter to strictly scope export to the authenticated teacher's grades
	teacherProfileID, err := s.resolveTeacherProfileID(context.Background(), teacherID)
	if err == nil && teacherProfileID != "" {
		filter.TeacherID = teacherProfileID
	} else {
		filter.TeacherID = teacherID
	}
	if schoolID == "" && s.userRepo != nil {
		if tUser, uErr := s.userRepo.GetByID(context.Background(), teacherID); uErr == nil && tUser != nil && tUser.SchoolID != nil {
			schoolID = *tUser.SchoolID
		}
	}
	if schoolID != "" {
		filter.SchoolID = schoolID
	}
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
		data, err := ExportToCSV(grades, ExportOptions{Format: "csv"})
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate CSV: %w", err)
		}
		return data, "text/csv", nil
	}

	if format == "pdf" {
		data, err := ExportToPDF(grades, ExportOptions{Format: "pdf"})
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate PDF: %w", err)
		}
		return data, "application/pdf", nil
	}

	return nil, "", fmt.Errorf("unsupported format: %s", format)
}

func (s *service) UpdateGrade(ctx context.Context, teacherID string, gradeID string, req UpdateGradeRequest) (*GradeResponse, error) {
	grade, err := s.repo.FindByID(gradeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch grade: %w", err)
	}
	if grade == nil {
		return nil, fmt.Errorf("grade not found")
	}

	// Resolve the teacher's profile ID (teachers.id) from the user ID (users.id)
	// so we compare the same domain — grade.TeacherID references teachers.id.
	teacherProfileID, err := s.resolveTeacherProfileID(ctx, teacherID)
	if err != nil {
		return nil, fmt.Errorf("could not resolve teacher profile: %w", err)
	}

	if grade.TeacherID != teacherProfileID {
		return nil, fmt.Errorf("%w: can only modify own grades", ErrUnauthorized)
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

func (s *service) DeleteGrade(ctx context.Context, teacherID string, gradeID string) error {
	grade, err := s.repo.FindByID(gradeID)
	if err != nil {
		return fmt.Errorf("failed to fetch grade: %w", err)
	}
	if grade == nil {
		return fmt.Errorf("grade not found")
	}

	// Resolve teacher profile ID for correct ownership comparison.
	teacherProfileID, err := s.resolveTeacherProfileID(ctx, teacherID)
	if err != nil {
		return fmt.Errorf("could not resolve teacher profile: %w", err)
	}

	if grade.TeacherID != teacherProfileID {
		return fmt.Errorf("%w: unauthorized to delete this grade", ErrUnauthorized)
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

func (s *service) ValidateParentGuardian(ctx context.Context, parentID, studentID string) error {
	if parentID == "" || studentID == "" {
		return ErrNotGuardian
	}
	isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
	if err != nil {
		return err
	}
	if !isGuardian {
		return ErrNotGuardian
	}
	return nil
}

func (s *service) GetChildGrades(ctx context.Context, parentID string, studentID string, filter GradeFilter) (*MyGradesResponse, error) {
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

	return s.GetMyGrades(ctx, studentID, filter)
}

func (s *service) GetChildAverages(ctx context.Context, parentID string, studentID string) (*StudentAveragesResponse, error) {
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

	return s.GetMyAverages(ctx, studentID)
}

func (s *service) GetMyGrades(ctx context.Context, studentID string, filter GradeFilter) (*MyGradesResponse, error) {
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

	var dbConn *sql.DB
	if s.validator != nil {
		dbConn = s.validator.db
	}
	sem1Start, sem1End, sem2Start, sem2End := academicYearDates(ctx, dbConn, filter.SchoolID)
	if gr, ok := semestersMap[1]; ok {
		sort.Slice(gr, func(i, j int) bool {
			return gr[i].Date.Before(gr[j].Date)
		})
		response.Semesters = append(response.Semesters, SemesterGradesSummary{
			Semester:  1,
			StartDate: sem1Start,
			EndDate:   sem1End,
			Grades:    gr,
		})
	}
	if gr, ok := semestersMap[2]; ok {
		sort.Slice(gr, func(i, j int) bool {
			return gr[i].Date.Before(gr[j].Date)
		})
		response.Semesters = append(response.Semesters, SemesterGradesSummary{
			Semester:  2,
			StartDate: sem2Start,
			EndDate:   sem2End,
			Grades:    gr,
		})
	}

	return response, nil
}

func (s *service) GetMyAverages(ctx context.Context, studentID string) (*StudentAveragesResponse, error) {
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

		overall := s.calculator.CalculateWeightedAverage(gs)

		cond := "OTTIMO"
		switch {
		case overall < 6.0:
			cond = "ATTENZIONE"
		case overall < 7.0:
			cond = "SUFFICIENTE"
		case overall < 8.0:
			cond = "BUONO"
		case overall < 9.0:
			cond = "DISTINTO"
		default:
			cond = "OTTIMO"
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

func (s *service) GetMyTrend(ctx context.Context, actorID string, actorRole string, studentID string, subjectID string) (*TrendResponse, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" {
		switch actorRole {
		case "student":
			if actorID != studentID {
				return nil, ErrUnauthorized
			}
		case "parent":
			isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, studentID)
			if err != nil {
				return nil, fmt.Errorf("guardian check failed: %w", err)
			}
			if !isGuardian {
				return nil, ErrNotGuardian
			}
		case "teacher":
			if s.validator != nil {
				assigned, err := s.validator.IsTeacherAssignedToStudent(ctx, actorID, studentID)
				if err != nil || !assigned {
					return nil, ErrUnauthorized
				}
			}
		default:
			return nil, ErrUnauthorized
		}
	}
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

	classAverage := -1.0
	var classID string
	currentSem := 1
	if len(relevant) > 0 && relevant[len(relevant)-1].Semester > 0 {
		currentSem = int(relevant[len(relevant)-1].Semester)
	}
	if s.repo != nil {
		sName, cName, cID, sID, _ := s.repo.GetStudentClassAndSchoolInfo(ctx, studentID)
		_ = sName
		_ = cName
		_ = sID
		classID = cID
		if classID != "" {
			avg, err := s.repo.GetClassSubjectAverage(ctx, classID, subjectID, currentSem, studentID)
			if err == nil && avg >= 0 {
				classAverage = avg
			}
		}
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
		position := "not_available"
		if classAverage >= 0 {
			if movingAvg > classAverage {
				position = "above_average"
			} else if movingAvg < classAverage {
				position = "below_average"
			} else {
				position = "at_average"
			}
		}

		pointClassAvg := classAverage
		if pointClassAvg < 0 {
			pointClassAvg = 0.0
		}

		points = append(points, TrendPoint{
			Date:         g.Date.Format("2006-01-02"),
			Grade:        g.GradeValue,
			MovingAvg3:   movingAvg,
			ClassAverage: pointClassAvg,
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

func (s *service) GetSemesterReport(ctx context.Context, actorID string, actorRole string, studentID string, semester int) (*SemesterReportResponse, error) {
	if actorRole == "" || actorID == "" {
		return nil, ErrUnauthorized
	}
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" {
		switch actorRole {
		case "student":
			if actorID != studentID {
				return nil, ErrUnauthorized
			}
		case "parent":
			isGuardian, err := s.userRepo.IsGuardian(ctx, actorID, studentID)
			if err != nil {
				return nil, fmt.Errorf("guardian check failed: %w", err)
			}
			if !isGuardian {
				return nil, ErrNotGuardian
			}
		case "teacher":
			if s.validator != nil {
				assigned, err := s.validator.IsTeacherAssignedToStudent(ctx, actorID, studentID)
				if err != nil || !assigned {
					return nil, ErrUnauthorized
				}
			}
		default:
			return nil, ErrUnauthorized
		}
	}

	grades, err := s.repo.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var studentName, className, classID, schoolYear, schoolID string
	if s.repo != nil {
		sName, cName, cID, sID, _ := s.repo.GetStudentClassAndSchoolInfo(ctx, studentID)
		studentName, className, classID, schoolID = sName, cName, cID, sID
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

	teacherMap := make(map[string]string)
	if s.repo != nil && classID != "" {
		tNames, tErr := s.repo.GetTeacherNamesByClass(ctx, classID)
		if tErr == nil {
			teacherMap = tNames
		}
	}

	subjectNameMap := make(map[string]string)
	if s.repo != nil {
		sNames, sErr := s.repo.GetSubjectNamesMap(ctx, schoolID)
		if sErr == nil {
			subjectNameMap = sNames
		}
	}

	behaviorGrade := 0.0
	scholasticCredit := 0.0
	totalAbsenceDays := 0

	if s.repo != nil {
		bg, sc, found, _ := s.repo.GetScrutinyRecordSummary(ctx, studentID, semester)
		if found {
			behaviorGrade = bg
			scholasticCredit = sc
		}

		var dbConn *sql.DB
		if s.validator != nil {
			dbConn = s.validator.db
		}
		sem1Start, sem1End, sem2Start, sem2End := academicYearDates(ctx, dbConn, schoolID)
		var startD, endD string
		if semester == 1 {
			startD, endD = sem1Start, sem1End
		} else {
			startD, endD = sem2Start, sem2End
		}
		if absCount, err := s.repo.GetStudentAbsenceCountForPeriod(ctx, studentID, startD, endD); err == nil {
			totalAbsenceDays = absCount
		}
	}

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

		// Final grade rounded to integer on the Italian 1-10 school report scale
		finalGrade := math.Round(avg)
		if finalGrade < 1 {
			finalGrade = 1
		} else if finalGrade > 10 {
			finalGrade = 10
		}

		tName := teacherMap[subID]
		if tName == "" {
			tName = "Docente"
		}

		subName := subjectNameMap[subID]
		if subName == "" {
			subName = subID
		}

		subjects = append(subjects, SubjectReport{
			Subject:        subName,
			SubjectID:      subID,
			Teacher:        tName,
			FinalGrade:     finalGrade,
			SubjectAverage: math.Round(avg*100) / 100,
			GradeCount:     len(gs),
			AbsenceDays:    totalAbsenceDays,
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
			tName := teacherMap[subID]
			if tName == "" {
				tName = "Docente"
			}
			subName := subjectNameMap[subID]
			if subName == "" {
				subName = subID
			}
			subjects = append(subjects, SubjectReport{
				Subject:        subName,
				SubjectID:      subID,
				Teacher:        tName,
				FinalGrade:     0,
				SubjectAverage: 0,
				GradeCount:     0,
				AbsenceDays:    totalAbsenceDays,
				Notes:          "",
				Grades:         []GradeVal{},
				Passed:         false,
			})
		}
	}

	gradedCount := len(subMap)
	totalEnrolled := len(enrolledSubjects)

	denominator := gradedCount
	if totalEnrolled > gradedCount {
		denominator = totalEnrolled
	}
	overall := 0.0
	if denominator > 0 {
		overall = math.Round((totalSum/float64(denominator))*100) / 100
	}

	// Promotion evaluation: all enrolled subjects must be graded, passed, and overall >= 6.0
	promoted := "NO"
	requiredSubjects := totalEnrolled
	if requiredSubjects == 0 && enrollErr == nil {
		requiredSubjects = gradedCount
	}
	if requiredSubjects > 0 && passedCount == requiredSubjects && gradedCount >= requiredSubjects && overall >= 6.0 {
		promoted = "SÌ"
	}

	return &SemesterReportResponse{
		Semester:         semester,
		StudentName:      studentName,
		ClassName:        className,
		SchoolYear:       schoolYear,
		BehaviorGrade:    behaviorGrade,
		ScholasticCredit: scholasticCredit,
		OverallAverage:   overall,
		TotalAbsenceDays: totalAbsenceDays,
		Subjects:         subjects,
		Promoted:         promoted,
		Status:           "OK",
		LastUpdate:       time.Now(),
	}, nil
}

func (s *service) GetChildSemesterReport(ctx context.Context, parentID, studentID string, semester int) (*SemesterReportResponse, error) {
	// FIX: removed redundant explicit IsGuardian call.
	// GetSemesterReport is the single gate-keeper for all roles (including "parent"),
	// so duplicating the check here caused two divergent guardian lookups.
	return s.GetSemesterReport(ctx, parentID, "parent", studentID, semester)
}

func (s *service) GenerateSemesterReportPDF(ctx context.Context, actorID, actorRole, studentID string, semester int) ([]byte, error) {
	report, err := s.GetSemesterReport(ctx, actorID, actorRole, studentID, semester)
	if err != nil {
		return nil, err
	}
	if s.pdfExporter == nil {
		s.pdfExporter = NewReportCardPDFExporter()
	}
	return s.pdfExporter.ExportReportCard(report, semester)
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
		Weight:         g.Weight,
		IsPublished:    g.IsPublished,
		TestID:         g.TestID,
	}
}

// academicYearDates returns semester date boundaries for the current school year,
// reading custom dates from school_settings if configured.
func academicYearDates(ctx context.Context, db *sql.DB, schoolID string) (sem1Start, sem1End, sem2Start, sem2End string) {
	if db != nil && schoolID != "" {
		var s1S, s1E, s2S, s2E sql.NullString
		_ = db.QueryRowContext(ctx,
			`SELECT sem1_start_date, sem1_end_date, sem2_start_date, sem2_end_date FROM school_settings WHERE school_id = $1`, schoolID,
		).Scan(&s1S, &s1E, &s2S, &s2E)
		if s1S.Valid && s1E.Valid && s2S.Valid && s2E.Valid && s1S.String != "" {
			return s1S.String, s1E.String, s2S.String, s2E.String
		}
	}
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
	logger.Log.Warnf("academicYearDates: no custom semester dates configured for schoolID=%q, using hardcoded defaults (%s to %s, %s to %s). Configure school_settings to define accurate term boundaries.", schoolID, sem1Start, sem1End, sem2Start, sem2End)
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

func (s *service) CreateTestWithGrades(teacherID string, req CreateClassTestRequest) (*ClassTest, error) {
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

	testDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid test date format '%s': %w", req.Date, err)
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
		return nil, fmt.Errorf("failed to create test: %w", err)
	}

	validSemester := func(sem int) int {
		if sem == 1 || sem == 2 {
			return sem
		}
		return 1
	}

	var gradesList []*Grade
	for _, gInput := range req.Grades {
		// FIX: -1 is the conventional sentinel for "absent / not evaluated" (see ValidateGradeValue).
		// In the context of CreateTestWithGrades we skip it rather than inserting
		// a -1 grade row, which would be meaningless in aggregate calculations.
		// Any value < 0 is treated as absent/skip.
		if gInput.GradeValue == nil || *gInput.GradeValue < 0 {
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
			IsPublished:    false,
			GradeCategory:  "summative",
			EvaluationType: evalType,
			CreatedBy:      teacherID,
			TestID:         &test.ID,
		}
		gradesList = append(gradesList, grade)
	}

	if len(gradesList) > 0 {
		if err := s.repo.BatchCreate(gradesList); err != nil {
			return nil, fmt.Errorf("failed to insert grades for test: %w", err)
		}
		if s.broadcaster != nil {
			s.broadcaster.BroadcastToSchool(schoolID, "TEST_CREATED", test)
			for _, g := range gradesList {
				if g.IsPublished {
					s.broadcaster.BroadcastToUser(g.StudentID, schoolID, "GRADE_ADDED", s.mapSingleResponse(*g))
				}
			}
		}
	}

	return test, nil
}

// checkClassAccessPermission enforces role-based access control for class-scoped test queries.
// Admins, secretaries, principals, vice-principals and system_auditors bypass the check.
// Teachers must be assigned to the class; students must be enrolled; parents must have a guardian link.
func (s *service) checkClassAccessPermission(ctx context.Context, actorID, actorRole, classID string) error {
	if actorRole == "" || actorID == "" {
		return ErrUnauthorized
	}
	if actorRole == "admin" || actorRole == "superadmin" || actorRole == "secretary" ||
		actorRole == "principal" || actorRole == "vice_principal" || actorRole == "system_auditor" {
		return nil
	}
	if s.validator == nil || s.validator.db == nil {
		return fmt.Errorf("%w: authorization service unavailable", ErrUnauthorized)
	}
	switch actorRole {
	case "teacher":
		var exists bool
		if err := s.validator.db.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM class_subjects cs LEFT JOIN teachers t ON cs.teacher_id = t.id OR cs.teacher_id = t.user_id WHERE cs.class_id::text = $1 AND (cs.teacher_id::text = $2 OR t.user_id::text = $2))`,
			classID, actorID,
		).Scan(&exists); err != nil || !exists {
			return ErrUnauthorized
		}
	case "student":
		var isEnrolled bool
		if err := s.validator.db.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM class_students WHERE class_id::text = $1 AND student_id::text = $2)`,
			classID, actorID,
		).Scan(&isEnrolled); err != nil || !isEnrolled {
			return ErrUnauthorized
		}
	case "parent":
		var isParentGuardian bool
		if err := s.validator.db.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM parent_student_guardians psg JOIN class_students cs ON psg.student_id = cs.student_id WHERE psg.parent_id::text = $1 AND cs.class_id::text = $2)`,
			actorID, classID,
		).Scan(&isParentGuardian); err != nil || !isParentGuardian {
			return ErrUnauthorized
		}
	default:
		return ErrUnauthorized
	}
	return nil
}

func (s *service) GetClassTests(ctx context.Context, actorID string, actorRole string, classID string, subjectID string) ([]ClassTestResponse, error) {
	if err := s.checkClassAccessPermission(ctx, actorID, actorRole, classID); err != nil {
		return nil, err
	}

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

func (s *service) GetUpcomingTestsByClass(ctx context.Context, actorID string, actorRole string, classID string) ([]ClassTestResponse, error) {
	// FIX: replaced duplicated 35-line role-check block with shared helper.
	if err := s.checkClassAccessPermission(ctx, actorID, actorRole, classID); err != nil {
		return nil, err
	}

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
	// FIX: propagate the error instead of silently ignoring it.
	// A DB failure previously set teacherProfileID="" which made the ownership
	// check always deny valid teachers (teacherProfileID == "" branch was true).
	teacherProfileID, err := s.resolveTeacherProfileID(context.Background(), teacherID)
	if err != nil {
		return fmt.Errorf("could not resolve teacher profile for ownership check: %w", err)
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

	test, err := s.repo.FindTestByID(testID)
	if err != nil {
		return fmt.Errorf("could not resolve test details: %w", err)
	}

	if test.TeacherID != teacherID && (teacherProfileID == "" || test.TeacherID != teacherProfileID) {
		return ErrUnauthorized
	}

	// FIX: validate req.Date before attempting to parse to avoid a confusing
	// time.Parse error message when the field is empty.
	if req.Date == "" {
		return fmt.Errorf("test date is required")
	}
	testDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return fmt.Errorf("invalid test date format '%s': %w", req.Date, err)
	}

	test.Title = req.Title
	test.Date = testDate
	test.TeacherNotes = req.TeacherNotes
	test.ParentNotes = req.ParentNotes
	test.EvaluationType = req.EvaluationType

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
				SubjectID:      test.SubjectID,
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
func (s *service) UpsertWeightConfig(actorID, actorRole, schoolID string, req UpsertWeightConfigRequest) (*GradeWeightConfig, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" {
		return nil, ErrUnauthorized
	}
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
func (s *service) DeleteWeightConfig(actorID, actorRole, schoolID, configID string) error {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" {
		return ErrUnauthorized
	}
	if configID == "" {
		return fmt.Errorf("config id required")
	}
	return s.repo.DeleteWeightConfig(configID)
}

func (s *service) resolveTeacherProfileID(ctx context.Context, userID string) (string, error) {
	if s.validator == nil || s.validator.db == nil {
		return userID, nil
	}
	var teacherProfileID string
	err := s.validator.db.QueryRowContext(ctx, `SELECT id FROM teachers WHERE user_id = $1`, userID).Scan(&teacherProfileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userID, nil
		}
		return "", fmt.Errorf("failed to resolve teacher profile ID: %w", err)
	}
	return teacherProfileID, nil
}
