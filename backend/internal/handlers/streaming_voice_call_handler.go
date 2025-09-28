package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"seven-ai-backend/internal/models"
	"seven-ai-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketUpgrader WebSocket升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 在生产环境中应该检查来源
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// StreamingVoiceCallHandler 流式语音通话处理器
type StreamingVoiceCallHandler struct {
	streamingService  *services.StreamingVoiceCallService
	activeConnections map[string]*websocket.Conn
	mu                sync.RWMutex
}

// NewStreamingVoiceCallHandler 创建流式语音通话处理器
func NewStreamingVoiceCallHandler(streamingService *services.StreamingVoiceCallService) *StreamingVoiceCallHandler {
	handler := &StreamingVoiceCallHandler{
		streamingService:  streamingService,
		activeConnections: make(map[string]*websocket.Conn),
	}

	// 设置连接错误回调
	streamingService.SetConnectionErrorCallback(handler.handleConnectionError)

	// 设置AI回复回调
	streamingService.SetResponseCallback(handler.handleAIResponse)

	return handler
}

// WebSocketMessage WebSocket消息
type WebSocketMessage struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id"`
	Data      interface{} `json:"data"`
}

// AudioChunkMessage 音频分片消息
type AudioChunkMessage struct {
	AudioData []int `json:"audio_data"` // 前端发送的是int数组
}

// StartCallMessage 开始通话消息
type StartCallMessage struct {
	UserID      int64 `json:"user_id"`
	CharacterID int64 `json:"character_id"`
}

// ResponseMessage 响应消息
type ResponseMessage struct {
	UserText      string `json:"user_text"`
	TextResponse  string `json:"text_response"`
	AudioResponse string `json:"audio_response"`
	IsComplete    bool   `json:"is_complete"`
}

// handleAIResponse 处理AI回复
func (h *StreamingVoiceCallHandler) handleAIResponse(sessionID, userText, aiText string, audioData []byte) {
	h.mu.RLock()
	conn, exists := h.activeConnections[sessionID]
	h.mu.RUnlock()

	if !exists || conn == nil {
		fmt.Printf("WebSocket连接不存在: %s\n", sessionID)
		return
	}

	// 发送AI回复消息
	response := WebSocketMessage{
		Type:      "ai_response",
		SessionID: sessionID,
		Data: map[string]interface{}{
			"user_text":   userText,
			"ai_text":     aiText,
			"audio_data":  audioData,
			"is_complete": true,
		},
	}

	h.sendMessage(conn, response)
	fmt.Printf("发送AI回复: sessionID=%s, userText=%s, aiText=%s\n", sessionID, userText, aiText)
}

// handleConnectionError 处理连接错误
func (h *StreamingVoiceCallHandler) handleConnectionError(sessionID string, err error) {
	h.mu.RLock()
	conn, exists := h.activeConnections[sessionID]
	h.mu.RUnlock()

	if exists && conn != nil {
		// 发送错误消息给前端
		h.sendError(conn, sessionID, err.Error())

		// 关闭连接
		conn.Close()

		// 从活跃连接中移除
		h.mu.Lock()
		delete(h.activeConnections, sessionID)
		h.mu.Unlock()
	}
}

// getGreetingText 获取角色打招呼文本（复制普通语音通话的逻辑）
func getGreetingText(character *models.PresetCharacter) string {
	switch character.Name {
	case "林黛玉":
		return "哟，你来了？今日倒是比往日早了些，莫非是有什么心事？"
	case "孙悟空":
		return "兄弟，你来了！俺老孙正闲着，有什么话尽管说，别客气！"
	case "李白":
		return "人生得意须尽欢，今日得与君通话，当浮一大白！有什么想聊的？"
	case "赫敏·格兰杰", "赫敏":
		return "你好，这里是赫敏，找我有什么事吗？"
	default:
		return fmt.Sprintf("你好！我是%s，很高兴和你通话！", character.Name)
	}
}

