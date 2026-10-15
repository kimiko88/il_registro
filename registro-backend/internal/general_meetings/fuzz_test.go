package general_meetings

import (
	"context"
	"testing"
)

func FuzzCreateMeetingDateParsing(f *testing.F) {
	f.Add("2026-10-15T14:30:00Z", "Collegio Docenti")
	f.Add("2026-10-15", "Consiglio di Istituto")
	f.Add("invalid-date", "Test Title")
	f.Add("", "")

	f.Fuzz(func(t *testing.T, dateStr, title string) {
		repo := newMockRepo()
		svc := NewService(repo)

		req := CreateGeneralMeetingRequest{
			Title:       title,
			MeetingDate: dateStr,
		}

		_, err := svc.CreateMeeting(context.Background(), "user-1", "school-1", req)
		if dateStr == "invalid-date" || dateStr == "" {
			if err == nil {
				t.Fatalf("expected error for invalid date format %q", dateStr)
			}
		}
	})
}
