package services

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"seven-ai-backend/internal/models"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// QiniuASRClient 七牛云WebSocket ASR客户端
type QiniuASRClient struct {
	conn   *websocket.Conn
	seq    int
	apiKey string
}

// NewQiniuASRClient 创建七牛云ASR客户端
func NewQiniuASRClient(apiKey string) (*QiniuASRClient, error) {
	u := url.URL{Scheme: "wss", Host: "openai.qiniu.com", Path: "/v1/voice/asr"}
	header := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", apiKey)},
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return nil, fmt.Errorf("连接七牛云ASR失败: %v", err)
	}

	client := &QiniuASRClient{
		conn:   conn,
		seq:    1,
		apiKey: apiKey,
	}

	// 发送配置
	if err := client.sendConfig(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("发送ASR配置失败: %v", err)
	}

	return client, nil
}

// sendConfig 发送ASR配置
func (c *QiniuASRClient) sendConfig() error {
	req := map[string]interface{}{
		"user": map[string]interface{}{
			"uid": "streaming_client",
		},
		"audio": map[string]interface{}{
			"format":      "pcm",
			"sample_rate": 16000,
			"bits":        16,
			"channel":     1,
			"codec":       "raw",
		},
		"request": map[string]interface{}{
			"model_name":  "asr",
			"enable_punc": true,
		},
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	compressed := &bytes.Buffer{}
	gz := gzip.NewWriter(compressed)
	gz.Write(payload)
	gz.Close()

	// 生成协议头
	header := []byte{0x11, 0x11, 0x11, 0x00} // 版本1, 消息类型1, 序列1, 压缩1
	seqBytes := make([]byte, 4)
	seqBytes[0] = byte(c.seq >> 24)
	seqBytes[1] = byte(c.seq >> 16)
	seqBytes[2] = byte(c.seq >> 8)
	seqBytes[3] = byte(c.seq)

	sizeBytes := make([]byte, 4)
	sizeBytes[0] = byte(len(compressed.Bytes()) >> 24)
	sizeBytes[1] = byte(len(compressed.Bytes()) >> 16)
	sizeBytes[2] = byte(len(compressed.Bytes()) >> 8)
	sizeBytes[3] = byte(len(compressed.Bytes()))

	message := append(header, seqBytes...)
	message = append(message, sizeBytes...)
	message = append(message, compressed.Bytes()...)

	return c.conn.WriteMessage(websocket.BinaryMessage, message)
}

// SendAudio 发送音频数据
func (c *QiniuASRClient) SendAudio(audioData []byte) error {
	c.seq++

	compressed := &bytes.Buffer{}
	gz := gzip.NewWriter(compressed)
	gz.Write(audioData)
	gz.Close()

	// 生成协议头 (消息类型2 = 音频数据)
	header := []byte{0x11, 0x21, 0x11, 0x00} // 版本1, 消息类型2, 序列1, 压缩1
	seqBytes := make([]byte, 4)
	seqBytes[0] = byte(c.seq >> 24)
	seqBytes[1] = byte(c.seq >> 16)
	seqBytes[2] = byte(c.seq >> 8)
	seqBytes[3] = byte(c.seq)

	sizeBytes := make([]byte, 4)
	sizeBytes[0] = byte(len(compressed.Bytes()) >> 24)
	sizeBytes[1] = byte(len(compressed.Bytes()) >> 16)
	sizeBytes[2] = byte(len(compressed.Bytes()) >> 8)
	sizeBytes[3] = byte(len(compressed.Bytes()))

	message := append(header, seqBytes...)
	message = append(message, sizeBytes...)
	message = append(message, compressed.Bytes()...)

	return c.conn.WriteMessage(websocket.BinaryMessage, message)
}

// ReadResponse 读取ASR响应
func (c *QiniuASRClient) ReadResponse() (string, error) {
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		return "", err
	}

	// 解析响应
	if len(message) < 4 {
		return "", fmt.Errorf("响应消息太短")
	}

	headerSize := message[0] & 0x0f
	messageType := message[1] >> 4
	messageFlags := message[1] & 0x0f
	serialMethod := message[2] >> 4
	compression := message[2] & 0x0f

	payload := message[headerSize*4:]

	// 跳过序列号
	if messageFlags&0x01 != 0 {
		payload = payload[4:]
	}

	// 解析payload大小
	if messageType == 0x09 { // FULL_SERVER_RESPONSE
		if len(payload) < 4 {
			return "", fmt.Errorf("payload太短")
		}
		payloadSize := int(payload[0])<<24 | int(payload[1])<<16 | int(payload[2])<<8 | int(payload[3])
		payload = payload[4 : 4+payloadSize]
	}

	// 解压缩
	if compression == 0x01 {
		reader, err := gzip.NewReader(bytes.NewReader(payload))
		if err != nil {
			return "", err
		}
		decompressed := &bytes.Buffer{}
		decompressed.ReadFrom(reader)
		payload = decompressed.Bytes()
	}

	// 解析JSON
	if serialMethod == 0x01 {
		var result map[string]interface{}
		if err := json.Unmarshal(payload, &result); err != nil {
			return "", err
		}

		// 提取文本
		if result["result"] != nil {
			if resultMap, ok := result["result"].(map[string]interface{}); ok {
				if text, ok := resultMap["text"].(string); ok {
					return text, nil
				}
			}
		}
	}

	return "", nil
}

