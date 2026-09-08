package classes

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepository_GetMonthlyJournalData_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()
	classID := uuid.New().String()
	studentID := uuid.New().String()

	// 1. Class & School metadata query
	classRows := sqlmock.NewRows([]string{"class_name", "school_name", "coordinator_name"}).
		AddRow("5A", "Liceo Scientifico", "Mario Rossi")
	mock.ExpectQuery(`SELECT COALESCE\(al\.name`).
		WithArgs(classID).
		WillReturnRows(classRows)

	// 2. Students query
	stRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(studentID, "Bianchi Luigi")
	mock.ExpectQuery(`SELECT s\.id::text`).
		WithArgs(classID).
		WillReturnRows(stRows)

	// 3. Attendance query
	attRows := sqlmock.NewRows([]string{"student_id", "day", "status"}).
		AddRow(studentID, 10, "Present").
		AddRow(studentID, 11, "Absent").
		AddRow(studentID, 12, "Late").
		AddRow(studentID, 13, "LeftEarly")
	mock.ExpectQuery(`SELECT a\.student_id::text`).
		WithArgs(classID, "2026-09-01", "2026-10-01").
		WillReturnRows(attRows)

	// 4. Lessons query
	lessonsRows := sqlmock.NewRows([]string{"date", "hour", "teacher_name", "subject_name", "topic", "type"}).
		AddRow("10/09/2026", 1, "Mario Rossi", "Matematica", "Derivate", "Lezione")
	mock.ExpectQuery(`SELECT TO_CHAR\(l\.date`).
		WithArgs(classID, "2026-09-01", "2026-10-01").
		WillReturnRows(lessonsRows)

	// 5. Disciplinary notes query
	notesRows := sqlmock.NewRows([]string{"date", "student_name", "teacher_name", "description", "note_type"}).
		AddRow("10/09/2026", "Bianchi Luigi", "Mario Rossi", "Disturba", "Nota")
	mock.ExpectQuery(`SELECT TO_CHAR\(sn\.date`).
		WithArgs(classID, "2026-09-01", "2026-10-01").
		WillReturnRows(notesRows)

	data, err := repo.GetMonthlyJournalData(ctx, classID, 2026, 9)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "5A", data.ClassName)
	assert.Equal(t, "Liceo Scientifico", data.SchoolName)
	assert.Equal(t, 1, len(data.Students))
	assert.Equal(t, "Bianchi Luigi", data.Students[0].Name)
	assert.Equal(t, 1, len(data.Lessons))
	assert.Equal(t, 1, len(data.DisciplinaryNotes))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetMonthlyJournalData_RowsErr(t *testing.T) {
	t.Run("Students rows error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := NewRepository(db)
		ctx := context.Background()
		classID := uuid.New().String()

		classRows := sqlmock.NewRows([]string{"class_name", "school_name", "coordinator_name"}).
			AddRow("5A", "Liceo Scientifico", "Mario Rossi")
		mock.ExpectQuery(`SELECT COALESCE\(al\.name`).
			WithArgs(classID).
			WillReturnRows(classRows)

		stRows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("st-1", "Student One").
			RowError(0, errors.New("st row error"))
		mock.ExpectQuery(`SELECT s\.id::text`).
			WithArgs(classID).
			WillReturnRows(stRows)

		data, err := repo.GetMonthlyJournalData(ctx, classID, 2026, 9)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "st row error", err.Error())
	})

	t.Run("Attendance rows error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := NewRepository(db)
		ctx := context.Background()
		classID := uuid.New().String()

		classRows := sqlmock.NewRows([]string{"class_name", "school_name", "coordinator_name"}).
			AddRow("5A", "Liceo Scientifico", "Mario Rossi")
		mock.ExpectQuery(`SELECT COALESCE\(al\.name`).
			WithArgs(classID).
			WillReturnRows(classRows)

		stRows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("st-1", "Student One")
		mock.ExpectQuery(`SELECT s\.id::text`).
			WithArgs(classID).
			WillReturnRows(stRows)

		attRows := sqlmock.NewRows([]string{"student_id", "day", "status"}).
			AddRow("st-1", 1, "Present").
			RowError(0, errors.New("att row error"))
		mock.ExpectQuery(`SELECT a\.student_id::text`).
			WithArgs(classID, "2026-09-01", "2026-10-01").
			WillReturnRows(attRows)

		data, err := repo.GetMonthlyJournalData(ctx, classID, 2026, 9)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "att row error", err.Error())
	})

	t.Run("Lessons rows error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := NewRepository(db)
		ctx := context.Background()
		classID := uuid.New().String()

		classRows := sqlmock.NewRows([]string{"class_name", "school_name", "coordinator_name"}).
			AddRow("5A", "Liceo Scientifico", "Mario Rossi")
		mock.ExpectQuery(`SELECT COALESCE\(al\.name`).
			WithArgs(classID).
			WillReturnRows(classRows)

		stRows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("st-1", "Student One")
		mock.ExpectQuery(`SELECT s\.id::text`).
			WithArgs(classID).
			WillReturnRows(stRows)

		attRows := sqlmock.NewRows([]string{"student_id", "day", "status"}).
			AddRow("st-1", 1, "Present")
		mock.ExpectQuery(`SELECT a\.student_id::text`).
			WithArgs(classID, "2026-09-01", "2026-10-01").
			WillReturnRows(attRows)

		lessonsRows := sqlmock.NewRows([]string{"date", "hour", "teacher_name", "subject_name", "topic", "type"}).
			AddRow("10/09/2026", 1, "Teacher", "Subject", "Topic", "Type").
			RowError(0, errors.New("lessons row error"))
		mock.ExpectQuery(`SELECT TO_CHAR\(l\.date`).
			WithArgs(classID, "2026-09-01", "2026-10-01").
			WillReturnRows(lessonsRows)

		data, err := repo.GetMonthlyJournalData(ctx, classID, 2026, 9)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "lessons row error", err.Error())
	})

	t.Run("Notes rows error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := NewRepository(db)
		ctx := context.Background()
		classID := uuid.New().String()

		classRows := sqlmock.NewRows([]string{"class_name", "school_name", "coordinator_name"}).
			AddRow("5A", "Liceo Scientifico", "Mario Rossi")
		mock.ExpectQuery(`SELECT COALESCE\(al\.name`).
			WithArgs(classID).
			WillReturnRows(classRows)

		stRows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("st-1", "Student One")
		mock.ExpectQuery(`SELECT s\.id::text`).
			WithArgs(classID).
			WillReturnRows(stRows)

		attRows := sqlmock.NewRows([]string{"student_id", "day", "status"}).
			AddRow("st-1", 1, "Present")
		mock.ExpectQuery(`SELECT a\.student_id::text`).
			WithArgs(classID, "2026-09-01", "2026-10-01").
			WillReturnRows(attRows)

		lessonsRows := sqlmock.NewRows([]string{"date", "hour", "teacher_name", "subject_name", "topic", "type"}).
			AddRow("10/09/2026", 1, "Teacher", "Subject", "Topic", "Type")
		mock.ExpectQuery(`SELECT TO_CHAR\(l\.date`).
			WithArgs(classID, "2026-09-01", "2026-10-01").
			WillReturnRows(lessonsRows)

		notesRows := sqlmock.NewRows([]string{"date", "student_name", "teacher_name", "description", "note_type"}).
			AddRow("10/09/2026", "Student One", "Teacher", "Desc", "Note").
			RowError(0, errors.New("notes row error"))
		mock.ExpectQuery(`SELECT TO_CHAR\(sn\.date`).
			WithArgs(classID, "2026-09-01", "2026-10-01").
			WillReturnRows(notesRows)

		data, err := repo.GetMonthlyJournalData(ctx, classID, 2026, 9)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "notes row error", err.Error())
	})
}
