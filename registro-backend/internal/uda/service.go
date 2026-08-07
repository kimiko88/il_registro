package uda

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"net/http"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUda(ctx context.Context, schoolID, teacherID string, req CreateUdaRequest) (*UdaPlan, error) {
	var start, end *time.Time
	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			start = &t
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			end = &t
		}
	}
	if req.Period == "" {
		req.Period = "annuale"
	}
	if req.Competencies == nil {
		req.Competencies = []string{}
	}

	plan := &UdaPlan{
		SchoolID:           schoolID,
		ClassID:            req.ClassID,
		SubjectID:          req.SubjectID,
		TeacherID:          teacherID,
		Title:              req.Title,
		Description:        req.Description,
		Period:             req.Period,
		StartDate:          start,
		EndDate:            end,
		Competencies:       req.Competencies,
		Objectives:         req.Objectives,
		Methodologies:      req.Methodologies,
		EvaluationCriteria: req.EvaluationCriteria,
		Status:             "draft",
	}

	if err := s.repo.Create(ctx, plan); err != nil {
		return nil, fmt.Errorf("failed to create UDA plan: %w", err)
	}
	return plan, nil
}

func (s *Service) ListAll(ctx context.Context) ([]*UdaPlan, error) {
	return s.repo.ListAll(ctx)
}

func (s *Service) ListByClass(ctx context.Context, classID string) ([]*UdaPlan, error) {
	return s.repo.ListByClass(ctx, classID)
}

func (s *Service) UpdateUda(ctx context.Context, id string, req UpdateUdaRequest) (*UdaPlan, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *Service) DeleteUda(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// Handler HTTP
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	uda := rg.Group("/uda")
	{
		uda.GET("", h.ListAll)
		uda.GET("/class/:classID", h.ListByClass)
		uda.POST("", h.CreateUda)
		uda.PUT("/:id", h.UpdateUda)
		uda.DELETE("/:id", h.DeleteUda)
	}
}

func (h *Handler) ListAll(c *gin.Context) {
	plans, err := h.service.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *Handler) ListByClass(c *gin.Context) {
	classID := c.Param("classID")
	plans, err := h.service.ListByClass(c.Request.Context(), classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *Handler) CreateUda(c *gin.Context) {
	teacherID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateUdaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.service.CreateUda(c.Request.Context(), schoolID, teacherID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, plan)
}

func (h *Handler) UpdateUda(c *gin.Context) {
	id := c.Param("id")
	var req UpdateUdaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.service.UpdateUda(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plan)
}

func (h *Handler) DeleteUda(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteUda(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "uda deleted successfully"})
}
