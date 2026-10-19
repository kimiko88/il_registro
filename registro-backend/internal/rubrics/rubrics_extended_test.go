package rubrics

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRubricsRouter(svc *Service, userID, role, schoolID string) *gin.Engine {
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

func TestRubricsHandlers_AllRoutes(t *testing.T) {
	mockRepo := new(MockRubricsRepo)
	mockUsers := new(MockUsersRepo)
	svc := NewService(mockRepo, mockUsers)

	rTeacher := setupTestRubricsRouter(svc, "teacher-1", "teacher", "sch-1")

	// 1. POST /rubrics - Create Rubric
	mockRepo.On("CreateRubric", mock.Anything, mock.AnythingOfType("*rubrics.Rubric")).Return(nil).Once()
	createPayload := `{"title":"Rubrica Storia","description":"Esposizione orale","subject_id":"sub-1","criteria":[{"name":"Chiarezza","max_score":10,"levels":[{"score":10,"label":"Ottimo"}]}]}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/rubrics", bytes.NewBufferString(createPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. GET /rubrics - List Rubrics
	mockRepo.On("ListRubrics", mock.Anything, "sch-1", "teacher-1").Return([]*Rubric{{ID: "rub-1", Title: "Rubrica Storia"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/rubrics", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. GET /rubrics/:id - Get Rubric
	mockRepo.On("GetRubricByID", mock.Anything, "rub-1").Return(&Rubric{ID: "rub-1", Title: "Rubrica Storia"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/rubrics/rub-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. PUT /rubrics/:id - Update Rubric
	mockRepo.On("GetRubricByID", mock.Anything, "rub-1").Return(&Rubric{ID: "rub-1", SchoolID: "sch-1", TeacherID: "teacher-1"}, nil).Once()
	mockRepo.On("UpdateRubric", mock.Anything, mock.AnythingOfType("*rubrics.Rubric")).Return(nil).Once()
	updatePayload := `{"subject_id":"sub-1","title":"Rubrica Storia Modificata"}`
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/rubrics/rub-1", bytes.NewBufferString(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. POST /rubrics/:id/assessments - Assess Student
	mockRepo.On("CreateAssessment", mock.Anything, mock.AnythingOfType("*rubrics.RubricAssessment")).Return(nil).Once()
	assessPayload := `{"student_id":"std-1","class_id":"cls-1","date":"2026-05-15","scores":[{"criterion_id":"c1","level_id":"l1","score":9.5}],"notes":"Ottima prova"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/rubrics/rub-1/assessments", bytes.NewBufferString(assessPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 6. GET /rubrics/assessments/student/:studentID
	mockRepo.On("ListAssessmentsByStudent", mock.Anything, "std-1").Return([]*RubricAssessment{{ID: "ass-1", StudentID: "std-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/rubrics/assessments/student/std-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. GET /rubrics/assessments/class/:classID
	mockRepo.On("ListAssessmentsByClass", mock.Anything, "cls-1").Return([]*RubricAssessment{{ID: "ass-1", ClassID: "cls-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/rubrics/assessments/class/cls-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 8. DELETE /rubrics/:id - Delete Rubric
	mockRepo.On("GetRubricByID", mock.Anything, "rub-1").Return(&Rubric{ID: "rub-1", SchoolID: "sch-1", TeacherID: "teacher-1"}, nil).Once()
	mockRepo.On("DeleteRubric", mock.Anything, "rub-1").Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/rubrics/rub-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 9. Unauthorized / Forbidden tests
	rStudent := setupTestRubricsRouter(svc, "std-2", "student", "sch-1")
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/rubrics", bytes.NewBufferString(createPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Student looking at another student's assessments
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/rubrics/assessments/student/std-1", nil)
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Parent is guardian -> OK
	rParent := setupTestRubricsRouter(svc, "parent-1", "parent", "sch-1")
	mockUsers.On("IsGuardian", mock.Anything, "parent-1", "std-1").Return(true, nil).Once()
	mockRepo.On("ListAssessmentsByStudent", mock.Anything, "std-1").Return([]*RubricAssessment{{ID: "ass-1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/rubrics/assessments/student/std-1", nil)
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// GetByID not found
	mockRepo.On("GetRubricByID", mock.Anything, "missing").Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/rubrics/missing", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Delete forbidden by other teacher
	mockRepo.On("GetRubricByID", mock.Anything, "rub-2").Return(&Rubric{ID: "rub-2", TeacherID: "other-teacher", SchoolID: "sch-1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/rubrics/rub-2", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	mockRepo.AssertExpectations(t)
	mockUsers.AssertExpectations(t)
}

func TestRubricsService_ExtendedEdgeCases(t *testing.T) {
	mockRepo := new(MockRubricsRepo)
	svc := NewService(mockRepo)
	ctx := context.Background()

	// 1. GetRubric
	_, err := svc.GetRubric(ctx, "")
	assert.ErrorContains(t, err, "required")

	// 2. UpdateRubric edge cases
	_, err = svc.UpdateRubric(ctx, "usr-1", "teacher", "sch-1", "", CreateRubricRequest{})
	assert.ErrorContains(t, err, "required")

	mockRepo.On("GetRubricByID", ctx, "missing").Return(nil, errors.New("not found")).Once()
	_, err = svc.UpdateRubric(ctx, "usr-1", "teacher", "sch-1", "missing", CreateRubricRequest{})
	assert.Error(t, err)

	mockRepo.On("GetRubricByID", ctx, "nil-rub").Return(nil, nil).Once()
	_, err = svc.UpdateRubric(ctx, "usr-1", "teacher", "sch-1", "nil-rub", CreateRubricRequest{})
	assert.Equal(t, ErrRubricNotFound, err)

	// Other teacher editing
	mockRepo.On("GetRubricByID", ctx, "rub-1").Return(&Rubric{ID: "rub-1", TeacherID: "teacher-1", SchoolID: "sch-1"}, nil).Once()
	_, err = svc.UpdateRubric(ctx, "other-teacher", "teacher", "sch-1", "rub-1", CreateRubricRequest{})
	assert.Equal(t, ErrUnauthorized, err)

	// Admin of different school
	mockRepo.On("GetRubricByID", ctx, "rub-1").Return(&Rubric{ID: "rub-1", TeacherID: "teacher-1", SchoolID: "sch-1"}, nil).Once()
	_, err = svc.UpdateRubric(ctx, "admin-2", "admin", "sch-other", "rub-1", CreateRubricRequest{})
	assert.Equal(t, ErrUnauthorized, err)

	// 3. DeleteRubric edge cases
	err = svc.DeleteRubric(ctx, "usr-1", "teacher", "sch-1", "")
	assert.ErrorContains(t, err, "required")

	mockRepo.On("GetRubricByID", ctx, "missing-del").Return(nil, errors.New("not found")).Once()
	err = svc.DeleteRubric(ctx, "usr-1", "teacher", "sch-1", "missing-del")
	assert.Error(t, err)

	mockRepo.On("GetRubricByID", ctx, "nil-del").Return(nil, nil).Once()
	err = svc.DeleteRubric(ctx, "usr-1", "teacher", "sch-1", "nil-del")
	assert.Equal(t, ErrRubricNotFound, err)

	// 4. AssessStudent edge cases
	_, err = svc.AssessStudent(ctx, "teacher-1", "", CreateAssessmentRequest{})
	assert.ErrorContains(t, err, "rubric_id required")

	_, err = svc.AssessStudent(ctx, "teacher-1", "rub-1", CreateAssessmentRequest{StudentID: ""})
	assert.ErrorContains(t, err, "student_id required")

	_, err = svc.AssessStudent(ctx, "teacher-1", "rub-1", CreateAssessmentRequest{
		StudentID: "std-1",
		Scores:    []CriterionScore{{Score: -5}},
	})
	assert.ErrorContains(t, err, "cannot be negative")

	// 5. ListAssessmentsByStudent & ListAssessmentsByClass edge cases
	_, err = svc.ListAssessmentsByStudent(ctx, "")
	assert.ErrorContains(t, err, "student_id required")

	_, err = svc.ListAssessmentsByClass(ctx, "")
	assert.ErrorContains(t, err, "class_id required")

	// Success cases
	now := time.Now()
	mockRepo.On("ListAssessmentsByStudent", ctx, "std-1").Return([]*RubricAssessment{{ID: "ass-1", Date: now}}, nil).Once()
	assList, err := svc.ListAssessmentsByStudent(ctx, "std-1")
	assert.NoError(t, err)
	assert.Len(t, assList, 1)

	mockRepo.On("ListAssessmentsByClass", ctx, "cls-1").Return([]*RubricAssessment{{ID: "ass-2", Date: now}}, nil).Once()
	clsList, err := svc.ListAssessmentsByClass(ctx, "cls-1")
	assert.NoError(t, err)
	assert.Len(t, clsList, 1)

	mockRepo.AssertExpectations(t)
}
