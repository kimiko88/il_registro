package timetablegen

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrJobNotFound  = errors.New("timetable generation job not found")
	ErrJobNotReady  = errors.New("cannot publish schedule: job has not completed successfully")
	ErrEmptySlots   = errors.New("no slots available to publish")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("insufficient permissions")
)

type Service interface {
	// Preferences
	SaveTeacherPreferences(ctx context.Context, schoolID, teacherID string, req SavePreferencesRequest) error
	GetTeacherPreferences(ctx context.Context, schoolID, teacherID string, academicYearID *string) ([]TeacherPreference, error)
	IsDesiderataWindowOpen(ctx context.Context, schoolID string) (bool, error)
	SetDesiderataWindowOpen(ctx context.Context, schoolID string, isOpen bool) error

	// Room Requirements
	ListRoomRequirements(ctx context.Context, schoolID string) ([]SubjectRoomRequirement, error)
	SaveRoomRequirement(ctx context.Context, schoolID string, req SaveRoomRequirementRequest) (*SubjectRoomRequirement, error)
	DeleteRoomRequirement(ctx context.Context, id string) error

	// Constraints
	ListConstraints(ctx context.Context, schoolID string) ([]TimetableConstraint, error)
	SaveConstraint(ctx context.Context, schoolID string, req SaveConstraintRequest) (*TimetableConstraint, error)
	DeleteConstraint(ctx context.Context, id string) error

	// Generation Jobs
	StartGeneration(ctx context.Context, schoolID, userID string, req GenerateTimetableRequest) (string, error)
	GetJobStatus(ctx context.Context, schoolID, jobID string) (*TimetableJob, error)
	PublishSchedule(ctx context.Context, schoolID, userID, jobID string) error
	AdjustJobSlots(ctx context.Context, schoolID, userID, jobID string, req AdjustTimetableRequest) (*TimetableGenerationResult, error)

	// Academic Years, Curriculum Plans & Class Daily Limits
	ListAcademicYears(ctx context.Context, schoolID string) ([]string, error)
	ListClassesCurriculumPlans(ctx context.Context, schoolID, academicYear string) ([]ClassCurriculumPlan, error)
	GetClassCurriculumPlan(ctx context.Context, schoolID, classID string) (*ClassCurriculumPlan, error)
	SaveClassCurriculumPlan(ctx context.Context, schoolID, classID string, req SaveClassCurriculumPlanRequest) error
	InheritClassCurriculumPlan(ctx context.Context, schoolID, targetClassID, sourceAcademicYear string) (*ClassCurriculumPlan, error)
	InheritAllClassesCurriculumPlans(ctx context.Context, schoolID, sourceAcademicYear string) (*InheritAllResult, error)

	// Teacher Quick Preferences
	GetTeachersQuickPreferences(ctx context.Context, schoolID string, academicYearID *string) (*TeacherQuickPreferencesOverviewResponse, error)
	SaveTeacherQuickPreferences(ctx context.Context, schoolID string, req SaveTeacherQuickPreferencesRequest) error
}

type service struct {
	repo      Repository
	generator *TimetableGenerator
}

func NewService(repo Repository, generator *TimetableGenerator) Service {
	if repo == nil {
		panic("timetablegen.NewService: repo must not be nil")
	}
	if generator == nil {
		generator = NewGenerator(DefaultConfig())
	}
	return &service{
		repo:      repo,
		generator: generator,
	}
}

// ----------------- Preferences -----------------

func (s *service) SaveTeacherPreferences(ctx context.Context, schoolID, teacherID string, req SavePreferencesRequest) error {
	if schoolID == "" || teacherID == "" {
		return errors.New("school_id and teacher_id are required")
	}
	return s.repo.SavePreferencesBatch(ctx, schoolID, teacherID, req.AcademicYearID, req.Preferences)
}

func (s *service) GetTeacherPreferences(ctx context.Context, schoolID, teacherID string, academicYearID *string) ([]TeacherPreference, error) {
	return s.repo.GetTeacherPreferences(ctx, schoolID, teacherID, academicYearID)
}

func (s *service) IsDesiderataWindowOpen(ctx context.Context, schoolID string) (bool, error) {
	return s.repo.GetDesiderataWindow(ctx, schoolID)
}

func (s *service) SetDesiderataWindowOpen(ctx context.Context, schoolID string, isOpen bool) error {
	return s.repo.SetDesiderataWindow(ctx, schoolID, isOpen)
}

