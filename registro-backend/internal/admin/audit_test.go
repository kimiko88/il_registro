package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListAdminUsers_NonSuperAdmin_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupTestHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/admin/users/admins", nil)
	c.Request = req
	c.Set("user_id", "admin-user-1")
	c.Set("role", "admin") // regular admin, not superadmin

	handler.ListAdminUsers(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateSchoolSetting_ExpandedAllowlist(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	mockRepo.On("UpdateSetting", mock.Anything, "school-1", "enable_pcto", "true").Return(nil).Once()
	mockRepo.On("UpdateSetting", mock.Anything, "school-1", "timetable_visible_to_parents", "true").Return(nil).Once()

	err := svc.UpdateSchoolSetting(context.Background(), "school-1", "enable_pcto", "true")
	assert.NoError(t, err)

	err = svc.UpdateSchoolSetting(context.Background(), "school-1", "timetable_visible_to_parents", "true")
	assert.NoError(t, err)

	err = svc.UpdateSchoolSetting(context.Background(), "school-1", "unauthorized_key", "value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or unauthorized setting key")

	mockRepo.AssertExpectations(t)
}

func TestUpdateSchool_PassesSchoolFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	handler := NewHandler(svc)

	schoolID := "school-1"
	schoolFilter := "school-1"
	mockRepo.On("UpdateSchool", mock.Anything, schoolID, mock.Anything).Return(&SchoolResponse{
		ID:   schoolID,
		Name: "Updated School",
	}, nil).Once()
	mockRepo.On("LogAdminAction", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("PUT", "/admin/schools/"+schoolID, strings.NewReader(`{"name":"Updated School"}`))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: schoolID}}
	c.Set("user_id", "admin-1")
	c.Set("role", "admin")
	c.Set("school_id", schoolFilter)

	handler.UpdateSchool(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
