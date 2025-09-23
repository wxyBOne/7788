package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AIService struct {
	apiKey       string
	baseURL      string
	asrAPIKey    string
	ttsAPIKey    string
	visionAPIKey string
	client       *http.Client
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type VisionRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type ASRRequest struct {
	Audio string `json:"audio"` // base64编码的音频
}

type ASRResponse struct {
	Text string `json:"text"`
}

type TTSRequest struct {
	Text  string  `json:"text"`
	Voice string  `json:"voice"`
	Speed float64 `json:"speed"`
	Pitch float64 `json:"pitch"`
}

type TTSResponse struct {
	AudioURL string `json:"audio_url"`
}

func NewAIService(apiKey, baseURL, asrAPIKey, ttsAPIKey, visionAPIKey string) *AIService {
	return &AIService{
		apiKey:       apiKey,
		baseURL:      baseURL,
		asrAPIKey:    asrAPIKey,
		ttsAPIKey:    ttsAPIKey,
		visionAPIKey: visionAPIKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ChatWithLLM 与LLM对话
func (s *AIService) ChatWithLLM(messages []Message, model string, temperature float64) (string, error) {
	// 暂时返回模拟响应，等API配置完成后替换
	if s.apiKey == "" {
		return s.generateMockResponse(messages), nil
	}

	req := ChatRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   1000,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", s.baseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from AI")
}

// AnalyzeImage 分析图片
func (s *AIService) AnalyzeImage(imageBase64 string, prompt string) (string, error) {
	// 暂时返回模拟响应
	if s.visionAPIKey == "" {
		return "我看到了一张图片，但暂时无法详细分析。", nil
	}

	messages := []Message{
		{
			Role:    "user",
			Content: fmt.Sprintf("%s\n\n图片数据: %s", prompt, imageBase64),
		},
	}

	return s.ChatWithLLM(messages, "qwen-vl-max", 0.7)
}

// SpeechToText 语音转文字
func (s *AIService) SpeechToText(audioData []byte) (string, error) {
	// 暂时返回模拟响应
	if s.asrAPIKey == "" {
		return "我听到了您的声音，但暂时无法识别具体内容。", nil
	}

	// TODO: 实现真实的ASR API调用
	return "模拟语音识别结果", nil
}

// TextToSpeech 文字转语音
func (s *AIService) TextToSpeech(text, voiceType string) (string, error) {
	// 暂时返回模拟响应
	if s.ttsAPIKey == "" {
		return "", fmt.Errorf("TTS服务暂未配置")
	}

	// TODO: 实现真实的TTS API调用
	return "", fmt.Errorf("TTS服务暂未实现")
}

// 生成模拟响应（用于开发阶段）
func (s *AIService) generateMockResponse(messages []Message) string {
	if len(messages) == 0 {
		return "你好！我是AI助手，很高兴为您服务。"
	}

	lastMessage := messages[len(messages)-1].Content

	// 根据消息内容生成简单的模拟响应
	if contains(lastMessage, "你好") || contains(lastMessage, "hello") {
		return "你好！很高兴见到你！"
	}

	if contains(lastMessage, "你是谁") || contains(lastMessage, "who are you") {
		return "我是一个AI助手，正在学习和成长中。"
	}

	return "我明白了，让我想想怎么回答你..."
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
