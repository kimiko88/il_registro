package timetablegen

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockTimetableRepo struct {
	preferences []TeacherPreference
	reqs        []SubjectRoomRequirement
	constraints []TimetableConstraint
	jobs        map[string]*TimetableJob
	published   []GeneratedSlot
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

func (m *mockTimetableRepo) PublishGeneratedSchedule(ctx context.Context, schoolID string, slots []GeneratedSlot) error {
	m.published = slots
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
}
