package visitors_test

import (
	"context"
	"testing"

	"registro-backend/internal/visitors"
)

type mockVisitorsRepo struct {
	createdVisitor     *visitors.Visitor
	createdEarlyExit   *visitors.EarlyExit
	createdMaintenance *visitors.MaintenanceReport
}

func (m *mockVisitorsRepo) CreateVisitor(ctx context.Context, v *visitors.Visitor) error {
	m.createdVisitor = v
	v.ID = "vis-123"
	return nil
}

func (m *mockVisitorsRepo) RecordVisitorExit(ctx context.Context, id, schoolID, notes string) error {
	return nil
}

func (m *mockVisitorsRepo) ListTodayVisitors(ctx context.Context, schoolID, date string) ([]*visitors.Visitor, error) {
	return []*visitors.Visitor{
		{
			ID:       "vis-1",
			SchoolID: schoolID,
			Name:     "Mario Rossi",
			Purpose:  visitors.PurposeParent,
		},
	}, nil
}

func (m *mockVisitorsRepo) GetVisitor(ctx context.Context, id, schoolID string) (*visitors.Visitor, error) {
	return &visitors.Visitor{
		ID:       id,
		SchoolID: schoolID,
		Name:     "Mario Rossi",
	}, nil
}

func (m *mockVisitorsRepo) CreateEarlyExit(ctx context.Context, e *visitors.EarlyExit) error {
	m.createdEarlyExit = e
	e.ID = "exit-1"
	return nil
}

func (m *mockVisitorsRepo) RecordStudentReturn(ctx context.Context, id, schoolID, notes string) error {
	return nil
}

func (m *mockVisitorsRepo) ListTodayEarlyExits(ctx context.Context, schoolID, date string) ([]*visitors.EarlyExit, error) {
	return []*visitors.EarlyExit{
		{
			ID:            "exit-1",
			SchoolID:      schoolID,
			StudentID:     "student-1",
			DelegateeName: "Anna Verdi",
		},
	}, nil
}

func (m *mockVisitorsRepo) CreateMaintenanceReport(ctx context.Context, r *visitors.MaintenanceReport) error {
	m.createdMaintenance = r
	r.ID = "maint-1"
	return nil
}

func (m *mockVisitorsRepo) ListMaintenanceReports(ctx context.Context, schoolID, statusFilter string) ([]*visitors.MaintenanceReport, error) {
	return []*visitors.MaintenanceReport{
		{
			ID:          "maint-1",
			SchoolID:    schoolID,
			Location:    "Aula 10",
			Category:    "elettrico",
			Status:      visitors.MaintenanceOpen,
			Priority:    visitors.PriorityHigh,
			Description: "Presa non funzionante",
		},
	}, nil
}

func (m *mockVisitorsRepo) UpdateMaintenanceStatus(ctx context.Context, id, schoolID string, req visitors.UpdateMaintenanceStatusRequest) error {
	return nil
}

func TestVisitorsService(t *testing.T) {
	mockRepo := &mockVisitorsRepo{}
	svc := visitors.NewService(mockRepo)
	ctx := context.Background()

	t.Run("RegisterVisitor", func(t *testing.T) {
		req := visitors.RegisterVisitorRequest{
			Name:       "Giuseppe Garibaldi",
			DocumentID: "AB123456C",
			Purpose:    visitors.PurposeSupplier,
			HostName:   "DSGA",
		}
		v, err := svc.RegisterVisitor(ctx, "sch-1", "user-coll", req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.ID != "vis-123" {
			t.Errorf("expected ID vis-123, got %s", v.ID)
		}
		if v.Name != "Giuseppe Garibaldi" {
			t.Errorf("expected Giuseppe Garibaldi, got %s", v.Name)
		}
	})

	t.Run("RecordEarlyExit", func(t *testing.T) {
		req := visitors.RecordEarlyExitRequest{
			StudentID:     "stud-99",
			DelegateeName: "Luigi Bianchi",
			DelegateRel:   "Padre",
			ReasonCode:    "visita_medica",
		}
		exit, err := svc.RecordEarlyExit(ctx, "sch-1", "user-coll", req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exit.StudentID != "stud-99" {
			t.Errorf("expected stud-99, got %s", exit.StudentID)
		}
	})

	t.Run("CreateMaintenanceReport", func(t *testing.T) {
		req := visitors.CreateMaintenanceReportRequest{
			Location:    "Palestra",
			Category:    "strutturale",
			Description: "Vetro rotto porta emergenza",
			Priority:    visitors.PriorityUrgent,
		}
		rep, err := svc.CreateMaintenanceReport(ctx, "sch-1", "user-coll", req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rep.Priority != visitors.PriorityUrgent {
			t.Errorf("expected urgent priority, got %s", rep.Priority)
		}
	})

	t.Run("RolesAllowed", func(t *testing.T) {
		if !visitors.CanAccessVisitorRegistry("collaboratore_scolastico") {
			t.Error("collaboratore_scolastico should have access")
		}
		if !visitors.CanAccessVisitorRegistry("collaboratore_ds") {
			t.Error("collaboratore_ds should have access")
		}
		if visitors.CanAccessVisitorRegistry("student") {
			t.Error("student should NOT have access")
		}
	})
}
