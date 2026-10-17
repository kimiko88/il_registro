package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/albopretorio"
	"registro-backend/internal/interpelli"
	"registro-backend/internal/inventory"
	"registro-backend/internal/maturita"
	"registro-backend/internal/meals"
	"registro-backend/internal/payments"
	"registro-backend/internal/privacy"
	"registro-backend/internal/psychology"
	"registro-backend/internal/sidi"
	"registro-backend/internal/signatures"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper router setup for enterprise dealbreakers integration tests
func setupEnterpriseIntegrationTestRouter() (
	*gin.Engine,
	*mockPaymentsRepo,
	*payments.Service,
	*interpelli.Service,
	*albopretorio.Service,
	*signatures.Handler,
	*maturita.Service,
	*meals.Service,
	*inventory.Service,
	*privacy.Service,
	sidi.Service,
	*psychology.Service,
) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// 1. Payments
	payRepo := newMockPaymentsRepo()
	paySvc := payments.NewService(payRepo, &mockUserRepoForPayments{})
	payH := payments.NewHandler(paySvc)

	// 2. Interpelli
	intSvc := interpelli.NewService(nil)
	intH := interpelli.NewHandler(intSvc)

	// 3. Albo Pretorio
	alboSvc := albopretorio.NewService(nil)
	alboH := albopretorio.NewHandler(alboSvc)

	// 4. Signatures
	sigRepo := &mockSigRepo{}
	feqRepo := newMockFeqRepo()
	userLookup := &mockUserLookupForSig{mfaEnabled: true, mfaSecret: "JBSWY3DPEHPK3PXP"}
	docSvc := &mockDocServiceForSig{}
	baseSigSvc := signatures.NewService(sigRepo, docSvc, userLookup)
	feqSvc := signatures.NewQualifiedServiceWithProviders(nil, feqRepo, userLookup, &signatures.SoftwareQSCD{}, &mockTSAClient{})
	sidiExportSvc := signatures.NewSidiExportService()
	sigH := signatures.NewHandler(baseSigSvc, feqSvc, sidiExportSvc)

	// 5. Maturita
	matSvc := maturita.NewService(nil)
	matH := maturita.NewHandler(matSvc)

	// 6. Meals
	mealsSvc := meals.NewService(nil)
	mealsH := meals.NewHandler(mealsSvc)

	// 7. Inventory
	invSvc := inventory.NewService(nil)
	invH := inventory.NewHandler(invSvc)

	// 8. Privacy
	privSvc := privacy.NewService(nil)
	privH := privacy.NewHandler(privSvc)

	// 9. SIDI
	sidiSvc := sidi.NewService(nil)
	sidiH := sidi.NewHandler(sidiSvc)

	// 10. Psychology
	psySvc := psychology.NewService(nil)
	psyH := psychology.NewHandler(psySvc)

	authMiddleware := func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = "admin"
		}
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "admin-1"
		}
		schoolID := c.GetHeader("X-School-ID")
		if schoolID == "" {
			schoolID = "school-1"
		}
		c.Set("role", role)
		c.Set("user_id", userID)
		c.Set("school_id", schoolID)
		c.Next()
	}

	api := r.Group("/api/v1")

	// Public routes
	intH.RegisterPublicRoutes(api)
	alboH.RegisterPublicRoutes(api)
	sigH.RegisterPublicRoutes(api)

	// Protected routes
	protected := api.Group("", authMiddleware)
	{
		payH.RegisterRoutes(protected)
		intH.RegisterProtectedRoutes(protected)
		alboH.RegisterProtectedRoutes(protected)
		sigH.RegisterRoutes(protected)
		matH.RegisterRoutes(protected)
		mealsH.RegisterRoutes(protected)
		invH.RegisterRoutes(protected)
		privH.RegisterRoutes(protected)
		sidiH.RegisterRoutes(protected)
		psyH.RegisterRoutes(protected)
	}

	return r, payRepo, paySvc, intSvc, alboSvc, sigH, matSvc, mealsSvc, invSvc, privSvc, sidiSvc, psySvc
}

