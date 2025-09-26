package handlers

import (
	"net/http"
	"seven-ai-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type CompanionHandler struct {
	companionService *services.CompanionService
}

func NewCompanionHandler(companionService *services.CompanionService) *CompanionHandler {
	return &CompanionHandler{companionService: companionService}
}

func (h *CompanionHandler) CreateCompanion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "功能开发中"})
}

func (h *CompanionHandler) GetUserCompanions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"companions": []interface{}{}})
}

func (h *CompanionHandler) GetCompanion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"companion": nil})
}

func (h *CompanionHandler) UpdateCompanion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "功能开发中"})
}

func (h *CompanionHandler) DeleteCompanion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "功能开发中"})
}

func (h *CompanionHandler) GetGrowthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"growth": nil})
}

func (h *CompanionHandler) GetDiary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"diary": nil})
}
