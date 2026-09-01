package reports

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportsInitialization(t *testing.T) {
	svc := NewService(nil)
	assert.NotNil(t, svc)
}

func TestExportSidiStudentsXML(t *testing.T) {
	svc := NewService(nil)
	ctx := context.Background()
	data, err := svc.ExportSidiStudentsXML(ctx, "class-uuid-1")
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	assert.Contains(t, string(data), "<FlussoAnagrafeSIDI>")
	assert.Contains(t, string(data), "<ClassID>class-uuid-1</ClassID>")
	assert.Contains(t, string(data), "<CodiceFiscale>RSSMRA08A01H501U</CodiceFiscale>")
}

func TestExportSidiScrutiniXML(t *testing.T) {
	svc := NewService(nil)
	ctx := context.Background()
	data, err := svc.ExportSidiScrutiniXML(ctx, "class-uuid-1", 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "<FlussoScrutiniSIDI>")
	assert.Contains(t, string(data), "<Quadrimestre>1</Quadrimestre>")
	assert.Contains(t, string(data), "<Voto>8</Voto>")

	// Default quadrimestre when 0 or negative
	dataDefault, err := svc.ExportSidiScrutiniXML(ctx, "class-uuid-1", 0)
	assert.NoError(t, err)
	assert.Contains(t, string(dataDefault), "<Quadrimestre>2</Quadrimestre>")
}

func TestExportSidiAttendanceCSV(t *testing.T) {
	svc := NewService(nil)
	ctx := context.Background()
	data, err := svc.ExportSidiAttendanceCSV(ctx, "class-uuid-1")
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	csvStr := string(data)
	assert.Contains(t, csvStr, "CODICE_FISCALE;COGNOME;NOME;CLASSE;ORE_ASSENZA_GIUSTIFICATE")
	assert.Contains(t, csvStr, "RSSMRA08A01H501U;Rossi;Mario")
}
