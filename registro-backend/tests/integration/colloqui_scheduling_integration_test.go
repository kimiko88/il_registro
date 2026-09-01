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

	"registro-backend/internal/scheduling"
	"registro-backend/internal/teachers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Teachers Repository ──────────────────────────────────────────────

type mockTeacherRepoForScheduling struct{}

func (m *mockTeacherRepoForScheduling) List(_ context.Context, _ string) ([]teachers.Teacher, error) {
	return nil, nil
}
func (m *mockTeacherRepoForScheduling) Get(_ context.Context, id string) (*teachers.Teacher, error) {
	return &teachers.Teacher{ID: id, UserID: id, SchoolID: "school-1"}, nil
}
func (m *mockTeacherRepoForScheduling) GetByUserID(_ context.Context, userID string) (*teachers.Teacher, error) {
	return &teachers.Teacher{ID: "prof-" + userID, UserID: userID, SchoolID: "school-1"}, nil
}
func (m *mockTeacherRepoForScheduling) GetBySubject(_ context.Context, _ string) ([]teachers.Teacher, error) {
	return nil, nil
}
func (m *mockTeacherRepoForScheduling) GetSubjects(_ context.Context, _ string) ([]teachers.TeacherSubject, error) {
	return nil, nil
}
func (m *mockTeacherRepoForScheduling) AssignSubject(_ context.Context, _, _ string) error {
	return nil
}
func (m *mockTeacherRepoForScheduling) RemoveSubject(_ context.Context, _, _ string) error {
	return nil
}
func (m *mockTeacherRepoForScheduling) GetDashboardStats(_ context.Context, _ string) (map[string]interface{}, error) {
	return nil, nil
}

// ─── Mock Scheduling Repository ───────────────────────────────────────────

type mockSchedulingRepo struct {
	slots    map[string]*scheduling.ColloquioSlot
	bookings map[string]*scheduling.ColloquioBooking
}

func newMockSchedulingRepo() *mockSchedulingRepo {
	return &mockSchedulingRepo{
		slots:    make(map[string]*scheduling.ColloquioSlot),
		bookings: make(map[string]*scheduling.ColloquioBooking),
	}
}

func (m *mockSchedulingRepo) CreateSlot(_ context.Context, slot *scheduling.ColloquioSlot) error {
	if slot.ID == "" {
		slot.ID = "slot-" + time.Now().Format("150405.000")
	}
	m.slots[slot.ID] = slot
	return nil
}

func (m *mockSchedulingRepo) GetSlots(_ context.Context, teacherID string, _, _ time.Time) ([]scheduling.ColloquioSlot, error) {
	var list []scheduling.ColloquioSlot
	for _, s := range m.slots {
		if teacherID != "" && s.TeacherID != teacherID {
			continue
		}
		list = append(list, *s)
	}
	return list, nil
}

func (m *mockSchedulingRepo) GetAvailableSlots(_ context.Context, _, teacherID string, _, _ time.Time) ([]scheduling.ColloquioSlot, error) {
	var list []scheduling.ColloquioSlot
	for _, s := range m.slots {
		if teacherID != "" && s.TeacherID != teacherID {
			continue
		}
		if !s.IsCancelled && s.BookingCount < s.MaxBookings {
			list = append(list, *s)
		}
	}
	return list, nil
}

func (m *mockSchedulingRepo) GetSlotByID(_ context.Context, id string) (*scheduling.ColloquioSlot, error) {
	s, ok := m.slots[id]
	if !ok {
		return nil, errors.New("slot not found")
	}
	return s, nil
}

func (m *mockSchedulingRepo) UpdateSlot(_ context.Context, slot *scheduling.ColloquioSlot) error {
	m.slots[slot.ID] = slot
	return nil
}

func (m *mockSchedulingRepo) DeleteSlot(_ context.Context, id string) error {
	delete(m.slots, id)
	return nil
}

func (m *mockSchedulingRepo) CreateSlotsBatch(_ context.Context, slots []scheduling.ColloquioSlot) error {
	for i := range slots {
		_ = m.CreateSlot(context.Background(), &slots[i])
	}
	return nil
}

func (m *mockSchedulingRepo) CreateBooking(_ context.Context, b *scheduling.ColloquioBooking) error {
	if b.ID == "" {
		b.ID = "booking-" + time.Now().Format("150405.000")
	}
	b.BookedAt = time.Now()
	b.Status = scheduling.StatusConfirmed
	m.bookings[b.ID] = b
	return nil
}

func (m *mockSchedulingRepo) GetBooking(_ context.Context, id string) (*scheduling.ColloquioBooking, error) {
	b, ok := m.bookings[id]
	if !ok {
		return nil, errors.New("booking not found")
	}
	return b, nil
}

func (m *mockSchedulingRepo) GetBookingsByParent(_ context.Context, parentID string) ([]scheduling.ColloquioBooking, error) {
	var list []scheduling.ColloquioBooking
	for _, b := range m.bookings {
		if b.ParentID != nil && *b.ParentID == parentID {
			list = append(list, *b)
		}
	}
	return list, nil
}

func (m *mockSchedulingRepo) GetBookingsByTeacher(_ context.Context, teacherID string) ([]scheduling.ColloquioBooking, error) {
	var list []scheduling.ColloquioBooking
	for _, b := range m.bookings {
		slot := m.slots[b.SlotID]
		if slot != nil && slot.TeacherID == teacherID {
			list = append(list, *b)
		}
	}
	return list, nil
}

func (m *mockSchedulingRepo) UpdateBooking(_ context.Context, booking *scheduling.ColloquioBooking) error {
	m.bookings[booking.ID] = booking
	return nil
}

