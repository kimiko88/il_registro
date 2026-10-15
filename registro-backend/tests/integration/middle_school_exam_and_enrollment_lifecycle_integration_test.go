package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/enrollment"
	"registro-backend/internal/middleschoolexam"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockLifecycleEnrollmentRepo struct {
	apps   map[string]enrollment.EnrollmentApplication
	drafts map[string]*enrollment.ClassFormationDraft
}

func newMockLifecycleEnrollmentRepo() *mockLifecycleEnrollmentRepo {
	return &mockLifecycleEnrollmentRepo{
		apps:   make(map[string]enrollment.EnrollmentApplication),
		drafts: make(map[string]*enrollment.ClassFormationDraft),
	}
}

func (m *mockLifecycleEnrollmentRepo) SaveApplications(ctx context.Context, apps []enrollment.EnrollmentApplication, schoolID, academicYear string) (int, error) {
	for _, a := range apps {
		a.SchoolID = schoolID
		a.AcademicYear = academicYear
		m.apps[a.SidiApplicationID] = a
	}
	return len(apps), nil
}

func (m *mockLifecycleEnrollmentRepo) ListApplications(ctx context.Context, schoolID, academicYear, status string) ([]enrollment.EnrollmentApplication, error) {
	var list []enrollment.EnrollmentApplication
	for _, a := range m.apps {
		if a.SchoolID == schoolID && a.AcademicYear == academicYear {
			list = append(list, a)
		}
	}
	return list, nil
}

func (m *mockLifecycleEnrollmentRepo) SaveDraft(ctx context.Context, draft *enrollment.ClassFormationDraft) error {
	m.drafts[draft.ID] = draft
	return nil
}

func (m *mockLifecycleEnrollmentRepo) GetDraft(ctx context.Context, id string) (*enrollment.ClassFormationDraft, error) {
	return m.drafts[id], nil
}

func (m *mockLifecycleEnrollmentRepo) ListDrafts(ctx context.Context, schoolID string) ([]enrollment.ClassFormationDraft, error) {
	var list []enrollment.ClassFormationDraft
	for _, d := range m.drafts {
		if d.SchoolID == schoolID {
			list = append(list, *d)
		}
	}
	return list, nil
}

func (m *mockLifecycleEnrollmentRepo) UpdateDraftAssignments(ctx context.Context, id string, assignments enrollment.ClassFormationDraftResult) error {
	if d, ok := m.drafts[id]; ok {
		d.Assignments = assignments
	}
	return nil
}

func (m *mockLifecycleEnrollmentRepo) FinalizeDraft(ctx context.Context, id string) error {
	if d, ok := m.drafts[id]; ok {
		d.IsFinalized = true
	}
	return nil
}

type mockLifecycleExamRepo struct {
	exams      map[string]*middleschoolexam.MiddleSchoolExam
	candidates map[string]*middleschoolexam.MiddleSchoolExamCandidate
}

func newMockLifecycleExamRepo() *mockLifecycleExamRepo {
	return &mockLifecycleExamRepo{
		exams:      make(map[string]*middleschoolexam.MiddleSchoolExam),
		candidates: make(map[string]*middleschoolexam.MiddleSchoolExamCandidate),
	}
}

func (m *mockLifecycleExamRepo) GetOrCreateExamForClass(ctx context.Context, schoolID, classID, academicYear, president string) (*middleschoolexam.MiddleSchoolExam, error) {
	e := &middleschoolexam.MiddleSchoolExam{
		ID:                  "exam-sec-1",
		SchoolID:            schoolID,
		ClassID:             classID,
		AcademicYear:        academicYear,
		SubcommissionNumber: 1,
		PresidentName:       president,
		Status:              "admission",
		CreatedAt:           time.Now(),
	}
	m.exams[e.ID] = e
	return e, nil
}

func (m *mockLifecycleExamRepo) GetExamByID(ctx context.Context, id string) (*middleschoolexam.MiddleSchoolExam, error) {
	return m.exams[id], nil
}

func (m *mockLifecycleExamRepo) UpdateExamStatus(ctx context.Context, id, status string) error {
	if e, ok := m.exams[id]; ok {
		e.Status = status
	}
	return nil
}

func (m *mockLifecycleExamRepo) ListCandidates(ctx context.Context, examID string) ([]middleschoolexam.MiddleSchoolExamCandidate, error) {
	var list []middleschoolexam.MiddleSchoolExamCandidate
	for _, c := range m.candidates {
		if c.ExamID == examID {
			list = append(list, *c)
		}
	}
	return list, nil
}

