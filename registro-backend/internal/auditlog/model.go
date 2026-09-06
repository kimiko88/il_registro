package auditlog

import (
	"time"
)

type AuditEvent struct {
	ID         string    `json:"id" db:"id"`
	SchoolID   string    `json:"school_id" db:"school_id"`
	ActorID    string    `json:"actor_id" db:"actor_id"`
	ActorRole  string    `json:"actor_role" db:"actor_role"`
	ActorName  string    `json:"actor_name" db:"actor_name"`
	Action     string    `json:"action" db:"action"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   string    `json:"entity_id" db:"entity_id"`
	Details    string    `json:"details" db:"details"`
	IPAddress  string    `json:"ip_address" db:"ip_address"`
	UserAgent  string    `json:"user_agent" db:"user_agent"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type FilterParams struct {
	SchoolID   string
	ActorID    string
	Action     string
	EntityType string
	IPAddress  string
	From       string
	To         string
	Page       int
	Limit      int
}

type PaginatedAuditLogs struct {
	Data       []AuditEvent `json:"data"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	TotalPages int          `json:"total_pages"`
}