// ----------------- Room Requirements -----------------

func (s *service) ListRoomRequirements(ctx context.Context, schoolID string) ([]SubjectRoomRequirement, error) {
	return s.repo.ListRoomRequirements(ctx, schoolID)
}

func (s *service) SaveRoomRequirement(ctx context.Context, schoolID string, req SaveRoomRequirementRequest) (*SubjectRoomRequirement, error) {
	return s.repo.SaveRoomRequirement(ctx, schoolID, req)
}

func (s *service) DeleteRoomRequirement(ctx context.Context, id string) error {
	return s.repo.DeleteRoomRequirement(ctx, id)
}

// ----------------- Constraints -----------------

func (s *service) ListConstraints(ctx context.Context, schoolID string) ([]TimetableConstraint, error) {
	return s.repo.ListConstraints(ctx, schoolID)
}

func (s *service) SaveConstraint(ctx context.Context, schoolID string, req SaveConstraintRequest) (*TimetableConstraint, error) {
	return s.repo.SaveConstraint(ctx, schoolID, req)
}

func (s *service) DeleteConstraint(ctx context.Context, id string) error {
	return s.repo.DeleteConstraint(ctx, id)
}

// ----------------- Generation Jobs -----------------

func (s *service) StartGeneration(ctx context.Context, schoolID, userID string, req GenerateTimetableRequest) (string, error) {
	jobID := uuid.New().String()
	paramBytes, _ := json.Marshal(req)

	job := &TimetableJob{
		ID:             jobID,
		SchoolID:       schoolID,
		AcademicYearID: req.AcademicYearID,
		TriggeredBy:    &userID,
		Status:         JobStatusPending,
		Algorithm:      "greedy_local_search",
		Parameters:     paramBytes,
	}

	createdJob, err := s.repo.CreateJob(ctx, job)
	if err != nil {
		return "", fmt.Errorf("failed to create generation job: %w", err)
	}

	// Launch async goroutine with configurable timeout (up to 10 minutes)
	go func(jID, sID string, aYearID *string, timeLimit int, maxIter int) {
		timeout := 60 * time.Second
		if timeLimit > 0 && timeLimit <= 600 {
			timeout = time.Duration(timeLimit) * time.Second
		}
		bgCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		now := time.Now()
		runningJob := &TimetableJob{
			ID:        jID,
			Status:    JobStatusRunning,
			StartedAt: &now,
		}
		_ = s.repo.UpdateJob(bgCtx, runningJob)

		// 1. Load Data
		assignments, err := s.repo.LoadAssignments(bgCtx, sID, aYearID)
		if err != nil {
			s.failJob(bgCtx, jID, fmt.Sprintf("Failed loading assignments: %v", err))
			return
		}

		if len(assignments) == 0 {
			s.failJob(bgCtx, jID, "Nessuna cattedra trovata per questo anno accademico (tabella class_subjects vuota)")
			return
		}

		rooms, err := s.repo.LoadRooms(bgCtx, sID)
		if err != nil {
			log.Printf("Warning: failed to load rooms: %v", err)
			rooms = []RoomData{}
		}

		roomReqs, err := s.repo.LoadRoomRequirements(bgCtx, sID)
		if err != nil {
			log.Printf("Warning: failed to load room requirements: %v", err)
			roomReqs = make(map[string]SubjectRoomRequirement)
		}

		preferences, err := s.repo.LoadAllPreferences(bgCtx, sID, aYearID)
		if err != nil {
			log.Printf("Warning: failed to load preferences: %v", err)
			preferences = []TeacherPreference{}
		}

		constraints, err := s.repo.ListConstraints(bgCtx, sID)
		if err != nil {
			log.Printf("Warning: failed to load constraints: %v", err)
			constraints = []TimetableConstraint{}
		}

		// 2. Run Algorithm with configured limits
		gen := s.generator
		if timeLimit > 0 || maxIter > 0 {
			cfg := s.generator.config
			if timeLimit > 0 {
				cfg.TimeLimitSeconds = timeLimit
			}
			if maxIter > 0 {
				cfg.MaxIterations = maxIter
			}
			gen = NewGenerator(cfg)
		}

		res, err := gen.Generate(bgCtx, assignments, rooms, roomReqs, preferences, constraints)
		if err != nil {
			s.failJob(bgCtx, jID, fmt.Sprintf("Generazione fallita: %v", err))
			return
		}
		res.JobID = jID

		// 3. Mark completed
		resBytes, _ := json.Marshal(res)
		completedAt := time.Now()
		completedJob := &TimetableJob{
			ID:            jID,
			Status:        JobStatusCompleted,
			ResultSummary: resBytes,
			CompletedAt:   &completedAt,
		}
		_ = s.repo.UpdateJob(context.Background(), completedJob)
	}(createdJob.ID, schoolID, req.AcademicYearID, req.TimeLimitSeconds, req.MaxIterations)

	return createdJob.ID, nil
}

