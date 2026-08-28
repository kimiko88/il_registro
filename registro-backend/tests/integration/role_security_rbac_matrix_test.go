package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/attendance"
	"registro-backend/internal/classes"
	"registro-backend/internal/grades"
	"registro-backend/internal/scrutiny"
	"registro-backend/internal/substitutions"
	"registro-backend/internal/users"
	"registro-backend/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSubRepo struct {
	mock.Mock
}

func (m *mockSubRepo) Create(ctx context.Context, sub *substitutions.Substitution) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}
func (m *mockSubRepo) GetByID(ctx context.Context, id string) (*substitutions.Substitution, error) {
	return nil, nil
}
func (m *mockSubRepo) ListBySchool(ctx context.Context, schoolID, date string) ([]*substitutions.Substitution, error) {
	return nil, nil
}
func (m *mockSubRepo) ListByTeacher(ctx context.Context, teacherID, date string) ([]*substitutions.Substitution, error) {
	return nil, nil
}
func (m *mockSubRepo) AssignSubstitute(ctx context.Context, id string, substituteTeacherID string, notes string) error {
	return nil
}
func (m *mockSubRepo) ConfirmSubstitution(ctx context.Context, id string, substituteTeacherID string) error {
	return nil
}
func (m *mockSubRepo) SignRegister(ctx context.Context, id string, sigHash string, notes string) error {
	return nil
}
func (m *mockSubRepo) GetAvailableTeachers(ctx context.Context, schoolID string) ([]substitutions.TeacherCandidate, error) {
	return nil, nil
}
func (m *mockSubRepo) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *mockSubRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	return true, nil
}
func (m *mockSubRepo) IsTeacherAssignedToSubject(ctx context.Context, teacherID, subjectID string) (bool, error) {
	return true, nil
}
func (m *mockSubRepo) GetWeeklySubstitutionCount(ctx context.Context, teacherID string) (int, error) {
	return 0, nil
}