// HandleFirstCall 处理第一次流式通话请求，让AI主动打招呼
func (h *StreamingVoiceCallHandler) HandleFirstCall(c *gin.Context) {
	// 获取用户ID
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user ID"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		CharacterID int    `json:"character_id"`
		SessionID   string `json:"session_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("收到first-call请求: UserID=%d, CharacterID=%d, SessionID=%s", userID, req.CharacterID, req.SessionID)

	// 获取角色信息
	character, err := h.streamingService.GetCharacterByID(req.CharacterID)
	if err != nil {
		log.Printf("获取角色信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色信息失败"})
		return
	}

	log.Printf("获取角色成功: %s", character.Name)

	// 生成AI的打招呼文本（复制普通语音通话的逻辑）
	greetingText := getGreetingText(character)
	log.Printf("生成打招呼文本: %s", greetingText)

	// 调用TTS生成音频
	audioData, err := h.streamingService.GenerateTTS(greetingText, character.Name)
	if err != nil {
		log.Printf("生成AI音频失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成AI音频失败"})
		return
	}

	log.Printf("TTS生成成功，音频长度: %d bytes", len(audioData))

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"text_response":  greetingText,
		"audio_response": audioData,
		"session_id":     req.SessionID,
	})
}

// HandleWebSocket 处理WebSocket连接
func (h *StreamingVoiceCallHandler) HandleWebSocket(c *gin.Context) {
	// 简单的认证检查（生产环境中应该使用更安全的认证方式）
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// 获取用户ID（从查询参数获取，因为WebSocket不支持自定义头部）
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user ID"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	log.Printf("WebSocket连接用户ID: %d", userID)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	defer conn.Close()

	// 设置连接参数 - 增加超时时间到5分钟，给ASR+TTS处理留足够时间
	conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		return nil
	})

	// 启动ping定时器 - 每2分钟发送一次ping
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()

	var sessionID string
	var isCallActive bool

	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		switch msg.Type {
		case "start_call":
			// 开始通话
			var startMsg StartCallMessage
			if dataBytes, err := json.Marshal(msg.Data); err == nil {
				json.Unmarshal(dataBytes, &startMsg)
			}

			req := &services.StreamingVoiceCallRequest{
				UserID:      int64(userID), // 转换为int64类型
				CharacterID: startMsg.CharacterID,
				SessionID:   msg.SessionID,
			}

			resp, err := h.streamingService.StartStreamingCall(req)
			if err != nil {
				h.sendError(conn, msg.SessionID, err.Error())
				continue
			}

			sessionID = resp.SessionID
			isCallActive = true

			// 注册连接
			h.mu.Lock()
			h.activeConnections[sessionID] = conn
			h.mu.Unlock()

			h.sendMessage(conn, WebSocketMessage{
				Type:      "call_started",
				SessionID: sessionID,
				Data:      resp,
			})

		case "audio_chunk":
			// 处理音频分片
			if !isCallActive {
				h.sendError(conn, sessionID, "Call not active")
				continue
			}

			var audioMsg AudioChunkMessage
			if dataBytes, err := json.Marshal(msg.Data); err == nil {
				json.Unmarshal(dataBytes, &audioMsg)
			}

			fmt.Printf("收到音频分片: sessionID=%s, 数据长度=%d\n", sessionID, len(audioMsg.AudioData))

			// 将int数组转换为byte数组
			audioBytes := make([]byte, len(audioMsg.AudioData))
			for i, v := range audioMsg.AudioData {
				audioBytes[i] = byte(v)
			}

			err = h.streamingService.ProcessAudioChunk(sessionID, audioBytes)
			if err != nil {
				fmt.Printf("处理音频分片失败: %v\n", err)
				h.sendError(conn, sessionID, err.Error())
				continue
			}

			fmt.Printf("音频分片处理成功: sessionID=%s\n", sessionID)

		case "voice_end":
			// 处理语音结束信号
			if isCallActive {
				err := h.streamingService.ProcessVoiceEnd(sessionID)
				if err != nil {
					h.sendError(conn, sessionID, err.Error())
					continue
				}
			}

		case "stop_call":
			// 停止通话
			if isCallActive {
				err := h.streamingService.StopStreamingCall(sessionID)
				if err != nil {
					h.sendError(conn, sessionID, err.Error())
				} else {
					h.sendMessage(conn, WebSocketMessage{
						Type:      "call_stopped",
						SessionID: sessionID,
						Data:      map[string]bool{"success": true},
					})
				}
				isCallActive = false
			}

		case "ping":
			// 心跳
			h.sendMessage(conn, WebSocketMessage{
				Type: "pong",
				Data: map[string]string{"timestamp": time.Now().Format(time.RFC3339)},
			})
		}
	}

	// 清理会话和连接
	if isCallActive && sessionID != "" {
		h.streamingService.StopStreamingCall(sessionID)

		// 从活跃连接中移除
		h.mu.Lock()
		delete(h.activeConnections, sessionID)
		h.mu.Unlock()
	}
}

// sendMessage 发送消息
func (h *StreamingVoiceCallHandler) sendMessage(conn *websocket.Conn, msg WebSocketMessage) {
	err := conn.WriteJSON(msg)
	if err != nil {
		log.Printf("Failed to send WebSocket message: %v", err)
	}
}

// sendError 发送错误消息
func (h *StreamingVoiceCallHandler) sendError(conn *websocket.Conn, sessionID, errorMsg string) {
	h.sendMessage(conn, WebSocketMessage{
		Type:      "error",
		SessionID: sessionID,
		Data:      map[string]string{"error": errorMsg},
	})
}

// GetSessionStatus 获取会话状态
func (h *StreamingVoiceCallHandler) GetSessionStatus(c *gin.Context) {
	sessionID := c.Param("sessionId")

	session, exists := h.streamingService.GetSessionStatus(sessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":   session.ID,
		"is_active":    session.IsActive,
		"user_id":      session.UserID,
		"character_id": session.CharacterID,
	})
}
