package religion_alternative

import (
	"strings"
	"testing"
)

func FuzzIsValidOptionType(f *testing.F) {
	seeds := []string{
		OptionIRC,
		OptionMateriaAlternativa,
		OptionStudioAssistito,
		OptionStudioLibero,
		OptionUscitaScuola,
		"",
		"IRC",
		"materia_alternativa ",
		"unknown",
		"religione",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, option string) {
		isValid := IsValidOptionType(option)
		trimmed := strings.TrimSpace(option)

		switch trimmed {
		case OptionIRC, OptionMateriaAlternativa, OptionStudioAssistito, OptionStudioLibero, OptionUscitaScuola:
			if trimmed == option && !isValid {
				t.Fatalf("expected valid for exact match %q", option)
			}
		default:
			if isValid {
				t.Fatalf("expected invalid for non-canonical %q", option)
			}
		}
	})
}

func FuzzIsValidJudgmentLevel(f *testing.F) {
	seeds := []string{
		JudgmentOttimo,
		JudgmentDistinto,
		JudgmentBuono,
		JudgmentSufficiente,
		JudgmentNonSufficiente,
		"insufficiente",
		"eccellente",
		"",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, judgment string) {
		isValid := IsValidJudgmentLevel(judgment)
		switch judgment {
		case JudgmentOttimo, JudgmentDistinto, JudgmentBuono, JudgmentSufficiente, JudgmentNonSufficiente:
			if !isValid {
				t.Fatalf("expected valid for %q", judgment)
			}
		default:
			if isValid {
				t.Fatalf("expected invalid for %q", judgment)
			}
		}
	})
}

func FuzzFilterStudentsForAlternativeGroup(f *testing.F) {
	f.Add(OptionMateriaAlternativa, OptionIRC, OptionUscitaScuola)
	f.Add(OptionStudioAssistito, OptionStudioLibero, OptionMateriaAlternativa)

	f.Fuzz(func(t *testing.T, opt1, opt2, opt3 string) {
		input := []StudentOptionSummary{
			{StudentID: "s-1", OptionType: opt1},
			{StudentID: "s-2", OptionType: opt2},
			{StudentID: "s-3", OptionType: opt3},
		}

		filtered := FilterStudentsForAlternativeGroup(input)

		for _, s := range filtered {
			if s.OptionType != OptionMateriaAlternativa {
				t.Fatalf("filtered list contains non-alternative option: %q", s.OptionType)
			}
		}
	})
}
