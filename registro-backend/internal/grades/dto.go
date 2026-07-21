package grades

import "time"

// --- Requests ---

type CreateGradeRequest struct {
	StudentID     string  `json:"student_id" binding:"required"`
	SubjectID     string  `json:"subject_id" binding:"required"`
	GradeValue    float64 `json:"grade_value" binding:"required,min=-1,max=10"`
	GradeType     string  `json:"grade_type" binding:"required"`
	Semester      int     `json:"semester" binding:"required,min=1,max=2"` // 1 or 2
	Description   string  `json:"description"`
	RubricID      *string `json:"rubric_id"`
	Weight        float64 `json:"weight"`
	IsPublished   bool    `json:"is_published"`
	GradeCategory  string  `json:"grade_category"`          // formative/summative/practical
	EvaluationType *string `json:"evaluation_type"`         // Written/Oral/Practical
	Date           string  `json:"date" binding:"required"` // ISO date string preferred for input
}

type UpdateGradeRequest struct {
	GradeValue     *float64 `json:"grade_value"`
	GradeType      *string  `json:"grade_type"`
	Description    *string  `json:"description"`
	Date           *string  `json:"date"` // YYYY-MM-DD
	Weight         *float64 `json:"weight"`
	IsPublished    *bool    `json:"is_published"`
	GradeCategory  *string  `json:"grade_category"`
	EvaluationType *string  `json:"evaluation_type"`
	Reason         string   `json:"reason" binding:"required"` // Reason is mandatory for updates
}

type BulkImportRequest struct {
	Semester int `form:"semester" binding:"required"`
}

// --- Filters ---

type GradeFilter struct {
	Semester    int    `form:"semester"`
	SubjectID   string `form:"subject_id"`
	GradeType   string `form:"grade_type"`
	IsPublished *bool  `form:"published"` // Pointer to distinguish false from missing
	SortBy      string `form:"sort"`
}

// --- Responses ---

type GradeResponse struct {
	ID            string    `json:"id"`
	StudentID     string    `json:"student_id"`
	SubjectID     string    `json:"subject_id"`
	TeacherID     string    `json:"teacher_id"`
	GradeValue    float64   `json:"grade_value"`
	GradeType     string    `json:"grade_type"`
	Semester      int       `json:"semester"`
	Description   string    `json:"description"`
	Date           time.Time `json:"date"`
	GradeCategory  string    `json:"grade_category"`
	EvaluationType *string   `json:"evaluation_type,omitempty"`
	IsPublished    bool      `json:"is_published"`
	TestID         *string   `json:"test_id,omitempty"`
}

// --- Student/Parent Specific Responses ---

type MyGradesResponse struct {
	Student          StudentInfo             `json:"student"`
	RegistrationDate string                  `json:"registration_date"`
	Semesters        []SemesterGradesSummary `json:"semesters"`
}

type SemesterGradesSummary struct {
	Semester  int             `json:"semester"`
	StartDate string          `json:"start_date"`
	EndDate   string          `json:"end_date"`
	Grades    []GradeResponse `json:"grades"`
}

type StudentAveragesResponse struct {
	Semester1 SemesterAverageSummary `json:"semester1"`
	Semester2 SemesterAverageSummary `json:"semester2"`
}

type SemesterAverageSummary struct {
	Subjects       []SubjectAverage `json:"subjects"`
	OverallAverage float64          `json:"overall_average"`
	Condition      string           `json:"condition"` // "OK" | "MONITOR" | "ALERT"
}

type SubjectAverage struct {
	Subject         string  `json:"subject"`
	Average         float64 `json:"average"`
	WeightedAverage float64 `json:"weighted_average"`
	TotalGrades     int     `json:"total_grades"`
	LastUpdateDate  string  `json:"last_update_date"`
}

type TrendResponse struct {
	Subject string       `json:"subject"`
	Trends  []TrendPoint `json:"trends"`
	Summary TrendSummary `json:"summary"`
}

type TrendPoint struct {
	Date         string  `json:"date"`
	Grade        float64 `json:"grade"`
	MovingAvg3   float64 `json:"moving_avg_3"`
	ClassAverage float64 `json:"class_average"`
	Position     string  `json:"position"` // "above_average", "at_average", "below_average"
}

type TrendSummary struct {
	TrendDirection string `json:"trend_direction"` // "improving", "stable", "declining"
	OverallTrend   string `json:"overall_trend"`
	Recommendation string `json:"recommendation"`
}

type SemesterReportResponse struct {
	Semester       int             `json:"semester"`
	Period         PeriodDate      `json:"period"`
	Subjects       []SubjectReport `json:"subjects"`
	OverallAverage float64         `json:"overall_average"`
	Promoted       string          `json:"promoted"`
	Status         string          `json:"status"`
	LastUpdate     time.Time       `json:"last_update"`
}

