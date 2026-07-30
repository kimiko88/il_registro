package students

import (
	"database/sql"
	"net/http"
	"strconv"
	"sync"
	"time"

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

type VotoItem struct {
	ID        string    `json:"id"`
	SubjectID string    `json:"subject_id"`
	Value     float64   `json:"value"`
	Category  string    `json:"category"`
	Date      time.Time `json:"date"`
}

type PresenzaItem struct {
	ID     string    `json:"id"`
	Date   time.Time `json:"date"`
	Status string    `json:"status"`
	Notes  string    `json:"notes"`
}

type NotaItem struct {
	ID          string    `json:"id"`
	AuthorName  string    `json:"author_name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type PCTOItem struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Company string    `json:"company"`
	Hours   int       `json:"hours"`
	Date    time.Time `json:"date"`
}

type CompitoItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	SubjectID   string    `json:"subject_id"`
	DueDate     time.Time `json:"due_date"`
	Description string    `json:"description"`
}

type DocumentoItem struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type FascicoloHandler struct {
	db *sql.DB
}

func NewFascicoloHandler(db *sql.DB) *FascicoloHandler {
	return &FascicoloHandler{db: db}
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

	voti := []VotoItem{}
	presenze := []PresenzaItem{}
	note := []NotaItem{}
	pcto := []PCTOItem{}
	compiti := []CompitoItem{}
	documenti := []DocumentoItem{}

	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(6)

	// 1. Voti
	go func() {
		defer wg.Done()
		if h.db == nil {
			return
		}
		rows, err := h.db.QueryContext(c.Request.Context(),
			`SELECT id, subject_id, grade_value, grade_category, date
			 FROM grades
			 WHERE student_id = $1 AND semester = $2 AND is_published = true AND deleted_at IS NULL
			 ORDER BY date DESC`, studentID, semester,
		)
		if err != nil {
			return
		}
		defer rows.Close()
		var items []VotoItem
		for rows.Next() {
			var v VotoItem
			if err := rows.Scan(&v.ID, &v.SubjectID, &v.Value, &v.Category, &v.Date); err == nil {
				items = append(items, v)
			}
		}
		if len(items) > 0 {
			mu.Lock()
			voti = items
			mu.Unlock()
		}
	}()

	// 2. Presenze
	go func() {
		defer wg.Done()
		if h.db == nil {
			return
		}
		rows, err := h.db.QueryContext(c.Request.Context(),
			`SELECT id, date, status, COALESCE(notes, '')
			 FROM attendance
			 WHERE student_id = $1
			 ORDER BY date DESC`, studentID,
		)
		if err != nil {
			return
		}
		defer rows.Close()
		var items []PresenzaItem
		for rows.Next() {
			var p PresenzaItem
			if err := rows.Scan(&p.ID, &p.Date, &p.Status, &p.Notes); err == nil {
				items = append(items, p)
			}
		}
		if len(items) > 0 {
			mu.Lock()
			presenze = items
			mu.Unlock()
		}
	}()

	// 3. Note
	go func() {
		defer wg.Done()
		if h.db == nil {
			return
		}
		rows, err := h.db.QueryContext(c.Request.Context(),
			`SELECT dn.id, COALESCE(u.first_name || ' ' || u.last_name, 'Docente'), dn.description, dn.created_at
			 FROM disciplinary_notes dn
			 LEFT JOIN users u ON dn.author_id = u.id
			 WHERE dn.student_id = $1
			 ORDER BY dn.created_at DESC`, studentID,
		)
		if err != nil {
			return
		}
		defer rows.Close()
		var items []NotaItem
		for rows.Next() {
			var n NotaItem
			if err := rows.Scan(&n.ID, &n.AuthorName, &n.Description, &n.Date); err == nil {
				items = append(items, n)
			}
		}
		if len(items) > 0 {
			mu.Lock()
			note = items
			mu.Unlock()
		}
	}()

	// 4. PCTO
	go func() {
		defer wg.Done()
		if h.db == nil {
			return
		}
		rows, err := h.db.QueryContext(c.Request.Context(),
			`SELECT id, title, COALESCE(company_name, ''), hours, date
			 FROM pcto_activities
			 WHERE student_id = $1
			 ORDER BY date DESC`, studentID,
		)
		if err != nil {
			return
		}
		defer rows.Close()
		var items []PCTOItem
		for rows.Next() {
			var p PCTOItem
			if err := rows.Scan(&p.ID, &p.Title, &p.Company, &p.Hours, &p.Date); err == nil {
				items = append(items, p)
			}
		}
		if len(items) > 0 {
			mu.Lock()
			pcto = items
			mu.Unlock()
		}
	}()

	// 5. Compiti
	go func() {
		defer wg.Done()
		if h.db == nil {
			return
		}
		rows, err := h.db.QueryContext(c.Request.Context(),
			`SELECT hw.id, hw.title, hw.subject_id, hw.due_date, COALESCE(hw.description, '')
			 FROM homeworks hw
			 JOIN class_students cs ON hw.class_id = cs.class_id
			 WHERE cs.student_id = $1
			 ORDER BY hw.due_date DESC`, studentID,
		)
		if err != nil {
			return
		}
		defer rows.Close()
		var items []CompitoItem
		for rows.Next() {
			var c CompitoItem
			if err := rows.Scan(&c.ID, &c.Title, &c.SubjectID, &c.DueDate, &c.Description); err == nil {
				items = append(items, c)
			}
		}
		if len(items) > 0 {
			mu.Lock()
			compiti = items
			mu.Unlock()
		}
	}()

	// 6. Documenti
	go func() {
		defer wg.Done()
		if h.db == nil {
			return
		}
		rows, err := h.db.QueryContext(c.Request.Context(),
			`SELECT id, title, type, status, created_at
			 FROM student_documents
			 WHERE student_id = $1
			 ORDER BY created_at DESC`, studentID,
		)
		if err != nil {
			return
		}
		defer rows.Close()
		var items []DocumentoItem
		for rows.Next() {
			var d DocumentoItem
			if err := rows.Scan(&d.ID, &d.Title, &d.Type, &d.Status, &d.CreatedAt); err == nil {
				items = append(items, d)
			}
		}
		if len(items) > 0 {
			mu.Lock()
			documenti = items
			mu.Unlock()
		}
	}()

	wg.Wait()

	fascicolo := StudentFascicolo{
		StudentID: studentID,
		Semester:  semester,
		Voti:      voti,
		Presenze:  presenze,
		Note:      note,
		PCTO:      pcto,
		Compiti:   compiti,
		Documenti: documenti,
	}

	c.JSON(http.StatusOK, fascicolo)
}
