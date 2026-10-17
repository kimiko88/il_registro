package textbooks

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockAIEFullRepo struct {
	mockRepo
	catalog   []AIECatalogBook
	limits    map[string]*SpendingLimit
	adoptions map[string][]ClassAdoptionItem
}

func newMockAIEFullRepo() *mockAIEFullRepo {
	return &mockAIEFullRepo{
		mockRepo:  *newMockRepo(),
		catalog:   []AIECatalogBook{},
		limits:    make(map[string]*SpendingLimit),
		adoptions: make(map[string][]ClassAdoptionItem),
	}
}

func (m *mockAIEFullRepo) UpsertAIECatalog(ctx context.Context, books []AIECatalogBook) (int, error) {
	m.catalog = append(m.catalog, books...)
	return len(books), nil
}

func (m *mockAIEFullRepo) SearchAIECatalog(ctx context.Context, queryStr, subject, schoolOrder string, limit int) ([]AIECatalogBook, error) {
	return m.catalog, nil
}

func (m *mockAIEFullRepo) GetSpendingLimit(ctx context.Context, schoolID string, classYear int, schoolOrder, academicYear string) (*SpendingLimit, error) {
	key := schoolID + "-1"
	if l, ok := m.limits[key]; ok {
		return l, nil
	}
	return &SpendingLimit{
		SchoolID:            schoolID,
		ClassYear:           1,
		SchoolOrder:         "secondaria_2",
		MaxAmount:           300.0,
		AllowedTolerancePct: 10.0,
		AcademicYear:        "2026/2027",
	}, nil
}

func (m *mockAIEFullRepo) UpsertSpendingLimit(ctx context.Context, limit *SpendingLimit) error {
	m.limits[limit.SchoolID+"-1"] = limit
	return nil
}

func (m *mockAIEFullRepo) ListClassAdoptions(ctx context.Context, classID string) ([]ClassAdoptionItem, error) {
	return m.adoptions[classID], nil
}

func (m *mockAIEFullRepo) SaveClassAdoption(ctx context.Context, item *ClassAdoptionItem) error {
	item.ID = "adopt-1"
	m.adoptions[item.ClassID] = append(m.adoptions[item.ClassID], *item)
	return nil
}

func (m *mockAIEFullRepo) DeleteClassAdoption(ctx context.Context, id string) error {
	for k, v := range m.adoptions {
		var filtered []ClassAdoptionItem
		for _, it := range v {
			if it.ID != id {
				filtered = append(filtered, it)
			}
		}
		m.adoptions[k] = filtered
	}
	return nil
}

func (m *mockAIEFullRepo) GetClassInfo(ctx context.Context, classID string) (int, string, string, string, string, error) {
	return 1, "secondaria_2", "2026/2027", "RMPS010004", "1A", nil
}

func setupAIETestRouter(svc *Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/api/v1")
	rg.Use(func(c *gin.Context) {
		c.Set("user_id", "user-sec-1")
		c.Set("role", "secretary")
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler := NewHandler(svc)
	handler.RegisterRoutes(rg)
	return r
}

func TestAIE_Integration_FullWorkflow(t *testing.T) {
	repo := newMockAIEFullRepo()
	svc := NewService(repo)
	router := setupAIETestRouter(svc)

	// 1. Import AIE CSV file via POST /api/v1/textbooks/aie/import
	csvBody := &bytes.Buffer{}
	writer := multipart.NewWriter(csvBody)
	part, err := writer.CreateFormFile("file", "catalog.csv")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	_, _ = part.Write([]byte(`CODICE_ISBN;TITOLO;AUTORI;EDITORE;DISCIPLINA;PREZZO;VOLUME;ANNO_EDIZIONE;ORDINE_SCUOLA;DIGITALE
9788808836243;Matematica.blu 2.0;Bergamini;Zanichelli;Matematica;28.90;1;2023;secondaria_2;false
9788804724567;Divina Commedia;Dante;Mondadori;Italiano;22.00;1;2021;secondaria_2;false
`))
	_ = writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/textbooks/aie/import", csvBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on import, got %d: %s", w.Code, w.Body.String())
	}

	// 2. Search catalog via GET /api/v1/textbooks/aie/catalog
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/textbooks/aie/catalog?q=Matematica", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on catalog search, got %d", w.Code)
	}
	var books []AIECatalogBook
	_ = json.Unmarshal(w.Body.Bytes(), &books)
	if len(books) != 2 {
		t.Errorf("expected 2 books, got %d", len(books))
	}

	// 3. Adopt a book via POST /api/v1/textbooks/classes/class-1/adoptions
	adoptPayload := ClassAdoptionItem{
		SubjectID:    "subj-math",
		SubjectName:  "Matematica",
		BookID:       "book-1",
		BookTitle:    "Matematica.blu 2.0",
		Price:        28.90,
		AdoptionType: "nuova_adozione",
	}
	jsonBytes, _ := json.Marshal(adoptPayload)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/textbooks/classes/class-1/adoptions", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 on adoption, got %d: %s", w.Code, w.Body.String())
	}

	// 4. Check spending report via GET /api/v1/textbooks/classes/class-1/spending-report
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/textbooks/classes/class-1/spending-report", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on spending report, got %d", w.Code)
	}
	var rep SpendingReport
	_ = json.Unmarshal(w.Body.Bytes(), &rep)
	if rep.TotalSpending != 28.90 {
		t.Errorf("expected 28.90 total spending, got %.2f", rep.TotalSpending)
	}
	if rep.Status != StatusWithinLimit {
		t.Errorf("expected WITHIN_LIMIT, got %s", rep.Status)
	}

	// 5. Export AIE file via GET /api/v1/textbooks/classes/class-1/aie-export
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/textbooks/classes/class-1/aie-export", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on AIE export, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("RMPS010004")) {
		t.Errorf("expected export to contain school code RMPS010004")
	}
}
