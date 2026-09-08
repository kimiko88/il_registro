package teachers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateTeacherRegisterPDF(t *testing.T) {
	data := &TeacherRegisterData{
		TeacherName:  "Prof. Giovanni Falcone",
		SchoolName:   "Liceo Scientifico Statale Leonardo",
		ClassName:    "3ª A",
		SubjectName:  "Matematica e Fisica",
		AcademicYear: "2023/2024",
		Students: []TeacherRegisterStudent{
			{
				ID:          "st-1",
				Name:        "Mario Rossi",
				Q1Written:   "7.5",
				Q1Oral:      "8.0",
				Q1Practical: "-",
				Q1Avg:       "7.75",
				Q2Written:   "8.0",
				Q2Oral:      "8.5",
				Q2Practical: "-",
				Q2Avg:       "8.25",
				FinalAvg:    "8.00",
				Absences:    4,
			},
			{
				ID:          "st-2",
				Name:        "Luigi Bianchi",
				Q1Written:   "6.0",
				Q1Oral:      "6.5",
				Q1Practical: "-",
				Q1Avg:       "6.25",
				Q2Written:   "7.0",
				Q2Oral:      "7.0",
				Q2Practical: "-",
				Q2Avg:       "7.00",
				FinalAvg:    "6.63",
				Absences:    2,
			},
		},
		Lessons: []TeacherRegisterLesson{
			{
				Date:     "15/09/2023",
				Hour:     1,
				Topic:    "Introduzione alle funzioni reali di variabile reale",
				Type:     "Frontale",
				SignedBy: "Prof. Giovanni Falcone (FEQ)",
			},
			{
				Date:     "16/09/2023",
				Hour:     2,
				Topic:    "Esercitazione guidata su dominio e iniettività",
				Type:     "Laboratorio",
				SignedBy: "Prof. Giovanni Falcone (FEQ)",
			},
		},
	}

	pdfBytes, err := GenerateTeacherRegisterPDF(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, pdfBytes)
	assert.True(t, len(pdfBytes) > 1000)
	// PDF magic bytes %PDF-
	assert.Equal(t, "%PDF-", string(pdfBytes[0:5]))
}

func TestHandler_GetPersonalRegisterPDF(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	handler := NewHandler(svc)

	mockRepo.On("GetPersonalRegisterData", mock.Anything, "t-1", "c-1", "sub-1").Return(&TeacherRegisterData{
		TeacherName:  "Prof. Falcone",
		SchoolName:   "Liceo Leonardo",
		ClassName:    "3A",
		SubjectName:  "Matematica",
		AcademicYear: "2023/2024",
		Students:     []TeacherRegisterStudent{},
		Lessons:      []TeacherRegisterLesson{},
	}, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", "t-1")
	c.Request = httptest.NewRequest("GET", "/teachers/registro-personale/pdf?teacher_id=t-1&class_id=c-1&subject_id=sub-1", nil)

	handler.GetPersonalRegisterPDF(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "registro_personale_docente.pdf")
	assert.True(t, w.Body.Len() > 500)
}
