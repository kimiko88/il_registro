package sidi

import (
	"encoding/xml"
	"fmt"
	"time"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) BuildSidiXML(schoolCode, schoolName, schoolYear, exportType string, students []StudenteSIDI, scrutini []ScrutinioSIDI) ([]byte, error) {
	flusso := FlussoSIDI{
		Versione: "2.4",
		Testata: Testata{
			CodiceScuola: schoolCode,
			AnnoSco:      schoolYear,
			DataExport:   time.Now(),
			TipoFlusso:   exportType,
		},
		DatiScuola: DatiScuola{
			Denominazione: schoolName,
			Studenti:      students,
			Scrutini:      scrutini,
		},
	}

	data, err := xml.MarshalIndent(flusso, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SIDI XML: %w", err)
	}

	header := []byte(xml.Header)
	return append(header, data...), nil
}

func (b *Builder) ValidateSidiData(students []StudenteSIDI) SidiValidationResult {
	result := SidiValidationResult{
		Valid:        true,
		TotalRecords: len(students),
	}

	for _, s := range students {
		if s.CodiceSIDI == "" {
			result.MissingSidiIDs = append(result.MissingSidiIDs, fmt.Sprintf("%s %s (CF: %s)", s.Cognome, s.Nome, s.CodiceFiscale))
		}
		if len(s.CodiceFiscale) != 16 {
			result.Errors = append(result.Errors, fmt.Sprintf("Codice Fiscale non valido per %s %s: '%s'", s.Cognome, s.Nome, s.CodiceFiscale))
		}
	}

	if len(result.MissingSidiIDs) > 0 || len(result.Errors) > 0 {
		result.Valid = false
	}
	return result
}
