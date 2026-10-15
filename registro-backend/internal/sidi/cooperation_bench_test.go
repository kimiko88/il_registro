package sidi

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkSyncStudentCodes(b *testing.B) {
	svc := NewService(nil)
	ctx := context.Background()

	const count = 500
	students := make([]StudenteSIDI, count)
	for i := 0; i < count; i++ {
		students[i] = StudenteSIDI{
			CodiceFiscale: fmt.Sprintf("CF%014d", i),
			CodiceSIDI:    "",
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = svc.SyncStudentCodes(ctx, "school-1", students)
	}
}
