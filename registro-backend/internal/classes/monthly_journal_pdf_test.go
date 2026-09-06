package classes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateMonthlyJournalPDF(t *testing.T) {
	data := &MonthlyJournalData{
		SchoolName:      "Liceo Statale Scientifico G. Galilei",
		ClassName:       "4ª B Scientifico",
		CoordinatorName: "Prof.ssa Maria Rossi",
		Month:           10,
		Year:            2023,
		MonthName:       "Ottobre",
		DaysInMonth:     31,
		Students: []MonthlyStudentAttendance{
			{
				ID:   "st-1",
				Name: "Mario Rossi",
				DailyStatus: map[int]string{
					1: "P", 2: "P", 3: "A", 4: "P", 5: "R", 6: "P",
				},
				TotalP: 20,
				TotalA: 1,
				TotalR: 1,
				TotalU: 0,
			},
			{
				ID:   "st-2",
				Name: "Giulia Verdi",
				DailyStatus: map[int]string{
					1: "P", 2: "P", 3: "P", 4: "P", 5: "P", 6: "U",
				},
				TotalP: 21,
				TotalA: 0,
				TotalR: 0,
				TotalU: 1,
			},
		},
		Lessons: []MonthlyLesson{
			{
				Date:        "02/10/2023",
				Hour:        1,
				TeacherName: "Prof.ssa Maria Rossi",
				SubjectName: "Matematica",
				Topic:       "Limiti notevoli di funzioni goniometriche ed esponenziali",
				Type:        "Frontale",
			},
			{
				Date:        "03/10/2023",
				Hour:        2,
				TeacherName: "Prof. Marco Bianchi",
				SubjectName: "Italiano",
				Topic:       "Dante: Inferno, Canto XXVI e la figura di Ulisse",
				Type:        "Frontale",
			},
		},
		DisciplinaryNotes: []MonthlyDisciplinaryNote{
			{
				Date:        "05/10/2023",
				StudentName: "Mario Rossi",
				TeacherName: "Prof. Marco Bianchi",
				Description: "Disturba ripetutamente la spiegazione della lezione",
				NoteType:    "Comportamentale",
			},
		},
		Stats: MonthlyJournalStats{
			TotalSchoolDays: 26,
			TotalAbsences:   1,
			TotalLates:      1,
			TotalEarlyExits: 1,
			AttendanceRate:  97.7,
		},
	}

	pdfBytes, err := GenerateMonthlyJournalPDF(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, pdfBytes)
	assert.True(t, len(pdfBytes) > 1000)
	assert.Equal(t, "%PDF-", string(pdfBytes[0:5]))
}

func TestHandler_GetMonthlyJournalPDF(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	handler := NewHandler(svc)

	mockRepo.On("Get", mock.Anything, "class-1").Return(&Class{
		ID:       "class-1",
		SchoolID: "school-1",
		Name:     "4B",
	}, nil).Once()

	mockRepo.On("GetMonthlyJournalData", mock.Anything, "class-1", 2023, 10).Return(&MonthlyJournalData{
		SchoolName:      "Liceo Galilei",
		ClassName:       "4B",
		CoordinatorName: "Prof.ssa Rossi",
		Month:           10,
		Year:            2023,
		MonthName:       "Ottobre",
		DaysInMonth:     31,
		Students:        []MonthlyStudentAttendance{},
		Lessons:         []MonthlyLesson{},
		Stats:           MonthlyJournalStats{TotalSchoolDays: 26},
	}, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", "u-1")
	c.Set("school_id", "school-1")
	c.Params = gin.Params{{Key: "id", Value: "class-1"}}
	c.Request = httptest.NewRequest("GET", "/classes/class-1/giornale-mensile/pdf?year=2023&month=10", nil)

	handler.GetMonthlyJournalPDF(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "giornale_classe_class-1_2023_10.pdf")
	assert.True(t, w.Body.Len() > 500)
}