func (m *mockLifecycleExamRepo) SaveCandidateAdmission(ctx context.Context, examID, studentID string, grade int, judgment string, admitted bool) error {
	key := examID + ":" + studentID
	m.candidates[key] = &middleschoolexam.MiddleSchoolExamCandidate{
		ID:                "cand-" + studentID,
		ExamID:            examID,
		StudentID:         studentID,
		StudentName:       "Candidato " + studentID,
		AdmissionGrade:    grade,
		AdmissionJudgment: judgment,
		IsAdmitted:        admitted,
		CreatedAt:         time.Now(),
	}
	return nil
}

func (m *mockLifecycleExamRepo) SaveCandidateGrades(ctx context.Context, c *middleschoolexam.MiddleSchoolExamCandidate) error {
	key := c.ExamID + ":" + c.StudentID
	m.candidates[key] = c
	return nil
}

func (m *mockLifecycleExamRepo) GetCandidate(ctx context.Context, examID, studentID string) (*middleschoolexam.MiddleSchoolExamCandidate, error) {
	key := examID + ":" + studentID
	if c, ok := m.candidates[key]; ok {
		return c, nil
	}
	return &middleschoolexam.MiddleSchoolExamCandidate{
		ID:             "cand-" + studentID,
		ExamID:         examID,
		StudentID:      studentID,
		StudentName:    "Candidato " + studentID,
		AdmissionGrade: 9,
		IsAdmitted:     true,
	}, nil
}

func (m *mockLifecycleExamRepo) GetCandidateDiplomaData(ctx context.Context, candidateID string) (*middleschoolexam.DiplomaData, error) {
	return &middleschoolexam.DiplomaData{
		SchoolName:      "Scuola Secondaria Statale Carducci",
		SchoolCode:      "RMIC822001",
		AcademicYear:    "2026/2027",
		StudentFullName: "Matteo Bianchi",
		BirthDate:       "15/05/2012",
		BirthPlace:      "Milano",
		FinalGrade:      10,
		HasHonors:       true,
		PresidentName:   "Dott.ssa Clara Bellini",
		IssueDate:       time.Now(),
	}, nil
}

