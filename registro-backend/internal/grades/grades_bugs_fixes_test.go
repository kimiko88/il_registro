package grades

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGrades_Handler_TeacherNilValidator_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForDeprecationTest)
	h := NewHandler(mockSvc, nil) // validator is nil

	r := gin.New()
	r.GET("/grades/my-grades/trend", func(c *gin.Context) {
		c.Set("user_id", "t1")
		c.Set("role", "teacher")
		h.GetMyTrend(c)
	})

	req := httptest.NewRequest("GET", "/grades/my-grades/trend?student_id=s2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code, "Teacher accessing another student's trend with nil validator must be rejected with 403")
}

func TestGrades_Handler_GetStudentGrades_DeprecationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForDeprecationTest)
	h := NewHandler(mockSvc, nil)

	mockSvc.On("GetStudentGradesWithFilter", mock.Anything, "s1", "student", "s1", mock.Anything).Return([]GradeResponse{}, nil).Once()

	r := gin.New()
	r.GET("/grades/student/:studentID", func(c *gin.Context) {
		c.Set("user_id", "s1")
		c.Set("role", "student")
		h.GetStudentGrades(c)
	})

	req := httptest.NewRequest("GET", "/grades/student/s1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true", w.Header().Get("Deprecation"))
	assert.Contains(t, w.Header().Get("Link"), "paged")
}

func TestGrades_Service_CoordinatorRoleAccess(t *testing.T) {
	s := &service{validator: nil}

	// Coordinator role should pass permission check
	err := s.checkGradeAccessPermissions(context.Background(), "coord1", "coordinator", "s1")
	assert.NoError(t, err)
}

func TestCalculator_DetectOutliers_UnclampedThresholds(t *testing.T) {
	calc := NewCalculator()

	// High mean cluster with low outlier
	// Mean around 8.8, stdDev around 2.4 -> low bound around 4.0, high bound around 13.6
	// Grade 2.0 should be detected as outlier (< low bound)
	grades := []Grade{
		{ID: "g1", GradeValue: 9.0, GradeType: GradeTypeNumeric},
		{ID: "g2", GradeValue: 9.5, GradeType: GradeTypeNumeric},
		{ID: "g3", GradeValue: 9.0, GradeType: GradeTypeNumeric},
		{ID: "g4", GradeValue: 8.5, GradeType: GradeTypeNumeric},
		{ID: "g5", GradeValue: 9.0, GradeType: GradeTypeNumeric},
		{ID: "g6", GradeValue: 10.0, GradeType: GradeTypeNumeric},
		{ID: "g-outlier", GradeValue: 2.0, GradeType: GradeTypeNumeric},
	}

	outliers := calc.DetectOutliers(grades)
	assert.Contains(t, outliers, "g-outlier")
	assert.Len(t, outliers, 1)
}

func TestCalculator_CalculatePercentile_InvalidScore(t *testing.T) {
	calc := NewCalculator()
	grades := []Grade{
		{ID: "g1", GradeValue: 6.0, GradeType: GradeTypeNumeric},
		{ID: "g2", GradeValue: 8.0, GradeType: GradeTypeNumeric},
	}

	assert.Equal(t, 0.0, calc.CalculatePercentile(grades, 0.0))
	assert.Equal(t, 0.0, calc.CalculatePercentile(grades, -1.0))
	assert.Equal(t, 0.0, calc.CalculatePercentile(grades, 11.0))
	assert.InDelta(t, 50.0, calc.CalculatePercentile(grades, 7.0), 0.01)
}

func TestValidator_CreditUpperLimit(t *testing.T) {
	v := NewValidator(nil)
	assert.NoError(t, v.ValidateGradeValue(10, string(GradeTypeCredit)))
	assert.NoError(t, v.ValidateGradeValue(25, string(GradeTypeCredit)))
	assert.Error(t, v.ValidateGradeValue(26, string(GradeTypeCredit)))
	assert.Error(t, v.ValidateGradeValue(0, string(GradeTypeCredit)))
}

func TestGrades_GetMyAverages_GravementeInsufficiente(t *testing.T) {
	mockRepo := new(MockRepository)
	calc := NewCalculator()
	svc := &service{repo: mockRepo, calculator: calc}

	mockRepo.On("FindByStudent", "s-low").Return([]Grade{
		{
			ID:            "g1",
			StudentID:     "s-low",
			SubjectID:     "sub1",
			GradeValue:    3.5,
			GradeCategory: GradeCategorySummative,
			Semester:      1,
			IsPublished:   true,
		},
	}, nil).Once()

	res, err := svc.GetMyAverages(context.Background(), "s-low")
	assert.NoError(t, err)
	assert.Equal(t, "GRAVEMENTE INSUFFICIENTE", res.Semester1.Condition)
}

