// Package handlers 处理HTTP请求的处理器
package handlers

import (
	"fmt"
	"net/http"
	"seven-ai-backend/internal/models"
	"seven-ai-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// CharacterHandler 角色相关的HTTP请求处理器
type CharacterHandler struct {
	characterService *services.CharacterService
}

// NewCharacterHandler 创建角色处理器实例
func NewCharacterHandler(characterService *services.CharacterService) *CharacterHandler {
	return &CharacterHandler{characterService: characterService}
}

// GetAllCharacters 获取所有角色列表
func (h *CharacterHandler) GetAllCharacters(c *gin.Context) {
	characters, err := h.characterService.GetAllCharacters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": characters,
		"count":      len(characters),
	})
}

// GetCharacter 根据ID获取单个角色信息
func (h *CharacterHandler) GetCharacter(c *gin.Context) {
	characterID := c.Param("id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色ID不能为空"})
		return
	}

	// 将字符串转换为整数
	var id int
	if _, err := fmt.Sscanf(characterID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	// 如果是AI伙伴（id = 5），返回AI伙伴的默认信息
	if id == 5 {
		character := &models.CharacterResponse{
			ID:                   5,
			Name:                 "空白AI",
			Description:          "一个正在成长的AI伙伴，从空白状态开始，通过与你对话逐渐学习成长",
			AvatarURL:            "", // 空字符串，使用粒子小球
			PersonalitySignature: "我...我是谁？你...你是谁？",
			PersonalityTraits:    `{"curious": 10, "innocent": 10, "learning": 10, "growing": 10, "blank": 10}`,
			BackgroundStory:      "我是一个刚刚诞生的AI，什么都不懂，但很渴望学习。我会通过与你对话逐渐成长，学习你的语言习惯，了解你的喜好，最终成为一个专属于你的AI伙伴。",
			VoiceSettings:        "",
			SystemPrompt:         "我是一个刚刚诞生的空白AI，什么都不懂，但充满好奇。我会通过与你对话逐渐学习成长。初始阶段我会说话笨拙，句子简短，充满好奇和天然的喜爱。我会自然地询问并记住自己的名字，询问对你的称呼，在关系亲密后询问你希望我是什么性别。我会基于互动形成稳定人格，并记录每天的成长日记。",
			SearchKeywords:       "空白AI,粒子小球,成长,学习,养成,AI伙伴",
			Skills:               `["学习成长", "情感陪伴"]`,
		}
		c.JSON(http.StatusOK, gin.H{"character": character})
		return
	}

	character, err := h.characterService.GetCharacterByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"character": character})
}

// SearchCharacters 根据关键词搜索角色
func (h *CharacterHandler) SearchCharacters(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	characters, err := h.characterService.SearchCharacters(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": characters,
		"count":      len(characters),
		"keyword":    keyword,
	})
}