// Close 关闭连接
func (c *QiniuASRClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Character 角色信息结构体
type Character struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// StreamingVoiceCallService 流式语音通话服务
type StreamingVoiceCallService struct {
	aiService *AIService
	asrClient *RealtimeASRClient
	db        *sql.DB
	mu        sync.RWMutex
	sessions  map[string]*VoiceCallSession
	// 添加WebSocket连接通知回调
	onConnectionError func(sessionID string, err error)
	// 添加AI回复回调
	onResponseCallback func(sessionID, userText, aiText string, audioData []byte)
}

// VoiceCallSession 语音通话会话
type VoiceCallSession struct {
	ID           string
	UserID       int64
	CharacterID  int64
	ASRClient    *RealtimeASRClient
	IsActive     bool
	LastText     string
	TextOffset   int
	SilenceTimer *time.Timer
	AudioBuffer  []byte // 添加音频缓冲区
	mu           sync.RWMutex
}

// StreamingVoiceCallRequest 流式语音通话请求
type StreamingVoiceCallRequest struct {
	UserID      int64  `json:"user_id"`
	CharacterID int64  `json:"character_id"`
	SessionID   string `json:"session_id"`
}

// StreamingVoiceCallResponse 流式语音通话响应
type StreamingVoiceCallResponse struct {
	SessionID     string `json:"session_id"`
	UserText      string `json:"user_text"`
	TextResponse  string `json:"text_response"`
	AudioResponse string `json:"audio_response"`
	IsComplete    bool   `json:"is_complete"`
	Error         string `json:"error,omitempty"`
}

// NewStreamingVoiceCallService 创建流式语音通话服务
func NewStreamingVoiceCallService(aiService *AIService, db *sql.DB) *StreamingVoiceCallService {
	return &StreamingVoiceCallService{
		aiService: aiService,
		db:        db,
		sessions:  make(map[string]*VoiceCallSession),
	}
}

// SetConnectionErrorCallback 设置连接错误回调
func (s *StreamingVoiceCallService) SetConnectionErrorCallback(callback func(sessionID string, err error)) {
	s.onConnectionError = callback
}

// GetCharacterByID 根据ID获取角色信息
func (s *StreamingVoiceCallService) GetCharacterByID(characterID int) (*models.PresetCharacter, error) {
	query := `SELECT id, name, personality_signature, avatar_url FROM preset_characters WHERE id = ?`

	var character models.PresetCharacter
	err := s.db.QueryRow(query, characterID).Scan(
		&character.ID,
		&character.Name,
		&character.PersonalitySignature,
		&character.AvatarURL,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, fmt.Errorf("获取角色信息失败: %v", err)
	}

	return &character, nil
}

// GenerateTTS 生成TTS音频
func (s *StreamingVoiceCallService) GenerateTTS(text string, characterName string) ([]byte, error) {
	return s.aiService.TextToSpeech(text, characterName)
}

// SetResponseCallback 设置AI回复回调函数
func (s *StreamingVoiceCallService) SetResponseCallback(callback func(sessionID, userText, aiText string, audioData []byte)) {
	s.onResponseCallback = callback
}

// StartStreamingCall 开始流式语音通话
func (s *StreamingVoiceCallService) StartStreamingCall(req *StreamingVoiceCallRequest) (*StreamingVoiceCallResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否已存在会话
	if session, exists := s.sessions[req.SessionID]; exists {
		if session.IsActive {
			return &StreamingVoiceCallResponse{
				SessionID:  req.SessionID,
				IsComplete: false,
			}, nil
		}
	}

	// 创建会话（不使用WebSocket ASR，改用HTTP ASR）
	session := &VoiceCallSession{
		ID:          req.SessionID,
		UserID:      req.UserID,
		CharacterID: req.CharacterID,
		ASRClient:   nil, // 不使用WebSocket ASR
		IsActive:    true,
		LastText:    "",
		TextOffset:  0,
	}

	s.sessions[req.SessionID] = session

	return &StreamingVoiceCallResponse{
		SessionID:  req.SessionID,
		IsComplete: false,
	}, nil
}

// ProcessAudioChunk 处理音频分片
func (s *StreamingVoiceCallService) ProcessAudioChunk(sessionID string, audioData []byte) error {
	s.mu.RLock()
	session, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists || !session.IsActive {
		return fmt.Errorf("session not found or inactive")
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	fmt.Printf("处理音频分片: sessionID=%s, 数据长度=%d bytes\n", sessionID, len(audioData))

	// 累积音频数据
	if session.AudioBuffer == nil {
		session.AudioBuffer = make([]byte, 0)
	}
	session.AudioBuffer = append(session.AudioBuffer, audioData...)

	// 不设置固定阈值，而是等待前端发送"语音结束"信号
	// 前端会在检测到语音结束时发送一个特殊的结束信号

	return nil
}

// ProcessVoiceEnd 处理语音结束信号
func (s *StreamingVoiceCallService) ProcessVoiceEnd(sessionID string) error {
	s.mu.RLock()
	session, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists || !session.IsActive {
		return fmt.Errorf("session not found or inactive")
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	// 如果有累积的音频数据，立即处理
	if len(session.AudioBuffer) > 0 {
		fmt.Printf("语音结束，处理累积音频: sessionID=%s, 数据长度=%d bytes\n", sessionID, len(session.AudioBuffer))

		// 复制音频数据用于处理
		audioDataCopy := make([]byte, len(session.AudioBuffer))
		copy(audioDataCopy, session.AudioBuffer)
		session.AudioBuffer = make([]byte, 0) // 重置缓冲区

		// 异步处理音频数据
		go s.processAccumulatedAudio(session.ID, audioDataCopy)
	}

	return nil
}

// processAccumulatedAudio 处理累积的音频数据
func (s *StreamingVoiceCallService) processAccumulatedAudio(sessionID string, audioData []byte) {
	if len(audioData) == 0 {
		return
	}

	fmt.Printf("处理累积音频: sessionID=%s, 数据长度=%d bytes\n", sessionID, len(audioData))

	// 调用ASR进行语音识别
	text, err := s.aiService.SpeechToText(audioData)
	if err != nil {
		log.Printf("ASR识别失败: %v", err)

		// ASR失败时，获取角色信息并回复噪音响应
		s.mu.RLock()
		session, exists := s.sessions[sessionID]
		s.mu.RUnlock()

		if !exists || !session.IsActive {
			return
		}

		// 获取角色信息
		character, err := s.getCharacterByID(session.CharacterID)
		if err != nil {
			log.Printf("Failed to get character: %v", err)
			return
		}

		// 获取噪音响应
		noiseResponse := s.getNoiseResponseForCharacter(character.Name)
		log.Printf("ASR失败，使用噪音响应: %s", noiseResponse)

		// 处理噪音响应
		s.processAIResponse(session, "", noiseResponse, character)
		return
	}

	fmt.Printf("ASR识别结果: %s\n", text)

	if text == "" {
		log.Printf("ASR识别结果为空，跳过处理")
		return
	}

	// 获取会话信息
	s.mu.RLock()
	session, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists || !session.IsActive {
		log.Printf("会话不存在或已失效: %s", sessionID)
		return
	}

	// 处理识别结果
	s.processCompleteText(session, text)
}

// listenForResponses 监听ASR响应
func (s *StreamingVoiceCallService) listenForResponses(session *VoiceCallSession) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in listenForResponses: %v", r)
			// 通知WebSocket连接错误
			if s.onConnectionError != nil {
				s.onConnectionError(session.ID, fmt.Errorf("panic recovered: %v", r))
			}
		}
		session.ASRClient.Close()
		s.mu.Lock()
		delete(s.sessions, session.ID)
		s.mu.Unlock()
	}()

	for session.IsActive {
		// 检查连接状态
		if !session.ASRClient.IsConnected() {
			log.Printf("ASR client disconnected, stopping listener")
			// 通知WebSocket连接错误
			if s.onConnectionError != nil {
				s.onConnectionError(session.ID, fmt.Errorf("ASR client disconnected"))
			}
			break
		}

		// 设置读取超时
		err := session.ASRClient.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			log.Printf("Failed to set read deadline: %v", err)
			// 通知WebSocket连接错误
			if s.onConnectionError != nil {
				s.onConnectionError(session.ID, fmt.Errorf("failed to set read deadline: %w", err))
			}
			break
		}

		text, err := session.ASRClient.ReadResponse()
		if err != nil {
			log.Printf("Failed to read ASR response: %v", err)
			// 如果是连接错误，停止监听并通知
			if !session.ASRClient.IsConnected() {
				if s.onConnectionError != nil {
					s.onConnectionError(session.ID, fmt.Errorf("ASR connection failed: %w", err))
				}
				break
			}
			continue
		}

		if text == "" {
			continue
		}

		session.mu.Lock()

		// 检查是否有新文本
		if len(text) > session.TextOffset {
			session.TextOffset = len(text)

			// 重置静音计时器
			if session.SilenceTimer != nil {
				session.SilenceTimer.Stop()
			}

			// 设置静音超时（1.5秒）
			session.SilenceTimer = time.AfterFunc(1500*time.Millisecond, func() {
				s.processCompleteText(session, text)
			})
		}

		session.mu.Unlock()
	}
}

