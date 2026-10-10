package general_meetings

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func BenchmarkListMeetings(b *testing.B) {
	repo := newMockRepo()
	svc := NewService(repo)
	ctx := context.Background()

	// Seed 50 meetings
	for i := 0; i < 50; i++ {
		id := fmt.Sprintf("m-%d", i)
		repo.meetings[id] = &GeneralMeeting{
			ID:          id,
			SchoolID:    "school-1",
			Title:       fmt.Sprintf("Riunione n.%d", i),
			MeetingDate: time.Now().AddDate(0, 0, i),
		}
	}

	for b.Loop() {
		_, _ = svc.ListMeetings(ctx, "school-1", "user-1", "teacher")
	}
}

func BenchmarkRegisterUser(b *testing.B) {
	repo := newMockRepo()
	svc := NewService(repo)
	ctx := context.Background()

	deadline := time.Now().Add(24 * time.Hour)
	maxParts := 100000000
	repo.meetings["m-bench"] = &GeneralMeeting{
		ID:                   "m-bench",
		SchoolID:             "school-1",
		RegistrationDeadline: &deadline,
		MaxParticipants:      &maxParts,
	}

	i := 0
	for b.Loop() {
		_ = svc.RegisterUser(ctx, "m-bench", fmt.Sprintf("u-%d", i))
		i++
	}
}
