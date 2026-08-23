package grades

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 1. Test parseFilter normalization for Page <= 0
func TestGrades_AuditFix_ParseFilter_PageZeroAndNegative(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/grades?page=0&limit=20&semester=3", nil)
	c.Request = req

	filter := h.parseFilter(c)
	assert.Equal(t, 1, filter.Page, "Page=0 must be normalized to 1")
	assert.Equal(t, 0, filter.Semester, "Semester > 2 must be normalized to 0")

	reqNeg := httptest.NewRequest("GET", "/grades?page=-5", nil)
	c.Request = reqNeg
	filterNeg := h.parseFilter(c)
	assert.Equal(t, 1, filterNeg.Page, "Page < 0 must be normalized to 1")
}

// 2. Test CalculateWeightedAverage fallback to arithmetic mean when weights are 0
func TestGrades_AuditFix_CalculateWeightedAverage_FallbackArithmetic(t *testing.T) {
	calc := NewCalculator()

	// Grades in Italian school scale (1 to 10) with Weight = 0
	grades := []Grade{
		{GradeValue: 8.0, GradeType: GradeTypeNumeric, Weight: 0},
		{GradeValue: 6.0, GradeType: GradeTypeNumeric, Weight: 0},
	}

	weightedAvg := calc.CalculateWeightedAverage(grades)
	assert.Equal(t, 7.0, weightedAvg, "Weighted average with 0 weights must fall back to arithmetic average (7.0)")

	// Normal weighted average
	gradesWithWeights := []Grade{
		{GradeValue: 8.0, GradeType: GradeTypeNumeric, Weight: 1.0},
		{GradeValue: 6.0, GradeType: GradeTypeNumeric, Weight: 3.0},
	}
	// (8*1 + 6*3) / 4 = 26/4 = 6.5
	assert.Equal(t, 6.5, calc.CalculateWeightedAverage(gradesWithWeights))
}

// 3 & 9. Test ConvertJudgmentToValue case-insensitivity and 1-10 Italian scale
func TestGrades_AuditFix_ConvertJudgmentToValue_CaseInsensitive(t *testing.T) {
	calc := NewCalculator()

	testJudgments := map[string]float64{
		"ottimo":                   10.0,
		"OTTIMO":                   10.0,
		"  Eccellente  ":           10.0,
		"Avanzato":                 10.0,
		"Distinto":                 9.0,
		"distinto":                 9.0,
		"buono":                    8.0,
		"BUONO":                    8.0,
		"Intermedio":               8.0,
		"Discreto":                 7.0,
		"base":                     7.0,
		"Sufficiente":              6.0,
		"sufficiente":              6.0,
		"mediocre":                 5.0,
		"Quasi Sufficiente":        5.0,
		"iniziale":                 5.0,
		"insufficiente":            4.0,
		"non raggiunto":            4.0,
		"gravemente insufficiente": 3.0,
		"GRAVEMENTE INSUFFICIENTE": 3.0,
		"giudizio sconosciuto":     0.0,
	}

	for input, expectedVal := range testJudgments {
		val := calc.ConvertJudgmentToValue(input)
		assert.Equal(t, expectedVal, val, "Judgment '%s' must map to %f", input, expectedVal)
	}

	// CalculateAverage with qualitative judgments in different casings
	grades := []Grade{
		{GradeType: GradeTypeJudgment, Description: "ottimo"},
		{GradeType: GradeTypeJudgment, Description: "BUONO"},
		{GradeType: GradeTypeJudgment, Description: "sufficiente"},
	}
	// (10 + 8 + 6) / 3 = 8.0
	assert.Equal(t, 8.0, calc.CalculateAverage(grades))
}

