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

	"registro-backend/internal/personnel_desk"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type memoryPersonnelDeskRepo struct {
	requests map[string]*personnel_desk.DeskRequest
}

func newMemoryPersonnelDeskRepo() *memoryPersonnelDeskRepo {
	return &memoryPersonnelDeskRepo{
		requests: make(map[string]*personnel_desk.DeskRequest),
	}
}

func (m *memoryPersonnelDeskRepo) Create(ctx context.Context, r *personnel_desk.DeskRequest) error {
	if r.ID == "" {
		r.ID = "req-" + time.Now().Format("150405.000000")
	}
	m.requests[r.ID] = r
	return nil
}

func (m *memoryPersonnelDeskRepo) GetByID(ctx context.Context, schoolID, id string) (*personnel_desk.DeskRequest, error) {
	req, ok := m.requests[id]
	if !ok || (schoolID != "" && req.SchoolID != "" && req.SchoolID != schoolID) {
		return nil, errors.New("richiesta non trovata")
	}
	return req, nil
}

func (m *memoryPersonnelDeskRepo) List(ctx context.Context, schoolID, applicantID, status string) ([]*personnel_desk.DeskRequest, error) {
	var out []*personnel_desk.DeskRequest
	for _, r := range m.requests {
		if schoolID != "" && r.SchoolID != "" && r.SchoolID != schoolID {
			continue
		}
		if applicantID != "" && r.ApplicantID != applicantID {
			continue
		}
		if status != "" && string(r.Status) != status {
			continue
		}
		out = append(out, r)
	}
	return out, nil
}

func (m *memoryPersonnelDeskRepo) Update(ctx context.Context, r *personnel_desk.DeskRequest) error {
	m.requests[r.ID] = r
	return nil
}

func (m *memoryPersonnelDeskRepo) Delete(ctx context.Context, schoolID, id, applicantID string) error {
	req, ok := m.requests[id]
	if !ok || req.ApplicantID != applicantID {
		return errors.New("richiesta non trovata")
	}
	delete(m.requests, id)
	return nil
}

func TestIntegration_PersonnelDesk_FullWorkflow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMemoryPersonnelDeskRepo()

	rWithAuth := gin.New()
	rWithAuth.Use(func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-User-ID"))
		c.Set("role", c.GetHeader("X-User-Role"))
		c.Set("school_id", c.GetHeader("X-School-ID"))
		c.Next()
	})
	svc := personnel_desk.NewService(repo)
	handler := personnel_desk.NewHandler(svc)
	api := rWithAuth.Group("/api/v1")
	handler.RegisterRoutes(api)

	schoolID := "school-test-1"
	employeeID := "employee-123"

	// Step 1: Employee creates draft
	createBody := personnel_desk.CreateDeskRequestInput{
		Category:    personnel_desk.CatPermessoBreve,
		SubCategory: "permesso_breve",
		StartDate:   "2026-10-01",
		EndDate:     "2026-10-01",
		Hours:       2.0,
		Description: "Permesso breve per motivi personali",
		SubmitNow:   false,
	}
	b, _ := json.Marshal(createBody)
	req1, _ := http.NewRequest("POST", "/api/v1/personnel-desk/requests", bytes.NewBuffer(b))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-User-ID", employeeID)
	req1.Header.Set("X-User-Role", "teacher")
	req1.Header.Set("X-School-ID", schoolID)
	w1 := httptest.NewRecorder()
	rWithAuth.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	var created personnel_desk.DeskRequest
	err := json.Unmarshal(w1.Body.Bytes(), &created)
	assert.NoError(t, err)
	assert.Equal(t, personnel_desk.StatusDraft, created.Status)
	reqID := created.ID

	// Step 2: Employee submits request
	req2, _ := http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/"+reqID+"/submit", nil)
	req2.Header.Set("X-User-ID", employeeID)
	req2.Header.Set("X-User-Role", "teacher")
	req2.Header.Set("X-School-ID", schoolID)
	w2 := httptest.NewRecorder()
	rWithAuth.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var submitted personnel_desk.DeskRequest
	err = json.Unmarshal(w2.Body.Bytes(), &submitted)
	assert.NoError(t, err)
	assert.Equal(t, personnel_desk.StatusSubmitted, submitted.Status)

	// Step 3: Unauthorized role (student) attempts AA Review -> 403 Forbidden
	aaInput := personnel_desk.AAReviewInput{
		Note:    "Documentazione verificata con successo",
		Approve: true,
	}
	bAA, _ := json.Marshal(aaInput)
	req3Unauth, _ := http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/"+reqID+"/aa-review", bytes.NewBuffer(bAA))
	req3Unauth.Header.Set("Content-Type", "application/json")
	req3Unauth.Header.Set("X-User-ID", "student-1")
	req3Unauth.Header.Set("X-User-Role", "student")
	req3Unauth.Header.Set("X-School-ID", schoolID)
	w3Unauth := httptest.NewRecorder()
	rWithAuth.ServeHTTP(w3Unauth, req3Unauth)
	assert.Equal(t, http.StatusForbidden, w3Unauth.Code)

	// Step 4: Administrative Assistant (AA) performs review -> status dsga_review
	req3, _ := http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/"+reqID+"/aa-review", bytes.NewBuffer(bAA))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("X-User-ID", "aa-user-1")
	req3.Header.Set("X-User-Role", "assistente_amministrativo")
	req3.Header.Set("X-School-ID", schoolID)
	w3 := httptest.NewRecorder()
	rWithAuth.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	var reviewed personnel_desk.DeskRequest
	err = json.Unmarshal(w3.Body.Bytes(), &reviewed)
	assert.NoError(t, err)
	assert.Equal(t, personnel_desk.StatusDSGAReview, reviewed.Status)

	// Step 5: DSGA affixes visto contabile -> status ds_review
	dsgaInput := personnel_desk.DSGASignInput{
		Note:    "Copertura finanziaria e monte ore verificati",
		Approve: true,
	}
	bDSGA, _ := json.Marshal(dsgaInput)
	req4, _ := http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/"+reqID+"/dsga-sign", bytes.NewBuffer(bDSGA))
	req4.Header.Set("Content-Type", "application/json")
	req4.Header.Set("X-User-ID", "dsga-user-1")
	req4.Header.Set("X-User-Role", "dsga")
	req4.Header.Set("X-School-ID", schoolID)
	w4 := httptest.NewRecorder()
	rWithAuth.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)

	var signed personnel_desk.DeskRequest
	err = json.Unmarshal(w4.Body.Bytes(), &signed)
	assert.NoError(t, err)
	assert.Equal(t, personnel_desk.StatusDSReview, signed.Status)

	// Step 6: Principal issues official administrative decree -> status approved
	dsInput := personnel_desk.DSApproveInput{
		DecreeNum: "DEC-2026/042",
		Note:      "Si autorizza il permesso breve richiesto.",
		Approve:   true,
	}
	bDS, _ := json.Marshal(dsInput)
	req5, _ := http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/"+reqID+"/ds-approve", bytes.NewBuffer(bDS))
	req5.Header.Set("Content-Type", "application/json")
	req5.Header.Set("X-User-ID", "principal-1")
	req5.Header.Set("X-User-Role", "principal")
	req5.Header.Set("X-School-ID", schoolID)
	w5 := httptest.NewRecorder()
	rWithAuth.ServeHTTP(w5, req5)
	assert.Equal(t, http.StatusOK, w5.Code)

	var approved personnel_desk.DeskRequest
	err = json.Unmarshal(w5.Body.Bytes(), &approved)
	assert.NoError(t, err)
	assert.Equal(t, personnel_desk.StatusApproved, approved.Status)
	assert.NotNil(t, approved.DSDecreeNum)
	assert.Equal(t, "DEC-2026/042", *approved.DSDecreeNum)

	// Step 7: Employee queries their approved request
	req6, _ := http.NewRequest("GET", "/api/v1/personnel-desk/requests/"+reqID, nil)
	req6.Header.Set("X-User-ID", employeeID)
	req6.Header.Set("X-User-Role", "teacher")
	req6.Header.Set("X-School-ID", schoolID)
	w6 := httptest.NewRecorder()
	rWithAuth.ServeHTTP(w6, req6)
	assert.Equal(t, http.StatusOK, w6.Code)
}

