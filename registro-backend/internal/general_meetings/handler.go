package general_meetings

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	gm := r.Group("/general-meetings")
	{
		gm.GET("", h.List)
		gm.POST("", h.Create)
		gm.GET("/:id", h.GetByID)
		gm.DELETE("/:id", h.Delete)

		gm.POST("/:id/register", h.Register)
		gm.POST("/:id/unregister", h.Unregister)
		gm.GET("/:id/registrations", h.ListRegistrations)
	}
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	meetings, err := h.service.ListMeetings(c.Request.Context(), schoolID, userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if meetings == nil {
		meetings = []*GeneralMeeting{}
	}
	c.JSON(http.StatusOK, meetings)
}

func (h *Handler) Create(c *gin.Context) {
	actorID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateGeneralMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := h.service.CreateMeeting(c.Request.Context(), actorID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (h *Handler) GetByID(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")

	m, err := h.service.GetMeeting(c.Request.Context(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "meeting not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *Handler) Delete(c *gin.Context) {
	actorID := c.GetString("user_id")
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id := c.Param("id")
	if err := h.service.DeleteMeeting(c.Request.Context(), id, actorID, role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Register(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	meetingID := c.Param("id")

	if err := h.service.RegisterUser(c.Request.Context(), meetingID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registered successfully"})
}

func (h *Handler) Unregister(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	meetingID := c.Param("id")

	if err := h.service.UnregisterUser(c.Request.Context(), meetingID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unregistered successfully"})
}

func (h *Handler) ListRegistrations(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" && role != "secretary" && role != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	meetingID := c.Param("id")

	regs, err := h.service.ListRegistrations(c.Request.Context(), meetingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if regs == nil {
		regs = []*GeneralMeetingRegistration{}
	}
	c.JSON(http.StatusOK, regs)
}
