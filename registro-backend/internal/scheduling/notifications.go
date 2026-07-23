package scheduling

import "registro-backend/pkg/logger"

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (n *NotificationService) NotifyBooking(booking *ColloquioBooking, slot *ColloquioSlot, role string) {
	// Mock: Send email
	logger.Log.Infof("Sending notification for %s to %s", booking.ID, role)
}
