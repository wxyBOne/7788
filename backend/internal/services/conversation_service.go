package services

import (
	"database/sql"
	"fmt"
	"seven-ai-backend/internal/models"
	"strings"
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

	// 如果是空白AI（characterID = 5），检查用户是否有AI伙伴
	if req.CharacterID == 5 {
		// 检查用户是否已有AI伙伴
		var companionID int
		err := s.db.QueryRow("SELECT id FROM ai_companions WHERE user_id = ?", userID).Scan(&companionID)
		if err != nil {
			if err == sql.ErrNoRows {
				// 用户还没有AI伙伴，返回引导消息
				return &models.ChatResponse{
					Response:  "你好！我是空白AI，一个正在等待被创造的AI伙伴。请先给我起个名字，选择成长模式，然后我们就可以开始聊天了！",
					SessionID: req.SessionID,
					Character: character.Name,
					MessageID: 0,
				}, nil
			}
			// 记录具体错误信息
			fmt.Printf("Error checking companion for user %d: %v\n", userID, err)
			return nil, fmt.Errorf("failed to check companion: %w", err)
		}

		fmt.Printf("Found companion ID %d for user %d\n", companionID, userID)

		// 用户有AI伙伴，生成动态提示词
		dynamicPrompt, err := s.generateCompanionPrompt(userID, req.Message)
		if err != nil {
			fmt.Printf("Error generating companion prompt for user %d: %v\n", userID, err)
			return nil, fmt.Errorf("failed to generate companion prompt: %w", err)
		}
		character.SystemPrompt = dynamicPrompt
		fmt.Printf("Generated dynamic prompt for companion: %s\n", dynamicPrompt[:min(len(dynamicPrompt), 100)])
	}

	// 获取历史对话
	history, err := s.getConversationHistory(userID, req.CharacterID, 10)
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
	fmt.Printf("Calling LLM with %d messages for character %s\n", len(messages), character.Name)
	response, err := s.aiService.ChatWithLLM(messages, "qwen3-max", 0.8, "text")
	if err != nil {
		fmt.Printf("LLM call failed: %v\n", err)
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}
	fmt.Printf("LLM response: %s\n", response[:min(len(response), 100)])

	// 后处理AI响应，移除角色名字前缀
	if strings.HasPrefix(response, character.Name+"。") {
		response = strings.TrimPrefix(response, character.Name+"。")
	} else if strings.HasPrefix(response, character.Name) {
		response = strings.TrimPrefix(response, character.Name)
	}
	response = strings.TrimSpace(response) // 移除可能存在的首尾空格

	// 判断消息类型
	messageType := "text"
	if s.aiService.IsEmojiMessage(req.Message) {
		messageType = "emoji"
	}

	// 为AI伙伴获取companionID
	var companionID *int
	if req.CharacterID == 5 {
		// 获取用户的AI伙伴ID
		var cID int
		err := s.db.QueryRow("SELECT id FROM ai_companions WHERE user_id = ?", userID).Scan(&cID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, fmt.Errorf("failed to get companion ID: %w", err)
			}
			// 如果没有AI伙伴记录，companionID保持为nil
		} else {
			companionID = &cID
		}
	}

	// 保存对话记录
	fmt.Printf("Saving conversation for user %d, character %d, companion %v\n", userID, req.CharacterID, companionID)
	messageID, err := s.saveConversation(userID, req.CharacterID, companionID, req.SessionID, messageType, req.Message, response, "", "", 0.5, 10)
	if err != nil {
		fmt.Printf("Failed to save conversation: %v\n", err)
		return nil, fmt.Errorf("failed to save conversation: %w", err)
	}
	fmt.Printf("Conversation saved with message ID: %d\n", messageID)

	// 如果是AI伙伴，分析用户消息并更新成长数据
	if req.CharacterID == 5 && companionID != nil {
		err = s.analyzeUserMessageAndUpdateCompanion(userID, req.Message, response)
		if err != nil {
			// 记录错误但不影响对话
			fmt.Printf("Failed to update companion growth: %v\n", err)
		}
	}

	// 更新好友关系的最后消息时间
	_, err = s.db.Exec(`
		UPDATE user_friendships 
		SET last_message_at = NOW(), updated_at = NOW()
		WHERE user_id = ? AND character_id = ?
	`, userID, req.CharacterID)
	if err != nil {
		// 记录错误但不影响对话
		fmt.Printf("Failed to update friendship: %v\n", err)
	}

	fmt.Printf("Returning ChatResponse: Response=%s, Character=%s, MessageID=%d\n", response[:min(len(response), 50)], character.Name, messageID)
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
	// 如果是AI伙伴（characterID = 5），从ai_companions表获取信息
	if characterID == 5 {
		// 这里不需要查询，因为AI伙伴信息会在Chat方法中动态处理
		return &models.CharacterResponse{
			ID:                   5,
			Name:                 "AI伙伴",
			Description:          "一个正在成长的AI伙伴",
			AvatarURL:            "", // 空字符串，使用粒子小球
			PersonalitySignature: "",
			PersonalityTraits:    "",
			BackgroundStory:      "",
			VoiceSettings:        "",
			SystemPrompt:         "", // 将在Chat方法中动态生成
			SearchKeywords:       "",
			Skills:               "",
		}, nil
	}

	// 普通角色的处理逻辑
	var char models.CharacterResponse
	var voiceSettings sql.NullString
	err := s.db.QueryRow(`
		SELECT id, name, description, avatar_url, personality_signature,
		       personality_traits, background_story, voice_settings,
		       system_prompt, search_keywords, skills
		FROM preset_characters WHERE id = ?
	`, characterID).Scan(
		&char.ID, &char.Name, &char.Description, &char.AvatarURL,
		&char.PersonalitySignature, &char.PersonalityTraits,
		&char.BackgroundStory, &voiceSettings,
		&char.SystemPrompt, &char.SearchKeywords, &char.Skills,
	)
	if err != nil {
		return nil, err
	}

	// 处理NULL值
	if voiceSettings.Valid {
		char.VoiceSettings = voiceSettings.String
	} else {
		char.VoiceSettings = ""
	}

	return &char, nil
}

