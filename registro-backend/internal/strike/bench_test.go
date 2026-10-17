package strike

import (
	"math"
	"testing"
)

func BenchmarkIsValidIntention(b *testing.B) {
	intentions := []StrikeIntention{
		IntentionParticipates,
		IntentionNotParticipates,
		IntentionUndecided,
		IntentionUnanswered,
		"invalid",
	}

	idx := 0
	for b.Loop() {
		_ = IsValidIntention(intentions[idx%len(intentions)])
		idx++
	}
}

func BenchmarkAggregateStrikePercentages(b *testing.B) {
	totalStaff := 120
	participates := 45
	notParticipates := 50
	undecided := 15
	unanswered := totalStaff - (participates + notParticipates + undecided)

	for b.Loop() {
		pPct := math.Round((float64(participates)/float64(totalStaff))*10000) / 100
		npPct := math.Round((float64(notParticipates)/float64(totalStaff))*10000) / 100
		uPct := math.Round((float64(undecided)/float64(totalStaff))*10000) / 100
		unansPct := math.Round((float64(unanswered)/float64(totalStaff))*10000) / 100
		_ = pPct + npPct + uPct + unansPct
	}
}
