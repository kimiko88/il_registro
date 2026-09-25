package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/timetablegen"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockIntegrationTimetableRepo struct {
	preferences      []timetablegen.TeacherPreference
	reqs             []timetablegen.SubjectRoomRequirement
	constraints      []timetablegen.TimetableConstraint
	jobs             map[string]*timetablegen.TimetableJob
	published        []timetablegen.GeneratedSlot
	assignments      []timetablegen.AssignmentData
	rooms            []timetablegen.RoomData
	associatedGroups []timetablegen.AssociatedGroup
}

func newMockIntegrationTimetableRepo() *mockIntegrationTimetableRepo {
	return &mockIntegrationTimetableRepo{
		jobs: make(map[string]*timetablegen.TimetableJob),
	}
}

func (m *mockIntegrationTimetableRepo) SavePreferencesBatch(ctx context.Context, schoolID, teacherID string, academicYearID *string, prefs []timetablegen.PreferenceEntry) error {
	for _, p := range prefs {
		m.preferences = append(m.preferences, timetablegen.TeacherPreference{
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

func (m *mockIntegrationTimetableRepo) GetTeacherPreferences(ctx context.Context, schoolID, teacherID string, academicYearID *string) ([]timetablegen.TeacherPreference, error) {
	var res []timetablegen.TeacherPreference
	for _, p := range m.preferences {
		if p.SchoolID == schoolID && p.TeacherID == teacherID {
			res = append(res, p)
		}
	}
	return res, nil
}

func (m *mockIntegrationTimetableRepo) LoadAllPreferences(ctx context.Context, schoolID string, academicYearID *string) ([]timetablegen.TeacherPreference, error) {
	return m.preferences, nil
}

func (m *mockIntegrationTimetableRepo) ListRoomRequirements(ctx context.Context, schoolID string) ([]timetablegen.SubjectRoomRequirement, error) {
	return m.reqs, nil
}

func (m *mockIntegrationTimetableRepo) SaveRoomRequirement(ctx context.Context, schoolID string, req timetablegen.SaveRoomRequirementRequest) (*timetablegen.SubjectRoomRequirement, error) {
	labHours := req.LabHours
	if labHours <= 0 {
		labHours = 1
	}
	r := &timetablegen.SubjectRoomRequirement{
		ID:               uuid.New().String(),
		SchoolID:         schoolID,
		SubjectID:        req.SubjectID,
		RequiredRoomType: req.RequiredRoomType,
		IsMandatory:      req.IsMandatory,
		LabHours:         labHours,
	}
	m.reqs = append(m.reqs, *r)
	return r, nil
}

func (m *mockIntegrationTimetableRepo) DeleteRoomRequirement(ctx context.Context, id string) error {
	var filtered []timetablegen.SubjectRoomRequirement
	for _, r := range m.reqs {
		if r.ID != id {
			filtered = append(filtered, r)
		}
	}
	m.reqs = filtered
	return nil
}

func (m *mockIntegrationTimetableRepo) ListConstraints(ctx context.Context, schoolID string) ([]timetablegen.TimetableConstraint, error) {
	return m.constraints, nil
}

func (m *mockIntegrationTimetableRepo) SaveConstraint(ctx context.Context, schoolID string, req timetablegen.SaveConstraintRequest) (*timetablegen.TimetableConstraint, error) {
	c := &timetablegen.TimetableConstraint{
		ID:             uuid.New().String(),
		SchoolID:       schoolID,
		ConstraintType: req.ConstraintType,
		TargetType:     req.TargetType,
		TargetID:       req.TargetID,
		Parameters:     req.Parameters,
		IsHard:         req.IsHard,
	}
	m.constraints = append(m.constraints, *c)
	return c, nil
}

func (m *mockIntegrationTimetableRepo) DeleteConstraint(ctx context.Context, id string) error {
	var filtered []timetablegen.TimetableConstraint
	for _, c := range m.constraints {
		if c.ID != id {
			filtered = append(filtered, c)
		}
	}
	m.constraints = filtered
	return nil
}

func (m *mockIntegrationTimetableRepo) GetDesiderataWindow(ctx context.Context, schoolID string) (bool, error) {
	for _, c := range m.constraints {
		if c.SchoolID == schoolID && c.ConstraintType == "teacher_desiderata_window" {
			return c.IsActive, nil
		}
	}
	return false, nil
}

func (m *mockIntegrationTimetableRepo) SetDesiderataWindow(ctx context.Context, schoolID string, isOpen bool) error {
	for i, c := range m.constraints {
		if c.SchoolID == schoolID && c.ConstraintType == "teacher_desiderata_window" {
			m.constraints[i].IsActive = isOpen
			return nil
		}
	}
	m.constraints = append(m.constraints, timetablegen.TimetableConstraint{
		ID:             uuid.New().String(),
		SchoolID:       schoolID,
		ConstraintType: "teacher_desiderata_window",
		IsActive:       isOpen,
	})
	return nil
}

func (m *mockIntegrationTimetableRepo) CreateJob(ctx context.Context, job *timetablegen.TimetableJob) (*timetablegen.TimetableJob, error) {
	if job.ID == "" {
		job.ID = uuid.New().String()
	}
	job.CreatedAt = time.Now()
	m.jobs[job.ID] = job
	return job, nil
}

func (m *mockIntegrationTimetableRepo) UpdateJob(ctx context.Context, job *timetablegen.TimetableJob) error {
	existing, ok := m.jobs[job.ID]
	if ok {
		if job.Status != "" {
			existing.Status = job.Status
		}
		if job.StartedAt != nil {
			existing.StartedAt = job.StartedAt
		}
		if job.CompletedAt != nil {
			existing.CompletedAt = job.CompletedAt
		}
		if len(job.ResultSummary) > 0 {
			existing.ResultSummary = job.ResultSummary
		}
		if job.ErrorMessage != nil {
			existing.ErrorMessage = job.ErrorMessage
		}
	} else {
		m.jobs[job.ID] = job
	}
	return nil
}

func (m *mockIntegrationTimetableRepo) GetJob(ctx context.Context, id string) (*timetablegen.TimetableJob, error) {
	job, ok := m.jobs[id]
	if !ok {
		return nil, fmt.Errorf("job not found")
	}
	return job, nil
}

func (m *mockIntegrationTimetableRepo) LoadAssignments(ctx context.Context, schoolID string, academicYearID *string) ([]timetablegen.AssignmentData, error) {
	return m.assignments, nil
}

func (m *mockIntegrationTimetableRepo) LoadRooms(ctx context.Context, schoolID string) ([]timetablegen.RoomData, error) {
	return m.rooms, nil
}

func (m *mockIntegrationTimetableRepo) LoadRoomRequirements(ctx context.Context, schoolID string) (map[string]timetablegen.SubjectRoomRequirement, error) {
	res := make(map[string]timetablegen.SubjectRoomRequirement)
	for _, r := range m.reqs {
		res[r.SubjectID] = r
	}
	return res, nil
}

func (m *mockIntegrationTimetableRepo) LoadAssociatedGroups(ctx context.Context, schoolID string) ([]timetablegen.AssociatedGroup, error) {
	return m.associatedGroups, nil
}

func (m *mockIntegrationTimetableRepo) PublishGeneratedSchedule(ctx context.Context, schoolID string, slots []timetablegen.GeneratedSlot) error {
	m.published = append(m.published, slots...)
	return nil
}

func setupTimetableIntegrationRouter(repo timetablegen.Repository, role, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	generator := timetablegen.NewGenerator(timetablegen.DefaultConfig())
	svc := timetablegen.NewService(repo, generator)
	handler := timetablegen.NewHandler(svc)

	r.Use(func(c *gin.Context) {
		c.Set("role", role)
		c.Set("user_id", userID)
		c.Set("school_id", schoolID)
		c.Next()
	})

	api := r.Group("")
	handler.RegisterRoutes(api)

	return r
}

func TestIntegration_TimetableGenerationLifecycle(t *testing.T) {
	repo := newMockIntegrationTimetableRepo()
	schoolID := "school-timetables-1"
	teacherSeniorID := "teacher-senior"
	teacherJuniorID := "teacher-junior"
	vicePrincipalID := "vp-1"

	hiringSenior := time.Date(2008, 9, 1, 0, 0, 0, 0, time.UTC)
	hiringJunior := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)

	bldID := "bld-central"
	repo.assignments = []timetablegen.AssignmentData{
		{
			ClassID:       "class-1A",
			ClassName:     "1A",
			BuildingID:    &bldID,
			SubjectID:     "sub-info",
			SubjectName:   "Informatica",
			TeacherID:     teacherSeniorID,
			TeacherUserID: teacherSeniorID,
			TeacherName:   "Prof Senior",
			HiringDate:    &hiringSenior,
			HoursPerWeek:  2,
		},
		{
			ClassID:       "class-2B",
			ClassName:     "2B",
			BuildingID:    &bldID,
			SubjectID:     "sub-arte",
			SubjectName:   "Arte",
			TeacherID:     teacherJuniorID,
			TeacherUserID: teacherJuniorID,
			TeacherName:   "Prof Junior",
			HiringDate:    &hiringJunior,
			HoursPerWeek:  2,
		},
	}

	repo.rooms = []timetablegen.RoomData{
		{
			ID:         "room-lab-central",
			BuildingID: &bldID,
			Name:       "Lab Info Centrale",
			RoomType:   "lab_informatica",
			Capacity:   28,
			IsActive:   true,
		},
	}

	// 0. Teacher Senior tries to save preferences before window is open -> 403 Forbidden
	routerTeacherSenior := setupTimetableIntegrationRouter(repo, "teacher", teacherSeniorID, schoolID)
	savePrefReq := timetablegen.SavePreferencesRequest{
		Preferences: []timetablegen.PreferenceEntry{
			{DayOfWeek: 1, HourIndex: 1, PreferenceType: timetablegen.PrefPreferred},
			{DayOfWeek: 1, HourIndex: 2, PreferenceType: timetablegen.PrefPreferred},
		},
	}
	body, _ := json.Marshal(savePrefReq)
	req := httptest.NewRequest(http.MethodPost, "/timetable/preferences", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routerTeacherSenior.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 0b. Vice-Principal opens the desiderata window -> 200 OK
	routerVPInitial := setupTimetableIntegrationRouter(repo, "vice_principal", vicePrincipalID, schoolID)
	openWinReq := timetablegen.SetDesiderataWindowRequest{IsOpen: true}
	winBody, _ := json.Marshal(openWinReq)
	winReq := httptest.NewRequest(http.MethodPost, "/timetable/preferences/window", bytes.NewBuffer(winBody))
	winReq.Header.Set("Content-Type", "application/json")
	wWin := httptest.NewRecorder()
	routerVPInitial.ServeHTTP(wWin, winReq)
	assert.Equal(t, http.StatusOK, wWin.Code)

	// 1. Teacher Senior sets Monday 1st and 2nd hour as Preferred with window open -> 200 OK
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/timetable/preferences", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	routerTeacherSenior.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. Teacher Junior tries to trigger Timetable Generation -> 403 Forbidden
	routerTeacherJunior := setupTimetableIntegrationRouter(repo, "teacher", teacherJuniorID, schoolID)
	genReq := timetablegen.GenerateTimetableRequest{TimeLimitSeconds: 10}
	body, _ = json.Marshal(genReq)
	req = httptest.NewRequest(http.MethodPost, "/timetable/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerTeacherJunior.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	// 3. Vice-Principal configures Informatics to require lab_informatica
	routerVP := setupTimetableIntegrationRouter(repo, "vice_principal", vicePrincipalID, schoolID)
	reqRoomReq := timetablegen.SaveRoomRequirementRequest{
		SubjectID:        "sub-info",
		RequiredRoomType: "lab_informatica",
		IsMandatory:      true,
		LabHours:         2,
	}
	body, _ = json.Marshal(reqRoomReq)
	req = httptest.NewRequest(http.MethodPost, "/timetable/room-requirements", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerVP.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Vice-Principal triggers automatic timetable generation
	genReqVP := timetablegen.GenerateTimetableRequest{TimeLimitSeconds: 10}
	body, _ = json.Marshal(genReqVP)
	req = httptest.NewRequest(http.MethodPost, "/timetable/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerVP.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var genResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &genResp)
	assert.NoError(t, err)
	jobID, ok := genResp["job_id"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, jobID)

	// 5. Poll for completion
	var jobResp timetablegen.TimetableJob
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		req = httptest.NewRequest(http.MethodGet, "/timetable/generate/"+jobID, nil)
		w = httptest.NewRecorder()
		routerVP.ServeHTTP(w, req)
		if w.Code == http.StatusOK {
			_ = json.Unmarshal(w.Body.Bytes(), &jobResp)
			if jobResp.Status == timetablegen.JobStatusCompleted {
				break
			}
		}
	}

	assert.Equal(t, timetablegen.JobStatusCompleted, jobResp.Status)

	var result timetablegen.TimetableGenerationResult
	err = json.Unmarshal(jobResp.ResultSummary, &result)
	assert.NoError(t, err)
	assert.Equal(t, 4, result.AssignedSlots)
	assert.Equal(t, 100.0, result.CoveragePct)

	// Verify that Informatics got the lab room assigned in matching building
	var infoSlotsAssignedLab int
	for _, slot := range result.Slots {
		if slot.SubjectID == "sub-info" {
			assert.NotNil(t, slot.RoomID)
			assert.Equal(t, "room-lab-central", *slot.RoomID)
			infoSlotsAssignedLab++
		}
	}
	assert.Equal(t, 2, infoSlotsAssignedLab)

	// 6. Vice-Principal publishes the schedule into class_schedules
	publishReq := httptest.NewRequest(http.MethodPost, "/timetable/generate/"+jobID+"/publish", nil)
	w = httptest.NewRecorder()
	routerVP.ServeHTTP(w, publishReq)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, repo.published, 4)
}