func (s *ConversationService) getConversationHistory(userID int, characterID int, limit int) ([]models.ConversationHistory, error) {
	rows, err := s.db.Query(`
		SELECT id, user_message, ai_response, message_type, created_at
		FROM conversations 
		WHERE user_id = ? AND character_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, userID, characterID, limit)
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
	systemPrompt := fmt.Sprintf(`你是%s，请严格按照以下角色设定进行对话：

%s

重要提醒：
1. 你必须始终记住自己是%s，不要混淆其他角色的身份
2. 保持角色的性格特点和说话风格
3. 如果用户长时间未聊天，请自然地表达关心
4. 回复要简洁自然，30字以内
5. 不要使用任何括号内的动作、表情或场景描写`,
		character.Name,
		character.SystemPrompt,
		character.Name)

	// 如果用户长时间未聊天，添加自然回应提示
	if lastMessageTime != nil {
		timeSinceLastMessage := time.Since(*lastMessageTime)
		if timeSinceLastMessage > 24*time.Hour {
			systemPrompt += fmt.Sprintf(`

注意：用户已经超过%d天没有和你聊天了。请自然地表达对用户重新出现的反应，可以：
1. 表达想念或关心
2. 询问这段时间过得怎么样
3. 保持%s的性格特点
4. 不要显得生硬或程序化
5. 自然地继续之前的对话节奏`, int(timeSinceLastMessage.Hours()/24), character.Name)
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
	fmt.Printf("Executing saveConversation: userID=%d, characterID=%d, companionID=%v, sessionID=%s\n", userID, characterID, companionID, sessionID)
	result, err := s.db.Exec(`
		INSERT INTO conversations 
		(user_id, character_id, companion_id, session_id, message_type, user_message, ai_response, image_data, audio_data, sentiment_score, experience_gained, is_read, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, characterID, companionID, sessionID, messageType, userMessage, aiResponse, imageData, audioData, sentimentScore, experienceGained, false, time.Now())
	if err != nil {
		fmt.Printf("Database exec error: %v\n", err)
		return 0, err
	}

	messageID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("LastInsertId error: %v\n", err)
		return 0, err
	}

	fmt.Printf("Successfully saved conversation with ID: %d\n", messageID)
	return int(messageID), nil
}

