package support

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"registro-backend/internal/users"
)

var (
	ErrEmptyTopic = errors.New("argomenti e attività del diario di sostegno obbligatori")
)

type Service interface {
	CreateDiaryEntry(ctx context.Context, schoolID, teacherID string, req CreateDiaryEntryRequest) (*SupportDiaryEntry, error)
	ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]SupportDiaryEntry, error)
	DeleteDiaryEntry(ctx context.Context, id, teacherID string) error

	CreatePeiGoal(ctx context.Context, schoolID string, req CreatePeiGoalRequest) (*SupportPeiGoal, error)
	UpdateGoalProgress(ctx context.Context, id, status string) error
	ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]SupportPeiGoal, error)
	DeletePeiGoal(ctx context.Context, id string) error
	GetUserRepo() users.Repository

	CreateGloMeeting(ctx context.Context, schoolID string, req CreateGloMeetingRequest) (*SupportGloMeeting, error)
	GetGloMeeting(ctx context.Context, id string) (*SupportGloMeeting, error)
	ListGloMeetings(ctx context.Context, schoolID, studentID string) ([]SupportGloMeeting, error)
	SignGloMeeting(ctx context.Context, memberID, memberName string, req SignGloMeetingRequest) (*GloSignature, error)
	SaveSupportHours(ctx context.Context, schoolID string, req SaveSupportHoursRequest) (*SupportHoursQuantification, error)
	GetSupportHours(ctx context.Context, schoolID, studentID, schoolYear, peiStage string) (*SupportHoursQuantification, error)
}

type service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(repo Repository, uRepo ...users.Repository) Service {
	s := &service{repo: repo}
	if len(uRepo) > 0 && uRepo[0] != nil {
		s.userRepo = uRepo[0]
	}
	return s
}

func (s *service) GetUserRepo() users.Repository {
	return s.userRepo
}

func (s *service) CreateDiaryEntry(ctx context.Context, schoolID, teacherID string, req CreateDiaryEntryRequest) (*SupportDiaryEntry, error) {
	if req.TopicAndActivities == "" {
		return nil, ErrEmptyTopic
	}
	if req.ActivityType == "" {
		req.ActivityType = "in_classe"
	}

	entry := &SupportDiaryEntry{
		SchoolID:           schoolID,
		TeacherID:          teacherID,
		StudentID:          req.StudentID,
		ClassID:            req.ClassID,
		EntryDate:          req.EntryDate,
		TimeSlot:           req.TimeSlot,
		CoTeacherID:        req.CoTeacherID,
		ActivityType:       req.ActivityType,
		TopicAndActivities: req.TopicAndActivities,
		StudentResponses:   req.StudentResponses,
		EducatorNotes:      req.EducatorNotes,
		IsSharedWithFamily: req.IsSharedWithFamily,
	}

	if err := s.repo.CreateDiaryEntry(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *service) ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]SupportDiaryEntry, error) {
	return s.repo.ListDiaryEntries(ctx, schoolID, studentID, teacherID, classID, isFamily)
}

func (s *service) DeleteDiaryEntry(ctx context.Context, id, teacherID string) error {
	return s.repo.DeleteDiaryEntry(ctx, id, teacherID)
}

func (s *service) CreatePeiGoal(ctx context.Context, schoolID string, req CreatePeiGoalRequest) (*SupportPeiGoal, error) {
	if req.PeiType == "" {
		req.PeiType = "equipollente"
	}
	if req.Axis == "" {
		req.Axis = "autonomia"
	}
	if req.ProgressStatus == "" {
		req.ProgressStatus = "non_avviato"
	}
	if req.ExpectedTerm == "" {
		req.ExpectedTerm = "annuale"
	}

	dimension := req.MinisterialDimension
	if dimension == "" {
		switch req.Axis {
		case "relazionale":
			dimension = "dimensione_relazione"
		case "comunicazionale", "linguistica":
			dimension = "dimensione_comunicazione"
		case "cognitiva":
			dimension = "dimensione_cognitiva"
		default:
			dimension = "dimensione_autonomia"
		}
	}

	pathway := req.PathwayType
	if pathway == "" {
		if req.PeiType == "differenziato" {
			pathway = "percorso_c_differenziato"
		} else {
			pathway = "percorso_b_personalizzato"
		}
	}

	goal := &SupportPeiGoal{
		SchoolID:             schoolID,
		StudentID:            req.StudentID,
		PeiType:              req.PeiType,
		Axis:                 req.Axis,
		MinisterialDimension: dimension,
		PathwayType:          pathway,
		Title:                req.Title,
		Description:          req.Description,
		ExpectedTerm:         req.ExpectedTerm,
		ProgressStatus:       req.ProgressStatus,
		GloNotes:             req.GloNotes,
	}

	if err := s.repo.CreatePeiGoal(ctx, goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *service) UpdateGoalProgress(ctx context.Context, id, status string) error {
	return s.repo.UpdatePeiGoalProgress(ctx, id, status)
}

func (s *service) ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]SupportPeiGoal, error) {
	return s.repo.ListPeiGoals(ctx, schoolID, studentID)
}

