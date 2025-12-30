package scheduling

import (
	"strings"
)

type CalendarService struct{}

func NewCalendarService() *CalendarService {
	return &CalendarService{}
}

func (c *CalendarService) GenerateICS(slots []ColloquioSlot, bookings []ColloquioBooking) string {
	var sb strings.Builder
	sb.WriteString("BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//RegistroElettronico//NONSGML v1.0//EN\n")

	// Logic to add events
	// ...

	sb.WriteString("END:VCALENDAR")
	return sb.String()
}
