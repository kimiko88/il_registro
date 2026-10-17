package enrollment

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type mockEnrollmentRepo struct {
	apps   map[string]EnrollmentApplication
	drafts map[string]*ClassFormationDraft
}

func newMockEnrollmentRepo() *mockEnrollmentRepo {
	return &mockEnrollmentRepo{
		apps:   make(map[string]EnrollmentApplication),
		drafts: make(map[string]*ClassFormationDraft),
	}
}

func (m *mockEnrollmentRepo) SaveApplications(ctx context.Context, apps []EnrollmentApplication, schoolID, academicYear string) (int, error) {
	for _, a := range apps {
		a.SchoolID = schoolID
		a.AcademicYear = academicYear
		m.apps[a.SidiApplicationID] = a
	}
	return len(apps), nil
}

func (m *mockEnrollmentRepo) ListApplications(ctx context.Context, schoolID, academicYear, status string) ([]EnrollmentApplication, error) {
	var list []EnrollmentApplication
	for _, a := range m.apps {
		if a.SchoolID == schoolID && a.AcademicYear == academicYear {
			list = append(list, a)
		}
	}
	return list, nil
}

func (m *mockEnrollmentRepo) SaveDraft(ctx context.Context, draft *ClassFormationDraft) error {
	m.drafts[draft.ID] = draft
	return nil
}

func (m *mockEnrollmentRepo) GetDraft(ctx context.Context, id string) (*ClassFormationDraft, error) {
	if d, ok := m.drafts[id]; ok {
		return d, nil
	}
	return nil, nil
}

func (m *mockEnrollmentRepo) ListDrafts(ctx context.Context, schoolID string) ([]ClassFormationDraft, error) {
	var list []ClassFormationDraft
	for _, d := range m.drafts {
		if d.SchoolID == schoolID {
			list = append(list, *d)
		}
	}
	return list, nil
}

func (m *mockEnrollmentRepo) UpdateDraftAssignments(ctx context.Context, id string, assignments ClassFormationDraftResult) error {
	if d, ok := m.drafts[id]; ok {
		d.Assignments = assignments
	}
	return nil
}

func (m *mockEnrollmentRepo) FinalizeDraft(ctx context.Context, id string) error {
	if d, ok := m.drafts[id]; ok {
		d.IsFinalized = true
	}
	return nil
}

func setupEnrollmentTestRouter(svc *Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/api/v1")
	rg.Use(func(c *gin.Context) {
		c.Set("user_id", "user-sec-1")
		c.Set("role", "secretary")
		c.Set("school_id", "school-1")
		c.Next()
	})
	h := NewHandler(svc)
	h.RegisterRoutes(rg)
	return r
}

func TestEnrollment_Integration_Flow(t *testing.T) {
	repo := newMockEnrollmentRepo()
	svc := NewService(repo)
	router := setupEnrollmentTestRouter(svc)

	// 1. Import SIDI CSV file
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)
	part, _ := writer.CreateFormFile("file", "domande_sidi.csv")
	_, _ = part.Write([]byte(`ID_DOMANDA;NOME;COGNOME;CODICE_FISCALE;DATA_NASCITA;SESSO;VOTO_LICENZA;INDIRIZZO;SECONDA_LINGUA;L104;DSA;RELIGIONE;COMPAGNI_RICHIESTI;GENITORE1_NOME;GENITORE1_COGNOME;GENITORE1_EMAIL;GENITORE1_TEL
SIDI-01;Lorenzo;Galli;GLLLNZ10A01H501K;2012-04-10;M;9;Scientifico;Spagnolo;false;false;irc;;Paolo;Galli;paolo@example.com;3330001111
SIDI-02;Chiara;Fontana;FNTCHR10A41H501W;2012-06-25;F;8;Scientifico;Spagnolo;false;false;irc;;Marco;Fontana;marco@example.com;3330002222
`))
	_ = writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/enrollment/import-sidi?academic_year=2026/2027", bodyBuf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on SIDI import, got %d: %s", w.Code, w.Body.String())
	}

	// 2. Generate class formation draft
	genPayload := struct {
		Title        string          `json:"title"`
		AcademicYear string          `json:"academic_year"`
		Parameters   FormationParams `json:"parameters"`
	}{
		Title:        "Formazione Sezioni Prime 2026/2027",
		AcademicYear: "2026/2027",
		Parameters: FormationParams{
			TargetClassCount: 2,
			MaxL104PerClass:  1,
			BalanceGender:    true,
			BalanceGrades:    true,
		},
	}
	genBytes, _ := json.Marshal(genPayload)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/enrollment/formation-drafts/generate", bytes.NewReader(genBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 on generate draft, got %d: %s", w.Code, w.Body.String())
	}
	var draft ClassFormationDraft
	_ = json.Unmarshal(w.Body.Bytes(), &draft)
	if len(draft.Assignments.Classes) != 2 {
		t.Errorf("expected 2 classes formed, got %d", len(draft.Assignments.Classes))
	}

	// 3. Finalize draft
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/enrollment/formation-drafts/"+draft.ID+"/finalize", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on finalize draft, got %d: %s", w.Code, w.Body.String())
	}

	_ = time.Now()
}
