package auth

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	adminPassword string
	apiToken string
}

func NewHandler(adminPassword, apiToken string) *Handler {
	return &Handler{adminPassword: adminPassword, apiToken: apiToken}
}

type loginRequest struct {
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if subtle.ConstantTimeCompare([]byte(req.Password), []byte(h.adminPassword)) != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": h.apiToken})
}