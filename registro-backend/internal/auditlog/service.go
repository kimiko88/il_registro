package auditlog

import (
	"context"
	"fmt"
	"math"
	"strings"
)

// truncationLimit is the maximum number of audit log entries returned by ExportCSV.
// Entries beyond this limit are silently dropped; a WARN is logged when this occurs.
const truncationLimit = 5000

type Service interface {
	Log(ctx context.Context, event AuditEvent)
	List(ctx context.Context, p FilterParams) (*PaginatedAuditLogs, error)
	ExportCSV(ctx context.Context, p FilterParams) ([]byte, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Log(ctx context.Context, event AuditEvent) {
	s.repo.InsertAsync(event)
}

func (s *service) List(ctx context.Context, p FilterParams) (*PaginatedAuditLogs, error) {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Page <= 0 {
		p.Page = 1
	}

	data, total, err := s.repo.List(ctx, p)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []AuditEvent{}
	}

	totalPages := int(math.Ceil(float64(total) / float64(p.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginatedAuditLogs{
		Data:       data,
		Total:      total,
		Page:       p.Page,
		Limit:      p.Limit,
		TotalPages: totalPages,
	}, nil
}

// sanitizeCSVField prevents CSV injection (formula injection) by prepending a tab
// to values that start with formula-triggering characters (=, +, -, @).
// See: https://owasp.org/www-community/attacks/CSV_Injection
func sanitizeCSVField(v string) string {
	if len(v) > 0 {
		switch v[0] {
		case '=', '+', '-', '@':
			return "\t" + v
		}
	}
	return v
}

func (s *service) ExportCSV(ctx context.Context, p FilterParams) ([]byte, error) {
	p.Limit = truncationLimit
	p.Page = 1
	data, total, err := s.repo.List(ctx, p)
	if err != nil {
		return nil, err
	}

	var sb strings.Builder
	if total > truncationLimit {
		sb.WriteString(fmt.Sprintf("# WARNING: Export truncated to %d records (total matching: %d). Narrow filter criteria for more records.\n", truncationLimit, total))
	}
	sb.WriteString("ID,Data/Ora,Utente,Ruolo,Azione,Tipo Entita,ID Entita,IP,Dettagli\n")
	for _, ev := range data {
		// Bug 121: sanitize all CSV fields to prevent formula injection
		sb.WriteString(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
			sanitizeCSVField(ev.ID),
			sanitizeCSVField(ev.CreatedAt.Format("2006-01-02 15:04:05")),
			sanitizeCSVField(ev.ActorName),
			sanitizeCSVField(ev.ActorRole),
			sanitizeCSVField(ev.Action),
			sanitizeCSVField(ev.EntityType),
			sanitizeCSVField(ev.EntityID),
			sanitizeCSVField(ev.IPAddress),
			sanitizeCSVField(ev.Details),
		))
	}
	return []byte(sb.String()), nil
}
