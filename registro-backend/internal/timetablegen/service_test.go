package timetablegen

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockTimetableRepo struct {
	preferences          []TeacherPreference
	reqs                 []SubjectRoomRequirement
	constraints          []TimetableConstraint
	jobs                 map[string]*TimetableJob
	published            []GeneratedSlot
	desiderataWindowOpen bool
}

func newMockTimetableRepo() *mockTimetableRepo {
	return &mockTimetableRepo{
		jobs: make(map[string]*TimetableJob),
	}
}

func (m *mockTimetableRepo) SavePreferencesBatch(ctx context.Context, schoolID, teacherID string, academicYearID *string, prefs []PreferenceEntry) error {
	for _, p := range prefs {
		m.preferences = append(m.preferences, TeacherPreference{
			ID:             uuid.New().String(),
			SchoolID:       schoolID,
			TeacherID:      teacherID,
			AcademicYearID: academicYearID,
			DayOfWeek:      p.DayOfWeek,
			HourIndex:      p.HourIndex,
			PreferenceType: p.PreferenceType,
			Reason:         p.Reason,
		})
	}
	return nil
}

func (m *mockTimetableRepo) GetTeacherPreferences(ctx context.Context, schoolID, teacherID string, academicYearID *string) ([]TeacherPreference, error) {
	var res []TeacherPreference
	for _, p := range m.preferences {
		if p.SchoolID == schoolID && p.TeacherID == teacherID {
			res = append(res, p)
		}
	}
	return res, nil
}

func (m *mockTimetableRepo) LoadAllPreferences(ctx context.Context, schoolID string, academicYearID *string) ([]TeacherPreference, error) {
	return m.preferences, nil
}

func (m *mockTimetableRepo) ListRoomRequirements(ctx context.Context, schoolID string) ([]SubjectRoomRequirement, error) {
	return m.reqs, nil
}

func (m *mockTimetableRepo) SaveRoomRequirement(ctx context.Context, schoolID string, req SaveRoomRequirementRequest) (*SubjectRoomRequirement, error) {
	sr := SubjectRoomRequirement{
		ID:               uuid.New().String(),
		SchoolID:         schoolID,
		SubjectID:        req.SubjectID,
		RequiredRoomType: req.RequiredRoomType,
		LabHours:         req.LabHours,
		IsMandatory:      req.IsMandatory,
		CreatedAt:        time.Now(),
	}
	m.reqs = append(m.reqs, sr)
	return &sr, nil
}

func (m *mockTimetableRepo) DeleteRoomRequirement(ctx context.Context, id string) error {
	return nil
}

func (m *mockTimetableRepo) ListConstraints(ctx context.Context, schoolID string) ([]TimetableConstraint, error) {
	return m.constraints, nil
}

func (m *mockTimetableRepo) SaveConstraint(ctx context.Context, schoolID string, req SaveConstraintRequest) (*TimetableConstraint, error) {
	c := TimetableConstraint{
		ID:             uuid.New().String(),
		SchoolID:       schoolID,
		ConstraintType: req.ConstraintType,
		TargetType:     req.TargetType,
		TargetID:       req.TargetID,
		Parameters:     req.Parameters,
		IsHard:         req.IsHard,
		Priority:       req.Priority,
		IsActive:       true,
		CreatedAt:      time.Now(),
	}
	m.constraints = append(m.constraints, c)
	return &c, nil
}

func (m *mockTimetableRepo) DeleteConstraint(ctx context.Context, id string) error {
	return nil
}

func (m *mockTimetableRepo) GetDesiderataWindow(ctx context.Context, schoolID string) (bool, error) {
	return m.desiderataWindowOpen, nil
}

func (m *mockTimetableRepo) SetDesiderataWindow(ctx context.Context, schoolID string, isOpen bool) error {
	m.desiderataWindowOpen = isOpen
	return nil
}

func (m *mockTimetableRepo) CreateJob(ctx context.Context, job *TimetableJob) (*TimetableJob, error) {
	if job.ID == "" {
		job.ID = uuid.New().String()
	}
	m.jobs[job.ID] = job
	return job, nil
}

func (m *mockTimetableRepo) UpdateJob(ctx context.Context, job *TimetableJob) error {
	if existing, ok := m.jobs[job.ID]; ok {
		if job.Status != "" {
			existing.Status = job.Status
		}
		if len(job.ResultSummary) > 0 {
			existing.ResultSummary = job.ResultSummary
		}
		if job.ErrorMessage != nil {
			existing.ErrorMessage = job.ErrorMessage
		}
		if job.CompletedAt != nil {
			existing.CompletedAt = job.CompletedAt
		}
	}
	return nil
}