func (s *service) failJob(ctx context.Context, jobID string, errMsg string) {
	now := time.Now()
	failedJob := &TimetableJob{
		ID:           jobID,
		Status:       JobStatusFailed,
		ErrorMessage: &errMsg,
		CompletedAt:  &now,
	}
	_ = s.repo.UpdateJob(ctx, failedJob)
}

func (s *service) GetJobStatus(ctx context.Context, schoolID, jobID string) (*TimetableJob, error) {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return nil, ErrJobNotFound
	}
	if job.SchoolID != schoolID {
		return nil, ErrForbidden
	}
	return job, nil
}

func (s *service) PublishSchedule(ctx context.Context, schoolID, userID, jobID string) error {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return ErrJobNotFound
	}
	if job.SchoolID != schoolID {
		return ErrForbidden
	}
	if job.Status != JobStatusCompleted {
		return ErrJobNotReady
	}

	var result TimetableGenerationResult
	if err := json.Unmarshal(job.ResultSummary, &result); err != nil {
		return fmt.Errorf("failed to parse job result summary: %w", err)
	}

	if len(result.Slots) == 0 {
		return ErrEmptySlots
	}

	return s.repo.PublishGeneratedSchedule(ctx, schoolID, result.Slots)
}

func (s *service) AdjustJobSlots(ctx context.Context, schoolID, userID, jobID string, req AdjustTimetableRequest) (*TimetableGenerationResult, error) {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return nil, ErrJobNotFound
	}
	if job.SchoolID != schoolID {
		return nil, ErrForbidden
	}
	if job.Status != JobStatusCompleted {
		return nil, errors.New("cannot adjust an uncompleted timetable generation job")
	}

	var result TimetableGenerationResult
	if len(job.ResultSummary) > 0 {
		_ = json.Unmarshal(job.ResultSummary, &result)
	}

	type slotKey struct {
		day  int
		hour int
	}
	classSlots := make(map[string]map[slotKey]string)
	teacherSlots := make(map[string]map[slotKey]string)
	roomSlots := make(map[string]map[slotKey]string)

	var hardConflicts []HardConflict

	for _, slot := range req.Slots {
		sk := slotKey{day: slot.DayOfWeek, hour: slot.HourIndex}

		// Class overlap check
		if slot.ClassID != "" {
			if _, ok := classSlots[slot.ClassID]; !ok {
				classSlots[slot.ClassID] = make(map[slotKey]string)
			}
			if prevSubj, exists := classSlots[slot.ClassID][sk]; exists {
				hardConflicts = append(hardConflicts, HardConflict{
					Type:        "class_overlap",
					Description: fmt.Sprintf("Sovrapposizione nella classe %s: sia %s che %s assegnate al giorno %d, ora %d", slot.ClassName, prevSubj, slot.SubjectName, slot.DayOfWeek, slot.HourIndex),
					ClassID:     slot.ClassID,
				})
			} else {
				classSlots[slot.ClassID][sk] = slot.SubjectName
			}
		}

		// Teacher overlap check
		if slot.TeacherID != nil && *slot.TeacherID != "" && !strings.HasPrefix(*slot.TeacherID, "unassigned-") && !strings.HasPrefix(*slot.TeacherID, "spezzone-") {
			tID := *slot.TeacherID
			if _, ok := teacherSlots[tID]; !ok {
				teacherSlots[tID] = make(map[slotKey]string)
			}
			if prevSubj, exists := teacherSlots[tID][sk]; exists {
				if prevSubj != slot.SubjectName {
					hardConflicts = append(hardConflicts, HardConflict{
						Type:        "teacher_overlap",
						Description: fmt.Sprintf("Docente %s sovrapposto: assegnato contemporaneamente a %s e %s al giorno %d, ora %d", slot.TeacherName, prevSubj, slot.SubjectName, slot.DayOfWeek, slot.HourIndex),
						TeacherID:   tID,
						DayOfWeek:   slot.DayOfWeek,
						HourIndex:   slot.HourIndex,
					})
				}
			} else {
				teacherSlots[tID][sk] = slot.SubjectName
			}
		}

		// Room overlap check
		if slot.RoomID != nil && *slot.RoomID != "" {
			rID := *slot.RoomID
			if _, ok := roomSlots[rID]; !ok {
				roomSlots[rID] = make(map[slotKey]string)
			}
			if prevClass, exists := roomSlots[rID][sk]; exists {
				hardConflicts = append(hardConflicts, HardConflict{
					Type:        "room_overlap",
					Description: fmt.Sprintf("Aula %s occupata contemporaneamente da %s e %s al giorno %d, ora %d", slot.RoomName, prevClass, slot.ClassName, slot.DayOfWeek, slot.HourIndex),
				})
			} else {
				roomSlots[rID][sk] = slot.ClassName
			}
		}
	}

	result.Slots = req.Slots
	result.AssignedSlots = len(req.Slots)
	result.HardConflicts = hardConflicts
	if result.TotalSlots > 0 {
		result.CoveragePct = float64(result.AssignedSlots) / float64(result.TotalSlots) * 100.0
	}

	resBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize adjusted timetable: %w", err)
	}

	job.ResultSummary = resBytes
	if err := s.repo.UpdateJob(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to update adjusted timetable job: %w", err)
	}

	return &result, nil
}

