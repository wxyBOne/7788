package handlers

import (
	"net/http"
	"seven-ai-backend/internal/models"
	"seven-ai-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConversationHandler struct {
	conversationService *services.ConversationService
}

func NewConversationHandler(conversationService *services.ConversationService) *ConversationHandler {
	return &ConversationHandler{conversationService: conversationService}
}

func (h *ConversationHandler) Chat(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.conversationService.Chat(userID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "聊天失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ConversationHandler) VoiceChat(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req models.VoiceChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.conversationService.VoiceChat(userID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "语音聊天失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ConversationHandler) ImageChat(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req models.ImageChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.conversationService.ImageChat(userID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片聊天失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ConversationHandler) GetHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	characterIDStr := c.Query("character_id")
	if characterIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色ID不能为空"})
		return
	}

	characterID, err := strconv.Atoi(characterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	history, err := h.conversationService.GetHistory(userID.(int), characterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取聊天记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    history,
	})
}

func (h *ConversationHandler) GetSessionHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	characterIDStr := c.Query("character_id")
	if characterIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色ID不能为空"})
		return
	}

	characterID, err := strconv.Atoi(characterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	history, err := h.conversationService.GetHistory(userID.(int), characterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取会话历史失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    history,
	})
}
