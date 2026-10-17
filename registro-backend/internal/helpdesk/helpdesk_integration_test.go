package helpdesk

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

type mockHelpDeskRepo struct {
	slots    map[string]*HelpDeskSlot
	bookings map[string]*HelpDeskBooking
}

func newMockHelpDeskRepo() *mockHelpDeskRepo {
	return &mockHelpDeskRepo{
		slots:    make(map[string]*HelpDeskSlot),
		bookings: make(map[string]*HelpDeskBooking),
	}
}

func (m *mockHelpDeskRepo) CreateSlot(ctx context.Context, slot *HelpDeskSlot) error {
	slot.ID = "slot-test-1"
	slot.CreatedAt = time.Now()
	m.slots[slot.ID] = slot
	return nil
}

func (m *mockHelpDeskRepo) GetSlot(ctx context.Context, id string) (*HelpDeskSlot, error) {
	if s, ok := m.slots[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (m *mockHelpDeskRepo) ListSlots(ctx context.Context, schoolID, subjectID, date, status string) ([]HelpDeskSlot, error) {
	var list []HelpDeskSlot
	for _, s := range m.slots {
		list = append(list, *s)
	}
	return list, nil
}

func (m *mockHelpDeskRepo) BookSlot(ctx context.Context, booking *HelpDeskBooking) error {
	booking.ID = "book-test-1"
	booking.BookedAt = time.Now()
	booking.Status = "booked"
	m.bookings[booking.ID] = booking
	return nil
}

func (m *mockHelpDeskRepo) ListBookings(ctx context.Context, slotID string) ([]HelpDeskBooking, error) {
	var list []HelpDeskBooking
	for _, b := range m.bookings {
		if b.SlotID == slotID {
			list = append(list, *b)
		}
	}
	return list, nil
}

func (m *mockHelpDeskRepo) MarkAttendance(ctx context.Context, bookingID, status string) error {
	if b, ok := m.bookings[bookingID]; ok {
		b.Status = status
		now := time.Now()
		b.AttendedAt = &now
	}
	return nil
}

func (m *mockHelpDeskRepo) CompleteSlot(ctx context.Context, slotID string) error {
	if s, ok := m.slots[slotID]; ok {
		s.Status = "completed"
	}
	return nil
}

func (m *mockHelpDeskRepo) GetFISReport(ctx context.Context, schoolID string) ([]FISAccountingSummary, error) {
	return []FISAccountingSummary{
		{
			TeacherID:      "teacher-1",
			TeacherName:    "Prof. Bianchi",
			CompletedSlots: 2,
			TotalHours:     3.0,
			AttendedCount:  5,
		},
	}, nil
}

func setupHelpDeskRouter(repo Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		c.Set("user_id", "teacher-1")
		c.Next()
	})

	svc := NewService(repo)
	handler := NewHandler(svc)
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)

	return r
}

func TestHelpDeskWorkflow(t *testing.T) {
	repo := newMockHelpDeskRepo()
	router := setupHelpDeskRouter(repo)

	// Step 1: Create Slot
	slotPayload := map[string]interface{}{
		"subject_id":   "subj-math",
		"slot_date":    "2026-10-15",
		"start_time":   "15:00",
		"end_time":     "16:30",
		"max_capacity": 4,
	}
	body, _ := json.Marshal(slotPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/help-desk/slots", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on create slot, got %d: %s", w.Code, w.Body.String())
	}

	// Step 2: Book Slot
	bookPayload := map[string]interface{}{
		"topic_description": "Dubbi su integrali per parti",
	}
	body, _ = json.Marshal(bookPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/help-desk/slots/slot-test-1/book", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on book slot, got %d: %s", w.Code, w.Body.String())
	}

	// Step 3: Complete Slot
	req = httptest.NewRequest(http.MethodPut, "/api/v1/help-desk/slots/slot-test-1/complete", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on complete slot, got %d", w.Code)
	}

	// Step 4: Get FIS Accounting Report
	req = httptest.NewRequest(http.MethodGet, "/api/v1/help-desk/fis-report", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on fis report, got %d", w.Code)
	}
	var res struct {
		Data []FISAccountingSummary `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	if len(res.Data) != 1 || res.Data[0].TotalHours != 3.0 {
		t.Errorf("unexpected FIS report data: %+v", res)
	}
}