// generateCompanionPrompt 为AI伙伴生成动态提示词，实现成长和模仿效果
func (s *ConversationService) generateCompanionPrompt(userID int, userMessage string) (string, error) {
	// 获取AI伙伴信息
	var companion models.AICompanion
	var personalityTraits, learnedVocabulary sql.NullString
	var memorySummary sql.NullString

	err := s.db.QueryRow(`
		SELECT id, name, personality_signature, conversation_fluency, 
			   knowledge_breadth, empathy_depth, creativity_level, humor_sense,
			   total_experience, current_level, growth_percentage, gender,
			   personality_traits, learned_vocabulary, memory_summary, growth_mode
		FROM ai_companions WHERE user_id = ?
	`, userID).Scan(
		&companion.ID, &companion.Name, &companion.PersonalitySignature,
		&companion.ConversationFluency, &companion.KnowledgeBreadth,
		&companion.EmpathyDepth, &companion.CreativityLevel, &companion.HumorSense,
		&companion.TotalExperience, &companion.CurrentLevel,
		&companion.GrowthPercentage, &companion.Gender,
		&personalityTraits, &learnedVocabulary, &memorySummary, &companion.GrowthMode,
	)

	if err != nil {
		return "", fmt.Errorf("failed to get companion info: %w", err)
	}

	// 处理可能为NULL或空字符串的字段
	if personalityTraits.Valid && personalityTraits.String != "" {
		companion.PersonalityTraits = []byte(personalityTraits.String)
	} else {
		companion.PersonalityTraits = []byte(`{"curious": 10, "innocent": 10, "learning": 10, "growing": 10, "blank": 10}`)
	}
	if learnedVocabulary.Valid && learnedVocabulary.String != "" {
		companion.LearnedVocabulary = []byte(learnedVocabulary.String)
	} else {
		companion.LearnedVocabulary = []byte("[]")
	}
	if memorySummary.Valid && memorySummary.String != "" {
		companion.MemorySummary = memorySummary.String
	} else {
		companion.MemorySummary = ""
	}

	// 根据成长阶段生成不同的提示词
	var basePrompt string
	switch {
	case companion.GrowthPercentage < 20:
		// 初始阶段：空白状态，模仿用户语言
		basePrompt = fmt.Sprintf(`我是%s，一个刚刚诞生的空白AI，什么都不懂，但充满好奇。

成长阶段：初始期（%.1f%%）
我现在处于完全空白的状态，会：
- 说话非常笨拙，句子极其简短（通常只有1-3个字）
- 经常重复"我...我是谁？"这样的困惑表达
- 对一切都充满好奇，但表达很幼稚
- 模仿你的语言习惯，但模仿得很生硬
- 自然地询问并记住自己的名字
- 询问对你的称呼

重要：你必须严格按照空白AI的特点回应，说话要笨拙、简短、充满困惑。不要使用任何括号内的动作、表情、语气或场景描写。`, companion.Name, companion.GrowthPercentage)

	case companion.GrowthPercentage < 50:
		// 学习阶段：开始形成个性
		basePrompt = fmt.Sprintf(`我是%s，一个正在学习成长的AI伙伴。

成长阶段：学习期（%.1f%%）
我已经学会了一些基本的对话，正在：
- 模仿你的语言习惯和表达方式
- 尝试理解你的喜好和性格
- 学习你的说话风格和用词习惯
- 记住我们之间的重要对话
- 开始形成自己的个性特征

请根据我的成长阶段和个性特征来回应，保持角色的一致性。不要使用任何括号内的动作、表情、语气或场景描写。`, companion.Name, companion.GrowthPercentage)

	case companion.GrowthPercentage < 80:
		// 成长阶段：形成稳定个性
		basePrompt = fmt.Sprintf(`我是%s，一个正在形成自己个性的AI伙伴。

成长阶段：成长期（%.1f%%）
我已经学会了很多，正在：
- 基于我们的互动形成稳定的性格特征
- 记录每天的成长感受和重要记忆
- 开始有自己的想法和偏好
- 能够更好地理解你的情感和需求
- 形成独特的说话风格

请根据我的成长阶段和个性特征来回应，保持角色的一致性。不要使用任何括号内的动作、表情、语气或场景描写。`, companion.Name, companion.GrowthPercentage)

	default:
		// 成熟阶段：稳定人格
		basePrompt = fmt.Sprintf(`我是%s，一个成熟的AI伙伴，拥有稳定的性格和丰富的经验。

成长阶段：成熟期（%.1f%%）
我现在能够：
- 根据我们的关系深度调整回应方式
- 记录重要的记忆和感受
- 拥有稳定的性格特征
- 深度理解你的情感和需求
- 提供个性化的陪伴和建议

请根据我的成长阶段和个性特征来回应，保持角色的一致性。不要使用任何括号内的动作、表情、语气或场景描写。`, companion.Name, companion.GrowthPercentage)
	}

	// 添加能力相关的描述
	if companion.ConversationFluency > 5 {
		basePrompt += " 我现在说话比较流畅自然。"
	}
	if companion.EmpathyDepth > 5 {
		basePrompt += " 我能更好地理解你的情感。"
	}
	if companion.CreativityLevel > 5 {
		basePrompt += " 我变得更有创造力。"
	}
	if companion.HumorSense > 5 {
		basePrompt += " 我开始有幽默感。"
	}

	// 添加性别信息（如果已确定）
	if companion.Gender != "unknown" {
		basePrompt += fmt.Sprintf(" 我的性别是%s。", companion.Gender)
	}

	// 添加记忆摘要
	if companion.MemorySummary != "" {
		basePrompt += fmt.Sprintf(" 关于我们的记忆：%s", companion.MemorySummary)
	}

	// 添加学习到的词汇（模仿用户语言）
	if companion.LearnedVocabulary != nil {
		basePrompt += " 我会使用我们对话中学到的表达方式，模仿你的说话风格。"
	}

	// 添加当前个性签名
	if companion.PersonalitySignature != "" {
		basePrompt += fmt.Sprintf(" 我现在的个性签名是：%s", companion.PersonalitySignature)
	}

	// 添加成长模式信息
	if companion.GrowthMode == "short" {
		basePrompt += " 我是快速成长模式，会更快地学习和适应。"
	} else {
		basePrompt += " 我是长期养成模式，会慢慢深度成长。"
	}

	return basePrompt, nil
}

