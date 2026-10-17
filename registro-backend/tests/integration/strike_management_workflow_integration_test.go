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

	"registro-backend/internal/strike"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type memoryStrikeRepo struct {
	notices      map[string]*strike.StrikeNotice
	declarations map[string]*strike.StrikeDeclaration // key: noticeID + ":" + userID
}

func newMemoryStrikeRepo() *memoryStrikeRepo {
	return &memoryStrikeRepo{
		notices:      make(map[string]*strike.StrikeNotice),
		declarations: make(map[string]*strike.StrikeDeclaration),
	}
}

func (m *memoryStrikeRepo) CreateNotice(ctx context.Context, n *strike.StrikeNotice) error {
	if n.ID == "" {
		n.ID = "notice-" + time.Now().Format("150405.000000")
	}
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	m.notices[n.ID] = n
	return nil
}

func (m *memoryStrikeRepo) GetNoticeByID(ctx context.Context, id, schoolID string) (*strike.StrikeNotice, error) {
	n, ok := m.notices[id]
	if !ok || (schoolID != "" && n.SchoolID != "" && n.SchoolID != schoolID) {
		return nil, errors.New("avviso non trovato")
	}
	return n, nil
}

func (m *memoryStrikeRepo) ListNotices(ctx context.Context, schoolID string) ([]*strike.StrikeNotice, error) {
	var out []*strike.StrikeNotice
	for _, n := range m.notices {
		if schoolID != "" && n.SchoolID != "" && n.SchoolID != schoolID {
			continue
		}
		out = append(out, n)
	}
	return out, nil
}

func (m *memoryStrikeRepo) DeleteNotice(ctx context.Context, id, schoolID string) error {
	n, ok := m.notices[id]
	if !ok || (schoolID != "" && n.SchoolID != "" && n.SchoolID != schoolID) {
		return errors.New("avviso non trovato")
	}
	delete(m.notices, id)
	return nil
}

func (m *memoryStrikeRepo) UpsertDeclaration(ctx context.Context, d *strike.StrikeDeclaration) error {
	key := d.StrikeNoticeID + ":" + d.UserID
	m.declarations[key] = d
	return nil
}

func (m *memoryStrikeRepo) GetDeclaration(ctx context.Context, noticeID, userID string) (*strike.StrikeDeclaration, error) {
	key := noticeID + ":" + userID
	d, ok := m.declarations[key]
	if !ok {
		return nil, nil
	}
	return d, nil
}

func (m *memoryStrikeRepo) GetNoticeSummary(ctx context.Context, noticeID, schoolID string) (*strike.StrikeNoticeSummaryResponse, error) {
	n, err := m.GetNoticeByID(ctx, noticeID, schoolID)
	if err != nil {
		return nil, err
	}

	totalAnswered := 0
	participates := 0
	notParticipates := 0
	undecided := 0

	for _, d := range m.declarations {
		if d.StrikeNoticeID == noticeID {
			totalAnswered++
			switch d.Intention {
			case strike.IntentionParticipates:
				participates++
			case strike.IntentionNotParticipates:
				notParticipates++
			case strike.IntentionUndecided:
				undecided++
			}
		}
	}

	return &strike.StrikeNoticeSummaryResponse{
		Notice:               n,
		TotalStaff:           10,
		TotalAnswered:        totalAnswered,
		ParticipatesCount:    participates,
		NotParticipatesCount: notParticipates,
		UndecidedCount:       undecided,
		UnansweredCount:      10 - totalAnswered,
		ParticipatesPercent:  float64(participates) / 10.0 * 100.0,
	}, nil
}

func (m *memoryStrikeRepo) CreateBachecaCommunication(ctx context.Context, schoolID, senderID, title, content string, deadline time.Time) (string, error) {
	return "comm-123", nil
}

func (m *memoryStrikeRepo) ResolveSchoolID(ctx context.Context, userID string) string {
	return "school-1"
}

