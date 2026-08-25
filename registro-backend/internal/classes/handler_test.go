package classes

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupClassesHandlerTest() (*Handler, *MockRepository) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	handler := NewHandler(service)
	return handler, mockRepo
}

func TestHandler_GetLessonTopics(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "unauthorized user",
			userID: "",
			setupMock: func(m *MockRepository) {
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "successful get lesson topics",
			userID: "user-123",
			setupMock: func(m *MockRepository) {
				m.On("Get", mock.Anything, "class-1").Maybe().Return(&Class{ID: "class-1", SchoolID: "162737ff-081f-436c-8874-11cd57bc60f1"}, nil)
				m.On("GetLessonTopics", mock.Anything, "class-1").Return([]LessonTopic{
					{
						ID:               "lt-1",
						Date:             time.Now(),
						SubjectName:      "Matematica",
						Topic:            "Algebre di Boole",
						TeacherFirstName: "Mario",
						TeacherLastName:  "Rossi",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var topics []LessonTopic
				err := json.Unmarshal(w.Body.Bytes(), &topics)
				assert.NoError(t, err)
				assert.Len(t, topics, 1)
				assert.Equal(t, "Algebre di Boole", topics[0].Topic)
			},
		},
		{
			name:   "repository error",
			userID: "user-123",
			setupMock: func(m *MockRepository) {
				m.On("Get", mock.Anything, "class-1").Maybe().Return(&Class{ID: "class-1", SchoolID: "162737ff-081f-436c-8874-11cd57bc60f1"}, nil)
				m.On("GetLessonTopics", mock.Anything, "class-1").Return(nil, errors.New("db failure"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo := setupClassesHandlerTest()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tt.userID != "" {
				c.Set("user_id", tt.userID)
				c.Set("role", "teacher")
			}
			c.Params = gin.Params{{Key: "id", Value: "class-1"}}
			c.Request = httptest.NewRequest("GET", "/classes/class-1/lesson-topics", nil)

			handler.GetLessonTopics(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetDisciplinaryNotes(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "unauthorized user",
			userID: "",
			setupMock: func(m *MockRepository) {
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "successful get disciplinary notes",
			userID: "user-123",
			setupMock: func(m *MockRepository) {
				m.On("Get", mock.Anything, "class-1").Maybe().Return(&Class{ID: "class-1", SchoolID: "162737ff-081f-436c-8874-11cd57bc60f1"}, nil)
				m.On("GetDisciplinaryNotes", mock.Anything, "class-1").Return([]DisciplinaryNoteReport{
					{
						ID:               "dn-1",
						Date:             time.Now(),
						StudentFirstName: "Luigi",
						StudentLastName:  "Verdi",
						NoteType:         "disciplinary",
						Description:      "Nota di condotta",
						TeacherFirstName: "Mario",
						TeacherLastName:  "Rossi",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var notes []DisciplinaryNoteReport
				err := json.Unmarshal(w.Body.Bytes(), &notes)
				assert.NoError(t, err)
				assert.Len(t, notes, 1)
				assert.Equal(t, "Nota di condotta", notes[0].Description)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo := setupClassesHandlerTest()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tt.userID != "" {
				c.Set("user_id", tt.userID)
				c.Set("role", "teacher")
			}
			c.Params = gin.Params{{Key: "id", Value: "class-1"}}
			c.Request = httptest.NewRequest("GET", "/classes/class-1/disciplinary-notes", nil)

			handler.GetDisciplinaryNotes(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_List_SuperadminAllowedWithoutSchoolID(t *testing.T) {
	handler, mockRepo := setupClassesHandlerTest()
	mockRepo.On("List", mock.Anything, "", "").Return([]Class{
		{ID: "class-1", Name: "1A", SchoolID: "school-1"},
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", "superadmin-user-id")
	c.Set("role", "superadmin")
	c.Request = httptest.NewRequest("GET", "/classes", nil)

	handler.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var res []Class
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "1A", res[0].Name)
	mockRepo.AssertExpectations(t)
}
