package competencies_invalsi_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/competencies_invalsi"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) UpsertCertification(c *competencies_invalsi.CompetencyCertificationDM14) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *mockRepo) GetCertification(id string) (*competencies_invalsi.CompetencyCertificationDM14, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*competencies_invalsi.CompetencyCertificationDM14), args.Error(1)
}

func (m *mockRepo) GetCertificationByStudent(studentID, schoolYear, schoolTier string) (*competencies_invalsi.CompetencyCertificationDM14, error) {
	args := m.Called(studentID, schoolYear, schoolTier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*competencies_invalsi.CompetencyCertificationDM14), args.Error(1)
}

func (m *mockRepo) ListCertificationsByClass(classID, schoolYear string) ([]competencies_invalsi.CompetencyCertificationDM14, error) {
	args := m.Called(classID, schoolYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]competencies_invalsi.CompetencyCertificationDM14), args.Error(1)
}

func (m *mockRepo) SignCertification(id, principalID string) error {
	args := m.Called(id, principalID)
	return args.Error(0)
}

func (m *mockRepo) InsertInvalsiBatch(results []competencies_invalsi.InvalsiResult) error {
	args := m.Called(results)
	return args.Error(0)
}

func (m *mockRepo) ListInvalsiByClass(classID, schoolYear string) ([]competencies_invalsi.InvalsiResult, error) {
	args := m.Called(classID, schoolYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]competencies_invalsi.InvalsiResult), args.Error(1)
}

func (m *mockRepo) GetSchoolInvalsiAnalytics(schoolID, schoolYear string) ([]competencies_invalsi.InvalsiResult, error) {
	args := m.Called(schoolID, schoolYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]competencies_invalsi.InvalsiResult), args.Error(1)
}

func setupTestRouter(svc competencies_invalsi.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "sch-test-123")
		c.Set("user_id", "principal-test-789")
		c.Next()
	})
	h := competencies_invalsi.NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestCompetenciesInvalsiService_EdgeCases(t *testing.T) {
	repo := new(mockRepo)
	svc := competencies_invalsi.NewService(repo)

	// Nil safety
	assert.Error(t, svc.SaveCertification(nil))

	// GetCertification empty id
	_, err := svc.GetCertification("")
	assert.Error(t, err)

	// GetCertification success & not found
	repo.On("GetCertification", "cert-1").Return(&competencies_invalsi.CompetencyCertificationDM14{ID: "cert-1"}, nil).Once()
	c1, err := svc.GetCertification("cert-1")
	assert.NoError(t, err)
	assert.Equal(t, "cert-1", c1.ID)

	// GetCertificationByStudent empty student
	_, err = svc.GetCertificationByStudent("", "2025/2026", "primo_ciclo")
	assert.Error(t, err)

	repo.On("GetCertificationByStudent", "stud-1", "2025/2026", "primo_ciclo").Return(&competencies_invalsi.CompetencyCertificationDM14{StudentID: "stud-1"}, nil).Once()
	c2, err := svc.GetCertificationByStudent("stud-1", "2025/2026", "primo_ciclo")
	assert.NoError(t, err)
	assert.Equal(t, "stud-1", c2.StudentID)

	// ListCertificationsByClass
	_, err = svc.ListCertificationsByClass("", "2025/2026")
	assert.Error(t, err)

	repo.On("ListCertificationsByClass", "cls-1", "2025/2026").Return([]competencies_invalsi.CompetencyCertificationDM14{{ID: "cert-1"}}, nil).Once()
	list, err := svc.ListCertificationsByClass("cls-1", "2025/2026")
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// SignCertification
	assert.Error(t, svc.SignCertification("", "princ-1"))
	assert.Error(t, svc.SignCertification("cert-1", ""))

	repo.On("SignCertification", "cert-1", "princ-1").Return(nil).Once()
	assert.NoError(t, svc.SignCertification("cert-1", "princ-1"))

	// ImportInvalsiBatch empty
	assert.Error(t, svc.ImportInvalsiBatch(nil))
	assert.Error(t, svc.ImportInvalsiBatch([]competencies_invalsi.InvalsiResult{}))

	// ListInvalsiByClass
	_, err = svc.ListInvalsiByClass("", "2025/2026")
	assert.Error(t, err)

	repo.On("ListInvalsiByClass", "cls-1", "2025/2026").Return([]competencies_invalsi.InvalsiResult{{ID: "inv-1"}}, nil).Once()
	invList, err := svc.ListInvalsiByClass("cls-1", "2025/2026")
	assert.NoError(t, err)
	assert.Len(t, invList, 1)

	// GetSchoolAnalytics
	_, err = svc.GetSchoolAnalytics("", "2025/2026")
	assert.Error(t, err)

	repo.On("GetSchoolInvalsiAnalytics", "sch-1", "2025/2026").Return([]competencies_invalsi.InvalsiResult{{ID: "inv-1"}}, nil).Once()
	analytics, err := svc.GetSchoolAnalytics("sch-1", "2025/2026")
	assert.NoError(t, err)
	assert.Len(t, analytics, 1)

	repo.AssertExpectations(t)
}