type PeriodDate struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// --- Analytics DTOs ---

type AnalyticsClassResponse struct {
	Class         ClassInfo     `json:"class"`
	Semester      int           `json:"semester"`
	AnalysisDate  string        `json:"analysisDate"`
	StudentStats  StudentStats  `json:"studentStats"`
	GradeAnalysis GradeAnalysis `json:"gradeAnalysis"`
	Subjects      []SubjectPerf `json:"subjectsPerformance"`
	Credits       CreditsInfo   `json:"credits,omitempty"`
	AtRisk        []RiskStudent `json:"atRisk"`
	TopPerformers []TopStudent  `json:"topPerformers"`
}

type ClassInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Year    string `json:"year"`
	Section string `json:"section"`
}

type StudentStats struct {
	Total     int `json:"total"`
	Present   int `json:"present"`
	Promoted  int `json:"promoted"`
	Suspended int `json:"suspended"`
	Dropped   int `json:"dropped"`
}

type GradeAnalysis struct {
	Average      float64             `json:"averageClass"`
	Median       float64             `json:"median"`
	StdDev       float64             `json:"stdDeviation"`
	Min          float64             `json:"min"`
	Max          float64             `json:"max"`
	Distribution map[string]CountPct `json:"distribution"`
	BellCurve    BellCurveInfo       `json:"bellCurve"`
}

type CountPct struct {
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type BellCurveInfo struct {
	Mean          float64 `json:"mean"`
	StdDev        float64 `json:"standardDeviation"`
	Skewness      float64 `json:"skewness"`
	Kurtosis      float64 `json:"kurtosis"`
	NormalityTest string  `json:"normalityTest"`
}

type SubjectPerf struct {
	Subject  string        `json:"subject"`
	Teacher  string        `json:"teacher"`
	AvgGrade float64       `json:"avgGrade"`
	Passed   int           `json:"passed"`
	Failed   int           `json:"failed"`
	PassRate float64       `json:"passRate"`
	Outliers []OutlierInfo `json:"outliers"`
	Trend    string        `json:"trend"`
}

type OutlierInfo struct {
	StudentID   string  `json:"studentId"`
	StudentName string  `json:"studentName"`
	Grade       float64 `json:"grade"`
}

type CreditsInfo struct {
	AvgCredits     float64 `json:"avgCredits"`
	TotalAssigned  int     `json:"totalCreditsAssigned"`
	PctoCredits    int     `json:"pctoCredits"`
	ProjectCredits int     `json:"projectCredits"`
	Behavior       string  `json:"behavior"`
}

type RiskStudent struct {
	StudentID      string   `json:"studentId"`
	StudentName    string   `json:"studentName"`
	FailedSubjects []string `json:"failedSubjects"`
	AvgFailing     float64  `json:"avgFailing"`
	Recommendation string   `json:"recommendation"`
	ActionPlan     string   `json:"actionPlan"`
}

type TopStudent struct {
	StudentID         string   `json:"studentId"`
	StudentName       string   `json:"studentName"`
	AvgGrade          float64  `json:"avgGrade"`
	ExcellentSubjects []string `json:"excellentSubjects"`
}

type AnalyticsSubjectResponse struct {
	Subject         SubjectMeta       `json:"subject"`
	Semester        int               `json:"semester"`
	ByClass         []ClassPerf       `json:"performanceByClass"`
	Aggregated      SubjectAggregated `json:"aggregated"`
	Trend           TrendAnalysis     `json:"trendAnalysis"`
	TopPerformers   []TopStudent      `json:"topPerformers"`
	UnderPerformers []RiskStudent     `json:"underPerformers"` // Reusing RiskStudent struct structure roughly
	Recommendations []string          `json:"schoolRecommendations"`
}

type SubjectMeta struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Hours int    `json:"hours"`
	Type  string `json:"type"`
}

type ClassPerf struct {
	ClassID   string  `json:"classId"`
	ClassName string  `json:"className"`
	Count     int     `json:"studentCount"`
	AvgGrade  float64 `json:"avgGrade"`
	PassRate  float64 `json:"passRate"`
	// simplistic top/weak
	Trend string `json:"trend"`
}

type SubjectAggregated struct {
	TotalStudents int                `json:"totalStudents"`
	TotalGrades   int                `json:"totalGrades"`
	AvgGrade      float64            `json:"avgGrade"`
	PassRate      float64            `json:"passRate"`
	Difficulty    string             `json:"difficulty"`
	Categories    map[string]CatStat `json:"gradesByCategory"`
}

type CatStat struct {
	Count    int     `json:"count"`
	AvgGrade float64 `json:"avgGrade"`
}

