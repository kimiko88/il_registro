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

	"registro-backend/internal/verbali"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockVerbaliRepo struct {
	meetings  map[string]*verbali.CouncilMeeting
	verbali   map[string]*verbali.MeetingVerbale
	sigs      map[string][]verbali.VerbaleSignature
	templates map[string]*verbali.MeetingVerbaleTemplate
}

func (m *mockVerbaliRepo) CreateMeeting(ctx context.Context, cm *verbali.CouncilMeeting) error {
	if cm.ID == "" {
		cm.ID = "meeting-1"
	}
	cm.CreatedAt = time.Now()
	m.meetings[cm.ID] = cm
	return nil
}

func (m *mockVerbaliRepo) ListMeetings(ctx context.Context, schoolID, classID string) ([]*verbali.CouncilMeeting, error) {
	res := make([]*verbali.CouncilMeeting, 0)
	for _, cm := range m.meetings {
		res = append(res, cm)
	}
	return res, nil
}

func (m *mockVerbaliRepo) GetMeetingByID(ctx context.Context, id string) (*verbali.CouncilMeeting, error) {
	cm, ok := m.meetings[id]
	if !ok {
		return nil, errors.New("meeting not found")
	}
	return cm, nil
}

func (m *mockVerbaliRepo) GetMeetingWithCoordinator(ctx context.Context, meetingID string) (*verbali.CouncilMeeting, string, error) {
	cm, ok := m.meetings[meetingID]
	if !ok {
		return nil, "", errors.New("meeting not found")
	}
	return cm, "coord-1", nil
}

func (m *mockVerbaliRepo) CreateVerbale(ctx context.Context, v *verbali.MeetingVerbale) error {
	if v.ID == "" {
		v.ID = "verb-1"
	}
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	m.verbali[v.ID] = v
	return nil
}

func (m *mockVerbaliRepo) GetVerbaleByID(ctx context.Context, id string, userID string) (*verbali.MeetingVerbale, error) {
	v, ok := m.verbali[id]
	if !ok {
		return nil, errors.New("verbale not found")
	}
	return v, nil
}

func (m *mockVerbaliRepo) UpdateVerbale(ctx context.Context, v *verbali.MeetingVerbale) error {
	v.UpdatedAt = time.Now()
	m.verbali[v.ID] = v
	return nil
}

func (m *mockVerbaliRepo) DeleteVerbale(ctx context.Context, id string) error {
	delete(m.verbali, id)
	return nil
}

func (m *mockVerbaliRepo) ListVerbali(ctx context.Context, meetingID string, userID string, onlySigned bool) ([]*verbali.MeetingVerbale, error) {
	res := make([]*verbali.MeetingVerbale, 0)
	for _, v := range m.verbali {
		if v.MeetingID == meetingID {
			if !onlySigned || v.IsSigned {
				res = append(res, v)
			}
		}
	}
	return res, nil
}

func (m *mockVerbaliRepo) ListAllVerbali(ctx context.Context, schoolID, classID, userID string, onlySigned bool) ([]*verbali.MeetingVerbale, error) {
	res := make([]*verbali.MeetingVerbale, 0)
	for _, v := range m.verbali {
		if !onlySigned || v.IsSigned {
			res = append(res, v)
		}
	}
	return res, nil
}

func (m *mockVerbaliRepo) SignVerbale(ctx context.Context, verbaleID, userID, ipAddress string) error {
	m.sigs[verbaleID] = append(m.sigs[verbaleID], verbali.VerbaleSignature{
		ID:        "sig-1",
		VerbaleID: verbaleID,
		UserID:    userID,
		SignedAt:  time.Now(),
		IPAddress: ipAddress,
	})
	return nil
}

func (m *mockVerbaliRepo) MarkVerbaleSigned(ctx context.Context, verbaleID string) error {
	if v, ok := m.verbali[verbaleID]; ok {
		v.IsSigned = true
		v.Status = "signed"
		now := time.Now()
		v.SignedAt = &now
	}
	return nil
}

