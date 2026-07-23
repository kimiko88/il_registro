package auditlog

import (
	"context"
	"fmt"
	"math"
)

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

func (s *service) ExportCSV(ctx context.Context, p FilterParams) ([]byte, error) {
	p.Limit = 5000
	p.Page = 1
	data, _, err := s.repo.List(ctx, p)
	if err != nil {
		return nil, err
	}

	csvStr := "ID,Data/Ora,Utente,Ruolo,Azione,Tipo Entita,ID Entita,IP,Dettagli\n"
	for _, ev := range data {
		csvStr += fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
			ev.ID, ev.CreatedAt.Format("2006-01-02 15:04:05"), ev.ActorName, ev.ActorRole,
			ev.Action, ev.EntityType, ev.EntityID, ev.IPAddress, ev.Details,
		)
	}
	return []byte(csvStr), nil
}
