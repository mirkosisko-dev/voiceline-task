package handler

import (
	"github.com/gin-gonic/gin"
)

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type HealthHandler struct{}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "voiceline-api",
	})
}
