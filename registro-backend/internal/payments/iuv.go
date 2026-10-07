package payments

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// CalculateMod97 calculates the ISO 7064 Mod 97-10 check digits.
// According to AgID / PagoPA specs:
// Check digits = 98 - ((numericString * 100) % 97)
// If result is 1 digit, prepend with "0".
func CalculateMod97(numericString string) (string, error) {
	if len(numericString) == 0 {
		return "", errors.New("input string cannot be empty")
	}
	for _, ch := range numericString {
		if ch < '0' || ch > '9' {
			return "", fmt.Errorf("invalid character '%c', only numeric digits allowed", ch)
		}
	}

	// We append "00" to simulate (numericString * 100)
	prepended := numericString + "00"
	num := new(big.Int)
	num, ok := num.SetString(prepended, 10)
	if !ok {
		return "", errors.New("failed to parse numeric string")
	}

	modVal := new(big.Int).Mod(num, big.NewInt(97)).Int64()
	checkVal := 98 - modVal
	if checkVal == 97 {
		checkVal = 0
	}

	return fmt.Sprintf("%02d", checkVal), nil
}

// GenerateIUV creates a compliant 17-digit IUV (Identificativo Univoco Versamento):
// Structure:
// - Digit 1: Auxiliary Digit ('0', '1', '2', or '3', default '0')
// - Digits 2-3: Application Code (e.g., '12' for school voluntary/mandatory fees)
// - Digits 4-15: Progressive numeric string (12 digits, zero-padded timestamp/sequence)
// - Digits 16-17: ISO 7064 Mod 97-10 check digits
func GenerateIUV(auxiliaryDigit int, applicationCode string, progressive int64) (string, error) {
	if auxiliaryDigit < 0 || auxiliaryDigit > 3 {
		return "", fmt.Errorf("auxiliary digit must be between 0 and 3, got %d", auxiliaryDigit)
	}
	if len(applicationCode) != 2 {
		return "", fmt.Errorf("application code must be exactly 2 digits, got '%s'", applicationCode)
	}
	if progressive < 0 {
		return "", fmt.Errorf("progressive sequence must be non-negative, got %d", progressive)
	}

	// Format base 15 digits
	base := fmt.Sprintf("%d%s%012d", auxiliaryDigit, applicationCode, progressive%1000000000000)

	checkDigits, err := CalculateMod97(base)
	if err != nil {
		return "", fmt.Errorf("error generating mod97 check digits: %w", err)
	}

	return base + checkDigits, nil
}

// ValidateIUV verifies if a given 17-digit IUV has valid check digits.
func ValidateIUV(iuv string) bool {
	clean := strings.TrimSpace(iuv)
	if len(clean) != 17 {
		return false
	}
	for _, ch := range clean {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	base := clean[:15]
	expectedCheck := clean[15:]

	calculated, err := CalculateMod97(base)
	if err != nil {
		return false
	}

	return calculated == expectedCheck
}

// GeneratePagoPAQRCodePayload formats the standard PagoPA notice QR code payload:
// Format: PAGOPA|002|{IUV}|{cfEnte}|{importoCentesimi}
func GeneratePagoPAQRCodePayload(cfEnte string, iuv string, amountCents int64) string {
	cf := strings.TrimSpace(cfEnte)
	cleanIUV := strings.TrimSpace(iuv)
	return fmt.Sprintf("PAGOPA|002|%s|%s|%d", cleanIUV, cf, amountCents)
}

// ParsePagoPAQRCodePayload parses a PagoPA QR code string back into its constituent components.
func ParsePagoPAQRCodePayload(payload string) (iuv string, cfEnte string, amountCents int64, err error) {
	parts := strings.Split(payload, "|")
	if len(parts) != 5 || parts[0] != "PAGOPA" || parts[1] != "002" {
		return "", "", 0, errors.New("invalid PagoPA QR payload format")
	}

	iuv = parts[2]
	cfEnte = parts[3]
	amountCents, err = strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid amount in payload: %w", err)
	}

	return iuv, cfEnte, amountCents, nil
}

// GenerateBollettinoNotice builds the PagoPA notice data for printable PDF or in-app download.
func GenerateBollettinoNotice(payment *SchoolPayment, schoolCF, schoolName string) (*BollettinoNotice, error) {
	if payment == nil {
		return nil, errors.New("payment cannot be nil")
	}
	if payment.IUV == "" {
		// Auto-generate if missing
		iuv, err := GenerateIUV(0, "01", time.Now().UnixNano()%1000000000000)
		if err != nil {
			return nil, err
		}
		payment.IUV = iuv
	}

	amountCents := int64(payment.Amount * 100)
	qrPayload := GeneratePagoPAQRCodePayload(schoolCF, payment.IUV, amountCents)
	barcode128 := fmt.Sprintf("(415)%s(8020)%s(3902)%08d", schoolCF, payment.IUV, amountCents)

	causale := fmt.Sprintf("%s - %s", payment.Title, payment.StudentName)
	if len(causale) > 140 {
		causale = causale[:140]
	}

	return &BollettinoNotice{
		PaymentID:         payment.ID,
		IUV:               payment.IUV,
		SchoolCF:          schoolCF,
		SchoolName:        schoolName,
		StudentName:       payment.StudentName,
		Title:             payment.Title,
		Amount:            payment.Amount,
		DueDate:           payment.DueDate,
		QRCodePayload:     qrPayload,
		Barcode128:        barcode128,
		CausaleVersamento: causale,
		GeneratedAt:       time.Now(),
	}, nil
}
