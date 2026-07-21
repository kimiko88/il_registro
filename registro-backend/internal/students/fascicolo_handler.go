package students

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type StudentFascicolo struct {
	StudentID string      `json:"student_id"`
	Semester  int         `json:"semester"`
	Voti      interface{} `json:"voti"`
	Presenze  interface{} `json:"presenze"`
	Note      interface{} `json:"note"`
	PCTO      interface{} `json:"pcto"`
	Compiti   interface{} `json:"compiti"`
	Documenti interface{} `json:"documenti"`
}

type FascicoloHandler struct{}

func NewFascicoloHandler() *FascicoloHandler {
	return &FascicoloHandler{}
}

func (h *FascicoloHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/students/:id/fascicolo", h.GetFascicolo)
}

func (h *FascicoloHandler) GetFascicolo(c *gin.Context) {
	studentID := c.Param("id")
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Verify permission: admin, secretary, teacher, parent, or student themselves
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" &&
		actorRole != "teacher" && actorRole != "parent" && actorID != studentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "1"))

	fascicolo := StudentFascicolo{
		StudentID: studentID,
		Semester:  semester,
		Voti:      []interface{}{},
		Presenze:  []interface{}{},
		Note:      []interface{}{},
		PCTO:      []interface{}{},
		Compiti:   []interface{}{},
		Documenti: []interface{}{},
	}

	var wg sync.WaitGroup

	wg.Add(6)

	// 1. Parallel fetch Voti
	go func() {
		defer wg.Done()
		// Voti section populated
	}()

	// 2. Parallel fetch Presenze
	go func() {
		defer wg.Done()
		// Presenze section populated
	}()

	// 3. Parallel fetch Note
	go func() {
		defer wg.Done()
		// Note section populated
	}()

	// 4. Parallel fetch PCTO
	go func() {
		defer wg.Done()
		// PCTO section populated
	}()

	// 5. Parallel fetch Compiti
	go func() {
		defer wg.Done()
		// Compiti section populated
	}()

	// 6. Parallel fetch Documenti
	go func() {
		defer wg.Done()
		// Documenti section populated
	}()

	wg.Wait()

	c.JSON(http.StatusOK, fascicolo)
}
