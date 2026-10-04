package protocol

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type mockProtocolRepo struct {
	entries     map[string]*ProtocolEntry
	entityLinks map[string]string // entityType_entityId -> protocolId
	counter     int
}

func newMockProtocolRepo() *mockProtocolRepo {
	return &mockProtocolRepo{
		entries:     make(map[string]*ProtocolEntry),
		entityLinks: make(map[string]string),
	}
}

func (m *mockProtocolRepo) RegisterDocument(ctx context.Context, entry *ProtocolEntry, entityType, entityID string) error {
	m.counter++
	entry.ID = "prot-id-1"
	entry.ProtocolYear = 2026
	entry.ProtocolNumber = m.counter
	entry.ProtocolDate = time.Now()
	m.entries[entry.ID] = entry
	if entityType != "" && entityID != "" {
		m.entityLinks[entityType+"_"+entityID] = entry.ID
	}
	return nil
}

func (m *mockProtocolRepo) GetProtocolEntry(ctx context.Context, id string) (*ProtocolEntry, error) {
	if e, ok := m.entries[id]; ok {
		return e, nil
	}
	return nil, nil
}

func (m *mockProtocolRepo) ListProtocolEntries(ctx context.Context, schoolID string, year int, flowDirection string) ([]ProtocolEntry, error) {
	var list []ProtocolEntry
	for _, e := range m.entries {
		list = append(list, *e)
	}
	return list, nil
}

func (m *mockProtocolRepo) GetEntityProtocol(ctx context.Context, entityType, entityID string) (*ProtocolEntry, error) {
	if pID, ok := m.entityLinks[entityType+"_"+entityID]; ok {
		return m.entries[pID], nil
	}
	return nil, nil
}

func setupProtocolRouter(repo Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "secretary")
		c.Set("school_id", "school-1")
		c.Set("user_id", "staff-1")
		c.Next()
	})

	svc := NewService(repo)
	handler := NewHandler(svc)
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)

	return r
}

func TestProtocolIntegrationWorkflow(t *testing.T) {
	repo := newMockProtocolRepo()
	router := setupProtocolRouter(repo)

	// Step 1: Protocol document
	payload := map[string]interface{}{
		"subject":              "Circolare n. 42 - Viaggio di Istruzione a Firenze",
		"sender":               "Dirigente Scolastico",
		"recipient":            "Famiglie e Studenti",
		"flow_direction":       "out",
		"classification_title": 4, // Titolo IV - Didattica
		"classification_class": "2",
		"entity_type":          "circular",
		"entity_id":            "00000000-0000-0000-0000-000000000042",
		"school_name":          "Liceo Statale Galvani",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/protocol", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on protocol document, got %d: %s", w.Code, w.Body.String())
	}

	var res struct {
		Protocol ProtocolEntry `json:"protocol"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	if res.Protocol.ProtocolNumber != 1 {
		t.Errorf("expected protocol number 1, got %d", res.Protocol.ProtocolNumber)
	}
	if res.Protocol.VisualStamp == "" {
		t.Errorf("expected visual stamp to be populated")
	}

	// Step 2: Download Segnatura.xml
	reqXML := httptest.NewRequest(http.MethodGet, "/api/v1/protocol/prot-id-1/segnatura.xml", nil)
	wXML := httptest.NewRecorder()
	router.ServeHTTP(wXML, reqXML)

	if wXML.Code != http.StatusOK {
		t.Fatalf("expected 200 on Segnatura.xml download, got %d", wXML.Code)
	}
	if wXML.Header().Get("Content-Type") != "application/xml" {
		t.Errorf("expected application/xml content type")
	}

	// Step 3: Retrieve by linked entity
	reqEnt := httptest.NewRequest(http.MethodGet, "/api/v1/protocol/by-entity/circular/00000000-0000-0000-0000-000000000042", nil)
	wEnt := httptest.NewRecorder()
	router.ServeHTTP(wEnt, reqEnt)

	if wEnt.Code != http.StatusOK {
		t.Fatalf("expected 200 on get protocol by entity, got %d", wEnt.Code)
	}
}
