package grades

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
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
	assert.False(t, isVotableGrade(0.0), "0.0 must be excluded as unrated/unset")
	assert.False(t, isVotableGrade(-1.0), "Negative values must be excluded")
	assert.False(t, isVotableGrade(10.5), "Values above 10.0 must be excluded")
	assert.True(t, isVotableGrade(1.0), "1.0 is valid minimum grade")
	assert.True(t, isVotableGrade(6.0), "6.0 is valid passing grade")
	assert.True(t, isVotableGrade(10.0), "10.0 is valid maximum grade")
	assert.True(t, isVotableGrade(8.0), "Converted judgment 8.0 is valid")
	assert.False(t, isVotableGrade(0.0), "Unrecognized judgment 0.0 is excluded")
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

// 14. Test ValidateGradeValue Italian 1-10 scale
func TestGrades_AuditFix_ValidateGradeValue_ItalianScale(t *testing.T) {
	v := NewValidator(nil)

	// In Italian system: 0.0 is NOT a valid grade (must be between 1 and 10, or -1 for absence)
	assert.Error(t, v.ValidateGradeValue(0.0, "numeric"), "0.0 must be rejected")
	assert.Error(t, v.ValidateGradeValue(0.5, "numeric"), "0.5 must be rejected")
	assert.Error(t, v.ValidateGradeValue(10.5, "numeric"), "10.5 must be rejected")
	assert.Error(t, v.ValidateGradeValue(-2.0, "numeric"), "-2.0 must be rejected")

	assert.NoError(t, v.ValidateGradeValue(1.0, "numeric"), "1.0 must be accepted")
	assert.NoError(t, v.ValidateGradeValue(6.0, "numeric"), "6.0 must be accepted")
	assert.NoError(t, v.ValidateGradeValue(10.0, "numeric"), "10.0 must be accepted")
	assert.NoError(t, v.ValidateGradeValue(-1.0, "numeric"), "-1.0 must be accepted as absence")
}

// 15. Test helpers ConvertNumericToJudgment and ConvertJudgmentToNumeric
func TestGrades_AuditFix_Helpers_ConvertNumericAndJudgment(t *testing.T) {
	// Out of bounds / non-evaluated
	assert.Equal(t, "Non valutato", ConvertNumericToJudgment(0.0))
	assert.Equal(t, "Non valutato", ConvertNumericToJudgment(10.5))
	assert.Equal(t, "Gravemente Insufficiente", ConvertNumericToJudgment(3.0))
	assert.Equal(t, "Insufficiente", ConvertNumericToJudgment(5.0))
	assert.Equal(t, "Sufficiente", ConvertNumericToJudgment(6.0))
	assert.Equal(t, "Discreto", ConvertNumericToJudgment(7.5))
	assert.Equal(t, "Buono", ConvertNumericToJudgment(8.5))
	assert.Equal(t, "Distinto", ConvertNumericToJudgment(9.5))
	assert.Equal(t, "Ottimo", ConvertNumericToJudgment(10.0))

	// Judgment string to numeric
	assert.Equal(t, 3.0, ConvertJudgmentToNumeric("Gravemente Insufficiente"))
	assert.Equal(t, 5.0, ConvertJudgmentToNumeric("Quasi Sufficiente"))
	assert.Equal(t, 6.0, ConvertJudgmentToNumeric("Sufficiente"))
	assert.Equal(t, 10.0, ConvertJudgmentToNumeric("Ottimo"))
	assert.Equal(t, 0.0, ConvertJudgmentToNumeric("InvalidJudgment"))
}

// 16. Test ParseCSVGrades validation and date fallback
func TestGrades_AuditFix_ParseCSVGrades_ValidationAndDateSafety(t *testing.T) {
	csvData := `StudentID,SubjectID,Value,Date,Category,Description
s1,sub1,8.5,2026-03-15,summative,Compito
s2,sub1,0.0,2026-03-15,summative,InvalidZeroGrade
s3,sub1,15.0,2026-03-15,summative,OutOfRangeGrade
s4,sub1,not_a_number,2026-03-15,summative,MalformedValue
s5,sub1,7.0,,summative,EmptyDateShouldFallbackToNow`

	reqs, err := ParseCSVGrades(strings.NewReader(csvData), 2)
	assert.NoError(t, err)
	assert.Len(t, reqs, 2, "Only valid 1-10 grades should be parsed (s1 and s5)")

	assert.Equal(t, "s1", reqs[0].StudentID)
	assert.Equal(t, 8.5, reqs[0].GradeValue)
	assert.Equal(t, 2026, reqs[0].Date.Year())

	assert.Equal(t, "s5", reqs[1].StudentID)
	assert.Equal(t, 7.0, reqs[1].GradeValue)
	assert.False(t, reqs[1].Date.IsZero(), "Date must fallback to time.Now() and not be zero")
}

// 17. Test calcSemesterAverages ignores Semester == 0
func TestGrades_AuditFix_CalcSemesterAverages_SemesterZeroExcluded(t *testing.T) {
	grades := []GradeResponse{
		{GradeValue: 8.0, Semester: 1, GradeCategory: string(GradeCategorySummative), Weight: 1.0},
		{GradeValue: 9.0, Semester: 2, GradeCategory: string(GradeCategorySummative), Weight: 1.0},
		{GradeValue: 4.0, Semester: 0, GradeCategory: string(GradeCategorySummative), Weight: 1.0}, // Unassigned semester
	}

	avg1, avg2 := calcSemesterAverages(grades)
	assert.Equal(t, 8.0, avg1, "Semester 1 average must only include Semester == 1")
	assert.Equal(t, 9.0, avg2, "Semester 2 average must only include Semester == 2, not Semester == 0")
}

// 18. Test GetMyAverages returns NV for empty semester
func TestGrades_AuditFix_GetMyAverages_EmptySemester_NV(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{
		repo:       mockRepo,
		calculator: NewCalculator(),
	}

	// Student with grades ONLY in semester 1
	grades := []Grade{
		{
			ID:            "g1",
			StudentID:     "s1",
			SubjectID:     "sub1",
			GradeValue:    8.0,
			Weight:        1.0,
			Semester:      1,
			IsPublished:   true,
			GradeCategory: GradeCategorySummative,
		},
	}

	mockRepo.On("FindByStudent", "s1").Return(grades, nil).Once()

	res, err := svc.GetMyAverages(context.Background(), "s1")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, "DISTINTO", res.Semester1.Condition)
	assert.Equal(t, "N.V.", res.Semester2.Condition, "Empty semester must return N.V.")
}
