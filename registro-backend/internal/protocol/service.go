package protocol

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ProtocolDocument(ctx context.Context, entry *ProtocolEntry, entityType, entityID, schoolName string) (*ProtocolEntry, []byte, error) {
	if entry.SchoolID == "" || entry.Subject == "" || entry.Sender == "" || entry.Recipient == "" {
		return nil, nil, errors.New("school_id, oggetto, mittente e destinatario sono obbligatori")
	}
	if entry.ClassificationTitle < 1 || entry.ClassificationTitle > 10 {
		entry.ClassificationTitle = 7 // Default Titolo VII - Alunni
	}
	if entry.FlowDirection != "in" && entry.FlowDirection != "out" && entry.FlowDirection != "internal" {
		entry.FlowDirection = "in"
	}
	if entry.DocumentHashSHA256 == "" {
		entry.DocumentHashSHA256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	}

	if err := s.repo.RegisterDocument(ctx, entry, entityType, entityID); err != nil {
		return nil, nil, err
	}

	if schoolName == "" {
		schoolName = "ISTITUTO SCOLASTICO STATALE"
	}
	entry.VisualStamp = FormatVisualStamp(*entry, schoolName)

	xmlBytes, err := GenerateSegnaturaXML(*entry)
	if err != nil {
		return nil, nil, err
	}

	return entry, xmlBytes, nil
}

func (s *Service) ListEntries(ctx context.Context, schoolID string, year int, flowDirection string) ([]ProtocolEntry, error) {
	if schoolID == "" {
		return nil, errors.New("school_id obbligatorio")
	}
	if year == 0 {
		year = time.Now().Year()
	}
	return s.repo.ListProtocolEntries(ctx, schoolID, year, flowDirection)
}

func (s *Service) GetEntry(ctx context.Context, id string) (*ProtocolEntry, error) {
	if id == "" {
		return nil, errors.New("id obbligatorio")
	}
	return s.repo.GetProtocolEntry(ctx, id)
}

func (s *Service) GetEntityProtocol(ctx context.Context, entityType, entityID string) (*ProtocolEntry, error) {
	if entityType == "" || entityID == "" {
		return nil, errors.New("entity_type ed entity_id obbligatori")
	}
	return s.repo.GetEntityProtocol(ctx, entityType, entityID)
}
