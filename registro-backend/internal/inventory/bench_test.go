package inventory

import (
	"testing"
	"time"
)

func BenchmarkCalculateDepreciation(b *testing.B) {
	acq := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = CalculateDepreciation(1200.0, 20.0, acq, now)
	}
}
