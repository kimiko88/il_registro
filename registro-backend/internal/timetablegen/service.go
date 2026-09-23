package timetablegen

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

	// Launch async goroutine with 30-second context
	go func(jID, sID string, aYearID *string, timeLimit int) {
		timeout := 25 * time.Second
		if timeLimit > 0 && timeLimit <= 30 {
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

		// 2. Run Algorithm
		res, err := s.generator.Generate(bgCtx, assignments, rooms, roomReqs, preferences, constraints)
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
	}(createdJob.ID, schoolID, req.AcademicYearID, req.TimeLimitSeconds)

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