func (m *mockTimetableRepo) GetJob(ctx context.Context, id string) (*TimetableJob, error) {
	if j, ok := m.jobs[id]; ok {
		return j, nil
	}
	return nil, ErrJobNotFound
}

func (m *mockTimetableRepo) LoadAssignments(ctx context.Context, schoolID string, academicYearID *string) ([]AssignmentData, error) {
	return []AssignmentData{
		{
			ClassID:      "c-1",
			ClassName:    "1A",
			SubjectID:    "s-1",
			SubjectName:  "Matematica",
			TeacherID:    "t-1",
			TeacherName:  "Prof Uno",
			HoursPerWeek: 2,
		},
	}, nil
}

func (m *mockTimetableRepo) LoadRooms(ctx context.Context, schoolID string) ([]RoomData, error) {
	return nil, nil
}

func (m *mockTimetableRepo) LoadRoomRequirements(ctx context.Context, schoolID string) (map[string]SubjectRoomRequirement, error) {
	return make(map[string]SubjectRoomRequirement), nil
}

func (m *mockTimetableRepo) LoadAssociatedGroups(ctx context.Context, schoolID string) ([]AssociatedGroup, error) {
	return nil, nil
}

func (m *mockTimetableRepo) PublishGeneratedSchedule(ctx context.Context, schoolID string, slots []GeneratedSlot) error {
	m.published = slots
	return nil
}

func (m *mockTimetableRepo) ListAcademicYears(ctx context.Context, schoolID string) ([]string, error) {
	return []string{"2024/2025", "2023/2024"}, nil
}

func (m *mockTimetableRepo) ListClassesCurriculumPlans(ctx context.Context, schoolID, academicYear string) ([]ClassCurriculumPlan, error) {
	return []ClassCurriculumPlan{
		{
			ClassID:        "c-1",
			ClassName:      "1A",
			AcademicYear:   "2024/2025",
			MinHoursPerDay: 4,
			MaxHoursPerDay: 6,
			TotalHoursWeek: 5,
			Subjects: []ClassSubjectPlanItem{
				{SubjectID: "s-1", SubjectName: "Matematica", HoursPerWeek: 5},
			},
		},
	}, nil
}

func (m *mockTimetableRepo) GetClassCurriculumPlan(ctx context.Context, schoolID, classID string) (*ClassCurriculumPlan, error) {
	return &ClassCurriculumPlan{
		ClassID:        classID,
		ClassName:      "1A",
		AcademicYear:   "2024/2025",
		MinHoursPerDay: 4,
		MaxHoursPerDay: 6,
		TotalHoursWeek: 5,
		Subjects: []ClassSubjectPlanItem{
			{SubjectID: "s-1", SubjectName: "Matematica", HoursPerWeek: 5},
		},
	}, nil
}

func (m *mockTimetableRepo) SaveClassCurriculumPlan(ctx context.Context, schoolID, classID string, req SaveClassCurriculumPlanRequest) error {
	return nil
}

func (m *mockTimetableRepo) InheritClassCurriculumPlan(ctx context.Context, schoolID, targetClassID, sourceAcademicYear string) (*ClassCurriculumPlan, error) {
	return &ClassCurriculumPlan{
		ClassID:        targetClassID,
		ClassName:      "1A",
		AcademicYear:   "2024/2025",
		MinHoursPerDay: 4,
		MaxHoursPerDay: 6,
		TotalHoursWeek: 6,
		Subjects: []ClassSubjectPlanItem{
			{SubjectID: "s-1", SubjectName: "Matematica", HoursPerWeek: 4},
			{SubjectID: "s-2", SubjectName: "Italiano", HoursPerWeek: 2},
		},
	}, nil
}

func (m *mockTimetableRepo) InheritAllClassesCurriculumPlans(ctx context.Context, schoolID, sourceAcademicYear string) (*InheritAllResult, error) {
	return &InheritAllResult{
		ClassesUpdated: 1,
		SubjectsCopied: 2,
		Message:        "Ereditato con successo",
	}, nil
}

func (m *mockTimetableRepo) GetTeachersQuickPreferences(ctx context.Context, schoolID string, academicYearID *string) ([]TeacherQuickPreferenceItem, map[int]int, error) {
	return []TeacherQuickPreferenceItem{
		{
			TeacherID:    "t-1",
			TeacherName:  "Prof Rossi",
			SubjectName:  "Matematica",
			DayOff:       1,
			TimeSlotPref: "early_hours",
		},
		{
			TeacherID:    "t-2",
			TeacherName:  "Prof Bianchi",
			SubjectName:  "Italiano",
			DayOff:       3,
			TimeSlotPref: "late_hours",
		},
	}, map[int]int{1: 1, 3: 1}, nil
}

