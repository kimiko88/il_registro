package notifications

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type PushProvider interface {
	SendPush(ctx context.Context, token PushToken, title, body string, payload map[string]interface{}) error
}

type Service interface {
	RegisterToken(ctx context.Context, userID string, req RegisterTokenRequest) error
	UnregisterToken(ctx context.Context, userID, deviceToken string) error
	SendPushNotification(ctx context.Context, req SendNotificationRequest) (int, error)
	CreateInAppNotification(ctx context.Context, userID, title, body, notifType string, payload map[string]interface{}) (*DBNotification, error)
	ListDBNotifications(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]DBNotification, error)
	MarkAsRead(ctx context.Context, userID, notificationID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	GetPWAManifest() PWAConfig
}

type serviceImpl struct {
	repo         Repository
	pushProvider PushProvider
}

func NewService(repo Repository) Service {
	if repo == nil {
		panic("notifications.NewService: repo must not be nil")
	}
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) SetPushProvider(provider PushProvider) {
	s.pushProvider = provider
}

var validPlatforms = map[string]bool{
	"ios":     true,
	"android": true,
	"web":     true,
}

func (s *serviceImpl) RegisterToken(ctx context.Context, userID string, req RegisterTokenRequest) error {
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

func (s *serviceImpl) UnregisterToken(ctx context.Context, userID, deviceToken string) error {
	return s.repo.DeleteToken(ctx, userID, deviceToken)
}

func (s *serviceImpl) SendPushNotification(ctx context.Context, req SendNotificationRequest) (int, error) {
	tokens, err := s.repo.GetUserTokens(ctx, req.UserID)
	if err != nil {
		return 0, err
	}
	if len(tokens) == 0 {
		return 0, nil
	}

	payload := req.Data
	if payload == nil {
		payload = map[string]interface{}{}
	}

	sentCount := 0
	for _, t := range tokens {
		if s.pushProvider != nil {
			if err := s.pushProvider.SendPush(ctx, t, req.Title, req.Body, payload); err != nil {
				log.Printf("[WARN] notifications.SendPushNotification: push failed to %s (%s): %v", t.DeviceToken, t.Platform, err)
			} else {
				sentCount++
			}
		} else {
			log.Printf("[WARN] notifications.SendPushNotification: no push provider configured — logging dispatch to %s (%s): %q — %q",
				t.DeviceToken, t.Platform, req.Title, req.Body)
		}
	}
	return sentCount, nil
}

func (s *serviceImpl) CreateInAppNotification(ctx context.Context, userID, title, body, notifType string, payload map[string]interface{}) (*DBNotification, error) {
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
	if n.ID == "" {
		n.ID = fmt.Sprintf("notif-%d", time.Now().UnixNano())
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}

	// Async FCM / Push dispatch with timeout context and error logging
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[ERROR] panic recovered in async push dispatch for user %s: %v", userID, r)
			}
		}()
		asyncCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req := SendNotificationRequest{
			UserID: userID,
			Title:  title,
			Body:   body,
			Data:   payload,
		}
		var err error
		for attempt := 1; attempt <= 2; attempt++ {
			select {
			case <-asyncCtx.Done():
				return
			default:
			}
			_, err = s.SendPushNotification(asyncCtx, req)
			if err == nil {
				break
			}
			log.Printf("[WARN] notifications.CreateInAppNotification: attempt %d push dispatch failed for user %s: %v", attempt, userID, err)
			select {
			case <-asyncCtx.Done():
				return
			case <-time.After(200 * time.Millisecond):
			}
		}
	}()

	return n, nil
}

func (s *serviceImpl) ListDBNotifications(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]DBNotification, error) {
	if userID == "" {
		return nil, errors.New("unauthorized")
	}
	return s.repo.ListDBNotifications(ctx, userID, unreadOnly, limit, offset)
}

func (s *serviceImpl) MarkAsRead(ctx context.Context, userID, notificationID string) error {
	if userID == "" {
		return errors.New("unauthorized")
	}
	return s.repo.MarkAsRead(ctx, userID, notificationID)
}

func (s *serviceImpl) MarkAllAsRead(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("unauthorized")
	}
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *serviceImpl) GetPWAManifest() PWAConfig {
	return PWAConfig{
		Name:       "Registro Elettronico Scolastico",
		ShortName:  "Registro",
		StartURL:   "/",
		Display:    "standalone",
		ThemeColor: "#1e3a8a",
		BGColor:    "#ffffff",
	}
}
