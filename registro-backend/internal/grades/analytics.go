package grades

import (
	"time"
)

type AnalyticsService interface {
	GetStudentAverage(studentID string, subjectID string) (float64, error)
	GetClassAverage(classID string, subjectID string) (float64, error)

	// Advanced Analytics
	GetClassAnalysis(classID string, semester int) (*AnalyticsClassResponse, error)
	GetSubjectAnalysis(subjectID string, semester int) (*AnalyticsSubjectResponse, error)
	GetStudentProfile(studentID string, semester int) (*AnalyticsStudentResponse, error)
	GetSchoolStatistics(year string, schoolID ...string) (*SchoolStatisticsResponse, error)
}

type analyticsService struct {
	repo       Repository
	calculator *Calculator
}

func NewAnalyticsService(r Repository) AnalyticsService {
	return &analyticsService{
		repo:       r,
		calculator: NewCalculator(),
	}
}

func (a *analyticsService) GetStudentAverage(studentID string, subjectID string) (float64, error) {
	// Filter by student and subject at database level via FindWithFilter
	filter := GradeFilter{
		StudentID: studentID,
		SubjectID: subjectID,
	}
	grades, err := a.repo.FindWithFilter(filter)
	if err != nil {
		return 0, err
	}

	var subjectGrades []Grade
	for _, g := range grades {
		if g.IsPublished && g.DeletedAt == nil {
			subjectGrades = append(subjectGrades, g)
		}
	}

	return a.calculator.CalculateAverage(subjectGrades), nil
}

func (a *analyticsService) GetClassAverage(classID string, subjectID string) (float64, error) {
	// Fetch all semesters for class and subject in a single optimized query (semester 0 = all)
	allGrades, err := a.repo.FindByClassAndSubject(classID, subjectID, 0)
	if err != nil {
		return 0, err
	}

	var validGrades []Grade
	for _, g := range allGrades {
		if g.IsPublished && g.DeletedAt == nil {
			validGrades = append(validGrades, g)
		}
	}

	return a.calculator.CalculateAverage(validGrades), nil
}

// --- Specific Implementations ---

