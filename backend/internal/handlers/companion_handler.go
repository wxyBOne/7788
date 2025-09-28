// Package handlers 处理HTTP请求的处理器
package handlers

import (
	"net/http"
	"seven-ai-backend/internal/models"
	"seven-ai-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CompanionHandler AI伙伴相关的HTTP请求处理器
type CompanionHandler struct {
	companionService *services.CompanionService
}

// NewCompanionHandler 创建AI伙伴处理器实例
func NewCompanionHandler(companionService *services.CompanionService) *CompanionHandler {
	return &CompanionHandler{companionService: companionService}
}

// CreateCompanion 创建AI伙伴
func (h *CompanionHandler) CreateCompanion(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req models.CreateCompanionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	companion, err := h.companionService.CreateCompanion(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "AI伙伴创建成功",
		"companion": companion,
	})
}

// GetUserCompanions 获取用户的AI伙伴列表
func (h *CompanionHandler) GetUserCompanions(c *gin.Context) {
	userID := c.GetInt("user_id")

	companions, err := h.companionService.GetUserCompanions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"companions": companions,
	})
}

// GetCompanion 获取单个AI伙伴信息
func (h *CompanionHandler) GetCompanion(c *gin.Context) {
	companionIDStr := c.Param("id")
	companionID, err := strconv.Atoi(companionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的AI伙伴ID"})
		return
	}

	companion, err := h.companionService.GetCompanion(companionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"companion": companion,
	})
}

// UpdateCompanion 更新AI伙伴信息
func (h *CompanionHandler) UpdateCompanion(c *gin.Context) {
	companionIDStr := c.Param("id")
	companionID, err := strconv.Atoi(companionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的AI伙伴ID"})
		return
	}

	var req models.UpdateCompanionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	err = h.companionService.UpdateCompanion(companionID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "AI伙伴信息更新成功",
	})
}


// GetGrowthStatus 获取AI伙伴成长状态
func (h *CompanionHandler) GetGrowthStatus(c *gin.Context) {
	companionIDStr := c.Param("id")
	companionID, err := strconv.Atoi(companionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的AI伙伴ID"})
		return
	}

	growth, err := h.companionService.GetGrowthStatus(companionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"growth": growth,
	})
}

// GetDiary 获取AI伙伴日记
func (h *CompanionHandler) GetDiary(c *gin.Context) {
	companionIDStr := c.Param("id")
	companionID, err := strconv.Atoi(companionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的AI伙伴ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	diaries, err := h.companionService.GetDiary(companionID, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"diaries": diaries,
	})
}

// GetEmotionState 获取AI伙伴情绪状态（用于粒子小球外观）
func (h *CompanionHandler) GetEmotionState(c *gin.Context) {
	companionIDStr := c.Param("id")
	companionID, err := strconv.Atoi(companionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的AI伙伴ID"})
		return
	}

	emotion, err := h.companionService.GetEmotionState(companionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"emotion": emotion,
	})
}
