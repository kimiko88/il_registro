package scrutiny

import (
	"fmt"
	"math"
	"testing"
)

func BenchmarkMapGradeToReligionJudgment(b *testing.B) {
	values := []float64{10.0, 9.2, 8.5, 7.0, 6.2, 5.0, 4.0, 0.0}
	idx := 0
	for b.Loop() {
		_ = mapGradeToReligionJudgment(values[idx%len(values)])
		idx++
	}
}

func BenchmarkSubjectAveragesRounding(b *testing.B) {
	averages := []float64{7.3, 7.5, 6.0, 5.4, 9.8, 8.25, 4.9}
	idx := 0
	for b.Loop() {
		avg := averages[idx%len(averages)]
		idx++
		_ = math.Min(10, math.Max(1, math.Round(avg)))
	}
}

func BenchmarkMatrixStudentRowCalculation(b *testing.B) {
	numStudents := 30
	numSubjects := 12

	type mockGradeEntry struct {
		sum   float64
		count int
	}

	gradeIndex := make(map[string]map[string]*mockGradeEntry)
	for s := 0; s < numStudents; s++ {
		sID := fmt.Sprintf("stu-%d", s)
		gradeIndex[sID] = make(map[string]*mockGradeEntry)
		for sub := 0; sub < numSubjects; sub++ {
			subID := fmt.Sprintf("sub-%d", sub)
			gradeIndex[sID][subID] = &mockGradeEntry{
				sum:   35.0 + float64((s+sub)%15),
				count: 5,
			}
		}
	}

	subjects := make([]SubjectInfo, numSubjects)
	for sub := 0; sub < numSubjects; sub++ {
		subjects[sub] = SubjectInfo{
			ID:         fmt.Sprintf("sub-%d", sub),
			Name:       fmt.Sprintf("Materia %d", sub),
			IsReligion: sub == 0,
		}
	}

	for b.Loop() {
		matrix := ScrutinyMatrix{
			ClassID:  "class-1",
			Semester: 1,
			Subjects: subjects,
		}

		for s := 0; s < numStudents; s++ {
			sID := fmt.Sprintf("stu-%d", s)
			row := StudentScrutinyRow{
				StudentID:   sID,
				StudentName: fmt.Sprintf("Studente %d", s),
				SubjectData: make(map[string]SubjectAverages),
			}

			smap := gradeIndex[sID]
			for _, sub := range subjects {
				e := smap[sub.ID]
				avg := e.sum / float64(e.count)
				proposed := math.Min(10, math.Max(1, math.Round(avg)))
				proposedJudgment := ""
				if sub.IsReligion {
					proposedJudgment = mapGradeToReligionJudgment(avg)
				}
				row.SubjectData[sub.ID] = SubjectAverages{
					Average:          avg,
					GradeCount:       e.count,
					Proposed:         proposed,
					ProposedJudgment: proposedJudgment,
				}
			}
			matrix.Students = append(matrix.Students, row)
		}
	}
}