// ─────────────────────────────────────────────────────────────────────────────
// 1. PagoPA & OPI Integration Test
// ─────────────────────────────────────────────────────────────────────────────
func TestIntegration_PagoPA_And_OPI_Reconciliation(t *testing.T) {
	router, payRepo, _, _, _, _, _, _, _, _, _, _ := setupEnterpriseIntegrationTestRouter()

	ctx := context.Background()

	// 1. Pre-populate a SchoolPayment with IUV
	paymentID := "pay-integration-101"
	iuv, err := payments.GenerateIUV(0, "12", 100000000142)
	require.NoError(t, err)
	qrCodePayload := payments.GeneratePagoPAQRCodePayload("80012345678", iuv, int64(125.50*100))

	payment := &payments.SchoolPayment{
		ID:            paymentID,
		SchoolID:      "school-1",
		StudentID:     "student-100",
		Amount:        125.50,
		Title:         "Assicurazione Scolastica e Quota Gite 2024/2025",
		Status:        "pending",
		IUV:           iuv,
		QRCodePayload: qrCodePayload,
		DueDate:       time.Now().Add(30 * 24 * time.Hour),
	}
	err = payRepo.Create(ctx, payment)
	require.NoError(t, err)

	// 2. Fetch Bollettino Notice
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/payments/"+paymentID+"/bollettino", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var bollettinoResp payments.BollettinoNotice
	err = json.Unmarshal(w.Body.Bytes(), &bollettinoResp)
	require.NoError(t, err)
	assert.Equal(t, iuv, bollettinoResp.IUV)
	assert.Contains(t, bollettinoResp.QRCodePayload, iuv)

	// 3. Online Direct Checkout Session
	checkoutReq := map[string]string{
		"return_url": "https://school.example.edu/payments/return",
	}
	body, _ := json.Marshal(checkoutReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/payments/"+paymentID+"/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var checkoutResp payments.CheckoutSession
	err = json.Unmarshal(w.Body.Bytes(), &checkoutResp)
	require.NoError(t, err)
	assert.NotEmpty(t, checkoutResp.RedirectURL)
	assert.Equal(t, paymentID, checkoutResp.PaymentID)

	// 4. OPI Reconciliation Stream (XML Bank Output)
	opiXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<flusso_giornale_di_cassa>
    <codice_abi_bt>01005</codice_abi_bt>
    <identificativo_flusso>OPI-2026-00981</identificativo_flusso>
    <data_ora_creazione>2026-10-07T10:00:00</data_ora_creazione>
    <dettaglio_quietanza>
        <numero_quietanza>Q1001</numero_quietanza>
        <data_accredito>2026-10-06</data_accredito>
        <importo>125.50</importo>
        <identificativo_univoco_versamento>%s</identificativo_univoco_versamento>
        <codice_fiscale_debitore>RSSMRA80A01H501U</codice_fiscale_debitore>
        <anagrafica_debitore>Mario Rossi</anagrafica_debitore>
        <esito_accredito>ESEGUITO</esito_accredito>
    </dettaglio_quietanza>
</flusso_giornale_di_cassa>`, iuv)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/payments/reconcile-opi?format=OPI_XML", bytes.NewBufferString(opiXML))
	req.Header.Set("Content-Type", "application/xml")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var report payments.ReconciliationReport
	err = json.Unmarshal(w.Body.Bytes(), &report)
	require.NoError(t, err)
	assert.Equal(t, 1, report.TotalProcessed)
	assert.Equal(t, 1, report.TotalReconciled)
	assert.Equal(t, 125.50, report.TotalAmount)

	// Verify Payment is marked paid
	p, err := payRepo.GetByID(ctx, paymentID)
	require.NoError(t, err)
	assert.Equal(t, "paid", p.Status)
	assert.NotNil(t, p.PaidAt)

	// 5. Overdue solleciti listing
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/payments/solleciti", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─────────────────────────────────────────────────────────────────────────────
// 2. Interpelli Supplenze (O.M. 88/2024) Integration Test
// ─────────────────────────────────────────────────────────────────────────────
func TestIntegration_Interpelli_Supplenze(t *testing.T) {
	router, _, _, _, _, _, _, _, _, _, _, _ := setupEnterpriseIntegrationTestRouter()

	// 1. Create Interpello Notice
	noticeReq := interpelli.CreateNoticeRequest{
		Title:         "Docente Matematica Scuola Secondaria di II Grado",
		ConcorsoClass: "A026",
		PostType:      "comune",
		WeeklyHours:   18,
		StartDate:     "2026-11-01",
		EndDate:       "2027-06-30",
		Deadline:      "2026-10-31",
		Description:   "Supplenza per cattedra intera A026",
	}
	body, _ := json.Marshal(noticeReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/interpelli", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var notice interpelli.InterpelloNotice
	err := json.Unmarshal(w.Body.Bytes(), &notice)
	require.NoError(t, err)
	assert.NotEmpty(t, notice.ID)
	assert.Equal(t, "A026", notice.ConcorsoClass)

	// 2. Public candidate applies with DPR 445 flag
	candReq := interpelli.SubmitCandidaturaRequest{
		CandidateName:    "Elena",
		CandidateSurname: "Conti",
		FiscalCode:       "CNTLNE88D41H501Y",
		Email:            "elena.conti@example.com",
		Phone:            "+39 340 1234567",
		GraduationGrade:  110.0,
		GraduationLode:   true,
		HasHabilitation:  true,
		MonthsOfService:  12,
		CertLanguage:     "C1",
		DPR445Declared:   true,
	}
	body, _ = json.Marshal(candReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/public/interpelli/"+notice.ID+"/candidatura", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var cand interpelli.InterpelloCandidatura
	err = json.Unmarshal(w.Body.Bytes(), &cand)
	require.NoError(t, err)
	assert.NotEmpty(t, cand.ID)
	assert.True(t, cand.TotalScore > 0)

	// 3. Admin views Graduatoria
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/interpelli/"+notice.ID+"/graduatoria", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var gradRes struct {
		NoticeID    string                             `json:"notice_id"`
		Total       int                                `json:"total"`
		Graduatoria []interpelli.InterpelloCandidatura `json:"graduatoria"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &gradRes)
	require.NoError(t, err)
	require.Len(t, gradRes.Graduatoria, 1)
	assert.Equal(t, cand.ID, gradRes.Graduatoria[0].ID)

	// 4. Admin convokes candidate
	convocaReq := interpelli.ConvocaRequest{
		HoursToRespond: 24,
	}
	body, _ = json.Marshal(convocaReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/interpelli/candidature/"+cand.ID+"/convoca", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. Candidate responds
	respReq := interpelli.RispondiRequest{
		Risposta: "accettata",
		Notes:    "Accetto la proposta di supplenza",
	}
	body, _ = json.Marshal(respReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/interpelli/candidature/"+cand.ID+"/risposta", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─────────────────────────────────────────────────────────────────────────────
// 3. Albo Pretorio Online (L. 69/2009 & D.Lgs. 33/2013) Integration Test
// ─────────────────────────────────────────────────────────────────────────────
func TestIntegration_Albo_Pretorio_And_ANAC(t *testing.T) {
	router, _, _, _, _, _, _, _, _, _, _, _ := setupEnterpriseIntegrationTestRouter()

	// 1. Publish Act
	pubReq := albopretorio.PublishActRequest{
		Category:              albopretorio.CategoryBandiGarePNRR,
		Subject:               "Determina a contrarre affidamento fornitura arredi innovativi PNRR 4.0",
		DocumentFileURL:       "https://storage.scuola.edu.it/atti/determina-pnrr-01.pdf",
		DocumentSHA256:        "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		DaysDuration:          15,
		IsTransparencySection: true,
		CIGCode:               "B12948271A",
		AwardedAmount:         45000.0,
	}
	body, _ := json.Marshal(pubReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/albo-pretorio", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var atto albopretorio.AlboItem
	err := json.Unmarshal(w.Body.Bytes(), &atto)
	require.NoError(t, err)
	assert.NotEmpty(t, atto.ID)
	assert.Equal(t, "in_pubblicazione", atto.Status)

	// 2. Public checks Albo Pretorio (without auth)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/public/albo-pretorio?school_id=school-1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var listRes struct {
		Total int                     `json:"total"`
		Items []albopretorio.AlboItem `json:"items"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &listRes)
	require.NoError(t, err)
	assert.NotEmpty(t, listRes.Items)

	// 3. Defiggi Atto & generate Relata di Pubblicazione (using force=true to simulate expiration)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/albo-pretorio/"+atto.ID+"/defissione?force=true", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var defisso albopretorio.AlboItem
	err = json.Unmarshal(w.Body.Bytes(), &defisso)
	require.NoError(t, err)
	assert.Equal(t, "defisso", defisso.Status)
	assert.NotEmpty(t, defisso.RelataText)

	// 4. ANAC XML Export
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/public/amministrazione-trasparente/anac.xml?school_id=school-1&anno=%d", time.Now().Year()), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/xml; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "<legge190:dati")
	assert.Contains(t, w.Body.String(), "<cig>B12948271A</cig>")
}

// ─────────────────────────────────────────────────────────────────────────────
// 4. FEQ CSC Remote Batch Sign & Digital Stamp (CAD Art. 23) Integration Test
// ─────────────────────────────────────────────────────────────────────────────
func TestIntegration_CSC_Batch_And_Digital_Stamp(t *testing.T) {
	router, _, _, _, _, _, _, _, _, _, _, _ := setupEnterpriseIntegrationTestRouter()

	// 1. CSC Batch Signature
	batchReq := signatures.CSCBatchSignRequest{
		DocumentIDs: []string{"doc-pagella-1", "doc-pagella-2"},
		PIN:         "123456",
		OTP:         "987654",
	}
	body, _ := json.Marshal(batchReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/signatures/csc/batch-sign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var batchResp signatures.CSCBatchSignResponse
	err := json.Unmarshal(w.Body.Bytes(), &batchResp)
	require.NoError(t, err)
	assert.Equal(t, 2, batchResp.TotalRequested)
	assert.Equal(t, 2, batchResp.TotalSigned)

	// 2. Create Digital Stamp (Glifo CAD Art. 23)
	stampReq := map[string]string{
		"document_id":     "doc-pagella-1",
		"document_type":   "PAGELLA_ELETTRONICA",
		"document_sha256": "4b227777d4dd1fc61c6f884f48641d02b4d121d3fd328cb08b5531fcacdabf8a",
		"signer_name":     "Dott.ssa Laura Merini",
		"signer_role":     "Dirigente Scolastico",
	}
	body, _ = json.Marshal(stampReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/signatures/digital-stamp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var stamp signatures.DigitalStamp
	err = json.Unmarshal(w.Body.Bytes(), &stamp)
	require.NoError(t, err)
	assert.NotEmpty(t, stamp.GlifoToken)
	assert.NotEmpty(t, stamp.QRCodePayload)

	// 3. Verify Digital Stamp via Public Token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/public/verifica-glifo/"+stamp.GlifoToken, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var tokenRes map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &tokenRes)
	require.NoError(t, err)
	assert.Equal(t, "VALIDO_CONFORME_CAD_ART_23", tokenRes["status"])

	// 4. Verify Digital Stamp directly with payload & sha256
	verifyDirectReq := map[string]string{
		"payload":         stamp.QRCodePayload,
		"original_sha256": "4b227777d4dd1fc61c6f884f48641d02b4d121d3fd328cb08b5531fcacdabf8a",
	}
	body, _ = json.Marshal(verifyDirectReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/signatures/digital-stamp/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var verif signatures.DigitalStampVerification
	err = json.Unmarshal(w.Body.Bytes(), &verif)
	require.NoError(t, err)
	assert.True(t, verif.IsValid)
	assert.Equal(t, "doc-pagella-1", verif.DocumentID)
}

// ─────────────────────────────────────────────────────────────────────────────
// 5. Maturità Esame di Stato II Ciclo & Curriculum dello Studente
// ─────────────────────────────────────────────────────────────────────────────
func TestIntegration_Maturita_And_Curriculum(t *testing.T) {
	router, _, _, _, _, _, _, _, _, _, _, _ := setupEnterpriseIntegrationTestRouter()

	// 1. Save Commission
	commReq := maturita.CommissioneMaturita{
		CommissionCode: "RMIS09900B-COMM01",
		SchoolYear:     "2024/2025",
		ClassID:        "class-5A",
		PresidentName:  "Prof. Giovanni Valenti",
		Commissioners: []maturita.Commissioner{
			{Name: "Maria Rossi", Role: "commissario_interno", Subject: "Italiano"},
			{Name: "Roberto Verdi", Role: "commissario_esterno", Subject: "Matematica"},
		},
	}
	body, _ := json.Marshal(commReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/maturita/commission", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. Calculate Credits via Allegato A Table
	calcReq := map[string]interface{}{
		"grade_3rd": 8.5,
		"high_3rd":  true,
		"grade_4th": 8.2,
		"high_4th":  true,
		"grade_5th": 8.7,
		"high_5th":  true,
	}
	body, _ = json.Marshal(calcReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/maturita/calculate-credits", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var credits maturita.TrienniumCredits
	err := json.Unmarshal(w.Body.Bytes(), &credits)
	require.NoError(t, err)
	assert.True(t, credits.TotalCredits > 0)

	// 3. Save Student Record and Exam Scores
	studentRec := maturita.StudentMaturitaRecord{
		StudentID:  "student-mat-1",
		SchoolID:   "school-1",
		ClassID:    "class-5A",
		SchoolYear: "2024/2025",
		Credits:    credits,
		Scores: maturita.ExamScores{
			Written1Score: 18.5,
			Written2Score: 19.0,
			OralScore:     19.5,
			BonusPoints:   2,
			FinalScore:    95,
		},
		CurriculumStudente: &maturita.CurriculumStudente{
			StudentCF:       "RSSMRA06A01H501Z",
			StudentFullName: "Mario Rossi",
			DiplomaTitle:    "Diploma di Liceo Scientifico",
			FinalScore:      95,
			PCTOHoursTotal:  120,
		},
	}
	body, _ = json.Marshal(studentRec)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/maturita/student-record", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Export Curriculum dello Studente XML
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/maturita/curriculum-xml?student_id=student-mat-1&school_year=2024/2025", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/xml; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "<curriculum_dello_studente")
}

// ─────────────────────────────────────────────────────────────────────────────
// 6. Refezione Scolastica & Gestione Diete Speciali
// ─────────────────────────────────────────────────────────────────────────────
func TestIntegration_Meals_Refezione_And_Diete(t *testing.T) {
	router, _, _, _, _, _, _, _, _, _, _, _ := setupEnterpriseIntegrationTestRouter()

	// 1. Register Special Diet
	dietReq := meals.SpecialDietRequest{
		StudentID:             "student-celiac-01",
		DietCategory:          "sanitaria",
		SpecificDiet:          "celiachia",
		MedicalCertificateURL: "https://storage.scuola.edu.it/certificati/med-celiac-01.pdf",
		CertificateExpiryDate: "2026-06-30",
		AllergensList:         "glutine",
		Notes:                 "Celiachia accertata",
	}
	body, _ := json.Marshal(dietReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/meals/special-diet", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Family Wallet Top-Up (55.00)
	topupReq := map[string]interface{}{
		"amount":     55.00,
		"pagopa_iuv": "00100000000000142",
	}
	body, _ = json.Marshal(topupReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/meals/wallet/student-celiac-01/topup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var wallet meals.MealWallet
	err := json.Unmarshal(w.Body.Bytes(), &wallet)
	require.NoError(t, err)
	assert.Equal(t, 55.00, wallet.Balance)

	// 3. Teacher takes Morning Roll Call (Debits 5.50)
	rollReq := meals.MealRollCallRequest{
		ClassID:    "class-2A",
		Date:       "2026-10-07",
		PresentIDs: []string{"student-celiac-01", "student-standard-02", "student-standard-03"},
	}
	body, _ = json.Marshal(rollReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/meals/roll-call", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var cateringReport meals.CateringOrderReport
	err = json.Unmarshal(w.Body.Bytes(), &cateringReport)
	require.NoError(t, err)
	assert.Equal(t, 3, cateringReport.TotalMeals)
	assert.Equal(t, 1, cateringReport.TotalDiets)

	// 4. Wallet after roll call: 55.00 - 5.50 = 49.50
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/meals/wallet/student-celiac-01", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &wallet)
	require.NoError(t, err)
	assert.Equal(t, 49.50, wallet.Balance)
}

// ─────────────────────────────────────────────────────────────────────────────
// 7. Inventario Scolastico, Cespiti e Comodato d'Uso (D.I. 129/2018)
// ─────────────────────────────────────────────────────────────────────────────
func TestIntegration_Inventory_And_Device_Loan(t *testing.T) {
	router, _, _, _, _, _, _, _, _, _, _, _ := setupEnterpriseIntegrationTestRouter()

	// 1. Create School Asset
	assetReq := inventory.CreateAssetRequest{
		Category:         inventory.CategoryInformatica,
		Description:      "Notebook Lenovo ThinkPad PNRR 4.0",
		SerialNumber:     "LEN-SN-449102",
		BuildingLocation: "Sede Centrale",
		RoomLocation:     "Laboratorio Informatica 2",
		AcquisitionDate:  "2026-01-15",
		InitialValue:     650.00,
		DepreciationRate: 20.0,
	}
	body, _ := json.Marshal(assetReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/inventory/assets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var asset inventory.AssetItem
	err := json.Unmarshal(w.Body.Bytes(), &asset)
	require.NoError(t, err)
	assert.NotEmpty(t, asset.ID)
	assert.NotEmpty(t, asset.InventoryNumber)

	// 2. Barcode Label
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/inventory/assets/"+asset.ID+"/label", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. Device Loan Contract
	loanReq := inventory.CreateLoanRequest{
		AssetID:              asset.ID,
		StudentID:            "student-loan-01",
		ParentUserID:         "parent-01",
		ExpectedReturnDate:   "2027-06-15",
		DeviceConditionNotes: "Perfette condizioni con alimentatore originale",
	}
	body, _ = json.Marshal(loanReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/inventory/loans", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var loan inventory.LoanContract
	err = json.Unmarshal(w.Body.Bytes(), &loan)
	require.NoError(t, err)
	assert.Equal(t, "attivo", loan.Status)

	// 4. Return Device
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/inventory/loans/"+loan.ID+"/return", bytes.NewBufferString(`{"condition":"ottime"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─────────────────────────────────────────────────────────────────────────────
// 8. Registro Trattamenti Privacy & Semaforo Consensi GDPR
// ─────────────────────────────────────────────────────────────────────────────
func TestIntegration_Privacy_Treatments_And_Traffic_Light(t *testing.T) {
	router, _, _, _, _, _, _, _, _, _, _, _ := setupEnterpriseIntegrationTestRouter()

	// 1. Register GDPR Treatment
	treatReq := privacy.TreatmentRequest{
		ActivityName:         "Pubblicazione foto/video attività didattiche su portale istituzionale",
		LegalBasis:           "consenso",
		InterestedCategories: "Studenti minori di anni 14 e 18",
		RetentionPeriod:      "Durata del ciclo scolastico",
		DPOContact:           "dpo@scuola.edu.it",
		SecurityMeasures:     "Misure adeguate di pseudonimizzazione e crittografia",
	}
	body, _ := json.Marshal(treatReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/privacy/treatments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Save Family Consents
	consentReq := privacy.SaveConsentRequest{
		StudentID:               "student-priv-01",
		SchoolYear:              "2024/2025",
		PhotoVideoSocialConsent: true,
		CloudWorkspaceConsent:   true,
		WalkingTripsConsent:     true,
	}
	body, _ = json.Marshal(consentReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/privacy/consent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. Class Traffic Light Badge Query
	badgeReq := map[string]interface{}{
		"student_ids": []string{"student-priv-01"},
		"school_year": "2024/2025",
	}
	body, _ = json.Marshal(badgeReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/privacy/class-badges", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Badges map[string]privacy.TrafficLightBadge `json:"badges"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	badge, exists := res.Badges["student-priv-01"]
	require.True(t, exists)
	assert.Equal(t, privacy.TrafficLightGreen, badge)
}

// ─────────────────────────────────────────────────────────────────────────────
// 9. Cooperazione Applicativa Diretta con i WebService SIDI (MIM)
// ─────────────────────────────────────────────────────────────────────────────
func TestIntegration_SIDI_Cooperation(t *testing.T) {
	router, _, _, _, _, _, _, _, _, _, _, _ := setupEnterpriseIntegrationTestRouter()

	// 1. Get Cooperation Config
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/sidi/cooperation-config", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var cfg sidi.SidiCooperationConfig
	err := json.Unmarshal(w.Body.Bytes(), &cfg)
	require.NoError(t, err)
	assert.NotEmpty(t, cfg.EndpointURL)
	assert.NotEmpty(t, cfg.CertificatoPostazione)

	// 2. Synchronize SIDI Student Codes
	syncReq := map[string]interface{}{
		"students": []sidi.StudenteSIDI{
			{CodiceFiscale: "BNCGLI08B41H501A", Cognome: "Bianchi", Nome: "Giulia", Classe: "3ª A"},
		},
	}
	body, _ := json.Marshal(syncReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/sidi/sync-student-codes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var syncResp sidi.SyncSidiCodesResponse
	err = json.Unmarshal(w.Body.Bytes(), &syncResp)
	require.NoError(t, err)
	assert.Equal(t, 1, syncResp.TotalProcessed)
	assert.Equal(t, 1, syncResp.TotalUpdated)

	// 3. 1-Click Push Scrutiny Results
	scrutinyReq := sidi.PushScrutinyResultsRequest{
		ClassID:    "class-3A",
		SchoolYear: "2024/2025",
		Results: []sidi.ScrutinioSIDI{
			{CodiceSIDI: "SIDI-1000001", Classe: "3ª A", EsitoFinale: "AMMESSO", CreditiFormat: 12},
		},
	}
	body, _ = json.Marshal(scrutinyReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/sidi/push-scrutiny-results", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var pushResp sidi.PushScrutinyResultsResponse
	err = json.Unmarshal(w.Body.Bytes(), &pushResp)
	require.NoError(t, err)
	assert.Equal(t, "TRASMESSO_CON_SUCCESSO", pushResp.Status)
	assert.NotEmpty(t, pushResp.ProtocolloMIM)
}

// ─────────────────────────────────────────────────────────────────────────────
// 10. Sportello d'Ascolto Psicologico (CIC) & Segreto Professionale L. 56/1989
// ─────────────────────────────────────────────────────────────────────────────
func TestIntegration_Psychology_Secrecy_And_Dual_Consent(t *testing.T) {
	router, _, _, _, _, _, _, _, _, _, _, _ := setupEnterpriseIntegrationTestRouter()

	slotTime := time.Now().Add(48 * time.Hour).Truncate(time.Minute)

	// 1. Attempt minor booking WITHOUT dual parent consent -> Must Fail
	bookReq := psychology.BookSessionRequest{
		SlotTime:        slotTime.Format("2006-01-02 15:04"),
		DurationMinutes: 45,
		PsychologistID:  "psy-01",
	}
	body, _ := json.Marshal(bookReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/psychology/book", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Role", "student")
	req.Header.Set("X-User-ID", "minor-student-1")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "entrambi i genitori")

	// 2. Both parents sign informed consent
	consentReq := psychology.SaveConsentRequest{
		StudentID:  "minor-student-1",
		SchoolYear: "2025/2026",
	}
	body, _ = json.Marshal(consentReq)

	// Parent 1
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/psychology/consent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Role", "parent")
	req.Header.Set("X-User-ID", "parent-1")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Parent 2
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/psychology/consent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Role", "parent")
	req.Header.Set("X-User-ID", "parent-2")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. Now minor booking succeeds with Anonymous Alias
	body, _ = json.Marshal(bookReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/psychology/book", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Role", "student")
	req.Header.Set("X-User-ID", "minor-student-1")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var session psychology.PsychologySession
	err := json.Unmarshal(w.Body.Bytes(), &session)
	require.NoError(t, err)
	assert.NotEmpty(t, session.ID)
	assert.NotEmpty(t, session.AnonymousAlias)

	// 4. Clinical notes update rejected for principal/admin (L. 56/1989 Secrecy)
	notesReq := map[string]string{
		"notes": "Note riservate di colloquio clinico",
	}
	body, _ = json.Marshal(notesReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/psychology/sessions/"+session.ID+"/notes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Role", "principal")
	req.Header.Set("X-User-ID", "ds-1")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 5. Clinical notes update allowed for Psychologist role
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/psychology/sessions/"+session.ID+"/notes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Role", "psychologist")
	req.Header.Set("X-User-ID", "psy-01")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