func setupStrikeRouter(repo strike.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-User-ID"))
		c.Set("role", c.GetHeader("X-User-Role"))
		c.Set("school_id", c.GetHeader("X-School-ID"))
		c.Next()
	})

	svc := strike.NewService(repo)
	handler := strike.NewHandler(svc)
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_StrikeManagement_FullLifecycle(t *testing.T) {
	repo := newMemoryStrikeRepo()
	r := setupStrikeRouter(repo)

	schoolID := "school-1"
	dsgaID := "dsga-10"

	// 1. DSGA creates a strike notice
	createNoticeReq := strike.CreateStrikeNoticeRequest{
		Title:               "Sciopero Generale del Comparto Istruzione e Ricerca",
		ProclaimedBy:        "FLC-CGIL, CISL Scuola, UIL Scuola",
		StrikeDate:          "2026-10-15",
		DeclarationDeadline: time.Now().Add(48 * time.Hour).Format(time.RFC3339),
		Content:             "Si comunica lo sciopero intera giornata per tutto il personale docente e ATA.",
		Notes:               "Servizi minimi garantiti ai sensi dell'accordo ARAN 2/12/2020",
		PublishToBacheca:    true,
	}
	b, _ := json.Marshal(createNoticeReq)
	req1, _ := http.NewRequest("POST", "/api/v1/strike-notices", bytes.NewBuffer(b))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-User-ID", dsgaID)
	req1.Header.Set("X-User-Role", "dsga")
	req1.Header.Set("X-School-ID", schoolID)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	var createdNotice strike.StrikeNotice
	err := json.Unmarshal(w1.Body.Bytes(), &createdNotice)
	assert.NoError(t, err)
	noticeID := createdNotice.ID
	assert.NotEmpty(t, noticeID)

	// 2. Student attempts to create notice -> 403 Forbidden
	reqStudent, _ := http.NewRequest("POST", "/api/v1/strike-notices", bytes.NewBuffer(b))
	reqStudent.Header.Set("Content-Type", "application/json")
	reqStudent.Header.Set("X-User-ID", "student-1")
	reqStudent.Header.Set("X-User-Role", "student")
	reqStudent.Header.Set("X-School-ID", schoolID)
	wStudent := httptest.NewRecorder()
	r.ServeHTTP(wStudent, reqStudent)
	assert.Equal(t, http.StatusForbidden, wStudent.Code)

	// 3. Teacher submits declaration (participates)
	declareReq := strike.SubmitDeclarationRequest{
		Intention: strike.IntentionParticipates,
		Notes:     "Aderisco allo sciopero",
	}
	bDeclare, _ := json.Marshal(declareReq)
	req2, _ := http.NewRequest("POST", "/api/v1/strike-notices/"+noticeID+"/declare", bytes.NewBuffer(bDeclare))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-User-ID", "teacher-1")
	req2.Header.Set("X-User-Role", "teacher")
	req2.Header.Set("X-School-ID", schoolID)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// 4. ATA staff submits declaration (not_participates)
	declareReq2 := strike.SubmitDeclarationRequest{
		Intention: strike.IntentionNotParticipates,
		Notes:     "In servizio per contingenti minimi",
	}
	bDeclare2, _ := json.Marshal(declareReq2)
	req3, _ := http.NewRequest("POST", "/api/v1/strike-notices/"+noticeID+"/declare", bytes.NewBuffer(bDeclare2))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("X-User-ID", "ata-1")
	req3.Header.Set("X-User-Role", "collaboratore_scolastico")
	req3.Header.Set("X-School-ID", schoolID)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	// 5. DSGA requests strike summary
	reqSummary, _ := http.NewRequest("GET", "/api/v1/strike-notices/"+noticeID+"/summary", nil)
	reqSummary.Header.Set("X-User-ID", dsgaID)
	reqSummary.Header.Set("X-User-Role", "dsga")
	reqSummary.Header.Set("X-School-ID", schoolID)
	wSummary := httptest.NewRecorder()
	r.ServeHTTP(wSummary, reqSummary)
	assert.Equal(t, http.StatusOK, wSummary.Code)

	var summary strike.StrikeNoticeSummaryResponse
	err = json.Unmarshal(wSummary.Body.Bytes(), &summary)
	assert.NoError(t, err)
	assert.Equal(t, 2, summary.TotalAnswered)
	assert.Equal(t, 1, summary.ParticipatesCount)
	assert.Equal(t, 1, summary.NotParticipatesCount)
	assert.Equal(t, 10.0, summary.ParticipatesPercent)

	// 6. Teacher attempts to delete notice -> 403 Forbidden
	reqDelUnauth, _ := http.NewRequest("DELETE", "/api/v1/strike-notices/"+noticeID, nil)
	reqDelUnauth.Header.Set("X-User-ID", "teacher-1")
	reqDelUnauth.Header.Set("X-User-Role", "teacher")
	reqDelUnauth.Header.Set("X-School-ID", schoolID)
	wDelUnauth := httptest.NewRecorder()
	r.ServeHTTP(wDelUnauth, reqDelUnauth)
	assert.Equal(t, http.StatusForbidden, wDelUnauth.Code)

	// 7. DSGA deletes notice -> 200 OK
	reqDel, _ := http.NewRequest("DELETE", "/api/v1/strike-notices/"+noticeID, nil)
	reqDel.Header.Set("X-User-ID", dsgaID)
	reqDel.Header.Set("X-User-Role", "dsga")
	reqDel.Header.Set("X-School-ID", schoolID)
	wDel := httptest.NewRecorder()
	r.ServeHTTP(wDel, reqDel)
	assert.Equal(t, http.StatusOK, wDel.Code)
}
