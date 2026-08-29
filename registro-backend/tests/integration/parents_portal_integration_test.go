package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/attendance"
	"registro-backend/internal/communications"
	"registro-backend/internal/grades"
	"registro-backend/internal/parents"
	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Parents Repositories ─────────────────────────────────────────────

type mockParentsRepo struct {
	childrenMap map[string][]string // parentUserID -> [studentUserIDs]
}

func (m *mockParentsRepo) GetChildrenByParentUserID(_ context.Context, parentUserID string) ([]string, error) {
	return m.childrenMap[parentUserID], nil
}

type mockUsersForParentsRepo struct {
	users.Repository
	childrenMap map[string][]users.StudentChild
	guardians   map[string]map[string]bool // parentUserID -> studentUserID -> true
}

func (m *mockUsersForParentsRepo) GetChildren(_ context.Context, parentUserID string) ([]users.StudentChild, error) {
	return m.childrenMap[parentUserID], nil
}

func (m *mockUsersForParentsRepo) GetByID(_ context.Context, id string) (*users.User, error) {
	schoolID := "school-1"
	return &users.User{
		ID:       id,
		SchoolID: &schoolID,
	}, nil
}

func (m *mockUsersForParentsRepo) IsGuardian(_ context.Context, parentUserID string, studentUserID string) (bool, error) {
	if m.guardians[parentUserID] != nil {
		return m.guardians[parentUserID][studentUserID], nil
	}
	return false, nil
}

type mockGradesForParentsRepo struct {
	grades.Repository
	studentGrades map[string][]grades.Grade
}

func (m *mockGradesForParentsRepo) FindByStudent(studentID string) ([]grades.Grade, error) {
	return m.studentGrades[studentID], nil
}

type mockAttendanceForParentsRepo struct {
	attendance.Repository
}

func (m *mockAttendanceForParentsRepo) GetStats(_ string) (*attendance.SummaryResponse, error) {
	return &attendance.SummaryResponse{
		TotalAbsences: 2,
		TotalLates:    1,
	}, nil
}

type mockCommsForParentsRepo struct {
	communications.Repository
}

func (m *mockCommsForParentsRepo) ListBacheca(_ context.Context, _, _ string) ([]*communications.Message, error) {
	return []*communications.Message{}, nil
}

// ─── Setup Router Helper ───────────────────────────────────────────────────

func setupParentsRouter(
	pRepo *mockParentsRepo,
	uRepo *mockUsersForParentsRepo,
	gRepo *mockGradesForParentsRepo,
	parentUserID string,
) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := parents.NewService(
		pRepo,
		uRepo,
		gRepo,
		&mockAttendanceForParentsRepo{},
		&mockCommsForParentsRepo{},
	)
	h := parents.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		if parentUserID != "" {
			c.Set("user_id", parentUserID)
			c.Set("role", "parent")
			c.Set("school_id", "school-1")
		}
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

// ─── Tests ─────────────────────────────────────────────────────────────────

// PAR01 — Parent retrieves dashboard containing linked children with grades summary
func TestParents_Dashboard_Success(t *testing.T) {
	pRepo := &mockParentsRepo{childrenMap: map[string][]string{
		"parent-1": {"student-u1"},
	}}
	uRepo := &mockUsersForParentsRepo{
		childrenMap: map[string][]users.StudentChild{
			"parent-1": {
				{
					ID:         "prof-1",
					UserID:     "student-u1",
					FirstName:  "Marco",
					LastName:   "Rossi",
					Class:      "3A",
					SchoolName: "Liceo Scientifico",
				},
			},
		},
		guardians: map[string]map[string]bool{
			"parent-1": {"student-u1": true},
		},
	}
	gRepo := &mockGradesForParentsRepo{
		studentGrades: map[string][]grades.Grade{
			"student-u1": {
				{ID: "g-1", GradeValue: 8.0, IsPublished: true},
				{ID: "g-2", GradeValue: 7.0, IsPublished: true},
			},
		},
	}

	r := setupParentsRouter(pRepo, uRepo, gRepo, "parent-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/parents/dashboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp parents.ParentDashboardResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "parent-1", resp.ParentID)
	require.Len(t, resp.Children, 1)
	assert.Equal(t, "Marco", resp.Children[0].Student.FirstName)
	assert.Equal(t, 7.5, resp.Children[0].AverageGrade)
}

// PAR02 — Parent with no linked children gets empty children list
func TestParents_Dashboard_EmptyWhenNoChildren(t *testing.T) {
	pRepo := &mockParentsRepo{childrenMap: map[string][]string{}}
	uRepo := &mockUsersForParentsRepo{
		childrenMap: map[string][]users.StudentChild{},
		guardians:   map[string]map[string]bool{},
	}
	gRepo := &mockGradesForParentsRepo{studentGrades: map[string][]grades.Grade{}}

	r := setupParentsRouter(pRepo, uRepo, gRepo, "parent-no-kids")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/parents/dashboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp parents.ParentDashboardResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Empty(t, resp.Children)
}

// PAR03 — Parent retrieves specific child grades average
func TestParents_GetChildGradesAverage_Authorized(t *testing.T) {
	pRepo := &mockParentsRepo{}
	uRepo := &mockUsersForParentsRepo{
		guardians: map[string]map[string]bool{
			"parent-1": {"student-u1": true},
		},
	}
	gRepo := &mockGradesForParentsRepo{
		studentGrades: map[string][]grades.Grade{
			"student-u1": {
				{ID: "g-1", GradeValue: 9.0, IsPublished: true},
				{ID: "g-2", GradeValue: 9.0, IsPublished: true},
			},
		},
	}

	r := setupParentsRouter(pRepo, uRepo, gRepo, "parent-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/parents/child/student-u1/grades-average", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, float64(9.0), body["average"])
}

// PAR04 — Parent cannot query student they are not a guardian of
func TestParents_GetChildGradesAverage_UnauthorizedForbidden(t *testing.T) {
	pRepo := &mockParentsRepo{}
	uRepo := &mockUsersForParentsRepo{
		guardians: map[string]map[string]bool{
			"parent-1": {"other-student-u99": false},
		},
	}
	gRepo := &mockGradesForParentsRepo{studentGrades: map[string][]grades.Grade{}}

	r := setupParentsRouter(pRepo, uRepo, gRepo, "parent-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/parents/child/other-student-u99/grades-average", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// PAR05 — Unauthenticated requests return 401
func TestParents_Unauthorized_NoAuth(t *testing.T) {
	pRepo := &mockParentsRepo{}
	uRepo := &mockUsersForParentsRepo{}
	gRepo := &mockGradesForParentsRepo{}

	r := setupParentsRouter(pRepo, uRepo, gRepo, "")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/parents/dashboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
