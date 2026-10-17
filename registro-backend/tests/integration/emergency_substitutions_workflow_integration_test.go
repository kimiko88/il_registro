package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/substitutions"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type memorySubstitutionsIntegrationRepo struct {
	subs             map[string]*substitutions.Substitution
	available        []substitutions.TeacherCandidate
	classTeachers    map[string]bool
	subjectTeachers  map[string]bool
	subCounts        map[string]int
	profileIDMapping map[string]string
}

func newMemorySubstitutionsIntegrationRepo() *memorySubstitutionsIntegrationRepo {
	return &memorySubstitutionsIntegrationRepo{
		subs:             make(map[string]*substitutions.Substitution),
		classTeachers:    make(map[string]bool),
		subjectTeachers:  make(map[string]bool),
		subCounts:        make(map[string]int),
		profileIDMapping: make(map[string]string),
	}
}

func (m *memorySubstitutionsIntegrationRepo) Create(ctx context.Context, sub *substitutions.Substitution) error {
	if sub.ID == "" {
		sub.ID = fmt.Sprintf("sub-%d", len(m.subs)+1)
	}
	sub.CreatedAt = time.Now()
	m.subs[sub.ID] = sub
	return nil
}

func (m *memorySubstitutionsIntegrationRepo) GetByID(ctx context.Context, id string) (*substitutions.Substitution, error) {
	sub, ok := m.subs[id]
	if !ok {
		return nil, errors.New("sostituzione non trovata")
	}
	return sub, nil
}

func (m *memorySubstitutionsIntegrationRepo) ListBySchool(ctx context.Context, schoolID, date string) ([]*substitutions.Substitution, error) {
	var list []*substitutions.Substitution
	for _, s := range m.subs {
		if schoolID != "" && s.SchoolID != "" && s.SchoolID != schoolID {
			continue
		}
		list = append(list, s)
	}
	return list, nil
}

func (m *memorySubstitutionsIntegrationRepo) ListByTeacher(ctx context.Context, teacherID string, date string) ([]*substitutions.Substitution, error) {
	var list []*substitutions.Substitution
	for _, s := range m.subs {
		if s.SubstituteTeacherID != nil && *s.SubstituteTeacherID == teacherID {
			list = append(list, s)
		}
	}
	return list, nil
}

func (m *memorySubstitutionsIntegrationRepo) AssignSubstitute(ctx context.Context, id string, substituteTeacherID string, notes string) error {
	sub, ok := m.subs[id]
	if !ok {
		return errors.New("sostituzione non trovata")
	}
	sub.SubstituteTeacherID = &substituteTeacherID
	sub.Status = substitutions.StatusAssigned
	sub.Notes = notes
	return nil
}

func (m *memorySubstitutionsIntegrationRepo) ConfirmSubstitution(ctx context.Context, id string, substituteTeacherID string) error {
	sub, ok := m.subs[id]
	if !ok {
		return errors.New("sostituzione non trovata")
	}
	sub.Status = substitutions.StatusConfirmed
	return nil
}

func (m *memorySubstitutionsIntegrationRepo) SignRegister(ctx context.Context, id string, sigHash string, notes string) error {
	sub, ok := m.subs[id]
	if !ok {
		return errors.New("sostituzione non trovata")
	}
	now := time.Now()
	sub.SignedBySubstitute = true
	sub.SignatureHash = sigHash
	sub.SignatureTimestamp = &now
	sub.OfficialRegisterNotes = notes
	return nil
}

func (m *memorySubstitutionsIntegrationRepo) GetAvailableTeachers(ctx context.Context, schoolID string) ([]substitutions.TeacherCandidate, error) {
	return m.available, nil
}

func (m *memorySubstitutionsIntegrationRepo) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	if p, ok := m.profileIDMapping[userID]; ok {
		return p, nil
	}
	return userID, nil
}

func (m *memorySubstitutionsIntegrationRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	return m.classTeachers[teacherID+":"+classID], nil
}

func (m *memorySubstitutionsIntegrationRepo) IsTeacherAssignedToSubject(ctx context.Context, teacherID, subjectID string) (bool, error) {
	return m.subjectTeachers[teacherID+":"+subjectID], nil
}

func (m *memorySubstitutionsIntegrationRepo) GetWeeklySubstitutionCount(ctx context.Context, teacherID string) (int, error) {
	return m.subCounts[teacherID], nil
}

func setupSubstitutionsIntegrationRouter(repo substitutions.Repository, currentUserID, currentRole, currentSchoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if currentUserID != "" {
			c.Set("user_id", currentUserID)
		}
		if currentRole != "" {
			c.Set("role", currentRole)
		}
		if currentSchoolID != "" {
			c.Set("school_id", currentSchoolID)
		}
		c.Next()
	})

	svc := substitutions.NewService(repo)
	handler := substitutions.NewHandler(svc)

	v1 := r.Group("/api/v1")
	handler.RegisterRoutes(v1)

	return r
}