func (m *mockSchedulingRepo) CountBookingsForParent(_ context.Context, _ string, _ time.Time, _, _ time.Time) (int, error) {
	return 0, nil
}

func (m *mockSchedulingRepo) GetSettings(_ context.Context, _ string) (*scheduling.ColloquioSettings, error) {
	return &scheduling.ColloquioSettings{BookingWindowDays: 14, BookingBufferHours: 2}, nil
}

func (m *mockSchedulingRepo) UpdateSettings(_ context.Context, _ *scheduling.ColloquioSettings) error {
	return nil
}

func (m *mockSchedulingRepo) GetAnalytics(_ context.Context, _ string) (*scheduling.AnalyticsResponse, error) {
	return &scheduling.AnalyticsResponse{TotalSlots: len(m.slots), UtilizationRate: 0.5}, nil
}

func (m *mockSchedulingRepo) ResolveParentUserID(_ context.Context, userID string) (string, error) {
	return userID, nil
}

func (m *mockSchedulingRepo) IsGuardian(_ context.Context, _, _ string) (bool, error) {
	return true, nil
}

// ─── Setup Router Helper ───────────────────────────────────────────────────

func setupSchedulingRouter(repo *mockSchedulingRepo, role, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	tRepo := &mockTeacherRepoForScheduling{}
	svc := scheduling.NewService(repo, tRepo, nil, nil, nil)
	h := scheduling.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

// ─── Tests ─────────────────────────────────────────────────────────────────

// SCH01 — Teacher creates scheduling slots
func TestScheduling_CreateSlot_TeacherSuccess(t *testing.T) {
	repo := newMockSchedulingRepo()
	r := setupSchedulingRouter(repo, "teacher", "teacher-1", "school-1")

	futureDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	body, _ := json.Marshal(scheduling.CreateSlotRequest{
		Dates:       []string{futureDate},
		StartTime:   "15:00",
		EndTime:     "16:00",
		Duration:    15,
		MaxBookings: 1,
		Type:        scheduling.SlotIndividual,
		Location:    "Aula Ricevimento / Google Meet",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/scheduling/slots", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "slot created")
}

// SCH02 — Student cannot create scheduling slots
func TestScheduling_CreateSlot_StudentForbidden(t *testing.T) {
	repo := newMockSchedulingRepo()
	r := setupSchedulingRouter(repo, "student", "student-1", "school-1")

	futureDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	body, _ := json.Marshal(scheduling.CreateSlotRequest{
		Dates:     []string{futureDate},
		StartTime: "15:00",
		EndTime:   "16:00",
		Type:      scheduling.SlotIndividual,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/scheduling/slots", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// SCH03 — Teacher retrieves own slots
func TestScheduling_GetMySlots_Success(t *testing.T) {
	repo := newMockSchedulingRepo()
	start, _ := time.Parse("15:04", "10:00")
	end, _ := time.Parse("15:04", "10:30")
	_ = repo.CreateSlot(context.Background(), &scheduling.ColloquioSlot{
		ID:          "slot-101",
		TeacherID:   "prof-teacher-1",
		Date:        time.Now().AddDate(0, 0, 2),
		StartTime:   start,
		EndTime:     end,
		MaxBookings: 1,
	})

	r := setupSchedulingRouter(repo, "teacher", "teacher-1", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/scheduling/slots/my", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var slots []scheduling.SlotResponse
	err := json.Unmarshal(w.Body.Bytes(), &slots)
	require.NoError(t, err)
	assert.Len(t, slots, 1)
}

// SCH04 — Parent books an available slot
func TestScheduling_BookSlot_Success(t *testing.T) {
	repo := newMockSchedulingRepo()
	start, _ := time.Parse("15:04", "15:00")
	end, _ := time.Parse("15:04", "15:30")
	slotDate := time.Now().AddDate(0, 0, 3)
	_ = repo.CreateSlot(context.Background(), &scheduling.ColloquioSlot{
		ID:           "slot-parent-book",
		TeacherID:    "teacher-1",
		Date:         slotDate,
		StartTime:    start,
		EndTime:      end,
		MaxBookings:  1,
		BookingCount: 0,
	})

	r := setupSchedulingRouter(repo, "parent", "parent-1", "school-1")

	stdID := "student-child-1"
	body, _ := json.Marshal(scheduling.BookSlotRequest{
		SlotID:    "slot-parent-book",
		StudentID: &stdID,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/scheduling/book", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var resp scheduling.BookingResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, scheduling.StatusConfirmed, resp.Status)
}

// SCH05 — Parent cancels a booking
func TestScheduling_CancelBooking_Success(t *testing.T) {
	repo := newMockSchedulingRepo()
	slotDate := time.Now().AddDate(0, 0, 3)
	_ = repo.CreateSlot(context.Background(), &scheduling.ColloquioSlot{
		ID:           "slot-to-cancel",
		TeacherID:    "teacher-1",
		Date:         slotDate,
		MaxBookings:  1,
		BookingCount: 1,
	})
	parentID := "parent-1"
	_ = repo.CreateBooking(context.Background(), &scheduling.ColloquioBooking{
		ID:       "booking-to-cancel",
		SlotID:   "slot-to-cancel",
		ParentID: &parentID,
		Status:   scheduling.StatusConfirmed,
	})

	r := setupSchedulingRouter(repo, "parent", "parent-1", "school-1")
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/scheduling/bookings/booking-to-cancel/cancel", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	booking, _ := repo.GetBooking(context.Background(), "booking-to-cancel")
	assert.Equal(t, scheduling.StatusCancelled, booking.Status)
}
