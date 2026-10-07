package psychology

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrProfessionalSecrecyViolation = errors.New("accesso vietato ex L. 56/1989: le note cliniche dello sportello psicologico sono protette da segreto professionale e inaccessibili a docenti, presidenza e segreteria")
	ErrBothParentsConsentRequired   = errors.New("per gli studenti minorenni è obbligatorio il consenso informato preventivo sottoscritto da entrambi i genitori prima di accedere al colloquio")
	ErrSessionNotFound              = errors.New("colloquio psicologico non trovato")
	ErrSlotAlreadyBooked            = errors.New("orario non disponibile o già prenotato")
)

type Repository interface {
	SaveConsent(ctx context.Context, c *PsychologyConsent) error
	GetConsent(ctx context.Context, studentID, schoolYear string) (*PsychologyConsent, error)
	SaveSession(ctx context.Context, s *PsychologySession) error
	GetSessionByID(ctx context.Context, id string) (*PsychologySession, error)
	ListSessionsByStudent(ctx context.Context, studentID string) ([]*PsychologySession, error)
	ListSessionsByPsychologist(ctx context.Context, psychologistID string) ([]*PsychologySession, error)
}

type MemoryRepository struct {
	mu       sync.RWMutex
	consents map[string]*PsychologyConsent // key: studentID:year
	sessions map[string]*PsychologySession // key: sessionID
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		consents: make(map[string]*PsychologyConsent),
		sessions: make(map[string]*PsychologySession),
	}
}

func (m *MemoryRepository) SaveConsent(_ context.Context, c *PsychologyConsent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%s:%s", c.StudentID, c.SchoolYear)
	m.consents[key] = c
	return nil
}

func (m *MemoryRepository) GetConsent(_ context.Context, studentID, schoolYear string) (*PsychologyConsent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := fmt.Sprintf("%s:%s", studentID, schoolYear)
	c, ok := m.consents[key]
	if !ok {
		return nil, errors.New("consenso non trovato")
	}
	return c, nil
}

func (m *MemoryRepository) SaveSession(_ context.Context, s *PsychologySession) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[s.ID] = s
	return nil
}

func (m *MemoryRepository) GetSessionByID(_ context.Context, id string) (*PsychologySession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return s, nil
}

func (m *MemoryRepository) ListSessionsByStudent(_ context.Context, studentID string) ([]*PsychologySession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*PsychologySession
	for _, s := range m.sessions {
		if s.StudentID == studentID {
			list = append(list, s)
		}
	}
	return list, nil
}

func (m *MemoryRepository) ListSessionsByPsychologist(_ context.Context, psychologistID string) ([]*PsychologySession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*PsychologySession
	for _, s := range m.sessions {
		if s.PsychologistID == psychologistID {
			list = append(list, s)
		}
	}
	return list, nil
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		repo = NewMemoryRepository()
	}
	return &Service{repo: repo}
}

// SignParentConsent signs informed consent from one parent. Once both parents sign, status becomes ConsentApproved.
func (s *Service) SignParentConsent(ctx context.Context, schoolID, studentID, schoolYear, parentUserID string) (*PsychologyConsent, error) {
	consent, err := s.repo.GetConsent(ctx, studentID, schoolYear)
	now := time.Now()
	if err != nil {
		consent = &PsychologyConsent{
			ID:              fmt.Sprintf("psy-cons-%s-%s", studentID, schoolYear),
			SchoolID:        schoolID,
			StudentID:       studentID,
			SchoolYear:      schoolYear,
			Parent1ID:       parentUserID,
			Parent1SignedAt: &now,
			Status:          ConsentPending,
			CreatedAt:       now,
		}
	} else {
		if consent.Parent1ID == "" || consent.Parent1ID == parentUserID {
			consent.Parent1ID = parentUserID
			consent.Parent1SignedAt = &now
		} else if consent.Parent2ID == "" || consent.Parent2ID == parentUserID {
			consent.Parent2ID = parentUserID
			consent.Parent2SignedAt = &now
		}

		if consent.Parent1SignedAt != nil && consent.Parent2SignedAt != nil {
			consent.Status = ConsentApproved
		}
	}

	if err := s.repo.SaveConsent(ctx, consent); err != nil {
		return nil, err
	}
	return consent, nil
}

// BookSession creates a confidential booking. If the student is a minor, BOTH parental signatures are strictly verified.
func (s *Service) BookSession(ctx context.Context, schoolID, studentID, schoolYear string, isMinor bool, req BookSessionRequest) (*PsychologySession, error) {
	if isMinor {
		consent, err := s.repo.GetConsent(ctx, studentID, schoolYear)
		if err != nil || consent.Status != ConsentApproved {
			return nil, ErrBothParentsConsentRequired
		}
	}

	slotTime, err := time.Parse("2006-01-02 15:04", req.SlotTime)
	if err != nil {
		slotTime, err = time.Parse(time.RFC3339, req.SlotTime)
		if err != nil {
			return nil, fmt.Errorf("formato data orario colloquio non valido: %w", err)
		}
	}

	dur := req.DurationMinutes
	if dur <= 0 {
		dur = 45 // Standard 45 minutes
	}

	// Generate pseudorandom alias for student privacy
	randBytes := make([]byte, 3)
	_, _ = rand.Read(randBytes)
	alias := fmt.Sprintf("ALUNNO-CIC-%s", hex.EncodeToString(randBytes))

	session := &PsychologySession{
		ID:              fmt.Sprintf("sess-%d", time.Now().UnixNano()),
		SchoolID:        schoolID,
		StudentID:       studentID,
		PsychologistID:  req.PsychologistID,
		AnonymousAlias:  alias,
		SlotTime:        slotTime,
		DurationMinutes: dur,
		Status:          "prenotato",
		CreatedAt:       time.Now(),
	}

	if err := s.repo.SaveSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

// GetSessionSecure enforces strict L. 56/1989 professional secrecy.
// Only the designated psychologist can access clinical notes.
func (s *Service) GetSessionSecure(ctx context.Context, sessionID, requestingUserID, requestingRole string) (*PsychologySession, error) {
	sess, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// If the requester is the student, return without clinical notes
	if sess.StudentID == requestingUserID {
		safeCopy := *sess
		safeCopy.EncryptedClinicalNotes = ""
		return &safeCopy, nil
	}

	// Strictly prohibit principals, admins, teachers, secretaries from reading clinical notes
	if requestingRole != "psychologist" || sess.PsychologistID != requestingUserID {
		return nil, ErrProfessionalSecrecyViolation
	}

	return sess, nil
}

func (s *Service) UpdateClinicalNotes(ctx context.Context, sessionID, psychologistID, notes string) (*PsychologySession, error) {
	sess, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if sess.PsychologistID != psychologistID {
		return nil, ErrProfessionalSecrecyViolation
	}

	sess.EncryptedClinicalNotes = notes
	sess.Status = "svolto"
	if err := s.repo.SaveSession(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Service) ListMySessions(ctx context.Context, userID, role string) ([]*PsychologySession, error) {
	if role == "psychologist" {
		return s.repo.ListSessionsByPsychologist(ctx, userID)
	}
	// For students, sanitize notes
	sessions, err := s.repo.ListSessionsByStudent(ctx, userID)
	if err != nil {
		return nil, err
	}
	var sanitized []*PsychologySession
	for _, sess := range sessions {
		safe := *sess
		safe.EncryptedClinicalNotes = ""
		sanitized = append(sanitized, &safe)
	}
	return sanitized, nil
}
