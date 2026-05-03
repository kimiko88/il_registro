package scrutiny

import (
	"context"
	"fmt"
	"math"
	"registro-backend/internal/attendance"
	"registro-backend/internal/classes"
	"registro-backend/internal/grades"
	"registro-backend/internal/users"
)

type Service struct {
	repo       Repository
	gradeRepo  grades.Repository
	classRepo  classes.Repository
	userRepo   users.Repository
	attRepo    attendance.Repository
}

func NewService(repo Repository, gr grades.Repository, cr classes.Repository, ur users.Repository, ar attendance.Repository) *Service {
	return &Service{
		repo:      repo,
		gradeRepo: gr,
		classRepo: cr,
		userRepo:  ur,
		attRepo:   ar,
	}
}

func (s *Service) GetMatrix(ctx context.Context, classID string, semester int) (*ScrutinyMatrix, error) {
	// 1. Get Class Subjects
	subjects, err := s.classRepo.GetClassSubjects(ctx, classID)
	if err != nil {
		return nil, err
	}

	// 2. Get Students in Class
	allStudents, err := s.userRepo.GetStudentsByClass(ctx, classID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("DEBUG: Found %d students for class %s\n", len(allStudents), classID)

	// 3. Get Existing Records
	records, err := s.repo.ListRecordsByClass(ctx, classID, semester)
	if err != nil {
		return nil, err
	}
	recordMap := make(map[string]ScrutinyRecord)
	for _, r := range records {
		recordMap[r.StudentID] = r
	}

	matrix := &ScrutinyMatrix{
		ClassID:  classID,
		Semester: semester,
	}

	for _, sub := range subjects {
		matrix.Subjects = append(matrix.Subjects, SubjectInfo{ID: sub.SubjectID, Name: sub.SubjectName})
	}

	for _, stu := range allStudents {
		row := StudentScrutinyRow{
			StudentID:   stu.ID,
			StudentName: stu.LastName + " " + stu.FirstName,
			SubjectData: make(map[string]SubjectAverages),
		}

		// Calculate Averages for each subject
		for _, sub := range subjects {
			gradesList, err := s.gradeRepo.FindByClassAndSubject(classID, sub.SubjectID, semester)
			if err != nil {
				continue
			}

			var sum float64
			var count int
			for _, g := range gradesList {
				if g.StudentID == stu.StudentID && g.IsPublished && g.DeletedAt == nil {
					sum += g.GradeValue
					count++
				}
			}

			avg := 0.0
			if count > 0 {
				avg = sum / float64(count)
			}

			row.SubjectData[sub.SubjectID] = SubjectAverages{
				Average:    avg,
				GradeCount: count,
				Proposed:   math.Round(avg), // Simplified rounding
			}
		}

		if _, ok := recordMap[stu.ID]; ok {
			// Load full record with grades
			full, _ := s.repo.GetRecord(ctx, stu.ID, classID, semester)
			row.Record = full
		}

		// 4. Get Attendance Stats
		stats, err := s.attRepo.GetStats(stu.StudentID) // Note: GetStats should ideally take semester dates
		if err == nil && stats != nil {
			row.AttendanceStats = AttendanceSummary{
				Absences:   stats.TotalAbsences,
				Lates:      stats.TotalLates,
				EarlyExits: stats.TotalEarlyExits,
			}
		}

		matrix.Students = append(matrix.Students, row)
	}

	return matrix, nil
}

func (s *Service) SaveScrutiny(ctx context.Context, coordinatorID string, req SaveScrutinyRequest) error {
	rec := &ScrutinyRecord{
		StudentID:     req.StudentID,
		ClassID:       req.ClassID,
		Semester:      req.Semester,
		ConductGrade:  req.ConductGrade,
		FinalDecision: req.FinalDecision,
		Notes:         req.Notes,
		CoordinatorID: coordinatorID,
	}

	for _, g := range req.Grades {
		rec.Grades = append(rec.Grades, ScrutinyGrade{
			SubjectID:  g.SubjectID,
			FinalGrade: g.FinalGrade,
			TeacherID:  coordinatorID, // Usually the coordinator finalized it
		})
	}

	return s.repo.SaveRecord(ctx, rec)
}
