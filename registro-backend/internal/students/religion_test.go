package students

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestReligionChoice_IsValid(t *testing.T) {
	assert.True(t, ReligionChoiceAvvalente.IsValid())
	assert.True(t, ReligionChoiceNonAvvalente.IsValid())
	assert.True(t, ReligionChoiceAttivitaAlternativa.IsValid())

	assert.False(t, ReligionChoice("").IsValid())
	assert.False(t, ReligionChoice("altro").IsValid())
	assert.False(t, ReligionChoice("esonero").IsValid())
}

func TestReligionRepository_GetChoice(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewReligionRepository(db)

	t.Run("existing explicit choice", func(t *testing.T) {
		studentID := "stu-100"
		schoolID := "sch-1"

		mock.ExpectQuery("SELECT s.id, s.school_id FROM students s").
			WithArgs(studentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "school_id"}).AddRow(studentID, schoolID))

		now := time.Now()
		mock.ExpectQuery("SELECT COALESCE\\(src.id::text").
			WithArgs(studentID).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "student_id", "user_id", "school_id", "choice",
				"first_name", "last_name", "class_id", "class_name",
				"enrollment_number", "updated_by", "updated_at",
			}).AddRow("rec-1", studentID, "usr-1", schoolID, "non_avvalente", "Mario", "Rossi", "cls-1", "3A", "MAT123", "sec-1", now))

		item, err := repo.GetChoice(context.Background(), studentID, schoolID)
		require.NoError(t, err)
		assert.Equal(t, ReligionChoiceNonAvvalente, item.Choice)
		assert.Equal(t, "Mario", item.FirstName)
		assert.Equal(t, "Rossi", item.LastName)
		assert.Equal(t, "3A", item.ClassName)
	})

	t.Run("default choice when no explicit row exists", func(t *testing.T) {
		studentID := "stu-200"
		schoolID := "sch-1"

		mock.ExpectQuery("SELECT s.id, s.school_id FROM students s").
			WithArgs(studentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "school_id"}).AddRow(studentID, schoolID))

		mock.ExpectQuery("SELECT COALESCE\\(src.id::text").
			WithArgs(studentID).
			WillReturnError(sql.ErrNoRows)

		item, err := repo.GetChoice(context.Background(), studentID, schoolID)
		require.NoError(t, err)
		assert.Equal(t, ReligionChoiceAvvalente, item.Choice)
	})

	t.Run("invalid student resolution", func(t *testing.T) {
		mock.ExpectQuery("SELECT s.id, s.school_id FROM students s").
			WithArgs("unknown").
			WillReturnError(sql.ErrNoRows)

		item, err := repo.GetChoice(context.Background(), "unknown", "")
		assert.Error(t, err)
		assert.Nil(t, item)
	})
}

func TestReligionRepository_SetChoice(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewReligionRepository(db)

	t.Run("invalid choice returns error", func(t *testing.T) {
		_, err := repo.SetChoice(context.Background(), "stu-1", "sch-1", "sec-1", "invalid_choice")
		assert.Error(t, err)
	})

	t.Run("valid choice upsert succeeds", func(t *testing.T) {
		studentID := "11111111-1111-1111-1111-111111111111"
		schoolID := "22222222-2222-2222-2222-222222222222"
		updatedBy := "33333333-3333-3333-3333-333333333333"

		mock.ExpectQuery("SELECT s.id, s.school_id FROM students s").
			WithArgs(studentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "school_id"}).AddRow(studentID, schoolID))

		mock.ExpectExec("INSERT INTO student_religion_choices").
			WithArgs(studentID, schoolID, "attivita_alternativa", updatedBy).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT s.id, s.school_id FROM students s").
			WithArgs(studentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "school_id"}).AddRow(studentID, schoolID))

		now := time.Now()
		mock.ExpectQuery("SELECT COALESCE\\(src.id::text").
			WithArgs(studentID).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "student_id", "user_id", "school_id", "choice",
				"first_name", "last_name", "class_id", "class_name",
				"enrollment_number", "updated_by", "updated_at",
			}).AddRow("rec-1", studentID, "usr-1", schoolID, "attivita_alternativa", "Luigi", "Verdi", "cls-1", "4B", "MAT456", updatedBy, now))

		item, err := repo.SetChoice(context.Background(), studentID, schoolID, updatedBy, ReligionChoiceAttivitaAlternativa)
		require.NoError(t, err)
		assert.Equal(t, ReligionChoiceAttivitaAlternativa, item.Choice)
	})
}

func TestReligionRepository_BatchSetChoices(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewReligionRepository(db)

	t.Run("empty list does nothing", func(t *testing.T) {
		err := repo.BatchSetChoices(context.Background(), []string{}, "sch-1", "sec-1", ReligionChoiceNonAvvalente)
		assert.NoError(t, err)
	})

	t.Run("invalid choice errors", func(t *testing.T) {
		err := repo.BatchSetChoices(context.Background(), []string{"s1"}, "sch-1", "sec-1", "bad_choice")
		assert.Error(t, err)
	})

	t.Run("executes transaction for all students", func(t *testing.T) {
		studentIDs := []string{"stu-1", "stu-2"}
		schoolID := "22222222-2222-2222-2222-222222222222"
		updatedBy := "33333333-3333-3333-3333-333333333333"

		mock.ExpectBegin()
		mock.ExpectPrepare("INSERT INTO student_religion_choices")
		mock.ExpectExec("INSERT INTO student_religion_choices").
			WithArgs("stu-1", schoolID, "non_avvalente", updatedBy).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO student_religion_choices").
			WithArgs("stu-2", schoolID, "non_avvalente", updatedBy).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.BatchSetChoices(context.Background(), studentIDs, schoolID, updatedBy, ReligionChoiceNonAvvalente)
		assert.NoError(t, err)
	})
}

