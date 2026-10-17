package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/visitors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type memoryVisitorsRepo struct {
	visitors    map[string]*visitors.Visitor
	earlyExits  map[string]*visitors.EarlyExit
	maintenance map[string]*visitors.MaintenanceReport
}

func newMemoryVisitorsRepo() *memoryVisitorsRepo {
	return &memoryVisitorsRepo{
		visitors:    make(map[string]*visitors.Visitor),
		earlyExits:  make(map[string]*visitors.EarlyExit),
		maintenance: make(map[string]*visitors.MaintenanceReport),
	}
}

func (m *memoryVisitorsRepo) CreateVisitor(ctx context.Context, v *visitors.Visitor) error {
	if v.ID == "" {
		v.ID = "vis-" + time.Now().Format("150405.000000")
	}
	v.EntryTime = time.Now()
	v.CreatedAt = time.Now()
	m.visitors[v.ID] = v
	return nil
}

func (m *memoryVisitorsRepo) RecordVisitorExit(ctx context.Context, id, schoolID string, notes string) error {
	v, ok := m.visitors[id]
	if !ok || (schoolID != "" && v.SchoolID != "" && v.SchoolID != schoolID) {
		return errors.New("visitatore non trovato")
	}
	now := time.Now()
	v.ExitTime = &now
	if notes != "" {
		v.Notes = notes
	}
	return nil
}

func (m *memoryVisitorsRepo) ListTodayVisitors(ctx context.Context, schoolID, date string) ([]*visitors.Visitor, error) {
	var out []*visitors.Visitor
	for _, v := range m.visitors {
		if schoolID != "" && v.SchoolID != "" && v.SchoolID != schoolID {
			continue
		}
		out = append(out, v)
	}
	return out, nil
}

func (m *memoryVisitorsRepo) GetVisitor(ctx context.Context, id, schoolID string) (*visitors.Visitor, error) {
	v, ok := m.visitors[id]
	if !ok || (schoolID != "" && v.SchoolID != "" && v.SchoolID != schoolID) {
		return nil, errors.New("visitatore non trovato")
	}
	return v, nil
}

func (m *memoryVisitorsRepo) CreateEarlyExit(ctx context.Context, e *visitors.EarlyExit) error {
	if e.ID == "" {
		e.ID = "exit-" + time.Now().Format("150405.000000")
	}
	e.ExitTime = time.Now()
	e.CreatedAt = time.Now()
	m.earlyExits[e.ID] = e
	return nil
}

func (m *memoryVisitorsRepo) RecordStudentReturn(ctx context.Context, id, schoolID, notes string) error {
	e, ok := m.earlyExits[id]
	if !ok || (schoolID != "" && e.SchoolID != "" && e.SchoolID != schoolID) {
		return errors.New("uscita non trovata")
	}
	now := time.Now()
	e.ReturnTime = &now
	if notes != "" {
		e.Notes = notes
	}
	return nil
}

func (m *memoryVisitorsRepo) ListTodayEarlyExits(ctx context.Context, schoolID, date string) ([]*visitors.EarlyExit, error) {
	var out []*visitors.EarlyExit
	for _, e := range m.earlyExits {
		if schoolID != "" && e.SchoolID != "" && e.SchoolID != schoolID {
			continue
		}
		out = append(out, e)
	}
	return out, nil
}

func (m *memoryVisitorsRepo) CreateMaintenanceReport(ctx context.Context, r *visitors.MaintenanceReport) error {
	if r.ID == "" {
		r.ID = "maint-" + time.Now().Format("150405.000000")
	}
	r.Status = visitors.MaintenanceOpen
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	m.maintenance[r.ID] = r
	return nil
}

func (m *memoryVisitorsRepo) ListMaintenanceReports(ctx context.Context, schoolID string, statusFilter string) ([]*visitors.MaintenanceReport, error) {
	var out []*visitors.MaintenanceReport
	for _, mr := range m.maintenance {
		if schoolID != "" && mr.SchoolID != "" && mr.SchoolID != schoolID {
			continue
		}
		if statusFilter != "" && string(mr.Status) != statusFilter {
			continue
		}
		out = append(out, mr)
	}
	return out, nil
}

