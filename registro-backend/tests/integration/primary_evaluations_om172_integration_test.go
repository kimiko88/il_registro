package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/primaryeval"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockPrimaryIntegrationRepo struct {
	objectives  map[string]*primaryeval.LearningObjective
	evaluations []primaryeval.PrimaryEvaluation
}

func newMockPrimaryIntegrationRepo() *mockPrimaryIntegrationRepo {
	return &mockPrimaryIntegrationRepo{
		objectives:  make(map[string]*primaryeval.LearningObjective),
		evaluations: []primaryeval.PrimaryEvaluation{},
	}
}

func (m *mockPrimaryIntegrationRepo) CreateObjective(ctx context.Context, obj *primaryeval.LearningObjective) (*primaryeval.LearningObjective, error) {
	if obj.ID == "" {
		obj.ID = "obj-" + obj.Title
	}
	m.objectives[obj.ID] = obj
	return obj, nil
}

func (m *mockPrimaryIntegrationRepo) ListObjectives(ctx context.Context, schoolID, subjectID string, classID *string, yearGrade int, academicYear string) ([]primaryeval.LearningObjective, error) {
	var res []primaryeval.LearningObjective
	for _, o := range m.objectives {
		if o.SchoolID == schoolID {
			if subjectID != "" && o.SubjectID != subjectID {
				continue
			}
			if yearGrade > 0 && o.YearGrade != yearGrade {
				continue
			}
			res = append(res, *o)
		}
	}
	return res, nil
}

func (m *mockPrimaryIntegrationRepo) GetObjective(ctx context.Context, id string) (*primaryeval.LearningObjective, error) {
	return m.objectives[id], nil
}

func (m *mockPrimaryIntegrationRepo) DeleteObjective(ctx context.Context, id string) error {
	delete(m.objectives, id)
	return nil
}

func (m *mockPrimaryIntegrationRepo) SaveEvaluationsBatch(ctx context.Context, schoolID, teacherID string, req *primaryeval.SaveEvaluationsBatchRequest) error {
	for _, item := range req.Evaluations {
		eval := primaryeval.PrimaryEvaluation{
			ID:                   "eval-" + item.StudentID + "-" + req.ObjectiveID,
			SchoolID:             schoolID,
			StudentID:            item.StudentID,
			StudentName:          "Studente " + item.StudentID,
			ClassID:              req.ClassID,
			SubjectID:            req.SubjectID,
			TeacherID:            teacherID,
			ObjectiveID:          req.ObjectiveID,
			Level:                item.Level,
			DimensionAutonomy:    item.DimensionAutonomy,
			DimensionContinuity:  item.DimensionContinuity,
			DimensionFamiliarity: item.DimensionFamiliarity,
			DimensionResources:   item.DimensionResources,
			Semester:             req.Semester,
			Notes:                item.Notes,
			Date:                 time.Now(),
		}
		m.evaluations = append(m.evaluations, eval)
	}
	return nil
}

func (m *mockPrimaryIntegrationRepo) ListEvaluationsByClassAndSubject(ctx context.Context, classID, subjectID string, semester int) ([]primaryeval.PrimaryEvaluation, error) {
	var res []primaryeval.PrimaryEvaluation
	for _, e := range m.evaluations {
		if e.ClassID == classID && e.SubjectID == subjectID {
			if semester == 0 || e.Semester == semester {
				res = append(res, e)
			}
		}
	}
	return res, nil
}

func (m *mockPrimaryIntegrationRepo) GetStudentEvaluations(ctx context.Context, studentID string, semester int) ([]primaryeval.PrimaryEvaluation, error) {
	var res []primaryeval.PrimaryEvaluation
	for _, e := range m.evaluations {
		if e.StudentID == studentID {
			if semester == 0 || e.Semester == semester {
				res = append(res, e)
			}
		}
	}
	return res, nil
}

