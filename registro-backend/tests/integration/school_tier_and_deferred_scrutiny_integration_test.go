package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/schools"
	"registro-backend/internal/scrutiny"
	"registro-backend/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockScrutinyDeferredRepo implements scrutiny.Repository for deferred scrutiny testing
type mockScrutinyDeferredRepo struct {
	mock.Mock
}

var _ scrutiny.Repository = (*mockScrutinyDeferredRepo)(nil)

func (m *mockScrutinyDeferredRepo) SaveRecord(ctx context.Context, record *scrutiny.ScrutinyRecord) error {
	return nil
}

func (m *mockScrutinyDeferredRepo) GetRecord(ctx context.Context, studentID, classID string, semester int) (*scrutiny.ScrutinyRecord, error) {
	return nil, nil
}

func (m *mockScrutinyDeferredRepo) ListRecordsByClass(ctx context.Context, classID string, semester int) ([]scrutiny.ScrutinyRecord, error) {
	return nil, nil
}

func (m *mockScrutinyDeferredRepo) ValidateClassScrutiny(ctx context.Context, classID string, semester int, validatorID string) error {
	return nil
}

func (m *mockScrutinyDeferredRepo) UpdateClassScrutinyStatus(ctx context.Context, classID string, semester int, status string) error {
	return nil
}

func (m *mockScrutinyDeferredRepo) SaveDeficiency(ctx context.Context, def *scrutiny.StudentDeficiency) error {
	args := m.Called(ctx, def)
	return args.Error(0)
}

func (m *mockScrutinyDeferredRepo) GetDeficienciesByStudent(ctx context.Context, studentID string) ([]scrutiny.StudentDeficiency, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]scrutiny.StudentDeficiency), args.Error(1)
}

func (m *mockScrutinyDeferredRepo) GetDeficienciesByClass(ctx context.Context, classID string, semester int) ([]scrutiny.StudentDeficiency, error) {
	args := m.Called(ctx, classID, semester)
	return args.Get(0).([]scrutiny.StudentDeficiency), args.Error(1)
}

