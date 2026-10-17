package maturita

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

var (
	ErrInvalidGradeYear = errors.New("l'anno di corso per i crediti deve essere 3, 4 o 5")
	ErrInvalidAverage   = errors.New("la media dei voti deve essere compresa tra 6.0 e 10.0")
	ErrExamScoreRange   = errors.New("il punteggio di ciascuna prova d'esame deve essere compreso tra 0 e 20")
	ErrBonusIneligible  = errors.New("il candidato non possiede i requisiti per il bonus d'esame (richiesti almeno 30 crediti di triennio e 50 punti nelle prove)")
	ErrLodeIneligible   = errors.New("la lode può essere conferita solo con 40 crediti di triennio e 60 punti nelle prove d'esame senza bonus")
	ErrRecordNotFound   = errors.New("record di maturità non trovato")
)

// CalculateYearCredits calculates credit points based on D.Lgs. 62/2017 Allegato A.
func CalculateYearCredits(yearGrade int, average float64, highBand bool) (int, error) {
	if yearGrade < 3 || yearGrade > 5 {
		return 0, ErrInvalidGradeYear
	}
	if average < 6.0 || average > 10.0 {
		return 0, ErrInvalidAverage
	}

	switch yearGrade {
	case 3: // Class 3ª (Max 12)
		switch {
		case average == 6.0:
			if highBand {
				return 8, nil
			}
			return 7, nil
		case average <= 7.0:
			if highBand {
				return 9, nil
			}
			return 8, nil
		case average <= 8.0:
			if highBand {
				return 10, nil
			}
			return 9, nil
		case average <= 9.0:
			if highBand {
				return 11, nil
			}
			return 10, nil
		default: // > 9.0
			if highBand {
				return 12, nil
			}
			return 11, nil
		}

	case 4: // Class 4ª (Max 13)
		switch {
		case average == 6.0:
			if highBand {
				return 9, nil
			}
			return 8, nil
		case average <= 7.0:
			if highBand {
				return 10, nil
			}
			return 9, nil
		case average <= 8.0:
			if highBand {
				return 11, nil
			}
			return 10, nil
		case average <= 9.0:
			if highBand {
				return 12, nil
			}
			return 11, nil
		default: // > 9.0
			if highBand {
				return 13, nil
			}
			return 12, nil
		}

	case 5: // Class 5ª (Max 15)
		switch {
		case average == 6.0:
			if highBand {
				return 10, nil
			}
			return 9, nil
		case average <= 7.0:
			if highBand {
				return 11, nil
			}
			return 10, nil
		case average <= 8.0:
			if highBand {
				return 12, nil
			}
			return 11, nil
		case average <= 9.0:
			if highBand {
				return 13, nil
			}
			return 12, nil
		default: // > 9.0
			if highBand {
				return 15, nil
			}
			return 14, nil
		}
	}

	return 0, nil
}

// CalculateTrienniumCredits aggregates 3rd, 4th and 5th year credits into total credits (Max 40).
func CalculateTrienniumCredits(c3, c4, c5 int) (int, error) {
	if c3 < 7 || c3 > 12 {
		return 0, fmt.Errorf("crediti 3° anno fuori intervallo consentito (7-12): %d", c3)
	}
	if c4 < 8 || c4 > 13 {
		return 0, fmt.Errorf("crediti 4° anno fuori intervallo consentito (8-13): %d", c4)
	}
	if c5 < 9 || c5 > 15 {
		return 0, fmt.Errorf("crediti 5° anno fuori intervallo consentito (9-15): %d", c5)
	}
	tot := c3 + c4 + c5
	if tot > 40 {
		tot = 40
	}
	return tot, nil
}

// CalculateFinalExamScores computes the overall score and verifies bonus & lode rules.
func CalculateFinalExamScores(credits TrienniumCredits, written1, written2, oral float64, bonus int, wantLode bool) (*ExamScores, error) {
	if written1 < 0 || written1 > 20 || written2 < 0 || written2 > 20 || oral < 0 || oral > 20 {
		return nil, ErrExamScoreRange
	}

	examTotal := written1 + written2 + oral
	bonusEligible := credits.TotalCredits >= 30 && examTotal >= 50.0

	actualBonus := 0
	if bonus > 0 {
		if !bonusEligible {
			return nil, ErrBonusIneligible
		}
		if bonus > 5 {
			bonus = 5
		}
		actualBonus = bonus
	}

	finalSum := float64(credits.TotalCredits) + examTotal + float64(actualBonus)
	finalScore := int(math.Round(finalSum))
	if finalScore > 100 {
		finalScore = 100
	}

	lode := false
	if wantLode {
		if credits.TotalCredits != 40 || examTotal < 59.99 || actualBonus > 0 {
			return nil, ErrLodeIneligible
		}
		lode = true
	}

	return &ExamScores{
		Written1Score: written1,
		Written2Score: written2,
		OralScore:     oral,
		ExamTotal:     examTotal,
		BonusEligible: bonusEligible,
		BonusPoints:   actualBonus,
		FinalScore:    finalScore,
		Lode:          lode,
	}, nil
}