// processCompleteText 处理完整文本
func (s *StreamingVoiceCallService) processCompleteText(session *VoiceCallSession, text string) {
	if text == "" || len(text) < 2 {
		return
	}

	log.Printf("Processing complete text: %s", text)

	// 获取角色信息
	character, err := s.getCharacterByID(session.CharacterID)
	if err != nil {
		log.Printf("Failed to get character: %v", err)
		return
	}

	// 获取对话历史（像普通语音通话一样）
	history, err := s.getConversationHistory(session.UserID, session.CharacterID, 5)
	if err != nil {
		log.Printf("Failed to get conversation history: %v", err)
		// 如果获取历史失败，使用简单的消息
		messages := []Message{
			{Role: "system", Content: character.PersonalitySignature},
			{Role: "user", Content: text},
		}
		aiText, err := s.aiService.ChatWithLLM(messages, "", 0.7, "voice")
		if err != nil {
			log.Printf("Failed to get LLM response: %v", err)
			return
		}
		// 处理AI回复
		s.processAIResponse(session, text, aiText, character)
	} else {
		// 构建消息历史（像普通语音通话一样）
		messageHistory := s.buildMessageHistoryWithMemory(history, 4)

		// 构建消息列表
		messages := append(messageHistory, Message{
			Role:    "user",
			Content: text,
		})

		// AI回复
		aiText, err := s.aiService.ChatWithLLM(messages, "", 0.7, "voice")
		if err != nil {
			log.Printf("Failed to get LLM response: %v", err)
			return
		}

		// 处理AI回复
		s.processAIResponse(session, text, aiText, character)
	}
}