func TestValidator_ValidateJudgmentString_Comprehensive(t *testing.T) {
	v := NewValidator(nil)
	validJudgments := []string{
		"Insufficiente", "Mediocre", "Sufficiente", "Discreto", "Buono",
		"Distinto", "Ottimo", "Eccellente", "Gravemente Insufficiente", "Quasi Sufficiente",
		"Avanzato", "Intermedio", "Base", "Iniziale", "Non Raggiunto",
		"avanzato", "OTTIMO", "sufficiente", "non raggiunto",
	}

	for _, j := range validJudgments {
		err := v.ValidateJudgmentString(j)
		assert.NoError(t, err, "Judgment %q should be valid", j)
	}

	invalidJudgments := []string{
		"NonValido", "Unknown", "10", "Super", "", "InvalidText",
	}
	for _, j := range invalidJudgments {
		err := v.ValidateJudgmentString(j)
		assert.Error(t, err, "Judgment %q should be invalid", j)
	}
}

func TestAnalytics_GetClassAnalysis_ExcludesAbsencesFromDistributionAndCalculates1To3(t *testing.T) {
	mockRepo := new(MockRepository)
	analytics := NewAnalyticsService(mockRepo)

	mockRepo.On("FindByClass", "c1", 1).Return([]Grade{
		{ID: "g1", StudentID: "s1", SubjectID: "math", GradeValue: -1, GradeCategory: GradeCategorySummative, IsPublished: true},
		{ID: "g2", StudentID: "s1", SubjectID: "math", GradeValue: 3.0, GradeCategory: GradeCategorySummative, IsPublished: true},
		{ID: "g3", StudentID: "s2", SubjectID: "math", GradeValue: 5.0, GradeCategory: GradeCategorySummative, IsPublished: true},
		{ID: "g4", StudentID: "s2", SubjectID: "history", GradeValue: 4.0, GradeCategory: GradeCategorySummative, IsPublished: true},
		{ID: "g5", StudentID: "s3", SubjectID: "math", GradeValue: 8.0, GradeCategory: GradeCategorySummative, IsPublished: true},
		{ID: "g6", StudentID: "s3", SubjectID: "history", GradeValue: 9.0, GradeCategory: GradeCategorySummative, IsPublished: true},
	}, nil).Once()

	res, err := analytics.GetClassAnalysis("c1", 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	// Distribution should have "1-3" and not count -1 as a low grade
	assert.Contains(t, res.GradeAnalysis.Distribution, "1-3")
	assert.Equal(t, 1, res.GradeAnalysis.Distribution["1-3"].Count) // only g2 (3.0), g1 (-1) excluded
	assert.Equal(t, 3.0, res.GradeAnalysis.Min)
	assert.Equal(t, 9.0, res.GradeAnalysis.Max)

	// Check at risk students - s2 has failed math (5.0) and history (4.0), avgFailing should be 4.5
	var s2Risk *RiskStudent
	for _, r := range res.AtRisk {
		if r.StudentID == "s2" {
			s2Risk = &r
			break
		}
	}
	assert.NotNil(t, s2Risk)
	assert.Equal(t, 4.5, s2Risk.AvgFailing)
}

func TestAnalytics_GetClassAnalysis_NoSummativeGrades_ZeroMinMax(t *testing.T) {
	mockRepo := new(MockRepository)
	analytics := NewAnalyticsService(mockRepo)

	mockRepo.On("FindByClass", "c-empty", 1).Return([]Grade{
		{ID: "g1", StudentID: "s1", SubjectID: "math", GradeValue: -1, GradeCategory: GradeCategorySummative, IsPublished: true},
	}, nil).Once()

	res, err := analytics.GetClassAnalysis("c-empty", 1)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, res.GradeAnalysis.Min)
	assert.Equal(t, 0.0, res.GradeAnalysis.Max)
	assert.Equal(t, 0, res.GradeAnalysis.Distribution["1-3"].Count)
}

func TestAnalytics_GetSubjectAnalysis_ExcludesAbsencesFromPassRate(t *testing.T) {
	mockRepo := new(MockRepository)
	analytics := NewAnalyticsService(mockRepo)

	mockRepo.On("FindBySubject", "math", 1).Return([]Grade{
		{ID: "g1", GradeValue: -1, GradeCategory: GradeCategorySummative, IsPublished: true},
		{ID: "g2", GradeValue: 6.0, GradeCategory: GradeCategorySummative, IsPublished: true},
		{ID: "g3", GradeValue: 8.0, GradeCategory: GradeCategorySummative, IsPublished: true},
	}, nil).Once()

	res, err := analytics.GetSubjectAnalysis("math", 1)
	assert.NoError(t, err)
	assert.Equal(t, 2, res.Aggregated.TotalGrades)
	assert.Equal(t, 100.0, res.Aggregated.PassRate) // 2 out of 2 passed (absence excluded)
}

