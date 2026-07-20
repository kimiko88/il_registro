package notes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if userID == "" || schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.service.CreateNote(c.Request.Context(), userID, schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, note)
}

func (h *Handler) Approve(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	noteID := c.Param("id")

	if err := h.service.ApproveNote(c.Request.Context(), actorID, actorRole, noteID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "note approved successfully"})
}

func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	noteID := c.Param("id")

	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.service.UpdateNote(c.Request.Context(), userID, noteID, req)
	if err != nil {
		if err.Error() == "unauthorized: can only edit own notes" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, note)
}

func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	noteID := c.Param("id")

	if err := h.service.DeleteNote(c.Request.Context(), userID, noteID); err != nil {
		if err.Error() == "unauthorized: can only delete own notes" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) List(c *gin.Context) {
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	filter := NoteFilter{
		StudentID: c.Query("student_id"),
		ClassID:   c.Query("class_id"),
		TeacherID: c.Query("teacher_id"),
		Type:      NoteType(c.Query("type")),
		DateFrom:  c.Query("date_from"),
		DateTo:    c.Query("date_to"),
		ActorID:   actorID,
		ActorRole: actorRole,
	}

	notes, err := h.service.ListNotes(c.Request.Context(), filter)
	if err != nil {
		if err.Error() == "unauthorized: not a guardian of this student" || err.Error() == "unauthorized: parent must specify student_id" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notes)
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/notes")
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.POST("/:id/approve", h.Approve)
		group.PATCH("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
	}
}
