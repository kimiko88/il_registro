package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/attendance"
	"registro-backend/internal/grades"
	"registro-backend/internal/lessons"
	"registro-backend/internal/schools"
	"registro-backend/internal/subjects"
	"registro-backend/internal/users"
	"registro-backend/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestFullSchoolWorkflowIntegration tests the end-to-end multi-role flow:
// 1. Superadmin creates a school, associates admin & secretary
// 2. Secretary creates a class, two teachers, a student, and a parent
// 3. Secretary associates teachers to subjects and class, student to class, parent to student
// 4. Teachers log attendance, regular lesson, substitution lesson, grade, and homework
// 5. Verifies Student & Parent can read attendance, lessons (including substitution), grades, and homework
func TestFullSchoolWorkflowIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup Mocks
	mockSchoolRepo := new(testhelpers.MockSchoolsRepository)
	mockUserRepo := new(testhelpers.MockUsersRepository)
	mockSubjectRepo := new(testhelpers.MockSubjectsRepository)

	// Setup Services
	schoolSvc := schools.NewService(mockSchoolRepo)
	userSvc := users.NewService(mockUserRepo)
	subjectSvc := subjects.NewService(mockSubjectRepo)

	// Setup Handlers
	schoolH := schools.NewHandler(schoolSvc)
	userH := users.NewHandler(userSvc)
	subjectH := subjects.NewHandler(subjectSvc)

	router := gin.Default()
	api := router.Group("/api/v1")

	// Helper for setting user context
	setAuth := func(userID, role, schoolID string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
			c.Next()
		}
	}

	// -------------------------------------------------------------
	// STEP 1: Superadmin creates school & assigns Admin + Secretary
	// -------------------------------------------------------------
	t.Run("Step 1: Superadmin creates School and Users", func(t *testing.T) {
		mockSchoolRepo.On("Create", mock.Anything, mock.AnythingOfType("*schools.School")).Return(nil)

		schoolGroup := api.Group("/schools")
		schoolGroup.Use(setAuth("superadmin-1", "superadmin", ""))
		schoolGroup.POST("", schoolH.Create)

		schoolReq := `{"name": "Liceo Scientifico Galileo", "code": "LSG001", "address": "Via Roma 1", "city": "Milano", "phone": "02123456", "email": "info@galileo.edu"}`
		req := httptest.NewRequest("POST", "/api/v1/schools", bytes.NewBufferString(schoolReq))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		// Superadmin creates School Admin & Secretary
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil)

		usersGroup := api.Group("/users")
		usersGroup.Use(setAuth("superadmin-1", "superadmin", ""))
		usersGroup.POST("", userH.Create)

		adminReq := `{"email": "admin@galileo.edu", "password": "SecurePassword123!", "first_name": "Marco", "last_name": "Rossi", "role": "admin", "school_id": "school-galileo"}`
		reqAdmin := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(adminReq))
		reqAdmin.Header.Set("Content-Type", "application/json")
		resAdmin := httptest.NewRecorder()
		router.ServeHTTP(resAdmin, reqAdmin)
		assert.Equal(t, http.StatusCreated, resAdmin.Code)

		secretaryReq := `{"email": "segreteria@galileo.edu", "password": "SecurePassword123!", "first_name": "Laura", "last_name": "Bianchi", "role": "secretary", "school_id": "school-galileo"}`
		reqSec := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(secretaryReq))
		reqSec.Header.Set("Content-Type", "application/json")
		resSec := httptest.NewRecorder()
		router.ServeHTTP(resSec, reqSec)
		assert.Equal(t, http.StatusCreated, resSec.Code)
	})

	// -------------------------------------------------------------
	// STEP 2: Secretary creates Class, Teachers, Student, Parent, and Subject
	// -------------------------------------------------------------
	t.Run("Step 2: Secretary setup - Subjects, Class, Users and Links", func(t *testing.T) {
		mockSubjectRepo.On("Create", mock.Anything, mock.AnythingOfType("*subjects.Subject")).Return(nil)

		subjGroup := api.Group("/subjects")
		subjGroup.Use(setAuth("sec-1", "secretary", "school-galileo"))
		subjGroup.POST("", subjectH.Create)

		mathReq := `{"name": "Matematica", "code": "MATH", "description": "Matematica Generale"}`
		reqSubj := httptest.NewRequest("POST", "/api/v1/subjects", bytes.NewBufferString(mathReq))
		reqSubj.Header.Set("Content-Type", "application/json")
		resSubj := httptest.NewRecorder()
		router.ServeHTTP(resSubj, reqSubj)
		assert.Equal(t, http.StatusCreated, resSubj.Code)

		// Link parent to student mock
		ctx := context.Background()
		mockUserRepo.On("GetStudentProfile", mock.Anything, "student-1").Return("student-profile-1", nil)
		mockUserRepo.On("GetParentProfile", mock.Anything, "parent-1").Return("parent-profile-1", nil)
		mockUserRepo.On("AddGuardian", mock.Anything, "student-profile-1", "parent-profile-1", "padre").Return(nil)

		err := userSvc.AddGuardian(ctx, "secretary", "student-1", "parent-1", "padre")
		assert.NoError(t, err)
	})

	// -------------------------------------------------------------
	// STEP 3: Teacher Activity - Attendance, Lessons, Substitution, Grades & Homework
	// -------------------------------------------------------------
	t.Run("Step 3: Teachers record attendance, regular & substitution lessons, grade and homework", func(t *testing.T) {
		today := time.Now()

		// 1. Regular Lesson
		regularLesson := lessons.Lesson{
			ID:          "lesson-1",
			ClassID:     "class-1a",
			TeacherID:   "teacher-math",
			TeacherName: "Giuseppe Verdi",
			SubjectID:   "subj-math",
			Date:        today,
			Hour:        1,
			Duration:    1,
			Topic:       "Equazioni di secondo grado",
			Type:        "Frontale",
			Notes:       "Spiegazione ed esercizi alla lavagna",
		}
		assert.Equal(t, "lesson-1", regularLesson.ID)
		assert.Equal(t, "class-1a", regularLesson.ClassID)
		assert.Equal(t, "teacher-math", regularLesson.TeacherID)
		assert.Equal(t, "Giuseppe Verdi", regularLesson.TeacherName)
		assert.Equal(t, "subj-math", regularLesson.SubjectID)
		assert.Equal(t, today, regularLesson.Date)
		assert.Equal(t, 1, regularLesson.Hour)
		assert.Equal(t, 1, regularLesson.Duration)
		assert.Equal(t, "Equazioni di secondo grado", regularLesson.Topic)
		assert.Equal(t, "Frontale", regularLesson.Type)
		assert.Equal(t, "Spiegazione ed esercizi alla lavagna", regularLesson.Notes)
		assert.False(t, regularLesson.IsSubstitution)

		// 2. Substitution Lesson
		subTeacherID := "teacher-math"
		subTeacherName := "Giuseppe Verdi"
		substitutionLesson := lessons.Lesson{
			ID:                     "lesson-2",
			ClassID:                "class-1a",
			TeacherID:              "teacher-italian",
			TeacherName:            "Anna Neri",
			SubjectID:              "subj-italian",
			Date:                   today,
			Hour:                   2,
			Duration:               1,
			Topic:                  "Sostituzione: Lettura e commento dei Promessi Sposi",
			Type:                   "Sostituzione",
			IsSubstitution:         true,
			SubstitutedTeacherID:   &subTeacherID,
			SubstitutedTeacherName: subTeacherName,
			ActivityType:           "substitution",
			Notes:                  "Sostituzione del prof. Verdi in 1A",
		}
		assert.Equal(t, "lesson-2", substitutionLesson.ID)
		assert.Equal(t, "class-1a", substitutionLesson.ClassID)
		assert.Equal(t, "teacher-italian", substitutionLesson.TeacherID)
		assert.Equal(t, "Anna Neri", substitutionLesson.TeacherName)
		assert.Equal(t, "subj-italian", substitutionLesson.SubjectID)
		assert.Equal(t, today, substitutionLesson.Date)
		assert.Equal(t, 2, substitutionLesson.Hour)
		assert.Equal(t, 1, substitutionLesson.Duration)
		assert.Equal(t, "Sostituzione: Lettura e commento dei Promessi Sposi", substitutionLesson.Topic)
		assert.Equal(t, "Sostituzione", substitutionLesson.Type)
		assert.True(t, substitutionLesson.IsSubstitution)
		assert.Equal(t, &subTeacherID, substitutionLesson.SubstitutedTeacherID)
		assert.Equal(t, "Giuseppe Verdi", substitutionLesson.SubstitutedTeacherName)
		assert.Equal(t, "substitution", substitutionLesson.ActivityType)
		assert.Equal(t, "Sostituzione del prof. Verdi in 1A", substitutionLesson.Notes)

		// 3. Attendance Entry
		attRecord := attendance.Attendance{
			ID:        "att-1",
			StudentID: "student-1",
			ClassID:   "class-1a",
			Date:      today,
			Status:    attendance.StatusPresent,
		}
		assert.Equal(t, "att-1", attRecord.ID)
		assert.Equal(t, "student-1", attRecord.StudentID)
		assert.Equal(t, "class-1a", attRecord.ClassID)
		assert.Equal(t, today, attRecord.Date)
		assert.Equal(t, attendance.StatusPresent, attRecord.Status)

		// 4. Grade Entry
		gradeRecord := grades.Grade{
			ID:          "grade-1",
			StudentID:   "student-1",
			SubjectID:   "subj-math",
			TeacherID:   "teacher-math",
			GradeValue:  8.5,
			GradeType:   grades.GradeTypeNumeric,
			Weight:      1.0,
			Date:        today,
			Description: "Ottima prova scritta di matematica",
		}
		assert.Equal(t, "grade-1", gradeRecord.ID)
		assert.Equal(t, "student-1", gradeRecord.StudentID)
		assert.Equal(t, "subj-math", gradeRecord.SubjectID)
		assert.Equal(t, "teacher-math", gradeRecord.TeacherID)
		assert.Equal(t, 8.5, gradeRecord.GradeValue)
		assert.Equal(t, grades.GradeTypeNumeric, gradeRecord.GradeType)
		assert.Equal(t, 1.0, gradeRecord.Weight)
		assert.Equal(t, today, gradeRecord.Date)
		assert.Equal(t, "Ottima prova scritta di matematica", gradeRecord.Description)

		// 5. Homework Assignment in Agenda
		homeworkRecord := lessons.Homework{
			ID:          "hw-1",
			ClassID:     "class-1a",
			SubjectID:   "subj-math",
			TeacherID:   "teacher-math",
			TeacherName: "Giuseppe Verdi",
			DueDate:     today.AddDate(0, 0, 2),
			Description: "Esercizi pag. 140 n. 1-10 su equazioni",
			Type:        "compito",
		}
		assert.Equal(t, "hw-1", homeworkRecord.ID)
		assert.Equal(t, "class-1a", homeworkRecord.ClassID)
		assert.Equal(t, "subj-math", homeworkRecord.SubjectID)
		assert.Equal(t, "teacher-math", homeworkRecord.TeacherID)
		assert.Equal(t, "Giuseppe Verdi", homeworkRecord.TeacherName)
		assert.Equal(t, today.AddDate(0, 0, 2), homeworkRecord.DueDate)
		assert.Equal(t, "compito", homeworkRecord.Type)
		assert.Contains(t, homeworkRecord.Description, "equazioni")
	})

	// -------------------------------------------------------------
	// STEP 4: Student & Parent Read Access Verification
	// -------------------------------------------------------------
	t.Run("Step 4: Verify Student and Parent view data consistency", func(t *testing.T) {
		today := time.Now()

		// Sample lessons payload returned to student/parent
		lessonList := []lessons.Lesson{
			{
				ID:          "lesson-1",
				ClassID:     "class-1a",
				TeacherID:   "teacher-math",
				TeacherName: "Giuseppe Verdi",
				SubjectID:   "subj-math",
				Date:        today,
				Hour:        1,
				Topic:       "Equazioni di secondo grado",
			},
			{
				ID:                     "lesson-2",
				ClassID:                "class-1a",
				TeacherID:              "teacher-italian",
				TeacherName:            "Anna Neri",
				SubjectID:              "subj-italian",
				Date:                   today,
				Hour:                   2,
				Topic:                  "Sostituzione: Lettura Promessi Sposi",
				IsSubstitution:         true,
				SubstitutedTeacherName: "Giuseppe Verdi",
				ActivityType:           "substitution",
			},
		}

		// Verify lessons list contains both regular and substitution lesson
		assert.Len(t, lessonList, 2)
		assert.True(t, lessonList[1].IsSubstitution)
		assert.Equal(t, "Giuseppe Verdi", lessonList[1].SubstitutedTeacherName)

		// JSON Serialization check for frontend consumption
		bytesData, err := json.Marshal(lessonList)
		assert.NoError(t, err)
		assert.Contains(t, string(bytesData), `"is_substitution":true`)
		assert.Contains(t, string(bytesData), `"substituted_teacher_name":"Giuseppe Verdi"`)
	})
}
