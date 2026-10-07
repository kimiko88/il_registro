package payments

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func BenchmarkGenerateIUV(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateIUV(0, "01", int64(i))
	}
}

func BenchmarkValidateIUV(b *testing.B) {
	iuv, _ := GenerateIUV(0, "01", 123456789)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ValidateIUV(iuv)
	}
}

func BenchmarkReconcileQuietanze(b *testing.B) {
	const count = 500
	payments := make([]*SchoolPayment, count)
	quietanze := make([]OPIQuietanza, count)

	for i := 0; i < count; i++ {
		iuv, _ := GenerateIUV(0, "01", int64(i))
		payments[i] = &SchoolPayment{
			ID:     fmt.Sprintf("pay-%d", i),
			IUV:    iuv,
			Amount: 50.0,
			Status: "pending",
		}
		quietanze[i] = OPIQuietanza{
			NumeroQuietanza: fmt.Sprintf("Q-%d", i),
			IUV:             iuv,
			Importo:         50.0,
			DataAccredito:   time.Now(),
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ReconcileQuietanze(payments, quietanze)
	}
}

func BenchmarkParseOPIStreamXML(b *testing.B) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<flusso_giornale_di_cassa>
    <codice_abi_bt>01005</codice_abi_bt>
    <identificativo_flusso>OPI-2026-00981</identificativo_flusso>
    <data_ora_creazione>2026-10-07T10:00:00</data_ora_creazione>
    <dettaglio_quietanza>
        <numero_quietanza>Q1001</numero_quietanza>
        <data_accredito>2026-10-06</data_accredito>
        <importo>15.50</importo>
        <identificativo_univoco_versamento>00100000000000142</identificativo_univoco_versamento>
        <codice_fiscale_debitore>RSSMRA08A01H501Z</codice_fiscale_debitore>
        <anagrafica_debitore>Mario Rossi</anagrafica_debitore>
        <esito_accredito>ESEGUITO</esito_accredito>
    </dettaglio_quietanza>
</flusso_giornale_di_cassa>`

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseOPIStream(strings.NewReader(xmlData), "OPI_XML")
	}
}
