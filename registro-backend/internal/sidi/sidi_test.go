package sidi

import (
	"context"
	"strings"
	"testing"
)

func TestSidiXmlGenerationAndValidation(t *testing.T) {
	svc := NewService(nil)
	ctx := context.Background()

	req := GenerateSidiRequest{
		ExportType: "SCRUTINIO_SETTEMBRE_DEBITI",
		SchoolYear: "2025/2026",
	}

	record, validation, err := svc.GenerateExport(ctx, "school-123", "sec-user-1", req)
	if err != nil {
		t.Fatalf("GenerateExport failed: %v", err)
	}

	if record == nil || record.XMLContent == "" {
		t.Fatalf("expected non-empty XML content")
	}

	if !strings.Contains(record.XMLContent, "<FlussoSIDI") {
		t.Errorf("XML does not contain root <FlussoSIDI>")
	}

	if !strings.Contains(record.XMLContent, "SCRUTINIO_SETTEMBRE_DEBITI") {
		t.Errorf("XML does not contain export type")
	}

	if validation == nil || !validation.Valid {
		t.Errorf("expected sample data to be valid, got errors: %v", validation.Errors)
	}
}

func TestSidiValidationWithInvalidData(t *testing.T) {
	builder := NewBuilder()
	invalidStudents := []StudenteSIDI{
		{
			CodiceSIDI:    "",
			CodiceFiscale: "INVALID_CF",
			Cognome:       "Rossi",
			Nome:          "Mario",
		},
	}

	val := builder.ValidateSidiData(invalidStudents)
	if val.Valid {
		t.Errorf("expected validation to fail for invalid student data")
	}
	if len(val.MissingSidiIDs) != 1 {
		t.Errorf("expected 1 missing SIDI ID, got %d", len(val.MissingSidiIDs))
	}
	if len(val.Errors) != 1 {
		t.Errorf("expected 1 CF format error, got %d", len(val.Errors))
	}
}
