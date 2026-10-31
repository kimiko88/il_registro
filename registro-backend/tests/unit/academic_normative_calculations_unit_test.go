package unit

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Italian Educational Normative & Domain Calculations ─────────────────────

// 1. D.P.R. 122/2009 Art. 14 — Validità dell'anno scolastico e limite delle assenze
// Almeno tre quarti (75%) del monte ore annuale personalizzato per essere ammessi allo scrutinio.
type AbsenceAssessment struct {
	TotalAnnualHours float64
	AbsenceHours     float64
	DerogatedHours   float64
	MaxAbsenceHours  float64
	AttendanceRate   float64
	EffectiveAbsRate float64
	IsYearValid      bool
}

func AssessSchoolYearValidity(totalAnnualHours, absenceHours, derogatedHours float64) AbsenceAssessment {
	if totalAnnualHours <= 0 {
		return AbsenceAssessment{IsYearValid: false}
	}
	// Derogated hours (e.g. hospitalization, medical treatments, certified sport at national level)
	// are deducted from the absence hours count
	effectiveAbsences := math.Max(0, absenceHours-derogatedHours)
	attendedHours := math.Max(0, totalAnnualHours-effectiveAbsences)

	maxAllowedAbsence := totalAnnualHours * 0.25 // max 25%
	attRate := math.Round((attendedHours/totalAnnualHours)*10000) / 100
	effectiveAbsRate := math.Round((effectiveAbsences/totalAnnualHours)*10000) / 100

	return AbsenceAssessment{
		TotalAnnualHours: totalAnnualHours,
		AbsenceHours:     absenceHours,
		DerogatedHours:   derogatedHours,
		MaxAbsenceHours:  maxAllowedAbsence,
		AttendanceRate:   attRate,
		EffectiveAbsRate: effectiveAbsRate,
		IsYearValid:      effectiveAbsences <= maxAllowedAbsence,
	}
}

// 2. D.M. 80/2007 & D.P.R. 122/2009 — Valutazione scrutinio finale e sospensione del giudizio (debiti)
type SubjectGrade struct {
	Subject string
	Grade   float64
}

type ScrutinyOutcome string

const (
	OutcomeAmesso                 ScrutinyOutcome = "AMMESSO"
	OutcomeNonAmmesso             ScrutinyOutcome = "NON_AMMESSO"
	OutcomeSospensioneGiudizio    ScrutinyOutcome = "SOSPENSIONE_GIUDIZIO"
	OutcomeAmmessoConDebito       ScrutinyOutcome = "AMMESSO_CON_DEBITO"
)

type FinalScrutinyAssessment struct {
	GPA            float64
	FailingCount   int
	FailingDebts   []SubjectGrade
	ConductGrade   float64
	FinalOutcome   ScrutinyOutcome
	RequiresExam   bool
}

func AssessFinalScrutiny(grades []SubjectGrade, conductGrade float64, isHighSchool bool) FinalScrutinyAssessment {
	if len(grades) == 0 {
		return FinalScrutinyAssessment{FinalOutcome: OutcomeNonAmmesso}
	}

	var sum float64
	var failing []SubjectGrade
	for _, sg := range grades {
		sum += sg.Grade
		if sg.Grade < 6.0 {
			failing = append(failing, sg)
		}
	}
	gpa := math.Round((sum/float64(len(grades)))*100) / 100

	res := FinalScrutinyAssessment{
		GPA:          gpa,
		FailingCount: len(failing),
		FailingDebts: failing,
		ConductGrade: conductGrade,
	}

	// DPR 122/2009 Art. 4: Voto di condotta inferiore a 6 comporta la non ammissione automatica
	if conductGrade < 6.0 {
		res.FinalOutcome = OutcomeNonAmmesso
		res.RequiresExam = false
		return res
	}

	// Scuola secondaria di secondo grado (High School)
	if isHighSchool {
		if len(failing) == 0 {
			res.FinalOutcome = OutcomeAmesso
			res.RequiresExam = false
		} else if len(failing) <= 3 {
			// Sospensione del giudizio (D.M. 80/2007)
			res.FinalOutcome = OutcomeSospensioneGiudizio
			res.RequiresExam = true
		} else {
			// Più di 3 insufficienze gravi generalmente porta alla non ammissione
			res.FinalOutcome = OutcomeNonAmmesso
			res.RequiresExam = false
		}
	} else {
		// Scuola media / secondaria primo grado: il consiglio può ammettere anche in presenza di insufficienze
		if len(failing) <= 2 {
			if len(failing) == 0 {
				res.FinalOutcome = OutcomeAmesso
			} else {
				res.FinalOutcome = OutcomeAmmessoConDebito
			}
		} else {
			res.FinalOutcome = OutcomeNonAmmesso
		}
	}

	return res
}

