package handlers

import (
	"fmt"
	"net/http"
	"seven-ai-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type CharacterHandler struct {
	characterService *services.CharacterService
}

func NewCharacterHandler(characterService *services.CharacterService) *CharacterHandler {
	return &CharacterHandler{characterService: characterService}
}

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

func (h *CharacterHandler) GetCharacter(c *gin.Context) {
	characterID := c.Param("id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色ID不能为空"})
		return
	}

	// 这里需要将字符串转换为整数，简化处理
	var id int
	if _, err := fmt.Sscanf(characterID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	character, err := h.characterService.GetCharacterByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"character": character})
}

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