func TestCompetenciesInvalsiHandler_AllEndpoints(t *testing.T) {
	repo := new(mockRepo)
	svc := competencies_invalsi.NewService(repo)
	r := setupTestRouter(svc)

	// 1. POST /certifications
	repo.On("UpsertCertification", mock.Anything).Return(nil).Once()
	payload := `{"student_id":"stud-1","class_id":"cls-1","school_tier":"primaria"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/competencies-dm14/certifications", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// POST bad json
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/competencies-dm14/certifications", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 2. GET /certifications/:id
	repo.On("GetCertification", "cert-1").Return(&competencies_invalsi.CompetencyCertificationDM14{ID: "cert-1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/competencies-dm14/certifications/cert-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetCertification", "missing").Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/competencies-dm14/certifications/missing", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 3. GET /certifications/by-student
	repo.On("GetCertificationByStudent", "stud-1", "2025/2026", "primo_ciclo").Return(&competencies_invalsi.CompetencyCertificationDM14{ID: "cert-1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/competencies-dm14/certifications/by-student?student_id=stud-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetCertificationByStudent", "missing", "2025/2026", "primo_ciclo").Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/competencies-dm14/certifications/by-student?student_id=missing", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 4. GET /certifications/by-class
	repo.On("ListCertificationsByClass", "cls-1", "2025/2026").Return([]competencies_invalsi.CompetencyCertificationDM14{{ID: "c1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/competencies-dm14/certifications/by-class?class_id=cls-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListCertificationsByClass", "cls-err", "2025/2026").Return(nil, errors.New("err")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/competencies-dm14/certifications/by-class?class_id=cls-err", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 5. POST /certifications/:id/sign
	repo.On("SignCertification", "cert-1", "principal-test-789").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/competencies-dm14/certifications/cert-1/sign", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. POST /invalsi/batch
	repo.On("InsertInvalsiBatch", mock.Anything).Return(nil).Once()
	invalsiPayload := `{"results":[{"student_id":"s1","subject":"matematica","level":3}]}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/competencies-dm14/invalsi/batch", bytes.NewBufferString(invalsiPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/competencies-dm14/invalsi/batch", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 7. GET /invalsi/by-class
	repo.On("ListInvalsiByClass", "cls-1", "2025/2026").Return([]competencies_invalsi.InvalsiResult{{ID: "inv-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/competencies-dm14/invalsi/by-class?class_id=cls-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListInvalsiByClass", "cls-err", "2025/2026").Return(nil, errors.New("err")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/competencies-dm14/invalsi/by-class?class_id=cls-err", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 8. GET /invalsi/analytics
	repo.On("GetSchoolInvalsiAnalytics", "sch-test-123", "2025/2026").Return([]competencies_invalsi.InvalsiResult{{ID: "inv-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/competencies-dm14/invalsi/analytics", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetSchoolInvalsiAnalytics", "sch-test-123", "2025/2026").Return(nil, errors.New("err")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/competencies-dm14/invalsi/analytics", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	repo.AssertExpectations(t)
}