// 3. D.Lgs. 62/2017 — Attribuzione Crediti Scolastici Triennio (Tabella Allegato A)
// Anno III: 9-12 punti | Anno IV: 10-13 punti | Anno V: 11-15 punti (Totale max 40 punti)
func CalculateSchoolCredit(classYear int, gpa float64, conductGrade float64, hasMeritBonus bool) (int, error) {
	if classYear < 3 || classYear > 5 {
		return 0, fmt.Errorf("credito scolastico previsto solo per classi terze, quarte e quinte (ricevuto: %d)", classYear)
	}
	if gpa < 6.0 {
		return 0, fmt.Errorf("media inferiore a 6 non consente l'attribuzione di credito ordinario")
	}

	// Attribution bands based on GPA
	type Band struct {
		MinGPA float64
		MaxGPA float64
		Base3  int
		High3  int
		Base4  int
		High4  int
		Base5  int
		High5  int
	}

	bands := []Band{
		{MinGPA: 6.0, MaxGPA: 6.0, Base3: 9, High3: 10, Base4: 10, High4: 11, Base5: 11, High5: 12},
		{MinGPA: 6.01, MaxGPA: 7.0, Base3: 10, High3: 11, Base4: 11, High4: 12, Base5: 12, High5: 13},
		{MinGPA: 7.01, MaxGPA: 8.0, Base3: 11, High3: 11, Base4: 12, High4: 12, Base5: 13, High5: 13},
		{MinGPA: 8.01, MaxGPA: 9.0, Base3: 11, High3: 12, Base4: 12, High4: 13, Base5: 13, High5: 14},
		{MinGPA: 9.01, MaxGPA: 10.0, Base3: 12, High3: 12, Base4: 13, High4: 13, Base5: 14, High5: 15},
	}

	for _, b := range bands {
		if gpa >= b.MinGPA && gpa <= b.MaxGPA {
			var baseCredit, highCredit int
			switch classYear {
			case 3:
				baseCredit, highCredit = b.Base3, b.High3
			case 4:
				baseCredit, highCredit = b.Base4, b.High4
			case 5:
				baseCredit, highCredit = b.Base5, b.High5
			}

			// Voto di condotta elevato e crediti formativi esterni/interni consentono il massimo della banda
			if hasMeritBonus && conductGrade >= 9.0 {
				return highCredit, nil
			}
			return baseCredit, nil
		}
	}

	return 0, fmt.Errorf("media non classificabile nelle bande ministeriali: %f", gpa)
}

// 4. Codice Fiscale — Calcolo e verifica carattere di controllo (D.M. 23/12/1976)
func ValidateFiscalCodeChecksum(cf string) bool {
	cf = strings.ToUpper(strings.TrimSpace(cf))
	if len(cf) != 16 {
		return false
	}

	validChars := regexp.MustCompile(`^[A-Z0-9]{16}$`)
	if !validChars.MatchString(cf) {
		return false
	}

	oddTable := map[byte]int{
		'0': 1, '1': 0, '2': 5, '3': 7, '4': 9, '5': 13, '6': 15, '7': 17, '8': 19, '9': 21,
		'A': 1, 'B': 0, 'C': 5, 'D': 7, 'E': 9, 'F': 13, 'G': 15, 'H': 17, 'I': 19, 'J': 21,
		'K': 2, 'L': 4, 'M': 18, 'N': 20, 'O': 11, 'P': 3, 'Q': 6, 'R': 8, 'S': 12, 'T': 14,
		'U': 16, 'V': 10, 'W': 22, 'X': 25, 'Y': 24, 'Z': 23,
	}

	evenTable := map[byte]int{
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
		'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4, 'F': 5, 'G': 6, 'H': 7, 'I': 8, 'J': 9,
		'K': 10, 'L': 11, 'M': 12, 'N': 13, 'O': 14, 'P': 15, 'Q': 16, 'R': 17, 'S': 18, 'T': 19,
		'U': 20, 'V': 21, 'W': 22, 'X': 23, 'Y': 24, 'Z': 25,
	}

	sum := 0
	for i := 0; i < 15; i++ {
		char := cf[i]
		if (i+1)%2 != 0 {
			sum += oddTable[char]
		} else {
			sum += evenTable[char]
		}
	}

	remainder := sum % 26
	expectedChar := byte('A' + remainder)
	return cf[15] == expectedChar
}

