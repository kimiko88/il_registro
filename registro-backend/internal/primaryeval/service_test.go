package primaryeval

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockRepository struct {
	objectives  map[string]*LearningObjective
	evaluations []PrimaryEvaluation
	errToReturn error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		objectives:  make(map[string]*LearningObjective),
		evaluations: []PrimaryEvaluation{},
	}
}

func (m *mockRepository) CreateObjective(ctx context.Context, obj *LearningObjective) (*LearningObjective, error) {
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	if obj.ID == "" {
		obj.ID = "obj-mock-id"
	}
	m.objectives[obj.ID] = obj
	return obj, nil
}

func (m *mockRepository) ListObjectives(ctx context.Context, schoolID, subjectID string, classID *string, yearGrade int, academicYear string) ([]LearningObjective, error) {
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	var res []LearningObjective
	for _, o := range m.objectives {
		if o.SchoolID == schoolID {
			res = append(res, *o)
		}
	}
	return res, nil
}

func (m *mockRepository) GetObjective(ctx context.Context, id string) (*LearningObjective, error) {
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	obj, exists := m.objectives[id]
	if !exists {
		return nil, errors.New("objective not found")
	}
	return obj, nil
}

func (m *mockRepository) DeleteObjective(ctx context.Context, id string) error {
	if m.errToReturn != nil {
		return m.errToReturn
	}
	delete(m.objectives, id)
	return nil
}

func (m *mockRepository) SaveEvaluationsBatch(ctx context.Context, schoolID, teacherID string, req *SaveEvaluationsBatchRequest) error {
	if m.errToReturn != nil {
		return m.errToReturn
	}
	for _, item := range req.Evaluations {
		eval := PrimaryEvaluation{
			ID:                   "eval-" + item.StudentID,
			SchoolID:             schoolID,
			StudentID:            item.StudentID,
			ClassID:              req.ClassID,
			SubjectID:            req.SubjectID,
			TeacherID:            teacherID,
			ObjectiveID:          req.ObjectiveID,
			Level:                item.Level,
			DimensionAutonomy:    item.DimensionAutonomy,
			DimensionContinuity:  item.DimensionContinuity,
			DimensionFamiliarity: item.DimensionFamiliarity,
			DimensionResources:   item.DimensionResources,
			Semester:             req.Semester,
			Notes:                item.Notes,
			Date:                 time.Now(),
		}
		m.evaluations = append(m.evaluations, eval)
	}
	return nil
}

func (m *mockRepository) ListEvaluationsByClassAndSubject(ctx context.Context, classID, subjectID string, semester int) ([]PrimaryEvaluation, error) {
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	var res []PrimaryEvaluation
	for _, e := range m.evaluations {
		if e.ClassID == classID && e.SubjectID == subjectID {
			if semester == 0 || e.Semester == semester {
				res = append(res, e)
			}
		}
	}
	return res, nil
}

func (m *mockRepository) GetStudentEvaluations(ctx context.Context, studentID string, semester int) ([]PrimaryEvaluation, error) {
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	var res []PrimaryEvaluation
	for _, e := range m.evaluations {
		if e.StudentID == studentID {
			if semester == 0 || e.Semester == semester {
				res = append(res, e)
			}
		}
	}
	return res, nil
}

func TestIsValidPrimaryLevel(t *testing.T) {
	valid := []string{
		LevelAvanzato,
		LevelIntermedio,
		LevelBase,
		LevelInViaPrimaAcquisiz,
	}
	for _, l := range valid {
		if !IsValidPrimaryLevel(l) {
			t.Errorf("expected level %s to be valid", l)
		}
	}

	invalid := []string{
		"ottimo",
		"insufficiente",
		"",
		"10",
		"avanzato_plus",
	}
	for _, l := range invalid {
		if IsValidPrimaryLevel(l) {
			t.Errorf("expected level %s to be invalid", l)
		}
	}
}

