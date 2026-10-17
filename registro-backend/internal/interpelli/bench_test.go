package interpelli

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func BenchmarkCalculateScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = CalculateScore(105.0, true, true, 18, "C1", 3)
	}
}

func BenchmarkGetGraduatoria(b *testing.B) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	svc := NewService(repo)

	noticeID := "bench-notice"
	_ = repo.SaveNotice(ctx, &InterpelloNotice{
		ID:       noticeID,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   "aperto",
	})

	const numCandidates = 1000
	for i := 0; i < numCandidates; i++ {
		grade := 60.0 + float64(i%51)
		_, _, _, total := CalculateScore(grade, i%10 == 0, i%3 == 0, i%24, "B2", i%4)
		_ = repo.SaveCandidatura(ctx, &InterpelloCandidatura{
			ID:              fmt.Sprintf("cand-%d", i),
			NoticeID:        noticeID,
			CandidateName:   fmt.Sprintf("Name-%d", i),
			TotalScore:      total,
			GraduationGrade: grade,
			HasHabilitation: i%3 == 0,
			ScoreService:    float64(i % 30),
		})
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = svc.GetGraduatoria(ctx, noticeID)
	}
}