func TestRoleSecurityRBACMatrix(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const schoolID = "school-rbac-01"
	const classID = "class-rbac-1a"
	const studentID = "student-rbac-01"
	const teacherID = "teacher-creator-01"
	const otherTeacherID = "teacher-other-02"
	todayStr := time.Now().Format("2006-01-02")
	schoolIDPtr := schoolID

	setAuthHeader := func(userID, role, sID string) gin.HandlerFunc {
		return func(c *gin.Context) {
			if role != "anonymous" {
				c.Set("user_id", userID)
				c.Set("role", role)
				c.Set("school_id", sID)
			}
			c.Next()
		}
	}

	allRoles := []string{
		"superadmin", "admin", "secretary", "principal", "vice_principal",
		"staff", "coordinator", "teacher", "student", "parent", "anonymous",
	}

	t.Run("1. Grade Addition RBAC Matrix", func(t *testing.T) {
		allowed := map[string]bool{
			"superadmin": true, "admin": true, "teacher": true,
		}

		for _, role := range allRoles {
			currentRole := role
			t.Run("role_"+currentRole, func(t *testing.T) {
				mUser := new(testhelpers.MockUsersRepository)
				mUser.On("GetByID", mock.Anything, mock.Anything).Return(&users.User{
					ID:       "user-1",
					Role:     currentRole,
					SchoolID: &schoolIDPtr,
				}, nil).Maybe()

				mGrade := new(testhelpers.MockGradesRepository)
				mGrade.On("Create", mock.Anything).Return(nil).Maybe()

				gSvc := grades.NewService(mGrade, mUser, nil, nil)
				gH := grades.NewHandler(gSvc, nil)

				r := gin.New()
				g := r.Group("/api/v1")
				g.Use(setAuthHeader("user-1", currentRole, schoolID))
				g.POST("/grades", gH.AddGrade)

				body := fmt.Sprintf(`{"student_id":"%s","subject_id":"subj-1","grade_value":8.0,"grade_type":"numeric","semester":1,"date":"2025-10-15"}`, studentID)
				req := httptest.NewRequest("POST", "/api/v1/grades", bytes.NewBufferString(body))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				r.ServeHTTP(res, req)

				if allowed[currentRole] {
					assert.Equal(t, http.StatusCreated, res.Code, "Role %s should be allowed to add grades", currentRole)
				} else {
					assert.Contains(t, []int{http.StatusForbidden, http.StatusUnauthorized}, res.Code, "Role %s should be denied from adding grades", currentRole)
				}
			})
		}
	})

	t.Run("2. Educazione Civica Shared Grade Creator Restrictions", func(t *testing.T) {
		sharedGrade := &grades.Grade{
			ID:          "grade-civica-101",
			StudentID:   studentID,
			SubjectID:   "subj-civica",
			TeacherID:   teacherID,
			SchoolID:    schoolID,
			GradeValue:  9.0,
			GradeType:   "numeric",
			IsPublished: true,
		}

		mUser := new(testhelpers.MockUsersRepository)
		mUser.On("GetByID", mock.Anything, teacherID).Return(&users.User{
			ID:       teacherID,
			Role:     "teacher",
			SchoolID: &schoolIDPtr,
		}, nil).Maybe()
		mUser.On("GetByID", mock.Anything, otherTeacherID).Return(&users.User{
			ID:       otherTeacherID,
			Role:     "teacher",
			SchoolID: &schoolIDPtr,
		}, nil).Maybe()

		mGrade := new(testhelpers.MockGradesRepository)
		mGrade.On("FindByID", "grade-civica-101").Return(sharedGrade, nil)
		mGrade.On("Update", mock.Anything, mock.Anything).Return(nil)

		gSvc := grades.NewService(mGrade, mUser, nil, nil)
		gH := grades.NewHandler(gSvc, nil)

		// Creator teacher updating grade -> PASS
		rCreator := gin.New()
		gC := rCreator.Group("/api/v1")
		gC.Use(setAuthHeader(teacherID, "teacher", schoolID))
		gC.PUT("/grades/:id", gH.UpdateGrade)

		bodyUp := `{"grade_value":9.5,"reason":"Correzione voto"}`
		reqC := httptest.NewRequest("PUT", "/api/v1/grades/grade-civica-101", bytes.NewBufferString(bodyUp))
		reqC.Header.Set("Content-Type", "application/json")
		resC := httptest.NewRecorder()
		rCreator.ServeHTTP(resC, reqC)
		assert.Equal(t, http.StatusOK, resC.Code)

		// Other teacher updating shared grade created by teacherID -> FORBIDDEN
		rOther := gin.New()
		gO := rOther.Group("/api/v1")
		gO.Use(setAuthHeader(otherTeacherID, "teacher", schoolID))
		gO.PUT("/grades/:id", gH.UpdateGrade)

		reqO := httptest.NewRequest("PUT", "/api/v1/grades/grade-civica-101", bytes.NewBufferString(bodyUp))
		reqO.Header.Set("Content-Type", "application/json")
		resO := httptest.NewRecorder()
		rOther.ServeHTTP(resO, reqO)
		assert.Equal(t, http.StatusForbidden, resO.Code)
	})

	t.Run("3. Attendance Marking RBAC Matrix", func(t *testing.T) {
		allowed := map[string]bool{
			"superadmin": true, "admin": true, "teacher": true,
		}

		for _, role := range allRoles {
			currentRole := role
			t.Run("role_"+currentRole, func(t *testing.T) {
				mUser := new(testhelpers.MockUsersRepository)
				mUser.On("GetByID", mock.Anything, mock.Anything).Return(&users.User{
					ID:       "user-1",
					Role:     currentRole,
					SchoolID: &schoolIDPtr,
				}, nil).Maybe()

				mAtt := new(mockAttRepo)
				mAtt.On("IsTeacherAssignedToClass", mock.Anything, "user-1", classID).Return(true, nil).Maybe()
				mAtt.On("BatchCreate", mock.Anything).Return(nil).Maybe()

				attSvc := attendance.NewService(mAtt, mUser, nil, nil)
				attH := attendance.NewHandler(attSvc)

				r := gin.New()
				g := r.Group("/api/v1")
				g.Use(setAuthHeader("user-1", currentRole, schoolID))
				g.POST("/attendance/mark-bulk", attH.MarkBulk)

				body := fmt.Sprintf(`{"class_id":"%s","date":"%s","hour":1,"statuses":[{"student_id":"%s","status":"Present"}]}`, classID, todayStr, studentID)
				req := httptest.NewRequest("POST", "/api/v1/attendance/mark-bulk", bytes.NewBufferString(body))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				r.ServeHTTP(res, req)

				if allowed[currentRole] {
					assert.Equal(t, http.StatusOK, res.Code, "Role %s should be allowed to mark attendance", currentRole)
				} else {
					assert.Contains(t, []int{http.StatusForbidden, http.StatusUnauthorized}, res.Code, "Role %s should be denied from marking attendance", currentRole)
				}
			})
		}
	})

	t.Run("4. Substitutions Management RBAC Matrix", func(t *testing.T) {
		allowed := map[string]bool{
			"superadmin": true, "admin": true, "secretary": true, "principal": true, "vice_principal": true, "staff": true,
		}

		for _, role := range allRoles {
			currentRole := role
			t.Run("role_"+currentRole, func(t *testing.T) {
				mUser := new(testhelpers.MockUsersRepository)
				mUser.On("GetByID", mock.Anything, mock.Anything).Return(&users.User{
					ID:       "user-1",
					Role:     currentRole,
					SchoolID: &schoolIDPtr,
					IsStaff:  currentRole == "staff" || currentRole == "vice_principal" || currentRole == "principal",
				}, nil).Maybe()

				mSub := new(mockSubRepo)
				mSub.On("Create", mock.Anything, mock.Anything).Return(nil).Maybe()

				subSvc := substitutions.NewService(mSub)
				subH := substitutions.NewHandler(subSvc)

				r := gin.New()
				g := r.Group("/api/v1")
				g.Use(func(c *gin.Context) {
					if currentRole != "anonymous" {
						c.Set("user_id", "user-1")
						c.Set("role", currentRole)
						c.Set("school_id", schoolID)
						if currentRole == "staff" || currentRole == "vice_principal" || currentRole == "principal" {
							c.Set("is_staff", true)
						}
					}
					c.Next()
				})
				g.POST("/substitutions", subH.Create)

				body := fmt.Sprintf(`{"class_id":"%s","date":"%s","hour":2,"absent_teacher_id":"%s"}`, classID, todayStr, teacherID)
				req := httptest.NewRequest("POST", "/api/v1/substitutions", bytes.NewBufferString(body))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				r.ServeHTTP(res, req)

				if allowed[currentRole] {
					assert.Equal(t, http.StatusCreated, res.Code, "Role %s should be allowed to create substitution", currentRole)
				} else {
					assert.Contains(t, []int{http.StatusForbidden, http.StatusUnauthorized}, res.Code, "Role %s should be denied from creating substitution", currentRole)
				}
			})
		}
	})

	t.Run("5. Scrutiny Overview Access RBAC Matrix", func(t *testing.T) {
		allowed := map[string]bool{
			"superadmin": true, "admin": true, "secretary": true, "principal": true,
		}

		for _, role := range allRoles {
			currentRole := role
			t.Run("role_"+currentRole, func(t *testing.T) {
				mUser := new(testhelpers.MockUsersRepository)
				mUser.On("GetByID", mock.Anything, mock.Anything).Return(&users.User{
					ID:       "user-1",
					Role:     currentRole,
					SchoolID: &schoolIDPtr,
				}, nil).Maybe()

				mGrade := new(testhelpers.MockGradesRepository)
				mAtt := new(mockAttRepo)
				mClass := new(testhelpers.MockClassesRepository)
				mClass.On("List", mock.Anything, mock.Anything, mock.Anything).Return([]classes.Class{}, nil).Maybe()

				mScrutiny := new(mockScrutinyRepo)
				mScrutiny.On("GetOverview", mock.Anything, schoolID).Return(map[string]interface{}{"status": "ok"}, nil).Maybe()

				scrutinySvc := scrutiny.NewService(mScrutiny, mGrade, mClass, mUser, mAtt)
				scrutinyH := scrutiny.NewHandler(scrutinySvc, nil)

				r := gin.New()
				g := r.Group("/api/v1")
				g.Use(setAuthHeader("user-1", currentRole, schoolID))
				g.GET("/scrutiny/overview", scrutinyH.GetOverview)

				req := httptest.NewRequest("GET", "/api/v1/scrutiny/overview", nil)
				res := httptest.NewRecorder()
				r.ServeHTTP(res, req)

				if allowed[currentRole] {
					assert.Equal(t, http.StatusOK, res.Code, "Role %s should be allowed scrutiny overview access", currentRole)
				} else {
					assert.Contains(t, []int{http.StatusForbidden, http.StatusUnauthorized}, res.Code, "Role %s should be denied scrutiny overview access", currentRole)
				}
			})
		}
	})
}
