package interpelli

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	candSeq atomic.Int64

	ErrNoticeNotFound      = errors.New("avviso di interpello non trovato")
	ErrNoticeExpired       = errors.New("il termine perentorio per la presentazione delle domande è scaduto")
	ErrCandidaturaNotFound = errors.New("candidatura non trovata")
	ErrDPR445Mandatory     = errors.New("è obbligatorio accettare la dichiarazione di veridicità ai sensi del D.P.R. 445/2000")
	ErrInvalidGraduation   = errors.New("il voto di laurea/diploma deve essere compreso tra 60 e 110")
)

type Repository interface {
	SaveNotice(ctx context.Context, n *InterpelloNotice) error
	GetNoticeByID(ctx context.Context, id string) (*InterpelloNotice, error)
	ListNotices(ctx context.Context, schoolID, concorsoClass string, onlyActive bool) ([]*InterpelloNotice, error)
	SaveCandidatura(ctx context.Context, c *InterpelloCandidatura) error
	GetCandidaturaByID(ctx context.Context, id string) (*InterpelloCandidatura, error)
	ListCandidatureByNotice(ctx context.Context, noticeID string) ([]*InterpelloCandidatura, error)
}

type MemoryRepository struct {
	mu          sync.RWMutex
	notices     map[string]*InterpelloNotice
	candidature map[string]*InterpelloCandidatura
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		notices:     make(map[string]*InterpelloNotice),
		candidature: make(map[string]*InterpelloCandidatura),
	}
}

func (m *MemoryRepository) SaveNotice(_ context.Context, n *InterpelloNotice) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notices[n.ID] = n
	return nil
}

func (m *MemoryRepository) GetNoticeByID(_ context.Context, id string) (*InterpelloNotice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n, ok := m.notices[id]
	if !ok {
		return nil, ErrNoticeNotFound
	}
	return n, nil
}

func (m *MemoryRepository) ListNotices(_ context.Context, schoolID, concorsoClass string, onlyActive bool) ([]*InterpelloNotice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*InterpelloNotice
	now := time.Now()
	for _, n := range m.notices {
		if schoolID != "" && n.SchoolID != schoolID {
			continue
		}
		if concorsoClass != "" && !strings.EqualFold(n.ConcorsoClass, concorsoClass) {
			continue
		}
		if onlyActive && (n.Status != "aperto" || n.Deadline.Before(now)) {
			continue
		}
		list = append(list, n)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.After(list[j].CreatedAt)
	})
	return list, nil
}

func (m *MemoryRepository) SaveCandidatura(_ context.Context, c *InterpelloCandidatura) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.candidature[c.ID] = c
	return nil
}

func (m *MemoryRepository) GetCandidaturaByID(_ context.Context, id string) (*InterpelloCandidatura, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.candidature[id]
	if !ok {
		return nil, ErrCandidaturaNotFound
	}
	return c, nil
}

