package grades

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidateReligionJudgment(t *testing.T) {
	v := NewValidator(nil)

	// Valid cases
	cases := []struct {
		input       string
		expectedVal float64
	}{
		{"Non classificabile", 0},
		{"non classificabile", 0},
		{"NON CLASSIFICABILE", 0},
		{"Insufficiente", 4},
		{"insufficiente", 4},
		{"Sufficiente", 6},
		{"sufficiente", 6},
		{"Buono", 7},
		{"buono", 7},
		{"Distinto", 8},
		{"distinto", 8},
		{"Ottimo", 10},
		{"ottimo", 10},
		{"OTTIMO", 10},
		{"Ottimo - verifica 1° quadrimestre", 10},
		{"Distinto: interrogazione", 8},
		{"Buono, partecipazione attiva", 7},
		{"Insufficiente - non ha svolto il compito", 4},
		{"Non classificabile - assenze prolungate", 0},
	}

	for _, tc := range cases {
		val, err := v.ValidateReligionJudgment(tc.input)
		assert.NoError(t, err, "input '%s' should be valid", tc.input)
		assert.Equal(t, tc.expectedVal, val)
	}

	// Invalid cases
	invalidCases := []string{
		"Discreto",
		"Eccellente",
		"10",
		"6",
		"Avanzato",
		"",
		"Non Giudicabile",
		"Gravemente insufficiente",
	}

	for _, inv := range invalidCases {
		_, err := v.ValidateReligionJudgment(inv)
		assert.Error(t, err, "input '%s' should be invalid for religion", inv)
	}
}

func TestValidateGradeValue_JudgmentZero(t *testing.T) {
	v := NewValidator(nil)

	// 0.0 should be permitted for GradeTypeJudgment (represents "Non classificabile")
	err := v.ValidateGradeValue(0.0, string(GradeTypeJudgment))
	assert.NoError(t, err)

	// Standard values 1 to 10
	assert.NoError(t, v.ValidateGradeValue(4.0, string(GradeTypeJudgment)))
	assert.NoError(t, v.ValidateGradeValue(6.0, string(GradeTypeJudgment)))
	assert.NoError(t, v.ValidateGradeValue(10.0, string(GradeTypeJudgment)))

	// Out of bounds
	assert.Error(t, v.ValidateGradeValue(-1.0, string(GradeTypeJudgment)))
	assert.Error(t, v.ValidateGradeValue(10.5, string(GradeTypeJudgment)))
}

func TestIsReligionSubject(t *testing.T) {
	t.Run("nil db returns false", func(t *testing.T) {
		v := NewValidator(nil)
		isRel, err := v.IsReligionSubject(context.Background(), "sub-1")
		assert.NoError(t, err)
		assert.False(t, isRel)
	})

	t.Run("empty subjectID returns false", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		v := NewValidator(db)
		isRel, err := v.IsReligionSubject(context.Background(), "")
		assert.NoError(t, err)
		assert.False(t, isRel)
	})

	t.Run("db returns true for religion subject", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT COALESCE\\(is_religion, FALSE\\) FROM subjects").
			WithArgs("sub-religion").
			WillReturnRows(sqlmock.NewRows([]string{"is_religion"}).AddRow(true))

		v := NewValidator(db)
		isRel, err := v.IsReligionSubject(context.Background(), "sub-religion")
		assert.NoError(t, err)
		assert.True(t, isRel)
	})

	t.Run("db returns false for standard subject", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT COALESCE\\(is_religion, FALSE\\) FROM subjects").
			WithArgs("sub-math").
			WillReturnRows(sqlmock.NewRows([]string{"is_religion"}).AddRow(false))

		v := NewValidator(db)
		isRel, err := v.IsReligionSubject(context.Background(), "sub-math")
		assert.NoError(t, err)
		assert.False(t, isRel)
	})

	t.Run("db returns ErrNoRows defaults to false without error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT COALESCE\\(is_religion, FALSE\\) FROM subjects").
			WithArgs("sub-unknown").
			WillReturnError(sql.ErrNoRows)

		v := NewValidator(db)
		isRel, err := v.IsReligionSubject(context.Background(), "sub-unknown")
		assert.NoError(t, err)
		assert.False(t, isRel)
	})
}