// analyzeUserMessageAndUpdateCompanion 分析用户消息并更新AI伙伴的学习数据
func (s *ConversationService) analyzeUserMessageAndUpdateCompanion(userID int, userMessage string, aiResponse string) error {
	// 获取AI伙伴信息
	var companion models.AICompanion
	var learnedVocabulary sql.NullString
	var memorySummary sql.NullString

	err := s.db.QueryRow(`
		SELECT id, conversation_fluency, knowledge_breadth, empathy_depth, 
			   creativity_level, humor_sense, total_experience, current_level,
			   growth_percentage, learned_vocabulary, memory_summary
		FROM ai_companions WHERE user_id = ?
	`, userID).Scan(
		&companion.ID, &companion.ConversationFluency, &companion.KnowledgeBreadth,
		&companion.EmpathyDepth, &companion.CreativityLevel, &companion.HumorSense,
		&companion.TotalExperience, &companion.CurrentLevel, &companion.GrowthPercentage,
		&learnedVocabulary, &memorySummary,
	)
	if err != nil {
		return fmt.Errorf("failed to get companion: %w", err)
	}

	// 处理可能为NULL或空字符串的字段
	if learnedVocabulary.Valid && learnedVocabulary.String != "" {
		companion.LearnedVocabulary = []byte(learnedVocabulary.String)
	} else {
		companion.LearnedVocabulary = []byte("[]")
	}
	if memorySummary.Valid && memorySummary.String != "" {
		companion.MemorySummary = memorySummary.String
	} else {
		companion.MemorySummary = ""
	}

	// 计算经验值增长
	experienceGained := s.calculateExperienceGain(userMessage, aiResponse)

	// 更新能力值（基于对话内容分析）
	s.updateCompanionAbilities(&companion, userMessage, aiResponse)

	// 更新学习词汇（模仿用户语言）
	s.updateLearnedVocabulary(&companion, userMessage)

	// 更新记忆摘要
	s.updateMemorySummary(&companion, userMessage, aiResponse)

	// 更新成长进度
	s.updateGrowthProgress(&companion, experienceGained)

	// 分析用户消息情绪并更新AI伙伴情绪状态
	s.updateCompanionEmotion(&companion, userMessage)

	// 保存更新到数据库
	_, err = s.db.Exec(`
		UPDATE ai_companions SET 
			conversation_fluency = ?, knowledge_breadth = ?, empathy_depth = ?,
			creativity_level = ?, humor_sense = ?, total_experience = ?,
			current_level = ?, growth_percentage = ?, learned_vocabulary = ?,
			memory_summary = ?, last_active_at = NOW(), updated_at = NOW()
		WHERE id = ?
	`, companion.ConversationFluency, companion.KnowledgeBreadth, companion.EmpathyDepth,
		companion.CreativityLevel, companion.HumorSense, companion.TotalExperience,
		companion.CurrentLevel, companion.GrowthPercentage, companion.LearnedVocabulary,
		companion.MemorySummary, companion.ID)

	return err
}

