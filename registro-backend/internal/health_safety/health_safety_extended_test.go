package health_safety_test

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
	"registro-backend/internal/health_safety"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreateMedication(med *health_safety.StudentMedication) error {
	args := m.Called(med)
	return args.Error(0)
}

func (m *mockRepo) ListMedicationsBySchool(schoolID, classID string) ([]health_safety.StudentMedication, error) {
	args := m.Called(schoolID, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]health_safety.StudentMedication), args.Error(1)
}

func (m *mockRepo) GetMedication(id string) (*health_safety.StudentMedication, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*health_safety.StudentMedication), args.Error(1)
}

func (m *mockRepo) UpdateMedicationStatus(id, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *mockRepo) LogMedicationAdministration(a *health_safety.MedicationAdministration) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m *mockRepo) ListMedicationAdministrations(medicationID string) ([]health_safety.MedicationAdministration, error) {
	args := m.Called(medicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]health_safety.MedicationAdministration), args.Error(1)
}

func (m *mockRepo) CreateSchoolAccident(a *health_safety.SchoolAccident) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m *mockRepo) ListSchoolAccidents(schoolID string) ([]health_safety.SchoolAccident, error) {
	args := m.Called(schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]health_safety.SchoolAccident), args.Error(1)
}

func (m *mockRepo) GetSchoolAccident(id string) (*health_safety.SchoolAccident, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*health_safety.SchoolAccident), args.Error(1)
}

func (m *mockRepo) UpdateAccidentInailTransmission(id, protocolNumber string) error {
	args := m.Called(id, protocolNumber)
	return args.Error(0)
}

func (m *mockRepo) UpdateAccidentStatus(id, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func setupTestRouter(svc health_safety.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "sch-test-123")
		c.Set("user_id", "teacher-test-456")
		c.Next()
	})
	h := health_safety.NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestHealthSafetyService_Full(t *testing.T) {
	repo := new(mockRepo)
	svc := health_safety.NewService(repo)

	// RegisterMedication validations
	assert.Error(t, svc.RegisterMedication(nil))
	assert.Error(t, svc.RegisterMedication(&health_safety.StudentMedication{SchoolID: ""}))
	assert.Error(t, svc.RegisterMedication(&health_safety.StudentMedication{SchoolID: "sch", StudentID: ""}))
	assert.Error(t, svc.RegisterMedication(&health_safety.StudentMedication{SchoolID: "sch", StudentID: "s1", MedicationName: ""}))
	assert.Error(t, svc.RegisterMedication(&health_safety.StudentMedication{SchoolID: "sch", StudentID: "s1", MedicationName: "Insulina", Dosage: ""}))
	assert.Error(t, svc.RegisterMedication(&health_safety.StudentMedication{SchoolID: "sch", StudentID: "s1", MedicationName: "Insulina", Dosage: "10ui", AdministrationMode: ""}))

	repo.On("CreateMedication", mock.Anything).Return(nil).Once()
	med := &health_safety.StudentMedication{SchoolID: "sch", StudentID: "s1", MedicationName: "Insulina", Dosage: "10ui", AdministrationMode: "Sottocute"}
	assert.NoError(t, svc.RegisterMedication(med))
	assert.Equal(t, "attivo", med.Status)

	// ListMedications
	_, err := svc.ListMedications("", "cls-1")
	assert.Error(t, err)

	repo.On("ListMedicationsBySchool", "sch-1", "cls-1").Return([]health_safety.StudentMedication{{ID: "m1"}}, nil).Once()
	meds, err := svc.ListMedications("sch-1", "cls-1")
	assert.NoError(t, err)
	assert.Len(t, meds, 1)

	// GetMedication
	_, err = svc.GetMedication("")
	assert.Error(t, err)

	repo.On("GetMedication", "m1").Return(&health_safety.StudentMedication{ID: "m1"}, nil).Once()
	m1, err := svc.GetMedication("m1")
	assert.NoError(t, err)
	assert.Equal(t, "m1", m1.ID)

	// LogAdministration
	assert.Error(t, svc.LogAdministration(nil))
	assert.Error(t, svc.LogAdministration(&health_safety.MedicationAdministration{MedicationID: ""}))
	assert.Error(t, svc.LogAdministration(&health_safety.MedicationAdministration{MedicationID: "m1", AdministeredBy: ""}))
	assert.Error(t, svc.LogAdministration(&health_safety.MedicationAdministration{MedicationID: "m1", AdministeredBy: "t1", DoseGiven: ""}))

	repo.On("LogMedicationAdministration", mock.Anything).Return(nil).Once()
	assert.NoError(t, svc.LogAdministration(&health_safety.MedicationAdministration{MedicationID: "m1", AdministeredBy: "t1", DoseGiven: "10ui"}))

	// ListAdministrations
	_, err = svc.ListAdministrations("")
	assert.Error(t, err)

	repo.On("ListMedicationAdministrations", "m1").Return([]health_safety.MedicationAdministration{{ID: "a1"}}, nil).Once()
	adms, err := svc.ListAdministrations("m1")
	assert.NoError(t, err)
	assert.Len(t, adms, 1)

	// RecordAccident
	assert.Error(t, svc.RecordAccident(nil))
	assert.Error(t, svc.RecordAccident(&health_safety.SchoolAccident{SchoolID: ""}))
	assert.Error(t, svc.RecordAccident(&health_safety.SchoolAccident{SchoolID: "sch", PersonID: ""}))
	assert.Error(t, svc.RecordAccident(&health_safety.SchoolAccident{SchoolID: "sch", PersonID: "p1", DynamicsDescription: ""}))

	repo.On("CreateSchoolAccident", mock.Anything).Return(nil).Once()
	acc := &health_safety.SchoolAccident{SchoolID: "sch", PersonID: "p1", DynamicsDescription: "Caduta durante educazione fisica"}
	assert.NoError(t, svc.RecordAccident(acc))
	assert.Equal(t, "aperto", acc.Status)
	assert.NotEmpty(t, acc.AccidentNumber)
	assert.False(t, acc.InailDeadline.IsZero())

	// ListAccidents
	_, err = svc.ListAccidents("")
	assert.Error(t, err)

	repo.On("ListSchoolAccidents", "sch-1").Return([]health_safety.SchoolAccident{{ID: "acc-1"}}, nil).Once()
	acclist, err := svc.ListAccidents("sch-1")
	assert.NoError(t, err)
	assert.Len(t, acclist, 1)

	// GetAccident
	_, err = svc.GetAccident("")
	assert.Error(t, err)

	repo.On("GetSchoolAccident", "acc-1").Return(&health_safety.SchoolAccident{ID: "acc-1"}, nil).Once()
	acc1, err := svc.GetAccident("acc-1")
	assert.NoError(t, err)
	assert.Equal(t, "acc-1", acc1.ID)

	// TransmitInail
	assert.Error(t, svc.TransmitInail("", "INAIL-123"))
	assert.Error(t, svc.TransmitInail("acc-1", ""))

	repo.On("UpdateAccidentInailTransmission", "acc-1", "INAIL-123").Return(nil).Once()
	assert.NoError(t, svc.TransmitInail("acc-1", "INAIL-123"))

	repo.AssertExpectations(t)
}