func TestEnrollmentAndMiddleSchoolExam_FullLifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup Enrollment
	enrRepo := newMockLifecycleEnrollmentRepo()
	enrSvc := enrollment.NewService(enrRepo)
	enrH := enrollment.NewHandler(enrSvc)

	// Setup Middle School Exam
	examRepo := newMockLifecycleExamRepo()
	examSvc := middleschoolexam.NewService(examRepo)
	examH := middleschoolexam.NewHandler(examSvc)

	router := gin.New()
	api := router.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("school_id", "school-sec-1")
		c.Set("user_id", "dirigente-1")
		c.Set("role", "principal")
		c.Next()
	})
	enrH.RegisterRoutes(api)
	examH.RegisterRoutes(api)

	// --- PHASE 1: SIDI Import & Class Formation ---
	sidiCSV := `ID_DOMANDA;NOME;COGNOME;CODICE_FISCALE;DATA_NASCITA;SESSO;VOTO_LICENZA;INDIRIZZO;SECONDA_LINGUA;L104;DSA;RELIGIONE;COMPAGNI_RICHIESTI;GENITORE1_NOME;GENITORE1_COGNOME;GENITORE1_EMAIL;GENITORE1_TEL
SIDI-101;Matteo;Bianchi;BNCMTT12E15F205A;2012-05-15;M;9;Musicale;Francese;false;false;irc;;Andrea;Bianchi;andrea@example.com;3331112222
SIDI-102;Giulia;Verdi;VRDGLI12F45H501Z;2012-06-18;F;10;Musicale;Francese;false;false;irc;;Luca;Verdi;luca@example.com;3332223333
SIDI-103;Davide;Neri;NRIDVD12G01H501R;2012-07-01;M;8;Musicale;Francese;true;false;irc;;Sara;Neri;sara@example.com;3333334444
`
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)
	part, _ := writer.CreateFormFile("file", "iscrizioni.csv")
	_, _ = part.Write([]byte(sidiCSV))
	_ = writer.Close()

	reqImport, _ := http.NewRequest(http.MethodPost, "/api/v1/enrollment/import-sidi?academic_year=2026/2027", bodyBuf)
	reqImport.Header.Set("Content-Type", writer.FormDataContentType())
	wImport := httptest.NewRecorder()
	router.ServeHTTP(wImport, reqImport)
	assert.Equal(t, http.StatusOK, wImport.Code)

	var importResp map[string]interface{}
	_ = json.Unmarshal(wImport.Body.Bytes(), &importResp)
	assert.Equal(t, float64(3), importResp["imported"])

	// Generate Draft
	genReq := struct {
		Title        string                     `json:"title"`
		AcademicYear string                     `json:"academic_year"`
		Parameters   enrollment.FormationParams `json:"parameters"`
	}{
		Title:        "Formazione Classi Prime Musicali",
		AcademicYear: "2026/2027",
		Parameters: enrollment.FormationParams{
			TargetClassCount: 2,
			MaxL104PerClass:  1,
			BalanceGender:    true,
			BalanceGrades:    true,
		},
	}
	genBytes, _ := json.Marshal(genReq)
	reqDraft, _ := http.NewRequest(http.MethodPost, "/api/v1/enrollment/formation-drafts/generate", bytes.NewReader(genBytes))
	reqDraft.Header.Set("Content-Type", "application/json")
	wDraft := httptest.NewRecorder()
	router.ServeHTTP(wDraft, reqDraft)
	assert.Equal(t, http.StatusCreated, wDraft.Code)

	var draft enrollment.ClassFormationDraft
	_ = json.Unmarshal(wDraft.Body.Bytes(), &draft)
	assert.Len(t, draft.Assignments.Classes, 2)

	// Finalize draft
	reqFin, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/enrollment/formation-drafts/%s/finalize", draft.ID), nil)
	wFin := httptest.NewRecorder()
	router.ServeHTTP(wFin, reqFin)
	assert.Equal(t, http.StatusOK, wFin.Code)

	// --- PHASE 2: Middle School Exam Commission & Deliberation ---
	reqExam, _ := http.NewRequest(http.MethodPost, "/api/v1/middle-school-exam/classes/class-3m", bytes.NewBuffer([]byte(`{"president_name":"Dott.ssa Clara Bellini"}`)))
	reqExam.Header.Set("Content-Type", "application/json")
	wExam := httptest.NewRecorder()
	router.ServeHTTP(wExam, reqExam)
	assert.Equal(t, http.StatusOK, wExam.Code)

	var exam middleschoolexam.MiddleSchoolExam
	_ = json.Unmarshal(wExam.Body.Bytes(), &exam)
	assert.Equal(t, "exam-sec-1", exam.ID)

	// Save Admission for candidate Matteo Bianchi (Grade: 9)
	admissionData := map[string]interface{}{
		"grade":    9,
		"judgment": "Percorso triennale solido e maturo",
		"admitted": true,
	}
	admBytes, _ := json.Marshal(admissionData)
	reqAdm, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/middle-school-exam/%s/candidates/BNCMTT12E15F205A/admission", exam.ID), bytes.NewReader(admBytes))
	reqAdm.Header.Set("Content-Type", "application/json")
	wAdm := httptest.NewRecorder()
	router.ServeHTTP(wAdm, reqAdm)
	assert.Equal(t, http.StatusOK, wAdm.Code)

	// Evaluate Matteo with all 10s and proposed honors
	evalPayload := middleschoolexam.ExamCandidateGrades{
		GradeItalian:    10.0,
		GradeMath:       10.0,
		GradeEnglish:    10.0,
		GradeSecondLang: 10.0,
		GradeInterview:  10.0,
		ProposedHonors:  true,
	}
	evalReqBody := map[string]interface{}{
		"grades": evalPayload,
		"notes":  "Prova d'esame eccellente, approvata lode all'unanimità",
	}
	evalBytes, _ := json.Marshal(evalReqBody)
	reqEval, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/middle-school-exam/%s/candidates/BNCMTT12E15F205A/evaluate", exam.ID), bytes.NewReader(evalBytes))
	reqEval.Header.Set("Content-Type", "application/json")
	wEval := httptest.NewRecorder()
	router.ServeHTTP(wEval, reqEval)
	assert.Equal(t, http.StatusOK, wEval.Code)

	var evalResult middleschoolexam.MiddleSchoolExamCandidate
	_ = json.Unmarshal(wEval.Body.Bytes(), &evalResult)
	// (9 admission + 10 exam)/2 = 9.5 -> rounds to 10
	assert.Equal(t, 10, evalResult.FinalGrade)
	assert.True(t, evalResult.HasHonors)
	assert.Equal(t, "licenziato", evalResult.Outcome)

	// Generate Official Diploma
	reqDip, _ := http.NewRequest(http.MethodGet, "/api/v1/middle-school-exam/candidates/BNCMTT12E15F205A/diploma", nil)
	wDip := httptest.NewRecorder()
	router.ServeHTTP(wDip, reqDip)
	assert.Equal(t, http.StatusOK, wDip.Code)

	diplomaBody := wDip.Body.String()
	assert.Contains(t, diplomaBody, "REPUBBLICA ITALIANA")
	assert.Contains(t, diplomaBody, "MATTEO BIANCHI")
	assert.Contains(t, diplomaBody, "10/10 E LODE")
	assert.Contains(t, diplomaBody, "Dott.ssa Clara Bellini")
}
