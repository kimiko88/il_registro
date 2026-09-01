package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/substitutions"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSubstitutionsFullRepo struct {
	subs map[string]*substitutions.Substitution
}

func newMockSubstitutionsFullRepo() *mockSubstitutionsFullRepo {
	return &mockSubstitutionsFullRepo{
		subs: make(map[string]*substitutions.Substitution),
	}
}

func (m *mockSubstitutionsFullRepo) Create(ctx context.Context, s *substitutions.Substitution) error {
	if s.ID == "" {
		s.ID = "sub-" + time.Now().Format("150405.000000")
	}
	s.CreatedAt = time.Now()
	m.subs[s.ID] = s
	return nil
}

func (m *mockSubstitutionsFullRepo) GetByID(ctx context.Context, id string) (*substitutions.Substitution, error) {
	s, ok := m.subs[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *mockSubstitutionsFullRepo) ListBySchool(ctx context.Context, schoolID, date string) ([]*substitutions.Substitution, error) {
	var res []*substitutions.Substitution
	for _, s := range m.subs {
		if schoolID != "" && s.SchoolID != schoolID {
			continue
		}
		res = append(res, s)
	}
	return res, nil
}

func (m *mockSubstitutionsFullRepo) ListByTeacher(ctx context.Context, teacherID string, date string) ([]*substitutions.Substitution, error) {
	var res []*substitutions.Substitution
	for _, s := range m.subs {
		if s.SubstituteTeacherID != nil && *s.SubstituteTeacherID == teacherID {
			res = append(res, s)
		}
	}
	return res, nil
}

func (m *mockSubstitutionsFullRepo) AssignSubstitute(ctx context.Context, id string, substituteTeacherID string, notes string) error {
	if s, ok := m.subs[id]; ok {
		s.SubstituteTeacherID = &substituteTeacherID
		s.Status = substitutions.StatusAssigned
		s.Notes = notes
	}
	return nil
}

func (m *mockSubstitutionsFullRepo) ConfirmSubstitution(ctx context.Context, id string, substituteTeacherID string) error {
	if s, ok := m.subs[id]; ok {
		s.Status = substitutions.StatusConfirmed
	}
	return nil
}

func (m *mockSubstitutionsFullRepo) SignRegister(ctx context.Context, id string, sigHash string, notes string) error {
	if s, ok := m.subs[id]; ok {
		s.SignedBySubstitute = true
		now := time.Now()
		s.SignatureTimestamp = &now
		s.SignatureHash = sigHash
		s.OfficialRegisterNotes = notes
	}
	return nil
}

func (m *mockSubstitutionsFullRepo) GetAvailableTeachers(ctx context.Context, schoolID string) ([]substitutions.TeacherCandidate, error) {
	return []substitutions.TeacherCandidate{
		{TeacherID: "t-cand-1", UserID: "u-cand-1", TeacherName: "Prof. Sostituto"},
	}, nil
}

func (m *mockSubstitutionsFullRepo) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	return userID, nil
}

func (m *mockSubstitutionsFullRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	return true, nil
}

func (m *mockSubstitutionsFullRepo) IsTeacherAssignedToSubject(ctx context.Context, teacherID, subjectID string) (bool, error) {
	return true, nil
}

func (m *mockSubstitutionsFullRepo) GetWeeklySubstitutionCount(ctx context.Context, teacherID string) (int, error) {
	return 2, nil
}

func setupSubstitutionsTestRouter(repo substitutions.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := substitutions.NewService(repo, nil, nil, nil)
	handler := substitutions.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = "admin"
		}
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "admin-1"
		}
		c.Set("role", role)
		c.Set("user_id", userID)
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_Substitutions_FullLifecycle_Workflow(t *testing.T) {
	repo := newMockSubstitutionsFullRepo()
	r := setupSubstitutionsTestRouter(repo)

	today := time.Now().Format("2006-01-02")

	// 1. Admin/Presidenza creates a substitution request for absent teacher
	createReq := substitutions.CreateSubstitutionRequest{
		ClassID:         "class-3B",
		AbsentTeacherID: "teacher-absent-1",
		Date:            today,
		Hour:            3,
		SubjectID:       "sub-history",
		Notes:           "Docente assente per malattia",
	}
	body, _ := json.Marshal(createReq)
	reqCreate, _ := http.NewRequest("POST", "/api/v1/substitutions", bytes.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("X-Role", "admin")
	wCreate := httptest.NewRecorder()
	r.ServeHTTP(wCreate, reqCreate)
	require.Equal(t, http.StatusCreated, wCreate.Code)

	var created substitutions.Substitution
	err := json.Unmarshal(wCreate.Body.Bytes(), &created)
	require.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, substitutions.StatusPending, created.Status)

	// 2. Assign substitute teacher
	assignReq := substitutions.AssignSubstituteRequest{
		SubstituteTeacherID: "teacher-sub-2",
		Notes:               "Assegnata ora di supplenza da orario a disposizione",
	}
	bodyAssign, _ := json.Marshal(assignReq)
	reqAssign, _ := http.NewRequest("PUT", "/api/v1/substitutions/"+created.ID+"/assign", bytes.NewReader(bodyAssign))
	reqAssign.Header.Set("Content-Type", "application/json")
	reqAssign.Header.Set("X-Role", "admin")
	wAssign := httptest.NewRecorder()
	r.ServeHTTP(wAssign, reqAssign)
	require.Equal(t, http.StatusOK, wAssign.Code)

	// 3. Substitute teacher confirms acceptance
	reqConfirm, _ := http.NewRequest("PATCH", "/api/v1/substitutions/"+created.ID+"/confirm", nil)
	reqConfirm.Header.Set("X-Role", "teacher")
	reqConfirm.Header.Set("X-User-ID", "teacher-sub-2")
	wConfirm := httptest.NewRecorder()
	r.ServeHTTP(wConfirm, reqConfirm)
	require.Equal(t, http.StatusOK, wConfirm.Code)

	// 4. Substitute signs class register with official notes and digital hash
	signReq := map[string]string{
		"notes": "Sostituzione svolta regolarmente: completamento esercitazione su Carlo Magno",
	}
	bodySign, _ := json.Marshal(signReq)
	reqSign, _ := http.NewRequest("POST", "/api/v1/substitutions/"+created.ID+"/sign-register", bytes.NewReader(bodySign))
	reqSign.Header.Set("Content-Type", "application/json")
	reqSign.Header.Set("X-Role", "teacher")
	reqSign.Header.Set("X-User-ID", "teacher-sub-2")
	wSign := httptest.NewRecorder()
	r.ServeHTTP(wSign, reqSign)
	require.Equal(t, http.StatusOK, wSign.Code)

	// 5. Verify substitution record in list
	reqList, _ := http.NewRequest("GET", "/api/v1/substitutions", nil)
	reqList.Header.Set("X-Role", "admin")
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)
	require.Equal(t, http.StatusOK, wList.Code)

	var list []*substitutions.Substitution
	err = json.Unmarshal(wList.Body.Bytes(), &list)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.True(t, list[0].SignedBySubstitute)
}
