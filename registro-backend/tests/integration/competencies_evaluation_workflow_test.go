package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/competencies"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCompetenciesService struct {
	evals map[string]*competencies.Evaluation
}

func newMockCompetenciesService() *mockCompetenciesService {
	return &mockCompetenciesService{
		evals: make(map[string]*competencies.Evaluation),
	}
}

func setupCompetenciesRouter(mock *mockCompetenciesService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

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

	comp := api.Group("/competencies")
	{
		comp.GET("/student/:studentID", func(c *gin.Context) {
			studentID := c.Param("studentID")
			var res []*competencies.Evaluation
			for _, e := range mock.evals {
				if e.StudentID == studentID {
					res = append(res, e)
				}
			}
			c.JSON(http.StatusOK, res)
		})
		comp.POST("/evaluations", func(c *gin.Context) {
			var req competencies.SaveEvaluationRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if !competencies.IsValidCompetencyLevel(req.Level) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level"})
				return
			}
			eval := &competencies.Evaluation{
				ID:             "eval-" + time.Now().Format("150405.000000"),
				SchoolID:       c.GetString("school_id"),
				StudentID:      req.StudentID,
				ClassID:        req.ClassID,
				EvaluatorID:    c.GetString("user_id"),
				Semester:       req.Semester,
				CompetenceCode: req.CompetenceCode,
				CompetenceName: req.CompetenceName,
				Level:          req.Level,
				Descriptor:     req.Descriptor,
				Notes:          req.Notes,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			mock.evals[eval.ID] = eval
			c.JSON(http.StatusOK, eval)
		})
		comp.POST("/batch", func(c *gin.Context) {
			var req competencies.BatchSaveRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			for _, st := range req.Evaluations {
				for code, lvl := range st.Evaluations {
					if competencies.IsValidCompetencyLevel(lvl) {
						id := "eval-" + st.StudentID + "-" + code
						mock.evals[id] = &competencies.Evaluation{
							ID:             id,
							SchoolID:       c.GetString("school_id"),
							StudentID:      st.StudentID,
							ClassID:        req.ClassID,
							CompetenceCode: code,
							CompetenceName: code,
							Level:          lvl,
							CreatedAt:      time.Now(),
							UpdatedAt:      time.Now(),
						}
					}
				}
			}
			c.JSON(http.StatusOK, gin.H{"message": "valutazioni per competenze salvate con successo"})
		})
	}
	return r
}

func TestIntegration_Competencies_Evaluation_Workflow(t *testing.T) {
	mock := newMockCompetenciesService()
	r := setupCompetenciesRouter(mock)

	// 1. Teacher assigns a single competency evaluation (A_Avanzato)
	singleReq := competencies.SaveEvaluationRequest{
		StudentID:      "student-1",
		ClassID:        "class-3B",
		Semester:       2,
		CompetenceCode: "EU_DIGITAL",
		CompetenceName: "Competenza Digitale e Pensiero Computazionale",
		Level:          "A_Avanzato",
		Descriptor:     "Utilizza con piena padronanza e spirito critico le tecnologie digitali",
		Notes:          "Progetto di coding svolto in autonomia",
	}
	body, _ := json.Marshal(singleReq)
	reqSingle, _ := http.NewRequest("POST", "/api/v1/competencies/evaluations", bytes.NewReader(body))
	reqSingle.Header.Set("Content-Type", "application/json")
	reqSingle.Header.Set("X-Role", "teacher")
	wSingle := httptest.NewRecorder()
	r.ServeHTTP(wSingle, reqSingle)
	require.Equal(t, http.StatusOK, wSingle.Code)

	var evalResp competencies.Evaluation
	err := json.Unmarshal(wSingle.Body.Bytes(), &evalResp)
	require.NoError(t, err)
	assert.Equal(t, "A_Avanzato", evalResp.Level)
	assert.Equal(t, "EU_DIGITAL", evalResp.CompetenceCode)

	// 2. Batch save competencies for class
	batchReq := competencies.BatchSaveRequest{
		ClassID: "class-3B",
		Period:  "finale",
		Evaluations: []competencies.StudentCompetencyEvaluation{
			{
				StudentID:   "student-2",
				StudentName: "Giulia Bianchi",
				Evaluations: map[string]string{
					"EU_LANG":  "B_Intermedio",
					"EU_STEM":  "A_Avanzato",
					"EU_CIVIC": "A_Avanzato",
				},
			},
		},
	}
	bodyBatch, _ := json.Marshal(batchReq)
	reqBatch, _ := http.NewRequest("POST", "/api/v1/competencies/batch", bytes.NewReader(bodyBatch))
	reqBatch.Header.Set("Content-Type", "application/json")
	reqBatch.Header.Set("X-Role", "teacher")
	wBatch := httptest.NewRecorder()
	r.ServeHTTP(wBatch, reqBatch)
	require.Equal(t, http.StatusOK, wBatch.Code)

	// 3. Query student's competency portfolio
	reqGet, _ := http.NewRequest("GET", "/api/v1/competencies/student/student-2", nil)
	reqGet.Header.Set("X-Role", "teacher")
	wGet := httptest.NewRecorder()
	r.ServeHTTP(wGet, reqGet)
	require.Equal(t, http.StatusOK, wGet.Code)

	var studentEvals []*competencies.Evaluation
	err = json.Unmarshal(wGet.Body.Bytes(), &studentEvals)
	require.NoError(t, err)
	assert.Len(t, studentEvals, 3)
}