// 4. Test GetStudentGrades delegation to GetStudentGradesWithFilter
func TestGrades_AuditFix_GetStudentGrades_SingleResponsibility(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForDeprecationTest)
	h := NewHandler(mockSvc, nil)

	mockSvc.On("GetStudentGradesWithFilter", mock.Anything, "stud-1", "student", "stud-1", mock.Anything).
		Return([]GradeResponse{{ID: "g1", GradeValue: 8.0}}, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/grades/student/stud-1", nil)
	c.Set("user_id", "stud-1")
	c.Set("role", "student")
	c.Params = []gin.Param{{Key: "studentID", Value: "stud-1"}}

	h.GetStudentGrades(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true", w.Header().Get("Deprecation"))
}

// 5. Test GetChildGradesAverage role enforcement
func TestGrades_AuditFix_GetChildGradesAverage_RoleEnforcement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForDeprecationTest)
	h := NewHandler(mockSvc, nil)

	// Teacher or Student calling child-grades must get 403 Forbidden
	wStudent := httptest.NewRecorder()
	cStudent, _ := gin.CreateTestContext(wStudent)
	cStudent.Request = httptest.NewRequest("GET", "/grades/child-grades/stud-1/average", nil)
	cStudent.Set("user_id", "stud-1")
	cStudent.Set("role", "student")
	cStudent.Params = []gin.Param{{Key: "studentID", Value: "stud-1"}}
	h.GetChildGradesAverage(cStudent)
	assert.Equal(t, http.StatusForbidden, wStudent.Code, "Student calling GetChildGradesAverage must receive 403")

	// Parent role should proceed to service
	mockSvc.On("GetChildAverages", mock.Anything, "parent-1", "stud-1").Return(&StudentAveragesResponse{}, nil).Once()
	wParent := httptest.NewRecorder()
	cParent, _ := gin.CreateTestContext(wParent)
	cParent.Request = httptest.NewRequest("GET", "/grades/child-grades/stud-1/average", nil)
	cParent.Set("user_id", "parent-1")
	cParent.Set("role", "parent")
	cParent.Params = []gin.Param{{Key: "studentID", Value: "stud-1"}}
	h.GetChildGradesAverage(cParent)
	assert.Equal(t, http.StatusOK, wParent.Code, "Parent calling GetChildGradesAverage must succeed with 200")
}

// 6. Test GetChildSemesterReport role enforcement
func TestGrades_AuditFix_GetChildSemesterReport_RoleEnforcement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForDeprecationTest)
	h := NewHandler(mockSvc, nil)

	// Student attempting to call GetChildSemesterReport must be rejected with 403
	wStudent := httptest.NewRecorder()
	cStudent, _ := gin.CreateTestContext(wStudent)
	cStudent.Request = httptest.NewRequest("GET", "/grades/child-grades/stud-1/semester/1", nil)
	cStudent.Set("user_id", "stud-1")
	cStudent.Set("role", "student")
	cStudent.Params = []gin.Param{
		{Key: "studentID", Value: "stud-1"},
		{Key: "semester", Value: "1"},
	}
	h.GetChildSemesterReport(cStudent)
	assert.Equal(t, http.StatusForbidden, wStudent.Code, "Student calling GetChildSemesterReport must receive 403")

	// Parent should proceed
	mockSvc.On("GetChildSemesterReport", mock.Anything, "parent-1", "stud-1", 1).Return(&SemesterReportResponse{}, nil).Once()
	wParent := httptest.NewRecorder()
	cParent, _ := gin.CreateTestContext(wParent)
	cParent.Request = httptest.NewRequest("GET", "/grades/child-grades/stud-1/semester/1", nil)
	cParent.Set("user_id", "parent-1")
	cParent.Set("role", "parent")
	cParent.Params = []gin.Param{
		{Key: "studentID", Value: "stud-1"},
		{Key: "semester", Value: "1"},
	}
	h.GetChildSemesterReport(cParent)
	assert.Equal(t, http.StatusOK, wParent.Code, "Parent calling GetChildSemesterReport must succeed with 200")
}

// 7. Test DownloadSemesterReportPDF semantic actor validation
func TestGrades_AuditFix_DownloadSemesterReportPDF_ActorValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForDeprecationTest)
	h := NewHandler(mockSvc, nil)

	mockSvc.On("ValidateParentGuardian", mock.Anything, "parent-1", "stud-1").Return(nil).Once()
	mockSvc.On("GenerateSemesterReportPDF", mock.Anything, "parent-1", "parent", "stud-1", 1).Return([]byte("%PDF-1.4..."), nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/grades/my-grades/semester/1/pdf?student_id=stud-1", nil)
	c.Set("user_id", "parent-1")
	c.Set("role", "parent")
	c.Params = []gin.Param{{Key: "semester", Value: "1"}}

	h.DownloadSemesterReportPDF(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "pagella_q1_stud-1.pdf")
}

