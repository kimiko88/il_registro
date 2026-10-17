package pdp

import (
	"encoding/json"
	"testing"

	"registro-backend/pkg/crypto"
)

func BenchmarkPDPDiagnosisEncryption(b *testing.B) {
	diagnosis := "Diagnosi clinica specialistica di Dislessia Evolutiva (F81.0) redatta ai sensi della Legge 170/2010"

	for b.Loop() {
		_, _ = crypto.EncryptString(diagnosis)
	}
}

func BenchmarkPDPDiagnosisDecryption(b *testing.B) {
	diagnosis := "Diagnosi clinica specialistica di Dislessia Evolutiva (F81.0) redatta ai sensi della Legge 170/2010"
	ciphertext, _ := crypto.EncryptString(diagnosis)

	for b.Loop() {
		_, _ = crypto.DecryptString(ciphertext)
	}
}

func BenchmarkPDPContentMarshalling(b *testing.B) {
	content := PdpContent{
		Objectives: []string{
			"Raggiungimento degli obiettivi minimi disciplinari",
			"Miglioramento dell'autonomia nello studio individuale",
		},
		Compensative: []string{
			"calcolatrice",
			"sintesi_vocale",
			"mappe_concettuali",
			"tempo_aggiuntivo_30",
		},
		Dispensative: []string{
			"lettura_ad_alta_voce",
			"scrittura_veloce",
			"copiatura_lavagna",
		},
		EvaluationTools: []string{
			"Interrogazioni programmate",
			"Verifiche scritte con riduzione del numero di quesiti",
		},
		Notes:      "Verifica periodica al termine del primo quadrimestre",
		ReviewDate: "2026-01-31",
	}

	for b.Loop() {
		data, _ := json.Marshal(content)
		var c PdpContent
		_ = json.Unmarshal(data, &c)
	}
}
