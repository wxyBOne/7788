package services

import (
	"database/sql"
	"fmt"
	"seven-ai-backend/internal/models"
	"strings"
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

// GetUserFriends 获取用户的好友列表（包括AI伙伴）
func (s *FriendshipService) GetUserFriends(userID int) ([]models.FriendInfo, error) {
	var friends []models.FriendInfo

	// 获取普通角色好友
	query := `
		SELECT 
			uf.id,
			uf.character_id,
			pc.name,
			pc.avatar_url,
			pc.personality_signature,
			COALESCE(
				CASE 
					WHEN c.ai_response != '' THEN c.ai_response
					WHEN c.user_message != '' THEN c.user_message
					ELSE ''
				END, ''
			) as last_message,
			uf.last_message_at,
			true as is_online,
			'character' as type
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

	for rows.Next() {
		var friend models.FriendInfo
		err := rows.Scan(
			&friend.ID,
			&friend.CharacterID,
			&friend.Name,
			&friend.AvatarURL,
			&friend.PersonalitySignature,
			&friend.LastMessage,
			&friend.LastMessageAt,
			&friend.IsOnline,
			&friend.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan friend: %w", err)
		}
		friends = append(friends, friend)
	}

	// 获取AI伙伴
	aiQuery := `
		SELECT 
			ac.id,
			5 as character_id, -- AI伙伴使用5作为特殊character_id
			ac.name,
			ac.avatar_url,
			ac.personality_signature,
			COALESCE(
				CASE 
					WHEN c.ai_response != '' THEN c.ai_response
					WHEN c.user_message != '' THEN c.user_message
					ELSE ''
				END, ''
			) as last_message,
			ac.last_active_at as last_message_at,
			true as is_online,
			'companion' as type,
			ac.growth_percentage,
			ac.current_level,
			ac.total_experience
		FROM ai_companions ac
		LEFT JOIN conversations c ON c.user_id = ac.user_id 
			AND c.character_id = 5 
			AND c.created_at = (
				SELECT MAX(created_at) 
				FROM conversations c2 
				WHERE c2.user_id = ac.user_id 
				AND c2.character_id = 5
			)
		WHERE ac.user_id = ?
		ORDER BY ac.last_active_at DESC
	`

	aiRows, err := s.db.Query(aiQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("查询AI伙伴失败: %v", err)
	}
	defer aiRows.Close()

	for aiRows.Next() {
		var friend models.FriendInfo
		var growthPercentage sql.NullFloat64
		var currentLevel sql.NullInt64
		var totalExperience sql.NullInt64

		err := aiRows.Scan(
			&friend.ID,
			&friend.CharacterID,
			&friend.Name,
			&friend.AvatarURL,
			&friend.PersonalitySignature,
			&friend.LastMessage,
			&friend.LastMessageAt,
			&friend.IsOnline,
			&friend.Type,
			&growthPercentage,
			&currentLevel,
			&totalExperience,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描AI伙伴数据失败: %v", err)
		}

		// 设置AI伙伴特有的字段
		if growthPercentage.Valid {
			friend.GrowthPercentage = growthPercentage.Float64
		}
		if currentLevel.Valid {
			friend.CurrentLevel = int(currentLevel.Int64)
		}
		if totalExperience.Valid {
			friend.TotalExperience = int(totalExperience.Int64)
		}

		friends = append(friends, friend)
	}

	// 添加空白AI到好友列表（如果用户还没有AI伙伴）
	var hasCompanion bool
	err = s.db.QueryRow("SELECT COUNT(*) > 0 FROM ai_companions WHERE user_id = ?", userID).Scan(&hasCompanion)
	if err != nil {
		return nil, fmt.Errorf("failed to check companion: %w", err)
	}

	if !hasCompanion {
		// 用户还没有AI伙伴，添加空白AI到好友列表
		blankAI := models.FriendInfo{
			ID:                   0,
			CharacterID:          5,
			Name:                 "空白AI",
			AvatarURL:            "",
			PersonalitySignature: "我...我是谁？你...你是谁？",
			LastMessage:          "点击创建你的专属AI伙伴",
			LastMessageAt:        &time.Time{},
			IsOnline:             true,
			Type:                 "blank", // 特殊类型标识
		}
		friends = append(friends, blankAI)
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
			pc.search_keywords,
			pc.avatar_url,
			pc.personality_signature
		FROM preset_characters pc
		LEFT JOIN user_friendships uf ON pc.id = uf.character_id AND uf.user_id = ? AND uf.is_active = true
		WHERE uf.id IS NULL
		AND (pc.name LIKE ? OR pc.description LIKE ? OR pc.search_keywords LIKE ?)
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
		var searchKeywords string
		err := rows.Scan(
			&char.ID,
			&char.Name,
			&searchKeywords,
			&char.AvatarURL,
			&char.PersonalitySignature,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan character: %w", err)
		}

		// 从search_keywords中提取第二段作为描述
		char.Description = extractSecondKeyword(searchKeywords)
		char.IsAdded = false // 只返回未添加的角色

		characters = append(characters, char)
	}

	return characters, nil
}

// extractSecondKeyword 从search_keywords中提取第二段
func extractSecondKeyword(searchKeywords string) string {
	if searchKeywords == "" {
		return ""
	}

	// 按逗号分割
	keywords := strings.Split(searchKeywords, ",")
	if len(keywords) < 2 {
		return keywords[0] // 如果只有一段，返回第一段
	}

	// 返回第二段，去除首尾空格
	return strings.TrimSpace(keywords[1])
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
	var voiceSettings sql.NullString
	err := s.db.QueryRow(`
		SELECT id, name, description, avatar_url, personality_signature,
		       personality_traits, background_story, voice_settings,
		       system_prompt, search_keywords
		FROM preset_characters WHERE id = ?
	`, characterID).Scan(
		&character.ID, &character.Name, &character.Description, &character.AvatarURL,
		&character.PersonalitySignature, &character.PersonalityTraits,
		&character.BackgroundStory, &voiceSettings,
		&character.SystemPrompt, &character.SearchKeywords,
	)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}

	// 处理NULL值
	if voiceSettings.Valid {
		character.VoiceSettings = voiceSettings.String
	} else {
		character.VoiceSettings = ""
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

	response, err := s.aiService.ChatWithLLM(messages, "qwen3-max", 0.8, "text")
	if err != nil {
		return fmt.Errorf("failed to generate welcome message: %w", err)
	}

	// 后处理AI响应，移除角色名字前缀
	if strings.HasPrefix(response, character.Name+"。") {
		response = strings.TrimPrefix(response, character.Name+"。")
	} else if strings.HasPrefix(response, character.Name) {
		response = strings.TrimPrefix(response, character.Name)
	}
	response = strings.TrimSpace(response) // 移除可能存在的首尾空格

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
