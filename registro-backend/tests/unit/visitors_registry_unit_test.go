package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"registro-backend/internal/visitors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockVisitorsRepo struct {
	visitors.Repository
	visitorsList     []*visitors.Visitor
	earlyExitsList   []*visitors.EarlyExit
	maintenanceList  []*visitors.MaintenanceReport
	createVisitorErr error
	recordExitErr    error
	createExitErr    error
	recordReturnErr  error
	createMaintErr   error
	updateMaintErr   error
}

func (m *mockVisitorsRepo) CreateVisitor(ctx context.Context, v *visitors.Visitor) error {
	if m.createVisitorErr != nil {
		return m.createVisitorErr
	}
	if v.ID == "" {
		v.ID = "v-test-1"
	}
	v.EntryTime = time.Now()
	v.CreatedAt = time.Now()
	m.visitorsList = append(m.visitorsList, v)
	return nil
}

func (m *mockVisitorsRepo) RecordVisitorExit(ctx context.Context, id, schoolID string, notes string) error {
	if m.recordExitErr != nil {
		return m.recordExitErr
	}
	now := time.Now()
	for _, v := range m.visitorsList {
		if v.ID == id {
			v.ExitTime = &now
			if notes != "" {
				v.Notes = notes
			}
			return nil
		}
	}
	return errors.New("visitatore non trovato")
}

func (m *mockVisitorsRepo) ListTodayVisitors(ctx context.Context, schoolID, date string) ([]*visitors.Visitor, error) {
	return m.visitorsList, nil
}

func (m *mockVisitorsRepo) CreateEarlyExit(ctx context.Context, e *visitors.EarlyExit) error {
	if m.createExitErr != nil {
		return m.createExitErr
	}
	if e.ID == "" {
		e.ID = "exit-test-1"
	}
	e.ExitTime = time.Now()
	e.CreatedAt = time.Now()
	m.earlyExitsList = append(m.earlyExitsList, e)
	return nil
}

func (m *mockVisitorsRepo) RecordStudentReturn(ctx context.Context, id, schoolID, notes string) error {
	if m.recordReturnErr != nil {
		return m.recordReturnErr
	}
	now := time.Now()
	for _, e := range m.earlyExitsList {
		if e.ID == id {
			e.ReturnTime = &now
			if notes != "" {
				e.Notes = notes
			}
			return nil
		}
	}
	return errors.New("uscita non trovata")
}

func (m *mockVisitorsRepo) ListTodayEarlyExits(ctx context.Context, schoolID, date string) ([]*visitors.EarlyExit, error) {
	return m.earlyExitsList, nil
}

func (m *mockVisitorsRepo) CreateMaintenanceReport(ctx context.Context, r *visitors.MaintenanceReport) error {
	if m.createMaintErr != nil {
		return m.createMaintErr
	}
	if r.ID == "" {
		r.ID = "maint-test-1"
	}
	r.Status = visitors.MaintenanceOpen
	r.CreatedAt = time.Now()
	m.maintenanceList = append(m.maintenanceList, r)
	return nil
}

