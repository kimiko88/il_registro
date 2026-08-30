package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/grades"
	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockGradeWeightsRepo struct {
	grades.Repository
	configs []grades.GradeWeightConfig
	tests   map[string]*grades.ClassTest
	grades  map[string]*grades.Grade
}

func newMockGradeWeightsRepo() *mockGradeWeightsRepo {
	return &mockGradeWeightsRepo{
		configs: make([]grades.GradeWeightConfig, 0),
		tests:   make(map[string]*grades.ClassTest),
		grades:  make(map[string]*grades.Grade),
	}
}

func (m *mockGradeWeightsRepo) UpsertWeightConfig(ctx context.Context, cfg *grades.GradeWeightConfig) (*grades.GradeWeightConfig, error) {
	if cfg.ID == "" {
		cfg.ID = "cfg-1"
	}
	m.configs = append(m.configs, *cfg)
	return cfg, nil
}

func (m *mockGradeWeightsRepo) GetWeightConfigs(ctx context.Context, schoolID, subjectID, classID string) ([]grades.GradeWeightConfig, error) {
	return m.configs, nil
}

func (m *mockGradeWeightsRepo) ListWeightConfigs(schoolID, subjectID string) ([]grades.GradeWeightConfig, error) {
	return m.configs, nil
}

func (m *mockGradeWeightsRepo) CreateTest(ctx context.Context, test *grades.ClassTest) error {
	if test.ID == "" {
		test.ID = "test-matrix-1"
	}
	test.CreatedAt = time.Now()
	m.tests[test.ID] = test
	return nil
}

func (m *mockGradeWeightsRepo) BatchCreate(ctx context.Context, gList []*grades.Grade) error {
	for _, g := range gList {
		if g.ID == "" {
			g.ID = "grade-" + g.StudentID
		}
		m.grades[g.ID] = g
	}
	return nil
}

func (m *mockGradeWeightsRepo) CreateTestWithGrades(ctx context.Context, test *grades.ClassTest, gList []*grades.Grade) error {
	if err := m.CreateTest(ctx, test); err != nil {
		return err
	}
	return m.BatchCreate(ctx, gList)
}

func (m *mockGradeWeightsRepo) FindByStudent(ctx context.Context, studentID string) ([]grades.Grade, error) {
	var res []grades.Grade
	for _, g := range m.grades {
		if g.StudentID == studentID {
			res = append(res, *g)
		}
	}
	return res, nil
}

func (m *mockGradeWeightsRepo) FindBySubject(ctx context.Context, subjectID string, limit int) ([]grades.Grade, error) {
	var res []grades.Grade
	for _, g := range m.grades {
		if g.SubjectID == subjectID {
			res = append(res, *g)
		}
	}
	return res, nil
}


type mockUserRepoForWeights struct {
	users.Repository
}

func (m *mockUserRepoForWeights) GetByID(ctx context.Context, id string) (*users.User, error) {
	schoolID := "school-1"
	return &users.User{
		ID:       id,
		SchoolID: &schoolID,
		Role:     "teacher",
	}, nil
}

func (m *mockUserRepoForWeights) IsParentOf(parentID, studentID string) (bool, error) {
	return true, nil
}

func (m *mockUserRepoForWeights) IsTeacherOf(teacherID, classID, subjectID string) (bool, error) {
	return true, nil
}

func setupGradeWeightsRouter(repo grades.Repository, uRepo users.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := grades.NewService(repo, uRepo, nil, nil)
	handler := grades.NewHandler(svc, nil)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = "teacher"
		}
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "teacher-1"
		}
		c.Set("role", role)
		c.Set("user_id", userID)
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_GradeWeights_And_MatrixSubmission_Workflow(t *testing.T) {
	repo := newMockGradeWeightsRepo()
	uRepo := &mockUserRepoForWeights{}
	r := setupGradeWeightsRouter(repo, uRepo)

	// 1. Configure Grade Weights for Italian Subject as Admin/Secretary (Written weight: 1.5)
	subID := "sub-ita"
	weightPayload := grades.UpsertWeightConfigRequest{
		SubjectID:     &subID,
		GradeCategory: "summative",
		Weight:        1.5,
	}
	bodyWeight, _ := json.Marshal(weightPayload)
	reqWeight, _ := http.NewRequest("PUT", "/api/v1/grades/weight-configs", bytes.NewReader(bodyWeight))
	reqWeight.Header.Set("Content-Type", "application/json")
	reqWeight.Header.Set("X-Role", "admin")
	reqWeight.Header.Set("X-User-ID", "admin-1")
	wWeight := httptest.NewRecorder()
	r.ServeHTTP(wWeight, reqWeight)
	require.Equal(t, http.StatusOK, wWeight.Code)

	// 2. Fetch configured weights
	reqGetWeights, _ := http.NewRequest("GET", "/api/v1/grades/weight-configs?subject_id=sub-ita", nil)
	reqGetWeights.Header.Set("X-Role", "teacher")
	wGetWeights := httptest.NewRecorder()
	r.ServeHTTP(wGetWeights, reqGetWeights)
	require.Equal(t, http.StatusOK, wGetWeights.Code)

	var weights []grades.GradeWeightConfig
	err := json.Unmarshal(wGetWeights.Body.Bytes(), &weights)
	require.NoError(t, err)
	assert.Len(t, weights, 1)
	assert.Equal(t, 1.5, weights[0].Weight)

	// 3. Matrix Grade Entry for Class (Create Test with Multi-Student Grades as Teacher)
	testDate := time.Now().Format("2006-01-02")
	val1 := 8.5
	val2 := 5.5
	matrixPayload := grades.CreateClassTestRequest{
		ClassID:        "550e8400-e29b-41d4-a716-446655440000",
		SubjectID:      "sub-ita",
		Title:          "Verifica Scritta Divina Commedia - Inferno Canto V",
		Date:           testDate,
		EvaluationType: "Written",
		Grades: []grades.StudentGradeInput{
			{
				StudentID:  "550e8400-e29b-41d4-a716-446655440001",
				GradeValue: &val1,
				Notes:      "Ottima analisi del testo e contestualizzazione",
			},
			{
				StudentID:  "550e8400-e29b-41d4-a716-446655440002",
				GradeValue: &val2,
				Notes:      "Comprensione parziale, da consolidare",
			},
		},
	}
	bodyMatrix, _ := json.Marshal(matrixPayload)
	reqMatrix, _ := http.NewRequest("POST", "/api/v1/grades/tests", bytes.NewReader(bodyMatrix))
	reqMatrix.Header.Set("Content-Type", "application/json")
	reqMatrix.Header.Set("X-Role", "teacher")
	reqMatrix.Header.Set("X-User-ID", "teacher-1")
	wMatrix := httptest.NewRecorder()
	r.ServeHTTP(wMatrix, reqMatrix)
	require.Equal(t, http.StatusCreated, wMatrix.Code)

	var createdTest struct {
		Message string `json:"message"`
		TestID  string `json:"test_id"`
	}
	err = json.Unmarshal(wMatrix.Body.Bytes(), &createdTest)
	require.NoError(t, err)
	assert.NotEmpty(t, createdTest.TestID)
	assert.Len(t, repo.grades, 2)
}
