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

	"registro-backend/internal/trips"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockTripsRepo struct {
	tripsMap map[string]*trips.EducationalTrip
	consents map[string][]*trips.TripConsent
}

func (m *mockTripsRepo) CreateTrip(ctx context.Context, t *trips.EducationalTrip) error {
	if t.ID == "" {
		t.ID = "trip-1"
	}
	t.CreatedAt = time.Now()
	m.tripsMap[t.ID] = t
	return nil
}

func (m *mockTripsRepo) GetTripByID(ctx context.Context, id string) (*trips.EducationalTrip, error) {
	t, ok := m.tripsMap[id]
	if !ok {
		return nil, errors.New("trip not found")
	}
	return t, nil
}

func (m *mockTripsRepo) ListTrips(ctx context.Context, schoolID, studentID string) ([]*trips.EducationalTrip, error) {
	res := make([]*trips.EducationalTrip, 0)
	for _, t := range m.tripsMap {
		res = append(res, t)
	}
	return res, nil
}

func (m *mockTripsRepo) SubmitConsent(ctx context.Context, c *trips.TripConsent) error {
	if c.ID == "" {
		c.ID = "consent-1"
	}
	c.SignedAt = time.Now()
	m.consents[c.TripID] = append(m.consents[c.TripID], c)
	return nil
}

func (m *mockTripsRepo) ListConsents(ctx context.Context, tripID string) ([]*trips.TripConsent, error) {
	return m.consents[tripID], nil
}

func TestIntegration_Trips_Educational_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockTripsRepo{
		tripsMap: make(map[string]*trips.EducationalTrip),
		consents: make(map[string][]*trips.TripConsent),
	}
	svc := trips.NewService(repo)
	handler := trips.NewHandler(svc)

	r := gin.New()
	r.POST("/trips", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.CreateTrip(c)
	})
	r.POST("/trips/consent", func(c *gin.Context) {
		c.Set("user_id", "parent-1")
		c.Set("role", "parent")
		c.Set("school_id", "school-1")
		handler.SubmitConsent(c)
	})

	// 1. Teacher creates educational trip
	tripReq := trips.CreateTripRequest{
		Title:                "Visita Didattica Museo della Scienza",
		Destination:          "Milano",
		DepartureDate:        "2026-10-20",
		ReturnDate:           "2026-10-20",
		Description:          "Uscita didattica per laboratorio di fisica e chimica",
		AccompanyingTeachers: "Prof. Rossi, Prof. Bianchi",
		ClassIDs:             []string{"class-1"},
	}
	body1, _ := json.Marshal(tripReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/trips", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Parent signs digital consent
	consentReq := trips.SubmitConsentRequest{
		TripID:    "trip-1",
		StudentID: "st-1",
		Status:    "granted",
	}
	body2, _ := json.Marshal(consentReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/trips/consent", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
