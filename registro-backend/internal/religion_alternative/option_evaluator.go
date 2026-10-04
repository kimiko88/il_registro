package religion_alternative

import (
	"time"
)

// Ministerial options for students regarding Catholic Religion Teaching (IRC)
const (
	OptionIRC                = "irc"
	OptionMateriaAlternativa = "materia_alternativa"
	OptionStudioAssistito    = "studio_assistito"
	OptionStudioLibero       = "studio_libero"
	OptionUscitaScuola       = "uscita_scuola"
)

// Synthetic judgments for IRC and Alternative Subject (Art. 309 D.Lgs. 297/1994)
const (
	JudgmentOttimo         = "ottimo"
	JudgmentDistinto       = "distinto"
	JudgmentBuono          = "buono"
	JudgmentSufficiente    = "sufficiente"
	JudgmentNonSufficiente = "non_sufficiente"
)

type StudentReligionOption struct {
	ID           string    `json:"id"`
	StudentID    string    `json:"student_id"`
	SchoolID     string    `json:"school_id"`
	AcademicYear string    `json:"academic_year"`
	OptionType   string    `json:"option_type"`
	Notes        string    `json:"notes,omitempty"`
	ChosenAt     time.Time `json:"chosen_at"`
}

type AlternativeEvaluation struct {
	ID               string    `json:"id"`
	StudentID        string    `json:"student_id"`
	SchoolID         string    `json:"school_id"`
	GroupID          *string   `json:"group_id,omitempty"`
	ClassID          *string   `json:"class_id,omitempty"`
	Period           string    `json:"period"` // q1, q2, finale
	SubjectKind      string    `json:"subject_kind"` // irc, materia_alternativa
	JudgmentLevel    string    `json:"judgment_level"`
	DescriptiveNotes string    `json:"descriptive_notes,omitempty"`
	CreatedBy        *string   `json:"created_by,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type StudentOptionSummary struct {
	StudentID  string `json:"student_id"`
	OptionType string `json:"option_type"`
}

// IsValidOptionType verifies if the ministerial option is allowed
func IsValidOptionType(optionType string) bool {
	switch optionType {
	case OptionIRC, OptionMateriaAlternativa, OptionStudioAssistito, OptionStudioLibero, OptionUscitaScuola:
		return true
	default:
		return false
	}
}

// IsValidJudgmentLevel verifies if the judgment matches the Art. 309 D.Lgs. 297/1994 scale
func IsValidJudgmentLevel(level string) bool {
	switch level {
	case JudgmentOttimo, JudgmentDistinto, JudgmentBuono, JudgmentSufficiente, JudgmentNonSufficiente:
		return true
	default:
		return false
	}
}

// IsExcusedFromAbsenceLimit checks if the student option implies exemption from the 25% max absence limit
func IsExcusedFromAbsenceLimit(optionType string) bool {
	return optionType == OptionUscitaScuola
}

// FilterStudentsForAlternativeGroup extracts students who must be enrolled in open-class alternative subject groups
func FilterStudentsForAlternativeGroup(students []StudentOptionSummary) []StudentOptionSummary {
	var results []StudentOptionSummary
	for _, s := range students {
		if s.OptionType == OptionMateriaAlternativa {
			results = append(results, s)
		}
	}
	return results
}