func (m *mockVerbaliRepo) GetSignatures(ctx context.Context, verbaleID string) ([]verbali.VerbaleSignature, error) {
	return m.sigs[verbaleID], nil
}

func (m *mockVerbaliRepo) ClassBelongsToSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	return true, nil
}

func (m *mockVerbaliRepo) CreateTemplate(ctx context.Context, t *verbali.MeetingVerbaleTemplate) error {
	if t.ID == "" {
		t.ID = "tpl-1"
	}
	m.templates[t.ID] = t
	return nil
}

func (m *mockVerbaliRepo) ListTemplates(ctx context.Context, schoolID, meetingType string) ([]*verbali.MeetingVerbaleTemplate, error) {
	res := make([]*verbali.MeetingVerbaleTemplate, 0)
	for _, t := range m.templates {
		res = append(res, t)
	}
	return res, nil
}

func (m *mockVerbaliRepo) GetTemplateByID(ctx context.Context, id string) (*verbali.MeetingVerbaleTemplate, error) {
	t, ok := m.templates[id]
	if !ok {
		return nil, errors.New("template not found")
	}
	return t, nil
}

func (m *mockVerbaliRepo) UpdateTemplate(ctx context.Context, t *verbali.MeetingVerbaleTemplate) error {
	m.templates[t.ID] = t
	return nil
}

func (m *mockVerbaliRepo) DeleteTemplate(ctx context.Context, id, schoolID string) error {
	delete(m.templates, id)
	return nil
}

