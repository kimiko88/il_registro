package security_passes_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/security_passes"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreateHallPass(p *security_passes.StudentHallPass) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *mockRepo) CloseHallPass(id string, returnTime time.Time) error {
	args := m.Called(id, returnTime)
	return args.Error(0)
}

func (m *mockRepo) ListActiveHallPassesBySchool(schoolID string) ([]security_passes.StudentHallPass, error) {
	args := m.Called(schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]security_passes.StudentHallPass), args.Error(1)
}

func (m *mockRepo) ListHallPassesByClass(classID string, limit int) ([]security_passes.StudentHallPass, error) {
	args := m.Called(classID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]security_passes.StudentHallPass), args.Error(1)
}

func (m *mockRepo) CreatePickupDelegate(d *security_passes.StudentPickupDelegate) error {
	args := m.Called(d)
	return args.Error(0)
}

func (m *mockRepo) ListDelegatesByStudent(studentID string) ([]security_passes.StudentPickupDelegate, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]security_passes.StudentPickupDelegate), args.Error(1)
}

func (m *mockRepo) RevokeDelegate(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter(svc security_passes.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "sch-test-123")
		c.Set("user_id", "teacher-test-456")
		c.Next()
	})
	h := security_passes.NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestSecurityPassesService_Full(t *testing.T) {
	repo := new(mockRepo)
	svc := security_passes.NewService(repo)

	// IssueHallPass validations
	assert.Error(t, svc.IssueHallPass(nil))
	assert.Error(t, svc.IssueHallPass(&security_passes.StudentHallPass{SchoolID: ""}))
	assert.Error(t, svc.IssueHallPass(&security_passes.StudentHallPass{SchoolID: "sch", ClassID: ""}))
	assert.Error(t, svc.IssueHallPass(&security_passes.StudentHallPass{SchoolID: "sch", ClassID: "cls", StudentID: ""}))
	assert.Error(t, svc.IssueHallPass(&security_passes.StudentHallPass{SchoolID: "sch", ClassID: "cls", StudentID: "s1", TeacherID: ""}))

	repo.On("CreateHallPass", mock.Anything).Return(nil).Once()
	p := &security_passes.StudentHallPass{SchoolID: "sch", ClassID: "cls", StudentID: "s1", TeacherID: "t1"}
	assert.NoError(t, svc.IssueHallPass(p))
	assert.Equal(t, "bagno", p.PassType)
	assert.Equal(t, "fuori_aula", p.Status)

	// ReturnHallPass
	assert.Error(t, svc.ReturnHallPass(""))
	repo.On("CloseHallPass", "p1", mock.Anything).Return(nil).Once()
	assert.NoError(t, svc.ReturnHallPass("p1"))

	// GetActivePasses
	_, err := svc.GetActivePasses("")
	assert.Error(t, err)

	repo.On("ListActiveHallPassesBySchool", "sch-1").Return([]security_passes.StudentHallPass{{ID: "p1"}}, nil).Once()
	passes, err := svc.GetActivePasses("sch-1")
	assert.NoError(t, err)
	assert.Len(t, passes, 1)

	// GetClassPasses
	_, err = svc.GetClassPasses("", 10)
	assert.Error(t, err)

	repo.On("ListHallPassesByClass", "cls-1", 30).Return([]security_passes.StudentHallPass{{ID: "p1"}}, nil).Once()
	cpasses, err := svc.GetClassPasses("cls-1", 0)
	assert.NoError(t, err)
	assert.Len(t, cpasses, 1)

	// AddPickupDelegate
	assert.Error(t, svc.AddPickupDelegate(nil))
	assert.Error(t, svc.AddPickupDelegate(&security_passes.StudentPickupDelegate{StudentID: ""}))
	assert.Error(t, svc.AddPickupDelegate(&security_passes.StudentPickupDelegate{StudentID: "s1", DelegateName: ""}))
	assert.Error(t, svc.AddPickupDelegate(&security_passes.StudentPickupDelegate{StudentID: "s1", DelegateName: "Nonna", FiscalCode: ""}))
	assert.Error(t, svc.AddPickupDelegate(&security_passes.StudentPickupDelegate{StudentID: "s1", DelegateName: "Nonna", FiscalCode: "RSSMRA...", Phone: ""}))

	repo.On("CreatePickupDelegate", mock.Anything).Return(nil).Once()
	del := &security_passes.StudentPickupDelegate{StudentID: "s1", DelegateName: "Nonna", FiscalCode: "RSSMRA...", Phone: "123456"}
	assert.NoError(t, svc.AddPickupDelegate(del))
	assert.Equal(t, "attivo", del.Status)

	// GetStudentDelegates
	_, err = svc.GetStudentDelegates("")
	assert.Error(t, err)

	repo.On("ListDelegatesByStudent", "s1").Return([]security_passes.StudentPickupDelegate{{ID: "d1"}}, nil).Once()
	delList, err := svc.GetStudentDelegates("s1")
	assert.NoError(t, err)
	assert.Len(t, delList, 1)

	// RevokePickupDelegate
	assert.Error(t, svc.RevokePickupDelegate(""))

	repo.On("RevokeDelegate", "d1").Return(nil).Once()
	assert.NoError(t, svc.RevokePickupDelegate("d1"))

	repo.AssertExpectations(t)
}

func TestSecurityPassesHandler_Endpoints(t *testing.T) {
	repo := new(mockRepo)
	svc := security_passes.NewService(repo)
	r := setupTestRouter(svc)

	// 1. POST /hall-passes
	repo.On("CreateHallPass", mock.Anything).Return(nil).Once()
	payload := `{"class_id":"cls-1","student_id":"s1","pass_type":"infermeria"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/security-passes/hall-passes", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/security-passes/hall-passes", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 2. POST /hall-passes/:id/return
	repo.On("CloseHallPass", "p1", mock.Anything).Return(nil).Once()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/security-passes/hall-passes/p1/return", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. GET /hall-passes/active
	repo.On("ListActiveHallPassesBySchool", "sch-test-123").Return([]security_passes.StudentHallPass{{ID: "p1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/security-passes/hall-passes/active", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListActiveHallPassesBySchool", "sch-test-123").Return(nil, errors.New("err")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/security-passes/hall-passes/active", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 4. GET /hall-passes/by-class
	repo.On("ListHallPassesByClass", "cls-1", 30).Return([]security_passes.StudentHallPass{{ID: "p1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/security-passes/hall-passes/by-class?class_id=cls-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListHallPassesByClass", "cls-err", 30).Return(nil, errors.New("err")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/security-passes/hall-passes/by-class?class_id=cls-err", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 5. POST /delegates
	repo.On("CreatePickupDelegate", mock.Anything).Return(nil).Once()
	delPayload := `{"student_id":"s1","delegate_name":"Zio Mario","fiscal_code":"CF123","phone":"3331234567"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/security-passes/delegates", bytes.NewBufferString(delPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/security-passes/delegates", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 6. GET /delegates/by-student
	repo.On("ListDelegatesByStudent", "s1").Return([]security_passes.StudentPickupDelegate{{ID: "d1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/security-passes/delegates/by-student?student_id=s1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListDelegatesByStudent", "serr").Return(nil, errors.New("err")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/security-passes/delegates/by-student?student_id=serr", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 7. POST /delegates/:id/revoke
	repo.On("RevokeDelegate", "d1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/security-passes/delegates/d1/revoke", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.AssertExpectations(t)
}