type TrendAnalysis struct {
	Semester1     SemStat `json:"semester1"`
	Semester2     SemStat `json:"semester2"`
	YearTrend     string  `json:"yearTrend"`
	ChangePercent string  `json:"changePercent"` // "+1.4%"
}

type SemStat struct {
	Avg      float64 `json:"avg"`
	PassRate float64 `json:"passRate"`
}

type AnalyticsStudentResponse struct {
	Student         StudentInfo        `json:"student"`
	OverallProfile  OverallProfile     `json:"overallProfile"`
	SubjectAverages []SubjectAvgDetail `json:"subjectAverages"`
	Strengths       []StrengthWeakness `json:"strengths"`
	Weaknesses      []StrengthWeakness `json:"weaknesses"`
	TrendAnalysis   StudentTrend       `json:"trendAnalysis"`
	Credits         StudentCredits     `json:"credits"`
	Recommendations Recommendations    `json:"recommendations"`
	PeerComparison  PeerComp           `json:"peerComparison"`
}

type OverallProfile struct {
	OverallAverage  float64 `json:"overallAverage"`
	AverageTrend    string  `json:"averageTrend"`
	PromotionStatus string  `json:"promotionStatus"`
	ClassRank       string  `json:"classRank"`
	PercentileRank  float64 `json:"percentileRank"`
}

type SubjectAvgDetail struct {
	Subject      string  `json:"subject"`
	Average      float64 `json:"average"`
	ClassAverage float64 `json:"classAverage"`
	Comparison   string  `json:"comparison"`
	GradesCount  int     `json:"grades"`
	Trend        string  `json:"trend"`
	FormativeAvg float64 `json:"formativeAvg"`
	SummativeAvg float64 `json:"summativeAvg"`
}

type StrengthWeakness struct {
	Subject string  `json:"subject"`
	Reason  string  `json:"reason"`
	Diff    float64 `json:"diff,omitempty"` // above/below class
}

type StudentTrend struct {
	OverallTrend  string    `json:"overallTrend"`
	LastQuarter   []float64 `json:"lastQuarter"`
	MovingAvg     float64   `json:"movingAverage"`
	Volatility    string    `json:"volatility"`
	Semester1Avg  float64   `json:"semester1Avg"`
	Semester2Avg  float64   `json:"semester2Avg"`
	SemesterTrend string    `json:"semesterTrend"`
}

type StudentCredits struct {
	TotalCredits   int    `json:"totalCredits"`
	PctoHours      int    `json:"pctoHours"`
	PctoCredits    int    `json:"pctoCredits"`
	ProjectCredits int    `json:"projectCredits"`
	CreditTrend    string `json:"creditTrend"`
}

type Recommendations struct {
	Academic  []string `json:"academic"`
	Behavior  string   `json:"behavior"`
	NextSteps string   `json:"nextSteps"`
}

type PeerComp struct {
	PercentileRank   float64 `json:"percentileRank"`
	TopStudentAvg    float64 `json:"topStudentAvg"`
	BottomStudentAvg float64 `json:"bottomStudentAvg"`
	StudentPosition  string  `json:"studentPosition"`
}

type SchoolStatisticsResponse struct {
	School          SchoolMeta         `json:"school"`
	ReportDate      string             `json:"reportDate"`
	DataVolume      DataVolume         `json:"dataVolume"`
	Overall         GradeAnalysis      `json:"overallStatistics"`
	ByYear          []YearStat         `json:"byYear"`
	BySemester      []SemStatGeneral   `json:"bySemester"`
	ByCategory      map[string]CatStat `json:"byCategory"`
	CriticalAreas   []AreaStat         `json:"criticalAreas"`
	ExcellenceAreas []AreaStat         `json:"excellenceAreas"`
	ScrutinyDates   []ScrutinyDate     `json:"scrutinyDates"`
	Recommendations []string           `json:"strategicRecommendations"`
}

type SchoolMeta struct {
	Name      string `json:"name"`
	City      string `json:"city"`
	YearStart string `json:"yearStart"`
	YearEnd   string `json:"yearEnd"`
}

type DataVolume struct {
	TotalGrades   int     `json:"totalGrades"`
	TotalStudents int     `json:"totalStudents"`
	TotalTeachers int     `json:"totalTeachers"`
	TotalClasses  int     `json:"totalClasses"`
	TotalSubjects int     `json:"totalSubjects"`
	GradingRate   float64 `json:"gradingRate"`
}

type YearStat struct {
	Year          string  `json:"year"`
	AvgGrade      float64 `json:"avgGrade"`
	StudentCount  int     `json:"studentCount"`
	PromotionRate float64 `json:"promotionRate"`
	FailureRate   float64 `json:"failureRate"`
	AvgCredits    float64 `json:"avgCredits,omitempty"`
}

