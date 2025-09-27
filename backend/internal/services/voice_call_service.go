package services

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"seven-ai-backend/internal/models"
	"time"
)

type VoiceCallService struct {
	db        *sql.DB
	aiService *AIService
}

func NewVoiceCallService(db *sql.DB, aiService *AIService) *VoiceCallService {
	return &VoiceCallService{
		db:        db,
		aiService: aiService,
	}
}

// ProcessVoiceCall 处理语音通话
func (s *VoiceCallService) ProcessVoiceCall(userID int, req models.VoiceCallRequest) (*models.VoiceCallResponse, error) {
	var userText string
	var aiText string
	var messageID int64

	// 如果不是第一次通话，处理用户语音
	if !req.IsFirstCall && req.AudioData != "" {
		// 解码音频数据
		audioData, err := base64.StdEncoding.DecodeString(req.AudioData)
		if err != nil {
			return nil, fmt.Errorf("音频数据解码失败: %v", err)
		}

		// 语音转文字
		userText, err = s.aiService.SpeechToText(audioData)
		if err != nil {
			return nil, fmt.Errorf("语音识别失败: %v", err)
		}

		// 检查是否为噪音或空内容
		if userText == "" || len(userText) < 2 {
			// 获取角色信息用于噪音响应
			character, err := s.getCharacterByID(req.CharacterID)
			if err != nil {
				return nil, fmt.Errorf("获取角色信息失败: %v", err)
			}

			// 返回噪音响应
			noiseResponse := s.getNoiseResponseForCharacter(character.Name)

			// 生成噪音响应的语音
			aiAudioData, err := s.aiService.TextToSpeech(noiseResponse, character.Name)
			if err != nil {
				return nil, fmt.Errorf("生成噪音响应语音失败: %v", err)
			}

			aiAudioBase64 := base64.StdEncoding.EncodeToString(aiAudioData)

			return &models.VoiceCallResponse{
				TextResponse:  noiseResponse,
				AudioResponse: aiAudioBase64, // 返回音频给前端播放
				SessionID:     req.SessionID,
				Character:     character.Name,
				MessageID:     0,
			}, nil
		}

		// 获取对话历史
		history, err := s.getConversationHistory(userID, req.CharacterID, 5)
		if err != nil {
			return nil, fmt.Errorf("获取对话历史失败: %v", err)
		}

		// 构建消息历史
		messageHistory := s.buildMessageHistoryWithMemory(history, 4)

		// 构建消息列表
		messages := append(messageHistory, Message{
			Role:    "user",
			Content: userText,
		})

		// AI回复
		aiText, err = s.aiService.ChatWithLLM(messages, "", 0.7, "voice")
		if err != nil {
			return nil, fmt.Errorf("AI回复生成失败: %v", err)
		}

		// 保存对话记录
		messageID, err = s.saveVoiceCall(userID, req.CharacterID, userText, aiText, req.SessionID)
		if err != nil {
			return nil, fmt.Errorf("保存对话记录失败: %v", err)
		}
	} else {
		// 第一次通话，AI主动打招呼
		character, err := s.getCharacterByID(req.CharacterID)
		if err != nil {
			return nil, fmt.Errorf("获取角色信息失败: %v", err)
		}

		// 根据角色个性生成个性化打招呼
		aiText = s.generateGreetingForCharacter(character)

		// 保存AI打招呼记录
		messageID, err = s.saveVoiceCall(userID, req.CharacterID, "", aiText, req.SessionID)
		if err != nil {
			return nil, fmt.Errorf("保存对话记录失败: %v", err)
		}
	}

	// 生成AI语音回复（返回给前端播放）
	var aiAudioBase64 string
	if aiText != "" {
		character, err := s.getCharacterByID(req.CharacterID)
		if err != nil {
			return nil, fmt.Errorf("获取角色信息失败: %v", err)
		}

		aiAudioData, err := s.aiService.TextToSpeech(aiText, character.Name)
		if err != nil {
			return nil, fmt.Errorf("生成AI语音回复失败: %v", err)
		}

		aiAudioBase64 = base64.StdEncoding.EncodeToString(aiAudioData)
	}

	return &models.VoiceCallResponse{
		TextResponse:  aiText,
		AudioResponse: aiAudioBase64, // 返回音频给前端播放
		SessionID:     req.SessionID,
		Character:     "",
		MessageID:     int(messageID),
	}, nil
}

