package handlers

import (
	"net/http"
	"seven-ai-backend/internal/models"
	"seven-ai-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VoiceCallHandler struct {
	voiceCallService *services.VoiceCallService
}

func NewVoiceCallHandler(voiceCallService *services.VoiceCallService) *VoiceCallHandler {
	return &VoiceCallHandler{
		voiceCallService: voiceCallService,
	}
}

// ProcessVoiceCall 处理语音通话
func (h *VoiceCallHandler) ProcessVoiceCall(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.VoiceCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.voiceCallService.ProcessVoiceCall(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetVoiceCallHistory 获取语音通话历史
func (h *VoiceCallHandler) GetVoiceCallHistory(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	characterIDStr := c.Query("character_id")
	characterID, err := strconv.Atoi(characterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	history, err := h.voiceCallService.GetVoiceCallHistory(userID, characterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    history,
	})
}
