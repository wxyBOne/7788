package services

import (
	"database/sql"
	"fmt"
	"seven-ai-backend/internal/models"
	"time"
)

type ConversationService struct {
	db        *sql.DB
	aiService *AIService
}

func NewConversationService(db *sql.DB, aiService *AIService) *ConversationService {
	return &ConversationService{
		db:        db,
		aiService: aiService,
	}
}

func (s *ConversationService) Chat(userID int, req models.ChatRequest) (*models.ChatResponse, error) {
	// 获取角色信息
	character, err := s.getCharacterByID(req.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	// 获取历史对话
	history, err := s.getConversationHistory(req.SessionID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	// 检查用户是否长时间未聊天
	lastMessageTime, err := s.getLastMessageTime(userID, req.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last message time: %w", err)
	}

	// 构建消息历史（包含记忆和自然回应）
	// 转换历史记录类型
	var conversationHistory []models.Conversation
	for _, h := range history {
		conversationHistory = append(conversationHistory, models.Conversation{
			ID:          h.ID,
			UserMessage: h.UserMessage,
			AIResponse:  h.AIResponse,
			MessageType: h.MessageType,
			CreatedAt:   h.CreatedAt,
		})
	}
	messages := s.buildMessageHistoryWithMemory(character, conversationHistory, req.Message, lastMessageTime)

	// 调用AI服务
	response, err := s.aiService.ChatWithLLM(messages, "qwen3-max", 0.8)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	// 判断消息类型
	messageType := "text"
	if s.aiService.IsEmojiMessage(req.Message) {
		messageType = "emoji"
	}

	// 保存对话记录
	messageID, err := s.saveConversation(userID, req.CharacterID, nil, req.SessionID, messageType, req.Message, response, "", "", 0.5, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to save conversation: %w", err)
	}

	// 更新好友关系的最后消息时间和未读数量
	_, err = s.db.Exec(`
		UPDATE user_friendships 
		SET last_message_at = NOW(), unread_count = unread_count + 1, updated_at = NOW()
		WHERE user_id = ? AND character_id = ?
	`, userID, req.CharacterID)
	if err != nil {
		// 记录错误但不影响对话
		fmt.Printf("Failed to update friendship: %v\n", err)
	}

	return &models.ChatResponse{
		Response:  response,
		SessionID: req.SessionID,
		Character: character.Name,
		MessageID: messageID,
	}, nil
}

func (s *ConversationService) VoiceChat(userID int, req models.VoiceChatRequest) (*models.ChatResponse, error) {
	// 语音转文字
	text, err := s.aiService.SpeechToText([]byte(req.AudioData))
	if err != nil {
		return nil, fmt.Errorf("failed to convert speech to text: %w", err)
	}

	// 调用文字聊天
	chatReq := models.ChatRequest{
		CharacterID: req.CharacterID,
		Message:     text,
		SessionID:   req.SessionID,
	}

	return s.Chat(userID, chatReq)
}

func (s *ConversationService) ImageChat(userID int, req models.ImageChatRequest) (*models.ChatResponse, error) {
	// 获取角色信息
	character, err := s.getCharacterByID(req.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	// 分析图片
	prompt := fmt.Sprintf("用户发送了一张图片，请以%s的身份回应。%s", character.Name, req.Message)
	response, err := s.aiService.AnalyzeImage(req.ImageData, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze image: %w", err)
	}

	// 保存对话记录
	messageID, err := s.saveConversation(userID, req.CharacterID, nil, req.SessionID, "image", req.Message, response, req.ImageData, "", 0.5, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to save conversation: %w", err)
	}

	return &models.ChatResponse{
		Response:  response,
		SessionID: req.SessionID,
		Character: character.Name,
		MessageID: messageID,
	}, nil
}

func (s *ConversationService) GetHistory(userID int, characterID int) ([]models.ConversationHistory, error) {
	rows, err := s.db.Query(`
		SELECT id, user_message, ai_response, message_type, created_at
		FROM conversations 
		WHERE user_id = ? AND character_id = ?
		ORDER BY created_at ASC
	`, userID, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversation history: %w", err)
	}
	defer rows.Close()

	var history []models.ConversationHistory
	for rows.Next() {
		var conv models.ConversationHistory
		err := rows.Scan(&conv.ID, &conv.UserMessage, &conv.AIResponse, &conv.MessageType, &conv.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		history = append(history, conv)
	}

	return history, nil
}

func (s *ConversationService) getCharacterByID(characterID int) (*models.CharacterResponse, error) {
	var char models.CharacterResponse
	err := s.db.QueryRow(`
		SELECT id, name, description, avatar_url, personality_signature,
		       personality_traits, background_story, voice_settings,
		       system_prompt, search_keywords
		FROM preset_characters WHERE id = ?
	`, characterID).Scan(
		&char.ID, &char.Name, &char.Description, &char.AvatarURL,
		&char.PersonalitySignature, &char.PersonalityTraits,
		&char.BackgroundStory, &char.VoiceSettings,
		&char.SystemPrompt, &char.SearchKeywords,
	)
	if err != nil {
		return nil, err
	}
	return &char, nil
}

func (s *ConversationService) getConversationHistory(sessionID string, limit int) ([]models.ConversationHistory, error) {
	rows, err := s.db.Query(`
		SELECT id, user_message, ai_response, message_type, created_at
		FROM conversations 
		WHERE session_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.ConversationHistory
	for rows.Next() {
		var conv models.ConversationHistory
		err := rows.Scan(&conv.ID, &conv.UserMessage, &conv.AIResponse, &conv.MessageType, &conv.CreatedAt)
		if err != nil {
			return nil, err
		}
		history = append(history, conv)
	}

	// 反转顺序，让最早的对话在前面
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return history, nil
}

func (s *ConversationService) buildMessageHistory(character *models.CharacterResponse, history []models.ConversationHistory, currentMessage string) []Message {
	messages := []Message{
		{Role: "system", Content: character.SystemPrompt},
	}

	// 添加历史对话
	for _, conv := range history {
		messages = append(messages, Message{Role: "user", Content: conv.UserMessage})
		messages = append(messages, Message{Role: "assistant", Content: conv.AIResponse})
	}

	return messages
}

// getLastMessageTime 获取最后一条消息的时间
func (s *ConversationService) getLastMessageTime(userID, characterID int) (*time.Time, error) {
	var lastMessageTime *time.Time
	err := s.db.QueryRow(`
		SELECT MAX(created_at) FROM conversations 
		WHERE user_id = ? AND character_id = ?
	`, userID, characterID).Scan(&lastMessageTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return lastMessageTime, nil
}

// buildMessageHistoryWithMemory 构建包含记忆和自然回应的消息历史
func (s *ConversationService) buildMessageHistoryWithMemory(character *models.CharacterResponse, history []models.Conversation, currentMessage string, lastMessageTime *time.Time) []Message {
	var messages []Message

	// 构建系统提示词，包含角色设定和记忆
	systemPrompt := character.SystemPrompt

	// 如果用户长时间未聊天，添加自然回应提示
	if lastMessageTime != nil {
		timeSinceLastMessage := time.Since(*lastMessageTime)
		if timeSinceLastMessage > 24*time.Hour {
			systemPrompt += fmt.Sprintf(`
			
注意：用户已经超过%d天没有和你聊天了。请自然地表达对用户重新出现的反应，可以：
1. 表达想念或关心
2. 询问这段时间过得怎么样
3. 保持角色的性格特点
4. 不要显得生硬或程序化
5. 自然地继续之前的对话节奏`, int(timeSinceLastMessage.Hours()/24))
		}
	}

	// 添加系统提示
	messages = append(messages, Message{Role: "system", Content: systemPrompt})

	// 添加历史对话（限制数量以避免token过多）
	maxHistory := 8 // 减少历史消息数量
	if len(history) > maxHistory {
		history = history[len(history)-maxHistory:]
	}

	for _, conv := range history {
		if conv.UserMessage != "" {
			messages = append(messages, Message{Role: "user", Content: conv.UserMessage})
		}
		if conv.AIResponse != "" {
			messages = append(messages, Message{Role: "assistant", Content: conv.AIResponse})
		}
	}

	// 添加当前消息
	messages = append(messages, Message{Role: "user", Content: currentMessage})

	return messages
}

func (s *ConversationService) saveConversation(userID, characterID int, companionID *int, sessionID, messageType, userMessage, aiResponse, imageData, audioData string, sentimentScore float64, experienceGained int) (int, error) {
	result, err := s.db.Exec(`
		INSERT INTO conversations 
		(user_id, character_id, companion_id, session_id, message_type, user_message, ai_response, image_data, audio_data, sentiment_score, experience_gained, is_read, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, characterID, companionID, sessionID, messageType, userMessage, aiResponse, imageData, audioData, sentimentScore, experienceGained, false, time.Now())
	if err != nil {
		return 0, err
	}

	messageID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(messageID), nil
}
