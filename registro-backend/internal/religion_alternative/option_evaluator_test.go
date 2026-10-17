package religion_alternative

import (
	"testing"
)

func TestOptionTypeValidation(t *testing.T) {
	validOptions := []string{
		OptionIRC,
		OptionMateriaAlternativa,
		OptionStudioAssistito,
		OptionStudioLibero,
		OptionUscitaScuola,
	}

	for _, opt := range validOptions {
		if !IsValidOptionType(opt) {
			t.Errorf("expected option %s to be valid", opt)
		}
	}

	invalidOptions := []string{"invalid", "calcio", "dormire", ""}
	for _, opt := range invalidOptions {
		if IsValidOptionType(opt) {
			t.Errorf("expected option %s to be invalid", opt)
		}
	}
}

func TestJudgmentLevelValidation(t *testing.T) {
	validJudgments := []string{
		JudgmentOttimo,
		JudgmentDistinto,
		JudgmentBuono,
		JudgmentSufficiente,
		JudgmentNonSufficiente,
	}

	for _, j := range validJudgments {
		if !IsValidJudgmentLevel(j) {
			t.Errorf("expected judgment %s to be valid", j)
		}
	}

	invalidJudgments := []string{"10", "6", "insufficiente", "gravemente_insufficiente", ""}
	for _, j := range invalidJudgments {
		if IsValidJudgmentLevel(j) {
			t.Errorf("expected judgment %s to be invalid", j)
		}
	}
}

func TestAttendanceExemption(t *testing.T) {
	// Uscita scuola must NOT be counted as unexcused absence (D.P.R. 122/2009)
	if !IsExcusedFromAbsenceLimit(OptionUscitaScuola) {
		t.Errorf("OptionUscitaScuola should be excused from 25%% absence limit calculation")
	}

	if IsExcusedFromAbsenceLimit(OptionIRC) {
		t.Errorf("OptionIRC should not have automatic absence limit exemption")
	}

	if IsExcusedFromAbsenceLimit(OptionMateriaAlternativa) {
		t.Errorf("OptionMateriaAlternativa must be attended in assigned group, not excused from limit")
	}
}

func TestFilterStudentsForAlternativeGroup(t *testing.T) {
	options := []StudentOptionSummary{
		{StudentID: "s1", OptionType: OptionMateriaAlternativa},
		{StudentID: "s2", OptionType: OptionIRC},
		{StudentID: "s3", OptionType: OptionMateriaAlternativa},
		{StudentID: "s4", OptionType: OptionUscitaScuola},
		{StudentID: "s5", OptionType: OptionStudioAssistito},
	}

	alternatives := FilterStudentsForAlternativeGroup(options)
	if len(alternatives) != 2 {
		t.Fatalf("expected 2 students for alternative group, got %d", len(alternatives))
	}
	if alternatives[0].StudentID != "s1" || alternatives[1].StudentID != "s3" {
		t.Errorf("unexpected students filtered: %+v", alternatives)
	}
}