func TestIsStudentAvvalente(t *testing.T) {
	t.Run("nil db returns true by default", func(t *testing.T) {
		v := NewValidator(nil)
		isAvv, err := v.IsStudentAvvalente(context.Background(), "stu-1", "sch-1")
		assert.NoError(t, err)
		assert.True(t, isAvv)
	})

	t.Run("student is avvalente", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT choice::text FROM student_religion_choices").
			WithArgs("stu-1", "sch-1").
			WillReturnRows(sqlmock.NewRows([]string{"choice"}).AddRow("avvalente"))

		v := NewValidator(db)
		isAvv, err := v.IsStudentAvvalente(context.Background(), "stu-1", "sch-1")
		assert.NoError(t, err)
		assert.True(t, isAvv)
	})

	t.Run("student is non_avvalente", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT choice::text FROM student_religion_choices").
			WithArgs("stu-2", "sch-1").
			WillReturnRows(sqlmock.NewRows([]string{"choice"}).AddRow("non_avvalente"))

		v := NewValidator(db)
		isAvv, err := v.IsStudentAvvalente(context.Background(), "stu-2", "sch-1")
		assert.NoError(t, err)
		assert.False(t, isAvv)
	})

	t.Run("student is attivita_alternativa", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT choice::text FROM student_religion_choices").
			WithArgs("stu-3", "sch-1").
			WillReturnRows(sqlmock.NewRows([]string{"choice"}).AddRow("attivita_alternativa"))

		v := NewValidator(db)
		isAvv, err := v.IsStudentAvvalente(context.Background(), "stu-3", "sch-1")
		assert.NoError(t, err)
		assert.False(t, isAvv)
	})

	t.Run("no explicit row defaults to true (avvalente)", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT choice::text FROM student_religion_choices").
			WithArgs("stu-new", "sch-1").
			WillReturnError(sql.ErrNoRows)

		v := NewValidator(db)
		isAvv, err := v.IsStudentAvvalente(context.Background(), "stu-new", "sch-1")
		assert.NoError(t, err)
		assert.True(t, isAvv)
	})
}

func TestGetClassReligionChoices(t *testing.T) {
	t.Run("nil db returns empty map", func(t *testing.T) {
		v := NewValidator(nil)
		choices, err := v.GetClassReligionChoices(context.Background(), "cls-1")
		assert.NoError(t, err)
		assert.Empty(t, choices)
	})

	t.Run("empty classID returns empty map", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		v := NewValidator(db)
		choices, err := v.GetClassReligionChoices(context.Background(), "")
		assert.NoError(t, err)
		assert.Empty(t, choices)
	})

	t.Run("returns mapped student choices and user choices", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT s.id::text, s.user_id::text, COALESCE").
			WithArgs("cls-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "choice"}).
				AddRow("stu-1", "usr-1", "avvalente").
				AddRow("stu-2", "usr-2", "non_avvalente"))

		v := NewValidator(db)
		choices, err := v.GetClassReligionChoices(context.Background(), "cls-1")
		require.NoError(t, err)
		assert.Equal(t, "avvalente", choices["stu-1"])
		assert.Equal(t, "avvalente", choices["usr-1"])
		assert.Equal(t, "non_avvalente", choices["stu-2"])
		assert.Equal(t, "non_avvalente", choices["usr-2"])
	})

	t.Run("query error propagates", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT s.id::text, s.user_id::text, COALESCE").
			WithArgs("cls-1").
			WillReturnError(errors.New("db connection failure"))

		v := NewValidator(db)
		_, err = v.GetClassReligionChoices(context.Background(), "cls-1")
		assert.Error(t, err)
	})

	t.Run("rows.Err error propagates after iteration", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "user_id", "choice"}).
			AddRow("stu-1", "usr-1", "avvalente").
			RowError(0, errors.New("stream broken"))

		mock.ExpectQuery("SELECT s.id::text, s.user_id::text, COALESCE").
			WithArgs("cls-1").
			WillReturnRows(rows)

		v := NewValidator(db)
		_, err = v.GetClassReligionChoices(context.Background(), "cls-1")
		assert.Error(t, err)
	})
}

