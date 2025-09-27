package models

import "time"

// VoiceCallRequest 语音通话请求
type VoiceCallRequest struct {
	CharacterID int    `json:"character_id" binding:"required"`
	AudioData   string `json:"audio_data"` // base64编码的音频数据，第一次通话时可以为空
	SessionID   string `json:"session_id"`
	IsFirstCall bool   `json:"is_first_call"` // 是否是第一次通话请求
}

// VoiceCallResponse 语音通话响应
type VoiceCallResponse struct {
	TextResponse  string `json:"text_response"`  // AI的文字回复
	AudioResponse string `json:"audio_response"` // base64编码的AI语音回复
	SessionID     string `json:"session_id"`
	Character     string `json:"character"`
	MessageID     int    `json:"message_id"`
}

// VoiceCall 语音通话记录
type VoiceCall struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	CharacterID     int       `json:"character_id"`
	CompanionID     int       `json:"companion_id"`
	CallType        string    `json:"call_type"`
	Status          string    `json:"status"`
	DurationSeconds int       `json:"duration_seconds"`
	AudioFileURL    string    `json:"audio_file_url"`
	StartedAt       time.Time `json:"started_at"`
	EndedAt         time.Time `json:"ended_at"`
}