func (m *mockScrutinyDeferredRepo) SaveDeferredScrutiny(ctx context.Context, req *scrutiny.SaveDeferredScrutinyRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestSchoolTierAndDeferredScrutinyIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup Mocks
	mockSchoolRepo := new(testhelpers.MockSchoolsRepository)
	mockScrutinyRepo := new(mockScrutinyDeferredRepo)

	// Setup Services
	schoolSvc := schools.NewService(mockSchoolRepo)
	scrutinySvc := scrutiny.NewService(mockScrutinyRepo, nil, nil, nil, nil)

	// Setup Handlers
	schoolHandler := schools.NewHandler(schoolSvc)
	scrutinyHandler := scrutiny.NewHandler(scrutinySvc, nil)

	// Setup Router
	router := gin.New()
	api := router.Group("/api/v1")

	// Schools routes with auth context
	schoolsGroup := api.Group("/schools")
	schoolsGroup.Use(func(c *gin.Context) {
		c.Set("user_id", "admin-user-1")
		c.Set("role", "superadmin")
		c.Next()
	})
	{
		schoolsGroup.GET("/tiers", schoolHandler.ListTiers)
		schoolsGroup.POST("", schoolHandler.Create)
		schoolsGroup.GET("/:id", schoolHandler.Get)
		schoolsGroup.GET("/:id/tier-features", schoolHandler.GetTierFeatures)
	}

	// Scrutiny routes with auth context
	scrutinyGroup := api.Group("/scrutiny")
	scrutinyGroup.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-coord-1")
		c.Set("role", "teacher")
		c.Next()
	})
	{
		scrutinyGroup.POST("/deferred", scrutinyHandler.SaveDeferredScrutiny)
	}

	// ─── Test 1: GET /api/v1/schools/tiers ────────────────────────────────────
	t.Run("GET /schools/tiers returns canonical Italian school tiers and normative metadata", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/schools/tiers", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var res struct {
			Items []schools.TierFeatures `json:"items"`
			Total int                    `json:"total"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.Equal(t, 6, res.Total)
		assert.Len(t, res.Items, 6)

		// Verify specific tier features
		var infanzia, primaria, sec2 *schools.TierFeatures
		for i := range res.Items {
			switch res.Items[i].Tier {
			case schools.TierInfanzia:
				infanzia = &res.Items[i]
			case schools.TierPrimaria:
				primaria = &res.Items[i]
			case schools.TierSecondariaSecondoGrado:
				sec2 = &res.Items[i]
			}
		}

		require.NotNil(t, infanzia)
		assert.False(t, infanzia.HasGrades)
		assert.True(t, infanzia.HasCampiEsperienza)
		assert.False(t, infanzia.HasDeferredScrutiny)

		require.NotNil(t, primaria)
		assert.False(t, primaria.HasGrades)
		assert.True(t, primaria.HasPrimaryLevels)
		assert.False(t, primaria.HasDeferredScrutiny)
		assert.Equal(t, "obiettivi_descrittivi", primaria.EvaluationType)

		require.NotNil(t, sec2)
		assert.True(t, sec2.HasGrades)
		assert.True(t, sec2.HasDeferredScrutiny)
		assert.True(t, sec2.HasSchoolCredits)
		assert.True(t, sec2.HasPCTO)
		assert.Equal(t, "numerica_decimale", sec2.EvaluationType)
	})

	// ─── Test 2: POST /api/v1/schools with school_level (Primaria) ───────────
	t.Run("POST /schools normalizes and persists school_level and type", func(t *testing.T) {
		mockSchoolRepo.On("Create", mock.Anything, mock.MatchedBy(func(s *schools.School) bool {
			return s.Name == "IC Leonardo da Vinci - Plesso Primaria" &&
				s.Code == "RMIC801001" &&
				s.SchoolLevel == schools.TierPrimaria &&
				s.Type == schools.TierPrimaria
		})).Return(nil).Once()

		payload := schools.CreateSchoolRequest{
			Name:        "IC Leonardo da Vinci - Plesso Primaria",
			Code:        "RMIC801001",
			SchoolLevel: "scuola_primaria", // raw string to test normalization
			Address:     "Via Roma 1",
			City:        "Roma",
			Phone:       "06123456",
			Email:       "rmic801001@istruzione.it",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/schools", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)

		var created schools.School
		err := json.Unmarshal(rec.Body.Bytes(), &created)
		require.NoError(t, err)
		assert.Equal(t, schools.TierPrimaria, created.SchoolLevel)
		assert.Equal(t, schools.TierPrimaria, created.Type)
	})

	// ─── Test 3: GET /api/v1/schools/:id/tier-features (Secondaria II Grado) ──
	t.Run("GET /schools/:id/tier-features returns deferred scrutiny support for high school", func(t *testing.T) {
		schoolID := "high-school-uuid-99"
		mockSchoolRepo.On("GetByID", mock.Anything, schoolID).Return(&schools.School{
			ID:          schoolID,
			Name:        "Liceo Scientifico A. Righi",
			Code:        "RMPS010004",
			SchoolLevel: schools.TierSecondariaSecondoGrado,
			Type:        schools.TierSecondariaSecondoGrado,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/schools/"+schoolID+"/tier-features", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var features schools.TierFeatures
		err := json.Unmarshal(rec.Body.Bytes(), &features)
		require.NoError(t, err)
		assert.Equal(t, schools.TierSecondariaSecondoGrado, features.Tier)
		assert.True(t, features.HasDeferredScrutiny)
		assert.True(t, features.HasSchoolCredits)
		assert.False(t, features.HasPrimaryLevels)
		assert.False(t, features.HasCampiEsperienza)
	})

	// ─── Test 4: GET /api/v1/schools/:id/tier-features (Infanzia) ─────────────
	t.Run("GET /schools/:id/tier-features returns campi d'esperienza for infanzia", func(t *testing.T) {
		schoolID := "infanzia-uuid-42"
		mockSchoolRepo.On("GetByID", mock.Anything, schoolID).Return(&schools.School{
			ID:          schoolID,
			Name:        "Scuola dell'Infanzia Il Girotondo",
			Code:        "RMAA02000X",
			SchoolLevel: schools.TierInfanzia,
			Type:        schools.TierInfanzia,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/schools/"+schoolID+"/tier-features", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var features schools.TierFeatures
		err := json.Unmarshal(rec.Body.Bytes(), &features)
		require.NoError(t, err)
		assert.Equal(t, schools.TierInfanzia, features.Tier)
		assert.True(t, features.HasCampiEsperienza)
		assert.False(t, features.HasGrades)
		assert.False(t, features.HasDeferredScrutiny)
	})

	// ─── Test 5: POST /api/v1/scrutiny/deferred ───────────────────────────────
	t.Run("POST /scrutiny/deferred resolves reservation for student with suspended judgment", func(t *testing.T) {
		mockScrutinyRepo.On("SaveDeferredScrutiny", mock.Anything, mock.MatchedBy(func(req *scrutiny.SaveDeferredScrutinyRequest) bool {
			return req.StudentID == "student-deferred-1" &&
				req.ClassID == "class-4b-scientifico" &&
				req.FinalDecision == "promosso_con_debiti_saldati"
		})).Return(nil).Once()

		reqBody := scrutiny.SaveDeferredScrutinyRequest{
			StudentID:     "student-deferred-1",
			ClassID:       "class-4b-scientifico",
			FinalDecision: "promosso_con_debiti_saldati",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/scrutiny/deferred", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "deferred scrutiny saved successfully", resp["message"])
	})

	// ─── Test 6: POST /api/v1/scrutiny/deferred validation error ──────────────
	t.Run("POST /scrutiny/deferred returns 400 when missing required fields", func(t *testing.T) {
		invalidBody := `{"student_id": "student-only"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/scrutiny/deferred", bytes.NewBufferString(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