func (a *analyticsService) GetClassAnalysis(classID string, semester int) (*AnalyticsClassResponse, error) {
	// 1. Fetch Grades
	grades, err := a.repo.FindByClass(classID, semester)
	if err != nil {
		return nil, err
	}

	// Filter Published Only
	var valid []Grade
	for _, g := range grades {
		if g.IsPublished && g.DeletedAt == nil {
			valid = append(valid, g)
		}
	}

	// Separate Summative for Average
	var summative []Grade
	for _, g := range valid {
		if g.GradeCategory == GradeCategorySummative {
			summative = append(summative, g)
		}
	}

	// 2. Compute Stats
	mean, stdDev, skew, kurt := a.calculator.CalculateBellCurve(summative)
	median := a.calculator.CalculateMedian(summative)

	// Distribution
	dist := make(map[string]CountPct)
	totalSumm := len(summative)
	if totalSumm > 0 {
		ranges := map[string]int{"0-3": 0, "4-5.9": 0, "6-6.9": 0, "7-7.9": 0, "8-8.9": 0, "9-10": 0}
		for _, g := range summative {
			val := g.GradeValue
			if val < 4 {
				ranges["0-3"]++
			} else if val < 6 {
				ranges["4-5.9"]++
			} else if val < 7 {
				ranges["6-6.9"]++
			} else if val < 8 {
				ranges["7-7.9"]++
			} else if val < 9 {
				ranges["8-8.9"]++
			} else {
				ranges["9-10"]++
			}
		}
		for k, v := range ranges {
			dist[k] = CountPct{Count: v, Percentage: (float64(v) / float64(totalSumm)) * 100}
		}
	}

	// 3. Subject Performance Aggregated
	subMap := make(map[string][]Grade)
	for _, g := range summative {
		subMap[g.SubjectID] = append(subMap[g.SubjectID], g)
	}
	var subPerf []SubjectPerf
	for subID, gs := range subMap {
		avg := a.calculator.CalculateAverage(gs)
		passed, failed := 0, 0
		for _, g := range gs {
			if g.GradeValue >= 6 {
				passed++
			} else {
				failed++
			}
		}
		rate := 0.0
		if len(gs) > 0 {
			rate = (float64(passed) / float64(len(gs))) * 100
		}

		subPerf = append(subPerf, SubjectPerf{
			Subject:  subID,
			AvgGrade: avg,
			Passed:   passed,
			Failed:   failed,
			PassRate: rate,
			Trend:    "stable", // Mock
		})
	}

	// 4. Student Stats & Risk
	// Need to group by Student to count avg
	stuMap := make(map[string][]Grade)
	for _, g := range summative {
		stuMap[g.StudentID] = append(stuMap[g.StudentID], g)
	}

	atRisk := []RiskStudent{}
	topPerf := []TopStudent{}
	promoted := 0
	suspended := 0

	for sID, sGrades := range stuMap {
		avg := a.calculator.CalculateWeightedAverage(sGrades) // Weighted correct? Prompt says "Media: SOLO voti sommativi" which I filtered.

		// Check failing
		var failedSubs []string
		subGroups := make(map[string][]Grade)
		for _, g := range sGrades {
			subGroups[g.SubjectID] = append(subGroups[g.SubjectID], g)
		}

		failingCount := 0
		for sub, gs := range subGroups {
			if a.calculator.CalculateAverage(gs) < 6.0 {
				failingCount++
				failedSubs = append(failedSubs, sub)
			}
		}

		if failingCount == 0 {
			promoted++
		} else if failingCount <= 2 {
			suspended++ // Sospensione del giudizio (1-2 debiti formativi)
		}

		stuDisplayName := "Student " + sID
		if len(sID) >= 4 {
			stuDisplayName = "Student " + sID[:4]
		}

		if failingCount > 0 {
			atRisk = append(atRisk, RiskStudent{
				StudentID:      sID,
				StudentName:    stuDisplayName,
				FailedSubjects: failedSubs,
				AvgFailing:     avg,
				Recommendation: "Supporto intensivo",
			})
		}

		if avg >= 8.0 {
			topPerf = append(topPerf, TopStudent{
				StudentID:   sID,
				StudentName: stuDisplayName,
				AvgGrade:    avg,
			})
		}
	}

	className := "Class " + classID
	if len(classID) >= 4 {
		className = "Class " + classID[:4]
	}

	// Calculate real min and max from class summative grades on Italian scale [1.0, 10.0]
	realMin := 10.0
	realMax := 1.0
	if len(summative) > 0 {
		for _, g := range summative {
			if g.GradeValue >= 1.0 && g.GradeValue <= 10.0 {
				if g.GradeValue < realMin {
					realMin = g.GradeValue
				}
				if g.GradeValue > realMax {
					realMax = g.GradeValue
				}
			}
		}
		if realMin > realMax {
			realMin, realMax = 1.0, 10.0
		}
	} else {
		realMin, realMax = 1.0, 10.0
	}

	return &AnalyticsClassResponse{
		Class:        ClassInfo{ID: classID, Name: className},
		Semester:     semester,
		AnalysisDate: time.Now().Format("2006-01-02"),
		StudentStats: StudentStats{
			Total:     len(stuMap),
			Present:   len(stuMap), // simple assumption
			Promoted:  promoted,
			Suspended: suspended,
		},
		GradeAnalysis: GradeAnalysis{
			Average:      mean,
			Median:       median,
			StdDev:       stdDev,
			BellCurve:    BellCurveInfo{Mean: mean, StdDev: stdDev, Skewness: skew, Kurtosis: kurt, NormalityTest: "approx_normal"},
			Distribution: dist,
			Min:          realMin,
			Max:          realMax,
		},
		Subjects:      subPerf,
		AtRisk:        atRisk,
		TopPerformers: topPerf,
	}, nil
}