func TestConvertJudgmentToNumeric_AllCases(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"gravemente insufficiente", 3.0},
		{"insufficiente", 4.0},
		{"non raggiunto", 4.0},
		{"mediocre", 5.0},
		{"quasi sufficiente", 5.0},
		{"iniziale", 5.0},
		{"sufficiente", 6.0},
		{"base", 6.0},
		{"discreto", 7.0},
		{"buono", 8.0},
		{"intermedio", 8.0},
		{"distinto", 9.0},
		{"ottimo", 10.0},
		{"eccellente", 10.0},
		{"avanzato", 10.0},
		{"sconosciuto", 0.0},
		{"", 0.0},
	}

	for _, tc := range tests {
		got := ConvertJudgmentToNumeric(tc.input)
		assert.Equal(t, tc.expected, got, "input: %s", tc.input)
	}
}

func TestConvertNumericToJudgment_AllRanges(t *testing.T) {
	tests := []struct {
		val      float64
		expected string
	}{
		{0.5, "Non valutato"},
		{11.0, "Non valutato"},
		{2.5, "Gravemente Insufficiente"},
		{5.0, "Insufficiente"},
		{6.5, "Sufficiente"},
		{7.5, "Discreto"},
		{8.5, "Buono"},
		{9.5, "Distinto"},
		{10.0, "Ottimo"},
	}

	for _, tc := range tests {
		got := ConvertNumericToJudgment(tc.val)
		assert.Equal(t, tc.expected, got, "val: %v", tc.val)
	}
}

func TestValidator_ValidateGradeDateAndDescription(t *testing.T) {
	v := NewValidator(nil)

	now := time.Now()
	// Future date > 24h
	futureDate := now.Add(48 * time.Hour)
	assert.Error(t, v.ValidateGradeDate(futureDate, 1))

	// Date too old > 10 years
	oldDate := now.AddDate(-11, 0, 0)
	assert.Error(t, v.ValidateGradeDate(oldDate, 1))

	// Valid date within semester 1
	startSem1, _ := GetSemesterDateRange(1)
	validSem1Date := startSem1.AddDate(0, 1, 0)
	if validSem1Date.Before(now.Add(24 * time.Hour)) {
		assert.NoError(t, v.ValidateGradeDate(validSem1Date, 1))
	}

	// Date outside of semester 2
	assert.Error(t, v.ValidateGradeDate(validSem1Date, 2))

	// Description validation
	assert.NoError(t, v.ValidateDescription("Verifica orale su Dante"))
	assert.NoError(t, v.ValidateDescription(""))
	longDesc := ""
	for i := 0; i < 505; i++ {
		longDesc += "a"
	}
	assert.Error(t, v.ValidateDescription(longDesc))
}

