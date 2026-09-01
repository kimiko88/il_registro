package competencies

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
)

type Evaluation struct {
	ID             string    `json:"id" db:"id"`
	SchoolID       string    `json:"school_id" db:"school_id"`
	StudentID      string    `json:"student_id" db:"student_id"`
	ClassID        string    `json:"class_id" db:"class_id"`
	SubjectID      *string   `json:"subject_id,omitempty" db:"subject_id"`
	SubjectName    string    `json:"subject_name,omitempty" db:"subject_name"`
	EvaluatorID    string    `json:"evaluator_id" db:"evaluator_id"`
	EvaluatorName  string    `json:"evaluator_name,omitempty" db:"evaluator_name"`
	Semester       int       `json:"semester" db:"semester"`
	CompetenceCode string    `json:"competence_code" db:"competence_code"`
	CompetenceName string    `json:"competence_name" db:"competence_name"`
	Level          string    `json:"level" db:"level"` // A_Avanzato, B_Intermedio, C_Base, D_Iniziale
	Descriptor     string    `json:"descriptor" db:"descriptor"`
	Notes          string    `json:"notes" db:"notes"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type SaveEvaluationRequest struct {
	StudentID      string  `json:"student_id" binding:"required"`
	ClassID        string  `json:"class_id" binding:"required"`
	SubjectID      *string `json:"subject_id"`
	Semester       int     `json:"semester"`
	CompetenceCode string  `json:"competence_code" binding:"required"`
	CompetenceName string  `json:"competence_name" binding:"required"`
	Level          string  `json:"level" binding:"required"`
	Descriptor     string  `json:"descriptor"`
	Notes          string  `json:"notes"`
}

type StudentCompetencyEvaluation struct {
	StudentID   string            `json:"student_id"`
	StudentName string            `json:"student_name"`
	Evaluations map[string]string `json:"evaluations"`
}

type BatchSaveRequest struct {
	ClassID     string                        `json:"class_id"`
	SubjectID   string                        `json:"subject_id"`
	Period      string                        `json:"period"`
	Evaluations []StudentCompetencyEvaluation `json:"evaluations"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(ctx context.Context, schoolID, evaluatorID string, req SaveEvaluationRequest) (*Evaluation, error) {
	if req.Semester <= 0 {
		req.Semester = 1
	}
	query := `
		INSERT INTO competence_evaluations (school_id, student_id, class_id, subject_id, evaluator_id, semester, competence_code, competence_name, level, descriptor, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		ON CONFLICT (student_id, competence_code, semester) DO UPDATE
		SET level = EXCLUDED.level,
		    descriptor = EXCLUDED.descriptor,
		    notes = EXCLUDED.notes,
		    updated_at = NOW()
		RETURNING id, created_at, updated_at`

	eval := &Evaluation{
		SchoolID:       schoolID,
		StudentID:      req.StudentID,
		ClassID:        req.ClassID,
		SubjectID:      req.SubjectID,
		EvaluatorID:    evaluatorID,
		Semester:       req.Semester,
		CompetenceCode: req.CompetenceCode,
		CompetenceName: req.CompetenceName,
		Level:          req.Level,
		Descriptor:     req.Descriptor,
		Notes:          req.Notes,
	}

	err := r.db.QueryRowContext(ctx, query,
		schoolID, req.StudentID, req.ClassID, req.SubjectID, evaluatorID,
		req.Semester, req.CompetenceCode, req.CompetenceName, req.Level, req.Descriptor, req.Notes,
	).Scan(&eval.ID, &eval.CreatedAt, &eval.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return eval, nil
}

func (r *Repository) GetByStudent(ctx context.Context, studentID string, semester int) ([]*Evaluation, error) {
	query := `
		SELECT ce.id, ce.school_id, ce.student_id, ce.class_id, ce.subject_id,
		       COALESCE(s.name, '') AS subject_name, ce.evaluator_id,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS evaluator_name,
		       ce.semester, ce.competence_code, ce.competence_name, ce.level,
		       COALESCE(ce.descriptor, ''), COALESCE(ce.notes, ''), ce.created_at, ce.updated_at
		FROM competence_evaluations ce
		LEFT JOIN subjects s ON ce.subject_id = s.id
		LEFT JOIN users u ON ce.evaluator_id = u.id
		WHERE ce.student_id = $1::uuid AND ($2 = 0 OR ce.semester = $2)
		ORDER BY ce.competence_code ASC`

	rows, err := r.db.QueryContext(ctx, query, studentID, semester)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Evaluation
	for rows.Next() {
		e := &Evaluation{}
		if err := rows.Scan(
			&e.ID, &e.SchoolID, &e.StudentID, &e.ClassID, &e.SubjectID,
			&e.SubjectName, &e.EvaluatorID, &e.EvaluatorName,
			&e.Semester, &e.CompetenceCode, &e.CompetenceName, &e.Level,
			&e.Descriptor, &e.Notes, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *Repository) GetByClass(ctx context.Context, classID, subjectID string) ([]StudentCompetencyEvaluation, error) {
	query := `
		SELECT s.id, COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE ($1 = '' OR s.class_id = $1::uuid)
		ORDER BY u.last_name ASC, u.first_name ASC`

	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []StudentCompetencyEvaluation
	for rows.Next() {
		var sc StudentCompetencyEvaluation
		if err := rows.Scan(&sc.StudentID, &sc.StudentName); err != nil {
			return nil, err
		}
		sc.Evaluations = make(map[string]string)
		result = append(result, sc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if result == nil {
		return []StudentCompetencyEvaluation{}, nil
	}

	if len(result) > 0 {
		evalQuery := `
			SELECT student_id, competence_code, level
			FROM competence_evaluations
			WHERE ($1 = '' OR class_id = $1::uuid) AND ($2 = '' OR subject_id::text = $2)`
		evalRows, err := r.db.QueryContext(ctx, evalQuery, classID, subjectID)
		if err != nil {
			return nil, err
		}
		defer func() { _ = evalRows.Close() }()
		evalMap := make(map[string]map[string]string)
		for evalRows.Next() {
			var stID, code, lvl string
			if err := evalRows.Scan(&stID, &code, &lvl); err != nil {
				return nil, err
			}
			if evalMap[stID] == nil {
				evalMap[stID] = make(map[string]string)
			}
			evalMap[stID][code] = lvl
		}
		if err := evalRows.Err(); err != nil {
			return nil, err
		}
		for i := range result {
			if m, ok := evalMap[result[i].StudentID]; ok {
				result[i].Evaluations = m
			}
		}
	}

	return result, nil
}

type Service struct {
	repo     *Repository
	userRepo users.Repository
}

func NewService(repo *Repository, uRepo ...users.Repository) *Service {
	svc := &Service{repo: repo}
	if len(uRepo) > 0 && uRepo[0] != nil {
		svc.userRepo = uRepo[0]
	}
	return svc
}

func (s *Service) GetUserRepo() users.Repository {
	return s.userRepo
}

func IsValidCompetencyLevel(level string) bool {
	switch level {
	case "A_Avanzato", "B_Intermedio", "C_Base", "D_Iniziale", "A", "B", "C", "D", "beginner", "intermediate", "advanced", "expert":
		return true
	default:
		return false
	}
}

func (s *Service) SaveEvaluation(ctx context.Context, schoolID, evaluatorID string, req SaveEvaluationRequest) (*Evaluation, error) {
	if req.StudentID == "" || req.ClassID == "" || req.CompetenceCode == "" {
		return nil, fmt.Errorf("student_id, class_id, and competence_code are required")
	}
	if !IsValidCompetencyLevel(req.Level) {
		return nil, fmt.Errorf("invalid competency level: %s", req.Level)
	}
	return s.repo.Upsert(ctx, schoolID, evaluatorID, req)
}

func (s *Service) GetStudentEvaluations(ctx context.Context, studentID string, semester int) ([]*Evaluation, error) {
	if s.repo == nil {
		return []*Evaluation{}, nil
	}
	return s.repo.GetByStudent(ctx, studentID, semester)
}

func (s *Service) GetClassEvaluations(ctx context.Context, classID, subjectID string) ([]StudentCompetencyEvaluation, error) {
	if s.repo == nil {
		return []StudentCompetencyEvaluation{}, nil
	}
	return s.repo.GetByClass(ctx, classID, subjectID)
}

// Handler HTTP
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	comp := rg.Group("/competencies")
	{
		comp.GET("", h.GetClassEvaluations)
		comp.GET("/student/:studentID", h.GetStudentEvaluations)
		comp.POST("/evaluations", h.SaveEvaluation)
		comp.POST("/batch", h.BatchSave)
	}
}

func (h *Handler) GetClassEvaluations(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "secretary" && role != "coordinator" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: class competency evaluations restricted to teaching staff and administration"})
		return
	}
	classID := c.Query("class_id")
	subjectID := c.Query("subject_id")
	evals, err := h.service.GetClassEvaluations(c.Request.Context(), classID, subjectID)
	if err != nil {
		c.JSON(http.StatusOK, []StudentCompetencyEvaluation{})
		return
	}
	c.JSON(http.StatusOK, evals)
}

func (h *Handler) GetStudentEvaluations(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	studentID := c.Param("studentID")
	if role == "student" && userID != studentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: cannot view other students evaluations"})
		return
	}
	if role == "parent" {
		if uRepo := h.service.GetUserRepo(); uRepo != nil {
			isG, err := uRepo.IsGuardian(c.Request.Context(), userID, studentID)
			if err != nil || !isG {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: non sei il tutore legale di questo studente"})
				return
			}
		}
	}
	sem := 0
	evals, err := h.service.GetStudentEvaluations(c.Request.Context(), studentID, sem)
	if err != nil {
		c.JSON(http.StatusOK, []*Evaluation{})
		return
	}
	c.JSON(http.StatusOK, evals)
}

func (h *Handler) SaveEvaluation(c *gin.Context) {
	evaluatorID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	if evaluatorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req SaveEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eval, err := h.service.SaveEvaluation(c.Request.Context(), schoolID, evaluatorID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, eval)
}

func (h *Handler) BatchSave(c *gin.Context) {
	evaluatorID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	if evaluatorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req BatchSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, studentEval := range req.Evaluations {
		for code, level := range studentEval.Evaluations {
			if !IsValidCompetencyLevel(level) {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("livello non valido per competenza %s: %s", code, level)})
				return
			}
			var subj *string
			if req.SubjectID != "" {
				subj = &req.SubjectID
			}
			if _, err := h.service.SaveEvaluation(c.Request.Context(), schoolID, evaluatorID, SaveEvaluationRequest{
				StudentID:      studentEval.StudentID,
				ClassID:        req.ClassID,
				SubjectID:      subj,
				CompetenceCode: code,
				CompetenceName: code,
				Level:          level,
			}); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("errore nel salvataggio della valutazione per studente %s: %v", studentEval.StudentID, err)})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "valutazioni per competenze salvate con successo"})
}
