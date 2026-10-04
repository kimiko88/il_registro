package enrollment

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

type FormationSolver struct{}

func NewFormationSolver() *FormationSolver {
	return &FormationSolver{}
}

func (s *FormationSolver) Solve(pool []EnrollmentApplication, params FormationParams) (*ClassFormationDraftResult, error) {
	if len(pool) == 0 {
		return nil, errors.New("empty student pool")
	}

	targetClasses := params.TargetClassCount
	if targetClasses <= 0 {
		targetClasses = int(math.Ceil(float64(len(pool)) / 25.0))
		if targetClasses <= 0 {
			targetClasses = 1
		}
	}

	classNames := make([]string, targetClasses)
	for i := 0; i < targetClasses; i++ {
		classNames[i] = fmt.Sprintf("1%c", 'A'+i)
	}

	classes := make([]ClassGroupAssignment, targetClasses)
	for i := 0; i < targetClasses; i++ {
		classes[i] = ClassGroupAssignment{
			ClassName: classNames[i],
			Students:  []EnrollmentApplication{},
		}
	}

	// 1. Separate students by special traits
	var l104Students []EnrollmentApplication
	var standardStudents []EnrollmentApplication

	for _, app := range pool {
		if app.HasDisabilityL104 {
			l104Students = append(l104Students, app)
		} else {
			standardStudents = append(standardStudents, app)
		}
	}

	// Distribute L.104 students evenly across target classes (max 1 per class ideally)
	for i, l104App := range l104Students {
		targetIdx := i % targetClasses
		classes[targetIdx].Students = append(classes[targetIdx].Students, l104App)
		classes[targetIdx].L104Count++
	}

	// 2. Identify paired friends (mutual requested classmates)
	placedMap := make(map[string]bool)
	taxCodeToApp := make(map[string]EnrollmentApplication)
	for _, app := range standardStudents {
		taxCodeToApp[app.StudentTaxCode] = app
	}

	for _, app := range standardStudents {
		if placedMap[app.StudentTaxCode] {
			continue
		}

		var group []EnrollmentApplication
		group = append(group, app)
		placedMap[app.StudentTaxCode] = true

		// Find mutual friends
		for _, reqCF := range app.RequestedClassmates {
			if friend, ok := taxCodeToApp[reqCF]; ok && !placedMap[reqCF] {
				// Check if friend also requested app (mutual)
				isMutual := false
				for _, backReq := range friend.RequestedClassmates {
					if backReq == app.StudentTaxCode {
						isMutual = true
						break
					}
				}
				if isMutual {
					group = append(group, friend)
					placedMap[reqCF] = true
				}
			}
		}

		// Find best class to place this group (least populated class that respects DPR 81/2009)
		bestClassIdx := 0
		minSize := 999999
		for cIdx, cls := range classes {
			currentSize := len(cls.Students)
			// If class has L.104, limit capacity to 20 per DPR 81/2009
			maxAllowed := 27
			if cls.L104Count > 0 {
				maxAllowed = 20
			}

			if currentSize+len(group) <= maxAllowed && currentSize < minSize {
				minSize = currentSize
				bestClassIdx = cIdx
			}
		}

		// If all classes at max allowed, fallback to smallest
		if minSize == 999999 {
			for cIdx, cls := range classes {
				if len(cls.Students) < minSize {
					minSize = len(cls.Students)
					bestClassIdx = cIdx
				}
			}
		}

		classes[bestClassIdx].Students = append(classes[bestClassIdx].Students, group...)
	}

	// Sort and distribute remaining standard students balancing gender and grade
	var remaining []EnrollmentApplication
	for _, app := range standardStudents {
		if !placedMap[app.StudentTaxCode] {
			remaining = append(remaining, app)
		}
	}

	// Sort remaining by grade descending to interleave
	sort.Slice(remaining, func(i, j int) bool {
		return remaining[i].MiddleSchoolGrade > remaining[j].MiddleSchoolGrade
	})

	for _, app := range remaining {
		bestClassIdx := 0
		minScore := math.MaxFloat64

		for cIdx, cls := range classes {
			curTotal := len(cls.Students)
			maxAllowed := 27
			if cls.L104Count > 0 {
				maxAllowed = 20
			}

			if curTotal >= maxAllowed {
				continue
			}

			// Score function: penalty for size, penalty for gender imbalance
			males := 0
			females := 0
			for _, s := range cls.Students {
				if s.Gender == "M" {
					males++
				} else {
					females++
				}
			}

			score := float64(curTotal) * 10.0
			if app.Gender == "M" {
				score += float64(males) * 5.0
			} else {
				score += float64(females) * 5.0
			}

			if score < minScore {
				minScore = score
				bestClassIdx = cIdx
			}
		}

		classes[bestClassIdx].Students = append(classes[bestClassIdx].Students, app)
	}

	// Calculate metrics for each class
	for i := range classes {
		cls := &classes[i]
		cls.TotalStudents = len(cls.Students)
		totalGrades := 0
		cls.MalesCount = 0
		cls.FemalesCount = 0
		cls.L104Count = 0
		cls.DSACount = 0

		for _, s := range cls.Students {
			if s.Gender == "M" {
				cls.MalesCount++
			} else {
				cls.FemalesCount++
			}
			if s.HasDisabilityL104 {
				cls.L104Count++
			}
			if s.HasDSA {
				cls.DSACount++
			}
			totalGrades += s.MiddleSchoolGrade
		}

		if cls.TotalStudents > 0 {
			cls.AverageGrade = math.Round((float64(totalGrades)/float64(cls.TotalStudents))*100) / 100
		}
	}

	return &ClassFormationDraftResult{
		Classes: classes,
	}, nil
}
