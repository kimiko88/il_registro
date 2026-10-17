package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/strike"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_StrikeDeclarations_EdgeCasesAndValidation(t *testing.T) {
	repo := newMemoryStrikeRepo()
	r := setupStrikeRouter(repo)

	schoolID := "school-strike-edge-1"
	adminID := "admin-strike-1"

	// 1. Create a notice that has a valid future deadline
	futureDeadline := time.Now().Add(72 * time.Hour)
	createReq := strike.CreateStrikeNoticeRequest{
		Title:               "Sciopero Nazionale Trasporti e Scuola",
		ProclaimedBy:        "Sindacato Autonomo",
		StrikeDate:          "2026-11-20",
		DeclarationDeadline: futureDeadline.Format(time.RFC3339),
		Content:             "Comunicazione preventiva ai sensi di legge",
	}
	b, _ := json.Marshal(createReq)
	reqCreate, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices", bytes.NewBuffer(b))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("X-User-ID", adminID)
	reqCreate.Header.Set("X-User-Role", "admin")
	reqCreate.Header.Set("X-School-ID", schoolID)
	wCreate := httptest.NewRecorder()
	r.ServeHTTP(wCreate, reqCreate)
	require.Equal(t, http.StatusCreated, wCreate.Code)

	var activeNotice strike.StrikeNotice
	err := json.Unmarshal(wCreate.Body.Bytes(), &activeNotice)
	require.NoError(t, err)
	activeNoticeID := activeNotice.ID

	// 2. Declaration with invalid intention -> 400 Bad Request
	invalidDeclareReq := map[string]string{
		"intention": "maybe",
	}
	bInvalid, _ := json.Marshal(invalidDeclareReq)
	reqInvalid, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices/"+activeNoticeID+"/declare", bytes.NewBuffer(bInvalid))
	reqInvalid.Header.Set("Content-Type", "application/json")
	reqInvalid.Header.Set("X-User-ID", "teacher-user-1")
	reqInvalid.Header.Set("X-User-Role", "teacher")
	reqInvalid.Header.Set("X-School-ID", schoolID)
	wInvalid := httptest.NewRecorder()
	r.ServeHTTP(wInvalid, reqInvalid)
	assert.Equal(t, http.StatusBadRequest, wInvalid.Code)

	// 3. Declaration for non-existent notice ID -> 404 Not Found
	reqNotFound, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices/non-existent-id/declare", bytes.NewBufferString(`{"intention":"participates"}`))
	reqNotFound.Header.Set("Content-Type", "application/json")
	reqNotFound.Header.Set("X-User-ID", "teacher-user-1")
	reqNotFound.Header.Set("X-User-Role", "teacher")
	reqNotFound.Header.Set("X-School-ID", schoolID)
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)
	assert.Equal(t, http.StatusNotFound, wNotFound.Code)

	// 4. Initial declaration: teacher-user-1 declares 'undecided'
	reqDec1, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices/"+activeNoticeID+"/declare", bytes.NewBufferString(`{"intention":"undecided"}`))
	reqDec1.Header.Set("Content-Type", "application/json")
	reqDec1.Header.Set("X-User-ID", "teacher-user-1")
	reqDec1.Header.Set("X-User-Role", "teacher")
	reqDec1.Header.Set("X-School-ID", schoolID)
	wDec1 := httptest.NewRecorder()
	r.ServeHTTP(wDec1, reqDec1)
	assert.Equal(t, http.StatusOK, wDec1.Code)

	// 5. Update intention before deadline: teacher-user-1 changes to 'participates'
	reqDec2, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices/"+activeNoticeID+"/declare", bytes.NewBufferString(`{"intention":"participates","notes":"Adesione confermata"}`))
	reqDec2.Header.Set("Content-Type", "application/json")
	reqDec2.Header.Set("X-User-ID", "teacher-user-1")
	reqDec2.Header.Set("X-User-Role", "teacher")
	reqDec2.Header.Set("X-School-ID", schoolID)
	wDec2 := httptest.NewRecorder()
	r.ServeHTTP(wDec2, reqDec2)
	assert.Equal(t, http.StatusOK, wDec2.Code)

	// 6. ATA secretary roles can access summary
	ataRoles := []string{"secretary", "assistente_amministrativo", "assistente_personale"}
	for _, role := range ataRoles {
		reqSum, _ := http.NewRequest(http.MethodGet, "/api/v1/strike-notices/"+activeNoticeID+"/summary", nil)
		reqSum.Header.Set("X-User-ID", "staff-id-"+role)
		reqSum.Header.Set("X-User-Role", role)
		reqSum.Header.Set("X-School-ID", schoolID)
		wSum := httptest.NewRecorder()
		r.ServeHTTP(wSum, reqSum)
		assert.Equal(t, http.StatusOK, wSum.Code, "Role %s should have access to summary", role)

		var summResp strike.StrikeNoticeSummaryResponse
		err = json.Unmarshal(wSum.Body.Bytes(), &summResp)
		require.NoError(t, err)
		// Exactly 1 answered declaration (updated from undecided to participates)
		assert.Equal(t, 1, summResp.TotalAnswered)
		assert.Equal(t, 1, summResp.ParticipatesCount)
		assert.Equal(t, 0, summResp.UndecidedCount)
	}

	// 7. Forbidden roles blocked from summary
	forbiddenRoles := []string{"teacher", "student", "parent", "collaboratore_scolastico"}
	for _, role := range forbiddenRoles {
		reqForbidden, _ := http.NewRequest(http.MethodGet, "/api/v1/strike-notices/"+activeNoticeID+"/summary", nil)
		reqForbidden.Header.Set("X-User-ID", "user-"+role)
		reqForbidden.Header.Set("X-User-Role", role)
		reqForbidden.Header.Set("X-School-ID", schoolID)
		wForbidden := httptest.NewRecorder()
		r.ServeHTTP(wForbidden, reqForbidden)
		assert.Equal(t, http.StatusForbidden, wForbidden.Code, "Role %s must be forbidden from summary", role)
	}

	// 8. Create an EXPIRED notice and verify submitting declaration is rejected with 400
	pastDeadline := time.Now().Add(-2 * time.Hour)
	expiredNotice := strike.StrikeNotice{
		ID:                  "notice-expired-1",
		SchoolID:            schoolID,
		Title:               "Sciopero Scaduto",
		ProclaimedBy:        "Organizzazione Sindacale",
		StrikeDate:          "2026-09-01",
		DeclarationDeadline: pastDeadline,
		IsPublished:         true,
	}
	_ = repo.CreateNotice(reqCreate.Context(), &expiredNotice)

	reqExpiredDec, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices/notice-expired-1/declare", bytes.NewBufferString(`{"intention":"participates"}`))
	reqExpiredDec.Header.Set("Content-Type", "application/json")
	reqExpiredDec.Header.Set("X-User-ID", "teacher-user-1")
	reqExpiredDec.Header.Set("X-User-Role", "teacher")
	reqExpiredDec.Header.Set("X-School-ID", schoolID)
	wExpiredDec := httptest.NewRecorder()
	r.ServeHTTP(wExpiredDec, reqExpiredDec)
	assert.Equal(t, http.StatusForbidden, wExpiredDec.Code, "Submitting to expired notice must return 403 Forbidden")
	var errResp map[string]interface{}
	err = json.Unmarshal(wExpiredDec.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "DEADLINE_PASSED", errResp["code"])
}