func (m *mockTimetableRepo) SaveTeacherQuickPreferences(ctx context.Context, schoolID string, academicYearID *string, items []TeacherQuickPreferenceItem) error {
	return nil
}

func TestTimetableService(t *testing.T) {
	repo := newMockTimetableRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()

	schoolID := "school-1"
	teacherID := "teacher-1"

	// 1. Save and Get Preferences
	err := svc.SaveTeacherPreferences(ctx, schoolID, teacherID, SavePreferencesRequest{
		Preferences: []PreferenceEntry{
			{DayOfWeek: 1, HourIndex: 1, PreferenceType: PrefPreferred},
			{DayOfWeek: 2, HourIndex: 6, PreferenceType: PrefUnavailable},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error saving preferences: %v", err)
	}

	prefs, err := svc.GetTeacherPreferences(ctx, schoolID, teacherID, nil)
	if err != nil {
		t.Fatalf("unexpected error getting preferences: %v", err)
	}
	if len(prefs) != 2 {
		t.Errorf("expected 2 preferences, got %d", len(prefs))
	}

	// 2. Start generation job
	jobID, err := svc.StartGeneration(ctx, schoolID, "user-admin", GenerateTimetableRequest{
		TimeLimitSeconds: 5,
	})
	if err != nil {
		t.Fatalf("unexpected error starting generation: %v", err)
	}
	if jobID == "" {
		t.Fatalf("expected non-empty jobID")
	}

	// Wait briefly for background goroutine to execute
	time.Sleep(200 * time.Millisecond)

	job, err := svc.GetJobStatus(ctx, schoolID, jobID)
	if err != nil {
		t.Fatalf("unexpected error getting job status: %v", err)
	}
	if job.Status != JobStatusCompleted {
		t.Fatalf("expected job to complete, got %s (err: %v)", job.Status, job.ErrorMessage)
	}

	// 3. Publish schedule
	err = svc.PublishSchedule(ctx, schoolID, "user-admin", jobID)
	if err != nil {
		t.Fatalf("unexpected error publishing schedule: %v", err)
	}

	if len(repo.published) != 2 {
		t.Errorf("expected 2 published slots, got %d", len(repo.published))
	}

	// Verify ResultSummary JSON was stored properly
	var res TimetableGenerationResult
	if err := json.Unmarshal(job.ResultSummary, &res); err != nil {
		t.Fatalf("failed to unmarshal job summary: %v", err)
	}
	if res.AssignedSlots != 2 {
		t.Errorf("expected 2 assigned slots in summary, got %d", res.AssignedSlots)
	}

	// 4. Test Desiderata Window Toggle
	isOpen, err := svc.IsDesiderataWindowOpen(ctx, schoolID)
	if err != nil {
		t.Fatalf("unexpected error checking window: %v", err)
	}
	if isOpen {
		t.Errorf("expected window to be closed initially")
	}

	if err := svc.SetDesiderataWindowOpen(ctx, schoolID, true); err != nil {
		t.Fatalf("unexpected error setting window open: %v", err)
	}

	isOpen, err = svc.IsDesiderataWindowOpen(ctx, schoolID)
	if err != nil || !isOpen {
		t.Fatalf("expected window to be open now, got %v (err: %v)", isOpen, err)
	}

	// 5. Test AdjustJobSlots
	// A: Adjust with non-overlapping valid slots
	adjustedSlots := []GeneratedSlot{
		{
			ClassID:     "class-1",
			ClassName:   "1A",
			SubjectID:   "subj-1",
			SubjectName: "Matematica",
			TeacherID:   &teacherID,
			TeacherName: "Prof. Rossi",
			DayOfWeek:   1,
			HourIndex:   1,
		},
		{
			ClassID:     "class-1",
			ClassName:   "1A",
			SubjectID:   "subj-1",
			SubjectName: "Matematica",
			TeacherID:   &teacherID,
			TeacherName: "Prof. Rossi",
			DayOfWeek:   1,
			HourIndex:   2,
		},
	}
	adjRes, err := svc.AdjustJobSlots(ctx, schoolID, "user-admin", jobID, AdjustTimetableRequest{Slots: adjustedSlots})
	if err != nil {
		t.Fatalf("unexpected error adjusting slots: %v", err)
	}
	if len(adjRes.HardConflicts) != 0 {
		t.Errorf("expected 0 hard conflicts, got %d", len(adjRes.HardConflicts))
	}
	if adjRes.AssignedSlots != 2 {
		t.Errorf("expected 2 assigned slots, got %d", adjRes.AssignedSlots)
	}

	// B: Adjust with overlapping slots (collision on same class & time)
	collidingSlots := []GeneratedSlot{
		{
			ClassID:     "class-1",
			ClassName:   "1A",
			SubjectID:   "subj-1",
			SubjectName: "Matematica",
			DayOfWeek:   2,
			HourIndex:   1,
		},
		{
			ClassID:     "class-1",
			ClassName:   "1A",
			SubjectID:   "subj-2",
			SubjectName: "Fisica",
			DayOfWeek:   2,
			HourIndex:   1,
		},
	}
	conflictRes, err := svc.AdjustJobSlots(ctx, schoolID, "user-admin", jobID, AdjustTimetableRequest{Slots: collidingSlots})
	if err != nil {
		t.Fatalf("unexpected error adjusting colliding slots: %v", err)
	}
	if len(conflictRes.HardConflicts) == 0 {
		t.Errorf("expected collision to be flagged in HardConflicts")
	}
}

func TestCurriculumPlansAndDailyLimits(t *testing.T) {
	repo := newMockTimetableRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()
	schoolID := "school-1"

	// 1. List Academic Years
	years, err := svc.ListAcademicYears(ctx, schoolID)
	if err != nil {
		t.Fatalf("failed to list academic years: %v", err)
	}
	if len(years) < 2 {
		t.Errorf("expected at least 2 academic years, got %d", len(years))
	}

	// 2. List Classes Plans
	plans, err := svc.ListClassesCurriculumPlans(ctx, schoolID, "2024/2025")
	if err != nil {
		t.Fatalf("failed to list classes curriculum plans: %v", err)
	}
	if len(plans) != 1 || plans[0].ClassName != "1A" {
		t.Errorf("unexpected classes plans: %+v", plans)
	}

	// 3. Save Curriculum Plan
	saveReq := SaveClassCurriculumPlanRequest{
		MinHoursPerDay: 4,
		MaxHoursPerDay: 6,
		Subjects: []ClassSubjectPlanItem{
			{SubjectID: "s-1", SubjectName: "Matematica", HoursPerWeek: 5},
			{SubjectID: "s-2", SubjectName: "Italiano", HoursPerWeek: 4},
		},
	}
	if err := svc.SaveClassCurriculumPlan(ctx, schoolID, "c-1", saveReq); err != nil {
		t.Fatalf("failed to save class curriculum plan: %v", err)
	}

	// 4. Inherit Class Plan
	inherited, err := svc.InheritClassCurriculumPlan(ctx, schoolID, "c-1", "2023/2024")
	if err != nil {
		t.Fatalf("failed to inherit class curriculum plan: %v", err)
	}
	if len(inherited.Subjects) != 2 {
		t.Errorf("expected 2 inherited subjects, got %d", len(inherited.Subjects))
	}

	// 5. Inherit All Classes
	allRes, err := svc.InheritAllClassesCurriculumPlans(ctx, schoolID, "2023/2024")
	if err != nil {
		t.Fatalf("failed to inherit all classes: %v", err)
	}
	if allRes.ClassesUpdated != 1 {
		t.Errorf("expected 1 class updated, got %d", allRes.ClassesUpdated)
	}
}

func TestTeacherQuickPreferences(t *testing.T) {
	ctx := context.Background()
	repo := newMockTimetableRepo()
	svc := NewService(repo, nil)
	schoolID := "school-test-1"

	// 1. Get Teachers Quick Preferences
	res, err := svc.GetTeachersQuickPreferences(ctx, schoolID, nil)
	if err != nil {
		t.Fatalf("failed to get teachers quick preferences: %v", err)
	}
	if len(res.Teachers) != 2 {
		t.Errorf("expected 2 teachers, got %d", len(res.Teachers))
	}
	if res.Teachers[0].DayOff != 1 || res.Teachers[0].TimeSlotPref != "early_hours" {
		t.Errorf("unexpected teacher 0 quick preferences: %+v", res.Teachers[0])
	}
	if res.DayOffCounts[1] != 1 || res.DayOffCounts[3] != 1 {
		t.Errorf("unexpected day off counts: %+v", res.DayOffCounts)
	}

	// 2. Save Teachers Quick Preferences
	saveReq := SaveTeacherQuickPreferencesRequest{
		Preferences: []TeacherQuickPreferenceItem{
			{
				TeacherID:    "t-1",
				DayOff:       2,
				TimeSlotPref: "late_hours",
			},
		},
	}
	if err := svc.SaveTeacherQuickPreferences(ctx, schoolID, saveReq); err != nil {
		t.Fatalf("failed to save teacher quick preferences: %v", err)
	}
}