func (m *MemoryRepository) ListCandidatureByNotice(_ context.Context, noticeID string) ([]*InterpelloCandidatura, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*InterpelloCandidatura
	for _, c := range m.candidature {
		if c.NoticeID == noticeID {
			list = append(list, c)
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

// CalculateScore computes the score based on O.M. 88/2024 GPS & Interpello criteria.
func CalculateScore(grade float64, lode bool, habilitation bool, monthsService int, langCert string, digitalCount int) (float64, float64, float64, float64) {
	if grade < 60 {
		grade = 60
	}
	if grade > 110 {
		grade = 110
	}

	// 1. Base graduation points: 12 + 0.5 * (grade - 60)
	scoreDegree := 12.0 + 0.5*(grade-60.0)
	if lode {
		scoreDegree += 4.0
	}

	// 2. Habilitation points
	if habilitation {
		scoreDegree += 24.0
	}

	// 3. Service points: 2 points per month, cap 36 points
	if monthsService < 0 {
		monthsService = 0
	}
	scoreService := float64(monthsService) * 2.0
	if scoreService > 36.0 {
		scoreService = 36.0
	}

	// 4. Cultural & digital certifications
	scoreCerts := 0.0
	switch strings.ToUpper(strings.TrimSpace(langCert)) {
	case "B2":
		scoreCerts += 3.0
	case "C1":
		scoreCerts += 4.0
	case "C2":
		scoreCerts += 6.0
	}

	if digitalCount < 0 {
		digitalCount = 0
	}
	digitalScore := float64(digitalCount) * 0.5
	if digitalScore > 2.0 {
		digitalScore = 2.0
	}
	scoreCerts += digitalScore

	total := scoreDegree + scoreService + scoreCerts
	return scoreDegree, scoreService, scoreCerts, total
}

func (s *Service) CreateNotice(ctx context.Context, schoolID string, req CreateNoticeRequest) (*InterpelloNotice, error) {
	if schoolID == "" {
		return nil, errors.New("school_id is required")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format, expected YYYY-MM-DD: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format, expected YYYY-MM-DD: %w", err)
	}

	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		deadline, err = time.Parse("2006-01-02 15:04:05", req.Deadline)
		if err != nil {
			deadline, err = time.Parse("2006-01-02", req.Deadline)
			if err != nil {
				return nil, fmt.Errorf("invalid deadline format: %w", err)
			}
			deadline = deadline.Add(23*time.Hour + 59*time.Minute)
		}
	}

	proto := req.ProtocolNumber
	if proto == "" {
		proto = fmt.Sprintf("INTERPELLO-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)
	}

	notice := &InterpelloNotice{
		ID:             fmt.Sprintf("not-%d", time.Now().UnixNano()),
		SchoolID:       schoolID,
		ProtocolNumber: proto,
		Title:          req.Title,
		ConcorsoClass:  strings.ToUpper(strings.TrimSpace(req.ConcorsoClass)),
		PostType:       req.PostType,
		WeeklyHours:    req.WeeklyHours,
		StartDate:      startDate,
		EndDate:        endDate,
		Deadline:       deadline,
		Status:         "aperto",
		Description:    req.Description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.SaveNotice(ctx, notice); err != nil {
		return nil, err
	}
	return notice, nil
}

func (s *Service) ListPublicNotices(ctx context.Context, schoolID, concorsoClass string) ([]*InterpelloNotice, error) {
	return s.repo.ListNotices(ctx, schoolID, concorsoClass, true)
}

func (s *Service) GetNotice(ctx context.Context, id string) (*InterpelloNotice, error) {
	return s.repo.GetNoticeByID(ctx, id)
}

func (s *Service) SubmitCandidatura(ctx context.Context, noticeID string, req SubmitCandidaturaRequest) (*InterpelloCandidatura, error) {
	notice, err := s.repo.GetNoticeByID(ctx, noticeID)
	if err != nil {
		return nil, err
	}

	if notice.Status != "aperto" || time.Now().After(notice.Deadline) {
		return nil, ErrNoticeExpired
	}

	if !req.DPR445Declared {
		return nil, ErrDPR445Mandatory
	}

	if req.GraduationGrade < 60 || req.GraduationGrade > 110 {
		return nil, ErrInvalidGraduation
	}

	_, scoreService, scoreCerts, total := CalculateScore(
		req.GraduationGrade,
		req.GraduationLode,
		req.HasHabilitation,
		req.MonthsOfService,
		req.CertLanguage,
		req.CertDigitalCount,
	)

	cand := &InterpelloCandidatura{
		ID:               fmt.Sprintf("cand-%d-%d", time.Now().UnixNano(), candSeq.Add(1)),
		NoticeID:         noticeID,
		CandidateName:    strings.TrimSpace(req.CandidateName),
		CandidateSurname: strings.TrimSpace(req.CandidateSurname),
		FiscalCode:       strings.ToUpper(strings.TrimSpace(req.FiscalCode)),
		Email:            strings.TrimSpace(req.Email),
		PEC:              strings.TrimSpace(req.PEC),
		Phone:            strings.TrimSpace(req.Phone),
		GraduationGrade:  req.GraduationGrade,
		GraduationLode:   req.GraduationLode,
		HasHabilitation:  req.HasHabilitation,
		ScoreService:     scoreService,
		ScoreCerts:       scoreCerts,
		TotalScore:       total,
		CVUrl:            req.CVUrl,
		DPR445Declared:   req.DPR445Declared,
		Status:           "valutata",
		CreatedAt:        time.Now(),
	}

	if err := s.repo.SaveCandidatura(ctx, cand); err != nil {
		return nil, err
	}
	return cand, nil
}

func (s *Service) GetGraduatoria(ctx context.Context, noticeID string) ([]*InterpelloCandidatura, error) {
	cands, err := s.repo.ListCandidatureByNotice(ctx, noticeID)
	if err != nil {
		return nil, err
	}

	// Sort according to O.M. 88/2024:
	// 1. Total score descending
	// 2. Tie-breaker: Habilitation precedence
	// 3. Tie-breaker: Service score
	// 4. Tie-breaker: Graduation grade
	sort.Slice(cands, func(i, j int) bool {
		if cands[i].TotalScore != cands[j].TotalScore {
			return cands[i].TotalScore > cands[j].TotalScore
		}
		if cands[i].HasHabilitation != cands[j].HasHabilitation {
			return cands[i].HasHabilitation
		}
		if cands[i].ScoreService != cands[j].ScoreService {
			return cands[i].ScoreService > cands[j].ScoreService
		}
		return cands[i].GraduationGrade > cands[j].GraduationGrade
	})

	return cands, nil
}

func (s *Service) ConvocaCandidato(ctx context.Context, candidaturaID string, hoursToRespond int) (*InterpelloCandidatura, error) {
	cand, err := s.repo.GetCandidaturaByID(ctx, candidaturaID)
	if err != nil {
		return nil, err
	}

	if hoursToRespond <= 0 {
		hoursToRespond = 24 // Standard 24 hours per O.M. 88/2024
	}

	now := time.Now()
	deadline := now.Add(time.Duration(hoursToRespond) * time.Hour)

	cand.Status = "convocata"
	cand.ConvocationSentAt = &now
	cand.ConvocationDeadline = &deadline

	if err := s.repo.SaveCandidatura(ctx, cand); err != nil {
		return nil, err
	}
	return cand, nil
}

func (s *Service) RispondiConvocazione(ctx context.Context, candidaturaID string, risposta string, notes string) (*InterpelloCandidatura, error) {
	cand, err := s.repo.GetCandidaturaByID(ctx, candidaturaID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	rispostaNorm := strings.ToLower(strings.TrimSpace(risposta))
	if rispostaNorm != "accettata" && rispostaNorm != "rifiutata" {
		return nil, fmt.Errorf("risposta non valida: deve essere 'accettata' o 'rifiutata'")
	}

	cand.Status = rispostaNorm
	cand.ResponseAt = &now
	cand.ResponseNotes = notes

	if err := s.repo.SaveCandidatura(ctx, cand); err != nil {
		return nil, err
	}
	return cand, nil
}