// TestPrimaryEvaluationsOM172_IntegrationFlow tests the complete ministerial flow of Primary School descriptive evaluations (O.M. 172/2020)
func TestPrimaryEvaluationsOM172_IntegrationFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockPrimaryIntegrationRepo()
	service := primaryeval.NewService(repo)
	handler := primaryeval.NewHandler(service)

	router := gin.New()
	api := router.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("school_id", "school-primaria-1")
		c.Set("user_id", "teacher-maestra-anna")
		c.Set("role", "teacher")
		c.Next()
	})
	handler.RegisterRoutes(api)

	// Step 1: Create 2 Ministerial Learning Objectives (Italiano, Year Grade 3)
	objPayload1 := primaryeval.CreateObjectiveRequest{
		SubjectID:    "sub-italiano",
		YearGrade:    3,
		Title:        "Ascolto e comprensione orale",
		Description:  "Comprende il tema principale e le informazioni essenziali di un testo ascoltato",
		AcademicYear: "2025/2026",
	}
	body1, _ := json.Marshal(objPayload1)
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/primary/objectives", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	objPayload2 := primaryeval.CreateObjectiveRequest{
		SubjectID:    "sub-italiano",
		YearGrade:    3,
		Title:        "Lettura e comprensione del testo",
		Description:  "Legge testi di varia natura cogliendo il senso globale e dettagli",
		AcademicYear: "2025/2026",
	}
	body2, _ := json.Marshal(objPayload2)
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/primary/objectives", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)

	// Step 2: List objectives and verify they were stored
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/primary/objectives?subject_id=sub-italiano&year_grade=3", nil)
	wList := httptest.NewRecorder()
	router.ServeHTTP(wList, reqList)
	assert.Equal(t, http.StatusOK, wList.Code)

	var listResp struct {
		Data []primaryeval.LearningObjective `json:"data"`
	}
	_ = json.Unmarshal(wList.Body.Bytes(), &listResp)
	assert.Len(t, listResp.Data, 2)
	objID1 := listResp.Data[0].ID
	objID2 := listResp.Data[1].ID

	// Step 3: Register batch evaluations for Class 3A (3 Students) with 4 Ministerial Dimensions (O.M. 172/2020)
	batchPayload := primaryeval.SaveEvaluationsBatchRequest{
		ClassID:     "class-3a",
		SubjectID:   "sub-italiano",
		ObjectiveID: objID1,
		Date:        "2026-01-20",
		Semester:    1,
		Evaluations: []primaryeval.StudentObjectiveLevelItem{
			{
				StudentID:            "stu-lorenzo",
				Level:                primaryeval.LevelAvanzato,
				DimensionAutonomy:    "completamente_autonomo",
				DimensionContinuity:  "continuo",
				DimensionFamiliarity: "situazioni_note_e_non_note",
				DimensionResources:   "risorse_proprie_e_fornite",
				Notes:                "Ottima partecipazione attiva",
			},
			{
				StudentID:            "stu-matteo",
				Level:                primaryeval.LevelIntermedio,
				DimensionAutonomy:    "autonomo",
				DimensionContinuity:  "discontinuo",
				DimensionFamiliarity: "situazioni_note",
				DimensionResources:   "risorse_fornite",
			},
			{
				StudentID:            "stu-sofia",
				Level:                primaryeval.LevelInViaPrimaAcquisiz,
				DimensionAutonomy:    "solo_se_guidato",
				DimensionContinuity:  "raro",
				DimensionFamiliarity: "situazioni_note",
				DimensionResources:   "supporto_docente",
				Notes:                "In fase di consolidamento",
			},
		},
	}
	batchBytes, _ := json.Marshal(batchPayload)
	reqBatch, _ := http.NewRequest(http.MethodPost, "/api/v1/primary/evaluations/batch", bytes.NewBuffer(batchBytes))
	reqBatch.Header.Set("Content-Type", "application/json")
	wBatch := httptest.NewRecorder()
	router.ServeHTTP(wBatch, reqBatch)
	assert.Equal(t, http.StatusOK, wBatch.Code)

	// Step 4: Register second objective for Student Lorenzo
	batchPayload2 := primaryeval.SaveEvaluationsBatchRequest{
		ClassID:     "class-3a",
		SubjectID:   "sub-italiano",
		ObjectiveID: objID2,
		Date:        "2026-01-22",
		Semester:    1,
		Evaluations: []primaryeval.StudentObjectiveLevelItem{
			{
				StudentID: "stu-lorenzo",
				Level:     primaryeval.LevelIntermedio,
			},
		},
	}
	batchBytes2, _ := json.Marshal(batchPayload2)
	reqBatch2, _ := http.NewRequest(http.MethodPost, "/api/v1/primary/evaluations/batch", bytes.NewBuffer(batchBytes2))
	reqBatch2.Header.Set("Content-Type", "application/json")
	wBatch2 := httptest.NewRecorder()
	router.ServeHTTP(wBatch2, reqBatch2)
	assert.Equal(t, http.StatusOK, wBatch2.Code)

	// Step 5: Query individual student evaluations for Lorenzo
	reqStudent, _ := http.NewRequest(http.MethodGet, "/api/v1/primary/student/stu-lorenzo?semester=1", nil)
	wStudent := httptest.NewRecorder()
	router.ServeHTTP(wStudent, reqStudent)
	assert.Equal(t, http.StatusOK, wStudent.Code)

	var stuResp struct {
		Data []primaryeval.PrimaryEvaluation `json:"data"`
	}
	_ = json.Unmarshal(wStudent.Body.Bytes(), &stuResp)
	assert.Len(t, stuResp.Data, 2)
	assert.Equal(t, primaryeval.LevelAvanzato, stuResp.Data[0].Level)
	assert.Equal(t, primaryeval.LevelIntermedio, stuResp.Data[1].Level)

	// Step 6: Query class matrix response
	reqMatrix, _ := http.NewRequest(http.MethodGet, "/api/v1/primary/matrix?class_id=class-3a&subject_id=sub-italiano&semester=1", nil)
	wMatrix := httptest.NewRecorder()
	router.ServeHTTP(wMatrix, reqMatrix)
	assert.Equal(t, http.StatusOK, wMatrix.Code)

	var matrixResp struct {
		Data primaryeval.PrimaryMatrixResponse `json:"data"`
	}
	_ = json.Unmarshal(wMatrix.Body.Bytes(), &matrixResp)
	assert.Equal(t, "class-3a", matrixResp.Data.ClassID)
	assert.Equal(t, 1, matrixResp.Data.Semester)
	assert.Len(t, matrixResp.Data.Students, 3)

	// Verify that Lorenzo's matrix row has both evaluations mapped
	var lorenzoRow *primaryeval.PrimaryStudentMatrixRow
	for _, s := range matrixResp.Data.Students {
		if s.StudentID == "stu-lorenzo" {
			r := s
			lorenzoRow = &r
			break
		}
	}
	assert.NotNil(t, lorenzoRow)
	assert.Equal(t, primaryeval.LevelAvanzato, lorenzoRow.Evaluations[objID1].Level)
	assert.Equal(t, primaryeval.LevelIntermedio, lorenzoRow.Evaluations[objID2].Level)
}
