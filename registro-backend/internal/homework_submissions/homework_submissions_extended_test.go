package homework_submissions_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/homework_submissions"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) UpsertSubmission(s *homework_submissions.HomeworkSubmission) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *mockRepo) GetSubmission(id string) (*homework_submissions.HomeworkSubmission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*homework_submissions.HomeworkSubmission), args.Error(1)
}

func (m *mockRepo) GetSubmissionByHomeworkAndStudent(homeworkID, studentID string) (*homework_submissions.HomeworkSubmission, error) {
	args := m.Called(homeworkID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*homework_submissions.HomeworkSubmission), args.Error(1)
}

func (m *mockRepo) ListSubmissionsByHomework(homeworkID string) ([]homework_submissions.HomeworkSubmission, error) {
	args := m.Called(homeworkID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]homework_submissions.HomeworkSubmission), args.Error(1)
}

func (m *mockRepo) ListSubmissionsByStudent(studentID string) ([]homework_submissions.HomeworkSubmission, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]homework_submissions.HomeworkSubmission), args.Error(1)
}

func (m *mockRepo) SaveGrading(id string, req homework_submissions.GradeSubmissionRequest, teacherID string) error {
	args := m.Called(id, req, teacherID)
	return args.Error(0)
}

func (m *mockRepo) SetClassroomSync(id, submissionID string) error {
	args := m.Called(id, submissionID)
	return args.Error(0)
}