func TestService_AddGrade_ReligionValidations(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockRepo := new(MockRepository)
	mockUserRepo := new(MockUserRepo)
	validator := NewValidator(db)

	svc := &service{
		repo:       mockRepo,
		userRepo:   mockUserRepo,
		validator:  validator,
		calculator: NewCalculator(),
	}

	schoolID := "sch-1"
	teacherUser := &users.User{
		ID:       "teach-1",
		SchoolID: &schoolID,
		Role:     "teacher",
	}

	validDate := time.Now().Format("2006-01-02")

	// 1. Rejects if GradeType is not "judgment" for religion subject
	t.Run("Rejects non-judgment grade type for religion", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, "teach-1").Return(teacherUser, nil).Once()

		sqlMock.ExpectQuery("SELECT COALESCE\\(is_religion, FALSE\\) FROM subjects").
			WithArgs("sub-rel").
			WillReturnRows(sqlmock.NewRows([]string{"is_religion"}).AddRow(true))

		req := CreateGradeRequest{
			StudentID:   "stu-1",
			SubjectID:   "sub-rel",
			GradeValue:  8.0,
			GradeType:   "decimal",
			Semester:    1,
			Date:        validDate,
			Description: "Ottimo",
		}

		_, err := svc.AddGrade(context.Background(), "teach-1", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "solo voti di tipo 'judgment'")
	})

	// 2. Rejects if Description is not a valid IRC judgment
	t.Run("Rejects invalid IRC judgment text", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, "teach-1").Return(teacherUser, nil).Once()

		sqlMock.ExpectQuery("SELECT COALESCE\\(is_religion, FALSE\\) FROM subjects").
			WithArgs("sub-rel").
			WillReturnRows(sqlmock.NewRows([]string{"is_religion"}).AddRow(true))

		req := CreateGradeRequest{
			StudentID:   "stu-1",
			SubjectID:   "sub-rel",
			GradeValue:  0,
			GradeType:   "judgment",
			Semester:    1,
			Date:        validDate,
			Description: "Discreto", // Not permitted in religion
		}

		_, err := svc.AddGrade(context.Background(), "teach-1", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "giudizio IRC non valido")
	})

	// 3. Rejects if student is non_avvalente
	t.Run("Rejects grade for non_avvalente student", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, "teach-1").Return(teacherUser, nil).Once()

		sqlMock.ExpectQuery("SELECT COALESCE\\(is_religion, FALSE\\) FROM subjects").
			WithArgs("sub-rel").
			WillReturnRows(sqlmock.NewRows([]string{"is_religion"}).AddRow(true))

		sqlMock.ExpectQuery("SELECT choice::text FROM student_religion_choices").
			WithArgs("stu-exempt", schoolID).
			WillReturnRows(sqlmock.NewRows([]string{"choice"}).AddRow("non_avvalente"))

		req := CreateGradeRequest{
			StudentID:   "stu-exempt",
			SubjectID:   "sub-rel",
			GradeValue:  0,
			GradeType:   "judgment",
			Semester:    1,
			Date:        validDate,
			Description: "Ottimo",
		}

		_, err := svc.AddGrade(context.Background(), "teach-1", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "lo studente non si avvale dell'IRC")
	})

	// 4. Accepts valid judgment for avvalente student
	t.Run("Accepts valid judgment for avvalente student", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, "teach-1").Return(teacherUser, nil).Once()

		sqlMock.ExpectQuery("SELECT COALESCE\\(is_religion, FALSE\\) FROM subjects").
			WithArgs("sub-rel").
			WillReturnRows(sqlmock.NewRows([]string{"is_religion"}).AddRow(true))

		sqlMock.ExpectQuery("SELECT choice::text FROM student_religion_choices").
			WithArgs("stu-avv", schoolID).
			WillReturnRows(sqlmock.NewRows([]string{"choice"}).AddRow("avvalente"))

		sqlMock.ExpectQuery("SELECT id FROM teachers WHERE user_id = \\$1").
			WithArgs("teach-1").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("prof-1"))

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(g *Grade) bool {
			return g.StudentID == "stu-avv" && g.GradeType == GradeTypeJudgment && g.GradeValue == 10.0
		})).Return(nil).Once()

		req := CreateGradeRequest{
			StudentID:   "stu-avv",
			SubjectID:   "sub-rel",
			GradeValue:  0,
			GradeType:   "judgment",
			Semester:    1,
			Date:        validDate,
			Description: "Ottimo",
		}

		resp, err := svc.AddGrade(context.Background(), "teach-1", req)
		assert.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 10.0, resp.GradeValue)
	})
}
