package payments

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOPIStreamXML(t *testing.T) {
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

	quietanze, err := ParseOPIStream(strings.NewReader(xmlData), "OPI_XML")
	require.NoError(t, err)
	require.Len(t, quietanze, 1)
	assert.Equal(t, "Q1001", quietanze[0].NumeroQuietanza)
	assert.Equal(t, "00100000000000142", quietanze[0].IUV)
	assert.Equal(t, 15.50, quietanze[0].Importo)
	assert.Equal(t, "Mario Rossi", quietanze[0].NomeDebitore)
}

func TestParseOPIStreamCSV(t *testing.T) {
	csvData := `IUV,Importo,Data,Quietanza,CodiceFiscale,NomeDebitore
00100000000000142,45.00,2026-10-05,Q9988,BNCGLI08B41H501A,Giulia Bianchi
00100000000000143,120.00,2026-10-06,Q9989,VRDLGI08C12H501B,Luigi Verdi`

	quietanze, err := ParseOPIStream(strings.NewReader(csvData), "CSV_BANCA")
	require.NoError(t, err)
	require.Len(t, quietanze, 2)
	assert.Equal(t, "00100000000000142", quietanze[0].IUV)
	assert.Equal(t, 45.00, quietanze[0].Importo)
	assert.Equal(t, "Q9988", quietanze[0].NumeroQuietanza)
	assert.Equal(t, "Giulia Bianchi", quietanze[0].NomeDebitore)
}

func TestReconcileQuietanze(t *testing.T) {
	payments := []*SchoolPayment{
		{
			ID:     "pay-1",
			IUV:    "00100000000000142",
			Amount: 45.00,
			Status: "pending",
			Title:  "Gita Scolastica",
		},
		{
			ID:     "pay-2",
			IUV:    "00100000000000143",
			Amount: 80.00,
			Status: "pending",
			Title:  "Contributo Volontario",
		},
	}

	quietanze := []OPIQuietanza{
		{
			NumeroQuietanza: "Q9988",
			IUV:             "00100000000000142",
			Importo:         45.00,
			DataAccredito:   time.Now(),
		},
		{
			NumeroQuietanza: "Q9999",
			IUV:             "99999999999999999", // Unmatched
			Importo:         10.00,
			DataAccredito:   time.Now(),
		},
	}

	report := ReconcileQuietanze(payments, quietanze)
	assert.Equal(t, 2, report.TotalProcessed)
	assert.Equal(t, 1, report.TotalReconciled)
	assert.Equal(t, 1, report.TotalUnmatched)
	assert.Equal(t, 45.00, report.TotalAmount)
	assert.Equal(t, "paid", payments[0].Status)
	assert.Equal(t, "pending", payments[1].Status)
}

func TestCheckDelinquentPayments(t *testing.T) {
	now := time.Now()
	payments := []*SchoolPayment{
		{
			ID:          "pay-overdue",
			Title:       "Assicurazione Obbligatoria",
			Amount:      12.00,
			DueDate:     now.Add(-10 * 24 * time.Hour),
			Status:      "pending",
			StudentID:   "std-1",
			StudentName: "Paolo Neri",
		},
		{
			ID:          "pay-ontime",
			Title:       "Quota Esami",
			Amount:      50.00,
			DueDate:     now.Add(15 * 24 * time.Hour),
			Status:      "pending",
			StudentID:   "std-2",
			StudentName: "Anna Ferrari",
		},
	}

	solleciti := CheckDelinquentPayments(payments, now)
	require.Len(t, solleciti, 1)
	assert.Equal(t, "pay-overdue", solleciti[0].PaymentID)
	assert.Equal(t, 10, solleciti[0].DaysOverdue)
	assert.Contains(t, solleciti[0].SollecitoText, "SOLLECITO DI PAGAMENTO N. 1")
}

func TestInitiateCheckoutSession(t *testing.T) {
	p := &SchoolPayment{
		ID:     "pay-online",
		Amount: 85.00,
		Status: "pending",
	}

	session, err := InitiateCheckoutSession(p, "https://scuola.edu.it/esito-pagamento")
	require.NoError(t, err)
	assert.NotEmpty(t, session.SessionToken)
	assert.Contains(t, session.RedirectURL, "https://checkout.pagopa.it")
	assert.Equal(t, 85.00, session.Amount)
	assert.NotEmpty(t, p.IUV)
}