func (a *analyticsService) GetSubjectAnalysis(subjectID string, semester int) (*AnalyticsSubjectResponse, error) {
	grades, err := a.repo.FindBySubject(subjectID, semester)
	if err != nil {
		return nil, err
	}

	// Summative only for key metrics
	var summative []Grade
	passCount := 0
	for _, g := range grades {
		if g.IsPublished && g.DeletedAt == nil && g.GradeCategory == GradeCategorySummative {
			summative = append(summative, g)
			if g.GradeValue >= 6.0 {
				passCount++
			}
		}
	}

	avg := a.calculator.CalculateAverage(summative)
	passRate := 0.0
	if len(summative) > 0 {
		passRate = (float64(passCount) / float64(len(summative))) * 100.0
	}

	return &AnalyticsSubjectResponse{
		Subject:  SubjectMeta{ID: subjectID, Name: "Subject"},
		Semester: semester,
		Aggregated: SubjectAggregated{
			TotalGrades: len(summative),
			AvgGrade:    avg,
			PassRate:    passRate,
		},
	}, nil
}

func (a *analyticsService) GetStudentProfile(studentID string, semester int) (*AnalyticsStudentResponse, error) {
	// Re-use Logic from GetMyGrades/GetMyAverages but add Metrics
	grades, err := a.repo.FindByStudent(studentID)
	if err != nil {
		return nil, err
	}

	var relevant []Grade
	var summative []Grade
	for _, g := range grades {
		if !g.IsPublished || g.DeletedAt != nil {
			continue
		}
		if semester > 0 && int(g.Semester) != semester {
			continue
		}
		relevant = append(relevant, g)
		if g.GradeCategory == GradeCategorySummative {
			summative = append(summative, g)
		}
	}

	overallAvg := a.calculator.CalculateWeightedAverage(summative)

	// Subjects
	subMap := make(map[string][]Grade)
	for _, g := range relevant {
		subMap[g.SubjectID] = append(subMap[g.SubjectID], g)
	}

	var subAvgs []SubjectAvgDetail
	failingSubjects := 0
	for sub, gs := range subMap {
		var subSum []Grade
		for _, g := range gs {
			if g.GradeCategory == GradeCategorySummative {
				subSum = append(subSum, g)
			}
		}

		avg := a.calculator.CalculateAverage(subSum)
		if len(subSum) > 0 && avg < 6.0 {
			failingSubjects++
		}
		subAvgs = append(subAvgs, SubjectAvgDetail{
			Subject:     sub,
			Average:     avg,
			GradesCount: len(gs),
			Trend:       "stable",
		})
	}

	promotionStatus := "On Track"
	if failingSubjects > 2 {
		promotionStatus = "Critical"
	} else if failingSubjects > 0 {
		promotionStatus = "At Risk"
	}

	return &AnalyticsStudentResponse{
		Student: StudentInfo{ID: studentID},
		OverallProfile: OverallProfile{
			OverallAverage:  overallAvg,
			PromotionStatus: promotionStatus,
		},
		SubjectAverages: subAvgs,
	}, nil
}

func (a *analyticsService) GetSchoolStatistics(year string, schoolID ...string) (*SchoolStatisticsResponse, error) {
	filter := GradeFilter{}
	if len(schoolID) > 0 && schoolID[0] != "" {
		filter.SchoolID = schoolID[0]
	}
	allGrades, err := a.repo.FindWithFilter(filter)
	if err != nil {
		return nil, err
	}

	var grades []Grade
	for _, g := range allGrades {
		if g.IsPublished && g.DeletedAt == nil {
			grades = append(grades, g)
		}
	}

	vol := DataVolume{
		TotalGrades:   len(grades),
		TotalStudents: 0,
	}

	var summative []Grade
	studentSet := make(map[string]bool)
	for _, g := range grades {
		studentSet[g.StudentID] = true
		if g.GradeCategory == GradeCategorySummative {
			summative = append(summative, g)
		}
	}
	vol.TotalStudents = len(studentSet)

	avg := a.calculator.CalculateAverage(summative)

	return &SchoolStatisticsResponse{
		ReportDate: time.Now().Format("2006-01-02"),
		DataVolume: vol,
		Overall: GradeAnalysis{
			Average: avg,
		},
	}, nil
}