func (m *memoryVisitorsRepo) UpdateMaintenanceStatus(ctx context.Context, id, schoolID string, req visitors.UpdateMaintenanceStatusRequest) error {
	mr, ok := m.maintenance[id]
	if !ok || (schoolID != "" && mr.SchoolID != "" && mr.SchoolID != schoolID) {
		return errors.New("segnalazione non trovata")
	}
	mr.Status = req.Status
	if req.AssignedTo != "" {
		mr.AssignedTo = &req.AssignedTo
	}
	mr.UpdatedAt = time.Now()
	return nil
}

func setupVisitorsRouter(repo visitors.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-User-ID"))
		c.Set("role", c.GetHeader("X-User-Role"))
		c.Set("school_id", c.GetHeader("X-School-ID"))
		c.Next()
	})

	svc := visitors.NewService(repo)
	handler := visitors.NewHandler(svc)
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_VisitorsRegistry_FullLifecycle(t *testing.T) {
	repo := newMemoryVisitorsRepo()
	r := setupVisitorsRouter(repo)

	schoolID := "school-1"
	staffID := "staff-reception-1"

	// 1. Collaboratore scolastico registers a visitor with badge B-101
	visitorReq := visitors.RegisterVisitorRequest{
		Name:        "Mario Rossi",
		DocumentID:  "CI-CA12345XY",
		Purpose:     visitors.PurposeParent,
		HostName:    "Prof. Bianchi",
		BadgeNumber: "B-101",
		Notes:       "Colloquio individuale straordinario",
	}
	b, _ := json.Marshal(visitorReq)
	req1, _ := http.NewRequest("POST", "/api/v1/visitors", bytes.NewBuffer(b))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-User-ID", staffID)
	req1.Header.Set("X-User-Role", "collaboratore_scolastico")
	req1.Header.Set("X-School-ID", schoolID)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	var createdVisitor visitors.Visitor
	err := json.Unmarshal(w1.Body.Bytes(), &createdVisitor)
	assert.NoError(t, err)
	visitorID := createdVisitor.ID
	assert.NotEmpty(t, visitorID)

	// 2. Unauthorized role (student or teacher) attempting access -> 403 Forbidden
	reqTeacher, _ := http.NewRequest("GET", "/api/v1/visitors", nil)
	reqTeacher.Header.Set("X-User-ID", "teacher-1")
	reqTeacher.Header.Set("X-User-Role", "teacher")
	reqTeacher.Header.Set("X-School-ID", schoolID)
	wTeacher := httptest.NewRecorder()
	r.ServeHTTP(wTeacher, reqTeacher)
	assert.Equal(t, http.StatusForbidden, wTeacher.Code)

	// 3. List active visitors today
	reqList, _ := http.NewRequest("GET", "/api/v1/visitors", nil)
	reqList.Header.Set("X-User-ID", staffID)
	reqList.Header.Set("X-User-Role", "collaboratore_scolastico")
	reqList.Header.Set("X-School-ID", schoolID)
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)
	assert.Equal(t, http.StatusOK, wList.Code)

	var list []*visitors.Visitor
	err = json.Unmarshal(wList.Body.Bytes(), &list)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Nil(t, list[0].ExitTime)

	// 4. Record student early exit
	earlyExitReq := visitors.RecordEarlyExitRequest{
		StudentID:     "student-42",
		DelegateeName: "Mario Rossi",
		DelegateRel:   "Padre",
		ReasonCode:    "visita_medica",
		Notes:         "Uscita ore 11:30 per visita odontoiatrica",
	}
	bEarly, _ := json.Marshal(earlyExitReq)
	reqEarly, _ := http.NewRequest("POST", "/api/v1/visitors/early-exits", bytes.NewBuffer(bEarly))
	reqEarly.Header.Set("Content-Type", "application/json")
	reqEarly.Header.Set("X-User-ID", staffID)
	reqEarly.Header.Set("X-User-Role", "assistente_amministrativo")
	reqEarly.Header.Set("X-School-ID", schoolID)
	wEarly := httptest.NewRecorder()
	r.ServeHTTP(wEarly, reqEarly)
	assert.Equal(t, http.StatusCreated, wEarly.Code)

	var createdExit visitors.EarlyExit
	err = json.Unmarshal(wEarly.Body.Bytes(), &createdExit)
	assert.NoError(t, err)
	exitID := createdExit.ID

	// 5. Record student return
	returnReq := visitors.RecordStudentReturnRequest{
		Notes: "Rientrato a scuola ore 12:45 con certificato",
	}
	bReturn, _ := json.Marshal(returnReq)
	reqReturn, _ := http.NewRequest("PATCH", "/api/v1/visitors/early-exits/"+exitID+"/return", bytes.NewBuffer(bReturn))
	reqReturn.Header.Set("Content-Type", "application/json")
	reqReturn.Header.Set("X-User-ID", staffID)
	reqReturn.Header.Set("X-User-Role", "assistente_amministrativo")
	reqReturn.Header.Set("X-School-ID", schoolID)
	wReturn := httptest.NewRecorder()
	r.ServeHTTP(wReturn, reqReturn)
	assert.Equal(t, http.StatusOK, wReturn.Code)

	// 6. Record visitor exit and badge return
	exitVisitorReq := visitors.RecordExitRequest{
		Notes: "Uscita regolare, badge B-101 riconsegnato",
	}
	bExitVis, _ := json.Marshal(exitVisitorReq)
	reqExitVis, _ := http.NewRequest("PATCH", "/api/v1/visitors/"+visitorID+"/exit", bytes.NewBuffer(bExitVis))
	reqExitVis.Header.Set("Content-Type", "application/json")
	reqExitVis.Header.Set("X-User-ID", staffID)
	reqExitVis.Header.Set("X-User-Role", "collaboratore_scolastico")
	reqExitVis.Header.Set("X-School-ID", schoolID)
	wExitVis := httptest.NewRecorder()
	r.ServeHTTP(wExitVis, reqExitVis)
	assert.Equal(t, http.StatusOK, wExitVis.Code)

	// 7. Maintenance incident reporting
	maintReq := visitors.CreateMaintenanceReportRequest{
		Location:    "Laboratorio Informatica 2",
		Category:    "elettrico",
		Description: "Presa elettrica danneggiata alla postazione 5",
		Priority:    visitors.PriorityHigh,
	}
	bMaint, _ := json.Marshal(maintReq)
	reqMaint, _ := http.NewRequest("POST", "/api/v1/visitors/maintenance", bytes.NewBuffer(bMaint))
	reqMaint.Header.Set("Content-Type", "application/json")
	reqMaint.Header.Set("X-User-ID", staffID)
	reqMaint.Header.Set("X-User-Role", "assistente_tecnico")
	reqMaint.Header.Set("X-School-ID", schoolID)
	wMaint := httptest.NewRecorder()
	r.ServeHTTP(wMaint, reqMaint)
	assert.Equal(t, http.StatusCreated, wMaint.Code)

	var createdMaint visitors.MaintenanceReport
	err = json.Unmarshal(wMaint.Body.Bytes(), &createdMaint)
	assert.NoError(t, err)
	assert.Equal(t, visitors.MaintenanceOpen, createdMaint.Status)

	// 8. Update maintenance report status to in_lavorazione
	updateMaintReq := visitors.UpdateMaintenanceStatusRequest{
		Status:     visitors.MaintenanceInProgress,
		AssignedTo: "Tecnico Elettricista Ditta",
	}
	bUpMaint, _ := json.Marshal(updateMaintReq)
	reqUpMaint, _ := http.NewRequest("PATCH", "/api/v1/visitors/maintenance/"+createdMaint.ID+"/status", bytes.NewBuffer(bUpMaint))
	reqUpMaint.Header.Set("Content-Type", "application/json")
	reqUpMaint.Header.Set("X-User-ID", staffID)
	reqUpMaint.Header.Set("X-User-Role", "dsga")
	reqUpMaint.Header.Set("X-School-ID", schoolID)
	wUpMaint := httptest.NewRecorder()
	r.ServeHTTP(wUpMaint, reqUpMaint)
	assert.Equal(t, http.StatusOK, wUpMaint.Code)
}
