package students

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ReligionChoice rappresenta l'opzione di avvalimento IRC dello studente.
type ReligionChoice string

const (
	ReligionChoiceAvvalente           ReligionChoice = "avvalente"
	ReligionChoiceNonAvvalente        ReligionChoice = "non_avvalente"
	ReligionChoiceAttivitaAlternativa ReligionChoice = "attivita_alternativa"
)

func (r ReligionChoice) IsValid() bool {
	switch r {
	case ReligionChoiceAvvalente, ReligionChoiceNonAvvalente, ReligionChoiceAttivitaAlternativa:
		return true
	default:
		return false
	}
}

// StudentReligionChoiceItem descrive la scelta di avvalimento con i dettagli dello studente.
type StudentReligionChoiceItem struct {
	ID               string         `json:"id"`
	StudentID        string         `json:"student_id"`
	UserID           string         `json:"user_id"`
	SchoolID         string         `json:"school_id"`
	Choice           ReligionChoice `json:"choice"`
	FirstName        string         `json:"first_name"`
	LastName         string         `json:"last_name"`
	ClassID          string         `json:"class_id"`
	ClassName        string         `json:"class_name"`
	EnrollmentNumber string         `json:"enrollment_number"`
	UpdatedBy        *string        `json:"updated_by,omitempty"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type SetReligionChoiceRequest struct {
	Choice ReligionChoice `json:"choice" binding:"required"`
}

type BatchSetReligionChoiceRequest struct {
	StudentIDs []string       `json:"student_ids" binding:"required"`
	Choice     ReligionChoice `json:"choice" binding:"required"`
}

// ReligionRepository gestisce l'accesso al database per le scelte di avvalimento.
type ReligionRepository struct {
	db *sql.DB
}

func NewReligionRepository(db *sql.DB) *ReligionRepository {
	return &ReligionRepository{db: db}
}

// resolveStudentProfileID recupera l'ID del record students(id) e school_id dato un studentID o userID.
func (r *ReligionRepository) resolveStudentProfileID(ctx context.Context, idOrUserID, fallbackSchoolID string) (studentProfileID, schoolID string, err error) {
	query := `
		SELECT s.id, s.school_id
		FROM students s
		WHERE s.id::text = $1 OR s.user_id::text = $1
		LIMIT 1
	`
	err = r.db.QueryRowContext(ctx, query, idOrUserID).Scan(&studentProfileID, &schoolID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if fallbackSchoolID != "" {
				return idOrUserID, fallbackSchoolID, nil
			}
			return "", "", fmt.Errorf("studente non trovato per id %s", idOrUserID)
		}
		return "", "", fmt.Errorf("errore lookup studente: %w", err)
	}
	return studentProfileID, schoolID, nil
}

// GetChoice recupera la scelta per un singolo studente.
// Se non è ancora presente un record esplicito, restituisce la scelta di default ('avvalente').
func (r *ReligionRepository) GetChoice(ctx context.Context, idOrUserID, callerSchoolID string) (*StudentReligionChoiceItem, error) {
	studentProfileID, schoolID, err := r.resolveStudentProfileID(ctx, idOrUserID, callerSchoolID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			COALESCE(src.id::text, ''),
			s.id::text,
			s.user_id::text,
			s.school_id::text,
			COALESCE(src.choice::text, 'avvalente'),
			COALESCE(u.first_name, ''),
			COALESCE(u.last_name, ''),
			COALESCE(c.id::text, ''),
			COALESCE(c.name || COALESCE(c.section, ''), ''),
			COALESCE(s.enrollment_number, ''),
			src.updated_by::text,
			COALESCE(src.updated_at, s.created_at)
		FROM students s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN classes c ON s.class_id = c.id
		LEFT JOIN student_religion_choices src
			ON src.student_id = s.id AND src.school_id = s.school_id
		WHERE s.id::text = $1
	`

	var item StudentReligionChoiceItem
	var rawChoice string
	var updatedBy sql.NullString

	err = r.db.QueryRowContext(ctx, query, studentProfileID).Scan(
		&item.ID,
		&item.StudentID,
		&item.UserID,
		&item.SchoolID,
		&rawChoice,
		&item.FirstName,
		&item.LastName,
		&item.ClassID,
		&item.ClassName,
		&item.EnrollmentNumber,
		&updatedBy,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &StudentReligionChoiceItem{
				StudentID: studentProfileID,
				SchoolID:  schoolID,
				Choice:    ReligionChoiceAvvalente,
				UpdatedAt: time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("query scelta religione fallita: %w", err)
	}

	item.Choice = ReligionChoice(rawChoice)
	if updatedBy.Valid {
		val := updatedBy.String
		item.UpdatedBy = &val
	}

	return &item, nil
}

// SetChoice inserisce o aggiorna la scelta di avvalimento (upsert).
func (r *ReligionRepository) SetChoice(ctx context.Context, idOrUserID, callerSchoolID, updatedBy string, choice ReligionChoice) (*StudentReligionChoiceItem, error) {
	if !choice.IsValid() {
		return nil, fmt.Errorf("scelta di avvalimento non valida: %s", choice)
	}

	studentProfileID, schoolID, err := r.resolveStudentProfileID(ctx, idOrUserID, callerSchoolID)
	if err != nil {
		return nil, err
	}

	var updatedByParam interface{}
	if updatedBy != "" {
		updatedByParam = updatedBy
	}

	upsertQuery := `
		INSERT INTO student_religion_choices (
			student_id, school_id, choice, updated_by, updated_at
		) VALUES (
			$1::uuid, $2::uuid, $3, $4::uuid, NOW()
		)
		ON CONFLICT (student_id, school_id)
		DO UPDATE SET
			choice = EXCLUDED.choice,
			updated_by = EXCLUDED.updated_by,
			updated_at = NOW()
	`

	_, err = r.db.ExecContext(ctx, upsertQuery, studentProfileID, schoolID, string(choice), updatedByParam)
	if err != nil {
		return nil, fmt.Errorf("upsert scelta religione fallita: %w", err)
	}

	return r.GetChoice(ctx, studentProfileID, schoolID)
}

// ListChoices recupera le scelte di avvalimento per una scuola con eventuale filtro per classe.
func (r *ReligionRepository) ListChoices(ctx context.Context, schoolID, classID string) ([]StudentReligionChoiceItem, error) {
	baseQuery := `
		SELECT
			COALESCE(src.id::text, ''),
			s.id::text,
			s.user_id::text,
			s.school_id::text,
			COALESCE(src.choice::text, 'avvalente'),
			COALESCE(u.first_name, ''),
			COALESCE(u.last_name, ''),
			COALESCE(c.id::text, ''),
			COALESCE(c.name || COALESCE(c.section, ''), ''),
			COALESCE(s.enrollment_number, ''),
			src.updated_by::text,
			COALESCE(src.updated_at, s.created_at)
		FROM students s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN classes c ON s.class_id = c.id
		LEFT JOIN student_religion_choices src
			ON src.student_id = s.id AND src.school_id = s.school_id
		WHERE s.school_id::text = $1
		  AND (s.deleted_at IS NULL)
	`
	var args []interface{}
	args = append(args, schoolID)

	if classID != "" {
		baseQuery += " AND s.class_id::text = $2"
		args = append(args, classID)
	}

	baseQuery += " ORDER BY c.name ASC NULLS LAST, c.section ASC NULLS LAST, u.last_name ASC, u.first_name ASC"

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query lista scelte religione fallita: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []StudentReligionChoiceItem
	for rows.Next() {
		var item StudentReligionChoiceItem
		var rawChoice string
		var updatedBy sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.StudentID,
			&item.UserID,
			&item.SchoolID,
			&rawChoice,
			&item.FirstName,
			&item.LastName,
			&item.ClassID,
			&item.ClassName,
			&item.EnrollmentNumber,
			&updatedBy,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan riga scelta religione fallito: %w", err)
		}

		item.Choice = ReligionChoice(rawChoice)
		if updatedBy.Valid {
			val := updatedBy.String
			item.UpdatedBy = &val
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterazione righe scelte fallita: %w", err)
	}

	return results, nil
}

// BatchSetChoices aggiorna la scelta di avvalimento per una lista di studenti.
func (r *ReligionRepository) BatchSetChoices(ctx context.Context, studentIDs []string, schoolID, updatedBy string, choice ReligionChoice) error {
	if !choice.IsValid() {
		return fmt.Errorf("scelta di avvalimento non valida: %s", choice)
	}
	if len(studentIDs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("transazione fallita: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO student_religion_choices (
			student_id, school_id, choice, updated_by, updated_at
		) VALUES (
			(SELECT id FROM students WHERE id::text = $1 OR user_id::text = $1 LIMIT 1),
			$2::uuid, $3, $4::uuid, NOW()
		)
		ON CONFLICT (student_id, school_id)
		DO UPDATE SET
			choice = EXCLUDED.choice,
			updated_by = EXCLUDED.updated_by,
			updated_at = NOW()
	`)
	if err != nil {
		return fmt.Errorf("preparazione statement batch fallita: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	var updatedByParam interface{}
	if updatedBy != "" {
		updatedByParam = updatedBy
	}

	for _, sID := range studentIDs {
		if sID == "" {
			continue
		}
		_, err := stmt.ExecContext(ctx, sID, schoolID, string(choice), updatedByParam)
		if err != nil {
			return fmt.Errorf("aggiornamento studente %s fallito: %w", sID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transazione batch fallito: %w", err)
	}

	return nil
}

// ReligionHandler gestisce le rotte HTTP per l'avvalimento IRC.
type ReligionHandler struct {
	repo *ReligionRepository
}

func NewReligionHandler(db *sql.DB) *ReligionHandler {
	return &ReligionHandler{repo: NewReligionRepository(db)}
}

func (h *ReligionHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/students/:id/religion-choice", h.GetStudentReligionChoice)
	r.PUT("/students/:id/religion-choice", h.SetStudentReligionChoice)
	r.GET("/secretary/religion-choices", h.ListReligionChoices)
	r.POST("/secretary/religion-choices/batch", h.BatchSetReligionChoices)
	r.GET("/classes/:id/religion-choices", h.GetClassReligionChoices)
}

func (h *ReligionHandler) GetStudentReligionChoice(c *gin.Context) {
	studentID := c.Param("id")
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")

	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	item, err := h.repo.GetChoice(c.Request.Context(), studentID, schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Permesso di visualizzazione: admin/segreteria/docente o studente stesso/genitore
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" &&
		actorRole != "principal" && actorRole != "vice_principal" && actorRole != "teacher" {
		if item.UserID != actorID && item.StudentID != actorID {
			c.JSON(http.StatusForbidden, gin.H{"error": "accesso non consentito"})
			return
		}
	}

	c.JSON(http.StatusOK, item)
}

func (h *ReligionHandler) SetStudentReligionChoice(c *gin.Context) {
	studentID := c.Param("id")
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")

	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Solo segreteria o amministratori scolastici possono aggiornare la scelta
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" &&
		actorRole != "principal" && actorRole != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo la segreteria o l'amministrazione possono modificare l'avvalimento IRC"})
		return
	}

	var req SetReligionChoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload non valido: " + err.Error()})
		return
	}

	if !req.Choice.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scelta non valida. Valori consentiti: 'avvalente', 'non_avvalente', 'attivita_alternativa'"})
		return
	}

	item, err := h.repo.SetChoice(c.Request.Context(), studentID, schoolID, actorID, req.Choice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ReligionHandler) ListReligionChoices(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")

	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" &&
		actorRole != "principal" && actorRole != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso riservato alla segreteria e dirigenza"})
		return
	}

	if targetSchool := c.Query("school_id"); targetSchool != "" && (actorRole == "superadmin" || targetSchool == schoolID) {
		schoolID = targetSchool
	}
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id mancante"})
		return
	}

	classID := c.Query("class_id")
	list, err := h.repo.ListChoices(c.Request.Context(), schoolID, classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []StudentReligionChoiceItem{}
	}

	c.JSON(http.StatusOK, list)
}

func (h *ReligionHandler) BatchSetReligionChoices(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")

	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" &&
		actorRole != "principal" && actorRole != "vice_principal" {
		c.JSON(http.StatusForbidden, gin.H{"error": "azione riservata alla segreteria"})
		return
	}

	var req BatchSetReligionChoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload non valido: " + err.Error()})
		return
	}

	if !req.Choice.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scelta non valida"})
		return
	}

	if len(req.StudentIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nessuno studente specificato"})
		return
	}

	err := h.repo.BatchSetChoices(c.Request.Context(), req.StudentIDs, schoolID, actorID, req.Choice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "scelte di avvalimento aggiornate con successo", "count": len(req.StudentIDs)})
}

func (h *ReligionHandler) GetClassReligionChoices(c *gin.Context) {
	classID := c.Param("id")
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	schoolID := c.GetString("school_id")

	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Docenti, segreteria, admin possono accedere alle scelte della classe
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "secretary" &&
		actorRole != "principal" && actorRole != "vice_principal" && actorRole != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "accesso non consentito"})
		return
	}

	list, err := h.repo.ListChoices(c.Request.Context(), schoolID, classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []StudentReligionChoiceItem{}
	}

	c.JSON(http.StatusOK, list)
}