func TestAnalytics_GetSchoolStatistics_YearFilter(t *testing.T) {
	mockRepo := new(MockRepository)
	analytics := NewAnalyticsService(mockRepo)

	inYearDate, _ := time.Parse("2006-01-02", "2025-11-15")
	outYearDate, _ := time.Parse("2006-01-02", "2024-05-10")

	mockRepo.On("FindWithFilter", GradeFilter{SchoolID: "sch1"}).Return([]Grade{
		{ID: "g1", StudentID: "s1", GradeValue: 8.0, GradeCategory: GradeCategorySummative, Date: inYearDate, IsPublished: true},
		{ID: "g2", StudentID: "s2", GradeValue: 4.0, GradeCategory: GradeCategorySummative, Date: outYearDate, IsPublished: true},
	}, nil).Once()

	res, err := analytics.GetSchoolStatistics("2025/2026", "sch1")
	assert.NoError(t, err)
	assert.Equal(t, 1, res.DataVolume.TotalGrades)
	assert.Equal(t, 1, res.DataVolume.TotalStudents)
	assert.Equal(t, 8.0, res.Overall.Average)
}

func TestGrades_GetStudentGradesWithFilter_StudentRole_FiltersUnpublished(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	delTime := time.Now()
	mockRepo.On("FindByStudent", "s1").Return([]Grade{
		{ID: "g1", StudentID: "s1", GradeValue: 8.0, IsPublished: true, DeletedAt: nil},
		{ID: "g2", StudentID: "s1", GradeValue: 9.0, IsPublished: false, DeletedAt: nil},     // unpublished draft
		{ID: "g3", StudentID: "s1", GradeValue: 7.0, IsPublished: true, DeletedAt: &delTime}, // deleted
	}, nil).Once()

	res, err := svc.GetStudentGradesWithFilter(context.Background(), "s1", "student", "s1", GradeFilter{})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "g1", res[0].ID)
}

func TestGrades_GetSemesterReport_PromotionFallbackOnEnrollErr(t *testing.T) {
	mockRepo := new(MockRepository)
	calc := NewCalculator()
	svc := &service{repo: mockRepo, calculator: calc}

	mockRepo.On("FindByStudent", "s1").Return([]Grade{
		{ID: "g1", StudentID: "s1", SubjectID: "math", GradeValue: 8.0, Semester: 1, IsPublished: true, GradeCategory: GradeCategorySummative},
		{ID: "g2", StudentID: "s1", SubjectID: "history", GradeValue: 7.0, Semester: 1, IsPublished: true, GradeCategory: GradeCategorySummative},
	}, nil).Once()

	mockRepo.On("GetStudentClassAndSchoolInfo", mock.Anything, "s1").Return("Student 1", "Class 1A", "c1", "sch1", nil).Once()
	mockRepo.On("FindEnrolledSubjects", "s1", 1).Return(nil, assert.AnError).Once()
	mockRepo.On("GetTeacherNamesByClass", mock.Anything, "c1").Return(map[string]string{}, nil).Once()
	mockRepo.On("GetSubjectNamesMap", mock.Anything, "sch1").Return(map[string]string{}, nil).Once()
	mockRepo.On("GetScrutinyRecordSummary", mock.Anything, "s1", 1).Return(0.0, 0.0, false, nil).Once()

	rep, err := svc.GetSemesterReport(context.Background(), "admin-1", "admin", "s1", 1)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	assert.Equal(t, "SÌ", rep.Promoted)
}

func TestGrades_GetMyTrend_ExcludesAbsencesFromTrend(t *testing.T) {
	mockRepo := new(MockRepository)
	calc := NewCalculator()
	svc := &service{repo: mockRepo, calculator: calc}

	mockRepo.On("FindByStudent", "s1").Return([]Grade{
		{ID: "g1", StudentID: "s1", SubjectID: "math", GradeValue: 8.0, Date: time.Now().AddDate(0, 0, -5), IsPublished: true},
		{ID: "g2", StudentID: "s1", SubjectID: "math", GradeValue: -1.0, Date: time.Now().AddDate(0, 0, -3), IsPublished: true}, // absence
		{ID: "g3", StudentID: "s1", SubjectID: "math", GradeValue: 9.0, Date: time.Now().AddDate(0, 0, -1), IsPublished: true},
	}, nil).Once()
	mockRepo.On("GetStudentClassAndSchoolInfo", mock.Anything, "s1").Return("Student 1", "Class 1A", "c1", "sch1", nil).Once()
	mockRepo.On("GetClassSubjectAverage", mock.Anything, "c1", "math", mock.Anything, "s1").Return(7.5, nil).Once()

	trend, err := svc.GetMyTrend(context.Background(), "s1", "student", "s1", "math")
	assert.NoError(t, err)
	assert.NotNil(t, trend)
	assert.Len(t, trend.Trends, 2) // g2 (-1) excluded
	assert.Equal(t, 8.0, trend.Trends[0].Grade)
	assert.Equal(t, 9.0, trend.Trends[1].Grade)
}