func TestIntegration_PersonnelDesk_NegativeAndRejection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMemoryPersonnelDeskRepo()
	rWithAuth := gin.New()
	rWithAuth.Use(func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-User-ID"))
		c.Set("role", c.GetHeader("X-User-Role"))
		c.Set("school_id", c.GetHeader("X-School-ID"))
		c.Next()
	})
	svc := personnel_desk.NewService(repo)
	handler := personnel_desk.NewHandler(svc)
	api := rWithAuth.Group("/api/v1")
	handler.RegisterRoutes(api)

	schoolID := "school-1"
	employeeID := "emp-reject"

	// 1. Employee creates and submits request immediately
	createBody := personnel_desk.CreateDeskRequestInput{
		Category:    personnel_desk.CatFerie,
		SubCategory: "ferie_estive",
		StartDate:   "2026-07-01",
		EndDate:     "2026-07-15",
		Days:        15,
		Description: "Richiesta ferie anticipate",
		SubmitNow:   true,
	}
	b, _ := json.Marshal(createBody)
	req, _ := http.NewRequest("POST", "/api/v1/personnel-desk/requests", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", employeeID)
	req.Header.Set("X-User-Role", "collaboratore_scolastico")
	req.Header.Set("X-School-ID", schoolID)
	w := httptest.NewRecorder()
	rWithAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var created personnel_desk.DeskRequest
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	reqID := created.ID

	// 2. AA rejects request -> status rejected
	aaInput := personnel_desk.AAReviewInput{
		Note:    "Periodo non concedibile per esigenze di servizio",
		Approve: false,
	}
	bAA, _ := json.Marshal(aaInput)
	reqAA, _ := http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/"+reqID+"/aa-review", bytes.NewBuffer(bAA))
	reqAA.Header.Set("Content-Type", "application/json")
	reqAA.Header.Set("X-User-ID", "aa-user")
	reqAA.Header.Set("X-User-Role", "secretary")
	reqAA.Header.Set("X-School-ID", schoolID)
	wAA := httptest.NewRecorder()
	rWithAuth.ServeHTTP(wAA, reqAA)
	assert.Equal(t, http.StatusOK, wAA.Code)

	var rejected personnel_desk.DeskRequest
	_ = json.Unmarshal(wAA.Body.Bytes(), &rejected)
	assert.Equal(t, personnel_desk.StatusRejected, rejected.Status)
}