func setupTestRouter(svc homework_submissions.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "sch-test-123")
		c.Set("user_id", "student-test-456")
		c.Next()
	})
	h := homework_submissions.NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestHomeworkSubmissionsService_Full(t *testing.T) {
	repo := new(mockRepo)
	svc := homework_submissions.NewService(repo)

	// SubmitHomework validations
	assert.Error(t, svc.SubmitHomework(nil))
	assert.Error(t, svc.SubmitHomework(&homework_submissions.HomeworkSubmission{HomeworkID: ""}))
	assert.Error(t, svc.SubmitHomework(&homework_submissions.HomeworkSubmission{HomeworkID: "h1", StudentID: ""}))
	assert.Error(t, svc.SubmitHomework(&homework_submissions.HomeworkSubmission{HomeworkID: "h1", StudentID: "s1", FileURL: ""}))

	repo.On("UpsertSubmission", mock.Anything).Return(nil).Once()
	sub := &homework_submissions.HomeworkSubmission{HomeworkID: "h1", StudentID: "s1", FileURL: "https://example.com/file.pdf"}
	assert.NoError(t, svc.SubmitHomework(sub))
	assert.Equal(t, "pdf", sub.FileType)

	// GetSubmission
	_, err := svc.GetSubmission("")
	assert.Error(t, err)

	repo.On("GetSubmission", "sub-1").Return(&homework_submissions.HomeworkSubmission{ID: "sub-1"}, nil).Once()
	s1, err := svc.GetSubmission("sub-1")
	assert.NoError(t, err)
	assert.Equal(t, "sub-1", s1.ID)

	// GetStudentSubmission
	_, err = svc.GetStudentSubmission("", "s1")
	assert.Error(t, err)
	_, err = svc.GetStudentSubmission("h1", "")
	assert.Error(t, err)

	repo.On("GetSubmissionByHomeworkAndStudent", "h1", "s1").Return(&homework_submissions.HomeworkSubmission{ID: "sub-1"}, nil).Once()
	s2, err := svc.GetStudentSubmission("h1", "s1")
	assert.NoError(t, err)
	assert.Equal(t, "sub-1", s2.ID)

	// ListHomeworkSubmissions
	_, err = svc.ListHomeworkSubmissions("")
	assert.Error(t, err)

	repo.On("ListSubmissionsByHomework", "h1").Return([]homework_submissions.HomeworkSubmission{{ID: "sub-1"}}, nil).Once()
	list, err := svc.ListHomeworkSubmissions("h1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// ListStudentSubmissions
	_, err = svc.ListStudentSubmissions("")
	assert.Error(t, err)

	repo.On("ListSubmissionsByStudent", "s1").Return([]homework_submissions.HomeworkSubmission{{ID: "sub-1"}}, nil).Once()
	studList, err := svc.ListStudentSubmissions("s1")
	assert.NoError(t, err)
	assert.Len(t, studList, 1)

	// GradeSubmission
	assert.Error(t, svc.GradeSubmission("", homework_submissions.GradeSubmissionRequest{}, "teach-1"))
	assert.Error(t, svc.GradeSubmission("sub-1", homework_submissions.GradeSubmissionRequest{}, ""))

	gradeTooHigh := 11.0
	assert.Error(t, svc.GradeSubmission("sub-1", homework_submissions.GradeSubmissionRequest{TeacherGrade: &gradeTooHigh}, "teach-1"))

	gradeTooLow := 0.5
	assert.Error(t, svc.GradeSubmission("sub-1", homework_submissions.GradeSubmissionRequest{TeacherGrade: &gradeTooLow}, "teach-1"))

	validGrade := 8.5
	repo.On("SaveGrading", "sub-1", mock.Anything, "teach-1").Return(nil).Once()
	assert.NoError(t, svc.GradeSubmission("sub-1", homework_submissions.GradeSubmissionRequest{TeacherGrade: &validGrade}, "teach-1"))

	// SyncGoogleClassroom
	assert.Error(t, svc.SyncGoogleClassroom("", "csub-1"))
	assert.Error(t, svc.SyncGoogleClassroom("sub-1", ""))

	repo.On("SetClassroomSync", "sub-1", "csub-1").Return(nil).Once()
	assert.NoError(t, svc.SyncGoogleClassroom("sub-1", "csub-1"))

	repo.AssertExpectations(t)
}

func TestHomeworkSubmissionsHandler_Endpoints(t *testing.T) {
	repo := new(mockRepo)
	svc := homework_submissions.NewService(repo)
	r := setupTestRouter(svc)

	// 1. POST /homework-submissions
	repo.On("UpsertSubmission", mock.Anything).Return(nil).Once()
	payload := `{"homework_id":"hw-1","file_url":"https://storage/hw1.pdf"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/homework-submissions", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/homework-submissions", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 2. GET /homework-submissions/:id
	repo.On("GetSubmission", "sub-1").Return(&homework_submissions.HomeworkSubmission{ID: "sub-1"}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/homework-submissions/sub-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("GetSubmission", "missing").Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/homework-submissions/missing", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 3. GET /homework-submissions/by-homework
	repo.On("ListSubmissionsByHomework", "hw-1").Return([]homework_submissions.HomeworkSubmission{{ID: "s1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/homework-submissions/by-homework?homework_id=hw-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListSubmissionsByHomework", "hw-err").Return(nil, errors.New("db error")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/homework-submissions/by-homework?homework_id=hw-err", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 4. GET /homework-submissions/my-submissions
	repo.On("ListSubmissionsByStudent", "student-test-456").Return([]homework_submissions.HomeworkSubmission{{ID: "s1"}}, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/homework-submissions/my-submissions", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.On("ListSubmissionsByStudent", "student-test-456").Return(nil, errors.New("db error")).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/homework-submissions/my-submissions", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 5. POST /homework-submissions/:id/grade
	repo.On("SaveGrading", "sub-1", mock.Anything, "student-test-456").Return(nil).Once()
	gradePayload := `{"teacher_grade":9.0,"teacher_feedback":"Ottimo lavoro"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/homework-submissions/sub-1/grade", bytes.NewBufferString(gradePayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/homework-submissions/sub-1/grade", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 6. POST /homework-submissions/:id/sync-classroom
	repo.On("SetClassroomSync", "sub-1", "gclass-123").Return(nil).Once()
	syncPayload := `{"google_classroom_submission_id":"gclass-123"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/homework-submissions/sub-1/sync-classroom", bytes.NewBufferString(syncPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	repo.AssertExpectations(t)
}
