package activity

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type activityRequest struct {
	Name string	`json:"name" binding:"required, max=255"`
	Notes *string `json:"notes"`
	Day string `json:"day" binding:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime string `json:"start_time" bidning:"required"`
	EndTime string `json:"end_time" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req activityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	a := &Activity{
		Name: req.Name,
		Notes: req.Notes,
		Day: req.Day,
		StartTime: req.StartTime,
		EndTime: req.EndTime,
	}

	if err := h.service.Create(c.Request.Context(), a); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data":a})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid id"})
		return
	}

	var req activityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	a := &Activity{
		ID: id,
		Name: req.Name,
		Notes: req.Notes,
		Day: req.Day,
		StartTime: req.StartTime,
		EndTime: req.EndTime,
	}

	if err := h.service.Update(c.Request.Context(), a); err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data":a})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete activity"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListByDay(c *gin.Context) {
	day := c.Param("day")
	activities, err := h.service.ListByDay(c.Request.Context(), day)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list activities"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data":activities})
}

func respondServiceError(c *gin.Context, err error) {
	if errors.Is(err, ErrOverlap) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save activity"})
}