// GetVoiceCallHistory 获取语音通话历史
func (s *VoiceCallService) GetVoiceCallHistory(userID, characterID int) ([]models.VoiceCall, error) {
	query := `
		SELECT id, user_id, character_id, companion_id, message_type, created_at
		FROM conversations 
		WHERE user_id = ? AND character_id = ? AND message_type = 'voice'
		ORDER BY created_at DESC 
		LIMIT 50
	`

	rows, err := s.db.Query(query, userID, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calls []models.VoiceCall
	for rows.Next() {
		var call models.VoiceCall
		var messageType string
		var createdAt time.Time

		err := rows.Scan(&call.ID, &call.UserID, &call.CharacterID, &call.CompanionID, &messageType, &createdAt)
		if err != nil {
			continue
		}

		// 设置默认值
		call.CallType = "voice"
		call.Status = "completed"
		call.DurationSeconds = 0
		call.AudioFileURL = ""
		call.StartedAt = createdAt
		call.EndedAt = time.Time{}

		calls = append(calls, call)
	}

	return calls, nil
}

// getCharacterByID 根据ID获取角色信息
func (s *VoiceCallService) getCharacterByID(characterID int) (*models.PresetCharacter, error) {
	query := `SELECT id, name, personality_signature, avatar_url FROM preset_characters WHERE id = ?`

	var character models.PresetCharacter
	err := s.db.QueryRow(query, characterID).Scan(
		&character.ID,
		&character.Name,
		&character.PersonalitySignature,
		&character.AvatarURL,
	)

	if err != nil {
		return nil, err
	}

	return &character, nil
}

// getConversationHistory 获取对话历史
func (s *VoiceCallService) getConversationHistory(userID, characterID int, limit int) ([]models.Conversation, error) {
	query := `
		SELECT id, user_id, character_id, user_message, ai_response, message_type, created_at
		FROM conversations 
		WHERE user_id = ? AND character_id = ?
		ORDER BY created_at DESC 
		LIMIT ?
	`

	rows, err := s.db.Query(query, userID, characterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		var characterID int
		err := rows.Scan(
			&conv.ID,
			&conv.UserID,
			&characterID,
			&conv.UserMessage,
			&conv.AIResponse,
			&conv.MessageType,
			&conv.CreatedAt,
		)
		if err != nil {
			continue
		}
		conv.CharacterID = &characterID
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

// buildMessageHistoryWithMemory 构建带记忆的消息历史
func (s *VoiceCallService) buildMessageHistoryWithMemory(conversations []models.Conversation, maxHistory int) []Message {
	var messages []Message

	// 限制历史消息数量
	if len(conversations) > maxHistory {
		conversations = conversations[:maxHistory]
	}

	// 反转顺序，让最新的消息在后面
	for i := len(conversations) - 1; i >= 0; i-- {
		conv := conversations[i]

		if conv.UserMessage != "" {
			messages = append(messages, Message{
				Role:    "user",
				Content: conv.UserMessage,
			})
		}

		if conv.AIResponse != "" {
			messages = append(messages, Message{
				Role:    "assistant",
				Content: conv.AIResponse,
			})
		}
	}

	return messages
}

// saveVoiceCall 保存语音通话记录到conversations表
func (s *VoiceCallService) saveVoiceCall(userID, characterID int, userText, aiText, sessionID string) (int64, error) {
	query := `
		INSERT INTO conversations (user_id, character_id, user_message, ai_response, message_type, session_id, created_at)
		VALUES (?, ?, ?, ?, 'voice', ?, ?)
	`

	result, err := s.db.Exec(query, userID, characterID, userText, aiText, sessionID, time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// getNoiseResponseForCharacter 根据角色人设返回噪音响应
func (s *VoiceCallService) getNoiseResponseForCharacter(characterName string) string {
	switch characterName {
	case "林黛玉":
		return "哎呀，这声音怎么听不清呢？你那边是不是太吵了？"
	case "孙悟空":
		return "兄弟，你这声音俺老孙听不清啊，是不是环境太吵了？"
	case "李白":
		return "这声音如雾里看花，听不真切，莫非是环境嘈杂？"
	case "赫敏·格兰杰":
		return "抱歉，环境噪音太大，我听不清楚你说什么。"
	default:
		return "抱歉，我听不清楚你说什么，可能是环境太吵了。"
	}
}

// generateGreetingForCharacter 根据角色个性生成个性化打招呼
func (s *VoiceCallService) generateGreetingForCharacter(character *models.PresetCharacter) string {
	switch character.Name {
	case "林黛玉":
		return "哟，你来了？今日倒是比往日早了些，莫非是有什么心事想与我聊聊？"
	case "孙悟空":
		return "兄弟，你来了！俺老孙正闲着，有什么话尽管说，别客气！"
	case "李白":
		return "人生得意须尽欢，今日得与君通话，当浮一大白！有什么想聊的？"
	case "赫敏·格兰杰":
		return "你好，很高兴接到你的通话。有什么问题需要我帮助分析的吗？"
	default:
		return fmt.Sprintf("你好！我是%s，很高兴和你通话！", character.Name)
	}
}