// ----------------- Curriculum Plans & Class Daily Limits -----------------

func (s *service) ListAcademicYears(ctx context.Context, schoolID string) ([]string, error) {
	return s.repo.ListAcademicYears(ctx, schoolID)
}

func (s *service) ListClassesCurriculumPlans(ctx context.Context, schoolID, academicYear string) ([]ClassCurriculumPlan, error) {
	return s.repo.ListClassesCurriculumPlans(ctx, schoolID, academicYear)
}

func (s *service) GetClassCurriculumPlan(ctx context.Context, schoolID, classID string) (*ClassCurriculumPlan, error) {
	return s.repo.GetClassCurriculumPlan(ctx, schoolID, classID)
}

func (s *service) SaveClassCurriculumPlan(ctx context.Context, schoolID, classID string, req SaveClassCurriculumPlanRequest) error {
	if req.MaxHoursPerDay <= 0 {
		req.MaxHoursPerDay = 6
	}
	if req.MinHoursPerDay <= 0 {
		req.MinHoursPerDay = 4
	}
	if req.MinHoursPerDay > req.MaxHoursPerDay {
		req.MinHoursPerDay = req.MaxHoursPerDay
	}
	return s.repo.SaveClassCurriculumPlan(ctx, schoolID, classID, req)
}

func (s *service) InheritClassCurriculumPlan(ctx context.Context, schoolID, targetClassID, sourceAcademicYear string) (*ClassCurriculumPlan, error) {
	return s.repo.InheritClassCurriculumPlan(ctx, schoolID, targetClassID, sourceAcademicYear)
}

func (s *service) InheritAllClassesCurriculumPlans(ctx context.Context, schoolID, sourceAcademicYear string) (*InheritAllResult, error) {
	return s.repo.InheritAllClassesCurriculumPlans(ctx, schoolID, sourceAcademicYear)
}

// ----------------- Teacher Quick Preferences -----------------

func (s *service) GetTeachersQuickPreferences(ctx context.Context, schoolID string, academicYearID *string) (*TeacherQuickPreferencesOverviewResponse, error) {
	teachers, dayOffCounts, err := s.repo.GetTeachersQuickPreferences(ctx, schoolID, academicYearID)
	if err != nil {
		return nil, err
	}
	return &TeacherQuickPreferencesOverviewResponse{
		AcademicYearID: academicYearID,
		Teachers:       teachers,
		DayOffCounts:   dayOffCounts,
	}, nil
}

func (s *service) SaveTeacherQuickPreferences(ctx context.Context, schoolID string, req SaveTeacherQuickPreferencesRequest) error {
	return s.repo.SaveTeacherQuickPreferences(ctx, schoolID, req.AcademicYearID, req.Preferences)
}
