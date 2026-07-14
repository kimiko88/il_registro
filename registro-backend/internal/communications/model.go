package communications

import (
	"time"
)

type Message struct {
	ID          string     `json:"id"`
	SenderID    string     `json:"sender_id"`
	ReceiverIDs []string   `json:"receiver_ids"` // JSON array in DB
	Subject     string     `json:"subject"`
	Body        string     `json:"body"`
	Type        string     `json:"type"` // 'email', 'notification', 'internal'
	CreatedAt   time.Time  `json:"created_at"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	IsSigned    bool       `json:"is_signed"`
}

type CreateMessageRequest struct {
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	Type       string   `json:"type"`
}
