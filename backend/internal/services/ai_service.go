package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type AIService struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
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
	Index   int     `json:"index"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ASR相关结构体
type ASRRequest struct {
	Model string `json:"model"`
	Audio struct {
		Format string `json:"format"`
		URL    string `json:"url"`
	} `json:"audio"`
}

type ASRResponse struct {
	Reqid     string `json:"reqid"`
	Operation string `json:"operation"`
	Data      struct {
		Result struct {
			Text string `json:"text"`
		} `json:"result"`
	} `json:"data"`
}

// TTS相关结构体
type TTSRequest struct {
	Audio struct {
		VoiceType  string  `json:"voice_type"`
		Encoding   string  `json:"encoding"`
		SpeedRatio float64 `json:"speed_ratio"`
	} `json:"audio"`
	Request struct {
		Text      string `json:"text"`
		MaxTokens int    `json:"max_tokens,omitempty"`
	} `json:"request"`
}

type TTSResponse struct {
	Reqid     string `json:"reqid"`
	Operation string `json:"operation"`
	Sequence  int    `json:"sequence"`
	Data      string `json:"data"` // base64编码的音频数据
}

type VisionRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type VisionResponse struct {
	Choices []VisionChoice `json:"choices"`
}

type VisionChoice struct {
	Message VisionMessage `json:"message"`
}

type VisionMessage struct {
	Content string `json:"content"`
}

func NewAIService(apiKey, baseURL, model string) *AIService {
	return &AIService{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 6 * time.Second, // 进一步减少超时时间
		},
	}
}

// ChatWithLLM 与LLM对话
func (s *AIService) ChatWithLLM(messages []Message, model string, temperature float64, messageType string) (string, error) {
	// 检查消息类型，如果是表情消息，添加表情识别提示
	messages = s.processEmojiMessages(messages)

	// 添加消息类型标识
	if messageType == "voice" {
		// 在系统消息中添加语音通话标识
		systemMessage := Message{
			Role:    "system",
			Content: "这是一次语音通话，请用符合角色的口吻以及语气自然地和用户对话，请口语化而不是书面语。绝不使用任何括号内的动作、表情、语气或场景描写。",
		}
		messages = append([]Message{systemMessage}, messages...)
	}

	// 暂时返回模拟响应，等API配置完成后替换
	if s.apiKey == "" {
		return "你好！我是AI助手，很高兴为您服务。", nil
	}

	// 使用配置的模型，如果传入的model为空则使用默认模型
	usedModel := model
	if usedModel == "" {
		usedModel = s.model
	}

	req := ChatRequest{
		Model:       usedModel,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   30, // 进一步减少到30，冲刺5秒内
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
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("AI API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from AI")
}

// AnalyzeImage 分析图片
func (s *AIService) AnalyzeImage(imageBase64 string, prompt string) (string, error) {
	// 暂时返回模拟响应
	return "我看到了一张图片，但暂时无法详细分析。", nil
}

// SpeechToText 语音转文字 (ASR)
func (s *AIService) SpeechToText(audioData []byte) (string, error) {
	if len(audioData) == 0 || s.apiKey == "" {
		return "", nil
	}

	start := time.Now()

	// 上传到公网URL并调用ASR API
	publicURL, err := s.uploadToPublicURL(audioData)
	if err != nil {
		return "", err
	}
	fmt.Printf("上传耗时: %v\n", time.Since(start))

	// 调用真实的ASR API
	asrStart := time.Now()
	result, err := s.callQiniuASRAPI(publicURL)
	if err != nil {
		return "", err
	}
	fmt.Printf("ASR耗时: %v\n", time.Since(asrStart))
	fmt.Printf("总耗时: %v\n", time.Since(start))

	return result, nil
}

// callQiniuASRAPI 调用七牛云ASR API
func (s *AIService) callQiniuASRAPI(audioURL string) (string, error) {
	// 尝试多种音频格式
	formats := []string{"mp4", "webm", "wav"}

	for _, format := range formats {
		// 构建ASR请求
		req := ASRRequest{
			Model: "asr",
			Audio: struct {
				Format string `json:"format"`
				URL    string `json:"url"`
			}{
				Format: format,
				URL:    audioURL,
			},
		}

		reqBody, err := json.Marshal(req)
		if err != nil {
			continue
		}

		// 尝试主要接入点
		result, err := s.tryASRRequest(s.baseURL+"/voice/asr", reqBody)
		if err == nil {
			return result, nil
		}

		// 如果主要接入点失败，尝试备用接入点
		fmt.Printf("格式 %s 主要接入点失败，尝试备用接入点: %v\n", format, err)
		backupURL := "https://api.qnaigc.com/v1/voice/asr"
		result, err = s.tryASRRequest(backupURL, reqBody)
		if err == nil {
			return result, nil
		}

		fmt.Printf("格式 %s 也失败了: %v\n", format, err)
	}

	return "", fmt.Errorf("所有格式的ASR请求都失败了")
}

// tryASRRequest 尝试ASR请求
func (s *AIService) tryASRRequest(url string, reqBody []byte) (string, error) {
	start := time.Now()

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: 4 * time.Second, // 进一步减少ASR超时时间
	}

	// 发送HTTP请求
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create ASR request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("ASR HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("ASR请求耗时: %v (URL: %s)\n", time.Since(start), url)

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ASR API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var asrResp ASRResponse
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read ASR response: %w", err)
	}

	if err := json.Unmarshal(respBody, &asrResp); err != nil {
		return "", fmt.Errorf("failed to parse ASR response: %w", err)
	}

	// 提取识别文本
	if asrResp.Data.Result.Text != "" {
		return asrResp.Data.Result.Text, nil
	}

	return "", nil
}

// uploadToPublicURL 上传音频到公网URL
func (s *AIService) uploadToPublicURL(audioData []byte) (string, error) {
	start := time.Now()

	// 如果文件太大，尝试压缩
	if len(audioData) > 100000 { // 100KB
		fmt.Printf("文件较大 (%d bytes)，尝试优化上传\n", len(audioData))
	}

	// 创建multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件字段 - 动态检测音频格式
	filename := fmt.Sprintf("audio_%d.webm", time.Now().Unix())
	if len(audioData) > 4 {
		// 检测音频格式
		if audioData[0] == 0x00 && audioData[1] == 0x00 && audioData[2] == 0x00 && audioData[3] == 0x20 {
			filename = fmt.Sprintf("audio_%d.mp4", time.Now().Unix())
		} else if audioData[0] == 0x52 && audioData[1] == 0x49 && audioData[2] == 0x46 && audioData[3] == 0x46 {
			filename = fmt.Sprintf("audio_%d.wav", time.Now().Unix())
		}
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = part.Write(audioData)
	if err != nil {
		return "", fmt.Errorf("failed to write file data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// 使用0x0.st文件上传服务
	req, err := http.NewRequest("POST", "https://0x0.st", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 为上传创建专门的客户端，合理超时时间
	uploadClient := &http.Client{
		Timeout: 10 * time.Second, // 上传需要一定时间，但不要太长
	}

	resp, err := uploadClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read upload response: %w", err)
	}

	// 0x0.st返回的是直接的URL
	publicURL := strings.TrimSpace(string(respBody))

	// 验证URL格式
	if !strings.HasPrefix(publicURL, "http") {
		return "", fmt.Errorf("invalid URL returned: %s", publicURL)
	}

	fmt.Printf("上传耗时: %v (文件大小: %d bytes)\n", time.Since(start), len(audioData))
	return publicURL, nil
}

// SpeechToTextWithCharacter 带角色人设的语音转文字
func (s *AIService) SpeechToTextWithCharacter(audioData []byte, characterName string) (string, error) {
	// 调用ASR
	text, err := s.SpeechToText(audioData)
	if err != nil {
		return "", nil // 返回空字符串，让后端处理
	}

	// 如果识别结果为空或太短，返回空字符串
	if strings.TrimSpace(text) == "" || len(strings.TrimSpace(text)) < 2 {
		return "", nil
	}

	return text, nil
}

// getNoiseResponseForCharacter 根据角色人设返回噪音响应
func (s *AIService) getNoiseResponseForCharacter(characterName string) string {
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

// TextToSpeech 文字转语音 (TTS)
func (s *AIService) TextToSpeech(text string, characterName string) ([]byte, error) {
	start := time.Now()

	// 限制文本长度，提高TTS速度
	if len([]rune(text)) > 20 {
		text = string([]rune(text)[:20]) + "..."
	}

	// 根据角色选择音色
	voiceType := s.getVoiceTypeForCharacter(characterName)

	req := TTSRequest{
		Audio: struct {
			VoiceType  string  `json:"voice_type"`
			Encoding   string  `json:"encoding"`
			SpeedRatio float64 `json:"speed_ratio"`
		}{
			VoiceType:  voiceType,
			Encoding:   "mp3",
			SpeedRatio: 0.9, // 稍微慢一点，更自然
		},
		Request: struct {
			Text      string `json:"text"`
			MaxTokens int    `json:"max_tokens,omitempty"`
		}{
			Text:      text,
			MaxTokens: 60, // 与LLM保持一致，减少TTS处理时间
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TTS request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.baseURL+"/voice/tts", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create TTS request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	// 为TTS创建专门的客户端，合理超时时间
	ttsClient := &http.Client{
		Timeout: 8 * time.Second, // TTS需要一定时间处理，但不要太长
	}

	resp, err := ttsClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("TTS HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TTS API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var ttsResp TTSResponse
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read TTS response: %w", err)
	}

	if err := json.Unmarshal(respBody, &ttsResp); err != nil {
		return nil, fmt.Errorf("failed to parse TTS response: %w", err)
	}

	// 解码base64音频数据
	audioData, err := base64.StdEncoding.DecodeString(ttsResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode TTS audio data: %w", err)
	}

	fmt.Printf("TTS耗时: %v\n", time.Since(start))
	return audioData, nil
}

// getVoiceTypeForCharacter 根据角色名称获取对应的音色
func (s *AIService) getVoiceTypeForCharacter(characterName string) string {
	switch characterName {
	case "林黛玉":
		return "qiniu_zh_female_tmjxxy" // 甜美教学小源 - 温柔细腻
	case "孙悟空":
		return "qiniu_zh_male_mzjsxg" // 名著角色猴哥 - 豪爽有力
	case "李白":
		return "qiniu_zh_male_tyygjs" // 通用阳光讲师 - 豪迈洒脱
	case "赫敏·格兰杰":
		return "qiniu_zh_female_zxjxnjs" // 知性教学女教师 - 清晰理性
	default:
		return "qiniu_zh_female_tmjxxy" // 默认甜美教学小源
	}
}

// processEmojiMessages 处理表情消息，添加表情识别提示
func (s *AIService) processEmojiMessages(messages []Message) []Message {
	for i, msg := range messages {
		if msg.Role == "user" && s.isEmojiMessage(msg.Content) {
			// 为用户消息添加表情识别提示
			messages[i].Content = fmt.Sprintf("用户发送了表情：%s\n请根据这个表情和你的角色性格，自然地回应。", msg.Content)
		}
	}
	return messages
}

// isEmojiMessage 检查消息是否只包含表情
func (s *AIService) isEmojiMessage(content string) bool {
	// 简单的表情检测：如果消息只包含表情符号（Unicode表情范围）
	emojiRanges := []struct {
		start, end rune
	}{
		{0x1F600, 0x1F64F}, // Emoticons
		{0x1F300, 0x1F5FF}, // Misc Symbols and Pictographs
		{0x1F680, 0x1F6FF}, // Transport and Map
		{0x1F1E0, 0x1F1FF}, // Regional indicator symbols
		{0x2600, 0x26FF},   // Misc symbols
		{0x2700, 0x27BF},   // Dingbats
		{0xFE00, 0xFE0F},   // Variation Selectors
		{0x1F900, 0x1F9FF}, // Supplemental Symbols and Pictographs
		{0x1F018, 0x1F0F5}, // Playing cards
		{0x1F200, 0x1F2FF}, // Enclosed characters
	}

	for _, r := range content {
		isEmoji := false
		for _, emojiRange := range emojiRanges {
			if r >= emojiRange.start && r <= emojiRange.end {
				isEmoji = true
				break
			}
		}
		if !isEmoji && r != ' ' && r != '\n' && r != '\t' {
			return false
		}
	}
	return len(content) > 0
}

// IsEmojiMessage 检查消息是否只包含表情（公开方法）
func (s *AIService) IsEmojiMessage(content string) bool {
	return s.isEmojiMessage(content)
}
