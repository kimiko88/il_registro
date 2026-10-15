package albopretorio

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	ErrActNotFound    = errors.New("atto di albo pretorio non trovato")
	ErrActNotExpired  = errors.New("l'atto non può essere defisso prima del decorso dei 15 giorni prescritti ex L. 69/2009")
	ErrInvalidSubject = errors.New("l'oggetto dell'atto è obbligatorio")
)

type Repository interface {
	Save(ctx context.Context, item *AlboItem) error
	GetByID(ctx context.Context, id string) (*AlboItem, error)
	List(ctx context.Context, schoolID string) ([]*AlboItem, error)
	GetNextRepertoryNumber(ctx context.Context, schoolID string, year int) (int, error)
}

type MemoryRepository struct {
	mu        sync.RWMutex
	items     map[string]*AlboItem
	repertory map[string]int // key: schoolID:year
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items:     make(map[string]*AlboItem),
		repertory: make(map[string]int),
	}
}

func (m *MemoryRepository) Save(_ context.Context, item *AlboItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[item.ID] = item
	return nil
}

func (m *MemoryRepository) GetByID(_ context.Context, id string) (*AlboItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	item, ok := m.items[id]
	if !ok {
		return nil, ErrActNotFound
	}
	return item, nil
}

func (m *MemoryRepository) List(_ context.Context, schoolID string) ([]*AlboItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*AlboItem
	for _, item := range m.items {
		if schoolID != "" && item.SchoolID != schoolID {
			continue
		}
		list = append(list, item)
	}
	return list, nil
}

func (m *MemoryRepository) GetNextRepertoryNumber(_ context.Context, schoolID string, year int) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%s:%d", schoolID, year)
	m.repertory[key]++
	return m.repertory[key], nil
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

func (s *Service) PublishAct(ctx context.Context, schoolID, publishedBy string, req PublishActRequest) (*AlboItem, error) {
	if schoolID == "" {
		return nil, errors.New("school_id is required")
	}
	subj := strings.TrimSpace(req.Subject)
	if subj == "" {
		return nil, ErrInvalidSubject
	}

	year := time.Now().Year()
	nextNum, err := s.repo.GetNextRepertoryNumber(ctx, schoolID, year)
	if err != nil {
		return nil, err
	}

	days := req.DaysDuration
	if days <= 0 {
		days = 15 // Mandated by Legge 69/2009
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(days) * 24 * time.Hour)
	repertoryCode := fmt.Sprintf("%d/%05d", year, nextNum)

	item := &AlboItem{
		ID:                      fmt.Sprintf("albo-%d-%05d", year, nextNum),
		SchoolID:                schoolID,
		RepertoryYear:           year,
		RepertoryNumber:         nextNum,
		RepertoryCode:           repertoryCode,
		Category:                req.Category,
		Subject:                 subj,
		PublishedAt:             now,
		ExpiresAt:               expiresAt,
		Status:                  "in_pubblicazione",
		DocumentFileURL:         req.DocumentFileURL,
		DocumentSHA256:          req.DocumentSHA256,
		PublishedBy:             publishedBy,
		IsTransparencySection:   req.IsTransparencySection,
		TransparencyMacroFamily: req.TransparencyMacroFamily,
		TransparencySubFamily:   req.TransparencySubFamily,
		CIGCode:                 req.CIGCode,
		AwardedAmount:           req.AwardedAmount,
		CreatedAt:               now,
		UpdatedAt:               now,
	}

	if err := s.repo.Save(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) ListPublicAlbo(ctx context.Context, schoolID, category, query string, year int, archived bool) ([]*AlboItem, error) {
	items, err := s.repo.List(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var filtered []*AlboItem
	for _, item := range items {
		// Auto-defission check: if expiresAt is past, mark defisso for public historical archive
		if item.Status == "in_pubblicazione" && now.After(item.ExpiresAt) {
			item.Status = "defisso"
			_ = s.repo.Save(ctx, item)
		}

		if archived {
			if item.Status != "defisso" {
				continue
			}
		} else {
			if item.Status != "in_pubblicazione" {
				continue
			}
		}

		if category != "" && !strings.EqualFold(item.Category, category) {
			continue
		}
		if year > 0 && item.RepertoryYear != year {
			continue
		}
		if query != "" {
			q := strings.ToLower(query)
			if !strings.Contains(strings.ToLower(item.Subject), q) && !strings.Contains(strings.ToLower(item.RepertoryCode), q) {
				continue
			}
		}
		filtered = append(filtered, item)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].RepertoryNumber > filtered[j].RepertoryNumber
	})

	return filtered, nil
}

