package religion_alternative

import (
	"fmt"
	"testing"
)

func BenchmarkIsValidOptionType(b *testing.B) {
	options := []string{
		OptionIRC,
		OptionMateriaAlternativa,
		OptionStudioAssistito,
		OptionStudioLibero,
		OptionUscitaScuola,
		"invalid_option",
	}

	idx := 0
	for b.Loop() {
		_ = IsValidOptionType(options[idx%len(options)])
		idx++
	}
}

func BenchmarkIsValidJudgmentLevel(b *testing.B) {
	judgments := []string{
		JudgmentOttimo,
		JudgmentDistinto,
		JudgmentBuono,
		JudgmentSufficiente,
		JudgmentNonSufficiente,
		"invalid",
	}

	idx := 0
	for b.Loop() {
		_ = IsValidJudgmentLevel(judgments[idx%len(judgments)])
		idx++
	}
}

func BenchmarkFilterStudentsForAlternativeGroup(b *testing.B) {
	students := make([]StudentOptionSummary, 100)
	opts := []string{
		OptionIRC,
		OptionMateriaAlternativa,
		OptionStudioAssistito,
		OptionStudioLibero,
		OptionUscitaScuola,
	}
	for i := 0; i < 100; i++ {
		students[i] = StudentOptionSummary{
			StudentID:  fmt.Sprintf("stu-%d", i),
			OptionType: opts[i%len(opts)],
		}
	}

	for b.Loop() {
		_ = FilterStudentsForAlternativeGroup(students)
	}
}
