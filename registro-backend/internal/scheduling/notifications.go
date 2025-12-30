package scheduling

import "fmt"

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (n *NotificationService) NotifyBooking(booking *ColloquioBooking, slot *ColloquioSlot, role string) {
	// Mock: Send email
	fmt.Printf("Sending notification for %s to %s\n", booking.ID, role)
}
