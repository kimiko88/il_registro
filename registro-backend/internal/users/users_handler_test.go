package users

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func setupTestUsersRouter(svc *Service, userID, role, schoolID string) *gin.Engine {
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
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestUsersHandler_AllEndpoints(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	schID := "sch-1"
	rAdmin := setupTestUsersRouter(svc, "admin-1", "admin", "sch-1")

	// 1. POST /users - Create User
	repo.On("Create", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil).Once()
	repo.On("LogAudit", mock.Anything, mock.Anything).Return(nil).Maybe()
	repo.On("RevokeAllUserTokens", mock.Anything, mock.Anything).Return(nil).Maybe()

	createPayload := `{"email":"mario.rossi@scuola.it","password":"Password123!","first_name":"Mario","last_name":"Rossi","role":"teacher","school_id":"sch-1"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(createPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. GET /users - List Users
	repo.On("List", mock.Anything, mock.AnythingOfType("users.UserFilter")).Return([]User{{ID: "u1", FirstName: "Mario"}}, 1, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. GET /users/:id - Get User
	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID, FirstName: "Mario", Role: "teacher"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/u1", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. PATCH /users/:id - Update User
	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID, FirstName: "Mario", Role: "teacher"}, nil).Once()
	repo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil).Once()
	updatePayload := `{"first_name":"Mario Giuseppe"}`
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/users/u1", bytes.NewBufferString(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. POST /users/:id/change-password
	hashed, _ := bcrypt.GenerateFromPassword([]byte("OldPassword123!"), bcrypt.DefaultCost)
	repo.On("GetByID", mock.Anything, "admin-1").Return(&User{ID: "admin-1", PasswordHash: string(hashed)}, nil).Once()
	repo.On("GetPasswordHistory", mock.Anything, "admin-1").Return([]string{}, nil).Once()
	repo.On("ChangePasswordTx", mock.Anything, "admin-1", mock.Anything).Return(nil).Once()

	pwdPayload := `{"current_password":"OldPassword123!","new_password":"NewPassword123!"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/admin-1/change-password", bytes.NewBufferString(pwdPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. POST /users/:id/reset-password
	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID}, nil).Once()
	repo.On("GetPasswordHistory", mock.Anything, "u1").Return([]string{}, nil).Once()
	repo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil).Once()
	repo.On("AddPasswordHistory", mock.Anything, "u1", mock.Anything).Return(nil).Once()

	resetPayload := `{"new_password":"ResetPassword123!"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/u1/reset-password", bytes.NewBufferString(resetPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. PATCH /users/:id/roles
	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID, Role: "teacher"}, nil).Once()
	repo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil).Once()

	rolesPayload := `{"role":"coordinator"}`
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/users/u1/roles", bytes.NewBufferString(rolesPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 8. GET /users/:id/audit-log
	repo.On("GetAuditLogs", mock.Anything, "u1", 50, 0).Return([]AuditLog{{ID: "a1", Action: "CREATE"}}, 1, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/u1/audit-log", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 9. POST /users/:id/gdpr-export
	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID, Email: "mario@scuola.it"}, nil).Once()
	repo.On("GetAuditLogs", mock.Anything, "u1", mock.Anything, mock.Anything).Return([]AuditLog{}, 0, nil).Once()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/u1/gdpr-export", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 10. PATCH /users/:id/disable-mfa
	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID}, nil).Once()
	repo.On("ClearTempMFASecret", mock.Anything, "u1").Return(nil).Once()
	repo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil).Once()
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/users/u1/disable-mfa", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 11. POST /users/gdpr/retention
	repo.On("ApplyDataRetention", mock.Anything, &schID, mock.AnythingOfType("time.Time")).Return(10, nil).Once()
	retentionPayload := `{"retention_years": 5}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/gdpr/retention", bytes.NewBufferString(retentionPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 12. DELETE /users/:id - Soft Delete
	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID}, nil).Once()
	repo.On("Delete", mock.Anything, "u1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/users/u1", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// 13. POST /users/:id/restore
	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID}, nil).Once()
	repo.On("Restore", mock.Anything, "u1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/u1/restore", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 14. GET /users/me/children & switch child
	rParent := setupTestUsersRouter(svc, "parent-1", "parent", "sch-1")
	repo.On("GetChildren", mock.Anything, "parent-1").Return([]StudentChild{{ID: "s1", FirstName: "Luca"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/me/children", nil)
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("IsGuardian", mock.Anything, "parent-1", "s1").Return(true, nil).Once()
	repo.On("GetByID", mock.Anything, "s1").Return(&User{ID: "s1", SchoolID: &schID, FirstName: "Luca", Role: "student"}, nil).Once()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/me/switch-child/s1", nil)
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 15. DELETE /users/:id/gdpr-delete
	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID}, nil).Once()
	repo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/users/u1/gdpr-delete", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 16. Assignments: GET & POST & DELETE /users/:id/assignments
	repo.On("GetAssignments", mock.Anything, "u1").Return([]UserAssignment{{ID: "as-1", UserID: "u1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/u1/assignments", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID}, nil).Once()
	repo.On("CreateAssignment", mock.Anything, mock.AnythingOfType("*users.UserAssignment")).Return(nil).Once()
	asPayload := `{"assignment_type":"coordinatore_classe","title":"Coordinatore 3A"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/u1/assignments", bytes.NewBufferString(asPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID}, nil).Once()
	repo.On("GetAssignments", mock.Anything, "u1").Return([]UserAssignment{{ID: "as-1", AssignmentType: "coordinatore_classe"}}, nil).Once()
	repo.On("DeleteAssignment", mock.Anything, "as-1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/users/u1/assignments/as-1", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 17. Coordinated Classes: PUT /users/:id/coordinated-classes
	repo.On("GetByID", mock.Anything, "u1").Return(&User{ID: "u1", SchoolID: &schID}, nil).Once()
	repo.On("SetCoordinatedClasses", mock.Anything, "sch-1", "u1", []string{"c1", "c2"}).Return(nil).Once()
	coordPayload := `{"class_ids":["c1","c2"]}`
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/users/u1/coordinated-classes", bytes.NewBufferString(coordPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 18. Guardians: GET & POST & DELETE /users/:id/guardians
	repo.On("GetStudentProfile", mock.Anything, "s1").Return("prof-s1", nil).Times(3)
	repo.On("GetGuardians", mock.Anything, "prof-s1").Return([]GuardianInfo{{ParentUserID: "p1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/s1/guardians", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetParentProfile", mock.Anything, "p1").Return("prof-p1", nil).Twice()
	repo.On("AddGuardian", mock.Anything, "prof-s1", "prof-p1", "Father").Return(nil).Once()
	guardianPayload := `{"parent_user_id":"p1","relationship_type":"Father"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/s1/guardians", bytes.NewBufferString(guardianPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("RemoveGuardian", mock.Anything, "prof-s1", "prof-p1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/users/s1/guardians/p1", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 19. GET /users/students/:id/fascicolo
	repo.On("GetByID", mock.Anything, "s1").Return(&User{ID: "s1", SchoolID: &schID, IsActive: true}, nil).Once()
	repo.On("GetStudentProfile", mock.Anything, "s1").Return("prof-s1", nil).Once()
	repo.On("GetGuardians", mock.Anything, "prof-s1").Return([]GuardianInfo{}, nil).Once()
	repo.On("GetFascicoloSummary", mock.Anything, "s1", true).Return(map[string]interface{}{"status": "Active"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/students/s1/fascicolo", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 20. POST /users/bulk-delete
	repo.On("ListByIDs", mock.Anything, []string{"u2", "u3"}).Return([]User{{ID: "u2", SchoolID: &schID}, {ID: "u3", SchoolID: &schID}}, nil).Once()
	repo.On("BulkDelete", mock.Anything, []string{"u2", "u3"}).Return(2, nil).Once()
	bulkDelPayload := `{"user_ids":["u2","u3"]}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/bulk-delete", bytes.NewBufferString(bulkDelPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 21. GET /users/search
	repo.On("List", mock.Anything, mock.AnythingOfType("users.UserFilter")).Return([]User{{ID: "u1"}}, 1, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/search?role=teacher&school_id=sch-1", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 22. Error validations (Invalid JSON body)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/users/u1", bytes.NewBufferString(`{invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	repo.AssertExpectations(t)
}

func TestUsersHandler_Unauthenticated(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	rUnauth := setupTestUsersRouter(svc, "", "", "")

	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/v1/users"},
		{http.MethodGet, "/api/v1/users"},
		{http.MethodGet, "/api/v1/users/u1"},
		{http.MethodPatch, "/api/v1/users/u1"},
		{http.MethodDelete, "/api/v1/users/u1"},
		{http.MethodPost, "/api/v1/users/bulk-import"},
		{http.MethodPost, "/api/v1/users/bulk-delete"},
		{http.MethodGet, "/api/v1/users/me/children"},
		{http.MethodPost, "/api/v1/users/me/switch-child/s1"},
		{http.MethodGet, "/api/v1/users/s1/guardians"},
		{http.MethodPost, "/api/v1/users/s1/guardians"},
		{http.MethodDelete, "/api/v1/users/s1/guardians/p1"},
		{http.MethodGet, "/api/v1/users/students/s1/fascicolo"},
	}

	for _, ep := range endpoints {
		req, _ := http.NewRequest(ep.method, ep.path, nil)
		w := httptest.NewRecorder()
		rUnauth.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected unauthorized for %s %s", ep.method, ep.path)
	}
}

func TestUsersHandler_Errors(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	rAdmin := setupTestUsersRouter(svc, "admin-1", "admin", "sch-1")

	// 1. Get not found
	repo.On("GetByID", mock.Anything, "not-found").Return(nil, ErrUserNotFound).Once()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/not-found", nil)
	w := httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 2. Change password of another user (forbidden)
	pwdPayload := `{"current_password":"OldPassword123!","new_password":"NewPassword123!"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/other-user/change-password", bytes.NewBufferString(pwdPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 3. Reset password with invalid body
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/u1/reset-password", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 4. Assign roles with invalid body
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/users/u1/roles", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 5. Add assignment invalid body
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/u1/assignments", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 6. Set coordinated classes invalid body
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/users/u1/coordinated-classes", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 7. Create with invalid role
	invalidRolePayload := `{"email":"test@test.it","password":"Password123!","first_name":"T","last_name":"T","role":"astronaut","school_id":"sch-1"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(invalidRolePayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 8. Delete user not found
	repo.On("GetByID", mock.Anything, "not-found-del").Return(nil, ErrUserNotFound).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/users/not-found-del", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 9. Restore user not found
	repo.On("GetByID", mock.Anything, "not-found-res").Return(nil, ErrUserNotFound).Once()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/not-found-res/restore", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 10. Disable MFA by non-admin (forbidden)
	rTeacher := setupTestUsersRouter(svc, "t-1", "teacher", "sch-1")
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/users/u1/disable-mfa", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 11. Add guardian by unauthorized role (e.g., student)
	rStudent := setupTestUsersRouter(svc, "s-1", "student", "sch-1")
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/s1/guardians", bytes.NewBufferString(`{"parent_user_id":"p1"}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