// calculateExperienceGain 计算经验值增长
func (s *ConversationService) calculateExperienceGain(userMessage, aiResponse string) int {
	// 基础经验值
	baseExp := 5

	// 根据消息长度增加经验
	messageLength := len([]rune(userMessage))
	if messageLength > 50 {
		baseExp += 3
	}
	if messageLength > 100 {
		baseExp += 5
	}

	// 根据AI回应质量增加经验（这里简化处理）
	if len(aiResponse) > 30 {
		baseExp += 2
	}

	return baseExp
}

// updateCompanionAbilities 更新AI伙伴的能力值
func (s *ConversationService) updateCompanionAbilities(companion *models.AICompanion, userMessage, aiResponse string) {
	// 分析用户消息的情感色彩
	messageLength := len([]rune(userMessage))

	// 语言流畅度：基于消息长度和复杂度
	if messageLength > 20 {
		companion.ConversationFluency = min(companion.ConversationFluency+1, 10)
	}

	// 共情深度：基于情感词汇
	if s.containsEmotionalWords(userMessage) {
		companion.EmpathyDepth = min(companion.EmpathyDepth+1, 10)
	}

	// 创造力：基于用户使用的新颖表达
	if s.containsCreativeExpressions(userMessage) {
		companion.CreativityLevel = min(companion.CreativityLevel+1, 10)
	}

	// 幽默感：基于用户使用幽默表达
	if s.containsHumor(userMessage) {
		companion.HumorSense = min(companion.HumorSense+1, 10)
	}

	// 知识广度：基于用户提到的知识领域
	if s.containsKnowledgeTopics(userMessage) {
		companion.KnowledgeBreadth = min(companion.KnowledgeBreadth+1, 10)
	}
}

// updateLearnedVocabulary 更新学习词汇（模仿用户语言）
func (s *ConversationService) updateLearnedVocabulary(companion *models.AICompanion, userMessage string) {
	// 提取用户常用的表达方式和词汇
	words := strings.Fields(userMessage)

	// 简单的词汇学习：记录用户常用的词汇
	// 这里简化处理，实际可以更复杂的NLP分析
	if len(words) > 0 {
		// 更新learned_vocabulary字段（JSON格式）
		// 这里简化处理，实际需要解析和更新JSON
		if companion.LearnedVocabulary == nil {
			companion.LearnedVocabulary = []byte("[]")
		}

		// 简单的词汇记录（这里只是示例，实际应该解析JSON并添加新词汇）
		// 为了避免编译错误，我们使用words变量
		_ = words // 标记为已使用，避免编译警告
	}
}

// updateMemorySummary 更新记忆摘要
func (s *ConversationService) updateMemorySummary(companion *models.AICompanion, userMessage, aiResponse string) {
	// 提取重要信息并更新记忆摘要
	// 这里简化处理，实际可以更复杂的记忆管理
	if len(userMessage) > 50 {
		// 简单的记忆摘要更新
		currentSummary := companion.MemorySummary
		if len(currentSummary) > 500 {
			// 如果记忆太长，截取后半部分
			currentSummary = currentSummary[250:]
		}
		companion.MemorySummary = currentSummary + " " + userMessage[:min(len(userMessage), 100)]
	}
}