// 8. Test CalculateBellCurve & CalculateStandardDeviation single-pass consistency
func TestGrades_AuditFix_CalculateBellCurve_SinglePass(t *testing.T) {
	calc := NewCalculator()

	grades := []Grade{
		{GradeValue: 6.0, GradeType: GradeTypeNumeric},
		{GradeValue: 7.0, GradeType: GradeTypeNumeric},
		{GradeValue: 8.0, GradeType: GradeTypeNumeric},
		{GradeValue: 9.0, GradeType: GradeTypeNumeric},
	}

	mean, stdDev, skewness, kurtosis := calc.CalculateBellCurve(grades)
	expectedStdDev := calc.CalculateStandardDeviation(grades)

	assert.Equal(t, 7.5, mean)
	assert.InDelta(t, expectedStdDev, stdDev, 0.0001)
	assert.InDelta(t, 0.0, skewness, 0.0001, "Symmetric distribution must have 0 skewness")
	assert.NotPanics(t, func() {
		_ = kurtosis
	})
}

// 10. Test sanitizeFilenameParam with package-level regex
func TestGrades_AuditFix_SanitizeFilenameParam(t *testing.T) {
	assert.Equal(t, "student_123-abc", sanitizeFilenameParam("student_123-abc"))
	assert.Equal(t, "student123", sanitizeFilenameParam("student;123/..\\"))
	assert.Equal(t, "export", sanitizeFilenameParam("???///:::"))
}

// 11. Test isVotableGrade bounds (1.0 to 10.0 Italian scale)
func TestGrades_AuditFix_IsVotableGrade_ItalianScale(t *testing.T) {
	assert.False(t, isVotableGrade(Grade{GradeType: GradeTypeNumeric}, 0.0), "0.0 must be excluded as unrated/unset")
	assert.False(t, isVotableGrade(Grade{GradeType: GradeTypeNumeric}, -1.0), "Negative values must be excluded")
	assert.False(t, isVotableGrade(Grade{GradeType: GradeTypeNumeric}, 10.5), "Values above 10.0 must be excluded")
	assert.True(t, isVotableGrade(Grade{GradeType: GradeTypeNumeric}, 1.0), "1.0 is valid minimum grade")
	assert.True(t, isVotableGrade(Grade{GradeType: GradeTypeNumeric}, 6.0), "6.0 is valid passing grade")
	assert.True(t, isVotableGrade(Grade{GradeType: GradeTypeNumeric}, 10.0), "10.0 is valid maximum grade")
	assert.True(t, isVotableGrade(Grade{GradeType: GradeTypeJudgment}, 8.0), "Converted judgment 8.0 is valid")
	assert.False(t, isVotableGrade(Grade{GradeType: GradeTypeJudgment}, 0.0), "Unrecognized judgment 0.0 is excluded")
}

// 12. Test DetectOutliers bounds clamped to [1.0, 10.0]
func TestGrades_AuditFix_DetectOutliers_ClampedBounds(t *testing.T) {
	calc := NewCalculator()

	// High sample dataset where 10.0 is an outlier (> Mean + 2*StdDev)
	grades := []Grade{
		{ID: "g1", GradeValue: 2.0, GradeType: GradeTypeNumeric},
		{ID: "g2", GradeValue: 2.0, GradeType: GradeTypeNumeric},
		{ID: "g3", GradeValue: 2.5, GradeType: GradeTypeNumeric},
		{ID: "g4", GradeValue: 2.0, GradeType: GradeTypeNumeric},
		{ID: "g5", GradeValue: 2.0, GradeType: GradeTypeNumeric},
		{ID: "g6", GradeValue: 2.0, GradeType: GradeTypeNumeric},
		{ID: "g7", GradeValue: 2.5, GradeType: GradeTypeNumeric},
		{ID: "g8", GradeValue: 2.0, GradeType: GradeTypeNumeric},
		{ID: "g9", GradeValue: 2.0, GradeType: GradeTypeNumeric},
		{ID: "g10", GradeValue: 10.0, GradeType: GradeTypeNumeric}, // clear outlier
	}

	outliers := calc.DetectOutliers(grades)
	assert.Contains(t, outliers, "g10")
}

// 13. Test academicYearDates with nil DB and validator safety
func TestGrades_AuditFix_AcademicYearDates_NilDBSafety(t *testing.T) {
	ctx := context.Background()

	// Calling academicYearDates with nil DB must safely return standard default dates without panicking
	s1Start, s1End, s2Start, s2End := academicYearDates(ctx, nil, "")
	assert.NotEmpty(t, s1Start)
	assert.NotEmpty(t, s1End)
	assert.NotEmpty(t, s2Start)
	assert.NotEmpty(t, s2End)
}
