package communications

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/communications")
	{
		g.GET("", h.List)
		g.GET("/bacheca", h.ListBacheca)
		g.GET("/:id", h.GetByID)
		g.PUT("/:id", h.Update)
		g.POST("", h.Send)
		g.DELETE("/:id", h.Delete)
		g.POST("/:id/sign", h.Sign)
		g.POST("/:id/read", h.MarkAsRead)
		g.GET("/:id/signatures", h.GetSignatures)
		g.GET("/:id/signature-report", h.GetSignatureReport)
		g.GET("/:id/unread-users", h.GetUnreadUsers)
	}
}

func (h *Handler) List(c *gin.Context) {
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	msgs, err := h.service.ListMessages(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, msgs)
}

func (h *Handler) ListBacheca(c *gin.Context) {
	uid := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	msgs, err := h.service.ListBacheca(c.Request.Context(), schoolID, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, msgs)
}

func (h *Handler) Send(c *gin.Context) {
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg, err := h.service.SendMessage(c.Request.Context(), uid, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, msg)
}

// Delete verifies that the caller is the author of the message before deleting.
func (h *Handler) Delete(c *gin.Context) {
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.GetString("role")
	id := c.Param("id")
	if err := h.service.DeleteMessage(c.Request.Context(), uid, role, id); err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) Sign(c *gin.Context) {
	id := c.Param("id")
	uid := c.GetString("user_id")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ipAddress := c.ClientIP()
	if err := h.service.SignMessageWithIP(c.Request.Context(), id, uid, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "signed"})
}

func (h *Handler) GetSignatures(c *gin.Context) {
	id := c.Param("id")
	names, err := h.service.GetMessageSignatures(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, names)
}

func (h *Handler) GetSignatureReport(c *gin.Context) {
	id := c.Param("id")
	role := c.GetString("role")
	report, err := h.service.GetSignatureReport(c.Request.Context(), role, id)
	if err != nil {
		if err.Error() == "forbidden: signature reports are restricted to staff" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	msg, err := h.service.GetMessageByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, msg)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	uid := c.GetString("user_id")
	role := c.GetString("role")
	var req struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateMessage(c.Request.Context(), uid, role, id, req.Subject, req.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	uid := c.GetString("user_id")
	ipAddress := c.ClientIP()
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.MarkAsRead(c.Request.Context(), id, uid, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *Handler) GetUnreadUsers(c *gin.Context) {
	id := c.Param("id")
	unread, err := h.service.GetUnreadUsers(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, unread)
}