type Repository interface {
	SaveCommission(ctx context.Context, c *CommissioneMaturita) error
	GetCommissionByClass(ctx context.Context, classID, schoolYear string) (*CommissioneMaturita, error)
	SaveRecord(ctx context.Context, r *StudentMaturitaRecord) error
	GetRecord(ctx context.Context, studentID, schoolYear string) (*StudentMaturitaRecord, error)
	ListRecordsByClass(ctx context.Context, classID, schoolYear string) ([]*StudentMaturitaRecord, error)
}

type MemoryRepository struct {
	mu          sync.RWMutex
	commissions map[string]*CommissioneMaturita   // key: classID:year
	records     map[string]*StudentMaturitaRecord // key: studentID:year
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		commissions: make(map[string]*CommissioneMaturita),
		records:     make(map[string]*StudentMaturitaRecord),
	}
}

func (m *MemoryRepository) SaveCommission(_ context.Context, c *CommissioneMaturita) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%s:%s", c.ClassID, c.SchoolYear)
	m.commissions[key] = c
	return nil
}

func (m *MemoryRepository) GetCommissionByClass(_ context.Context, classID, schoolYear string) (*CommissioneMaturita, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := fmt.Sprintf("%s:%s", classID, schoolYear)
	c, ok := m.commissions[key]
	if !ok {
		return nil, errors.New("commissione non trovata")
	}
	return c, nil
}

func (m *MemoryRepository) SaveRecord(_ context.Context, r *StudentMaturitaRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%s:%s", r.StudentID, r.SchoolYear)
	m.records[key] = r
	return nil
}

func (m *MemoryRepository) GetRecord(_ context.Context, studentID, schoolYear string) (*StudentMaturitaRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := fmt.Sprintf("%s:%s", studentID, schoolYear)
	r, ok := m.records[key]
	if !ok {
		return nil, ErrRecordNotFound
	}
	return r, nil
}

func (m *MemoryRepository) ListRecordsByClass(_ context.Context, classID, schoolYear string) ([]*StudentMaturitaRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*StudentMaturitaRecord
	for _, r := range m.records {
		if r.ClassID == classID && r.SchoolYear == schoolYear {
			list = append(list, r)
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

func (s *Service) SaveCommission(ctx context.Context, c *CommissioneMaturita) error {
	return s.repo.SaveCommission(ctx, c)
}

func (s *Service) GetCommission(ctx context.Context, classID, schoolYear string) (*CommissioneMaturita, error) {
	return s.repo.GetCommissionByClass(ctx, classID, schoolYear)
}

func (s *Service) SaveStudentRecord(ctx context.Context, r *StudentMaturitaRecord) error {
	r.UpdatedAt = time.Now()
	return s.repo.SaveRecord(ctx, r)
}

func (s *Service) GetStudentRecord(ctx context.Context, studentID, schoolYear string) (*StudentMaturitaRecord, error) {
	return s.repo.GetRecord(ctx, studentID, schoolYear)
}

func (s *Service) GetTabelloneClasse(ctx context.Context, classID, schoolYear string) ([]*StudentMaturitaRecord, error) {
	return s.repo.ListRecordsByClass(ctx, classID, schoolYear)
}

func (s *Service) ExportCurriculumXML(record *StudentMaturitaRecord) (string, error) {
	if record == nil {
		return "", errors.New("record cannot be nil")
	}

	cf := "RSSMRA06A01H501Z"
	name := record.StudentName
	pctoHours := 150
	if record.CurriculumStudente != nil {
		if record.CurriculumStudente.StudentCF != "" {
			cf = record.CurriculumStudente.StudentCF
		}
		if record.CurriculumStudente.StudentFullName != "" {
			name = record.CurriculumStudente.StudentFullName
		}
		if record.CurriculumStudente.PCTOHoursTotal > 0 {
			pctoHours = record.CurriculumStudente.PCTOHoursTotal
		}
	}

	xmlStruct := MinisterialCurriculumXML{
		CodiceFiscale:   cf,
		Nominativo:      name,
		EsitoMaturita:   record.Scores.FinalScore,
		LodeConferita:   record.Scores.Lode,
		CreditiTriennio: record.Credits.TotalCredits,
		OrePCTO:         pctoHours,
	}

	out, err := xml.MarshalIndent(xmlStruct, "", "  ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(out), nil
}