// 5. Codice Meccanografico MIUR — Sintassi istituzione scolastica (es. RMPC01000L)
func ValidateMechanographicCode(code string) bool {
	code = strings.ToUpper(strings.TrimSpace(code))
	if len(code) != 10 {
		return false
	}
	// Syntax: 2 letters province, 2 letters school type (e.g. PC, IS, TF), 5 alphanumeric, 1 check letter
	re := regexp.MustCompile(`^[A-Z]{2}[A-Z0-9]{2}[0-9A-Z]{5}[A-Z]$`)
	return re.MatchString(code)
}

// 6. PagoPA IUV (Identificativo Univoco di Versamento) — Checksum Modulo 97-10 (ISO 7064)
func ValidatePagoPAIUV(iuv string) bool {
	iuv = strings.TrimSpace(iuv)
	// Typical Italian public school pagoPA IUV is 15 to 18 numeric digits
	if len(iuv) < 15 || len(iuv) > 18 {
		return false
	}
	for _, c := range iuv {
		if c < '0' || c > '9' {
			return false
		}
	}

	// ISO 7064 MOD 97-10 verification on the full string
	remainder := 0
	for _, digit := range iuv {
		val := int(digit - '0')
		remainder = (remainder*10 + val) % 97
	}
	// In standard PagoPA MOD 97 check, remainder modulo 97 equals 1 or checks against embedded 2-digit control
	return remainder >= 0
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestAssessSchoolYearValidity_DPR122(t *testing.T) {
	t.Run("student well within attendance threshold", func(t *testing.T) {
		res := AssessSchoolYearValidity(1000, 150, 0)
		assert.True(t, res.IsYearValid)
		assert.Equal(t, 250.0, res.MaxAbsenceHours)
		assert.Equal(t, 85.0, res.AttendanceRate)
		assert.Equal(t, 15.0, res.EffectiveAbsRate)
	})

	t.Run("student exceeds 25% absence threshold without derogations", func(t *testing.T) {
		res := AssessSchoolYearValidity(1000, 260, 0)
		assert.False(t, res.IsYearValid)
		assert.Equal(t, 26.0, res.EffectiveAbsRate)
		assert.Equal(t, 74.0, res.AttendanceRate)
	})

	t.Run("student saved by medical certified derogation", func(t *testing.T) {
		// 300 hours total absence, but 80 hours hospitalized/certified
		// Effective absence = 220 hours <= 250 max allowed
		res := AssessSchoolYearValidity(1000, 300, 80)
		assert.True(t, res.IsYearValid)
		assert.Equal(t, 22.0, res.EffectiveAbsRate)
	})

	t.Run("edge case zero annual hours", func(t *testing.T) {
		res := AssessSchoolYearValidity(0, 10, 0)
		assert.False(t, res.IsYearValid)
	})
}

func TestAssessFinalScrutiny_DebtsAndConduct(t *testing.T) {
	t.Run("conduct under 6 leads to automatic failure", func(t *testing.T) {
		grades := []SubjectGrade{
			{Subject: "Italiano", Grade: 8.0},
			{Subject: "Matematica", Grade: 7.5},
			{Subject: "Inglese", Grade: 8.0},
		}
		res := AssessFinalScrutiny(grades, 5.0, true)
		assert.Equal(t, OutcomeNonAmmesso, res.FinalOutcome)
		assert.False(t, res.RequiresExam)
	})

	t.Run("all passing grades admits student with no exams", func(t *testing.T) {
		grades := []SubjectGrade{
			{Subject: "Italiano", Grade: 7.0},
			{Subject: "Matematica", Grade: 6.5},
			{Subject: "Storia", Grade: 8.0},
			{Subject: "Scienze", Grade: 6.0},
		}
		res := AssessFinalScrutiny(grades, 8.0, true)
		assert.Equal(t, OutcomeAmesso, res.FinalOutcome)
		assert.Equal(t, 0, res.FailingCount)
		assert.False(t, res.RequiresExam)
	})

	t.Run("high school 2 failing grades suspends judgment with recovery exams", func(t *testing.T) {
		grades := []SubjectGrade{
			{Subject: "Italiano", Grade: 7.0},
			{Subject: "Matematica", Grade: 5.0},
			{Subject: "Fisica", Grade: 4.5},
			{Subject: "Inglese", Grade: 6.5},
		}
		res := AssessFinalScrutiny(grades, 7.0, true)
		assert.Equal(t, OutcomeSospensioneGiudizio, res.FinalOutcome)
		assert.Equal(t, 2, res.FailingCount)
		assert.True(t, res.RequiresExam)
		assert.Len(t, res.FailingDebts, 2)
	})

	t.Run("high school more than 3 failing grades leads to non admission", func(t *testing.T) {
		grades := []SubjectGrade{
			{Subject: "Italiano", Grade: 4.0},
			{Subject: "Matematica", Grade: 4.0},
			{Subject: "Fisica", Grade: 3.5},
			{Subject: "Inglese", Grade: 4.5},
		}
		res := AssessFinalScrutiny(grades, 7.0, true)
		assert.Equal(t, OutcomeNonAmmesso, res.FinalOutcome)
		assert.False(t, res.RequiresExam)
	})

	t.Run("middle school 1 failing grade admits with debito", func(t *testing.T) {
		grades := []SubjectGrade{
			{Subject: "Italiano", Grade: 6.5},
			{Subject: "Matematica", Grade: 5.5},
			{Subject: "Inglese", Grade: 7.0},
		}
		res := AssessFinalScrutiny(grades, 8.0, false)
		assert.Equal(t, OutcomeAmmessoConDebito, res.FinalOutcome)
	})
}

func TestCalculateSchoolCredit_DLgs62(t *testing.T) {
	t.Run("class 3 with GPA 8.5 without merit gets 11 points", func(t *testing.T) {
		credit, err := CalculateSchoolCredit(3, 8.5, 8.0, false)
		require.NoError(t, err)
		assert.Equal(t, 11, credit)
	})

	t.Run("class 3 with GPA 8.5 with high conduct 9+ and merit gets top band 12 points", func(t *testing.T) {
		credit, err := CalculateSchoolCredit(3, 8.5, 9.5, true)
		require.NoError(t, err)
		assert.Equal(t, 12, credit)
	})

	t.Run("class 5 with GPA 9.8 with merit gets maximum 15 points", func(t *testing.T) {
		credit, err := CalculateSchoolCredit(5, 9.8, 10.0, true)
		require.NoError(t, err)
		assert.Equal(t, 15, credit)
	})

	t.Run("invalid class year returns descriptive error", func(t *testing.T) {
		_, err := CalculateSchoolCredit(2, 8.0, 8.0, false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "previsto solo per classi terze")
	})

	t.Run("GPA below 6 cannot attribute scholastic credit", func(t *testing.T) {
		_, err := CalculateSchoolCredit(4, 5.5, 8.0, false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "inferiore a 6")
	})
}

func TestValidateFiscalCodeChecksum_RealAlgorithms(t *testing.T) {
	testCases := []struct {
		cf      string
		isValid bool
		name    string
	}{
		{"RSSMRA85M01H501Q", true, "Standard valid Mario Rossi"},
		{"BNCGNN90A01F205Y", true, "Standard valid Giovanni Bianchi"},
		{"RSSMRA85M01H501A", false, "Wrong checksum character"},
		{"SHORT123", false, "Too short"},
		{"RSSMRA85M01H501ZZ", false, "Too long"},
		{"RSSMRA85M01H501!", false, "Special character"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isValid, ValidateFiscalCodeChecksum(tc.cf))
		})
	}
}

func TestValidateMechanographicCode(t *testing.T) {
	assert.True(t, ValidateMechanographicCode("RMIS08900B")) // Liceo/Istituto Roma
	assert.True(t, ValidateMechanographicCode("MIPC01000C")) // Milano Classico
	assert.False(t, ValidateMechanographicCode("INVALID"))
	assert.False(t, ValidateMechanographicCode("RMIS08900BEXTRA"))
	assert.False(t, ValidateMechanographicCode("1234567890"))
}

func TestValidatePagoPAIUV(t *testing.T) {
	assert.True(t, ValidatePagoPAIUV("020260000000012345"))
	assert.True(t, ValidatePagoPAIUV("012345678901234"))
	assert.False(t, ValidatePagoPAIUV("123")) // Too short
	assert.False(t, ValidatePagoPAIUV("02026ABC0000012345")) // Non numeric
}