func (s *Service) DefiggiAtto(ctx context.Context, id, dsName string, force bool) (*AlboItem, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if !force && now.Before(item.ExpiresAt) {
		return nil, ErrActNotExpired
	}

	item.Status = "defisso"
	item.UpdatedAt = now

	// Generate and seal the Relata di Pubblicazione
	relata := fmt.Sprintf(
		"RELATA DI PUBBLICAZIONE (L. 69/2009)\nSi attesta che il presente atto Reg. %s è stato pubblicato all'Albo Pretorio Informatico dal %s al %s per 15 giorni consecutivi, senza opposizioni.\nImpronta informatica documento (SHA-256): %s\nIl Dirigente Scolastico: %s",
		item.RepertoryCode,
		item.PublishedAt.Format("02/01/2006"),
		now.Format("02/01/2006"),
		item.DocumentSHA256,
		dsName,
	)
	item.RelataText = relata
	item.RelataSignedBy = dsName
	item.RelataSignedAt = &now

	if err := s.repo.Save(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GeneraCertificato(ctx context.Context, id, dsName string) (*CertificatoPubblicazione, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	solarDays := int(item.ExpiresAt.Sub(item.PublishedAt).Hours() / 24)
	attestation := fmt.Sprintf(
		"CERTIFICATO DI AVVENUTA PUBBLICAZIONE A NORMA DELL'ART. 32 L. 69/2009\nIl sottoscritto %s, in qualità di Dirigente Scolastico, certifica che il documento '%s' (Repertorio N. %s) è stato ritualmente affisso all'Albo Pretorio telematico per il prescritto periodo di legge (%d giorni solari) garantendo la piena pubblicità legale e l'integrità del documento telematico conforme al CAD.",
		dsName,
		item.Subject,
		item.RepertoryCode,
		solarDays,
	)

	return &CertificatoPubblicazione{
		RepertoryCode:    item.RepertoryCode,
		Subject:          item.Subject,
		Category:         item.Category,
		PublishedAt:      item.PublishedAt,
		DefissoAt:        item.ExpiresAt,
		SolarDays:        solarDays,
		DocumentSHA256:   item.DocumentSHA256,
		LegalAttestation: attestation,
		DirigenteName:    dsName,
		CertifiedAt:      time.Now(),
	}, nil
}

func (s *Service) GenerateANACXML(ctx context.Context, schoolID string, anno int) (string, error) {
	items, err := s.repo.List(ctx, schoolID)
	if err != nil {
		return "", err
	}

	dataset := ANACDataset{
		XMLNS: "legge190_1_0",
		Metadata: ANACMetadata{
			Titolo:            "Adempimenti L. 190/2012 e D.Lgs. 33/2013 - Bandi di gara e contratti",
			CodiceFiscale:     "97123456789",
			AnnoRiferimento:   anno,
			DataPubblicazione: time.Now().Format("2006-01-02"),
		},
		Lotti: make([]ANACLotto, 0),
	}

	for _, item := range items {
		if !item.IsTransparencySection || item.CIGCode == "" {
			continue
		}
		if anno > 0 && item.RepertoryYear != anno {
			continue
		}

		dataset.Lotti = append(dataset.Lotti, ANACLotto{
			CIG:                      item.CIGCode,
			StrutturaProponente:      "Istituto Scolastico Statale",
			Oggetto:                  item.Subject,
			SceltaContraente:         "AFFIDAMENTO DIRETTO EX ART. 50 D.LGS. 36/2023",
			ImportoAggiudicazione:    item.AwardedAmount,
			TempiCompletamentoInizio: item.PublishedAt.Format("2006-01-02"),
			TempiCompletamentoFine:   item.ExpiresAt.Format("2006-01-02"),
		})
	}

	output, err := xml.MarshalIndent(dataset, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(output), nil
}