func TestReligionHandler_HTTP(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	handler := NewReligionHandler(db)

	r := gin.New()
	// Add mock auth middleware to inject user and school
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "sec-user-1")
		c.Set("school_id", "sch-1")
		c.Set("role", "secretary")
		c.Next()
	})
	handler.RegisterRoutes(r.Group("/api/v1"))

	t.Run("GET /students/:id/religion-choice", func(t *testing.T) {
		studentID := "stu-1"
		mock.ExpectQuery("SELECT s.id, s.school_id FROM students s").
			WithArgs(studentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "school_id"}).AddRow(studentID, "sch-1"))

		mock.ExpectQuery("SELECT COALESCE\\(src.id::text").
			WithArgs(studentID).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "student_id", "user_id", "school_id", "choice",
				"first_name", "last_name", "class_id", "class_name",
				"enrollment_number", "updated_by", "updated_at",
			}).AddRow("rec-1", studentID, "usr-1", "sch-1", "avvalente", "Anna", "Bianchi", "cls-1", "1A", "MAT789", nil, time.Now()))

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/students/"+studentID+"/religion-choice", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp StudentReligionChoiceItem
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, ReligionChoiceAvvalente, resp.Choice)
	})

	t.Run("PUT /students/:id/religion-choice invalid choice", func(t *testing.T) {
		payload := map[string]string{"choice": "invalid_choice"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, "/api/v1/students/stu-1/religion-choice", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /secretary/religion-choices/batch empty student_ids", func(t *testing.T) {
		payload := BatchSetReligionChoiceRequest{
			StudentIDs: []string{},
			Choice:     ReligionChoiceNonAvvalente,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/secretary/religion-choices/batch", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GET /secretary/religion-choices success", func(t *testing.T) {
		mock.ExpectQuery("SELECT COALESCE\\(src.id::text").
			WithArgs("sch-1").
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "student_id", "user_id", "school_id", "choice",
				"first_name", "last_name", "class_id", "class_name",
				"enrollment_number", "updated_by", "updated_at",
			}).AddRow("rec-1", "stu-1", "usr-1", "sch-1", "avvalente", "Mario", "Rossi", "cls-1", "1A", "MAT1", nil, time.Now()))

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/secretary/religion-choices", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var list []StudentReligionChoiceItem
		err := json.Unmarshal(w.Body.Bytes(), &list)
		require.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, "Mario", list[0].FirstName)
	})

	t.Run("GET /secretary/religion-choices forbidden for student role", func(t *testing.T) {
		forbiddenRouter := gin.New()
		forbiddenRouter.Use(func(c *gin.Context) {
			c.Set("user_id", "stu-user-1")
			c.Set("role", "student")
			c.Set("school_id", "sch-1")
			c.Next()
		})
		handler.RegisterRoutes(forbiddenRouter.Group("/api/v1"))

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/secretary/religion-choices", nil)
		w := httptest.NewRecorder()
		forbiddenRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestReligionRepository_ListChoices(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewReligionRepository(db)

	t.Run("list with class filter", func(t *testing.T) {
		mock.ExpectQuery("SELECT COALESCE\\(src.id::text").
			WithArgs("sch-1", "cls-1").
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "student_id", "user_id", "school_id", "choice",
				"first_name", "last_name", "class_id", "class_name",
				"enrollment_number", "updated_by", "updated_at",
			}).AddRow("rec-1", "stu-1", "usr-1", "sch-1", "avvalente", "Mario", "Rossi", "cls-1", "1A", "MAT1", nil, time.Now()))

		items, err := repo.ListChoices(context.Background(), "sch-1", "cls-1")
		require.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, ReligionChoiceAvvalente, items[0].Choice)
	})

	t.Run("database error propagates", func(t *testing.T) {
		mock.ExpectQuery("SELECT COALESCE\\(src.id::text").
			WithArgs("sch-1").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.ListChoices(context.Background(), "sch-1", "")
		assert.Error(t, err)
	})

	t.Run("rows.Err error propagates", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "student_id", "user_id", "school_id", "choice",
			"first_name", "last_name", "class_id", "class_name",
			"enrollment_number", "updated_by", "updated_at",
		}).
			AddRow("rec-1", "stu-1", "usr-1", "sch-1", "avvalente", "Mario", "Rossi", "cls-1", "1A", "MAT1", nil, time.Now()).
			RowError(0, errors.New("read error"))

		mock.ExpectQuery("SELECT COALESCE\\(src.id::text").
			WithArgs("sch-1").
			WillReturnRows(rows)

		_, err := repo.ListChoices(context.Background(), "sch-1", "")
		assert.Error(t, err)
	})
}

