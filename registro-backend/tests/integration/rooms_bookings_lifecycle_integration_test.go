package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/rooms"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockIntegrationRoomsRepo struct {
	buildings []rooms.SchoolBuilding
	rooms     []rooms.BookableRoom
	bookings  []rooms.RoomBooking
}

func newMockIntegrationRoomsRepo() *mockIntegrationRoomsRepo {
	return &mockIntegrationRoomsRepo{
		buildings: make([]rooms.SchoolBuilding, 0),
		rooms:     make([]rooms.BookableRoom, 0),
		bookings:  make([]rooms.RoomBooking, 0),
	}
}

func (m *mockIntegrationRoomsRepo) CreateBuilding(ctx context.Context, b *rooms.SchoolBuilding) (*rooms.SchoolBuilding, error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	m.buildings = append(m.buildings, *b)
	return b, nil
}

func (m *mockIntegrationRoomsRepo) GetBuildingByID(ctx context.Context, id string) (*rooms.SchoolBuilding, error) {
	for _, b := range m.buildings {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, rooms.ErrNotFound
}

func (m *mockIntegrationRoomsRepo) ListBuildings(ctx context.Context, schoolID string) ([]rooms.SchoolBuilding, error) {
	var res []rooms.SchoolBuilding
	for _, b := range m.buildings {
		if b.SchoolID == schoolID {
			res = append(res, b)
		}
	}
	return res, nil
}

func (m *mockIntegrationRoomsRepo) UpdateBuilding(ctx context.Context, b *rooms.SchoolBuilding) (*rooms.SchoolBuilding, error) {
	for i, existing := range m.buildings {
		if existing.ID == b.ID {
			m.buildings[i] = *b
			return b, nil
		}
	}
	return nil, rooms.ErrNotFound
}

func (m *mockIntegrationRoomsRepo) DeleteBuilding(ctx context.Context, id string) error {
	var filtered []rooms.SchoolBuilding
	for _, b := range m.buildings {
		if b.ID != id {
			filtered = append(filtered, b)
		}
	}
	m.buildings = filtered
	return nil
}

func (m *mockIntegrationRoomsRepo) CreateRoom(ctx context.Context, r *rooms.BookableRoom) (*rooms.BookableRoom, error) {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	m.rooms = append(m.rooms, *r)
	return r, nil
}

func (m *mockIntegrationRoomsRepo) GetRoomByID(ctx context.Context, id string) (*rooms.BookableRoom, error) {
	for _, r := range m.rooms {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, rooms.ErrNotFound
}

func (m *mockIntegrationRoomsRepo) ListRooms(ctx context.Context, schoolID string, buildingID *string, roomType *string, activeOnly bool) ([]rooms.BookableRoom, error) {
	var res []rooms.BookableRoom
	for _, r := range m.rooms {
		if r.SchoolID != schoolID {
			continue
		}
		if buildingID != nil && (r.BuildingID == nil || *r.BuildingID != *buildingID) {
			continue
		}
		if roomType != nil && r.RoomType != *roomType {
			continue
		}
		if activeOnly && !r.IsActive {
			continue
		}
		res = append(res, r)
	}
	return res, nil
}

func (m *mockIntegrationRoomsRepo) UpdateRoom(ctx context.Context, r *rooms.BookableRoom) (*rooms.BookableRoom, error) {
	for i, existing := range m.rooms {
		if existing.ID == r.ID {
			m.rooms[i] = *r
			return r, nil
		}
	}
	return nil, rooms.ErrNotFound
}

func (m *mockIntegrationRoomsRepo) DeleteRoom(ctx context.Context, id string) error {
	var filtered []rooms.BookableRoom
	for _, r := range m.rooms {
		if r.ID != id {
			filtered = append(filtered, r)
		}
	}
	m.rooms = filtered
	return nil
}

func (m *mockIntegrationRoomsRepo) CheckSlotConflict(ctx context.Context, roomID string, date string, hourIndex int, excludeBookingID *string) (bool, error) {
	for _, b := range m.bookings {
		if excludeBookingID != nil && b.ID == *excludeBookingID {
			continue
		}
		if b.RoomID == roomID && b.BookingDate == date && b.HourIndex == hourIndex && b.Status == rooms.BookingStatusConfirmed {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockIntegrationRoomsRepo) CreateBooking(ctx context.Context, b *rooms.RoomBooking) (*rooms.RoomBooking, error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	b.Status = rooms.BookingStatusConfirmed
	m.bookings = append(m.bookings, *b)
	return b, nil
}

func (m *mockIntegrationRoomsRepo) CreateBookingBatch(ctx context.Context, bookings []rooms.RoomBooking) ([]rooms.RoomBooking, error) {
	var created []rooms.RoomBooking
	for _, b := range bookings {
		if b.ID == "" {
			b.ID = uuid.New().String()
		}
		b.CreatedAt = time.Now()
		b.UpdatedAt = time.Now()
		b.Status = rooms.BookingStatusConfirmed
		m.bookings = append(m.bookings, b)
		created = append(created, b)
	}
	return created, nil
}

func (m *mockIntegrationRoomsRepo) GetBookingByID(ctx context.Context, id string) (*rooms.RoomBooking, error) {
	for _, b := range m.bookings {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, rooms.ErrNotFound
}

func (m *mockIntegrationRoomsRepo) ListBookings(ctx context.Context, filter rooms.BookingFilter) ([]rooms.RoomBooking, error) {
	var res []rooms.RoomBooking
	for _, b := range m.bookings {
		if filter.SchoolID != "" && b.SchoolID != filter.SchoolID {
			continue
		}
		if filter.RoomID != nil && b.RoomID != *filter.RoomID {
			continue
		}
		if filter.TeacherID != nil && b.TeacherID != *filter.TeacherID {
			continue
		}
		res = append(res, b)
	}
	return res, nil
}

func (m *mockIntegrationRoomsRepo) CancelBooking(ctx context.Context, id string) error {
	for i, b := range m.bookings {
		if b.ID == id {
			m.bookings[i].Status = rooms.BookingStatusCancelled
			return nil
		}
	}
	return rooms.ErrNotFound
}

func (m *mockIntegrationRoomsRepo) CancelRecurringSeries(ctx context.Context, parentBookingID string) error {
	for i, b := range m.bookings {
		if b.ID == parentBookingID || (b.ParentBookingID != nil && *b.ParentBookingID == parentBookingID) {
			m.bookings[i].Status = rooms.BookingStatusCancelled
		}
	}
	return nil
}

func (m *mockIntegrationRoomsRepo) GetRoomBookingsForDateRange(ctx context.Context, roomID string, fromDate, toDate string) ([]rooms.RoomBooking, error) {
	var res []rooms.RoomBooking
	for _, b := range m.bookings {
		if b.RoomID == roomID && b.Status == rooms.BookingStatusConfirmed {
			res = append(res, b)
		}
	}
	return res, nil
}

func setupRoomsIntegrationRouter(repo rooms.Repository, currentRole, currentUserID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := rooms.NewService(repo)
	handler := rooms.NewHandler(svc)

	r.Use(func(c *gin.Context) {
		c.Set("role", currentRole)
		c.Set("user_id", currentUserID)
		c.Set("school_id", schoolID)
		c.Next()
	})

	api := r.Group("/rooms")
	{
		api.GET("/buildings", handler.ListBuildings)
		api.POST("/buildings", handler.CreateBuilding)
		api.PUT("/buildings/:id", handler.UpdateBuilding)
		api.DELETE("/buildings/:id", handler.DeleteBuilding)

		api.GET("", handler.ListRooms)
		api.POST("", handler.CreateRoom)
		api.GET("/:id/availability", handler.GetRoomAvailability)

		api.GET("/bookings", handler.ListBookings)
		api.POST("/bookings", handler.CreateBooking)
		api.DELETE("/bookings/:id", handler.CancelBooking)
	}

	return r
}

func TestIntegration_RoomsAndBookingsLifecycle(t *testing.T) {
	repo := newMockIntegrationRoomsRepo()
	schoolID := "school-test-1"
	secretaryID := "sec-1"
	teacherID := "teacher-1"

	// 1. Secretary creates a Building (Plesso)
	routerSec := setupRoomsIntegrationRouter(repo, "secretary", secretaryID, schoolID)

	bldReq := rooms.CreateBuildingRequest{
		Name:    "Plesso Fermi",
		Address: func(s string) *string { return &s }("Via Roma 10, Milano"),
	}
	body, _ := json.Marshal(bldReq)
	req := httptest.NewRequest(http.MethodPost, "/rooms/buildings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routerSec.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdBld rooms.SchoolBuilding
	err := json.Unmarshal(w.Body.Bytes(), &createdBld)
	assert.NoError(t, err)
	assert.Equal(t, "Plesso Fermi", createdBld.Name)

	// 2. Secretary creates a Bookable Computer Lab in this building
	roomReq := rooms.CreateRoomRequest{
		BuildingID: &createdBld.ID,
		Name:       "Laboratorio Multimediale",
		RoomType:   rooms.RoomTypeLabInformatica,
		Capacity:   30,
		Equipment:  []string{"30 PC", "LIM", "Fibra 1Gbps"},
	}
	body, _ = json.Marshal(roomReq)
	req = httptest.NewRequest(http.MethodPost, "/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerSec.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdRoom rooms.BookableRoom
	err = json.Unmarshal(w.Body.Bytes(), &createdRoom)
	assert.NoError(t, err)
	assert.Equal(t, rooms.RoomTypeLabInformatica, createdRoom.RoomType)

	// 3. Teacher lists rooms filtered by building
	routerTeacher := setupRoomsIntegrationRouter(repo, "teacher", teacherID, schoolID)
	req = httptest.NewRequest(http.MethodGet, "/rooms?building_id="+createdBld.ID, nil)
	w = httptest.NewRecorder()
	routerTeacher.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var roomsList []rooms.BookableRoom
	err = json.Unmarshal(w.Body.Bytes(), &roomsList)
	assert.NoError(t, err)
	assert.Len(t, roomsList, 1)
	assert.Equal(t, "Laboratorio Multimediale", roomsList[0].Name)

	// 4. Teacher books a Spot Slot (2026-10-15 hour 2)
	bookingReq := rooms.CreateBookingRequest{
		RoomID:      createdRoom.ID,
		ClassID:     func(s string) *string { return &s }("class-3A"),
		BookingDate: "2026-10-15",
		HourIndex:   2,
		Notes:       func(s string) *string { return &s }("Lezione coding Python"),
	}
	body, _ = json.Marshal(bookingReq)
	req = httptest.NewRequest(http.MethodPost, "/rooms/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerTeacher.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var spotBooking rooms.RoomBooking
	err = json.Unmarshal(w.Body.Bytes(), &spotBooking)
	assert.NoError(t, err)
	assert.Equal(t, 2, spotBooking.HourIndex)

	// 5. Another teacher tries to book the SAME slot -> Conflict 409
	routerTeacher2 := setupRoomsIntegrationRouter(repo, "teacher", "teacher-2", schoolID)
	conflictReq := rooms.CreateBookingRequest{
		RoomID:      createdRoom.ID,
		ClassID:     func(s string) *string { return &s }("class-4B"),
		BookingDate: "2026-10-15",
		HourIndex:   2,
	}
	body, _ = json.Marshal(conflictReq)
	req = httptest.NewRequest(http.MethodPost, "/rooms/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerTeacher2.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	// 6. Teacher books Recurring Weekly slots (2026-11-03 to 2026-11-17, 3 occurrences)
	recReq := rooms.CreateBookingRequest{
		RoomID:            createdRoom.ID,
		ClassID:           func(s string) *string { return &s }("class-3A"),
		BookingDate:       "2026-11-03",
		HourIndex:         4,
		IsRecurring:       true,
		RecurrencePattern: func(s string) *string { return &s }(rooms.RecurrenceWeekly),
		RecurringUntil:    func(s string) *string { return &s }("2026-11-17"),
	}
	body, _ = json.Marshal(recReq)
	req = httptest.NewRequest(http.MethodPost, "/rooms/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerTeacher.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var recResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &recResp)
	assert.NoError(t, err)
	occurrences, ok := recResp["occurrences"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, occurrences, 3)

	// 7. Teacher cancels the entire recurring series
	firstOcc := occurrences[0].(map[string]interface{})
	firstID := firstOcc["id"].(string)

	cancelURL := "/rooms/bookings/" + firstID + "?cancel_series=true"
	req = httptest.NewRequest(http.MethodDelete, cancelURL, nil)
	w = httptest.NewRecorder()
	routerTeacher.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify all 3 are cancelled in repo
	for _, item := range occurrences {
		occMap := item.(map[string]interface{})
		id := occMap["id"].(string)
		b, err := repo.GetBookingByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, rooms.BookingStatusCancelled, b.Status)
	}
}