type SemStatGeneral struct {
	Semester    int     `json:"semester"`
	AvgGrade    float64 `json:"avgGrade"`
	PassRate    float64 `json:"passRate"`
	TotalGrades int     `json:"totalGrades"`
}

type AreaStat struct {
	Area           string  `json:"area"`
	FailureRate    float64 `json:"failureRate"` // or PassRate for excellence?
	PassRate       float64 `json:"passRate,omitempty"`
	AvgGrade       float64 `json:"avgGrade"`
	Recommendation string  `json:"recommendation"`
}

type ScrutinyDate struct {
	Semester int    `json:"semester"`
	Date     string `json:"date"`
	Status   string `json:"status"`
}

type SubjectReport struct {
	Subject        string     `json:"subject"`
	Teacher        string     `json:"teacher"`
	Grades         []GradeVal `json:"grades"`
	SubjectAverage float64    `json:"subject_average"`
	Passed         bool       `json:"passed"`
}

type GradeVal struct {
	Value    float64   `json:"value"`
	Date     time.Time `json:"date"`
	Category string    `json:"category"`
}

// --- Teacher Dashboard Responses (Existing) ---

type StudentGradesSummary struct {
	Student StudentInfo     `json:"student"`
	Grades  []GradeResponse `json:"grades"`
	Summary GradeSummary    `json:"summary"`
}

type StudentInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	ClassID string `json:"class_id"`
}

type GradeSummary struct {
	TotalGrades  int     `json:"total_grades"`
	Semester1Avg float64 `json:"semester1_avg"`
	Semester2Avg float64 `json:"semester2_avg"`
}

type ClassGradesResponse struct {
	ClassID  string                `json:"class_id"`
	Students []StudentGradeSummary `json:"students"`
}

type StudentGradeSummary struct {
	StudentID    string          `json:"student_id"`
	FullName     string          `json:"full_name"`
	AvgSemester1 float64         `json:"avg_semester1"`
	AvgSemester2 float64         `json:"avg_semester2"`
	Grades       []GradeResponse `json:"grades"`
}

type SubjectStatsResponse struct {
	Subject SubjectInfo `json:"subject"`
	Classes []ClassStat `json:"classes"`
}

type SubjectInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsVotable bool   `json:"is_votable"`
}

type ClassStat struct {
	ClassID           string         `json:"class_id"`
	ClassName         string         `json:"class_name"`
	TotalStudents     int            `json:"total_students"`
	AvgGrade          float64        `json:"avg_grade"`
	GradeDistribution map[string]int `json:"grade_distribution"`
}

type ImportResult struct {
	Imported int           `json:"imported"`
	Skipped  int           `json:"skipped"`
	Errors   []ImportError `json:"errors"`
}

type ImportError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

type CreateClassTestRequest struct {
	ClassID        string              `json:"class_id" binding:"required"`
	SubjectID      string              `json:"subject_id" binding:"required"`
	Title          string              `json:"title" binding:"required"`
	Date           string              `json:"date" binding:"required"`
	TeacherNotes   string              `json:"teacher_notes"`
	ParentNotes    string              `json:"parent_notes"`
	EvaluationType string              `json:"evaluation_type" binding:"required"` // Written/Oral/Practical
	Semester       int                 `json:"semester"`
	Grades         []StudentGradeInput `json:"grades"`
}

type StudentGradeInput struct {
	StudentID  string   `json:"student_id" binding:"required"`
	GradeValue *float64 `json:"grade_value" binding:"required"`
	Notes      string   `json:"notes"`
}

type ClassTestResponse struct {
	ID             string `json:"id"`
	ClassID        string `json:"class_id"`
	SubjectID      string `json:"subject_id"`
	TeacherID      string `json:"teacher_id"`
	Title          string `json:"title"`
	Date           string `json:"date"`
	TeacherNotes   string `json:"teacher_notes,omitempty"`
	ParentNotes    string `json:"parent_notes,omitempty"`
	EvaluationType string `json:"evaluation_type"`
}

type UpdateClassTestRequest struct {
	Title          string                    `json:"title" binding:"required"`
	Date           string                    `json:"date" binding:"required"`
	TeacherNotes   string                    `json:"teacher_notes"`
	ParentNotes    string                    `json:"parent_notes"`
	EvaluationType string                    `json:"evaluation_type" binding:"required"`
	Semester       int                       `json:"semester"`
	Grades         []UpdateStudentGradeInput `json:"grades"`
}

type UpdateStudentGradeInput struct {
	StudentID  string   `json:"student_id" binding:"required"`
	GradeValue *float64 `json:"grade_value"`
	Notes      string   `json:"notes"`
}
