package colloqui

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestColloquiRouter(svc Service, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		if role != "" {
			c.Set("role", role)
		}
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestParseFlexibleDate(t *testing.T) {
	assert.True(t, parseFlexibleDate("").IsZero())
	assert.True(t, parseFlexibleDate("null").IsZero())
	assert.True(t, parseFlexibleDate("undefined").IsZero())

	d1 := parseFlexibleDate("2026-05-15")
	assert.Equal(t, 2026, d1.Year())
	assert.Equal(t, time.Month(5), d1.Month())
	assert.Equal(t, 15, d1.Day())

	d2 := parseFlexibleDate("2026-05-15T14:30:00Z")
	assert.Equal(t, 2026, d2.Year())

	assert.True(t, parseFlexibleDate("not-a-date").IsZero())
}

func TestColloquiHandler_Endpoints(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	// Teacher router
	rTeacher := setupTestColloquiRouter(svc, "teacher-1", "teacher", "sch-1")

	// 1. GET /colloqui/slots/my
	repo.On("ListSlots", mock.Anything, mock.Anything).Return([]*ColloquioSlot{{ID: "slot-1", TeacherID: "teacher-1"}}, nil).Once()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/colloqui/slots/my", nil)
	w := httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. GET /colloqui/slots (Available slots)
	repo.On("ListSlots", mock.Anything, mock.Anything).Return([]*ColloquioSlot{{ID: "slot-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/colloqui/slots?date_from=2026-05-01&date_to=2026-05-31", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. GET /colloqui/availability/:teacherID
	repo.On("ListSlots", mock.Anything, mock.Anything).Return([]*ColloquioSlot{{ID: "slot-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/colloqui/availability/teacher-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. PATCH /colloqui/slots/:id
	repo.On("GetSlotByID", mock.Anything, "slot-1").Return(&ColloquioSlot{ID: "slot-1", TeacherID: "teacher-1"}, nil).Once()
	repo.On("PatchSlot", mock.Anything, "slot-1", "10:30", "11:30").Return(nil).Once()
	patchPayload := `{"start_time":"10:30","end_time":"11:30"}`
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/colloqui/slots/slot-1", bytes.NewBufferString(patchPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. DELETE /colloqui/slots/:id (CancelSlot)
	repo.On("GetSlotByID", mock.Anything, "slot-1").Return(&ColloquioSlot{ID: "slot-1", TeacherID: "teacher-1"}, nil).Once()
	repo.On("CancelSlot", mock.Anything, "slot-1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/colloqui/slots/slot-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. GET /colloqui/slots/:id/bookings
	repo.On("GetSlotByID", mock.Anything, "slot-1").Return(&ColloquioSlot{ID: "slot-1", TeacherID: "teacher-1"}, nil).Once()
	repo.On("ListSlotBookings", mock.Anything, "slot-1").Return([]*ColloquioBooking{{ID: "b1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/colloqui/slots/slot-1/bookings", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. GET /colloqui/bookings (Parent router)
	rParent := setupTestColloquiRouter(svc, "parent-1", "parent", "sch-1")
	repo.On("ListUserBookings", mock.Anything, "parent-1", "sch-1").Return([]*ColloquioBooking{{ID: "b1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/colloqui/bookings", nil)
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 8. GET /colloqui/bookings/:id
	pID := "parent-1"
	repo.On("GetBookingByID", mock.Anything, "b1").Return(&ColloquioBooking{ID: "b1", ParentID: &pID}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/colloqui/bookings/b1", nil)
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 9. PUT /colloqui/bookings/:id/status
	repo.On("GetBookingByID", mock.Anything, "b1").Return(&ColloquioBooking{ID: "b1", ParentID: &pID, SlotID: "slot-1"}, nil).Once()
	repo.On("GetSlotByID", mock.Anything, "slot-1").Return(&ColloquioSlot{ID: "slot-1", TeacherID: "teacher-1"}, nil).Once()
	repo.On("UpdateBookingStatus", mock.Anything, "b1", StatusCancelled, "parent-1", "Impegno imprevisto").Return(nil).Once()
	statusPayload := `{"status":"Cancelled","reason":"Impegno imprevisto"}`
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/colloqui/bookings/b1/status", bytes.NewBufferString(statusPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 10. PATCH /colloqui/bookings/:id/cancel
	repo.On("GetBookingByID", mock.Anything, "b1").Return(&ColloquioBooking{ID: "b1", ParentID: &pID, SlotID: "slot-1"}, nil).Once()
	repo.On("GetSlotByID", mock.Anything, "slot-1").Return(&ColloquioSlot{ID: "slot-1", TeacherID: "teacher-1"}, nil).Once()
	repo.On("UpdateBookingStatus", mock.Anything, "b1", StatusCancelled, "parent-1", "").Return(nil).Once()
	cancelPayload := `{"reason":"Annullamento"}`
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/colloqui/bookings/b1/cancel", bytes.NewBufferString(cancelPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 11. General Meetings endpoints
	rAdmin := setupTestColloquiRouter(svc, "admin-1", "admin", "sch-1")

	// ListGeneralMeetings
	repo.On("ListGeneralMeetings", mock.Anything, "sch-1").Return([]GeneralParentMeeting{{ID: "gm-1", Title: "Incontro Generale"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/colloqui/general-meetings", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// GetGeneralMeeting
	repo.On("GetGeneralMeeting", mock.Anything, "gm-1").Return(&GeneralParentMeeting{ID: "gm-1", Title: "Incontro Generale"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/colloqui/general-meetings/gm-1", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// BookQueueTicket
	repo.On("BookQueueTicket", mock.Anything, mock.AnythingOfType("*colloqui.GeneralMeetingQueueTicket")).Return(nil).Once()
	ticketPayload := `{"meeting_id":"gm-1","teacher_id":"teacher-1","student_id":"s1","notes":"Chiarimento voti"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/colloqui/general-meetings/book-ticket", bytes.NewBufferString(ticketPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// ListQueueTickets
	repo.On("ListQueueTickets", mock.Anything, "gm-1", "", "parent-1").Return([]GeneralMeetingQueueTicket{{ID: "t1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/colloqui/general-meetings/gm-1/tickets", nil)
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// UpdateTicketStatus
	repo.On("UpdateTicketStatus", mock.Anything, "t1", "called", "").Return(nil).Once()
	ticketStatusPayload := `{"status":"called"}`
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/colloqui/general-meetings/tickets/t1/status", bytes.NewBufferString(ticketStatusPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Error handling
	repo.On("GetGeneralMeeting", mock.Anything, "missing").Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/colloqui/general-meetings/missing", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// CreateGeneralMeeting
	repo.On("CreateGeneralMeeting", mock.Anything, mock.AnythingOfType("*colloqui.GeneralParentMeeting"), mock.Anything).Return(nil).Once()
	repo.On("GetGeneralMeeting", mock.Anything, mock.Anything).Return(&GeneralParentMeeting{ID: "gm-2", Title: "Nuovo Incontro"}, nil).Once()
	gmPayload := `{"title":"Nuovo Incontro","event_date":"2026-11-15","start_time":"15:00","end_time":"18:00","slot_duration_minutes":10}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/colloqui/general-meetings", bytes.NewBufferString(gmPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// CreateAssembly
	repo.On("ExistsOverlappingSlot", mock.Anything, "admin-1", "2026-10-30", "10:00", "12:00").Return(false, nil).Once()
	repo.On("CreateSlot", mock.Anything, mock.AnythingOfType("*colloqui.ColloquioSlot")).Return(nil).Once()
	assemblyPayload := `{"class_id":"cls-1","title":"Assemblea Generale","date":"2026-10-30","start_time":"10:00","end_time":"12:00","max_bookings":50,"location":"Aula Magna"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/colloqui/assemblies", bytes.NewBufferString(assemblyPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// ListSlotBookings
	repo.On("GetSlotByID", mock.Anything, "slot-1").Return(&ColloquioSlot{ID: "slot-1", TeacherID: "teacher-1"}, nil).Once()
	repo.On("ListSlotBookings", mock.Anything, "slot-1").Return([]*ColloquioBooking{{ID: "b1", SlotID: "slot-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/colloqui/slots/slot-1/bookings", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// GetAvailabilityByTeacher
	repo.On("ListSlots", mock.Anything, mock.AnythingOfType("colloqui.SlotFilter")).Return([]*ColloquioSlot{{ID: "slot-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/colloqui/availability/teacher-1?from=2026-06-01&to=2026-06-15", nil)
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.AssertExpectations(t)
}