func (m *mockVisitorsRepo) ListMaintenanceReports(ctx context.Context, schoolID string, statusFilter string) ([]*visitors.MaintenanceReport, error) {
	if statusFilter == "" {
		return m.maintenanceList, nil
	}
	var filtered []*visitors.MaintenanceReport
	for _, r := range m.maintenanceList {
		if string(r.Status) == statusFilter {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}

func (m *mockVisitorsRepo) UpdateMaintenanceStatus(ctx context.Context, id, schoolID string, req visitors.UpdateMaintenanceStatusRequest) error {
	if m.updateMaintErr != nil {
		return m.updateMaintErr
	}
	for _, r := range m.maintenanceList {
		if r.ID == id {
			r.Status = req.Status
			if req.AssignedTo != "" {
				assigned := req.AssignedTo
				r.AssignedTo = &assigned
			}
			return nil
		}
	}
	return errors.New("segnalazione non trovata")
}

func TestVisitors_RegisterVisitor_Unit(t *testing.T) {
	ctx := context.Background()
	schoolID := "school-unit-1"
	recordedBy := "user-ata-1"

	t.Run("rejects visitor with empty name", func(t *testing.T) {
		repo := &mockVisitorsRepo{}
		svc := visitors.NewService(repo)

		req := visitors.RegisterVisitorRequest{
			Name:    "",
			Purpose: visitors.PurposeParent,
		}

		v, err := svc.RegisterVisitor(ctx, schoolID, recordedBy, req)
		require.Error(t, err)
		assert.Nil(t, v)
		assert.Contains(t, err.Error(), "obbligatorio")
	})

	t.Run("registers visitor with badge and purpose correctly", func(t *testing.T) {
		repo := &mockVisitorsRepo{}
		svc := visitors.NewService(repo)

		req := visitors.RegisterVisitorRequest{
			Name:        "Marco Rossi",
			DocumentID:  "CI-998877",
			Purpose:     visitors.PurposeSupplier,
			HostName:    "DSGA",
			BadgeNumber: "B-12",
			Notes:       "Controllo impianto riscaldamento",
		}

		v, err := svc.RegisterVisitor(ctx, schoolID, recordedBy, req)
		require.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, "Marco Rossi", v.Name)
		assert.Equal(t, visitors.PurposeSupplier, v.Purpose)
		require.NotNil(t, v.BadgeNumber)
		assert.Equal(t, "B-12", *v.BadgeNumber)
		assert.Equal(t, schoolID, v.SchoolID)
		assert.Equal(t, recordedBy, v.RecordedBy)
	})

	t.Run("handles empty badge as nil pointer", func(t *testing.T) {
		repo := &mockVisitorsRepo{}
		svc := visitors.NewService(repo)

		req := visitors.RegisterVisitorRequest{
			Name:        "Giulia Neri",
			Purpose:     visitors.PurposeInstitution,
			BadgeNumber: "",
		}

		v, err := svc.RegisterVisitor(ctx, schoolID, recordedBy, req)
		require.NoError(t, err)
		assert.Nil(t, v.BadgeNumber)
	})

	t.Run("records visitor exit with timestamp", func(t *testing.T) {
		repo := &mockVisitorsRepo{}
		svc := visitors.NewService(repo)

		v, _ := svc.RegisterVisitor(ctx, schoolID, recordedBy, visitors.RegisterVisitorRequest{
			Name: "Luca Verdi",
		})
		require.NotNil(t, v)
		assert.Nil(t, v.ExitTime)

		err := svc.RecordVisitorExit(ctx, v.ID, schoolID, "Badge restituito")
		require.NoError(t, err)
		assert.NotNil(t, v.ExitTime)
		assert.Equal(t, "Badge restituito", v.Notes)
	})
}

func TestVisitors_EarlyExitsAndReturns_Unit(t *testing.T) {
	ctx := context.Background()
	schoolID := "school-unit-1"
	recordedBy := "user-ata-1"

	t.Run("records early exit and subsequent student return", func(t *testing.T) {
		repo := &mockVisitorsRepo{}
		svc := visitors.NewService(repo)

		req := visitors.RecordEarlyExitRequest{
			StudentID:     "student-42",
			DelegateeName: "Giovanni Bianchi",
			DelegateRel:   "Padre",
			ReasonCode:    "visita_medica",
			Notes:         "Uscita autorizzata con delega scritta",
		}

		exit, err := svc.RecordEarlyExit(ctx, schoolID, recordedBy, req)
		require.NoError(t, err)
		require.NotNil(t, exit)
		assert.Equal(t, "student-42", exit.StudentID)
		assert.Equal(t, "Giovanni Bianchi", exit.DelegateeName)
		assert.Equal(t, "Padre", exit.DelegateRel)
		assert.Nil(t, exit.ReturnTime)

		// Record student return
		err = svc.RecordStudentReturn(ctx, exit.ID, schoolID, "Rientrato con certificato")
		require.NoError(t, err)
		assert.NotNil(t, exit.ReturnTime)
		assert.Equal(t, "Rientrato con certificato", exit.Notes)
	})
}

func TestVisitors_MaintenanceReports_Lifecycle_Unit(t *testing.T) {
	ctx := context.Background()
	schoolID := "school-unit-1"
	reportedBy := "user-cs-1"

	t.Run("creates maintenance report with open status and updates to in_lavorazione", func(t *testing.T) {
		repo := &mockVisitorsRepo{}
		svc := visitors.NewService(repo)

		req := visitors.CreateMaintenanceReportRequest{
			Location:    "Plesso Centrale - Laboratorio Informatica 2",
			Category:    "informatica",
			Description: "Proiettore non si accende, spia rossa lampeggiante",
			Priority:    visitors.PriorityHigh,
		}

		report, err := svc.CreateMaintenanceReport(ctx, schoolID, reportedBy, req)
		require.NoError(t, err)
		require.NotNil(t, report)
		assert.Equal(t, visitors.MaintenanceOpen, report.Status)
		assert.Equal(t, visitors.PriorityHigh, report.Priority)

		// Update status to in_lavorazione
		err = svc.UpdateMaintenanceStatus(ctx, report.ID, schoolID, visitors.UpdateMaintenanceStatusRequest{
			Status:     visitors.MaintenanceInProgress,
			AssignedTo: "Mario Tecnico",
		})
		require.NoError(t, err)
		assert.Equal(t, visitors.MaintenanceInProgress, report.Status)

		// List filtered by status
		inProgList, err := svc.ListMaintenanceReports(ctx, schoolID, string(visitors.MaintenanceInProgress))
		require.NoError(t, err)
		assert.Len(t, inProgList, 1)

		closedList, err := svc.ListMaintenanceReports(ctx, schoolID, string(visitors.MaintenanceClosed))
		require.NoError(t, err)
		assert.Empty(t, closedList)
	})
}
