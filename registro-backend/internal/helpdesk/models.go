package helpdesk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type HelpDeskSlot struct {
	ID           string    `json:"id"`
	SchoolID     string    `json:"school_id"`
	TeacherID    string    `json:"teacher_id"`
	SubjectID    string    `json:"subject_id"`
	SlotDate     string    `json:"slot_date"` // YYYY-MM-DD
	StartTime    string    `json:"start_time"` // HH:MM
	EndTime      string    `json:"end_time"`   // HH:MM
	RoomID       *string   `json:"room_id,omitempty"`
	MaxCapacity  int       `json:"max_capacity"`
	Status       string    `json:"status"` // open, fully_booked, completed, cancelled
	BookingsCount int      `json:"bookings_count,omitempty"`
	TeacherName  string    `json:"teacher_name,omitempty"`
	SubjectName  string    `json:"subject_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type HelpDeskBooking struct {
	ID               string     `json:"id"`
	SlotID           string     `json:"slot_id"`
	StudentID        string     `json:"student_id"`
	TopicDescription string     `json:"topic_description"`
	Status           string     `json:"status"` // booked, attended, absent, cancelled
	BookedAt         time.Time  `json:"booked_at"`
	AttendedAt       *time.Time `json:"attended_at,omitempty"`
	StudentName      string     `json:"student_name,omitempty"`
}

type FISAccountingSummary struct {
	TeacherID     string  `json:"teacher_id"`
	TeacherName   string  `json:"teacher_name"`
	CompletedSlots int     `json:"completed_slots"`
	TotalHours    float64 `json:"total_hours"`
	AttendedCount int     `json:"attended_count"`
}

func CalculateSlotHours(startTime, endTime string) (float64, error) {
	parseMinutes := func(t string) (int, error) {
		parts := strings.Split(t, ":")
		if len(parts) < 2 {
			return 0, fmt.Errorf("formato orario non valido: %s", t)
		}
		h, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		return h*60 + m, nil
	}

	startMins, err := parseMinutes(startTime)
	if err != nil {
		return 0, err
	}
	endMins, err := parseMinutes(endTime)
	if err != nil {
		return 0, err
	}

	diff := endMins - startMins
	if diff <= 0 {
		return 0, errors.New("l'orario di fine deve essere successivo all'orario di inizio")
	}

	return float64(diff) / 60.0, nil
}

func CanBookSlot(slot HelpDeskSlot, currentBookings int) bool {
	return currentBookings < slot.MaxCapacity
}
