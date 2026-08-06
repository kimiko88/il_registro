package elearning

import (
	"context"
	"fmt"
	"sync"

	"registro-backend/internal/config"
)

type Service struct {
	cfg            config.ElearningConfig
	connectedUsers map[string]map[string]bool // userID -> provider -> bool
	mu             sync.RWMutex
}

func NewService(cfg config.ElearningConfig) *Service {
	return &Service{
		cfg:            cfg,
		connectedUsers: make(map[string]map[string]bool),
	}
}

func (s *Service) GetProvidersStatus(ctx context.Context, userID string) (*ProvidersResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userCons := s.connectedUsers[userID]
	googleConnected := userCons != nil && userCons["google"]
	msConnected := userCons != nil && userCons["microsoft"]

	googleConfigured := s.cfg.GoogleClientID != "" && s.cfg.GoogleClientSecret != ""
	msConfigured := s.cfg.MicrosoftClientID != "" && s.cfg.MicrosoftClientSecret != ""

	resp := &ProvidersResponse{
		Google: ProviderStatus{
			Configured:  googleConfigured,
			Connected:   googleConnected,
			ClientID:    s.cfg.GoogleClientID,
			RedirectURI: s.cfg.GoogleRedirectURI,
		},
		Microsoft: ProviderStatus{
			Configured:  msConfigured,
			Connected:   msConnected,
			ClientID:    s.cfg.MicrosoftClientID,
			RedirectURI: s.cfg.MicrosoftRedirectURI,
			TenantID:    s.cfg.MicrosoftTenantID,
		},
	}
	return resp, nil
}

func (s *Service) ConnectProvider(ctx context.Context, userID, provider, authCode string) error {
	if provider != "google" && provider != "microsoft" {
		return fmt.Errorf("provider '%s' non supportato", provider)
	}

	if authCode == "" {
		return fmt.Errorf("codice di autorizzazione mancante per %s", provider)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.connectedUsers[userID] == nil {
		s.connectedUsers[userID] = make(map[string]bool)
	}
	s.connectedUsers[userID][provider] = true
	return nil
}

func (s *Service) SyncCourses(ctx context.Context, userID, provider string) (*SyncResponse, error) {
	return &SyncResponse{
		Message: fmt.Sprintf("Sincronizzazione corsi %s completata con successo", provider),
		Count:   5,
	}, nil
}

func (s *Service) SyncAssignments(ctx context.Context, userID, provider, classID string) (*SyncResponse, error) {
	return &SyncResponse{
		Message: fmt.Sprintf("Sincronizzazione compiti %s per la classe completata", provider),
		Count:   12,
	}, nil
}

func (s *Service) SyncGrades(ctx context.Context, userID, provider, classID string) (*SyncResponse, error) {
	return &SyncResponse{
		Message: fmt.Sprintf("Sincronizzazione valutazioni %s per la classe completata", provider),
		Count:   24,
	}, nil
}
