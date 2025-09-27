package handlers

import (
	"net/http"
	"seven-ai-backend/internal/models"
	"seven-ai-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FriendshipHandler struct {
	friendshipService *services.FriendshipService
}

func NewFriendshipHandler(friendshipService *services.FriendshipService) *FriendshipHandler {
	return &FriendshipHandler{
		friendshipService: friendshipService,
	}
}

// GetUserFriends 获取用户好友列表
func (h *FriendshipHandler) GetUserFriends(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	friends, err := h.friendshipService.GetUserFriends(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    friends,
	})
}

// SearchAvailableCharacters 搜索可添加的角色
func (h *FriendshipHandler) SearchAvailableCharacters(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	keyword := c.Query("keyword")
	// 允许空关键词，显示所有角色

	characters, err := h.friendshipService.SearchAvailableCharacters(userID.(int), keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    characters,
	})
}

// AddFriend 添加好友
func (h *FriendshipHandler) AddFriend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.AddFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.friendshipService.AddFriend(userID.(int), req.CharacterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Friend added successfully",
	})
}

// RemoveFriend 移除好友
func (h *FriendshipHandler) RemoveFriend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	characterIDStr := c.Param("character_id")
	characterID, err := strconv.Atoi(characterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	err = h.friendshipService.RemoveFriend(userID.(int), characterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Friend removed successfully",
	})
}
