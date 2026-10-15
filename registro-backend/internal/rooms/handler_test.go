package rooms

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRoomsRouter(svc Service, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "school-1")
		c.Set("user_id", "user-teacher-1")
		c.Set("role", role)
		c.Next()
	})
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestRoomsHandler_RoleAccess(t *testing.T) {
	repo := newMockRoomsRepo()
	svc := NewService(repo)
	rStudent := setupRoomsRouter(svc, "student")

	// Student cannot create a room
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/rooms", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	rStudent.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden 403, got %d", w.Code)
	}

	// Admin can create building & room
	rAdmin := setupRoomsRouter(svc, "admin")
	body, _ := json.Marshal(CreateRoomRequest{
		Name:     "Lab Informatica A",
		RoomType: RoomTypeLabInformatica,
		Capacity: 30,
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rAdmin.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRoomsHandler_ListAndBook(t *testing.T) {
	repo := newMockRoomsRepo()
	svc := NewService(repo)
	rTeacher := setupRoomsRouter(svc, "teacher")

	// Seed room
	room := &BookableRoom{
		ID:              "room-101",
		SchoolID:        "school-1",
		Name:            "Laboratorio di Fisica",
		RoomType:        RoomTypeLabFisica,
		Capacity:        25,
		RequiresBooking: true,
		IsActive:        true,
	}
	repo.rooms = append(repo.rooms, *room)

	// List
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms", nil)
	rTeacher.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 ok, got %d", w.Code)
	}

	// Book room via /api/v1/rooms/bookings
	bookBody, _ := json.Marshal(CreateBookingRequest{
		RoomID:      "room-101",
		BookingDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		HourIndex:   3,
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/rooms/bookings", bytes.NewBuffer(bookBody))
	req.Header.Set("Content-Type", "application/json")
	rTeacher.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 created on booking, got %d: %s", w.Code, w.Body.String())
	}
}
