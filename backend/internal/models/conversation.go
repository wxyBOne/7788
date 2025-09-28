package models

import "time"

type Conversation struct {
	ID               int       `json:"id" db:"id"`
	UserID           int       `json:"user_id" db:"user_id"`
	CharacterID      *int      `json:"character_id" db:"character_id"`
	CompanionID      *int      `json:"companion_id" db:"companion_id"`
	SessionID        string    `json:"session_id" db:"session_id"`
	MessageType      string    `json:"message_type" db:"message_type"`
	UserMessage      string    `json:"user_message" db:"user_message"`
	AIResponse       string    `json:"ai_response" db:"ai_response"`
	ImageData        string    `json:"image_data" db:"image_data"`
	AudioData        string    `json:"audio_data" db:"audio_data"`
	SentimentScore   float64   `json:"sentiment_score" db:"sentiment_score"`
	ExperienceGained int       `json:"experience_gained" db:"experience_gained"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

type ChatRequest struct {
	CharacterID int    `json:"character_id" binding:"min=0"` // 允许0值（AI伙伴）
	Message     string `json:"message" binding:"required"`
	SessionID   string `json:"session_id"`
}

type ChatResponse struct {
	Response  string `json:"response"`
	SessionID string `json:"session_id"`
	Character string `json:"character"`
	MessageID int    `json:"message_id"`
}

type VoiceChatRequest struct {
	CharacterID int    `json:"character_id" binding:"required"`
	AudioData   string `json:"audio_data" binding:"required"`
	SessionID   string `json:"session_id"`
}

type ImageChatRequest struct {
	CharacterID int    `json:"character_id" binding:"required"`
	ImageData   string `json:"image_data" binding:"required"`
	Message     string `json:"message"`
	SessionID   string `json:"session_id"`
}

type ConversationHistory struct {
	ID          int       `json:"id"`
	UserMessage string    `json:"user_message"`
	AIResponse  string    `json:"ai_response"`
	MessageType string    `json:"message_type"`
	CreatedAt   time.Time `json:"created_at"`
}
