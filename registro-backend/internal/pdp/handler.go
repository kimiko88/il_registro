package pdp

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// contextKey for JWT claims passed by the middleware.
const (
	ContextKeyUserID   = "user_id"
	ContextKeyUserRole = "role"
	ContextKeySchoolID = "school_id"
)

func getRole(c *gin.Context) string {
	return c.GetString(ContextKeyUserRole)
}

// Handler exposes PDP/PEI HTTP endpoints.
type Handler struct {
	svc *Service
}

// NewHandler creates a new PDP handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes mounts all PDP routes onto a router group.
// Expected prefix: /api/v1  (caller passes the authenticated group)
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/pdp")
	{
		// List by student
		g.GET("/student/:studentId", h.GetByStudent)
		// List by class
		g.GET("/class/:classId", h.GetByClass)
		// Single plan
		g.GET("/:id", h.GetByID)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
		// Sharing & approval flow
		g.POST("/:id/share", h.Share)
		g.POST("/:id/approve", h.ApproveByFamily)
		// Vocabulary
		g.GET("/measures/compensative", h.ListCompensative)
		g.GET("/measures/dispensative", h.ListDispensative)
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func (h *Handler) GetByStudent(c *gin.Context) {
	actorID := c.GetString(ContextKeyUserID)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	actorRole := getRole(c)
	schoolID := c.GetString(ContextKeySchoolID)
	studentID := c.Param("studentId")
	year := c.DefaultQuery("year", "")

	plans, err := h.svc.GetByStudent(c.Request.Context(), actorID, actorRole, schoolID, studentID, year)
	if handleErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func (h *Handler) GetByClass(c *gin.Context) {
	actorID := c.GetString(ContextKeyUserID)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	actorRole := getRole(c)
	schoolID := c.GetString(ContextKeySchoolID)
	classID := c.Param("classId")
	year := c.DefaultQuery("year", "")

	plans, err := h.svc.GetByClass(c.Request.Context(), actorRole, schoolID, classID, year)
	if handleErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func (h *Handler) GetByID(c *gin.Context) {
	actorID := c.GetString(ContextKeyUserID)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	actorRole := getRole(c)
	schoolID := c.GetString(ContextKeySchoolID)
	plan, err := h.svc.GetByID(c.Request.Context(), actorID, actorRole, schoolID, c.Param("id"))
	if handleErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, plan)
}

func (h *Handler) Create(c *gin.Context) {
	actorID := c.GetString(ContextKeyUserID)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	actorRole := getRole(c)
	schoolID := c.GetString(ContextKeySchoolID)

	var req CreatePdpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.svc.CreatePlan(c.Request.Context(), actorID, actorRole, schoolID, &req)
	if handleErr(c, err) {
		return
	}
	c.JSON(http.StatusCreated, plan)
}

func (h *Handler) Update(c *gin.Context) {
	actorID := c.GetString(ContextKeyUserID)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	actorRole := getRole(c)
	schoolID := c.GetString(ContextKeySchoolID)
	var req UpdatePdpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plan, err := h.svc.UpdatePlan(c.Request.Context(), actorRole, schoolID, c.Param("id"), &req)
	if handleErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, plan)
}

func (h *Handler) Delete(c *gin.Context) {
	actorID := c.GetString(ContextKeyUserID)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	actorRole := getRole(c)
	schoolID := c.GetString(ContextKeySchoolID)
	if err := h.svc.DeletePlan(c.Request.Context(), actorRole, schoolID, c.Param("id")); handleErr(c, err) {
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) Share(c *gin.Context) {
	actorID := c.GetString(ContextKeyUserID)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	actorRole := getRole(c)
	schoolID := c.GetString(ContextKeySchoolID)
	var req ShareWithFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plan, err := h.svc.ShareWithFamily(c.Request.Context(), actorRole, schoolID, c.Param("id"), req.Share)
	if handleErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, plan)
}

func (h *Handler) ApproveByFamily(c *gin.Context) {
	actorID := c.GetString(ContextKeyUserID)
	actorRole := getRole(c)
	schoolID := c.GetString(ContextKeySchoolID)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if actorRole != "parent" && actorRole != "admin" && actorRole != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo i genitori o gli amministratori possono approvare i piani PDP"})
		return
	}
	if err := h.svc.ApproveByFamily(c.Request.Context(), actorRole, actorID, schoolID, c.Param("id")); handleErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Piano approvato dalla famiglia"})
}

func (h *Handler) ListCompensative(c *gin.Context) {
	actorID := c.GetString(ContextKeyUserID)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"measures": StandardCompensativeMeasures})
}

func (h *Handler) ListDispensative(c *gin.Context) {
	actorID := c.GetString(ContextKeyUserID)
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"measures": StandardDispensativeMeasures})
}

// ── Error Helper ──────────────────────────────────────────────────────────────

func handleErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, ErrPlanNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, ErrUnauthorized):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotGuardian):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotSharedYet):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, ErrAlreadyApproved):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
	return true
}
