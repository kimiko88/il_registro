package payments

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Standard OPI (Ordinativo Informatico di Pagamento e Incasso) SIOPE+ XML Structures
type FlussoGiornaleDiCassa struct {
	XMLName              xml.Name          `xml:"flusso_giornale_di_cassa"`
	CodiceABI            string            `xml:"codice_abi_bt"`
	IdentificativoFlusso string            `xml:"identificativo_flusso"`
	DataOraCreazione     string            `xml:"data_ora_creazione"`
	Quietanze            []OPIXMLQuietanza `xml:"dettaglio_quietanza"`
}

type OPIXMLQuietanza struct {
	NumeroQuietanza string  `xml:"numero_quietanza"`
	DataAccredito   string  `xml:"data_accredito"`
	Importo         float64 `xml:"importo"`
	IUV             string  `xml:"identificativo_univoco_versamento"`
	CodiceFiscale   string  `xml:"codice_fiscale_debitore"`
	Denominazione   string  `xml:"anagrafica_debitore"`
	Esito           string  `xml:"esito_accredito"`
}

// ParseOPIStream parses either OPI XML or CSV quietanze streams.
func ParseOPIStream(r io.Reader, format string) ([]OPIQuietanza, error) {
	if r == nil {
		return nil, errors.New("reader cannot be nil")
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading stream: %w", err)
	}

	trimmed := bytes.TrimSpace(buf)
	if len(trimmed) == 0 {
		return nil, errors.New("empty quietanze data stream")
	}

	normalizedFormat := strings.ToUpper(strings.TrimSpace(format))
	if normalizedFormat == "OPI_XML" || normalizedFormat == "XML" || bytes.HasPrefix(trimmed, []byte("<")) {
		return parseOPIXML(trimmed)
	}

	return parseOPICSV(trimmed)
}

func parseOPIXML(data []byte) ([]OPIQuietanza, error) {
	var flusso FlussoGiornaleDiCassa
	if err := xml.Unmarshal(data, &flusso); err != nil {
		return nil, fmt.Errorf("xml unmarshal error: %w", err)
	}

	var result []OPIQuietanza
	for _, q := range flusso.Quietanze {
		var accreditoDate time.Time
		if q.DataAccredito != "" {
			var errDate error
			accreditoDate, errDate = time.Parse("2006-01-02", q.DataAccredito)
			if errDate != nil {
				accreditoDate = time.Now()
			}
		} else {
			accreditoDate = time.Now()
		}

		cleanIUV := strings.TrimSpace(q.IUV)
		result = append(result, OPIQuietanza{
			CodiceFlusso:    flusso.IdentificativoFlusso,
			NumeroQuietanza: q.NumeroQuietanza,
			IUV:             cleanIUV,
			Importo:         q.Importo,
			DataAccredito:   accreditoDate,
			CodiceDebitore:  q.CodiceFiscale,
			NomeDebitore:    q.Denominazione,
			Esito:           q.Esito,
		})
	}
	return result, nil
}

func parseOPICSV(data []byte) ([]OPIQuietanza, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv parsing error: %w", err)
	}
	if len(records) < 2 {
		return nil, errors.New("csv contains insufficient rows (requires header + at least 1 row)")
	}

	header := records[0]
	iuvIdx := -1
	amountIdx := -1
	dateIdx := -1
	receiptIdx := -1
	cfIdx := -1
	nameIdx := -1

	for i, col := range header {
		lower := strings.ToLower(strings.TrimSpace(col))
		switch {
		case strings.Contains(lower, "iuv"):
			iuvIdx = i
		case strings.Contains(lower, "importo") || strings.Contains(lower, "amount"):
			amountIdx = i
		case strings.Contains(lower, "data") || strings.Contains(lower, "date"):
			dateIdx = i
		case strings.Contains(lower, "quietanza") || strings.Contains(lower, "ricevuta"):
			receiptIdx = i
		case strings.Contains(lower, "fiscale") || strings.Contains(lower, "cf"):
			cfIdx = i
		case strings.Contains(lower, "debitore") || strings.Contains(lower, "nominativo") || strings.Contains(lower, "nome"):
			nameIdx = i
		}
	}

	if iuvIdx == -1 || amountIdx == -1 {
		return nil, errors.New("csv missing mandatory columns: IUV and Importo")
	}

	var results []OPIQuietanza
	for rowIdx, row := range records[1:] {
		if len(row) <= iuvIdx || len(row) <= amountIdx {
			continue
		}
		rawIUV := strings.TrimSpace(row[iuvIdx])
		if rawIUV == "" {
			continue
		}

		rawAmount := strings.ReplaceAll(strings.TrimSpace(row[amountIdx]), ",", ".")
		amount, errAmount := strconv.ParseFloat(rawAmount, 64)
		if errAmount != nil {
			amount = 0.0
		}

		accredito := time.Now()
		if dateIdx != -1 && len(row) > dateIdx && row[dateIdx] != "" {
			if d, errD := time.Parse("2006-01-02", strings.TrimSpace(row[dateIdx])); errD == nil {
				accredito = d
			}
		}

		receiptNum := fmt.Sprintf("CSV-Q-%d", rowIdx+1)
		if receiptIdx != -1 && len(row) > receiptIdx && row[receiptIdx] != "" {
			receiptNum = row[receiptIdx]
		}

		cf := ""
		if cfIdx != -1 && len(row) > cfIdx {
			cf = row[cfIdx]
		}

		name := ""
		if nameIdx != -1 && len(row) > nameIdx {
			name = row[nameIdx]
		}

		results = append(results, OPIQuietanza{
			CodiceFlusso:    "CSV-IMPORT",
			NumeroQuietanza: receiptNum,
			IUV:             rawIUV,
			Importo:         amount,
			DataAccredito:   accredito,
			CodiceDebitore:  cf,
			NomeDebitore:    name,
			Esito:           "ACCREDITATO",
		})
	}

	return results, nil
}