// processAIResponse 处理AI回复
func (s *StreamingVoiceCallService) processAIResponse(session *VoiceCallSession, userText, aiText string, character *models.PresetCharacter) {

	// 调用TTS
	aiAudioData, err := s.aiService.TextToSpeech(aiText, character.Name)
	if err != nil {
		log.Printf("Failed to get TTS response: %v", err)
		return
	}

	// 保存对话记录
	err = s.saveVoiceCall(session.UserID, session.CharacterID, session.ID, userText, aiText)
	if err != nil {
		log.Printf("Failed to save voice call: %v", err)
	}

	// 通过WebSocket将结果发送给前端
	log.Printf("AI Response: %s", aiText)
	log.Printf("Audio data length: %d bytes", len(aiAudioData))

	// 发送AI回复给前端
	if s.onResponseCallback != nil {
		s.onResponseCallback(session.ID, userText, aiText, aiAudioData)
	}

	// 重置文本偏移
	session.mu.Lock()
	session.TextOffset = 0
	session.LastText = ""
	session.mu.Unlock()
}

// StopStreamingCall 停止流式语音通话
func (s *StreamingVoiceCallService) StopStreamingCall(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	session.mu.Lock()
	session.IsActive = false
	if session.SilenceTimer != nil {
		session.SilenceTimer.Stop()
	}
	session.mu.Unlock()

	return nil
}