func TestService_CreateObjective(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()
	svc := NewService(repo)

	// Missing title
	_, err := svc.CreateObjective(ctx, "school-1", &CreateObjectiveRequest{
		YearGrade: 3,
		SubjectID: "sub-1",
	})
	if !errors.Is(err, ErrMissingTitle) {
		t.Fatalf("expected ErrMissingTitle, got %v", err)
	}

	// Invalid yearGrade < 1
	_, err = svc.CreateObjective(ctx, "school-1", &CreateObjectiveRequest{
		Title:     "Ascolto e parlato",
		YearGrade: 0,
		SubjectID: "sub-1",
	})
	if !errors.Is(err, ErrInvalidYearGrade) {
		t.Fatalf("expected ErrInvalidYearGrade, got %v", err)
	}

	// Invalid yearGrade > 5
	_, err = svc.CreateObjective(ctx, "school-1", &CreateObjectiveRequest{
		Title:     "Ascolto e parlato",
		YearGrade: 6,
		SubjectID: "sub-1",
	})
	if !errors.Is(err, ErrInvalidYearGrade) {
		t.Fatalf("expected ErrInvalidYearGrade, got %v", err)
	}

	// Valid creation with default academic year
	obj, err := svc.CreateObjective(ctx, "school-1", &CreateObjectiveRequest{
		Title:       "Comprensione del testo narrativo",
		Description: "Legge e comprende testi di varia tipologia",
		YearGrade:   4,
		SubjectID:   "sub-italiano",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.AcademicYear != "2025/2026" {
		t.Errorf("expected default academic year 2025/2026, got %s", obj.AcademicYear)
	}
	if obj.Title != "Comprensione del testo narrativo" {
		t.Errorf("expected title to match, got %s", obj.Title)
	}
}

func TestService_SaveEvaluationsBatch(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()
	svc := NewService(repo)

	// Missing mandatory headers
	err := svc.SaveEvaluationsBatch(ctx, "school-1", "teacher-1", &SaveEvaluationsBatchRequest{
		ClassID:   "",
		SubjectID: "sub-1",
	})
	if err == nil {
		t.Fatal("expected error for empty class id")
	}

	// Invalid semester
	err = svc.SaveEvaluationsBatch(ctx, "school-1", "teacher-1", &SaveEvaluationsBatchRequest{
		ClassID:     "class-1",
		SubjectID:   "sub-1",
		ObjectiveID: "obj-1",
		Semester:    3,
		Evaluations: []StudentObjectiveLevelItem{
			{StudentID: "st-1", Level: LevelAvanzato},
		},
	})
	if err == nil {
		t.Fatal("expected error for invalid semester")
	}

	// Invalid ministerial level in evaluation item
	err = svc.SaveEvaluationsBatch(ctx, "school-1", "teacher-1", &SaveEvaluationsBatchRequest{
		ClassID:     "class-1",
		SubjectID:   "sub-1",
		ObjectiveID: "obj-1",
		Semester:    1,
		Evaluations: []StudentObjectiveLevelItem{
			{StudentID: "st-1", Level: "ottimo"},
		},
	})
	if !errors.Is(err, ErrInvalidLevel) {
		t.Fatalf("expected ErrInvalidLevel, got %v", err)
	}

	// Valid batch
	err = svc.SaveEvaluationsBatch(ctx, "school-1", "teacher-1", &SaveEvaluationsBatchRequest{
		ClassID:     "class-1",
		SubjectID:   "sub-1",
		ObjectiveID: "obj-1",
		Semester:    1,
		Date:        "2026-01-15",
		Evaluations: []StudentObjectiveLevelItem{
			{StudentID: "st-1", Level: LevelAvanzato, DimensionAutonomy: "autonomo"},
			{StudentID: "st-2", Level: LevelIntermedio, DimensionContinuity: "continuo"},
			{StudentID: "st-3", Level: LevelInViaPrimaAcquisiz, Notes: "necessita supporto"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error saving batch: %v", err)
	}

	if len(repo.evaluations) != 3 {
		t.Fatalf("expected 3 evaluations saved, got %d", len(repo.evaluations))
	}
}

func TestService_GetPrimaryMatrix(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()
	svc := NewService(repo)

	// Seed objectives
	obj1 := &LearningObjective{ID: "obj-1", SchoolID: "school-1", Title: "Lettura", YearGrade: 2}
	obj2 := &LearningObjective{ID: "obj-2", SchoolID: "school-1", Title: "Scrittura", YearGrade: 2}
	repo.objectives["obj-1"] = obj1
	repo.objectives["obj-2"] = obj2

	// Seed evaluations for student 1 & student 2
	repo.evaluations = append(repo.evaluations,
		PrimaryEvaluation{
			ID:          "e1",
			StudentID:   "st-1",
			StudentName: "Rossi Mario",
			ClassID:     "cls-2a",
			SubjectID:   "sub-ita",
			ObjectiveID: "obj-1",
			Level:       LevelAvanzato,
			Semester:    1,
			Date:        time.Now(),
		},
		PrimaryEvaluation{
			ID:          "e2",
			StudentID:   "st-1",
			StudentName: "Rossi Mario",
			ClassID:     "cls-2a",
			SubjectID:   "sub-ita",
			ObjectiveID: "obj-2",
			Level:       LevelIntermedio,
			Semester:    1,
			Date:        time.Now(),
		},
		PrimaryEvaluation{
			ID:          "e3",
			StudentID:   "st-2",
			StudentName: "Bianchi Anna",
			ClassID:     "cls-2a",
			SubjectID:   "sub-ita",
			ObjectiveID: "obj-1",
			Level:       LevelBase,
			Semester:    1,
			Date:        time.Now(),
		},
	)

	matrix, err := svc.GetPrimaryMatrix(ctx, "school-1", "cls-2a", "sub-ita", 1)
	if err != nil {
		t.Fatalf("unexpected matrix error: %v", err)
	}

	if len(matrix.Students) != 2 {
		t.Fatalf("expected 2 students in matrix, got %d", len(matrix.Students))
	}
	if len(matrix.Objectives) != 2 {
		t.Fatalf("expected 2 objectives in matrix, got %d", len(matrix.Objectives))
	}

	// Verify Rossi Mario has obj-1 = avanzato and obj-2 = intermedio
	marioRow := matrix.Students[0]
	if marioRow.StudentID != "st-1" {
		t.Errorf("expected student st-1 first, got %s", marioRow.StudentID)
	}
	if marioRow.Evaluations["obj-1"].Level != LevelAvanzato {
		t.Errorf("expected LevelAvanzato for obj-1, got %s", marioRow.Evaluations["obj-1"].Level)
	}
	if marioRow.Evaluations["obj-2"].Level != LevelIntermedio {
		t.Errorf("expected LevelIntermedio for obj-2, got %s", marioRow.Evaluations["obj-2"].Level)
	}
}
