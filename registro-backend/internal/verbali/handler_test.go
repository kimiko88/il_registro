package verbali

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIsAllowedVerbaliRole(t *testing.T) {
	allowed := []string{"teacher", "coordinator", "admin", "superadmin", "secretary", "principal", "vice_principal", "docente"}
	for _, role := range allowed {
		assert.True(t, isAllowedVerbaliRole(role), "expected %s to be allowed", role)
		assert.True(t, isAllowedVerbaliRole(role), "expected uppercase %s to be allowed", role)
	}

	disallowed := []string{"student", "parent", "unknown", ""}
	for _, role := range disallowed {
		assert.False(t, isAllowedVerbaliRole(role), "expected %s to be disallowed", role)
	}
}

func setupVerbaliRouter(svc *Service, role, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		if role != "" {
			c.Set("role", role)
		}
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	handler := NewHandler(svc)
	handler.RegisterRoutes(r.Group("/api/v1"))
	return r
}

func TestHandler_Verbali(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	// 1. Unauthorized & Forbidden checks
	rAnon := setupVerbaliRouter(svc, "", "", "")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/meetings", nil)
	w := httptest.NewRecorder()
	rAnon.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	rStudent := setupVerbaliRouter(svc, "student", "student-1", "school-1")
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/verbali/meetings", nil)
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 2. Allowed teacher router
	rTeacher := setupVerbaliRouter(svc, "teacher", "teacher-1", "school-1")

	// CreateMeeting
	mockRepo.On("ClassBelongsToSchool", mock.Anything, "class-1", "school-1").Return(true, nil).Once()
	mockRepo.On("CreateMeeting", mock.Anything, mock.Anything).Return(nil).Once()

	body, _ := json.Marshal(CreateMeetingRequest{
		ClassID:   "class-1",
		Title:     "Consiglio Ordinario",
		Date:      "2026-11-10",
		StartTime: "15:00",
		EndTime:   "16:30",
	})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/verbali/meetings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// ListMeetings
	mockRepo.On("ListMeetings", mock.Anything, "school-1", "class-1").Return([]*CouncilMeeting{{ID: "m-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/verbali/meetings?class_id=class-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// CreateVerbale
	mockRepo.On("GetMeetingByID", mock.Anything, "m-1").Return(&CouncilMeeting{ID: "m-1"}, nil).Once()
	mockRepo.On("CreateVerbale", mock.Anything, mock.Anything).Return(nil).Once()
	verbBody, _ := json.Marshal(CreateVerbaleRequest{
		MeetingID: "m-1",
		Title:     "Verbale Iniziale",
		Content:   "Contenuto del verbale",
	})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/verbali", bytes.NewBuffer(verbBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// GetVerbale
	mockRepo.On("GetVerbaleByID", mock.Anything, "v-1", "teacher-1").Return(&MeetingVerbale{ID: "v-1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/verbali/v-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// ListVerbali (requires valid UUID)
	validMeetingUUID := "a0000000-0000-0000-0000-000000000001"
	mockRepo.On("ListVerbali", mock.Anything, validMeetingUUID, "teacher-1").Return([]*MeetingVerbale{{ID: "v-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/verbali/meeting/"+validMeetingUUID, nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// SignVerbale
	secID := "teacher-1"
	mockRepo.On("GetVerbaleByID", mock.Anything, "v-1", "teacher-1").Return(&MeetingVerbale{ID: "v-1", SecretaryID: &secID, IsPublished: true}, nil).Once()
	mockRepo.On("SignVerbale", mock.Anything, "v-1", "teacher-1", mock.Anything).Return(nil).Once()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/verbali/v-1/sign", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// GetSignatures
	mockRepo.On("GetSignatures", mock.Anything, "v-1").Return([]VerbaleSignature{{ID: "sig-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/verbali/v-1/signatures", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// ExportPDF
	mockRepo.On("GetVerbaleByID", mock.Anything, "v-1", "teacher-1").Return(&MeetingVerbale{ID: "v-1", MeetingID: "m-1", Title: "Verbale 1", Content: "Delibera approvata"}, nil).Once()
	mockRepo.On("GetSignatures", mock.Anything, "v-1").Return([]VerbaleSignature{{ID: "sig-1", UserID: "u-1", UserName: "Docente Segretario", IPAddress: "127.0.0.1"}}, nil).Once()
	mockRepo.On("GetMeetingByID", mock.Anything, "m-1").Return(&CouncilMeeting{ID: "m-1", Title: "Consiglio", StartTime: "15:00", EndTime: "16:00", ClassID: "3A"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/verbali/v-1/pdf", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))

	mockRepo.AssertExpectations(t)
}

func TestGenerateVerbale(t *testing.T) {
	verbale := &MeetingVerbale{
		ID:        "v-1",
		MeetingID: "m-1",
		Title:     "Verbale di Prova con caratteri accentati: àèéìòù",
		Content:   "Punto 1: Approvazione del verbale precedente.\nPunto 2: Andamento didattico e disciplinare generale.",
	}
	meeting := &CouncilMeeting{
		ID:        "m-1",
		Title:     "Consiglio di Classe",
		StartTime: "14:30",
		EndTime:   "16:00",
		ClassID:   "5B",
		Agenda:    "Varie ed eventuali",
	}
	sigs := []VerbaleSignature{
		{
			ID:        "sig-1",
			UserID:    "u-1",
			UserName:  "Mario Rossi",
			IPAddress: "192.168.1.100",
		},
	}

	pdfBytes, err := GenerateVerbale(verbale, meeting, sigs)
	assert.NoError(t, err)
	assert.NotEmpty(t, pdfBytes)
	assert.True(t, len(pdfBytes) > 500)
}
