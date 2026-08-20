package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/didactic_materials"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockDidacticRepo struct {
	mats map[string]*didactic_materials.DidacticMaterial
}

func (m *mockDidacticRepo) Create(dm *didactic_materials.DidacticMaterial) error {
	if dm.ID == "" {
		dm.ID = "mat-1"
	}
	m.mats[dm.ID] = dm
	return nil
}

func (m *mockDidacticRepo) GetByClass(classID string) ([]didactic_materials.DidacticMaterial, error) {
	res := make([]didactic_materials.DidacticMaterial, 0)
	for _, dm := range m.mats {
		if dm.ClassID == classID {
			res = append(res, *dm)
		}
	}
	return res, nil
}

func (m *mockDidacticRepo) Delete(id, teacherID string) error {
	delete(m.mats, id)
	return nil
}

func TestIntegration_Didactic_Materials_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockDidacticRepo{mats: make(map[string]*didactic_materials.DidacticMaterial)}
	svc := didactic_materials.NewService(repo, nil)
	handler := didactic_materials.NewHandler(svc)

	r := gin.New()
	r.POST("/didactic-materials", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.Create(c)
	})
	r.GET("/didactic-materials/class/:class_id", func(c *gin.Context) {
		c.Set("user_id", "student-1")
		c.Set("role", "student")
		c.Set("school_id", "school-1")
		handler.GetByClass(c)
	})

	// 1. Teacher uploads courseware / study material
	matReq := didactic_materials.CreateMaterialRequest{
		ClassID:       "class-1",
		SubjectID:     "subj-1",
		Title:         "Dispense di Fisica Quantistica - Cap. 1",
		Description:   "Materiale integrativo di supporto alle lezioni del primo trimestre",
		AttachmentURL: "https://storage.school.it/dispense/fisica_cap1.pdf",
	}
	body1, _ := json.Marshal(matReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/didactic-materials", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Student views class materials
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/didactic-materials/class/class-1", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
