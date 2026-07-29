package notifications

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("notifications.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

var validPlatforms = map[string]bool{
	"ios":     true,
	"android": true,
	"web":     true,
}

func (s *Service) RegisterToken(ctx context.Context, userID string, req RegisterTokenRequest) error {
	if userID == "" {
		return errors.New("unauthorized")
	}
	if req.DeviceToken == "" {
		return errors.New("device_token cannot be empty")
	}
	if !validPlatforms[req.Platform] {
		return fmt.Errorf("invalid platform '%s': allowed platforms are ios, android, web", req.Platform)
	}
	token := &PushToken{
		UserID:      userID,
		DeviceToken: req.DeviceToken,
		Platform:    req.Platform,
	}
	return s.repo.SaveToken(ctx, token)
}

func (s *Service) UnregisterToken(ctx context.Context, userID, deviceToken string) error {
	return s.repo.DeleteToken(ctx, userID, deviceToken)
}

func (s *Service) SendPushNotification(ctx context.Context, req SendNotificationRequest) (int, error) {
	tokens, err := s.repo.GetUserTokens(ctx, req.UserID)
	if err != nil {
		return 0, err
	}
	if len(tokens) == 0 {
		return 0, nil
	}

	// Bug 125: FCM / APNs integration not yet implemented.
	// Log at WARN level so operators know push notifications are stubs.
	sentCount := 0
	for _, t := range tokens {
		log.Printf("[WARN] notifications.SendPushNotification: push stub — would dispatch to %s (%s): %q — %q",
			t.DeviceToken, t.Platform, req.Title, req.Body)
		sentCount++
	}
	return sentCount, nil
}

func (s *Service) CreateInAppNotification(ctx context.Context, userID, title, body, notifType string, payload map[string]interface{}) (*DBNotification, error) {
	n := &DBNotification{
		UserID:  userID,
		Title:   title,
		Body:    body,
		Type:    notifType,
		Payload: payload,
	}
	if err := s.repo.CreateDBNotification(ctx, n); err != nil {
		return nil, err
	}

	// Async FCM / Push dispatch with timeout context and error logging
	go func() {
		asyncCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := s.SendPushNotification(asyncCtx, SendNotificationRequest{
			UserID: userID,
			Title:  title,
			Body:   body,
		}); err != nil {
			log.Printf("[WARN] notifications.CreateInAppNotification: async push dispatch failed for user %s: %v", userID, err)
		}
	}()

	return n, nil
}

func (s *Service) ListDBNotifications(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]DBNotification, error) {
	if userID == "" {
		return nil, errors.New("unauthorized")
	}
	return s.repo.ListDBNotifications(ctx, userID, unreadOnly, limit, offset)
}

func (s *Service) MarkAsRead(ctx context.Context, userID, notificationID string) error {
	if userID == "" {
		return errors.New("unauthorized")
	}
	return s.repo.MarkAsRead(ctx, userID, notificationID)
}

func (s *Service) MarkAllAsRead(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("unauthorized")
	}
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *Service) GetPWAManifest() PWAConfig {
	return PWAConfig{
		Name:       "Registro Elettronico Scolastico",
		ShortName:  "Registro",
		StartURL:   "/",
		Display:    "standalone",
		ThemeColor: "#1e3a8a",
		BGColor:    "#ffffff",
	}
}