func TestIntegration_Verbali_Meeting_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockVerbaliRepo{
		meetings:  make(map[string]*verbali.CouncilMeeting),
		verbali:   make(map[string]*verbali.MeetingVerbale),
		sigs:      make(map[string][]verbali.VerbaleSignature),
		templates: make(map[string]*verbali.MeetingVerbaleTemplate),
	}
	svc := verbali.NewService(repo)
	handler := verbali.NewHandler(svc)

	var currentUserID, currentRole string
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", currentUserID)
		c.Set("role", currentRole)
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(r.Group(""))

	// 1. Dirigente creates a template with ODG
	currentUserID = "ds-1"
	currentRole = "principal"
	tplReq := map[string]interface{}{
		"title":            "Modello Consiglio di Classe",
		"meeting_type":     "consiglio_classe",
		"default_agenda":   "1. Andamento didattico\n2. Provvedimenti BES",
		"template_content": "Verbale del consiglio...",
	}
	bTpl, _ := json.Marshal(tplReq)
	wTpl := httptest.NewRecorder()
	reqTpl, _ := http.NewRequest("POST", "/verbali/templates", bytes.NewBuffer(bTpl))
	reqTpl.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wTpl, reqTpl)
	assert.Equal(t, http.StatusCreated, wTpl.Code)

	// 2. Coordinator creates council meeting
	currentUserID = "coord-1"
	currentRole = "teacher"
	mReq := map[string]interface{}{
		"class_id":     "class-1",
		"meeting_type": "consiglio_classe",
		"title":        "Consiglio di Classe 1° Trimestre - Classe 2A",
		"date":         "2026-10-15",
		"start_time":   "15:00",
		"end_time":     "17:00",
		"agenda":       "Andamento didattico-disciplinare e approvazione PDP",
	}
	bM, _ := json.Marshal(mReq)
	wM := httptest.NewRecorder()
	reqM, _ := http.NewRequest("POST", "/verbali/meetings", bytes.NewBuffer(bM))
	reqM.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wM, reqM)
	assert.Equal(t, http.StatusCreated, wM.Code)

	// 3. Verbalista drafts verbale
	currentUserID = "sec-1"
	currentRole = "teacher"
	secID := "sec-1"
	presID := "coord-1"
	vReq := map[string]interface{}{
		"meeting_id":   "meeting-1",
		"title":        "Verbale Consiglio di Classe n. 1",
		"content":      "Alle ore 15:00 si insedia il Consiglio di Classe. Bozza iniziale.",
		"secretary_id": secID,
		"president_id": presID,
	}
	bV, _ := json.Marshal(vReq)
	wV := httptest.NewRecorder()
	reqV, _ := http.NewRequest("POST", "/verbali", bytes.NewBuffer(bV))
	reqV.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wV, reqV)
	assert.Equal(t, http.StatusCreated, wV.Code)

	// 4. Third teacher tries to edit draft -> Blocked with 403
	currentUserID = "other-teacher"
	currentRole = "teacher"
	putReq := map[string]interface{}{
		"title":   "Tentativo non autorizzato",
		"content": "Modifica da terzo docente",
	}
	bPutFail, _ := json.Marshal(putReq)
	wPutFail := httptest.NewRecorder()
	reqPutFail, _ := http.NewRequest("PUT", "/verbali/verb-1", bytes.NewBuffer(bPutFail))
	reqPutFail.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wPutFail, reqPutFail)
	assert.Equal(t, http.StatusForbidden, wPutFail.Code)

	// 5. Coordinator edits draft -> Allowed
	currentUserID = "coord-1"
	currentRole = "teacher"
	putReqCoord := map[string]interface{}{
		"title":        "Verbale Consiglio di Classe n. 1 (Revisione Coordinatore)",
		"content":      "Testo revisionato dal coordinatore.",
		"secretary_id": secID,
		"president_id": presID,
	}
	bPutCoord, _ := json.Marshal(putReqCoord)
	wPutCoord := httptest.NewRecorder()
	reqPutCoord, _ := http.NewRequest("PUT", "/verbali/verb-1", bytes.NewBuffer(bPutCoord))
	reqPutCoord.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wPutCoord, reqPutCoord)
	assert.Equal(t, http.StatusOK, wPutCoord.Code)

	// 6. Dirigente tries to view draft -> Hidden (404)
	currentUserID = "ds-1"
	currentRole = "principal"
	wGetDs := httptest.NewRecorder()
	reqGetDs, _ := http.NewRequest("GET", "/verbali/verb-1", nil)
	r.ServeHTTP(wGetDs, reqGetDs)
	assert.Equal(t, http.StatusNotFound, wGetDs.Code)

	// 7. Verbalista signs verbale -> becomes officially signed and locked
	currentUserID = "sec-1"
	currentRole = "teacher"
	wSign := httptest.NewRecorder()
	reqSign, _ := http.NewRequest("POST", "/verbali/verb-1/sign", nil)
	r.ServeHTTP(wSign, reqSign)
	assert.Equal(t, http.StatusOK, wSign.Code)

	// 8. Attempt to modify after signing -> Blocked (403 VERBALE_LOCKED)
	wPostSignEdit := httptest.NewRecorder()
	reqPostSignEdit, _ := http.NewRequest("PUT", "/verbali/verb-1", bytes.NewBuffer(bPutCoord))
	reqPostSignEdit.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wPostSignEdit, reqPostSignEdit)
	assert.Equal(t, http.StatusForbidden, wPostSignEdit.Code)

	// 9. Dirigente can now view the signed verbale!
	currentUserID = "ds-1"
	currentRole = "principal"
	wGetDsSigned := httptest.NewRecorder()
	reqGetDsSigned, _ := http.NewRequest("GET", "/verbali/verb-1", nil)
	r.ServeHTTP(wGetDsSigned, reqGetDsSigned)
	assert.Equal(t, http.StatusOK, wGetDsSigned.Code)

	var dsResult verbali.MeetingVerbale
	_ = json.Unmarshal(wGetDsSigned.Body.Bytes(), &dsResult)
	assert.True(t, dsResult.IsSigned)
	assert.False(t, dsResult.CanEdit) // Read-only for everyone
}
