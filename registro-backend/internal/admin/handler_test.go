package admin

import (
	"bytes"
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

func setupTestHandler() (*Handler, *MockRepository) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	handler := NewHandler(service)
	return handler, mockRepo
}

func TestHandler_ListAuditLogs(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful list logs",
			setupMock: func(m *MockRepository) {
				m.On("ListAuditLogs", mock.Anything, mock.MatchedBy(func(req *AuditLogListRequest) bool {
					return req.PageSize == 20 || req.PageSize == 0
				}), 0).Return([]ActivityLogEntry{
					{
						ID:         "log-1",
						AdminID:    "admin-1",
						AdminName:  "Admin One",
						ActionType: "create",
						Target:     "school",
						Details:    "Created school Foo",
						CreatedAt:  time.Now(),
					},
				}, int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response AuditLogListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), response.Total)
				assert.Len(t, response.Items, 1)
				assert.Equal(t, "log-1", response.Items[0].ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo := setupTestHandler()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/admin/audit-logs", nil)

			handler.ListAuditLogs(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetDashboardStats(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		setupMock      func(*MockRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful dashboard stats for superadmin",
			setupContext: func(c *gin.Context) {
				c.Set("role", "superadmin")
			},
			setupMock: func(m *MockRepository) {
				m.On("CountSchools", mock.Anything, (*string)(nil)).Return(int64(5), nil)
				m.On("CountUsers", mock.Anything, (*string)(nil)).Return(int64(100), nil)
				m.On("CountUsersByRole", mock.Anything, "teacher", (*string)(nil)).Return(int64(20), nil)
				m.On("CountUsersByRole", mock.Anything, "student", (*string)(nil)).Return(int64(70), nil)
				m.On("CountDocuments", mock.Anything, (*string)(nil)).Return(int64(150), nil)
				m.On("CountPendingDocuments", mock.Anything, (*string)(nil)).Return(int64(10), nil)
				m.On("CountCommunications", mock.Anything, (*string)(nil)).Return(int64(25), nil)
				m.On("CountActiveUsers24h", mock.Anything, (*string)(nil)).Return(int64(50), nil)
				m.On("GetRecentEvents", mock.Anything, 10, (*string)(nil)).Return([]RecentEvent{}, nil)
				m.On("GetSystemHealth", mock.Anything).Return(&SystemHealthStatus{OverallStatus: "healthy"}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var stats DashboardStatsResponse
				err := json.Unmarshal(w.Body.Bytes(), &stats)
				assert.NoError(t, err)
				assert.Equal(t, int64(5), stats.TotalSchools)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo := setupTestHandler()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/admin/dashboard/stats", nil)
			if tt.setupContext != nil {
				tt.setupContext(c)
			}

			handler.GetDashboardStats(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_ListSchools(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		setupContext   func(*gin.Context)
		setupMock      func(*MockRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "successful list schools",
			queryParams: "?page=1&page_size=10",
			setupContext: func(c *gin.Context) {
				c.Set("filtered_school_id", "")
			},
			setupMock: func(m *MockRepository) {
				m.On("ListSchools", mock.Anything, mock.AnythingOfType("*admin.SchoolListRequest"), 0, (*string)(nil)).
					Return([]SchoolResponse{{ID: "school-1", Name: "Test School"}}, int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response SchoolListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Items, 1)
			},
		},
		{
			name: "repository error",
			setupContext: func(c *gin.Context) {
				c.Set("filtered_school_id", "")
			},
			setupMock: func(m *MockRepository) {
				m.On("ListSchools", mock.Anything, mock.AnythingOfType("*admin.SchoolListRequest"), 0, (*string)(nil)).
					Return(nil, int64(0), errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo := setupTestHandler()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			url := "/admin/schools"
			if tt.queryParams != "" {
				url += tt.queryParams
			}
			c.Request = httptest.NewRequest("GET", url, nil)
			if tt.setupContext != nil {
				tt.setupContext(c)
			}

			handler.ListSchools(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_CreateSchool(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupContext   func(*gin.Context)
		setupMock      func(*MockRepository)
		expectedStatus int
	}{
		{
			name: "successful creation",
			requestBody: CreateSchoolRequest{
				Name:     "New School",
				Code:     "NEW01",
				Address:  "123 Main St",
				City:     "Rome",
				Province: "RM",
				ZipCode:  "00100",
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", "1")
				c.Set("role", "superadmin")
			},
			setupMock: func(m *MockRepository) {
				m.On("SchoolCodeExists", mock.Anything, "NEW01").Return(false, nil)
				m.On("CreateSchool", mock.Anything, mock.AnythingOfType("*admin.CreateSchoolRequest")).
					Return(&SchoolResponse{ID: "school-new", Name: "New School"}, nil)
				m.On("LogAdminAction", mock.Anything, "1", "create", "school", mock.Anything, (*string)(nil), mock.AnythingOfType("string")).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockRepository) {},
		},
		{
			name: "repository error",
			requestBody: CreateSchoolRequest{
				Name:     "New School",
				Code:     "NEW01",
				Address:  "123 Main St",
				City:     "Rome",
				Province: "RM",
				ZipCode:  "00100",
			},
			setupContext: func(c *gin.Context) {
				c.Set("userID", "1")
				c.Set("role", "superadmin")
			},
			setupMock: func(m *MockRepository) {
				m.On("SchoolCodeExists", mock.Anything, "NEW01").Return(false, nil)
				m.On("CreateSchool", mock.Anything, mock.AnythingOfType("*admin.CreateSchoolRequest")).
					Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo := setupTestHandler()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			c.Request = httptest.NewRequest("POST", "/admin/schools", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")
			if tt.setupContext != nil {
				tt.setupContext(c)
			}

			handler.CreateSchool(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetSchool(t *testing.T) {
	tests := []struct {
		name           string
		schoolID       string
		setupContext   func(*gin.Context)
		setupMock      func(*MockRepository)
		expectedStatus int
	}{
		{
			name:     "successful get school",
			schoolID: "school-123",
			setupContext: func(c *gin.Context) {
				c.Set("role", "superadmin")
			},
			setupMock: func(m *MockRepository) {
				m.On("GetSchool", mock.Anything, "school-123").
					Return(&SchoolResponse{ID: "school-123", Name: "Test School"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "school not found",
			schoolID: "school-999",
			setupContext: func(c *gin.Context) {
				c.Set("role", "superadmin")
			},
			setupMock: func(m *MockRepository) {
				m.On("GetSchool", mock.Anything, "school-999").
					Return(nil, ErrSchoolNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo := setupTestHandler()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/admin/schools/"+tt.schoolID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.schoolID}}
			if tt.setupContext != nil {
				tt.setupContext(c)
			}

			handler.GetSchool(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteSchool(t *testing.T) {
	tests := []struct {
		name           string
		schoolID       string
		setupContext   func(*gin.Context)
		setupMock      func(*MockRepository)
		expectedStatus int
	}{
		{
			name:     "successful deletion",
			schoolID: "school-123",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", "1")
				c.Set("role", "superadmin")
			},
			setupMock: func(m *MockRepository) {
				m.On("DeleteSchool", mock.Anything, "school-123").Return(nil)
				m.On("LogAdminAction", mock.Anything, "1", "delete", "school", mock.Anything, mock.Anything, mock.AnythingOfType("string")).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "repository error",
			schoolID: "school-123",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", "1")
				c.Set("role", "superadmin")
			},
			setupMock: func(m *MockRepository) {
				m.On("DeleteSchool", mock.Anything, "school-123").Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo := setupTestHandler()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("DELETE", "/admin/schools/"+tt.schoolID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.schoolID}}
			if tt.setupContext != nil {
				tt.setupContext(c)
			}

			handler.DeleteSchool(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_ListAdminUsers(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		setupMock      func(*MockRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "successful list",
			queryParams: "?page=1&page_size=10",
			setupMock: func(m *MockRepository) {
				m.On("ListAdminUsers", mock.Anything, 0, 10, (*string)(nil)).
					Return([]AdminUserResponse{{ID: "admin-1", Email: "admin@test.com"}}, int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response AdminUserListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Items, 1)
			},
		},
		{
			name: "repository error",
			setupMock: func(m *MockRepository) {
				m.On("ListAdminUsers", mock.Anything, 0, 20, (*string)(nil)).
					Return(nil, int64(0), errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockRepo := setupTestHandler()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			url := "/admin/users/admins"
			if tt.queryParams != "" {
				url += tt.queryParams
			}
			c.Request = httptest.NewRequest("GET", url, nil)

			handler.ListAdminUsers(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
