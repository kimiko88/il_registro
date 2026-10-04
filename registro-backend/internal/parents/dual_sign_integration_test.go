package parents

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupDualSignTestRouter(parentID string, svc *Service, ds *DualSignatureService, ag *AccessGuard) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/api/v1")
	rg.Use(func(c *gin.Context) {
		c.Set("user_id", parentID)
		c.Set("role", "parent")
		c.Next()
	})
	h := NewHandler(svc)
	h.SetDualSignatureService(ds)
	h.SetAccessGuard(ag)
	h.RegisterRoutes(rg)
	return r
}

func TestDualSign_Integration_HTTPWorkflow(t *testing.T) {
	mockRepo := newMockDualSignRepo()
	ds := NewDualSignatureService(mockRepo)
	ag := NewAccessGuard(mockRepo)
	svc := &Service{}

	// 1. Parent 1 creates authorization request
	routerP1 := setupDualSignTestRouter("parent-user-1", svc, ds, ag)

	createPayload := CreateDualAuthParams{
		StudentID:     "student-1",
		DocumentType:  "trip_consent",
		DocumentRefID: "trip-999",
		Title:         "Gita d'Istruzione a Roma",
		Parent1ID:     "parent-user-1",
		Parent2ID:     "parent-user-2",
		IsShared:      true,
		Deadline:      time.Now().Add(48 * time.Hour),
	}
	body, _ := json.Marshal(createPayload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/parents/dual-authorizations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routerP1.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 on create dual auth, got %d: %s", w.Code, w.Body.String())
	}
	var created DualParentalAuthorization
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	if created.Status != StatusPendingFirst {
		t.Fatalf("expected pending_first, got %s", created.Status)
	}

	// 2. Parent 1 signs with PIN
	signBody1, _ := json.Marshal(SignDualAuthRequest{PIN: "1234"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/parents/dual-authorizations/"+created.ID+"/sign", bytes.NewReader(signBody1))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerP1.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on parent 1 sign, got %d: %s", w.Code, w.Body.String())
	}
	var signedP1 DualParentalAuthorization
	_ = json.Unmarshal(w.Body.Bytes(), &signedP1)
	if signedP1.Status != StatusPendingSecond {
		t.Fatalf("expected pending_second, got %s", signedP1.Status)
	}

	// 3. Parent 2 accesses system and signs with PIN -> Document completed
	routerP2 := setupDualSignTestRouter("parent-user-2", svc, ds, ag)

	signBody2, _ := json.Marshal(SignDualAuthRequest{PIN: "5678"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/parents/dual-authorizations/"+created.ID+"/sign", bytes.NewReader(signBody2))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerP2.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on parent 2 sign, got %d: %s", w.Code, w.Body.String())
	}
	var completed DualParentalAuthorization
	_ = json.Unmarshal(w.Body.Bytes(), &completed)
	if completed.Status != StatusCompleted {
		t.Fatalf("expected completed, got %s", completed.Status)
	}

	// 4. Restricted parent test: blocked by access guard
	routerRestricted := setupDualSignTestRouter("parent-restricted", svc, ds, ag)
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/parents/child/student-1/custody", nil)
	w = httptest.NewRecorder()
	routerRestricted.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for restricted parent, got %d", w.Code)
	}
}