func TestHealthSafetyHandler_Endpoints(t *testing.T) {
	repo := new(mockRepo)
	svc := health_safety.NewService(repo)
	r := setupTestRouter(svc)

	// 1. POST /medications
	repo.On("CreateMedication", mock.Anything).Return(nil).Once()
	medPayload := `{"student_id":"s1","medication_name":"Ventolin","dosage":"2 puff","administration_mode":"inalatore"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/health-safety/medications", bytes.NewBufferString(medPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/health-safety/medications", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 2. GET /medications
	repo.On("ListMedicationsBySchool", "sch-test-123", "cls-1").Return([]health_safety.StudentMedication{{ID: "m1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/health-safety/medications?class_id=cls-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListMedicationsBySchool", "sch-test-123", "cls-err").Return(nil, errors.New("err")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/health-safety/medications?class_id=cls-err", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 3. GET /medications/:id
	repo.On("GetMedication", "m1").Return(&health_safety.StudentMedication{ID: "m1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/health-safety/medications/m1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetMedication", "missing").Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/health-safety/medications/missing", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 4. POST /medications/:id/administer
	repo.On("LogMedicationAdministration", mock.Anything).Return(nil).Once()
	admPayload := `{"dose_given":"2 puff","notes":"Pre-palestra"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/health-safety/medications/m1/administer", bytes.NewBufferString(admPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/health-safety/medications/m1/administer", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 5. GET /medications/:id/administrations
	repo.On("ListMedicationAdministrations", "m1").Return([]health_safety.MedicationAdministration{{ID: "a1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/health-safety/medications/m1/administrations", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. POST /accidents
	repo.On("CreateSchoolAccident", mock.Anything).Return(nil).Once()
	accPayload := `{"person_id":"p1","person_role":"student","dynamics_description":"Distorsione caviglia","accident_time":"` + time.Now().Format(time.RFC3339) + `"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/health-safety/accidents", bytes.NewBufferString(accPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/health-safety/accidents", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 7. GET /accidents
	repo.On("ListSchoolAccidents", "sch-test-123").Return([]health_safety.SchoolAccident{{ID: "acc1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/health-safety/accidents", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 8. GET /accidents/:id
	repo.On("GetSchoolAccident", "acc1").Return(&health_safety.SchoolAccident{ID: "acc1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/health-safety/accidents/acc1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetSchoolAccident", "missing").Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/health-safety/accidents/missing", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 9. POST /accidents/:id/transmit-inail
	repo.On("UpdateAccidentInailTransmission", "acc1", "INAIL-PROT-456").Return(nil).Once()
	inailPayload := `{"protocol":"INAIL-PROT-456"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/health-safety/accidents/acc1/transmit-inail", bytes.NewBufferString(inailPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/health-safety/accidents/acc1/transmit-inail", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	repo.AssertExpectations(t)
}
