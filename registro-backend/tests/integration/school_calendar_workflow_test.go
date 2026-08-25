package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/schoolcalendar"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockSchoolCalendarRepo struct {
	year    *schoolcalendar.SchoolYearSettings
	nonDays map[string]*schoolcalendar.NonTeachingDay
	periods map[string]*schoolcalendar.AcademicPeriod
}

func (m *mockSchoolCalendarRepo) GetYear(schoolID string) (*schoolcalendar.SchoolYearSettings, error) {
	return m.year, nil
}

func (m *mockSchoolCalendarRepo) UpsertYear(s *schoolcalendar.SchoolYearSettings) error {
	m.year = s
	return nil
}

func (m *mockSchoolCalendarRepo) AddNonTeachingDay(d *schoolcalendar.NonTeachingDay) error {
	if d.ID == "" {
		d.ID = "ntd-1"
	}
	m.nonDays[d.ID] = d
	return nil
}

func (m *mockSchoolCalendarRepo) DeleteNonTeachingDay(schoolID, id string) error {
	delete(m.nonDays, id)
	return nil
}

func (m *mockSchoolCalendarRepo) ListNonTeachingDays(schoolID string) ([]schoolcalendar.NonTeachingDay, error) {
	res := make([]schoolcalendar.NonTeachingDay, 0)
	for _, n := range m.nonDays {
		res = append(res, *n)
	}
	return res, nil
}

func (m *mockSchoolCalendarRepo) CountTeachingDays(schoolID string, from, to time.Time) (int, error) {
	return 200, nil
}

func (m *mockSchoolCalendarRepo) CreateAcademicPeriod(p *schoolcalendar.AcademicPeriod) error {
	if p.ID == "" {
		p.ID = "period-1"
	}
	m.periods[p.ID] = p
	return nil
}

func (m *mockSchoolCalendarRepo) ListAcademicPeriods(schoolID string) ([]schoolcalendar.AcademicPeriod, error) {
	res := make([]schoolcalendar.AcademicPeriod, 0)
	for _, p := range m.periods {
		res = append(res, *p)
	}
	return res, nil
}

func TestIntegration_School_Calendar_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockSchoolCalendarRepo{
		nonDays: make(map[string]*schoolcalendar.NonTeachingDay),
		periods: make(map[string]*schoolcalendar.AcademicPeriod),
	}
	svc := schoolcalendar.NewService(repo)
	handler := schoolcalendar.NewHandler(svc)

	r := gin.New()
	r.PUT("/school-calendar/year", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.SetYear(c)
	})
	r.POST("/school-calendar/non-teaching-days", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.AddNonTeachingDay(c)
	})

	// 1. Admin configures Academic Year bounds (2026/2027)
	yearReq := schoolcalendar.SetYearRequest{
		YearLabel: "2026/2027",
		StartDate: "2026-09-12",
		EndDate:   "2027-06-08",
	}
	body1, _ := json.Marshal(yearReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("PUT", "/school-calendar/year", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 2. Admin adds regional holiday / non-teaching day
	holidayReq := schoolcalendar.AddNonTeachingDayRequest{
		Date:  "2026-11-01",
		Label: "Festa di Tutti i Santi",
	}
	body2, _ := json.Marshal(holidayReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/school-calendar/non-teaching-days", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)
}
