package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestHandler() (*Handler, *MockRepository) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	handler := NewHandler(service)
	return handler, mockRepo
}

func TestHandler_ListAuditLogs(t *testing.T) {
	gin.SetMode(gin.TestMode)

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
					// Check default pagination
					return req.PageSize == 20 || req.PageSize == 0 // Service defaults it
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