func TestEmergencySubstitutionsWorkflow_Integration(t *testing.T) {
	repo := newMemorySubstitutionsIntegrationRepo()
	schoolID := "school-liceo-scientifico"

	// Preset available candidate teachers
	repo.available = []substitutions.TeacherCandidate{
		{TeacherID: "t-rossi", TeacherName: "Prof. Mario Rossi"},
		{TeacherID: "t-bianchi", TeacherName: "Prof.ssa Laura Bianchi"},
	}
	// Rossi teaches the same class (bonus +25)
	repo.classTeachers["t-rossi:class-3B"] = true
	repo.subjectTeachers["t-rossi:subj-fisica"] = true
	repo.subCounts["t-rossi"] = 0

	// 1. Unauthorized user (student) attempts to create an emergency substitution -> 403 Forbidden
	{
		rStudent := setupSubstitutionsIntegrationRouter(repo, "student-1", "student", schoolID)
		body, _ := json.Marshal(substitutions.CreateSubstitutionRequest{
			ClassID:         "class-3B",
			AbsentTeacherID: "t-verdi",
			Date:            "2026-10-20",
			Slot:            3,
			Subject:         "Fisica",
		})
		req := httptest.NewRequest("POST", "/api/v1/substitutions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rStudent.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	}

	// 2. Vice Principal posts an emergency substitution
	var createdSubID string
	{
		rVicePrincipal := setupSubstitutionsIntegrationRouter(repo, "vp-1", "vice_principal", schoolID)
		body, _ := json.Marshal(substitutions.CreateSubstitutionRequest{
			ClassID:         "class-3B",
			AbsentTeacherID: "t-verdi",
			Date:            "2026-10-20",
			Slot:            3,
			Subject:         "Fisica",
			SubjectID:       "subj-fisica",
			Notes:           "Docente assente per convocazione assemblea",
		})
		req := httptest.NewRequest("POST", "/api/v1/substitutions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rVicePrincipal.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		var sub substitutions.Substitution
		err := json.Unmarshal(w.Body.Bytes(), &sub)
		require.NoError(t, err)
		createdSubID = sub.ID
		assert.NotEmpty(t, createdSubID)
		assert.Equal(t, substitutions.StatusPending, sub.Status)
	}

	// 3. Query Recommendation Scoring Algorithm
	{
		rSecretary := setupSubstitutionsIntegrationRouter(repo, "sec-1", "secretary", schoolID)
		req := httptest.NewRequest("GET", "/api/v1/substitutions/recommend-substitutes?class_id=class-3B&subject_id=subj-fisica&date=2026-10-20&hour=3", nil)
		w := httptest.NewRecorder()
		rSecretary.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var recs []substitutions.SubstituteRecommendation
		err := json.Unmarshal(w.Body.Bytes(), &recs)
		require.NoError(t, err)
		require.Len(t, recs, 2)
		// Top recommendation should be Prof. Rossi with 95 score
		assert.Equal(t, "t-rossi", recs[0].TeacherID)
		assert.Equal(t, 95, recs[0].Score)
		assert.True(t, recs[0].TeachesClass)
		assert.True(t, recs[0].TeachesSubject)
	}

	// 4. Assign Substitute Teacher
	{
		rSecretary := setupSubstitutionsIntegrationRouter(repo, "sec-1", "secretary", schoolID)
		body, _ := json.Marshal(substitutions.AssignSubstituteRequest{
			SubstituteTeacherID: "t-rossi",
			Notes:               "Assegnato per competenza disciplinare",
		})
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/substitutions/%s/assign", createdSubID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rSecretary.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	}

	// 5. Today Summary Check
	{
		rVP := setupSubstitutionsIntegrationRouter(repo, "vp-1", "vice_principal", schoolID)
		req := httptest.NewRequest("GET", "/api/v1/substitutions/today-summary", nil)
		w := httptest.NewRecorder()
		rVP.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var summary map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &summary)
		require.NoError(t, err)
		assert.Equal(t, float64(1), summary["total"])
		assert.Equal(t, float64(1), summary["assigned"])
		assert.Equal(t, float64(0), summary["pending"])
	}

	// 6. Assigned Teacher Electronic Register Signing
	{
		// First try unauthorized teacher
		rOtherTeacher := setupSubstitutionsIntegrationRouter(repo, "t-bianchi", "teacher", schoolID)
		body, _ := json.Marshal(map[string]string{"notes": "Tentativo non autorizzato"})
		req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/substitutions/%s/sign-register", createdSubID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rOtherTeacher.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Legitimate assigned teacher signs
		rAssignedTeacher := setupSubstitutionsIntegrationRouter(repo, "t-rossi", "teacher", schoolID)
		bodyValid, _ := json.Marshal(map[string]string{"notes": "Svolto esperimento sul moto rettilineo uniforme in laboratorio"})
		reqValid := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/substitutions/%s/sign-register", createdSubID), bytes.NewReader(bodyValid))
		reqValid.Header.Set("Content-Type", "application/json")
		wValid := httptest.NewRecorder()
		rAssignedTeacher.ServeHTTP(wValid, reqValid)

		require.Equal(t, http.StatusOK, wValid.Code)
		var res map[string]interface{}
		_ = json.Unmarshal(wValid.Body.Bytes(), &res)
		assert.Equal(t, "registro supplenze firmato con successo", res["message"])

		// Verify substitution record is marked signed
		sub, err := repo.GetByID(context.Background(), createdSubID)
		require.NoError(t, err)
		assert.True(t, sub.SignedBySubstitute)
		assert.NotEmpty(t, sub.SignatureHash)
		assert.Equal(t, "Svolto esperimento sul moto rettilineo uniforme in laboratorio", sub.OfficialRegisterNotes)
	}
}