// GetSessionStatus 获取会话状态
func (s *StreamingVoiceCallService) GetSessionStatus(sessionID string) (*VoiceCallSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	return session, exists
}

// 辅助方法
func (s *StreamingVoiceCallService) getCharacterByID(characterID int64) (*models.PresetCharacter, error) {
	query := `SELECT id, name, personality_signature, avatar_url FROM preset_characters WHERE id = ?`

	var character models.PresetCharacter
	err := s.db.QueryRow(query, characterID).Scan(
		&character.ID,
		&character.Name,
		&character.PersonalitySignature,
		&character.AvatarURL,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, fmt.Errorf("获取角色信息失败: %v", err)
	}

	return &character, nil
}

// getConversationHistory 获取对话历史
func (s *StreamingVoiceCallService) getConversationHistory(userID, characterID int64, limit int) ([]models.Conversation, error) {
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
func (s *StreamingVoiceCallService) buildMessageHistoryWithMemory(conversations []models.Conversation, maxHistory int) []Message {
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

func (s *StreamingVoiceCallService) saveVoiceCall(userID, characterID int64, sessionID, userText, aiText string) error {
	// 使用 conversations 表保存对话记录
	query := `
		INSERT INTO conversations (user_id, character_id, user_message, ai_response, message_type, session_id, created_at)
		VALUES (?, ?, ?, ?, 'voice', ?, NOW())
	`

	result, err := s.db.Exec(query, userID, characterID, userText, aiText, sessionID)
	if err != nil {
		return fmt.Errorf("保存语音通话记录失败: %v", err)
	}

	// 获取插入的记录ID
	insertID, _ := result.LastInsertId()
	fmt.Printf("Successfully saved voice call: ID=%d, UserID=%d, CharacterID=%d\n", insertID, userID, characterID)

	// 更新好友关系的最后消息时间
	_, err = s.db.Exec(`
		UPDATE user_friendships 
		SET last_message_at = NOW(), updated_at = NOW()
		WHERE user_id = ? AND character_id = ?
	`, userID, characterID)
	if err != nil {
		// 记录错误但不影响语音通话
		fmt.Printf("Failed to update friendship last_message_at: %v\n", err)
	}

	return nil
}

// getNoiseResponseForCharacter 根据角色人设返回噪音响应
func (s *StreamingVoiceCallService) getNoiseResponseForCharacter(characterName string) string {
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
