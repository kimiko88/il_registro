package strike

import (
	"strings"
	"testing"
)

func FuzzIsValidIntention(f *testing.F) {
	seeds := []string{
		string(IntentionParticipates),
		string(IntentionNotParticipates),
		string(IntentionUndecided),
		string(IntentionUnanswered),
		"participates",
		"not_participates",
		"undecided",
		"yes",
		"no",
		"",
		"PARTICIPATES",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		val := StrikeIntention(input)
		isValid := IsValidIntention(val)

		trimmed := strings.TrimSpace(input)
		switch StrikeIntention(trimmed) {
		case IntentionParticipates, IntentionNotParticipates, IntentionUndecided:
			if trimmed == input && !isValid {
				t.Fatalf("expected valid for %q", input)
			}
		default:
			if isValid {
				t.Fatalf("expected invalid for %q", input)
			}
		}
	})
}
