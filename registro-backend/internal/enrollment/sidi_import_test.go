package enrollment

import (
	"strings"
	"testing"
)

func TestParseSIDIExportCSV(t *testing.T) {
	csvData := `ID_DOMANDA;NOME;COGNOME;CODICE_FISCALE;DATA_NASCITA;SESSO;VOTO_LICENZA;INDIRIZZO;SECONDA_LINGUA;L104;DSA;RELIGIONE;COMPAGNI_RICHIESTI;GENITORE1_NOME;GENITORE1_COGNOME;GENITORE1_EMAIL;GENITORE1_TEL
SIDI-2026-001;Mario;Rossi;RSSMRA10A01H501U;2012-05-14;M;8;Scientifico;Spagnolo;false;false;irc;BNCGPP10A01H501V;Giuseppe;Rossi;giuseppe.rossi@example.com;3331112223
SIDI-2026-002;Giuseppe;Bianchi;BNCGPP10A01H501V;2012-03-22;M;7;Scientifico;Spagnolo;false;false;irc;RSSMRA10A01H501U;Luigi;Bianchi;luigi.bianchi@example.com;3334445556
SIDI-2026-003;Lucia;Verdi;VRDLCU10B41H501Z;2012-09-08;F;9;Scientifico;Francese;true;false;materia_alternativa;;Anna;Verdi;anna.verdi@example.com;3337778889
`
	apps, err := ParseSIDIExport(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error parsing SIDI CSV: %v", err)
	}

	if len(apps) != 3 {
		t.Fatalf("expected 3 applications, got %d", len(apps))
	}

	app1 := apps[0]
	if app1.SidiApplicationID != "SIDI-2026-001" || app1.StudentFirstName != "Mario" || app1.Gender != "M" {
		t.Errorf("unexpected data for student 1: %+v", app1)
	}
	if app1.MiddleSchoolGrade != 8 {
		t.Errorf("expected grade 8, got %d", app1.MiddleSchoolGrade)
	}
	if len(app1.RequestedClassmates) == 0 || app1.RequestedClassmates[0] != "BNCGPP10A01H501V" {
		t.Errorf("expected requested classmate BNCGPP10A01H501V, got %v", app1.RequestedClassmates)
	}

	app3 := apps[2]
	if !app3.HasDisabilityL104 {
		t.Errorf("expected student 3 to have L.104 disability")
	}
	if app3.ReligionChoice != "materia_alternativa" {
		t.Errorf("expected materia_alternativa, got %s", app3.ReligionChoice)
	}
}

func TestParseSIDIExport_Validation(t *testing.T) {
	_, err := ParseSIDIExport(strings.NewReader(""))
	if err == nil {
		t.Errorf("expected error for empty stream")
	}

	invalidCSV := `COL1;COL2
1;2`
	_, err = ParseSIDIExport(strings.NewReader(invalidCSV))
	if err == nil {
		t.Errorf("expected error for invalid SIDI headers")
	}
}
