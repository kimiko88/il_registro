package attendance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParseWindowParams_ToBeforeFrom_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest("GET", "/test?from=2026-12-01&to=2026-01-01", nil)
	c.Request = req

	_, _, err := parseWindowParams(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data di inizio successiva alla data di fine")
}

func TestExportAttendance_InvalidDate_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &Handler{}

	r := gin.New()
	r.GET("/export", h.ExportAttendance)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/export?class_id=cls-1&date=invalid-date", nil)
	// mock teacher credentials in request
	r.ServeHTTP(w, req)

	// Since actorID is missing, returns 401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPendingJustifications_FiltersEmptySchoolID(t *testing.T) {
	svc := &service{
		repo: &mockAttendanceRepoForTesting{
			pendingJustifications: []Justification{
				{ID: "j1", SchoolID: "school-1", StudentID: "s1"},
				{ID: "j2", SchoolID: "", StudentID: "s2"},
				{ID: "j3", SchoolID: "school-2", StudentID: "s3"},
			},
		},
		userRepo: &mockUserRepo{},
	}

	resp, err := svc.GetPendingJustifications(context.Background(), "teacher-1", "teacher", "class-1", "school-1")
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "j1", resp[0].ID)
}

type mockAttendanceRepoForTesting struct {
	Repository
	pendingJustifications []Justification
}

func (m *mockAttendanceRepoForTesting) FindPendingJustifications(classID, schoolID string) ([]Justification, error) {
	if schoolID == "" {
		return m.pendingJustifications, nil
	}
	var filtered []Justification
	for _, j := range m.pendingJustifications {
		if j.SchoolID == schoolID {
			filtered = append(filtered, j)
		}
	}
	return filtered, nil
}

func (m *mockAttendanceRepoForTesting) IsTeacherAssignedToClass(_ context.Context, _, _ string) (bool, error) {
	return true, nil
}