func (s *service) DeletePeiGoal(ctx context.Context, id string) error {
	return s.repo.DeletePeiGoal(ctx, id)
}

func (s *service) CreateGloMeeting(ctx context.Context, schoolID string, req CreateGloMeetingRequest) (*SupportGloMeeting, error) {
	if req.MeetingDate == "" {
		req.MeetingDate = time.Now().Format("2006-01-02")
	}
	if req.SchoolYear == "" {
		req.SchoolYear = "2025/2026"
	}
	if req.MeetingType == "" {
		req.MeetingType = "INIZIALE"
	}

	meeting := &SupportGloMeeting{
		SchoolID:          schoolID,
		StudentID:         req.StudentID,
		MeetingDate:       req.MeetingDate,
		MeetingType:       req.MeetingType,
		SchoolYear:        req.SchoolYear,
		MembersPresent:    req.MembersPresent,
		DiscussionSummary: req.DiscussionSummary,
		DecisionsTaken:    req.DecisionsTaken,
		Status:            "IN_FIRMA",
	}

	if err := s.repo.CreateGloMeeting(ctx, meeting); err != nil {
		return nil, err
	}

	return meeting, nil
}

func (s *service) GetGloMeeting(ctx context.Context, id string) (*SupportGloMeeting, error) {
	return s.repo.GetGloMeeting(ctx, id)
}

func (s *service) ListGloMeetings(ctx context.Context, schoolID, studentID string) ([]SupportGloMeeting, error) {
	return s.repo.ListGloMeetings(ctx, schoolID, studentID)
}

func (s *service) SignGloMeeting(ctx context.Context, memberID, memberName string, req SignGloMeetingRequest) (*GloSignature, error) {
	meeting, err := s.repo.GetGloMeeting(ctx, req.MeetingID)
	if err != nil {
		return nil, fmt.Errorf("verbale GLO non trovato: %w", err)
	}

	sigType := req.SignatureType
	if sigType == "" {
		sigType = "DIGITAL_AGID"
	}

	// Compute AgID-compliant hash representation
	hashInput := fmt.Sprintf("%s:%s:%s:%s:%d", meeting.ID, memberID, req.MemberRole, sigType, time.Now().UnixNano())
	hashBytes := sha256.Sum256([]byte(hashInput))
	sigHash := fmt.Sprintf("sha256:%x", hashBytes)

	now := time.Now()
	sig := &GloSignature{
		MeetingID:     req.MeetingID,
		MemberID:      memberID,
		MemberName:    memberName,
		MemberRole:    req.MemberRole,
		SignedAt:      &now,
		SignatureType: sigType,
		SignatureHash: sigHash,
		Status:        "SIGNED",
	}

	if err := s.repo.SignGloMeeting(ctx, req.MeetingID, sig); err != nil {
		return nil, err
	}

	return sig, nil
}

func (s *service) SaveSupportHours(ctx context.Context, schoolID string, req SaveSupportHoursRequest) (*SupportHoursQuantification, error) {
	if req.SchoolYear == "" {
		req.SchoolYear = "2025/2026"
	}
	if req.PeiStage == "" {
		req.PeiStage = "DEFINITIVO"
	}

	// Range validation (0 to 4 ICF)
	if req.DebtRelazione < 0 || req.DebtRelazione > 4 {
		req.DebtRelazione = 0
	}
	if req.DebtComunicazione < 0 || req.DebtComunicazione > 4 {
		req.DebtComunicazione = 0
	}
	if req.DebtAutonomia < 0 || req.DebtAutonomia > 4 {
		req.DebtAutonomia = 0
	}
	if req.DebtCognitiva < 0 || req.DebtCognitiva > 4 {
		req.DebtCognitiva = 0
	}

	h := &SupportHoursQuantification{
		SchoolID:                    schoolID,
		StudentID:                   req.StudentID,
		SchoolYear:                  req.SchoolYear,
		PeiStage:                    req.PeiStage,
		DebtRelazione:               req.DebtRelazione,
		DebtComunicazione:           req.DebtComunicazione,
		DebtAutonomia:               req.DebtAutonomia,
		DebtCognitiva:               req.DebtCognitiva,
		WeeklySupportHoursRequested: req.WeeklySupportHoursRequested,
		WeeklySupportHoursAssigned:  req.WeeklySupportHoursAssigned,
		AssistantHoursRequested:     req.AssistantHoursRequested,
		AssistantHoursAssigned:      req.AssistantHoursAssigned,
		TransportationNeeded:        req.TransportationNeeded,
		PersonalCareNeeded:          req.PersonalCareNeeded,
		JustificationSummary:        req.JustificationSummary,
	}

	if err := s.repo.SaveSupportHours(ctx, h); err != nil {
		return nil, err
	}

	return h, nil
}

func (s *service) GetSupportHours(ctx context.Context, schoolID, studentID, schoolYear, peiStage string) (*SupportHoursQuantification, error) {
	if schoolYear == "" {
		schoolYear = "2025/2026"
	}
	if peiStage == "" {
		peiStage = "DEFINITIVO"
	}
	return s.repo.GetSupportHours(ctx, schoolID, studentID, schoolYear, peiStage)
}
