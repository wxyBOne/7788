package services

import (
	"database/sql"
	"fmt"
	"seven-ai-backend/internal/models"
	"time"
)

type FriendshipService struct {
	db        *sql.DB
	aiService *AIService
}

func NewFriendshipService(db *sql.DB, aiService *AIService) *FriendshipService {
	return &FriendshipService{
		db:        db,
		aiService: aiService,
	}
}

// GetUserFriends 获取用户的好友列表
func (s *FriendshipService) GetUserFriends(userID int) ([]models.FriendInfo, error) {
	query := `
		SELECT 
			uf.id,
			pc.name,
			pc.avatar_url,
			pc.personality_signature,
			COALESCE(
				CASE 
					WHEN c.user_message != '' THEN c.user_message
					WHEN c.ai_response != '' THEN c.ai_response
					ELSE ''
				END, ''
			) as last_message,
			uf.last_message_at,
			uf.unread_count,
			true as is_online
		FROM user_friendships uf
		JOIN preset_characters pc ON uf.character_id = pc.id
		LEFT JOIN conversations c ON c.user_id = uf.user_id 
			AND c.character_id = uf.character_id 
			AND c.created_at = (
				SELECT MAX(created_at) 
				FROM conversations c2 
				WHERE c2.user_id = uf.user_id 
				AND c2.character_id = uf.character_id
			)
		WHERE uf.user_id = ? AND uf.is_active = true
		ORDER BY uf.last_message_at DESC, uf.created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user friends: %w", err)
	}
	defer rows.Close()

	var friends []models.FriendInfo
	for rows.Next() {
		var friend models.FriendInfo
		err := rows.Scan(
			&friend.ID,
			&friend.Name,
			&friend.AvatarURL,
			&friend.PersonalitySignature,
			&friend.LastMessage,
			&friend.LastMessageAt,
			&friend.UnreadCount,
			&friend.IsOnline,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan friend: %w", err)
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

// SearchAvailableCharacters 搜索可添加的角色
func (s *FriendshipService) SearchAvailableCharacters(userID int, keyword string) ([]models.AvailableCharacter, error) {
	searchPattern := "%" + keyword + "%"
	query := `
		SELECT 
			pc.id,
			pc.name,
			pc.description,
			pc.avatar_url,
			pc.personality_signature,
			CASE WHEN uf.id IS NOT NULL THEN true ELSE false END as is_added
		FROM preset_characters pc
		LEFT JOIN user_friendships uf ON pc.id = uf.character_id AND uf.user_id = ? AND uf.is_active = true
		WHERE pc.name LIKE ? OR pc.description LIKE ? OR pc.search_keywords LIKE ?
		ORDER BY pc.id
	`

	rows, err := s.db.Query(query, userID, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search characters: %w", err)
	}
	defer rows.Close()

	var characters []models.AvailableCharacter
	for rows.Next() {
		var char models.AvailableCharacter
		err := rows.Scan(
			&char.ID,
			&char.Name,
			&char.Description,
			&char.AvatarURL,
			&char.PersonalitySignature,
			&char.IsAdded,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan character: %w", err)
		}
		characters = append(characters, char)
	}

	return characters, nil
}

// AddFriend 添加好友
func (s *FriendshipService) AddFriend(userID int, characterID int) error {
	// 检查是否已经是好友
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM user_friendships 
		WHERE user_id = ? AND character_id = ? AND is_active = true
	`, userID, characterID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check friendship: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("already friends with this character")
	}

	// 添加好友关系
	_, err = s.db.Exec(`
		INSERT INTO user_friendships (user_id, character_id, is_active, created_at, updated_at)
		VALUES (?, ?, true, NOW(), NOW())
	`, userID, characterID)
	if err != nil {
		return fmt.Errorf("failed to add friend: %w", err)
	}

	// 生成AI欢迎消息
	err = s.generateWelcomeMessage(userID, characterID)
	if err != nil {
		// 记录错误但不影响添加好友
		fmt.Printf("Failed to generate welcome message: %v\n", err)
	}

	return nil
}

// RemoveFriend 移除好友
func (s *FriendshipService) RemoveFriend(userID int, characterID int) error {
	_, err := s.db.Exec(`
		UPDATE user_friendships 
		SET is_active = false, updated_at = NOW()
		WHERE user_id = ? AND character_id = ?
	`, userID, characterID)
	if err != nil {
		return fmt.Errorf("failed to remove friend: %w", err)
	}
	return nil
}

// generateWelcomeMessage 生成AI欢迎消息
func (s *FriendshipService) generateWelcomeMessage(userID int, characterID int) error {
	// 获取角色信息
	var character models.CharacterResponse
	err := s.db.QueryRow(`
		SELECT id, name, description, avatar_url, personality_signature,
		       personality_traits, background_story, voice_settings,
		       system_prompt, search_keywords
		FROM preset_characters WHERE id = ?
	`, characterID).Scan(
		&character.ID, &character.Name, &character.Description, &character.AvatarURL,
		&character.PersonalitySignature, &character.PersonalityTraits,
		&character.BackgroundStory, &character.VoiceSettings,
		&character.SystemPrompt, &character.SearchKeywords,
	)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}

	// 构建欢迎消息提示
	welcomePrompt := fmt.Sprintf(`
%s

现在用户刚刚添加你为好友，请以%s的身份说一句欢迎的话。要求：
1. 符合角色的性格特点
2. 简短自然，不超过50字
3. 体现角色的个性特征
4. 不要过于热情或冷淡，保持角色设定

请直接回复欢迎消息，不要包含任何解释。
`, character.SystemPrompt, character.Name)

	// 调用AI生成欢迎消息
	messages := []Message{
		{Role: "system", Content: welcomePrompt},
	}

	response, err := s.aiService.ChatWithLLM(messages, "qwen3-max", 0.8)
	if err != nil {
		return fmt.Errorf("failed to generate welcome message: %w", err)
	}

	// 保存欢迎消息到数据库
	sessionID := fmt.Sprintf("welcome_%d_%d_%d", userID, characterID, time.Now().Unix())
	_, err = s.db.Exec(`
		INSERT INTO conversations 
		(user_id, character_id, session_id, message_type, user_message, ai_response, is_ai_initiated, created_at)
		VALUES (?, ?, ?, 'text', '', ?, true, NOW())
	`, userID, characterID, sessionID, response)
	if err != nil {
		return fmt.Errorf("failed to save welcome message: %w", err)
	}

	// 更新好友关系的最后消息时间和未读数量
	_, err = s.db.Exec(`
		UPDATE user_friendships 
		SET last_message_at = NOW(), unread_count = unread_count + 1, updated_at = NOW()
		WHERE user_id = ? AND character_id = ?
	`, userID, characterID)
	if err != nil {
		return fmt.Errorf("failed to update friendship: %w", err)
	}

	return nil
}

// MarkMessagesAsRead 标记消息为已读
func (s *FriendshipService) MarkMessagesAsRead(userID int, characterID int) error {
	// 更新未读数量为0
	_, err := s.db.Exec(`
		UPDATE user_friendships 
		SET unread_count = 0, updated_at = NOW()
		WHERE user_id = ? AND character_id = ?
	`, userID, characterID)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	// 标记对话记录为已读
	_, err = s.db.Exec(`
		UPDATE conversations 
		SET is_read = true
		WHERE user_id = ? AND character_id = ? AND is_read = false
	`, userID, characterID)
	if err != nil {
		return fmt.Errorf("failed to mark conversations as read: %w", err)
	}

	return nil
}