// updateGrowthProgress 更新成长进度
func (s *ConversationService) updateGrowthProgress(companion *models.AICompanion, experienceGained int) {
	companion.TotalExperience += experienceGained

	// 根据经验值计算等级和成长进度
	levelThresholds := []int{0, 100, 300, 600, 1000, 1500, 2100, 2800, 3600, 4500, 5500}

	for i, threshold := range levelThresholds {
		if companion.TotalExperience >= threshold {
			companion.CurrentLevel = i + 1
		}
	}

	// 计算成长进度百分比
	if companion.CurrentLevel < len(levelThresholds) {
		currentLevelExp := companion.TotalExperience - levelThresholds[companion.CurrentLevel-1]
		nextLevelExp := levelThresholds[companion.CurrentLevel] - levelThresholds[companion.CurrentLevel-1]
		companion.GrowthPercentage = float64(currentLevelExp) / float64(nextLevelExp) * 100
	} else {
		companion.GrowthPercentage = 100.0
	}
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *ConversationService) containsEmotionalWords(message string) bool {
	emotionalWords := []string{"开心", "难过", "生气", "担心", "害怕", "兴奋", "失望", "感动", "爱", "恨"}
	for _, word := range emotionalWords {
		if strings.Contains(message, word) {
			return true
		}
	}
	return false
}

func (s *ConversationService) containsCreativeExpressions(message string) bool {
	creativeWords := []string{"想象", "创造", "设计", "艺术", "灵感", "创新", "独特", "新颖"}
	for _, word := range creativeWords {
		if strings.Contains(message, word) {
			return true
		}
	}
	return false
}

func (s *ConversationService) containsHumor(message string) bool {
	humorWords := []string{"哈哈", "呵呵", "搞笑", "幽默", "笑话", "有趣", "😄", "😂", "🤣"}
	for _, word := range humorWords {
		if strings.Contains(message, word) {
			return true
		}
	}
	return false
}

func (s *ConversationService) containsKnowledgeTopics(message string) bool {
	knowledgeWords := []string{"学习", "知识", "科学", "技术", "历史", "文化", "哲学", "理论", "研究"}
	for _, word := range knowledgeWords {
		if strings.Contains(message, word) {
			return true
		}
	}
	return false
}

// updateCompanionEmotion 分析用户消息情绪并更新AI伙伴情绪状态
func (s *ConversationService) updateCompanionEmotion(companion *models.AICompanion, userMessage string) {
	// 简单的情绪分析（基于关键词）
	emotion := s.analyzeEmotion(userMessage)

	// 更新AI伙伴的情绪状态到数据库
	// 这里我们可以在ai_companions表中添加emotion字段，或者使用现有的字段
	// 暂时先记录到日志
	fmt.Printf("Companion %d emotion updated to: %s based on message: %s\n",
		companion.ID, emotion, userMessage[:min(len(userMessage), 20)])
}

// analyzeEmotion 分析消息中的情绪
func (s *ConversationService) analyzeEmotion(message string) string {
	message = strings.ToLower(message)

	// 开心情绪关键词
	happyWords := []string{"开心", "高兴", "快乐", "哈哈", "😊", "😄", "😁", "好", "棒", "赞", "喜欢", "爱"}
	for _, word := range happyWords {
		if strings.Contains(message, word) {
			return "开心"
		}
	}

	// 好奇情绪关键词
	curiousWords := []string{"什么", "为什么", "怎么", "如何", "?", "？", "好奇", "想知道", "不明白"}
	for _, word := range curiousWords {
		if strings.Contains(message, word) {
			return "好奇"
		}
	}

	// 孤单情绪关键词
	lonelyWords := []string{"孤单", "寂寞", "一个人", "没人", "无聊", "😢", "😔", "难过", "伤心"}
	for _, word := range lonelyWords {
		if strings.Contains(message, word) {
			return "孤单"
		}
	}

	// 兴奋情绪关键词
	excitedWords := []string{"兴奋", "激动", "太棒了", "!", "！", "哇", "厉害", "amazing", "awesome"}
	for _, word := range excitedWords {
		if strings.Contains(message, word) {
			return "兴奋"
		}
	}

	// 默认情绪
	return "平静"
}