// ReconcileQuietanze matches parsed quietanze against existing payments.
func ReconcileQuietanze(payments []*SchoolPayment, quietanze []OPIQuietanza) *ReconciliationReport {
	report := &ReconciliationReport{
		CodiceFlusso:       "RECONCILIATION-BATCH",
		TotalProcessed:     len(quietanze),
		ReconciledPayments: make([]*SchoolPayment, 0),
		UnmatchedQuietanze: make([]OPIQuietanza, 0),
	}

	if len(quietanze) > 0 {
		report.CodiceFlusso = quietanze[0].CodiceFlusso
	}

	// Map payments by IUV and by ID
	byIUV := make(map[string]*SchoolPayment)
	byID := make(map[string]*SchoolPayment)
	for _, p := range payments {
		if p.IUV != "" {
			byIUV[strings.TrimSpace(p.IUV)] = p
		}
		byID[p.ID] = p
	}

	for _, q := range quietanze {
		matchedPayment, exists := byIUV[q.IUV]
		if !exists {
			// Fallback: match by ID if IUV was stored as payment ID
			matchedPayment, exists = byID[q.IUV]
		}

		if exists && matchedPayment.Status != "paid" {
			matchedPayment.Status = "paid"
			matchedPayment.PaidAt = &q.DataAccredito
			matchedPayment.PaymentMethod = "PagoPA / OPI SIOPE+"
			matchedPayment.TransactionID = q.NumeroQuietanza
			matchedPayment.ReceiptNumber = fmt.Sprintf("OPI-%s", q.NumeroQuietanza)
			now := time.Now()
			matchedPayment.ReconciledAt = &now

			report.TotalReconciled++
			report.TotalAmount += q.Importo
			report.ReconciledPayments = append(report.ReconciledPayments, matchedPayment)
		} else {
			report.TotalUnmatched++
			report.UnmatchedQuietanze = append(report.UnmatchedQuietanze, q)
		}
	}

	return report
}

// CheckDelinquentPayments identifies payments past their due date and generates solleciti.
func CheckDelinquentPayments(payments []*SchoolPayment, now time.Time) []*PaymentSollecito {
	var solleciti []*PaymentSollecito

	for _, p := range payments {
		if p.Status == "pending" && p.DueDate.Before(now) {
			daysOverdue := int(now.Sub(p.DueDate).Hours() / 24)
			sollecitoNum := p.SollecitoCount + 1

			text := fmt.Sprintf(
				"SOLLECITO DI PAGAMENTO N. %d\nGentile genitore, si ricorda che il pagamento '%s' dell'importo di €%.2f risultava dovuto entro il %s ed è scaduto da %d giorni. Si invita a provvedere al versamento tramite PagoPA.",
				sollecitoNum,
				p.Title,
				p.Amount,
				p.DueDate.Format("02/01/2006"),
				daysOverdue,
			)

			solleciti = append(solleciti, &PaymentSollecito{
				PaymentID:      p.ID,
				StudentID:      p.StudentID,
				StudentName:    p.StudentName,
				Title:          p.Title,
				Amount:         p.Amount,
				DueDate:        p.DueDate,
				DaysOverdue:    daysOverdue,
				SollecitoCount: sollecitoNum,
				SollecitoText:  text,
			})
		}
	}

	return solleciti
}

// InitiateCheckoutSession creates a direct parent PagoPA online checkout session.
func InitiateCheckoutSession(p *SchoolPayment, returnURL string) (*CheckoutSession, error) {
	if p == nil {
		return nil, errors.New("payment cannot be nil")
	}
	if p.Status == "paid" {
		return nil, errors.New("payment is already paid")
	}
	if p.IUV == "" {
		iuv, err := GenerateIUV(0, "02", time.Now().UnixNano()%1000000000000)
		if err != nil {
			return nil, err
		}
		p.IUV = iuv
	}

	token := fmt.Sprintf("SES-PAGOPA-%s-%d", p.IUV, time.Now().Unix())
	redirect := fmt.Sprintf("https://checkout.pagopa.it/v2/session/%s?redirect=%s", token, returnURL)

	p.CheckoutSessionToken = token
	return &CheckoutSession{
		PaymentID:    p.ID,
		SessionToken: token,
		RedirectURL:  redirect,
		ExpiresAt:    time.Now().Add(30 * time.Minute),
		Amount:       p.Amount,
		IUV:          p.IUV,
	}, nil
}
